# 📊 k8s-cli Command Structure & Flow Diagrams

## 🎯 Complete Command Overview

```mermaid
graph TB
    subgraph "k8s-cli Command Hierarchy"
        CLI[k8s-cli]
        
        subgraph "Core Analysis Commands"
            ALL[all<br/>📊 Complete Analysis]
            VERSION[version<br/>ℹ️ Cluster Info]
            RESOURCES[resources<br/>📦 Basic Resources]
            RECOMMEND[recommend<br/>💡 Suggestions]
        end
        
        subgraph "Advanced Analytics"
            METRICS[metrics<br/>📈 Real-time Metrics]
            COST[cost<br/>💰 Cost Analysis]
            WORKLOAD[workload<br/>🔍 Health Scoring]
            LOGS[logs<br/>📋 Event Analysis]
        end
        
        subgraph "Enterprise Integration"
            EXPORT[export<br/>📤 Data Export]
        end
        
        subgraph "Data Flow"
            K8S_API[(Kubernetes API)]
            METRICS_SERVER[(Metrics Server)]
            EVENT_STORE[(Event Store)]
        end
        
        subgraph "Output Formats"
            CONSOLE[🖥️ Terminal]
            JSON[📄 JSON Files]
            CSV[📊 CSV Reports]
            PROMETHEUS[📈 Prometheus]
        end
    end
    
    CLI --> ALL
    CLI --> VERSION
    CLI --> RESOURCES
    CLI --> RECOMMEND
    CLI --> METRICS
    CLI --> COST
    CLI --> WORKLOAD
    CLI --> LOGS
    CLI --> EXPORT
    
    ALL --> K8S_API
    ALL --> METRICS_SERVER
    ALL --> EVENT_STORE
    
    METRICS --> METRICS_SERVER
    COST --> K8S_API
    WORKLOAD --> K8S_API
    LOGS --> EVENT_STORE
    VERSION --> K8S_API
    RESOURCES --> K8S_API
    RECOMMEND --> K8S_API
    
    ALL --> CONSOLE
    METRICS --> CONSOLE
    COST --> CONSOLE
    WORKLOAD --> CONSOLE
    LOGS --> CONSOLE
    
    EXPORT --> JSON
    EXPORT --> CSV
    EXPORT --> PROMETHEUS
    
    classDef coreCmd fill:#e3f2fd,stroke:#1976d2,stroke-width:3px
    classDef advancedCmd fill:#f3e5f5,stroke:#7b1fa2,stroke-width:3px
    classDef enterpriseCmd fill:#e8f5e8,stroke:#388e3c,stroke-width:3px
    classDef dataSource fill:#fff3e0,stroke:#f57c00,stroke-width:2px
    classDef output fill:#fce4ec,stroke:#c2185b,stroke-width:2px
    
    class ALL,VERSION,RESOURCES,RECOMMEND coreCmd
    class METRICS,COST,WORKLOAD,LOGS advancedCmd
    class EXPORT enterpriseCmd
    class K8S_API,METRICS_SERVER,EVENT_STORE dataSource
    class CONSOLE,JSON,CSV,PROMETHEUS output
```

## 🔄 Command Execution Flow

```mermaid
sequenceDiagram
    participant User
    participant CLI as k8s-cli
    participant K8s as Kubernetes API
    participant Metrics as Metrics Server
    participant Analysis as Analysis Engine
    participant Export as Export Engine
    participant Output as Output/Files

    Note over User,Output: Complete Analysis Flow (k8s-cli all)
    
    User->>CLI: k8s-cli all
    CLI->>K8s: Get cluster info
    K8s-->>CLI: Version, nodes, pods
    CLI->>Metrics: Get real-time metrics
    Metrics-->>CLI: CPU/Memory usage
    CLI->>Analysis: Analyze costs & workloads
    Analysis-->>CLI: Health scores, costs
    CLI->>K8s: Get recent events
    K8s-->>CLI: Critical events
    CLI->>Output: Formatted report
    Output-->>User: Complete analysis

    Note over User,Output: Metrics Analysis Flow
    
    User->>CLI: k8s-cli metrics --utilization
    CLI->>Metrics: Get node metrics
    Metrics-->>CLI: Real-time usage
    CLI->>Analysis: Calculate utilization
    Analysis-->>CLI: Efficiency scores
    CLI->>Output: Metrics report
    Output-->>User: Utilization analysis

    Note over User,Output: Export Flow
    
    User->>CLI: k8s-cli export --format json
    CLI->>K8s: Collect all data
    CLI->>Metrics: Get metrics data
    CLI->>Analysis: Process data
    Analysis-->>CLI: Analyzed data
    CLI->>Export: Format as JSON
    Export->>Output: Write JSON file
    Output-->>User: Export complete
```

