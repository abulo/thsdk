// Code generated from the Python THSDK 2.0.2 source (thsdk/_apis/account.py). DO NOT EDIT BY HAND.
// 本文件对应 Python 包的 thsdk/_apis/account.py：13 个 Canonical 方法的一一对应包装。
// 返回值为原生 envelope 结构（any）；表格形态请使用同目录 api_dataframe.go 中的
// <GoName>Table 包装（对应 Python 公开层的 pandas.DataFrame 结果）。

package thsdk

// AccountPermissions 调用原生 AccountPermissions；请求 None，返回 AccountPermissions。
// 对应 Python thsdk._apis.account.account_permissions。
func AccountPermissions(opts ...*CallOptions) (any, error) {
	return callNone("account_permissions", opts...)
}

// AddAccountWatchlistGroupSecurities 调用原生 AddAccountWatchlistGroupSecurities；请求 AccountWatchlistGroupSecuritiesRequest，返回 AccountWatchlistGroups。 会修改当前账号的远端数据。
// 对应 Python thsdk._apis.account.add_account_watchlist_group_securities。
func AddAccountWatchlistGroupSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("add_account_watchlist_group_securities", params,
		[]string{"group_id", "securities"},
		nil, opts...)
}

// AddAccountWatchlistSecurities 调用原生 AddAccountWatchlistSecurities；请求 AddAccountWatchlistSecuritiesRequest，返回 AccountWatchlistUpdate。 会修改当前账号的远端数据。
// 对应 Python thsdk._apis.account.add_account_watchlist_securities。
func AddAccountWatchlistSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("add_account_watchlist_securities", params,
		[]string{"securities"},
		nil, opts...)
}

// ClearAccountWatchlist 调用原生 ClearAccountWatchlist；请求 None，返回 AccountWatchlistUpdate。 会修改当前账号的远端数据。
// 对应 Python thsdk._apis.account.clear_account_watchlist。
func ClearAccountWatchlist(opts ...*CallOptions) (any, error) {
	return callNone("clear_account_watchlist", opts...)
}

// CreateAccountWatchlistGroup 调用原生 CreateAccountWatchlistGroup；请求 CreateAccountWatchlistGroupRequest，返回 AccountWatchlistGroups。 会修改当前账号的远端数据。
// 对应 Python thsdk._apis.account.create_account_watchlist_group。
func CreateAccountWatchlistGroup(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("create_account_watchlist_group", params,
		[]string{"name"},
		nil, opts...)
}

// DeleteAccountWatchlistGroup 调用原生 DeleteAccountWatchlistGroup；请求 DeleteAccountWatchlistGroupRequest，返回 AccountWatchlistGroups。 会修改当前账号的远端数据。
// 对应 Python thsdk._apis.account.delete_account_watchlist_group。
func DeleteAccountWatchlistGroup(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("delete_account_watchlist_group", params,
		[]string{"group_id"},
		nil, opts...)
}

// GetAccountWatchlist 调用原生 GetAccountWatchlist；请求 None，返回 AccountWatchlistSnapshot。
// 对应 Python thsdk._apis.account.get_account_watchlist。
func GetAccountWatchlist(opts ...*CallOptions) (any, error) {
	return callNone("get_account_watchlist", opts...)
}

// GetAccountWatchlistGroups 调用原生 GetAccountWatchlistGroups；请求 None，返回 AccountWatchlistGroups。
// 对应 Python thsdk._apis.account.get_account_watchlist_groups。
func GetAccountWatchlistGroups(opts ...*CallOptions) (any, error) {
	return callNone("get_account_watchlist_groups", opts...)
}

// RemoveAccountWatchlistGroupSecurities 调用原生 RemoveAccountWatchlistGroupSecurities；请求 AccountWatchlistGroupSecuritiesRequest，返回 AccountWatchlistGroups。 会修改当前账号的远端数据。
// 对应 Python thsdk._apis.account.remove_account_watchlist_group_securities。
func RemoveAccountWatchlistGroupSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("remove_account_watchlist_group_securities", params,
		[]string{"group_id", "securities"},
		nil, opts...)
}

// RemoveAccountWatchlistSecurities 调用原生 RemoveAccountWatchlistSecurities；请求 RemoveAccountWatchlistSecuritiesRequest，返回 AccountWatchlistUpdate。 会修改当前账号的远端数据。
// 对应 Python thsdk._apis.account.remove_account_watchlist_securities。
func RemoveAccountWatchlistSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("remove_account_watchlist_securities", params,
		[]string{"securities"},
		nil, opts...)
}

// RenameAccountWatchlistGroup 调用原生 RenameAccountWatchlistGroup；请求 RenameAccountWatchlistGroupRequest，返回 AccountWatchlistGroups。 会修改当前账号的远端数据。
// 对应 Python thsdk._apis.account.rename_account_watchlist_group。
func RenameAccountWatchlistGroup(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("rename_account_watchlist_group", params,
		[]string{"group_id", "name"},
		nil, opts...)
}

// ReplaceAccountWatchlistGroupSecurities 调用原生 ReplaceAccountWatchlistGroupSecurities；请求 AccountWatchlistGroupSecuritiesRequest，返回 AccountWatchlistGroups。 会修改当前账号的远端数据。
// 对应 Python thsdk._apis.account.replace_account_watchlist_group_securities。
func ReplaceAccountWatchlistGroupSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("replace_account_watchlist_group_securities", params,
		[]string{"group_id", "securities"},
		nil, opts...)
}

// ReplaceAccountWatchlistSecurities 调用原生 ReplaceAccountWatchlistSecurities；请求 ReplaceAccountWatchlistSecuritiesRequest，返回 AccountWatchlistUpdate。 会修改当前账号的远端数据。
// 对应 Python thsdk._apis.account.replace_account_watchlist_securities。
func ReplaceAccountWatchlistSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("replace_account_watchlist_securities", params,
		[]string{"version", "securities"},
		nil, opts...)
}
