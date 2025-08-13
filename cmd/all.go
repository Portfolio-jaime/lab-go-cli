package cmd

import (
	"fmt"
	"strings"

	"k8s-cli/pkg/kubernetes"
	"k8s-cli/pkg/recommendations"
	"k8s-cli/pkg/table"

	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Show complete cluster analysis (version, components, resources, and recommendations)",
	Long:  `Run a comprehensive analysis of your Kubernetes cluster including version information, installed components, resource consumption, and optimization recommendations.`,
	RunE:  runAllCommand,
}

func init() {
	rootCmd.AddCommand(allCmd)
}

func runAllCommand(cmd *cobra.Command, args []string) error {
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")

	client, err := kubernetes.NewClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	fmt.Println("🚀 Running complete Kubernetes cluster analysis...")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	if err := showVersionInfo(client); err != nil {
		fmt.Printf("Warning: Could not retrieve version info: %v\n", err)
	}

	if err := showComponentsInfo(client); err != nil {
		fmt.Printf("Warning: Could not retrieve components info: %v\n", err)
	}

	if err := showResourcesInfo(client); err != nil {
		fmt.Printf("Warning: Could not retrieve resources info: %v\n", err)
	}

	if err := showRecommendationsInfo(client); err != nil {
		fmt.Printf("Warning: Could not retrieve recommendations: %v\n", err)
	}

	if err := showRealTimeMetrics(client); err != nil {
		fmt.Printf("Warning: Could not retrieve real-time metrics: %v\n", err)
	}

	if err := showCostOverview(client); err != nil {
		fmt.Printf("Warning: Could not retrieve cost overview: %v\n", err)
	}

	if err := showWorkloadHealth(client); err != nil {
		fmt.Printf("Warning: Could not retrieve workload health: %v\n", err)
	}

	if err := showCriticalEvents(client); err != nil {
		fmt.Printf("Warning: Could not retrieve critical events: %v\n", err)
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("✅ Comprehensive cluster analysis complete!")
	fmt.Println("\n💡 For detailed analysis, use:")
	fmt.Println("  • k8s-cli metrics --pods --utilization")
	fmt.Println("  • k8s-cli cost")
	fmt.Println("  • k8s-cli workload")
	fmt.Println("  • k8s-cli logs")
	fmt.Println("  • k8s-cli export --format json")

	return nil
}

func showVersionInfo(client *kubernetes.Client) error {
	fmt.Println("📊 CLUSTER VERSION INFORMATION")
	fmt.Println(strings.Repeat("-", 40))

	clusterInfo, err := client.GetClusterVersion()
	if err != nil {
		return err
	}

	versionTable := table.NewTable([]string{"Property", "Value"})
	versionTable.AddRow([]string{"Kubernetes Version", clusterInfo.GitVersion})
	versionTable.AddRow([]string{"Platform", clusterInfo.Platform})
	versionTable.AddRow([]string{"Build Date", clusterInfo.BuildDate})
	versionTable.Render()
	fmt.Println()

	return nil
}

func showComponentsInfo(client *kubernetes.Client) error {
	fmt.Println("🔧 INSTALLED COMPONENTS")
	fmt.Println(strings.Repeat("-", 40))

	components, err := client.GetInstalledComponents()
	if err != nil {
		return err
	}

	if len(components) == 0 {
		fmt.Println("No common components detected.")
		fmt.Println()
		return nil
	}

	componentTable := table.NewTable([]string{"Component", "Namespace", "Status", "Version"})
	for _, comp := range components {
		componentTable.AddRow([]string{comp.Name, comp.Namespace, comp.Status, comp.Version})
	}
	componentTable.Render()
	fmt.Println()

	return nil
}

func showResourcesInfo(client *kubernetes.Client) error {
	fmt.Println("📈 CLUSTER RESOURCES")
	fmt.Println(strings.Repeat("-", 40))

	summary, err := client.GetSimpleClusterSummary()
	if err != nil {
		return err
	}

	summaryTable := table.NewTable([]string{"Metric", "Value"})
	summaryTable.AddRow([]string{"Total Nodes", fmt.Sprintf("%d", summary.TotalNodes)})
	summaryTable.AddRow([]string{"Total Pods", fmt.Sprintf("%d", summary.TotalPods)})
	summaryTable.AddRow([]string{"CPU Capacity", summary.TotalCPUCapacity + " cores"})
	summaryTable.AddRow([]string{"Memory Capacity", summary.TotalMemCapacity})
	summaryTable.Render()
	fmt.Println()

	nodes, err := client.GetSimpleNodesInfo()
	if err == nil && len(nodes) > 0 {
		fmt.Println("🖥️  Node Summary:")
		nodeTable := table.NewTable([]string{"Node", "Status", "Role", "Age"})
		for _, node := range nodes {
			nodeTable.AddRow([]string{node.Name, node.Status, node.Roles, node.Age})
		}
		nodeTable.Render()
		fmt.Println()
	}

	return nil
}

func showRecommendationsInfo(client *kubernetes.Client) error {
	fmt.Println("💡 RECOMMENDATIONS")
	fmt.Println(strings.Repeat("-", 40))

	analyzer := recommendations.NewRecommendationAnalyzer(client)
	recs, err := analyzer.AnalyzeCluster()
	if err != nil {
		return err
	}

	if len(recs) == 0 {
		fmt.Println("✅ No recommendations - cluster looks good!")
		fmt.Println()
		return nil
	}

	highPriority := 0
	mediumPriority := 0
	lowPriority := 0

	for _, rec := range recs {
		switch rec.Severity {
		case "High":
			highPriority++
		case "Medium":
			mediumPriority++
		case "Low":
			lowPriority++
		}
	}

	recSummaryTable := table.NewTable([]string{"Severity", "Count"})
	if highPriority > 0 {
		recSummaryTable.AddRow([]string{"High Priority", fmt.Sprintf("%d", highPriority)})
	}
	if mediumPriority > 0 {
		recSummaryTable.AddRow([]string{"Medium Priority", fmt.Sprintf("%d", mediumPriority)})
	}
	if lowPriority > 0 {
		recSummaryTable.AddRow([]string{"Low Priority", fmt.Sprintf("%d", lowPriority)})
	}
	recSummaryTable.Render()

	fmt.Printf("\n💡 Run 'k8s-cli recommend' for detailed recommendations.\n")
	fmt.Println()

	return nil
}

func showRealTimeMetrics(client *kubernetes.Client) error {
	fmt.Println("📊 REAL-TIME METRICS OVERVIEW")
	fmt.Println(strings.Repeat("-", 40))

	clusterMetrics, err := client.GetClusterMetrics()
	if err != nil {
		return err
	}

	metricsTable := table.NewTable([]string{"Resource", "Usage", "Capacity", "Utilization"})
	metricsTable.AddRow([]string{
		"CPU",
		clusterMetrics.TotalCPUUsage,
		clusterMetrics.TotalCPUCapacity,
		fmt.Sprintf("%.1f%%", clusterMetrics.CPUUsagePercent),
	})
	metricsTable.AddRow([]string{
		"Memory",
		clusterMetrics.TotalMemoryUsage,
		clusterMetrics.TotalMemoryCapacity,
		fmt.Sprintf("%.1f%%", clusterMetrics.MemoryUsagePercent),
	})
	metricsTable.Render()
	fmt.Println()

	return nil
}

func showCostOverview(client *kubernetes.Client) error {
	fmt.Println("💰 COST OVERVIEW")
	fmt.Println(strings.Repeat("-", 40))

	analysis, err := client.GetCostAnalysis()
	if err != nil {
		return err
	}

	totalSavings := 0.0
	for _, resource := range analysis.UnderutilizedResources {
		totalSavings += resource.EstimatedSavings
	}

	costTable := table.NewTable([]string{"Metric", "Value"})
	costTable.AddRow([]string{"Monthly Cost", fmt.Sprintf("$%.2f", analysis.TotalMonthlyCost)})
	costTable.AddRow([]string{"Potential Savings", fmt.Sprintf("$%.2f", totalSavings)})
	costTable.AddRow([]string{"Underutilized Resources", fmt.Sprintf("%d", len(analysis.UnderutilizedResources))})
	costTable.Render()
	fmt.Println()

	return nil
}

func showWorkloadHealth(client *kubernetes.Client) error {
	fmt.Println("🔍 WORKLOAD HEALTH SUMMARY")
	fmt.Println(strings.Repeat("-", 40))

	analysis, err := client.GetWorkloadAnalysis("")
	if err != nil {
		return err
	}

	workloadTable := table.NewTable([]string{"Workload Type", "Total", "Healthy", "Issues"})
	workloadTable.AddRow([]string{
		"Deployments",
		fmt.Sprintf("%d", analysis.WorkloadSummary.TotalDeployments),
		fmt.Sprintf("%d", analysis.WorkloadSummary.HealthyDeployments),
		fmt.Sprintf("%d", analysis.WorkloadSummary.TotalDeployments-analysis.WorkloadSummary.HealthyDeployments),
	})
	workloadTable.AddRow([]string{
		"StatefulSets",
		fmt.Sprintf("%d", analysis.WorkloadSummary.TotalStatefulSets),
		fmt.Sprintf("%d", analysis.WorkloadSummary.HealthyStatefulSets),
		fmt.Sprintf("%d", analysis.WorkloadSummary.TotalStatefulSets-analysis.WorkloadSummary.HealthyStatefulSets),
	})
	workloadTable.AddRow([]string{
		"DaemonSets",
		fmt.Sprintf("%d", analysis.WorkloadSummary.TotalDaemonSets),
		fmt.Sprintf("%d", analysis.WorkloadSummary.HealthyDaemonSets),
		fmt.Sprintf("%d", analysis.WorkloadSummary.TotalDaemonSets-analysis.WorkloadSummary.HealthyDaemonSets),
	})

	healthStatus := fmt.Sprintf("%d/100", analysis.WorkloadSummary.OverallHealthScore)
	if analysis.WorkloadSummary.CriticalIssues > 0 {
		healthStatus += fmt.Sprintf(" (%d critical)", analysis.WorkloadSummary.CriticalIssues)
	}

	workloadTable.AddRow([]string{"Overall Health", healthStatus, "", ""})
	workloadTable.Render()
	fmt.Println()

	return nil
}

func showCriticalEvents(client *kubernetes.Client) error {
	fmt.Println("🚨 RECENT CRITICAL EVENTS")
	fmt.Println(strings.Repeat("-", 40))

	events, err := client.GetClusterEvents("", 1)
	if err != nil {
		return err
	}

	criticalEvents := []kubernetes.ClusterEvent{}
	for _, event := range events {
		if event.Severity == "Critical" {
			criticalEvents = append(criticalEvents, event)
		}
	}

	if len(criticalEvents) == 0 {
		fmt.Println("✅ No critical events in the last hour")
		fmt.Println()
		return nil
	}

	eventsTable := table.NewTable([]string{"Object", "Reason", "Message", "Count"})
	for i, event := range criticalEvents {
		if i >= 5 {
			break
		}

		message := event.Message
		if len(message) > 40 {
			message = message[:37] + "..."
		}

		eventsTable.AddRow([]string{
			event.Object,
			event.Reason,
			message,
			fmt.Sprintf("%d", event.Count),
		})
	}
	eventsTable.Render()

	if len(criticalEvents) > 5 {
		fmt.Printf("... and %d more critical events\n", len(criticalEvents)-5)
	}
	fmt.Println()

	return nil
}
