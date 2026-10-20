// Package thsdk 是同花顺 THS 行情/交易 SDK 的 Go 实现。
//
// 包提供三层调用能力：
//
//   - 认证：Auth（账号密码/环境变量/临时账户）、AuthQRCode（二维码登录）、
//     Logout（退出）；
//   - 80 个 Canonical 方法（以 <GoName> 命名，返回原生 envelope 结构，
//     与底层协议方法一一对应）；
//   - 63 个只读方法的表格包装（<GoName>Table，返回 *Table，
//     提供列/行/索引/元数据结构的表格结果）；
//   - 便捷函数：GetPrice / Klines / IntradayData / TickLevel1 / Depth /
//     CorporateAction / SearchSymbols / CompleteTHSCode / WencaiNLP / News /
//     MarketSecurities / BlockConstituents。
//
// 原生运行组件（hq.dylib / hq.so / hq.dll）默认通过 cgo 专用线程桥加载
// （推荐路径，运行最可靠），支持 macOS / Linux / Windows 三平台。
package thsdk

import "sync"

// Version 是 SDK 的对外版本号。
const Version = "3.0.0"

// 进程级默认客户端的全局状态：
//   - defaultClient：惰性创建的默认客户端实例；
//   - defaultClientLock：保护 defaultClient 的互斥锁，保证并发安全。
var (
	defaultClient     *Client
	defaultClientLock sync.Mutex
)

// getClient 返回进程级默认客户端，首次调用时惰性创建。
//
// 入参：无。
//
// 出参：*Client 默认客户端实例。会话目录为首次调用时的进程工作目录。
//
// 说明：所有公开便捷函数（Auth/GetPrice 等）都经由它转发请求，
// 保证进程内单例语义，避免多实例各自维护认证状态。
func getClient() *Client {
	defaultClientLock.Lock()
	defer defaultClientLock.Unlock()
	if defaultClient == nil {
		defaultClient = NewClient("", nil)
	}
	return defaultClient
}

// SetDefaultClient 替换进程级默认客户端。
//
// 入参：
//   - client：新的客户端实例；传 nil 会恢复惰性创建行为（下次调用时重建）。
//
// 出参：无。
//
// 说明：主要用于测试隔离（替换为 stub 客户端）或应用层注入自定义配置；
// 正常使用无需调用。
func SetDefaultClient(client *Client) {
	defaultClientLock.Lock()
	defer defaultClientLock.Unlock()
	defaultClient = client
}

// Close 关闭默认客户端的运行组件连接并清除认证状态。
//
// 入参：无。
//
// 出参：
//   - error：关闭失败时返回错误，成功返回 nil。
//
// 说明：可在程序结束前 defer 调用，用于释放原生运行组件资源。
func Close() error {
	defaultClientLock.Lock()
	client := defaultClient
	defaultClientLock.Unlock()
	if client == nil {
		return nil
	}
	return client.Close()
}
