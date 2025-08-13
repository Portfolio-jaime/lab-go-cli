# 🔧 k8s-cli Make Commands Guide

> **Complete reference for building, testing, and installing k8s-cli**

---

## 📋 Quick Reference

### ⚡ Essential Commands
```bash
# Build and install (recommended for users)
make -f Makefile.dev install-user

# Complete development cycle
make -f Makefile.dev dev-cycle

# Get help and see all available commands
make -f Makefile.dev help
```

---

## 🚀 Installation Commands

### 📍 User Installation (Recommended)
```bash
# Install to ~/bin directory (no sudo required)
make -f Makefile.dev install-user

# Verify installation
k8s-cli --version
k8s-cli version
```

**Prerequisites:**
- Ensure `$HOME/bin` is in your PATH
- If not, add this to your `~/.bashrc` or `~/.zshrc`:
  ```bash
  export PATH="$HOME/bin:$PATH"
  ```

### 🌍 System-Wide Installation
```bash
# Install to /usr/local/bin (requires sudo)
make -f Makefile.dev install

# Verify installation
k8s-cli --version
```

### 🗑️ Uninstallation
```bash
# Remove from user directory
make -f Makefile.dev uninstall-user

# Remove from system
make -f Makefile.dev uninstall
```

---

## 🔨 Build Commands

### 🏗️ Standard Build
```bash
# Basic build - creates ./bin/k8s-cli
make -f Makefile.dev build

# Run directly from bin
./bin/k8s-cli --version
./bin/k8s-cli version
```

### 🐛 Debug Build
```bash
# Build with debug information
make -f Makefile.dev dev-build

# Creates ./bin/k8s-cli-debug
./bin/k8s-cli-debug version
```

### 🌎 Multi-Platform Release Build
```bash
# Build for all platforms (Linux, macOS, Windows)
make -f Makefile.dev release-build

# Package releases
make -f Makefile.dev release-package

# Check ./bin/release/ directory for binaries
ls ./bin/release/
```

---

## 🧪 Testing Commands

### ✅ Basic Testing
```bash
# Run all unit tests
make -f Makefile.dev test

# Run tests with coverage
make -f Makefile.dev test-coverage

# View coverage report (opens browser)
open coverage.html  # macOS
xdg-open coverage.html  # Linux
```

### 🔄 Continuous Testing
```bash
# Auto-run tests when files change
make -f Makefile.dev test-watch
```

### 🌐 Integration & E2E Testing
```bash
# Run integration tests (requires cluster)
make -f Makefile.dev test-integration

# Run end-to-end tests
make -f Makefile.dev test-e2e
```

### ⚡ Performance Testing
```bash
# Run benchmarks
make -f Makefile.dev benchmark
```

---

## 📊 Code Quality Commands

### 🎨 Code Formatting
```bash
# Format all Go code
make -f Makefile.dev fmt
```

### 🔍 Code Analysis
```bash
# Run linter
make -f Makefile.dev lint

# Run go vet
make -f Makefile.dev vet

# Run all quality checks
make -f Makefile.dev check-all
```

### 🛡️ Security Scanning
```bash
# Run security vulnerability scan
make -f Makefile.dev security-scan
```

---

## 🚀 Development Workflow Commands

### 🔄 Complete Development Cycle
```bash
# Format + Test + Build (recommended before commits)
make -f Makefile.dev dev-cycle
```

### 👀 File Watching & Auto-rebuild
```bash
# Auto-rebuild on file changes (requires 'air' or file watching tools)
make -f Makefile.dev watch

# Alternative file watcher
make -f Makefile.dev dev-watch

# Auto-update with full cycle on changes
make -f Makefile.dev auto-update

# Smart watch (incremental builds based on changed files)
make -f Makefile.dev smart-watch
```

### 🔍 Pre-commit & Pre-push Checks
```bash
# Pre-commit checks (format, vet, lint, test)
make -f Makefile.dev pre-commit

# Pre-push checks (includes integration tests)
make -f Makefile.dev pre-push
```

---

## 🏃 Running Commands

### 🎯 Quick Test Runs
```bash
# Show CLI help
make -f Makefile.dev run-help

# Test version command
make -f Makefile.dev run-version

# Run complete analysis (requires cluster)
make -f Makefile.dev run-all
```

### 📊 Feature Testing
```bash
# Test metrics command
make -f Makefile.dev run-metrics

# Test cost analysis
make -f Makefile.dev run-cost

# Test workload analysis
make -f Makefile.dev run-workload

# Test logs analysis
make -f Makefile.dev run-logs

# Test export functionality
make -f Makefile.dev run-export
```

### 🎬 Demo
```bash
# Run interactive demo
make -f Makefile.dev demo
```

---

## 📚 Documentation Commands

### 📖 Generate Documentation
```bash
# Generate command documentation
make -f Makefile.dev docs-generate

# Update all documentation
make -f Makefile.dev docs-update

# Validate documentation links
make -f Makefile.dev docs-validate
```

### 🌐 Serve Documentation
```bash
# Serve docs locally at http://localhost:8000
make -f Makefile.dev docs-serve
```

---

## 🧹 Maintenance Commands

### 🗑️ Cleanup
```bash
# Clean build artifacts
make -f Makefile.dev clean
```

### 📦 Dependency Management
```bash
# Update Go dependencies
make -f Makefile.dev deps-update

# Tidy Go modules
make -f Makefile.dev deps-tidy

# Audit dependencies for vulnerabilities
make -f Makefile.dev deps-audit
```

### ℹ️ Version Information
```bash
# Show version information
make -f Makefile.dev version
```

---

## 🚀 Release Automation Commands

