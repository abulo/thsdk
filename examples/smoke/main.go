// 命令 smoke 验证本机原生运行组件（hq.dylib / hq.so / hq.dll）能否被
// purego 加载并通过 ABI 握手（版本、严格探测、方法清单），无需外部网络。
package main

import (
	"fmt"
	"log"
	"runtime"

	thsdk "thsdk"
)

func main() {
	bridge, err := thsdk.NewRuntimeBridge("")
	if err != nil {
		log.Fatalf("运行组件加载/握手失败: %v", err)
	}
	methods := thsdk.MethodInventory()
	fmt.Printf("THSDK Go 运行组件握手成功\n")
	fmt.Printf("  平台: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("  组件: %s\n", bridge.Path())
	fmt.Printf("  方法清单: %d 项\n", len(methods))
}
