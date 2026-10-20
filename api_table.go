// 本文件为 63 个只读方法生成表格结果包装：
// <GoName>Table 返回 *Table（列/行结构），
// 原生 envelope 的其余字段保存在 Table.Attrs["thsdk"]["metadata"]。

package thsdk

// CalculateSecurityRealtimeStatisticsTable(params map[string]any) 返回 RealtimeStatsResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func CalculateSecurityRealtimeStatisticsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := CalculateSecurityRealtimeStatistics(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("calculate_security_realtime_statistics", result)
}

// GetAccountWatchlistTable 返回 AccountWatchlistSnapshot 表格结果（*Table）。
// 入参：opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func GetAccountWatchlistTable(opts ...*CallOptions) (*Table, error) {
	result, err := GetAccountWatchlist(opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("get_account_watchlist", result)
}

// GetAccountWatchlistGroupsTable 返回 AccountWatchlistGroups 表格结果（*Table）。
// 入参：opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func GetAccountWatchlistGroupsTable(opts ...*CallOptions) (*Table, error) {
	result, err := GetAccountWatchlistGroups(opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("get_account_watchlist_groups", result)
}

// GetMarketMetadataTable(params map[string]any) 返回 MarketMetadata 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func GetMarketMetadataTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := GetMarketMetadata(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("get_market_metadata", result)
}

// GetMarketSecurityNamesTable(params map[string]any) 返回 MarketSecurityNamesResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func GetMarketSecurityNamesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := GetMarketSecurityNames(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("get_market_security_names", result)
}

// GetSecurityConceptTagsTable(params map[string]any) 返回 SecurityConceptTagsResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func GetSecurityConceptTagsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := GetSecurityConceptTags(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("get_security_concept_tags", result)
}

// GetSecurityPriceLimitEventsTable(params map[string]any) 返回 PriceLimitEventsResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func GetSecurityPriceLimitEventsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := GetSecurityPriceLimitEvents(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("get_security_price_limit_events", result)
}

// GetSecurityTradingTimelineTable(params map[string]any) 返回 SecurityTradingTimeline 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func GetSecurityTradingTimelineTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := GetSecurityTradingTimeline(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("get_security_trading_timeline", result)
}

// ListBlockConstituentsTable(params map[string]any) 返回 list[BlockConstituent] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListBlockConstituentsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListBlockConstituents(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_block_constituents", result)
}

// ListBlockDescriptionsTable(params map[string]any) 返回 list[BlockDescription] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListBlockDescriptionsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListBlockDescriptions(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_block_descriptions", result)
}

// ListCommodityStockLinkageNewsTable(params map[string]any) 返回 CommodityStockLinkageNewsResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListCommodityStockLinkageNewsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListCommodityStockLinkageNews(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_commodity_stock_linkage_news", result)
}

// ListFuturesRelatedSecuritiesTable 返回 list[FuturesRelatedSecurity] 表格结果（*Table）。
// 入参：opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListFuturesRelatedSecuritiesTable(opts ...*CallOptions) (*Table, error) {
	result, err := ListFuturesRelatedSecurities(opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_futures_related_securities", result)
}

// ListHotBlockCalendarEventsTable 返回 HotBlockCalendarResult 表格结果（*Table）。
// 入参：opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListHotBlockCalendarEventsTable(opts ...*CallOptions) (*Table, error) {
	result, err := ListHotBlockCalendarEvents(opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_hot_block_calendar_events", result)
}

// ListHotEventNewsTable(params map[string]any) 返回 HotEventNewsResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListHotEventNewsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListHotEventNews(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_hot_event_news", result)
}

// ListIndexCallAuctionQuotesTable(params map[string]any) 返回 list[IndexCallAuctionQuote] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListIndexCallAuctionQuotesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListIndexCallAuctionQuotes(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_index_call_auction_quotes", result)
}

// ListIndustryChildrenTable(params map[string]any) 返回 list[IndustryChild] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListIndustryChildrenTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListIndustryChildren(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_industry_children", result)
}

// ListMarketCallAuctionUnusualEventsTable(params map[string]any) 返回 list[ShortTermEvent] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListMarketCallAuctionUnusualEventsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListMarketCallAuctionUnusualEvents(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_market_call_auction_unusual_events", result)
}

// ListMarketChangeEventsTable 返回 list[MarketChangeEvent] 表格结果（*Table）。
// 入参：opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListMarketChangeEventsTable(opts ...*CallOptions) (*Table, error) {
	result, err := ListMarketChangeEvents(opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_market_change_events", result)
}

// ListMarketSecuritiesTable(params map[string]any) 返回 MarketSecuritiesResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListMarketSecuritiesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListMarketSecurities(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_market_securities", result)
}

