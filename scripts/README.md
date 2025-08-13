# 🛠️ k8s-cli Automation Scripts

> **Collection of automation scripts for release management, version control, and quality assurance**

---

## 📁 Scripts Overview

### 🚀 [`release.sh`](release.sh)
**Main release automation script**

Complete end-to-end release process including version bumping, changelog generation, quality checks, building, and git operations.

```bash
# Usage
./scripts/release.sh [major|minor|patch] [skip-checks]

# Examples
./scripts/release.sh patch           # Create patch release with full checks
./scripts/release.sh minor           # Create minor release
./scripts/release.sh major skip-checks  # Major release, skip quality checks
```

**What it does:**
- ✅ Validates git repository status
- 🔢 Increments version automatically
- 📝 Generates smart changelog entries
- 📊 Updates documentation (README badges, etc.)
- 🔍 Runs comprehensive quality checks
- 🏗️ Builds multi-platform releases
- 🏷️ Creates git commits and tags

---

### 🔢 [`bump-version.sh`](bump-version.sh)
**Standalone version management utility**

Comprehensive version management with validation, history, and flexible operations.

```bash
# Usage
./scripts/bump-version.sh [COMMAND] [OPTIONS]

# Examples
./scripts/bump-version.sh current        # Show current version
./scripts/bump-version.sh bump patch     # Increment patch version
./scripts/bump-version.sh set v2.5.0     # Set specific version
./scripts/bump-version.sh next minor     # Preview next minor version
./scripts/bump-version.sh history        # Show version history from git
./scripts/bump-version.sh validate v2.0.1 # Validate version format
```

**Features:**
- 📈 Semantic versioning support
- 🔍 Version format validation
- 📊 Version history from git tags
- 👀 Preview mode (show without changing)
- 🔒 Dry-run capability

---

### 📝 [`generate-changelog.sh`](generate-changelog.sh)
**Advanced changelog generator**

Smart changelog generation with commit analysis, metrics, and templating.

```bash
# Usage  
./scripts/generate-changelog.sh [VERSION] [CURRENT_VERSION]

# Examples
./scripts/generate-changelog.sh v2.0.2 v2.0.1  # Generate entry for v2.0.2
./scripts/generate-changelog.sh v2.1.0         # Auto-detect current version
```

**Smart Features:**
- 🤖 **Commit Analysis**: Automatically categorizes commits by patterns
- 📊 **Metrics Collection**: Counts commits, contributors, files changed
- 🏗️ **Build Information**: Go version, binary size, test coverage
- 📋 **Business Impact**: Generates DevOps, SRE, and DX benefits
- 🎨 **Template System**: Uses customizable templates

**Commit Patterns:**
- `feat:`, `add:`, `new:` → ✨ Added section
- `enhance:`, `improve:`, `update:`, `refactor:` → 🔧 Enhanced section  
- `fix:`, `bug:`, `hotfix:` → 🐛 Fixed section
- `docs:`, `doc:`, `documentation:` → 📚 Documentation section

---

### 🔍 [`pre-release-checks.sh`](pre-release-checks.sh)
**Comprehensive quality validation suite**

Extensive pre-release checks covering code quality, security, documentation, and build validation.

```bash
# Usage
./scripts/pre-release-checks.sh [skip-slow]

# Examples  
./scripts/pre-release-checks.sh              # Full validation suite
./scripts/pre-release-checks.sh skip-slow    # Skip slow checks (tests, coverage, etc.)
```

**Validation Categories:**

#### 🏗️ Environment Checks
- Git repository status and cleanliness
- Branch status and remote synchronization
- Go environment and dependencies validation
- Working directory state

#### 📋 Code Quality Checks  
- Code compilation verification
- Unit tests execution and results
- Test coverage analysis and reporting
- Code linting (golangci-lint/go vet)
- Security vulnerability scanning (govulncheck)

#### 📄 Documentation Checks
- VERSION file format validation
- Changelog entries verification
- Documentation completeness assessment
- Makefile targets validation

#### 🎯 Build Validation
- Release build capability testing
- Binary size and artifact validation
- Multi-platform build support verification

**Example Output:**
```
🎯 Pre-release Check Summary
============================
✅ Passed: 15
⚠️  Warnings: 2
❌ Failed: 0

📊 Overall Score: 100%
🎉 Ready for release! (with 2 warnings)
```

---

### 📋 [`changelog-template.md`](changelog-template.md)
**Changelog entry template**

Template file used by the changelog generator with placeholders for dynamic content.

**Template Variables:**
- `{{VERSION}}` - Release version
- `{{DATE}}` - Release date
- `{{RELEASE_TYPE}}` - Major/Minor/Patch
- `{{GIT_COMMIT}}` - Git commit hash
- `{{BUILD_DATE}}` - Build timestamp
- `{{COMMIT_COUNT}}` - Number of commits since last release
- `{{CONTRIBUTORS}}` - Number of contributors
- `{{TEST_STATUS}}` - Test execution status
- `{{COVERAGE}}` - Test coverage percentage

---

## 🎯 Integration with Makefile

All scripts are integrated into the main Makefile.dev for easy access:

