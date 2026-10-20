// Code generated from the Python THSDK 2.0.2 source (thsdk/_apis/market_data.py). DO NOT EDIT BY HAND.
// 本文件对应 Python 包的 thsdk/_apis/market_data.py：17 个 Canonical 方法的一一对应包装。
// 返回值为原生 envelope 结构（any）；表格形态请使用同目录 api_dataframe.go 中的
// <GoName>Table 包装（对应 Python 公开层的 pandas.DataFrame 结果）。

package thsdk

// CheckMarketTimelineReadiness 调用原生 CheckMarketTimelineReadiness；请求 MarketTimelineReadinessRequest，返回 MarketTimelineReadinessResult。
// 对应 Python thsdk._apis.market_data.check_market_timeline_readiness。
func CheckMarketTimelineReadiness(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("check_market_timeline_readiness", params,
		[]string{"market", "time_index_day"},
		nil, opts...)
}

// GetSecurityTradingTimeline 调用原生 GetSecurityTradingTimeline；请求 SecurityTradingTimelineRequest，返回 SecurityTradingTimeline。
// 对应 Python thsdk._apis.market_data.get_security_trading_timeline。
func GetSecurityTradingTimeline(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("get_security_trading_timeline", params,
		[]string{"security"},
		nil, opts...)
}

// ListIndexCallAuctionQuotes 调用原生 ListIndexCallAuctionQuotes；请求 IndexCallAuctionQuotesRequest，返回 list[IndexCallAuctionQuote]。
// 对应 Python thsdk._apis.market_data.list_index_call_auction_quotes。
func ListIndexCallAuctionQuotes(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_index_call_auction_quotes", params,
		[]string{"index", "phase"},
		nil, opts...)
}

// ListSecurityCallAuctionQuotes 调用原生 ListSecurityCallAuctionQuotes；请求 SecurityCallAuctionQuotesRequest，返回 list[CallAuctionQuote]。
// 对应 Python thsdk._apis.market_data.list_security_call_auction_quotes。
func ListSecurityCallAuctionQuotes(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_call_auction_quotes", params,
		[]string{"security", "phase"},
		[]string{"date"}, opts...)
}

// ListSecurityCorporateActions 调用原生 ListSecurityCorporateActions；请求 SecurityCorporateActionsRequest，返回 list[CorporateAction]。
// 对应 Python thsdk._apis.market_data.list_security_corporate_actions。
func ListSecurityCorporateActions(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_corporate_actions", params,
		[]string{"security", "start", "end", "adjust", "period"},
		nil, opts...)
}

// ListSecurityDailyCapitalFlows 调用原生 ListSecurityDailyCapitalFlows；请求 SecurityDailyCapitalFlowsRequest，返回 list[DailyCapitalFlow]。
// 对应 Python thsdk._apis.market_data.list_security_daily_capital_flows。
func ListSecurityDailyCapitalFlows(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_daily_capital_flows", params,
		[]string{"securities"},
		nil, opts...)
}

// ListSecurityDailyKLinesWithPreviousClose 调用原生 ListSecurityDailyKLinesWithPreviousClose；请求 DailyKLineRequest，返回 list[DailyKLineWithPreviousClose]。
// 对应 Python thsdk._apis.market_data.list_security_daily_k_lines_with_previous_close。
func ListSecurityDailyKLinesWithPreviousClose(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_daily_k_lines_with_previous_close", params,
		[]string{"security", "start", "end", "adjust"},
		nil, opts...)
}

// ListSecurityExtendedHoursIntradayBars 调用原生 ListSecurityExtendedHoursIntradayBars；请求 SecurityExtendedHoursIntradayBarsRequest，返回 list[IntradayBar]。
// 对应 Python thsdk._apis.market_data.list_security_extended_hours_intraday_bars。
func ListSecurityExtendedHoursIntradayBars(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_extended_hours_intraday_bars", params,
		[]string{"security", "session", "date"},
		nil, opts...)
}

// ListSecurityExtendedHoursTicks 调用原生 ListSecurityExtendedHoursTicks；请求 SecurityExtendedHoursTicksRequest，返回 list[Tick]。
// 对应 Python thsdk._apis.market_data.list_security_extended_hours_ticks。
func ListSecurityExtendedHoursTicks(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_extended_hours_ticks", params,
		[]string{"security", "session", "window"},
		nil, opts...)
}

// ListSecurityFinancialSnapshots 调用原生 ListSecurityFinancialSnapshots；请求 SecurityFinancialSnapshotsRequest，返回 list[SecurityFinancialSnapshot]。
// 对应 Python thsdk._apis.market_data.list_security_financial_snapshots。
func ListSecurityFinancialSnapshots(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_financial_snapshots", params,
		[]string{"securities", "fields"},
		nil, opts...)
}

// ListSecurityIntradayBars 调用原生 ListSecurityIntradayBars；请求 IntradayBarsRequest，返回 list[IntradayBar]。
// 对应 Python thsdk._apis.market_data.list_security_intraday_bars。
func ListSecurityIntradayBars(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_intraday_bars", params,
		[]string{"security", "date"},
		nil, opts...)
}

// ListSecurityKLines 调用原生 ListSecurityKLines；请求 KLineRequest，返回 list[KLine]。
// 对应 Python thsdk._apis.market_data.list_security_k_lines。
func ListSecurityKLines(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_k_lines", params,
		[]string{"security", "start", "end", "adjust", "period"},
		nil, opts...)
}

// ListSecurityOptionGreeks 调用原生 ListSecurityOptionGreeks；请求 OptionGreeksRequest，返回 list[OptionGreeksRow]。
// 对应 Python thsdk._apis.market_data.list_security_option_greeks。
func ListSecurityOptionGreeks(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_option_greeks", params,
		[]string{"securities"},
		nil, opts...)
}

// ListSecurityOrderBooks 调用原生 ListSecurityOrderBooks；请求 SecurityOrderBooksRequest，返回 list[SecurityOrderBook]。
// 对应 Python thsdk._apis.market_data.list_security_order_books。
func ListSecurityOrderBooks(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_order_books", params,
		[]string{"securities", "depth"},
		nil, opts...)
}

// ListSecurityPriceVolumeLevels 调用原生 ListSecurityPriceVolumeLevels；请求 SecurityPriceVolumeLevelsRequest，返回 list[PriceVolumeLevel]。
// 对应 Python thsdk._apis.market_data.list_security_price_volume_levels。
func ListSecurityPriceVolumeLevels(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_price_volume_levels", params,
		[]string{"security"},
		nil, opts...)
}

// ListSecurityTicks 调用原生 ListSecurityTicks；请求 SecurityTicksRequest，返回 list[Tick]。
// 对应 Python thsdk._apis.market_data.list_security_ticks。
func ListSecurityTicks(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_ticks", params,
		[]string{"security", "window"},
		nil, opts...)
}

// ListSecurityTimeAndSalesTicks 调用原生 ListSecurityTimeAndSalesTicks；请求 SecurityTimeAndSalesTicksRequest，返回 list[Tick]。
// 对应 Python thsdk._apis.market_data.list_security_time_and_sales_ticks。
func ListSecurityTimeAndSalesTicks(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_time_and_sales_ticks", params,
		[]string{"security", "count"},
		nil, opts...)
}
