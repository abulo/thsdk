//go:build !cgo && windows

package thsdk

// 本文件是 CGO_ENABLED=0 构建在 windows 下的原生桥实现：
// 用标准库 syscall 加载 hq.dll（purego 在 windows 上不提供 Dlopen），
// 再交给 purego 注册 CallAlloc / FreeResult 调用。
//
// 注意：对 Go 语言构建的 c-shared 运行组件（本发行包即如此），在 Go 进程
// 中直接调用会与宿主 Go 运行时发生 TLS 冲突而崩溃，因此运行时请使用默认的
// cgo 构建（bridge_cgo.go）。本实现保留用于无 C 工具链环境下的编译验证。

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/ebitengine/purego"
)

// initNative 加载原生运行组件并注册 CallAlloc / FreeResult。
//
// 入参：无（组件路径来自 RuntimeBridge.path）。
//
// 出参：
//   - error：组件加载失败或符号缺失时返回错误。
//
// 说明：Windows 下使用 LoadLibrary 打开 DLL 句柄，随后由 RegisterLibFunc
// 覆盖两个闭包为原生函数；找不到符号时的 panic 统一转换为错误返回。
func (b *RuntimeBridge) initNative() error {
	// 用系统 API 加载 DLL。
	library, err := syscall.LoadLibrary(b.path)
	if err != nil {
		return newTHSDKError("THSDK 运行组件加载失败，请确认安装包与当前环境匹配")
	}
	// 先以空实现占位，随后由 RegisterLibFunc 覆盖为原生函数。
	b.callAlloc = func(string) unsafe.Pointer { return nil }
	b.freeResult = func(unsafe.Pointer) {}
	func() {
		defer func() {
			// RegisterLibFunc 找不到符号时 panic；转换为统一错误。
			if recovered := recover(); recovered != nil {
				err = fmt.Errorf("%v", recovered)
			}
		}()
		// 注册两个导出符号到闭包变量（使用 DLL 句柄）。
		purego.RegisterLibFunc(&b.callAlloc, uintptr(library), "CallAlloc")
		purego.RegisterLibFunc(&b.freeResult, uintptr(library), "FreeResult")
	}()
	if err != nil {
		return newTHSDKError("THSDK 运行组件版本过旧，请重新安装")
	}
	return nil
}
