#!/usr/bin/env bash
# AutoDev Installer — https://github.com/HEETMEHTA18/autodev
# Usage: curl -fsSL https://raw.githubusercontent.com/HEETMEHTA18/autodev/main/scripts/install.sh | bash
set -euo pipefail

AUTODEV_VERSION="${AUTODEV_VERSION:-latest}"
INSTALL_DIR="${AUTODEV_INSTALL_DIR:-/usr/local/bin}"
REPO="heetmehta18/autodev"
BINARY_NAME="autodev"
GITHUB_API="https://api.github.com"
RELEASES_URL="https://github.com/${REPO}/releases"

# ── Parse arguments ───────────────────────────────────────────────────────────
INSTALL_PACKAGES=()
while [[ $# -gt 0 ]]; do
  case "$1" in
    --install)
      shift
      while [[ $# -gt 0 && ! "$1" =~ ^- ]]; do
        INSTALL_PACKAGES+=("$1")
        shift
      done
      ;;
    *)
      shift
      ;;
  esac
done

# ── Colours ──────────────────────────────────────────────────────────────────
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
BLUE='\033[0;34m'; BOLD='\033[1m'; RESET='\033[0m'

info()    { echo -e "${BLUE}[autodev]${RESET} $*"; }
success() { echo -e "${GREEN}[autodev]${RESET} $*"; }
warn()    { echo -e "${YELLOW}[autodev]${RESET} $*"; }
error()   { echo -e "${RED}[autodev] ERROR:${RESET} $*" >&2; exit 1; }

# ── Banner ────────────────────────────────────────────────────────────────────
echo -e "${YELLOW}"
cat <<'EOF'
  AutoDev Installer
  Clone. Scan. Install. Build. — The App Store for Developers
  https://github.com/HEETMEHTA18/autodev
EOF
echo -e "${RESET}"

# ── OS / Arch detection ───────────────────────────────────────────────────────
detect_os() {
  case "$(uname -s)" in
    Linux*)   echo "linux" ;;
    Darwin*)  echo "darwin" ;;
    MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
    *) error "Unsupported OS: $(uname -s)" ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64)  echo "amd64" ;;
    arm64|aarch64) echo "arm64" ;;
    armv7*)        echo "armv7" ;;
    *) error "Unsupported architecture: $(uname -m)" ;;
  esac
}

OS=$(detect_os)
ARCH=$(detect_arch)
info "Detected: ${OS}/${ARCH}"

# ── PATH setup ────────────────────────────────────────────────────────────────
# If no install dir was chosen and the user cannot write to /usr/local/bin,
# fall back to a user-local bin dir and register it on PATH for all shells.
if [ -z "${AUTODEV_INSTALL_DIR:-}" ]; then
  if [ "${OS}" = "linux" ] || [ "${OS}" = "darwin" ]; then
    if [ -d "${HOME}/.local/bin" ] && [ -w "${HOME}/.local/bin" ]; then
      INSTALL_DIR="${HOME}/.local/bin"
      info "Using ${INSTALL_DIR} (user-writable, no sudo needed)"
    elif ! touch /usr/local/bin/.autodev-write-test 2>/dev/null; then
      INSTALL_DIR="${HOME}/.local/bin"
      mkdir -p "${INSTALL_DIR}"
      info "Using ${INSTALL_DIR} (registered on PATH for you)"
    else
      rm -f /usr/local/bin/.autodev-write-test
    fi
  elif [ "${OS}" = "windows" ]; then
    # Windows: prefer a user-local dir under LOCALAPPDATA.
    INSTALL_DIR="${LOCALAPPDATA:-$HOME/AppData/Local}/autodev/bin"
    mkdir -p "${INSTALL_DIR}"
    info "Using ${INSTALL_DIR} (Windows user install)"
  fi
fi

