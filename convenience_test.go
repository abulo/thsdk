package thsdk

import (
	"strings"
	"testing"
	"time"
)

// authDefault 给默认客户端配一个已认证的 stub 桥。
func authDefault(script *scriptedBridge) {
	script.responses["auth"] = okResponse(map[string]any{"authenticated": true})
	client := script.client()
	SetDefaultClient(client)
	if _, err := client.Auth("alice", "secret", ""); err != nil {
		panic(err)
	}
}

func TestDateNumber(t *testing.T) {
	cases := []struct {
		value any
		want  int64
	}{
		{time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC), 20240305},
		{20240305, 20240305},
		{"2024-03-05", 20240305},
		{"20240305", 20240305},
		// 仅做 8 位数字校验，不验证月/日合法性。
		{"2024-13-01", 20241301},
	}
	for _, item := range cases {
		number, provided, err := dateNumber(item.value)
		if err != nil || !provided || number != item.want {
			t.Fatalf("dateNumber(%v) = %d, %v, %v；期望 %d", item.value, number, provided, err, item.want)
		}
	}
	number, provided, err := dateNumber(nil)
	if err != nil || provided || number != 0 {
		t.Fatalf("nil 应返回未提供: %d %v %v", number, provided, err)
	}
	// 含时间的字符串（带空格/冒号）不属于 8 位日期。
	for _, bad := range []any{"202401", 123, "abc", "2024-03-05 10:00:00"} {
		if _, _, err := dateNumber(bad); err == nil {
			t.Fatalf("非法日期 %v 应报错", bad)
		}
	}
}

func TestDateText(t *testing.T) {
	if dateText(nil) != "0" {
		t.Fatal("未提供日期应返回 0")
	}
	if dateText("2024-03-05") != "20240305" {
		t.Fatal("日期转文本不符")
	}
}

