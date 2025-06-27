#!/usr/bin/env bash
set -euxo pipefail

# ==============================================================================
# 安装脚本 install.sh
#
# 功能：
#   自动下载并安装 tlsctl 二进制文件，支持 GitHub 和 CNB 镜像源。
#   自动识别系统架构和操作系统。
#   支持通过环境变量自定义版本号、镜像源和安装目录。
#
# 环境变量：
#   VERSION     安装版本，默认 v0.0.2
#   MIRROR      下载源，可选 github（默认）或 cnb
#   INSTALL_DIR 安装路径，默认 /usr/local/bin
#
# 示例用法：
#   curl -sSL https://cnb.cool/zhiqiangwang/tlsctl/-/git/raw/main/install.sh | MIRROR=cnb bash
#   curl -sSL https://raw.githubusercontent.com/chihqiang/tlsctl/main/install.sh | bash
#
# 依赖项：
#   jq、wget、tar、gzip
#
# 退出码：
#   0   安装成功
#   1   安装失败或参数错误
# ==============================================================================

# ========== CONFIG ==========
BIN_NAME="tlsctl"
VERSION="${VERSION:-"v0.0.2"}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
MIRROR="${MIRROR:-github}"

# ========== CHECK DEPENDENCIES ==========
if ! command -v jq &>/dev/null; then
    echo "❌ Error: jq is not installed. Please install jq."
    exit 1
fi

# ========== DETECT OS & ARCH ==========
ARCH=$(uname -m)
OS=$(uname | tr '[:upper:]' '[:lower:]')
case "$ARCH" in
    "x86_64") ARCH="amd64" ;;
    "aarch64"|"arm64") ARCH="arm64" ;;
    *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

# ========== DOWNLOAD URL ==========
version_filename=${VERSION}/${BIN_NAME}_${OS}_${ARCH}.tar.gz
case "$MIRROR" in
    cnb)
        DOWNLOAD_URL="https://cnb.cool/zhiqiangwang/tlsctl/-/releases/download/${version_filename}"
        ;;
    github)
        DOWNLOAD_URL="https://github.com/chihqiang/tlsctl/releases/download/${version_filename}"
        ;;
    *)
        echo "❌ Unknown MIRROR: '$MIRROR'. Use 'github' or 'cnb'."
        exit 1
        ;;
esac

echo "📦 Downloading ${BIN_NAME} ${VERSION} from ${MIRROR}..."

# ========== DOWNLOAD & INSTALL ==========
TEMP="$(mktemp -d)"
trap 'rm -rf "$TEMP"' EXIT INT

wget --progress=dot:mega "${DOWNLOAD_URL}" -O "$TEMP/${BIN_NAME}.tar.gz" || {
    echo "❌ Failed to download from ${DOWNLOAD_URL}"
    exit 1
}

cd "$TEMP"
tar -zxf "${BIN_NAME}.tar.gz"
chmod +x "${BIN_NAME}"

SUDO=""
if [ "$(id -u)" -ne 0 ]; then
    SUDO="sudo"
fi

if [ -f "${INSTALL_DIR}/${BIN_NAME}" ]; then
    echo "🧹 Removing old version from ${INSTALL_DIR}..."
    $SUDO rm -f "${INSTALL_DIR}/${BIN_NAME}"
fi

echo "🚀 Installing ${BIN_NAME} to ${INSTALL_DIR}..."
$SUDO mv "${BIN_NAME}" "${INSTALL_DIR}/"

# ========== VERIFY ==========
echo "✅ Verifying installation..."
if ! "${BIN_NAME}" --help >/dev/null 2>&1; then
    echo "❌ Error: ${BIN_NAME} installation failed"
    exit 1
fi

echo "🎉 ${BIN_NAME} ${VERSION} installed successfully."
