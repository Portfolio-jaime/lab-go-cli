package cmd

import (
	"fmt"
	"strings"

	"k8s-cli/pkg/kubernetes"
	"k8s-cli/pkg/table"

	"github.com/spf13/cobra"
)

var workloadCmd = &cobra.Command{
	Use:   "workload",
	Short: "Analyze workload health and performance across the cluster",
	Long:  `Comprehensive analysis of deployments, statefulsets, daemonsets, and pods to identify health issues, performance problems, and optimization opportunities.`,
	RunE:  runWorkloadCommand,
}

var (
	showWorkloadDeployments  bool
	showWorkloadStatefulSets bool
	showWorkloadDaemonSets   bool
	showWorkloadPods         bool
	showWorkloadSummary      bool
	workloadNamespace        string
	onlyUnhealthy            bool
)

func init() {
	rootCmd.AddCommand(workloadCmd)
	workloadCmd.Flags().BoolVar(&showWorkloadDeployments, "deployments", true, "Show deployment analysis")
	workloadCmd.Flags().BoolVar(&showWorkloadStatefulSets, "statefulsets", true, "Show statefulset analysis")
	workloadCmd.Flags().BoolVar(&showWorkloadDaemonSets, "daemonsets", true, "Show daemonset analysis")
	workloadCmd.Flags().BoolVar(&showWorkloadPods, "pods", false, "Show detailed pod analysis")
	workloadCmd.Flags().BoolVar(&showWorkloadSummary, "summary", true, "Show workload summary")
	workloadCmd.Flags().StringVarP(&workloadNamespace, "namespace", "n", "", "Namespace to analyze (empty for all)")
	workloadCmd.Flags().BoolVar(&onlyUnhealthy, "unhealthy-only", false, "Show only unhealthy workloads")
}

func runWorkloadCommand(cmd *cobra.Command, args []string) error {
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")

	client, err := kubernetes.NewClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	fmt.Println("🔍 Workload Health Analysis")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	analysis, err := client.GetWorkloadAnalysis(workloadNamespace)
	if err != nil {
		return fmt.Errorf("failed to get workload analysis: %w", err)
	}

	if showWorkloadSummary {
		showWorkloadOverview(&analysis.WorkloadSummary)
	}

	if showWorkloadDeployments {
		showDeploymentAnalysis(analysis.DeploymentAnalysis)
	}

	if showWorkloadStatefulSets {
		showStatefulSetAnalysis(analysis.StatefulSetAnalysis)
	}

	if showWorkloadDaemonSets {
		showDaemonSetAnalysis(analysis.DaemonSetAnalysis)
	}

	if showWorkloadPods {
		showPodsAnalysis(analysis.PodAnalysis)
	}

	return nil
}

