package thsdk

// 本文件实现会话层：进程级客户端、认证与数据请求调度。
//
//   - Client：进程内逻辑客户端，持有运行桥（惰性创建）与认证状态；
//   - NewClient：创建客户端实例；
//   - Auth / AuthQRCode / Logout：三种认证入口与退出；
//   - Request：已认证后的统一数据请求入口（含限流重试与超时控制）；
//   - 辅助：账号 MAC 派生、环境变量凭据、秒/毫秒换算与参数校验。
//
// 认证是显式行为：数据请求绝不会在调用方不知情的情况下静默恢复
// 磁盘会话，必须先 Auth / AuthQRCode 成功。

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
	// accountMACSeed 是账号 MAC 派生的固定种子字节串（末尾含 NUL），
	// 与原生组件约定逐字一致，不可改动。
	accountMACSeed = "thsdk-account-mac-v1\x00"
	// envUsername 是账号环境变量名。
	envUsername = "THS_USERNAME"
	// envPassword 是密码环境变量名。
	envPassword = "THS_PASSWORD"
	// maxTimeoutSeconds 是所有调用超时的上限（600 秒）。
	maxTimeoutSeconds = 600.0
	// minQRPollSeconds / maxQRPollSeconds 是二维码轮询间隔的边界（秒）。
	minQRPollSeconds = 0.5
	maxQRPollSeconds = 30.0
	// maxRateLimitRetryMS 是限流重试允许的最大等待毫秒数（1 秒）。
	maxRateLimitRetryMS = 1000
)

// validateSeconds 校验秒级参数范围。
//
// 入参：
//   - value：待校验的秒数；
//   - name：参数名（用于错误提示）；
//   - minimum / maximum：允许的闭区间边界。
//
// 出参：
//   - float64：校验通过的原值；
//   - error：NaN / Infinity 或越界时返回错误。
func validateSeconds(value float64, name string, minimum, maximum float64) (float64, error) {
	// NaN/Infinity 与非闭区间值一律拒绝。
	if math.IsNaN(value) || math.IsInf(value, 0) || value < minimum || value > maximum {
		return 0, fmt.Errorf("%s 必须在 %g 到 %g 秒之间", name, minimum, maximum)
	}
	return value, nil
}

// timeoutMilliseconds 把秒向上取整为毫秒，且至少 1ms。
//
// 入参：
//   - value：秒数。
//
// 出参：
//   - int：毫秒数（最小 1）。
//
// 说明：原生协议要求正整数毫秒，0 会视为未指定，因此向上取整并保底 1。
func timeoutMilliseconds(value float64) int {
	milliseconds := int(math.Ceil(value * 1000))
	if milliseconds < 1 {
		return 1
	}
	return milliseconds
}

// deriveAccountMAC 把一个账号名映射为稳定的本地管理单播 MAC。
//
// 入参：
//   - username：账号名。
//
// 出参：
//   - string：形如 "02:xx:xx:xx:xx:xx" 的 MAC 字符串。
//
// 规则：SHA-256(种子+账号) 取前 6 字节；首字节强制本地管理单播位
// （清除 U/L 位、置 L 位），保证 MAC 稳定且不冲突。
func deriveAccountMAC(username string) string {
	// 计算摘要并取前 6 字节。
	digest := sha256.Sum256([]byte(accountMACSeed + username))
	octets := append([]byte(nil), digest[:6]...)
	// 首字节置为本地管理单播（0x02 最低两位置位/清零规则）。
	octets[0] = (octets[0] & 0xFC) | 0x02
	// 格式化为冒号分隔的十六进制。
	parts := make([]string, len(octets))
	for i, value := range octets {
		parts[i] = fmt.Sprintf("%02x", value)
	}
	return strings.Join(parts, ":")
}

// environmentCredentials 读取环境变量凭据。
//
// 入参：无。
//
// 出参：
//   - *string：账号指针；未设置环境变量时为 nil；
//   - *string：密码指针；未设置时为 nil；
//   - error：只设置了其中一个变量或值为空时返回认证错误。
//
// 规则：两个变量必须同时存在且非空，否则报错。
func environmentCredentials() (*string, *string, error) {
	username, hasUsername := os.LookupEnv(envUsername)
	password, hasPassword := os.LookupEnv(envPassword)
	// 两者都未设置：返回空凭据（不视为错误）。
	if !hasUsername && !hasPassword {
		return nil, nil, nil
	}
	// 只设置其一或值不合法：视为配置错误。
	if !hasUsername || !hasPassword || strings.TrimSpace(username) == "" || password == "" {
		return nil, nil, &AuthenticationError{
			Message: "环境变量 THS_USERNAME 和 THS_PASSWORD 必须同时设置且不能为空",
		}
	}
	return &username, &password, nil
}