## 🎯 Command Feature Matrix

```mermaid
gitgraph
    commit id: "v1.0 - Basic CLI"
    branch core-commands
    checkout core-commands
    commit id: "all - Basic analysis"
    commit id: "version - Cluster info"
    commit id: "resources - Resource list"
    commit id: "recommend - Basic tips"
    
    checkout main
    merge core-commands
    commit id: "v1.5 - Enhanced Core"
    
    branch advanced-analytics
    checkout advanced-analytics
    commit id: "metrics - Real-time data"
    commit id: "cost - Cost analysis"
    commit id: "workload - Health scoring"
    commit id: "logs - Event analysis"
    
    checkout main
    merge advanced-analytics
    commit id: "v2.0 - Enterprise Platform"
    
    branch enterprise-features
    checkout enterprise-features
    commit id: "export - Multi-format"
    commit id: "integration - APIs"
    commit id: "automation - CI/CD"
    
    checkout main
    merge enterprise-features
    commit id: "v2.0 - Production Ready"
```

## 📊 Command Complexity & User Journey

```mermaid
journey
    title k8s-cli User Journey
    section Discovery
      Try basic command: 5: User
      See cluster info: 4: User
      Understand capabilities: 5: User
    section Analysis
      Run complete analysis: 5: User
      Check real-time metrics: 4: User
      Analyze costs: 5: User
      Review workload health: 4: User
    section Optimization
      Identify issues: 5: User
      Get recommendations: 5: User
      Export data: 4: User
      Implement changes: 3: User
    section Integration
      Setup automation: 4: User
      Monitor continuously: 5: User
      Share with team: 5: User
```

## 🔧 Command Flag Hierarchy

```mermaid
mindmap
  root((Command Flags))
    Global Flags
      --kubeconfig
        Custom kubeconfig path
        Default: ~/.kube/config
      --config
        CLI configuration file
        Default: ~/.k8s-cli.yaml
      --namespace -n
        Target namespace
        Default: all namespaces
      --output -o
        Output format
        table, json, yaml
    
    Metrics Flags
      --nodes
        Show node metrics
        Real-time CPU/Memory
      --pods
        Show pod metrics
        Per-namespace filtering
      --utilization
        Utilization analysis
        Efficiency scoring
    
    Cost Flags
      --nodes
        Node cost breakdown
        Instance type pricing
      --namespaces
        Namespace costs
        Per-pod calculations
      --underutilized
        Wasted resources
        Savings potential
      --optimizations
        Cost recommendations
        Priority scoring
    
    Workload Flags
      --deployments
        Deployment health
        Configuration issues
      --statefulsets
        StatefulSet analysis
        Storage considerations
      --daemonsets
        DaemonSet status
        Node coverage
      --pods
        Pod-level details
        Restart patterns
      --unhealthy-only
        Problem workloads
        Issue filtering
    
    Logs Flags
      --critical
        Critical events only
        High-priority issues
      --patterns
        Error pattern analysis
        Frequency tracking
      --security-events
        Security-related
        Risk assessment
      --hours
        Time window
        Default: 24 hours
    
    Export Flags
      --format -f
        json, csv, prometheus
        Multi-format support
      --output -o
        Output directory
        File organization
      --filename
        Custom filename
        Timestamp options
      --costs
        Include cost data
        Financial analysis
      --metrics
        Include metrics
        Performance data
      --logs
        Include log analysis
        Event correlation
```

## 🎯 Use Case Mapping

