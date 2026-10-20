// 本文件实现 securities 模块的 18 个 Canonical 方法包装（请求结构见方法清单）。
// 返回值为原生 envelope 结构（any）；表格形态请使用同目录 api_table.go 中的
// <GoName>Table 包装返回表格结果。

package thsdk

// GetMarketMetadata 调用原生方法 GetMarketMetadata，请求类型 MarketMetadataRequest，返回类型 MarketMetadata。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func GetMarketMetadata(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("get_market_metadata", params,
		[]string{"market"},
		nil, opts...)
}

// GetMarketSecurityNames 调用原生方法 GetMarketSecurityNames，请求类型 MarketSecurityNamesRequest，返回类型 MarketSecurityNamesResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func GetMarketSecurityNames(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("get_market_security_names", params,
		[]string{"market"},
		nil, opts...)
}

// GetSecurityConceptTags 调用原生方法 GetSecurityConceptTags，请求类型 SecurityConceptTagsRequest，返回类型 SecurityConceptTagsResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func GetSecurityConceptTags(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("get_security_concept_tags", params,
		[]string{"security"},
		nil, opts...)
}

// GetSecurityIndustry 调用原生方法 GetSecurityIndustry，请求类型 SecurityRequest，返回类型 SecurityIndustryMapping。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func GetSecurityIndustry(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("get_security_industry", params,
		[]string{"security"},
		nil, opts...)
}

// ListBlockConstituents 调用原生 ListBlockConstituents；请求 BlockConstituentsRequest，返回 list[BlockConstituent]。
func ListBlockConstituents(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_block_constituents", params,
		[]string{"block"},
		nil, opts...)
}

// ListBlockDescriptions 调用原生 ListBlockDescriptions；请求 BlockDescriptionsRequest，返回 list[BlockDescription]。
func ListBlockDescriptions(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_block_descriptions", params,
		[]string{"blocks"},
		nil, opts...)
}

// ListFuturesRelatedSecurities 调用原生 ListFuturesRelatedSecurities；请求 None，返回 list[FuturesRelatedSecurity]。
func ListFuturesRelatedSecurities(opts ...*CallOptions) (any, error) {
	return callNone("list_futures_related_securities", opts...)
}

// ListIndustryChildren 调用原生 ListIndustryChildren；请求 IndustryChildrenRequest，返回 list[IndustryChild]。
func ListIndustryChildren(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_industry_children", params,
		[]string{"industry"},
		nil, opts...)
}

// ListMarketSecurities 调用原生方法 ListMarketSecurities，请求类型 MarketSecuritiesRequest，返回类型 MarketSecuritiesResult。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func ListMarketSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_market_securities", params,
		[]string{"market"},
		nil, opts...)
}

// ListRelatedSecurityPerformances 调用原生 ListRelatedSecurityPerformances；请求 RelatedSecurityPerformancesRequest，返回 list[RelatedSecurityPerformance]。
func ListRelatedSecurityPerformances(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_related_security_performances", params,
		[]string{"security"},
		nil, opts...)
}

// ListSecurityAHRelations 调用原生 ListSecurityAHRelations；请求 SecurityAHRelationsRequest，返回 list[SecurityAHRelation]。
func ListSecurityAHRelations(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_ah_relations", params,
		[]string{"security"},
		nil, opts...)
}

// ListSecurityBlockMemberships 调用原生 ListSecurityBlockMemberships；请求 SecurityBlockMembershipsRequest，返回 list[SecurityBlockMembership]。
func ListSecurityBlockMemberships(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_block_memberships", params,
		[]string{"security"},
		nil, opts...)
}

// ListSecurityFuturesRelations 调用原生 ListSecurityFuturesRelations；请求 list[Security]，返回 list[SecurityFuturesRelation]。
func ListSecurityFuturesRelations(securities []any, opts ...*CallOptions) (any, error) {
	return callSecurities("list_security_futures_relations", securities, opts...)
}

// ListSecurityIndustryMappings 调用原生 ListSecurityIndustryMappings；请求 SecurityIndustryMappingsRequest，返回 list[SecurityIndustryMapping]。
func ListSecurityIndustryMappings(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("list_security_industry_mappings", params,
		[]string{"securities"},
		nil, opts...)
}

// ListSecurityLinkRelations 调用原生 ListSecurityLinkRelations；请求 linkKey string，返回 list[SecurityLinkRelation]。
func ListSecurityLinkRelations(link_key string, opts ...*CallOptions) (any, error) {
	return callString("list_security_link_relations", link_key, "link_key", opts...)
}

// ResolveBlock 调用原生方法 ResolveBlock，请求类型 ResolveBlockRequest，返回类型 Block。
// 入参：params 为请求参数映射（必填字段见方法清单）；opts 为调用选项（超时）。
// 出参：any 为原生返回的 data；error 为调用失败时的错误。
func ResolveBlock(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("resolve_block", params,
		[]string{"name"},
		nil, opts...)
}

// ResolveSecurities 调用原生 ResolveSecurities；请求 ResolveSecuritiesRequest，返回 list[ResolvedSecurity]。
func ResolveSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("resolve_securities", params,
		[]string{"codes"},
		nil, opts...)
}

// SearchSecurities 调用原生 SearchSecurities；请求 SearchSecuritiesRequest，返回 list[SecurityCandidate]。
func SearchSecurities(params map[string]any, opts ...*CallOptions) (any, error) {
	return callTyped("search_securities", params,
		[]string{"pattern"},
		nil, opts...)
}
