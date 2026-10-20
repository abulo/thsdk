# thsdk-go

`thsdk-go` 是 Python 版 `thsdk`（同花顺 THS 行情/交易 SDK，v2.0.2）的 **Go 逐行直译移植版**，作为全新项目独立交付：

- **80 个 Canonical 原生方法**全部实现（`List*` / `Get*` 命名，与 Python 包装一一对应）；
- **63 个 DataFrame Table 包装**（`<方法名>Table`）等价实现 `pandas.DataFrame`；
- **15 个便捷函数**（`GetPrice` / `Klines` / `IntradayData` / `TickLevel1` / `Depth` / `CorporateAction` / `SearchSymbols` / `CompleteTHSCode` / `WencaiNLP` / `News` / `MarketSecurities` / `BlockConstituents` 等）；
- **三平台运行组件**：macOS（x86_64 / arm64）、Linux（x86_64 / arm64）、Windows（x86_64）；
- 业务逻辑、常量定义、错误文案与 Python 源码逐行对齐；模块 `Version = "3.0.0"`、ABI `2.0.0`。

## 目录结构

```
thsdk-go/
├── go.mod / go.sum          # module thsdk；依赖 purego（CGO_ENABLED=0 回退路径）
├── errors.go                # THSDKError / AuthenticationError / NotAuthenticatedError / APIError
├── constants.go             # 20 个枚举（Adjust / Period / CallAuctionPhase / …）
├── colors.go                # 终端渲染 ANSI 颜色常量（二维码等输出统一取用）
├── method_manifest.go       # NativeMethodSpec + 80 newSpec + requestContracts + MethodInventory()
├── jsonutil.go              # JSON 严格校验 / 编解码 / parseResponse / 数字归一化
├── runtime.go               # RuntimeBridge：请求封装、限流重试、握手（version/strict/methods）
├── bridge.c / bridge.h      # 原生桥（专用线程承载 dlopen 与 CallAlloc/FreeResult）
├── bridge_cgo.go            # cgo 构建的原生桥实现（推荐，运行时唯一可靠路径）
├── bridge_purego*.go        # CGO_ENABLED=0 回退（purego FFI；编译验证用）
├── session.go               # Client / Auth（账号密码+MAC/验证码+临时账户）/ AuthQRCode / Request / Logout / Close
├── temporary_accounts.go    # 52 个临时账户
├── qr_display.go            # 二维码终端渲染（ZXing 兼容解码）
├── base.go                  # CallOptions / normalizeSecurity / normalizeBlock / normalizeNamedFields
├── api_account.go           # 13 个账户方法
├── api_securities.go        # 18 个证券方法
├── api_market_data.go       # 17 个行情方法
├── api_news.go              # 14 个快讯方法
├── api_wencai.go            # 7 个问财方法
├── api_rankings.go          # 11 个排行方法
├── api_dataframe.go         # 63 个 Table 包装
├── api_convenience.go       # 15 个便捷函数 + periods/adjustments/security/date 等助手
├── table.go                 # Table（pandas.DataFrame 等价实现）
├── thsdk.go                 # Version / getClient 单例 / SetDefaultClient / Close
├── libs/                    # 三平台原生运行组件副本
├── examples/                # demo（完整调用示例）/ table / smoke（本机组件握手）
└── scripts/build_all.sh     # 构建 / 测试 / 五目标交叉编译
```

## 构建与测试

要求：Go ≥ 1.21（本工程在 go1.27.1 验证）。构建 cgo 桥需要平台 C 编译器
（macOS: clang；Linux: gcc；Windows: MinGW-w64），这也是**运行时**必需的。

```bash
# 一键：go build + go vet + go test + 五目标交叉编译
./scripts/build_all.sh

# 或手动：
go build ./...
go vet ./...
go test -count=1 .        # 96 个单元测试（含内存边界、路径缓存、ANSI 颜色用例）
go run ./examples/smoke   # 本机原生组件握手冒烟（需 libs/ 与当前平台匹配）
```

交叉编译（无 C 工具链的编译验证路径，`CGO_ENABLED=0`）：

```bash
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build ./...
CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build ./...
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ./...
CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build ./...
```

