package cmd

import (
	"fmt"

	"k8s-cli/pkg/kubernetes"
	"k8s-cli/pkg/table"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show Kubernetes cluster version and installed components",
	Long:  `Display detailed information about the Kubernetes cluster version and detect installed components like metrics-server, ArgoCD, Kuma, etc.`,
	RunE:  runVersionCommand,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersionCommand(cmd *cobra.Command, args []string) error {
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
	
	client, err := kubernetes.NewClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	fmt.Println("🔍 Analyzing Kubernetes cluster...")
	fmt.Println()

	clusterInfo, err := client.GetClusterVersion()
	if err != nil {
		return fmt.Errorf("failed to get cluster version: %w", err)
	}

	fmt.Println("📊 Cluster Version Information:")
	versionTable := table.NewTable([]string{"Property", "Value"})
	versionTable.AddRow([]string{"Kubernetes Version", clusterInfo.GitVersion})
	versionTable.AddRow([]string{"Major Version", clusterInfo.Major})
	versionTable.AddRow([]string{"Minor Version", clusterInfo.Minor})
	versionTable.AddRow([]string{"Platform", clusterInfo.Platform})
	versionTable.AddRow([]string{"Build Date", clusterInfo.BuildDate})
	versionTable.AddRow([]string{"Go Version", clusterInfo.GoVersion})
	versionTable.AddRow([]string{"Git Commit", clusterInfo.GitCommit})
	versionTable.Render()

	fmt.Println()
	fmt.Println("🔧 Installed Components:")
	
	components, err := client.GetInstalledComponents()
	if err != nil {
		return fmt.Errorf("failed to get installed components: %w", err)
	}

	if len(components) == 0 {
		fmt.Println("No common components detected in the cluster.")
		return nil
	}

	componentTable := table.NewTable([]string{"Component", "Namespace", "Status", "Version", "Ready"})
	for _, comp := range components {
		componentTable.AddRow([]string{comp.Name, comp.Namespace, comp.Status, comp.Version, comp.Ready})
	}
	componentTable.Render()

	return nil
}