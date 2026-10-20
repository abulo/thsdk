package thsdk

// 本文件实现二维码挑战的终端显示。
//
// 尽力而为地显示一次二维码挑战：下载 HTTPS PNG 图片、手工解码 PNG
// （zlib 解压 + 滤镜还原 + 亮度计算）、采样二维码模块并渲染到终端。
// 显示失败绝不影响认证流程本身。

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	pngSignature     = "\x89PNG\r\n\x1a\n"
	maxImageBytes    = 2 * 1024 * 1024 // 下载图片体积上限（2 MiB）
	maxImageSide     = 1024            // 图片边长上限（1024 像素）
	maxRenderSeconds = 5.0             // 渲染超时上限（5 秒）
	quietZone        = 4               // 二维码静区宽度（4 模块）
)

// printQRChallenge 尽力显示一次二维码挑战。
// stream 为空时写入标准输出。
// printQRChallenge 尽力显示一次二维码挑战。
// stream 为空时写入标准输出。
// 入参：
//   - challenge：入参说明见函数用途。
//   - renderTimeout：入参说明见函数用途。
//   - stream：入参说明见函数用途。
func printQRChallenge(challenge map[string]any, renderTimeout float64, stream io.Writer) {
	output := stream
	if output == nil {
		output = os.Stdout
	}
	imageURL := challengeString(challenge, "image_url")
	scanURL := challengeString(challenge, "scan_url")

	writeLine(output, "请使用同花顺手机客户端扫码确认登录：", "Scan with the mobile app to confirm login:")
	flushStream(output)

	rendered := false
	if imageURL != "" && renderTimeout > 0 {
		modules, err := loadQRModules(imageURL, math.Min(renderTimeout, maxRenderSeconds))
		if err == nil && modules != nil {
			rendered = writeTerminalQR(modules, output)
		}
	}
	if rendered {
		if scanURL != "" {
			writeLine(output, "扫码地址："+scanURL, "Scan URL: "+scanURL)
		} else if imageURL != "" {
			writeLine(output, "二维码图片："+imageURL, "QR image: "+imageURL)
		}
	} else if imageURL != "" {
		writeLine(output, "当前终端未能显示二维码，请打开二维码图片地址。", "The QR code could not be rendered; open the QR image URL instead.")
		writeLine(output, "二维码图片："+imageURL, "QR image: "+imageURL)
	} else if scanURL != "" {
		writeLine(output, "扫码地址："+scanURL, "Scan URL: "+scanURL)
	}
	flushStream(output)
}

// challengeString 取字典字段并转为字符串；缺失或 nil 返回空串。
// str(value.get(key) or "")。
func challengeString(value map[string]any, key string) string {
	item := value[key]
	if item == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(item))
}

func writeLine(stream io.Writer, value, asciiFallback string) bool {
	if _, err := fmt.Fprintln(stream, value); err == nil {
		return true
	}
	if _, err := fmt.Fprintln(stream, asciiFallback); err == nil {
		return true
	}
	return false
}

// 入参：
//   - stream：入参说明见函数用途。
func flushStream(stream io.Writer) {
	if flusher, ok := stream.(interface{ Flush() error }); ok {
		_ = flusher.Flush()
		return
	}
	if file, ok := stream.(*os.File); ok {
		_ = file.Sync()
	}
}

// loadQRModules 带超时地加载二维码模块。
// 即使服务器一直保持 socket 打开也不会无限阻塞。
// loadQRModules 带超时地加载二维码模块。
// 即使服务器一直保持 socket 打开也不会无限阻塞。
// 入参：
//   - imageURL：入参说明见函数用途。
//   - timeout：入参说明见函数用途。
//
// 出参：[][]bool,error。
func loadQRModules(imageURL string, timeout float64) ([][]bool, error) {
	if timeout <= 0 {
		return nil, fmt.Errorf("二维码渲染超时")
	}
	type result struct {
		modules [][]bool
		err     error
	}
	channel := make(chan result, 1)
	go func() {
		data, err := downloadPNG(imageURL)
		if err != nil {
			channel <- result{err: err}
			return
		}
		pixels, err := decodePNG(data)
		if err != nil {
			channel <- result{err: err}
			return
		}
		modules, err := qrModules(pixels)
		channel <- result{modules: modules, err: err}
	}()
	select {
	case <-time.After(time.Duration(timeout * float64(time.Second))):
		return nil, fmt.Errorf("二维码渲染超时")
	case outcome := <-channel:
		if outcome.err != nil {
			return nil, outcome.err
		}
		if outcome.modules == nil {
			return nil, fmt.Errorf("二维码模块为空")
		}
		return outcome.modules, nil
	}
}

