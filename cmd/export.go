package cmd

import (
	"fmt"
	"strings"
	"time"

	"k8s-cli/pkg/export"
	"k8s-cli/pkg/kubernetes"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export cluster data to various formats (JSON, CSV, Prometheus)",
	Long:  `Export comprehensive cluster data including metrics, costs, events, and analysis to files for external processing, reporting, or monitoring integration.`,
	RunE:  runExportCommand,
}

var (
	exportFormat     string
	exportOutput     string
	exportFilename   string
	includeMetrics   bool
	includeCosts     bool
	includeLogs      bool
	includeEvents    bool
	exportNamespace  string
	exportHours      int
)

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "json", "Export format: json, csv, prometheus")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "./exports", "Output directory")
	exportCmd.Flags().StringVar(&exportFilename, "filename", "", "Custom filename (without extension)")
	exportCmd.Flags().BoolVar(&includeMetrics, "metrics", true, "Include cluster and pod metrics")
	exportCmd.Flags().BoolVar(&includeCosts, "costs", true, "Include cost analysis")
	exportCmd.Flags().BoolVar(&includeLogs, "logs", true, "Include log analysis")
	exportCmd.Flags().BoolVar(&includeEvents, "events", true, "Include cluster events")
	exportCmd.Flags().StringVarP(&exportNamespace, "namespace", "n", "", "Namespace to export (empty for all)")
	exportCmd.Flags().IntVar(&exportHours, "hours", 24, "Hours of events/logs to include")
}

func runExportCommand(cmd *cobra.Command, args []string) error {
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
	
	client, err := kubernetes.NewClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	exporter := export.NewExporter(exportOutput)

	fmt.Printf("📤 Exporting cluster data to %s format...\n", strings.ToUpper(exportFormat))
	fmt.Println(strings.Repeat("=", 50))

	data := &export.ExportData{
		Timestamp: time.Now(),
	}

	if includeMetrics {
		fmt.Println("📊 Collecting cluster metrics...")
		if metrics, err := client.GetClusterMetrics(); err == nil {
			data.ClusterMetrics = metrics
		}

		if nodeMetrics, err := client.GetRealTimeNodeMetrics(); err == nil {
			data.NodeMetrics = nodeMetrics
		}

		if podMetrics, err := client.GetRealTimePodMetrics(exportNamespace); err == nil {
			data.PodMetrics = podMetrics
		}

		if utilizations, err := client.GetResourceUtilization(); err == nil {
			data.Utilizations = utilizations
		}
	}

	if includeCosts {
		fmt.Println("💰 Collecting cost analysis...")
		if costAnalysis, err := client.GetCostAnalysis(); err == nil {
			data.CostAnalysis = costAnalysis
		}
	}

	if includeLogs {
		fmt.Println("📋 Collecting log analysis...")
		if logAnalysis, err := client.GetLogAnalysis(exportNamespace, exportHours); err == nil {
			data.LogAnalysis = logAnalysis
		}
	}

	if includeEvents {
		fmt.Println("📅 Collecting cluster events...")
		if events, err := client.GetClusterEvents(exportNamespace, exportHours); err == nil {
			data.Events = events
		}
	}

	switch exportFormat {
	case "json":
		if err := exportToJSON(exporter, data); err != nil {
			return err
		}
	case "csv":
		if err := exportToCSV(exporter, data); err != nil {
			return err
		}
	case "prometheus":
		if err := exportToPrometheus(exporter, data); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported export format: %s", exportFormat)
	}

	fmt.Println("✅ Export completed successfully!")
	return nil
}

func exportToJSON(exporter *export.Exporter, data *export.ExportData) error {
	filename := exportFilename
	if filename == "" {
		filename = fmt.Sprintf("k8s-cluster-export-%s", time.Now().Format("2006-01-02-15-04-05"))
	}

	if err := exporter.ExportToJSON(data, filename); err != nil {
		return fmt.Errorf("failed to export JSON: %w", err)
	}

	fullPath := exporter.GetExportPath(filename + ".json")
	fmt.Printf("📄 JSON export saved to: %s\n", fullPath)
	return nil
}

func exportToCSV(exporter *export.Exporter, data *export.ExportData) error {
	baseFilename := exportFilename
	if baseFilename == "" {
		baseFilename = fmt.Sprintf("k8s-export-%s", time.Now().Format("2006-01-02-15-04-05"))
	}

	var exportedFiles []string

	if data.NodeMetrics != nil && len(data.NodeMetrics) > 0 {
		filename := fmt.Sprintf("%s-node-metrics", baseFilename)
		if err := exporter.ExportNodeMetricsToCSV(data.NodeMetrics, filename); err != nil {
			return fmt.Errorf("failed to export node metrics CSV: %w", err)
		}
		exportedFiles = append(exportedFiles, exporter.GetExportPath(filename+".csv"))
	}

	if data.PodMetrics != nil && len(data.PodMetrics) > 0 {
		filename := fmt.Sprintf("%s-pod-metrics", baseFilename)
		if err := exporter.ExportPodMetricsToCSV(data.PodMetrics, filename); err != nil {
			return fmt.Errorf("failed to export pod metrics CSV: %w", err)
		}
		exportedFiles = append(exportedFiles, exporter.GetExportPath(filename+".csv"))
	}

	if data.CostAnalysis != nil {
		filename := fmt.Sprintf("%s-cost-analysis", baseFilename)
		if err := exporter.ExportCostAnalysisToCSV(data.CostAnalysis, filename); err != nil {
			return fmt.Errorf("failed to export cost analysis CSV: %w", err)
		}
		exportedFiles = append(exportedFiles, exporter.GetExportPath(filename+".csv"))
	}

	if data.Utilizations != nil && len(data.Utilizations) > 0 {
		filename := fmt.Sprintf("%s-utilization", baseFilename)
		if err := exporter.ExportUtilizationToCSV(data.Utilizations, filename); err != nil {
			return fmt.Errorf("failed to export utilization CSV: %w", err)
		}
		exportedFiles = append(exportedFiles, exporter.GetExportPath(filename+".csv"))
	}

	if data.Events != nil && len(data.Events) > 0 {
		filename := fmt.Sprintf("%s-events", baseFilename)
		if err := exporter.ExportEventsToCSV(data.Events, filename); err != nil {
			return fmt.Errorf("failed to export events CSV: %w", err)
		}
		exportedFiles = append(exportedFiles, exporter.GetExportPath(filename+".csv"))
	}

	if len(exportedFiles) == 0 {
		return fmt.Errorf("no data available to export")
	}

	fmt.Printf("📊 CSV exports saved to:\n")
	for _, file := range exportedFiles {
		fmt.Printf("  - %s\n", file)
	}

	return nil
}

func exportToPrometheus(exporter *export.Exporter, data *export.ExportData) error {
	filename := exportFilename
	if filename == "" {
		filename = fmt.Sprintf("k8s-prometheus-metrics-%s", time.Now().Format("2006-01-02-15-04-05"))
	}

	if err := exporter.ExportPrometheusMetrics(data, filename); err != nil {
		return fmt.Errorf("failed to export Prometheus metrics: %w", err)
	}

	fullPath := exporter.GetExportPath(filename + ".txt")
	fmt.Printf("📈 Prometheus metrics saved to: %s\n", fullPath)
	fmt.Println("\n💡 You can now scrape these metrics with Prometheus or import into your monitoring system.")
	return nil
}