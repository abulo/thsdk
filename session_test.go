package thsdk

import (
	"os"
	"strings"
	"testing"
)

func TestValidateSeconds(t *testing.T) {
	value, err := validateSeconds(5, "t", 0.5, 600)
	if err != nil || value != 5 {
		t.Fatalf("正常值应通过: %v %v", value, err)
	}
	if _, err := validateSeconds(-1, "t", 0.5, 600); err == nil {
		t.Fatal("小于最小值应报错")
	}
	if _, err := validateSeconds(601, "t", 0.5, 600); err == nil {
		t.Fatal("大于最大值应报错")
	}
	if _, err := validateSeconds(0.1, "t", 0.5, 600); err == nil {
		t.Fatal("边界值应报错")
	}
}

func TestTimeoutMilliseconds(t *testing.T) {
	if timeoutMilliseconds(0.0001) != 1 {
		t.Fatal("至少 1ms")
	}
	if timeoutMilliseconds(1.2) != 1200 {
		t.Fatal("向上取整")
	}
}

func TestDeriveAccountMAC(t *testing.T) {
	mac := deriveAccountMAC("test_user")
	parts := strings.Split(mac, ":")
	if len(parts) != 6 {
		t.Fatalf("MAC 应有 6 组: %s", mac)
	}
	first, _ := strconvParseInt(parts[0], 16, 8)
	if first&0x02 != 0x02 || first&0x01 != 0 {
		t.Fatalf("MAC 应为本地管理单播: %s", mac)
	}
	// 幂等性
	if deriveAccountMAC("test_user") != mac {
		t.Fatal("MAC 派生必须稳定")
	}
	if deriveAccountMAC("other") == mac {
		t.Fatal("不同账号必须派生不同 MAC")
	}
}

func TestEnvironmentCredentials(t *testing.T) {
	// 显式清除两个环境变量（t.Setenv 空串会使键存在，模拟不出"未设置"）。
	os.Unsetenv("THS_USERNAME")
	os.Unsetenv("THS_PASSWORD")
	user, pass, err := environmentCredentials()
	if err != nil || user != nil || pass != nil {
		t.Fatalf("双空应返回 nil: %v %v %v", user, pass, err)
	}

	t.Setenv("THS_USERNAME", "u")
	t.Setenv("THS_PASSWORD", "")
	if _, _, err := environmentCredentials(); err == nil {
		t.Fatal("只设置一个变量应报错")
	}

	t.Setenv("THS_USERNAME", "u")
	t.Setenv("THS_PASSWORD", "p")
	user, pass, err = environmentCredentials()
	if err != nil || user == nil || *user != "u" || pass == nil || *pass != "p" {
		t.Fatalf("完整凭据读取失败: %v %v %v", user, pass, err)
	}
}

func TestAuthWithStubBridge(t *testing.T) {
	script := &scriptedBridge{responses: map[string]string{
		"auth": okResponse(map[string]any{"authenticated": true}),
	}}
	client := script.client()
	ok, err := client.Auth("alice", "secret", "")
	if err != nil || !ok {
		t.Fatalf("认证失败: %v %v", ok, err)
	}
	request := script.lastRequest()
	params := request["params"].(map[string]any)
	if params["username"] != "alice" || params["password"] != "secret" {
		t.Fatalf("认证参数不符: %v", params)
	}
	if params["session_dir"] != client.SessionDir() {
		t.Fatalf("会话目录不符: %v", params["session_dir"])
	}
	if params["mac"] != deriveAccountMAC("alice") {
		t.Fatalf("默认 MAC 不符: %v", params["mac"])
	}
}

func TestAuthExplicitMAC(t *testing.T) {
	script := &scriptedBridge{responses: map[string]string{
		"auth": okResponse(map[string]any{"authenticated": true}),
	}}
	client := script.client()
	if _, err := client.Auth("alice", "secret", "02:00:00:00:00:01"); err != nil {
		t.Fatal(err)
	}
	params := script.lastRequest()["params"].(map[string]any)
	if params["mac"] != "02:00:00:00:00:01" {
		t.Fatalf("显式 MAC 未生效: %v", params["mac"])
	}
}

func TestAuthValidation(t *testing.T) {
	script := &scriptedBridge{responses: map[string]string{"auth": okResponse(map[string]any{"authenticated": true})}}
	client := script.client()
	if _, err := client.Auth("alice", "", ""); err == nil {
		t.Fatal("密码缺失应报错")
	}
	if _, err := client.Auth("", "secret", ""); err == nil {
		t.Fatal("账号缺失应报错")
	}
	if _, err := client.Auth("", "", "02:00:00:00:00:01"); err == nil {
		t.Fatal("mac 无账号应报错")
	}
}

