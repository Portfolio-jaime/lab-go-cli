package kubernetes

import (
	"fmt"
	"sort"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CostAnalysis struct {
	TotalMonthlyCost       float64
	NodeCosts              []NodeCost
	NamespaceCosts         []NamespaceCost
	UnderutilizedResources []UnderutilizedResource
	CostOptimizations      []CostOptimization
}

type NodeCost struct {
	Name           string
	Type           string
	MonthlyCost    float64
	CPUCapacity    string
	MemoryCapacity string
	CPUUtilization float64
	MemUtilization float64
	Efficiency     string
}

type NamespaceCost struct {
	Name           string
	MonthlyCost    float64
	CPURequests    string
	MemoryRequests string
	PodsCount      int
	CostPerPod     float64
}

type UnderutilizedResource struct {
	Type             string
	Name             string
	Namespace        string
	CPUWaste         string
	MemoryWaste      string
	EstimatedSavings float64
	Recommendation   string
}

type CostOptimization struct {
	Type             string
	Description      string
	PotentialSavings float64
	Priority         string
	Action           string
}

// AWS EC2 pricing estimates (simplified)
var nodeTypeCosts = map[string]float64{
	"t3.micro":  0.0104 * 24 * 30, // $7.49/month
	"t3.small":  0.0208 * 24 * 30, // $14.98/month
	"t3.medium": 0.0416 * 24 * 30, // $29.97/month
	"t3.large":  0.0832 * 24 * 30, // $59.94/month
	"t3.xlarge": 0.1664 * 24 * 30, // $119.88/month
	"m5.large":  0.096 * 24 * 30,  // $69.12/month
	"m5.xlarge": 0.192 * 24 * 30,  // $138.24/month
	"c5.large":  0.085 * 24 * 30,  // $61.20/month
	"c5.xlarge": 0.17 * 24 * 30,   // $122.40/month
	"default":   0.10 * 24 * 30,   // $72/month (default estimate)
}

func (c *Client) GetCostAnalysis() (*CostAnalysis, error) {
	nodes, err := c.Clientset.CoreV1().Nodes().List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	var nodeMetrics []NodeMetrics
	if metrics, err := c.GetRealTimeNodeMetrics(); err == nil {
		nodeMetrics = metrics
	}

	nodeCosts := c.calculateNodeCosts(nodes.Items, nodeMetrics)

	namespaceCosts, err := c.calculateNamespaceCosts()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate namespace costs: %w", err)
	}

	underutilized, err := c.findUnderutilizedResources()
	if err != nil {
		return nil, fmt.Errorf("failed to find underutilized resources: %w", err)
	}

	optimizations := c.generateCostOptimizations(nodeCosts, namespaceCosts, underutilized)

	totalCost := 0.0
	for _, nc := range nodeCosts {
		totalCost += nc.MonthlyCost
	}

	return &CostAnalysis{
		TotalMonthlyCost:       totalCost,
		NodeCosts:              nodeCosts,
		NamespaceCosts:         namespaceCosts,
		UnderutilizedResources: underutilized,
		CostOptimizations:      optimizations,
	}, nil
}

func (c *Client) calculateNodeCosts(nodes []corev1.Node, nodeMetrics []NodeMetrics) []NodeCost {
	metricsMap := make(map[string]NodeMetrics)
	for _, metric := range nodeMetrics {
		metricsMap[metric.Name] = metric
	}

	var nodeCosts []NodeCost
	for _, node := range nodes {
		nodeType := c.extractNodeType(&node)
		cost := c.getNodeCost(nodeType)

		var cpuUtil, memUtil float64
		var efficiency string

		if metric, exists := metricsMap[node.Name]; exists {
			cpuUtil = metric.CPUUsagePercent
			memUtil = metric.MemoryUsagePercent
			efficiency = c.calculateEfficiency(cpuUtil, memUtil)
		} else {
			efficiency = "No metrics"
		}

		cpuCapacity := node.Status.Capacity[corev1.ResourceCPU]
		memCapacity := node.Status.Capacity[corev1.ResourceMemory]

		nodeCosts = append(nodeCosts, NodeCost{
			Name:           node.Name,
			Type:           nodeType,
			MonthlyCost:    cost,
			CPUCapacity:    formatCPU(cpuCapacity.MilliValue()),
			MemoryCapacity: formatBytes(memCapacity.Value()),
			CPUUtilization: cpuUtil,
			MemUtilization: memUtil,
			Efficiency:     efficiency,
		})
	}

	return nodeCosts
}

