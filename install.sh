#!/bin/sh
set -e

REPO="born1337/hyperliquid-terminal"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  darwin) OS="darwin" ;;
  linux)  OS="linux" ;;
  mingw*|msys*|cygwin*) OS="windows" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect arch
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64)  ARCH="amd64" ;;
  arm64|aarch64)  ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest version
VERSION=$(curl -sSf "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
if [ -z "$VERSION" ]; then
  echo "Failed to fetch latest version. Falling back to go install..."
  if command -v go >/dev/null 2>&1; then
    go install "github.com/${REPO}@latest"
    echo "Installed hltui via go install"
    exit 0
  fi
  echo "Error: no releases found and Go is not installed"
  exit 1
fi

echo "Installing hltui ${VERSION} (${OS}/${ARCH})..."

# Download
EXT="tar.gz"
if [ "$OS" = "windows" ]; then
  EXT="zip"
fi
URL="https://github.com/${REPO}/releases/download/${VERSION}/hltui_${OS}_${ARCH}.${EXT}"

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

curl -sSfL "$URL" -o "${TMPDIR}/hltui.${EXT}"

# Extract
cd "$TMPDIR"
if [ "$EXT" = "zip" ]; then
  unzip -q "hltui.${EXT}"
else
  tar xzf "hltui.${EXT}"
fi

# Install
if [ -w "$INSTALL_DIR" ]; then
  mv hltui "$INSTALL_DIR/hltui"
else
  sudo mv hltui "$INSTALL_DIR/hltui"
fi

echo "hltui ${VERSION} installed to ${INSTALL_DIR}/hltui"
