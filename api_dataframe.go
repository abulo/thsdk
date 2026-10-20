// Code generated from the Python THSDK 2.0.2 source (thsdk/api.py). DO NOT EDIT BY HAND.
// 本文件对应 Python 公开层为 63 个只读 DataFrame 方法生成的
// pandas.DataFrame 包装：<GoName>Table 返回 *Table（DataFrame 等价实现），
// 原生 envelope 的其余字段保存在 Table.Attrs["thsdk"]["metadata"]。

package thsdk

// CalculateSecurityRealtimeStatisticsTable(params map[string]any) 返回 RealtimeStatsResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.calculate_security_realtime_statistics（返回 pandas.DataFrame）。
func CalculateSecurityRealtimeStatisticsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := CalculateSecurityRealtimeStatistics(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("calculate_security_realtime_statistics", result)
}

// GetAccountWatchlistTable 返回 AccountWatchlistSnapshot 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.get_account_watchlist（返回 pandas.DataFrame）。
func GetAccountWatchlistTable(opts ...*CallOptions) (*Table, error) {
	result, err := GetAccountWatchlist(opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("get_account_watchlist", result)
}

// GetAccountWatchlistGroupsTable 返回 AccountWatchlistGroups 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.get_account_watchlist_groups（返回 pandas.DataFrame）。
func GetAccountWatchlistGroupsTable(opts ...*CallOptions) (*Table, error) {
	result, err := GetAccountWatchlistGroups(opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("get_account_watchlist_groups", result)
}

// GetMarketMetadataTable(params map[string]any) 返回 MarketMetadata 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.get_market_metadata（返回 pandas.DataFrame）。
func GetMarketMetadataTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := GetMarketMetadata(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("get_market_metadata", result)
}

// GetMarketSecurityNamesTable(params map[string]any) 返回 MarketSecurityNamesResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.get_market_security_names（返回 pandas.DataFrame）。
func GetMarketSecurityNamesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := GetMarketSecurityNames(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("get_market_security_names", result)
}

// GetSecurityConceptTagsTable(params map[string]any) 返回 SecurityConceptTagsResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.get_security_concept_tags（返回 pandas.DataFrame）。
func GetSecurityConceptTagsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := GetSecurityConceptTags(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("get_security_concept_tags", result)
}

// GetSecurityPriceLimitEventsTable(params map[string]any) 返回 PriceLimitEventsResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.get_security_price_limit_events（返回 pandas.DataFrame）。
func GetSecurityPriceLimitEventsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := GetSecurityPriceLimitEvents(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("get_security_price_limit_events", result)
}

// GetSecurityTradingTimelineTable(params map[string]any) 返回 SecurityTradingTimeline 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.get_security_trading_timeline（返回 pandas.DataFrame）。
func GetSecurityTradingTimelineTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := GetSecurityTradingTimeline(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("get_security_trading_timeline", result)
}

// ListBlockConstituentsTable(params map[string]any) 返回 list[BlockConstituent] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_block_constituents（返回 pandas.DataFrame）。
func ListBlockConstituentsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListBlockConstituents(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_block_constituents", result)
}

// ListBlockDescriptionsTable(params map[string]any) 返回 list[BlockDescription] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_block_descriptions（返回 pandas.DataFrame）。
func ListBlockDescriptionsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListBlockDescriptions(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_block_descriptions", result)
}

// ListCommodityStockLinkageNewsTable(params map[string]any) 返回 CommodityStockLinkageNewsResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_commodity_stock_linkage_news（返回 pandas.DataFrame）。
func ListCommodityStockLinkageNewsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListCommodityStockLinkageNews(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_commodity_stock_linkage_news", result)
}

// ListFuturesRelatedSecuritiesTable 返回 list[FuturesRelatedSecurity] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_futures_related_securities（返回 pandas.DataFrame）。
func ListFuturesRelatedSecuritiesTable(opts ...*CallOptions) (*Table, error) {
	result, err := ListFuturesRelatedSecurities(opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_futures_related_securities", result)
}

// ListHotBlockCalendarEventsTable 返回 HotBlockCalendarResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_hot_block_calendar_events（返回 pandas.DataFrame）。
func ListHotBlockCalendarEventsTable(opts ...*CallOptions) (*Table, error) {
	result, err := ListHotBlockCalendarEvents(opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_hot_block_calendar_events", result)
}

// ListHotEventNewsTable(params map[string]any) 返回 HotEventNewsResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_hot_event_news（返回 pandas.DataFrame）。
func ListHotEventNewsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListHotEventNews(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_hot_event_news", result)
}

// ListIndexCallAuctionQuotesTable(params map[string]any) 返回 list[IndexCallAuctionQuote] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_index_call_auction_quotes（返回 pandas.DataFrame）。
func ListIndexCallAuctionQuotesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListIndexCallAuctionQuotes(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_index_call_auction_quotes", result)
}

// ListIndustryChildrenTable(params map[string]any) 返回 list[IndustryChild] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_industry_children（返回 pandas.DataFrame）。
func ListIndustryChildrenTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListIndustryChildren(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_industry_children", result)
}

// ListMarketCallAuctionUnusualEventsTable(params map[string]any) 返回 list[ShortTermEvent] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_market_call_auction_unusual_events（返回 pandas.DataFrame）。
func ListMarketCallAuctionUnusualEventsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListMarketCallAuctionUnusualEvents(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_market_call_auction_unusual_events", result)
}

// ListMarketChangeEventsTable 返回 list[MarketChangeEvent] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_market_change_events（返回 pandas.DataFrame）。
func ListMarketChangeEventsTable(opts ...*CallOptions) (*Table, error) {
	result, err := ListMarketChangeEvents(opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_market_change_events", result)
}

// ListMarketSecuritiesTable(params map[string]any) 返回 MarketSecuritiesResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_market_securities（返回 pandas.DataFrame）。
func ListMarketSecuritiesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListMarketSecurities(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_market_securities", result)
}

// ListMarketShortTermEventsTable(params map[string]any) 返回 list[ShortTermEvent] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_market_short_term_events（返回 pandas.DataFrame）。
func ListMarketShortTermEventsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListMarketShortTermEvents(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_market_short_term_events", result)
}

// ListNewsFlashItemsTable(params map[string]any) 返回 NewsFlashItemsResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_news_flash_items（返回 pandas.DataFrame）。
func ListNewsFlashItemsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListNewsFlashItems(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_news_flash_items", result)
}

// ListNewsItemsTable(params map[string]any) 返回 NewsItemsResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_news_items（返回 pandas.DataFrame）。
func ListNewsItemsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListNewsItems(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_news_items", result)
}

// ListRelatedSecurityPerformancesTable(params map[string]any) 返回 list[RelatedSecurityPerformance] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_related_security_performances（返回 pandas.DataFrame）。
func ListRelatedSecurityPerformancesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListRelatedSecurityPerformances(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_related_security_performances", result)
}

// ListSecuritiesByRealtimeSignalTable(params map[string]any) 返回 list[RealtimeSignalSecurity] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_securities_by_realtime_signal（返回 pandas.DataFrame）。
func ListSecuritiesByRealtimeSignalTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecuritiesByRealtimeSignal(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_securities_by_realtime_signal", result)
}

// ListSecurityAHRelationsTable(params map[string]any) 返回 list[SecurityAHRelation] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_ah_relations（返回 pandas.DataFrame）。
func ListSecurityAHRelationsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityAHRelations(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_ah_relations", result)
}

// ListSecurityBlockMembershipsTable(params map[string]any) 返回 list[SecurityBlockMembership] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_block_memberships（返回 pandas.DataFrame）。
func ListSecurityBlockMembershipsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityBlockMemberships(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_block_memberships", result)
}

// ListSecurityCallAuctionQuotesTable(params map[string]any) 返回 list[CallAuctionQuote] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_call_auction_quotes（返回 pandas.DataFrame）。
func ListSecurityCallAuctionQuotesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityCallAuctionQuotes(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_call_auction_quotes", result)
}

// ListSecurityCorporateActionsTable(params map[string]any) 返回 list[CorporateAction] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_corporate_actions（返回 pandas.DataFrame）。
func ListSecurityCorporateActionsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityCorporateActions(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_corporate_actions", result)
}

// ListSecurityDailyCapitalFlowsTable(params map[string]any) 返回 list[DailyCapitalFlow] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_daily_capital_flows（返回 pandas.DataFrame）。
func ListSecurityDailyCapitalFlowsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityDailyCapitalFlows(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_daily_capital_flows", result)
}

// ListSecurityDailyKLinesWithPreviousCloseTable(params map[string]any) 返回 list[DailyKLineWithPreviousClose] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_daily_k_lines_with_previous_close（返回 pandas.DataFrame）。
func ListSecurityDailyKLinesWithPreviousCloseTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityDailyKLinesWithPreviousClose(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_daily_k_lines_with_previous_close", result)
}

// ListSecurityExtendedHoursIntradayBarsTable(params map[string]any) 返回 list[IntradayBar] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_extended_hours_intraday_bars（返回 pandas.DataFrame）。
func ListSecurityExtendedHoursIntradayBarsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityExtendedHoursIntradayBars(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_extended_hours_intraday_bars", result)
}

// ListSecurityExtendedHoursTicksTable(params map[string]any) 返回 list[Tick] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_extended_hours_ticks（返回 pandas.DataFrame）。
func ListSecurityExtendedHoursTicksTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityExtendedHoursTicks(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_extended_hours_ticks", result)
}

// ListSecurityFinancialSnapshotsTable(params map[string]any) 返回 list[SecurityFinancialSnapshot] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_financial_snapshots（返回 pandas.DataFrame）。
func ListSecurityFinancialSnapshotsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityFinancialSnapshots(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_financial_snapshots", result)
}

// ListSecurityFuturesRelationsTable(securities []any) 返回 list[SecurityFuturesRelation] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_futures_relations（返回 pandas.DataFrame）。
func ListSecurityFuturesRelationsTable(securities []any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityFuturesRelations(securities, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_futures_relations", result)
}

// ListSecurityIndustryMappingsTable(params map[string]any) 返回 list[SecurityIndustryMapping] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_industry_mappings（返回 pandas.DataFrame）。
func ListSecurityIndustryMappingsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityIndustryMappings(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_industry_mappings", result)
}

// ListSecurityIntradayBarsTable(params map[string]any) 返回 list[IntradayBar] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_intraday_bars（返回 pandas.DataFrame）。
func ListSecurityIntradayBarsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityIntradayBars(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_intraday_bars", result)
}

// ListSecurityKLinesTable(params map[string]any) 返回 list[KLine] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_k_lines（返回 pandas.DataFrame）。
func ListSecurityKLinesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityKLines(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_k_lines", result)
}

// ListSecurityLinkRelationsTable(value string) 返回 list[SecurityLinkRelation] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_link_relations（返回 pandas.DataFrame）。
func ListSecurityLinkRelationsTable(value string, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityLinkRelations(value, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_link_relations", result)
}

// ListSecurityNewsEventsTable(params map[string]any) 返回 SecurityNewsItemsResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_news_events（返回 pandas.DataFrame）。
func ListSecurityNewsEventsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityNewsEvents(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_news_events", result)
}

// ListSecurityNewsMarkersTable(params map[string]any) 返回 SecurityNewsItemsResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_news_markers（返回 pandas.DataFrame）。
func ListSecurityNewsMarkersTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityNewsMarkers(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_news_markers", result)
}

// ListSecurityOptionGreeksTable(params map[string]any) 返回 list[OptionGreeksRow] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_option_greeks（返回 pandas.DataFrame）。
func ListSecurityOptionGreeksTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityOptionGreeks(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_option_greeks", result)
}

// ListSecurityOrderBooksTable(params map[string]any) 返回 list[SecurityOrderBook] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_order_books（返回 pandas.DataFrame）。
func ListSecurityOrderBooksTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityOrderBooks(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_order_books", result)
}

// ListSecurityPriceVolumeLevelsTable(params map[string]any) 返回 list[PriceVolumeLevel] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_price_volume_levels（返回 pandas.DataFrame）。
func ListSecurityPriceVolumeLevelsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityPriceVolumeLevels(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_price_volume_levels", result)
}

// ListSecurityResearchReportsTable(params map[string]any) 返回 SecurityResearchReportsResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_research_reports（返回 pandas.DataFrame）。
func ListSecurityResearchReportsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityResearchReports(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_research_reports", result)
}

// ListSecurityTicksTable(params map[string]any) 返回 list[Tick] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_ticks（返回 pandas.DataFrame）。
func ListSecurityTicksTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityTicks(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_ticks", result)
}

// ListSecurityTimeAndSalesTicksTable(params map[string]any) 返回 list[Tick] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_security_time_and_sales_ticks（返回 pandas.DataFrame）。
func ListSecurityTimeAndSalesTicksTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityTimeAndSalesTicks(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_security_time_and_sales_ticks", result)
}

// ListWencaiExpressionFieldsTable(value string) 返回 list[WencaiField] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_wencai_expression_fields（返回 pandas.DataFrame）。
func ListWencaiExpressionFieldsTable(value string, opts ...*CallOptions) (*Table, error) {
	result, err := ListWencaiExpressionFields(value, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_wencai_expression_fields", result)
}

// ListWencaiHotBlocksTable(params map[string]any) 返回 WencaiHotBlocksResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.list_wencai_hot_blocks（返回 pandas.DataFrame）。
func ListWencaiHotBlocksTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListWencaiHotBlocks(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("list_wencai_hot_blocks", result)
}

// QueryWencaiTable(params map[string]any) 返回 WencaiQueryResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.query_wencai（返回 pandas.DataFrame）。
func QueryWencaiTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := QueryWencai(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("query_wencai", result)
}

// QueryWencaiRealtimeFieldsTable(params map[string]any) 返回 WencaiRealtimeFieldsResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.query_wencai_realtime_fields（返回 pandas.DataFrame）。
func QueryWencaiRealtimeFieldsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := QueryWencaiRealtimeFields(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("query_wencai_realtime_fields", result)
}

// QueryWencaiSecuritiesTable(params map[string]any) 返回 WencaiSecuritiesResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.query_wencai_securities（返回 pandas.DataFrame）。
func QueryWencaiSecuritiesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := QueryWencaiSecurities(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("query_wencai_securities", result)
}

// RankBlockSecuritiesTable(params map[string]any) 返回 SortedSecuritiesResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.rank_block_securities（返回 pandas.DataFrame）。
func RankBlockSecuritiesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := RankBlockSecurities(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("rank_block_securities", result)
}

// RankBlockSecuritiesByIndustryTable(params map[string]any) 返回 SortedSecuritiesResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.rank_block_securities_by_industry（返回 pandas.DataFrame）。
func RankBlockSecuritiesByIndustryTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := RankBlockSecuritiesByIndustry(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("rank_block_securities_by_industry", result)
}

// RankBlockSecuritiesByWencaiFieldTable(params map[string]any) 返回 WencaiBlockRankingResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.rank_block_securities_by_wencai_field（返回 pandas.DataFrame）。
func RankBlockSecuritiesByWencaiFieldTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := RankBlockSecuritiesByWencaiField(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("rank_block_securities_by_wencai_field", result)
}

// RankConstituentsByPerformanceContributionTable(params map[string]any) 返回 list[PerformanceContributionItem] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.rank_constituents_by_performance_contribution（返回 pandas.DataFrame）。
func RankConstituentsByPerformanceContributionTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := RankConstituentsByPerformanceContribution(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("rank_constituents_by_performance_contribution", result)
}

// RankRelatedSecuritiesTable(params map[string]any) 返回 SortedSecuritiesResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.rank_related_securities（返回 pandas.DataFrame）。
func RankRelatedSecuritiesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := RankRelatedSecurities(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("rank_related_securities", result)
}

// RankSecuritiesByPopularityTable(params map[string]any) 返回 PopularityRankTopResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.rank_securities_by_popularity（返回 pandas.DataFrame）。
func RankSecuritiesByPopularityTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := RankSecuritiesByPopularity(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("rank_securities_by_popularity", result)
}

// RankSecuritiesByRealtimeStatisticTable(params map[string]any) 返回 RealtimeStatsResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.rank_securities_by_realtime_statistic（返回 pandas.DataFrame）。
func RankSecuritiesByRealtimeStatisticTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := RankSecuritiesByRealtimeStatistic(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("rank_securities_by_realtime_statistic", result)
}

// ResolveSecuritiesTable(params map[string]any) 返回 list[ResolvedSecurity] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.resolve_securities（返回 pandas.DataFrame）。
func ResolveSecuritiesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ResolveSecurities(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("resolve_securities", result)
}

// ResolveWencaiSecuritiesTable(params map[string]any) 返回 ResolvedWencaiSecuritiesResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.resolve_wencai_securities（返回 pandas.DataFrame）。
func ResolveWencaiSecuritiesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ResolveWencaiSecurities(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("resolve_wencai_securities", result)
}

// SearchSecuritiesTable(params map[string]any) 返回 list[SecurityCandidate] 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.search_securities（返回 pandas.DataFrame）。
func SearchSecuritiesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := SearchSecurities(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("search_securities", result)
}

// SortSecuritiesTable(params map[string]any) 返回 SortedSecuritiesResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.sort_securities（返回 pandas.DataFrame）。
func SortSecuritiesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := SortSecurities(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("sort_securities", result)
}

// SortSecuritiesByIndustryTable(params map[string]any) 返回 SortedSecuritiesResult 的 DataFrame 等价表格（*Table）。
// 对应 Python 公开层 thsdk.sort_securities_by_industry（返回 pandas.DataFrame）。
func SortSecuritiesByIndustryTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := SortSecuritiesByIndustry(params, opts...)
	if err != nil {
		return nil, err
	}
	return dataframeResult("sort_securities_by_industry", result)
}
