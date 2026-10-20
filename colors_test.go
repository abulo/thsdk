package thsdk

// colors_test.go 验证公共颜色/样式常量的完整性与互斥性，
// 确保二维码终端渲染使用的 ANSI 转义码符合预期。

import "testing"

// TestAnsiColorConstants 验证 ANSI 常量与标准转义序列一致。
func TestAnsiColorConstants(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"ansiReset", ansiReset, "\x1b[0m"},
		{"ansiForegroundBlack", ansiForegroundBlack, "\x1b[30m"},
		{"ansiBackgroundWhite", ansiBackgroundWhite, "\x1b[47m"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}
}

// TestAnsiRenderComposition 验证二维码渲染的颜色组合序列（黑字白底+重置）。
func TestAnsiRenderComposition(t *testing.T) {
	want := "\x1b[30m\x1b[47m█\x1b[0m"
	got := ansiForegroundBlack + ansiBackgroundWhite + "█" + ansiReset
	if got != want {
		t.Errorf("渲染序列 = %q, want %q", got, want)
	}
}
