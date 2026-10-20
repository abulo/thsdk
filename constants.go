package thsdk

// 本文件集中定义协议层全部枚举常量（整数枚举与字符串枚举）。
//
// Go 中枚举以带类型的常量表达，分为两类：
//   - 整数枚举（Adjust / Period / InfoType 等）→ 整型常量；
//   - 字符串枚举（CallAuctionPhase / NewsTraversalMode 等）→ 字符串常量。
//
// 所有常量值与原生协议逐字一致，可直接作为请求字段传入，不可改动。

// Adjust 复权方式。
type Adjust int

const (
	AdjustNone     Adjust = 0 // 不复权
	AdjustForward  Adjust = 1 // 前复权
	AdjustBackward Adjust = 2 // 后复权
)

// Period K 线周期。
type Period int

const (
	PeriodMinute1   Period = 0x3001 // 1 分钟
	PeriodMinute5   Period = 0x3005 // 5 分钟
	PeriodMinute15  Period = 0x300F // 15 分钟
	PeriodMinute30  Period = 0x301E // 30 分钟
	PeriodMinute60  Period = 0x303C // 60 分钟
	PeriodMinute120 Period = 0x3078 // 120 分钟
	PeriodDay       Period = 0x4000 // 日线
	PeriodWeek      Period = 0x5001 // 周线
	PeriodMonth     Period = 0x6001 // 月线
	PeriodQuarter   Period = 0x6003 // 季线
	PeriodYear      Period = 0x7001 // 年线
)

// CallAuctionPhase 集合竞价阶段。
type CallAuctionPhase string

const (
	CallAuctionPhaseOpening CallAuctionPhase = "opening"
	CallAuctionPhaseClosing CallAuctionPhase = "closing"
	CallAuctionPhaseCustom  CallAuctionPhase = "custom"
)

// IndexCallAuctionPhase 指数集合竞价阶段。
type IndexCallAuctionPhase string

const (
	IndexCallAuctionPhaseOpening IndexCallAuctionPhase = "opening"
	IndexCallAuctionPhaseClosing IndexCallAuctionPhase = "closing"
)

// ExtendedHoursSession 盘前盘后时段。
type ExtendedHoursSession string

const (
	ExtendedHoursSessionPreMarket  ExtendedHoursSession = "pre_market"
	ExtendedHoursSessionPostMarket ExtendedHoursSession = "post_market"
)

// OrderBookDepth 盘口档位。
type OrderBookDepth string

const (
	OrderBookDepthFive OrderBookDepth = "5"
)

// NewsFlashTraversalMode 快讯遍历模式。
type NewsFlashTraversalMode string

const (
	NewsFlashTraversalModePage  NewsFlashTraversalMode = "page"
	NewsFlashTraversalModeOlder NewsFlashTraversalMode = "older"
	NewsFlashTraversalModeNewer NewsFlashTraversalMode = "newer"
)

// NewsTimeline 资讯时间线。
type NewsTimeline string

const (
	NewsTimelineHistorical NewsTimeline = "historical"
	NewsTimelineRealtime   NewsTimeline = "realtime"
)

// NewsTraversalMode 资讯遍历模式。
type NewsTraversalMode string

const (
	NewsTraversalModePage  NewsTraversalMode = "page"
	NewsTraversalModeSince NewsTraversalMode = "since"
)

// NewsDirection 资讯方向。
type NewsDirection string

const (
	NewsDirectionFirst NewsDirection = "first"
	NewsDirectionOlder NewsDirection = "older"
	NewsDirectionNewer NewsDirection = "newer"
)

// InfoType 资讯类型。
type InfoType int

const (
	InfoTypeBad                    InfoType = -1
	InfoTypeHistoryTime            InfoType = 2053
	InfoTypeRealtime               InfoType = 2054
	InfoTypeNewsReport             InfoType = 8193
	InfoTypeShSzFund               InfoType = 10240
	InfoTypeForeignExchangeNews    InfoType = 10241
	InfoTypeForeignExchangeFx      InfoType = 10242
	InfoTypeForeignExchangeCj      InfoType = 10243
	InfoTypeGlobalMarket           InfoType = 10244
	InfoTypeGlobalFinance          InfoType = 10245
	InfoTypeGlobalEconomicData     InfoType = 10246
	InfoTypeStockInfo              InfoType = 14339
	InfoTypeIndustryInfo           InfoType = 14340
	InfoTypeAnnouncement           InfoType = 14341
	InfoTypeStockResearchReport    InfoType = 14342
	InfoTypeIndustryResearchReport InfoType = 14343
	InfoTypeSecurityTapeReading    InfoType = 14348
	InfoTypeIndexFutures           InfoType = 14349
	InfoTypeMarketTapeReading      InfoType = 14355
	InfoTypeMarketNews             InfoType = 14356
	InfoTypeMarketTrendReading     InfoType = 14359
	InfoTypeGltMessage             InfoType = 14363
	InfoTypeNewStockInfo           InfoType = 14364
	InfoTypeMultiplyFutures        InfoType = 14365
	InfoTypeNeeqNews               InfoType = 14366
	InfoTypeNeeqSecondaryNews      InfoType = 14367
	InfoTypeCommodityIndex         InfoType = 14368
	InfoTypeSecurityInsight        InfoType = 14371
	InfoTypeShToHkNews             InfoType = 32769
	InfoTypeHkNews                 InfoType = 32770
	InfoTypeUsaNews                InfoType = 32770
	InfoTypeHkPanInfo              InfoType = 32771
)

