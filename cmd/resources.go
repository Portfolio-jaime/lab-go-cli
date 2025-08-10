package cmd

import (
	"fmt"

	"k8s-cli/pkg/kubernetes"
	"k8s-cli/pkg/table"

	"github.com/spf13/cobra"
)

var resourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "Show cluster resource consumption and utilization",
	Long:  `Display detailed information about resource consumption across the cluster, nodes, and pods with utilization percentages.`,
	RunE:  runResourcesCommand,
}

var (
	namespace string
	nodeOnly  bool
	podOnly   bool
)

func init() {
	rootCmd.AddCommand(resourcesCmd)
	resourcesCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace to filter pods (default: all namespaces)")
	resourcesCmd.Flags().BoolVar(&nodeOnly, "nodes", false, "Show only node resources")
	resourcesCmd.Flags().BoolVar(&podOnly, "pods", false, "Show only pod resources")
}

func runResourcesCommand(cmd *cobra.Command, args []string) error {
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
	
	client, err := kubernetes.NewClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	fmt.Println("📊 Analyzing cluster resources...")
	fmt.Println()

	if !nodeOnly && !podOnly {
		if err := showClusterSummary(client); err != nil {
			return err
		}
	}

	if !podOnly {
		if err := showNodesResources(client); err != nil {
			return err
		}
	}

	if !nodeOnly {
		if err := showPodsResources(client); err != nil {
			return err
		}
	}

	return nil
}

func showClusterSummary(client *kubernetes.Client) error {
	summary, err := client.GetSimpleClusterSummary()
	if err != nil {
		return fmt.Errorf("failed to get cluster summary: %w", err)
	}

	fmt.Println("🏢 Cluster Resource Summary:")
	summaryTable := table.NewTable([]string{"Metric", "Value"})
	summaryTable.AddRow([]string{"Total Nodes", fmt.Sprintf("%d", summary.TotalNodes)})
	summaryTable.AddRow([]string{"Total Pods", fmt.Sprintf("%d", summary.TotalPods)})
	summaryTable.AddRow([]string{"CPU Capacity", summary.TotalCPUCapacity + " cores"})
	summaryTable.AddRow([]string{"Memory Capacity", summary.TotalMemCapacity})
	summaryTable.AddRow([]string{"Note", "Install metrics-server for usage data"})
	
	summaryTable.Render()
	fmt.Println()

	return nil
}

func showNodesResources(client *kubernetes.Client) error {
	nodes, err := client.GetSimpleNodesInfo()
	if err != nil {
		return fmt.Errorf("failed to get nodes info: %w", err)
	}

	if len(nodes) == 0 {
		fmt.Println("No nodes found in the cluster.")
		return nil
	}

	fmt.Println("🖥️  Node Resources:")
	nodeTable := table.NewTable([]string{"Node", "Status", "Role", "Age", "Version", "CPU Capacity", "Memory Capacity"})
	
	for _, node := range nodes {
		nodeTable.AddRow([]string{
			node.Name,
			node.Status,
			node.Roles,
			node.Age,
			node.Version,
			node.CPUCapacity,
			node.MemoryCapacity,
		})
	}
	
	nodeTable.Render()
	fmt.Println()

	return nil
}

func showPodsResources(client *kubernetes.Client) error {
	pods, err := client.GetSimplePodsInfo(namespace)
	if err != nil {
		return fmt.Errorf("failed to get pods info: %w", err)
	}

	if len(pods) == 0 {
		namespaceText := "cluster"
		if namespace != "" {
			namespaceText = fmt.Sprintf("namespace '%s'", namespace)
		}
		fmt.Printf("No pods found in %s.\n", namespaceText)
		return nil
	}

	namespaceText := "All Namespaces"
	if namespace != "" {
		namespaceText = fmt.Sprintf("Namespace: %s", namespace)
	}

	fmt.Printf("🚀 Pod Resources (%s):\n", namespaceText)
	
	if len(pods) > 20 {
		fmt.Printf("Showing first 20 pods out of %d total pods. Use --namespace to filter.\n", len(pods))
		pods = pods[:20]
	}

	podTable := table.NewTable([]string{"Pod", "Namespace", "Status", "Restarts", "Age", "Node"})
	
	for _, pod := range pods {
		podTable.AddRow([]string{
			pod.Name,
			pod.Namespace,
			pod.Status,
			pod.Restarts,
			pod.Age,
			pod.Node,
		})
	}
	
	podTable.Render()
	fmt.Println()

	return nil
}