// downloadPNG 下载二维码 PNG 图片。
// downloadPNG 下载二维码 PNG 图片。
// 入参：
//   - imageURL：入参说明见函数用途。
//
// 出参：[]byte,error。
func downloadPNG(imageURL string) ([]byte, error) {
	parsed, err := url.Parse(imageURL)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
		return nil, fmt.Errorf("二维码图片地址必须使用 HTTPS")
	}
	client := &http.Client{Timeout: 5 * time.Second}
	request, err := http.NewRequest(http.MethodGet, imageURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "THSDK/2 QR Login")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if final := response.Request.URL; final == nil || final.Scheme != "https" {
		return nil, fmt.Errorf("二维码图片重定向必须使用 HTTPS")
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, maxImageBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxImageBytes {
		return nil, fmt.Errorf("二维码图片过大")
	}
	return data, nil
}

// qrModules 从解码像素中识别二维码模块。
// qrModules 从解码像素中识别二维码模块。
// 入参：
//   - pixels：入参说明见函数用途。
//
// 出参：[][]bool,error。
func qrModules(pixels [][]bool) ([][]bool, error) {
	imageHeight := len(pixels)
	imageWidth := len(pixels[0])

	left, top := imageWidth, imageHeight
	right, bottom := -1, -1
	for y, row := range pixels {
		for x, isBlack := range row {
			if !isBlack {
				continue
			}
			if x < left {
				left = x
			}
			if y < top {
				top = y
			}
			if x > right {
				right = x
			}
			if y > bottom {
				bottom = y
			}
		}
	}
	if right < 0 {
		return nil, fmt.Errorf("二维码图片没有有效图案")
	}
	width := right - left + 1
	height := bottom - top + 1

	type candidate struct {
		score  float64
		scale  int
		matrix [][]bool
	}
	var candidates []candidate
	for size := 21; size < 178; size += 4 {
		scaleX := float64(width) / float64(size)
		scaleY := float64(height) / float64(size)
		scale := int(math.Round((scaleX + scaleY) / 2))
		if scale < 1 || math.Abs(scaleX-float64(scale)) > 0.08 || math.Abs(scaleY-float64(scale)) > 0.08 {
			continue
		}
		if left+size*scale > imageWidth || top+size*scale > imageHeight {
			continue
		}
		matrix := sampleModules(pixels, left, top, size, scale)
		score := finderScore(matrix)
		candidates = append(candidates, candidate{score: score, scale: scale, matrix: matrix})
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("无法识别二维码模块")
	}
	best := candidates[0]
	for _, item := range candidates[1:] {
		if item.score > best.score || (item.score == best.score && item.scale > best.scale) {
			best = item
		}
	}
	if best.score < 0.85 {
		return nil, fmt.Errorf("二维码定位图案无效")
	}
	return best.matrix, nil
}

