# 📚 k8s-cli Usage Examples

## 🎯 Quick Start Examples

### Basic Cluster Analysis
```bash
# Complete cluster overview
k8s-cli all

# Kubernetes cluster version and enhanced component discovery
k8s-cli version

# CLI tool version information
k8s-cli --version
k8s-cli -v  # short form

# Basic resource information
k8s-cli resources
```

## 🔍 Enhanced Component Detection Examples (New in v2.0.1)

### CLI Version vs Cluster Version
```bash
# Show CLI tool version (build info, commit, etc.)
k8s-cli --version
# Output:
# k8s-cli version v2.0.1
# Git commit: abc123
# Built: 2025-08-13T10:30:00
# Go version: go1.24.6
# OS/Arch: darwin/arm64

# Show Kubernetes cluster version and components
k8s-cli version
# Output includes:
# - Kubernetes cluster version information
# - All components across ALL namespaces
# - Helm releases with version info
# - Source identification (Helm/Deployment/StatefulSet/DaemonSet)
```

### Component Detection Improvements
```bash
# Enhanced component discovery with Helm integration
k8s-cli version

# Expected improvements in v2.0.1:
# ✅ Finds components in ALL namespaces (not just predefined ones)
# ✅ Detects Helm releases automatically  
# ✅ Shows 25+ component types
# ✅ Identifies installation source (Helm vs K8s resources)
# ✅ Smart deduplication (prefers Helm info when available)
```

### Example Output Comparison

#### Before v2.0.1 (Limited Detection)
```
🔧 Installed Components:
┌─────────────────┬─────────────┬─────────┬─────────┬─────────┐
│ Component       │ Namespace   │ Status  │ Version │ Ready   │
├─────────────────┼─────────────┼─────────┼─────────┼─────────┤
│ metrics-server  │ kube-system │ Running │ v0.6.2  │ 1/1     │
│ nginx-ingress   │ ingress     │ Running │ 1.8.2   │ 2/2     │
└─────────────────┴─────────────┴─────────┴─────────┴─────────┘
```

#### After v2.0.1 (Enhanced Detection)
```
🔧 Installed Components:
   Searching in all namespaces for components and Helm releases...

   Found 12 components:

┌─────────────────┬─────────────┬─────────┬─────────┬─────────┬─────────────┐
│ Component       │ Namespace   │ Status  │ Version │ Ready   │ Source      │
├─────────────────┼─────────────┼─────────┼─────────┼─────────┼─────────────┤
│ metrics-server  │ kube-system │ Running │ v0.6.2  │ 1/1     │ Deployment  │
│ nginx-ingress   │ ingress     │ Running │ 1.8.2   │ 2/2     │ Deployment  │
│ prometheus      │ monitoring  │ Deployed│ 45.7.1  │ Helm    │ Helm        │
│ grafana         │ monitoring  │ Deployed│ 6.52.4  │ Helm    │ Helm        │
│ redis           │ cache       │ Running │ 7.0.8   │ 1/1     │ StatefulSet │
│ elasticsearch   │ logging     │ Running │ 8.6.2   │ 3/3     │ StatefulSet │
│ fluentd         │ logging     │ Running │ 1.16    │ Helm    │ Helm        │
│ cert-manager    │ cert-mgr    │ Running │ v1.11.0 │ 1/1     │ Deployment  │
│ vault           │ security    │ Deployed│ 0.23.0  │ Helm    │ Helm        │
│ argocd-server   │ argocd      │ Running │ v2.6.7  │ 1/1     │ Deployment  │
│ istio-proxy     │ istio-sys   │ Running │ 1.17.1  │ 2/2     │ DaemonSet   │
│ postgres        │ database    │ Running │ 15.2    │ 1/1     │ StatefulSet │
└─────────────────┴─────────────┴─────────┴─────────┴─────────┴─────────────┘
```

