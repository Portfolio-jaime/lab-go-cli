#!/bin/bash
# test/e2e/test_devcontainer_tools.sh
# Verifica que las herramientas principales estén instaladas en el DevContainer

set -e

check_tool() {
    command -v "$1" >/dev/null 2>&1 && echo "[OK] $1 instalado" || { echo "[ERROR] $1 no está instalado"; exit 1; }
}

echo "[DevContainer] Verificando herramientas principales..."
check_tool docker
check_tool kubectl
check_tool helm
check_tool minikube
check_tool go
check_tool k8s-cli

echo "[DevContainer] Todas las herramientas principales están instaladas correctamente."
