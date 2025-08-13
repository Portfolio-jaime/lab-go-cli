package cmd

import (
	"fmt"
	"strings"
	"time"

	"k8s-cli/pkg/kubernetes"
	"k8s-cli/pkg/table"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Analyze cluster events, logs, and identify critical issues",
	Long:  `Provide comprehensive analysis of cluster events, error patterns, and critical issues affecting cluster health.`,
	RunE:  runLogsCommand,
}

var (
	timeWindow         int
	showLogsCritical   bool
	showLogsWarnings   bool
	showLogsPatterns   bool
	showLogsResourceEvents bool
	showLogsSecurityEvents bool
	showLogsPodAnalysis    bool
	logsNamespace          string
)

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().IntVar(&timeWindow, "hours", 24, "Time window in hours to analyze events")
	logsCmd.Flags().BoolVar(&showLogsCritical, "critical", true, "Show critical events")
	logsCmd.Flags().BoolVar(&showLogsWarnings, "warnings", true, "Show warning events")
	logsCmd.Flags().BoolVar(&showLogsPatterns, "patterns", true, "Show error patterns")
	logsCmd.Flags().BoolVar(&showLogsResourceEvents, "resource-events", true, "Show resource-related events")
	logsCmd.Flags().BoolVar(&showLogsSecurityEvents, "security-events", true, "Show security-related events")
	logsCmd.Flags().BoolVar(&showLogsPodAnalysis, "pod-analysis", false, "Show detailed pod log analysis")
	logsCmd.Flags().StringVarP(&logsNamespace, "namespace", "n", "", "Namespace to analyze (empty for all)")
}

func runLogsCommand(cmd *cobra.Command, args []string) error {
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
	
	client, err := kubernetes.NewClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	fmt.Printf("📋 Cluster Events & Logs Analysis (Last %d hours)\n", timeWindow)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	analysis, err := client.GetLogAnalysis(logsNamespace, timeWindow)
	if err != nil {
		return fmt.Errorf("failed to get log analysis: %w", err)
	}

	showEventsOverview(analysis)

	if showLogsCritical {
		showLogsCriticalEvents(analysis.CriticalEvents)
	}

	if showLogsWarnings {
		showWarningEvents(analysis.WarningEvents)
	}

	if showLogsPatterns {
		showErrorPatterns(analysis.ErrorPatterns)
	}

	if showLogsResourceEvents {
		showResourceEventsAnalysis(analysis.ResourceEvents)
	}

	if showLogsSecurityEvents {
		showSecurityEventsAnalysis(analysis.SecurityEvents)
	}

	if showLogsPodAnalysis {
		if err := showPodLogsAnalysis(client, logsNamespace); err != nil {
			fmt.Printf("Warning: Could not retrieve pod logs analysis: %v\n", err)
		}
	}

	return nil
}

func showEventsOverview(analysis *kubernetes.LogAnalysis) {
	fmt.Println("📊 EVENTS OVERVIEW")
	fmt.Println(strings.Repeat("-", 40))

	overviewTable := table.NewTable([]string{"Category", "Count"})
	overviewTable.AddRow([]string{"Critical Events", fmt.Sprintf("%d", len(analysis.CriticalEvents))})
	overviewTable.AddRow([]string{"Warning Events", fmt.Sprintf("%d", len(analysis.WarningEvents))})
	overviewTable.AddRow([]string{"Error Patterns", fmt.Sprintf("%d", len(analysis.ErrorPatterns))})
	overviewTable.AddRow([]string{"Resource Events", fmt.Sprintf("%d", len(analysis.ResourceEvents))})
	overviewTable.AddRow([]string{"Security Events", fmt.Sprintf("%d", len(analysis.SecurityEvents))})
	overviewTable.Render()
	fmt.Println()
}

func showLogsCriticalEvents(events []kubernetes.ClusterEvent) {
	if len(events) == 0 {
		fmt.Println("✅ No critical events found!")
		fmt.Println()
		return
	}

	fmt.Println("🚨 CRITICAL EVENTS")
	fmt.Println(strings.Repeat("-", 40))

	criticalTable := table.NewTable([]string{"Time", "Object", "Reason", "Message", "Count"})
	for i, event := range events {
		if i >= 10 {
			break
		}
		
		timeStr := event.LastTime.Format("15:04:05")
		if time.Since(event.LastTime) > 24*time.Hour {
			timeStr = event.LastTime.Format("01-02 15:04")
		}

		message := event.Message
		if len(message) > 50 {
			message = message[:47] + "..."
		}

		criticalTable.AddRow([]string{
			timeStr,
			event.Object,
			event.Reason,
			message,
			fmt.Sprintf("%d", event.Count),
		})
	}
	criticalTable.Render()
	
	if len(events) > 10 {
		fmt.Printf("... and %d more critical events\n", len(events)-10)
	}
	fmt.Println()
}

func showWarningEvents(events []kubernetes.ClusterEvent) {
	if len(events) == 0 {
		fmt.Println("✅ No warning events found!")
		fmt.Println()
		return
	}

	fmt.Println("⚠️  WARNING EVENTS")
	fmt.Println(strings.Repeat("-", 40))

	warningTable := table.NewTable([]string{"Time", "Object", "Reason", "Message", "Count"})
	for i, event := range events {
		if i >= 10 {
			break
		}
		
		timeStr := event.LastTime.Format("15:04:05")
		if time.Since(event.LastTime) > 24*time.Hour {
			timeStr = event.LastTime.Format("01-02 15:04")
		}

		message := event.Message
		if len(message) > 50 {
			message = message[:47] + "..."
		}

		warningTable.AddRow([]string{
			timeStr,
			event.Object,
			event.Reason,
			message,
			fmt.Sprintf("%d", event.Count),
		})
	}
	warningTable.Render()
	
	if len(events) > 10 {
		fmt.Printf("... and %d more warning events\n", len(events)-10)
	}
	fmt.Println()
}

