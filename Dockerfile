# thsdk-go 三平台验证镜像（go1.27）
# 用法：在工程根目录执行
#   docker build -t thsdk-verify .
#   docker run --rm -v "$PWD":/work -w /work thsdk-verify bash scripts/docker_verify.sh
# 说明：Linux 真机验证 + Windows(wine 模拟)运行验证 + macOS 交叉编译验证。
#       容器内 Linux 版 wine 6 缺少 bcryptprimitives.dll（Go>=1.21 Windows 运行时必需），
#       此处从 wine10 包提取该 DLL 补入；真实 Windows 系统自带，无需此步。
FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates curl git gcc libc6-dev \
        gcc-mingw-w64-x86-64 wine64 xvfb \
        && rm -rf /var/lib/apt/lists/*

# go1.27（官方二进制发布版）
ARG GO_VERSION=1.27.1
RUN curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
        | tar -C /usr/local -xz
ENV PATH=/usr/local/go/bin:$PATH \
    GOPROXY=https://goproxy.cn,direct \
    GOFLAGS=-mod=mod

# wine 6 缺少 bcryptprimitives.dll：从 wine10 包提取补齐（真实 Windows 不需要）
RUN curl -fsSL -o /tmp/libwine10.deb \
        https://mirrors.aliyun.com/ubuntu/pool/universe/w/wine/libwine_10.0~repack-12ubuntu1_amd64.deb \
    && dpkg-deb -x /tmp/libwine10.deb /tmp/wine10 \
    && mkdir -p /root/.wine/drive_c/windows/system32 \
    && cp /tmp/wine10/usr/lib/x86_64-linux-gnu/wine/x86_64-windows/bcryptprimitives.dll \
          /root/.wine/drive_c/windows/system32/ \
    && rm -rf /tmp/libwine10.deb /tmp/wine10

WORKDIR /work
CMD ["bash", "scripts/docker_verify.sh"]