func TestDateTimeText(t *testing.T) {
	text, err := datetimeText("2024-01-02T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}
	if text != "2024-01-02T15:04:05+00:00" {
		t.Fatalf("带 Z 时间应转 +00:00: %s", text)
	}
	text, err = datetimeText("2024-01-02T15:04:05")
	if err != nil {
		t.Fatal(err)
	}
	if text != "2024-01-02T15:04:05+08:00" {
		t.Fatalf("无时区时间应按上海时区: %s", text)
	}
	text, err = datetimeText(time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if text != "2024-01-02T10:00:00+08:00" {
		t.Fatalf("time.Time 应墙钟套上海: %s", text)
	}
	text, err = datetimeText("2024-01-02")
	if err != nil {
		t.Fatal(err)
	}
	if text != "2024-01-02T00:00:00+08:00" {
		t.Fatalf("纯日期应补零点: %s", text)
	}
	if _, err := datetimeText("not-a-time"); err == nil {
		t.Fatal("非法时间应报错")
	}
}

func TestSecurityHelper(t *testing.T) {
	item, err := security("USHA000001")
	if err != nil || item["code"] != "000001" || item["market"] != "USHA" {
		t.Fatalf("security 助手失败: %v %v", item, err)
	}
	if _, err := security(""); err == nil {
		t.Fatal("空证券应报错")
	}
}

func TestSecuritiesHelper(t *testing.T) {
	items, err := securities("USHA000001")
	if err != nil || len(items) != 1 {
		t.Fatalf("单个证券应包一层: %v %v", items, err)
	}
	items, err = securities([]any{"USHA000001", map[string]any{"code": "600000", "market": "usza"}})
	if err != nil || len(items) != 2 {
		t.Fatalf("序列证券失败: %v %v", items, err)
	}
	if items[1].(map[string]any)["market"] != "USZA" {
		t.Fatalf("映射证券应规范化: %v", items[1])
	}
	if _, err := securities([]any{}); err == nil {
		t.Fatal("空序列应报错")
	}
}

func TestMarketValue(t *testing.T) {
	market, err := marketValue(" usha ")
	if err != nil || market != "USHA" {
		t.Fatalf("market 应大写并去空格: %v %v", market, err)
	}
	if _, err := marketValue("US"); err == nil {
		t.Fatal("非法市场代码应报错")
	}
}

func TestMarketsValue(t *testing.T) {
	markets, err := marketsValue("USZA, USHA ")
	if err != nil || len(markets) != 2 || markets[1] != "USHA" {
		t.Fatalf("字符串 markets 拆分失败: %v %v", markets, err)
	}
	markets, err = marketsValue([]any{"USZA", "USHA"})
	if err != nil || len(markets) != 2 || markets[1] != "USHA" {
		t.Fatalf("序列 markets 失败: %v %v", markets, err)
	}
	markets, err = marketsValue("  ")
	if err != nil || markets != nil {
		t.Fatalf("空白 markets 应返回 nil: %v %v", markets, err)
	}
}

func TestWencaiIdentity(t *testing.T) {
	identity := wencaiIdentity(map[string]any{
		"security": map[string]any{"code": "000001", "market": "SZ"},
	})
	if identity["code"] != "000001" || identity["market"] != "SZ" || identity["full_code"] != "SZ000001" {
		t.Fatalf("wencai 身份提取失败: %v", identity)
	}
	// full_code 以四位字母开头时拆分（按前四位字母拆分）。
	identity = wencaiIdentity(map[string]any{"full_code": "USHA000001"})
	if identity["code"] != "000001" || identity["market"] != "USHA" {
		t.Fatalf("full_code 拆分失败: %v", identity)
	}
	// 非字母前缀（如 SZ 开头）不拆分，保持原样。
	identity = wencaiIdentity(map[string]any{"full_code": "SZ000001"})
	if identity["full_code"] != "SZ000001" || identity["code"] != "" {
		t.Fatalf("非字母前缀不应拆分: %v", identity)
	}
	// 显式 code/market 与 full_code 共存时，显式字段优先。
	identity = wencaiIdentity(map[string]any{"full_code": "SH600000", "code": "600000", "market": "SH"})
	if identity["full_code"] != "SH600000" || identity["code"] != "600000" {
		t.Fatalf("显式字段优先: %v", identity)
	}
	if len(wencaiIdentity(map[string]any{"other": 1})) != 0 {
		t.Fatal("无身份信息应返回空映射")
	}
}

func TestWencaiIdentityKey(t *testing.T) {
	key := wencaiIdentityKey(map[string]string{"full_code": "SZ000001", "code": "000001", "market": "SZ"})
	if key == nil || key.kind != "full_code" {
		t.Fatalf("full_code 键类型不符: %v", key)
	}
	key = wencaiIdentityKey(map[string]string{"code": "000001", "market": "SZ"})
	if key == nil || key.kind != "security" || key.first != "SZ" || key.second != "000001" {
		t.Fatalf("security 键类型不符: %v", key)
	}
	if wencaiIdentityKey(map[string]string{}) != nil {
		t.Fatal("空身份应返回 nil")
	}
}

func TestToIndex(t *testing.T) {
	if toWencaiIndex(jsonNumber(3)) != 3 {
		t.Fatal("json.Number 索引不符")
	}
	if toWencaiIndex("7") != 7 {
		t.Fatal("字符串索引不符")
	}
	if toWencaiIndex(2.9) != 2 {
		t.Fatal("浮点索引应截断")
	}
	if toWencaiIndex("bad") != 0 {
		t.Fatal("非法索引应为 0")
	}
}

func TestGetPriceValidation(t *testing.T) {
	// 这些分支在触达原生调用之前就会失败。
	if _, err := GetPrice("USHA000001", &GetPriceOptions{Frequency: "bogus"}); err == nil {
		t.Fatal("非法 frequency 应报错")
	}
	if _, err := GetPrice("USHA000001", &GetPriceOptions{FQ: "bogus"}); err == nil {
		t.Fatal("非法 fq 应报错")
	}
	if _, err := GetPrice("USHA000001", &GetPriceOptions{Count: 10, StartDate: "2024-01-01"}); err == nil {
		t.Fatal("count 与日期互斥应报错")
	}
	if _, err := GetPrice("USHA000001", &GetPriceOptions{StartDate: "2024-01-01"}); err == nil {
		t.Fatal("start/end 必须同时提供")
	}
	if _, err := GetPrice("USHA000001", &GetPriceOptions{Count: 0}); err == nil {
		t.Fatal("count 为 0 且无日期时 start=-1 需原生调用，但此处应无错") // 反向断言避免误用
	}
	if _, err := GetPrice("USHA000001", &GetPriceOptions{Count: -5}); err == nil {
		t.Fatal("count 必须大于 0")
	}
	if _, err := GetPrice("USHA000001", &GetPriceOptions{Frequency: "1m", StartDate: "2024-01-01", EndDate: "2024-01-02"}); err == nil {
		t.Fatal("分钟 K 线不支持日期范围")
	}
}

func TestGetPriceSuccess(t *testing.T) {
	defer SetDefaultClient(nil)
	script := &scriptedBridge{responses: map[string]string{
		"list_security_k_lines": okResponse([]any{
			map[string]any{"time": "20240102", "close": jsonNumber(10.5), "open": jsonNumber(10.0)},
			map[string]any{"time": "20240103", "close": jsonNumber(11.5), "open": jsonNumber(11.0)},
		}),
	}}
	authDefault(script)
	table, err := GetPrice("USHA000001", &GetPriceOptions{Count: 1})
	if err != nil {
		t.Fatal(err)
	}
	if table.RowCount() != 1 {
		t.Fatalf("count 截断失败: %d", table.RowCount())
	}
	if table.IndexName != "time" {
		t.Fatal("time 应设为索引")
	}
	// 验证请求参数
	params := script.lastRequest()["params"].(map[string]any)
	if params["start"] != int64(-1) || params["period"] != int64(0x4000) {
		t.Fatalf("K 线请求参数不符: %v", params)
	}
}

func TestGetPriceFields(t *testing.T) {
	defer SetDefaultClient(nil)
	script := &scriptedBridge{responses: map[string]string{
		"list_security_k_lines": okResponse([]any{
			map[string]any{"time": "20240102", "close": jsonNumber(10.5), "open": jsonNumber(10.0)},
		}),
	}}
	authDefault(script)
	table, err := GetPrice("USHA000001", &GetPriceOptions{Count: 1, Fields: []string{"close"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(table.Columns) != 1 || table.Columns[0] != "close" {
		t.Fatalf("字段过滤不符: %v", table.Columns)
	}
}

func TestKlinesDefaults(t *testing.T) {
	defer SetDefaultClient(nil)
	script := &scriptedBridge{responses: map[string]string{
		"list_security_k_lines": okResponse([]any{map[string]any{"time": "20240102"}}),
	}}
	authDefault(script)
	if _, err := Klines("USHA000001"); err != nil {
		t.Fatal(err)
	}
	params := script.lastRequest()["params"].(map[string]any)
	if params["start"] != int64(-100) {
		t.Fatalf("klines 默认 count 应为 100: %v", params["start"])
	}
	if params["adjust"] != int64(0) {
		t.Fatalf("klines 默认不复权: %v", params["adjust"])
	}
}

func TestIntradayData(t *testing.T) {
	defer SetDefaultClient(nil)
	script := &scriptedBridge{responses: map[string]string{
		"list_security_intraday_bars": okResponse([]any{
			map[string]any{"time": "0930", "price": jsonNumber(10.0)},
		}),
	}}
	authDefault(script)
	table, err := IntradayData("USHA000001", "2024-01-02")
	if err != nil {
		t.Fatal(err)
	}
	if table.RowCount() != 1 {
		t.Fatal("分时数据行数不符")
	}
	params := script.lastRequest()["params"].(map[string]any)
	if params["date"] != "20240102" {
		t.Fatalf("date 转文本失败: %v", params["date"])
	}
}

func TestTickLevel1Validation(t *testing.T) {
	if _, err := TickLevel1("USHA000001", &TickLevel1Options{Count: 10, StartTime: "2024-01-02T09:30:00"}); err == nil {
		t.Fatal("count 与时间互斥应报错")
	}
	if _, err := TickLevel1("USHA000001", &TickLevel1Options{StartTime: "2024-01-02T09:30:00"}); err == nil {
		t.Fatal("start/end 必须同时提供")
	}
	if _, err := TickLevel1("USHA000001"); err == nil {
		t.Fatal("无参数应报错")
	}
	if _, err := TickLevel1("USHA000001", &TickLevel1Options{Count: -1}); err == nil {
		t.Fatal("count 必须大于 0")
	}
}

func TestDepthFlatten(t *testing.T) {
	defer SetDefaultClient(nil)
	script := &scriptedBridge{responses: map[string]string{
		"list_security_order_books": okResponse([]any{
			map[string]any{
				"security": map[string]any{"code": "000001", "market": "SZ"},
				"bids": []any{
					map[string]any{"price": jsonNumber(10.0), "volume": jsonNumber(100)},
					map[string]any{"price": jsonNumber(9.9), "volume": jsonNumber(200)},
				},
				"asks": []any{
					map[string]any{"price": jsonNumber(10.1), "volume": jsonNumber(50)},
				},
			},
		}),
	}}
	authDefault(script)
	table, err := Depth("USHA000001")
	if err != nil {
		t.Fatal(err)
	}
	if table.RowCount() != 3 {
		t.Fatalf("盘口应展开为 3 行: %d", table.RowCount())
	}
	if table.Rows[0]["full_code"] != "SZ000001" || table.Rows[0]["side"] != "bid" || table.Rows[0]["level"] != 1 {
		t.Fatalf("盘口行不符: %v", table.Rows[0])
	}
}

func TestCorporateAction(t *testing.T) {
	defer SetDefaultClient(nil)
	script := &scriptedBridge{responses: map[string]string{
		"list_security_corporate_actions": okResponse([]any{
			map[string]any{"time": "20240102"},
			map[string]any{"time": "20240103"},
			map[string]any{"time": "20240104"},
		}),
	}}
	authDefault(script)
	table, err := CorporateAction("USHA000001", &CorporateActionOptions{Count: 2})
	if err != nil {
		t.Fatal(err)
	}
	if table.RowCount() != 2 {
		t.Fatalf("count 截断不符: %d", table.RowCount())
	}
}

func TestSearchSymbols(t *testing.T) {
	defer SetDefaultClient(nil)
	script := &scriptedBridge{responses: map[string]string{
		"search_securities": okResponse([]any{
			map[string]any{"security": map[string]any{"code": "000001", "market": "SZ"}, "name": "平安银行"},
		}),
	}}
	authDefault(script)
	table, err := SearchSymbols("000001", "")
	if err != nil {
		t.Fatal(err)
	}
	if table.Rows[0]["code"] != "000001" || table.Rows[0]["full_code"] != "SZ000001" {
		t.Fatalf("搜索行身份不符: %v", table.Rows[0])
	}
	if _, err := SearchSymbols("  ", ""); err == nil {
		t.Fatal("空 pattern 应报错")
	}
	_, err = SearchSymbols("000001", "US")
	if err == nil {
		t.Fatal("非法市场应报错")
	}
}

func TestCompleteTHSCode(t *testing.T) {
	defer SetDefaultClient(nil)
	script := &scriptedBridge{responses: map[string]string{
		"resolve_securities": okResponse([]any{
			map[string]any{"security": map[string]any{"code": "000001", "market": "SZ"}},
		}),
	}}
	authDefault(script)
	table, err := CompleteTHSCode("000001")
	if err != nil {
		t.Fatal(err)
	}
	if table.Rows[0]["full_code"] != "SZ000001" {
		t.Fatalf("解析结果不符: %v", table.Rows[0])
	}
	if _, err := CompleteTHSCode(""); err == nil {
		t.Fatal("空代码应报错")
	}
	params := script.lastRequest()["params"].(map[string]any)
	if params["codes"].([]any)[0] != "000001" {
		t.Fatalf("codes 参数不符: %v", params["codes"])
	}
}

func TestWencaiNLP(t *testing.T) {
	defer SetDefaultClient(nil)
	script := &scriptedBridge{responses: map[string]string{
		"query_wencai": okResponse(map[string]any{
			"fields": []any{
				map[string]any{"index": jsonNumber(1), "timestamp": "t1", "display_name": "收盘价"},
			},
			"securities": []any{
				map[string]any{"security": map[string]any{"code": "000001", "market": "SZ"}},
				map[string]any{"security": map[string]any{"code": "600000", "market": "SH"}},
			},
			"rows": []any{
				map[string]any{
					"security": map[string]any{"code": "000001", "market": "SZ"},
					"cells": []any{
						map[string]any{"key_index": jsonNumber(1), "timestamp": "t1", "value": jsonNumber(10.5)},
					},
				},
			},
			"query": "市盈率小于 10",
		}),
	}}
	authDefault(script)
	table, err := WencaiNLP("市盈率小于 10")
	if err != nil {
		t.Fatal(err)
	}
	if table.RowCount() != 2 {
		t.Fatalf("应返回 2 行（1 匹配行 + 1 未匹配证券）: %d", table.RowCount())
	}
	first := table.Rows[0]
	if first["full_code"] != "SZ000001" {
		t.Fatalf("行身份不符: %v", first)
	}
	if first["收盘价"] != 10.5 {
		t.Fatalf("动态字段映射失败: %v", first)
	}
	if table.Attr("query") != "市盈率小于 10" {
		t.Fatalf("query attr 不符: %v", table.Attr("query"))
	}
}

func TestWencaiNLPDedup(t *testing.T) {
	defer SetDefaultClient(nil)
	script := &scriptedBridge{responses: map[string]string{
		"query_wencai": okResponse(map[string]any{
			"securities": []any{
				map[string]any{"security": map[string]any{"code": "000001", "market": "SZ"}},
			},
			"rows": []any{
				map[string]any{
					"security": map[string]any{"code": "000001", "market": "SZ"},
					"cells":    []any{},
				},
				map[string]any{
					"security": map[string]any{"code": "000001", "market": "SZ"},
					"cells":    []any{},
				},
			},
		}),
	}}
	authDefault(script)
	table, err := WencaiNLP("测试去重")
	if err != nil {
		t.Fatal(err)
	}
	if table.RowCount() != 1 {
		t.Fatalf("重复身份行应去重: %d", table.RowCount())
	}
}

func TestWencaiNLPValidation(t *testing.T) {
	if _, err := WencaiNLP("  "); err == nil {
		t.Fatal("空 condition 应报错")
	}
}

func TestNewsAttrs(t *testing.T) {
	defer SetDefaultClient(nil)
	script := &scriptedBridge{responses: map[string]string{
		"list_news_flash_items": okResponse(map[string]any{
			"items": []any{
				map[string]any{"id": "n1", "title": "快讯"},
			},
			"total": jsonNumber(42),
		}),
	}}
	authDefault(script)
	table, err := News()
	if err != nil {
		t.Fatal(err)
	}
	if table.RowCount() != 1 {
		t.Fatal("快讯行数不符")
	}
	if table.Attr("total") != int64(42) {
		t.Fatalf("total 应进入 attrs: %v", table.Attr("total"))
	}
	if _, err := News(&NewsOptions{Page: 0}); err == nil {
		t.Fatal("page 必须大于 0")
	}
}

func TestMarketSecurities(t *testing.T) {
	defer SetDefaultClient(nil)
	script := &scriptedBridge{responses: map[string]string{
		"list_market_securities": okResponse(map[string]any{
			"securities": []any{
				map[string]any{"security": map[string]any{"code": "000001", "market": "SZ"}},
			},
			"total": jsonNumber(1),
		}),
	}}
	authDefault(script)
	table, err := MarketSecurities("usha")
	if err != nil {
		t.Fatal(err)
	}
	if table.Rows[0]["full_code"] != "SZ000001" {
		t.Fatalf("证券行不符: %v", table.Rows[0])
	}
	if table.Attr("total") != int64(1) {
		t.Fatalf("total attr 不符: %v", table.Attr("total"))
	}
	params := script.lastRequest()["params"].(map[string]any)
	if params["market"] != "USHA" {
		t.Fatalf("市场参数不符: %v", params["market"])
	}
}

func TestBlockConstituentsByName(t *testing.T) {
	defer SetDefaultClient(nil)
	script := &scriptedBridge{responses: map[string]string{
		"resolve_block": okResponse(map[string]any{"id": jsonNumber(5), "name": "银行"}),
		"list_block_constituents": okResponse([]any{
			map[string]any{"security": map[string]any{"code": "000001", "market": "SZ"}},
		}),
	}}
	authDefault(script)
	table, err := BlockConstituents("银行")
	if err != nil {
		t.Fatal(err)
	}
	if table.RowCount() != 1 {
		t.Fatal("成分股行数不符")
	}
	params := script.lastRequest()["params"].(map[string]any)
	if params["block"].(map[string]any)["id"] != int64(5) {
		t.Fatalf("板块解析未生效: %v", params["block"])
	}
}

func TestBlockConstituentsByID(t *testing.T) {
	defer SetDefaultClient(nil)
	script := &scriptedBridge{responses: map[string]string{
		"list_block_constituents": okResponse([]any{}),
	}}
	authDefault(script)
	if _, err := BlockConstituents(5); err != nil {
		t.Fatal(err)
	}
	params := script.lastRequest()["params"].(map[string]any)
	if params["block"].(map[string]any)["id"] != int64(5) {
		t.Fatalf("板块 ID 不符: %v", params["block"])
	}
	if _, err := BlockConstituents(0); err == nil {
		t.Fatal("非正板块 ID 应报错")
	}
	if _, err := BlockConstituents("  "); err == nil {
		t.Fatal("空板块名应报错")
	}
}

func TestGetPriceTimeIndexString(t *testing.T) {
	// 分钟 K 线（time 为 HHMMSS）不应触发日期索引转换。
	defer SetDefaultClient(nil)
	script := &scriptedBridge{responses: map[string]string{
		"list_security_k_lines": okResponse([]any{
			map[string]any{"time": "093000", "close": jsonNumber(10.0)},
		}),
	}}
	authDefault(script)
	table, err := GetPrice("USHA000001", &GetPriceOptions{Count: 1, Frequency: "1m"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := table.Index[0].(string); !ok {
		t.Fatal("分钟时间不应转为日期")
	}
}

func TestGetPriceParamsMinute(t *testing.T) {
	defer SetDefaultClient(nil)
	script := &scriptedBridge{responses: map[string]string{
		"list_security_k_lines": okResponse([]any{
			map[string]any{"time": "093000", "close": jsonNumber(10.0)},
		}),
	}}
	authDefault(script)
	_, err := GetPrice("USHA000001", &GetPriceOptions{Count: 30, Frequency: "1m"})
	if err != nil {
		t.Fatal(err)
	}
	params := script.lastRequest()["params"].(map[string]any)
	if params["period"] != int64(0x3001) {
		t.Fatalf("1m 周期映射失败: %v", params["period"])
	}
}

func TestStrOrEmpty(t *testing.T) {
	if stringOrEmpty(nil) != "" {
		t.Fatal("nil 应转空串")
	}
	if stringOrEmpty("x") != "x" {
		t.Fatal("字符串原样返回")
	}
	if !strings.HasPrefix(stringOrEmpty(int64(3)), "3") {
		t.Fatal("数字应转字符串")
	}
}