### Advanced Detection Features
```bash
# The enhanced version automatically detects:

# 1. Helm Releases in any namespace
#    - Reads Helm secrets with owner=helm label
#    - Extracts version from app.kubernetes.io/version or version labels
#    - Shows Helm release status (Deployed, Failed, etc.)

# 2. All Resource Types:
#    - Deployments (traditional apps)
#    - StatefulSets (databases, persistent apps)  
#    - DaemonSets (node agents, monitoring)
#    - Helm releases (chart-based deployments)

# 3. Extended Component Library:
#    metrics-server, argocd, argo, kuma, istio, traefik,
#    nginx, cert-manager, prometheus, grafana, jaeger,
#    kiali, fluentd, elasticsearch, kibana, vault,
#    consul, etcd, redis, postgres, mysql, mongodb,
#    kafka, zookeeper, rabbitmq, jenkins, sonarqube,
#    nexus, harbor, docker-registry, ingress, gateway

# 4. Smart Deduplication:
#    - If a component is found via both K8s resources AND Helm
#    - Prioritizes Helm information (more accurate versions)
#    - Prevents duplicate entries in the output
```

## 📊 Real-time Metrics Examples

### Node Metrics
```bash
# Show all node metrics
k8s-cli metrics --nodes

# Show node metrics with utilization analysis
k8s-cli metrics --nodes --utilization

# Example output:
# 🌐 CLUSTER OVERVIEW
# ┌─────────┬─────────┬─────────┬─────────────┐
# │ Metric  │ Usage   │ Capacity│ Utilization │
# ├─────────┼─────────┼─────────┼─────────────┤
# │ CPU     │ 2.1     │ 4.0     │ 52.5%       │
# │ Memory  │ 3.2 GiB │ 8.0 GiB │ 40.0%       │
# └─────────┴─────────┴─────────┴─────────────┘
```

### Pod Metrics
```bash
# Show pod metrics for all namespaces
k8s-cli metrics --pods

# Show pod metrics for specific namespace
k8s-cli metrics --pods --namespace production

# Show pods with high resource usage
k8s-cli metrics --pods --utilization | grep "⚠️"
```

### Combined Analysis
```bash
# Full metrics analysis
k8s-cli metrics --nodes --pods --utilization

# Focus on production workloads
k8s-cli metrics --namespace production --pods --utilization
```

## 💰 Cost Analysis Examples

### Basic Cost Analysis
```bash
# Complete cost analysis
k8s-cli cost

# Example output:
# 💰 COST OVERVIEW
# ┌─────────────────────────┬─────────┐
# │ Metric                  │ Value   │
# ├─────────────────────────┼─────────┤
# │ Total Monthly Cost      │ $324.50 │
# │ Potential Monthly Savings│ $89.20  │
# │ Cost Efficiency         │ 72.5%   │
# └─────────────────────────┴─────────┘
```

### Node Cost Breakdown
```bash
# Show only node costs
k8s-cli cost --nodes --no-namespaces --no-underutilized --no-optimizations

# Focus on underutilized resources
k8s-cli cost --underutilized
```

### Cost Optimization
```bash
# Show only optimization recommendations
k8s-cli cost --optimizations --no-nodes --no-namespaces --no-underutilized

# Export cost analysis for reporting
k8s-cli export --format csv --costs --output ./reports/
```

## 🔍 Workload Health Examples

### Basic Health Check
```bash
# Analyze all workloads
k8s-cli workload

# Example output:
# 📊 WORKLOAD SUMMARY
# ┌─────────────┬───────┬─────────┬─────────────┐
# │ Type        │ Total │ Healthy │ Health Rate │
# ├─────────────┼───────┼─────────┼─────────────┤
# │ Deployments │ 12    │ 10      │ 83.3%       │
# │ StatefulSets│ 3     │ 3       │ 100.0%      │
# │ DaemonSets  │ 4     │ 4       │ 100.0%      │
# └─────────────┴───────┴─────────┴─────────────┘
```

### Problem Detection
```bash
# Show only unhealthy workloads
k8s-cli workload --unhealthy-only

# Focus on specific namespace
k8s-cli workload --namespace production --unhealthy-only

# Detailed pod analysis
k8s-cli workload --pods --unhealthy-only
```

