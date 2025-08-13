# 🚀 Release Automation Guide

> **Complete automation system for k8s-cli releases with version management, changelog generation, and quality checks**

---

## 📋 Overview

The k8s-cli release automation system provides a comprehensive, automated workflow for creating releases with:

- 🔢 **Automated version bumping** (semantic versioning)
- 📝 **Smart changelog generation** from git commits
- 🔍 **Pre-release quality checks** and validations
- 🏗️ **Multi-platform builds** and packaging
- 📊 **Documentation updates** (README badges, etc.)
- 🏷️ **Git tagging and commit management**

---

## ⚡ Quick Start

### 🎯 One-Command Release
```bash
# Patch release (2.0.1 -> 2.0.2)
make -f Makefile.dev release-patch

# Minor release (2.0.1 -> 2.1.0)  
make -f Makefile.dev release-minor

# Major release (2.0.1 -> 3.0.0)
make -f Makefile.dev release-major
```

### 🔍 Preview Before Release
```bash
# See what would be released
make -f Makefile.dev release-dry-run

# Check if ready for release
./scripts/pre-release-checks.sh
```

---

## 🛠️ Release Automation Components

### 📁 Files Structure
```
scripts/
├── release.sh              # Main release automation script
├── bump-version.sh          # Version management utility
├── generate-changelog.sh    # Advanced changelog generator
├── pre-release-checks.sh    # Quality validation suite
└── changelog-template.md    # Changelog template
```

### 🎯 Makefile Targets
```bash
# Release commands
make -f Makefile.dev release-patch      # Patch release
make -f Makefile.dev release-minor      # Minor release  
make -f Makefile.dev release-major      # Major release

# Fast releases (skip some checks)
make -f Makefile.dev release-patch-fast
make -f Makefile.dev release-minor-fast
make -f Makefile.dev release-major-fast

# Utilities
make -f Makefile.dev release-dry-run    # Preview release
make -f Makefile.dev push-release       # Push to remote
```

---

## 🔢 Version Management