### 📦 Automated Releases
```bash
# Create patch release (2.0.1 -> 2.0.2)
make -f Makefile.dev release-patch

# Create minor release (2.0.1 -> 2.1.0)  
make -f Makefile.dev release-minor

# Create major release (2.0.1 -> 3.0.0)
make -f Makefile.dev release-major
```

### ⚡ Fast Releases (Skip Some Checks)
```bash
# Quick patch release
make -f Makefile.dev release-patch-fast

# Quick minor release
make -f Makefile.dev release-minor-fast

# Quick major release
make -f Makefile.dev release-major-fast
```

### 🔍 Release Utilities
```bash
# Preview what would be released
make -f Makefile.dev release-dry-run

# Push release to remote
make -f Makefile.dev push-release
```

### 🛠️ Manual Release Tools
```bash
# Version management
./scripts/bump-version.sh current          # Show current version
./scripts/bump-version.sh next patch       # Preview next patch version
./scripts/bump-version.sh bump patch       # Bump patch version
./scripts/bump-version.sh set v2.5.0       # Set specific version

# Quality checks
./scripts/pre-release-checks.sh            # Full validation suite
./scripts/pre-release-checks.sh skip-slow  # Quick checks only

# Changelog generation
./scripts/generate-changelog.sh v2.0.2 v2.0.1  # Generate changelog entry
```

---

## 🐳 Docker Commands

### 🏗️ Docker Build
```bash
# Build Docker image
make -f Makefile.dev docker-build

# Run CLI in Docker container
make -f Makefile.dev docker-run
```

---

## 🏗️ Development Environment Setup

### 🛠️ First-Time Setup
```bash
# Setup complete development environment
make -f Makefile.dev dev-setup

# This installs:
# - golangci-lint (linter)
# - air (hot reload)
# - goimports (import formatting)
# - swag (API documentation)
```

---

## ⚡ Recommended Workflows

### 🆕 For New Users (First Time)
```bash
# 1. Setup development environment
make -f Makefile.dev dev-setup

# 2. Run complete development cycle
make -f Makefile.dev dev-cycle

# 3. Install CLI
make -f Makefile.dev install-user

# 4. Test installation
k8s-cli --version
k8s-cli version  # (requires cluster)
```

### 🔄 Daily Development
```bash
# Start file watcher for auto-rebuild
make -f Makefile.dev watch

# In another terminal, test changes:
./bin/k8s-cli --version
./bin/k8s-cli version
```

### 🚢 Before Committing Code
```bash
# Run all pre-commit checks
make -f Makefile.dev pre-commit

# If successful, commit your changes
git add .
git commit -m "your changes"
```

### 📦 Creating Releases
```bash
# 1. Check if ready for release
make -f Makefile.dev release-dry-run

# 2. Run comprehensive quality checks
./scripts/pre-release-checks.sh

# 3. Create automated release
make -f Makefile.dev release-patch  # or release-minor, release-major

# 4. Push to remote
make -f Makefile.dev push-release
```

### 🔧 Manual Version Management
```bash
# Check current version
./scripts/bump-version.sh current

# Preview next version
./scripts/bump-version.sh next patch

# Set specific version
./scripts/bump-version.sh set v2.5.0
```

### 📦 For Distribution
```bash
# Build release packages
make -f Makefile.dev release-package

# Check release artifacts
ls ./bin/release/
```

---

## ❗ Troubleshooting

### 🚨 Common Issues

#### Build Fails
```bash
# Clean and rebuild
make -f Makefile.dev clean
make -f Makefile.dev build
```

#### Tests Fail
```bash
# Run specific test package
go test -v ./cmd/
go test -v ./pkg/kubernetes/
```

#### Missing Dependencies
```bash
# Reinstall development tools
make -f Makefile.dev dev-setup

# Update and tidy dependencies  
make -f Makefile.dev deps-tidy
```

#### File Watcher Not Working
```bash
# Install file watching tool manually
# macOS (using Homebrew):
brew install fswatch

# Linux:
sudo apt-get install inotify-tools  # Ubuntu/Debian
sudo yum install inotify-tools      # RHEL/CentOS
```

### 🔧 Environment Variables

```bash
# Override default paths
export GOPATH=/path/to/go
export PATH=$GOPATH/bin:$PATH

# Override kubeconfig for testing
export KUBECONFIG=/path/to/your/kubeconfig
```

---

## 📝 Make Target Categories

### 🎯 **Essential** (Most Common)
- `install-user` - Install to user directory
- `dev-cycle` - Complete dev cycle  
- `build` - Standard build
- `test` - Run tests
- `clean` - Clean artifacts

### 🔧 **Development**  
- `dev-setup` - First-time setup
- `watch` - Auto-rebuild
- `pre-commit` - Pre-commit checks
- `fmt` - Format code
- `lint` - Run linter

### 🚀 **Advanced**
- `release-build` - Multi-platform build
- `test-integration` - Integration tests
- `security-scan` - Security analysis
- `docs-generate` - Generate documentation
- `docker-build` - Docker image build

---

## 🎉 Success Indicators

After running installation commands, you should see:

```bash
$ k8s-cli --version
k8s-cli version v2.0.1
Git commit: abc123
Built: 2025-08-13T10:30:00
Go version: go1.24.6
OS/Arch: darwin/arm64

$ k8s-cli version
🔍 Analyzing Kubernetes cluster...

📊 Cluster Version Information:
[Kubernetes cluster details]

🔧 Installed Components:
   Searching in all namespaces for components and Helm releases...
   
   Found X components:
[Component table with Helm releases]
```

---

**📧 Questions?** Check the [Development Guide](DEVELOPMENT.md) or open an issue on GitHub.