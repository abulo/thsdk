package thsdk

// 本文件登记 80 项同步 Canonical API 的稳定元数据：请求字段、必填字段、
// 结果类型、表格主集合字段等，并在 init() 中执行模块加载时的完整性校验
// （80 项唯一、63 个只读表格方法、契约一致等）。

import (
	"sort"
	"strings"
)

// RequestKind 描述一个 Canonical 方法的请求形态。
type RequestKind int

const (
	RequestKindNone       RequestKind = iota // 无请求体
	RequestKindObject                        // 对象参数
	RequestKindString                        // 单个字符串参数
	RequestKindSecurities                    // 证券序列参数
)

// NativeMethodSpec 描述一个可调用 Canonical 方法的稳定元数据。
type NativeMethodSpec struct {
	Name             string
	GoName           string
	Category         string
	RequestKind      RequestKind
	RequestType      string
	ResultType       string
	Mutating         bool
	RequestFields    []string
	RequiredFields   []string
	ReturnsTabular   bool
	TabularRowsField string
}

// requestContract 是 _REQUEST_CONTRACTS 请求契约中的单条记录（全部字段 + 必填字段）。
type requestContract struct {
	fields   []string
	required []string
}

// requestContracts 是全部请求契约的静态表，与原生侧逐字一致。
var requestContracts = map[string]requestContract{
	"account_permissions":                             {fields: []string{}, required: []string{}},
	"add_account_watchlist_group_securities":          {fields: []string{"group_id", "securities"}, required: []string{"group_id", "securities"}},
	"add_account_watchlist_securities":                {fields: []string{"securities", "add_to_front"}, required: []string{"securities"}},
	"analyze_security_limit_up":                       {fields: []string{"security", "date"}, required: []string{"security"}},
	"calculate_security_realtime_statistics":          {fields: []string{"securities", "metrics"}, required: []string{"securities", "metrics"}},
	"check_market_timeline_readiness":                 {fields: []string{"market", "time_index_day"}, required: []string{"market", "time_index_day"}},
	"clear_account_watchlist":                         {fields: []string{}, required: []string{}},
	"create_account_watchlist_group":                  {fields: []string{"name", "securities"}, required: []string{"name"}},
	"delete_account_watchlist_group":                  {fields: []string{"group_id"}, required: []string{"group_id"}},
	"get_account_watchlist":                           {fields: []string{}, required: []string{}},
	"get_account_watchlist_groups":                    {fields: []string{}, required: []string{}},
	"get_market_metadata":                             {fields: []string{"market", "version_flag"}, required: []string{"market"}},
	"get_market_security_names":                       {fields: []string{"market", "realtime_version_ini", "history_version_ini", "base_version_ini"}, required: []string{"market"}},
	"get_security_concept_tags":                       {fields: []string{"security"}, required: []string{"security"}},
	"get_security_industry":                           {fields: []string{"security"}, required: []string{"security"}},
	"get_security_popularity_rank":                    {fields: []string{"security", "benchmark"}, required: []string{"security"}},
	"get_security_price_limit_events":                 {fields: []string{"security", "start_year", "end_year"}, required: []string{"security"}},
	"get_security_short_term_highlights":              {fields: []string{"security"}, required: []string{"security"}},
	"get_security_trading_timeline":                   {fields: []string{"security"}, required: []string{"security"}},
	"list_block_constituents":                         {fields: []string{"block", "sort_begin", "sort_count", "sort_order", "sort_id"}, required: []string{"block"}},
	"list_block_descriptions":                         {fields: []string{"blocks"}, required: []string{"blocks"}},
	"list_commodity_stock_linkage_news":               {fields: []string{"commodity_index"}, required: []string{"commodity_index"}},
	"list_futures_related_securities":                 {fields: []string{}, required: []string{}},
	"list_hot_block_calendar_events":                  {fields: []string{}, required: []string{}},
	"list_hot_event_news":                             {fields: []string{"block", "range"}, required: []string{"block"}},
	"list_index_call_auction_quotes":                  {fields: []string{"index", "phase"}, required: []string{"index", "phase"}},
	"list_industry_children":                          {fields: []string{"industry"}, required: []string{"industry"}},
	"list_market_call_auction_unusual_events":         {fields: []string{"market", "trade_date"}, required: []string{"market"}},
	"list_market_change_events":                       {fields: []string{}, required: []string{}},
	"list_market_securities":                          {fields: []string{"market", "sort_begin", "sort_count", "sort_id", "sort_order", "func_period"}, required: []string{"market"}},
	"list_market_short_term_events":                   {fields: []string{"market", "mode", "securities", "max_count", "before_timestamp_micros"}, required: []string{"market", "mode"}},
	"list_news_flash_items":                           {fields: []string{"traversal", "tag", "tag_id", "page", "page_size", "page_time", "sequence", "created_at", "environment"}, required: []string{"traversal"}},
	"list_news_items":                                 {fields: []string{"category", "security", "market", "traversal", "summary_mode", "advertorial"}, required: []string{"category", "traversal"}},
	"list_related_security_performances":              {fields: []string{"security"}, required: []string{"security"}},
	"list_securities_by_realtime_signal":              {fields: []string{"markets", "signal"}, required: []string{"markets", "signal"}},
	"list_security_ah_relations":                      {fields: []string{"security"}, required: []string{"security"}},
	"list_security_block_memberships":                 {fields: []string{"security"}, required: []string{"security"}},
	"list_security_call_auction_quotes":               {fields: []string{"security", "phase", "date", "window"}, required: []string{"security", "phase"}},
	"list_security_corporate_actions":                 {fields: []string{"security", "start", "end", "adjust", "period"}, required: []string{"security", "start", "end", "adjust", "period"}},
	"list_security_daily_capital_flows":               {fields: []string{"securities"}, required: []string{"securities"}},
	"list_security_daily_k_lines_with_previous_close": {fields: []string{"security", "start", "end", "adjust"}, required: []string{"security", "start", "end", "adjust"}},
	"list_security_extended_hours_intraday_bars":      {fields: []string{"security", "session", "date"}, required: []string{"security", "session", "date"}},
	"list_security_extended_hours_ticks":              {fields: []string{"security", "session", "window"}, required: []string{"security", "session", "window"}},
	"list_security_financial_snapshots":               {fields: []string{"securities", "fields"}, required: []string{"securities", "fields"}},
	"list_security_futures_relations":                 {fields: []string{"securities"}, required: []string{"securities"}},
	"list_security_industry_mappings":                 {fields: []string{"securities"}, required: []string{"securities"}},
	"list_security_intraday_bars":                     {fields: []string{"security", "date"}, required: []string{"security", "date"}},
	"list_security_k_lines":                           {fields: []string{"security", "start", "end", "adjust", "period"}, required: []string{"security", "start", "end", "adjust", "period"}},
	"list_security_link_relations":                    {fields: []string{"link_key"}, required: []string{"link_key"}},
	"list_security_news_events":                       {fields: []string{"security", "timeline", "window"}, required: []string{"security", "timeline", "window"}},
	"list_security_news_markers":                      {fields: []string{"security", "timeline"}, required: []string{"security", "timeline"}},
	"list_security_option_greeks":                     {fields: []string{"securities", "metrics"}, required: []string{"securities"}},
	"list_security_order_books":                       {fields: []string{"securities", "depth"}, required: []string{"securities", "depth"}},
	"list_security_price_volume_levels":               {fields: []string{"security"}, required: []string{"security"}},
	"list_security_research_reports":                  {fields: []string{"security", "traversal"}, required: []string{"security", "traversal"}},
	"list_security_ticks":                             {fields: []string{"security", "window"}, required: []string{"security", "window"}},
	"list_security_time_and_sales_ticks":              {fields: []string{"security", "count"}, required: []string{"security", "count"}},
	"list_wencai_expression_fields":                   {fields: []string{"query"}, required: []string{"query"}},
	"list_wencai_hot_blocks":                          {fields: []string{"kind"}, required: []string{"kind"}},
	"query_wencai":                                    {fields: []string{"query", "markets", "limit"}, required: []string{"query"}},
	"query_wencai_realtime_fields":                    {fields: []string{"securities", "fields"}, required: []string{"securities", "fields"}},
	"query_wencai_securities":                         {fields: []string{"query", "markets", "limit"}, required: []string{"query"}},
	"rank_block_securities":                           {fields: []string{"block", "sort_begin", "sort_count", "sort_id", "sort_order", "exclude_securities", "exclude_blocks", "func_period"}, required: []string{"block"}},
	"rank_block_securities_by_industry":               {fields: []string{"block", "sort_begin", "sort_count", "sort_order"}, required: []string{"block"}},
	"rank_block_securities_by_wencai_field":           {fields: []string{"block", "field", "offset", "limit", "sort_order"}, required: []string{"block", "field"}},
	"rank_constituents_by_performance_contribution":   {fields: []string{"code", "market", "target", "sort_id", "sort_order", "sort_begin", "sort_count", "valid_begin", "valid_end"}, required: []string{"code", "market"}},
	"rank_related_securities":                         {fields: []string{"security", "sort_begin", "sort_count", "sort_id", "sort_order", "func_period"}, required: []string{"security"}},
	"rank_securities_by_popularity":                   {fields: []string{"type", "benchmark"}, required: []string{}},
	"rank_securities_by_realtime_statistic":           {fields: []string{"markets", "securities", "sort_by", "sort_dir", "sort_begin", "sort_count", "exclude_securities"}, required: []string{"sort_by"}},
	"remove_account_watchlist_group_securities":       {fields: []string{"group_id", "securities"}, required: []string{"group_id", "securities"}},
	"remove_account_watchlist_securities":             {fields: []string{"securities"}, required: []string{"securities"}},
	"rename_account_watchlist_group":                  {fields: []string{"group_id", "name"}, required: []string{"group_id", "name"}},
	"replace_account_watchlist_group_securities":      {fields: []string{"group_id", "securities"}, required: []string{"group_id", "securities"}},
	"replace_account_watchlist_securities":            {fields: []string{"version", "securities"}, required: []string{"version", "securities"}},
	"resolve_block":                                   {fields: []string{"name"}, required: []string{"name"}},
	"resolve_securities":                              {fields: []string{"codes"}, required: []string{"codes"}},
	"resolve_wencai_securities":                       {fields: []string{"query", "limit", "stock_suffix"}, required: []string{"query"}},
	"search_securities":                               {fields: []string{"pattern", "market"}, required: []string{"pattern"}},
	"sort_securities":                                 {fields: []string{"securities", "sort_begin", "sort_count", "sort_id", "sort_order", "func_period"}, required: []string{"securities"}},
	"sort_securities_by_industry":                     {fields: []string{"securities", "sort_begin", "sort_count", "sort_order"}, required: []string{"securities"}},
}

