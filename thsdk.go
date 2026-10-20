// Package thsdk 是同花顺 THSDK 2.0.2 的 Go 移植实现。
//
// 对应 Python 包 thsdk/__init__.py 与 thsdk/api.py 的公开面：
//
//   - 认证：Auth / AuthQRCode / Logout；
//   - 80 个 Canonical 方法（<GoName>，返回原生 envelope 结构）；
//   - 63 个只读方法的表格包装（<GoName>Table，返回 *Table，
//     pandas.DataFrame 的等价实现）；
//   - 便捷函数：GetPrice / Klines / IntradayData / TickLevel1 / Depth /
//     CorporateAction / SearchSymbols / CompleteTHSCode / WencaiNLP / News /
//     MarketSecurities / BlockConstituents。
//
// 原生运行组件（hq.dylib / hq.so / hq.dll）通过 purego（纯 Go FFI）加载，
// 支持 macOS / Linux / Windows，无需 cgo。
package thsdk

import "sync"

// Version 对应 Python __version__ = "3.0.0"。
const Version = "3.0.0"

var (
	defaultClient     *Client
	defaultClientLock sync.Mutex
)

// getClient 返回进程级默认客户端（惰性创建），对应 Python _get_client。
// 会话目录为首次调用时的进程工作目录。
func getClient() *Client {
	defaultClientLock.Lock()
	defer defaultClientLock.Unlock()
	if defaultClient == nil {
		defaultClient = NewClient("", nil)
	}
	return defaultClient
}

// SetDefaultClient 替换进程级默认客户端（主要用于测试隔离）；
// 传 nil 会恢复惰性创建行为。对应 Python _set_client_for_tests。
func SetDefaultClient(client *Client) {
	defaultClientLock.Lock()
	defer defaultClientLock.Unlock()
	defaultClient = client
}

// Close 关闭默认客户端的运行组件连接并清除认证状态，
// 对应 Python 进程退出时的 _close_at_exit。可在程序结束前 defer 调用。
func Close() error {
	defaultClientLock.Lock()
	client := defaultClient
	defaultClientLock.Unlock()
	if client == nil {
		return nil
	}
	return client.Close()
}