func (c *Client) calculateNamespaceCosts() ([]NamespaceCost, error) {
	namespaces, err := c.Clientset.CoreV1().Namespaces().List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespaceCosts []NamespaceCost
	for _, ns := range namespaces.Items {
		if ns.Name == "kube-system" || ns.Name == "kube-public" || ns.Name == "kube-node-lease" {
			continue
		}

		pods, err := c.Clientset.CoreV1().Pods(ns.Name).List(c.Context, metav1.ListOptions{})
		if err != nil {
			continue
		}

		var totalCPURequests, totalMemRequests int64
		podsCount := len(pods.Items)

		for _, pod := range pods.Items {
			cpuReq, memReq := getPodResourceRequests(&pod)
			totalCPURequests += cpuReq
			totalMemRequests += memReq
		}

		estimatedCost := c.estimateNamespaceCost(totalCPURequests, totalMemRequests)
		costPerPod := 0.0
		if podsCount > 0 {
			costPerPod = estimatedCost / float64(podsCount)
		}

		namespaceCosts = append(namespaceCosts, NamespaceCost{
			Name:           ns.Name,
			MonthlyCost:    estimatedCost,
			CPURequests:    formatCPU(totalCPURequests),
			MemoryRequests: formatBytes(totalMemRequests),
			PodsCount:      podsCount,
			CostPerPod:     costPerPod,
		})
	}

	sort.Slice(namespaceCosts, func(i, j int) bool {
		return namespaceCosts[i].MonthlyCost > namespaceCosts[j].MonthlyCost
	})

	return namespaceCosts, nil
}

func (c *Client) findUnderutilizedResources() ([]UnderutilizedResource, error) {
	utilizations, err := c.GetResourceUtilization()
	if err != nil {
		return nil, err
	}

	var underutilized []UnderutilizedResource
	for _, util := range utilizations {
		if util.CPUUtilization < 20 || util.MemUtilization < 20 {
			pod, err := c.Clientset.CoreV1().Pods(util.Namespace).Get(c.Context, util.Name, metav1.GetOptions{})
			if err != nil {
				continue
			}

			cpuReq, memReq := getPodResourceRequests(pod)

			cpuWaste := int64(float64(cpuReq) * (100 - util.CPUUtilization) / 100)
			memWaste := int64(float64(memReq) * (100 - util.MemUtilization) / 100)

			estimatedSavings := c.estimateResourceSavings(cpuWaste, memWaste)

			recommendation := c.generateRightsizingRecommendation(util.CPUUtilization, util.MemUtilization)

			underutilized = append(underutilized, UnderutilizedResource{
				Type:             "Pod",
				Name:             util.Name,
				Namespace:        util.Namespace,
				CPUWaste:         formatCPU(cpuWaste),
				MemoryWaste:      formatBytes(memWaste),
				EstimatedSavings: estimatedSavings,
				Recommendation:   recommendation,
			})
		}
	}

	sort.Slice(underutilized, func(i, j int) bool {
		return underutilized[i].EstimatedSavings > underutilized[j].EstimatedSavings
	})

	return underutilized, nil
}