// decodePNG 解码 PNG 像素为黑白矩阵。
// decodePNG 解码 PNG 像素为黑白矩阵。
// 入参：
//   - data：入参说明见函数用途。
//
// 出参：[][]bool,error。
func decodePNG(data []byte) ([][]bool, error) {
	if len(data) < len(pngSignature) || string(data[:len(pngSignature)]) != pngSignature {
		return nil, fmt.Errorf("二维码图片不是 PNG")
	}
	var ihdr []byte
	var idat []byte
	position := len(pngSignature)
	for position+12 <= len(data) {
		length := int(binary.BigEndian.Uint32(data[position : position+4]))
		kind := string(data[position+4 : position+8])
		end := position + 12 + length
		if end > len(data) {
			return nil, fmt.Errorf("PNG 数据不完整")
		}
		payload := data[position+8 : end-4]
		switch kind {
		case "IHDR":
			ihdr = payload
		case "IDAT":
			idat = append(idat, payload...)
		case "IEND":
			position = end
			goto chunksDone
		}
		position = end
	}
chunksDone:

	if len(ihdr) != 13 || len(idat) == 0 {
		return nil, fmt.Errorf("PNG 缺少必要数据")
	}
	width := binary.BigEndian.Uint32(ihdr[0:4])
	height := binary.BigEndian.Uint32(ihdr[4:8])
	bitDepth := ihdr[8]
	colorType := ihdr[9]
	compression := ihdr[10]
	filterMethod := ihdr[11]
	interlace := ihdr[12]

	channels := 0
	switch colorType {
	case 0:
		channels = 1
	case 2:
		channels = 3
	case 4:
		channels = 2
	case 6:
		channels = 4
	}
	if channels == 0 || bitDepth != 8 || compression != 0 || filterMethod != 0 || interlace != 0 ||
		width < 1 || width > maxImageSide || height < 1 || height > maxImageSide {
		return nil, fmt.Errorf("PNG 格式不受支持")
	}

	stride := int(width) * channels
	expectedSize := (stride + 1) * int(height)
	raw, err := inflatePNG(idat, expectedSize)
	if err != nil {
		return nil, err
	}

	rows := make([][]bool, 0, height)
	previous := make([]byte, stride)
	offset := 0
	for range height {
		filterType := raw[offset]
		offset++
		current := append([]byte(nil), raw[offset:offset+stride]...)
		offset += stride
		if filterType > 4 {
			return nil, fmt.Errorf("PNG 滤镜不受支持")
		}
		for index, value := range current {
			var left, above, upperLeft byte
			if index >= channels {
				left = current[index-channels]
				upperLeft = previous[index-channels]
			}
			above = previous[index]
			switch filterType {
			case 1:
				current[index] = value + left
			case 2:
				current[index] = value + above
			case 3:
				current[index] = value + byte((int(left)+int(above))/2)
			case 4:
				current[index] = value + byte(paeth(int(left), int(above), int(upperLeft)))
			}
		}
		row := make([]bool, 0, width)
		for index := 0; index < stride; index += channels {
			var red, green, blue byte
			if colorType == 0 || colorType == 4 {
				red = current[index]
				green = current[index]
				blue = current[index]
			} else {
				red = current[index]
				green = current[index+1]
				blue = current[index+2]
			}
			alpha := byte(255)
			if colorType == 4 || colorType == 6 {
				alpha = current[index+channels-1]
			}
			luminance := (299*int(red) + 587*int(green) + 114*int(blue)) / 1000
			luminance = (luminance*int(alpha) + 255*(255-int(alpha))) / 255
			row = append(row, luminance < 128)
		}
		rows = append(rows, row)
		previous = current
	}
	return rows, nil
}

// inflatePNG 解压 IDAT 数据并严格校验长度与尾部。
// inflatePNG 解压 IDAT 数据并严格校验长度与尾部。
// 入参：
//   - compressed：入参说明见函数用途。
//   - expectedSize：入参说明见函数用途。
//
// 出参：[]byte,error。
func inflatePNG(compressed []byte, expectedSize int) ([]byte, error) {
	reader := bytes.NewReader(compressed)
	zlibReader, err := zlib.NewReader(reader)
	if err != nil {
		return nil, fmt.Errorf("PNG 像素数据无效")
	}
	defer zlibReader.Close()
	raw, err := io.ReadAll(io.LimitReader(zlibReader, int64(expectedSize)+1))
	if err != nil {
		return nil, fmt.Errorf("PNG 像素数据无效")
	}
	// 压缩流之后仍有多余数据（unconsumed_tail）时判定无效。
	if reader.Len() != 0 {
		return nil, fmt.Errorf("PNG 像素数据无效")
	}
	if len(raw) != expectedSize {
		return nil, fmt.Errorf("PNG 像素数据无效")
	}
	return raw, nil
}

