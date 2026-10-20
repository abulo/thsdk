// 本文件实现 market_data 模块的 17 个 Canonical 方法包装（请求结构见方法清单）。
// 返回值为原生 envelope 结构（any）；表格形态请使用同目录 api_table.go 中的
// <GoName>Table 包装返回表格结果。

package thsdk

// CheckMarketTimelineReadiness 调用原生方法 CheckMarketTimelineReadiness，请求类型 MarketTimelineReadinessRequest，返回类型 MarketTimelineReadinessResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func CheckMarketTimelineReadiness(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("check_market_timeline_readiness", params,
		[]string{"market", "time_index_day"},
		nil, opts...)
}

// GetSecurityTradingTimeline 调用原生方法 GetSecurityTradingTimeline，请求类型 SecurityTradingTimelineRequest，返回类型 SecurityTradingTimeline。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func GetSecurityTradingTimeline(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("get_security_trading_timeline", params,
		[]string{"security"},
		nil, opts...)
}

// ListIndexCallAuctionQuotes 调用原生 ListIndexCallAuctionQuotes；请求 IndexCallAuctionQuotesRequest，返回 list[IndexCallAuctionQuote]。
func ListIndexCallAuctionQuotes(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_index_call_auction_quotes", params,
		[]string{"index", "phase"},
		nil, opts...)
}

// ListSecurityCallAuctionQuotes 调用原生 ListSecurityCallAuctionQuotes；请求 SecurityCallAuctionQuotesRequest，返回 list[CallAuctionQuote]。
func ListSecurityCallAuctionQuotes(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_call_auction_quotes", params,
		[]string{"security", "phase"},
		[]string{"date"}, opts...)
}

// ListSecurityCorporateActions 调用原生 ListSecurityCorporateActions；请求 SecurityCorporateActionsRequest，返回 list[CorporateAction]。
func ListSecurityCorporateActions(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_corporate_actions", params,
		[]string{"security", "start", "end", "adjust", "period"},
		nil, opts...)
}

// ListSecurityDailyCapitalFlows 调用原生 ListSecurityDailyCapitalFlows；请求 SecurityDailyCapitalFlowsRequest，返回 list[DailyCapitalFlow]。
func ListSecurityDailyCapitalFlows(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_daily_capital_flows", params,
		[]string{"securities"},
		nil, opts...)
}

// ListSecurityDailyKLinesWithPreviousClose 调用原生 ListSecurityDailyKLinesWithPreviousClose；请求 DailyKLineRequest，返回 list[DailyKLineWithPreviousClose]。
func ListSecurityDailyKLinesWithPreviousClose(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_daily_k_lines_with_previous_close", params,
		[]string{"security", "start", "end", "adjust"},
		nil, opts...)
}

// ListSecurityExtendedHoursIntradayBars 调用原生 ListSecurityExtendedHoursIntradayBars；请求 SecurityExtendedHoursIntradayBarsRequest，返回 list[IntradayBar]。
func ListSecurityExtendedHoursIntradayBars(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_extended_hours_intraday_bars", params,
		[]string{"security", "session", "date"},
		nil, opts...)
}

// ListSecurityExtendedHoursTicks 调用原生 ListSecurityExtendedHoursTicks；请求 SecurityExtendedHoursTicksRequest，返回 list[Tick]。
func ListSecurityExtendedHoursTicks(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_extended_hours_ticks", params,
		[]string{"security", "session", "window"},
		nil, opts...)
}

// ListSecurityFinancialSnapshots 调用原生 ListSecurityFinancialSnapshots；请求 SecurityFinancialSnapshotsRequest，返回 list[SecurityFinancialSnapshot]。
func ListSecurityFinancialSnapshots(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_financial_snapshots", params,
		[]string{"securities", "fields"},
		nil, opts...)
}

// ListSecurityIntradayBars 调用原生 ListSecurityIntradayBars；请求 IntradayBarsRequest，返回 list[IntradayBar]。
func ListSecurityIntradayBars(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_intraday_bars", params,
		[]string{"security", "date"},
		nil, opts...)
}

// ListSecurityKLines 调用原生 ListSecurityKLines；请求 KLineRequest，返回 list[KLine]。
func ListSecurityKLines(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_k_lines", params,
		[]string{"security", "start", "end", "adjust", "period"},
		nil, opts...)
}

// ListSecurityOptionGreeks 调用原生 ListSecurityOptionGreeks；请求 OptionGreeksRequest，返回 list[OptionGreeksRow]。
func ListSecurityOptionGreeks(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_option_greeks", params,
		[]string{"securities"},
		nil, opts...)
}

// ListSecurityOrderBooks 调用原生 ListSecurityOrderBooks；请求 SecurityOrderBooksRequest，返回 list[SecurityOrderBook]。
func ListSecurityOrderBooks(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_order_books", params,
		[]string{"securities", "depth"},
		nil, opts...)
}

// ListSecurityPriceVolumeLevels 调用原生 ListSecurityPriceVolumeLevels；请求 SecurityPriceVolumeLevelsRequest，返回 list[PriceVolumeLevel]。
func ListSecurityPriceVolumeLevels(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_price_volume_levels", params,
		[]string{"security"},
		nil, opts...)
}

// ListSecurityTicks 调用原生 ListSecurityTicks；请求 SecurityTicksRequest，返回 list[Tick]。
func ListSecurityTicks(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_ticks", params,
		[]string{"security", "window"},
		nil, opts...)
}

// ListSecurityTimeAndSalesTicks 调用原生 ListSecurityTimeAndSalesTicks；请求 SecurityTimeAndSalesTicksRequest，返回 list[Tick]。
func ListSecurityTimeAndSalesTicks(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_time_and_sales_ticks", params,
		[]string{"security", "count"},
		nil, opts...)
}
