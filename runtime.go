package thsdk

// 本文件实现运行时桥接层：
//
// RuntimeBridge 负责加载当前平台的原生运行组件（hq.dylib / hq.so / hq.dll），
// 通过原生桥调用其 CallAlloc / FreeResult 完成严格的 JSON 请求/响应交换，
// 并在启动时执行版本与方法清单握手。
//
// 原生调用层由构建标记选择：
//   - cgo 构建（默认，运行最可靠）：bridge_cgo.go + bridge.c，用专用原生
//     线程承载 dlopen 与 CallAlloc/FreeResult，规避 Go c-shared 库与宿主
//     Go 运行时在线程本地存储（TLS）上的冲突；
//   - CGO_ENABLED=0 构建：bridge_purego*.go（purego 纯 Go FFI），
//     保留用于无 C 工具链环境的编译验证。

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
	"unsafe"
)

// 协议层常量。
const (
	// abiVersion 是原生运行组件要求的 ABI 版本号，握手时校验。
	abiVersion = "2.0.0"
	// maxRequestBytes 是单个请求 JSON 的最大字节数（1 MiB），防止超大请求。
	maxRequestBytes = 1 << 20
	// maxMethodBytes 是方法名的最大字节数（128）。
	maxMethodBytes = 128
	// maxJSONDepth 是请求 JSON 的最大嵌套深度（64），防循环引用。
	maxJSONDepth = 64
	// maxCallTimeoutMS 是单次调用的最大超时毫秒数（600 秒）。
	maxCallTimeoutMS = 600_000
	// maxResponseBytes 是读取原生返回 C 字符串的长度保护上限（64 MiB）。
	// 防止原生侧野指针导致宿主读取失控。
	maxResponseBytes = 1 << 26
)

// componentLayout 记录各操作系统对应的原生组件目录与文件名。
//
// 键：runtime.GOOS 值；值：[子目录名, 组件文件名]。
var componentLayout = map[string][2]string{
	"darwin":  {"darwin", "hq.dylib"},
	"linux":   {"linux", "hq.so"},
	"windows": {"windows", "hq.dll"},
}

// RuntimeBridge 是进程内客户端的私有请求桥。
//
// 字段说明：
//   - path：已加载原生组件的绝对路径；
//   - callAlloc：向原生组件提交一次 JSON 请求并返回响应字符串指针的函数；
//   - freeResult：释放 callAlloc 返回的响应字符串指针；
//   - lock：串行化原生调用（原生桥为单线程作业模型，调用必须互斥）。
type RuntimeBridge struct {
	path       string
	callAlloc  func(request string) unsafe.Pointer
	freeResult func(ptr unsafe.Pointer)
	lock       sync.Mutex
}

