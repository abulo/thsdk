package thsdk

// 本文件对应 Python 包的 thsdk/session.py。
//
// _Client 是进程内唯一 THSDK 客户端的所有者：默认实例惰性创建，
// 会话目录取首次调用时的进程工作目录。认证是显式行为：数据方法
// 绝不会在调用方不知情的情况下静默恢复磁盘会话。

import (
	"crypto/sha256"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	// accountMACSeed 对应 Python _ACCOUNT_MAC_SEED = b"thsdk-account-mac-v1\0"。
	accountMACSeed = "thsdk-account-mac-v1\x00"
	envUsername    = "THS_USERNAME"
	envPassword    = "THS_PASSWORD"
	// maxTimeoutSeconds 对应 Python _MAX_TIMEOUT_SECONDS = 600。
	maxTimeoutSeconds = 600.0
	// minQRPollSeconds / maxQRPollSeconds 对应 Python 二维码轮询间隔边界。
	minQRPollSeconds = 0.5
	maxQRPollSeconds = 30.0
	// maxRateLimitRetryMS 对应 Python _MAX_RATE_LIMIT_RETRY_MS = 1_000。
	maxRateLimitRetryMS = 1000
)

// validateSeconds 校验秒级参数范围，对应 Python _validate_seconds。
func validateSeconds(value float64, name string, minimum, maximum float64) (float64, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) || value < minimum || value > maximum {
		return 0, fmt.Errorf("%s 必须在 %g 到 %g 秒之间", name, minimum, maximum)
	}
	return value, nil
}

// timeoutMilliseconds 把秒向上取整为毫秒且至少 1ms，对应 Python _timeout_milliseconds。
func timeoutMilliseconds(value float64) int {
	milliseconds := int(math.Ceil(value * 1000))
	if milliseconds < 1 {
		return 1
	}
	return milliseconds
}

// deriveAccountMAC 把一个账号名映射为稳定的本地管理单播 MAC，
// 对应 Python _derive_account_mac。
func deriveAccountMAC(username string) string {
	digest := sha256.Sum256([]byte(accountMACSeed + username))
	octets := append([]byte(nil), digest[:6]...)
	octets[0] = (octets[0] & 0xFC) | 0x02
	parts := make([]string, len(octets))
	for i, value := range octets {
		parts[i] = fmt.Sprintf("%02x", value)
	}
	return strings.Join(parts, ":")
}

// environmentCredentials 读取环境变量凭据，对应 Python _environment_credentials。
func environmentCredentials() (*string, *string, error) {
	username, hasUsername := os.LookupEnv(envUsername)
	password, hasPassword := os.LookupEnv(envPassword)
	if !hasUsername && !hasPassword {
		return nil, nil, nil
	}
	if !hasUsername || !hasPassword || strings.TrimSpace(username) == "" || password == "" {
		return nil, nil, &AuthenticationError{
			Message: "环境变量 THS_USERNAME 和 THS_PASSWORD 必须同时设置且不能为空",
		}
	}
	return &username, &password, nil
}

// Client 是进程本地 THSDK 客户端，对应 Python _Client。
type Client struct {
	sessionDir    string
	runtime       *RuntimeBridge
	stateLock     sync.Mutex
	authLock      sync.Mutex
	authenticated bool
}

// NewClient 创建一个客户端实例。sessionDir 为空时使用进程当前工作目录。
func NewClient(sessionDir string, runtime *RuntimeBridge) *Client {
	dir := sessionDir
	if dir == "" {
		dir, _ = os.Getwd()
	}
	if absolute, err := filepath.Abs(dir); err == nil {
		dir = absolute
	}
	return &Client{sessionDir: filepath.Clean(dir), runtime: runtime}
}

// SessionDir 返回会话目录。
func (c *Client) SessionDir() string {
	return c.sessionDir
}

// Runtime 惰性创建并返回运行桥，对应 Python _Client.runtime 属性。
func (c *Client) Runtime() (*RuntimeBridge, error) {
	c.stateLock.Lock()
	defer c.stateLock.Unlock()
	if c.runtime == nil {
		bridge, err := NewRuntimeBridge("")
		if err != nil {
			return nil, err
		}
		c.runtime = bridge
	}
	return c.runtime, nil
}

