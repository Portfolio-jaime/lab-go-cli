# 📚 k8s-cli Documentation Index

## 🎯 Overview

Comprehensive documentation for k8s-cli - Enterprise Kubernetes Analysis Platform

---

## 📖 Core Documentation

### 🚀 Getting Started
- **[README.md](README.md)** - Main project overview and quick start
- **[CHANGELOG.md](CHANGELOG.md)** - Version history and release notes
- **[Development Setup Script](scripts/dev-setup.sh)** - Automated development environment setup

### 🏗️ Architecture & Design
- **[Architecture Guide](docs/ARCHITECTURE.md)** - System design and component overview
- **[API Documentation](docs/API.md)** - Internal API reference and data structures
- **[Development Guide](docs/DEVELOPMENT.md)** - Contributing and development workflows

### 📚 Usage & Examples
- **[Usage Examples](docs/EXAMPLES.md)** - Comprehensive command examples and use cases
- **[Demo Script](examples/demo_new_features.sh)** - Interactive feature demonstration

---

## 🛠️ Development Resources

### 🔧 Build & Development
- **[Makefile](Makefile)** - Basic build targets
- **[Advanced Makefile](Makefile.dev)** - Comprehensive development workflows
- **[Air Config](.air.toml)** - Hot reload configuration

### 🎯 IDE Configuration
- **[VS Code Settings](.vscode/settings.json)** - Optimized Go development settings
- **[VS Code Tasks](.vscode/tasks.json)** - Build, test, and run tasks
- **[VS Code Debug](.vscode/launch.json)** - Debug configurations for all commands

### 🐳 Container Development
- **[DevContainer Config](.devcontainer/devcontainer.json)** - VS Code Dev Container setup
- **[DevContainer Documentation](docs/DEVCONTAINER.md)** - Complete container development guide

---

## 📊 Command Documentation

