#!/bin/bash

# 🔢 Version Bumping Utility for k8s-cli
# Standalone script for version management

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

log() {
    echo -e "${BLUE}[VERSION]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Show usage
usage() {
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo
    echo "Commands:"
    echo "  current                    Show current version"
    echo "  bump <major|minor|patch>   Bump version by type"
    echo "  set <version>              Set specific version"
    echo "  next <major|minor|patch>   Show what next version would be"
    echo "  validate <version>         Validate version format"
    echo "  history                    Show version history from git tags"
    echo
    echo "Options:"
    echo "  --dry-run                  Show what would change without making changes"
    echo "  --no-prefix                Work with versions without 'v' prefix"
    echo "  --help                     Show this help message"
    echo
    echo "Examples:"
    echo "  $0 current                 # Shows v2.0.1"
    echo "  $0 bump patch              # v2.0.1 -> v2.0.2"
    echo "  $0 bump minor              # v2.0.1 -> v2.1.0"  
    echo "  $0 bump major              # v2.0.1 -> v3.0.0"
    echo "  $0 set v2.5.0              # Set to specific version"
    echo "  $0 next minor              # Show v2.1.0 (without changing)"
    echo "  $0 --dry-run bump patch    # Preview patch bump"
}

# Get current version
get_current_version() {
    if [ -f "$VERSION_FILE" ]; then
        cat "$VERSION_FILE"
    else
        echo "v0.0.0"
    fi
}

# Validate version format
validate_version() {
    local version=$1
    local no_prefix=${2:-false}
    
    if [ "$no_prefix" = "true" ]; then
        # Allow version without v prefix
        if [[ $version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            return 0
        fi
    else
        # Require v prefix
        if [[ $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            return 0
        fi
    fi
    
    return 1
}

# Parse version into components
parse_version() {
    local version=$1
    local prefix=""
    
    # Check if version has v prefix
    if [[ $version =~ ^v ]]; then
        prefix="v"
        version=${version#v}
    fi
    
    IFS='.' read -ra VERSION_PARTS <<< "$version"
    local major=${VERSION_PARTS[0]:-0}
    local minor=${VERSION_PARTS[1]:-0}
    local patch=${VERSION_PARTS[2]:-0}
    
    echo "$prefix|$major|$minor|$patch"
}

# Increment version
increment_version() {
    local current_version=$1
    local increment_type=$2
    
    local parsed
    parsed=$(parse_version "$current_version")
    IFS='|' read -ra PARTS <<< "$parsed"
    
    local prefix=${PARTS[0]}
    local major=${PARTS[1]}
    local minor=${PARTS[2]}  
    local patch=${PARTS[3]}
    
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
            error "Invalid increment type: $increment_type"
            ;;
    esac
    
    echo "${prefix}${major}.${minor}.${patch}"
}

# Set version in VERSION file
set_version_file() {
    local version=$1
    local dry_run=${2:-false}
    
    if [ "$dry_run" = "true" ]; then
        log "Would set VERSION file to: $version"
        return 0
    fi
    
    echo "$version" > "$VERSION_FILE"
    success "Updated VERSION file to $version"
}

# Show version history from git tags
show_version_history() {
    log "Version history from git tags:"
    echo
    
    if git tag --list "v*" --sort=-version:refname > /dev/null 2>&1; then
        local tags
        tags=$(git tag --list "v*" --sort=-version:refname | head -10)
        
        if [ -z "$tags" ]; then
            echo "No version tags found in git history"
            return 0
        fi
        
        echo "Recent versions:"
        while IFS= read -r tag; do
            local date
            date=$(git log -1 --format="%ai" "$tag" 2>/dev/null | cut -d' ' -f1 || echo "unknown")
            printf "  %-12s %s\n" "$tag" "$date"
        done <<< "$tags"
    else
        warning "Git repository not found or no tags available"
    fi
}

# Compare versions
compare_versions() {
    local version1=$1
    local version2=$2
    
    # Remove v prefix for comparison
    version1=${version1#v}
    version2=${version2#v}
    
    # Use sort to compare versions
    local sorted
    sorted=$(printf "%s\n%s\n" "$version1" "$version2" | sort -V)
    
    local first_line
    first_line=$(echo "$sorted" | head -n1)
    
    if [ "$first_line" = "$version1" ]; then
        if [ "$version1" = "$version2" ]; then
            echo "equal"
        else
            echo "less"
        fi
    else
        echo "greater"
    fi
}

# Show current version command
cmd_current() {
    local current
    current=$(get_current_version)
    echo "$current"
}

# Bump version command  
cmd_bump() {
    local increment_type=$1
    local dry_run=${2:-false}
    
    local current
    current=$(get_current_version)
    local new_version
    new_version=$(increment_version "$current" "$increment_type")
    
    log "Current version: $current"
    log "New version: $new_version"
    
    if [ "$dry_run" = "true" ]; then
        log "Dry run - no changes made"
        return 0
    fi
    
    set_version_file "$new_version"
}

# Set version command
cmd_set() {
    local version=$1
    local dry_run=${2:-false}
    local no_prefix=${3:-false}
    
    if ! validate_version "$version" "$no_prefix"; then
        error "Invalid version format: $version"
    fi
    
    # Ensure version has v prefix if not disabled
    if [ "$no_prefix" != "true" ] && [[ ! $version =~ ^v ]]; then
        version="v$version"
    fi
    
    local current
    current=$(get_current_version)
    
    log "Current version: $current"
    log "Setting version to: $version"
    
    if [ "$dry_run" = "true" ]; then
        log "Dry run - no changes made"
        return 0
    fi
    
    set_version_file "$version"
}

# Next version command (show only)
cmd_next() {
    local increment_type=$1
    
    local current
    current=$(get_current_version)
    local next
    next=$(increment_version "$current" "$increment_type")
    
    echo "$next"
}

# Validate version command
cmd_validate() {
    local version=$1
    local no_prefix=${2:-false}
    
    if validate_version "$version" "$no_prefix"; then
        success "Version $version is valid"
        return 0
    else
        error "Version $version is invalid"
        return 1
    fi
}

# Main function
main() {
    local dry_run=false
    local no_prefix=false
    local command=""
    local args=()
    
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --dry-run)
                dry_run=true
                shift
                ;;
            --no-prefix)
                no_prefix=true
                shift
                ;;
            --help|-h)
                usage
                exit 0
                ;;
            current|bump|set|next|validate|history)
                command=$1
                shift
                break
                ;;
            *)
                error "Unknown option: $1"
                ;;
        esac
    done
    
    # Collect remaining arguments
    args=("$@")
    
    if [ -z "$command" ]; then
        usage
        exit 1
    fi
    
    case $command in
        current)
            cmd_current
            ;;
        bump)
            if [ ${#args[@]} -eq 0 ]; then
                error "Bump command requires increment type: major, minor, or patch"
            fi
            cmd_bump "${args[0]}" "$dry_run"
            ;;
        set)
            if [ ${#args[@]} -eq 0 ]; then
                error "Set command requires version argument"
            fi
            cmd_set "${args[0]}" "$dry_run" "$no_prefix"
            ;;
        next)
            if [ ${#args[@]} -eq 0 ]; then
                error "Next command requires increment type: major, minor, or patch"
            fi
            cmd_next "${args[0]}"
            ;;
        validate)
            if [ ${#args[@]} -eq 0 ]; then
                error "Validate command requires version argument"
            fi
            cmd_validate "${args[0]}" "$no_prefix"
            ;;
        history)
            show_version_history
            ;;
        *)
            error "Unknown command: $command"
            ;;
    esac
}

# Execute if run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi