# k8s-cli - Enterprise Kubernetes Analysis Platform

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![Version](https://img.shields.io/badge/Version-2.0.6-green?style=for-the-badge)](VERSION)

> **Enterprise-grade Kubernetes cluster analysis, cost optimization, and monitoring platform**

## 🎯 Overview

k8s-cli is a comprehensive platform that transforms raw Kubernetes cluster data into actionable insights for DevOps, FinOps, and SRE teams. Get real-time metrics, cost optimization recommendations, and proactive health monitoring in a single tool.

## ✨ Key Features

- **💰 Cost Optimization** - Identify underutilized resources and potential savings
- **📊 Real-time Metrics** - CPU/Memory utilization with performance insights  
- **🔍 Health Monitoring** - Workload health scoring and issue detection
- **📤 Multi-format Export** - JSON, CSV, and Prometheus integration
- **🎯 Enterprise Ready** - Automated CI/CD, security scanning, and comprehensive documentation

## 🚀 Quick Start

### Installation

```bash
# Download latest release
curl -L https://github.com/Portfolio-jaime/lab-go-cli/releases/latest/download/k8s-cli-linux-amd64 -o k8s-cli
chmod +x k8s-cli
sudo mv k8s-cli /usr/local/bin/
```

### Basic Usage

```bash
# View cluster resources overview
k8s-cli resources

# Get cost analysis
k8s-cli cost --namespace production

# Monitor workload health
k8s-cli workload --health-check

# Export metrics
k8s-cli export --format json --output cluster-report.json
```

## 📊 Core Commands

| Command | Description | Example |
|---------|-------------|---------|
| `resources` | Display cluster resources overview | `k8s-cli resources --namespace kube-system` |
| `cost` | Analyze resource costs and optimization opportunities | `k8s-cli cost --recommendations` |
| `workload` | Monitor workload health and performance | `k8s-cli workload --health-score` |
| `metrics` | Gather detailed cluster metrics | `k8s-cli metrics --cpu --memory` |
| `logs` | Retrieve and analyze logs | `k8s-cli logs --errors --last 1h` |
| `export` | Export data in multiple formats | `k8s-cli export --format prometheus` |

## 🏗️ Architecture

k8s-cli is built with a modular architecture designed for enterprise environments:

```
k8s-cli/
├── cmd/                  # CLI commands and interfaces
├── pkg/
│   ├── kubernetes/      # Kubernetes client and operations
│   ├── export/         # Data export functionality
│   ├── recommendations/ # Cost and performance analysis
│   └── table/          # Data visualization
├── docs/               # Comprehensive documentation
└── scripts/           # Automation and release scripts
```

## 🎯 Use Cases

### FinOps Teams
- **Cost Monitoring**: Track resource spending across namespaces
- **Optimization**: Identify oversized or underutilized resources
- **Budgeting**: Export cost data for financial planning

### DevOps Engineers
- **Resource Management**: Monitor cluster capacity and usage
- **Health Checks**: Proactive workload monitoring
- **Troubleshooting**: Quick access to logs and metrics

### SRE Teams
- **Performance Analysis**: Deep dive into resource utilization
- **Capacity Planning**: Data-driven scaling decisions
- **Incident Response**: Fast diagnosis with comprehensive metrics

## 📈 Enterprise Features

### CI/CD Integration
- **Automated Builds**: Multi-platform releases
- **Security Scanning**: Integrated vulnerability checks
- **Quality Gates**: Automated testing and validation

### Monitoring & Observability
- **Prometheus Integration**: Native metrics export
- **Health Scoring**: Workload health assessment
- **Alerting**: Integration with monitoring systems

### Security & Compliance
- **RBAC Support**: Kubernetes role-based access
- **Audit Logging**: Complete operation tracking
- **Secure Defaults**: Security-first configuration

## 🔗 Documentation

| Section | Description | Link |
|---------|-------------|------|
| **Getting Started** | Installation and basic usage | [Quick Start](quickstart.md) |
| **Architecture** | System design and components | [Architecture](ARCHITECTURE.md) |
| **Development** | Contributing and development setup | [Development](DEVELOPMENT.md) |
| **API Reference** | Complete API documentation | [API Reference](API.md) |
| **CI/CD** | Build and deployment processes | [CI/CD Guide](CI_CD_CONSOLIDATED.md) |

## 💡 Examples

### Cost Analysis
```bash
# Analyze costs for specific namespace
k8s-cli cost --namespace production --recommendations

# Export cost data for FinOps team
k8s-cli cost --export csv --output monthly-costs.csv
```

### Health Monitoring
```bash
# Check workload health across cluster
k8s-cli workload --health-check --all-namespaces

# Get detailed health score for specific deployment
k8s-cli workload --name my-app --health-score --detailed
```

### Performance Analysis
```bash
# Gather comprehensive metrics
k8s-cli metrics --cpu --memory --network --storage

# Focus on high-utilization resources
k8s-cli metrics --threshold 80 --recommendations
```

## 🤝 Contributing

We welcome contributions! Please see our [Development Guide](DEVELOPMENT.md) for details on:

- Setting up the development environment
- Running tests and quality checks
- Submitting pull requests
- Release process

## 📞 Support

**Maintainer:** Jaime Henao <jaime.andres.henao.arbelaez@ba.com>  
**Organization:** British Airways DevOps Team  
**Repository:** [GitHub](https://github.com/Portfolio-jaime/lab-go-cli)

For enterprise support, please contact the British Airways DevOps team.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Built with ❤️ for Enterprise Kubernetes Environments**