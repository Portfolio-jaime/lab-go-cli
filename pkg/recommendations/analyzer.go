package recommendations

import (
	"fmt"
	"strconv"
	"strings"

	"k8s-cli/pkg/kubernetes"
)

type Recommendation struct {
	Type        string
	Severity    string
	Title       string
	Description string
	Action      string
	Link        string
}

type RecommendationAnalyzer struct {
	client *kubernetes.Client
}

func NewRecommendationAnalyzer(client *kubernetes.Client) *RecommendationAnalyzer {
	return &RecommendationAnalyzer{
		client: client,
	}
}

func (r *RecommendationAnalyzer) AnalyzeCluster() ([]Recommendation, error) {
	var recommendations []Recommendation

	clusterRecs, err := r.analyzeSimpleClusterResources()
	if err == nil {
		recommendations = append(recommendations, clusterRecs...)
	}

	nodeRecs, err := r.analyzeSimpleNodes()
	if err == nil {
		recommendations = append(recommendations, nodeRecs...)
	}

	podRecs, err := r.analyzeSimplePods()
	if err == nil {
		recommendations = append(recommendations, podRecs...)
	}

	componentRecs, err := r.analyzeComponents()
	if err == nil {
		recommendations = append(recommendations, componentRecs...)
	}

	versionRecs, err := r.analyzeVersions()
	if err == nil {
		recommendations = append(recommendations, versionRecs...)
	}

	return recommendations, nil
}

func (r *RecommendationAnalyzer) analyzeSimpleClusterResources() ([]Recommendation, error) {
	var recommendations []Recommendation

	summary, err := r.client.GetSimpleClusterSummary()
	if err != nil {
		return recommendations, err
	}

	if summary.TotalNodes < 3 {
		recommendations = append(recommendations, Recommendation{
			Type:        "Availability",
			Severity:    "Medium",
			Title:       "Low Node Count",
			Description: fmt.Sprintf("Cluster has only %d nodes, which may impact high availability.", summary.TotalNodes),
			Action:      "Consider adding more nodes for better fault tolerance.",
		})
	}

	if summary.TotalPods > summary.TotalNodes*50 {
		recommendations = append(recommendations, Recommendation{
			Type:        "Resource",
			Severity:    "Medium",
			Title:       "High Pod Density",
			Description: fmt.Sprintf("Cluster has %d pods across %d nodes (avg %.1f pods/node).", summary.TotalPods, summary.TotalNodes, float64(summary.TotalPods)/float64(summary.TotalNodes)),
			Action:      "Consider adding more nodes to reduce pod density and improve performance.",
		})
	}

	return recommendations, nil
}

func (r *RecommendationAnalyzer) analyzeSimpleNodes() ([]Recommendation, error) {
	var recommendations []Recommendation

	nodes, err := r.client.GetSimpleNodesInfo()
	if err != nil {
		return recommendations, err
	}

	notReadyNodes := 0
	oldNodes := 0

	for _, node := range nodes {
		if strings.ToLower(node.Status) != "ready" {
			notReadyNodes++
		}

		if strings.Contains(node.Age, "d") {
			daysStr := strings.TrimSuffix(node.Age, "d")
			if days, err := strconv.Atoi(daysStr); err == nil && days > 365 {
				oldNodes++
			}
		}
	}

	if notReadyNodes > 0 {
		recommendations = append(recommendations, Recommendation{
			Type:        "Availability",
			Severity:    "High",
			Title:       "Nodes Not Ready",
			Description: fmt.Sprintf("%d nodes are not in Ready state.", notReadyNodes),
			Action:      "Investigate and fix the nodes that are not ready.",
		})
	}

	if oldNodes > 0 {
		recommendations = append(recommendations, Recommendation{
			Type:        "Maintenance",
			Severity:    "Low",
			Title:       "Old Nodes Detected",
			Description: fmt.Sprintf("%d nodes are over 1 year old.", oldNodes),
			Action:      "Consider refreshing old nodes for better performance and security.",
		})
	}

	return recommendations, nil
}

