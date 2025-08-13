# 🔄 Pull Request

## 📋 Description

<!-- Provide a clear and concise description of the changes -->

**Type of change:**
- [ ] 🐛 Bug fix (non-breaking change which fixes an issue)
- [ ] ✨ New feature (non-breaking change which adds functionality)
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to change)
- [ ] 📚 Documentation update
- [ ] 🔧 Code refactoring
- [ ] ⚡ Performance improvement
- [ ] 🧪 Test addition or improvement

## 🎯 Related Issues

<!-- Link to related issues -->
Fixes #[issue_number]
Closes #[issue_number]
Related to #[issue_number]

## 🚀 Changes Made

<!-- Describe the changes in detail -->

### Added
- 

### Changed
- 

### Fixed
- 

### Removed
- 

## 🧪 Testing

<!-- Describe how you tested these changes -->

- [ ] Unit tests pass (`make -f Makefile.dev test`)
- [ ] Integration tests pass (if applicable)
- [ ] Manual testing completed
- [ ] Code builds successfully (`make -f Makefile.dev build`)

### Test Commands Run
```bash
# Add the commands you used to test
make -f Makefile.dev dev-cycle
./bin/k8s-cli --version
./bin/k8s-cli version
```

### Test Results
<!-- Paste test output or describe results -->
```
[Paste test output here if relevant]
```

## 📊 Performance Impact

<!-- If applicable, describe performance implications -->

- [ ] No performance impact
- [ ] Performance improvement (describe below)
- [ ] Performance regression (describe mitigation below)
- [ ] Performance impact unknown/needs investigation

## 💔 Breaking Changes

<!-- List any breaking changes and migration steps -->

- [ ] No breaking changes
- [ ] Breaking changes (describe below)

**Breaking Changes:**
<!-- Describe what breaks and how users should migrate -->

**Migration Guide:**
<!-- Provide migration steps for users -->

## 📋 Checklist

### Code Quality
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Code is properly commented
- [ ] No hardcoded values or secrets
- [ ] Error handling is appropriate

### Documentation
- [ ] Updated relevant documentation
- [ ] Added/updated code comments
- [ ] Updated CHANGELOG.md (if significant change)
- [ ] Updated examples (if applicable)

### Testing
- [ ] Added/updated unit tests
- [ ] Added/updated integration tests (if applicable)
- [ ] Tested on multiple platforms (if applicable)
- [ ] Verified backward compatibility

### Release Impact
- [ ] Version bump needed (patch/minor/major)
- [ ] Release notes impact considered
- [ ] Deployment impact assessed

## 🔍 Review Focus Areas

<!-- Guide reviewers on what to focus on -->

Please pay special attention to:
- 
- 
- 

## 📸 Screenshots

<!-- If applicable, add screenshots to help explain your changes -->

| Before | After |
|--------|-------|
|        |       |

## 🌍 Environment

**Development Environment:**
- OS: [e.g., macOS, Linux, Windows]
- Go version: [e.g., 1.21]
- k8s-cli version: [current version]

**Test Environment:**
- Kubernetes version: [e.g., v1.28.0]
- Cluster type: [e.g., minikube, kind, GKE]

## 📚 Additional Notes

<!-- Any additional information that reviewers should know -->

## 🔄 Post-Merge Actions

<!-- Actions needed after merge (if any) -->

- [ ] Update deployment
- [ ] Update documentation site
- [ ] Notify community
- [ ] Create follow-up issues

---

**📝 Note to Reviewers:**
- This PR will trigger automated quality checks
- CI will run tests on multiple platforms and Go versions
- Release automation may trigger based on commit message format

**🤖 Automated Checks:**
- [ ] CI build and test passed
- [ ] Quality checks passed
- [ ] Security scan passed
- [ ] Documentation validation passed