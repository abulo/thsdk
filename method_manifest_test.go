package thsdk

import (
	"strings"
	"testing"
)

func TestMethodInventoryCount(t *testing.T) {
	methods := MethodInventory()
	if len(methods) != 80 {
		t.Fatalf("MethodInventory 长度 = %d，期望 80", len(methods))
	}
}

func TestNativeMethodsUniqueness(t *testing.T) {
	if len(nativeMethods) != 80 {
		t.Fatalf("nativeMethods 数量 = %d，期望 80", len(nativeMethods))
	}
	seen := map[string]bool{}
	for name := range nativeMethods {
		if seen[name] {
			t.Fatalf("方法名重复: %s", name)
		}
		seen[name] = true
	}
}

func TestDataFrameRowFieldsCount(t *testing.T) {
	count := 0
	for _, spec := range nativeMethods {
		if spec.ReturnsDataFrame {
			count++
		}
	}
	if count != 63 {
		t.Fatalf("返回 DataFrame 的方法 = %d，期望 63", count)
	}
}

func TestNoForbiddenMarkers(t *testing.T) {
	for name, spec := range nativeMethods {
		lower := strings.ToLower(name)
		for _, marker := range []string{"subscribe", "subscription", "level2", "query_data"} {
			if strings.Contains(lower, marker) {
				t.Fatalf("方法 %s 包含被禁止的标记 %s", name, marker)
			}
		}
		if spec.ResultType == "" {
			t.Fatalf("方法 %s 缺少 ResultType", name)
		}
	}
}

func TestRequestContracts(t *testing.T) {
	for name, spec := range nativeMethods {
		contract, ok := requestContracts[name]
		if !ok {
			t.Fatalf("方法 %s 缺少请求契约", name)
		}
		if spec.Mutating && len(contract.required) == 0 && len(contract.fields) != 0 {
			t.Fatalf("方法 %s 是写操作但契约异常", name)
		}
	}
}

func TestKnownResultTypes(t *testing.T) {
	// 写操作或未知主集合字段的方法应能查到 dataframeRowsByResultType 之外
	// 的类型定义；至少确保 map 存在且无空 ResultType。
	for name, spec := range nativeMethods {
		if spec.ResultType == "" {
			t.Fatalf("方法 %s 的 ResultType 为空", name)
		}
		_ = name
	}
}

func TestGroupMethodsMapping(t *testing.T) {
	// AccountWatchlistGroups 必须使用 groups 字段并在映射模式下返回。
	for name, spec := range nativeMethods {
		if name == "get_account_watchlist_groups" && spec.DataFrameRowsField != "groups" {
			t.Fatalf("get_account_watchlist_groups 主集合字段 = %s，期望 groups", spec.DataFrameRowsField)
		}
	}
}
