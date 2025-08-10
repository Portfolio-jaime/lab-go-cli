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
