package thsdk

// 本文件实现公开便捷函数层：
// get_price / klines / intraday_data / tick_level1 / depth / corporate_action /
// search_symbols / complete_ths_code / wencai_nlp / news / market_securities /
// block_constituents，以及公开认证入口 auth / auth_qrcode / logout。
//
// 设计说明：关键字参数在 Go 中映射为「选项结构体」；零值表示
// 未指定并采用默认值。所有表格结果返回 *Table
// （列/行结构的表格结果）。

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// periods 是周期中文名 → 协议 Period 的映射表。
var periods = map[string]int64{
	"1m": 0x3001, "5m": 0x3005, "15m": 0x300F, "30m": 0x301E, "60m": 0x303C, "120m": 0x3078,
	"day": 0x4000, "daily": 0x4000, "1d": 0x4000,
	"week": 0x5001, "weekly": 0x5001,
	"month": 0x6001, "monthly": 0x6001,
	"quarter": 0x6003,
	"year":    0x7001,
}

// adjustments 是复权中文名 → 协议 Adjust 的映射表（不含"不复权"）。
var adjustments = map[string]int64{
	"": 0, "none": 0, "pre": 1, "forward": 1, "post": 2, "backward": 2,
}

// Auth 按 显式凭据 → 环境变量 → 已有会话 → 临时账户 的顺序完成认证。
// 入参：username、password、mac、opts；语义见函数用途。
// 出参：bool（成功时）；error 为调用失败时的错误。
func Auth(username, password, mac string, opts ...*AuthOptions) (bool, error) {
	return getClient().Auth(username, password, mac, opts...)
}

// AuthQRCode 优先恢复已有会话；需要时直接显示二维码并等待扫码。
// 入参：opts；语义见函数用途。
// 出参：bool（成功时）；error 为调用失败时的错误。
func AuthQRCode(opts ...*AuthQRCodeOptions) (bool, error) {
	return getClient().AuthQRCode(opts...)
}

// Logout 退出当前账号并删除当前目录中可复用的会话文件。
// 入参：；语义见函数用途。
// 出参：error 为调用失败时的错误。
func Logout() error {
	return getClient().Logout()
}

// Request 在显式认证后执行一次包定义的数据请求。
// 入参：method、params、opts；语义见函数用途。
// 出参：any（成功时）；error 为调用失败时的错误。
func Request(method string, params any, opts ...*CallOptions) (any, error) {
	return getClient().Request(method, params, opts...)
}

// security 规范化证券输入（完整代码字符串或 code/market 映射）。
func security(value any) (map[string]any, error) {
	if mapping, ok := value.(map[string]any); ok {
		code := strings.TrimSpace(fmt.Sprint(mapping["code"]))
		market := strings.ToUpper(strings.TrimSpace(fmt.Sprint(mapping["market"])))
		if code == "" || len(market) != 4 || !isAlpha(market) {
			return nil, fmt.Errorf("security 字典必须包含 code 和四位 market")
		}
		return map[string]any{"code": code, "market": market}, nil
	}
	text, ok := value.(string)
	if !ok || strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("证券代码必须是非空字符串")
	}
	return normalizeSecurity(strings.TrimSpace(text))
}

// securities 把单个或序列证券输入统一为列表。
func securities(value any) ([]any, error) {
	if _, isString := value.(string); isString {
		item, err := security(value)
		if err != nil {
			return nil, err
		}
		return []any{item}, nil
	}
	if _, isMap := value.(map[string]any); isMap {
		item, err := security(value)
		if err != nil {
			return nil, err
		}
		return []any{item}, nil
	}
	sequence, ok := toSequenceAny(value)
	if !ok || len(sequence) == 0 {
		return nil, fmt.Errorf("证券代码不能为空")
	}
	result := make([]any, 0, len(sequence))
	for _, item := range sequence {
		converted, err := security(item)
		if err != nil {
			return nil, err
		}
		result = append(result, converted)
	}
	return result, nil
}

