# 🚀 CI/CD Development Guide

> **Complete guide for developing with k8s-cli's automated CI/CD pipeline**

---

## 📋 Overview

This guide shows how to develop effectively with k8s-cli's automated CI/CD pipeline, including best practices, workflow integration, and troubleshooting.

---

## 🛠️ Development Workflow

### 🌟 Quick Start Development Cycle

```bash
# 1. Setup development environment
make -f Makefile.dev dev-setup

# 2. Start development with auto-rebuild
make -f Makefile.dev watch

# 3. Make changes and test locally
make -f Makefile.dev pre-commit

# 4. Commit with conventional commits
git commit -m "feat: add new component detection"

# 5. Push and let CI handle the rest
git push origin feature/component-detection
```

### 🎯 Pre-commit Checklist

Always run before committing:

```bash
# Complete quality check suite
make -f Makefile.dev pre-commit

# This runs:
# ✅ Code formatting (gofmt, goimports)
# ✅ Linting (golangci-lint)  
# ✅ Security scan (govulncheck)
# ✅ Unit tests with coverage
# ✅ Build verification
```

---

## 🔄 CI/CD Integration

### 📊 Understanding Workflow Triggers

#### Automatic Triggers

**Pull Request Creation/Update:**
```bash
git checkout -b feature/new-feature
git commit -m "feat: implement new feature"
git push origin feature/new-feature
# → Creates PR → Triggers quality checks
```

**Push to Main (Automatic Release):**
```bash
# Feature commit triggers minor release
git commit -m "feat: add Helm component detection"
git push origin main
# → Auto-release: v2.1.0

# Bug fix triggers patch release  
git commit -m "fix: resolve memory leak in scanning"
git push origin main
# → Auto-release: v2.0.1

# Breaking change triggers major release
git commit -m "feat!: restructure CLI interface

BREAKING CHANGE: command structure has changed"
git push origin main
# → Auto-release: v3.0.0
```

#### Manual Triggers

**Manual Release:**
```bash
# GitHub → Actions → "🚀 Automated Release" → Run workflow
# Select: patch/minor/major
# Options: skip checks (for hotfixes), create GitHub release
```

### 🎯 Commit Message Strategy

