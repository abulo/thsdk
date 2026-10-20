package thsdk

// 测试助手：用脚本化 RuntimeBridge 模拟原生运行组件，使全部测试
// 不依赖真实的 hq.dylib / hq.so / hq.dll，可在任何平台与 CI 环境运行。

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

type unsafePointer = unsafe.Pointer

// jsonNumber 构造 json.Number。
func jsonNumber(value float64) json.Number {
	return json.Number(strconv.FormatFloat(value, 'f', -1, 64))
}

// strconvParseInt 是 strconv.ParseInt 的薄封装。
func strconvParseInt(text string, base, bits int) (int64, error) {
	return strconv.ParseInt(text, base, bits)
}

// scriptedBridge 记录收到的请求，并按方法名返回预设的 JSON 响应。
type scriptedBridge struct {
	responses map[string]string
	requests  []string
	buffers   [][]byte
}

// bridge 构造可注入 Client 的 RuntimeBridge。
func (s *scriptedBridge) bridge() *RuntimeBridge {
	return &RuntimeBridge{
		path: "stub",
		callAlloc: func(request string) unsafe.Pointer {
			s.requests = append(s.requests, request)
			response := s.responses[requestMethod(request)]
			if response == "" {
				response = `{"ok":false,"error":{"code":"unknown_error","message":"no scripted response"}}`
			}
			buffer := append([]byte(response), 0)
			s.buffers = append(s.buffers, buffer)
			return unsafe.Pointer(&buffer[0])
		},
		freeResult: func(unsafe.Pointer) {},
	}
}

// client 用脚本桥创建客户端（独立会话目录）。
func (s *scriptedBridge) client() *Client {
	return NewClient(tTempDir(), s.bridge())
}

// lastRequest 返回最近一次请求的 JSON 映射（数字归一化为 int64/float64，
// 与 parseResponse 语义一致）。
func (s *scriptedBridge) lastRequest() map[string]any {
	if len(s.requests) == 0 {
		return nil
	}
	decoder := json.NewDecoder(strings.NewReader(s.requests[len(s.requests)-1]))
	decoder.UseNumber()
	var payload map[string]any
	if err := decoder.Decode(&payload); err != nil {
		return nil
	}
	return normalizeJSON(payload).(map[string]any)
}

// requestMethod 从编码请求中提取方法名（请求为严格 JSON，method 在最前）。
func requestMethod(request string) string {
	marker := `"method":"`
	index := strings.Index(request, marker)
	if index < 0 {
		return ""
	}
	rest := request[index+len(marker):]
	end := strings.IndexByte(rest, '"')
	if end < 0 {
		return ""
	}
	return rest[:end]
}

// okResponse 构造 ok 响应。
func okResponse(data any) string {
	encoded, _ := json.Marshal(map[string]any{"ok": true, "data": data})
	return string(encoded)
}

// tTempDir 返回测试临时目录。
func tTempDir() string {
	dir, err := os.MkdirTemp("", "thsdk-go-test-")
	if err != nil {
		return "/tmp"
	}
	return dir
}
