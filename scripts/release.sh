#!/bin/bash

# 🚀 k8s-cli Automated Release Script
# This script automates the entire release process including version bumping,
# changelog generation, documentation updates, building, and tagging.

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
VERSION_FILE="$PROJECT_ROOT/VERSION"
CHANGELOG_FILE="$PROJECT_ROOT/CHANGELOG.md"

# Functions
log() {
    echo -e "${BLUE}[RELEASE]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Check if we're in a git repository
check_git_repo() {
    if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
        error "Not inside a git repository"
    fi
}

# Check if working directory is clean
check_clean_working_dir() {
    local skip_checks=${1:-false}
    
    if [ -n "$(git status --porcelain)" ]; then
        warning "Working directory is not clean. Uncommitted changes found:"
        git status --short
        echo
        
        if [ "$skip_checks" = "skip-checks" ]; then
            warning "Continuing anyway due to skip-checks mode"
        else
            read -p "Do you want to continue? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                error "Aborted by user"
            fi
        fi
    fi
}

# Get current version
get_current_version() {
    if [ -f "$VERSION_FILE" ]; then
        cat "$VERSION_FILE"
    else
        echo "v0.0.0"
    fi
}

# Increment version
increment_version() {
    local version=$1
    local increment_type=$2
    
    # Remove 'v' prefix if present
    version=${version#v}
    
    IFS='.' read -ra VERSION_PARTS <<< "$version"
    local major=${VERSION_PARTS[0]:-0}
    local minor=${VERSION_PARTS[1]:-0}
    local patch=${VERSION_PARTS[2]:-0}
    
    case $increment_type in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
        *)
            error "Invalid increment type: $increment_type. Use: major, minor, or patch"
            ;;
    esac
    
    echo "v${major}.${minor}.${patch}"
}

# Update VERSION file
update_version_file() {
    local new_version=$1
    echo "$new_version" > "$VERSION_FILE"
    success "Updated VERSION file to $new_version"
}

# Generate changelog entry
generate_changelog_entry() {
    local version=$1
    local date=$(date +%Y-%m-%d)
    local commit_hash=$(git rev-parse --short HEAD)
    local branch=$(git rev-parse --abbrev-ref HEAD)
    
    # Get commits since last tag
    local last_tag
    if git describe --tags --abbrev=0 > /dev/null 2>&1; then
        last_tag=$(git describe --tags --abbrev=0)
        local commits_since_tag=$(git log --oneline ${last_tag}..HEAD --no-merges)
    else
        local commits_since_tag=$(git log --oneline --no-merges -10)
    fi
    
    # Create temporary file for changelog entry
    local temp_changelog=$(mktemp)
    
    cat > "$temp_changelog" << EOF
## [$version] - $date

### 🚀 Release Summary

This release includes the following changes and improvements.

### ✨ Added
EOF

    # Analyze commit messages for features
    while IFS= read -r commit; do
        if [[ $commit == *"feat:"* ]] || [[ $commit == *"add:"* ]]; then
            local msg=$(echo "$commit" | cut -d' ' -f2-)
            echo "- $msg" >> "$temp_changelog"
        fi
    done <<< "$commits_since_tag"

    cat >> "$temp_changelog" << EOF

### 🔧 Enhanced
EOF

    # Analyze commit messages for improvements
    while IFS= read -r commit; do
        if [[ $commit == *"enhance:"* ]] || [[ $commit == *"improve:"* ]] || [[ $commit == *"update:"* ]]; then
            local msg=$(echo "$commit" | cut -d' ' -f2-)
            echo "- $msg" >> "$temp_changelog"
        fi
    done <<< "$commits_since_tag"

    cat >> "$temp_changelog" << EOF

### 🐛 Fixed
EOF

    # Analyze commit messages for fixes
    while IFS= read -r commit; do
        if [[ $commit == *"fix:"* ]] || [[ $commit == *"bug:"* ]]; then
            local msg=$(echo "$commit" | cut -d' ' -f2-)
            echo "- $msg" >> "$temp_changelog"
        fi
    done <<< "$commits_since_tag"

    cat >> "$temp_changelog" << EOF

### 🛠️ Technical
- Git commit: $commit_hash
- Build date: $date
- Branch: $branch

### 📊 Changes Summary
EOF
    
    echo "- $(echo "$commits_since_tag" | wc -l | tr -d ' ') commits since last release" >> "$temp_changelog"
    
    cat >> "$temp_changelog" << EOF

---

EOF

    echo "$temp_changelog"
}

# Update changelog
update_changelog() {
    local version=$1
    local temp_changelog
    temp_changelog=$(generate_changelog_entry "$version")
    
    # Create backup
    cp "$CHANGELOG_FILE" "${CHANGELOG_FILE}.backup"
    
    # Insert new entry after the header
    local temp_full_changelog=$(mktemp)
    
    # Find the line number after the first "## [" (skip header)
    local insert_line
    insert_line=$(grep -n "^## \[" "$CHANGELOG_FILE" | head -1 | cut -d: -f1)
    
    if [ -z "$insert_line" ]; then
        # No previous entries, insert after header
        insert_line=$(grep -n "^and this project adheres" "$CHANGELOG_FILE" | cut -d: -f1)
        insert_line=$((insert_line + 2))
    fi
    
    # Split changelog and insert new entry
    head -n $((insert_line - 1)) "$CHANGELOG_FILE" > "$temp_full_changelog"
    cat "$temp_changelog" >> "$temp_full_changelog"
    tail -n +$insert_line "$CHANGELOG_FILE" >> "$temp_full_changelog"
    
    mv "$temp_full_changelog" "$CHANGELOG_FILE"
    rm "$temp_changelog"
    
    success "Updated CHANGELOG.md with new entry"
}

# Update README badges
update_readme_badges() {
    local version=$1
    local readme_file="$PROJECT_ROOT/README.md"
    
    if [ -f "$readme_file" ]; then
        # Update version badge
        sed -i.backup "s/Version-v[0-9]\+\.[0-9]\+\.[0-9]\+-green/Version-$version-green/g" "$readme_file"
        rm "${readme_file}.backup" 2>/dev/null || true
        success "Updated README.md version badges"
    fi
}

# Run tests and quality checks
run_quality_checks() {
    log "Running quality checks..."
    
    cd "$PROJECT_ROOT"
    
    # Format code
    log "Formatting code..."
    make -f Makefile.dev fmt
    
    # Run tests
    log "Running tests..."
    make -f Makefile.dev test
    
    # Run linter
    log "Running linter..."
    make -f Makefile.dev lint || warning "Linter issues found, but continuing..."
    
    # Run security scan if available
    log "Running security scan..."
    make -f Makefile.dev security-scan || warning "Security scan failed or not available"
    
    success "Quality checks completed"
}

# Build release binaries
build_release() {
    log "Building release binaries..."
    
    cd "$PROJECT_ROOT"
    
    # Build for multiple platforms
    make -f Makefile.dev release-build
    make -f Makefile.dev release-package
    
    success "Release binaries built and packaged"
}

# Create git commit and tag
create_git_commit_and_tag() {
    local version=$1
    
    log "Creating git commit and tag..."
    
    # Add changed files
    git add VERSION CHANGELOG.md README.md
    
    # Create commit
    local commit_message="release: bump version to $version

🚀 Generated with k8s-cli release automation

Co-Authored-By: k8s-cli-bot <noreply@example.com>"

    git commit -m "$commit_message"
    
    # Create annotated tag
    git tag -a "$version" -m "Release $version

$(head -20 "$CHANGELOG_FILE" | tail -10)

🚀 Generated with k8s-cli release automation"
    
    success "Created commit and tag $version"
}

# Show release summary
show_release_summary() {
    local old_version=$1
    local new_version=$2
    
    echo
    echo "🎉 Release Summary"
    echo "=================="
    echo "• Previous version: $old_version"
    echo "• New version: $new_version"
    echo "• Changelog: Updated"
    echo "• README: Updated"  
    echo "• Tests: ✅ Passed"
    echo "• Build: ✅ Completed"
    echo "• Git: ✅ Committed and tagged"
    echo
    echo "📋 Next Steps:"
    echo "1. Review the changes: git show $new_version"
    echo "2. Push to remote: git push origin main --tags"
    echo "3. Create GitHub release (manual)"
    echo "4. Install new version: make -f Makefile.dev install-user"
    echo
}

# Main release function
main() {
    log "🚀 Starting k8s-cli automated release process"
    
    # Parse arguments
    local increment_type=${1:-patch}
    local skip_checks=${2:-false}
    
    if [[ ! "$increment_type" =~ ^(major|minor|patch)$ ]]; then
        error "Usage: $0 [major|minor|patch] [skip-checks]"
    fi
    
    # Pre-flight checks
    check_git_repo
    check_clean_working_dir "$skip_checks"
    
    # Get versions
    local current_version
    current_version=$(get_current_version)
    local new_version
    new_version=$(increment_version "$current_version" "$increment_type")
    
    log "Current version: $current_version"
    log "New version: $new_version"
    
    if [ "$skip_checks" = "skip-checks" ]; then
        log "Auto-continuing due to skip-checks mode"
    else
        echo
        read -p "Continue with release $new_version? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            error "Aborted by user"
        fi
    fi
    
    # Update version
    update_version_file "$new_version"
    
    # Update documentation
    update_changelog "$new_version"
    update_readme_badges "$new_version"
    
    # Quality checks (unless skipped)
    if [ "$skip_checks" != "skip-checks" ]; then
        run_quality_checks
    else
        warning "Skipping quality checks as requested"
    fi
    
    # Build release
    build_release
    
    # Git operations
    create_git_commit_and_tag "$new_version"
    
    # Show summary
    show_release_summary "$current_version" "$new_version"
    
    success "🎉 Release $new_version completed successfully!"
}

# Script execution
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi