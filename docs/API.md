# 📡 k8s-cli API Documentation

## 🎯 Overview

Esta documentación describe la API interna de k8s-cli, las estructuras de datos, y cómo extender la funcionalidad.

## 🏗️ Core Data Structures

### 1. Metrics Data Structures

#### NodeMetrics
```go
type NodeMetrics struct {
    Name               string  // Nombre del nodo
    CPUUsage          string  // Uso actual de CPU
    CPUUsagePercent   float64 // Porcentaje de uso de CPU
    MemoryUsage       string  // Uso actual de memoria
    MemoryUsagePercent float64 // Porcentaje de uso de memoria
    CPUCapacity       string  // Capacidad total de CPU
    MemoryCapacity    string  // Capacidad total de memoria
    Status            string  // Estado del nodo
}
```

#### PodMetrics
```go
type PodMetrics struct {
    Name           string // Nombre del pod
    Namespace      string // Namespace
    CPUUsage       string // Uso de CPU
    MemoryUsage    string // Uso de memoria
    CPURequests    string // CPU solicitada
    MemoryRequests string // Memoria solicitada
    CPULimits      string // Límites de CPU
    MemoryLimits   string // Límites de memoria
    Node           string // Nodo donde corre
    RestartCount   int32  // Número de reinicios
}
```

#### ClusterMetrics
```go
type ClusterMetrics struct {
    TotalCPUUsage       string  // Uso total de CPU
    TotalMemoryUsage    string  // Uso total de memoria
    TotalCPUCapacity    string  // Capacidad total de CPU
    TotalMemoryCapacity string  // Capacidad total de memoria
    CPUUsagePercent     float64 // Porcentaje de uso de CPU
    MemoryUsagePercent  float64 // Porcentaje de uso de memoria
    NodesCount          int     // Número de nodos
    PodsCount           int     // Número de pods
    NamespacesCount     int     // Número de namespaces
}
```

### 2. Cost Analysis Data Structures

#### CostAnalysis
```go
type CostAnalysis struct {
    TotalMonthlyCost       float64                 // Costo mensual total
    NodeCosts              []NodeCost              // Costos por nodo
    NamespaceCosts         []NamespaceCost         // Costos por namespace
    UnderutilizedResources []UnderutilizedResource // Recursos subutilizados
    CostOptimizations      []CostOptimization      // Recomendaciones de optimización
}
```

#### NodeCost
```go
type NodeCost struct {
    Name           string  // Nombre del nodo
    Type           string  // Tipo de instancia
    MonthlyCost    float64 // Costo mensual estimado
    CPUCapacity    string  // Capacidad de CPU
    MemoryCapacity string  // Capacidad de memoria
    CPUUtilization float64 // Utilización de CPU
    MemUtilization float64 // Utilización de memoria
    Efficiency     string  // Eficiencia general
}
```

### 3. Workload Analysis Data Structures

#### WorkloadAnalysis
```go
type WorkloadAnalysis struct {
    DeploymentAnalysis  []DeploymentHealth  // Análisis de deployments
    StatefulSetAnalysis []StatefulSetHealth // Análisis de statefulsets
    DaemonSetAnalysis   []DaemonSetHealth   // Análisis de daemonsets
    PodAnalysis         []PodHealth         // Análisis de pods
    WorkloadSummary     WorkloadSummary     // Resumen general
}
```

#### DeploymentHealth
```go
type DeploymentHealth struct {
    Name                string   // Nombre del deployment
    Namespace           string   // Namespace
    Replicas            int32    // Réplicas configuradas
    ReadyReplicas       int32    // Réplicas listas
    AvailableReplicas   int32    // Réplicas disponibles
    UnavailableReplicas int32    // Réplicas no disponibles
    Status              string   // Estado general
    Age                 string   // Edad del deployment
    RestartRate         float64  // Tasa de reinicios
    ResourceEfficiency  string   // Eficiencia de recursos
    HealthScore         int      // Puntuación de salud (0-100)
    Issues              []string // Problemas detectados
    Recommendations     []string // Recomendaciones
}
```

