# 🚀 GitHub Actions Automation Guide

> **Complete CI/CD automation for k8s-cli with automated releases, quality checks, and multi-platform builds**

---

## 📋 Overview

The k8s-cli project includes comprehensive GitHub Actions automation that provides:

- 🚀 **Automated Releases** - Semantic versioning and changelog generation
- 🔍 **Quality Assurance** - Comprehensive testing and code quality checks  
- 🏗️ **Multi-platform Builds** - Automatic builds for Linux, macOS, and Windows
- 📊 **Continuous Integration** - Automated testing on every push and PR
- 📝 **Smart Documentation** - Auto-generated release notes and PR summaries

---

## 🛠️ Workflows Overview

### 📁 Workflow Files

```
.github/
├── workflows/
│   ├── release.yml        # 🚀 Automated release workflow
│   ├── pull-request.yml   # 🔍 PR quality checks
│   └── ci.yml            # 🔄 Continuous integration
├── ISSUE_TEMPLATE/
│   ├── bug_report.yml     # 🐛 Bug report template
│   ├── feature_request.yml # ✨ Feature request template
│   └── config.yml        # 📋 Issue template configuration
└── PULL_REQUEST_TEMPLATE.md # 🔄 PR template
```

---

## 🚀 Automated Release Workflow

### **File:** `.github/workflows/release.yml`

#### 🎯 Triggers

1. **Manual Trigger** (workflow_dispatch):
   ```
   Repository → Actions → "🚀 Automated Release" → Run workflow
   ```
   - Choose release type: patch/minor/major
   - Option to skip quality checks
   - Option to create GitHub release

2. **Automatic Trigger** (push to main):
   - Analyzes commit messages for conventional commit patterns
   - `feat:` → minor release
   - `fix:`, `perf:` → patch release  
   - `BREAKING CHANGE:` or `!:` → major release

#### 🔄 Release Process

1. **📋 Pre-flight Checks**
   - Git repository validation
   - Working directory status
   - Commit message analysis

2. **🔍 Quality Assurance**
   - Pre-release validation suite
   - Code compilation and testing
   - Security vulnerability scanning
   - Linting and code quality

3. **🔢 Version Management**
   - Automatic semantic version bumping
   - VERSION file updates
   - Git tag creation

4. **📝 Documentation Generation**
   - Smart changelog generation from commits
   - README badge updates
   - Release notes creation

5. **🏗️ Multi-platform Builds**
   - Linux (AMD64, ARM64)
   - macOS (AMD64, ARM64)
   - Windows (AMD64)
   - Binary packaging (tar.gz, zip)

6. **📦 GitHub Release Creation**
   - Automatic GitHub release
   - Release asset uploads
   - Generated release notes

#### 📊 Example Usage

**Manual Release:**
```bash
# Go to GitHub → Actions → "🚀 Automated Release"
# Select: Release type: "patch"
# Click: "Run workflow"
```

**Automatic Release via Commit:**
```bash
git commit -m "feat: add new component detection feature"
git push origin main
# → Triggers minor release automatically
```

---

## 🔍 Pull Request Quality Checks

### **File:** `.github/workflows/pull-request.yml`

#### 🎯 Triggers
- Pull requests to `main` or `develop`
- Pull request updates (synchronize)

#### 🔍 Quality Checks

1. **🧪 Comprehensive Testing**
   - Unit tests execution
   - Test coverage analysis  
   - Benchmark performance tests

2. **📊 Code Quality**
   - Go linting (golangci-lint)
   - Code formatting validation
   - Security vulnerability scanning

3. **🏗️ Build Verification**
   - Multi-platform build testing
   - Binary functionality verification

4. **📝 Conventional Commits**
   - Commit message format validation
   - Automatic categorization for changelog

5. **🔒 Dependency Audit**
   - Vulnerability scanning
   - Outdated dependency detection

6. **💔 Breaking Change Detection**
   - API change analysis
   - Breaking change warnings

#### 🤖 Automated PR Comments

The workflow automatically comments on PRs with:
- ✅ Quality check results summary
- 📊 Test coverage information  
- 🎯 Linting and security scan results
- 🚀 Next steps and recommendations

---

## 🔄 Continuous Integration

### **File:** `.github/workflows/ci.yml`

#### 🎯 Triggers
- Push to `main` or `develop`
- Scheduled runs (weekly dependency audit)

#### 🧪 Test Matrix

**Multi-Platform Testing:**
- **OS**: Ubuntu, macOS, Windows
- **Go Versions**: 1.22, 1.23, 1.24 (updated for k8s.io v0.33.3 compatibility)

**Multi-Architecture Builds:**
- Linux: AMD64, ARM64
- macOS: AMD64, ARM64 (Apple Silicon)
- Windows: AMD64

