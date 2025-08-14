package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"k8s-cli/pkg/kubernetes"
)

type ExportData struct {
	Timestamp      time.Time                        `json:"timestamp"`
	ClusterMetrics *kubernetes.ClusterMetrics       `json:"cluster_metrics,omitempty"`
	NodeMetrics    []kubernetes.NodeMetrics         `json:"node_metrics,omitempty"`
	PodMetrics     []kubernetes.PodMetrics          `json:"pod_metrics,omitempty"`
	CostAnalysis   *kubernetes.CostAnalysis         `json:"cost_analysis,omitempty"`
	LogAnalysis    *kubernetes.LogAnalysis          `json:"log_analysis,omitempty"`
	Utilizations   []kubernetes.ResourceUtilization `json:"utilizations,omitempty"`
	Events         []kubernetes.ClusterEvent        `json:"events,omitempty"`
}

type Exporter struct {
	OutputDir string
}

func NewExporter(outputDir string) *Exporter {
	if outputDir == "" {
		outputDir = "./exports"
	}
	return &Exporter{OutputDir: outputDir}
}

func (e *Exporter) ExportToJSON(data *ExportData, filename string) error {
	if err := e.ensureOutputDir(); err != nil {
		return err
	}

	if filename == "" {
		filename = fmt.Sprintf("k8s-cluster-data-%s.json", time.Now().Format("2006-01-02-15-04-05"))
	}

	if !strings.HasSuffix(filename, ".json") {
		filename += ".json"
	}

	filepath := filepath.Join(e.OutputDir, filename)
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filepath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

func (e *Exporter) ExportNodeMetricsToCSV(metrics []kubernetes.NodeMetrics, filename string) error {
	if err := e.ensureOutputDir(); err != nil {
		return err
	}

	if filename == "" {
		filename = fmt.Sprintf("node-metrics-%s.csv", time.Now().Format("2006-01-02-15-04-05"))
	}

	if !strings.HasSuffix(filename, ".csv") {
		filename += ".csv"
	}

	filepath := filepath.Join(e.OutputDir, filename)
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filepath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{
		"Node", "Status", "CPU_Usage", "CPU_Usage_Percent", "Memory_Usage",
		"Memory_Usage_Percent", "CPU_Capacity", "Memory_Capacity",
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV headers: %w", err)
	}

	for _, metric := range metrics {
		record := []string{
			metric.Name,
			metric.Status,
			metric.CPUUsage,
			fmt.Sprintf("%.2f", metric.CPUUsagePercent),
			metric.MemoryUsage,
			fmt.Sprintf("%.2f", metric.MemoryUsagePercent),
			metric.CPUCapacity,
			metric.MemoryCapacity,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

func (e *Exporter) ExportPodMetricsToCSV(metrics []kubernetes.PodMetrics, filename string) error {
	if err := e.ensureOutputDir(); err != nil {
		return err
	}

	if filename == "" {
		filename = fmt.Sprintf("pod-metrics-%s.csv", time.Now().Format("2006-01-02-15-04-05"))
	}

	if !strings.HasSuffix(filename, ".csv") {
		filename += ".csv"
	}

	filepath := filepath.Join(e.OutputDir, filename)
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filepath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{
		"Pod", "Namespace", "Node", "CPU_Usage", "Memory_Usage",
		"CPU_Requests", "Memory_Requests", "CPU_Limits", "Memory_Limits", "Restart_Count",
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV headers: %w", err)
	}

	for _, metric := range metrics {
		record := []string{
			metric.Name,
			metric.Namespace,
			metric.Node,
			metric.CPUUsage,
			metric.MemoryUsage,
			metric.CPURequests,
			metric.MemoryRequests,
			metric.CPULimits,
			metric.MemoryLimits,
			fmt.Sprintf("%d", metric.RestartCount),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

func (e *Exporter) ExportCostAnalysisToCSV(analysis *kubernetes.CostAnalysis, filename string) error {
	if err := e.ensureOutputDir(); err != nil {
		return err
	}

	if filename == "" {
		filename = fmt.Sprintf("cost-analysis-%s.csv", time.Now().Format("2006-01-02-15-04-05"))
	}

	if !strings.HasSuffix(filename, ".csv") {
		filename += ".csv"
	}

	filepath := filepath.Join(e.OutputDir, filename)
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filepath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"=== NODE COSTS ==="}); err != nil {
		return fmt.Errorf("failed to write node costs header: %w", err)
	}
	headers := []string{
		"Node", "Type", "Monthly_Cost", "CPU_Capacity", "Memory_Capacity",
		"CPU_Utilization", "Memory_Utilization", "Efficiency",
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV headers: %w", err)
	}

	for _, node := range analysis.NodeCosts {
		record := []string{
			node.Name,
			node.Type,
			fmt.Sprintf("%.2f", node.MonthlyCost),
			node.CPUCapacity,
			node.MemoryCapacity,
			fmt.Sprintf("%.2f", node.CPUUtilization),
			fmt.Sprintf("%.2f", node.MemUtilization),
			node.Efficiency,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	if err := writer.Write([]string{""}); err != nil {
		return fmt.Errorf("failed to write empty line: %w", err)
	}
	if err := writer.Write([]string{"=== NAMESPACE COSTS ==="}); err != nil {
		return fmt.Errorf("failed to write namespace costs header: %w", err)
	}
	nsHeaders := []string{
		"Namespace", "Monthly_Cost", "Pods_Count", "Cost_Per_Pod",
		"CPU_Requests", "Memory_Requests",
	}
	if err := writer.Write(nsHeaders); err != nil {
		return fmt.Errorf("failed to write namespace headers: %w", err)
	}

	for _, ns := range analysis.NamespaceCosts {
		record := []string{
			ns.Name,
			fmt.Sprintf("%.2f", ns.MonthlyCost),
			fmt.Sprintf("%d", ns.PodsCount),
			fmt.Sprintf("%.2f", ns.CostPerPod),
			ns.CPURequests,
			ns.MemoryRequests,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write namespace record: %w", err)
		}
	}

	return nil
}

func (e *Exporter) ExportUtilizationToCSV(utilizations []kubernetes.ResourceUtilization, filename string) error {
	if err := e.ensureOutputDir(); err != nil {
		return err
	}

	if filename == "" {
		filename = fmt.Sprintf("resource-utilization-%s.csv", time.Now().Format("2006-01-02-15-04-05"))
	}

	if !strings.HasSuffix(filename, ".csv") {
		filename += ".csv"
	}

	filepath := filepath.Join(e.OutputDir, filename)
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filepath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{
		"Type", "Name", "Namespace", "CPU_Utilization", "Memory_Utilization", "Recommendation",
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV headers: %w", err)
	}

	for _, util := range utilizations {
		record := []string{
			util.Type,
			util.Name,
			util.Namespace,
			fmt.Sprintf("%.2f", util.CPUUtilization),
			fmt.Sprintf("%.2f", util.MemUtilization),
			util.Recommendation,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

func (e *Exporter) ExportEventsToCSV(events []kubernetes.ClusterEvent, filename string) error {
	if err := e.ensureOutputDir(); err != nil {
		return err
	}

	if filename == "" {
		filename = fmt.Sprintf("cluster-events-%s.csv", time.Now().Format("2006-01-02-15-04-05"))
	}

	if !strings.HasSuffix(filename, ".csv") {
		filename += ".csv"
	}

	filepath := filepath.Join(e.OutputDir, filename)
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filepath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{
		"Timestamp", "Type", "Severity", "Reason", "Object", "Namespace",
		"Message", "Count", "Component",
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV headers: %w", err)
	}

	for _, event := range events {
		record := []string{
			event.LastTime.Format(time.RFC3339),
			event.Type,
			event.Severity,
			event.Reason,
			event.Object,
			event.Namespace,
			event.Message,
			fmt.Sprintf("%d", event.Count),
			event.Component,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

func (e *Exporter) ExportPrometheusMetrics(data *ExportData, filename string) error {
	if err := e.ensureOutputDir(); err != nil {
		return err
	}

	if filename == "" {
		filename = fmt.Sprintf("prometheus-metrics-%s.txt", time.Now().Format("2006-01-02-15-04-05"))
	}

	if !strings.HasSuffix(filename, ".txt") {
		filename += ".txt"
	}

	filepath := filepath.Join(e.OutputDir, filename)
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filepath, err)
	}
	defer file.Close()

	timestamp := time.Now().Unix()

	if data.ClusterMetrics != nil {
		fmt.Fprintf(file, "# HELP k8s_cluster_cpu_usage_percent Cluster CPU usage percentage\n")
		fmt.Fprintf(file, "# TYPE k8s_cluster_cpu_usage_percent gauge\n")
		fmt.Fprintf(file, "k8s_cluster_cpu_usage_percent %.2f %d\n", data.ClusterMetrics.CPUUsagePercent, timestamp)

		fmt.Fprintf(file, "# HELP k8s_cluster_memory_usage_percent Cluster memory usage percentage\n")
		fmt.Fprintf(file, "# TYPE k8s_cluster_memory_usage_percent gauge\n")
		fmt.Fprintf(file, "k8s_cluster_memory_usage_percent %.2f %d\n", data.ClusterMetrics.MemoryUsagePercent, timestamp)

		fmt.Fprintf(file, "# HELP k8s_cluster_nodes_total Total number of nodes\n")
		fmt.Fprintf(file, "# TYPE k8s_cluster_nodes_total gauge\n")
		fmt.Fprintf(file, "k8s_cluster_nodes_total %d %d\n", data.ClusterMetrics.NodesCount, timestamp)

		fmt.Fprintf(file, "# HELP k8s_cluster_pods_total Total number of pods\n")
		fmt.Fprintf(file, "# TYPE k8s_cluster_pods_total gauge\n")
		fmt.Fprintf(file, "k8s_cluster_pods_total %d %d\n", data.ClusterMetrics.PodsCount, timestamp)
	}

	if data.NodeMetrics != nil {
		fmt.Fprintf(file, "# HELP k8s_node_cpu_usage_percent Node CPU usage percentage\n")
		fmt.Fprintf(file, "# TYPE k8s_node_cpu_usage_percent gauge\n")
		for _, node := range data.NodeMetrics {
			fmt.Fprintf(file, "k8s_node_cpu_usage_percent{node=\"%s\"} %.2f %d\n", node.Name, node.CPUUsagePercent, timestamp)
		}

		fmt.Fprintf(file, "# HELP k8s_node_memory_usage_percent Node memory usage percentage\n")
		fmt.Fprintf(file, "# TYPE k8s_node_memory_usage_percent gauge\n")
		for _, node := range data.NodeMetrics {
			fmt.Fprintf(file, "k8s_node_memory_usage_percent{node=\"%s\"} %.2f %d\n", node.Name, node.MemoryUsagePercent, timestamp)
		}
	}

	if data.CostAnalysis != nil {
		fmt.Fprintf(file, "# HELP k8s_cluster_monthly_cost_usd Estimated monthly cost in USD\n")
		fmt.Fprintf(file, "# TYPE k8s_cluster_monthly_cost_usd gauge\n")
		fmt.Fprintf(file, "k8s_cluster_monthly_cost_usd %.2f %d\n", data.CostAnalysis.TotalMonthlyCost, timestamp)
	}

	return nil
}

func (e *Exporter) ensureOutputDir() error {
	if _, err := os.Stat(e.OutputDir); os.IsNotExist(err) {
		return os.MkdirAll(e.OutputDir, 0755)
	}
	return nil
}

func (e *Exporter) GetExportPath(filename string) string {
	return filepath.Join(e.OutputDir, filename)
}
