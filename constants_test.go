package thsdk

import (
	"strings"
	"testing"
)

func TestAdjustValues(t *testing.T) {
	if AdjustNone != 0 || AdjustForward != 1 || AdjustBackward != 2 {
		t.Fatal("Adjust 枚举值必须为 0/1/2")
	}
}

func TestPeriodValues(t *testing.T) {
	expected := map[Period]int64{
		PeriodDay: 0x4000, PeriodWeek: 0x5001, PeriodMonth: 0x6001,
		PeriodQuarter: 0x6003, PeriodYear: 0x7001,
	}
	for period, value := range expected {
		if int64(period) != value {
			t.Fatalf("Period %d 期望 %#x", period, value)
		}
	}
}

func TestStringEnums(t *testing.T) {
	for _, value := range []string{
		string(CallAuctionPhaseOpening), string(CallAuctionPhaseClosing),
		string(IndexCallAuctionPhaseOpening), string(IndexCallAuctionPhaseClosing),
		string(ExtendedHoursSessionPreMarket), string(ExtendedHoursSessionPostMarket),
		string(OrderBookDepthFive),
		string(NewsFlashTraversalModePage), string(NewsFlashTraversalModeOlder), string(NewsFlashTraversalModeNewer),
		string(NewsTimelineHistorical), string(NewsTimelineRealtime),
		string(NewsTraversalModePage), string(NewsTraversalModeSince),
		string(NewsDirectionFirst), string(NewsDirectionOlder), string(NewsDirectionNewer),
		string(RealtimeSecuritySignalLimitUp),
	} {
		if strings.TrimSpace(value) == "" {
			t.Fatal("字符串枚举值不能为空")
		}
	}
}

func TestInfoTypeBasics(t *testing.T) {
	if InfoTypeBad != -1 {
		t.Fatal("InfoTypeBad 必须为 -1")
	}
	if InfoTypeStockInfo != 14339 {
		t.Fatal("InfoTypeStockInfo 必须为 14339")
	}
}

func TestConvenienceConstants(t *testing.T) {
	if periods["daily"] != 0x4000 {
		t.Fatal("daily 周期必须映射到 0x4000")
	}
	if periods["1m"] != 0x3001 {
		t.Fatal("1m 周期必须映射到 0x3001")
	}
	for key, value := range map[string]int64{
		"": 0, "none": 0, "pre": 1, "forward": 1, "post": 2, "backward": 2,
	} {
		if adjustments[key] != value {
			t.Fatalf("adjustment %q 期望 %d", key, value)
		}
	}
}
