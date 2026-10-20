// 本文件实现 rankings 模块的 11 个 Canonical 方法包装（请求结构见方法清单）。
// 返回值为原生 envelope 结构（any）；表格形态请使用同目录 api_table.go 中的
// <GoName>Table 包装返回表格结果。

package thsdk

// CalculateSecurityRealtimeStatistics 调用原生方法 CalculateSecurityRealtimeStatistics，请求类型 RealtimeStatsCalculationRequest，返回类型 RealtimeStatsResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func CalculateSecurityRealtimeStatistics(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("calculate_security_realtime_statistics", params,
		[]string{"securities", "metrics"},
		nil, opts...)
}

// GetSecurityPopularityRank 调用原生方法 GetSecurityPopularityRank，请求类型 PopularityRankRequest，返回类型 PopularityRankResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func GetSecurityPopularityRank(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("get_security_popularity_rank", params,
		[]string{"security"},
		nil, opts...)
}

// ListSecuritiesByRealtimeSignal 调用原生 ListSecuritiesByRealtimeSignal；请求 RealtimeSignalSecuritiesRequest，返回 list[RealtimeSignalSecurity]。
func ListSecuritiesByRealtimeSignal(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_securities_by_realtime_signal", params,
		[]string{"markets", "signal"},
		nil, opts...)
}

// RankBlockSecurities 调用原生方法 RankBlockSecurities，请求类型 SortedBlockSecuritiesRequest，返回类型 SortedSecuritiesResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func RankBlockSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("rank_block_securities", params,
		[]string{"block"},
		nil, opts...)
}

// RankBlockSecuritiesByIndustry 调用原生方法 RankBlockSecuritiesByIndustry，请求类型 IndustrySortedBlockSecuritiesRequest，返回类型 SortedSecuritiesResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func RankBlockSecuritiesByIndustry(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("rank_block_securities_by_industry", params,
		[]string{"block"},
		nil, opts...)
}

// RankConstituentsByPerformanceContribution 调用原生 RankConstituentsByPerformanceContribution；请求 PerformanceContributionRankingRequest，返回 list[PerformanceContributionItem]。
func RankConstituentsByPerformanceContribution(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("rank_constituents_by_performance_contribution", params,
		[]string{"code", "market"},
		nil, opts...)
}

// RankRelatedSecurities 调用原生方法 RankRelatedSecurities，请求类型 RelatedSecuritiesRequest，返回类型 SortedSecuritiesResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func RankRelatedSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("rank_related_securities", params,
		[]string{"security"},
		nil, opts...)
}

// RankSecuritiesByPopularity 调用原生方法 RankSecuritiesByPopularity，请求类型 PopularityRankTopRequest，返回类型 PopularityRankTopResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func RankSecuritiesByPopularity(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("rank_securities_by_popularity", params,
		nil,
		nil, opts...)
}

// RankSecuritiesByRealtimeStatistic 调用原生方法 RankSecuritiesByRealtimeStatistic，请求类型 RealtimeStatsRankingRequest，返回类型 RealtimeStatsResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func RankSecuritiesByRealtimeStatistic(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("rank_securities_by_realtime_statistic", params,
		[]string{"sort_by"},
		nil, opts...)
}

// SortSecurities 调用原生方法 SortSecurities，请求类型 SortedSecuritiesRequest，返回类型 SortedSecuritiesResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func SortSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("sort_securities", params,
		[]string{"securities"},
		nil, opts...)
}

// SortSecuritiesByIndustry 调用原生方法 SortSecuritiesByIndustry，请求类型 IndustrySortedSecuritiesRequest，返回类型 SortedSecuritiesResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func SortSecuritiesByIndustry(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("sort_securities_by_industry", params,
		[]string{"securities"},
		nil, opts...)
}
