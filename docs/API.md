# 📡 k8s-cli API Documentation

## 🎯 Overview

This documentation describes the internal API of k8s-cli, data structures, interfaces, and how to extend functionality. The API is designed for extensibility, type safety, and performance.

## 🏗️ Core Interfaces

### 🔧 **Primary Client Interface**

```go
// KubernetesClient provides the main interface for cluster operations
type KubernetesClient interface {
    // Cluster Information
    GetClusterInfo() (*ClusterInfo, error)
    GetClusterVersion() (*version.Info, error)
    
    // Resource Operations
    GetNodes() (*v1.NodeList, error)
    GetPods(namespace string) (*v1.PodList, error)
    GetServices(namespace string) (*v1.ServiceList, error)
    GetDeployments(namespace string) (*appsv1.DeploymentList, error)
    
    // Metrics Operations
    GetNodeMetrics() ([]NodeMetrics, error)
    GetPodMetrics(namespace string) ([]PodMetrics, error)
    GetResourceUtilization() ([]ResourceUtilization, error)
    GetRealTimeNodeMetrics() ([]NodeMetrics, error)
    GetRealTimePodMetrics(namespace string) ([]PodMetrics, error)
    
    // Analysis Operations
    GetCostAnalysis() (*CostAnalysis, error)
    GetWorkloadHealth() ([]WorkloadHealth, error)
    GetLogAnalysis(namespace string, hours int) (*LogAnalysis, error)
    GetClusterEvents(namespace string, hours int) ([]ClusterEvent, error)
    
    // Component Detection
    GetComponents() ([]Component, error)
    GetHelmReleases() ([]HelmRelease, error)
}
```

### 📊 **Export Interface**

```go
// Exporter handles multi-format data export
type Exporter interface {
    // Core export methods
    ExportToJSON(data interface{}, filename string) error
    ExportToCSV(data interface{}, filename string) error
    ExportToPrometheus(data interface{}, filename string) error
    
    // Specialized export methods
    ExportNodeMetricsToCSV(metrics []NodeMetrics, filename string) error
    ExportPodMetricsToCSV(metrics []PodMetrics, filename string) error
    ExportCostAnalysisToCSV(analysis *CostAnalysis, filename string) error
    ExportEventsToCSV(events []ClusterEvent, filename string) error
    ExportUtilizationToCSV(utilization []ResourceUtilization, filename string) error
    ExportPrometheusMetrics(data *ExportData, filename string) error
    
    // Utility methods
    GetExportPath(filename string) string
    SetOutputDirectory(dir string)
}
```

### 🎯 **Analyzer Interface**

```go
// Analyzer provides extensible analysis capabilities
type Analyzer interface {
    Name() string
    Version() string
    Analyze(cluster *ClusterData) (*AnalysisResult, error)
    GetRecommendations(result *AnalysisResult) ([]Recommendation, error)
    SupportsClusterVersion(version string) bool
}

// Specific analyzer implementations
type CostAnalyzer interface {
    Analyzer
    CalculateNodeCosts(nodes []v1.Node) ([]NodeCost, error)
    EstimateResourceCosts(utilization []ResourceUtilization) (*CostEstimate, error)
    FindOptimizations(analysis *CostAnalysis) ([]CostOptimization, error)
}

type HealthAnalyzer interface {
    Analyzer
    ScoreWorkloadHealth(workloads []Workload) ([]WorkloadHealth, error)
    DetectIssues(workload Workload) ([]HealthIssue, error)
    GenerateHealthRecommendations(health WorkloadHealth) ([]string, error)
}
```

## 🗂️ Core Data Structures

### 📊 **Metrics Data Structures**

