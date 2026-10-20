#!/bin/bash
# thsdk-go docker 三平台验证：宿主机一键入口
# 前置：本机已安装 docker 且 daemon 可用（macOS/Windows/Linux 均可）
# 用法：./scripts/run_docker_verify.sh
set -euo pipefail
cd "$(dirname "$0")/.."

echo "== 构建验证镜像（ubuntu22.04 + go1.27 + gcc/mingw/wine） =="
docker build -t thsdk-verify .

echo "== 在容器内执行三平台验证（Linux 真机 / Windows-wine / macOS 交叉编译） =="
docker run --rm -v "$PWD":/work -w /work thsdk-verify bash scripts/docker_verify.sh

echo "== 完成：结果见上方输出 =="