// tabularRowsByResultType
// 包装结果类型中的主记录集合字段。原生直接返回 list[T] 的方法不需要在此声明。
var tabularRowsByResultType = map[string]string{
	"AccountWatchlistGroups":          "groups",
	"AccountWatchlistSnapshot":        "securities",
	"CommodityStockLinkageNewsResult": "items",
	"HotBlockCalendarResult":          "items",
	"HotEventNewsResult":              "items",
	"MarketMetadata":                  "sections",
	"MarketSecuritiesResult":          "securities",
	"MarketSecurityNamesResult":       "items",
	"NewsFlashItemsResult":            "items",
	"NewsItemsResult":                 "items",
	"PopularityRankTopResult":         "items",
	"PriceLimitEventsResult":          "events",
	"RealtimeStatsResult":             "rows",
	"ResolvedWencaiSecuritiesResult":  "items",
	"SecurityConceptTagsResult":       "tags",
	"SecurityNewsItemsResult":         "items",
	"SecurityResearchReportsResult":   "items",
	"SecurityTradingTimeline":         "items",
	"SortedSecuritiesResult":          "securities",
	"WencaiBlockRankingResult":        "items",
	"WencaiHotBlocksResult":           "items",
	"WencaiQueryResult":               "rows",
	"WencaiRealtimeFieldsResult":      "rows",
	"WencaiSecuritiesResult":          "items",
}