// marketValue 校验四位字母市场代码。
func marketValue(value string) (string, error) {
	market := strings.ToUpper(strings.TrimSpace(value))
	if len(market) != 4 || !isAlpha(market) {
		return "", fmt.Errorf("市场代码必须是四位字母，例如 USHA 或 USZA")
	}
	return market, nil
}

// marketsValue 规范化市场输入（字符串或序列），空结果返回 nil。
func marketsValue(value any) ([]string, error) {
	var values []string
	if text, isString := value.(string); isString {
		for _, part := range strings.FieldsFunc(text, func(r rune) bool { return r == ',' || r == ' ' }) {
			if part != "" {
				values = append(values, part)
			}
		}
	} else if sequence, ok := toSequenceAny(value); ok {
		for _, item := range sequence {
			values = append(values, fmt.Sprint(item))
		}
	} else {
		return nil, nil
	}
	if len(values) == 0 {
		return nil, nil
	}
	result := make([]string, 0, len(values))
	for _, item := range values {
		converted, err := marketValue(item)
		if err != nil {
			return nil, err
		}
		result = append(result, converted)
	}
	return result, nil
}

// dateNumber 规范化日期为整数：返回 (数值, 是否提供, 错误)。
func dateNumber(value any) (int64, bool, error) {
	switch typed := value.(type) {
	case nil:
		return 0, false, nil
	case time.Time:
		return int64(typed.Year()*10000 + int(typed.Month())*100 + typed.Day()), true, nil
	case int:
		text := strconv.Itoa(typed)
		if len(text) != 8 || !allDigits(text) {
			return 0, true, fmt.Errorf("日期必须为 YYYY-MM-DD、YYYYMMDD、date 或 datetime")
		}
		return int64(typed), true, nil
	case int64:
		text := strconv.FormatInt(typed, 10)
		if len(text) != 8 || !allDigits(text) {
			return 0, true, fmt.Errorf("日期必须为 YYYY-MM-DD、YYYYMMDD、date 或 datetime")
		}
		return typed, true, nil
	default:
		text := strings.TrimSpace(fmt.Sprint(value))
		text = strings.ReplaceAll(text, "-", "")
		if len(text) != 8 || !allDigits(text) {
			return 0, true, fmt.Errorf("日期必须为 YYYY-MM-DD、YYYYMMDD、date 或 datetime")
		}
		number, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return 0, true, fmt.Errorf("日期必须为 YYYY-MM-DD、YYYYMMDD、date 或 datetime")
		}
		return number, true, nil
	}
}

// dateText 规范化日期为文本，未提供时返回 "0"。
func dateText(value any) string {
	number, provided, err := dateNumber(value)
	if err != nil || !provided {
		return "0"
	}
	return strconv.FormatInt(number, 10)
}

// datetimeText 规范化日期时间为 ISO 8601 文本，无时区按上海时区。
func datetimeText(value any) (string, error) {
	var parsed time.Time
	explicitZone := false
	switch typed := value.(type) {
	case time.Time:
		parsed = typed
	default:
		text := strings.TrimSpace(fmt.Sprint(value))
		text = strings.ReplaceAll(text, "Z", "+00:00")
		// 解析 ISO 8601：显式时区（±hh:mm）保持原偏移；无时区输入按上海时区。
		// 注意日期分隔符 "-" 不算时区，只检查时间部分之后的符号。
		explicitZone = strings.Contains(text, "+")
		if !explicitZone {
			if marker := strings.IndexAny(text, "Tt "); marker >= 0 {
				explicitZone = strings.Contains(text[marker:], "-")
			}
		}
		parsed = parseISODateTime(text)
		if parsed.IsZero() {
			return "", fmt.Errorf("时间必须是 datetime、date 或 ISO 8601 字符串")
		}
	}
	if !explicitZone && parsed.Location() == time.UTC {
		parsed = assignShanghai(parsed)
	}
	return parsed.Format(isoFormat), nil
}

// parseISODateTime 尝试多种 ISO 8601 布局解析时间。
func parseISODateTime(text string) time.Time {
	for _, layout := range []string{isoFormat, time.RFC3339Nano, "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02"} {
		if parsed, err := time.Parse(layout, text); err == nil {
			return parsed
		}
	}
	return time.Time{}
}

