package thsdk

import (
	"testing"
	"time"
)

func TestNormalizeSecurityString(t *testing.T) {
	normalized, err := normalizeSecurity("USHA000001")
	if err != nil {
		t.Fatal(err)
	}
	if normalized["code"] != "000001" || normalized["market"] != "USHA" {
		t.Fatalf("规范化结果不符: %v", normalized)
	}
}

func TestNormalizeSecurityStringUppercase(t *testing.T) {
	normalized, err := normalizeSecurity("usza600000")
	if err != nil {
		t.Fatal(err)
	}
	if normalized["market"] != "USZA" {
		t.Fatalf("市场代码应大写: %v", normalized)
	}
}

func TestNormalizeSecurityMap(t *testing.T) {
	normalized, err := normalizeSecurity(map[string]any{"code": "000001", "market": "usza"})
	if err != nil {
		t.Fatal(err)
	}
	if normalized["market"] != "USZA" {
		t.Fatalf("映射市场代码应大写: %v", normalized)
	}
}

func TestNormalizeSecurityInvalid(t *testing.T) {
	for _, value := range []any{"", "123", "000001.SZ", 42, true, nil} {
		if _, err := normalizeSecurity(value); err == nil {
			t.Fatalf("非法证券 %v 应报错", value)
		}
	}
}

func TestNormalizeBlock(t *testing.T) {
	normalized, err := normalizeBlock(int64(12))
	if err != nil {
		t.Fatal(err)
	}
	if normalized["id"] != int64(12) {
		t.Fatalf("板块 ID 不符: %v", normalized)
	}
	if _, err := normalizeBlock(int64(0)); err == nil {
		t.Fatal("非正板块 ID 应报错")
	}
	if _, err := normalizeBlock(true); err == nil {
		t.Fatal("布尔板块应报错")
	}
}

func TestNormalizeNamedFields(t *testing.T) {
	payload, err := normalizeNamedFields(map[string]any{
		"security": "USHA000001",
		"markets":  "USZA, USHA ",
		"block":    int64(5),
	}, map[string]bool{"window_start": true})
	if err != nil {
		t.Fatal(err)
	}
	security := payload["security"].(map[string]any)
	if security["market"] != "USHA" {
		t.Fatalf("security 应被规范化: %v", security)
	}
	markets := payload["markets"].([]any)
	if len(markets) != 2 || markets[1] != "USHA" {
		t.Fatalf("markets 应拆分并大写: %v", markets)
	}
	if payload["block"].(map[string]any)["id"] != int64(5) {
		t.Fatalf("block 应被规范化: %v", payload["block"])
	}
}

func TestNormalizeNamedFieldsDatetime(t *testing.T) {
	start := time.Date(2024, 5, 6, 7, 8, 9, 0, time.UTC)
	payload, err := normalizeNamedFields(map[string]any{"window_start": start}, map[string]bool{"window_start": true})
	if err != nil {
		t.Fatal(err)
	}
	if payload["window_start"] != "2024-05-06T07:08:09+08:00" {
		t.Fatalf("datetime 字段应转为 ISO 字符串: %v", payload["window_start"])
	}
}

func TestNormalizeNamedFieldsInvalid(t *testing.T) {
	if _, err := normalizeNamedFields(map[string]any{"securities": "000001.SZ"}, nil); err == nil {
		t.Fatal("securities 传字符串应报错")
	}
	if _, err := normalizeNamedFields(map[string]any{"security": "bad"}, nil); err == nil {
		t.Fatal("非法 security 应报错")
	}
}

func TestToSequence(t *testing.T) {
	if _, ok := toSequence("text"); ok {
		t.Fatal("字符串不应视为序列")
	}
	if _, ok := toSequence([]byte("bytes")); ok {
		t.Fatal("字节串不应视为序列")
	}
	sequence, ok := toSequence([]any{int64(1), "a"})
	if !ok || len(sequence) != 2 {
		t.Fatal("[]any 应视为序列")
	}
	sequence, ok = toSequence([]string{"a", "b"})
	if !ok || len(sequence) != 2 {
		t.Fatal("[]string 应视为序列")
	}
}

func TestApplyTimeout(t *testing.T) {
	timeout, err := applyTimeout(nil)
	if err != nil || timeout != nil {
		t.Fatalf("无选项应返回 nil: %v %v", timeout, err)
	}
	timeout, err = applyTimeout([]*CallOptions{WithTimeout(3)})
	if err != nil {
		t.Fatal(err)
	}
	if timeout == nil || *timeout != 3000 {
		t.Fatalf("超时应转为毫秒: %v", timeout)
	}
	if _, err := applyTimeout([]*CallOptions{WithTimeout(601)}); err == nil {
		t.Fatal("超过 600 秒应报错")
	}
}

func TestIsAlpha(t *testing.T) {
	if !isAlpha("SZ") || isAlpha("S1") || isAlpha("") || !isAlpha("中文") {
		t.Fatal("isAlpha 行为不符")
	}
}