### 📈 Semantic Versioning
The system follows [Semantic Versioning](https://semver.org/):

- **MAJOR** (x.0.0): Breaking changes, major new features
- **MINOR** (x.y.0): New features, backwards compatible  
- **PATCH** (x.y.z): Bug fixes, minor improvements

### 🛠️ Version Commands
```bash
# Show current version
./scripts/bump-version.sh current

# Preview next version
./scripts/bump-version.sh next patch    # Shows v2.0.2
./scripts/bump-version.sh next minor    # Shows v2.1.0
./scripts/bump-version.sh next major    # Shows v3.0.0

# Set specific version
./scripts/bump-version.sh set v2.5.0

# Show version history
./scripts/bump-version.sh history

# Validate version format
./scripts/bump-version.sh validate v2.0.1
```

---

## 📝 Smart Changelog Generation

### 🤖 Automated Analysis
The changelog generator analyzes git commits and automatically categorizes them:

#### Commit Patterns
- **Features**: `feat:`, `add:`, `new:` → ✨ Added section
- **Enhancements**: `enhance:`, `improve:`, `update:`, `refactor:` → 🔧 Enhanced section
- **Bug Fixes**: `fix:`, `bug:`, `hotfix:` → 🐛 Fixed section
- **Documentation**: `docs:`, `doc:`, `documentation:` → 📚 Documentation section

#### Smart Metrics
- Commit count since last release
- Contributors count
- Files changed
- Build information (Go version, binary size)
- Test coverage and status

### 📋 Example Generated Entry
```markdown
## [v2.0.2] - 2025-08-13

### 🚀 Patch Release

This release includes targeted improvements and fixes.

### ✨ Added
- **Enhanced component detection with Helm integration**
- **CLI version management with --version flag**

### 🔧 Enhanced
- **Improved namespace scanning for comprehensive discovery** 
- **Better error handling and user feedback**

### 🛠️ Technical Details
- **Git commit**: abc123
- **Total commits**: 12 commits since last release
- **Contributors**: 2
- **Test coverage**: 85%
```

---

## 🔍 Pre-release Quality Checks

### ✅ Comprehensive Validation Suite
The automation runs extensive checks before any release:

#### 🏗️ Environment Checks
- Git repository status
- Working directory cleanliness  
- Branch status and remote sync
- Go environment and dependencies

#### 📋 Code Quality Checks
- Code compilation
- Unit tests execution
- Test coverage analysis
- Linting (golangci-lint/go vet)
- Security vulnerability scan

#### 📄 Documentation Checks
- VERSION file validity
- Changelog entries
- Documentation completeness
- Makefile targets

#### 🎯 Build Validation
- Release build capability
- Binary size validation
- Multi-platform build support

### 📊 Example Check Output
```
🎯 Pre-release Check Summary
============================
✅ Passed: 15
⚠️  Warnings: 2  
❌ Failed: 0

📊 Overall Score: 100%

🎉 Ready for release!
   (with 2 warnings)
```

---

## 🚀 Release Process Workflow

### 🔄 Complete Automation Flow

1. **📋 Pre-flight Checks**
   - Validate git repository status
   - Check working directory cleanliness
   - Verify current version

2. **🔢 Version Management**  
   - Increment version (patch/minor/major)
   - Update VERSION file
   - Validate new version format

3. **📝 Documentation Updates**
   - Generate smart changelog entry
   - Update README badges
   - Refresh version references

4. **🔍 Quality Assurance**
   - Code formatting and linting
   - Unit tests execution
   - Security vulnerability scan
   - Build validation

5. **🏗️ Build Process**
   - Multi-platform binary builds
   - Release packaging (tar.gz, zip)
   - Build artifact validation

6. **📦 Git Operations**
   - Commit all changes
   - Create annotated release tag
   - Prepare for remote push

7. **📋 Release Summary**
   - Show comprehensive release report
   - Provide next steps guidance
   - Display push commands

---

## 📋 Usage Examples

### 🎯 Standard Release Workflow

#### 1. Check Release Readiness
```bash
# Preview what will be released
make -f Makefile.dev release-dry-run

# Output:
# Current version: v2.0.1
# Current branch: main  
# Uncommitted changes: None
# Recent commits:
#   abc123 feat: enhanced component detection
#   def456 fix: improved error handling
```

#### 2. Run Quality Checks
```bash
# Comprehensive pre-release validation
./scripts/pre-release-checks.sh

# Quick checks (skip slow tests)
./scripts/pre-release-checks.sh skip-slow
```

#### 3. Create Release
```bash
# Create patch release (most common)
make -f Makefile.dev release-patch

# The system will:
# - Increment version: v2.0.1 -> v2.0.2  
# - Generate changelog entry
# - Update documentation
# - Run quality checks
# - Build release binaries
# - Create git commit and tag
```

#### 4. Push to Remote
```bash
# Push release to remote repository
make -f Makefile.dev push-release

# Manual alternative:
git push origin main --tags
```

### 🏃 Fast Release (Skip Some Checks)
```bash
# Quick patch release with minimal checks
make -f Makefile.dev release-patch-fast

# Useful for:
# - Hotfixes
# - Documentation updates  
# - Minor corrections
```

### 🎛️ Manual Version Control
```bash
# Check current version
./scripts/bump-version.sh current

# Set specific version
./scripts/bump-version.sh set v2.5.0

# Preview next versions
./scripts/bump-version.sh next patch  # v2.0.2
./scripts/bump-version.sh next minor  # v2.1.0
./scripts/bump-version.sh next major  # v3.0.0
```

---

## ⚙️ Configuration & Customization

### 📁 File Locations
```bash
PROJECT_ROOT/
├── VERSION                    # Current version
├── CHANGELOG.md              # Release history
├── scripts/
│   ├── release.sh           # Main automation
│   ├── bump-version.sh      # Version management  
│   ├── generate-changelog.sh # Changelog generation
│   ├── pre-release-checks.sh # Quality checks
│   └── changelog-template.md # Template
└── Makefile.dev             # Release targets
```

### 🎛️ Customization Options

#### Changelog Template
Edit `scripts/changelog-template.md` to customize:
- Section headers and emojis
- Business impact descriptions  
- Usage examples format
- Technical details layout

#### Commit Patterns
Modify `scripts/generate-changelog.sh` to change:
- Commit message patterns
- Category classifications
- Smart descriptions

#### Quality Checks
Customize `scripts/pre-release-checks.sh` for:
- Additional validation rules
- Custom quality thresholds  
- Environment-specific checks

---

## 🔧 Troubleshooting

### ❗ Common Issues

#### 🚫 "Working directory not clean"
```bash
# Check status
git status

# Commit changes
git add .
git commit -m "prepare for release"

# Or stash temporarily  
git stash
```

#### 🚫 "Tests failing"
```bash
# Run tests manually
make -f Makefile.dev test

# Check specific failures
go test -v ./...

# Fix issues and retry
```

#### 🚫 "Version format invalid"
```bash
# Check VERSION file
cat VERSION

# Fix format (must be vX.Y.Z)
echo "v2.0.1" > VERSION
```

#### 🚫 "Git tag already exists"
```bash
# Check existing tags
git tag -l "v*" | tail -10

# Delete local tag if needed
git tag -d v2.0.1

# Delete remote tag if needed  
git push origin :refs/tags/v2.0.1
```

### 🛠️ Recovery Actions

#### Rollback Release
```bash
# Reset to previous commit (before release)
git reset --hard HEAD^

# Remove release tag
git tag -d v2.0.2

# Restore VERSION file
git checkout HEAD^ -- VERSION
```

#### Partial Release Fix
```bash
# Fix issues without full release
git add .
git commit -m "fix: post-release corrections"

# Update tag to point to new commit  
git tag -d v2.0.2
git tag v2.0.2
```

---

## 📊 Best Practices

### 🎯 Release Planning

#### When to Release
- **Patch**: Bug fixes, minor improvements, security updates
- **Minor**: New features, enhancements, non-breaking changes
- **Major**: Breaking changes, major architecture updates

#### Pre-release Checklist
- [ ] All tests passing
- [ ] Documentation updated
- [ ] Breaking changes documented
- [ ] Performance impact assessed
- [ ] Security review completed

#### Commit Message Guidelines
Use conventional commit format for best changelog generation:
```bash
feat: add new component detection
fix: resolve memory leak in metrics
docs: update installation guide
enhance: improve error messages
```

### 🚀 Automation Tips

#### Regular Release Cadence
- **Weekly patches**: Bug fixes and minor improvements
- **Monthly minors**: New features and enhancements  
- **Quarterly majors**: Major updates and breaking changes

#### Quality Gates
- Always run full checks for minor/major releases
- Use fast releases only for urgent patches
- Monitor test coverage trends
- Review security scan results

#### Documentation Maintenance
- Keep changelog entries detailed
- Update examples with new features
- Maintain version compatibility matrix
- Document breaking changes clearly

---

## 🤖 CI/CD Integration

### 🔄 GitHub Actions Example
```yaml
name: Automated Release
on:
  push:
    branches: [ main ]
  workflow_dispatch:
    inputs:
      release_type:
        description: 'Release type'
        required: true
        default: 'patch'
        type: choice
        options: [ 'patch', 'minor', 'major' ]

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
      with:
        fetch-depth: 0
        
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
        
    - name: Run Pre-release Checks
      run: ./scripts/pre-release-checks.sh
      
    - name: Create Release  
      run: make -f Makefile.dev release-${{ inputs.release_type }}
      
    - name: Push Release
      run: make -f Makefile.dev push-release
```

### 🎯 Integration Benefits
- Automated quality gates
- Consistent release process
- Audit trail maintenance
- Multi-environment deployment

---

## 📚 Additional Resources

### 🔗 Related Documentation
- [Make Commands Guide](MAKE_GUIDE.md) - Complete Makefile reference
- [Development Guide](DEVELOPMENT.md) - Development workflows
- [Examples](EXAMPLES.md) - Usage examples and tutorials

### 🛠️ Required Tools
- **Git** (2.0+): Version control
- **Go** (1.21+): Build and testing
- **Make**: Automation execution
- **golangci-lint**: Code quality (optional)
- **govulncheck**: Security scanning (optional)

### 🎓 Learning Resources
- [Semantic Versioning](https://semver.org/)
- [Conventional Commits](https://conventionalcommits.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [Go Release Best Practices](https://golang.org/doc/modules/release-workflow)

---

**🎉 Happy Releasing!** 

The automation system is designed to make releases safe, consistent, and efficient. For questions or improvements, please open an issue or contribute to the automation scripts.