// assignShanghai 把无时区时间按墙钟时间套用上海时区（不改变墙上时钟读数）。
func assignShanghai(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(),
		value.Hour(), value.Minute(), value.Second(), value.Nanosecond(), shanghai)
}

func allDigits(text string) bool {
	if text == "" {
		return false
	}
	for _, r := range text {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// GetPriceOptions 控制 GetPrice 的查询参数。
type GetPriceOptions struct {
	StartDate any      // 日期字符串 / time.Time / int（YYYYMMDD）；nil 表示不指定
	EndDate   any      // 同上
	Frequency string   // 默认 "daily"
	Fields    []string // 返回字段过滤；nil 表示全部
	FQ        string   // 默认 "pre"；"none" 表示不复权
	Count     int      // 0 表示不指定；与 StartDate/EndDate 互斥
}

// GetPrice 查询一只证券的 K 线并返回表格结果（*Table）。
// 入参：securityValue、opts；语义见函数用途。
// 出参：*Table（成功时）；error 为调用失败时的错误。
func GetPrice(securityValue any, opts ...*GetPriceOptions) (*Table, error) {
	option := GetPriceOptions{}
	for _, item := range opts {
		if item != nil {
			option = *item
		}
	}
	frequency := option.Frequency
	if frequency == "" {
		frequency = "daily"
	}
	fq := option.FQ
	if fq == "" {
		fq = "pre"
	}
	period, ok := periods[strings.ToLower(frequency)]
	if !ok {
		return nil, fmt.Errorf("不支持的 frequency: %s", frequency)
	}
	adjustment, ok := adjustments[strings.ToLower(fq)]
	if !ok {
		return nil, fmt.Errorf("不支持的 fq: %s", fq)
	}
	count := option.Count
	hasStart := option.StartDate != nil
	hasEnd := option.EndDate != nil
	if count != 0 && (hasStart || hasEnd) {
		return nil, fmt.Errorf("count 不能与 start_date/end_date 同时使用")
	}
	if hasStart != hasEnd {
		return nil, fmt.Errorf("start_date 和 end_date 必须同时提供")
	}

	var start, end int64
	if count != 0 {
		if count <= 0 {
			return nil, fmt.Errorf("count 必须大于 0")
		}
		start, end = -int64(count), 0
	} else if hasStart {
		if strings.HasSuffix(strings.ToLower(frequency), "m") {
			return nil, fmt.Errorf("分钟 K 线的日期范围暂不支持；请改用 count")
		}
		var err error
		if start, _, err = dateNumber(option.StartDate); err != nil {
			return nil, err
		}
		if end, _, err = dateNumber(option.EndDate); err != nil {
			return nil, err
		}
	} else {
		start, end = -1, 0
	}

	result, err := ListSecurityKLines(map[string]any{
		"security": securityValue,
		"start":    start,
		"end":      end,
		"adjust":   adjustment,
		"period":   period,
	})
	if err != nil {
		return nil, err
	}
	rows := result
	if count != 0 {
		if list, ok := rows.([]any); ok {
			n := count
			if n > len(list) {
				n = len(list)
			}
			rows = list[len(list)-n:]
		}
	}
	table, err := frame(rows, true)
	if err != nil {
		return nil, err
	}
	if len(option.Fields) > 0 {
		return table.Select(option.Fields)
	}
	return table, nil
}

// KlinesOptions 控制 Klines 的查询参数。
type KlinesOptions struct {
	StartTime any    // 日期字符串 / time.Time / int
	EndTime   any    // 同上
	Adjust    string // 默认 ""（不复权）
	Interval  string // 默认 "day"
	Count     int    // 0 表示未指定（与 start/end 均为空时默认 100）
}

// Klines 使用传统参数名查询 K 线；返回 *Table。
// 入参：thsCode、opts；语义见函数用途。
// 出参：*Table（成功时）；error 为调用失败时的错误。
func Klines(thsCode any, opts ...*KlinesOptions) (*Table, error) {
	option := KlinesOptions{}
	for _, item := range opts {
		if item != nil {
			option = *item
		}
	}
	if option.Count == 0 && option.StartTime == nil && option.EndTime == nil {
		option.Count = 100
	}
	interval := option.Interval
	if interval == "" {
		interval = "day"
	}
	fq := option.Adjust
	if fq == "" {
		fq = "none"
	}
	return GetPrice(thsCode, &GetPriceOptions{
		StartDate: option.StartTime,
		EndDate:   option.EndTime,
		Frequency: interval,
		FQ:        fq,
		Count:     option.Count,
	})
}

// IntradayData 查询指定交易日的分时走势；date 为 nil 时查询当天。
// 入参：thsCode、dateValue；语义见函数用途。
// 出参：*Table（成功时）；error 为调用失败时的错误。
func IntradayData(thsCode any, dateValue any) (*Table, error) {
	result, err := ListSecurityIntradayBars(map[string]any{
		"security": thsCode,
		"date":     dateText(dateValue),
	})
	if err != nil {
		return nil, err
	}
	return frame(result, false)
}

// TickLevel1Options 控制 TickLevel1 的查询参数。
type TickLevel1Options struct {
	Count     int // 默认 100；与时间范围互斥
	StartTime any // 与 EndTime 同时提供
	EndTime   any // 与 StartTime 同时提供
}

// TickLevel1 查询普通逐笔成交，可按最近条数或明确时间范围读取。
// 入参：thsCode、opts；语义见函数用途。
// 出参：*Table（成功时）；error 为调用失败时的错误。
func TickLevel1(thsCode any, opts ...*TickLevel1Options) (*Table, error) {
	option := TickLevel1Options{}
	for _, item := range opts {
		if item != nil {
			option = *item
		}
	}
	count := option.Count
	hasStart := option.StartTime != nil
	hasEnd := option.EndTime != nil
	if count != 0 && (hasStart || hasEnd) {
		return nil, fmt.Errorf("count 不能与 start_time/end_time 同时使用")
	}
	if hasStart != hasEnd {
		return nil, fmt.Errorf("start_time 和 end_time 必须同时提供")
	}
	var window map[string]any
	if count != 0 {
		if count <= 0 {
			return nil, fmt.Errorf("count 必须大于 0")
		}
		window = map[string]any{"count": count}
	} else if hasStart && hasEnd {
		start, err := datetimeText(option.StartTime)
		if err != nil {
			return nil, err
		}
		end, err := datetimeText(option.EndTime)
		if err != nil {
			return nil, err
		}
		window = map[string]any{"start": start, "end": end}
	} else {
		return nil, fmt.Errorf("请提供 count，或同时提供 start_time 和 end_time")
	}
	result, err := ListSecurityTicks(map[string]any{
		"security": thsCode,
		"window":   window,
	})
	if err != nil {
		return nil, err
	}
	rows := result
	if count != 0 {
		if list, ok := rows.([]any); ok {
			n := count
			if n > len(list) {
				n = len(list)
			}
			rows = list[len(list)-n:]
		}
	}
	return frame(rows, false)
}

// Depth 查询一只或多只证券的五档盘口，按买卖方向和档位展开为表格。
// 入参：thsCode；语义见函数用途。
// 出参：*Table（成功时）；error 为调用失败时的错误。
func Depth(thsCode any) (*Table, error) {
	inputs, err := securities(thsCode)
	if err != nil {
		return nil, err
	}
	result, err := ListSecurityOrderBooks(map[string]any{
		"securities": inputs,
		"depth":      "5",
	})
	if err != nil {
		return nil, err
	}
	var rows []map[string]any
	for _, item := range toSequenceAnyOrEmpty(result) {
		book, ok := item.(map[string]any)
		if !ok {
			continue
		}
		securityValue, _ := book["security"].(map[string]any)
		fullCode := fmt.Sprint(securityValue["market"]) + fmt.Sprint(securityValue["code"])
		sides := []struct {
			name string
			key  string
		}{{"bid", "bids"}, {"ask", "asks"}}
		for _, side := range sides {
			quotes, _ := toSequenceAny(book[side.key])
			for level, quoteValue := range quotes {
				quote, ok := quoteValue.(map[string]any)
				if !ok {
					continue
				}
				rows = append(rows, map[string]any{
					"full_code": fullCode,
					"side":      side.name,
					"level":     level + 1,
					"price":     quote["price"],
					"volume":    quote["volume"],
				})
			}
		}
	}
	return frame(rows, false)
}

// CorporateActionOptions 控制 CorporateAction 的查询参数。
type CorporateActionOptions struct {
	Count int // 默认 100
}

// CorporateAction 查询最近的分红、送转、配股和除权除息记录。
// 入参：thsCode、opts；语义见函数用途。
// 出参：*Table（成功时）；error 为调用失败时的错误。
func CorporateAction(thsCode any, opts ...*CorporateActionOptions) (*Table, error) {
	count := 100
	for _, item := range opts {
		if item == nil {
			continue
		}
		count = item.Count
	}
	if count <= 0 {
		return nil, fmt.Errorf("count 必须大于 0")
	}
	result, err := ListSecurityCorporateActions(map[string]any{
		"security": thsCode,
		"start":    -int64(count),
		"end":      0,
		"adjust":   0,
		"period":   periods["day"],
	})
	if err != nil {
		return nil, err
	}
	rows := result
	if list, ok := rows.([]any); ok {
		n := count
		if n > len(list) {
			n = len(list)
		}
		rows = list[len(list)-n:]
	}
	return frame(rows, false)
}

// SearchSymbols 按代码、拼音或中文名称搜索证券。
// 入参：pattern、market；语义见函数用途。
// 出参：*Table（成功时）；error 为调用失败时的错误。
func SearchSymbols(pattern, market string) (*Table, error) {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return nil, fmt.Errorf("pattern 不能为空")
	}
	params := map[string]any{"pattern": pattern}
	if market != "" {
		converted, err := marketValue(market)
		if err != nil {
			return nil, err
		}
		params["market"] = converted
	}
	result, err := SearchSecurities(params)
	if err != nil {
		return nil, err
	}
	return frame(identityRows(result), false)
}

// CompleteTHSCode 把一个或多个短代码解析为可用于行情查询的完整证券代码。
// 入参：thsCode；语义见函数用途。
// 出参：*Table（成功时）；error 为调用失败时的错误。
func CompleteTHSCode(thsCode any) (*Table, error) {
	var codes []string
	if text, ok := thsCode.(string); ok {
		codes = []string{text}
	} else if sequence, ok := toSequenceAny(thsCode); ok {
		for _, item := range sequence {
			codes = append(codes, fmt.Sprint(item))
		}
	}
	var cleaned []string
	for _, code := range codes {
		text := strings.ToUpper(strings.TrimSpace(code))
		if text != "" {
			cleaned = append(cleaned, text)
		}
	}
	if len(cleaned) == 0 {
		return nil, fmt.Errorf("ths_code 不能为空")
	}
	result, err := ResolveSecurities(map[string]any{"codes": cleaned})
	if err != nil {
		return nil, err
	}
	return frame(identityRows(result), false)
}

// WencaiNLPOptions 控制 WencaiNLP 的查询参数。
type WencaiNLPOptions struct {
	Markets any // 字符串（逗号/空格分隔）或市场序列；nil 表示全部
	Limit   int // 默认 100
}

// WencaiNLP 执行自然语言选股，并把可用动态字段展开为表格列。
// 入参：condition、opts；语义见函数用途。
// 出参：*Table（成功时）；error 为调用失败时的错误。
func WencaiNLP(condition string, opts ...*WencaiNLPOptions) (*Table, error) {
	option := WencaiNLPOptions{}
	for _, item := range opts {
		if item != nil {
			option = *item
		}
	}
	condition = strings.TrimSpace(condition)
	if condition == "" {
		return nil, fmt.Errorf("condition 不能为空")
	}
	limit := option.Limit
	if limit == 0 {
		limit = 100
	}
	if limit <= 0 {
		return nil, fmt.Errorf("limit 必须大于 0")
	}
	params := map[string]any{"query": condition, "limit": limit}
	if option.Markets != nil {
		markets, err := marketsValue(option.Markets)
		if err != nil {
			return nil, err
		}
		if markets != nil {
			params["markets"] = markets
		}
	}
	data, err := QueryWencai(params)
	if err != nil {
		return nil, err
	}
	dataMap, _ := data.(map[string]any)
	if dataMap == nil {
		dataMap = map[string]any{}
	}
	return wencaiResultTable(condition, dataMap)
}

// wencaiFieldKey 构造字段的去重键：(index 数值, timestamp 文本)。
type wencaiFieldKey struct {
	index     int64
	timestamp string
}

// wencaiIdentity 从问财证券或明细行提取标准证券身份。
func wencaiIdentity(item any) map[string]string {
	row, ok := item.(map[string]any)
	if !ok {
		return map[string]string{}
	}
	securityValue, _ := row["security"].(map[string]any)
	if securityValue == nil {
		securityValue = map[string]any{}
	}
	fullCode := strings.ToUpper(strings.TrimSpace(stringOrEmpty(row["full_code"])))
	code := strings.ToUpper(strings.TrimSpace(firstNonEmptyString(stringOrEmpty(row["code"]), stringOrEmpty(securityValue["code"]))))
	market := strings.ToUpper(strings.TrimSpace(firstNonEmptyString(stringOrEmpty(row["market"]), stringOrEmpty(securityValue["market"]))))
	if fullCode != "" && len(fullCode) > 4 && isAlpha(fullCode[:4]) {
		if market == "" {
			market = fullCode[:4]
		}
		if code == "" {
			code = fullCode[4:]
		}
	} else if market != "" && code != "" {
		fullCode = market + code
	}
	if fullCode == "" && code == "" {
		return map[string]string{}
	}
	return map[string]string{"full_code": fullCode, "code": code, "market": market}
}

// wencaiIdentityKey 构造证券身份的去重键。
func wencaiIdentityKey(identity map[string]string) *wencaiKey {
	if fullCode := identity["full_code"]; fullCode != "" {
		return &wencaiKey{kind: "full_code", first: fullCode}
	}
	if code := identity["code"]; code != "" {
		return &wencaiKey{kind: "security", first: identity["market"], second: code}
	}
	return nil
}

type wencaiKey struct {
	kind   string
	first  string
	second string
}

// wencaiResultTable 组装问财查询结果的表格。
func wencaiResultTable(condition string, data map[string]any) (*Table, error) {
	fieldNames := map[wencaiFieldKey]string{}
	for _, item := range toSequenceAnyOrEmpty(data["fields"]) {
		field, ok := item.(map[string]any)
		if !ok {
			continue
		}
		key := wencaiFieldKey{index: toWencaiIndex(field["index"]), timestamp: stringOrEmpty(field["timestamp"])}
		name := strings.TrimSpace(firstNonEmptyString(
			stringOrEmpty(field["display_name"]),
			stringOrEmpty(field["show_name"]),
			stringOrEmpty(field["index_name"]),
			stringOrEmpty(field["query_key"]),
		))
		if name != "" {
			fieldNames[key] = name
		}
	}

	var securityIdentities []map[string]string
	securityIndexes := map[wencaiKey]int{}
	securityCodeIndexes := map[string]*int{}
	for _, item := range toSequenceAnyOrEmpty(data["securities"]) {
		identity := wencaiIdentity(item)
		if len(identity) == 0 {
			continue
		}
		key := wencaiIdentityKey(identity)
		if key != nil {
			if _, exists := securityIndexes[*key]; exists {
				continue
			}
		}
		index := len(securityIdentities)
		securityIdentities = append(securityIdentities, identity)
		if key != nil {
			securityIndexes[*key] = index
		}
		code := identity["code"]
		if code != "" {
			if _, exists := securityCodeIndexes[code]; exists {
				securityCodeIndexes[code] = nil
			} else {
				value := index
				securityCodeIndexes[code] = &value
			}
		}
	}

	var rows []map[string]any
	matchedSecurityIndexes := map[int]bool{}
	rowIdentityKeys := map[wencaiKey]bool{}
	for index, item := range toSequenceAnyOrEmpty(data["rows"]) {
		row, ok := item.(map[string]any)
		if !ok {
			continue
		}
		itemIdentity := wencaiIdentity(item)
		itemKey := wencaiIdentityKey(itemIdentity)
		var securityIndex *int
		if itemKey != nil {
			if value, exists := securityIndexes[*itemKey]; exists {
				securityIndex = &value
			}
		}
		if securityIndex == nil {
			if code := itemIdentity["code"]; code != "" {
				if value, exists := securityCodeIndexes[code]; exists {
					securityIndex = value
				}
			}
		}
		if securityIndex == nil && len(itemIdentity) == 0 && index < len(securityIdentities) {
			value := index
			securityIndex = &value
		}

		fallback := map[string]string{}
		if securityIndex != nil {
			fallback = securityIdentities[*securityIndex]
		}
		code := firstNonEmptyString(itemIdentity["code"], fallback["code"])
		market := firstNonEmptyString(itemIdentity["market"], fallback["market"])
		fullCode := itemIdentity["full_code"]
		if fullCode == "" && market != "" && code != "" {
			fullCode = market + code
		}
		if fullCode == "" {
			fullCode = fallback["full_code"]
		}
		rowOut := map[string]any{"full_code": fullCode, "code": code, "market": market}
		for _, cellItem := range toSequenceAnyOrEmpty(row["cells"]) {
			cell, ok := cellItem.(map[string]any)
			if !ok {
				continue
			}
			cellName := strings.TrimSpace(stringOrEmpty(cell["name"]))
			if cellName == "code" {
				continue
			}
			key := wencaiFieldKey{index: toWencaiIndex(cell["key_index"]), timestamp: stringOrEmpty(cell["timestamp"])}
			name := fieldNames[key]
			if name == "" {
				name = cellName
			}
			if name != "" {
				rowOut[name] = cell["value"]
			}
		}

		rowKey := wencaiIdentityKey(wencaiIdentity(rowOut))
		if rowKey != nil {
			if rowIdentityKeys[*rowKey] {
				continue
			}
			rowIdentityKeys[*rowKey] = true
		}
		rows = append(rows, rowOut)
		if securityIndex != nil {
			matchedSecurityIndexes[*securityIndex] = true
		}
	}

	// 某些条件只有证券池，没有可查询的动态字段；此时仍返回证券身份。
	for index, identity := range securityIdentities {
		key := wencaiIdentityKey(identity)
		if matchedSecurityIndexes[index] {
			continue
		}
		if key != nil && rowIdentityKeys[*key] {
			continue
		}
		rows = append(rows, map[string]any{
			"full_code": identity["full_code"],
			"code":      identity["code"],
			"market":    identity["market"],
		})
		if key != nil {
			rowIdentityKeys[*key] = true
		}
	}

	table, err := frame(rows, false)
	if err != nil {
		return nil, err
	}
	query, exists := data["query"]
	if !exists {
		query = condition
	}
	table.SetAttr("query", query)
	table.SetAttr("fields", data["fields"])
	return table, nil
}

// toIndex 宽松地把值转为索引整数。
func toWencaiIndex(value any) int64 {
	if number, ok := asInt64(value); ok {
		return number
	}
	if number, ok := value.(float64); ok {
		return int64(number)
	}
	if text, ok := value.(string); ok {
		if number, err := strconv.ParseInt(strings.TrimSpace(text), 10, 64); err == nil {
			return number
		}
	}
	return 0
}

// NewsOptions 控制 News 的查询参数。
type NewsOptions struct {
	Page     int    // 默认 1
	PageSize int    // 默认 20
	Tag      string // 默认 ""
}

// News 查询 7x24 快讯；分页信息保存在返回表格的 attrs 中。
// 入参：opts；语义见函数用途。
// 出参：*Table（成功时）；error 为调用失败时的错误。
func News(opts ...*NewsOptions) (*Table, error) {
	page, pageSize, tag := 1, 20, ""
	for _, item := range opts {
		if item == nil {
			continue
		}
		// 与 关键字参数一致：显式传入的值无条件生效，由下方校验拦截非法值。
		page = item.Page
		pageSize = item.PageSize
		tag = item.Tag
	}
	if page <= 0 || pageSize <= 0 {
		return nil, fmt.Errorf("page 和 page_size 必须大于 0")
	}
	data, err := ListNewsFlashItems(map[string]any{
		"traversal": "page",
		"page":      page,
		"page_size": pageSize,
		"tag":       strings.TrimSpace(tag),
	})
	if err != nil {
		return nil, err
	}
	dataMap, _ := data.(map[string]any)
	if dataMap == nil {
		dataMap = map[string]any{}
	}
	table, err := frame(dataMap["items"], false)
	if err != nil {
		return nil, err
	}
	for key, value := range dataMap {
		if key != "items" {
			table.SetAttr(key, value)
		}
	}
	return table, nil
}

// MarketSecuritiesOptions 控制 MarketSecurities 的查询参数。
type MarketSecuritiesOptions struct {
	Start int // 默认 0
	Count int // 默认 100
}

// MarketSecurities 按市场列出证券；总数等分页信息保存在返回表格的 attrs 中。
// 入参：market、opts；语义见函数用途。
// 出参：*Table（成功时）；error 为调用失败时的错误。
func MarketSecurities(market string, opts ...*MarketSecuritiesOptions) (*Table, error) {
	start, count := 0, 100
	for _, item := range opts {
		if item == nil {
			continue
		}
		start = item.Start
		count = item.Count
	}
	if start < 0 || count <= 0 {
		return nil, fmt.Errorf("start 不能小于 0，count 必须大于 0")
	}
	market, err := marketValue(market)
	if err != nil {
		return nil, err
	}
	data, err := ListMarketSecurities(map[string]any{
		"market":     market,
		"sort_begin": start,
		"sort_count": count,
	})
	if err != nil {
		return nil, err
	}
	dataMap, _ := data.(map[string]any)
	if dataMap == nil {
		dataMap = map[string]any{}
	}
	table, err := frame(identityRows(dataMap["securities"]), false)
	if err != nil {
		return nil, err
	}
	for key, value := range dataMap {
		if key != "securities" {
			table.SetAttr(key, value)
		}
	}
	return table, nil
}

// BlockConstituentsOptions 控制 BlockConstituents 的查询参数。
type BlockConstituentsOptions struct {
	Start int // 默认 0
	Count int // 默认 100
}

// BlockConstituents 按板块 ID 或名称查询成分证券。
// 入参：block、opts；语义见函数用途。
// 出参：*Table（成功时）；error 为调用失败时的错误。
func BlockConstituents(block any, opts ...*BlockConstituentsOptions) (*Table, error) {
	start, count := 0, 100
	for _, item := range opts {
		if item == nil {
			continue
		}
		start = item.Start
		count = item.Count
	}
	if start < 0 || count <= 0 {
		return nil, fmt.Errorf("start 不能小于 0，count 必须大于 0")
	}

	var identity map[string]any
	if id, ok := asInt64(block); ok {
		if id <= 0 {
			return nil, fmt.Errorf("板块 ID 必须大于 0")
		}
		identity = map[string]any{"id": id}
	} else {
		name := strings.TrimSpace(fmt.Sprint(block))
		if name == "" {
			return nil, fmt.Errorf("板块名称不能为空")
		}
		resolved, err := ResolveBlock(map[string]any{"name": name})
		if err != nil {
			return nil, err
		}
		identity, _ = resolved.(map[string]any)
		if identity == nil {
			identity = map[string]any{}
		}
	}
	result, err := ListBlockConstituents(map[string]any{
		"block":      identity,
		"sort_begin": start,
		"sort_count": count,
		"sort_order": "D",
	})
	if err != nil {
		return nil, err
	}
	return frame(identityRows(result), false)
}

// stringOrEmpty 返回映射中的字符串值，缺失时为空串。
func stringOrEmpty(value any) string {
	if value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

// firstNonEmpty 返回第一个非空字符串。
func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

// toSequenceAnyOrEmpty 把值转为序列，非序列时返回空序列。
func toSequenceAnyOrEmpty(value any) []any {
	if sequence, ok := toSequenceAny(value); ok {
		return sequence
	}
	return nil
}
