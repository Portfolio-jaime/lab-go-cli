# 🚀 k8s-cli - Enterprise Kubernetes Analysis Platform

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![Version](https://img.shields.io/badge/Version-2.0.6-green?style=for-the-badge)](VERSION)
[![CI/CD](https://img.shields.io/badge/CI%2FCD-Automated-brightgreen?style=for-the-badge)](docs/ops/ci-cd.md)
[![Documentation](https://img.shields.io/badge/Docs-Complete-success?style=for-the-badge)](#-documentation)

> **Enterprise-grade Kubernetes cluster analysis, cost optimization, and monitoring platform**

---

## 🎯 Overview

k8s-cli is a comprehensive platform that transforms raw Kubernetes cluster data into actionable insights for DevOps, FinOps, and SRE teams. Get real-time metrics, cost optimization recommendations, and proactive health monitoring in a single tool.

### ✨ Key Features

- **💰 Cost Optimization** - Identify underutilized resources and potential savings
- **📊 Real-time Metrics** - CPU/Memory utilization with performance insights  
- **🔍 Health Monitoring** - Workload health scoring and issue detection
- **📤 Multi-format Export** - JSON, CSV, and Prometheus integration
- **🎯 Enterprise Ready** - Automated CI/CD, security scanning, and comprehensive documentation

---

## 🚀 Quick Start

### 📦 Installation

```bash
# Download latest release
curl -L https://github.com/your-org/k8s-cli/releases/latest/download/k8s-cli-linux-amd64.tar.gz | tar xz

# Move to PATH
sudo mv k8s-cli /usr/local/bin/

# Verify installation
k8s-cli --version
```

### ⚡ Basic Usage

```bash
# Complete cluster analysis
k8s-cli all

# Real-time metrics
k8s-cli metrics --nodes --pods

# Cost analysis
k8s-cli cost --underutilized

# Export data
k8s-cli export --format json --costs --metrics
```

---

## 📊 Core Commands

| Command | Description | Example |
|---------|-------------|---------|
| `all` | Complete cluster analysis | `k8s-cli all` |
| `metrics` | Real-time metrics and utilization | `k8s-cli metrics --nodes --pods --utilization` |
| `cost` | Cost analysis and optimization | `k8s-cli cost --underutilized` |
| `workload` | Workload health analysis | `k8s-cli workload --unhealthy-only` |
| `logs` | Events and log analysis | `k8s-cli logs --critical --hours 24` |
| `export` | Multi-format data export | `k8s-cli export --format csv --output ./reports/` |
| `version` | Cluster version information | `k8s-cli version` |

---

## 💼 Enterprise Use Cases

### 🏦 FinOps (Financial Operations)
```bash
# Daily cost optimization
k8s-cli cost --underutilized > daily-savings-$(date +%Y%m%d).txt

# Weekly finance report
k8s-cli export --format csv --costs --output ./finance-reports/
```

### 🔧 DevOps Monitoring
```bash
# Real-time cluster dashboard
k8s-cli metrics --nodes --pods --utilization

# Prometheus integration
k8s-cli export --format prometheus --output /var/lib/prometheus/
```

### 🚨 SRE (Site Reliability Engineering)
```bash
# Incident response
k8s-cli logs --critical --patterns --hours 2

# Post-incident analysis  
k8s-cli export --format json --logs --events --hours 24
```

---

## 🛠️ Development

### 📋 Prerequisites

- **Go 1.24+** (required for k8s.io dependencies)
- **Kubernetes cluster** (local or remote)
- **kubectl** configured

### 🔧 Setup

```bash
# Clone repository
git clone https://github.com/your-org/k8s-cli.git
cd k8s-cli

# Setup development environment  
make -f Makefile.dev dev-setup

# Build and test
make -f Makefile.dev build
make -f Makefile.dev test

# Start development with auto-rebuild
make -f Makefile.dev watch
```

### 🧪 Quality Assurance

```bash
# Run all quality checks
make -f Makefile.dev pre-commit

# This includes:
# ✅ Code formatting (gofmt, goimports)
# ✅ Linting (golangci-lint)
# ✅ Security scanning (govulncheck)  
# ✅ Unit testing with coverage
# ✅ Build verification
```

---

## 📚 Documentation

### 📖 User Documentation
- **[Installation Guide](docs/user/installation.md)** - Complete installation methods
- **[Quick Start Guide](docs/user/quick-start.md)** - Get started in 5 minutes
- **[Command Reference](docs/user/commands.md)** - Complete command documentation
- **[Examples](docs/user/examples.md)** - Practical usage examples

### 👨‍💻 Developer Documentation  
- **[Architecture Guide](docs/developer/architecture.md)** - System design and components
- **[Development Guide](docs/developer/development.md)** - Development setup and workflow
- **[API Reference](docs/developer/api.md)** - Internal API documentation
- **[Testing Guide](docs/developer/testing.md)** - Testing strategies and utilities

### 🔧 Operations Documentation
- **[CI/CD Guide](docs/ops/ci-cd.md)** - Automated workflows and releases
- **[Release Process](docs/ops/release-process.md)** - Release management  
- **[Deployment Guide](docs/ops/deployment.md)** - Production deployment
- **[Monitoring](docs/ops/monitoring.md)** - Observability and monitoring

### 🆘 Support Documentation
- **[Troubleshooting Guide](docs/reference/troubleshooting.md)** - Common issues and solutions
- **[Configuration Reference](docs/reference/configuration.md)** - All configuration options
- **[Error Codes](docs/reference/error-codes.md)** - Error reference guide

---

## 🔧 Configuration

### 🌍 Environment Variables

```bash
export K8S_CLI_NAMESPACE=production      # Default namespace
export K8S_CLI_OUTPUT_FORMAT=json        # Default output format  
export K8S_CLI_KUBECONFIG=/path/to/config # Custom kubeconfig
```

### 📄 Configuration File

```yaml
# ~/.k8s-cli.yaml
output:
  format: table
  colors: true
  
metrics:
  cache_duration: 5m
  include_system_pods: false
  
cost:
  currency: USD
  default_node_cost: 72.0
```

---

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### 🔄 Development Workflow

1. **Fork and clone** the repository
2. **Create feature branch** (`git checkout -b feature/amazing-feature`)
3. **Make changes** and add tests
4. **Run quality checks** (`make -f Makefile.dev pre-commit`)
5. **Commit changes** (`git commit -m 'feat: add amazing feature'`)
6. **Push to branch** (`git push origin feature/amazing-feature`)
7. **Open Pull Request** with comprehensive description

### 📝 Commit Convention

We use [Conventional Commits](https://conventionalcommits.org/):

```bash
feat: add new feature      # Triggers minor release
fix: resolve bug          # Triggers patch release  
docs: update documentation # No release
chore: update dependencies # No release
```

---

## 📈 Performance & Scalability

### ⚡ Performance Characteristics
- **Concurrent operations** - Parallel data fetching and analysis
- **Memory efficient** - Optimized for large clusters (1000+ nodes)
- **Fast execution** - Sub-second response for basic operations
- **Smart caching** - Configurable cache for API responses

### 📊 Scalability
- **Cluster size**: Tested with 1000+ nodes, 5000+ pods
- **Concurrent users**: Designed for team environments
- **Multi-cluster**: Architecture ready for federation

---

## 🔒 Security

### 🛡️ Security Features
- **No credential storage** - Uses existing kubeconfig
- **Read-only access** - Minimal required permissions
- **Secure exports** - Configurable data retention
- **Audit logging** - Optional activity tracking

### 🔐 Required Permissions

```yaml
# Minimum RBAC permissions
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: k8s-cli-reader
rules:
- apiGroups: [""]
  resources: ["nodes", "pods", "services", "events"]
  verbs: ["get", "list"]
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets", "daemonsets"]
  verbs: ["get", "list"]
- apiGroups: ["metrics.k8s.io"]
  resources: ["nodes", "pods"]
  verbs: ["get", "list"]
```

---

## 📊 Roadmap

### 🎯 Current Focus (v2.0.x)
- ✅ **Complete CI/CD automation** 
- ✅ **Comprehensive documentation**
- ✅ **Cross-platform reliability**
- ✅ **Enterprise-grade quality**

### 🚀 Next Release (v2.1.0)
- 🔄 **Multi-cluster support** - Federation and comparison
- 🤖 **Machine learning** - Predictive analytics
- 🌐 **Web dashboard** - Visual cluster analysis
- 🔌 **Plugin system** - Extensible architecture

### 🌟 Future Vision (v3.0.0)
- **Cloud integration** - Native cloud provider cost APIs
- **Real-time streaming** - Live cluster monitoring
- **Advanced analytics** - Trend analysis and forecasting
- **Enterprise features** - SSO, RBAC, multi-tenancy

---

## 🆘 Support

### 📞 Getting Help
- **📖 Documentation** - Comprehensive guides and references
- **🐛 Issues** - [GitHub Issues](https://github.com/your-org/k8s-cli/issues) for bugs
- **💬 Discussions** - [GitHub Discussions](https://github.com/your-org/k8s-cli/discussions) for questions
- **🔧 Support** - [Troubleshooting Guide](docs/reference/troubleshooting.md)

### 🏷️ Issue Templates
We provide templates for:
- 🐛 **Bug reports** - Structured bug reporting
- ✨ **Feature requests** - New feature proposals  
- 📚 **Documentation** - Documentation improvements
- 🆘 **Support** - General help requests

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🎉 Acknowledgments

- **Kubernetes Community** - For the amazing ecosystem
- **Contributors** - For making this project better
- **Users** - For feedback and real-world testing
- **Go Team** - For the excellent programming language

---

## 📈 Project Stats

![GitHub stars](https://img.shields.io/github/stars/your-org/k8s-cli?style=social)
![GitHub forks](https://img.shields.io/github/forks/your-org/k8s-cli?style=social)
![GitHub issues](https://img.shields.io/github/issues/your-org/k8s-cli)
![GitHub pull requests](https://img.shields.io/github/issues-pr/your-org/k8s-cli)

---

**🚀 Ready to optimize your Kubernetes clusters?**

[**Download k8s-cli**](https://github.com/your-org/k8s-cli/releases/latest) • [**Read the Docs**](#-documentation) • [**Join the Community**](https://github.com/your-org/k8s-cli/discussions)

---

<div align="center">

**Built with ❤️ by the k8s-cli team**

[Website](https://your-org.github.io/k8s-cli) • [Documentation](docs/) • [Releases](https://github.com/your-org/k8s-cli/releases) • [Contributing](CONTRIBUTING.md)

</div>