func showWorkloadOverview(summary *kubernetes.WorkloadSummary) {
	fmt.Println("📊 WORKLOAD SUMMARY")
	fmt.Println(strings.Repeat("-", 40))

	summaryTable := table.NewTable([]string{"Workload Type", "Total", "Healthy", "Health Rate"})

	deploymentRate := "N/A"
	if summary.TotalDeployments > 0 {
		deploymentRate = fmt.Sprintf("%.1f%%", float64(summary.HealthyDeployments)/float64(summary.TotalDeployments)*100)
	}
	summaryTable.AddRow([]string{"Deployments", fmt.Sprintf("%d", summary.TotalDeployments), fmt.Sprintf("%d", summary.HealthyDeployments), deploymentRate})

	statefulSetRate := "N/A"
	if summary.TotalStatefulSets > 0 {
		statefulSetRate = fmt.Sprintf("%.1f%%", float64(summary.HealthyStatefulSets)/float64(summary.TotalStatefulSets)*100)
	}
	summaryTable.AddRow([]string{"StatefulSets", fmt.Sprintf("%d", summary.TotalStatefulSets), fmt.Sprintf("%d", summary.HealthyStatefulSets), statefulSetRate})

	daemonSetRate := "N/A"
	if summary.TotalDaemonSets > 0 {
		daemonSetRate = fmt.Sprintf("%.1f%%", float64(summary.HealthyDaemonSets)/float64(summary.TotalDaemonSets)*100)
	}
	summaryTable.AddRow([]string{"DaemonSets", fmt.Sprintf("%d", summary.TotalDaemonSets), fmt.Sprintf("%d", summary.HealthyDaemonSets), daemonSetRate})

	podRate := "N/A"
	if summary.TotalPods > 0 {
		podRate = fmt.Sprintf("%.1f%%", float64(summary.HealthyPods)/float64(summary.TotalPods)*100)
	}
	summaryTable.AddRow([]string{"Pods", fmt.Sprintf("%d", summary.TotalPods), fmt.Sprintf("%d", summary.HealthyPods), podRate})

	summaryTable.Render()

	overallTable := table.NewTable([]string{"Metric", "Value"})
	overallTable.AddRow([]string{"Overall Health Score", fmt.Sprintf("%d/100", summary.OverallHealthScore)})
	overallTable.AddRow([]string{"Critical Issues", fmt.Sprintf("%d", summary.CriticalIssues)})

	healthStatus := "🟢 Excellent"
	if summary.OverallHealthScore < 80 {
		healthStatus = "🟡 Good"
	}
	if summary.OverallHealthScore < 60 {
		healthStatus = "🟠 Fair"
	}
	if summary.OverallHealthScore < 40 {
		healthStatus = "🔴 Poor"
	}
	overallTable.AddRow([]string{"Overall Status", healthStatus})
	overallTable.Render()

	fmt.Println()
}

func showDeploymentAnalysis(deployments []kubernetes.DeploymentHealth) {
	if len(deployments) == 0 {
		return
	}

	fmt.Println("🚀 DEPLOYMENT ANALYSIS")
	fmt.Println(strings.Repeat("-", 40))

	deploymentTable := table.NewTable([]string{"Name", "Namespace", "Replicas", "Status", "Health", "Issues"})

	for _, deploy := range deployments {
		if onlyUnhealthy && deploy.HealthScore >= 80 {
			continue
		}

		replicas := fmt.Sprintf("%d/%d", deploy.ReadyReplicas, deploy.Replicas)
		if deploy.UnavailableReplicas > 0 {
			replicas += fmt.Sprintf(" (-%d)", deploy.UnavailableReplicas)
		}

		status := deploy.Status
		if deploy.Status == "Critical" {
			status = "🔴 " + status
		} else if deploy.Status == "Warning" {
			status = "🟡 " + status
		} else {
			status = "🟢 " + status
		}

		healthScore := fmt.Sprintf("%d/100", deploy.HealthScore)
		if deploy.HealthScore < 60 {
			healthScore += " ⚠️"
		}

		issues := fmt.Sprintf("%d", len(deploy.Issues))
		if len(deploy.Issues) > 0 {
			issues += " ⚠️"
		}

		deploymentTable.AddRow([]string{
			deploy.Name,
			deploy.Namespace,
			replicas,
			status,
			healthScore,
			issues,
		})
	}
	deploymentTable.Render()

	if len(deployments) > 0 && !onlyUnhealthy {
		unhealthyCount := 0
		for _, deploy := range deployments {
			if deploy.HealthScore < 80 {
				unhealthyCount++
			}
		}

		if unhealthyCount > 0 {
			fmt.Printf("\n⚠️  %d deployments need attention. Use --unhealthy-only for details.\n", unhealthyCount)
		}
	}

	fmt.Println()
}

