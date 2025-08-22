# Quick Start Guide

Get up and running with k8s-cli in 5 minutes.

## 🏃‍♂️ 5-Minute Quick Start

### 1. Install k8s-cli
```bash
# Download and install (Linux/macOS)
curl -L https://github.com/Portfolio-jaime/lab-go-cli/releases/latest/download/k8s-cli-linux-amd64 -o k8s-cli
chmod +x k8s-cli
sudo mv k8s-cli /usr/local/bin/

# Verify installation
k8s-cli --version
```

### 2. Basic Commands
```bash
# Get cluster overview
k8s-cli resources

# Check cluster metrics
k8s-cli metrics

# Analyze costs
k8s-cli cost

# Complete analysis
k8s-cli all
```

### 3. Export Results
```bash
# Export to JSON
k8s-cli export --format json --output cluster-report.json

# Export to CSV for Excel
k8s-cli export --format csv --output metrics.csv
```

## 📊 Common Use Cases

### FinOps Team: Cost Analysis
```bash
# Daily cost check
k8s-cli cost --underutilized

# Weekly cost report
k8s-cli cost --detailed --export csv --output weekly-costs-$(date +%Y%m%d).csv

# Identify savings opportunities
k8s-cli cost --recommendations > cost-recommendations.txt
```

### DevOps: Cluster Health
```bash
# Quick health check
k8s-cli workload --unhealthy-only

# Resource utilization
k8s-cli metrics --nodes --pods --utilization

# Export for monitoring dashboard
k8s-cli metrics --export prometheus --output /var/lib/prometheus/
```

### SRE: Incident Response
```bash
# Find issues quickly
k8s-cli logs --critical --hours 2

# Get complete cluster state
k8s-cli all --namespace production --export json --output incident-$(date +%Y%m%d-%H%M).json

# Focus on specific workloads
k8s-cli workload --name frontend --detailed
```

## 🎯 Essential Commands

### Resource Overview
```bash
# All namespaces overview
k8s-cli resources

# Specific namespace
k8s-cli resources --namespace production

# Filter by resource type
k8s-cli resources --type deployment,service

# Limit output
k8s-cli resources --limit 10
```

### Real-time Metrics
```bash
# Node metrics
k8s-cli metrics --nodes

# Pod metrics with utilization
k8s-cli metrics --pods --utilization

# Specific namespace metrics
k8s-cli metrics --namespace kube-system

# Historical metrics (if available)
k8s-cli metrics --hours 24
```

### Cost Analysis
```bash
# Basic cost analysis
k8s-cli cost

# Show underutilized resources
k8s-cli cost --underutilized

# Cost by namespace
k8s-cli cost --by-namespace

# Include recommendations
k8s-cli cost --recommendations
```

### Workload Health
```bash
# All workloads health
k8s-cli workload

# Only unhealthy workloads
k8s-cli workload --unhealthy-only

# Specific workload analysis
k8s-cli workload --name my-app --namespace production

# Include health score
k8s-cli workload --health-score
```

### Event and Log Analysis
```bash
# Recent critical events
k8s-cli logs --critical

# Events from last hour
k8s-cli logs --hours 1

# Filter by patterns
k8s-cli logs --patterns --keywords "error,failed,timeout"

# Specific namespace events
k8s-cli logs --namespace production --hours 24
```

## 📤 Data Export

### Export Formats
```bash
# JSON (programmatic use)
k8s-cli export --format json --output cluster-data.json

# CSV (Excel/spreadsheets)
k8s-cli export --format csv --output cluster-metrics.csv

# YAML (configuration)
k8s-cli export --format yaml --output cluster-config.yaml

# Prometheus (monitoring)
k8s-cli export --format prometheus --output metrics.prom
```

### What to Export
```bash
# Everything
k8s-cli export --all

# Only metrics
k8s-cli export --metrics --costs

# Only logs and events
k8s-cli export --logs --events

# Specific data types
k8s-cli export --resources --workloads --output combined-report.json
```

### Scheduled Exports
```bash
# Daily cost report
0 9 * * * k8s-cli cost --export csv --output ~/reports/daily-costs-$(date +\%Y\%m\%d).csv

# Weekly comprehensive report
0 9 * * 1 k8s-cli all --export json --output ~/reports/weekly-cluster-$(date +\%Y\%m\%d).json

# Hourly metrics for monitoring
0 * * * * k8s-cli metrics --export prometheus --output /var/lib/prometheus/k8s-cli-metrics.prom
```

## ⚙️ Configuration

### Environment Variables
```bash
# Set default namespace
export K8S_CLI_NAMESPACE=production

# Set default output format
export K8S_CLI_OUTPUT_FORMAT=json

# Set custom kubeconfig
export K8S_CLI_KUBECONFIG=/path/to/custom/config
```

### Config File
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
  default_node_cost: 72.0  # per month
