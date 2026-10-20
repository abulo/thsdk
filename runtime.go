package thsdk

// 本文件对应 Python 包的 thsdk/_runtime.py。
//
// _RuntimeBridge 负责加载当前平台的原生运行组件（hq.dylib / hq.so / hq.dll），
// 通过原生桥调用其 CallAlloc / FreeResult 完成严格的 JSON 请求/响应交换，
// 并在启动时执行版本与方法清单握手。

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
	"unsafe"
)

const (
	// abiVersion 对应 Python _ABI_VERSION = "2.0.0"。
	abiVersion = "2.0.0"
	// maxRequestBytes 对应 Python _MAX_REQUEST_BYTES = 1 << 20。
	maxRequestBytes = 1 << 20
	// maxMethodBytes 对应 Python _MAX_METHOD_BYTES = 128。
	maxMethodBytes = 128
	// maxJSONDepth 对应 Python _MAX_JSON_DEPTH = 64。
	maxJSONDepth = 64
	// maxCallTimeoutMS 对应 Python _MAX_CALL_TIMEOUT_MS = 600_000。
	maxCallTimeoutMS = 600_000
	// maxResponseBytes 是读取原生返回 C 字符串的长度保护上限（64 MiB）。
	// Python 通过 ctypes.string_at 无上限读取，这里增加防护避免野指针失控。
	maxResponseBytes = 1 << 26
)

// componentLayout 对应 Python _COMPONENT_LAYOUT：(系统, 文件名)。
var componentLayout = map[string][2]string{
	"darwin":  {"darwin", "hq.dylib"},
	"linux":   {"linux", "hq.so"},
	"windows": {"windows", "hq.dll"},
}

// RuntimeBridge 是进程内客户端的私有请求桥，对应 Python _RuntimeBridge。
// 原生调用层由构建标记选择：cgo 构建使用 bridge_cgo.go（专用原生线程承载
// dlopen 与 CallAlloc/FreeResult，规避 Go c-shared 库与宿主 Go 运行时的
// TLS 冲突）；CGO_ENABLED=0 构建使用 bridge_purego.go（purego 纯 Go FFI）。
type RuntimeBridge struct {
	path       string
	callAlloc  func(request string) unsafe.Pointer
	freeResult func(ptr unsafe.Pointer)
	lock       sync.Mutex
}

// NewRuntimeBridge 加载并校验原生运行组件，对应 Python _RuntimeBridge.__init__。
func NewRuntimeBridge(componentPath string) (*RuntimeBridge, error) {
	path, err := expandComponentPath(componentPath)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return nil, newTHSDKError("未找到适用于当前环境的 THSDK 运行组件，请重新安装匹配的发行包")
	}

	bridge := &RuntimeBridge{path: path}
	if err := bridge.initNative(); err != nil {
		return nil, err
	}

	version, err := bridge.handshakeCall("version", nil)
	if err != nil {
		return nil, err
	}
	runtimeVersion := ""
	if ok, _ := version["ok"].(bool); ok {
		if data, ok := version["data"].(map[string]any); ok {
			if value, ok := data["abi_version"].(string); ok {
				runtimeVersion = value
			}
		}
	}
	if runtimeVersion != abiVersion {
		return nil, newTHSDKError("THSDK 运行组件版本不兼容，请重新安装")
	}

	// ABI 2.0.0 早于严格请求解码；主动探测行为，避免旧组件静默接受
	// 拼写错误或未知请求字段。对应 Python 的 strict_probe。
	strictProbe, err := bridge.handshakeCall("version", map[string]any{"thsdk_protocol_probe": true})
	if err != nil {
		return nil, err
	}
	strictError, _ := strictProbe["error"].(map[string]any)
	if ok, _ := strictProbe["ok"].(bool); ok || strictError["code"] != "invalid_request" {
		return nil, newTHSDKError("THSDK 运行组件请求协议过旧，请安装最新版运行组件")
	}

	methods, err := bridge.handshakeCall("methods", nil)
	if err != nil {
		return nil, err
	}
	runtimeMethods, _ := methods["data"].([]any)
	runtimeOK, _ := methods["ok"].(bool)
	if !runtimeOK || !equalStringSlices(runtimeMethods, nativeMethodNames) {
		return nil, newTHSDKError("THSDK 运行组件方法清单不兼容，请重新安装")
	}

	return bridge, nil
}

// Path 返回已加载运行组件的绝对路径。
func (b *RuntimeBridge) Path() string {
	if b == nil {
		return ""
	}
	return b.path
}

// expandComponentPath 展开并规范化运行组件路径，对应 Python 构造器中的
// Path(component_path or self._resolve_component()).expanduser().resolve()。
func expandComponentPath(componentPath string) (string, error) {
	path := componentPath
	if path == "" {
		var err error
		path, err = resolveComponent()
		if err != nil {
			return "", err
		}
	}
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, strings.TrimPrefix(path, "~"))
		}
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(absolute), nil
}