// ListMarketShortTermEventsTable(params map[string]any) 返回 list[ShortTermEvent] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListMarketShortTermEventsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListMarketShortTermEvents(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_market_short_term_events", result)
}

// ListNewsFlashItemsTable(params map[string]any) 返回 NewsFlashItemsResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListNewsFlashItemsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListNewsFlashItems(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_news_flash_items", result)
}

// ListNewsItemsTable(params map[string]any) 返回 NewsItemsResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListNewsItemsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListNewsItems(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_news_items", result)
}

// ListRelatedSecurityPerformancesTable(params map[string]any) 返回 list[RelatedSecurityPerformance] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListRelatedSecurityPerformancesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListRelatedSecurityPerformances(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_related_security_performances", result)
}

// ListSecuritiesByRealtimeSignalTable(params map[string]any) 返回 list[RealtimeSignalSecurity] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecuritiesByRealtimeSignalTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecuritiesByRealtimeSignal(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_securities_by_realtime_signal", result)
}

// ListSecurityAHRelationsTable(params map[string]any) 返回 list[SecurityAHRelation] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityAHRelationsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityAHRelations(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_ah_relations", result)
}

// ListSecurityBlockMembershipsTable(params map[string]any) 返回 list[SecurityBlockMembership] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityBlockMembershipsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityBlockMemberships(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_block_memberships", result)
}

// ListSecurityCallAuctionQuotesTable(params map[string]any) 返回 list[CallAuctionQuote] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityCallAuctionQuotesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityCallAuctionQuotes(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_call_auction_quotes", result)
}

// ListSecurityCorporateActionsTable(params map[string]any) 返回 list[CorporateAction] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityCorporateActionsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityCorporateActions(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_corporate_actions", result)
}

// ListSecurityDailyCapitalFlowsTable(params map[string]any) 返回 list[DailyCapitalFlow] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityDailyCapitalFlowsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityDailyCapitalFlows(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_daily_capital_flows", result)
}

// ListSecurityDailyKLinesWithPreviousCloseTable(params map[string]any) 返回 list[DailyKLineWithPreviousClose] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityDailyKLinesWithPreviousCloseTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityDailyKLinesWithPreviousClose(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_daily_k_lines_with_previous_close", result)
}

// ListSecurityExtendedHoursIntradayBarsTable(params map[string]any) 返回 list[IntradayBar] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityExtendedHoursIntradayBarsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityExtendedHoursIntradayBars(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_extended_hours_intraday_bars", result)
}

// ListSecurityExtendedHoursTicksTable(params map[string]any) 返回 list[Tick] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityExtendedHoursTicksTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityExtendedHoursTicks(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_extended_hours_ticks", result)
}

// ListSecurityFinancialSnapshotsTable(params map[string]any) 返回 list[SecurityFinancialSnapshot] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityFinancialSnapshotsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityFinancialSnapshots(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_financial_snapshots", result)
}

// ListSecurityFuturesRelationsTable(securities []any) 返回 list[SecurityFuturesRelation] 表格结果（*Table）。
// 入参：opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityFuturesRelationsTable(securities []any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityFuturesRelations(securities, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_futures_relations", result)
}

// ListSecurityIndustryMappingsTable(params map[string]any) 返回 list[SecurityIndustryMapping] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityIndustryMappingsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityIndustryMappings(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_industry_mappings", result)
}

// ListSecurityIntradayBarsTable(params map[string]any) 返回 list[IntradayBar] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityIntradayBarsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityIntradayBars(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_intraday_bars", result)
}

// ListSecurityKLinesTable(params map[string]any) 返回 list[KLine] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityKLinesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityKLines(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_k_lines", result)
}

// ListSecurityLinkRelationsTable(value string) 返回 list[SecurityLinkRelation] 表格结果（*Table）。
// 入参：opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityLinkRelationsTable(value string, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityLinkRelations(value, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_link_relations", result)
}

// ListSecurityNewsEventsTable(params map[string]any) 返回 SecurityNewsItemsResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityNewsEventsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityNewsEvents(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_news_events", result)
}

// ListSecurityNewsMarkersTable(params map[string]any) 返回 SecurityNewsItemsResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityNewsMarkersTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityNewsMarkers(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_news_markers", result)
}

// ListSecurityOptionGreeksTable(params map[string]any) 返回 list[OptionGreeksRow] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityOptionGreeksTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityOptionGreeks(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_option_greeks", result)
}

// ListSecurityOrderBooksTable(params map[string]any) 返回 list[SecurityOrderBook] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityOrderBooksTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityOrderBooks(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_order_books", result)
}

// ListSecurityPriceVolumeLevelsTable(params map[string]any) 返回 list[PriceVolumeLevel] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityPriceVolumeLevelsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityPriceVolumeLevels(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_price_volume_levels", result)
}

