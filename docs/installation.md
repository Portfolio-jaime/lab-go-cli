# Installation Guide

Complete installation guide for k8s-cli across different platforms and environments.

## 🚀 Quick Installation

### Binary Release (Recommended)
```bash
# Download latest release for Linux
curl -L https://github.com/Portfolio-jaime/lab-go-cli/releases/latest/download/k8s-cli-linux-amd64 -o k8s-cli
chmod +x k8s-cli
sudo mv k8s-cli /usr/local/bin/

# Download for macOS
curl -L https://github.com/Portfolio-jaime/lab-go-cli/releases/latest/download/k8s-cli-darwin-amd64 -o k8s-cli
chmod +x k8s-cli
sudo mv k8s-cli /usr/local/bin/

# Verify installation
k8s-cli --version
```

### Package Managers

#### Homebrew (macOS/Linux)
```bash
# Add tap (when available)
brew tap Portfolio-jaime/k8s-cli
brew install k8s-cli
```

#### APT (Debian/Ubuntu)
```bash
# Add repository (when available)
curl -fsSL https://k8s-cli.dev/gpg | sudo apt-key add -
echo "deb https://k8s-cli.dev/apt stable main" | sudo tee /etc/apt/sources.list.d/k8s-cli.list
sudo apt update
sudo apt install k8s-cli
```

## 📋 Prerequisites

### System Requirements
- **OS**: Linux (Ubuntu 18.04+), macOS (10.15+), Windows (WSL2)
- **Architecture**: x86_64, ARM64
- **Memory**: 512MB available RAM
- **Storage**: 50MB free space

### Kubernetes Access
- **kubectl** installed and configured
- **kubeconfig** file with cluster access
- **RBAC permissions** for resource reading

### Required Tools
```bash
# Essential tools
kubectl version --client    # Kubernetes CLI
curl --version             # HTTP client

# Optional but recommended
helm version               # Package manager
jq --version              # JSON processor
```

## 🛠️ Installation Methods

### 1. Binary Installation

#### Linux (x86_64)
```bash
# Create installation directory
sudo mkdir -p /opt/k8s-cli

# Download and install
curl -L https://github.com/Portfolio-jaime/lab-go-cli/releases/latest/download/k8s-cli-linux-amd64 \
  -o /tmp/k8s-cli
sudo install /tmp/k8s-cli /usr/local/bin/k8s-cli

# Cleanup
rm /tmp/k8s-cli

# Verify
k8s-cli version
```

#### macOS
```bash
# Intel Macs
curl -L https://github.com/Portfolio-jaime/lab-go-cli/releases/latest/download/k8s-cli-darwin-amd64 \
  -o /tmp/k8s-cli

# Apple Silicon Macs  
curl -L https://github.com/Portfolio-jaime/lab-go-cli/releases/latest/download/k8s-cli-darwin-arm64 \
  -o /tmp/k8s-cli

# Install
chmod +x /tmp/k8s-cli
sudo mv /tmp/k8s-cli /usr/local/bin/k8s-cli

# Verify
k8s-cli version
```

#### Windows (WSL2)
```bash
# Download in WSL2 environment
curl -L https://github.com/Portfolio-jaime/lab-go-cli/releases/latest/download/k8s-cli-linux-amd64 \
  -o k8s-cli
chmod +x k8s-cli
sudo mv k8s-cli /usr/local/bin/

# Alternative: Windows binary (future release)
# Download k8s-cli-windows-amd64.exe
```

### 2. Source Installation

#### Prerequisites for Building
```bash
# Install Go 1.21+
go version

# Install build tools
make --version
git --version
```

#### Build from Source
```bash
# Clone repository
git clone https://github.com/Portfolio-jaime/lab-go-cli.git
cd lab-go-cli

# Build using Makefile
make build

# Install locally
make install

# Or build manually
go build -o k8s-cli .
sudo mv k8s-cli /usr/local/bin/
```

### 3. Container Installation

#### Docker
```bash
# Run as container
docker run --rm -v ~/.kube:/root/.kube \
  ghcr.io/portfolio-jaime/k8s-cli:latest \
  k8s-cli --help

# Create alias for easy use
alias k8s-cli='docker run --rm -v ~/.kube:/root/.kube ghcr.io/portfolio-jaime/k8s-cli:latest k8s-cli'
```

#### Kubernetes Job
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: k8s-cli-analysis
spec:
  template:
    spec:
      containers:
      - name: k8s-cli
        image: ghcr.io/portfolio-jaime/k8s-cli:latest
        command: ["k8s-cli", "all", "--export", "json"]
        volumeMounts:
        - name: output
          mountPath: /output
      volumes:
      - name: output
        emptyDir: {}
      restartPolicy: Never
```

## ⚙️ Configuration

### Environment Variables
```bash
# Create configuration file
mkdir -p ~/.config/k8s-cli

# Set environment variables
export K8S_CLI_CONFIG_PATH=~/.config/k8s-cli/config.yaml
export K8S_CLI_CACHE_DIR=~/.cache/k8s-cli
export K8S_CLI_OUTPUT_FORMAT=table
export K8S_CLI_NAMESPACE=default
```

### Configuration File
```yaml
# ~/.config/k8s-cli/config.yaml
output:
  format: table        # table, json, yaml, csv
  colors: true         # Enable colored output
  pager: less          # Pager for long output

cluster:
  context: ""          # Kubernetes context (empty = current)
  namespace: "default" # Default namespace
  timeout: 30s         # API timeout

metrics:
  cache_duration: 5m   # Cache duration for metrics
  include_system: false # Include system namespaces

cost:
  currency: "USD"      # Currency for cost calculations
  node_cost_per_hour: 0.1  # Default node cost

export:
  directory: "./exports"    # Default export directory
  timestamp: true           # Include timestamp in filenames
```

## 🔧 Post-Installation Setup

### Kubernetes Access Verification
```bash
# Test Kubernetes connectivity
kubectl cluster-info

# Test k8s-cli access
k8s-cli version

# Run basic analysis
k8s-cli resources --namespace kube-system
```

### RBAC Setup (if needed)
```yaml
# k8s-cli-rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: k8s-cli-reader
rules:
- apiGroups: [""]
  resources: ["nodes", "pods", "services", "events", "namespaces"]
  verbs: ["get", "list"]
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets", "daemonsets", "replicasets"]
  verbs: ["get", "list"]
- apiGroups: ["metrics.k8s.io"]
  resources: ["nodes", "pods"]
  verbs: ["get", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: k8s-cli-binding
subjects:
- kind: User
  name: your-username
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: k8s-cli-reader
  apiGroup: rbac.authorization.k8s.io
```

```bash
# Apply RBAC
kubectl apply -f k8s-cli-rbac.yaml
```

### Shell Completion
```bash
# Bash completion
k8s-cli completion bash | sudo tee /etc/bash_completion.d/k8s-cli

# Zsh completion
k8s-cli completion zsh > ~/.zsh/completions/_k8s-cli

# Fish completion
k8s-cli completion fish > ~/.config/fish/completions/k8s-cli.fish
```

## 🚀 Verification

### Installation Test
```bash
# Version check
k8s-cli version

# Help command
k8s-cli --help

# Quick cluster test
k8s-cli resources --limit 5

# Full functionality test
k8s-cli all --dry-run
```

### Performance Test
```bash
# Metrics collection test
time k8s-cli metrics --nodes

# Export functionality test
k8s-cli export --format json --output /tmp/test-export.json
ls -la /tmp/test-export.json
```

## 🔄 Updates

### Updating k8s-cli
```bash
# Download latest version
curl -L https://github.com/Portfolio-jaime/lab-go-cli/releases/latest/download/k8s-cli-linux-amd64 \
  -o /tmp/k8s-cli-new

# Backup current version
sudo cp /usr/local/bin/k8s-cli /usr/local/bin/k8s-cli.backup

# Install new version
sudo install /tmp/k8s-cli-new /usr/local/bin/k8s-cli

# Verify update
k8s-cli version
```

### Automated Updates
```bash
# Create update script
cat > update-k8s-cli.sh << 'EOF'
#!/bin/bash
CURRENT_VERSION=$(k8s-cli version --short)
LATEST_VERSION=$(curl -s https://api.github.com/repos/Portfolio-jaime/lab-go-cli/releases/latest | jq -r .tag_name)

if [ "$CURRENT_VERSION" != "$LATEST_VERSION" ]; then
  echo "Updating k8s-cli from $CURRENT_VERSION to $LATEST_VERSION"
  curl -L https://github.com/Portfolio-jaime/lab-go-cli/releases/latest/download/k8s-cli-linux-amd64 -o /tmp/k8s-cli
  sudo install /tmp/k8s-cli /usr/local/bin/k8s-cli
  echo "Update completed"
else
  echo "k8s-cli is already up to date ($CURRENT_VERSION)"
fi
EOF

chmod +x update-k8s-cli.sh
```

## 🐛 Troubleshooting

### Common Issues

#### Permission Denied
```bash
# Fix executable permissions
chmod +x /usr/local/bin/k8s-cli

# Fix directory permissions
sudo chown $USER:$USER ~/.config/k8s-cli
```

#### Kubernetes Access Issues
```bash
# Check kubeconfig
kubectl config view
kubectl config current-context

# Test cluster access
kubectl get nodes
kubectl auth can-i get pods
```

#### Missing Dependencies
```bash
# Install missing tools
sudo apt update
sudo apt install -y curl jq

# Verify installations
which kubectl curl jq
```

### Getting Help
- **Documentation**: [docs/](../docs/)
- **Issues**: [GitHub Issues](https://github.com/Portfolio-jaime/lab-go-cli/issues)
- **Support**: jaime.andres.henao.arbelaez@ba.com