// paeth
func paeth(left, above, upperLeft int) int {
	prediction := left + above - upperLeft
	leftDistance := absoluteInt(prediction - left)
	aboveDistance := absoluteInt(prediction - above)
	upperLeftDistance := absoluteInt(prediction - upperLeft)
	if leftDistance <= aboveDistance && leftDistance <= upperLeftDistance {
		return left
	}
	if aboveDistance <= upperLeftDistance {
		return above
	}
	return upperLeft
}

func absoluteInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

// sampleModules 按模块尺寸采样。
func sampleModules(pixels [][]bool, left, top, size, scale int) [][]bool {
	matrix := make([][]bool, 0, size)
	threshold := float64(scale*scale) / 2
	for moduleY := 0; moduleY < size; moduleY++ {
		row := make([]bool, 0, size)
		startY := top + moduleY*scale
		for moduleX := 0; moduleX < size; moduleX++ {
			startX := left + moduleX*scale
			black := 0
			for y := startY; y < startY+scale; y++ {
				for x := startX; x < startX+scale; x++ {
					if pixels[y][x] {
						black++
					}
				}
			}
			row = append(row, float64(black) >= threshold)
		}
		matrix = append(matrix, row)
	}
	return matrix
}

// finderScore 计算三个定位图案的匹配分数。
func finderScore(matrix [][]bool) float64 {
	size := len(matrix)
	matches := 0
	total := 0
	origins := [][2]int{{0, 0}, {size - 7, 0}, {0, size - 7}}
	for _, origin := range origins {
		originX, originY := origin[0], origin[1]
		for y := 0; y < 7; y++ {
			for x := 0; x < 7; x++ {
				edge := x == 0 || x == 6 || y == 0 || y == 6
				center := x >= 2 && x <= 4 && y >= 2 && y <= 4
				expected := edge || center
				if matrix[originY+y][originX+x] == expected {
					matches++
				}
				total++
			}
		}
	}
	return float64(matches) / float64(total)
}

// writeTerminalQR 渲染终端二维码。
func writeTerminalQR(matrix [][]bool, stream io.Writer) bool {
	size := len(matrix)
	whiteRow := make([]bool, size+quietZone*2)
	padded := make([][]bool, 0, len(matrix)+quietZone*2+1)
	for range quietZone {
		padded = append(padded, append([]bool(nil), whiteRow...))
	}
	for _, row := range matrix {
		item := make([]bool, 0, len(row)+quietZone*2)
		item = append(item, make([]bool, quietZone)...)
		item = append(item, row...)
		item = append(item, make([]bool, quietZone)...)
		padded = append(padded, item)
	}
	for range quietZone {
		padded = append(padded, append([]bool(nil), whiteRow...))
	}
	if len(padded)%2 == 1 {
		padded = append(padded, append([]bool(nil), whiteRow...))
	}

	ansi := isTerminal(stream) && os.Getenv("TERM") != "dumb"
	symbols := map[[2]bool]string{
		{false, false}: " ",
		{true, true}:   "█",
		{true, false}:  "▀",
		{false, true}:  "▄",
	}
	_, _ = fmt.Fprintln(stream)
	for index := 0; index < len(padded); index += 2 {
		topRow := padded[index]
		bottomRow := padded[index+1]
		var line strings.Builder
		for cell := range topRow {
			line.WriteString(symbols[[2]bool{topRow[cell], bottomRow[cell]}])
		}
		if ansi {
			// 黑字白底渲染二维码（颜色常量见 colors.go）。
			_, _ = fmt.Fprintf(stream, "%s%s%s%s\n", ansiForegroundBlack, ansiBackgroundWhite, line.String(), ansiReset)
		} else {
			_, _ = fmt.Fprintln(stream, line.String())
		}
	}
	return true
}

// isTerminal 判断输出流是否为终端。
func isTerminal(stream io.Writer) bool {
	file, ok := stream.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