#### **NodeMetrics**
```go
type NodeMetrics struct {
    Name               string    `json:"name" csv:"name"`
    CPUUsage          string    `json:"cpu_usage" csv:"cpu_usage"`
    CPUUsagePercent   float64   `json:"cpu_usage_percent" csv:"cpu_usage_percent"`
    MemoryUsage       string    `json:"memory_usage" csv:"memory_usage"`
    MemoryUsagePercent float64  `json:"memory_usage_percent" csv:"memory_usage_percent"`
    CPUCapacity       string    `json:"cpu_capacity" csv:"cpu_capacity"`
    MemoryCapacity    string    `json:"memory_capacity" csv:"memory_capacity"`
    Status            string    `json:"status" csv:"status"`
    Region            string    `json:"region,omitempty" csv:"region"`
    InstanceType      string    `json:"instance_type,omitempty" csv:"instance_type"`
    Age               string    `json:"age" csv:"age"`
    CreatedAt         time.Time `json:"created_at" csv:"created_at"`
}
```

#### **PodMetrics**
```go
type PodMetrics struct {
    Name           string    `json:"name" csv:"name"`
    Namespace      string    `json:"namespace" csv:"namespace"`
    CPUUsage       string    `json:"cpu_usage" csv:"cpu_usage"`
    MemoryUsage    string    `json:"memory_usage" csv:"memory_usage"`
    CPURequest     string    `json:"cpu_request" csv:"cpu_request"`
    MemoryRequest  string    `json:"memory_request" csv:"memory_request"`
    CPULimit       string    `json:"cpu_limit" csv:"cpu_limit"`
    MemoryLimit    string    `json:"memory_limit" csv:"memory_limit"`
    Node           string    `json:"node" csv:"node"`
    Status         string    `json:"status" csv:"status"`
    RestartCount   int32     `json:"restart_count" csv:"restart_count"`
    Age            string    `json:"age" csv:"age"`
    CreatedAt      time.Time `json:"created_at" csv:"created_at"`
}
```

#### **ResourceUtilization**
```go
type ResourceUtilization struct {
    ResourceType    string  `json:"resource_type" csv:"resource_type"`
    Total          string  `json:"total" csv:"total"`
    Used           string  `json:"used" csv:"used"`
    Available      string  `json:"available" csv:"available"`
    UsagePercent   float64 `json:"usage_percent" csv:"usage_percent"`
    Recommendation string  `json:"recommendation" csv:"recommendation"`
    Priority       string  `json:"priority" csv:"priority"`
}
```

### 💰 **Cost Analysis Data Structures**

#### **CostAnalysis**
```go
type CostAnalysis struct {
    Timestamp         time.Time         `json:"timestamp"`
    ClusterName       string           `json:"cluster_name"`
    TotalMonthlyCost  float64          `json:"total_monthly_cost"`
    TotalDailyCost    float64          `json:"total_daily_cost"`
    NodeCosts         []NodeCost       `json:"node_costs"`
    NamespaceCosts    []NamespaceCost  `json:"namespace_costs"`
    Optimizations     []CostOptimization `json:"optimizations"`
    PotentialSavings  float64          `json:"potential_savings"`
    Currency          string           `json:"currency"`
    AnalysisVersion   string           `json:"analysis_version"`
}
```

#### **NodeCost**
```go
type NodeCost struct {
    Name             string  `json:"name" csv:"name"`
    InstanceType     string  `json:"instance_type" csv:"instance_type"`
    Region           string  `json:"region" csv:"region"`
    HourlyCost       float64 `json:"hourly_cost" csv:"hourly_cost"`
    MonthlyCost      float64 `json:"monthly_cost" csv:"monthly_cost"`
    CPUCost          float64 `json:"cpu_cost" csv:"cpu_cost"`
    MemoryCost       float64 `json:"memory_cost" csv:"memory_cost"`
    StorageCost      float64 `json:"storage_cost" csv:"storage_cost"`
    NetworkCost      float64 `json:"network_cost" csv:"network_cost"`
    Utilization      float64 `json:"utilization" csv:"utilization"`
    Efficiency       string  `json:"efficiency" csv:"efficiency"`
    Recommendations  []string `json:"recommendations" csv:"recommendations"`
}
```

