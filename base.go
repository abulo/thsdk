package thsdk

// 本文件对应 Python 包的 thsdk/_apis/_base.py。
//
// 集中实现请求字段规范化（证券/板块/市场/日期）、必填与未知字段校验，
// 以及四类调用入口（none / object / string / securities）。
//
// 设计说明：Python 的 *_apis 包装函数同时支持「params 映射」与「关键字参数」
// 两种调用方式，其中关键字参数路径会执行完整的字段规范化。Go 没有关键字参数，
// 因此本移植统一走规范化路径：所有传入的 params 映射都会按已知字段清单执行
// 与 Python 关键字路径一致的规范化（security 完整代码 → {code, market}、
// market 大写化、datetime → ISO 字符串等）。

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unicode"
)

// CallOptions 控制一次 API 调用的超时（秒）。
// Timeout 为 0 表示不指定（原生默认），否则必须大于 0 且不超过 600 秒。
type CallOptions struct {
	Timeout float64
}

// WithTimeout 返回携带指定超时（秒）的调用选项。
func WithTimeout(seconds float64) *CallOptions {
	return &CallOptions{Timeout: seconds}
}

// applyTimeout 校验并转换调用选项为毫秒指针，对应 Python _timeout + _timeout_milliseconds。
func applyTimeout(opts []*CallOptions) (*int, error) {
	var timeoutMs *int
	for _, option := range opts {
		if option == nil || option.Timeout <= 0 {
			continue
		}
		if option.Timeout > maxTimeoutSeconds {
			return nil, fmt.Errorf("timeout 必须大于 0 且不超过 600 秒")
		}
		milliseconds := timeoutMilliseconds(option.Timeout)
		timeoutMs = &milliseconds
		break
	}
	return timeoutMs, nil
}

// normalizeSecurity 把完整代码字符串或 code/market 映射规范化为请求结构，
// 对应 Python _normalize_security。输入类型对应 Python 的
// SecurityInput = str | Mapping[str, Any]。
func normalizeSecurity(value any) (map[string]any, error) {
	var code, market string
	switch item := value.(type) {
	case map[string]any:
		code = strings.TrimSpace(fmt.Sprint(item["code"]))
		market = strings.ToUpper(strings.TrimSpace(fmt.Sprint(item["market"])))
	case string:
		fullCode := strings.TrimSpace(item)
		market = strings.ToUpper(fullCode[:min(len(fullCode), 4)])
		code = fullCode[min(len(fullCode), 4):]
	default:
		return nil, fmt.Errorf("证券必须是完整代码字符串或包含 code/market 的映射")
	}
	if len(market) != 4 || !isAlpha(market) || code == "" {
		return nil, fmt.Errorf("请使用完整证券代码，或包含 code 和四位 market 的映射")
	}
	return map[string]any{"code": code, "market": market}, nil
}

// normalizeBlock 把板块 ID 或标准板块映射规范化为请求结构，
// 对应 Python _normalize_block。输入类型对应 Python 的
// BlockInput = int | Mapping[str, Any]。
func normalizeBlock(value any) (map[string]any, error) {
	switch item := value.(type) {
	case bool:
		return nil, fmt.Errorf("板块必须是 ID 或标准板块映射")
	case int:
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
		copied := make(map[string]any, len(item))
		for key, value := range item {
			copied[key] = value
		}
		return copied, nil
	default:
		return nil, fmt.Errorf("板块必须是 ID 或标准板块映射")
	}
}

// isAlpha 对应 Python str.isalpha()（Unicode 字母判定，允许中文等字母）。
func isAlpha(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// normalizeNamedFields 对命名请求字段执行完整规范化，对应 Python _normalize_named_fields。
// datetimeFields 中的字段会按 ISO 8601 时间字符串处理。
func normalizeNamedFields(fields map[string]any, datetimeFields map[string]bool) (map[string]any, error) {
	payload := make(map[string]any, len(fields))
	for key, value := range fields {
		payload[key] = value
	}
	for _, key := range []string{"security", "etf"} {
		if value, ok := payload[key]; ok {
			normalized, err := normalizeSecurity(value)
			if err != nil {
				return nil, err
			}
			payload[key] = normalized
		}
	}
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
	if value, ok := payload["market"]; ok {
		payload["market"] = strings.ToUpper(strings.TrimSpace(fmt.Sprint(value)))
	}
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
	if value, ok := payload["block"]; ok {
		normalized, err := normalizeBlock(value)
		if err != nil {
			return nil, err
		}
		payload["block"] = normalized
	}
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
	for key := range datetimeFields {
		if value, ok := payload[key]; ok {
			converted, err := jsonValue(value, 0)
			if err != nil {
				return nil, err
			}
			payload[key] = converted
		}
	}
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

// toSequence 判断值是否为非字符串序列（对应 Python 的
// isinstance(value, (str, bytes)) or not isinstance(value, Sequence) 排除逻辑）。
func toSequence(value any) ([]any, bool) {
	if _, isString := value.(string); isString {
		return nil, false
	}
	if _, isBytes := value.([]byte); isBytes {
		return nil, false
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		result := make([]any, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			result[i] = rv.Index(i).Interface()
		}
		return result, true
	default:
		return nil, false
	}
}

// callTyped 执行一次对象参数方法调用，对应 Python _call_typed。
// params 的未知字段与必填字段按清单校验，然后执行规范化。
func callTyped(method string, params map[string]any, required []string, datetimeFields []string, opts ...*CallOptions) (any, error) {
	spec, ok := nativeMethods[method]
	if !ok {
		return nil, newTHSDKError("未知的 Canonical 方法: %s", method)
	}
	if params == nil {
		params = map[string]any{}
	}
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
	var missing []string
	for _, name := range required {
		if _, present := params[name]; !present {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("缺少必填请求字段: %s", strings.Join(missing, ", "))
	}

	datetimeSet := make(map[string]bool, len(datetimeFields))
	for _, name := range datetimeFields {
		datetimeSet[name] = true
	}
	payload, err := normalizeNamedFields(params, datetimeSet)
	if err != nil {
		return nil, err
	}
	if _, err := applyTimeout(opts); err != nil {
		return nil, err
	}
	return getClient().Request(method, payload, opts...)
}

// callNone 执行一次无参数方法调用，对应 Python _call_none。
func callNone(method string, opts ...*CallOptions) (any, error) {
	if _, err := applyTimeout(opts); err != nil {
		return nil, err
	}
	return getClient().Request(method, nil, opts...)
}

// callString 执行一次单字符串参数方法调用，对应 Python _call_string。
func callString(method, value, argument string, opts ...*CallOptions) (any, error) {
	if value == "" || strings.TrimSpace(value) == "" {
		return nil, fmt.Errorf("%s 必须是非空字符串", argument)
	}
	if _, err := applyTimeout(opts); err != nil {
		return nil, err
	}
	return getClient().Request(method, strings.TrimSpace(value), opts...)
}

// callSecurities 执行一次证券序列参数方法调用，对应 Python _call_securities。
func callSecurities(method string, securities []any, opts ...*CallOptions) (any, error) {
	if len(securities) == 0 {
		return nil, fmt.Errorf("securities 不能为空")
	}
	payload := make([]any, 0, len(securities))
	for _, item := range securities {
		normalized, err := normalizeSecurity(item)
		if err != nil {
			return nil, err
		}
		payload = append(payload, normalized)
	}
	if _, err := applyTimeout(opts); err != nil {
		return nil, err
	}
	return getClient().Request(method, payload, opts...)
}

func stringInSlice(value string, values []string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}
	return false
}
