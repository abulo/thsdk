//go:build cgo

package thsdk

// 本文件是 cgo 构建下的原生桥实现：把 dlopen / CallAlloc / FreeResult
// 全部放到一条专用原生线程（见 bridge.c）上执行。
//
// 背景：原生运行组件（hq.dylib / hq.so / hq.dll）是 Go 语言构建的
// c-shared 库。若在 Go 进程的普通线程上 dlopen（其构造函数会初始化第二个
// Go 运行时）或直接调用其导出函数，两个 Go 运行时会在线程本地存储（TLS）
// 上互相覆盖，导致宿主运行时堆损坏（bad sweepgen / SIGSEGV）。
// 桥线程用 pthread/Windows 线程创建，从未运行过 Go 代码，TLS 槽位空闲，
// 等价于独立的原生调用环境，因此可以安全承载全部原生调用。

/*
#cgo linux LDFLAGS: -lpthread -ldl
#include <stdlib.h>
#include "bridge.h"
*/
import "C"

import (
	"unsafe"
)

// initNative 启动专用原生线程并完成 dlopen/dlsym。
//
// 入参：无（组件路径来自 RuntimeBridge.path）。
//
// 出参：
//   - error：原生组件加载或符号解析失败时返回错误。
//
// 说明：加载成功后注册 callAlloc / freeResult 两个闭包，
// 后续 Call 流程只经由这两个闭包与原生侧交换 C 字符串。
func (b *RuntimeBridge) initNative() error {
	// 把 Go 路径字符串转为 C 字符串传给桥（用完即释放）。
	cPath := C.CString(b.path)
	defer C.free(unsafe.Pointer(cPath))
	// 启动桥线程并加载组件，非 0 表示失败。
	if rc := C.bridge_init(cPath); rc != 0 {
		return newTHSDKError("THSDK 运行组件加载失败，请确认安装包与当前环境匹配")
	}
	// callAlloc：把请求交给桥线程执行；返回的 C 字符串由 freeResult 释放。
	// 请求字符串转为 C 字符串传入（桥线程内 strdup 复制），返回响应副本。
	b.callAlloc = func(request string) unsafe.Pointer {
		cRequest := C.CString(request)
		defer C.free(unsafe.Pointer(cRequest))
		return unsafe.Pointer(C.bridge_call(cRequest))
	}
	// freeResult：释放桥线程返回的响应字符串副本。
	b.freeResult = func(pointer unsafe.Pointer) {
		if pointer != nil {
			C.bridge_free_string((*C.char)(pointer))
		}
	}
	return nil
}