// extractError 解析响应错误，对应 Python _Client._error。
// 返回 (code, message, retry_after_ms)。
func extractError(response map[string]any) (string, string, *int) {
	errorValue, _ := response["error"].(map[string]any)
	if errorValue == nil {
		return "unknown_error", fmt.Sprint(response["error"]), nil
	}
	var retryAfter *int
	if value, ok := asInt64(errorValue["retry_after_ms"]); ok && value > 0 {
		converted := int(value)
		retryAfter = &converted
	}
	code, _ := errorValue["code"].(string)
	if code == "" {
		code = "unknown_error"
	}
	message, _ := errorValue["message"].(string)
	if message == "" {
		message = "unknown error"
	}
	return code, message, retryAfter
}

// invoke 执行一次方法调用，含一次限流重试与截止时间控制，
// 对应 Python _Client._invoke。
func (c *Client) invoke(method string, params any, timeoutMs *int) (any, error) {
	var deadline time.Time
	if timeoutMs != nil {
		deadline = time.Now().Add(time.Duration(*timeoutMs) * time.Millisecond)
	}
	callTimeoutMs := timeoutMs
	for attempt := 0; attempt < 2; attempt++ {
		bridge, err := c.Runtime()
		if err != nil {
			return nil, err
		}
		response, err := bridge.Call(method, params, callTimeoutMs)
		if err != nil {
			return nil, err
		}
		if ok, _ := response["ok"].(bool); ok {
			return normalizeJSON(response["data"]), nil
		}
		code, message, retryAfterMs := extractError(response)
		canRetry := attempt == 0 &&
			code == "rate_limited" &&
			retryAfterMs != nil &&
			*retryAfterMs <= maxRateLimitRetryMS
		if canRetry && !deadline.IsZero() {
			canRetry = time.Now().Add(time.Duration(*retryAfterMs) * time.Millisecond).Before(deadline)
		}
		if !canRetry {
			return nil, &APIError{Code: code, Message: message, RetryAfterMs: retryAfterMs}
		}
		time.Sleep(time.Duration(*retryAfterMs) * time.Millisecond)
		if !deadline.IsZero() {
			remaining := time.Until(deadline)
			if remaining <= 0 {
				return nil, &APIError{Code: code, Message: message, RetryAfterMs: retryAfterMs}
			}
			milliseconds := timeoutMilliseconds(remaining.Seconds())
			callTimeoutMs = &milliseconds
		}
	}
	return nil, newTHSDKError("unreachable")
}

// AuthOptions 控制 Auth 调用的超时（秒），默认 60。
type AuthOptions struct {
	Timeout float64
}

// WithAuthTimeout 返回携带指定认证超时（秒）的选项。
func WithAuthTimeout(seconds float64) *AuthOptions {
	return &AuthOptions{Timeout: seconds}
}

func authTimeout(opts []*AuthOptions) float64 {
	for _, option := range opts {
		if option != nil && option.Timeout > 0 {
			return option.Timeout
		}
	}
	return 60
}

