package cmd

import (
	"fmt"
	"strings"

	"k8s-cli/pkg/kubernetes"
	"k8s-cli/pkg/table"

	"github.com/spf13/cobra"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Show real-time cluster metrics and resource utilization",
	Long:  `Display current CPU and memory usage for nodes and pods, along with utilization analysis and recommendations.`,
	RunE:  runMetricsCommand,
}

var (
	showMetricsNodes     bool
	showMetricsPods      bool
	showMetricsUtilization bool
	metricsNamespace     string
)

func init() {
	rootCmd.AddCommand(metricsCmd)
	metricsCmd.Flags().BoolVar(&showMetricsNodes, "nodes", true, "Show node metrics")
	metricsCmd.Flags().BoolVar(&showMetricsPods, "pods", false, "Show pod metrics")
	metricsCmd.Flags().BoolVar(&showMetricsUtilization, "utilization", false, "Show resource utilization analysis")
	metricsCmd.Flags().StringVarP(&metricsNamespace, "namespace", "n", "", "Namespace for pod metrics (empty for all)")
}

func runMetricsCommand(cmd *cobra.Command, args []string) error {
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
	
	client, err := kubernetes.NewClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	fmt.Println("📊 Real-time Cluster Metrics")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	if err := showClusterMetrics(client); err != nil {
		fmt.Printf("Warning: Could not retrieve cluster metrics: %v\n", err)
	}

	if showMetricsNodes {
		if err := showNodeMetrics(client); err != nil {
			fmt.Printf("Warning: Could not retrieve node metrics: %v\n", err)
		}
	}

	if showMetricsPods {
		if err := showPodMetrics(client, metricsNamespace); err != nil {
			fmt.Printf("Warning: Could not retrieve pod metrics: %v\n", err)
		}
	}

	if showMetricsUtilization {
		if err := showUtilizationAnalysis(client); err != nil {
			fmt.Printf("Warning: Could not retrieve utilization analysis: %v\n", err)
		}
	}

	return nil
}

func showClusterMetrics(client *kubernetes.Client) error {
	fmt.Println("🌐 CLUSTER OVERVIEW")
	fmt.Println(strings.Repeat("-", 40))
	
	metrics, err := client.GetClusterMetrics()
	if err != nil {
		return err
	}

	overviewTable := table.NewTable([]string{"Metric", "Usage", "Capacity", "Utilization"})
	overviewTable.AddRow([]string{
		"CPU",
		metrics.TotalCPUUsage,
		metrics.TotalCPUCapacity,
		fmt.Sprintf("%.1f%%", metrics.CPUUsagePercent),
	})
	overviewTable.AddRow([]string{
		"Memory",
		metrics.TotalMemoryUsage,
		metrics.TotalMemoryCapacity,
		fmt.Sprintf("%.1f%%", metrics.MemoryUsagePercent),
	})
	overviewTable.Render()

	summaryTable := table.NewTable([]string{"Resource", "Count"})
	summaryTable.AddRow([]string{"Nodes", fmt.Sprintf("%d", metrics.NodesCount)})
	summaryTable.AddRow([]string{"Pods", fmt.Sprintf("%d", metrics.PodsCount)})
	summaryTable.AddRow([]string{"Namespaces", fmt.Sprintf("%d", metrics.NamespacesCount)})
	summaryTable.Render()
	
	fmt.Println()
	return nil
}

