package kubernetes

import (
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SimpleNodeInfo struct {
	Name           string
	Status         string
	Roles          string
	Age            string
	Version        string
	InternalIP     string
	CPUCapacity    string
	MemoryCapacity string
}

type SimplePodInfo struct {
	Name      string
	Namespace string
	Status    string
	Restarts  string
	Age       string
	Node      string
}

type SimpleClusterSummary struct {
	TotalNodes       int
	TotalPods        int
	TotalCPUCapacity string
	TotalMemCapacity string
}

func (c *Client) GetSimpleNodesInfo() ([]SimpleNodeInfo, error) {
	nodes, err := c.Clientset.CoreV1().Nodes().List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	var nodeInfos []SimpleNodeInfo

	for _, node := range nodes.Items {
		status := "Ready"
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady && condition.Status != corev1.ConditionTrue {
				status = "NotReady"
				break
			}
		}

		roles := getSimpleRoles(&node)
		age := getSimpleAge(node.CreationTimestamp.Time)
		internalIP := getSimpleInternalIP(&node)

		cpuCapacity := node.Status.Capacity[corev1.ResourceCPU]
		memoryCapacity := node.Status.Capacity[corev1.ResourceMemory]

		nodeInfos = append(nodeInfos, SimpleNodeInfo{
			Name:           node.Name,
			Status:         status,
			Roles:          roles,
			Age:            age,
			Version:        node.Status.NodeInfo.KubeletVersion,
			InternalIP:     internalIP,
			CPUCapacity:    cpuCapacity.String(),
			MemoryCapacity: formatSimpleBytes(memoryCapacity.Value()),
		})
	}

	return nodeInfos, nil
}

func (c *Client) GetSimplePodsInfo(namespace string) ([]SimplePodInfo, error) {
	listOptions := metav1.ListOptions{}

	pods, err := c.Clientset.CoreV1().Pods(namespace).List(c.Context, listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	var podInfos []SimplePodInfo

	for _, pod := range pods.Items {
		status := string(pod.Status.Phase)
		if pod.DeletionTimestamp != nil {
			status = "Terminating"
		}

		restarts := getSimpleTotalRestarts(&pod)
		age := getSimpleAge(pod.CreationTimestamp.Time)

		podInfos = append(podInfos, SimplePodInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    status,
			Restarts:  restarts,
			Age:       age,
			Node:      pod.Spec.NodeName,
		})
	}

	return podInfos, nil
}

func (c *Client) GetSimpleClusterSummary() (*SimpleClusterSummary, error) {
	nodes, err := c.Clientset.CoreV1().Nodes().List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	pods, err := c.Clientset.CoreV1().Pods("").List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	var totalCPUCapacity, totalMemCapacity int64

	for _, node := range nodes.Items {
		cpuQuantity := node.Status.Capacity[corev1.ResourceCPU]
		memQuantity := node.Status.Capacity[corev1.ResourceMemory]
		totalCPUCapacity += cpuQuantity.MilliValue()
		totalMemCapacity += memQuantity.Value()
	}

	return &SimpleClusterSummary{
		TotalNodes:       len(nodes.Items),
		TotalPods:        len(pods.Items),
		TotalCPUCapacity: fmt.Sprintf("%.1f", float64(totalCPUCapacity)/1000),
		TotalMemCapacity: formatSimpleBytes(totalMemCapacity),
	}, nil
}

func getSimpleRoles(node *corev1.Node) string {
	roles := []string{}

	for label := range node.Labels {
		if strings.Contains(label, "node-role.kubernetes.io/") {
			role := strings.TrimPrefix(label, "node-role.kubernetes.io/")
			if role != "" {
				roles = append(roles, role)
			}
		}
	}

	if len(roles) == 0 {
		roles = append(roles, "worker")
	}

	return strings.Join(roles, ",")
}

func getSimpleInternalIP(node *corev1.Node) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			return addr.Address
		}
	}
	return "N/A"
}

func getSimpleAge(creationTime time.Time) string {
	duration := time.Since(creationTime)
	days := int(duration.Hours() / 24)

	if days > 0 {
		return fmt.Sprintf("%dd", days)
	}

	hours := int(duration.Hours())
	if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}

	minutes := int(duration.Minutes())
	if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	}

	return fmt.Sprintf("%.0fs", duration.Seconds())
}

func getSimpleTotalRestarts(pod *corev1.Pod) string {
	var totalRestarts int32
	for _, containerStatus := range pod.Status.ContainerStatuses {
		totalRestarts += containerStatus.RestartCount
	}
	return fmt.Sprintf("%d", totalRestarts)
}

func formatSimpleBytes(bytes int64) string {
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
