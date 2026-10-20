package thsdk

// 本文件对应 Python 包的 thsdk/exceptions.py。

import "fmt"

// THSDKError 是高阶 THSDK 门面的基础错误类型，对应 Python 的 THSDKError(RuntimeError)。
// 运行时组件加载失败、请求校验失败、返回数据不兼容等场景都会返回该类型。
type THSDKError struct {
	Message string
}

func (e *THSDKError) Error() string { return e.Message }

func newTHSDKError(format string, args ...any) *THSDKError {
	return &THSDKError{Message: fmt.Sprintf(format, args...)}
}

// AuthenticationError 表示认证或持久化会话恢复失败。
// 对应 Python 的 AuthenticationError(THSDKError)。
type AuthenticationError struct {
	Message string
}

func (e *AuthenticationError) Error() string { return e.Message }

// NotAuthenticatedError 表示没有可用的已认证会话。
// 对应 Python 的 NotAuthenticatedError(AuthenticationError)。
type NotAuthenticatedError struct {
	Message string
}

func (e *NotAuthenticatedError) Error() string { return e.Message }

// APIError 表示一次类型化的公开 API 调用失败，由原生运行组件返回错误时产生。
// 对应 Python 的 APIError(THSDKError)。
//
// RetryAfterMs 仅在原生组件返回有效的整数 retry_after_ms 时非空。
type APIError struct {
	Code         string
	Message      string
	RetryAfterMs *int
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return e.Code + ": " + e.Message
	}
	return e.Message
}
