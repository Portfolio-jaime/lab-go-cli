package cmd

import (
	"fmt"
	"strings"

	"k8s-cli/pkg/kubernetes"
	"k8s-cli/pkg/recommendations"
	"k8s-cli/pkg/table"

	"github.com/spf13/cobra"
)

var recommendCmd = &cobra.Command{
	Use:   "recommend",
	Short: "Get recommendations for cluster optimization and best practices",
	Long:  `Analyze your Kubernetes cluster and provide recommendations for optimization, security, resource management, and best practices.`,
	RunE:  runRecommendCommand,
}

var (
	severityFilter string
	typeFilter     string
)

func init() {
	rootCmd.AddCommand(recommendCmd)
	recommendCmd.Flags().StringVar(&severityFilter, "severity", "", "Filter by severity (High, Medium, Low)")
	recommendCmd.Flags().StringVar(&typeFilter, "type", "", "Filter by type (Resource, Node, Workload, etc.)")
}

func runRecommendCommand(cmd *cobra.Command, args []string) error {
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
	
	client, err := kubernetes.NewClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	fmt.Println("🔍 Analyzing cluster for recommendations...")
	fmt.Println()

	analyzer := recommendations.NewRecommendationAnalyzer(client)
	recs, err := analyzer.AnalyzeCluster()
	if err != nil {
		return fmt.Errorf("failed to analyze cluster: %w", err)
	}

	filteredRecs := filterRecommendations(recs, severityFilter, typeFilter)

	if len(filteredRecs) == 0 {
		fmt.Println("✅ Great! No recommendations found. Your cluster looks well configured!")
		return nil
	}

	fmt.Printf("💡 Found %d recommendations:\n\n", len(filteredRecs))

	showRecommendationsByCategory(filteredRecs)

	return nil
}

func filterRecommendations(recs []recommendations.Recommendation, severity, recType string) []recommendations.Recommendation {
	var filtered []recommendations.Recommendation

	for _, rec := range recs {
		if severity != "" && !strings.EqualFold(rec.Severity, severity) {
			continue
		}
		if recType != "" && !strings.EqualFold(rec.Type, recType) {
			continue
		}
		filtered = append(filtered, rec)
	}

	return filtered
}

func showRecommendationsByCategory(recs []recommendations.Recommendation) {
	categories := make(map[string][]recommendations.Recommendation)
	
	for _, rec := range recs {
		categories[rec.Type] = append(categories[rec.Type], rec)
	}

	priorityOrder := []string{"Security", "Availability", "Resource", "Node", "Workload", "Component", "Monitoring", "Stability", "Resource Management", "Maintenance"}
	
	for _, category := range priorityOrder {
		if categoryRecs, exists := categories[category]; exists {
			showCategoryRecommendations(category, categoryRecs)
			delete(categories, category)
		}
	}

	for category, categoryRecs := range categories {
		showCategoryRecommendations(category, categoryRecs)
	}
}

func showCategoryRecommendations(category string, recs []recommendations.Recommendation) {
	fmt.Printf("📋 %s Recommendations:\n", category)
	
	recTable := table.NewTable([]string{"Severity", "Title", "Description", "Recommended Action"})
	
	for _, rec := range recs {
		recTable.AddRow([]string{
			rec.Severity,
			rec.Title,
			rec.Description,
			rec.Action,
		})
	}
	
	recTable.Render()
	fmt.Println()
}