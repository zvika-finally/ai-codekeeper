#!/bin/bash

# Build and Test Script for Finally AI CodeKeeper
set -e

echo "🔨 Building Finally AI CodeKeeper..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go 1.21+ first.${NC}"
    echo "Visit: https://golang.org/dl/"
    exit 1
fi

echo -e "${GREEN}✅ Go found: $(go version)${NC}"

# Navigate to framework directory
cd "$(dirname "$0")"

echo -e "${BLUE}📦 Installing Go dependencies...${NC}"
go mod tidy

echo -e "${BLUE}🔨 Building binary...${NC}"
go build -o codekeeper cmd/codekeeper/main.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Build successful!${NC}"
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

echo -e "${BLUE}🧪 Testing basic functionality...${NC}"

# Test version
echo "Testing version command:"
./codekeeper --version

# Test help
echo -e "\nTesting help command:"
./codekeeper --help

# Test Cursor setup
echo -e "\nTesting Cursor setup:"
./codekeeper cursor --help

echo -e "\n${GREEN}🎉 Build complete!${NC}"
echo -e "${YELLOW}📋 Next steps:${NC}"
echo "1. Setup Cursor IDE: ./codekeeper cursor setup"
echo "2. Configure ecosystem: ./codekeeper mcp-ecosystem --domain fintech"
echo "3. Create new project: ./codekeeper create test-app"
echo "4. Check available commands: ./codekeeper --help"
echo ""
echo -e "${BLUE}📁 Binary location: $(pwd)/codekeeper${NC}"
echo "Add to PATH: export PATH=\"$(pwd):\$PATH\""