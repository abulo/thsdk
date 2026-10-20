package thsdk

// 本文件实现公开 API 的统一请求包装层。
//
// 核心职责：
//   - CallOptions / WithTimeout：以「秒」为单位的调用超时选项；
//   - applyTimeout：把选项校验并转换为原生协议要求的毫秒指针；
//   - 字段规范化：normalizeSecurity（证券完整代码 → {code, market}）、
//     normalizeBlock（板块 ID / 映射）、normalizeNamedFields（命名请求
//     字段的完整规范化，含 datetime 字段）；
//   - 四类调用入口：callNone（无参）/ callTyped（对象参数）/
//     callString（单字符串）/ callSecurities（证券序列）。
//
// 设计说明：所有传入的 params 映射都会按已知字段清单执行与公开 API
// 一致的规范化（security 完整代码拆分、market 大写化、datetime 转
// ISO 字符串等），保证发往原生组件的请求结构统一、合法。

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unicode"
)

// CallOptions 控制一次 API 调用的超时，单位为秒。
//
// 字段说明：
//   - Timeout：超时秒数；0 表示不指定（使用原生默认超时），
//     否则必须大于 0 且不超过 600 秒。
type CallOptions struct {
	Timeout float64
}

// WithTimeout 构造一个带指定超时的调用选项。
//
// 入参：
//   - seconds：超时秒数（大于 0 且不超过 600）。
//
// 出参：
//   - *CallOptions：可直接传入各公开方法的 opts 参数。
func WithTimeout(seconds float64) *CallOptions {
	return &CallOptions{Timeout: seconds}
}

// applyTimeout 校验并转换调用选项为毫秒指针。
//
// 入参：
//   - opts：调用选项列表；取第一个合法非零选项生效。
//
// 出参：
//   - *int：原生协议要求的 timeout_ms；全部选项无效时返回 nil；
//   - error：选项超时值超过上限时返回错误。
//
// 规则：跳过 nil 与 Timeout<=0 的选项；首个合法选项的超时按秒转毫秒，
// 超过 maxTimeoutSeconds（600 秒）直接报错。
func applyTimeout(opts []*CallOptions) (*int, error) {
	var timeoutMs *int
	for _, option := range opts {
		// 跳过无效选项（nil 或非正超时）。
		if option == nil || option.Timeout <= 0 {
			continue
		}
		// 超时上限校验（600 秒）。
		if option.Timeout > maxTimeoutSeconds {
			return nil, fmt.Errorf("timeout 必须大于 0 且不超过 600 秒")
		}
		// 秒转毫秒，仅取首个合法选项。
		milliseconds := timeoutMilliseconds(option.Timeout)
		timeoutMs = &milliseconds
		break
	}
	return timeoutMs, nil
}

// normalizeSecurity 把证券入参规范化为 {code, market} 请求结构。
//
// 入参：
//   - value：完整代码字符串（如 "600000.SH"）或含 code/market 的映射。
//
// 出参：
//   - map[string]any：规范化后的 {code, market} 结构；
//   - error：输入类型不合法或 market 不是四位字母时返回错误。
//
// 规则：
//   - 字符串：前 4 位（大写）作为 market，其余作为 code；
//   - 映射：分别取 code 与 market（market 转为大写）；
//   - market 必须为 4 位 Unicode 字母且 code 非空。
func normalizeSecurity(value any) (map[string]any, error) {
	// 分别存放解析出的代码与市场。
	var code, market string
	switch item := value.(type) {
	case map[string]any:
		// 映射形态：直接取字段并清理空白、大写化市场。
		code = strings.TrimSpace(fmt.Sprint(item["code"]))
		market = strings.ToUpper(strings.TrimSpace(fmt.Sprint(item["market"])))
	case string:
		// 完整代码字符串：前 4 位为市场，其余为代码。
		fullCode := strings.TrimSpace(item)
		market = strings.ToUpper(fullCode[:min(len(fullCode), 4)])
		code = fullCode[min(len(fullCode), 4):]
	default:
		// 其他类型一律拒绝。
		return nil, fmt.Errorf("证券必须是完整代码字符串或包含 code/market 的映射")
	}
	// 市场必须为 4 位字母，代码不能为空。
	if len(market) != 4 || !isAlpha(market) || code == "" {
		return nil, fmt.Errorf("请使用完整证券代码，或包含 code 和四位 market 的映射")
	}
	return map[string]any{"code": code, "market": market}, nil
}

