#!/bin/bash

# K8s-CLI Usage Examples
# This script shows how to use the various commands in k8s-cli

echo "=== K8s-CLI Usage Examples ==="
echo

echo "1. Show complete cluster analysis (recommended for first run):"
echo "./k8s-cli all"
echo

echo "2. Show cluster version and installed components:"
echo "./k8s-cli version"
echo

echo "3. Show cluster resource consumption:"
echo "./k8s-cli resources"
echo

echo "4. Show resource consumption for a specific namespace:"
echo "./k8s-cli resources --namespace kube-system"
echo

echo "5. Show only node resources:"
echo "./k8s-cli resources --nodes"
echo

echo "6. Show only pod resources:"
echo "./k8s-cli resources --pods"
echo

echo "7. Get cluster optimization recommendations:"
echo "./k8s-cli recommend"
echo

echo "8. Filter recommendations by severity:"
echo "./k8s-cli recommend --severity High"
echo

echo "9. Filter recommendations by type:"
echo "./k8s-cli recommend --type Resource"
echo

echo "10. Use a specific kubeconfig file:"
echo "./k8s-cli all --kubeconfig /path/to/kubeconfig"
echo

echo "=== Prerequisites ==="
echo "- kubectl configured with cluster access"
echo "- For resource metrics: metrics-server installed"
echo "- For full functionality: cluster-admin permissions"
echo

echo "=== Installation ==="
echo "go build -o k8s-cli"
echo "# or"
echo "make build"