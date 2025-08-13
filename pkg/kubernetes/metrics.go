package kubernetes

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeMetrics struct {
	Name            string
	CPUUsage        string
	CPUUsagePercent float64
	MemoryUsage     string
	MemoryUsagePercent float64
	CPUCapacity     string
	MemoryCapacity  string
	Status          string
}

type PodMetrics struct {
	Name            string
	Namespace       string
	CPUUsage        string
	MemoryUsage     string
	CPURequests     string
	MemoryRequests  string
	CPULimits       string
	MemoryLimits    string
	Node            string
	RestartCount    int32
}

type ClusterMetrics struct {
	TotalCPUUsage        string
	TotalMemoryUsage     string
	TotalCPUCapacity     string
	TotalMemoryCapacity  string
	CPUUsagePercent      float64
	MemoryUsagePercent   float64
	NodesCount           int
	PodsCount            int
	NamespacesCount      int
}

type ResourceUtilization struct {
	Type            string
	Name            string
	Namespace       string
	CPUUtilization  float64
	MemUtilization  float64
	Recommendation  string
}

func (c *Client) GetRealTimeNodeMetrics() ([]NodeMetrics, error) {
	nodeMetrics, err := c.MetricsClient.MetricsV1beta1().NodeMetricses().List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get node metrics: %w", err)
	}

	nodes, err := c.Clientset.CoreV1().Nodes().List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	nodeCapacity := make(map[string]corev1.Node)
	for _, node := range nodes.Items {
		nodeCapacity[node.Name] = node
	}

	var metrics []NodeMetrics
	for _, metric := range nodeMetrics.Items {
		node, exists := nodeCapacity[metric.Name]
		if !exists {
			continue
		}

		cpuUsage := metric.Usage[corev1.ResourceCPU]
		memUsage := metric.Usage[corev1.ResourceMemory]
		
		cpuCapacity := node.Status.Capacity[corev1.ResourceCPU]
		memCapacity := node.Status.Capacity[corev1.ResourceMemory]

		cpuUsagePercent := float64(cpuUsage.MilliValue()) / float64(cpuCapacity.MilliValue()) * 100
		memUsagePercent := float64(memUsage.Value()) / float64(memCapacity.Value()) * 100

		status := "Ready"
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status != corev1.ConditionTrue {
				status = "NotReady"
				break
			}
		}

		metrics = append(metrics, NodeMetrics{
			Name:               metric.Name,
			CPUUsage:          formatCPU(cpuUsage.MilliValue()),
			CPUUsagePercent:   cpuUsagePercent,
			MemoryUsage:       formatBytes(memUsage.Value()),
			MemoryUsagePercent: memUsagePercent,
			CPUCapacity:       formatCPU(cpuCapacity.MilliValue()),
			MemoryCapacity:    formatBytes(memCapacity.Value()),
			Status:            status,
		})
	}

	return metrics, nil
}

func (c *Client) GetRealTimePodMetrics(namespace string) ([]PodMetrics, error) {
	podMetrics, err := c.MetricsClient.MetricsV1beta1().PodMetricses(namespace).List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod metrics: %w", err)
	}

	pods, err := c.Clientset.CoreV1().Pods(namespace).List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	podInfo := make(map[string]corev1.Pod)
	for _, pod := range pods.Items {
		podInfo[pod.Name] = pod
	}

	var metrics []PodMetrics
	for _, metric := range podMetrics.Items {
		pod, exists := podInfo[metric.Name]
		if !exists {
			continue
		}

		var totalCPUUsage, totalMemUsage int64
		for _, container := range metric.Containers {
			cpuQuantity := container.Usage[corev1.ResourceCPU]
			memQuantity := container.Usage[corev1.ResourceMemory]
			totalCPUUsage += cpuQuantity.MilliValue()
			totalMemUsage += memQuantity.Value()
		}

		cpuRequests, memRequests := getPodResourceRequests(&pod)
		cpuLimits, memLimits := getPodResourceLimits(&pod)
		restartCount := getTotalRestarts(&pod)

		metrics = append(metrics, PodMetrics{
			Name:           metric.Name,
			Namespace:      metric.Namespace,
			CPUUsage:       formatCPU(totalCPUUsage),
			MemoryUsage:    formatBytes(totalMemUsage),
			CPURequests:    formatCPU(cpuRequests),
			MemoryRequests: formatBytes(memRequests),
			CPULimits:      formatCPU(cpuLimits),
			MemoryLimits:   formatBytes(memLimits),
			Node:           pod.Spec.NodeName,
			RestartCount:   restartCount,
		})
	}

	return metrics, nil
}