// normalizeBlock 把板块入参规范化为请求结构。
//
// 入参：
//   - value：板块 ID（int/int64/int32，须为正）或标准板块映射。
//
// 出参：
//   - map[string]any：规范化后的结构；
//   - error：bool 或非法 ID 时返回错误。
//
// 规则：int 型统一转 int64 并校验大于 0；映射原样复制返回。
func normalizeBlock(value any) (map[string]any, error) {
	switch item := value.(type) {
	case bool:
		// 布尔值不是合法板块。
		return nil, fmt.Errorf("板块必须是 ID 或标准板块映射")
	case int:
		// 板块 ID 必须为正数。
		if item <= 0 {
			return nil, fmt.Errorf("板块 ID 必须大于 0")
		}
		return map[string]any{"id": int64(item)}, nil
	case int64:
		if item <= 0 {
			return nil, fmt.Errorf("板块 ID 必须大于 0")
		}
		return map[string]any{"id": item}, nil
	case int32:
		if item <= 0 {
			return nil, fmt.Errorf("板块 ID 必须大于 0")
		}
		return map[string]any{"id": int64(item)}, nil
	case map[string]any:
		// 标准板块映射：整体复制一份返回，避免外部修改影响内部。
		copied := make(map[string]any, len(item))
		for key, value := range item {
			copied[key] = value
		}
		return copied, nil
	default:
		// 其他类型一律拒绝。
		return nil, fmt.Errorf("板块必须是 ID 或标准板块映射")
	}
}