```

## 🔍 Filtering and Selection

### Namespace Filtering
```bash
# Single namespace
k8s-cli resources --namespace production

# Multiple namespaces
k8s-cli resources --namespace "production,staging"

# Exclude namespaces
k8s-cli resources --exclude-namespace "kube-system,kube-public"

# All namespaces explicitly
k8s-cli resources --all-namespaces
```

### Resource Type Filtering
```bash
# Specific resource types
k8s-cli resources --type "deployment,service,ingress"

# Exclude resource types
k8s-cli resources --exclude-type "secret,configmap"

# Only workload resources
k8s-cli workload --type "deployment,statefulset,daemonset"
```

### Label and Field Selectors
```bash
# Filter by labels
k8s-cli resources --selector "app=frontend,env=production"

# Filter by field selectors
k8s-cli resources --field-selector "status.phase=Running"

# Combine filters
k8s-cli resources --selector "tier=web" --namespace production
```

## 📊 Understanding Output

### Table Format (Default)
```bash
NAMESPACE    NAME              TYPE         AGE    CPU     MEMORY   STATUS
production   frontend-app      Deployment   5d     150m    256Mi    Running
production   backend-api       Deployment   3d     300m    512Mi    Running
production   database          StatefulSet  10d    500m    1Gi      Running
```

### JSON Format (Programmatic)
```json
{
  "cluster_info": {
    "name": "production-cluster",
    "version": "v1.28.0",
    "nodes": 5,
    "namespaces": 12
  },
  "resources": [
    {
      "namespace": "production",
      "name": "frontend-app",
      "type": "Deployment",
      "age": "5d",
      "cpu": "150m",
      "memory": "256Mi",
      "status": "Running"
    }
  ]
}
```

### Cost Analysis Output
```bash
💰 Cost Analysis Report
=======================

Cluster Monthly Cost: $2,847.50

Top Cost Centers:
1. compute-intensive-app: $847.20 (29.7%)
2. database-cluster: $623.40 (21.9%)
3. monitoring-stack: $445.80 (15.7%)

💡 Optimization Opportunities:
- Right-size over-provisioned pods: Save ~$425/month
- Use spot instances for dev workloads: Save ~$230/month
- Optimize storage usage: Save ~$145/month

Total Potential Savings: $800/month (28.1%)
```

## 🚀 Advanced Usage

### Combining Commands
```bash
# Comprehensive analysis with export
k8s-cli all --namespace production --export json | \
  jq '.cost_analysis.recommendations[]' > recommendations.json

# Filter and process results
k8s-cli metrics --pods --utilization | \
  grep -E "(High|Critical)" > high-utilization-pods.txt

# Chain multiple analyses
k8s-cli cost --underutilized && \
k8s-cli workload --unhealthy-only && \
k8s-cli logs --critical --hours 1
```

### Automation Examples
```bash
#!/bin/bash
# Daily cluster health check

# Run comprehensive analysis
k8s-cli all --export json --output daily-report-$(date +%Y%m%d).json

# Check for critical issues
CRITICAL_ISSUES=$(k8s-cli logs --critical --hours 24 --count)
if [ "$CRITICAL_ISSUES" -gt 0 ]; then
    echo "⚠️ $CRITICAL_ISSUES critical issues found!"
    k8s-cli logs --critical --hours 24 | mail -s "Critical K8s Issues" devops@company.com
fi

# Cost optimization check
SAVINGS=$(k8s-cli cost --recommendations --potential-savings)
if [ "$SAVINGS" -gt 1000 ]; then
    echo "💰 Potential savings: $${SAVINGS}/month"
    k8s-cli cost --export csv --output cost-optimization-$(date +%Y%m%d).csv
fi
```

## 🆘 Getting Help

### Command Help
```bash
# General help
k8s-cli --help

# Command-specific help
k8s-cli resources --help
k8s-cli cost --help
k8s-cli export --help

# Examples for each command
k8s-cli resources --examples
```

### Troubleshooting
```bash
# Verbose output for debugging
k8s-cli --verbose resources

# Check cluster connectivity
k8s-cli version --cluster

# Validate configuration
k8s-cli config --validate
```

### Support Resources
- **Documentation**: [Complete Documentation](../docs/)
- **Examples**: [Usage Examples](../examples/)
- **Issues**: [GitHub Issues](https://github.com/Portfolio-jaime/lab-go-cli/issues)
- **Support**: jaime.andres.henao.arbelaez@ba.com

## 🎓 Next Steps

1. **Explore Advanced Features**: [API Reference](API.md)
2. **Set Up Automation**: [CI/CD Integration](CI_CD_CONSOLIDATED.md)
3. **Customize Configuration**: [Development Guide](DEVELOPMENT.md)
4. **Contribute**: [Contributing Guidelines](../CONTRIBUTING.md)

---

**🚀 You're now ready to optimize your Kubernetes clusters with k8s-cli!**