### 🚀 Release Commands
```bash
make -f Makefile.dev release-patch       # → ./scripts/release.sh patch
make -f Makefile.dev release-minor       # → ./scripts/release.sh minor  
make -f Makefile.dev release-major       # → ./scripts/release.sh major
make -f Makefile.dev release-dry-run     # → Preview release info
make -f Makefile.dev push-release        # → Push release to remote
```

### ⚡ Fast Release Commands
```bash
make -f Makefile.dev release-patch-fast  # → ./scripts/release.sh patch skip-checks
make -f Makefile.dev release-minor-fast  # → ./scripts/release.sh minor skip-checks
make -f Makefile.dev release-major-fast  # → ./scripts/release.sh major skip-checks
```

---

## 🔧 Prerequisites

### Required Tools
- **Git** (2.0+): Version control operations
- **Go** (1.21+): Building, testing, and tools
- **Make**: Makefile execution
- **bash** (4.0+): Script execution
- **bc**: Math calculations (for coverage)

### Optional Tools (Enhanced Features)
- **golangci-lint**: Advanced code linting
- **govulncheck**: Security vulnerability scanning  
- **air**: File watching and hot reload

### Installation of Optional Tools
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest

# Install air (file watcher)
go install github.com/air-verse/air@latest
```

---

## ⚙️ Configuration

### 🛠️ Environment Variables
```bash
# Override default paths
export PROJECT_ROOT=/path/to/project
export VERSION_FILE=/path/to/VERSION  
export CHANGELOG_FILE=/path/to/CHANGELOG.md

# Customize behavior
export SKIP_QUALITY_CHECKS=true     # Skip quality checks
export DRY_RUN=true                  # Preview mode
export NO_GIT_OPERATIONS=true       # Skip git commits/tags
```

### 📁 File Structure
```
scripts/
├── release.sh                 # Main release automation
├── bump-version.sh            # Version management utility  
├── generate-changelog.sh      # Smart changelog generator
├── pre-release-checks.sh      # Quality validation suite
├── changelog-template.md      # Changelog template
└── README.md                  # This documentation
```

---

## 🎯 Usage Examples

### 🚀 Complete Release Workflow
```bash
# 1. Check current status
./scripts/bump-version.sh current
# Output: v2.0.1

# 2. Preview what next version would be
./scripts/bump-version.sh next patch  
# Output: v2.0.2

# 3. Run pre-release validation
./scripts/pre-release-checks.sh
# Output: Comprehensive validation results

# 4. Create automated release
./scripts/release.sh patch
# Output: Complete release process with v2.0.2

# 5. Push to remote (manual)
git push origin main --tags
```

### 🔧 Manual Operations
```bash
# Just bump version (no release)
./scripts/bump-version.sh bump patch

# Generate changelog only
./scripts/generate-changelog.sh v2.0.2 v2.0.1 > new-entry.md

# Run specific checks
./scripts/pre-release-checks.sh skip-slow

# Preview release without changes
./scripts/release.sh patch dry-run
```

### ⚡ Quick Operations
```bash
# Fast patch release (skip slow checks)
./scripts/release.sh patch skip-checks

# Check version history
./scripts/bump-version.sh history

# Validate version format
./scripts/bump-version.sh validate v2.5.0
```

---

## 🐛 Troubleshooting

### Common Issues

#### 🚫 Permission Denied
```bash
# Make scripts executable
chmod +x scripts/*.sh
```

#### 🚫 Git Repository Issues
```bash
# Initialize git repository
git init
git remote add origin <repository-url>

# Commit current changes
git add .
git commit -m "initial commit"
```

#### 🚫 Version File Missing
```bash
# Create VERSION file
echo "v1.0.0" > VERSION
git add VERSION
git commit -m "add version file"
```

#### 🚫 Dependencies Missing
```bash
# Install Go dependencies
go mod tidy

# Install development tools
make -f Makefile.dev dev-setup
```

### Debug Mode
```bash
# Run scripts with debug output
bash -x ./scripts/release.sh patch

# Or set debug flag
export DEBUG=true
./scripts/release.sh patch
```

---

## 🤝 Contributing

### Adding New Checks
1. Edit `pre-release-checks.sh`
2. Add new check function following the pattern:
```bash
check_new_feature() {
    log "Checking new feature..."
    
    if [[ condition ]]; then
        success "New feature check passed"
        return 0
    else
        error "New feature check failed"
        return 1
    fi
}
```
3. Call function in `main()` section

### Customizing Templates
1. Edit `changelog-template.md`
2. Add new placeholder variables `{{NEW_VAR}}`
3. Update `generate-changelog.sh` to replace placeholders

### Extending Version Management
1. Edit `bump-version.sh`
2. Add new commands in `main()` function
3. Follow existing patterns for validation and output

---

## 📚 Additional Resources

- [Release Automation Guide](../docs/RELEASE_AUTOMATION.md) - Complete automation documentation
- [Make Commands Guide](../docs/MAKE_GUIDE.md) - Makefile reference
- [Development Guide](../docs/DEVELOPMENT.md) - Development workflows
- [Semantic Versioning](https://semver.org/) - Versioning guidelines
- [Conventional Commits](https://conventionalcommits.org/) - Commit message format

---

**🎉 Happy Automating!**

These scripts are designed to make release management safe, consistent, and efficient. For questions or improvements, please contribute back to the project.