#### **CostOptimization**
```go
type CostOptimization struct {
    Type            string  `json:"type" csv:"type"`
    Resource        string  `json:"resource" csv:"resource"`
    CurrentCost     float64 `json:"current_cost" csv:"current_cost"`
    OptimizedCost   float64 `json:"optimized_cost" csv:"optimized_cost"`
    PotentialSaving float64 `json:"potential_saving" csv:"potential_saving"`
    Priority        string  `json:"priority" csv:"priority"`
    Description     string  `json:"description" csv:"description"`
    Action          string  `json:"action" csv:"action"`
    Confidence      float64 `json:"confidence" csv:"confidence"`
}
```

### 🏥 **Health Analysis Data Structures**

#### **WorkloadHealth**
```go
type WorkloadHealth struct {
    Name            string        `json:"name" csv:"name"`
    Namespace       string        `json:"namespace" csv:"namespace"`
    Type            string        `json:"type" csv:"type"`
    HealthScore     float64       `json:"health_score" csv:"health_score"`
    Status          string        `json:"status" csv:"status"`
    Replicas        int32         `json:"replicas" csv:"replicas"`
    ReadyReplicas   int32         `json:"ready_replicas" csv:"ready_replicas"`
    Issues          []HealthIssue `json:"issues"`
    Recommendations []string      `json:"recommendations"`
    LastUpdated     time.Time     `json:"last_updated" csv:"last_updated"`
}
```

#### **HealthIssue**
```go
type HealthIssue struct {
    Type        string    `json:"type" csv:"type"`
    Severity    string    `json:"severity" csv:"severity"`
    Description string    `json:"description" csv:"description"`
    Resource    string    `json:"resource" csv:"resource"`
    Remediation string    `json:"remediation" csv:"remediation"`
    DetectedAt  time.Time `json:"detected_at" csv:"detected_at"`
}
```

### 📝 **Event and Log Data Structures**

#### **ClusterEvent**
```go
type ClusterEvent struct {
    Timestamp    time.Time `json:"timestamp" csv:"timestamp"`
    Type         string    `json:"type" csv:"type"`
    Reason       string    `json:"reason" csv:"reason"`
    Message      string    `json:"message" csv:"message"`
    Source       string    `json:"source" csv:"source"`
    Object       string    `json:"object" csv:"object"`
    Namespace    string    `json:"namespace" csv:"namespace"`
    Severity     string    `json:"severity" csv:"severity"`
    Count        int32     `json:"count" csv:"count"`
    FirstSeen    time.Time `json:"first_seen" csv:"first_seen"`
    LastSeen     time.Time `json:"last_seen" csv:"last_seen"`
}
```

#### **LogAnalysis**
```go
type LogAnalysis struct {
    Timestamp      time.Time         `json:"timestamp"`
    TimeRange      string           `json:"time_range"`
    TotalEvents    int              `json:"total_events"`
    CriticalEvents int              `json:"critical_events"`
    WarningEvents  int              `json:"warning_events"`
    InfoEvents     int              `json:"info_events"`
    Patterns       []LogPattern     `json:"patterns"`
    Anomalies      []LogAnomaly     `json:"anomalies"`
    TopErrors      []ErrorSummary   `json:"top_errors"`
}
```

### 🔧 **Component Detection Data Structures**

#### **Component**
```go
type Component struct {
    Name        string            `json:"name" csv:"name"`
    Namespace   string            `json:"namespace" csv:"namespace"`
    Type        string            `json:"type" csv:"type"`
    Version     string            `json:"version" csv:"version"`
    Status      string            `json:"status" csv:"status"`
    Source      string            `json:"source" csv:"source"`
    Labels      map[string]string `json:"labels,omitempty"`
    Annotations map[string]string `json:"annotations,omitempty"`
    Age         string            `json:"age" csv:"age"`
    CreatedAt   time.Time         `json:"created_at" csv:"created_at"`
}
```

#### **HelmRelease**
```go
type HelmRelease struct {
    Name         string            `json:"name" csv:"name"`
    Namespace    string            `json:"namespace" csv:"namespace"`
    Chart        string            `json:"chart" csv:"chart"`
    Version      string            `json:"version" csv:"version"`
    AppVersion   string            `json:"app_version" csv:"app_version"`
    Status       string            `json:"status" csv:"status"`
    Revision     int               `json:"revision" csv:"revision"`
    Updated      time.Time         `json:"updated" csv:"updated"`
    Values       map[string]interface{} `json:"values,omitempty"`
}
```