// NewRuntimeBridge 加载并校验原生运行组件。
//
// 入参：
//   - componentPath：原生组件路径；传空串时按环境自动查找
//     （THSDK_RUNTIME_PATH → libs 目录布局）。
//
// 出参：
//   - *RuntimeBridge：校验通过后的运行桥实例；
//   - error：组件缺失、加载失败、版本不兼容或方法清单不匹配时返回错误。
//
// 校验流程：
//  1. 解析并确认组件文件存在；
//  2. 启动原生桥并加载组件（dlopen）与两个导出符号；
//  3. 握手 version：校验 ABI 版本一致；
//  4. 严格协议探测：旧组件会对未知字段静默放行，探测其是否支持
//     严格请求解码；
//  5. 握手 methods：校验原生方法清单与本包登记完全一致。
func NewRuntimeBridge(componentPath string) (*RuntimeBridge, error) {
	// 展开并规范化组件路径。
	path, err := expandComponentPath(componentPath)
	if err != nil {
		return nil, err
	}
	// 确认组件文件真实存在且不是目录。
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return nil, newTHSDKError("未找到适用于当前环境的 THSDK 运行组件，请重新安装匹配的发行包")
	}

	// 创建桥实例并初始化原生加载（由构建标记选择实现）。
	bridge := &RuntimeBridge{path: path}
	if err := bridge.initNative(); err != nil {
		return nil, err
	}

	// 握手 1：查询组件版本，校验 ABI 版本一致。
	version, err := bridge.handshakeCall("version", nil)
	if err != nil {
		return nil, err
	}
	// 从响应的 data.abi_version 字段取出组件 ABI 版本。
	runtimeVersion := ""
	if ok, _ := version["ok"].(bool); ok {
		if data, ok := version["data"].(map[string]any); ok {
			if value, ok := data["abi_version"].(string); ok {
				runtimeVersion = value
			}
		}
	}
	// ABI 版本不一致视为组件过旧或错误发行包。
	if runtimeVersion != abiVersion {
		return nil, newTHSDKError("THSDK 运行组件版本不兼容，请重新安装")
	}

	// 握手 2：严格协议探测。
	// 早于严格请求解码的组件会静默接受拼写错误或未知请求字段；
	// 这里主动发送探测参数，要求组件返回 invalid_request 错误，
	// 从而识别并拒绝过旧组件。
	strictProbe, err := bridge.handshakeCall("version", map[string]any{"thsdk_protocol_probe": true})
	if err != nil {
		return nil, err
	}
	strictError, _ := strictProbe["error"].(map[string]any)
	// 期望得到 invalid_request 错误；若返回 ok 或非该错误码，说明协议过旧。
	if ok, _ := strictProbe["ok"].(bool); ok || strictError["code"] != "invalid_request" {
		return nil, newTHSDKError("THSDK 运行组件请求协议过旧，请安装最新版运行组件")
	}

	// 握手 3：查询原生方法清单，必须与本包登记清单完全一致。
	methods, err := bridge.handshakeCall("methods", nil)
	if err != nil {
		return nil, err
	}
	runtimeMethods, _ := methods["data"].([]any)
	runtimeOK, _ := methods["ok"].(bool)
	if !runtimeOK || !equalStringSlices(runtimeMethods, nativeMethodNames) {
		return nil, newTHSDKError("THSDK 运行组件方法清单不兼容，请重新安装")
	}

	// 全部校验通过，返回可用桥实例。
	return bridge, nil
}

// Path 返回已加载运行组件的绝对路径。
//
// 入参：无。
//
// 出参：
//   - string：组件绝对路径；接收者为 nil 时返回空串。
func (b *RuntimeBridge) Path() string {
	if b == nil {
		return ""
	}
	return b.path
}

// expandComponentPath 展开并规范化运行组件路径。
//
// 入参：
//   - componentPath：原始路径；为空时自动解析。
//
// 出参：
//   - string：规范化后的绝对路径；
//   - error：自动解析失败时返回错误。
//
// 处理规则：
//   - 空路径 → resolveComponent() 自动定位；
//   - 以 ~ 开头 → 展开为用户主目录；
//   - 相对路径 → 转为绝对路径并清理冗余成分。
func expandComponentPath(componentPath string) (string, error) {
	path := componentPath
	if path == "" {
		// 未指定时按环境自动解析。
		var err error
		path, err = resolveComponent()
		if err != nil {
			return "", err
		}
	}
	// 展开 ~ 前缀为用户主目录。
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, strings.TrimPrefix(path, "~"))
		}
	}
	// 转为绝对路径并清理。
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(absolute), nil
}

