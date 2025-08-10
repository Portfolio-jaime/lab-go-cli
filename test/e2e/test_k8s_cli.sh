#!/bin/bash
# test/e2e/test_k8s_cli.sh
# Prueba end-to-end para k8s-cli

set -e

echo "[E2E] Iniciando minikube..."
minikube start

echo "[E2E] Compilando CLI..."
go build -o k8s-cli /workspace

# Prueba comando 'version'
echo "[E2E] Prueba end-to-end finalizada con éxito."

echo "[E2E] Ejecutando k8s-cli version..."
output_version=$(./k8s-cli version)
echo "$output_version"
echo "$output_version" | grep -q "Cluster Version Information" || { echo "[E2E] ERROR: No se detectó información de versión"; exit 1; }

echo "[E2E] Ejecutando k8s-cli resources..."
output_resources=$(./k8s-cli resources)
echo "$output_resources"
echo "$output_resources" | grep -q "Cluster Resource Summary" || { echo "[E2E] ERROR: No se detectó resumen de recursos"; exit 1; }

echo "[E2E] Ejecutando k8s-cli recommend..."
output_recommend=$(./k8s-cli recommend)
echo "$output_recommend"
echo "$output_recommend" | grep -q "recommendations" || { echo "[E2E] ERROR: No se detectaron recomendaciones"; exit 1; }

echo "[E2E] Validando integración con kubectl..."
kubectl get nodes | grep -q "minikube" || { echo "[E2E] ERROR: minikube no está en la lista de nodos"; exit 1; }

echo "[E2E] Todas las validaciones pasaron. Prueba end-to-end finalizada con éxito."