> 说明：`GOOS=windows` 的 **cgo 运行时构建**需在 Windows（或带 MinGW 交叉链）
> 上执行；`CGO_ENABLED=0` 交叉编译仅验证纯 Go 代码路径的编译正确性。
> Windows 真机一键测试：`scripts\build_all_windows.bat`（构建 + vet + 96 项单测 + hq.dll 握手冒烟）。

## 内存与缓存设计

- **原生响应读取防护**：`readCString` 先逐字节定位 NUL 再一次性拷贝为
  `string`，读取上限 64 MiB（`maxResponseBytes`）；超长返回
  `THSDK 运行组件返回数据过大`，防止原生侧野指针失控。所有原生字符串
  由桥线程 `strdup` 后交宿主读取，读取完立即 `FreeResult`，无驻留。
- **请求长度防护**：请求 JSON 上限 1 MiB、方法名 ≤ 128 字节、JSON 嵌套
  ≤ 64 层（`maxRequestBytes` / `maxMethodBytes` / `maxJSONDepth`）。
- **二维码内存上限**：PNG 下载 ≤ 2 MiB、图像边长 ≤ 1024、解压严格校验
  长度与尾部（`maxImageBytes` / `maxImageSide`）。
- **组件路径缓存**：`THSDK_RUNTIME_PATH` 环境变量每次实时读取（可运行时
  切换）；按「可执行文件目录 / 工作目录 / 源码目录」的目录布局查找结果
  用 `sync.OnceValues` 进程内缓存，避免重复文件系统探测。
- **客户端单例**：`getClient()` 惰性创建进程级 `Client`，`RuntimeBridge`
  在其内部只构建一次（含 version/strict/methods 三次握手）；认证状态
  与会话目录均为显式行为，无隐式磁盘缓存增长。
- **终端颜色统一**：二维码等终端输出的 ANSI 颜色码集中在 `colors.go`
  公共常量（`ansiReset` / `ansiForegroundBlack` / `ansiBackgroundWhite`），
  不在渲染模块中散落魔法数字。

### 三平台实测记录（2026-10-19）

| 平台 | 方式 | 结果 |
|---|---|---|
| macOS x86_64 | 真机 cgo 构建 + hq.dylib 握手 | ✅ 96 单测 / 80 方法对齐 |
| Linux x86_64 | 真机 cgo 构建 + hq.so 握手 + 真实 API（二维码挑战） | ✅ 96 单测 / race 通过 / 80 方法对齐 |
| Windows x86_64 | MinGW cgo 交叉构建 + Wine 模拟运行 hq.dll | ✅ 96 单测 / 80 方法对齐 / 真实 API（二维码挑战） |

> Wine 模拟说明：Go ≥ 1.21 的 Windows 运行时与 Go c-shared 组件都要求
> `bcryptprimitives.dll`（Win10+ 系统 DLL）。本机 Wine 6 无此 DLL，实测时
> 从 Wine 10 包提取该 DLL 补入 system32 后，**直接用 `go 1.27.1` 构建的
> Windows 二进制即全绿**（96 单测 / hq.dll 握手 / 真实 API）；真实 Windows
> 系统自带该 DLL，无需任何处理。期间修复了 `bridge.c` 的 Windows 编译问题
> （`_WIN32_WINNT` 定义与 CONDITION_VARIABLE → 事件实现，兼容各 mingw-w64 SDK）。

### Docker 三平台验证（go1.27）

本仓库提供一键 Docker 验证（Linux 真机 + Windows-wine 模拟运行 + macOS 交叉编译）：

```bash
./scripts/run_docker_verify.sh        # 宿主机入口：构建镜像并执行验证
# 等价命令：
docker build -t thsdk-verify .
docker run --rm -v "$PWD":/work -w /work thsdk-verify bash scripts/docker_verify.sh
```

- 镜像：ubuntu 22.04 + **go1.27.1** + gcc/mingw-w64/wine64 + bcryptprimitives.dll（wine 补齐）
- 容器内验证项：Linux `build/vet/96 单测/hq.so 握手` → Windows `96 单测/hq.dll 握手`（wine）→ macOS `darwin amd64+arm64 交叉编译`
- macOS 真机运行（hq.dylib 加载）需在 Mac 上执行：`go test ./... && go run ./examples/smoke`（docker 无法运行 macOS 容器）

