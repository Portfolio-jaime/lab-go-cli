# 🛠️ Development Guide

## 🚀 Quick Start

### Prerequisites
- Go 1.24.5 or later
- Docker & Docker Compose
- VS Code with Dev Containers extension
- kubectl configured

### Development Setup
```bash
# 1. Clone and open in Dev Container
git clone <repository>
code lab-go-cli

# 2. Reopen in Container when prompted
# 3. Start Minikube in container
.devcontainer/scripts/start-minikube.sh

# 4. Build and test
make build
make test
./k8s-cli --help
```

## 🏗️ Development Workflow

### Daily Development Cycle
```bash
# 1. Pull latest changes
git pull origin main

# 2. Create feature branch
git checkout -b feature/new-analysis-type

# 3. Development cycle
make dev-watch    # Auto-rebuild on changes
make test-watch   # Auto-test on changes

# 4. Before commit
make check-all    # Lint, test, vet
make docs-update  # Update documentation

# 5. Commit and push
git add .
git commit -m "feat: add new analysis type"
git push origin feature/new-analysis-type
```

## 📁 Project Structure Deep Dive

### Command Structure (cmd/)
```
cmd/
├── root.go          # Base command, global flags
├── all.go           # Comprehensive analysis
├── metrics.go       # Real-time metrics
├── cost.go          # Cost analysis
├── workload.go      # Workload health
├── logs.go          # Events and logs
├── export.go        # Data export
├── recommend.go     # Recommendations
├── resources.go     # Basic resources
└── version.go       # Version info
```

Each command follows this pattern:
```go
// Command definition
var commandCmd = &cobra.Command{
    Use:   "command",
    Short: "Short description",
    Long:  "Long description",
    RunE:  runCommandFunction,
}

// Flags
var (
    flagOne   bool
    flagTwo   string
)

// Initialization
func init() {
    rootCmd.AddCommand(commandCmd)
    commandCmd.Flags().BoolVar(&flagOne, "flag-one", false, "Description")
    commandCmd.Flags().StringVar(&flagTwo, "flag-two", "", "Description")
}

// Main function
func runCommandFunction(cmd *cobra.Command, args []string) error {
    // Implementation
}
```

### Business Logic Structure (pkg/)
```
pkg/
├── kubernetes/           # Core Kubernetes logic
│   ├── client.go        # K8s client management
│   ├── metrics.go       # Metrics collection
│   ├── cost_analysis.go # Cost calculations
│   ├── workload_*.go    # Workload analysis
│   └── events_logs.go   # Event processing
├── export/              # Export functionality
│   └── exporter.go      # Multi-format export
├── recommendations/     # Recommendation engine
│   └── analyzer.go      # Analysis logic
└── table/              # Output formatting
    └── table.go        # Table formatting
```

## 🔧 Development Tools

### Make Targets
```bash
# Development
make build          # Build binary
make dev-build      # Build with debug info
make dev-watch      # Auto-rebuild on changes
make clean          # Clean build artifacts

# Testing
make test           # Run unit tests
make test-watch     # Auto-test on changes
make test-coverage  # Generate coverage report
make test-integration # Run integration tests
make test-e2e       # End-to-end tests

# Quality
make fmt            # Format code
make lint           # Run linter
make vet            # Run go vet
make check-all      # All quality checks

# Documentation
make docs-generate  # Generate docs
make docs-update    # Update all docs
make docs-serve     # Serve docs locally

# Deployment
make release        # Build release binaries
make docker-build   # Build Docker image
make docker-push    # Push to registry
```

### Advanced Make Targets
```bash
# Development helpers
make deps-update    # Update dependencies
make deps-tidy      # Tidy go.mod
make generate       # Run go generate
make mock-generate  # Generate mocks

# Performance
make bench          # Run benchmarks
make profile        # CPU profiling
make memory-profile # Memory profiling

# Security
make security-scan  # Security vulnerability scan
make deps-audit     # Dependency audit

# CI/CD
make ci-test        # CI test suite
make ci-build       # CI build
make ci-deploy      # CI deployment
```

## 🧪 Testing Strategy

### Testing Pyramid
```
    E2E Tests (Few)
       ↑
  Integration Tests (Some)
       ↑
   Unit Tests (Many)
```

### Unit Tests
```go
// Example: pkg/kubernetes/metrics_test.go
package kubernetes

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestGetRealTimeNodeMetrics(t *testing.T) {
    // Arrange
    client := &Client{/* mock setup */}
    
    // Act
    metrics, err := client.GetRealTimeNodeMetrics()
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, metrics)
    assert.Greater(t, metrics[0].CPUUsagePercent, 0.0)
}
```

### Integration Tests
```go
// Example: test/integration/metrics_integration_test.go
//go:build integration

func TestMetricsIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    client, err := kubernetes.NewClient("")
    require.NoError(t, err)
    
    metrics, err := client.GetRealTimeNodeMetrics()
    require.NoError(t, err)
    assert.NotEmpty(t, metrics)
}
```