func showNodeMetrics(client *kubernetes.Client) error {
	fmt.Println("🖥️  NODE METRICS")
	fmt.Println(strings.Repeat("-", 40))
	
	nodeMetrics, err := client.GetRealTimeNodeMetrics()
	if err != nil {
		return err
	}

	if len(nodeMetrics) == 0 {
		fmt.Println("No node metrics available. Make sure metrics-server is installed.")
		fmt.Println()
		return nil
	}

	nodeTable := table.NewTable([]string{"Node", "Status", "CPU Usage", "CPU %", "Memory Usage", "Memory %"})
	for _, node := range nodeMetrics {
		status := node.Status
		if node.CPUUsagePercent > 80 || node.MemoryUsagePercent > 80 {
			status += " ⚠️"
		}
		
		nodeTable.AddRow([]string{
			node.Name,
			status,
			fmt.Sprintf("%s / %s", node.CPUUsage, node.CPUCapacity),
			fmt.Sprintf("%.1f%%", node.CPUUsagePercent),
			fmt.Sprintf("%s / %s", node.MemoryUsage, node.MemoryCapacity),
			fmt.Sprintf("%.1f%%", node.MemoryUsagePercent),
		})
	}
	nodeTable.Render()
	fmt.Println()

	return nil
}

func showPodMetrics(client *kubernetes.Client, namespace string) error {
	fmt.Println("🚀 POD METRICS")
	if namespace != "" {
		fmt.Printf("Namespace: %s\n", namespace)
	} else {
		fmt.Println("All Namespaces")
	}
	fmt.Println(strings.Repeat("-", 40))
	
	podMetrics, err := client.GetRealTimePodMetrics(namespace)
	if err != nil {
		return err
	}

	if len(podMetrics) == 0 {
		fmt.Println("No pod metrics available. Make sure metrics-server is installed.")
		fmt.Println()
		return nil
	}

	podTable := table.NewTable([]string{"Pod", "Namespace", "CPU Usage", "Memory Usage", "Restarts", "Node"})
	for _, pod := range podMetrics {
		restartInfo := fmt.Sprintf("%d", pod.RestartCount)
		if pod.RestartCount > 5 {
			restartInfo += " ⚠️"
		}
		
		podTable.AddRow([]string{
			pod.Name,
			pod.Namespace,
			pod.CPUUsage,
			pod.MemoryUsage,
			restartInfo,
			pod.Node,
		})
	}
	podTable.Render()
	fmt.Println()

	return nil
}

func showUtilizationAnalysis(client *kubernetes.Client) error {
	fmt.Println("📈 RESOURCE UTILIZATION ANALYSIS")
	fmt.Println(strings.Repeat("-", 40))
	
	utilizations, err := client.GetResourceUtilization()
	if err != nil {
		return err
	}

	if len(utilizations) == 0 {
		fmt.Println("No utilization data available.")
		fmt.Println()
		return nil
	}

	underutilized := 0
	overutilized := 0
	optimal := 0

	utilizationTable := table.NewTable([]string{"Pod", "Namespace", "CPU %", "Memory %", "Recommendation"})
	for _, util := range utilizations {
		if strings.Contains(util.Recommendation, "underutilized") {
			underutilized++
		} else if strings.Contains(util.Recommendation, "overutilized") {
			overutilized++
		} else {
			optimal++
		}

		cpuPercent := "N/A"
		if util.CPUUtilization > 0 {
			cpuPercent = fmt.Sprintf("%.1f%%", util.CPUUtilization)
		}
		
		memPercent := "N/A"
		if util.MemUtilization > 0 {
			memPercent = fmt.Sprintf("%.1f%%", util.MemUtilization)
		}

		utilizationTable.AddRow([]string{
			util.Name,
			util.Namespace,
			cpuPercent,
			memPercent,
			util.Recommendation,
		})
	}
	utilizationTable.Render()

	fmt.Printf("\n📊 Utilization Summary:\n")
	summaryTable := table.NewTable([]string{"Category", "Count"})
	summaryTable.AddRow([]string{"Under-utilized", fmt.Sprintf("%d", underutilized)})
	summaryTable.AddRow([]string{"Over-utilized", fmt.Sprintf("%d", overutilized)})
	summaryTable.AddRow([]string{"Optimal", fmt.Sprintf("%d", optimal)})
	summaryTable.Render()

	fmt.Println()
	return nil
}