func newSpec(name, goName, category string, kind RequestKind, requestType, resultType string, mutating bool) NativeMethodSpec {
	contract := requestContracts[name]
	rowsField := tabularRowsByResultType[resultType]
	spec := NativeMethodSpec{
		Name:             name,
		GoName:           goName,
		Category:         category,
		RequestKind:      kind,
		RequestType:      requestType,
		ResultType:       resultType,
		Mutating:         mutating,
		RequestFields:    contract.fields,
		RequiredFields:   contract.required,
		TabularRowsField: rowsField,
	}
	// 写操作不参与表格包装。
	if mutating {
		spec.TabularRowsField = ""
	}
	spec.ReturnsTabular = !mutating && (len(resultType) >= 5 && resultType[:5] == "list[" || rowsField != "")
	return spec
}

// nativeMethodSpecs _SPECS，共 80 项。
var nativeMethodSpecs = []NativeMethodSpec{
	newSpec("account_permissions", "AccountPermissions", "account", RequestKindNone, "", "AccountPermissions", false),
	newSpec("add_account_watchlist_group_securities", "AddAccountWatchlistGroupSecurities", "account", RequestKindObject, "AccountWatchlistGroupSecuritiesRequest", "AccountWatchlistGroups", true),
	newSpec("add_account_watchlist_securities", "AddAccountWatchlistSecurities", "account", RequestKindObject, "AddAccountWatchlistSecuritiesRequest", "AccountWatchlistUpdate", true),
	newSpec("analyze_security_limit_up", "AnalyzeSecurityLimitUp", "news", RequestKindObject, "SecurityLimitUpAnalysisRequest", "SecurityLimitUpAnalysisResult", false),
	newSpec("calculate_security_realtime_statistics", "CalculateSecurityRealtimeStatistics", "rankings", RequestKindObject, "RealtimeStatsCalculationRequest", "RealtimeStatsResult", false),
	newSpec("check_market_timeline_readiness", "CheckMarketTimelineReadiness", "market_data", RequestKindObject, "MarketTimelineReadinessRequest", "MarketTimelineReadinessResult", false),
	newSpec("clear_account_watchlist", "ClearAccountWatchlist", "account", RequestKindNone, "", "AccountWatchlistUpdate", true),
	newSpec("create_account_watchlist_group", "CreateAccountWatchlistGroup", "account", RequestKindObject, "CreateAccountWatchlistGroupRequest", "AccountWatchlistGroups", true),
	newSpec("delete_account_watchlist_group", "DeleteAccountWatchlistGroup", "account", RequestKindObject, "DeleteAccountWatchlistGroupRequest", "AccountWatchlistGroups", true),
	newSpec("get_account_watchlist", "GetAccountWatchlist", "account", RequestKindNone, "", "AccountWatchlistSnapshot", false),
	newSpec("get_account_watchlist_groups", "GetAccountWatchlistGroups", "account", RequestKindNone, "", "AccountWatchlistGroups", false),
	newSpec("get_market_metadata", "GetMarketMetadata", "securities", RequestKindObject, "MarketMetadataRequest", "MarketMetadata", false),
	newSpec("get_market_security_names", "GetMarketSecurityNames", "securities", RequestKindObject, "MarketSecurityNamesRequest", "MarketSecurityNamesResult", false),
	newSpec("get_security_concept_tags", "GetSecurityConceptTags", "securities", RequestKindObject, "SecurityConceptTagsRequest", "SecurityConceptTagsResult", false),
	newSpec("get_security_industry", "GetSecurityIndustry", "securities", RequestKindObject, "SecurityRequest", "SecurityIndustryMapping", false),
	newSpec("get_security_popularity_rank", "GetSecurityPopularityRank", "rankings", RequestKindObject, "PopularityRankRequest", "PopularityRankResult", false),
	newSpec("get_security_price_limit_events", "GetSecurityPriceLimitEvents", "news", RequestKindObject, "PriceLimitEventsRequest", "PriceLimitEventsResult", false),
	newSpec("get_security_short_term_highlights", "GetSecurityShortTermHighlights", "news", RequestKindObject, "SecurityShortTermHighlightsRequest", "SecurityShortTermHighlightsResult", false),
	newSpec("get_security_trading_timeline", "GetSecurityTradingTimeline", "market_data", RequestKindObject, "SecurityTradingTimelineRequest", "SecurityTradingTimeline", false),
	newSpec("list_block_constituents", "ListBlockConstituents", "securities", RequestKindObject, "BlockConstituentsRequest", "list[BlockConstituent]", false),
	newSpec("list_block_descriptions", "ListBlockDescriptions", "securities", RequestKindObject, "BlockDescriptionsRequest", "list[BlockDescription]", false),
	newSpec("list_commodity_stock_linkage_news", "ListCommodityStockLinkageNews", "news", RequestKindObject, "CommodityStockLinkageNewsRequest", "CommodityStockLinkageNewsResult", false),
	newSpec("list_futures_related_securities", "ListFuturesRelatedSecurities", "securities", RequestKindNone, "", "list[FuturesRelatedSecurity]", false),
	newSpec("list_hot_block_calendar_events", "ListHotBlockCalendarEvents", "news", RequestKindNone, "", "HotBlockCalendarResult", false),
	newSpec("list_hot_event_news", "ListHotEventNews", "news", RequestKindObject, "HotEventNewsRequest", "HotEventNewsResult", false),
	newSpec("list_index_call_auction_quotes", "ListIndexCallAuctionQuotes", "market_data", RequestKindObject, "IndexCallAuctionQuotesRequest", "list[IndexCallAuctionQuote]", false),
	newSpec("list_industry_children", "ListIndustryChildren", "securities", RequestKindObject, "IndustryChildrenRequest", "list[IndustryChild]", false),
	newSpec("list_market_call_auction_unusual_events", "ListMarketCallAuctionUnusualEvents", "news", RequestKindObject, "CallAuctionUnusualEventsRequest", "list[ShortTermEvent]", false),
	newSpec("list_market_change_events", "ListMarketChangeEvents", "news", RequestKindNone, "", "list[MarketChangeEvent]", false),
	newSpec("list_market_securities", "ListMarketSecurities", "securities", RequestKindObject, "MarketSecuritiesRequest", "MarketSecuritiesResult", false),
	newSpec("list_market_short_term_events", "ListMarketShortTermEvents", "news", RequestKindObject, "MarketShortTermEventsRequest", "list[ShortTermEvent]", false),
	newSpec("list_news_flash_items", "ListNewsFlashItems", "news", RequestKindObject, "NewsFlashItemsRequest", "NewsFlashItemsResult", false),
	newSpec("list_news_items", "ListNewsItems", "news", RequestKindObject, "NewsItemsRequest", "NewsItemsResult", false),
	newSpec("list_related_security_performances", "ListRelatedSecurityPerformances", "securities", RequestKindObject, "RelatedSecurityPerformancesRequest", "list[RelatedSecurityPerformance]", false),
	newSpec("list_securities_by_realtime_signal", "ListSecuritiesByRealtimeSignal", "rankings", RequestKindObject, "RealtimeSignalSecuritiesRequest", "list[RealtimeSignalSecurity]", false),
	newSpec("list_security_ah_relations", "ListSecurityAHRelations", "securities", RequestKindObject, "SecurityAHRelationsRequest", "list[SecurityAHRelation]", false),
	newSpec("list_security_block_memberships", "ListSecurityBlockMemberships", "securities", RequestKindObject, "SecurityBlockMembershipsRequest", "list[SecurityBlockMembership]", false),
	newSpec("list_security_call_auction_quotes", "ListSecurityCallAuctionQuotes", "market_data", RequestKindObject, "SecurityCallAuctionQuotesRequest", "list[CallAuctionQuote]", false),
	newSpec("list_security_corporate_actions", "ListSecurityCorporateActions", "market_data", RequestKindObject, "SecurityCorporateActionsRequest", "list[CorporateAction]", false),
	newSpec("list_security_daily_capital_flows", "ListSecurityDailyCapitalFlows", "market_data", RequestKindObject, "SecurityDailyCapitalFlowsRequest", "list[DailyCapitalFlow]", false),
	newSpec("list_security_daily_k_lines_with_previous_close", "ListSecurityDailyKLinesWithPreviousClose", "market_data", RequestKindObject, "DailyKLineRequest", "list[DailyKLineWithPreviousClose]", false),
	newSpec("list_security_extended_hours_intraday_bars", "ListSecurityExtendedHoursIntradayBars", "market_data", RequestKindObject, "SecurityExtendedHoursIntradayBarsRequest", "list[IntradayBar]", false),
	newSpec("list_security_extended_hours_ticks", "ListSecurityExtendedHoursTicks", "market_data", RequestKindObject, "SecurityExtendedHoursTicksRequest", "list[Tick]", false),
	newSpec("list_security_financial_snapshots", "ListSecurityFinancialSnapshots", "market_data", RequestKindObject, "SecurityFinancialSnapshotsRequest", "list[SecurityFinancialSnapshot]", false),
	newSpec("list_security_futures_relations", "ListSecurityFuturesRelations", "securities", RequestKindSecurities, "list[Security]", "list[SecurityFuturesRelation]", false),
	newSpec("list_security_industry_mappings", "ListSecurityIndustryMappings", "securities", RequestKindObject, "SecurityIndustryMappingsRequest", "list[SecurityIndustryMapping]", false),
	newSpec("list_security_intraday_bars", "ListSecurityIntradayBars", "market_data", RequestKindObject, "IntradayBarsRequest", "list[IntradayBar]", false),
	newSpec("list_security_k_lines", "ListSecurityKLines", "market_data", RequestKindObject, "KLineRequest", "list[KLine]", false),
	newSpec("list_security_link_relations", "ListSecurityLinkRelations", "securities", RequestKindString, "linkKey string", "list[SecurityLinkRelation]", false),
	newSpec("list_security_news_events", "ListSecurityNewsEvents", "news", RequestKindObject, "SecurityNewsEventsRequest", "SecurityNewsItemsResult", false),
	newSpec("list_security_news_markers", "ListSecurityNewsMarkers", "news", RequestKindObject, "SecurityNewsMarkersRequest", "SecurityNewsItemsResult", false),
	newSpec("list_security_option_greeks", "ListSecurityOptionGreeks", "market_data", RequestKindObject, "OptionGreeksRequest", "list[OptionGreeksRow]", false),
	newSpec("list_security_order_books", "ListSecurityOrderBooks", "market_data", RequestKindObject, "SecurityOrderBooksRequest", "list[SecurityOrderBook]", false),
	newSpec("list_security_price_volume_levels", "ListSecurityPriceVolumeLevels", "market_data", RequestKindObject, "SecurityPriceVolumeLevelsRequest", "list[PriceVolumeLevel]", false),
	newSpec("list_security_research_reports", "ListSecurityResearchReports", "news", RequestKindObject, "SecurityResearchReportsRequest", "SecurityResearchReportsResult", false),
	newSpec("list_security_ticks", "ListSecurityTicks", "market_data", RequestKindObject, "SecurityTicksRequest", "list[Tick]", false),
	newSpec("list_security_time_and_sales_ticks", "ListSecurityTimeAndSalesTicks", "market_data", RequestKindObject, "SecurityTimeAndSalesTicksRequest", "list[Tick]", false),
	newSpec("list_wencai_expression_fields", "ListWencaiExpressionFields", "wencai", RequestKindString, "query string", "list[WencaiField]", false),
	newSpec("list_wencai_hot_blocks", "ListWencaiHotBlocks", "wencai", RequestKindObject, "WencaiHotBlocksRequest", "WencaiHotBlocksResult", false),
	newSpec("query_wencai", "QueryWencai", "wencai", RequestKindObject, "WencaiQueryRequest", "WencaiQueryResult", false),
	newSpec("query_wencai_realtime_fields", "QueryWencaiRealtimeFields", "wencai", RequestKindObject, "WencaiRealtimeFieldsRequest", "WencaiRealtimeFieldsResult", false),
	newSpec("query_wencai_securities", "QueryWencaiSecurities", "wencai", RequestKindObject, "WencaiSecuritiesRequest", "WencaiSecuritiesResult", false),
	newSpec("rank_block_securities", "RankBlockSecurities", "rankings", RequestKindObject, "SortedBlockSecuritiesRequest", "SortedSecuritiesResult", false),
	newSpec("rank_block_securities_by_industry", "RankBlockSecuritiesByIndustry", "rankings", RequestKindObject, "IndustrySortedBlockSecuritiesRequest", "SortedSecuritiesResult", false),
	newSpec("rank_block_securities_by_wencai_field", "RankBlockSecuritiesByWencaiField", "wencai", RequestKindObject, "RankBlockSecuritiesByWencaiFieldRequest", "WencaiBlockRankingResult", false),
	newSpec("rank_constituents_by_performance_contribution", "RankConstituentsByPerformanceContribution", "rankings", RequestKindObject, "PerformanceContributionRankingRequest", "list[PerformanceContributionItem]", false),
	newSpec("rank_related_securities", "RankRelatedSecurities", "rankings", RequestKindObject, "RelatedSecuritiesRequest", "SortedSecuritiesResult", false),
	newSpec("rank_securities_by_popularity", "RankSecuritiesByPopularity", "rankings", RequestKindObject, "PopularityRankTopRequest", "PopularityRankTopResult", false),
	newSpec("rank_securities_by_realtime_statistic", "RankSecuritiesByRealtimeStatistic", "rankings", RequestKindObject, "RealtimeStatsRankingRequest", "RealtimeStatsResult", false),
	newSpec("remove_account_watchlist_group_securities", "RemoveAccountWatchlistGroupSecurities", "account", RequestKindObject, "AccountWatchlistGroupSecuritiesRequest", "AccountWatchlistGroups", true),
	newSpec("remove_account_watchlist_securities", "RemoveAccountWatchlistSecurities", "account", RequestKindObject, "RemoveAccountWatchlistSecuritiesRequest", "AccountWatchlistUpdate", true),
	newSpec("rename_account_watchlist_group", "RenameAccountWatchlistGroup", "account", RequestKindObject, "RenameAccountWatchlistGroupRequest", "AccountWatchlistGroups", true),
	newSpec("replace_account_watchlist_group_securities", "ReplaceAccountWatchlistGroupSecurities", "account", RequestKindObject, "AccountWatchlistGroupSecuritiesRequest", "AccountWatchlistGroups", true),
	newSpec("replace_account_watchlist_securities", "ReplaceAccountWatchlistSecurities", "account", RequestKindObject, "ReplaceAccountWatchlistSecuritiesRequest", "AccountWatchlistUpdate", true),
	newSpec("resolve_block", "ResolveBlock", "securities", RequestKindObject, "ResolveBlockRequest", "Block", false),
	newSpec("resolve_securities", "ResolveSecurities", "securities", RequestKindObject, "ResolveSecuritiesRequest", "list[ResolvedSecurity]", false),
	newSpec("resolve_wencai_securities", "ResolveWencaiSecurities", "wencai", RequestKindObject, "ResolveWencaiSecuritiesRequest", "ResolvedWencaiSecuritiesResult", false),
	newSpec("search_securities", "SearchSecurities", "securities", RequestKindObject, "SearchSecuritiesRequest", "list[SecurityCandidate]", false),
	newSpec("sort_securities", "SortSecurities", "rankings", RequestKindObject, "SortedSecuritiesRequest", "SortedSecuritiesResult", false),
	newSpec("sort_securities_by_industry", "SortSecuritiesByIndustry", "rankings", RequestKindObject, "IndustrySortedSecuritiesRequest", "SortedSecuritiesResult", false),
}

