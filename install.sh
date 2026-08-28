#!/bin/sh
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

REPO="${REPO:-Gonanf/pj}"
BINARY_NAME="pj"

info() {
    printf "${BLUE}==>${NC} %s\n" "$1"
}

success() {
    printf "${GREEN}==>${NC} %s\n" "$1"
}

error() {
    printf "${RED}Error:${NC} %s\n" "$1" >&2
    exit 1
}

# 1. Detect OS
OS="$(uname -s)"
case "$OS" in
    Linux)  OS_TAG="linux" ;;
    Darwin) OS_TAG="darwin" ;;
    *)      error "Unsupported operating system: $OS" ;;
esac

# 2. Detect Architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64)   ARCH_TAG="amd64" ;;
    arm64|aarch64)  ARCH_TAG="arm64" ;;
    *)              error "Unsupported architecture: $ARCH" ;;
esac

info "Detected platform: ${OS_TAG}/${ARCH_TAG}"

# 3. Determine latest version
if [ -z "$VERSION" ]; then
    info "Resolving latest release version..."
    LATEST_URL="https://api.github.com/repos/${REPO}/releases/latest"
    VERSION=$(curl -sSL -H "Accept: application/vnd.github.v3+json" "$LATEST_URL" 2>/dev/null | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    if [ -z "$VERSION" ]; then
        VERSION="v0.1.0"
        info "Could not determine latest release from GitHub API, falling back to ${VERSION}"
    else
        info "Latest release: ${VERSION}"
    fi
else
    info "Using requested version: ${VERSION}"
fi

# Clean version string (remove leading 'v' if present for tarball name)
RAW_VERSION="${VERSION#v}"

# 4. Construct download URLs
# Format matches goreleaser: pj_{version}_{os}_{arch}.tar.gz
TARBALL="pj_${RAW_VERSION}_${OS_TAG}_${ARCH_TAG}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${TARBALL}"
CHECKSUMS_URL="https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

info "Downloading ${TARBALL}..."
if ! curl -fsSL "$DOWNLOAD_URL" -o "${TMP_DIR}/${TARBALL}"; then
    error "Failed to download binary from: $DOWNLOAD_URL"
fi

# 5. Checksum verification
if curl -fsSL "$CHECKSUMS_URL" -o "${TMP_DIR}/checksums.txt" 2>/dev/null; then
    info "Verifying SHA-256 checksum..."
    cd "$TMP_DIR"
    if command -v sha256sum >/dev/null 2>&1; then
        grep "$TARBALL" checksums.txt | sha256sum -c - >/dev/null 2>&1 || error "Checksum verification failed"
    elif command -v shasum >/dev/null 2>&1; then
        grep "$TARBALL" checksums.txt | shasum -a 256 -c - >/dev/null 2>&1 || error "Checksum verification failed"
    fi
    info "Checksum verified successfully."
    cd - >/dev/null
fi

# 6. Extract binary
info "Extracting ${BINARY_NAME}..."
tar -xzf "${TMP_DIR}/${TARBALL}" -C "$TMP_DIR"

if [ ! -f "${TMP_DIR}/${BINARY_NAME}" ]; then
    error "Binary not found in archive"
fi

# 7. Select installation destination
DEST_DIR="/usr/local/bin"
if [ ! -w "$DEST_DIR" ]; then
    if [ -n "$SUDO_USER" ] || [ "$(id -u)" = "0" ]; then
        DEST_DIR="/usr/local/bin"
    else
        DEST_DIR="$HOME/.local/bin"
        mkdir -p "$DEST_DIR"
    fi
fi

TARGET="${DEST_DIR}/${BINARY_NAME}"
info "Installing to ${TARGET}..."

if [ -w "$DEST_DIR" ]; then
    mv "${TMP_DIR}/${BINARY_NAME}" "$TARGET"
else
    sudo mv "${TMP_DIR}/${BINARY_NAME}" "$TARGET"
fi

chmod +x "$TARGET"
success "${BINARY_NAME} ${VERSION} installed successfully to ${TARGET}!"

# 8. Check PATH
case ":$PATH:" in
    *":$DEST_DIR:"*) ;;
    *)
        printf "\n${BLUE}Note:${NC} %s is not in your PATH.\n" "$DEST_DIR"
        printf "Add it by adding this line to your shell profile (~/.bashrc, ~/.zshrc, or ~/.config/fish/config.fish):\n"
        printf "  export PATH=\"%s:\$PATH\"\n\n" "$DEST_DIR"
        ;;
esac