// InfoSummaryMode 资讯摘要模式。
type InfoSummaryMode int

const (
	InfoSummaryModeDefault  InfoSummaryMode = 0
	InfoSummaryModeDisabled InfoSummaryMode = 1
	InfoSummaryModeEnabled  InfoSummaryMode = 2
)

// RealtimeStatsMetric 实时统计指标。
type RealtimeStatsMetric int

const (
	RealtimeStatsMetricFirstLimitUpTime RealtimeStatsMetric = 330323
	RealtimeStatsMetricFinalLimitUpTime RealtimeStatsMetric = 330324
	RealtimeStatsMetricLimitUpType      RealtimeStatsMetric = 330325
	RealtimeStatsMetricLimitUpState     RealtimeStatsMetric = 330329
)

// RealtimeSecuritySignal 实时证券信号。
type RealtimeSecuritySignal string

const (
	RealtimeSecuritySignalLimitUp        RealtimeSecuritySignal = "limit_up"
	RealtimeSecuritySignalLimitDown      RealtimeSecuritySignal = "limit_down"
	RealtimeSecuritySignalLimitUpBreak   RealtimeSecuritySignal = "limit_up_break"
	RealtimeSecuritySignalLimitDownBreak RealtimeSecuritySignal = "limit_down_break"
)

// FinancialSnapshotField 财务快照字段。
type FinancialSnapshotField string

const (
	FinancialSnapshotFieldTotalShares      FinancialSnapshotField = "total_shares"
	FinancialSnapshotFieldFloatShares      FinancialSnapshotField = "float_shares"
	FinancialSnapshotFieldNetAssetPerShare FinancialSnapshotField = "net_asset_per_share"
)

// OptionGreeksMetric 期权 Greeks 指标。
type OptionGreeksMetric int

const (
	OptionGreeksMetricImpliedVolatility OptionGreeksMetric = 920341
	OptionGreeksMetricHistoryVolatility OptionGreeksMetric = 723853
	OptionGreeksMetricRealLeverRatio    OptionGreeksMetric = 2980
	OptionGreeksMetricTheoryPrice       OptionGreeksMetric = 920342
	OptionGreeksMetricDelta             OptionGreeksMetric = 723854
	OptionGreeksMetricGamma             OptionGreeksMetric = 723855
	OptionGreeksMetricRho               OptionGreeksMetric = 723856
	OptionGreeksMetricTheta             OptionGreeksMetric = 723857
	OptionGreeksMetricVega              OptionGreeksMetric = 723858
)

// WencaiHotBlockKind 问财热门板块类型。
type WencaiHotBlockKind string

const (
	WencaiHotBlockKindConcept  WencaiHotBlockKind = "concept"
	WencaiHotBlockKindIndustry WencaiHotBlockKind = "industry"
)

// HotEventNewsRange 热点事件资讯范围。
type HotEventNewsRange string

const (
	HotEventNewsRangeRecent     HotEventNewsRange = "recent"
	HotEventNewsRangeHistorical HotEventNewsRange = "historical"
)

// ShortTermEventMode 短线事件模式。
type ShortTermEventMode string

const (
	ShortTermEventModeMarketRecent   ShortTermEventMode = "market_recent"
	ShortTermEventModeSecurityScoped ShortTermEventMode = "security_scoped"
)

// PerformanceContributionTarget 业绩贡献目标。
type PerformanceContributionTarget string

const (
	PerformanceContributionTargetIndex    PerformanceContributionTarget = "index"
	PerformanceContributionTargetSecurity PerformanceContributionTarget = "security"
)

// PopularityRankTopType 人气榜类型。
type PopularityRankTopType int

const (
	PopularityRankTopTypeAllStocks PopularityRankTopType = 1
)
