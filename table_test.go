package thsdk

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestRecordsFrameList(t *testing.T) {
	table, err := recordsFrame([]any{
		map[string]any{"code": "000001", "close": 10.0},
		map[string]any{"code": "000002", "close": 20.0, "volume": 100},
	}, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(table.Columns) != 3 {
		t.Fatalf("列应为 3: %v", table.Columns)
	}
	// 单记录内键按字典序排序（Go map 迭代无序，保证确定性）。
	if table.Columns[0] != "close" {
		t.Fatalf("列顺序应为记录内字典序: %v", table.Columns)
	}
	if table.Rows[1]["volume"] != 100 || table.Rows[0]["volume"] != nil {
		t.Fatal("缺失列应为 nil")
	}
}

func TestRecordsFrameNil(t *testing.T) {
	table, err := recordsFrame(nil, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !table.Empty() {
		t.Fatal("nil 输入应为空表")
	}
}

func TestRecordsFrameMappingKeys(t *testing.T) {
	table, err := recordsFrame(map[string]any{
		"a": map[string]any{"x": 1},
		"b": map[string]any{"x": 2},
	}, true, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(table.Rows) != 2 || table.IndexName != "group_key" {
		t.Fatalf("映射表格不符: %v %v", table.Rows, table.IndexName)
	}
	if table.Index[0] != "a" {
		t.Fatalf("索引顺序不符: %v", table.Index)
	}
}

func TestRecordsFrameMappingOrder(t *testing.T) {
	table, err := recordsFrame(map[string]any{
		"a": map[string]any{"x": 1},
		"b": map[string]any{"x": 2},
		"c": map[string]any{"x": 3},
	}, true, []any{"c", "a"})
	if err != nil {
		t.Fatal(err)
	}
	if table.Index[0] != "c" || table.Index[1] != "a" || table.Index[2] != "b" {
		t.Fatalf("group_order 排序不符: %v", table.Index)
	}
}

func TestFrameTimeIndex(t *testing.T) {
	table, err := frame([]any{
		map[string]any{"time": "20240102", "close": 10.0},
		map[string]any{"time": "20240103", "close": 11.0},
	}, true)
	if err != nil {
		t.Fatal(err)
	}
	if table.IndexName != "time" {
		t.Fatalf("time 应设为索引: %v", table.IndexName)
	}
	if parsed, ok := table.Index[0].(time.Time); !ok || parsed.Year() != 2024 {
		t.Fatalf("time 列应转为时间: %v", table.Index[0])
	}
	if table.HasColumn("time") {
		t.Fatal("time 应已从数据列移除")
	}
}

func TestFrameMixedTimeKeepsString(t *testing.T) {
	table, err := frame([]any{
		map[string]any{"time": "20240102"},
		map[string]any{"time": "not-a-date"},
	}, true)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := table.Index[1].(string); !ok {
		t.Fatal("非全日期时 time 应保持字符串")
	}
}

func TestIdentityRows(t *testing.T) {
	rows := identityRows([]any{
		map[string]any{"security": map[string]any{"code": "000001", "market": "SZ"}, "price": 10.0},
		map[string]any{"code": "600000", "market": "SH"},
	})
	if len(rows) != 2 {
		t.Fatal("行数不符")
	}
	first := rows[0]
	if first["code"] != "000001" || first["market"] != "SZ" || first["full_code"] != "SZ000001" {
		t.Fatalf("身份展平不符: %v", first)
	}
	if _, exists := first["security"]; exists {
		t.Fatal("security 嵌套应被移除")
	}
	// 没有嵌套 security 的行不补 full_code。
	if _, exists := rows[1]["full_code"]; exists {
		t.Fatalf("无嵌套 security 的行不应补 full_code: %v", rows[1])
	}
}

func TestSelect(t *testing.T) {
	table, _ := recordsFrame([]any{
		map[string]any{"code": "000001", "close": 10.0, "open": 9.5},
	}, false, nil)
	selected, err := table.Select([]string{"code", "close"})
	if err != nil {
		t.Fatal(err)
	}
	if len(selected.Columns) != 2 || len(selected.Rows[0]) != 2 {
		t.Fatalf("选取失败: %v", selected.Columns)
	}
	if _, err := table.Select([]string{"missing"}); err == nil {
		t.Fatal("缺失字段应报错")
	}
}

func TestTableJSON(t *testing.T) {
	table, _ := recordsFrame([]any{map[string]any{"code": "000001"}}, false, nil)
	table.SetAttr("query", "测试")
	data, err := table.JSON()
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded["columns"].([]any)[0] != "code" {
		t.Fatalf("JSON 结构不符: %s", data)
	}
	if !strings.Contains(string(data), "测试") {
		t.Fatal("attrs 应包含中文")
	}
}

func TestDataframeResultEnvelope(t *testing.T) {
	result := map[string]any{
		"securities": []any{
			map[string]any{"code": "000001", "market": "SZ"},
		},
		"total": int64(1),
	}
	table, err := tabularResult("list_market_securities", result)
	if err != nil {
		t.Fatal(err)
	}
	if len(table.Rows) != 1 {
		t.Fatal("主集合行数不符")
	}
	metadata := table.Attr("thsdk").(map[string]any)
	if metadata["method"] != "list_market_securities" {
		t.Fatalf("thsdk 元数据不符: %v", metadata)
	}
	inner := metadata["metadata"].(map[string]any)
	if inner["total"] != int64(1) {
		t.Fatalf("非主集合字段应进入 metadata: %v", inner)
	}
}

func TestDataframeResultList(t *testing.T) {
	table, err := tabularResult("list_security_k_lines", []any{
		map[string]any{"time": "20240102", "close": 10.0},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(table.Rows) != 1 {
		t.Fatal("list[T] 应直接作为行")
	}
}
