#!/bin/bash
set -e

# Finally AI CodeKeeper Installation Script
# Usage: curl -fsSL https://raw.githubusercontent.com/zvika-finally/ai-codekeeper/main/install.sh | bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Installing Finally AI CodeKeeper...${NC}"

# Detect OS and architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux)
        OS="linux"
        ;;
    Darwin)
        OS="darwin"
        ;;
    CYGWIN*|MINGW*|MSYS*)
        OS="windows"
        ;;
    *)
        echo -e "${RED}❌ Unsupported operating system: $OS${NC}"
        exit 1
        ;;
esac

case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}❌ Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

# Get the latest release version
echo -e "${BLUE}📡 Fetching latest release...${NC}"
LATEST_VERSION=$(curl -s https://api.github.com/repos/zvika-finally/ai-codekeeper/releases/latest | grep '"tag_name"' | cut -d'"' -f4)

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}❌ Failed to fetch latest release version${NC}"
    exit 1
fi

echo -e "${GREEN}📦 Latest version: $LATEST_VERSION${NC}"

# Construct download URL
if [ "$OS" = "windows" ]; then
    BINARY_NAME="codekeeper-${OS}-${ARCH}.exe"
    INSTALL_NAME="codekeeper.exe"
else
    BINARY_NAME="codekeeper-${OS}-${ARCH}"
    INSTALL_NAME="codekeeper"
fi

DOWNLOAD_URL="https://github.com/zvika-finally/ai-codekeeper/releases/download/${LATEST_VERSION}/${BINARY_NAME}"

echo -e "${BLUE}⬇️  Downloading ${BINARY_NAME}...${NC}"

# Create temporary directory
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Download binary
if ! curl -fsSL "$DOWNLOAD_URL" -o "$TEMP_DIR/$INSTALL_NAME"; then
    echo -e "${RED}❌ Failed to download binary from $DOWNLOAD_URL${NC}"
    exit 1
fi

# Make executable
chmod +x "$TEMP_DIR/$INSTALL_NAME"

# Determine install location
if [ -w "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
elif [ -w "$HOME/.local/bin" ]; then
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
else
    INSTALL_DIR="$HOME/bin"
    mkdir -p "$INSTALL_DIR"
fi

# Install binary
echo -e "${BLUE}📦 Installing to $INSTALL_DIR...${NC}"
mv "$TEMP_DIR/$INSTALL_NAME" "$INSTALL_DIR/$INSTALL_NAME"

# Verify installation
if command -v "$INSTALL_NAME" >/dev/null 2>&1; then
    echo -e "${GREEN}✅ Installation successful!${NC}"
    echo -e "${GREEN}🎉 Finally AI CodeKeeper is now available as '$INSTALL_NAME'${NC}"
    echo ""
    echo -e "${YELLOW}📋 Quick start:${NC}"
    echo "  $INSTALL_NAME --version"
    echo "  $INSTALL_NAME --help"
    echo "  $INSTALL_NAME cursor setup"
else
    echo -e "${YELLOW}⚠️  Installation completed but '$INSTALL_NAME' not found in PATH${NC}"
    echo -e "${YELLOW}📝 Add $INSTALL_DIR to your PATH:${NC}"
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
    echo ""
    echo -e "${YELLOW}📋 Then try:${NC}"
    echo "  $INSTALL_DIR/$INSTALL_NAME --version"
fi