### E2E Tests
```bash
# test/e2e/test_k8s_cli.sh
#!/bin/bash
set -e

echo "Testing k8s-cli end-to-end..."

# Build
make build

# Test basic commands
./k8s-cli --help
./k8s-cli version
./k8s-cli all --dry-run

# Test with real cluster
if kubectl cluster-info >/dev/null 2>&1; then
    ./k8s-cli metrics --nodes
    ./k8s-cli cost --nodes
    ./k8s-cli export --format json --output /tmp
fi

echo "E2E tests passed!"
```

## 📝 Code Style Guide

### Go Style
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use [golangci-lint](https://golangci-lint.run/) configuration
- Maximum line length: 100 characters
- Use meaningful variable names

### Error Handling
```go
// Good: Wrap errors with context
func (c *Client) GetMetrics() error {
    data, err := c.fetchData()
    if err != nil {
        return fmt.Errorf("failed to fetch metrics data: %w", err)
    }
    return nil
}

// Bad: Ignore or log without returning
func (c *Client) GetMetrics() error {
    data, err := c.fetchData()
    if err != nil {
        log.Println("Error:", err) // Don't do this
        return nil
    }
    return nil
}
```

### Logging
```go
// Use structured logging
import "log/slog"

func (c *Client) ProcessMetrics() {
    slog.Info("Starting metrics processing",
        "namespace", c.namespace,
        "node_count", len(c.nodes))
}
```

### Comments
```go
// Package comment describes the package purpose
package kubernetes

// Function comments describe what, not how
// GetRealTimeNodeMetrics retrieves current CPU and memory usage
// for all nodes in the cluster using the metrics-server API.
func (c *Client) GetRealTimeNodeMetrics() ([]NodeMetrics, error) {
    // Implementation
}
```

## 🔄 Adding New Features

### 1. Adding a New Command
```bash
# 1. Create command file
touch cmd/newcommand.go

# 2. Implement command structure
cat > cmd/newcommand.go << 'EOF'
package cmd

import (
    "github.com/spf13/cobra"
    "k8s-cli/pkg/kubernetes"
)

var newCommandCmd = &cobra.Command{
    Use:   "newcommand",
    Short: "Description of new command",
    RunE:  runNewCommand,
}

func init() {
    rootCmd.AddCommand(newCommandCmd)
}

func runNewCommand(cmd *cobra.Command, args []string) error {
    // Implementation
    return nil
}
EOF

# 3. Add business logic in pkg/
# 4. Add tests
# 5. Update documentation
```

### 2. Adding New Analysis Type
```bash
# 1. Create analysis file
touch pkg/kubernetes/security_analysis.go

# 2. Define data structures
# 3. Implement analysis logic
# 4. Add to export system
# 5. Create command interface
# 6. Add tests and documentation
```

### 3. Adding New Export Format
```bash
# 1. Extend exporter
# 2. Add format-specific logic
# 3. Update command flags
# 4. Add tests
# 5. Update documentation
```

## 🐛 Debugging

### Debug Build
```bash
# Build with debug symbols
make dev-build

# Run with verbose logging
./k8s-cli --verbose metrics

# Enable debug output
DEBUG=true ./k8s-cli metrics
```

### Profiling
```bash
# CPU profiling
go tool pprof cpu.prof

# Memory profiling
go tool pprof mem.prof

# Run with profiling
./k8s-cli -cpuprofile=cpu.prof -memprofile=mem.prof metrics
```

### Common Issues

#### 1. Metrics Server Not Available
```bash
# Check if metrics-server is running
kubectl get pods -n kube-system | grep metrics-server

# Install if missing
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

#### 2. RBAC Permissions
```bash
# Check current permissions
kubectl auth can-i list pods --as=system:serviceaccount:default:default

# Create proper service account
kubectl create serviceaccount k8s-cli
kubectl create clusterrolebinding k8s-cli --clusterrole=view --serviceaccount=default:k8s-cli
```

#### 3. Build Issues
```bash
# Clean and rebuild
make clean
go mod tidy
make build

# Check dependencies
go mod verify
```

## 📚 Documentation

### Auto-generated Docs
```bash
# Generate command documentation
make docs-generate

# Update API documentation
make docs-api

# Generate usage examples
make docs-examples
```

### Manual Documentation
- Keep README.md updated
- Document breaking changes in CHANGELOG.md
- Update architecture docs for major changes
- Add examples for new features

## 🚀 Release Process

### Versioning
- Follow [Semantic Versioning](https://semver.org/)
- Tag releases: `v1.2.3`
- Update version in code before release

### Release Checklist
```bash
# 1. Update version
echo "v1.2.3" > VERSION

# 2. Update changelog
# Edit CHANGELOG.md

# 3. Run full test suite
make check-all
make test-e2e

# 4. Build release
make release

# 5. Tag and push
git tag v1.2.3
git push origin v1.2.3

# 6. Create GitHub release
gh release create v1.2.3 --generate-notes
```

## 🤝 Contributing

### Pull Request Process
1. Fork repository
2. Create feature branch
3. Make changes with tests
4. Run quality checks
5. Submit PR with description
6. Address review feedback
7. Merge after approval

### Code Review Guidelines
- Check for proper error handling
- Verify test coverage
- Review documentation updates
- Ensure backwards compatibility
- Test with real clusters when possible