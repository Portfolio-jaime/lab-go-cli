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

# Create alias for common commands
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

# Create minikube start script
cat > /workspace/.devcontainer/scripts/start-minikube.sh << 'EOF'
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
EOF

chmod +x /workspace/.devcontainer/scripts/start-minikube.sh

echo "Development environment setup complete!"
echo "To start minikube, run: ./.devcontainer/scripts/start-minikube.sh"