// cachedComponentLayoutPath 缓存按「可执行文件/工作目录/源码目录」布局查找
// 运行组件的结果。
//
// 设计目的：进程内只需探测一次，避免每次 NewRuntimeBridge 重复多次文件
// 系统访问；THSDK_RUNTIME_PATH 环境变量仍每次实时读取（显式指定优先，
// 且允许调用方在运行时切换组件）。
var cachedComponentLayoutPath = sync.OnceValues(func() (string, error) {
	// 按当前操作系统查组件文件名。
	layout, ok := componentLayout[runtime.GOOS]
	if !ok {
		return "", newTHSDKError("当前环境暂不受此安装包支持，请安装匹配的 THSDK 发行包")
	}
	// 按当前 CPU 架构选择子目录。
	var architecture string
	switch runtime.GOARCH {
	case "arm64":
		architecture = "arm64"
	case "amd64":
		architecture = "amd64"
	default:
		return "", newTHSDKError("当前环境暂不受此安装包支持，请安装匹配的 THSDK 发行包")
	}
	// 组件相对路径：libs/<系统>/<架构>/<文件名>。
	relative := filepath.Join("libs", layout[0], architecture, layout[1])

	// 候选根目录依次为：可执行文件目录 → 当前工作目录 → 本源码目录（仓库布局）。
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
	// 所有候选位置都未找到组件。
	return "", newTHSDKError("未找到适用于当前环境的 THSDK 运行组件，请重新安装匹配的发行包")
})

// resolveComponent 定位当前环境适用的运行组件。
//
// 入参：无。
//
// 出参：
//   - string：组件绝对路径；
//   - error：找不到组件时返回错误。
//
// 查找优先级：
//  1. THSDK_RUNTIME_PATH 环境变量（每次实时读取，可运行时切换）；
//  2. 目录布局缓存结果（cachedComponentLayoutPath）。
func resolveComponent() (string, error) {
	// 环境变量显式指定时优先，且不做缓存。
	if override := os.Getenv("THSDK_RUNTIME_PATH"); override != "" {
		return override, nil
	}
	return cachedComponentLayoutPath()
}

// fileExists 判断路径是否为存在的常规文件。
//
// 入参：
//   - path：待检查路径。
//
// 出参：
//   - bool：true 表示存在且不是目录。
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// handshakeCall 在启动阶段调用一次原生方法，只对限流响应做一次重试。
//
// 入参：
//   - method：原生方法名；
//   - params：请求参数；可为 nil。
//
// 出参：
//   - map[string]any：原生响应结构；
//   - error：调用失败时返回错误。
//
// 说明：握手阶段偶发限流（rate_limited）且等待时间 ≤ 1 秒时，
// 等待后重试一次；其余错误原样返回。
func (b *RuntimeBridge) handshakeCall(method string, params any) (map[string]any, error) {
	response, err := b.Call(method, params, nil)
	if err != nil {
		return nil, err
	}
	// 无错误字段则直接返回。
	errorValue, _ := response["error"].(map[string]any)
	if errorValue == nil {
		return response, nil
	}
	// 命中限流且等待时间在 1 秒内：睡眠后重试一次。
	code, _ := errorValue["code"].(string)
	if retryAfterMs, ok := asInt64(errorValue["retry_after_ms"]); ok {
		if code == "rate_limited" && retryAfterMs > 0 && retryAfterMs <= 1000 {
			time.Sleep(time.Duration(retryAfterMs) * time.Millisecond)
			return b.Call(method, params, nil)
		}
	}
	return response, nil
}

