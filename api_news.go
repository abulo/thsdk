// Code generated from the Python THSDK 2.0.2 source (thsdk/_apis/news.py). DO NOT EDIT BY HAND.
// 本文件对应 Python 包的 thsdk/_apis/news.py：14 个 Canonical 方法的一一对应包装。
// 返回值为原生 envelope 结构（any）；表格形态请使用同目录 api_dataframe.go 中的
// <GoName>Table 包装（对应 Python 公开层的 pandas.DataFrame 结果）。

package thsdk

// AnalyzeSecurityLimitUp 调用原生 AnalyzeSecurityLimitUp；请求 SecurityLimitUpAnalysisRequest，返回 SecurityLimitUpAnalysisResult。
// 对应 Python thsdk._apis.news.analyze_security_limit_up。
func AnalyzeSecurityLimitUp(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("analyze_security_limit_up", params,
		[]string{"security"},
		[]string{"date"}, opts...)
}

// GetSecurityPriceLimitEvents 调用原生 GetSecurityPriceLimitEvents；请求 PriceLimitEventsRequest，返回 PriceLimitEventsResult。
// 对应 Python thsdk._apis.news.get_security_price_limit_events。
func GetSecurityPriceLimitEvents(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("get_security_price_limit_events", params,
		[]string{"security"},
		nil, opts...)
}

// GetSecurityShortTermHighlights 调用原生 GetSecurityShortTermHighlights；请求 SecurityShortTermHighlightsRequest，返回 SecurityShortTermHighlightsResult。
// 对应 Python thsdk._apis.news.get_security_short_term_highlights。
func GetSecurityShortTermHighlights(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("get_security_short_term_highlights", params,
		[]string{"security"},
		nil, opts...)
}

// ListCommodityStockLinkageNews 调用原生 ListCommodityStockLinkageNews；请求 CommodityStockLinkageNewsRequest，返回 CommodityStockLinkageNewsResult。
// 对应 Python thsdk._apis.news.list_commodity_stock_linkage_news。
func ListCommodityStockLinkageNews(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_commodity_stock_linkage_news", params,
		[]string{"commodity_index"},
		nil, opts...)
}

// ListHotBlockCalendarEvents 调用原生 ListHotBlockCalendarEvents；请求 None，返回 HotBlockCalendarResult。
// 对应 Python thsdk._apis.news.list_hot_block_calendar_events。
func ListHotBlockCalendarEvents(opts ...*CallOptions) (any, error) {
	return callNone("list_hot_block_calendar_events", opts...)
}

// ListHotEventNews 调用原生 ListHotEventNews；请求 HotEventNewsRequest，返回 HotEventNewsResult。
// 对应 Python thsdk._apis.news.list_hot_event_news。
func ListHotEventNews(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_hot_event_news", params,
		[]string{"block"},
		nil, opts...)
}

// ListMarketCallAuctionUnusualEvents 调用原生 ListMarketCallAuctionUnusualEvents；请求 CallAuctionUnusualEventsRequest，返回 list[ShortTermEvent]。
// 对应 Python thsdk._apis.news.list_market_call_auction_unusual_events。
func ListMarketCallAuctionUnusualEvents(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_market_call_auction_unusual_events", params,
		[]string{"market"},
		nil, opts...)
}

// ListMarketChangeEvents 调用原生 ListMarketChangeEvents；请求 None，返回 list[MarketChangeEvent]。
// 对应 Python thsdk._apis.news.list_market_change_events。
func ListMarketChangeEvents(opts ...*CallOptions) (any, error) {
	return callNone("list_market_change_events", opts...)
}

// ListMarketShortTermEvents 调用原生 ListMarketShortTermEvents；请求 MarketShortTermEventsRequest，返回 list[ShortTermEvent]。
// 对应 Python thsdk._apis.news.list_market_short_term_events。
func ListMarketShortTermEvents(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_market_short_term_events", params,
		[]string{"market", "mode"},
		nil, opts...)
}

// ListNewsFlashItems 调用原生 ListNewsFlashItems；请求 NewsFlashItemsRequest，返回 NewsFlashItemsResult。
// 对应 Python thsdk._apis.news.list_news_flash_items。
func ListNewsFlashItems(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_news_flash_items", params,
		[]string{"traversal"},
		nil, opts...)
}

// ListNewsItems 调用原生 ListNewsItems；请求 NewsItemsRequest，返回 NewsItemsResult。
// 对应 Python thsdk._apis.news.list_news_items。
func ListNewsItems(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_news_items", params,
		[]string{"category", "traversal"},
		nil, opts...)
}

// ListSecurityNewsEvents 调用原生 ListSecurityNewsEvents；请求 SecurityNewsEventsRequest，返回 SecurityNewsItemsResult。
// 对应 Python thsdk._apis.news.list_security_news_events。
func ListSecurityNewsEvents(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_news_events", params,
		[]string{"security", "timeline", "window"},
		nil, opts...)
}

// ListSecurityNewsMarkers 调用原生 ListSecurityNewsMarkers；请求 SecurityNewsMarkersRequest，返回 SecurityNewsItemsResult。
// 对应 Python thsdk._apis.news.list_security_news_markers。
func ListSecurityNewsMarkers(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_news_markers", params,
		[]string{"security", "timeline"},
		nil, opts...)
}

// ListSecurityResearchReports 调用原生 ListSecurityResearchReports；请求 SecurityResearchReportsRequest，返回 SecurityResearchReportsResult。
// 对应 Python thsdk._apis.news.list_security_research_reports。
func ListSecurityResearchReports(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_research_reports", params,
		[]string{"security", "traversal"},
		nil, opts...)
}
