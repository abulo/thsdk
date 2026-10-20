// Code generated from the Python THSDK 2.0.2 source (thsdk/_apis/wencai.py). DO NOT EDIT BY HAND.
// 本文件对应 Python 包的 thsdk/_apis/wencai.py：7 个 Canonical 方法的一一对应包装。
// 返回值为原生 envelope 结构（any）；表格形态请使用同目录 api_dataframe.go 中的
// <GoName>Table 包装（对应 Python 公开层的 pandas.DataFrame 结果）。

package thsdk

// ListWencaiExpressionFields 调用原生 ListWencaiExpressionFields；请求 query string，返回 list[WencaiField]。
// 对应 Python thsdk._apis.wencai.list_wencai_expression_fields。
func ListWencaiExpressionFields(query string, opts ...*CallOptions) (any, error) {
	return callString("list_wencai_expression_fields", query, "query", opts...)
}

// ListWencaiHotBlocks 调用原生 ListWencaiHotBlocks；请求 WencaiHotBlocksRequest，返回 WencaiHotBlocksResult。
// 对应 Python thsdk._apis.wencai.list_wencai_hot_blocks。
func ListWencaiHotBlocks(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_wencai_hot_blocks", params,
		[]string{"kind"},
		nil, opts...)
}

// QueryWencai 调用原生 QueryWencai；请求 WencaiQueryRequest，返回 WencaiQueryResult。
// 对应 Python thsdk._apis.wencai.query_wencai。
func QueryWencai(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("query_wencai", params,
		[]string{"query"},
		nil, opts...)
}

// QueryWencaiRealtimeFields 调用原生 QueryWencaiRealtimeFields；请求 WencaiRealtimeFieldsRequest，返回 WencaiRealtimeFieldsResult。
// 对应 Python thsdk._apis.wencai.query_wencai_realtime_fields。
func QueryWencaiRealtimeFields(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("query_wencai_realtime_fields", params,
		[]string{"securities", "fields"},
		nil, opts...)
}

// QueryWencaiSecurities 调用原生 QueryWencaiSecurities；请求 WencaiSecuritiesRequest，返回 WencaiSecuritiesResult。
// 对应 Python thsdk._apis.wencai.query_wencai_securities。
func QueryWencaiSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("query_wencai_securities", params,
		[]string{"query"},
		nil, opts...)
}

// RankBlockSecuritiesByWencaiField 调用原生 RankBlockSecuritiesByWencaiField；请求 RankBlockSecuritiesByWencaiFieldRequest，返回 WencaiBlockRankingResult。
// 对应 Python thsdk._apis.wencai.rank_block_securities_by_wencai_field。
func RankBlockSecuritiesByWencaiField(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("rank_block_securities_by_wencai_field", params,
		[]string{"block", "field"},
		nil, opts...)
}

// ResolveWencaiSecurities 调用原生 ResolveWencaiSecurities；请求 ResolveWencaiSecuritiesRequest，返回 ResolvedWencaiSecuritiesResult。
// 对应 Python thsdk._apis.wencai.resolve_wencai_securities。
func ResolveWencaiSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("resolve_wencai_securities", params,
		[]string{"query"},
		nil, opts...)
}
