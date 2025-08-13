#!/bin/bash

# 📝 Advanced Changelog Generator for k8s-cli
# Generates detailed changelog entries with smart commit analysis

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
TEMPLATE_FILE="$SCRIPT_DIR/changelog-template.md"

log() {
    echo -e "${BLUE}[CHANGELOG]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Get git information
get_git_info() {
    local info_type=$1
    
    case $info_type in
        commit)
            git rev-parse --short HEAD
            ;;
        branch)
            git rev-parse --abbrev-ref HEAD
            ;;
        commit_count)
            local last_tag
            if git describe --tags --abbrev=0 > /dev/null 2>&1; then
                last_tag=$(git describe --tags --abbrev=0)
                git rev-list --count ${last_tag}..HEAD
            else
                git rev-list --count HEAD
            fi
            ;;
        contributors)
            local last_tag
            if git describe --tags --abbrev=0 > /dev/null 2>&1; then
                last_tag=$(git describe --tags --abbrev=0)
                git log --format='%an' ${last_tag}..HEAD | sort -u | wc -l | tr -d ' '
            else
                git log --format='%an' | sort -u | wc -l | tr -d ' '
            fi
            ;;
        files_changed)
            local last_tag
            if git describe --tags --abbrev=0 > /dev/null 2>&1; then
                last_tag=$(git describe --tags --abbrev=0)
                git diff --name-only ${last_tag}..HEAD | wc -l | tr -d ' '
            else
                echo "N/A"
            fi
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

# Analyze commits and categorize them
analyze_commits() {
    local category=$1
    local last_tag
    
    if git describe --tags --abbrev=0 > /dev/null 2>&1; then
        last_tag=$(git describe --tags --abbrev=0)
        local commits=$(git log --oneline ${last_tag}..HEAD --no-merges)
    else
        local commits=$(git log --oneline --no-merges -20)
    fi
    
    local output=""
    
    case $category in
        features)
            while IFS= read -r commit; do
                if [[ $commit == *"feat:"* ]] || [[ $commit == *"add:"* ]] || [[ $commit == *"new:"* ]]; then
                    local msg=$(echo "$commit" | sed 's/^[a-f0-9]* //' | sed 's/^feat: \|^add: \|^new: //')
                    output="${output}- **$msg**\n"
                fi
            done <<< "$commits"
            ;;
        enhancements)
            while IFS= read -r commit; do
                if [[ $commit == *"enhance:"* ]] || [[ $commit == *"improve:"* ]] || [[ $commit == *"update:"* ]] || [[ $commit == *"refactor:"* ]]; then
                    local msg=$(echo "$commit" | sed 's/^[a-f0-9]* //' | sed 's/^enhance: \|^improve: \|^update: \|^refactor: //')
                    output="${output}- **$msg**\n"
                fi
            done <<< "$commits"
            ;;
        fixes)
            while IFS= read -r commit; do
                if [[ $commit == *"fix:"* ]] || [[ $commit == *"bug:"* ]] || [[ $commit == *"hotfix:"* ]]; then
                    local msg=$(echo "$commit" | sed 's/^[a-f0-9]* //' | sed 's/^fix: \|^bug: \|^hotfix: //')
                    output="${output}- **$msg**\n"
                fi
            done <<< "$commits"
            ;;
        docs)
            while IFS= read -r commit; do
                if [[ $commit == *"docs:"* ]] || [[ $commit == *"doc:"* ]] || [[ $commit == *"documentation:"* ]]; then
                    local msg=$(echo "$commit" | sed 's/^[a-f0-9]* //' | sed 's/^docs: \|^doc: \|^documentation: //')
                    output="${output}- **$msg**\n"
                fi
            done <<< "$commits"
            ;;
        *)
            output="- No specific changes in this category"
            ;;
    esac
    
    if [ -z "$output" ]; then
        echo "- No changes in this category"
    else
        echo -e "$output"
    fi
}

# Get build information
get_build_info() {
    local info_type=$1
    
    case $info_type in
        go_version)
            go version | awk '{print $3}'
            ;;
        binary_size)
            if [ -f "$PROJECT_ROOT/bin/k8s-cli" ]; then
                local size=$(stat -f%z "$PROJECT_ROOT/bin/k8s-cli" 2>/dev/null || stat -c%s "$PROJECT_ROOT/bin/k8s-cli" 2>/dev/null || echo "0")
                local mb=$((size / 1024 / 1024))
                echo "${mb}MB"
            else
                echo "Not built"
            fi
            ;;
        test_status)
            cd "$PROJECT_ROOT"
            if make -f Makefile.dev test > /dev/null 2>&1; then
                echo "✅ Passing"
            else
                echo "❌ Some failures"
            fi
            ;;
        coverage)
            cd "$PROJECT_ROOT"
            if make -f Makefile.dev test-coverage > /dev/null 2>&1; then
                if [ -f "coverage.out" ]; then
                    go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//'
                else
                    echo "Unknown"
                fi
            else
                echo "Unknown"
            fi
            ;;
        build_time)
            local start_time=$(date +%s)
            cd "$PROJECT_ROOT"
            make -f Makefile.dev build > /dev/null 2>&1 || true
            local end_time=$(date +%s)
            echo "$((end_time - start_time))s"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