# Register the install dir on PATH for the current OS.
register_path() {
  local dir="$1"
  case "${OS}" in
    windows)
      local cur
      cur=$(cmd //c "echo %PATH%" 2>/dev/null || echo "")
      case "${cur}" in
        *"${dir}"*) ;;
        *) setx PATH "${dir};${cur}" >/dev/null 2>&1 && info "Added ${dir} to Windows user PATH" || warn "Could not setx PATH — add ${dir} to PATH manually." ;;
      esac
      ;;
    *)
      # POSIX shells (bash/zsh/sh) + fish.
      local line="export PATH=\"${dir}:\$PATH\""
      local fish_line="set -gx PATH \"${dir}\" \$PATH"
      for rc in "${HOME}/.bashrc" "${HOME}/.zshrc" "${HOME}/.profile" "${HOME}/.config/fish/config.fish"; do
        if [ -f "${rc}" ]; then
          if ! grep -qF "${dir}" "${rc}" 2>/dev/null; then
            if printf '%s\n' "${rc}" | grep -q fish; then
              printf '%s\n' "${fish_line}" >> "${rc}"
            else
              printf '%s\n' "${line}" >> "${rc}"
            fi
            info "Added ${dir} to ${rc}"
          fi
        fi
      done
      # Ensure at least one rc file gets the line so PATH survives new shells.
      if [ -f "${HOME}/.profile" ] && ! grep -qF "${dir}" "${HOME}/.profile" 2>/dev/null; then
        printf '%s\n' "${line}" >> "${HOME}/.profile"
      fi
      ;;
  esac
}

if [ -n "${INSTALL_DIR}" ]; then
  register_path "${INSTALL_DIR}"
fi

# ── Fetch latest version ──────────────────────────────────────────────────────
if [ "${AUTODEV_VERSION}" = "latest" ]; then
  info "Fetching latest release..."
  GITHUB_RELEASE_URL="${GITHUB_API}/repos/${REPO}/releases/latest"
  
  # Fetch JSON response from GitHub API
  if command -v curl &>/dev/null; then
    RELEASE_JSON=$(curl -fsSL -s "${GITHUB_RELEASE_URL}" || echo "")
  elif command -v wget &>/dev/null; then
    RELEASE_JSON=$(wget -qO- "${GITHUB_RELEASE_URL}" || echo "")
  else
    error "curl or wget is required"
  fi

  # Check if repo exists or has releases
  if [ -z "${RELEASE_JSON}" ] || echo "${RELEASE_JSON}" | grep -q "Not Found"; then
    warn "No GitHub release found for ${REPO}."
    warn "This happens if you haven't created a Release on GitHub yet."
    echo ""
    echo "To fix this and make the installer work:"
    echo "  1. Compile binaries locally or via GitHub Actions"
    echo "  2. Go to https://github.com/${REPO}/releases and create a release"
    echo "  3. Upload the compiled binaries as release assets"
    echo ""
    error "Please create a GitHub release first, or run the installer with a specific version: AUTODEV_VERSION=0.1.0"
  fi

  # Parse version tag
  TAG_NAME=$(echo "${RELEASE_JSON}" | sed -n -E 's/.*"tag_name":[[:space:]]*"([^"]+)".*/\1/p' || echo "")

  if [ -z "${TAG_NAME}" ]; then
    error "Could not parse version tag from GitHub API response. Please check if the release exists."
  fi

  # Normalize version for display (remove leading 'v')
  AUTODEV_VERSION="${TAG_NAME#v}"
else
  # User manually specified AUTODEV_VERSION
  if [[ "${AUTODEV_VERSION}" == v* ]]; then
    TAG_NAME="${AUTODEV_VERSION}"
    AUTODEV_VERSION="${AUTODEV_VERSION#v}"
  else
    if [[ "${AUTODEV_VERSION}" =~ ^[0-9] ]]; then
      TAG_NAME="v${AUTODEV_VERSION}"
    else
      TAG_NAME="${AUTODEV_VERSION}"
    fi
  fi
fi

info "Installing AutoDev v${AUTODEV_VERSION}..."

# ── Download ──────────────────────────────────────────────────────────────────
ARCHIVE_NAME="autodev_${OS}_${ARCH}"
if [ "${OS}" = "windows" ]; then
  ARCHIVE_EXT="zip"
else
  ARCHIVE_EXT="tar.gz"
fi

DOWNLOAD_URL="${RELEASES_URL}/download/${TAG_NAME}/${ARCHIVE_NAME}.${ARCHIVE_EXT}"
TMP_DIR=$(mktemp -d)
ARCHIVE_PATH="${TMP_DIR}/autodev.${ARCHIVE_EXT}"

info "Downloading from: ${DOWNLOAD_URL}"
if command -v curl &>/dev/null; then
  curl -fsSL --progress-bar "${DOWNLOAD_URL}" -o "${ARCHIVE_PATH}"
else
  wget -q --show-progress "${DOWNLOAD_URL}" -O "${ARCHIVE_PATH}"
fi

# ── Extract ───────────────────────────────────────────────────────────────────
info "Extracting..."
if [ "${ARCHIVE_EXT}" = "tar.gz" ]; then
  tar -xzf "${ARCHIVE_PATH}" -C "${TMP_DIR}"
else
  unzip -q "${ARCHIVE_PATH}" -d "${TMP_DIR}"
fi

BINARY_PATH="${TMP_DIR}/${BINARY_NAME}"
if [ ! -f "${BINARY_PATH}" ]; then
  # The release archive nests the binary under a folder (e.g.
  # autodev_linux_amd64/autodev). Locate it recursively.
  BINARY_PATH=$(find "${TMP_DIR}" -type f -name "${BINARY_NAME}" 2>/dev/null | head -n1)
fi
if [ -z "${BINARY_PATH}" ] || [ ! -f "${BINARY_PATH}" ]; then
  error "Could not locate the ${BINARY_NAME} binary inside the downloaded archive."
fi
chmod +x "${BINARY_PATH}"

# ── Install ───────────────────────────────────────────────────────────────────
info "Installing to ${INSTALL_DIR}..."
mkdir -p "${INSTALL_DIR}"
if [ -w "${INSTALL_DIR}" ]; then
  cp "${BINARY_PATH}" "${INSTALL_DIR}/${BINARY_NAME}"
else
  sudo cp "${BINARY_PATH}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

# ── Cleanup ───────────────────────────────────────────────────────────────────
rm -rf "${TMP_DIR}"

# ── Verify ────────────────────────────────────────────────────────────────────
if command -v autodev &>/dev/null; then
  VERSION_OUT=$(autodev --version 2>&1 || true)
  success "[OK] AutoDev installed: ${VERSION_OUT}"
else
  warn "Binary installed to ${INSTALL_DIR}/${BINARY_NAME}"
  warn "Make sure ${INSTALL_DIR} is in your PATH."
fi

# ── Package installation ──────────────────────────────────────────────────────
if [ ${#INSTALL_PACKAGES[@]} -gt 0 ]; then
  echo ""
  info "Installing requested packages: ${INSTALL_PACKAGES[*]}"
  
  # Resolve path to the autodev binary we just installed
  AUTODEV_BIN="${INSTALL_DIR}/${BINARY_NAME}"
  if [ ! -f "${AUTODEV_BIN}" ]; then
    AUTODEV_BIN=$(command -v autodev || echo "")
  fi
  
  if [ -n "${AUTODEV_BIN}" ]; then
    for pkg in "${INSTALL_PACKAGES[@]}"; do
      info "Installing package: ${pkg}..."
      if ! "${AUTODEV_BIN}" install "${pkg}"; then
        warn "Failed to install ${pkg}"
      fi
    done
  else
    error "Could not locate autodev binary to install packages."
  fi
fi

echo ""
echo -e "${BOLD}Getting started:${RESET}"
echo -e "  ${GREEN}autodev${RESET}              — open interactive installer"
echo -e "  ${GREEN}autodev install nodejs${RESET}   — install a specific package"
echo -e "  ${GREEN}autodev profile web-dev${RESET}  — install a developer profile"
echo -e "  ${GREEN}autodev github USERNAME${RESET}  — scan a GitHub user's repos"
echo -e "  ${GREEN}autodev doctor${RESET}           — check environment health"
echo ""
echo -e "  Docs: ${BLUE}https://github.com/HEETMEHTA18/autodev${RESET}"
echo ""
