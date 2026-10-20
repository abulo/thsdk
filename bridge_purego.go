//go:build !cgo && (darwin || linux)

package thsdk

// 本文件是 CGO_ENABLED=0 构建在 darwin/linux 下的原生桥实现：使用 purego
// （纯 Go FFI，无 cgo）在当前线程直接调用 CallAlloc / FreeResult。
//
// 注意：对 Go 语言构建的 c-shared 运行组件（本发行包即如此），在 Go 进程中
// 直接 dlopen/调用会与宿主 Go 运行时发生 TLS 冲突而崩溃，因此运行时请使用
// 默认的 cgo 构建（bridge_cgo.go）。本实现保留用于无 C 工具链环境下的
// 编译验证与普通非 Go 运行组件的调用。

import (
	"fmt"
	"unsafe"

	"github.com/ebitengine/purego"
)

// initNative 加载原生运行组件并注册 CallAlloc / FreeResult，对应 Python
// 构造器中 CallAlloc/FreeResult 符号声明部分。
func (b *RuntimeBridge) initNative() error {
	handle, err := purego.Dlopen(b.path, purego.RTLD_NOW)
	if err != nil {
		return newTHSDKError("THSDK 运行组件加载失败，请确认安装包与当前环境匹配")
	}
	// 返回值使用 unsafe.Pointer（等价于 ctypes 的 c_void_p），便于按 C 字符串读取。
	b.callAlloc = func(string) unsafe.Pointer { return nil }
	b.freeResult = func(unsafe.Pointer) {}
	func() {
		defer func() {
			// RegisterLibFunc 找不到符号时 panic；转换为统一错误。
			if recovered := recover(); recovered != nil {
				err = fmt.Errorf("%v", recovered)
			}
		}()
		purego.RegisterLibFunc(&b.callAlloc, handle, "CallAlloc")
		purego.RegisterLibFunc(&b.freeResult, handle, "FreeResult")
	}()
	if err != nil {
		return newTHSDKError("THSDK 运行组件版本过旧，请重新安装")
	}
	return nil
}
