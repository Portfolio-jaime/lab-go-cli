#!/bin/bash

# 🔍 Pre-release Validation and Quality Checks for k8s-cli
# Comprehensive checks to ensure release readiness

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
VERSION_FILE="$PROJECT_ROOT/VERSION"
CHANGELOG_FILE="$PROJECT_ROOT/CHANGELOG.md"

# Tracking
CHECKS_PASSED=0
CHECKS_FAILED=0
WARNINGS=0
ERRORS=()
WARNINGS_LIST=()

log() {
    echo -e "${BLUE}[CHECK]${NC} $1"
}

success() {
    echo -e "${GREEN}[✅ PASS]${NC} $1"
    CHECKS_PASSED=$((CHECKS_PASSED + 1))
}

warning() {
    echo -e "${YELLOW}[⚠️  WARN]${NC} $1"
    WARNINGS=$((WARNINGS + 1))
    WARNINGS_LIST+=("$1")
}

error() {
    echo -e "${RED}[❌ FAIL]${NC} $1"
    CHECKS_FAILED=$((CHECKS_FAILED + 1))
    ERRORS+=("$1")
}

# Check if we're in a git repository
check_git_repository() {
    log "Checking git repository status..."
    
    if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
        error "Not inside a git repository"
        return 1
    fi
    
    success "Git repository detected"
    return 0
}

# Check working directory cleanliness
check_working_directory() {
    log "Checking working directory cleanliness..."
    
    local status
    status=$(git status --porcelain)
    
    if [ -n "$status" ]; then
        warning "Working directory has uncommitted changes:"
        echo "$status" | while IFS= read -r line; do
            echo "    $line"
        done
        return 1
    fi
    
    success "Working directory is clean"
    return 0
}

# Check branch status
check_branch_status() {
    log "Checking branch status..."
    
    local current_branch
    current_branch=$(git rev-parse --abbrev-ref HEAD)
    
    if [ "$current_branch" != "main" ] && [ "$current_branch" != "master" ]; then
        warning "Not on main/master branch (current: $current_branch)"
    else
        success "On main branch: $current_branch"
    fi
    
    # Check if branch is up to date with remote
    if git remote > /dev/null 2>&1; then
        local remote_branch="origin/$current_branch"
        
        if git show-ref --verify --quiet refs/remotes/$remote_branch; then
            local local_commit
            local_commit=$(git rev-parse HEAD)
            local remote_commit
            remote_commit=$(git rev-parse $remote_branch)
            
            if [ "$local_commit" != "$remote_commit" ]; then
                warning "Local branch is not up to date with remote"
                return 1
            else
                success "Branch is up to date with remote"
            fi
        else
            warning "Remote branch $remote_branch not found"
        fi
    fi
    
    return 0
}

