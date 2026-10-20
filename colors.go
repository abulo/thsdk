package thsdk

// 本文件集中定义终端渲染使用的 ANSI 转义码常量。
//
// 当前唯一使用方是二维码终端渲染（qr_display.go 的 writeTerminalQR）。
// 集中到公共文件的目的：颜色语义统一维护、避免各渲染模块散落魔法数字、
// 便于后续扩展（如深浅色主题切换、其他终端输出场景复用）。

// ANSI 转义码常量。
const (
	// ansiReset 重置所有样式（ESC[0m）。
	ansiReset = "\x1b[0m"
	// ansiForegroundBlack 设置前景色为黑色（ESC[30m）。
	ansiForegroundBlack = "\x1b[30m"
	// ansiBackgroundWhite 设置背景色为白色（ESC[47m）。
	ansiBackgroundWhite = "\x1b[47m"
)