### 📤 **Export Data Structures**

#### **ExportData**
```go
type ExportData struct {
    Timestamp        time.Time              `json:"timestamp"`
    ClusterName      string                 `json:"cluster_name"`
    Version          string                 `json:"version"`
    ClusterMetrics   *ClusterInfo           `json:"cluster_metrics,omitempty"`
    NodeMetrics      []NodeMetrics          `json:"node_metrics,omitempty"`
    PodMetrics       []PodMetrics           `json:"pod_metrics,omitempty"`
    CostAnalysis     *CostAnalysis          `json:"cost_analysis,omitempty"`
    WorkloadHealth   []WorkloadHealth       `json:"workload_health,omitempty"`
    LogAnalysis      *LogAnalysis           `json:"log_analysis,omitempty"`
    Events           []ClusterEvent         `json:"events,omitempty"`
    Utilizations     []ResourceUtilization  `json:"utilizations,omitempty"`
    Components       []Component            `json:"components,omitempty"`
    Recommendations  []Recommendation       `json:"recommendations,omitempty"`
    ExportMetadata   ExportMetadata         `json:"export_metadata"`
}
```

#### **ExportMetadata**
```go
type ExportMetadata struct {
    ExportedAt     time.Time         `json:"exported_at"`
    ExportVersion  string            `json:"export_version"`
    CLIVersion     string            `json:"cli_version"`
    Format         string            `json:"format"`
    Sections       []string          `json:"sections"`
    Options        map[string]interface{} `json:"options,omitempty"`
    FileSize       int64             `json:"file_size,omitempty"`
    Checksum       string            `json:"checksum,omitempty"`
}
```

## 🔌 Extension APIs

### 📊 **Custom Analyzer Registration**

```go
// AnalyzerRegistry manages custom analyzers
type AnalyzerRegistry interface {
    Register(analyzer Analyzer) error
    Unregister(name string) error
    Get(name string) (Analyzer, error)
    List() []string
    Execute(name string, data *ClusterData) (*AnalysisResult, error)
}

// Example custom analyzer implementation
type CustomSecurityAnalyzer struct {
    name    string
    version string
}

func (a *CustomSecurityAnalyzer) Name() string { return a.name }
func (a *CustomSecurityAnalyzer) Version() string { return a.version }

func (a *CustomSecurityAnalyzer) Analyze(cluster *ClusterData) (*AnalysisResult, error) {
    // Custom security analysis logic
    return &AnalysisResult{
        Type: "security",
        Score: calculateSecurityScore(cluster),
        Issues: findSecurityIssues(cluster),
        Recommendations: generateSecurityRecommendations(cluster),
    }, nil
}
```

### 📤 **Custom Export Format Registration**

```go
// ExportFormatRegistry manages custom export formats
type ExportFormatRegistry interface {
    RegisterFormat(format string, handler ExportHandler) error
    UnregisterFormat(format string) error
    GetHandler(format string) (ExportHandler, error)
    SupportedFormats() []string
}

// ExportHandler interface for custom formats
type ExportHandler interface {
    Export(data interface{}, writer io.Writer, options ExportOptions) error
    ValidateOptions(options ExportOptions) error
    GetMimeType() string
    GetFileExtension() string
}

// Example custom export format
type XMLExportHandler struct{}

func (h *XMLExportHandler) Export(data interface{}, writer io.Writer, options ExportOptions) error {
    // XML export implementation
    encoder := xml.NewEncoder(writer)
    return encoder.Encode(data)
}

func (h *XMLExportHandler) GetMimeType() string { return "application/xml" }
func (h *XMLExportHandler) GetFileExtension() string { return ".xml" }
```

## 🔧 Configuration APIs

### ⚙️ **Configuration Interface**

