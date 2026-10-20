// 本文件实现 account 模块的 13 个 Canonical 方法包装（请求结构见方法清单）。
// 返回值为原生 envelope 结构（any）；表格形态请使用同目录 api_table.go 中的
// <GoName>Table 包装返回表格结果。

package thsdk

// AccountPermissions 调用原生方法 AccountPermissions（无请求参数），返回 AccountPermissions。
// 入参：opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func AccountPermissions(opts ...*CallOptions) (any, error) {
	return callNone("account_permissions", opts...)
}

// AddAccountWatchlistGroupSecurities 调用原生方法 AddAccountWatchlistGroupSecurities，请求类型 AccountWatchlistGroupSecuritiesRequest，返回类型 AccountWatchlistGroups。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func AddAccountWatchlistGroupSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("add_account_watchlist_group_securities", params,
		[]string{"group_id", "securities"},
		nil, opts...)
}

// AddAccountWatchlistSecurities 调用原生方法 AddAccountWatchlistSecurities，请求类型 AddAccountWatchlistSecuritiesRequest，返回类型 AccountWatchlistUpdate。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func AddAccountWatchlistSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("add_account_watchlist_securities", params,
		[]string{"securities"},
		nil, opts...)
}

// ClearAccountWatchlist 调用原生方法 ClearAccountWatchlist（无请求参数），返回 AccountWatchlistUpdate。
// 入参：opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func ClearAccountWatchlist(opts ...*CallOptions) (any, error) {
	return callNone("clear_account_watchlist", opts...)
}

// CreateAccountWatchlistGroup 调用原生方法 CreateAccountWatchlistGroup，请求类型 CreateAccountWatchlistGroupRequest，返回类型 AccountWatchlistGroups。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func CreateAccountWatchlistGroup(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("create_account_watchlist_group", params,
		[]string{"name"},
		nil, opts...)
}

// DeleteAccountWatchlistGroup 调用原生方法 DeleteAccountWatchlistGroup，请求类型 DeleteAccountWatchlistGroupRequest，返回类型 AccountWatchlistGroups。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func DeleteAccountWatchlistGroup(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("delete_account_watchlist_group", params,
		[]string{"group_id"},
		nil, opts...)
}

// GetAccountWatchlist 调用原生方法 GetAccountWatchlist（无请求参数），返回 AccountWatchlistSnapshot。
// 入参：opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func GetAccountWatchlist(opts ...*CallOptions) (any, error) {
	return callNone("get_account_watchlist", opts...)
}

// GetAccountWatchlistGroups 调用原生方法 GetAccountWatchlistGroups（无请求参数），返回 AccountWatchlistGroups。
// 入参：opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func GetAccountWatchlistGroups(opts ...*CallOptions) (any, error) {
	return callNone("get_account_watchlist_groups", opts...)
}

// RemoveAccountWatchlistGroupSecurities 调用原生方法 RemoveAccountWatchlistGroupSecurities，请求类型 AccountWatchlistGroupSecuritiesRequest，返回类型 AccountWatchlistGroups。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func RemoveAccountWatchlistGroupSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("remove_account_watchlist_group_securities", params,
		[]string{"group_id", "securities"},
		nil, opts...)
}

// RemoveAccountWatchlistSecurities 调用原生方法 RemoveAccountWatchlistSecurities，请求类型 RemoveAccountWatchlistSecuritiesRequest，返回类型 AccountWatchlistUpdate。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func RemoveAccountWatchlistSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("remove_account_watchlist_securities", params,
		[]string{"securities"},
		nil, opts...)
}

// RenameAccountWatchlistGroup 调用原生方法 RenameAccountWatchlistGroup，请求类型 RenameAccountWatchlistGroupRequest，返回类型 AccountWatchlistGroups。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func RenameAccountWatchlistGroup(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("rename_account_watchlist_group", params,
		[]string{"group_id", "name"},
		nil, opts...)
}

// ReplaceAccountWatchlistGroupSecurities 调用原生方法 ReplaceAccountWatchlistGroupSecurities，请求类型 AccountWatchlistGroupSecuritiesRequest，返回类型 AccountWatchlistGroups。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func ReplaceAccountWatchlistGroupSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("replace_account_watchlist_group_securities", params,
		[]string{"group_id", "securities"},
		nil, opts...)
}

// ReplaceAccountWatchlistSecurities 调用原生方法 ReplaceAccountWatchlistSecurities，请求类型 ReplaceAccountWatchlistSecuritiesRequest，返回类型 AccountWatchlistUpdate。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func ReplaceAccountWatchlistSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("replace_account_watchlist_securities", params,
		[]string{"version", "securities"},
		nil, opts...)
}