// Call 执行一次原生方法调用。
//
// 入参：
//   - method：原生方法名，必须是非空且不超过 128 字节的 ASCII 小写
//     snake_case；
//   - params：请求参数；为 nil 时请求体不含 params 字段；
//   - timeoutMs：超时毫秒数指针；为 nil 时不发送 timeout_ms 字段，
//     否则取值必须为 0（默认）或 1 到 600000 之间。
//
// 出参：
//   - map[string]any：原生响应结构（已校验含 ok 字段）；
//   - error：参数非法、请求过大或原生返回无效数据时返回错误。
//
// 说明：调用全程持有互斥锁，保证与原生单线程桥的串行模型一致。
func (b *RuntimeBridge) Call(method string, params any, timeoutMs *int) (map[string]any, error) {
	// 方法名基本校验：非空、无首尾空白。
	if method == "" || method != strings.TrimSpace(method) {
		return nil, newTHSDKError("THSDK method 必须是非空规范名称")
	}
	// 方法名格式与长度校验。
	if len(method) > maxMethodBytes || !jsonFieldName.MatchString(method) {
		return nil, newTHSDKError("THSDK method 必须是不超过 128 字节的 ASCII 小写 snake_case")
	}
	// 超时值范围校验。
	if timeoutMs != nil {
		if *timeoutMs < 0 || *timeoutMs > maxCallTimeoutMS {
			return nil, newTHSDKError("timeout_ms 必须是 0（默认）或 1 到 600000 之间的整数")
		}
	}

	// 组装请求结构。
	request := map[string]any{"method": method}
	if params != nil {
		request["params"] = params
	}
	if timeoutMs != nil {
		request["timeout_ms"] = *timeoutMs
	}
	// 严格校验请求 JSON 规则（字段名/深度/数值）。
	if err := validateJSONValue(request, 0); err != nil {
		return nil, err
	}
	// 严格序列化为 JSON 字节。
	encoded, err := encodeStrictJSON(request)
	if err != nil {
		return nil, err
	}
	// 请求体积上限保护。
	if len(encoded) > maxRequestBytes {
		return nil, newTHSDKError("THSDK 请求超过 1 MiB 限制")
	}

	// 加锁执行原生调用（串行模型）。
	b.lock.Lock()
	raw, err := b.callAllocated(encoded)
	b.lock.Unlock()
	if err != nil {
		return nil, err
	}
	// 解析响应并校验基本结构。
	response, err := parseResponse(raw)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// callAllocated 把编码后的请求交给原生 CallAlloc，读取返回的 C 字符串后
// 立即 FreeResult。
//
// 入参：
//   - encoded：已序列化的请求 JSON 字节。
//
// 出参：
//   - string：原生返回的响应文本；
//   - error：原生未返回结果时返回错误。
//
// 说明：string 参数在原生桥层会自动补 NUL 终止符，与 C 字符串约定一致；
// 响应指针读取后立刻释放，保证不泄漏。
func (b *RuntimeBridge) callAllocated(encoded []byte) (string, error) {
	pointer := b.callAlloc(string(encoded))
	if pointer == nil {
		return "", newTHSDKError("THSDK 运行组件未返回结果")
	}
	defer b.freeResult(pointer)
	return readCString(pointer, maxResponseBytes)
}

// readCString 读取 NUL 结尾的 C 字符串。
//
// 入参：
//   - pointer：C 字符串起始地址；可为空指针。
//   - maxBytes：最大读取字节数，防止野指针失控。
//
// 出参：
//   - string：读到的内容（不含结尾 NUL）；
//   - error：超过 maxBytes 仍未见 NUL 时返回错误。
//
// 实现说明：分两步完成——先逐字节定位 NUL（真实 C 字符串保证 NUL 结尾，
// 扫描在分配内提前终止，对 checkptr 安全），再按实际长度一次性拷贝为
// string，避免边扫描边扩容拷贝。
func readCString(pointer unsafe.Pointer, maxBytes int) (string, error) {
	// 空指针直接返回空串。
	if pointer == nil {
		return "", nil
	}
	// 逐字节扫描定位 NUL，同时受 maxBytes 保护。
	length := 0
	for length < maxBytes {
		if *(*byte)(unsafe.Add(pointer, length)) == 0 {
			break
		}
		length++
	}
	// 扫描到上限仍未遇到 NUL，视为数据过大（防护）。
	if length == maxBytes {
		return "", newTHSDKError("THSDK 运行组件返回数据过大")
	}
	// 按实际长度一次性转换为 string。
	payload := unsafe.Slice((*byte)(pointer), length)
	return string(payload), nil
}

// equalStringSlices 比较原生 methods 返回与登记清单是否完全一致。
//
// 入参：
//   - values：原生返回的方法名列表（元素为 string）；
//   - wanted：本包登记的期望方法名列表。
//
// 出参：
//   - bool：长度与逐项内容完全一致返回 true。
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