// ListSecurityResearchReportsTable(params map[string]any) 返回 SecurityResearchReportsResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityResearchReportsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityResearchReports(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_research_reports", result)
}

// ListSecurityTicksTable(params map[string]any) 返回 list[Tick] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityTicksTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityTicks(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_ticks", result)
}

// ListSecurityTimeAndSalesTicksTable(params map[string]any) 返回 list[Tick] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListSecurityTimeAndSalesTicksTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListSecurityTimeAndSalesTicks(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_security_time_and_sales_ticks", result)
}

// ListWencaiExpressionFieldsTable(value string) 返回 list[WencaiField] 表格结果（*Table）。
// 入参：opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListWencaiExpressionFieldsTable(value string, opts ...*CallOptions) (*Table, error) {
	result, err := ListWencaiExpressionFields(value, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_wencai_expression_fields", result)
}

// ListWencaiHotBlocksTable(params map[string]any) 返回 WencaiHotBlocksResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ListWencaiHotBlocksTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ListWencaiHotBlocks(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("list_wencai_hot_blocks", result)
}

// QueryWencaiTable(params map[string]any) 返回 WencaiQueryResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func QueryWencaiTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := QueryWencai(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("query_wencai", result)
}

// QueryWencaiRealtimeFieldsTable(params map[string]any) 返回 WencaiRealtimeFieldsResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func QueryWencaiRealtimeFieldsTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := QueryWencaiRealtimeFields(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("query_wencai_realtime_fields", result)
}

// QueryWencaiSecuritiesTable(params map[string]any) 返回 WencaiSecuritiesResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func QueryWencaiSecuritiesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := QueryWencaiSecurities(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("query_wencai_securities", result)
}

// RankBlockSecuritiesTable(params map[string]any) 返回 SortedSecuritiesResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func RankBlockSecuritiesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := RankBlockSecurities(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("rank_block_securities", result)
}

// RankBlockSecuritiesByIndustryTable(params map[string]any) 返回 SortedSecuritiesResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func RankBlockSecuritiesByIndustryTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := RankBlockSecuritiesByIndustry(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("rank_block_securities_by_industry", result)
}

// RankBlockSecuritiesByWencaiFieldTable(params map[string]any) 返回 WencaiBlockRankingResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func RankBlockSecuritiesByWencaiFieldTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := RankBlockSecuritiesByWencaiField(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("rank_block_securities_by_wencai_field", result)
}

// RankConstituentsByPerformanceContributionTable(params map[string]any) 返回 list[PerformanceContributionItem] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func RankConstituentsByPerformanceContributionTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := RankConstituentsByPerformanceContribution(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("rank_constituents_by_performance_contribution", result)
}

// RankRelatedSecuritiesTable(params map[string]any) 返回 SortedSecuritiesResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func RankRelatedSecuritiesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := RankRelatedSecurities(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("rank_related_securities", result)
}

// RankSecuritiesByPopularityTable(params map[string]any) 返回 PopularityRankTopResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func RankSecuritiesByPopularityTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := RankSecuritiesByPopularity(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("rank_securities_by_popularity", result)
}

// RankSecuritiesByRealtimeStatisticTable(params map[string]any) 返回 RealtimeStatsResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func RankSecuritiesByRealtimeStatisticTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := RankSecuritiesByRealtimeStatistic(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("rank_securities_by_realtime_statistic", result)
}

// ResolveSecuritiesTable(params map[string]any) 返回 list[ResolvedSecurity] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ResolveSecuritiesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ResolveSecurities(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("resolve_securities", result)
}

// ResolveWencaiSecuritiesTable(params map[string]any) 返回 ResolvedWencaiSecuritiesResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func ResolveWencaiSecuritiesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := ResolveWencaiSecurities(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("resolve_wencai_securities", result)
}

// SearchSecuritiesTable(params map[string]any) 返回 list[SecurityCandidate] 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func SearchSecuritiesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := SearchSecurities(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("search_securities", result)
}

// SortSecuritiesTable(params map[string]any) 返回 SortedSecuritiesResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func SortSecuritiesTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := SortSecurities(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("sort_securities", result)
}

// SortSecuritiesByIndustryTable(params map[string]any) 返回 SortedSecuritiesResult 表格结果（*Table）。
// 入参：params 为请求参数映射（同对应方法，可为 nil）；opts 为调用选项（超时）。
// 出参：*Table 为表格结果；error 为调用失败时的错误。
func SortSecuritiesByIndustryTable(params map[string]any, opts ...*CallOptions) (*Table, error) {
	result, err := SortSecuritiesByIndustry(params, opts...)
	if err != nil {
		return nil, err
	}
	return tabularResult("sort_securities_by_industry", result)
}
