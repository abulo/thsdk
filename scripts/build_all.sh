#!/usr/bin/env bash
# thsdk-go 一键构建脚本：单元测试 + vet + 本机构建 + 五目标交叉编译。
#
# 用法：./scripts/build_all.sh
#
# 说明：
#   - 默认以 cgo 构建（运行时必需的原生桥）；C 工具链要求见 README。
#   - 交叉编译使用 CGO_ENABLED=0 的 purego 回退路径，验证纯 Go 代码
#     在三平台/两架构下均可编译；windows 的 cgo 运行时构建需在目标环境执行。
#   - 需要 Go 1.21+；如 go 不在 PATH，先 export PATH=$PATH:/usr/local/go/bin。

set -euo pipefail
cd "$(dirname "$0")/.."

echo "==> go build ./... (本机, cgo)"
go build ./...

echo "==> go vet ./..."
go vet ./...

echo "==> go test -count=1 ."
go test -count=1 .

echo "==> 交叉编译验证（CGO_ENABLED=0）"
cross_build() {
    local goos="$1" goarch="$2"
    echo "    GOOS=$goos GOARCH=$goarch"
    CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" go build ./...
}
cross_build darwin  amd64
cross_build darwin  arm64
cross_build linux   amd64
cross_build linux   arm64
cross_build windows amd64

echo "==> 构建示例程序（本机）"
go build -o /tmp/thsdk-demo ./examples/demo
go build -o /tmp/thsdk-table-demo ./examples/table
go build -o /tmp/thsdk-smoke ./examples/smoke
echo "    示例产物: /tmp/thsdk-demo /tmp/thsdk-table-demo /tmp/thsdk-smoke"

echo
echo "全部通过 ✅  可执行: go run ./examples/smoke（本机组件握手冒烟）"
