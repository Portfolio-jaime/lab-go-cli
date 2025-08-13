#!/bin/bash
# Development Environment Setup Script for k8s-cli

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 k8s-cli Development Environment Setup${NC}"
echo -e "${BLUE}=======================================${NC}"
echo ""

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to install Go tools
install_go_tool() {
    local tool=$1
    local package=$2
    
    if ! command_exists "$tool"; then
        echo -e "${YELLOW}Installing $tool...${NC}"
        go install "$package" || {
            echo -e "${RED}Failed to install $tool${NC}"
            return 1
        }
        echo -e "${GREEN}✓ $tool installed${NC}"
    else
        echo -e "${GREEN}✓ $tool already installed${NC}"
    fi
}

# Check if Go is installed
if ! command_exists go; then
    echo -e "${RED}❌ Go is not installed. Please install Go 1.24.5 or later.${NC}"
    echo -e "${YELLOW}Visit: https://golang.org/dl/${NC}"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "${BLUE}Go version: $GO_VERSION${NC}"

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌ Not in k8s-cli project directory. Please run from project root.${NC}"
    exit 1
fi

echo -e "${BLUE}📦 Installing development dependencies...${NC}"

# Install development tools
install_go_tool "golangci-lint" "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
install_go_tool "air" "github.com/air-verse/air@latest"
install_go_tool "goimports" "golang.org/x/tools/cmd/goimports@latest"
install_go_tool "govulncheck" "golang.org/x/vuln/cmd/govulncheck@latest"

# Download Go dependencies
echo -e "${YELLOW}📥 Downloading Go dependencies...${NC}"
go mod download
go mod tidy

# Create necessary directories
echo -e "${YELLOW}📁 Creating project directories...${NC}"
mkdir -p bin
mkdir -p tmp
mkdir -p exports
mkdir -p reports

# Set up git hooks (if git repo exists)
if [ -d ".git" ]; then
    echo -e "${YELLOW}🔗 Setting up git hooks...${NC}"
    mkdir -p .git/hooks
    
    # Pre-commit hook
    cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
echo "Running pre-commit checks..."
make pre-commit
EOF
    chmod +x .git/hooks/pre-commit
    echo -e "${GREEN}✓ Git pre-commit hook installed${NC}"
fi

# Build the CLI
echo -e "${YELLOW}🔨 Building k8s-cli...${NC}"
make build

# Test the build
if [ -f "bin/k8s-cli" ]; then
    echo -e "${GREEN}✓ Build successful${NC}"
    ./bin/k8s-cli --help > /dev/null
    echo -e "${GREEN}✓ CLI working correctly${NC}"
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

# Create development aliases (optional)
echo -e "${YELLOW}📝 Development setup complete!${NC}"
echo ""
echo -e "${BLUE}🎯 Quick Start Commands:${NC}"
echo -e "  ${GREEN}make dev-setup${NC}     - Re-run this setup"
echo -e "  ${GREEN}make build${NC}         - Build the CLI"
echo -e "  ${GREEN}make test${NC}          - Run tests"
echo -e "  ${GREEN}make watch${NC}         - Auto-rebuild on changes"
echo -e "  ${GREEN}make dev-cycle${NC}     - Format, test, and build"
echo -e "  ${GREEN}make auto-update${NC}   - Watch and auto-rebuild"
echo ""
echo -e "${BLUE}🔧 Development Workflow:${NC}"
echo -e "  1. ${YELLOW}make watch${NC}       - Start auto-rebuild"
echo -e "  2. Edit Go files in your editor"
echo -e "  3. CLI rebuilds automatically"
echo -e "  4. Test with: ${GREEN}./bin/k8s-cli <command>${NC}"
echo ""
echo -e "${BLUE}📚 Documentation:${NC}"
echo -e "  ${GREEN}docs/DEVELOPMENT.md${NC} - Development guide"
echo -e "  ${GREEN}docs/API.md${NC}         - API documentation"
echo -e "  ${GREEN}docs/EXAMPLES.md${NC}    - Usage examples"
echo ""

# Check for Kubernetes cluster
if command_exists kubectl; then
    if kubectl cluster-info >/dev/null 2>&1; then
        echo -e "${GREEN}✓ Kubernetes cluster available${NC}"
        echo -e "${YELLOW}Try: ./bin/k8s-cli all${NC}"
    else
        echo -e "${YELLOW}⚠ No Kubernetes cluster configured${NC}"
        echo -e "${YELLOW}To test with real cluster, configure kubectl${NC}"
    fi
else
    echo -e "${YELLOW}⚠ kubectl not found${NC}"
    echo -e "${YELLOW}Install kubectl to test against real clusters${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Development environment ready!${NC}"
echo -e "${BLUE}Happy coding! 🚀${NC}"