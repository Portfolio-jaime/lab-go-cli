#!/usr/bin/env python3
import os
import subprocess
import sys
from pathlib import Path

def check_ssh_installed():
    print("\n[1] Verificando si SSH está instalado...")
    result = subprocess.run(["which", "ssh"], capture_output=True, text=True)
    if result.returncode == 0:
        print("✅ SSH está instalado en:", result.stdout.strip())
        return True
    else:
        print("❌ SSH no está instalado.")
        return False

def check_ssh_keys():
    print("\n[2] Buscando llaves SSH para GitHub...")
    ssh_dir = Path.home() / ".ssh"
    keys = list(ssh_dir.glob("id_*"))
    github_keys = [k for k in keys if k.is_file() and not k.name.endswith(".pub")]
    if github_keys:
        print("✅ Llaves SSH encontradas:")
        for k in github_keys:
            print("  -", k)
        return True
    else:
        print("❌ No se encontraron llaves SSH en ~/.ssh/")
        gen = input("¿Deseas generar una nueva llave SSH ahora? (s/n): ").strip().lower()
        if gen == "s":
            email = input("Ingresa tu email para GitHub: ").strip()
            subprocess.run(["ssh-keygen", "-t", "ed25519", "-C", email])
            print("Llave generada. Verifica ~/.ssh/")
        return False
def check_github_pat():
    print("\n[3] Buscando Personal Access Token (PAT) de GitHub en variables de entorno...")
    pat_vars = ["GH_TOKEN", "GITHUB_TOKEN", "GH_PAT", "GITHUB_PAT"]
    found = False
    for var in pat_vars:
        if os.environ.get(var):
            print(f"✅ PAT encontrado en la variable de entorno: {var}")
            found = True
    if not found:
        print("❌ No se encontró PAT en variables de entorno comunes.")
        set_pat = input("¿Deseas configurar un PAT ahora? (s/n): ").strip().lower()
        if set_pat == "s":
            print("Puedes exportar tu token con: export GITHUB_TOKEN=tu_token")
    return found

def check_github_app():
    print("\n[4] Buscando credenciales de GitHub App en variables de entorno...")
    app_vars = ["GITHUB_APP_ID", "GITHUB_APP_PRIVATE_KEY", "GH_APP_ID", "GH_APP_KEY"]
    found = False
    for var in app_vars:
        if os.environ.get(var):
            print(f"✅ Credencial de GitHub App encontrada: {var}")
            found = True
    if not found:
        print("❌ No se encontraron credenciales de GitHub App en variables de entorno comunes.")
    return found
def show_ssh_public_key():
    print("\n[Mostrar llave pública SSH]")
    ssh_dir = Path.home() / ".ssh"
    pub_keys = list(ssh_dir.glob("*.pub"))
    if not pub_keys:
        print("❌ No se encontraron llaves públicas en ~/.ssh/")
        return
    print("Llaves públicas encontradas:")
    for idx, k in enumerate(pub_keys):
        print(f"{idx+1}. {k}")
    if len(pub_keys) == 1:
        choice = 1
    else:
        choice = input(f"Selecciona el número de la llave que deseas mostrar (1-{len(pub_keys)}): ").strip()
        try:
            choice = int(choice)
        except ValueError:
            print("Opción inválida.")
            return
        if not (1 <= choice <= len(pub_keys)):
            print("Opción fuera de rango.")
            return
    key_path = pub_keys[choice-1]
    print(f"\nContenido de {key_path}:")
    with open(key_path, "r") as f:
        print(f.read())
    print("\nCopia la llave y pégala en GitHub -> Settings -> SSH and GPG keys.")
def check_github_oauth():
    print("\n[5] Buscando credenciales OAuth de GitHub en variables de entorno...")
    oauth_vars = ["GITHUB_OAUTH_CLIENT_ID", "GITHUB_OAUTH_CLIENT_SECRET", "GH_OAUTH_ID", "GH_OAUTH_SECRET"]
    found = False
    for var in oauth_vars:
        if os.environ.get(var):
            print(f"✅ Credencial OAuth encontrada: {var}")
            found = True
    if not found:
        print("❌ No se encontraron credenciales OAuth en variables de entorno comunes.")
    return found

def check_git_config():
    print("\n[6] Verificando configuración básica de git...")
    result = subprocess.run(["git", "config", "--get", "user.name"], capture_output=True, text=True)
    if result.stdout.strip():
        print("✅ Nombre de usuario git configurado:", result.stdout.strip())
    else:
        print("❌ No se ha configurado el nombre de usuario en git.")
        set_name = input("¿Deseas configurarlo ahora? (s/n): ").strip().lower()
        if set_name == "s":
            name = input("Ingresa tu nombre de usuario para git: ").strip()
            subprocess.run(["git", "config", "--global", "user.name", name])
    result = subprocess.run(["git", "config", "--get", "user.email"], capture_output=True, text=True)
    if result.stdout.strip():
        print("✅ Email de git configurado:", result.stdout.strip())
    else:
        print("❌ No se ha configurado el email en git.")
        set_email = input("¿Deseas configurarlo ahora? (s/n): ").strip().lower()
        if set_email == "s":
            email = input("Ingresa tu email para git: ").strip()
            subprocess.run(["git", "config", "--global", "user.email", email])

def check_github_connectivity():
    print("\n[7] Probando conectividad con GitHub...")
    result = subprocess.run(["ssh", "-T", "git@github.com"], capture_output=True, text=True)
    if "successfully authenticated" in result.stderr or "Hi" in result.stderr:
        print("✅ Conexión SSH a GitHub exitosa.")
    else:
        print("❌ No se pudo conectar a GitHub por SSH. Mensaje:")
        print(result.stderr.strip())

def main():
    print("\n=== Validador Interactivo de Entorno GitHub ===")
    while True:
        print("\nMenú principal:")
        print("1. Verificar si SSH está instalado")
        print("2. Buscar llaves SSH para GitHub")
        print("3. Buscar Personal Access Token (PAT) de GitHub")
        print("4. Buscar credenciales de GitHub App")
        print("5. Buscar credenciales OAuth de GitHub")
        print("6. Verificar configuración básica de git")
        print("7. Probar conectividad con GitHub")
        print("8. Ejecutar todas las validaciones")
        print("9. Mostrar llave pública SSH para copiar en GitHub")
        print("0. Salir")
        opcion = input("\nSelecciona una opción (0-9): ").strip()
        if opcion == "1":
            check_ssh_installed()
        elif opcion == "2":
            check_ssh_keys()
        elif opcion == "3":
            check_github_pat()
        elif opcion == "4":
            check_github_app()
        elif opcion == "5":
            check_github_oauth()
        elif opcion == "6":
            check_git_config()
        elif opcion == "7":
            check_github_connectivity()
        elif opcion == "8":
            check_ssh_installed()
            check_ssh_keys()
            check_github_pat()
            check_github_app()
            check_github_oauth()
            check_git_config()
            check_github_connectivity()
            print("\nValidación completa. Revisa los resultados arriba.")
        elif opcion == "9":
            show_ssh_public_key()
        elif opcion == "0":
            print("Saliendo...")
            break
        else:
            print("Opción no válida. Intenta de nuevo.")

if __name__ == "__main__":
    main()