func (c *Client) GetClusterMetrics() (*ClusterMetrics, error) {
	nodeMetrics, err := c.MetricsClient.MetricsV1beta1().NodeMetricses().List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get node metrics: %w", err)
	}

	nodes, err := c.Clientset.CoreV1().Nodes().List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	pods, err := c.Clientset.CoreV1().Pods("").List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	namespaces, err := c.Clientset.CoreV1().Namespaces().List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get namespaces: %w", err)
	}

	var totalCPUUsage, totalMemUsage int64
	var totalCPUCapacity, totalMemCapacity int64

	for _, metric := range nodeMetrics.Items {
		cpuQuantity := metric.Usage[corev1.ResourceCPU]
		memQuantity := metric.Usage[corev1.ResourceMemory]
		totalCPUUsage += cpuQuantity.MilliValue()
		totalMemUsage += memQuantity.Value()
	}

	for _, node := range nodes.Items {
		cpuQuantity := node.Status.Capacity[corev1.ResourceCPU]
		memQuantity := node.Status.Capacity[corev1.ResourceMemory]
		totalCPUCapacity += cpuQuantity.MilliValue()
		totalMemCapacity += memQuantity.Value()
	}

	cpuUsagePercent := float64(totalCPUUsage) / float64(totalCPUCapacity) * 100
	memUsagePercent := float64(totalMemUsage) / float64(totalMemCapacity) * 100

	return &ClusterMetrics{
		TotalCPUUsage:       formatCPU(totalCPUUsage),
		TotalMemoryUsage:    formatBytes(totalMemUsage),
		TotalCPUCapacity:    formatCPU(totalCPUCapacity),
		TotalMemoryCapacity: formatBytes(totalMemCapacity),
		CPUUsagePercent:     cpuUsagePercent,
		MemoryUsagePercent:  memUsagePercent,
		NodesCount:          len(nodes.Items),
		PodsCount:           len(pods.Items),
		NamespacesCount:     len(namespaces.Items),
	}, nil
}

func (c *Client) GetResourceUtilization() ([]ResourceUtilization, error) {
	podMetrics, err := c.MetricsClient.MetricsV1beta1().PodMetricses("").List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod metrics: %w", err)
	}

	pods, err := c.Clientset.CoreV1().Pods("").List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	podInfo := make(map[string]corev1.Pod)
	for _, pod := range pods.Items {
		key := fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)
		podInfo[key] = pod
	}

	var utilizations []ResourceUtilization
	for _, metric := range podMetrics.Items {
		key := fmt.Sprintf("%s/%s", metric.Namespace, metric.Name)
		pod, exists := podInfo[key]
		if !exists {
			continue
		}

		var totalCPUUsage, totalMemUsage int64
		for _, container := range metric.Containers {
			cpuQuantity := container.Usage[corev1.ResourceCPU]
			memQuantity := container.Usage[corev1.ResourceMemory]
			totalCPUUsage += cpuQuantity.MilliValue()
			totalMemUsage += memQuantity.Value()
		}

		cpuRequests, memRequests := getPodResourceRequests(&pod)
		
		var cpuUtilization, memUtilization float64
		var recommendation string

		if cpuRequests > 0 {
			cpuUtilization = float64(totalCPUUsage) / float64(cpuRequests) * 100
		}

		if memRequests > 0 {
			memUtilization = float64(totalMemUsage) / float64(memRequests) * 100
		}

		if cpuUtilization < 20 && memUtilization < 20 {
			recommendation = "Consider reducing resource requests - underutilized"
		} else if cpuUtilization > 90 || memUtilization > 90 {
			recommendation = "Consider increasing resource requests - overutilized"
		} else {
			recommendation = "Resource allocation looks good"
		}

		utilizations = append(utilizations, ResourceUtilization{
			Type:           "Pod",
			Name:           metric.Name,
			Namespace:      metric.Namespace,
			CPUUtilization: cpuUtilization,
			MemUtilization: memUtilization,
			Recommendation: recommendation,
		})
	}

	return utilizations, nil
}

func getPodResourceRequests(pod *corev1.Pod) (int64, int64) {
	var cpuRequests, memRequests int64
	for _, container := range pod.Spec.Containers {
		if cpu, exists := container.Resources.Requests[corev1.ResourceCPU]; exists && !cpu.IsZero() {
			cpuRequests += cpu.MilliValue()
		}
		if mem, exists := container.Resources.Requests[corev1.ResourceMemory]; exists && !mem.IsZero() {
			memRequests += mem.Value()
		}
	}
	return cpuRequests, memRequests
}

func getPodResourceLimits(pod *corev1.Pod) (int64, int64) {
	var cpuLimits, memLimits int64
	for _, container := range pod.Spec.Containers {
		if cpu, exists := container.Resources.Limits[corev1.ResourceCPU]; exists && !cpu.IsZero() {
			cpuLimits += cpu.MilliValue()
		}
		if mem, exists := container.Resources.Limits[corev1.ResourceMemory]; exists && !mem.IsZero() {
			memLimits += mem.Value()
		}
	}
	return cpuLimits, memLimits
}

func getTotalRestarts(pod *corev1.Pod) int32 {
	var totalRestarts int32
	for _, containerStatus := range pod.Status.ContainerStatuses {
		totalRestarts += containerStatus.RestartCount
	}
	return totalRestarts
}

func formatCPU(milliCores int64) string {
	if milliCores == 0 {
		return "0m"
	}
	if milliCores < 1000 {
		return fmt.Sprintf("%dm", milliCores)
	}
	return fmt.Sprintf("%.2f", float64(milliCores)/1000)
}

func formatBytes(bytes int64) string {
	if bytes == 0 {
		return "0B"
	}
	
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	suffixes := []string{"Ki", "Mi", "Gi", "Ti", "Pi", "Ei"}
	return fmt.Sprintf("%.1f %sB", float64(bytes)/float64(div), suffixes[exp])
}