### 4. Event and Log Data Structures

#### ClusterEvent
```go
type ClusterEvent struct {
    Type       string    // Tipo de evento
    Reason     string    // Razón del evento
    Message    string    // Mensaje del evento
    Object     string    // Objeto afectado
    Namespace  string    // Namespace
    FirstTime  time.Time // Primera ocurrencia
    LastTime   time.Time // Última ocurrencia
    Count      int32     // Número de ocurrencias
    Severity   string    // Severidad (Critical, Warning, Info)
    Component  string    // Componente que generó el evento
}
```

#### ErrorPattern
```go
type ErrorPattern struct {
    Pattern        string    // Patrón de error
    Count          int       // Número de ocurrencias
    LastSeen       time.Time // Última vez visto
    Severity       string    // Severidad
    Description    string    // Descripción del problema
    Recommendation string    // Recomendación para resolver
}
```

## 🔧 Core API Functions

### 1. Kubernetes Client API

#### Connection Management
```go
// Crear cliente de Kubernetes
func NewClient(kubeconfig string) (*Client, error)

// Obtener información del cluster
func (c *Client) GetClusterVersion() (*ClusterVersion, error)
```

#### Metrics API
```go
// Obtener métricas de nodos en tiempo real
func (c *Client) GetRealTimeNodeMetrics() ([]NodeMetrics, error)

// Obtener métricas de pods en tiempo real
func (c *Client) GetRealTimePodMetrics(namespace string) ([]PodMetrics, error)

// Obtener métricas generales del cluster
func (c *Client) GetClusterMetrics() (*ClusterMetrics, error)

// Obtener análisis de utilización de recursos
func (c *Client) GetResourceUtilization() ([]ResourceUtilization, error)
```

#### Cost Analysis API
```go
// Obtener análisis completo de costos
func (c *Client) GetCostAnalysis() (*CostAnalysis, error)

// Encontrar recursos subutilizados
func (c *Client) findUnderutilizedResources() ([]UnderutilizedResource, error)

// Generar recomendaciones de optimización de costos
func (c *Client) generateCostOptimizations(...) []CostOptimization
```

#### Workload Analysis API
```go
// Obtener análisis completo de workloads
func (c *Client) GetWorkloadAnalysis(namespace string) (*WorkloadAnalysis, error)

// Analizar salud de deployments
func (c *Client) analyzeDeployments(namespace string) ([]DeploymentHealth, error)

// Analizar salud de pods
func (c *Client) analyzePods(namespace string) ([]PodHealth, error)
```

#### Events and Logs API
```go
// Obtener eventos del cluster
func (c *Client) GetClusterEvents(namespace string, hours int) ([]ClusterEvent, error)

// Obtener análisis de logs
func (c *Client) GetLogAnalysis(namespace string, hours int) (*LogAnalysis, error)

// Obtener análisis de logs de pods
func (c *Client) GetPodLogsAnalysis(namespace string) ([]PodLogSummary, error)
```

### 2. Export API

#### Export Manager
```go
// Crear exportador
func NewExporter(outputDir string) *Exporter

// Exportar a JSON
func (e *Exporter) ExportToJSON(data *ExportData, filename string) error

// Exportar métricas de nodos a CSV
func (e *Exporter) ExportNodeMetricsToCSV(metrics []NodeMetrics, filename string) error

// Exportar análisis de costos a CSV
func (e *Exporter) ExportCostAnalysisToCSV(analysis *CostAnalysis, filename string) error

// Exportar métricas de Prometheus
func (e *Exporter) ExportPrometheusMetrics(data *ExportData, filename string) error
```

#### Export Data Structure
```go
type ExportData struct {
    Timestamp      time.Time                    // Timestamp del export
    ClusterMetrics *ClusterMetrics             // Métricas del cluster
    NodeMetrics    []NodeMetrics               // Métricas de nodos
    PodMetrics     []PodMetrics                // Métricas de pods
    CostAnalysis   *CostAnalysis               // Análisis de costos
    LogAnalysis    *LogAnalysis                // Análisis de logs
    Utilizations   []ResourceUtilization       // Utilización de recursos
    Events         []ClusterEvent              // Eventos del cluster
}
```

