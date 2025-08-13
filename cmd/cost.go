package cmd

import (
	"fmt"
	"strings"

	"k8s-cli/pkg/kubernetes"
	"k8s-cli/pkg/table"

	"github.com/spf13/cobra"
)

var costCmd = &cobra.Command{
	Use:   "cost",
	Short: "Analyze cluster costs and identify optimization opportunities",
	Long:  `Provide detailed cost analysis including node costs, namespace costs, underutilized resources, and optimization recommendations.`,
	RunE:  runCostCommand,
}

var (
	showCostNodes        bool
	showCostNamespaces   bool
	showCostUnderutilized bool
	showCostOptimizations bool
)

func init() {
	rootCmd.AddCommand(costCmd)
	costCmd.Flags().BoolVar(&showCostNodes, "nodes", true, "Show node cost breakdown")
	costCmd.Flags().BoolVar(&showCostNamespaces, "namespaces", true, "Show namespace cost analysis")
	costCmd.Flags().BoolVar(&showCostUnderutilized, "underutilized", true, "Show underutilized resources")
	costCmd.Flags().BoolVar(&showCostOptimizations, "optimizations", true, "Show cost optimization recommendations")
}

func runCostCommand(cmd *cobra.Command, args []string) error {
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
	
	client, err := kubernetes.NewClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	fmt.Println("💰 Cluster Cost Analysis")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	analysis, err := client.GetCostAnalysis()
	if err != nil {
		return fmt.Errorf("failed to get cost analysis: %w", err)
	}

	showCostSummary(analysis)

	if showCostNodes {
		showNodeCosts(analysis.NodeCosts)
	}

	if showCostNamespaces {
		showNamespaceCosts(analysis.NamespaceCosts)
	}

	if showCostUnderutilized {
		showUnderutilizedResources(analysis.UnderutilizedResources)
	}

	if showCostOptimizations {
		showOptimizationRecommendations(analysis.CostOptimizations)
	}

	return nil
}

func showCostSummary(analysis *kubernetes.CostAnalysis) {
	fmt.Println("📊 COST OVERVIEW")
	fmt.Println(strings.Repeat("-", 40))

	totalSavings := 0.0
	for _, resource := range analysis.UnderutilizedResources {
		totalSavings += resource.EstimatedSavings
	}

	overviewTable := table.NewTable([]string{"Metric", "Value"})
	overviewTable.AddRow([]string{"Total Monthly Cost", fmt.Sprintf("$%.2f", analysis.TotalMonthlyCost)})
	overviewTable.AddRow([]string{"Potential Monthly Savings", fmt.Sprintf("$%.2f", totalSavings)})
	overviewTable.AddRow([]string{"Cost Efficiency", fmt.Sprintf("%.1f%%", (analysis.TotalMonthlyCost-totalSavings)/analysis.TotalMonthlyCost*100)})
	overviewTable.AddRow([]string{"Optimization Opportunities", fmt.Sprintf("%d", len(analysis.CostOptimizations))})
	overviewTable.Render()
	fmt.Println()
}

func showNodeCosts(nodeCosts []kubernetes.NodeCost) {
	if len(nodeCosts) == 0 {
		return
	}

	fmt.Println("🖥️  NODE COSTS")
	fmt.Println(strings.Repeat("-", 40))

	nodeTable := table.NewTable([]string{"Node", "Type", "Monthly Cost", "CPU Util", "Memory Util", "Efficiency"})
	for _, node := range nodeCosts {
		costDisplay := fmt.Sprintf("$%.2f", node.MonthlyCost)
		if node.CPUUtilization < 30 || node.MemUtilization < 30 {
			costDisplay += " ⚠️"
		}

		cpuUtil := "N/A"
		memUtil := "N/A"
		if node.CPUUtilization > 0 {
			cpuUtil = fmt.Sprintf("%.1f%%", node.CPUUtilization)
		}
		if node.MemUtilization > 0 {
			memUtil = fmt.Sprintf("%.1f%%", node.MemUtilization)
		}

		nodeTable.AddRow([]string{
			node.Name,
			node.Type,
			costDisplay,
			cpuUtil,
			memUtil,
			node.Efficiency,
		})
	}
	nodeTable.Render()
	fmt.Println()
}

func showNamespaceCosts(namespaceCosts []kubernetes.NamespaceCost) {
	if len(namespaceCosts) == 0 {
		return
	}

	fmt.Println("🏢 NAMESPACE COSTS")
	fmt.Println(strings.Repeat("-", 40))

	namespaceTable := table.NewTable([]string{"Namespace", "Monthly Cost", "Pods", "Cost/Pod", "CPU Requests", "Memory Requests"})
	for _, ns := range namespaceCosts {
		if ns.MonthlyCost < 1.0 {
			continue
		}

		namespaceTable.AddRow([]string{
			ns.Name,
			fmt.Sprintf("$%.2f", ns.MonthlyCost),
			fmt.Sprintf("%d", ns.PodsCount),
			fmt.Sprintf("$%.2f", ns.CostPerPod),
			ns.CPURequests,
			ns.MemoryRequests,
		})
	}
	namespaceTable.Render()
	fmt.Println()
}

func showUnderutilizedResources(resources []kubernetes.UnderutilizedResource) {
	if len(resources) == 0 {
		fmt.Println("✅ No significantly underutilized resources found!")
		fmt.Println()
		return
	}

	fmt.Println("📉 UNDERUTILIZED RESOURCES")
	fmt.Println(strings.Repeat("-", 40))

	resourceTable := table.NewTable([]string{"Pod", "Namespace", "CPU Waste", "Memory Waste", "Monthly Savings", "Recommendation"})
	totalSavings := 0.0
	
	for i, resource := range resources {
		if i >= 10 {
			break
		}
		
		totalSavings += resource.EstimatedSavings
		resourceTable.AddRow([]string{
			resource.Name,
			resource.Namespace,
			resource.CPUWaste,
			resource.MemoryWaste,
			fmt.Sprintf("$%.2f", resource.EstimatedSavings),
			resource.Recommendation,
		})
	}
	resourceTable.Render()

	if len(resources) > 10 {
		fmt.Printf("... and %d more underutilized resources\n", len(resources)-10)
	}
	
	fmt.Printf("\n💡 Total potential monthly savings from top resources: $%.2f\n", totalSavings)
	fmt.Println()
}

func showOptimizationRecommendations(optimizations []kubernetes.CostOptimization) {
	if len(optimizations) == 0 {
		return
	}

	fmt.Println("🎯 COST OPTIMIZATION RECOMMENDATIONS")
	fmt.Println(strings.Repeat("-", 40))

	optimizationTable := table.NewTable([]string{"Priority", "Type", "Description", "Potential Savings", "Action"})
	totalPotentialSavings := 0.0

	for _, opt := range optimizations {
		totalPotentialSavings += opt.PotentialSavings
		
		savingsDisplay := "N/A"
		if opt.PotentialSavings > 0 {
			savingsDisplay = fmt.Sprintf("$%.2f/mo", opt.PotentialSavings)
		}

		priority := opt.Priority
		if opt.Priority == "High" {
			priority = "🔴 " + priority
		} else if opt.Priority == "Medium" {
			priority = "🟡 " + priority
		} else {
			priority = "🟢 " + priority
		}

		optimizationTable.AddRow([]string{
			priority,
			opt.Type,
			opt.Description,
			savingsDisplay,
			opt.Action,
		})
	}
	optimizationTable.Render()

	if totalPotentialSavings > 0 {
		fmt.Printf("\n💰 Total potential monthly savings: $%.2f\n", totalPotentialSavings)
	}
	fmt.Println()
}