// cachedComponentLayoutPath 缓存按「可执行文件/工作目录/源码目录」布局查找
// 运行组件的结果。进程内只需探测一次，避免每次 NewRuntimeBridge 重复多次
// 文件系统访问；THSDK_RUNTIME_PATH 环境变量仍每次实时读取（显式指定优先，
// 且允许调用方在运行时切换组件）。
var cachedComponentLayoutPath = sync.OnceValues(func() (string, error) {
	layout, ok := componentLayout[runtime.GOOS]
	if !ok {
		return "", newTHSDKError("当前环境暂不受此安装包支持，请安装匹配的 THSDK 发行包")
	}
	var architecture string
	switch runtime.GOARCH {
	case "arm64":
		architecture = "arm64"
	case "amd64":
		architecture = "amd64"
	default:
		return "", newTHSDKError("当前环境暂不受此安装包支持，请安装匹配的 THSDK 发行包")
	}
	relative := filepath.Join("libs", layout[0], architecture, layout[1])

	// 候选根目录：可执行文件目录 → 当前工作目录 → 本源码目录（仓库布局）。
	executable, err := os.Executable()
	if err == nil {
		if candidate := filepath.Join(filepath.Dir(executable), relative); fileExists(candidate) {
			return candidate, nil
		}
	}
	if workingDir, err := os.Getwd(); err == nil {
		if candidate := filepath.Join(workingDir, relative); fileExists(candidate) {
			return candidate, nil
		}
	}
	if _, sourceFile, _, ok := runtime.Caller(0); ok {
		if candidate := filepath.Join(filepath.Dir(sourceFile), relative); fileExists(candidate) {
			return candidate, nil
		}
	}
	return "", newTHSDKError("未找到适用于当前环境的 THSDK 运行组件，请重新安装匹配的发行包")
})

// resolveComponent 定位当前环境适用的运行组件，对应 Python _resolve_component：
// 优先 THSDK_RUNTIME_PATH 环境变量，否则使用缓存的目录布局查找结果。
func resolveComponent() (string, error) {
	if override := os.Getenv("THSDK_RUNTIME_PATH"); override != "" {
		return override, nil
	}
	return cachedComponentLayoutPath()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// handshakeCall 在启动阶段调用一次原生方法，只重试限流响应，
// 对应 Python _handshake_call。
func (b *RuntimeBridge) handshakeCall(method string, params any) (map[string]any, error) {
	response, err := b.Call(method, params, nil)
	if err != nil {
		return nil, err
	}
	errorValue, _ := response["error"].(map[string]any)
	if errorValue == nil {
		return response, nil
	}
	code, _ := errorValue["code"].(string)
	if retryAfterMs, ok := asInt64(errorValue["retry_after_ms"]); ok {
		if code == "rate_limited" && retryAfterMs > 0 && retryAfterMs <= 1000 {
			time.Sleep(time.Duration(retryAfterMs) * time.Millisecond)
			return b.Call(method, params, nil)
		}
	}
	return response, nil
}

// Call 执行一次原生方法调用，对应 Python _RuntimeBridge.call。
// params 为 nil 时请求体不含 params 字段；timeoutMs 为 nil 时不发送 timeout_ms。
func (b *RuntimeBridge) Call(method string, params any, timeoutMs *int) (map[string]any, error) {
	if method == "" || method != strings.TrimSpace(method) {
		return nil, newTHSDKError("THSDK method 必须是非空规范名称")
	}
	if len(method) > maxMethodBytes || !jsonFieldName.MatchString(method) {
		return nil, newTHSDKError("THSDK method 必须是不超过 128 字节的 ASCII 小写 snake_case")
	}
	if timeoutMs != nil {
		if *timeoutMs < 0 || *timeoutMs > maxCallTimeoutMS {
			return nil, newTHSDKError("timeout_ms 必须是 0（默认）或 1 到 600000 之间的整数")
		}
	}

	request := map[string]any{"method": method}
	if params != nil {
		request["params"] = params
	}
	if timeoutMs != nil {
		request["timeout_ms"] = *timeoutMs
	}
	if err := validateJSONValue(request, 0); err != nil {
		return nil, err
	}
	encoded, err := encodeStrictJSON(request)
	if err != nil {
		return nil, err
	}
	if len(encoded) > maxRequestBytes {
		return nil, newTHSDKError("THSDK 请求超过 1 MiB 限制")
	}

	b.lock.Lock()
	raw, err := b.callAllocated(encoded)
	b.lock.Unlock()
	if err != nil {
		return nil, err
	}
	response, err := parseResponse(raw)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// callAllocated 把编码后的请求交给原生 CallAlloc，读取返回的 C 字符串后
// 立即 FreeResult，对应 Python _call_allocated。
func (b *RuntimeBridge) callAllocated(encoded []byte) (string, error) {
	// purego 的 string 参数会自动补 NUL 终止符（与 Python ctypes.c_char_p 一致）。
	pointer := b.callAlloc(string(encoded))
	if pointer == nil {
		return "", newTHSDKError("THSDK 运行组件未返回结果")
	}
	defer b.freeResult(pointer)
	return readCString(pointer, maxResponseBytes)
}

// readCString 读取 NUL 结尾的 C 字符串，等价于 ctypes.string_at(pointer).decode("utf-8")。
//
// 实现分两步：先逐字节定位 NUL（真实 C 字符串保证 NUL 结尾，扫描在分配内
// 提前终止，对 checkptr 安全），再按实际长度一次性拷贝为 string，避免
// 边扫描边扩容拷贝。总长度超过 maxBytes 时返回错误（防野指针失控）。
func readCString(pointer unsafe.Pointer, maxBytes int) (string, error) {
	if pointer == nil {
		return "", nil
	}
	length := 0
	for length < maxBytes {
		if *(*byte)(unsafe.Add(pointer, length)) == 0 {
			break
		}
		length++
	}
	if length == maxBytes {
		return "", newTHSDKError("THSDK 运行组件返回数据过大")
	}
	payload := unsafe.Slice((*byte)(pointer), length)
	return string(payload), nil
}

// equalStringSlices 比较原生 methods 返回与清单是否完全一致。
func equalStringSlices(values []any, wanted []string) bool {
	if len(values) != len(wanted) {
		return false
	}
	for i, value := range values {
		if value != wanted[i] {
			return false
		}
	}
	return true
}
