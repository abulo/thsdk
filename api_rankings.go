// Code generated from the Python THSDK 2.0.2 source (thsdk/_apis/rankings.py). DO NOT EDIT BY HAND.
// 本文件对应 Python 包的 thsdk/_apis/rankings.py：11 个 Canonical 方法的一一对应包装。
// 返回值为原生 envelope 结构（any）；表格形态请使用同目录 api_dataframe.go 中的
// <GoName>Table 包装（对应 Python 公开层的 pandas.DataFrame 结果）。

package thsdk

// CalculateSecurityRealtimeStatistics 调用原生 CalculateSecurityRealtimeStatistics；请求 RealtimeStatsCalculationRequest，返回 RealtimeStatsResult。
// 对应 Python thsdk._apis.rankings.calculate_security_realtime_statistics。
func CalculateSecurityRealtimeStatistics(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("calculate_security_realtime_statistics", params,
		[]string{"securities", "metrics"},
		nil, opts...)
}

// GetSecurityPopularityRank 调用原生 GetSecurityPopularityRank；请求 PopularityRankRequest，返回 PopularityRankResult。
// 对应 Python thsdk._apis.rankings.get_security_popularity_rank。
func GetSecurityPopularityRank(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("get_security_popularity_rank", params,
		[]string{"security"},
		nil, opts...)
}

// ListSecuritiesByRealtimeSignal 调用原生 ListSecuritiesByRealtimeSignal；请求 RealtimeSignalSecuritiesRequest，返回 list[RealtimeSignalSecurity]。
// 对应 Python thsdk._apis.rankings.list_securities_by_realtime_signal。
func ListSecuritiesByRealtimeSignal(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_securities_by_realtime_signal", params,
		[]string{"markets", "signal"},
		nil, opts...)
}

// RankBlockSecurities 调用原生 RankBlockSecurities；请求 SortedBlockSecuritiesRequest，返回 SortedSecuritiesResult。
// 对应 Python thsdk._apis.rankings.rank_block_securities。
func RankBlockSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("rank_block_securities", params,
		[]string{"block"},
		nil, opts...)
}

// RankBlockSecuritiesByIndustry 调用原生 RankBlockSecuritiesByIndustry；请求 IndustrySortedBlockSecuritiesRequest，返回 SortedSecuritiesResult。
// 对应 Python thsdk._apis.rankings.rank_block_securities_by_industry。
func RankBlockSecuritiesByIndustry(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("rank_block_securities_by_industry", params,
		[]string{"block"},
		nil, opts...)
}

// RankConstituentsByPerformanceContribution 调用原生 RankConstituentsByPerformanceContribution；请求 PerformanceContributionRankingRequest，返回 list[PerformanceContributionItem]。
// 对应 Python thsdk._apis.rankings.rank_constituents_by_performance_contribution。
func RankConstituentsByPerformanceContribution(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("rank_constituents_by_performance_contribution", params,
		[]string{"code", "market"},
		nil, opts...)
}

// RankRelatedSecurities 调用原生 RankRelatedSecurities；请求 RelatedSecuritiesRequest，返回 SortedSecuritiesResult。
// 对应 Python thsdk._apis.rankings.rank_related_securities。
func RankRelatedSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("rank_related_securities", params,
		[]string{"security"},
		nil, opts...)
}

// RankSecuritiesByPopularity 调用原生 RankSecuritiesByPopularity；请求 PopularityRankTopRequest，返回 PopularityRankTopResult。
// 对应 Python thsdk._apis.rankings.rank_securities_by_popularity。
func RankSecuritiesByPopularity(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("rank_securities_by_popularity", params,
		nil,
		nil, opts...)
}

// RankSecuritiesByRealtimeStatistic 调用原生 RankSecuritiesByRealtimeStatistic；请求 RealtimeStatsRankingRequest，返回 RealtimeStatsResult。
// 对应 Python thsdk._apis.rankings.rank_securities_by_realtime_statistic。
func RankSecuritiesByRealtimeStatistic(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("rank_securities_by_realtime_statistic", params,
		[]string{"sort_by"},
		nil, opts...)
}

// SortSecurities 调用原生 SortSecurities；请求 SortedSecuritiesRequest，返回 SortedSecuritiesResult。
// 对应 Python thsdk._apis.rankings.sort_securities。
func SortSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("sort_securities", params,
		[]string{"securities"},
		nil, opts...)
}

// SortSecuritiesByIndustry 调用原生 SortSecuritiesByIndustry；请求 IndustrySortedSecuritiesRequest，返回 SortedSecuritiesResult。
// 对应 Python thsdk._apis.rankings.sort_securities_by_industry。
func SortSecuritiesByIndustry(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("sort_securities_by_industry", params,
		[]string{"securities"},
		nil, opts...)
}