### Core Commands
| Command | Description | Documentation |
|---------|-------------|---------------|
| `all` | Complete cluster analysis | [Examples](docs/EXAMPLES.md#basic-cluster-analysis) |
| `version` | Cluster version info | [Basic Usage](README.md#quick-start) |
| `resources` | Resource overview | [Basic Usage](README.md#quick-start) |
| `recommend` | Optimization recommendations | [Basic Usage](README.md#quick-start) |

### Advanced Commands (New in v2.0)
| Command | Description | Documentation |
|---------|-------------|---------------|
| `metrics` | Real-time metrics & utilization | [Examples](docs/EXAMPLES.md#real-time-metrics-examples) |
| `cost` | Cost analysis & optimization | [Examples](docs/EXAMPLES.md#cost-analysis-examples) |
| `workload` | Workload health analysis | [Examples](docs/EXAMPLES.md#workload-health-examples) |
| `logs` | Events & log analysis | [Examples](docs/EXAMPLES.md#logs-and-events-examples) |
| `export` | Multi-format data export | [Examples](docs/EXAMPLES.md#export-examples) |

---

## 🎯 Use Case Documentation

### 💼 Enterprise Use Cases
- **[FinOps](docs/EXAMPLES.md#finops-cost-optimization)** - Financial operations and cost optimization
- **[DevOps](docs/EXAMPLES.md#devops-monitoring-pipeline)** - Development operations and monitoring
- **[SRE](docs/EXAMPLES.md#sre-incident-response)** - Site reliability engineering
- **[CI/CD](docs/EXAMPLES.md#cicd-integration)** - Continuous integration workflows

### 🔧 Technical Workflows
- **[Troubleshooting](docs/EXAMPLES.md#troubleshooting-common-issues)** - Common issues and solutions
- **[Performance](docs/EXAMPLES.md#performance-examples)** - Optimization and benchmarking
- **[Configuration](docs/EXAMPLES.md#configuration-examples)** - Custom configuration options

---

## 🚀 Quick Reference

### Development Commands
```bash
# Setup development environment
./scripts/dev-setup.sh

# Start auto-rebuild watcher
make watch

# Run complete development cycle
make dev-cycle

# Quality checks
make check-all
```

### Build Commands
```bash
# Basic build
make build

# Development build with debug
make dev-build

# Release build (multi-platform)
make release-build
```

### Testing Commands
```bash
# Unit tests
make test

# Integration tests
make test-integration

# E2E tests
make test-e2e

# Test coverage
make test-coverage
```

### CLI Usage
```bash
# Complete analysis
./bin/k8s-cli all

# Real-time metrics
./bin/k8s-cli metrics --nodes --pods --utilization

# Cost analysis
./bin/k8s-cli cost --underutilized

# Export data
./bin/k8s-cli export --format json --costs --metrics
```

---

## 📁 File Structure Reference

```
lab-go-cli/
├── 📄 README.md                      # Main documentation
├── 📄 CHANGELOG.md                   # Version history
├── 📄 DOCUMENTATION_INDEX.md         # This file
├── 📄 VERSION                        # Current version
├── 📄 Makefile                       # Basic build targets
├── 📄 Makefile.dev                   # Advanced development
├── 📄 .air.toml                      # Hot reload config
│
├── 📁 docs/                          # Documentation
│   ├── 📄 ARCHITECTURE.md            # System architecture
│   ├── 📄 API.md                     # API reference
│   ├── 📄 DEVELOPMENT.md             # Development guide
│   └── 📄 EXAMPLES.md                # Usage examples
│
├── 📁 .vscode/                       # VS Code configuration
│   ├── 📄 settings.json              # Editor settings
│   ├── 📄 tasks.json                 # Build tasks
│   └── 📄 launch.json                # Debug configs
│
├── 📁 cmd/                           # CLI commands
│   ├── 📄 all.go                     # Complete analysis
│   ├── 📄 metrics.go                 # Real-time metrics
│   ├── 📄 cost.go                    # Cost analysis
│   ├── 📄 workload.go                # Workload health
│   ├── 📄 logs.go                    # Events & logs
│   └── 📄 export.go                  # Data export
│
├── 📁 pkg/                           # Business logic
│   ├── 📁 kubernetes/                # K8s integration
│   ├── 📁 export/                    # Export functionality
│   ├── 📁 recommendations/           # Recommendation engine
│   └── 📁 table/                     # Output formatting
│
├── 📁 scripts/                       # Development scripts
│   └── 📄 dev-setup.sh               # Environment setup
│
└── 📁 examples/                      # Examples & demos
    ├── 📄 demo_new_features.sh       # Interactive demo
    └── 📄 sample_usage.sh            # Basic examples
```

---

## 🔄 Documentation Maintenance

### Auto-generated Content
- **Command help** - Generated from CLI help text
- **API reference** - Generated from Go code comments
- **Changelog** - Updated with each release

### Manual Content
- **Architecture documentation** - Updated for major changes
- **Usage examples** - Updated with new features
- **Development guides** - Updated with workflow changes

### Update Commands
```bash
# Update all documentation
make docs-update

# Generate API docs
make docs-generate

# Validate documentation
make docs-validate

# Serve docs locally
make docs-serve
```

---

## 🤝 Contributing to Documentation

### Documentation Standards
- Use clear, concise language
- Include practical examples
- Keep content up-to-date
- Follow established formatting

### Adding Documentation
1. Update relevant `.md` files
2. Add examples to `docs/EXAMPLES.md`
3. Update this index if needed
4. Run `make docs-update`
5. Test examples work correctly

### Documentation Review
- Check for broken links
- Verify examples work
- Ensure consistency
- Test with real clusters

---

## 📈 Documentation Metrics

### Coverage
- ✅ **Architecture** - Complete system design documentation
- ✅ **API Reference** - All public interfaces documented
- ✅ **Usage Examples** - Comprehensive command examples
- ✅ **Development** - Complete development workflows
- ✅ **IDE Integration** - VS Code optimized setup

### Quality
- ✅ **Practical Examples** - Real-world use cases
- ✅ **Troubleshooting** - Common issues and solutions
- ✅ **Performance** - Optimization guidelines
- ✅ **Enterprise** - Business use case documentation

---

## 🔗 External References

### Related Projects
- [Kubernetes](https://kubernetes.io/docs/) - Official Kubernetes documentation
- [Cobra](https://cobra.dev/) - CLI framework documentation
- [Go](https://golang.org/doc/) - Go programming language documentation

### Tools & Dependencies
- [golangci-lint](https://golangci-lint.run/) - Go linter
- [Air](https://github.com/air-verse/air) - Hot reload for Go
- [VS Code Go](https://code.visualstudio.com/docs/languages/go) - Go in VS Code

---

**Last Updated:** January 15, 2024  
**Version:** 2.0.0  
**Maintainer:** k8s-cli Development Team