## 运行组件

`libs/` 目录布局与 Python 包一致（`libs/<系统>/<架构>/<组件>`）：

| 系统 | 架构 | 文件 |
|---|---|---|
| darwin | amd64 / arm64 | `hq.dylib` |
| linux | amd64 / arm64 | `hq.so` |
| windows | amd64 | `hq.dll` |

查找顺序（对应 Python `_resolve_component`）：`THSDK_RUNTIME_PATH` 环境变量
优先；其次可执行文件目录 / 当前工作目录 / 本仓库源码目录下的 `libs/...`。

## 快速开始

```go
package main

import (
    "fmt"
    "log"
    thsdk "thsdk"
)

func main() {
    // 方式一：账号密码（也可用 THS_USERNAME / THS_PASSWORD 环境变量）
    ok, err := thsdk.Auth("账号", "密码", "")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("认证完成:", ok)

    // 方式二：二维码登录（终端打印二维码，扫码后轮询）
    // ok, err := thsdk.AuthQRCode()

    defer thsdk.Logout()

    // 便捷函数：最近 5 根日 K（完整证券代码 = 四位市场前缀 + 代码）
    price, err := thsdk.GetPrice("USHA600519", &thsdk.GetPriceOptions{Count: 5})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(price)

    // 问财选股
    wencai, err := thsdk.WencaiNLP("市盈率小于 10 的贵州股票", &thsdk.WencaiNLPOptions{Limit: 10})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("问财结果行数:", wencai.RowCount())
}
```

### 完整证券代码格式

与 README 真实用法一致：**4 位字母市场前缀 + 证券代码**，例如 `USHA600519`
（上证 600519）、`USZA300033`（创业板 300033）。带点的 `000001.SZ` 不是 SDK
的合法输入。市场前缀映射：`USHA` 上证 / `USZA` 深证 / `USBJ` 北证 / `USHK`
港股 / `USUS` 美股 / `USUK` 英股 / `USJP` 日股（另有全量映射见 `constants.go`）。

### 认证方式

- **账号密码**：`Auth(username, password, extra)`；`extra` 传 `"mac"` 使用
  账号 MAC 派生验证码（前 6 字节设本地管理单播位），传 `"captcha:<验证码>"`
  配合 `AuthOptions{CaptchaID: "…"}` 处理图形验证码；
- **临时账户**：`Auth("", "", "")` 自动从 52 个临时账户挑选；
- **二维码**：`AuthQRCode()` 终端打印二维码，`WithAuthQRCodeTimeout` /
  `WithAuthQRCodePollInterval` 控制超时与轮询间隔。

### 请求选项

所有 `<方法名>Table` 与 `Client.Request` 接受 `*CallOptions`：

```go
table, err := thsdk.ListSecurityKLinesTable(map[string]any{
    "security": "USHA600519",
    "start":    -1,          // 0/负数 = 最近 N 根（配合 period）
    "end":      0,
    "adjust":   0,           // 0 不复权 / 1 前复权 / 2 后复权
    "period":   int64(0x4000), // 日线；1m=0x3001 … 120m=0x3078，周=0x5001 …
}, thsdk.WithTimeout(30))
```

## 与 Python 版的差异说明

| 项 | Python | Go |
|---|---|---|
| DataFrame | `pandas.DataFrame` | `*thsdk.Table`（`Columns`/`Rows`/`Index`/`IndexName`/`Attrs`） |
| 关键字参数 | 命名参数 | `*Options` 结构体 + 默认值 |
| JSON 数字 | `json.loads` 自动 int/float | `json.Number` → 整数 `int64`、小数 `float64` |
| 原生调用 | ctypes | cgo 专用线程桥（默认）；purego 回退（`CGO_ENABLED=0`） |
| 单例 | 模块级全局客户端 | `getClient()` + `SetDefaultClient` |
| 列顺序 | dict 插入序 | 记录内按字典序（Go map 无序，保证确定性） |
| 语义保留 | — | 数据校验、错误文案、常量值、周期/复权映射与 Python 完全一致 |

## 许可证与声明

`thsdk-go` 是 `thsdk-2.0.2` 的独立移植实现，仅用于学习与内部量化研究；
行情/交易数据接口的使用须遵守同花顺的相关服务条款。