func TestAuthFromEnvironment(t *testing.T) {
	t.Setenv("THS_USERNAME", "env_user")
	t.Setenv("THS_PASSWORD", "env_pass")
	script := &scriptedBridge{responses: map[string]string{
		"auth": okResponse(map[string]any{"authenticated": true}),
	}}
	client := script.client()
	if _, err := client.Auth("", "", ""); err != nil {
		t.Fatal(err)
	}
	params := script.lastRequest()["params"].(map[string]any)
	if params["username"] != "env_user" || params["password"] != "env_pass" {
		t.Fatalf("环境变量凭据未生效: %v", params)
	}
}

func TestAuthNotCompleted(t *testing.T) {
	script := &scriptedBridge{responses: map[string]string{
		"auth": okResponse(map[string]any{"authenticated": false}),
	}}
	client := script.client()
	if _, err := client.Auth("alice", "secret", ""); err == nil {
		t.Fatal("authenticated=false 应报认证未完成")
	}
}

func TestRequestRoundTrip(t *testing.T) {
	script := &scriptedBridge{responses: map[string]string{
		"auth": okResponse(map[string]any{"authenticated": true}),
		"list_security_k_lines": okResponse([]any{
			map[string]any{"time": "20240102", "close": jsonNumber(10.5)},
		}),
	}}
	client := script.client()
	if _, err := client.Auth("alice", "secret", ""); err != nil {
		t.Fatal(err)
	}
	data, err := client.Request("list_security_k_lines", map[string]any{
		"security": "USHA000001",
		"start":    int64(-1),
		"end":      0,
		"adjust":   0,
		"period":   int64(0x4000),
	})
	if err != nil {
		t.Fatal(err)
	}
	rows := data.([]any)
	row := rows[0].(map[string]any)
	if row["close"] != 10.5 {
		t.Fatalf("close 应为 float64: %v", row["close"])
	}
	request := script.lastRequest()
	if request["method"] != "list_security_k_lines" {
		t.Fatalf("请求方法不符: %v", request["method"])
	}
	// Request 是原始入口（对应 Python _Client._request），params 不做字段规范化；
	// 规范化由各 API 包装层负责。
	params := request["params"].(map[string]any)
	if params["security"] != "USHA000001" || params["start"] != int64(-1) {
		t.Fatalf("请求参数不符: %v", params)
	}
}

func TestRequestRateLimitedRetry(t *testing.T) {
	script := &scriptedBridge{responses: map[string]string{
		"auth": okResponse(map[string]any{"authenticated": true}),
		"ping": `{"ok":false,"error":{"code":"rate_limited","message":"slow down","retry_after_ms":50}}`,
	}}
	client := script.client()
	if _, err := client.Auth("alice", "secret", ""); err != nil {
		t.Fatal(err)
	}
	original := client.runtime.callAlloc
	calls := 0
	// 第一次调用返回限流，随后切换为成功响应，验证 invoke 的重试路径。
	client.runtime.callAlloc = func(request string) unsafePointer {
		calls++
		result := original(request)
		if calls == 1 {
			script.responses["ping"] = `{"ok":true,"data":{"pong":true}}`
		}
		return result
	}
	data, err := client.Request("ping", nil, WithTimeout(5))
	if err != nil {
		t.Fatalf("限流后应重试成功: %v", err)
	}
	if data.(map[string]any)["pong"] != true {
		t.Fatalf("响应不符: %v", data)
	}
	if calls != 2 {
		t.Fatalf("应恰好调用 2 次: %d", calls)
	}
}

func TestRequestNonRetryableError(t *testing.T) {
	script := &scriptedBridge{responses: map[string]string{
		"auth": okResponse(map[string]any{"authenticated": true}),
		"ping": `{"ok":false,"error":{"code":"bad_request","message":"bad"}}`,
	}}
	client := script.client()
	if _, err := client.Auth("alice", "secret", ""); err != nil {
		t.Fatal(err)
	}
	_, err := client.Request("ping", nil)
	if err == nil {
		t.Fatal("应返回错误")
	}
	apiErr, ok := err.(*APIError)
	if !ok || apiErr.Code != "bad_request" {
		t.Fatalf("错误类型不符: %#v", err)
	}
}

func TestLogoutDeletesSession(t *testing.T) {
	script := &scriptedBridge{responses: map[string]string{
		"logout": okResponse(map[string]any{}),
	}}
	client := script.client()
	if err := client.Logout(); err != nil {
		t.Fatal(err)
	}
	if script.lastRequest()["method"] != "logout" {
		t.Fatal("logout 应调用原生 logout")
	}
}