func (r *RecommendationAnalyzer) analyzeSimplePods() ([]Recommendation, error) {
	var recommendations []Recommendation

	pods, err := r.client.GetSimplePodsInfo("")
	if err != nil {
		return recommendations, err
	}

	failedPods := 0
	highRestartPods := 0
	terminatingPods := 0

	for _, pod := range pods {
		status := strings.ToLower(pod.Status)
		if strings.Contains(status, "failed") || strings.Contains(status, "error") {
			failedPods++
		}

		if strings.Contains(status, "terminating") {
			terminatingPods++
		}

		if restarts, err := strconv.Atoi(pod.Restarts); err == nil && restarts > 10 {
			highRestartPods++
		}
	}

	if failedPods > 0 {
		recommendations = append(recommendations, Recommendation{
			Type:        "Workload",
			Severity:    "Medium",
			Title:       "Failed Pods Detected",
			Description: fmt.Sprintf("%d pods are in failed state.", failedPods),
			Action:      "Investigate and fix failed pods, check logs for root cause.",
		})
	}

	if highRestartPods > 0 {
		recommendations = append(recommendations, Recommendation{
			Type:        "Stability",
			Severity:    "Medium",
			Title:       "High Restart Count Pods",
			Description: fmt.Sprintf("%d pods have more than 10 restarts.", highRestartPods),
			Action:      "Investigate pods with high restart counts for stability issues.",
		})
	}

	if terminatingPods > 5 {
		recommendations = append(recommendations, Recommendation{
			Type:        "Workload",
			Severity:    "Low",
			Title:       "Many Terminating Pods",
			Description: fmt.Sprintf("%d pods are stuck in terminating state.", terminatingPods),
			Action:      "Check for stuck terminating pods and force delete if necessary.",
		})
	}

	return recommendations, nil
}

func (r *RecommendationAnalyzer) analyzeComponents() ([]Recommendation, error) {
	var recommendations []Recommendation

	components, err := r.client.GetInstalledComponents()
	if err != nil {
		return recommendations, err
	}

	hasMetricsServer := false
	notReadyComponents := 0

	for _, comp := range components {
		if strings.Contains(strings.ToLower(comp.Name), "metrics-server") {
			hasMetricsServer = true
		}

		if strings.Contains(strings.ToLower(comp.Status), "not ready") {
			notReadyComponents++
		}
	}

	if !hasMetricsServer {
		recommendations = append(recommendations, Recommendation{
			Type:        "Monitoring",
			Severity:    "Medium",
			Title:       "Metrics Server Not Found",
			Description: "Metrics server is not detected in the cluster.",
			Action:      "Install metrics-server for resource monitoring capabilities.",
		})
	}

	if notReadyComponents > 0 {
		recommendations = append(recommendations, Recommendation{
			Type:        "Component",
			Severity:    "Medium",
			Title:       "Components Not Ready",
			Description: fmt.Sprintf("%d components are not in ready state.", notReadyComponents),
			Action:      "Check and fix components that are not ready.",
		})
	}

	return recommendations, nil
}

func (r *RecommendationAnalyzer) analyzeVersions() ([]Recommendation, error) {
	var recommendations []Recommendation

	clusterInfo, err := r.client.GetClusterVersion()
	if err != nil {
		return recommendations, err
	}

	majorVersion, _ := strconv.Atoi(clusterInfo.Major)
	minorVersion, _ := strconv.Atoi(clusterInfo.Minor)

	if majorVersion == 1 && minorVersion < 25 {
		recommendations = append(recommendations, Recommendation{
			Type:        "Security",
			Severity:    "High",
			Title:       "Outdated Kubernetes Version",
			Description: fmt.Sprintf("Kubernetes version %s.%s is outdated and may have security vulnerabilities.", clusterInfo.Major, clusterInfo.Minor),
			Action:      "Plan to upgrade to a supported Kubernetes version (1.25+).",
		})
	} else if majorVersion == 1 && minorVersion < 27 {
		recommendations = append(recommendations, Recommendation{
			Type:        "Maintenance",
			Severity:    "Low",
			Title:       "Consider Version Upgrade",
			Description: fmt.Sprintf("Kubernetes version %s.%s could be updated to get latest features.", clusterInfo.Major, clusterInfo.Minor),
			Action:      "Consider upgrading to a newer version for better features and support.",
		})
	}

	return recommendations, nil
}

