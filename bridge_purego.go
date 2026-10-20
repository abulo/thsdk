//go:build !cgo && (darwin || linux)

package thsdk

// 本文件是 CGO_ENABLED=0 构建在 darwin/linux 下的原生桥实现：
// 使用 purego（纯 Go FFI，无 cgo）在当前线程直接调用 CallAlloc / FreeResult。
//
// 注意：对 Go 语言构建的 c-shared 运行组件（本发行包即如此），在 Go 进程
// 中直接 dlopen/调用会与宿主 Go 运行时发生 TLS 冲突而崩溃，因此运行时请
// 使用默认的 cgo 构建（bridge_cgo.go）。本实现保留用于无 C 工具链环境下
// 的编译验证与普通非 Go 运行组件的调用。

import (
	"fmt"
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
// 说明：返回值使用 unsafe.Pointer（等价于 C 的 void*），便于按 C 字符串
// 读取；RegisterLibFunc 找不到符号时会 panic，这里统一转换为错误返回。
func (b *RuntimeBridge) initNative() error {
	// 打开动态库句柄（立即解析全部符号）。
	handle, err := purego.Dlopen(b.path, purego.RTLD_NOW)
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
		// 注册两个导出符号到闭包变量。
		purego.RegisterLibFunc(&b.callAlloc, handle, "CallAlloc")
		purego.RegisterLibFunc(&b.freeResult, handle, "FreeResult")
	}()
	if err != nil {
		return newTHSDKError("THSDK 运行组件版本过旧，请重新安装")
	}
	return nil
}
