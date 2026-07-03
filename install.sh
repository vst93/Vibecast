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

# --- Timezone-based language detection ---
# Detect if system is in UTC+8 (China Standard Time) → output Chinese, otherwise English.
detect_lang() {
  local offset
  # Try /etc/timezone first (Debian/Ubuntu)
  if [[ -f /etc/timezone ]]; then
    local tz
    tz=$(cat /etc/timezone 2>/dev/null || echo "")
    if [[ "$tz" == "Asia/Shanghai" || "$tz" == "Asia/Urumqi" || "$tz" == "Asia/Hong_Kong" || "$tz" == "Asia/Taipei" ]]; then
      LANG_ZH=1
      return
    fi
  fi
  # Try timedatectl (systemd)
  if command -v timedatectl &>/dev/null; then
    local tz
    tz=$(timedatectl show -p Timezone --value 2>/dev/null || echo "")
    if [[ "$tz" == "Asia/Shanghai" || "$tz" == "Asia/Urumqi" || "$tz" == "Asia/Hong_Kong" || "$tz" == "Asia/Taipei" ]]; then
      LANG_ZH=1
      return
    fi
  fi
  # Fallback: check UTC offset via date
  offset=$(date +%z 2>/dev/null || echo "+0000")
  if [[ "$offset" == "+0800" ]]; then
    LANG_ZH=1
    return
  fi
  LANG_ZH=0
}

# i18n helper: $1 = English text, $2 = Chinese text
i18n() {
  if [[ "${LANG_ZH:-0}" == "1" ]]; then
    echo -n "$2"
  else
    echo -n "$1"
  fi
}

detect_lang

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
      echo "  --version  $(i18n "Specific release tag" "指定版本号") (default: latest)"
      echo "  --dir      $(i18n "Install directory" "安装目录") (default: /usr/local/bin)"
      exit 0 ;;
    *) error "$(i18n "Unknown option" "未知选项"): $1"; exit 1 ;;
  esac
done

# Detect OS and arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
  linux)  OS="linux" ;;
  darwin) OS="darwin" ;;
  *) error "$(i18n "Unsupported OS" "不支持的操作系统"): $OS"; exit 1 ;;
esac

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) error "$(i18n "Unsupported architecture" "不支持的架构"): $ARCH"; exit 1 ;;
esac

# If version not specified, get latest release tag.
# API calls: try DIRECT first (mirrors return 403 on api.github.com),
# then fall back to mirrors.
if [[ -z "$VERSION" ]]; then
  info "$(i18n "Fetching latest release..." "正在获取最新版本...")"
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
    error "$(i18n "Could not determine latest version. Specify with --version." "无法确定最新版本，请使用 --version 指定。")"
    exit 1
  fi
fi

info "$(i18n "Version" "版本"): ${VERSION}"
info "$(i18n "Platform" "平台"): ${OS}/${ARCH}"

# Find the matching asset
ASSET_NAME="vibecast-${VERSION}-${OS}-${ARCH}"
if [[ "$OS" == "windows" ]]; then
  ASSET_NAME="${ASSET_NAME}.exe"
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET_NAME}"

# Download via mirrors
TMP_FILE=$(mktemp)
info "$(i18n "Downloading" "下载中") ${ASSET_NAME}"
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
  error "$(i18n "Download failed from all mirrors" "所有镜像下载均失败")"
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
      error "$(i18n "Checksum mismatch!" "校验和不匹配！")"
      error "  Expected: $EXPECTED_HASH"
      error "  Got:      $ACTUAL_HASH"
      rm -f "$TMP_FILE"
      exit 1
    fi
    info "$(i18n "Checksum verified" "校验和验证通过")"
  else
    warn "$(i18n "Asset not found in SHA256SUMS, skipping verification" "未在 SHA256SUMS 中找到该文件，跳过验证")"
  fi
else
  warn "$(i18n "No SHA256SUMS found, skipping verification" "未找到 SHA256SUMS，跳过验证")"
fi

chmod +x "$TMP_FILE"

# Install
if [[ -w "$INSTALL_DIR" ]]; then
  mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
else
  warn "$(i18n "Needs sudo to install to" "需要 sudo 权限安装到") ${INSTALL_DIR}"
  sudo mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
  sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
fi

info "$(i18n "Installed to" "已安装至") ${INSTALL_DIR}/${BINARY_NAME}"
echo ""
echo -e "${GREEN}Vibecast ${VERSION} $(i18n "installed!" "已安装！")${NC}"
echo ""
echo "$(i18n "Quick start:" "快速开始：")"
echo "  vibecast                          # $(i18n "run with defaults" "使用默认配置运行")"
echo "  vibecast --addr :3000             # $(i18n "custom port" "自定义端口")"
echo "  vibecast --storage /var/lib/vibecast/sites --db /var/lib/vibecast/vibecast.db"
echo ""
echo "  vibecast update                   # $(i18n "check and apply updates" "检查并应用更新")"
echo ""
echo "$(i18n "Then open" "然后打开") http://localhost:8080/dashboard"
