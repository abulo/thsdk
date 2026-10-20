package thsdk

import (
	"encoding/json"
	"math"
	"strings"
	"testing"
	"time"
)

func TestEncodeStrictJSONCompact(t *testing.T) {
	encoded, err := encodeStrictJSON(map[string]any{
		"security": "000001.SZ",
		"count":    int64(3),
	})
	if err != nil {
		t.Fatal(err)
	}
	text := string(encoded)
	if strings.Contains(text, " ") || strings.Contains(text, "\n") {
		t.Fatalf("JSON 应使用紧凑分隔符: %s", text)
	}
	if !strings.Contains(text, `"security":"000001.SZ"`) {
		t.Fatalf("JSON 内容不符: %s", text)
	}
}

func TestEncodeStrictJSONNonASCII(t *testing.T) {
	encoded, err := encodeStrictJSON(map[string]any{"name": "贵州茅台"})
	if err != nil {
		t.Fatal(err)
	}
	text := string(encoded)
	if strings.Contains(text, `\u`) {
		t.Fatalf("ensure_ascii=False：中文不应被转义: %s", text)
	}
}

func TestParseResponseUsesNumber(t *testing.T) {
	response, err := parseResponse(`{"ok":true,"data":{"price":12.5,"count":7}}`)
	if err != nil {
		t.Fatal(err)
	}
	data := response["data"].(map[string]any)
	if _, isNumber := data["count"].(json.Number); !isNumber {
		t.Fatal("响应数字必须保持为 json.Number")
	}
	if _, isNumber := data["price"].(json.Number); !isNumber {
		t.Fatal("响应数字必须保持为 json.Number")
	}
}

func TestParseResponseInvalid(t *testing.T) {
	if _, err := parseResponse(`not-json`); err == nil {
		t.Fatal("非法 JSON 应报错")
	}
	if _, err := parseResponse(`{"data":1}`); err == nil {
		t.Fatal("缺少 ok 字段应报错")
	}
}

func TestNormalizeJSON(t *testing.T) {
	response, _ := parseResponse(`{"ok":true,"data":{"i":7,"f":1.5,"nested":[3,2.25],"s":"x"}}`)
	data := response["data"]
	normalized := normalizeJSON(data)
	mapping := normalized.(map[string]any)
	if value, ok := mapping["i"].(int64); !ok || value != 7 {
		t.Fatal("整数应归一化为 int64")
	}
	if value, ok := mapping["f"].(float64); !ok || value != 1.5 {
		t.Fatal("小数应归一化为 float64")
	}
	if value, ok := mapping["nested"].([]any); !ok || value[0] != int64(3) {
		t.Fatal("嵌套数组应递归归一化")
	}
}

func TestValidateJSONValue(t *testing.T) {
	if err := validateJSONValue(map[string]any{"a_b": []any{int64(1), "x"}}, 0); err != nil {
		t.Fatal(err)
	}
	if err := validateJSONValue(map[string]any{"BadKey": 1}, 0); err == nil {
		t.Fatal("非 snake_case 字段应报错")
	}
	if err := validateJSONValue(map[string]any{"a": math.NaN()}, 0); err == nil {
		t.Fatal("NaN 应报错")
	}
	if err := validateJSONValue(map[string]any{"a": math.Inf(1)}, 0); err == nil {
		t.Fatal("Infinity 应报错")
	}
	if err := validateJSONValue(map[string]any{"a": struct{}{}}, 0); err == nil {
		t.Fatal("不可序列化类型应报错")
	}
}

func TestJSONValueDateTime(t *testing.T) {
	naive := time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)
	converted, err := jsonValue(naive, 0)
	if err != nil {
		t.Fatal(err)
	}
	// replace(tzinfo=Asia/Shanghai)：墙钟不变。
	if text, ok := converted.(string); !ok || text != "2024-01-02T15:04:05+08:00" {
		t.Fatalf("无时区时间应按上海时区输出，得到 %v", converted)
	}
	zoned := time.Date(2024, 1, 2, 10, 0, 0, 0, time.FixedZone("UTC+9", 9*3600))
	converted, err = jsonValue(zoned, 0)
	if err != nil {
		t.Fatal(err)
	}
	if text, ok := converted.(string); !ok || text != "2024-01-02T10:00:00+09:00" {
		t.Fatalf("带时区时间应原样输出，得到 %v", converted)
	}
}

func TestJSONValueDropsNil(t *testing.T) {
	converted, err := jsonValue(map[string]any{"keep": int64(1), "drop": nil}, 0)
	if err != nil {
		t.Fatal(err)
	}
	mapping := converted.(map[string]any)
	if _, exists := mapping["drop"]; exists {
		t.Fatal("nil 键应被丢弃")
	}
	if _, exists := mapping["keep"]; !exists {
		t.Fatal("非 nil 键应保留")
	}
}

func TestAsInt64(t *testing.T) {
	if value, ok := asInt64(json.Number("42")); !ok || value != 42 {
		t.Fatal("json.Number 整数应可转换")
	}
	if _, ok := asInt64(json.Number("1.5")); ok {
		t.Fatal("json.Number 小数不应视为整数")
	}
	if _, ok := asInt64(true); ok {
		t.Fatal("布尔不应视为整数")
	}
	if _, ok := asInt64(1.5); ok {
		t.Fatal("浮点不应视为整数")
	}
	if value, ok := asInt64(int64(-7)); !ok || value != -7 {
		t.Fatal("负整数应可转换")
	}
}