### Deployment Analysis
```bash
# Show only deployments
k8s-cli workload --deployments --no-statefulsets --no-daemonsets --no-summary

# Get detailed deployment health
k8s-cli workload --deployments --namespace production
```

## 📋 Logs and Events Examples

### Recent Events
```bash
# Show critical events from last 24 hours
k8s-cli logs --critical

# Show events from last hour
k8s-cli logs --hours 1

# Example output:
# 🚨 CRITICAL EVENTS
# ┌──────────┬─────────────────┬─────────────────┬─────────┬───────┐
# │ Time     │ Object          │ Reason          │ Message │ Count │
# ├──────────┼─────────────────┼─────────────────┼─────────┼───────┤
# │ 14:23:45 │ Pod/nginx-xyz   │ FailedScheduling│ No nodes│ 3     │
# └──────────┴─────────────────┴─────────────────┴─────────┴───────┘
```

### Error Patterns
```bash
# Analyze error patterns
k8s-cli logs --patterns

# Show resource-related events
k8s-cli logs --resource-events

# Security events
k8s-cli logs --security-events
```

### Namespace-specific Logs
```bash
# Analyze logs for specific namespace
k8s-cli logs --namespace kube-system --critical --patterns

# Pod-level analysis
k8s-cli logs --pod-analysis --namespace production
```

## 📤 Export Examples

### JSON Export
```bash
# Export complete cluster data
k8s-cli export --format json

# Export specific data types
k8s-cli export --format json --metrics --costs --no-logs --no-events

# Custom filename and location
k8s-cli export --format json --filename cluster-report --output ./reports/
```

### CSV Export
```bash
# Export to CSV format
k8s-cli export --format csv

# This creates multiple CSV files:
# - exports/node-metrics-2024-01-15-14-30-00.csv
# - exports/pod-metrics-2024-01-15-14-30-00.csv
# - exports/cost-analysis-2024-01-15-14-30-00.csv
# - exports/events-2024-01-15-14-30-00.csv
```

### Prometheus Export
```bash
# Export Prometheus metrics
k8s-cli export --format prometheus

# Example output file content:
# # HELP k8s_cluster_cpu_usage_percent Cluster CPU usage percentage
# # TYPE k8s_cluster_cpu_usage_percent gauge
# k8s_cluster_cpu_usage_percent 52.50 1705330200
```

### Automated Exports
```bash
# Daily export script
#!/bin/bash
DATE=$(date +%Y-%m-%d)
k8s-cli export --format json --filename "daily-report-$DATE" --output ./daily-reports/

# Export for monitoring
k8s-cli export --format prometheus --output /var/lib/prometheus/
```

## 🔄 Advanced Use Cases

### DevOps Monitoring Pipeline
```bash
#!/bin/bash
# monitoring-pipeline.sh

echo "🔍 Daily Cluster Health Check"

# 1. Check cluster health
k8s-cli workload --unhealthy-only > /tmp/unhealthy.txt
if [ -s /tmp/unhealthy.txt ]; then
    echo "⚠️ Unhealthy workloads detected!"
    cat /tmp/unhealthy.txt
fi

# 2. Check for critical events
k8s-cli logs --critical --hours 24 > /tmp/critical.txt
if [ -s /tmp/critical.txt ]; then
    echo "🚨 Critical events in last 24h!"
    cat /tmp/critical.txt
fi

# 3. Cost optimization check
k8s-cli cost --underutilized | grep -q "underutilized"
if [ $? -eq 0 ]; then
    echo "💰 Cost optimization opportunities found"
    k8s-cli cost --underutilized
fi

# 4. Export daily report
k8s-cli export --format json --filename "daily-$(date +%Y%m%d)"
```

### FinOps Cost Optimization
```bash
#!/bin/bash
# cost-optimization.sh

echo "💰 Weekly Cost Optimization Report"

# Generate cost report
k8s-cli cost > weekly-cost-report.txt

# Find biggest cost savings
k8s-cli cost --underutilized --no-nodes --no-namespaces | \
    grep "Monthly Savings" | head -10

# Export for finance team
k8s-cli export --format csv --costs --filename "weekly-costs-$(date +%Y%U)"

# Calculate potential savings
SAVINGS=$(k8s-cli cost --underutilized | grep "Total potential" | awk '{print $5}')
echo "💡 Potential monthly savings: $SAVINGS"
```

