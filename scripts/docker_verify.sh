#!/bin/bash
# thsdk-go 三平台验证脚本（在 docker 容器内执行，go1.27）
#   Linux：本机 cgo 构建 + vet + 全量单测 + hq.so 真实握手
#   Windows：mingw cgo 交叉构建 + wine 运行（单测 / hq.dll 握手）
#   macOS：darwin amd64/arm64 交叉编译（CGO_ENABLED=0）
# 用法：docker run --rm -v "$PWD":/work -w /work thsdk-verify bash scripts/docker_verify.sh
set -euo pipefail

echo "============================================================"
echo "== go 版本: $(go version)"
echo "============================================================"

echo
echo "== [1/4] Linux 构建 + vet + 全量单测（go1.27 本机 cgo） =="
go build ./...
go vet ./...
go test -count=1 .

echo
echo "== [2/4] Linux 真实组件握手（libs/linux/amd64/hq.so） =="
go run ./examples/smoke

echo
echo "== [3/4] Windows：mingw 交叉构建 + wine 运行（go1.27） =="
export CC=x86_64-w64-mingw32-gcc
export CGO_ENABLED=1 GOOS=windows GOARCH=amd64
go test -c -o /tmp/thsdk.test.exe .
WINEDEBUG=-all wine64 /tmp/thsdk.test.exe
go build -o /tmp/thsdk-smoke.exe ./examples/smoke
WINEDEBUG=-all wine64 /tmp/thsdk-smoke.exe
unset CC CGO_ENABLED GOOS GOARCH

echo
echo "== [4/4] macOS 交叉编译（darwin amd64 / arm64） =="
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ./...
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build ./...

echo
echo "============================================================"
echo "== 三平台验证全部通过（go1.27） =="
echo "== macOS 真机运行请在 Mac 上执行: go test ./... && go run ./examples/smoke =="
echo "============================================================"
