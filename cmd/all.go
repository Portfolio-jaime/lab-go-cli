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

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("✅ Cluster analysis complete!")

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