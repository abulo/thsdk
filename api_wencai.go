// 本文件实现 wencai 模块的 7 个 Canonical 方法包装（请求结构见方法清单）。
// 返回值为原生 envelope 结构（any）；表格形态请使用同目录 api_table.go 中的
// <GoName>Table 包装返回表格结果。

package thsdk

// ListWencaiExpressionFields 调用原生 ListWencaiExpressionFields；请求 query string，返回 list[WencaiField]。
func ListWencaiExpressionFields(query string, opts ...*CallOptions) (any, error) {
	return callString("list_wencai_expression_fields", query, "query", opts...)
}

// ListWencaiHotBlocks 调用原生方法 ListWencaiHotBlocks，请求类型 WencaiHotBlocksRequest，返回类型 WencaiHotBlocksResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func ListWencaiHotBlocks(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_wencai_hot_blocks", params,
		[]string{"kind"},
		nil, opts...)
}

// QueryWencai 调用原生方法 QueryWencai，请求类型 WencaiQueryRequest，返回类型 WencaiQueryResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func QueryWencai(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("query_wencai", params,
		[]string{"query"},
		nil, opts...)
}

// QueryWencaiRealtimeFields 调用原生方法 QueryWencaiRealtimeFields，请求类型 WencaiRealtimeFieldsRequest，返回类型 WencaiRealtimeFieldsResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func QueryWencaiRealtimeFields(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("query_wencai_realtime_fields", params,
		[]string{"securities", "fields"},
		nil, opts...)
}

// QueryWencaiSecurities 调用原生方法 QueryWencaiSecurities，请求类型 WencaiSecuritiesRequest，返回类型 WencaiSecuritiesResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func QueryWencaiSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("query_wencai_securities", params,
		[]string{"query"},
		nil, opts...)
}

// RankBlockSecuritiesByWencaiField 调用原生方法 RankBlockSecuritiesByWencaiField，请求类型 RankBlockSecuritiesByWencaiFieldRequest，返回类型 WencaiBlockRankingResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func RankBlockSecuritiesByWencaiField(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("rank_block_securities_by_wencai_field", params,
		[]string{"block", "field"},
		nil, opts...)
}

// ResolveWencaiSecurities 调用原生方法 ResolveWencaiSecurities，请求类型 ResolveWencaiSecuritiesRequest，返回类型 ResolvedWencaiSecuritiesResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func ResolveWencaiSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("resolve_wencai_securities", params,
		[]string{"query"},
		nil, opts...)
}