```mermaid
graph LR
    subgraph "User Personas"
        DEVOPS[DevOps Engineer]
        FINOPS[FinOps Analyst]
        SRE[SRE Engineer]
        MANAGER[Engineering Manager]
    end
    
    subgraph "Primary Commands"
        METRICS_CMD[metrics]
        COST_CMD[cost]
        WORKLOAD_CMD[workload]
        LOGS_CMD[logs]
        EXPORT_CMD[export]
        ALL_CMD[all]
    end
    
    subgraph "Business Outcomes"
        OPTIMIZATION[Cost Optimization]
        RELIABILITY[System Reliability]
        VISIBILITY[Operational Visibility]
        COMPLIANCE[Compliance Reporting]
    end
    
    DEVOPS --> METRICS_CMD
    DEVOPS --> WORKLOAD_CMD
    DEVOPS --> LOGS_CMD
    
    FINOPS --> COST_CMD
    FINOPS --> EXPORT_CMD
    
    SRE --> LOGS_CMD
    SRE --> WORKLOAD_CMD
    SRE --> ALL_CMD
    
    MANAGER --> ALL_CMD
    MANAGER --> EXPORT_CMD
    
    METRICS_CMD --> VISIBILITY
    COST_CMD --> OPTIMIZATION
    WORKLOAD_CMD --> RELIABILITY
    LOGS_CMD --> RELIABILITY
    EXPORT_CMD --> COMPLIANCE
    ALL_CMD --> VISIBILITY
    
    classDef persona fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef command fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px
    classDef outcome fill:#e8f5e8,stroke:#388e3c,stroke-width:2px
    
    class DEVOPS,FINOPS,SRE,MANAGER persona
    class METRICS_CMD,COST_CMD,WORKLOAD_CMD,LOGS_CMD,EXPORT_CMD,ALL_CMD command
    class OPTIMIZATION,RELIABILITY,VISIBILITY,COMPLIANCE outcome
```

## 📈 Command Performance Profile

```mermaid
gantt
    title Command Execution Time Profile
    dateFormat X
    axisFormat %s
    
    section Quick Commands
    version     :0, 1s
    resources   :0, 2s
    recommend   :0, 3s
    
    section Analysis Commands
    metrics     :0, 5s
    workload    :0, 8s
    cost        :0, 10s
    logs        :0, 12s
    
    section Comprehensive
    all         :0, 15s
    export      :0, 20s
```

## 🔄 Integration Patterns

```mermaid
graph TB
    subgraph "CI/CD Integration"
        PIPELINE[CI/CD Pipeline]
        HEALTH_CHECK[k8s-cli workload --unhealthy-only]
        COST_CHECK[k8s-cli cost --underutilized]
        EXPORT_DATA[k8s-cli export --format json]
    end
    
    subgraph "Monitoring Integration"
        PROMETHEUS[Prometheus]
        GRAFANA[Grafana]
        METRICS_EXPORT[k8s-cli export --format prometheus]
    end
    
    subgraph "Business Intelligence"
        BI_TOOLS[BI Tools]
        CSV_EXPORT[k8s-cli export --format csv]
        REPORTS[Automated Reports]
    end
    
    subgraph "Incident Response"
        ALERTING[Alert System]
        LOGS_ANALYSIS[k8s-cli logs --critical]
        WORKLOAD_CHECK[k8s-cli workload --unhealthy-only]
    end
    
    PIPELINE --> HEALTH_CHECK
    PIPELINE --> COST_CHECK
    PIPELINE --> EXPORT_DATA
    
    METRICS_EXPORT --> PROMETHEUS
    PROMETHEUS --> GRAFANA
    
    CSV_EXPORT --> BI_TOOLS
    BI_TOOLS --> REPORTS
    
    ALERTING --> LOGS_ANALYSIS
    ALERTING --> WORKLOAD_CHECK
    
    classDef integration fill:#fff3e0,stroke:#f57c00,stroke-width:2px
    classDef tool fill:#e8eaf6,stroke:#3f51b5,stroke-width:2px
    classDef command fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px
    
    class PIPELINE,PROMETHEUS,BI_TOOLS,ALERTING integration
    class GRAFANA,REPORTS tool
    class HEALTH_CHECK,COST_CHECK,EXPORT_DATA,METRICS_EXPORT,CSV_EXPORT,LOGS_ANALYSIS,WORKLOAD_CHECK command
```

---

This diagram suite provides a comprehensive visual understanding of the k8s-cli command structure, execution flows, and integration patterns, making it easier for users to understand how to leverage the platform's capabilities.