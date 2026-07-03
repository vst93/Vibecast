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

# GitHub mirror proxies for China mainland (tried in order, then falls back to direct GitHub)
MIRRORS=(
  "https://ghfast.top/"
  "https://gh-proxy.com/"
  ""
)

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

# If version not specified, get latest release tag.
# API calls: try DIRECT first (mirrors return 403 on api.github.com),
# then fall back to mirrors.
if [[ -z "$VERSION" ]]; then
  info "Fetching latest release..."
  API_URL="https://api.github.com/repos/${REPO}/releases/latest"

  # Direct first
  if command -v curl &>/dev/null; then
    VERSION=$(curl -fsSL --connect-timeout 10 "$API_URL" 2>/dev/null | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
  else
    VERSION=$(wget -qO- --timeout=10 "$API_URL" 2>/dev/null | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
  fi

  # Fallback to mirrors if direct failed
  if [[ -z "$VERSION" ]]; then
    for mirror in "${MIRRORS[@]}"; do
      [[ -z "$mirror" ]] && continue  # skip empty (already tried direct above)
      url="${mirror}${API_URL}"
      if command -v curl &>/dev/null; then
        VERSION=$(curl -fsSL --connect-timeout 10 "$url" 2>/dev/null | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
      else
        VERSION=$(wget -qO- --timeout=10 "$url" 2>/dev/null | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
      fi
      if [[ -n "$VERSION" ]]; then
        break
      fi
    done
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

# Download via mirrors
TMP_FILE=$(mktemp)
info "Downloading ${ASSET_NAME}"
downloaded=false
for mirror in "${MIRRORS[@]}"; do
  url="${mirror}${DOWNLOAD_URL}"
  if command -v curl &>/dev/null; then
    if curl -fsSL --connect-timeout 10 -o "$TMP_FILE" "$url"; then
      downloaded=true
      break
    fi
  else
    if wget -qO "$TMP_FILE" --timeout=10 "$url"; then
      downloaded=true
      break
    fi
  fi
done

if [[ "$downloaded" != "true" ]] || [[ ! -s "$TMP_FILE" ]]; then
  error "Download failed from all mirrors"
  rm -f "$TMP_FILE"
  exit 1
fi

# Verify SHA256 checksum (if SHA256SUMS is available for this release)
SUMS_URL="https://github.com/${REPO}/releases/download/${VERSION}/SHA256SUMS"
sums_ok=false
for mirror in "${MIRRORS[@]}"; do
  url="${mirror}${SUMS_URL}"
  if command -v curl &>/dev/null; then
    SUMS_DATA=$(curl -fsSL --connect-timeout 10 "$url" 2>/dev/null) || true
  else
    SUMS_DATA=$(wget -qO- --timeout=10 "$url" 2>/dev/null) || true
  fi
  if [[ -n "$SUMS_DATA" ]]; then
    sums_ok=true
    break
  fi
done

if [[ "$sums_ok" == "true" ]]; then
  EXPECTED_HASH=$(echo "$SUMS_DATA" | grep "$ASSET_NAME" | awk '{print $1}')
  if [[ -n "$EXPECTED_HASH" ]]; then
    ACTUAL_HASH=$(sha256sum "$TMP_FILE" | awk '{print $1}')
    if [[ "${ACTUAL_HASH,,}" != "${EXPECTED_HASH,,}" ]]; then
      error "Checksum mismatch!"
      error "  Expected: $EXPECTED_HASH"
      error "  Got:      $ACTUAL_HASH"
      rm -f "$TMP_FILE"
      exit 1
    fi
    info "Checksum verified"
  else
    warn "Asset not found in SHA256SUMS, skipping verification"
  fi
else
  warn "No SHA256SUMS found, skipping verification"
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
echo "  vibecast update                   # check and apply updates"
echo ""
echo "Then open http://localhost:8080/dashboard"