#### 📊 CI Jobs

1. **🧪 Test Matrix**
   - Cross-platform compatibility testing
   - Multiple Go version compatibility
   - Binary functionality verification

2. **🎯 Quality Checks**
   - Comprehensive linting
   - Security vulnerability scanning
   - Test coverage reporting

3. **🏗️ Build Matrix**
   - Multi-platform binary builds
   - Build artifact creation
   - Cross-compilation testing

4. **🔄 Integration Tests**
   - Kubernetes cluster setup (kind)
   - End-to-end testing
   - Real-world scenario testing

5. **📦 Dependency Audit**
   - Weekly dependency update checks
   - Security vulnerability monitoring
   - Module integrity verification

6. **🏃 Performance Benchmarks**
   - Performance regression detection
   - Benchmark result tracking

---

## 🎯 Configuration & Setup

### 🛠️ Required Secrets

No additional secrets required! The workflows use:
- `GITHUB_TOKEN` (automatically provided)
- Standard GitHub Actions permissions

### ⚙️ Workflow Permissions

```yaml
permissions:
  contents: write      # For creating releases and pushing tags
  pull-requests: write # For PR comments and reviews
  checks: write        # For check status updates
  issues: write        # For issue creation and updates
```

### 🎛️ Customization Options

#### Release Workflow Customization

**Automatic Release Triggers:**
Edit `.github/workflows/release.yml` to modify commit patterns:
```yaml
# Current patterns:
# feat: → minor release
# fix|perf: → patch release  
# BREAKING CHANGE: → major release
```

**Quality Check Configuration:**
```yaml
# Skip slow checks for automatic releases
skip_checks: ${{ github.event_name == 'push' }}
```

#### CI Workflow Customization

**Go Version Matrix:**
```yaml
strategy:
  matrix:
    go-version: ['1.22', '1.23', '1.24']  # Updated for k8s.io compatibility
```

**Platform Matrix:**
```yaml
strategy:
  matrix:
    os: [ubuntu-latest, macos-latest, windows-latest]
    goos: [linux, windows, darwin]
    goarch: [amd64, arm64]
```

---

## 🚀 Usage Guide

### 🎯 Creating Releases

#### 1. Manual Release (Recommended)
```bash
# Go to GitHub Repository
# Navigate to: Actions → "🚀 Automated Release"
# Click: "Run workflow"
# Select:
#   - Branch: main
#   - Release type: patch|minor|major
#   - Skip checks: false (for full validation)
#   - Create GitHub release: true
# Click: "Run workflow"
```

#### 2. Automatic Release via Commits
```bash
# Feature release (minor)
git commit -m "feat: add Helm component detection"
git push origin main

# Bug fix release (patch)  
git commit -m "fix: resolve memory leak in metrics collection"
git push origin main

# Breaking change release (major)
git commit -m "feat!: restructure CLI command interface

BREAKING CHANGE: command structure has changed"
git push origin main
```

### 🔍 Quality Assurance

#### Automatic PR Checks
- Create pull request → Automatic quality checks run
- View results in PR "Checks" tab
- Address any failing checks before merge

#### Manual Quality Validation
```bash
# Run locally before pushing
make -f Makefile.dev pre-commit
./scripts/pre-release-checks.sh
```

### 📊 Monitoring Workflows

#### Workflow Status
- **GitHub Repository** → **Actions** tab
- View running/completed workflows
- Check workflow logs and results
- Download build artifacts

#### Notifications
- **Watch repository** for release notifications
- **Email notifications** for workflow failures
- **Slack/Discord integration** (if configured)

---

## 🎭 Workflow Scenarios

### 🚀 Feature Development Workflow

1. **🌿 Create feature branch**
   ```bash
   git checkout -b feature/helm-detection
   ```

2. **💻 Develop feature**
   ```bash
   # Make changes
   git commit -m "feat: implement Helm release detection"
   ```

3. **🔄 Create Pull Request**
   - PR triggers quality checks automatically
   - Review automated PR comments
   - Address any issues found

4. **✅ Merge to main**
   ```bash
   # After PR approval and merge
   # Automatic minor release triggered by "feat:" commit
   ```

### 🐛 Bug Fix Workflow

1. **🔥 Identify bug**
   ```bash
   git checkout -b hotfix/memory-leak
   ```

2. **🔧 Fix issue**
   ```bash
   git commit -m "fix: resolve memory leak in component scanning"
   ```

3. **⚡ Fast release** (if urgent)
   - Use manual workflow with "skip checks" enabled
   - Or merge and let automatic release trigger

### 📦 Release Preparation

1. **📋 Pre-release checklist**
   ```bash
   # Run comprehensive checks
   make -f Makefile.dev release-dry-run
   ./scripts/pre-release-checks.sh
   ```