// Client 是进程本地的 THSDK 客户端。
//
// 字段说明：
//   - sessionDir：会话持久化目录（默认进程工作目录）；
//   - runtime：原生运行桥（惰性创建，首次调用时加载）；
//   - stateLock：保护认证标记的互斥锁；
//   - authLock：串行化认证/请求流程的互斥锁；
//   - authenticated：当前是否持有有效会话的标记。
type Client struct {
	sessionDir    string
	runtime       *RuntimeBridge
	stateLock     sync.Mutex
	authLock      sync.Mutex
	authenticated bool
}

// NewClient 创建一个客户端实例。
//
// 入参：
//   - sessionDir：会话持久化目录；传空串时使用进程当前工作目录；
//   - runtime：预置的运行桥；传 nil 时惰性创建。
//
// 出参：
//   - *Client：创建好的客户端实例。
//
// 说明：会话目录会规范化为绝对路径。
func NewClient(sessionDir string, runtime *RuntimeBridge) *Client {
	dir := sessionDir
	if dir == "" {
		// 缺省为当前工作目录。
		dir, _ = os.Getwd()
	}
	if absolute, err := filepath.Abs(dir); err == nil {
		dir = absolute
	}
	return &Client{sessionDir: filepath.Clean(dir), runtime: runtime}
}

// SessionDir 返回会话目录。
//
// 入参：无。
//
// 出参：
//   - string：会话持久化目录路径。
func (c *Client) SessionDir() string {
	return c.sessionDir
}

// Runtime 惰性创建并返回运行桥。
//
// 入参：无。
//
// 出参：
//   - *RuntimeBridge：已握手校验的运行桥；
//   - error：组件加载或握手失败时返回错误。
//
// 说明：首次调用时创建并缓存；并发调用由 stateLock 保护。
func (c *Client) Runtime() (*RuntimeBridge, error) {
	c.stateLock.Lock()
	defer c.stateLock.Unlock()
	// 未创建时按自动查找路径加载。
	if c.runtime == nil {
		bridge, err := NewRuntimeBridge("")
		if err != nil {
			return nil, err
		}
		c.runtime = bridge
	}
	return c.runtime, nil
}

// extractError 从原生响应中解析错误三元组。
//
// 入参：
//   - response：原生解析后的响应结构。
//
// 出参：
//   - code：错误码；响应无 error 对象时返回 "unknown_error"；
//   - message：错误描述；为空时返回 "unknown error"；
//   - retryAfterMs：建议重试等待毫秒数；仅响应提供有效正整数值时非空。
func extractError(response map[string]any) (string, string, *int) {
	// 取 error 子对象；缺失时按未知错误处理。
	errorValue, _ := response["error"].(map[string]any)
	if errorValue == nil {
		return "unknown_error", fmt.Sprint(response["error"]), nil
	}
	// 解析可选的重试等待毫秒数（仅接受正整数）。
	var retryAfter *int
	if value, ok := asInt64(errorValue["retry_after_ms"]); ok && value > 0 {
		converted := int(value)
		retryAfter = &converted
	}
	// 解析错误码，缺省为 unknown_error。
	code, _ := errorValue["code"].(string)
	if code == "" {
		code = "unknown_error"
	}
	// 解析错误描述，缺省为 unknown error。
	message, _ := errorValue["message"].(string)
	if message == "" {
		message = "unknown error"
	}
	return code, message, retryAfter
}