// Auth 按 显式凭据 → 环境变量 → 已有会话 → 临时账户 的顺序完成认证，
// 对应 Python _Client.auth。
//
// username/password 同时为空时走环境变量/会话/临时账户；mac 仅在提供
// 显式账号密码时生效，为空时按账号派生稳定 MAC。
func (c *Client) Auth(username, password, mac string, opts ...*AuthOptions) (bool, error) {
	timeout, err := validateSeconds(authTimeout(opts), "timeout", 1e-9, maxTimeoutSeconds)
	if err != nil {
		return false, err
	}
	if (username == "") != (password == "") {
		return false, &AuthenticationError{Message: "账号和密码必须同时提供"}
	}
	if username != "" && (strings.TrimSpace(username) == "" || password == "") {
		return false, &AuthenticationError{Message: "账号和密码不能为空"}
	}
	if mac != "" && username == "" {
		return false, &AuthenticationError{Message: "mac 只能与账号和密码同时提供"}
	}

	type credentials struct {
		account string
		secret  string
		device  *string
	}
	var creds *credentials
	if username != "" {
		var device *string
		if mac != "" {
			device = &mac
		}
		creds = &credentials{account: username, secret: password, device: device}
	} else {
		envUser, envPass, envErr := environmentCredentials()
		if envErr != nil {
			return false, envErr
		}
		if envUser != nil {
			creds = &credentials{account: *envUser, secret: *envPass, device: nil}
		}
	}

	deadline := time.Now().Add(time.Duration(timeout * float64(time.Second)))

	authenticate := func(params map[string]any) (bool, error) {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return false, &AuthenticationError{Message: "认证超时"}
		}
		milliseconds := timeoutMilliseconds(remaining.Seconds())
		data, err := c.invoke("auth", params, &milliseconds)
		if err != nil {
			return false, err
		}
		dataMap, _ := data.(map[string]any)
		if dataMap == nil || dataMap["authenticated"] != true {
			return false, &AuthenticationError{Message: "认证未完成"}
		}
		c.stateLock.Lock()
		c.authenticated = true
		c.stateLock.Unlock()
		return true, nil
	}

	accountParams := func(account, secret string, device *string) map[string]any {
		derived := deriveAccountMAC(account)
		if device != nil {
			derived = *device
		}
		return map[string]any{
			"session_dir": c.sessionDir,
			"username":    account,
			"password":    secret,
			"mac":         derived,
		}
	}

	c.authLock.Lock()
	defer c.authLock.Unlock()

	if creds != nil {
		authenticated, err := authenticate(accountParams(creds.account, creds.secret, creds.device))
		if err != nil {
			if apiErr, ok := err.(*APIError); ok {
				return false, &AuthenticationError{Message: apiErr.Error()}
			}
			return false, err
		}
		return authenticated, nil
	}

	c.stateLock.Lock()
	authenticated := c.authenticated
	c.stateLock.Unlock()
	if authenticated {
		return true, nil
	}

	sessionParams := map[string]any{"session_dir": c.sessionDir}
	if _, err := os.Stat(filepath.Join(c.sessionDir, "account.session")); err == nil {
		authenticated, authErr := authenticate(sessionParams)
		if authErr != nil {
			apiErr, isAPI := authErr.(*APIError)
			if !isAPI || apiErr.Code != "session_unavailable" {
				return false, &AuthenticationError{Message: authErr.Error()}
			}
		} else {
			return authenticated, nil
		}
	}

	guestUser, guestPass, guestMac := acquireTemporaryAccount()
	authenticated, err = authenticate(accountParams(guestUser, guestPass, &guestMac))
	if err != nil {
		if apiErr, ok := err.(*APIError); ok {
			return false, &AuthenticationError{Message: apiErr.Error()}
		}
		return false, err
	}
	return authenticated, nil
}

// AuthQRCodeOptions 控制二维码登录的参数，对应 Python auth_qrcode 关键字参数。
type AuthQRCodeOptions struct {
	Timeout      float64 // 默认 120 秒
	PollInterval float64 // 默认 1 秒
}

// WithAuthQRCodeTimeout 设置二维码登录总超时（秒）。
func WithAuthQRCodeTimeout(seconds float64) *AuthQRCodeOptions {
	return &AuthQRCodeOptions{Timeout: seconds}
}

// WithAuthQRCodePollInterval 设置二维码轮询间隔（秒）。
func WithAuthQRCodePollInterval(seconds float64) *AuthQRCodeOptions {
	return &AuthQRCodeOptions{PollInterval: seconds}
}

func authQRCodeOptions(opts []*AuthQRCodeOptions) (timeout, pollInterval float64) {
	timeout = 120
	pollInterval = 1
	for _, option := range opts {
		if option == nil {
			continue
		}
		if option.Timeout > 0 {
			timeout = option.Timeout
		}
		if option.PollInterval > 0 {
			pollInterval = option.PollInterval
		}
	}
	return timeout, pollInterval
}

