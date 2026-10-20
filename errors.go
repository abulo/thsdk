package thsdk

// 本文件定义整个 SDK 对外暴露的错误类型体系。
//
// 错误类型分层：
//   - THSDKError：所有错误的基类，表示 SDK 运行层面的错误（组件加载失败、
//     请求校验失败、返回数据不兼容等）；
//   - AuthenticationError：认证相关错误（凭据无效、会话恢复失败等）；
//   - NotAuthenticatedError：尚未认证就发起数据请求时的错误；
//   - APIError：原生运行组件返回业务错误时的类型化错误。
//
// 所有错误都实现 error 接口，可直接作为普通 error 返回给调用方。

import "fmt"

// THSDKError 是 SDK 最基础的错误类型，所有高层错误都从它派生。
//
// 典型场景：运行组件加载失败、请求参数校验失败、原生返回数据不兼容、
// 请求超限等。对应调用方可通过 errors.As 判断错误类型。
type THSDKError struct {
	// Message 保存人类可读的错误描述。
	Message string
}

// Error 实现 error 接口，返回错误描述文本。
func (e *THSDKError) Error() string { return e.Message }

// newTHSDKError 按格式化字符串构造一个 THSDKError。
//
// 入参：
//   - format：格式化模板，与 fmt.Sprintf 一致；
//   - args：模板对应的参数列表。
//
// 出参：返回构造好的错误指针，供各层直接返回。
func newTHSDKError(format string, args ...any) *THSDKError {
	return &THSDKError{Message: fmt.Sprintf(format, args...)}
}

// AuthenticationError 表示认证失败或持久化会话恢复失败。
//
// 典型场景：账号密码错误、环境变量凭据缺失、二维码登录超时、
// 临时账户认证失败、会话文件不可用等。
type AuthenticationError struct {
	// Message 保存人类可读的错误描述。
	Message string
}

// Error 实现 error 接口，返回错误描述文本。
func (e *AuthenticationError) Error() string { return e.Message }

// NotAuthenticatedError 表示发起数据请求时没有可用的已认证会话。
//
// 典型场景：未调用 Auth/AuthQRCode 就直接调用 GetPrice、List* 等数据接口。
// 调用方应先完成认证再发请求。
type NotAuthenticatedError struct {
	// Message 保存人类可读的错误描述。
	Message string
}

// Error 实现 error 接口，返回错误描述文本。
func (e *NotAuthenticatedError) Error() string { return e.Message }

// APIError 表示一次类型化的公开 API 调用失败，由原生运行组件返回错误时产生。
//
// 区别于 THSDKError：APIError 携带原生侧的结构化错误码（Code），
// 调用方可根据 Code 做精细化处理（如限流重试、字段重查等）。
type APIError struct {
	// Code 是原生运行组件返回的稳定错误码（如 rate_limited、
	// invalid_request、session_unavailable 等）。
	Code string
	// Message 是原生侧返回的错误描述文本。
	Message string
	// RetryAfterMs 仅在原生组件返回有效的整数 retry_after_ms 时非空，
	// 表示建议等待的毫秒数后再重试。
	RetryAfterMs *int
}

// Error 实现 error 接口。
//
// 返回格式：优先返回 "Code: Message"，Code 为空时仅返回 Message。
// 便于日志中快速定位错误类别。
func (e *APIError) Error() string {
	if e.Code != "" {
		return e.Code + ": " + e.Message
	}
	return e.Message
}
