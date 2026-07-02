#!/usr/bin/env bash
set -euo pipefail

# Vibecast install script
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/vst93/Vibecast/main/install.sh | bash
#   curl -fsSL https://raw.githubusercontent.com/vst93/Vibecast/main/install.sh | bash -s -- --version 20260702
#   curl -fsSL https://raw.githubusercontent.com/vst93/Vibecast/main/install.sh | bash -s -- --dir /opt/vibecast

REPO="vst93/Vibecast"
INSTALL_DIR="/usr/local/bin"
VERSION=""
BINARY_NAME="vibecast"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

info()  { echo -e "${GREEN}✓${NC} $1"; }
warn()  { echo -e "${YELLOW}⚠${NC} $1"; }
error() { echo -e "${RED}✗${NC} $1"; }

# Parse args
while [[ $# -gt 0 ]]; do
  case "$1" in
    --version) VERSION="$2"; shift 2 ;;
    --dir)     INSTALL_DIR="$2"; shift 2 ;;
    --help)
      echo "Usage: install.sh [--version <tag>] [--dir <path>]"
      echo "  --version  Specific release tag (default: latest)"
      echo "  --dir      Install directory (default: /usr/local/bin)"
      exit 0 ;;
    *) error "Unknown option: $1"; exit 1 ;;
  esac
done

# Detect OS and arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
  linux)  OS="linux" ;;
  darwin) OS="darwin" ;;
  *) error "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) error "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# If version not specified, get latest release tag
if [[ -z "$VERSION" ]]; then
  info "Fetching latest release..."
  if command -v curl &>/dev/null; then
    VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
  else
    VERSION=$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
  fi
  if [[ -z "$VERSION" ]]; then
    error "Could not determine latest version. Specify with --version."
    exit 1
  fi
fi

info "Version: ${VERSION}"
info "Platform: ${OS}/${ARCH}"

# Find the matching asset
ASSET_NAME="vibecast-${VERSION}-${OS}-${ARCH}"
if [[ "$OS" == "windows" ]]; then
  ASSET_NAME="${ASSET_NAME}.exe"
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET_NAME}"

# Download
TMP_FILE=$(mktemp)
info "Downloading ${DOWNLOAD_URL}"
if command -v curl &>/dev/null; then
  curl -fsSL -o "$TMP_FILE" "$DOWNLOAD_URL"
else
  wget -qO "$TMP_FILE" "$DOWNLOAD_URL"
fi

if [[ ! -s "$TMP_FILE" ]]; then
  error "Download failed or file is empty"
  rm -f "$TMP_FILE"
  exit 1
fi

chmod +x "$TMP_FILE"

# Install
if [[ -w "$INSTALL_DIR" ]]; then
  mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
else
  warn "Needs sudo to install to ${INSTALL_DIR}"
  sudo mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
  sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
fi

info "Installed to ${INSTALL_DIR}/${BINARY_NAME}"
echo ""
echo -e "${GREEN}Vibecast ${VERSION} installed!${NC}"
echo ""
echo "Quick start:"
echo "  vibecast                          # run with defaults"
echo "  vibecast --addr :3000             # custom port"
echo "  vibecast --storage /var/lib/vibecast/sites --db /var/lib/vibecast/vibecast.db"
echo ""
echo "Then open http://localhost:8080/dashboard"