Use [Conventional Commits](https://conventionalcommits.org/) for automatic release management:

```bash
# Types that trigger releases:
git commit -m "feat: new feature"     # → Minor release (1.1.0)
git commit -m "fix: bug fix"         # → Patch release (1.0.1)  
git commit -m "perf: performance"    # → Patch release (1.0.1)

# Breaking changes:
git commit -m "feat!: breaking change"           # → Major release (2.0.0)
git commit -m "feat: change

BREAKING CHANGE: details"                        # → Major release (2.0.0)

# Types that don't trigger releases:
git commit -m "docs: update readme"    # → No release
git commit -m "chore: update deps"     # → No release
git commit -m "test: add tests"        # → No release
git commit -m "ci: fix workflow"       # → No release
git commit -m "style: formatting"      # → No release
git commit -m "refactor: cleanup"      # → No release
```

---

## 🧪 Testing Strategy

### 🎯 Test Levels

#### 1. Unit Tests (Local & CI)
```bash
# Run unit tests locally
make -f Makefile.dev test

# With coverage
make -f Makefile.dev test-coverage
```

#### 2. Integration Tests (CI)
```bash
# Run integration tests locally (requires cluster)
make -f Makefile.dev test-integration

# CI automatically:
# - Sets up kind cluster
# - Runs integration tests
# - Cleans up resources
```

#### 3. Cross-Platform Tests (CI Only)
```bash
# CI automatically tests on:
# - Ubuntu 22.04 (Linux AMD64)
# - macOS 13 (Darwin AMD64)  
# - Windows 2022 (Windows AMD64)
# - Multiple Go versions (1.22, 1.23, 1.24)
```

### 🔧 Writing CI-Friendly Tests

#### Cross-Platform Compatibility
```go
func TestVersionFlag(t *testing.T) {
    // Skip integration tests on Windows in CI
    if runtime.GOOS == "windows" && os.Getenv("GITHUB_ACTIONS") == "true" {
        t.Skip("Skipping integration test on Windows in CI")
    }
    
    // Use cross-platform binary path
    binaryPath := getBinaryPath()
    
    // Auto-build if needed
    if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
        buildCmd := exec.Command("make", "-f", "Makefile.dev", "build")
        buildCmd.Dir = "../"
        if err := buildCmd.Run(); err != nil {
            t.Fatalf("Failed to build binary: %v", err)
        }
    }
}

func getBinaryPath() string {
    binaryName := "k8s-cli"
    if runtime.GOOS == "windows" {
        binaryName += ".exe"
    }
    return filepath.Join("..", "bin", binaryName)
}
```

#### Kubeconfig-Dependent Tests
```go
func TestKubernetesIntegration(t *testing.T) {
    // Skip if no kubeconfig available
    if _, err := os.Stat(os.Getenv("HOME") + "/.kube/config"); os.IsNotExist(err) {
        t.Skip("Skipping test: no kubeconfig found")
    }
    
    // Test continues...
}
```

---

## 🎭 Development Scenarios

### 🌟 Feature Development

#### 1. Start New Feature
```bash
# Create feature branch
git checkout -b feature/helm-detection

# Set up development environment
make -f Makefile.dev dev-setup
```

#### 2. Develop with Auto-reload
```bash
# Start file watcher for auto-rebuild
make -f Makefile.dev watch

# Or smart watch (different actions based on file type)
make -f Makefile.dev smart-watch
```

#### 3. Implement Feature
```go
// pkg/kubernetes/components.go
func DetectHelmReleases(client kubernetes.Interface) ([]Component, error) {
    // Implementation...
}
```

#### 4. Add Tests
```go
// pkg/kubernetes/components_test.go
func TestDetectHelmReleases(t *testing.T) {
    // Test implementation...
}
```

#### 5. Verify Quality
```bash
# Run full quality checks
make -f Makefile.dev pre-commit

# Check specific areas
make -f Makefile.dev lint
make -f Makefile.dev test
make -f Makefile.dev security-scan
```

#### 6. Commit and Push
```bash
git add .
git commit -m "feat: add Helm release detection

- Implement Helm secret scanning
- Add version extraction from labels  
- Include release status metadata
- Add comprehensive tests"

git push origin feature/helm-detection
```

#### 7. Create Pull Request
- CI automatically runs quality checks
- Review automated PR comments
- Address any issues found
- Merge when checks pass

#### 8. Automatic Release
```bash
# After merge to main, automatic minor release triggered
# New version: v2.1.0 (due to "feat:" commit)
```

### 🐛 Bug Fix Workflow

#### 1. Identify and Reproduce
```bash
git checkout -b hotfix/memory-leak

# Reproduce issue locally
./bin/k8s-cli all --namespace kube-system
# → Memory usage increases continuously
```

#### 2. Fix Issue
```go
// pkg/kubernetes/client.go
func (c *Client) GetComponents() ([]Component, error) {
    // Add proper resource cleanup
    defer resourceCleanup()
    
    // Fixed implementation...
}
```

#### 3. Verify Fix
```bash
# Test the fix locally
make -f Makefile.dev test
./bin/k8s-cli all --namespace kube-system
# → Memory usage stable
```

#### 4. Emergency Release (if critical)
```bash
git commit -m "fix: resolve memory leak in component scanning

Critical fix for production environments.
Memory usage now properly released after scanning."

git push origin hotfix/memory-leak
# → Create PR with expedited review
# → Or use manual release with "skip checks" for emergency
```

### 📦 Release Management

#### 1. Planned Release
```bash
# Prepare release branch
git checkout -b release/v2.1.0

# Update documentation
git commit -m "docs: update installation guide for v2.1.0"

# Use manual release workflow:
# GitHub → Actions → "🚀 Automated Release"
# Type: minor, Skip checks: false, Create release: true
```

#### 2. Pre-release Testing
```bash
# Test release process without creating actual release
make -f Makefile.dev release-dry-run

# Review what would be released
./scripts/bump-version.sh current
git log --oneline $(git describe --tags --abbrev=0)..HEAD
```

#### 3. Release Validation
```bash
# After release, verify:
# ✅ GitHub release created
# ✅ Assets uploaded (Linux, macOS, Windows binaries)
# ✅ Changelog updated
# ✅ Version tags correct
# ✅ Release notes generated
```

---

## 🛡️ Code Quality Standards

### 🎯 Linting Requirements

All code must pass:

```bash
# Go linting
golangci-lint run --timeout=5m

# Format checking
gofmt -s -d .
goimports -d .

# Security scanning
govulncheck ./...

# Module verification
go mod verify
```

### 📊 Testing Requirements

```bash
# Minimum test coverage: 70%
make -f Makefile.dev test-coverage

# All tests must pass
make -f Makefile.dev test

# Integration tests (when applicable)
make -f Makefile.dev test-integration
```

### 🔒 Security Standards

```bash
# No security vulnerabilities
govulncheck ./...

# Dependency audit
go list -json -deps ./... | nancy sleuth

# No secrets in code
git-secrets --scan
```

---

## 🔧 Advanced Development

### 🎯 Custom Make Targets

Create custom development workflows:

```makefile
# Makefile.dev additions
dev-feature: ## Start feature development cycle
	@make -f Makefile.dev clean
	@make -f Makefile.dev build
	@make -f Makefile.dev test
	@echo "🚀 Ready for feature development!"

quick-test: ## Run fast tests only
	@go test -short ./...

deep-test: ## Run all tests including slow ones
	@make -f Makefile.dev test
	@make -f Makefile.dev test-integration
```

### 🔄 CI/CD Customization

#### Skip CI for Documentation
```bash
git commit -m "docs: update readme [skip ci]"
```

#### Custom Workflow Triggers
```yaml
# .github/workflows/custom.yml
on:
  push:
    paths:
      - 'pkg/**'
      - 'cmd/**'
  pull_request:
    paths:
      - '**.go'
```

#### Environment-Specific Testing
```bash
# Test against specific Kubernetes versions
export KUBERNETES_VERSION=v1.28.0
make -f Makefile.dev test-integration
```

---

## 📊 Monitoring & Metrics

### 🎯 CI/CD Health Indicators

Monitor these metrics:

```bash
# Success rate (target: >95%)
GitHub → Insights → Actions

# Build time (target: <10 minutes)
# Test coverage (target: >70%)  
# Security scan results (target: 0 vulnerabilities)
```

### 📈 Development Velocity

Track improvements:

```bash
# Feature delivery time
# Bug fix turnaround
# Release frequency
# CI/CD reliability
```

---

## 🆘 Common Development Issues

### 🐛 Local Development Problems

#### Go Version Mismatch
```bash
# Check Go version
go version
# Should be: go version go1.24+ 

# Update if needed
brew install go@1.24  # macOS
# or download from https://golang.org/dl/
```

#### Module Issues
```bash
# Clean and refresh modules
go clean -modcache
go mod download
go mod tidy
```

#### Build Issues
```bash
# Clean build environment
make -f Makefile.dev clean
make -f Makefile.dev build

# Verify tools are installed
make -f Makefile.dev dev-setup
```

### 🔄 CI/CD Issues

See [Troubleshooting Guide](TROUBLESHOOTING.md) for detailed CI/CD issue resolution.

---

## 🚀 Best Practices Summary

### ✅ Do's

1. **Always run pre-commit checks:**
   ```bash
   make -f Makefile.dev pre-commit
   ```

2. **Use conventional commits:**
   ```bash
   git commit -m "feat: descriptive message"
   ```

3. **Test cross-platform compatibility:**
   ```bash
   GOOS=windows go build .
   GOOS=darwin go build .
   ```

4. **Write CI-friendly tests:**
   ```go
   if os.Getenv("CI") == "true" {
       // CI-specific test adjustments
   }
   ```

5. **Monitor workflow health:**
   - Set up GitHub notifications
   - Review failed builds promptly
   - Keep dependencies updated

### ❌ Don'ts

1. **Don't skip quality checks:**
   ```bash
   # ❌ Avoid this
   git commit -m "quick fix" --no-verify
   ```

2. **Don't ignore CI failures:**
   - Always investigate and fix
   - Don't merge failing PRs

3. **Don't use non-conventional commits for releases:**
   ```bash
   # ❌ Won't trigger proper release
   git commit -m "added new feature"
   
   # ✅ Triggers minor release
   git commit -m "feat: add new feature"
   ```

4. **Don't hardcode paths or OS-specific code:**
   ```go
   // ❌ OS-specific
   path := "/usr/local/bin/k8s-cli"
   
   // ✅ Cross-platform
   path := filepath.Join("bin", "k8s-cli")
   ```

---

## 🎉 Success Indicators

You're developing effectively when you see:

### ✅ Green CI/CD Pipeline
- All workflows passing consistently
- Fast feedback on pull requests
- Automatic releases working smoothly

### 📊 High Code Quality
- Test coverage >70%
- No linting errors
- Zero security vulnerabilities

### 🚀 Smooth Development Experience
- Fast local development cycle
- Reliable builds across platforms
- Clear feedback from automation

### 📈 Productive Team Workflow
- Quick feature delivery
- Minimal manual intervention
- Consistent code quality

---

**🎊 Congratulations!**

You now have a comprehensive understanding of developing with k8s-cli's CI/CD pipeline. This automated system will help you:

- 🔄 Develop faster with immediate feedback
- 🛡️ Maintain high code quality automatically  
- 🚀 Release confidently with automated testing
- 📈 Scale development as your team grows

The CI/CD pipeline handles the complex automation so you can focus on building great features! 🚀

---

## 📚 Quick Reference

### 🔗 Essential Commands
```bash
# Setup development
make -f Makefile.dev dev-setup

# Development cycle  
make -f Makefile.dev pre-commit

# Auto-rebuild
make -f Makefile.dev watch

# Release testing
make -f Makefile.dev release-dry-run
```

### 📖 Related Documentation
- [GitHub Actions Guide](GITHUB_ACTIONS.md)
- [Troubleshooting Guide](TROUBLESHOOTING.md) 
- [Make Commands Guide](MAKE_GUIDE.md)
- [Release Automation](RELEASE_AUTOMATION.md)