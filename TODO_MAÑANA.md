# 📋 TODO para Mañana - k8s-cli Version Flag

## 🎯 Tarea Principal
**Implementar `k8s-cli --version` para mostrar versión del CLI**

### 📝 Problema Actual
- `k8s-cli version` muestra información del cluster Kubernetes
- No hay forma de ver la versión del CLI tool en sí
- Confusión entre versión del CLI vs versión del cluster

### 🎯 Objetivo
Agregar soporte para `k8s-cli --version` que muestre:
- Versión del CLI
- Commit Git
- Fecha de build
- Información del Go runtime

---

## 🛠️ Implementación Requerida

### 1. Modificar `main.go`
```go
// Agregar variables de build
var (
    Version   = "dev"
    GitCommit = "unknown"
    BuildTime = "unknown"
    GoVersion = runtime.Version()
)

// Agregar flag --version en rootCmd
rootCmd.Flags().BoolP("version", "v", false, "Show CLI version")
```

### 2. Agregar lógica de version en `main.go`
```go
// En la función Execute o main
if version, _ := cmd.Flags().GetBool("version"); version {
    fmt.Printf("k8s-cli version %s\n", Version)
    fmt.Printf("Git commit: %s\n", GitCommit)
    fmt.Printf("Built: %s\n", BuildTime)
    fmt.Printf("Go version: %s\n", GoVersion)
    fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
    return
}
```

### 3. Actualizar Build Flags en Makefile.dev
Verificar que las build flags estén correctas:
```makefile
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"
```

---

## 📁 Archivos a Modificar

### Archivos Principales
- [ ] `main.go` - Agregar flag --version y lógica
- [ ] `cmd/root.go` - Si la lógica está ahí en lugar de main.go

### Archivos de Testing
- [ ] Crear test para verificar que --version funciona
- [ ] Actualizar tests existentes si es necesario

### Documentación
- [ ] Actualizar `README.md` con ejemplo de --version
- [ ] Actualizar `docs/EXAMPLES.md` con nuevo flag
- [ ] Actualizar help text del CLI

---

## 🧪 Plan de Testing

### Tests Unitarios
```bash
# Test que --version muestra información correcta
go test -v ./cmd/ -run TestVersionFlag

# Test que no interfiere con comando version existente
./bin/k8s-cli version  # Debe mostrar info de K8s
./bin/k8s-cli --version  # Debe mostrar info del CLI
```

### Tests de Integración
```bash
# Verificar que build flags funcionan
make -f Makefile.dev build
./bin/k8s-cli --version

# Verificar que no rompe funcionalidad existente
./bin/k8s-cli version  # Comando K8s original
./bin/k8s-cli --help   # Help general
```

---

## 📊 Resultado Esperado

### Comportamiento Actual
```bash
$ k8s-cli version
# Muestra versión de Kubernetes cluster

$ k8s-cli --version
# Error: flag no existe
```

### Comportamiento Deseado
```bash
$ k8s-cli version
# Muestra versión de Kubernetes cluster (sin cambios)

$ k8s-cli --version
k8s-cli version 2.0.0
Git commit: a1b2c3d
Built: 2024-01-16T10:30:00
Go version: go1.24.5
OS/Arch: darwin/arm64
```

---

## 🔧 Comandos de Desarrollo

### Setup y Build
```bash
# Navegar al proyecto
cd /Users/jaime.henao/arheanja/lab-go-cli

# Desarrollo con hot reload
make -f Makefile.dev watch

# Build y test
make -f Makefile.dev dev-cycle
```

### Testing
```bash
# Test unitarios
make -f Makefile.dev test

# Test manual
./bin/k8s-cli --version
./bin/k8s-cli -v
./bin/k8s-cli version  # Verificar que no cambió
```

---

## 📚 Referencias

### Documentación a Consultar
- [Cobra CLI Documentation](https://cobra.dev/) - Para flags globales
- [Go build flags](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies) - Para ldflags
- Proyecto actual en `cmd/` para ver estructura

### Ejemplos de Implementación
- Revisar otros CLIs como `kubectl --version`
- Ver implementación actual en `cmd/version.go`

---

## ✅ Criterios de Aceptación

### Funcionalidad
- [ ] `k8s-cli --version` muestra información del CLI
- [ ] `k8s-cli -v` funciona como shorthand
- [ ] `k8s-cli version` sigue funcionando para Kubernetes
- [ ] Información incluye: version, commit, build time, go version, OS/arch

### Calidad
- [ ] Tests unitarios cubren nueva funcionalidad
- [ ] No hay regresiones en funcionalidad existente
- [ ] Help text actualizado
- [ ] Documentación actualizada

### Build
- [ ] Build flags se aplican correctamente
- [ ] Makefile.dev funciona con nuevos cambios
- [ ] Release build incluye información correcta

---

## ⏰ Estimación de Tiempo
- **Implementación**: 30-45 minutos
- **Testing**: 15-20 minutos  
- **Documentación**: 10-15 minutos
- **Total**: ~1 hora

---

## 🚀 Para Después de Completar

### Verificación Final
```bash
# Build final y test completo
make -f Makefile.dev dev-cycle

# Instalar y probar
make -f Makefile.dev install-user
k8s-cli --version
k8s-cli version

# Verificar que todo funciona
k8s-cli all
```

### Commit
```bash
git add .
git commit -m "feat: add --version flag for CLI version info

- Add --version/-v flag to show CLI version, commit, build info
- Preserve existing 'version' command for Kubernetes cluster info
- Update documentation and examples
- Add unit tests for version flag"
```

---

**📅 Fecha:** 15 de Enero, 2024  
**⏰ Prioridad:** Alta  
**🎯 Objetivo:** Mejorar UX del CLI con información de versión clara