# Check version file exists and is valid
check_version_file() {
    log "Checking VERSION file..."
    
    if [ ! -f "$VERSION_FILE" ]; then
        error "VERSION file not found: $VERSION_FILE"
        return 1
    fi
    
    local version
    version=$(cat "$VERSION_FILE")
    
    if [[ ! $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        error "Invalid version format in VERSION file: $version"
        return 1
    fi
    
    success "VERSION file is valid: $version"
    return 0
}

# Check changelog exists and has recent entry
check_changelog() {
    log "Checking CHANGELOG.md..."
    
    if [ ! -f "$CHANGELOG_FILE" ]; then
        error "CHANGELOG.md not found: $CHANGELOG_FILE"
        return 1
    fi
    
    local version
    version=$(cat "$VERSION_FILE" 2>/dev/null || echo "unknown")
    
    if ! grep -q "\[${version}\]" "$CHANGELOG_FILE"; then
        warning "No changelog entry found for version $version"
        return 1
    fi
    
    success "Changelog entry exists for version $version"
    return 0
}

# Check Go environment
check_go_environment() {
    log "Checking Go environment..."
    
    if ! command -v go > /dev/null 2>&1; then
        error "Go is not installed or not in PATH"
        return 1
    fi
    
    local go_version
    go_version=$(go version | awk '{print $3}')
    success "Go version: $go_version"
    
    # Check go mod
    if [ ! -f "$PROJECT_ROOT/go.mod" ]; then
        error "go.mod file not found"
        return 1
    fi
    
    success "go.mod file exists"
    
    # Check dependencies
    cd "$PROJECT_ROOT"
    if ! go mod verify > /dev/null 2>&1; then
        warning "Go module verification failed"
        return 1
    fi
    
    success "Go module dependencies verified"
    return 0
}

# Check code compilation
check_compilation() {
    log "Checking code compilation..."
    
    cd "$PROJECT_ROOT"
    
    if ! go build -o /tmp/k8s-cli-test . > /dev/null 2>&1; then
        error "Code compilation failed"
        return 1
    fi
    
    rm -f /tmp/k8s-cli-test
    success "Code compiles successfully"
    return 0
}

# Run unit tests
check_unit_tests() {
    log "Running unit tests..."
    
    cd "$PROJECT_ROOT"
    
    local test_output
    if ! test_output=$(go test ./... 2>&1); then
        error "Unit tests failed:"
        echo "$test_output" | tail -10
        return 1
    fi
    
    # Count tests
    local test_count
    test_count=$(echo "$test_output" | grep -c "^ok\|^PASS" || echo "0")
    
    success "Unit tests passed ($test_count packages)"
    return 0
}

# Check test coverage
check_test_coverage() {
    log "Checking test coverage..."
    
    cd "$PROJECT_ROOT"
    
    if go test -coverprofile=/tmp/coverage.out ./... > /dev/null 2>&1; then
        local coverage
        coverage=$(go tool cover -func=/tmp/coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
        
        if (( $(echo "$coverage < 50" | bc -l) )); then
            warning "Test coverage is low: ${coverage}%"
        else
            success "Test coverage: ${coverage}%"
        fi
        
        rm -f /tmp/coverage.out
    else
        warning "Could not calculate test coverage"
    fi
    
    return 0
}

# Run linting
check_linting() {
    log "Running code linting..."
    
    cd "$PROJECT_ROOT"
    
    # Try golangci-lint
    if command -v golangci-lint > /dev/null 2>&1; then
        if golangci-lint run > /dev/null 2>&1; then
            success "Linting passed (golangci-lint)"
        else
            warning "Linting issues found (golangci-lint)"
        fi
    # Fallback to go vet
    elif go vet ./... > /dev/null 2>&1; then
        success "Linting passed (go vet)"
    else
        warning "Linting issues found (go vet)"
    fi
    
    return 0
}

# Check security vulnerabilities
check_security() {
    log "Checking security vulnerabilities..."
    
    cd "$PROJECT_ROOT"
    
    if command -v govulncheck > /dev/null 2>&1; then
        if govulncheck ./... > /dev/null 2>&1; then
            success "No security vulnerabilities found"
        else
            warning "Security vulnerabilities detected"
        fi
    else
        warning "govulncheck not available (install with: go install golang.org/x/vuln/cmd/govulncheck@latest)"
    fi
    
    return 0
}

# Check documentation
check_documentation() {
    log "Checking documentation..."
    
    local docs_dir="$PROJECT_ROOT/docs"
    local readme="$PROJECT_ROOT/README.md"
    
    if [ ! -f "$readme" ]; then
        warning "README.md not found"
    else
        success "README.md exists"
    fi
    
    if [ ! -d "$docs_dir" ]; then
        warning "docs/ directory not found"
    else
        local doc_count
        doc_count=$(find "$docs_dir" -name "*.md" | wc -l | tr -d ' ')
        success "Documentation directory exists ($doc_count .md files)"
    fi
    
    return 0
}

# Check makefile targets
check_makefile() {
    log "Checking Makefile..."
    
    local makefile="$PROJECT_ROOT/Makefile.dev"
    
    if [ ! -f "$makefile" ]; then
        error "Makefile.dev not found"
        return 1
    fi
    
    success "Makefile.dev exists"
    
    # Check essential targets
    local essential_targets=("build" "test" "clean" "install-user")
    
    for target in "${essential_targets[@]}"; do
        if grep -q "^$target:" "$makefile"; then
            success "Makefile target '$target' found"
        else
            warning "Makefile target '$target' missing"
        fi
    done
    
    return 0
}

# Check release artifacts
check_release_artifacts() {
    log "Checking release build capability..."
    
    cd "$PROJECT_ROOT"
    
    if make -f Makefile.dev build > /dev/null 2>&1; then
        success "Release build successful"
        
        local binary="$PROJECT_ROOT/bin/k8s-cli"
        if [ -f "$binary" ]; then
            local size
            size=$(stat -f%z "$binary" 2>/dev/null || stat -c%s "$binary" 2>/dev/null || echo "0")
            local mb=$((size / 1024 / 1024))
            success "Binary created: ${mb}MB"
        else
            warning "Binary not found after build"
        fi
    else
        error "Release build failed"
        return 1
    fi
    
    return 0
}

# Check dependencies for known issues
check_dependencies() {
    log "Checking dependencies..."
    
    cd "$PROJECT_ROOT"
    
    # Check for outdated dependencies
    if command -v go > /dev/null 2>&1; then
        local outdated
        outdated=$(go list -u -m all 2>/dev/null | grep '\[' | wc -l | tr -d ' ')
        
        if [ "$outdated" -gt 0 ]; then
            warning "$outdated dependencies have updates available"
        else
            success "All dependencies are up to date"
        fi
    fi
    
    # Check for replace directives in go.mod
    if grep -q "replace " "$PROJECT_ROOT/go.mod"; then
        warning "go.mod contains replace directives (may indicate development dependencies)"
    fi
    
    return 0
}

# Generate summary report
generate_summary() {
    echo
    echo "🎯 Pre-release Check Summary"
    echo "============================"
    echo -e "✅ Passed: ${GREEN}$CHECKS_PASSED${NC}"
    echo -e "⚠️  Warnings: ${YELLOW}$WARNINGS${NC}"
    echo -e "❌ Failed: ${RED}$CHECKS_FAILED${NC}"
    echo
    
    if [ ${#ERRORS[@]} -gt 0 ]; then
        echo -e "${RED}❌ Critical Issues:${NC}"
        for error in "${ERRORS[@]}"; do
            echo "   • $error"
        done
        echo
    fi
    
    if [ ${#WARNINGS_LIST[@]} -gt 0 ]; then
        echo -e "${YELLOW}⚠️  Warnings:${NC}"
        for warn in "${WARNINGS_LIST[@]}"; do
            echo "   • $warn"
        done
        echo
    fi
    
    # Overall assessment
    local total_checks=$((CHECKS_PASSED + CHECKS_FAILED))
    local success_rate=0
    
    if [ $total_checks -gt 0 ]; then
        success_rate=$(( (CHECKS_PASSED * 100) / total_checks ))
    fi
    
    echo "📊 Overall Score: ${success_rate}%"
    echo
    
    if [ $CHECKS_FAILED -eq 0 ]; then
        echo -e "${GREEN}🎉 Ready for release!${NC}"
        if [ $WARNINGS -gt 0 ]; then
            echo -e "${YELLOW}   (with $WARNINGS warnings)${NC}"
        fi
        return 0
    else
        echo -e "${RED}🚫 NOT ready for release${NC}"
        echo "   Please fix critical issues before proceeding"
        return 1
    fi
}

# Main execution
main() {
    local skip_slow=${1:-false}
    
    echo "🚀 k8s-cli Pre-release Quality Checks"
    echo "====================================="
    echo
    
    # Quick checks
    check_git_repository || true
    check_working_directory || true
    check_branch_status || true
    check_version_file || true
    check_changelog || true
    check_go_environment || true
    check_compilation || true
    check_makefile || true
    check_documentation || true
    
    # Potentially slow checks
    if [ "$skip_slow" != "skip-slow" ]; then
        check_unit_tests || true
        check_test_coverage || true
        check_linting || true
        check_security || true
        check_release_artifacts || true
        check_dependencies || true
    else
        log "Skipping slow checks (unit tests, coverage, security, etc.)"
    fi
    
    generate_summary
}

# Execute if run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi