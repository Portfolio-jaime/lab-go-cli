# DevContainer Configuration Guide

## 🐳 Overview

This DevContainer provides a complete Kubernetes development environment with Docker-in-Docker (DinD) and minikube running entirely inside a container. This approach ensures complete isolation from your host system while providing all necessary tools for Kubernetes development.

## 📁 File Structure

```
.devcontainer/
├── devcontainer.json      # VS Code DevContainer configuration
├── Dockerfile             # Container image definition
├── docker-compose.yml     # Services orchestration
└── scripts/
    ├── setup.sh           # Environment initialization
    └── start-minikube.sh  # Kubernetes cluster startup
```

## ⚙️ Configuration Files

### devcontainer.json

```json
{
  "name": "Go Development with Kubernetes Tools",
  "dockerComposeFile": "docker-compose.yml",
  "service": "devcontainer",
  "workspaceFolder": "/workspace",
  "shutdownAction": "stopCompose",
  "customizations": {
    "vscode": {
      "extensions": [
        "golang.go",
        "ms-kubernetes-tools.vscode-kubernetes-tools",
        "ms-vscode.vscode-yaml",
        "redhat.vscode-yaml"
      ],
      "settings": {
        "go.toolsManagement.checkForUpdates": "local",
        "go.useLanguageServer": true,
        "go.gopath": "/go",
        "go.goroot": "/usr/local/go",
        "terminal.integrated.cwd": "/workspace"
      }
    }
  },
  "remoteUser": "arheanja",
  "postCreateCommand": "/workspace/.devcontainer/scripts/setup.sh",
  "privileged": true
}
```

**Key Settings:**
- **User**: `arheanja` with full sudo privileges
- **Workspace**: `/workspace` mapped to project root
- **Extensions**: Go and Kubernetes development tools
- **Privileged**: Required for Docker-in-Docker

### Dockerfile

Multi-stage container with all development tools:

```dockerfile
FROM ubuntu:22.04

# Create arheanja user with root privileges
RUN groupadd --gid 1000 arheanja \
    && useradd --uid 1000 --gid arheanja --shell /bin/bash --create-home arheanja \
    && usermod -aG root arheanja

# Add arheanja to sudoers with full root privileges
RUN echo arheanja ALL=\(ALL:ALL\) NOPASSWD:ALL > /etc/sudoers.d/arheanja \
    && chmod 0440 /etc/sudoers.d/arheanja

# Install base tools
RUN apt-get update && apt-get install -y \
    git curl wget unzip apt-transport-https ca-certificates \
    gnupg lsb-release sudo vim nano jq conntrack \
    && rm -rf /var/lib/apt/lists/*

# Install Docker Engine (full, not just CLI)
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null \
    && apt-get update \
    && apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin \
    && rm -rf /var/lib/apt/lists/*

# Install Go with architecture detection
ENV GO_VERSION=1.21.5
RUN ARCH=$(dpkg --print-architecture) \
    && wget https://golang.org/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-${ARCH}.tar.gz \
    && rm go${GO_VERSION}.linux-${ARCH}.tar.gz

# Install kubectl with architecture detection
RUN ARCH=$(dpkg --print-architecture | sed 's/aarch64/arm64/') \
    && curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/${ARCH}/kubectl" \
    && install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl \
    && rm kubectl

# Install Helm
RUN curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Install minikube with architecture detection
RUN ARCH=$(dpkg --print-architecture) \
    && curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-${ARCH} \
    && chmod +x minikube \
    && mv minikube /usr/local/bin/

# Set up Go environment
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Create directories with proper permissions
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chown -R arheanja:arheanja "$GOPATH"
RUN mkdir -p /workspace && chown arheanja:arheanja /workspace
RUN mkdir -p /home/arheanja/.minikube /home/arheanja/.kube \
    && chown -R arheanja:arheanja /home/arheanja/.minikube /home/arheanja/.kube \
    && chmod -R 755 /home/arheanja/.minikube /home/arheanja/.kube

# Create Docker daemon startup script
RUN cat > /usr/local/bin/start-docker.sh << 'EOF'
#!/bin/bash
# Start Docker daemon
sudo dockerd --host=unix:///var/run/docker.sock --host=tcp://0.0.0.0:2376 &
sleep 5
echo "Docker daemon started"
EOF
RUN chmod +x /usr/local/bin/start-docker.sh

USER arheanja
WORKDIR /workspace

# Install Go development tools
RUN go install -v golang.org/x/tools/gopls@latest \
    && go install -v github.com/ramya-rao-a/go-outline@latest \
    && go install -v github.com/cweill/gotests/...@latest \
    && go install -v github.com/fatih/gomodifytags@latest \
    && go install -v github.com/josharian/impl@latest \
    && go install -v github.com/haya14busa/goplay/cmd/goplay@latest \
    && go install -v github.com/go-delve/delve/cmd/dlv@latest \
    && go install -v honnef.co/go/tools/cmd/staticcheck@latest
```

**Features:**
- **Multi-architecture support**: ARM64 and AMD64
- **Complete Docker Engine**: Not just CLI
- **Development tools**: Go tools, debugging, linting
- **Proper permissions**: User setup with sudo access

### docker-compose.yml

```yaml
services:
  devcontainer:
    build:
      context: .
      dockerfile: Dockerfile
    volumes:
      - ..:/workspace:cached
      - minikube-data:/home/arheanja/.minikube
      - kube-data:/home/arheanja/.kube
      - docker-data:/var/lib/docker
    command: sleep infinity
    privileged: true
    network_mode: host
    environment:
      - MINIKUBE_IN_DOCKER=true
    cap_add:
      - SYS_ADMIN
    security_opt:
      - seccomp:unconfined
    tmpfs:
      - /tmp
      - /run

volumes:
  minikube-data:
  kube-data:
  docker-data:
```

**Key Features:**
- **Privileged mode**: Required for Docker-in-Docker
- **Host networking**: Better performance for minikube
- **Persistent volumes**: Data survives container restarts
- **Security options**: Necessary for container orchestration

## 🔧 Initialization Scripts

### setup.sh

Environment initialization script that runs on container creation:

```bash
#!/bin/bash
set -e

echo "Setting up development environment..."

# Make sure Docker is accessible
sudo usermod -aG docker arheanja 2>/dev/null || true

# Ensure proper permissions for minikube and kubectl directories
sudo chown -R arheanja:arheanja ~/.minikube ~/.kube 2>/dev/null || true
chmod -R u+wrx ~/.minikube ~/.kube 2>/dev/null || true

# Configure kubectl completion
echo 'source <(kubectl completion bash)' >> ~/.bashrc

# Configure helm completion
echo 'source <(helm completion bash)' >> ~/.bashrc

# Configure minikube completion
echo 'source <(minikube completion bash)' >> ~/.bashrc

# Add Go tools to path
echo 'export PATH=$PATH:$GOPATH/bin:$GOROOT/bin' >> ~/.bashrc

# Create useful aliases
cat >> ~/.bashrc << 'EOF'

# Kubernetes aliases
alias k=kubectl
alias mk=minikube

# Go aliases  
alias gob='go build'
alias gor='go run'
alias got='go test'
alias gom='go mod'

EOF

echo "Development environment setup complete!"
echo "To start minikube, run: ./.devcontainer/scripts/start-minikube.sh"
```

### start-minikube.sh

Kubernetes cluster startup script:

```bash
#!/bin/bash
echo "Starting Docker daemon inside container..."

# Start Docker daemon if not running
if ! docker info >/dev/null 2>&1; then
    echo "Starting Docker daemon..."
    start-docker.sh
    sleep 10
fi

echo "Starting minikube with Docker-in-Docker..."

# Start minikube with docker driver inside container
minikube start \
  --driver=docker \
  --container-runtime=docker \
  --cpus=2 \
  --memory=3g \
  --disk-size=20g \
  --kubernetes-version=stable \
  --embed-certs \
  --apiserver-ips=127.0.0.1 \
  --apiserver-name=localhost

echo "Minikube started successfully!"

# Enable useful addons
echo "Enabling minikube addons..."
minikube addons enable dashboard
minikube addons enable metrics-server
minikube addons enable ingress

echo "Setup complete! Use 'minikube dashboard' to open the Kubernetes dashboard."
echo "Use 'kubectl get nodes' to verify cluster is running."
```

## 🚀 Usage

### Opening the DevContainer

1. **Prerequisites:**
   - Docker Desktop running
   - VS Code with DevContainers extension

2. **Open Project:**
   ```bash
   code .
   # Select "Reopen in Container"
   ```

3. **Wait for Build:**
   - First time: ~10-15 minutes
   - Subsequent: ~30 seconds

### Starting Kubernetes

```bash
# Inside the DevContainer terminal
./.devcontainer/scripts/start-minikube.sh

# Verify cluster
kubectl get nodes
minikube status
```

### Development Workflow

```bash
# 1. Build Go CLI
go mod tidy
go build -o k8s-cli .

# 2. Test CLI
./k8s-cli --help
./k8s-cli all

# 3. Deploy test applications
kubectl create deployment nginx --image=nginx
kubectl get pods

# 4. Test with real workload
./k8s-cli resources
```

## 🔍 Architecture Details

### Docker-in-Docker (DinD)

The DevContainer uses DinD to run minikube completely inside the container:

```
┌─────────────────────────────────────────┐
│           Host Docker                   │
│  ┌───────────────────────────────────┐  │
│  │        DevContainer               │  │
│  │  ┌─────────────────────────────┐  │  │
│  │  │      Docker Daemon          │  │  │
│  │  │  ┌───────────────────────┐  │  │  │
│  │  │  │      Minikube         │  │  │  │
│  │  │  │   ┌─────────────────┐ │  │  │  │
│  │  │  │   │   K8s Pods      │ │  │  │  │
│  │  │  │   └─────────────────┘ │  │  │  │
│  │  │  └───────────────────────┘  │  │  │
│  │  └─────────────────────────────┘  │  │
│  │  Go CLI + Development Tools       │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

**Benefits:**
- **Complete Isolation**: Nothing runs on host
- **Reproducible Environment**: Same setup everywhere
- **Easy Cleanup**: Delete container, everything gone
- **Version Control**: Environment as code

### Persistent Storage

Three volumes ensure data persistence:

1. **minikube-data**: `/home/arheanja/.minikube`
   - Minikube configuration
   - Cluster certificates
   - Node data

2. **kube-data**: `/home/arheanja/.kube`
   - kubectl configuration
   - Cluster credentials
   - Context information

3. **docker-data**: `/var/lib/docker`
   - Docker images
   - Container data
   - Build cache

### Networking

- **Host Network Mode**: Optimal performance for minikube
- **Port Access**: Direct access to cluster services
- **DNS Resolution**: Automatic service discovery

## 🔧 Customization

### Resource Limits

Adjust minikube resources in `start-minikube.sh`:

```bash
# Reduce memory for low-resource systems
minikube start --memory=2g --cpus=1

# Increase for high-performance development
minikube start --memory=8g --cpus=4
```

### Additional Tools

Add tools to Dockerfile:

```dockerfile
# Install additional CLI tools
RUN curl -LO https://github.com/derailed/k9s/releases/latest/download/k9s_Linux_x86_64.tar.gz \
    && tar -xzf k9s_Linux_x86_64.tar.gz -C /usr/local/bin/ k9s \
    && rm k9s_Linux_x86_64.tar.gz

# Install Terraform
RUN curl -LO https://releases.hashicorp.com/terraform/1.6.0/terraform_1.6.0_linux_amd64.zip \
    && unzip terraform_1.6.0_linux_amd64.zip -d /usr/local/bin/ \
    && rm terraform_1.6.0_linux_amd64.zip
```

### VS Code Extensions

Add extensions in `devcontainer.json`:

```json
"extensions": [
  "golang.go",
  "ms-kubernetes-tools.vscode-kubernetes-tools",
  "ms-vscode.vscode-yaml",
  "redhat.vscode-yaml",
  "hashicorp.terraform",
  "ms-vscode.docker"
]
```

## 🐛 Troubleshooting

### Container Build Issues

```bash
# Clean Docker system
docker system prune -af

# Rebuild container
# In VS Code: Command Palette > "Dev Containers: Rebuild Container"
```

### Docker Daemon Issues

```bash
# Check Docker daemon status inside container
sudo service docker status

# Restart Docker daemon
sudo service docker restart

# Manual Docker start
sudo dockerd --host=unix:///var/run/docker.sock &
```

### Minikube Issues

```bash
# Check minikube status
minikube status

# View minikube logs
minikube logs

# Delete and recreate cluster
minikube delete
./.devcontainer/scripts/start-minikube.sh
```

### Permission Issues

```bash
# Fix ownership of kube directories
sudo chown -R arheanja:arheanja ~/.minikube ~/.kube
chmod -R u+wrx ~/.minikube ~/.kube

# Add user to docker group
sudo usermod -aG docker arheanja
```

### Resource Issues

```bash
# Check container resources
docker stats

# Check host resources
free -h
df -h
```

## 💡 Best Practices

### Development

1. **Use Persistent Volumes**: Don't store important data in container filesystem
2. **Regular Commits**: DevContainer changes should be version controlled
3. **Resource Management**: Monitor container resource usage
4. **Clean Rebuilds**: Periodically rebuild container from scratch

### Performance

1. **Resource Allocation**: Adjust minikube memory/CPU based on needs
2. **Image Optimization**: Use multi-stage builds, minimize layers
3. **Volume Mounts**: Use cached mounts for better performance
4. **Network Mode**: Host networking for better minikube performance

### Security

1. **Privileged Mode**: Only use when necessary (required for DinD)
2. **User Permissions**: Use non-root user when possible
3. **Secrets**: Never include secrets in Dockerfile
4. **Updates**: Regularly update base images and tools

## 📚 Additional Resources

- [DevContainers Documentation](https://containers.dev/)
- [Docker-in-Docker Guide](https://docs.docker.com/engine/security/userns-remap/)
- [Minikube Documentation](https://minikube.sigs.k8s.io/docs/)
- [VS Code DevContainers](https://code.visualstudio.com/docs/devcontainers/containers)

---

🚀 **Happy Containerized Development!** 🎉