// AuthQRCode 优先恢复本地会话，否则完成一次二维码登录任务，
// 对应 Python _Client.auth_qrcode。
func (c *Client) AuthQRCode(opts ...*AuthQRCodeOptions) (bool, error) {
	timeout, pollInterval := authQRCodeOptions(opts)
	var err error
	timeout, err = validateSeconds(timeout, "timeout", 1e-9, maxTimeoutSeconds)
	if err != nil {
		return false, err
	}
	pollInterval, err = validateSeconds(pollInterval, "poll_interval", minQRPollSeconds, maxQRPollSeconds)
	if err != nil {
		return false, err
	}

	deadline := time.Now().Add(time.Duration(timeout * float64(time.Second)))
	startParams := map[string]any{
		"session_dir":      c.sessionDir,
		"timeout_seconds":  int(math.Ceil(timeout)),
		"poll_interval_ms": timeoutMilliseconds(pollInterval),
	}

	c.authLock.Lock()
	defer c.authLock.Unlock()

	remainingMS := timeoutMilliseconds(time.Until(deadline).Seconds())
	if remainingMS > 20000 {
		remainingMS = 20000
	}
	state, err := c.invoke("auth_qrcode", startParams, &remainingMS)
	if err != nil {
		return false, &AuthenticationError{Message: err.Error()}
	}
	stateMap, _ := state.(map[string]any)

	if authenticatedByState(stateMap) {
		c.stateLock.Lock()
		c.authenticated = true
		c.stateLock.Unlock()
		return true, nil
	}

	var shownChallenge [2]string
	showChallenge := func(value map[string]any) {
		challenge := [2]string{fmt.Sprint(value["scan_url"]), fmt.Sprint(value["image_url"])}
		if (challenge[0] == "" && challenge[1] == "") || challenge == shownChallenge {
			return
		}
		shownChallenge = challenge
		printQRChallenge(value, time.Until(deadline).Seconds(), nil)
	}

	completed := false
	for {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			break
		}
		milliseconds := timeoutMilliseconds(remaining.Seconds())
		state, err = c.invoke("auth_qrcode", nil, &milliseconds)
		if err != nil {
			return false, &AuthenticationError{Message: err.Error()}
		}
		stateMap, _ = state.(map[string]any)
		showChallenge(stateMap)
		if authenticatedByState(stateMap) {
			c.stateLock.Lock()
			c.authenticated = true
			c.stateLock.Unlock()
			completed = true
			return true, nil
		}
		remaining = time.Until(deadline)
		if remaining <= 0 {
			break
		}
		sleepSeconds := pollInterval
		if sleepSeconds < 0.1 {
			sleepSeconds = 0.1
		}
		if sleepSeconds > remaining.Seconds() {
			sleepSeconds = remaining.Seconds()
		}
		time.Sleep(time.Duration(sleepSeconds * float64(time.Second)))
	}

	if !completed {
		_, _ = c.invoke("auth_qrcode", map[string]any{"cancel": true}, nil)
	}
	return false, &AuthenticationError{Message: "二维码登录超时"}
}

// authenticatedByState 判断认证状态，对应 Python 的
// state.get("authenticated") or state.get("state") == "authenticated"。
func authenticatedByState(state map[string]any) bool {
	if state == nil {
		return false
	}
	if ok, _ := state["authenticated"].(bool); ok {
		return true
	}
	value, _ := state["state"].(string)
	return value == "authenticated"
}

// Request 在显式认证后执行一次包定义的数据请求，对应 Python _Client._request。
func (c *Client) Request(method string, params any, opts ...*CallOptions) (any, error) {
	c.authLock.Lock()
	defer c.authLock.Unlock()

	c.stateLock.Lock()
	authenticated := c.authenticated
	c.stateLock.Unlock()
	if !authenticated {
		return nil, &NotAuthenticatedError{
			Message: "尚未认证；请先调用 AuthQRCode() 或 Auth(...)",
		}
	}

	var timeoutMs *int
	for _, option := range opts {
		if option != nil && option.Timeout > 0 {
			seconds, err := validateSeconds(option.Timeout, "timeout", 1e-9, maxTimeoutSeconds)
			if err != nil {
				return nil, err
			}
			milliseconds := timeoutMilliseconds(seconds)
			timeoutMs = &milliseconds
			break
		}
	}
	return c.invoke(method, params, timeoutMs)
}

// Logout 退出当前账号并删除会话，对应 Python _Client.logout。
func (c *Client) Logout() error {
	c.authLock.Lock()
	defer c.authLock.Unlock()
	err := func() error {
		_, invokeErr := c.invoke("logout", map[string]any{"session_dir": c.sessionDir}, nil)
		return invokeErr
	}()
	c.stateLock.Lock()
	c.authenticated = false
	c.stateLock.Unlock()
	return err
}

// Close 关闭运行组件连接，对应 Python _Client._close。
func (c *Client) Close() error {
	c.authLock.Lock()
	defer c.authLock.Unlock()
	err := func() error {
		_, invokeErr := c.invoke("_close", nil, nil)
		return invokeErr
	}()
	c.stateLock.Lock()
	c.authenticated = false
	c.stateLock.Unlock()
	return err
}