```go
// Configuration management
type Config interface {
    // Basic configuration
    GetKubeconfig() string
    GetNamespace() string
    GetOutputFormat() string
    
    // Advanced configuration
    GetCacheTimeout() time.Duration
    GetMaxConcurrency() int
    GetExportDirectory() string
    
    // Feature flags
    IsMetricsEnabled() bool
    IsCostAnalysisEnabled() bool
    IsExportEnabled() bool
    
    // Custom configuration
    GetCustomConfig(key string) (interface{}, error)
    SetCustomConfig(key string, value interface{}) error
}

// Configuration options
type ConfigOptions struct {
    Kubeconfig      string        `yaml:"kubeconfig"`
    Namespace       string        `yaml:"namespace"`
    OutputFormat    string        `yaml:"output_format"`
    CacheTimeout    time.Duration `yaml:"cache_timeout"`
    MaxConcurrency  int           `yaml:"max_concurrency"`
    ExportDirectory string        `yaml:"export_directory"`
    
    Features struct {
        Metrics     bool `yaml:"metrics"`
        CostAnalysis bool `yaml:"cost_analysis"`
        Export      bool `yaml:"export"`
    } `yaml:"features"`
    
    Custom map[string]interface{} `yaml:"custom,omitempty"`
}
```

## 🧪 Testing APIs

### 🎯 **Test Utilities**

```go
// MockClient for testing
type MockClient struct {
    // Mock implementations for all KubernetesClient methods
    ClusterInfo      *ClusterInfo
    Nodes            []v1.Node
    Pods             []v1.Pod
    NodeMetrics      []NodeMetrics
    PodMetrics       []PodMetrics
    CostAnalysis     *CostAnalysis
    WorkloadHealth   []WorkloadHealth
    Components       []Component
    Events           []ClusterEvent
}

// Test data generators
func GenerateTestClusterData(size ClusterSize) *ClusterData
func GenerateTestNodeMetrics(count int) []NodeMetrics
func GenerateTestPodMetrics(count int) []PodMetrics
func GenerateTestCostAnalysis() *CostAnalysis

// Cluster size options for testing
type ClusterSize int

const (
    SmallCluster  ClusterSize = iota // 3 nodes, 10 pods
    MediumCluster                    // 10 nodes, 100 pods
    LargeCluster                     // 50 nodes, 1000 pods
    XLCluster                        // 200 nodes, 5000 pods
)
```

### 🔍 **Validation APIs**

```go
// DataValidator provides validation utilities
type DataValidator interface {
    ValidateClusterData(data *ClusterData) error
    ValidateMetrics(metrics []NodeMetrics) error
    ValidateExportData(data *ExportData) error
    ValidateConfiguration(config *ConfigOptions) error
}

// Example validation implementation
func ValidateNodeMetrics(metrics []NodeMetrics) error {
    for _, metric := range metrics {
        if metric.Name == "" {
            return fmt.Errorf("node name cannot be empty")
        }
        if metric.CPUUsagePercent < 0 || metric.CPUUsagePercent > 100 {
            return fmt.Errorf("invalid CPU usage percentage: %f", metric.CPUUsagePercent)
        }
        if metric.MemoryUsagePercent < 0 || metric.MemoryUsagePercent > 100 {
            return fmt.Errorf("invalid memory usage percentage: %f", metric.MemoryUsagePercent)
        }
    }
    return nil
}
```

## 📊 Performance APIs

### ⚡ **Performance Monitoring**

```go
// PerformanceMetrics tracks API performance
type PerformanceMetrics interface {
    StartTimer(operation string) Timer
    RecordDuration(operation string, duration time.Duration)
    RecordCount(metric string, count int64)
    GetMetrics() map[string]interface{}
    Reset()
}

// Timer for measuring operation duration
type Timer interface {
    Stop() time.Duration
    Duration() time.Duration
}

// Performance configuration
type PerformanceConfig struct {
    EnableMetrics     bool          `yaml:"enable_metrics"`
    MetricsInterval   time.Duration `yaml:"metrics_interval"`
    MaxConcurrency    int           `yaml:"max_concurrency"`
    RequestTimeout    time.Duration `yaml:"request_timeout"`
    CacheSize         int           `yaml:"cache_size"`
    CacheTTL          time.Duration `yaml:"cache_ttl"`
}
```