func showErrorPatterns(patterns []kubernetes.ErrorPattern) {
	if len(patterns) == 0 {
		fmt.Println("✅ No significant error patterns detected!")
		fmt.Println()
		return
	}

	fmt.Println("🔍 ERROR PATTERNS")
	fmt.Println(strings.Repeat("-", 40))

	patternTable := table.NewTable([]string{"Pattern", "Count", "Severity", "Last Seen", "Recommendation"})
	for i, pattern := range patterns {
		if i >= 8 {
			break
		}
		
		lastSeenStr := pattern.LastSeen.Format("15:04:05")
		if time.Since(pattern.LastSeen) > 24*time.Hour {
			lastSeenStr = pattern.LastSeen.Format("01-02 15:04")
		}

		recommendation := pattern.Recommendation
		if len(recommendation) > 40 {
			recommendation = recommendation[:37] + "..."
		}

		severity := pattern.Severity
		if pattern.Severity == "Critical" {
			severity = "🔴 " + severity
		} else if pattern.Severity == "Warning" {
			severity = "🟡 " + severity
		}

		patternTable.AddRow([]string{
			pattern.Pattern,
			fmt.Sprintf("%d", pattern.Count),
			severity,
			lastSeenStr,
			recommendation,
		})
	}
	patternTable.Render()
	fmt.Println()
}

func showResourceEventsAnalysis(events []kubernetes.ResourceEvent) {
	if len(events) == 0 {
		fmt.Println("✅ No resource-related events found!")
		fmt.Println()
		return
	}

	fmt.Println("📈 RESOURCE EVENTS")
	fmt.Println(strings.Repeat("-", 40))

	resourceTable := table.NewTable([]string{"Type", "Resource", "Event", "Impact", "Time"})
	for i, event := range events {
		if i >= 10 {
			break
		}
		
		timeStr := event.Timestamp.Format("15:04:05")
		if time.Since(event.Timestamp) > 24*time.Hour {
			timeStr = event.Timestamp.Format("01-02 15:04")
		}

		impact := event.Impact
		if event.Impact == "High" {
			impact = "🔴 " + impact
		} else if event.Impact == "Medium" {
			impact = "🟡 " + impact
		} else {
			impact = "🟢 " + impact
		}

		resourceTable.AddRow([]string{
			event.Type,
			event.ResourceName,
			event.Event,
			impact,
			timeStr,
		})
	}
	resourceTable.Render()
	fmt.Println()
}

func showSecurityEventsAnalysis(events []kubernetes.SecurityEvent) {
	if len(events) == 0 {
		fmt.Println("✅ No security-related events found!")
		fmt.Println()
		return
	}

	fmt.Println("🔒 SECURITY EVENTS")
	fmt.Println(strings.Repeat("-", 40))

	securityTable := table.NewTable([]string{"Risk Level", "Object", "Description", "Action", "Time"})
	for i, event := range events {
		if i >= 10 {
			break
		}
		
		timeStr := event.Timestamp.Format("15:04:05")
		if time.Since(event.Timestamp) > 24*time.Hour {
			timeStr = event.Timestamp.Format("01-02 15:04")
		}

		riskLevel := event.RiskLevel
		if event.RiskLevel == "High" {
			riskLevel = "🔴 " + riskLevel
		} else if event.RiskLevel == "Medium" {
			riskLevel = "🟡 " + riskLevel
		} else {
			riskLevel = "🟢 " + riskLevel
		}

		description := event.Description
		if len(description) > 30 {
			description = description[:27] + "..."
		}

		action := event.Action
		if len(action) > 30 {
			action = action[:27] + "..."
		}

		securityTable.AddRow([]string{
			riskLevel,
			event.Object,
			description,
			action,
			timeStr,
		})
	}
	securityTable.Render()
	fmt.Println()
}

func showPodLogsAnalysis(client *kubernetes.Client, namespace string) error {
	fmt.Println("🚀 POD LOG ANALYSIS")
	fmt.Println(strings.Repeat("-", 40))

	summaries, err := client.GetPodLogsAnalysis(namespace)
	if err != nil {
		return err
	}

	if len(summaries) == 0 {
		fmt.Println("✅ No problematic pods found!")
		fmt.Println()
		return nil
	}

	podTable := table.NewTable([]string{"Pod", "Namespace", "Status", "Errors", "Warnings", "Last Restart"})
	for i, summary := range summaries {
		if i >= 15 {
			break
		}
		
		lastRestartStr := "Never"
		if !summary.LastRestart.IsZero() {
			lastRestartStr = summary.LastRestart.Format("01-02 15:04")
		}

		status := summary.Status
		if summary.ErrorCount > 0 {
			status += " ⚠️"
		}

		podTable.AddRow([]string{
			summary.PodName,
			summary.Namespace,
			status,
			fmt.Sprintf("%d", summary.ErrorCount),
			fmt.Sprintf("%d", summary.WarningCount),
			lastRestartStr,
		})
	}
	podTable.Render()

	if len(summaries) > 15 {
		fmt.Printf("... and %d more pods with issues\n", len(summaries)-15)
	}
	fmt.Println()

	return nil
}