// nativeMethods 是 name → spec 的索引。
var nativeMethods = func() map[string]NativeMethodSpec {
	result := make(map[string]NativeMethodSpec, len(nativeMethodSpecs))
	for _, spec := range nativeMethodSpecs {
		result[spec.Name] = spec
	}
	return result
}()

// nativeMethodNames 是所有 Canonical 方法名的字典序排列。
var nativeMethodNames = func() []string {
	names := make([]string, 0, len(nativeMethodSpecs))
	for _, spec := range nativeMethodSpecs {
		names = append(names, spec.Name)
	}
	sort.Strings(names)
	return names
}()

// mutatingMethodNames 是全部写操作方法名集合。
var mutatingMethodNames = func() map[string]bool {
	result := make(map[string]bool)
	for _, spec := range nativeMethodSpecs {
		if spec.Mutating {
			result[spec.Name] = true
		}
	}
	return result
}()

// tabularMethodNames 是全部只读表格方法名集合（共 63 个）。
var tabularMethodNames = func() map[string]bool {
	result := make(map[string]bool)
	for _, spec := range nativeMethodSpecs {
		if spec.ReturnsTabular {
			result[spec.Name] = true
		}
	}
	return result
}()

// init 复刻 模块加载时的全部一致性断言，任何一条不满足即拒绝启动。
// init 复刻 模块加载时的全部一致性断言，任何一条不满足即拒绝启动。
// 入参：无。
func init() {
	if len(nativeMethods) != 80 {
		panic("THSDK native method manifest must contain exactly 80 unique methods")
	}
	for name := range requestContracts {
		if _, ok := nativeMethods[name]; !ok {
			panic("THSDK request contracts must match the native method manifest")
		}
	}
	for _, spec := range nativeMethodSpecs {
		if len(spec.RequestFields) != len(uniqueStrings(spec.RequestFields)) {
			panic("THSDK request fields must be unique and contain every required field")
		}
		if !containsAll(spec.RequestFields, spec.RequiredFields) {
			panic("THSDK request fields must be unique and contain every required field")
		}
	}
	for _, spec := range nativeMethodSpecs {
		if (spec.RequestKind == RequestKindNone) != (len(spec.RequestFields) == 0) {
			panic("THSDK request fields do not match their request kinds")
		}
		if (spec.RequestKind == RequestKindString || spec.RequestKind == RequestKindSecurities) &&
			!equalStrings(spec.RequestFields, spec.RequiredFields) {
			panic("THSDK request fields do not match their request kinds")
		}
	}
	for _, spec := range nativeMethodSpecs {
		if len(spec.Name) >= 5 && spec.Name[:5] == "list_" && !spec.ReturnsTabular {
			panic("THSDK 表格契约与方法语义不符")
		}
		if spec.Mutating && spec.ReturnsTabular {
			panic("THSDK 表格契约与方法语义不符")
		}
		if spec.TabularRowsField != "" && !spec.ReturnsTabular {
			panic("THSDK 表格契约与方法语义不符")
		}
	}
	if len(tabularMethodNames) != 63 {
		panic("THSDK 必须恰好暴露 63 个只读表格方法")
	}
	for name := range tabularMethodNames {
		if mutatingMethodNames[name] {
			panic("THSDK 必须恰好暴露 63 个只读表格方法")
		}
	}
	// subscription、Level-2、query_data 方法不得进入清单。
	for _, spec := range nativeMethodSpecs {
		lower := spec.Name + " " + spec.GoName
		for _, marker := range []string{"subscribe", "subscription", "level2", "query_data"} {
			if stringsContainsFold(lower, marker) {
				panic("subscription, Level-2, and query_data methods must not enter the THSDK manifest")
			}
		}
	}
}

// MethodInventory 返回完整的 80 项同步 Canonical API 清单。
func MethodInventory() []NativeMethodSpec {
	result := make([]NativeMethodSpec, len(nativeMethodSpecs))
	copy(result, nativeMethodSpecs)
	return result
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]bool, len(values))
	result := make([]string, 0, len(values))
	for _, v := range values {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

func containsAll(values, wanted []string) bool {
	set := make(map[string]bool, len(values))
	for _, v := range values {
		set[v] = true
	}
	for _, w := range wanted {
		if !set[w] {
			return false
		}
	}
	return true
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func stringsContainsFold(value, substr string) bool {
	return strings.Contains(strings.ToLower(value), strings.ToLower(substr))
}