## 🔐 Security APIs

### 🛡️ **Security Interface**

```go
// SecurityManager handles security operations
type SecurityManager interface {
    ValidatePermissions(config *rest.Config) error
    SanitizeOutput(data interface{}) interface{}
    CheckCredentials(kubeconfig string) error
    GenerateAuditLog(operation string, user string, timestamp time.Time) error
}

// Security configuration
type SecurityConfig struct {
    EnableAuditLogging bool     `yaml:"enable_audit_logging"`
    AuditLogPath      string   `yaml:"audit_log_path"`
    SanitizeSecrets   bool     `yaml:"sanitize_secrets"`
    RequiredRBAC      []string `yaml:"required_rbac"`
}

// Audit log entry
type AuditEntry struct {
    Timestamp time.Time `json:"timestamp"`
    User      string    `json:"user"`
    Operation string    `json:"operation"`
    Resource  string    `json:"resource"`
    Success   bool      `json:"success"`
    Error     string    `json:"error,omitempty"`
}
```

## 📈 API Versioning

### 🔄 **Version Management**

```go
// APIVersion defines the API version
const (
    APIVersion     = "v1"
    MinAPIVersion  = "v1"
    MaxAPIVersion  = "v1"
)

// VersionInfo provides version information
type VersionInfo struct {
    APIVersion    string `json:"api_version"`
    CLIVersion    string `json:"cli_version"`
    BuildDate     string `json:"build_date"`
    GitCommit     string `json:"git_commit"`
    GoVersion     string `json:"go_version"`
    Platform      string `json:"platform"`
}

// Compatibility checker
type CompatibilityChecker interface {
    IsCompatible(clientVersion, serverVersion string) bool
    GetMinimumVersion() string
    GetRecommendedVersion() string
}
```

## 🚀 Usage Examples

### 📊 **Basic API Usage**

```go
// Initialize client
client, err := kubernetes.NewClient("")
if err != nil {
    log.Fatal(err)
}

// Get cluster metrics
metrics, err := client.GetNodeMetrics()
if err != nil {
    log.Fatal(err)
}

// Perform cost analysis
costAnalysis, err := client.GetCostAnalysis()
if err != nil {
    log.Fatal(err)
}

// Export data
exporter := export.NewExporter("./exports")
exportData := &export.ExportData{
    Timestamp:    time.Now(),
    NodeMetrics:  metrics,
    CostAnalysis: costAnalysis,
}

err = exporter.ExportToJSON(exportData, "cluster-analysis")
if err != nil {
    log.Fatal(err)
}
```

### 🔌 **Custom Analyzer Example**

```go
// Define custom analyzer
type CustomPerformanceAnalyzer struct{}

func (a *CustomPerformanceAnalyzer) Name() string { return "performance" }
func (a *CustomPerformanceAnalyzer) Version() string { return "1.0.0" }

func (a *CustomPerformanceAnalyzer) Analyze(cluster *ClusterData) (*AnalysisResult, error) {
    // Custom performance analysis logic
    score := calculatePerformanceScore(cluster)
    issues := findPerformanceIssues(cluster)
    recommendations := generatePerformanceRecommendations(cluster)
    
    return &AnalysisResult{
        Type:            "performance",
        Score:           score,
        Issues:          issues,
        Recommendations: recommendations,
    }, nil
}

// Register the analyzer
analyzerRegistry.Register(&CustomPerformanceAnalyzer{})
```

---

## 🎯 API Best Practices

### ✅ **Implementation Guidelines**

1. **Error Handling**
   - Always return descriptive errors
   - Use error wrapping for context
   - Implement proper error types

2. **Performance**
   - Use contexts for cancellation
   - Implement proper caching
   - Use goroutines for concurrent operations

3. **Testing**
   - Create comprehensive unit tests
   - Use mocks for external dependencies
   - Implement integration tests

4. **Documentation**
   - Document all public interfaces
   - Provide usage examples
   - Keep documentation up-to-date

---

**Last Updated:** 2025-08-14  
**API Version:** v1  
**Compatible CLI Version:** 2.0.6+