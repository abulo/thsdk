// 本文件实现 news 模块的 14 个 Canonical 方法包装（请求结构见方法清单）。
// 返回值为原生 envelope 结构（any）；表格形态请使用同目录 api_table.go 中的
// <GoName>Table 包装返回表格结果。

package thsdk

// AnalyzeSecurityLimitUp 调用原生方法 AnalyzeSecurityLimitUp，请求类型 SecurityLimitUpAnalysisRequest，返回类型 SecurityLimitUpAnalysisResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func AnalyzeSecurityLimitUp(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("analyze_security_limit_up", params,
		[]string{"security"},
		[]string{"date"}, opts...)
}

// GetSecurityPriceLimitEvents 调用原生方法 GetSecurityPriceLimitEvents，请求类型 PriceLimitEventsRequest，返回类型 PriceLimitEventsResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func GetSecurityPriceLimitEvents(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("get_security_price_limit_events", params,
		[]string{"security"},
		nil, opts...)
}

// GetSecurityShortTermHighlights 调用原生方法 GetSecurityShortTermHighlights，请求类型 SecurityShortTermHighlightsRequest，返回类型 SecurityShortTermHighlightsResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func GetSecurityShortTermHighlights(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("get_security_short_term_highlights", params,
		[]string{"security"},
		nil, opts...)
}

// ListCommodityStockLinkageNews 调用原生方法 ListCommodityStockLinkageNews，请求类型 CommodityStockLinkageNewsRequest，返回类型 CommodityStockLinkageNewsResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func ListCommodityStockLinkageNews(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_commodity_stock_linkage_news", params,
		[]string{"commodity_index"},
		nil, opts...)
}

// ListHotBlockCalendarEvents 调用原生方法 ListHotBlockCalendarEvents（无请求参数），返回 HotBlockCalendarResult。
// 入参：opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func ListHotBlockCalendarEvents(opts ...*CallOptions) (any, error) {
	return callNone("list_hot_block_calendar_events", opts...)
}

// ListHotEventNews 调用原生方法 ListHotEventNews，请求类型 HotEventNewsRequest，返回类型 HotEventNewsResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func ListHotEventNews(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_hot_event_news", params,
		[]string{"block"},
		nil, opts...)
}

// ListMarketCallAuctionUnusualEvents 调用原生 ListMarketCallAuctionUnusualEvents；请求 CallAuctionUnusualEventsRequest，返回 list[ShortTermEvent]。
func ListMarketCallAuctionUnusualEvents(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_market_call_auction_unusual_events", params,
		[]string{"market"},
		nil, opts...)
}

// ListMarketChangeEvents 调用原生 ListMarketChangeEvents；请求 None，返回 list[MarketChangeEvent]。
func ListMarketChangeEvents(opts ...*CallOptions) (any, error) {
	return callNone("list_market_change_events", opts...)
}

// ListMarketShortTermEvents 调用原生 ListMarketShortTermEvents；请求 MarketShortTermEventsRequest，返回 list[ShortTermEvent]。
func ListMarketShortTermEvents(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_market_short_term_events", params,
		[]string{"market", "mode"},
		nil, opts...)
}

// ListNewsFlashItems 调用原生方法 ListNewsFlashItems，请求类型 NewsFlashItemsRequest，返回类型 NewsFlashItemsResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func ListNewsFlashItems(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_news_flash_items", params,
		[]string{"traversal"},
		nil, opts...)
}

// ListNewsItems 调用原生方法 ListNewsItems，请求类型 NewsItemsRequest，返回类型 NewsItemsResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func ListNewsItems(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_news_items", params,
		[]string{"category", "traversal"},
		nil, opts...)
}

// ListSecurityNewsEvents 调用原生方法 ListSecurityNewsEvents，请求类型 SecurityNewsEventsRequest，返回类型 SecurityNewsItemsResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func ListSecurityNewsEvents(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_news_events", params,
		[]string{"security", "timeline", "window"},
		nil, opts...)
}

// ListSecurityNewsMarkers 调用原生方法 ListSecurityNewsMarkers，请求类型 SecurityNewsMarkersRequest，返回类型 SecurityNewsItemsResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func ListSecurityNewsMarkers(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_news_markers", params,
		[]string{"security", "timeline"},
		nil, opts...)
}

// ListSecurityResearchReports 调用原生方法 ListSecurityResearchReports，请求类型 SecurityResearchReportsRequest，返回类型 SecurityResearchReportsResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func ListSecurityResearchReports(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_research_reports", params,
		[]string{"security", "traversal"},
		nil, opts...)
}