// invoke 执行一次方法调用，含一次限流重试与截止时间控制。
//
// 入参：
//   - method：原生方法名；
//   - params：请求参数（可为 nil）；
//   - timeoutMs：调用超时毫秒数指针（可为 nil）。
//
// 出参：
//   - any：成功时返回归一化后的 data 字段；
//   - error：调用失败或重试不可行时返回 APIError。
//
// 重试规则：首次调用若返回 rate_limited 且等待时间不超过 1 秒、
// 且（若设了截止时间）等待后仍在截止时间内，则睡眠后重试一次。
func (c *Client) invoke(method string, params any, timeoutMs *int) (any, error) {
	// 计算调用截止时间（用于限流重试的剩余时间控制）。
	var deadline time.Time
	if timeoutMs != nil {
		deadline = time.Now().Add(time.Duration(*timeoutMs) * time.Millisecond)
	}
	// 最多尝试两次（含一次限流重试）。
	callTimeoutMs := timeoutMs
	for attempt := 0; attempt < 2; attempt++ {
		// 获取运行桥（惰性加载）。
		bridge, err := c.Runtime()
		if err != nil {
			return nil, err
		}
		// 发起原生调用。
		response, err := bridge.Call(method, params, callTimeoutMs)
		if err != nil {
			return nil, err
		}
		// 成功：返回归一化后的 data。
		if ok, _ := response["ok"].(bool); ok {
			return normalizeJSON(response["data"]), nil
		}
		// 失败：解析错误三元组。
		code, message, retryAfterMs := extractError(response)
		// 判断是否允许重试：首次尝试 + rate_limited + 等待 ≤ 1 秒。
		canRetry := attempt == 0 &&
			code == "rate_limited" &&
			retryAfterMs != nil &&
			*retryAfterMs <= maxRateLimitRetryMS
		// 设了截止时间时，重试后必须在截止时间内。
		if canRetry && !deadline.IsZero() {
			canRetry = time.Now().Add(time.Duration(*retryAfterMs) * time.Millisecond).Before(deadline)
		}
		// 不可重试：直接返回类型化错误。
		if !canRetry {
			return nil, &APIError{Code: code, Message: message, RetryAfterMs: retryAfterMs}
		}
		// 等待后重试，并收缩剩余超时。
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
	// 理论上不可达（循环固定两次）。
	return nil, newTHSDKError("unreachable")
}

// AuthOptions 控制 Auth 调用的超时（秒），默认 60。
type AuthOptions struct {
	Timeout float64
}

// WithAuthTimeout 构造带指定认证超时的选项。
//
// 入参：
//   - seconds：认证超时秒数。
//
// 出参：
//   - *AuthOptions：可直接传入 Auth 的 opts 参数。
func WithAuthTimeout(seconds float64) *AuthOptions {
	return &AuthOptions{Timeout: seconds}
}

// authTimeout 从选项列表取认证超时，缺省 60 秒。
//
// 入参：
//   - opts：选项列表。
//
// 出参：
//   - float64：首个合法正超时；全部无效时返回 60。
func authTimeout(opts []*AuthOptions) float64 {
	for _, option := range opts {
		if option != nil && option.Timeout > 0 {
			return option.Timeout
		}
	}
	return 60
}

// Auth 按「显式凭据 → 环境变量 → 已有会话 → 临时账户」顺序完成认证。
//
// 入参：
//   - username / password：账号与密码；可同时为空以走后续路径；
//   - mac：设备 MAC；仅与显式账号密码同时提供时生效，为空时按账号
//     派生稳定 MAC；
//   - opts：认证选项（超时）。
//
// 出参：
//   - bool：是否认证成功；
//   - error：认证失败时返回 AuthenticationError 或底层错误。
//
// 说明：username/password 只提供一个会报错；已认证状态下再次调用
// （无凭据）直接返回 true。
func (c *Client) Auth(username, password, mac string, opts ...*AuthOptions) (bool, error) {
	// 校验认证超时。
	timeout, err := validateSeconds(authTimeout(opts), "timeout", 1e-9, maxTimeoutSeconds)
	if err != nil {
		return false, err
	}
	// 账号密码必须同时提供或不提供。
	if (username == "") != (password == "") {
		return false, &AuthenticationError{Message: "账号和密码必须同时提供"}
	}
	// 显式账号不能为空。
	if username != "" && (strings.TrimSpace(username) == "" || password == "") {
		return false, &AuthenticationError{Message: "账号和密码不能为空"}
	}
	// mac 必须伴随显式账号密码。
	if mac != "" && username == "" {
		return false, &AuthenticationError{Message: "mac 只能与账号和密码同时提供"}
	}

	// 凭据结构：显式或环境变量来源。
	type credentials struct {
		account string
		secret  string
		device  *string
	}
	var creds *credentials
	if username != "" {
		// 显式凭据；mac 可选。
		var device *string
		if mac != "" {
			device = &mac
		}
		creds = &credentials{account: username, secret: password, device: device}
	} else {
		// 环境变量凭据（可能未设置）。
		envUser, envPass, envErr := environmentCredentials()
		if envErr != nil {
			return false, envErr
		}
		if envUser != nil {
			creds = &credentials{account: *envUser, secret: *envPass, device: nil}
		}
	}

	// 认证截止时间（用于限流重试与超时控制）。
	deadline := time.Now().Add(time.Duration(timeout * float64(time.Second)))

	// authenticate 执行一次 auth 调用并置认证标记。
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
		// 响应必须带 authenticated=true。
		dataMap, _ := data.(map[string]any)
		if dataMap == nil || dataMap["authenticated"] != true {
			return false, &AuthenticationError{Message: "认证未完成"}
		}
		// 置认证标记。
		c.stateLock.Lock()
		c.authenticated = true
		c.stateLock.Unlock()
		return true, nil
	}

	// accountParams 组装账号认证参数（含派生或显式 MAC）。
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

	// 串行化认证流程。
	c.authLock.Lock()
	defer c.authLock.Unlock()

	// 有凭据：直接认证。
	if creds != nil {
		authenticated, err := authenticate(accountParams(creds.account, creds.secret, creds.device))
		if err != nil {
			// 原生 API 错误统一转认证错误。
			if apiErr, ok := err.(*APIError); ok {
				return false, &AuthenticationError{Message: apiErr.Error()}
			}
			return false, err
		}
		return authenticated, nil
	}

	// 无凭据：先检查本地会话标记。
	c.stateLock.Lock()
	authenticated := c.authenticated
	c.stateLock.Unlock()
	if authenticated {
		return true, nil
	}

	// 尝试恢复磁盘会话（account.session 存在时）。
	sessionParams := map[string]any{"session_dir": c.sessionDir}
	if _, err := os.Stat(filepath.Join(c.sessionDir, "account.session")); err == nil {
		authenticated, authErr := authenticate(sessionParams)
		if authErr != nil {
			// 会话不可用时忽略，继续走临时账户。
			apiErr, isAPI := authErr.(*APIError)
			if !isAPI || apiErr.Code != "session_unavailable" {
				return false, &AuthenticationError{Message: authErr.Error()}
			}
		} else {
			return authenticated, nil
		}
	}

	// 兜底：注册临时账户完成认证。
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

// AuthQRCodeOptions 控制二维码登录的参数。
//
// 字段说明：
//   - Timeout：登录总超时（秒），默认 120；
//   - PollInterval：轮询间隔（秒），默认 1。
type AuthQRCodeOptions struct {
	Timeout      float64
	PollInterval float64
}

// WithAuthQRCodeTimeout 设置二维码登录总超时（秒）。
//
// 入参：
//   - seconds：超时秒数。
//
// 出参：
//   - *AuthQRCodeOptions：可直接传入 AuthQRCode 的 opts 参数。
func WithAuthQRCodeTimeout(seconds float64) *AuthQRCodeOptions {
	return &AuthQRCodeOptions{Timeout: seconds}
}

// WithAuthQRCodePollInterval 设置二维码轮询间隔（秒）。
//
// 入参：
//   - seconds：轮询秒数。
//
// 出参：
//   - *AuthQRCodeOptions：可直接传入 AuthQRCode 的 opts 参数。
func WithAuthQRCodePollInterval(seconds float64) *AuthQRCodeOptions {
	return &AuthQRCodeOptions{PollInterval: seconds}
}

// authQRCodeOptions 合并选项列表为超时与轮询间隔。
//
// 入参：
//   - opts：选项列表。
//
// 出参：
//   - timeout：登录总超时（默认 120 秒）；
//   - pollInterval：轮询间隔（默认 1 秒）。
//
// 规则：逐项取首个正值的超时/间隔。
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

// AuthQRCode 优先恢复本地会话，否则完成一次二维码登录任务。
//
// 入参：
//   - opts：二维码登录选项（超时 / 轮询间隔）。
//
// 出参：
//   - bool：是否登录成功；
//   - error：超时或调用失败时返回认证错误。
//
// 流程：发起 auth_qrcode 任务 → 轮询状态 → 扫码完成（authenticated）
// 即成功；超时前未完成则取消任务并返回错误。二维码在首次出现时
// 通过 printQRChallenge 渲染到终端。
func (c *Client) AuthQRCode(opts ...*AuthQRCodeOptions) (bool, error) {
	// 解析并校验选项。
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

	// 截止时间与启动参数。
	deadline := time.Now().Add(time.Duration(timeout * float64(time.Second)))
	startParams := map[string]any{
		"session_dir":      c.sessionDir,
		"timeout_seconds":  int(math.Ceil(timeout)),
		"poll_interval_ms": timeoutMilliseconds(pollInterval),
	}

	// 串行化认证流程。
	c.authLock.Lock()
	defer c.authLock.Unlock()

	// 启动二维码登录任务（首轮调用最多等待 20 秒出图）。
	remainingMS := timeoutMilliseconds(time.Until(deadline).Seconds())
	if remainingMS > 20000 {
		remainingMS = 20000
	}
	state, err := c.invoke("auth_qrcode", startParams, &remainingMS)
	if err != nil {
		return false, &AuthenticationError{Message: err.Error()}
	}
	stateMap, _ := state.(map[string]any)

	// 启动即已完成认证（如会话恢复）。
	if authenticatedByState(stateMap) {
		c.stateLock.Lock()
		c.authenticated = true
		c.stateLock.Unlock()
		return true, nil
	}

	// 展示二维码（仅首次出现或内容变化时打印一次）。
	var shownChallenge [2]string
	showChallenge := func(value map[string]any) {
		challenge := [2]string{fmt.Sprint(value["scan_url"]), fmt.Sprint(value["image_url"])}
		if (challenge[0] == "" && challenge[1] == "") || challenge == shownChallenge {
			return
		}
		shownChallenge = challenge
		printQRChallenge(value, time.Until(deadline).Seconds(), nil)
	}

	// 轮询直到扫码完成或超时。
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
		// 扫码完成即成功。
		if authenticatedByState(stateMap) {
			c.stateLock.Lock()
			c.authenticated = true
			c.stateLock.Unlock()
			completed = true
			return true, nil
		}
		// 未完成：按轮询间隔睡眠（不超过剩余时间）。
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

	// 超时：尝试取消任务，返回超时错误。
	if !completed {
		_, _ = c.invoke("auth_qrcode", map[string]any{"cancel": true}, nil)
	}
	return false, &AuthenticationError{Message: "二维码登录超时"}
}

// authenticatedByState 判断二维码登录状态是否为已认证。
//
// 入参：
//   - state：原生返回的状态对象。
//
// 出参：
//   - bool：authenticated 为 true 或 state == "authenticated" 时返回 true。
func authenticatedByState(state map[string]any) bool {
	// 空状态视为未认证。
	if state == nil {
		return false
	}
	// 布尔标记优先。
	if ok, _ := state["authenticated"].(bool); ok {
		return true
	}
	// 字符串状态兜底。
	value, _ := state["state"].(string)
	return value == "authenticated"
}

// Request 在显式认证后执行一次包定义的数据请求。
//
// 入参：
//   - method：原生方法名；
//   - params：请求参数（可为 nil）；
//   - opts：调用选项（超时，秒）。
//
// 出参：
//   - any：成功时返回归一化后的 data 字段；
//   - error：未认证返回 NotAuthenticatedError；其余为调用错误。
//
// 说明：所有公开数据方法最终都经由它发起请求；认证检查持锁保证
// 与认证流程互斥。
func (c *Client) Request(method string, params any, opts ...*CallOptions) (any, error) {
	// 与认证流程互斥。
	c.authLock.Lock()
	defer c.authLock.Unlock()

	// 未认证时拒绝数据请求。
	c.stateLock.Lock()
	authenticated := c.authenticated
	c.stateLock.Unlock()
	if !authenticated {
		return nil, &NotAuthenticatedError{
			Message: "尚未认证；请先调用 AuthQRCode() 或 Auth(...)",
		}
	}

	// 解析调用超时（秒 → 毫秒）。
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

// Logout 退出当前账号并删除会话。
//
// 入参：无。
//
// 出参：
//   - error：原生 logout 调用失败时返回错误；成功返回 nil。
//
// 说明：无论原生调用是否成功，都会清除本地认证标记。
func (c *Client) Logout() error {
	// 串行化认证流程。
	c.authLock.Lock()
	defer c.authLock.Unlock()
	// 调用原生 logout。
	err := func() error {
		_, invokeErr := c.invoke("logout", map[string]any{"session_dir": c.sessionDir}, nil)
		return invokeErr
	}()
	// 清除本地认证标记。
	c.stateLock.Lock()
	c.authenticated = false
	c.stateLock.Unlock()
	return err
}

// Close 关闭运行组件连接并清除认证状态。
//
// 入参：无。
//
// 出参：
//   - error：原生 _close 调用失败时返回错误；成功返回 nil。
//
// 说明：可在程序结束前 defer 调用，释放原生运行组件资源。
func (c *Client) Close() error {
	// 串行化认证流程。
	c.authLock.Lock()
	defer c.authLock.Unlock()
	// 调用原生 _close。
	err := func() error {
		_, invokeErr := c.invoke("_close", nil, nil)
		return invokeErr
	}()
	// 清除本地认证标记。
	c.stateLock.Lock()
	c.authenticated = false
	c.stateLock.Unlock()
	return err
}