## 🎯 Usage Examples

### 1. Basic Metrics Retrieval
```go
client, err := kubernetes.NewClient("")
if err != nil {
    return err
}

// Obtener métricas de nodos
nodeMetrics, err := client.GetRealTimeNodeMetrics()
if err != nil {
    return err
}

for _, node := range nodeMetrics {
    fmt.Printf("Node: %s, CPU: %.1f%%, Memory: %.1f%%\n", 
        node.Name, node.CPUUsagePercent, node.MemoryUsagePercent)
}
```

### 2. Cost Analysis
```go
client, err := kubernetes.NewClient("")
if err != nil {
    return err
}

// Obtener análisis de costos
analysis, err := client.GetCostAnalysis()
if err != nil {
    return err
}

fmt.Printf("Total Monthly Cost: $%.2f\n", analysis.TotalMonthlyCost)
fmt.Printf("Underutilized Resources: %d\n", len(analysis.UnderutilizedResources))
```

### 3. Export Data
```go
exporter := export.NewExporter("./exports")
data := &export.ExportData{
    Timestamp: time.Now(),
    // ... populate data
}

// Exportar a JSON
err := exporter.ExportToJSON(data, "cluster-analysis")
if err != nil {
    return err
}

// Exportar a CSV
err = exporter.ExportNodeMetricsToCSV(data.NodeMetrics, "node-metrics")
if err != nil {
    return err
}
```

## 🔧 Extending the API

### 1. Adding New Metrics
```go
// 1. Define new data structure
type CustomMetric struct {
    Name  string
    Value float64
}

// 2. Add method to Client
func (c *Client) GetCustomMetrics() ([]CustomMetric, error) {
    // Implementation
}

// 3. Add to export data
type ExportData struct {
    // ... existing fields
    CustomMetrics []CustomMetric `json:"custom_metrics,omitempty"`
}
```

### 2. Adding New Analysis Types
```go
// 1. Define analysis structure
type SecurityAnalysis struct {
    Vulnerabilities []Vulnerability
    PolicyViolations []PolicyViolation
    RiskScore       int
}

// 2. Implement analysis function
func (c *Client) GetSecurityAnalysis() (*SecurityAnalysis, error) {
    // Implementation
}

// 3. Add command in cmd/
var securityCmd = &cobra.Command{
    Use:   "security",
    Short: "Analyze cluster security",
    RunE:  runSecurityCommand,
}
```

### 3. Adding New Export Formats
```go
// Add to exporter
func (e *Exporter) ExportToXML(data *ExportData, filename string) error {
    // XML export implementation
}

func (e *Exporter) ExportToYAML(data *ExportData, filename string) error {
    // YAML export implementation
}
```

## 🔍 Error Handling

### Standard Error Types
```go
// Client errors
var (
    ErrClientNotInitialized = errors.New("kubernetes client not initialized")
    ErrInvalidNamespace     = errors.New("invalid namespace")
    ErrMetricsNotAvailable  = errors.New("metrics server not available")
)

// Export errors
var (
    ErrInvalidFormat     = errors.New("invalid export format")
    ErrExportPathInvalid = errors.New("invalid export path")
)
```

### Error Handling Pattern
```go
func (c *Client) SomeOperation() error {
    if c.Clientset == nil {
        return ErrClientNotInitialized
    }
    
    result, err := c.Clientset.CoreV1().Pods("").List(...)
    if err != nil {
        return fmt.Errorf("failed to list pods: %w", err)
    }
    
    return nil
}
```

## 📊 Performance Considerations

### API Rate Limiting
- Use ListOptions with pagination
- Implement exponential backoff
- Cache frequently accessed data

### Memory Management
- Stream large datasets
- Use context for cancellation
- Implement garbage collection hints

### Concurrency
- Use goroutines for parallel analysis
- Implement proper synchronization
- Handle context cancellation