package thsdk

// runtime_test.go 覆盖运行组件路径解析的缓存/环境变量行为，
// 以及 readCString 分块读取的边界与防护（内存安全）。
//
// 说明：readCString 读取的是 C 分配内存（Go 运行时不做越界跟踪）；
// 测试用 Go []byte 模拟 C 缓冲区，因此每个用例都构造足够大的分配，
// 确保任何分块读取都不超出分配边界（避免 checkptr 误报）。

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"
)

// TestResolveComponentEnvOverride 验证 THSDK_RUNTIME_PATH 始终优先于目录布局查找，
// 且允许运行时切换（不命中进程级缓存）。
func TestResolveComponentEnvOverride(t *testing.T) {
	first := t.TempDir() + string(os.PathSeparator) + "hq.so"
	second := t.TempDir() + string(os.PathSeparator) + "hq.so"
	t.Setenv("THSDK_RUNTIME_PATH", first)
	got, err := resolveComponent()
	if err != nil {
		t.Fatalf("resolveComponent() error = %v", err)
	}
	if got != first {
		t.Errorf("环境变量优先 = %q, want %q", got, first)
	}
	t.Setenv("THSDK_RUNTIME_PATH", second)
	got, err = resolveComponent()
	if err != nil {
		t.Fatalf("resolveComponent() error = %v", err)
	}
	if got != second {
		t.Errorf("环境变量可运行时切换 = %q, want %q", got, second)
	}
}

// TestExpandComponentPath 验证路径展开（~ 与相对路径转绝对）。
func TestExpandComponentPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("无用户主目录: %v", err)
	}
	got, err := expandComponentPath("~/thsdk-test-libs/hq.so")
	if err != nil {
		t.Fatalf("expandComponentPath() error = %v", err)
	}
	want := filepath.Join(home, "thsdk-test-libs", "hq.so")
	if got != want {
		t.Errorf("expandComponentPath = %q, want %q", got, want)
	}

	absolute := filepath.Join(t.TempDir(), "hq.so")
	got, err = expandComponentPath(absolute)
	if err != nil {
		t.Fatalf("expandComponentPath() error = %v", err)
	}
	if got != filepath.Clean(absolute) {
		t.Errorf("expandComponentPath = %q, want %q", got, filepath.Clean(absolute))
	}
}

// bigBuffer 构造 len>=size 的字节缓冲区（末尾补零，避免越界读取）。
func bigBuffer(size int) []byte {
	return bytes.Repeat([]byte{0}, size)
}

// TestReadCStringBasic 验证 NUL 截断读取。
func TestReadCStringBasic(t *testing.T) {
	const maxBytes = 1 << 20
	payload := bigBuffer(maxBytes)
	copy(payload, "hello")
	payload[5] = 0
	got, err := readCString(unsafe.Pointer(&payload[0]), maxBytes)
	if err != nil {
		t.Fatalf("readCString() error = %v", err)
	}
	if got != "hello" {
		t.Errorf("readCString = %q, want %q", got, "hello")
	}
}

// TestReadCStringEmpty 验证以 NUL 开头返回空串。
func TestReadCStringEmpty(t *testing.T) {
	const maxBytes = 1 << 20
	payload := bigBuffer(maxBytes)
	copy(payload, "trailing")
	payload[0] = 0
	got, err := readCString(unsafe.Pointer(&payload[0]), maxBytes)
	if err != nil {
		t.Fatalf("readCString() error = %v", err)
	}
	if got != "" {
		t.Errorf("readCString = %q, want empty", got)
	}
}

// TestReadCStringNil 验证空指针安全返回。
func TestReadCStringNil(t *testing.T) {
	got, err := readCString(nil, 1024)
	if err != nil {
		t.Fatalf("readCString(nil) error = %v", err)
	}
	if got != "" {
		t.Errorf("readCString(nil) = %q, want empty", got)
	}
}

// TestReadCStringTooLong 验证超长数据触发防护错误（防野指针失控）。
func TestReadCStringTooLong(t *testing.T) {
	const maxBytes = 64
	payload := bytes.Repeat([]byte("a"), 128) // 分配必须 >= maxBytes，避免 checkptr 误报
	_, err := readCString(unsafe.Pointer(&payload[0]), maxBytes)
	if err == nil {
		t.Fatal("readCString 超过 maxBytes 应返回错误")
	}
	if !strings.Contains(err.Error(), "过大") {
		t.Errorf("错误信息应说明数据过大, got %q", err.Error())
	}
}

// TestReadCStringChunkBoundary 验证 NUL 恰好落在 4 KiB 分块边界上时仍正确截断。
func TestReadCStringChunkBoundary(t *testing.T) {
	const chunkSize = 4096
	payload := bigBuffer(chunkSize * 2)
	copy(payload, bytes.Repeat([]byte("a"), chunkSize))
	payload[chunkSize] = 0 // NUL 在分块边界
	got, err := readCString(unsafe.Pointer(&payload[0]), chunkSize*2)
	if err != nil {
		t.Fatalf("readCString() error = %v", err)
	}
	if got != strings.Repeat("a", chunkSize) {
		t.Errorf("分块边界截断不正确: len(got)=%d, want %d", len(got), chunkSize)
	}
}

// TestReadCStringMultiChunk 验证跨多个分块的连续读取。
func TestReadCStringMultiChunk(t *testing.T) {
	const chunkSize = 4096
	payload := bigBuffer(chunkSize * 8)
	contentLen := chunkSize*3 + 17
	copy(payload, bytes.Repeat([]byte("b"), contentLen))
	payload[contentLen] = 0
	got, err := readCString(unsafe.Pointer(&payload[0]), chunkSize*8)
	if err != nil {
		t.Fatalf("readCString() error = %v", err)
	}
	want := strings.Repeat("b", contentLen)
	if got != want {
		t.Errorf("多分块读取长度错误: len(got)=%d, want %d", len(got), len(want))
	}
}