2. **📝 Update documentation**
   ```bash
   git commit -m "docs: update installation guide for v2.1.0"
   ```

3. **🚀 Create release**
   - Use manual workflow for controlled release
   - Review generated changelog before release
   - Verify multi-platform builds

---

## 🔧 Troubleshooting

### ❗ Common Issues

#### 🚫 Workflow Permission Errors
```yaml
# Add to workflow file
permissions:
  contents: write
  pull-requests: write
```

#### 🚫 Release Creation Fails
- Check if tag already exists
- Verify GITHUB_TOKEN permissions
- Ensure main branch is up to date

#### 🚫 Quality Checks Fail
- Review workflow logs in Actions tab
- Run checks locally: `./scripts/pre-release-checks.sh`
- Fix issues and re-run or commit fixes

#### 🚫 Build Matrix Failures
- Check Go version compatibility
- Verify platform-specific code issues
- Review build logs for specific errors

### 🛠️ Debug Actions

#### View Detailed Logs
```bash
# GitHub → Repository → Actions → Select workflow run
# Click on failing job
# Expand log sections to see details
```

#### Local Reproduction
```bash
# Reproduce CI environment locally
make -f Makefile.dev dev-cycle
./scripts/pre-release-checks.sh
GOOS=linux GOARCH=amd64 go build -o bin/k8s-cli-linux .
```

#### Manual Release Recovery
```bash
# If automated release fails mid-process
git tag -d v2.0.2  # Delete local tag
git push origin :refs/tags/v2.0.2  # Delete remote tag
# Fix issues and re-run workflow
```

---

## 🎯 Best Practices

### 📝 Commit Message Guidelines

Follow [Conventional Commits](https://conventionalcommits.org/):

```bash
# Feature (minor release)
git commit -m "feat(component): add StatefulSet detection"

# Bug fix (patch release)  
git commit -m "fix(metrics): resolve memory leak in collection"

# Breaking change (major release)
git commit -m "feat!: restructure command interface

BREAKING CHANGE: command structure changed"

# Documentation (no release)
git commit -m "docs: update installation guide"

# Chore (no release)
git commit -m "chore: update dependencies"
```

### 🚀 Release Strategy

#### Release Frequency
- **Patch releases**: Weekly for bug fixes
- **Minor releases**: Bi-weekly for features  
- **Major releases**: Quarterly for breaking changes

#### Release Validation
- Always run full quality checks for minor/major releases
- Use fast releases only for critical hotfixes
- Test releases in staging environment when possible

#### Communication
- Use GitHub Discussions for release planning
- Tag team members in critical release PRs
- Document breaking changes in release notes

### 🔒 Security Considerations

#### Dependency Management
- Weekly automated dependency audits
- Immediate security patch releases
- Regular Go version updates

#### Access Control
- Use branch protection rules
- Require PR reviews for main branch
- Enable automatic security fixes

#### Secrets Management
- No secrets required for basic workflows  
- Use GitHub secrets for external integrations
- Regular secret rotation for long-lived tokens

---

## 📚 Additional Resources

### 🔗 Related Documentation
- [Release Automation Guide](RELEASE_AUTOMATION.md)
- [Make Commands Guide](MAKE_GUIDE.md)
- [Development Guide](DEVELOPMENT.md)
- [Examples and Usage](EXAMPLES.md)
- [Troubleshooting Guide](TROUBLESHOOTING.md) - **🆕 NEW: Comprehensive error resolution guide**

### 🛠️ GitHub Actions Resources
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Workflow Syntax Reference](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)
- [Marketplace Actions](https://github.com/marketplace?type=actions)

### 🎓 Learning Resources
- [Conventional Commits](https://conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [GitHub Flow](https://guides.github.com/introduction/flow/)

---

## 🎉 Success Indicators

After setup, you should see:

### ✅ Successful Automation
- 🚀 Automatic releases on feature commits
- 🔍 PR quality checks on every pull request  
- 🏗️ Multi-platform builds without errors
- 📊 Coverage reports and quality metrics

### 📊 Workflow Health
- ✅ CI passes consistently across platforms
- 🎯 Quality checks catch issues early
- 🚀 Releases deploy without manual intervention
- 📝 Documentation stays up to date automatically

### 🎯 Developer Experience
- 🔄 Fast feedback on pull requests
- 📋 Clear guidance from automated comments
- 🎊 Smooth release process
- 🛡️ Confidence in code quality

---

**🎊 Congratulations!** 

Your k8s-cli project now has enterprise-grade automation that will:
- Save hours of manual work
- Ensure consistent quality  
- Enable rapid, safe releases
- Scale with your team growth

The GitHub Actions automation handles the complex parts so you can focus on building great features! 🚀