# Determine release type based on version change
determine_release_type() {
    local current_version=$1
    local new_version=$2
    
    # Remove 'v' prefix
    current_version=${current_version#v}
    new_version=${new_version#v}
    
    IFS='.' read -ra CURRENT <<< "$current_version"
    IFS='.' read -ra NEW <<< "$new_version"
    
    local current_major=${CURRENT[0]:-0}
    local current_minor=${CURRENT[1]:-0}
    local new_major=${NEW[0]:-0}
    local new_minor=${NEW[1]:-0}
    
    if [ "$new_major" -gt "$current_major" ]; then
        echo "Major"
    elif [ "$new_minor" -gt "$current_minor" ]; then
        echo "Minor"
    else
        echo "Patch"
    fi
}

# Generate smart descriptions based on commits
generate_smart_descriptions() {
    local category=$1
    local commit_count=$(get_git_info commit_count)
    
    case $category in
        release_description)
            if [ "$commit_count" -gt 10 ]; then
                echo "This release includes significant improvements and new functionality across multiple areas."
            elif [ "$commit_count" -gt 5 ]; then
                echo "This release includes several improvements and bug fixes."
            else
                echo "This release includes targeted improvements and fixes."
            fi
            ;;
        devops_benefits)
            echo "- Enhanced automation capabilities
- Improved monitoring and observability
- Better integration with CI/CD pipelines"
            ;;
        sre_benefits)
            echo "- More reliable component detection
- Better error handling and recovery
- Enhanced troubleshooting capabilities"
            ;;
        dx_benefits)
            echo "- Improved CLI usability
- Better documentation and examples
- Enhanced development workflows"
            ;;
        usage_examples)
            echo "k8s-cli --version
k8s-cli version  # See enhanced component detection"
            ;;
        upgrade_steps)
            echo "1. **Update CLI**: \`make -f Makefile.dev install-user\`
2. **Verify installation**: \`k8s-cli --version\`
3. **Test new features**: \`k8s-cli version\`"
            ;;
        breaking_changes)
            echo "- No breaking changes in this release"
            ;;
        known_issues)
            echo "- No known critical issues"
            ;;
        *)
            echo "Details will be provided in the full release notes."
            ;;
    esac
}

# Generate changelog entry
generate_changelog() {
    local version=$1
    local current_version=${2:-"v0.0.0"}
    
    log "Generating changelog for $version..."
    
    if [ ! -f "$TEMPLATE_FILE" ]; then
        error "Template file not found: $TEMPLATE_FILE"
    fi
    
    local date=$(date +%Y-%m-%d)
    local release_type=$(determine_release_type "$current_version" "$version")
    
    # Read template
    local changelog_content
    changelog_content=$(cat "$TEMPLATE_FILE")
    
    # Replace placeholders
    changelog_content=${changelog_content//\{\{VERSION\}\}/$version}
    changelog_content=${changelog_content//\{\{DATE\}\}/$date}
    changelog_content=${changelog_content//\{\{RELEASE_TYPE\}\}/$release_type}
    changelog_content=${changelog_content//\{\{GIT_COMMIT\}\}/$(get_git_info commit)}
    changelog_content=${changelog_content//\{\{BUILD_DATE\}\}/$date}
    changelog_content=${changelog_content//\{\{GIT_BRANCH\}\}/$(get_git_info branch)}
    changelog_content=${changelog_content//\{\{GO_VERSION\}\}/$(get_build_info go_version)}
    changelog_content=${changelog_content//\{\{COMMIT_COUNT\}\}/$(get_git_info commit_count)}
    changelog_content=${changelog_content//\{\{FILES_CHANGED\}\}/$(get_git_info files_changed)}
    changelog_content=${changelog_content//\{\{CONTRIBUTORS\}\}/$(get_git_info contributors)}
    changelog_content=${changelog_content//\{\{TEST_STATUS\}\}/$(get_build_info test_status)}
    changelog_content=${changelog_content//\{\{COVERAGE\}\}/$(get_build_info coverage)}
    changelog_content=${changelog_content//\{\{BUILD_TIME\}\}/$(get_build_info build_time)}
    changelog_content=${changelog_content//\{\{BINARY_SIZE\}\}/$(get_build_info binary_size)}
    
    # Replace commit-based content
    changelog_content=${changelog_content//\{\{RELEASE_DESCRIPTION\}\}/$(generate_smart_descriptions release_description)}
    changelog_content=${changelog_content//\{\{FEATURES_LIST\}\}/$(analyze_commits features)}
    changelog_content=${changelog_content//\{\{ENHANCEMENTS_LIST\}\}/$(analyze_commits enhancements)}
    changelog_content=${changelog_content//\{\{FIXES_LIST\}\}/$(analyze_commits fixes)}
    changelog_content=${changelog_content//\{\{DOCS_LIST\}\}/$(analyze_commits docs)}
    
    # Replace smart descriptions
    changelog_content=${changelog_content//\{\{DEVOPS_BENEFITS\}\}/$(generate_smart_descriptions devops_benefits)}
    changelog_content=${changelog_content//\{\{SRE_BENEFITS\}\}/$(generate_smart_descriptions sre_benefits)}
    changelog_content=${changelog_content//\{\{DX_BENEFITS\}\}/$(generate_smart_descriptions dx_benefits)}
    changelog_content=${changelog_content//\{\{USAGE_EXAMPLES\}\}/$(generate_smart_descriptions usage_examples)}
    changelog_content=${changelog_content//\{\{UPGRADE_STEPS\}\}/$(generate_smart_descriptions upgrade_steps)}
    changelog_content=${changelog_content//\{\{BREAKING_CHANGES\}\}/$(generate_smart_descriptions breaking_changes)}
    changelog_content=${changelog_content//\{\{KNOWN_ISSUES\}\}/$(generate_smart_descriptions known_issues)}
    
    echo "$changelog_content"
    
    success "Changelog generated successfully"
}

# Main function
main() {
    local version=${1:-"v0.0.1"}
    local current_version=${2:-""}
    
    if [ -z "$current_version" ] && [ -f "$PROJECT_ROOT/VERSION" ]; then
        current_version=$(cat "$PROJECT_ROOT/VERSION")
    fi
    
    generate_changelog "$version" "$current_version"
}

# Execute if run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi