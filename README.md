# Lab Go CLI - Kubernetes Development Environment

🚀 **Complete DevOps Development Environment** with a Go CLI for Kubernetes cluster analysis, running in a fully containerized DevContainer with Docker-in-Docker (DinD) and Minikube.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Quick Start](#quick-start)
- [📖 User Guides](#-user-guides)
  - [Guide: Managing the Minikube Cluster](#guide-managing-the-minikube-cluster)
  - [Guide: Using the `k8s-cli` Tool](#guide-using-the-k8s-cli-tool)
- [🛠️ Development Workflow](#️-development-workflow)
- [🔧 Troubleshooting](#-troubleshooting)

## 🎯 Overview

This project provides a complete, isolated Kubernetes development environment that includes:

1.  **🐳 DevContainer**: A fully containerized development environment using VS Code.
2.  **☸️ Minikube**: A local Kubernetes cluster that runs inside the container via Docker-in-Docker.
3.  **🔧 Go CLI**: A custom tool for Kubernetes cluster analysis and monitoring.
4.  **🛠️ DevOps Tools**: Comes pre-installed with `kubectl`, `helm`, `go`, `docker`, and more.

The entire environment runs inside a single container, providing complete isolation from your host system.

## ✨ Features

### DevContainer Environment
- **🔒 Isolated Development**: Complete separation from the host system.
- **🐳 Docker-in-Docker**: The Docker daemon starts automatically, enabling full Docker functionality inside the container.
- **☸️ Containerized Minikube**: Run a local Kubernetes cluster inside your development environment.
- **👤 Custom User**: A non-root `arheanja` user with passwordless `sudo` privileges is provided.
- **🚀 Auto-Setup**: The environment is configured automatically on build.

### Kubernetes CLI Tool
- **🔍 Cluster Analysis**: Detailed version and component information.
- **📊 Resource Monitoring**: Real-time resource utilization.
- **💡 Smart Recommendations**: Optimization suggestions.
- **🎨 Beautiful Output**: Formatted tables with colors.
- **🔧 Component Detection**: Auto-detects installed K8s components like Ingress controllers and service meshes.

## 🚀 Quick Start

1.  **Open in DevContainer**:
    - Clone this repository and open the folder in VS Code.
    - When prompted, select **"Reopen in Container"**.
    - VS Code will build the Docker image and start the dev container. The Docker daemon will start automatically.

2.  **Start Kubernetes Cluster**:
    - Once inside the DevContainer, open a terminal (`Terminal > New Terminal`) and start the Minikube cluster manually:
      ```bash
      minikube start
      ```

3.  **Install and Test the CLI**:
    - With the cluster running, install the Go CLI to make it available system-wide:
      ```bash
      go install .
      ```
    - Test the CLI with the most comprehensive command:
      ```bash
      k8s-cli all
      ```

## 📖 User Guides

### Guide: Managing the Minikube Cluster

All `minikube` commands should be run from the terminal inside the DevContainer.

- **Start the cluster:**
  ```bash
  minikube start
  ```

- **Check the status:**
  ```bash
  minikube status
  kubectl get nodes
  ```

- **Interact with the cluster:**
  ```bash
  # List all pods in all namespaces
  kubectl get pods -A

  # Create a test deployment
  kubectl create deployment hello-world --image=gcr.io/google-samples/hello-app:1.0
  ```

- **Open the Kubernetes Dashboard:**
  - Minikube's dashboard is a web-based UI for your cluster.
  ```bash
  # This command will run and print a URL. Copy the URL and open it in a browser on your host machine.
  minikube dashboard
  ```

- **Stop the cluster (to save resources):**
  ```bash
  minikube stop
  ```

- **Delete the cluster (to start fresh):**
  ```bash
  minikube delete
  ```

### Guide: Using the `k8s-cli` Tool

This tool helps you analyze the state of your Kubernetes cluster.

#### Installation

To make the `k8s-cli` command available system-wide inside the container, use `go install`. This is the recommended method.

```bash
# Run from the root of the project directory
go install .
```
After running this, you can execute `k8s-cli` from any directory.

#### Building (for Development)

If you are actively changing the code, you may prefer to build a local executable instead of installing it each time.

```bash
# This creates an executable file named 'k8s-cli' in the current directory
go build -o k8s-cli .

# You must run it using the relative path
./k8s-cli all
```

#### Usage Examples

- **Get a complete overview of the cluster:**
  ```bash
  k8s-cli all
  ```

- **Check cluster version and detected components:**
  ```bash
  k8s-cli version
  ```

- **Monitor resource usage:**
  ```bash
  k8s-cli resources
  ```

- **Get optimization recommendations:**
  ```bash
  k8s-cli recommend --severity Medium
  ```

## 🛠️ Development Workflow

1.  **Start Development**:
    - Open the project in VS Code and select **"Reopen in Container"**.

2.  **Initialize Kubernetes**:
    - Open a terminal and run `minikube start`.
    - Verify the cluster is running with `kubectl get nodes`.

3.  **Develop the CLI**:
    - Make changes to the Go source code in the `cmd/` and `pkg/` directories.
    - Rebuild the CLI after making changes:
      ```bash
      go build -o k8s-cli .
      ```
    - Test the changes against your running Minikube cluster:
      ```bash
      ./k8s-cli resources
      ```

## 🔧 Troubleshooting

### Minikube Fails to Start with `cgroup` Errors

- **Symptom**: `minikube start` fails with errors like `cannot enter cgroupv2...` or `...read-only file system`.
- **Cause**: An incompatibility between the host system's cgroup v2 implementation and the Docker-in-Docker environment.
- **Solution**: The final, working configuration requires specific settings in `.devcontainer/docker-compose.yml` to align the container's cgroup handling with the host's.

```yaml
# In .devcontainer/docker-compose.yml
services:
  devcontainer:
    # ... other settings
    volumes:
      # This volume mount must be read-write
      - /sys/fs/cgroup:/sys/fs/cgroup:rw
    # This setting is crucial
    cgroup: host
```

### DevContainer Build Fails on `go install`

- **Symptom**: The `docker build` process fails during the `RUN go install...` step with network errors.
- **Cause**: A transient network issue or a firewall/VPN interfering with the container's ability to download packages.
- **Solution**:
    1.  The `RUN go install...` step in the `Dockerfile` has been commented out by default.
    2.  After the container starts, you must **install these tools manually** by running the command provided in the terminal upon connection.