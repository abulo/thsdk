package thsdk

// 本文件实现原生协议要求的 JSON 文档规则层：
//   - validateJSONValue：跨入原生 ABI 之前的严格 JSON 文档规则校验；
//   - jsonValue：请求值序列化转换（时间 → ISO 8601、枚举 → 值、容器递归）；
//   - 响应数字归一化（json.Number → int64/float64）与整数判定助手。

import (
	"bytes"
	"encoding/json"
	"math"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// shanghai 是 Asia/Shanghai 时区，用于无时区时间的规范化。
var shanghai = func() *time.Location {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// 极少数系统缺少时区数据库时回退到固定 UTC+8。
		return time.FixedZone("Asia/Shanghai", 8*60*60)
	}
	return loc
}()

// jsonFieldName 匹配合法的请求字段名：仅允许 ASCII 小写字母、数字、下划线。
var jsonFieldName = regexp.MustCompile(`^[a-z0-9_]+$`)

// isoFormat 是时间序列化的 ISO 8601 布局。
// 秒以下最多 6 位，无小数时不输出小数部分（Go 的 999999 布局会裁剪尾部零）。
// 使用数值偏移（-07:00）而非 Z07:00，保证零时区也输出 "+00:00" 而不是 "Z"。
const isoFormat = "2006-01-02T15:04:05.999999-07:00"

// validateJSONValue 应用原生 ABI 的 JSON 文档规则。
func validateJSONValue(value any, depth int) error {
	if depth > maxJSONDepth {
		return newTHSDKError("THSDK 请求的 JSON 嵌套超过 64 层")
	}
	switch v := value.(type) {
	case map[string]any:
		for key, item := range v {
			if !jsonFieldName.MatchString(key) {
				return newTHSDKError("THSDK 请求字段必须使用 ASCII 小写 snake_case")
			}
			if err := validateJSONValue(item, depth+1); err != nil {
				return err
			}
		}
		return nil
	case []any:
		for _, item := range v {
			if err := validateJSONValue(item, depth+1); err != nil {
				return err
			}
		}
		return nil
	case float64:
		if math.IsInf(v, 0) || math.IsNaN(v) {
			return newTHSDKError("THSDK 请求不能包含 NaN 或 Infinity")
		}
		return nil
	case float32:
		if math.IsInf(float64(v), 0) || math.IsNaN(float64(v)) {
			return newTHSDKError("THSDK 请求不能包含 NaN 或 Infinity")
		}
		return nil
	case nil, string, int, int64, bool:
		return nil
	}

	// 反射兜底：任意 map / slice（如 []Security、map[string]int）按同一规则递归。
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Map:
		for _, key := range rv.MapKeys() {
			keyString, ok := key.Interface().(string)
			if !ok || !jsonFieldName.MatchString(keyString) {
				return newTHSDKError("THSDK 请求字段必须使用 ASCII 小写 snake_case")
			}
			if err := validateJSONValue(rv.MapIndex(key).Interface(), depth+1); err != nil {
				return err
			}
		}
		return nil
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			if err := validateJSONValue(rv.Index(i).Interface(), depth+1); err != nil {
				return err
			}
		}
		return nil
	}
	return newTHSDKError("THSDK 请求包含不可序列化的值类型: %s", typeName(value))
}

// jsonValue 把常见 Go 容器与时间值转换为可严格序列化的结构。
//
//   - time.Time：无时区（time.UTC）按 Asia/Shanghai 处理，输出 ISO 8601；
//   - map：字段名必须是 ASCII 小写 snake_case，值为 nil 的键会被丢弃；
//   - 切片/数组：逐项递归（nil 保留）；
//   - 整数类型统一转为 int64，浮点统一转为 float64（拒绝 NaN/Infinity）；
//   - 其余原始值原样返回。
//
// jsonValue 把常见 Go 容器与时间值转换为可严格序列化的结构。
//
//   - time.Time：无时区（time.UTC）按 Asia/Shanghai 处理，输出 ISO 8601；
//   - map：字段名必须是 ASCII 小写 snake_case，值为 nil 的键会被丢弃；
//   - 切片/数组：逐项递归（nil 保留）；
//   - 整数类型统一转为 int64，浮点统一转为 float64（拒绝 NaN/Infinity）；
//   - 其余原始值原样返回。
//
// 入参：
//   - value：入参说明见函数用途。
//   - depth：入参说明见函数用途。
//
// 出参：any,error。
func jsonValue(value any, depth int) (any, error) {
	if depth > maxJSONDepth {
		return nil, newTHSDKError("THSDK 请求的 JSON 嵌套超过 64 层")
	}
	switch v := value.(type) {
	case time.Time:
		// 无时区时间按上海墙钟时间处理；墙钟读数不变（非时区换算）。
		if v.Location() == time.UTC {
			v = assignShanghai(v)
		}
		return v.Format(isoFormat), nil
	case map[string]any:
		result := make(map[string]any, len(v))
		for key, item := range v {
			if !jsonFieldName.MatchString(key) {
				return nil, newTHSDKError("请求字段必须使用 ASCII 小写 snake_case")
			}
			if item == nil {
				continue
			}
			converted, err := jsonValue(item, depth+1)
			if err != nil {
				return nil, err
			}
			result[key] = converted
		}
		return result, nil
	case []any:
		result := make([]any, len(v))
		for i, item := range v {
			converted, err := jsonValue(item, depth+1)
			if err != nil {
				return nil, err
			}
			result[i] = converted
		}
		return result, nil
	case nil, string, bool:
		return value, nil
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case int32:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case float64:
		if math.IsInf(v, 0) || math.IsNaN(v) {
			return nil, newTHSDKError("THSDK 请求不能包含 NaN 或 Infinity")
		}
		return v, nil
	case float32:
		if math.IsInf(float64(v), 0) || math.IsNaN(float64(v)) {
			return nil, newTHSDKError("THSDK 请求不能包含 NaN 或 Infinity")
		}
		return float64(v), nil
	case json.Number:
		return v, nil
	}

	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			return nil, nil
		}
		return jsonValue(rv.Elem().Interface(), depth+1)
	case reflect.Map:
		result := make(map[string]any, rv.Len())
		iter := rv.MapRange()
		for iter.Next() {
			keyString, ok := iter.Key().Interface().(string)
			if !ok || !jsonFieldName.MatchString(keyString) {
				return nil, newTHSDKError("请求字段必须使用 ASCII 小写 snake_case")
			}
			item := iter.Value().Interface()
			if item == nil {
				continue
			}
			converted, err := jsonValue(item, depth+1)
			if err != nil {
				return nil, err
			}
			result[keyString] = converted
		}
		return result, nil
	case reflect.Slice, reflect.Array:
		result := make([]any, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			converted, err := jsonValue(rv.Index(i).Interface(), depth+1)
			if err != nil {
				return nil, err
			}
			result[i] = converted
		}
		return result, nil
	}
	return nil, newTHSDKError("THSDK 请求包含不可序列化的值类型: %s", typeName(value))
}