func showStatefulSetAnalysis(statefulSets []kubernetes.StatefulSetHealth) {
	if len(statefulSets) == 0 {
		return
	}

	fmt.Println("💾 STATEFULSET ANALYSIS")
	fmt.Println(strings.Repeat("-", 40))

	ssTable := table.NewTable([]string{"Name", "Namespace", "Replicas", "Status", "Health", "Issues"})

	for _, ss := range statefulSets {
		if onlyUnhealthy && ss.HealthScore >= 80 {
			continue
		}

		replicas := fmt.Sprintf("%d/%d", ss.ReadyReplicas, ss.Replicas)

		status := ss.Status
		if ss.Status == "Critical" {
			status = "🔴 " + status
		} else if ss.Status == "Warning" {
			status = "🟡 " + status
		} else {
			status = "🟢 " + status
		}

		healthScore := fmt.Sprintf("%d/100", ss.HealthScore)
		if ss.HealthScore < 60 {
			healthScore += " ⚠️"
		}

		issues := fmt.Sprintf("%d", len(ss.Issues))
		if len(ss.Issues) > 0 {
			issues += " ⚠️"
		}

		ssTable.AddRow([]string{
			ss.Name,
			ss.Namespace,
			replicas,
			status,
			healthScore,
			issues,
		})
	}
	ssTable.Render()
	fmt.Println()
}

func showDaemonSetAnalysis(daemonSets []kubernetes.DaemonSetHealth) {
	if len(daemonSets) == 0 {
		return
	}

	fmt.Println("⚙️  DAEMONSET ANALYSIS")
	fmt.Println(strings.Repeat("-", 40))

	dsTable := table.NewTable([]string{"Name", "Namespace", "Scheduled", "Ready", "Status", "Health", "Issues"})

	for _, ds := range daemonSets {
		if onlyUnhealthy && ds.HealthScore >= 80 {
			continue
		}

		scheduled := fmt.Sprintf("%d", ds.CurrentNumberScheduled)
		ready := fmt.Sprintf("%d/%d", ds.NumberReady, ds.DesiredNumberScheduled)
		if ds.NumberUnavailable > 0 {
			ready += fmt.Sprintf(" (-%d)", ds.NumberUnavailable)
		}

		status := ds.Status
		if ds.Status == "Critical" {
			status = "🔴 " + status
		} else if ds.Status == "Warning" {
			status = "🟡 " + status
		} else {
			status = "🟢 " + status
		}

		healthScore := fmt.Sprintf("%d/100", ds.HealthScore)
		if ds.HealthScore < 60 {
			healthScore += " ⚠️"
		}

		issues := fmt.Sprintf("%d", len(ds.Issues))
		if len(ds.Issues) > 0 {
			issues += " ⚠️"
		}

		dsTable.AddRow([]string{
			ds.Name,
			ds.Namespace,
			scheduled,
			ready,
			status,
			healthScore,
			issues,
		})
	}
	dsTable.Render()
	fmt.Println()
}

func showPodsAnalysis(pods []kubernetes.PodHealth) {
	if len(pods) == 0 {
		return
	}

	fmt.Println("🚀 POD ANALYSIS")
	fmt.Println(strings.Repeat("-", 40))

	podTable := table.NewTable([]string{"Name", "Namespace", "Status", "Restarts", "Health", "Issues", "Node"})

	displayed := 0
	for _, pod := range pods {
		if onlyUnhealthy && pod.HealthScore >= 80 {
			continue
		}

		if displayed >= 20 {
			break
		}

		status := pod.Status
		if pod.Status != "Running" {
			status = "🔴 " + status
		} else {
			status = "🟢 " + status
		}

		restarts := fmt.Sprintf("%d", pod.RestartCount)
		if pod.RestartCount > 5 {
			restarts += " ⚠️"
		}

		healthScore := fmt.Sprintf("%d/100", pod.HealthScore)
		if pod.HealthScore < 60 {
			healthScore += " ⚠️"
		}

		issues := fmt.Sprintf("%d", len(pod.Issues))
		if len(pod.Issues) > 0 {
			issues += " ⚠️"
		}

		podTable.AddRow([]string{
			pod.Name,
			pod.Namespace,
			status,
			restarts,
			healthScore,
			issues,
			pod.Node,
		})
		displayed++
	}
	podTable.Render()

	if len(pods) > 20 {
		fmt.Printf("... and %d more pods. Use --unhealthy-only to focus on problematic pods.\n", len(pods)-20)
	}

	fmt.Println()
}
