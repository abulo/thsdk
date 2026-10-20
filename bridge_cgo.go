//go:build cgo

package thsdk

// 本文件是 cgo 构建下的原生桥实现：把 dlopen / CallAlloc / FreeResult 全部
// 放到一条专用原生线程（见 bridge.c）上执行。
//
// 为什么要这样做：hq.dylib / hq.so / hq.dll 是 Go 语言构建的 c-shared 库。
// 若在 Go 进程的普通线程上 dlopen（其构造函数会初始化第二个 Go 运行时）或
// 直接调用其导出函数，两个 Go 运行时会在线程本地存储（TLS）上互相覆盖，
// 导致宿主运行时堆损坏（bad sweepgen / SIGSEGV）。Python 侧通过 ctypes 调用
// 没有问题，是因为解释器线程没有 Go TLS。桥线程用 pthread/Windows 线程创建，
// 从未运行过 Go 代码，TLS 槽位空闲，等价于 Python 的调用环境。

/*
#cgo linux LDFLAGS: -lpthread -ldl
#include <stdlib.h>
#include "bridge.h"
*/
import "C"

import (
	"unsafe"
)

// initNative 启动专用原生线程并完成 dlopen/dlsym，对应 Python 构造器中
// CallAlloc/FreeResult 符号声明部分。
func (b *RuntimeBridge) initNative() error {
	cPath := C.CString(b.path)
	defer C.free(unsafe.Pointer(cPath))
	if rc := C.bridge_init(cPath); rc != 0 {
		return newTHSDKError("THSDK 运行组件加载失败，请确认安装包与当前环境匹配")
	}
	// callAlloc 把请求交给桥线程执行；返回的 C 字符串由 freeResult 释放。
	// 与 Python ctypes 的 c_char_p 参数、c_void_p 返回值语义一致。
	b.callAlloc = func(request string) unsafe.Pointer {
		cRequest := C.CString(request)
		defer C.free(unsafe.Pointer(cRequest))
		return unsafe.Pointer(C.bridge_call(cRequest))
	}
	b.freeResult = func(pointer unsafe.Pointer) {
		if pointer != nil {
			C.bridge_free_string((*C.char)(pointer))
		}
	}
	return nil
}