// isAlpha 判断字符串是否为 Unicode 字母序列。
//
// 入参：
//   - value：待检查字符串。
//
// 出参：
//   - bool：非空且全部字符为字母时返回 true。
//
// 说明：Unicode 字母判定允许中文、英文等；空串返回 false。
func isAlpha(value string) bool {
	// 空串不构成字母串。
	if value == "" {
		return false
	}
	// 逐字符检查是否为字母。
	for _, r := range value {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// normalizeNamedFields 对命名请求字段执行完整规范化。
//
// 入参：
//   - fields：调用方传入的原始请求字段映射；
//   - datetimeFields：需要按 ISO 8601 时间字符串处理的字段名集合。
//
// 出参：
//   - map[string]any：规范化后的请求结构；
//   - error：任一步规范化失败时返回错误。
//
// 规范化项：
//   - security / etf：完整代码 → {code, market}；
//   - securities / exclude_securities：证券序列逐项规范化；
//   - market：转为大写；
//   - markets：字符串（逗号/空格分隔）或序列 → 大写序列；
//   - block：板块 ID/映射规范化；
//   - blocks / exclude_blocks：板块序列逐项规范化；
//   - datetimeFields 中出现的字段：转 ISO 8601 字符串；
//   - 最终整体过一遍 jsonValue 规则校验。
func normalizeNamedFields(fields map[string]any, datetimeFields map[string]bool) (map[string]any, error) {
	// 先复制一份，避免修改调用方的映射。
	payload := make(map[string]any, len(fields))
	for key, value := range fields {
		payload[key] = value
	}
	// 单证券字段规范化。
	for _, key := range []string{"security", "etf"} {
		if value, ok := payload[key]; ok {
			normalized, err := normalizeSecurity(value)
			if err != nil {
				return nil, err
			}
			payload[key] = normalized
		}
	}
	// 证券序列字段逐项规范化。
	for _, key := range []string{"securities", "exclude_securities"} {
		if value, ok := payload[key]; ok {
			items, ok := toSequence(value)
			if !ok {
				return nil, fmt.Errorf("%s 必须是证券序列", key)
			}
			normalized := make([]any, 0, len(items))
			for _, item := range items {
				converted, err := normalizeSecurity(item)
				if err != nil {
					return nil, err
				}
				normalized = append(normalized, converted)
			}
			payload[key] = normalized
		}
	}
	// 单市场字段大写化。
	if value, ok := payload["market"]; ok {
		payload["market"] = strings.ToUpper(strings.TrimSpace(fmt.Sprint(value)))
	}
	// 市场序列字段：字符串（逗号/空格分隔）或序列 → 大写序列。
	if value, ok := payload["markets"]; ok {
		var items []any
		if text, isString := value.(string); isString {
			for _, part := range strings.FieldsFunc(text, func(r rune) bool { return r == ',' || r == ' ' }) {
				if part != "" {
					items = append(items, part)
				}
			}
		} else {
			sequence, ok := toSequence(value)
			if !ok {
				return nil, fmt.Errorf("markets 必须是市场序列")
			}
			items = sequence
		}
		normalized := make([]any, 0, len(items))
		for _, item := range items {
			normalized = append(normalized, strings.ToUpper(strings.TrimSpace(fmt.Sprint(item))))
		}
		payload["markets"] = normalized
	}
	// 单板块字段规范化。
	if value, ok := payload["block"]; ok {
		normalized, err := normalizeBlock(value)
		if err != nil {
			return nil, err
		}
		payload["block"] = normalized
	}
	// 板块序列字段逐项规范化。
	for _, key := range []string{"blocks", "exclude_blocks"} {
		if value, ok := payload[key]; ok {
			items, ok := toSequence(value)
			if !ok {
				return nil, fmt.Errorf("%s 必须是板块序列", key)
			}
			normalized := make([]any, 0, len(items))
			for _, item := range items {
				converted, err := normalizeBlock(item)
				if err != nil {
					return nil, err
				}
				normalized = append(normalized, converted)
			}
			payload[key] = normalized
		}
	}
	// datetime 字段转 ISO 8601 字符串。
	for key := range datetimeFields {
		if value, ok := payload[key]; ok {
			converted, err := jsonValue(value, 0)
			if err != nil {
				return nil, err
			}
			payload[key] = converted
		}
	}
	// 整体过一遍严格 JSON 规则校验。
	converted, err := jsonValue(payload, 0)
	if err != nil {
		return nil, err
	}
	result, ok := converted.(map[string]any)
	if !ok {
		return nil, newTHSDKError("请求字段规范化结果类型异常")
	}
	return result, nil
}

// toSequence 判断值是否为「非字符串/字节」的序列并转成 []any。
//
// 入参：
//   - value：任意值。
//
// 出参：
//   - []any：切片/数组形态的元素列表；
//   - bool：true 表示成功转换；字符串、字节序列或其他类型返回 false。
func toSequence(value any) ([]any, bool) {
	// 字符串与字节序列不是证券/板块序列。
	if _, isString := value.(string); isString {
		return nil, false
	}
	if _, isBytes := value.([]byte); isBytes {
		return nil, false
	}
	// 反射取切片/数组元素。
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		result := make([]any, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			result[i] = rv.Index(i).Interface()
		}
		return result, true
	default:
		// 其他类型不支持。
		return nil, false
	}
}

// callTyped 执行一次对象参数方法调用。
//
// 入参：
//   - method：原生方法名；
//   - params：请求参数映射；nil 视为空映射；
//   - required：必填字段名清单，缺失时报错；
//   - datetimeFields：需要按时间字符串规范化的字段名清单；
//   - opts：调用选项列表。
//
// 出参：
//   - any：原生响应的 data 字段；
//   - error：未知字段、缺必填字段、规范化失败或调用失败时返回错误。
//
// 流程：按方法清单校验未知/必填字段 → 规范化 → 经默认客户端请求。
func callTyped(method string, params map[string]any, required []string, datetimeFields []string, opts ...*CallOptions) (any, error) {
	// 从方法清单取规格，方法必须已登记。
	spec, ok := nativeMethods[method]
	if !ok {
		return nil, newTHSDKError("未知的 Canonical 方法: %s", method)
	}
	// nil 参数按空映射处理。
	if params == nil {
		params = map[string]any{}
	}
	// 未知字段检测：不在方法声明字段清单内的键一律拒绝。
	var unknown []string
	for key := range params {
		if !stringInSlice(key, spec.RequestFields) {
			unknown = append(unknown, key)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		return nil, fmt.Errorf("未知请求字段: %s", strings.Join(unknown, ", "))
	}
	// 必填字段检测：缺失即报错。
	var missing []string
	for _, name := range required {
		if _, present := params[name]; !present {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("缺少必填请求字段: %s", strings.Join(missing, ", "))
	}

	// datetime 字段名集合化。
	datetimeSet := make(map[string]bool, len(datetimeFields))
	for _, name := range datetimeFields {
		datetimeSet[name] = true
	}
	// 规范化请求字段。
	payload, err := normalizeNamedFields(params, datetimeSet)
	if err != nil {
		return nil, err
	}
	// 校验超时选项。
	if _, err := applyTimeout(opts); err != nil {
		return nil, err
	}
	// 经默认客户端发起请求。
	return getClient().Request(method, payload, opts...)
}

// callNone 执行一次无参数方法调用。
//
// 入参：
//   - method：原生方法名；
//   - opts：调用选项列表。
//
// 出参：
//   - any：原生响应的 data 字段；
//   - error：调用失败时返回错误。
func callNone(method string, opts ...*CallOptions) (any, error) {
	// 校验超时选项。
	if _, err := applyTimeout(opts); err != nil {
		return nil, err
	}
	return getClient().Request(method, nil, opts...)
}

// callString 执行一次单字符串参数方法调用。
//
// 入参：
//   - method：原生方法名；
//   - value：字符串参数；
//   - argument：参数名（用于错误提示）；
//   - opts：调用选项列表。
//
// 出参：
//   - any：原生响应的 data 字段；
//   - error：参数为空或调用失败时返回错误。
//
// 规则：value 去首尾空白后必须非空，作为请求体整体发送。
func callString(method, value, argument string, opts ...*CallOptions) (any, error) {
	// 空值直接拒绝。
	if value == "" || strings.TrimSpace(value) == "" {
		return nil, fmt.Errorf("%s 必须是非空字符串", argument)
	}
	// 校验超时选项。
	if _, err := applyTimeout(opts); err != nil {
		return nil, err
	}
	return getClient().Request(method, strings.TrimSpace(value), opts...)
}

// callSecurities 执行一次证券序列参数方法调用。
//
// 入参：
//   - method：原生方法名；
//   - securities：证券入参序列（完整代码字符串或映射）；
//   - opts：调用选项列表。
//
// 出参：
//   - any：原生响应的 data 字段；
//   - error：序列为空、规范化失败或调用失败时返回错误。
//
// 规则：逐项经 normalizeSecurity 规范化后作为请求体数组发送。
func callSecurities(method string, securities []any, opts ...*CallOptions) (any, error) {
	// 空序列直接拒绝。
	if len(securities) == 0 {
		return nil, fmt.Errorf("securities 不能为空")
	}
	// 逐项规范化证券。
	payload := make([]any, 0, len(securities))
	for _, item := range securities {
		normalized, err := normalizeSecurity(item)
		if err != nil {
			return nil, err
		}
		payload = append(payload, normalized)
	}
	// 校验超时选项。
	if _, err := applyTimeout(opts); err != nil {
		return nil, err
	}
	return getClient().Request(method, payload, opts...)
}

// stringInSlice 判断字符串切片是否包含指定值。
//
// 入参：
//   - value：目标字符串；
//   - values：待检查切片。
//
// 出参：
//   - bool：命中返回 true。
func stringInSlice(value string, values []string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}