// typeName 返回值的 Go 类型名，用于错误信息。
func typeName(value any) string {
	t := reflect.TypeOf(value)
	if t == nil {
		return "<nil>"
	}
	return t.String()
}

// encodeStrictJSON 以严格 JSON 方式序列化请求：
// ensure_ascii=False（不转义非 ASCII）、allow_nan=False、紧凑分隔符。
// json.dumps(..., ensure_ascii=False, allow_nan=False, separators=(",", ":"))。
// encodeStrictJSON 以严格 JSON 方式序列化请求：
// ensure_ascii=False（不转义非 ASCII）、allow_nan=False、紧凑分隔符。
// json.dumps(..., ensure_ascii=False, allow_nan=False, separators=(",", ":"))。
// 入参：
//   - request：入参说明见函数用途。
//
// 出参：[]byte,error。
func encodeStrictJSON(request map[string]any) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(request); err != nil {
		return nil, newTHSDKError("THSDK 请求无法编码为严格 JSON")
	}
	return bytes.TrimSuffix(buffer.Bytes(), []byte("\n")), nil
}

// parseResponse 解析原生运行组件的 JSON 响应并验证基本结构。
// 数字保持为 json.Number。
// 由调用方按需转换。
// parseResponse 解析原生运行组件的 JSON 响应并验证基本结构。
// 数字保持为 json.Number。
// 由调用方按需转换。
// 入参：
//   - raw：入参说明见函数用途。
//
// 出参：map[string]any,error。
func parseResponse(raw string) (map[string]any, error) {
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.UseNumber()
	var response map[string]any
	if err := decoder.Decode(&response); err != nil {
		return nil, newTHSDKError("THSDK 运行组件返回了无效数据")
	}
	if _, ok := response["ok"].(bool); !ok {
		return nil, newTHSDKError("THSDK 运行组件返回了不兼容的数据")
	}
	return response, nil
}

// normalizeJSON 递归把 json.Number 转为 int64（整数）或 float64（小数）。
// 使返回给调用方的数据与 json.loads 的 int/float 语义一致。
func normalizeJSON(value any) any {
	switch v := value.(type) {
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return i
		}
		if f, err := v.Float64(); err == nil {
			return f
		}
		return v.String()
	case map[string]any:
		out := make(map[string]any, len(v))
		for key, item := range v {
			out[key] = normalizeJSON(item)
		}
		return out
	case []any:
		out := make([]any, len(v))
		for i, item := range v {
			out[i] = normalizeJSON(item)
		}
		return out
	default:
		return value
	}
}

// asInt64 仅当值可以表示为整数（且不是布尔、不是浮点）时返回 true。
// isinstance(value, int) and not isinstance(value, bool)。
// asInt64 仅当值可以表示为整数（且不是布尔、不是浮点）时返回 true。
// isinstance(value, int) and not isinstance(value, bool)。
// 入参：
//   - value：入参说明见函数用途。
//
// 出参：int64,bool。
func asInt64(value any) (int64, bool) {
	switch v := value.(type) {
	case json.Number:
		i, err := v.Int64()
		return i, err == nil
	case int:
		return int64(v), true
	case int64:
		return v, true
	case int32:
		return int64(v), true
	case int16:
		return int64(v), true
	case int8:
		return int64(v), true
	case uint:
		return int64(v), true
	case uint64:
		return int64(v), v <= math.MaxInt64
	case uint32:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint8:
		return int64(v), true
	default:
		return 0, false
	}
}