### SRE Incident Response
```bash
#!/bin/bash
# incident-response.sh

NAMESPACE=${1:-"production"}
echo "🚨 Incident Response for namespace: $NAMESPACE"

# 1. Quick workload health check
echo "📊 Workload Health:"
k8s-cli workload --namespace $NAMESPACE --unhealthy-only

# 2. Recent critical events
echo "🚨 Critical Events (last 2 hours):"
k8s-cli logs --namespace $NAMESPACE --critical --hours 2

# 3. Resource pressure
echo "📈 Resource Metrics:"
k8s-cli metrics --namespace $NAMESPACE --pods --utilization

# 4. Export incident data
k8s-cli export --format json --namespace $NAMESPACE --hours 2 \
    --filename "incident-$(date +%Y%m%d-%H%M%S)"

echo "📋 Incident data exported for analysis"
```

### CI/CD Integration
```bash
#!/bin/bash
# ci-cd-health-check.sh

# Pre-deployment health check
echo "🔍 Pre-deployment cluster health check"

# Check cluster capacity
CPU_USAGE=$(k8s-cli metrics | grep "CPU" | awk '{print $4}' | tr -d '%')
if [ "$CPU_USAGE" -gt 80 ]; then
    echo "❌ High CPU usage: ${CPU_USAGE}%. Deployment may fail."
    exit 1
fi

# Check for critical issues
CRITICAL_COUNT=$(k8s-cli logs --critical --hours 1 | wc -l)
if [ "$CRITICAL_COUNT" -gt 5 ]; then
    echo "❌ Too many critical events: $CRITICAL_COUNT"
    exit 1
fi

echo "✅ Cluster health check passed"

# Post-deployment verification
echo "🚀 Post-deployment verification"
sleep 30  # Wait for deployment

# Check new workload health
k8s-cli workload --namespace $DEPLOY_NAMESPACE --unhealthy-only
```

## 🎯 Troubleshooting Common Issues

### No Metrics Available
```bash
# Check if metrics-server is running
kubectl get pods -n kube-system | grep metrics-server

# If not available, install metrics-server
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Wait and retry
k8s-cli metrics --nodes
```

### Permission Denied
```bash
# Check current permissions
kubectl auth can-i list pods
kubectl auth can-i get nodes

# Create service account with proper permissions
kubectl create serviceaccount k8s-cli-sa
kubectl create clusterrolebinding k8s-cli-binding \
    --clusterrole=view --serviceaccount=default:k8s-cli-sa

# Use service account token
kubectl get secret k8s-cli-sa-token -o jsonpath='{.data.token}' | base64 -d
```

### Large Cluster Performance
```bash
# Use namespace filtering for large clusters
k8s-cli metrics --namespace production

# Export data instead of displaying
k8s-cli export --format json --namespace production

# Limit time windows
k8s-cli logs --hours 1 --namespace production
```

## 🔧 Configuration Examples

### Custom Configuration File
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
  
export:
  default_format: json
  default_output: ./exports
  
filters:
  exclude_namespaces:
    - kube-system
    - kube-public
```

### Environment Variables
```bash
# Set default namespace
export K8S_CLI_NAMESPACE=production

# Set output format
export K8S_CLI_OUTPUT_FORMAT=json

# Set custom kubeconfig
export K8S_CLI_KUBECONFIG=/path/to/custom/kubeconfig

# Enable debug mode
export K8S_CLI_DEBUG=true
```

## 📈 Performance Examples

### Benchmarking
```bash
# Time analysis execution
time k8s-cli all

# Benchmark metrics collection
time k8s-cli metrics --nodes --pods

# Large cluster optimization
k8s-cli metrics --namespace production  # Focus on one namespace
```

### Monitoring Resource Usage
```bash
# Monitor CLI resource usage
/usr/bin/time -v k8s-cli all

# Memory-efficient export
k8s-cli export --format csv  # Uses less memory than JSON
```