func (c *Client) generateCostOptimizations(nodeCosts []NodeCost, namespaceCosts []NamespaceCost, underutilized []UnderutilizedResource) []CostOptimization {
	var optimizations []CostOptimization

	totalWastedCost := 0.0
	for _, resource := range underutilized {
		totalWastedCost += resource.EstimatedSavings
	}

	if totalWastedCost > 50 {
		optimizations = append(optimizations, CostOptimization{
			Type:             "Resource Rightsizing",
			Description:      fmt.Sprintf("Reduce resource requests for %d underutilized workloads", len(underutilized)),
			PotentialSavings: totalWastedCost,
			Priority:         "High",
			Action:           "Review and adjust CPU/Memory requests for underutilized pods",
		})
	}

	inefficientNodes := 0
	for _, node := range nodeCosts {
		if node.CPUUtilization < 30 && node.MemUtilization < 30 {
			inefficientNodes++
		}
	}

	if inefficientNodes > 0 && len(nodeCosts) > 1 {
		potentialSavings := 0.0
		for _, node := range nodeCosts {
			if node.CPUUtilization < 30 && node.MemUtilization < 30 {
				potentialSavings += node.MonthlyCost * 0.7
			}
		}

		optimizations = append(optimizations, CostOptimization{
			Type:             "Node Consolidation",
			Description:      fmt.Sprintf("Consolidate workloads from %d underutilized nodes", inefficientNodes),
			PotentialSavings: potentialSavings,
			Priority:         "Medium",
			Action:           "Consider using node affinity to consolidate workloads",
		})
	}

	if len(namespaceCosts) > 0 {
		highCostNamespaces := 0
		for _, ns := range namespaceCosts {
			if ns.MonthlyCost > 100 {
				highCostNamespaces++
			}
		}

		if highCostNamespaces > 0 {
			optimizations = append(optimizations, CostOptimization{
				Type:             "Namespace Optimization",
				Description:      fmt.Sprintf("Review resource allocation in %d high-cost namespaces", highCostNamespaces),
				PotentialSavings: 0.0,
				Priority:         "Medium",
				Action:           "Implement resource quotas and limits in expensive namespaces",
			})
		}
	}

	optimizations = append(optimizations, CostOptimization{
		Type:             "Monitoring",
		Description:      "Set up cost monitoring and alerting",
		PotentialSavings: 0.0,
		Priority:         "Low",
		Action:           "Implement resource usage monitoring and cost alerts",
	})

	return optimizations
}

func (c *Client) extractNodeType(node *corev1.Node) string {
	if instanceType, exists := node.Labels["node.kubernetes.io/instance-type"]; exists {
		return instanceType
	}
	if instanceType, exists := node.Labels["beta.kubernetes.io/instance-type"]; exists {
		return instanceType
	}
	return "default"
}

func (c *Client) getNodeCost(nodeType string) float64 {
	if cost, exists := nodeTypeCosts[nodeType]; exists {
		return cost
	}
	return nodeTypeCosts["default"]
}

func (c *Client) calculateEfficiency(cpuUtil, memUtil float64) string {
	avgUtil := (cpuUtil + memUtil) / 2
	if avgUtil > 70 {
		return "Excellent"
	} else if avgUtil > 50 {
		return "Good"
	} else if avgUtil > 30 {
		return "Fair"
	} else {
		return "Poor"
	}
}

func (c *Client) estimateNamespaceCost(cpuRequests, memRequests int64) float64 {
	cpuCostPerCore := 20.0
	memCostPerGB := 5.0

	cpuCores := float64(cpuRequests) / 1000
	memGB := float64(memRequests) / (1024 * 1024 * 1024)

	return (cpuCores * cpuCostPerCore) + (memGB * memCostPerGB)
}

func (c *Client) estimateResourceSavings(cpuWaste, memWaste int64) float64 {
	cpuCostPerCore := 20.0
	memCostPerGB := 5.0

	cpuCores := float64(cpuWaste) / 1000
	memGB := float64(memWaste) / (1024 * 1024 * 1024)

	return (cpuCores * cpuCostPerCore) + (memGB * memCostPerGB)
}

func (c *Client) generateRightsizingRecommendation(cpuUtil, memUtil float64) string {
	if cpuUtil < 10 && memUtil < 10 {
		return "Consider reducing requests by 50-70%"
	} else if cpuUtil < 20 && memUtil < 20 {
		return "Consider reducing requests by 30-50%"
	} else {
		return "Consider reducing requests by 10-30%"
	}
}
