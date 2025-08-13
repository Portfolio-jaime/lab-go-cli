package kubernetes

import (
	"fmt"
	"sort"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type WorkloadAnalysis struct {
	DeploymentAnalysis []DeploymentHealth
	StatefulSetAnalysis []StatefulSetHealth
	DaemonSetAnalysis   []DaemonSetHealth
	PodAnalysis         []PodHealth
	WorkloadSummary     WorkloadSummary
}

type DeploymentHealth struct {
	Name                string
	Namespace           string
	Replicas            int32
	ReadyReplicas       int32
	AvailableReplicas   int32
	UnavailableReplicas int32
	Status              string
	Age                 string
	RestartRate         float64
	ResourceEfficiency  string
	HealthScore         int
	Issues              []string
	Recommendations     []string
}

type StatefulSetHealth struct {
	Name            string
	Namespace       string
	Replicas        int32
	ReadyReplicas   int32
	CurrentReplicas int32
	Status          string
	Age             string
	HealthScore     int
	Issues          []string
	Recommendations []string
}

type DaemonSetHealth struct {
	Name                 string
	Namespace            string
	DesiredNumberScheduled int32
	CurrentNumberScheduled int32
	NumberReady          int32
	NumberUnavailable    int32
	Status               string
	Age                  string
	HealthScore          int
	Issues               []string
	Recommendations      []string
}

type PodHealth struct {
	Name            string
	Namespace       string
	Status          string
	RestartCount    int32
	Age             string
	Node            string
	CPUUsage        string
	MemoryUsage     string
	HealthScore     int
	Issues          []string
	LastRestartTime time.Time
}

type WorkloadSummary struct {
	TotalDeployments    int
	HealthyDeployments  int
	TotalStatefulSets   int
	HealthyStatefulSets int
	TotalDaemonSets     int
	HealthyDaemonSets   int
	TotalPods           int
	HealthyPods         int
	CriticalIssues      int
	OverallHealthScore  int
}

func (c *Client) GetWorkloadAnalysis(namespace string) (*WorkloadAnalysis, error) {
	analysis := &WorkloadAnalysis{}

	deployments, err := c.analyzeDeployments(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze deployments: %w", err)
	}
	analysis.DeploymentAnalysis = deployments

	statefulSets, err := c.analyzeStatefulSets(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze statefulsets: %w", err)
	}
	analysis.StatefulSetAnalysis = statefulSets

	daemonSets, err := c.analyzeDaemonSets(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze daemonsets: %w", err)
	}
	analysis.DaemonSetAnalysis = daemonSets

	pods, err := c.analyzePods(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze pods: %w", err)
	}
	analysis.PodAnalysis = pods

	analysis.WorkloadSummary = c.calculateWorkloadSummary(deployments, statefulSets, daemonSets, pods)

	return analysis, nil
}

func (c *Client) analyzeDeployments(namespace string) ([]DeploymentHealth, error) {
	deployments, err := c.Clientset.AppsV1().Deployments(namespace).List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var analysis []DeploymentHealth
	for _, deploy := range deployments.Items {
		health := c.analyzeDeploymentHealth(&deploy)
		analysis = append(analysis, health)
	}

	sort.Slice(analysis, func(i, j int) bool {
		return analysis[i].HealthScore < analysis[j].HealthScore
	})

	return analysis, nil
}

func (c *Client) analyzeStatefulSets(namespace string) ([]StatefulSetHealth, error) {
	statefulSets, err := c.Clientset.AppsV1().StatefulSets(namespace).List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var analysis []StatefulSetHealth
	for _, ss := range statefulSets.Items {
		health := c.analyzeStatefulSetHealth(&ss)
		analysis = append(analysis, health)
	}

	sort.Slice(analysis, func(i, j int) bool {
		return analysis[i].HealthScore < analysis[j].HealthScore
	})

	return analysis, nil
}

func (c *Client) analyzeDaemonSets(namespace string) ([]DaemonSetHealth, error) {
	daemonSets, err := c.Clientset.AppsV1().DaemonSets(namespace).List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var analysis []DaemonSetHealth
	for _, ds := range daemonSets.Items {
		health := c.analyzeDaemonSetHealth(&ds)
		analysis = append(analysis, health)
	}

	sort.Slice(analysis, func(i, j int) bool {
		return analysis[i].HealthScore < analysis[j].HealthScore
	})

	return analysis, nil
}

func (c *Client) analyzePods(namespace string) ([]PodHealth, error) {
	pods, err := c.Clientset.CoreV1().Pods(namespace).List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var analysis []PodHealth
	for _, pod := range pods.Items {
		if c.shouldSkipPod(&pod) {
			continue
		}
		health := c.analyzePodHealth(&pod)
		analysis = append(analysis, health)
	}

	sort.Slice(analysis, func(i, j int) bool {
		return analysis[i].HealthScore < analysis[j].HealthScore
	})

	return analysis, nil
}

func (c *Client) analyzeDeploymentHealth(deploy *appsv1.Deployment) DeploymentHealth {
	health := DeploymentHealth{
		Name:                deploy.Name,
		Namespace:           deploy.Namespace,
		Replicas:            *deploy.Spec.Replicas,
		ReadyReplicas:       deploy.Status.ReadyReplicas,
		AvailableReplicas:   deploy.Status.AvailableReplicas,
		UnavailableReplicas: deploy.Status.UnavailableReplicas,
		Age:                 time.Since(deploy.CreationTimestamp.Time).Truncate(time.Second).String(),
		Issues:              []string{},
		Recommendations:     []string{},
	}

	score := 100

	if health.Replicas != health.ReadyReplicas {
		health.Issues = append(health.Issues, fmt.Sprintf("Not all replicas ready (%d/%d)", health.ReadyReplicas, health.Replicas))
		score -= 30
	}

	if health.UnavailableReplicas > 0 {
		health.Issues = append(health.Issues, fmt.Sprintf("%d replicas unavailable", health.UnavailableReplicas))
		score -= 20
	}

	if health.Replicas == 1 {
		health.Issues = append(health.Issues, "Single replica - no high availability")
		health.Recommendations = append(health.Recommendations, "Consider increasing replicas for HA")
		score -= 10
	}

	if deploy.Spec.Template.Spec.Containers[0].Resources.Requests == nil {
		health.Issues = append(health.Issues, "No resource requests defined")
		health.Recommendations = append(health.Recommendations, "Define CPU and memory requests")
		score -= 15
	}

	if deploy.Spec.Template.Spec.Containers[0].Resources.Limits == nil {
		health.Issues = append(health.Issues, "No resource limits defined")
		health.Recommendations = append(health.Recommendations, "Define CPU and memory limits")
		score -= 10
	}

	if deploy.Spec.Template.Spec.Containers[0].LivenessProbe == nil {
		health.Issues = append(health.Issues, "No liveness probe configured")
		health.Recommendations = append(health.Recommendations, "Add liveness probe for better health monitoring")
		score -= 10
	}

	if deploy.Spec.Template.Spec.Containers[0].ReadinessProbe == nil {
		health.Issues = append(health.Issues, "No readiness probe configured")
		health.Recommendations = append(health.Recommendations, "Add readiness probe for better traffic management")
		score -= 10
	}

	if score < 0 {
		score = 0
	}
	health.HealthScore = score

	if score >= 80 {
		health.Status = "Healthy"
	} else if score >= 60 {
		health.Status = "Warning"
	} else {
		health.Status = "Critical"
	}

	return health
}

func (c *Client) analyzeStatefulSetHealth(ss *appsv1.StatefulSet) StatefulSetHealth {
	health := StatefulSetHealth{
		Name:            ss.Name,
		Namespace:       ss.Namespace,
		Replicas:        *ss.Spec.Replicas,
		ReadyReplicas:   ss.Status.ReadyReplicas,
		CurrentReplicas: ss.Status.CurrentReplicas,
		Age:             time.Since(ss.CreationTimestamp.Time).Truncate(time.Second).String(),
		Issues:          []string{},
		Recommendations: []string{},
	}

	score := 100

	if health.Replicas != health.ReadyReplicas {
		health.Issues = append(health.Issues, fmt.Sprintf("Not all replicas ready (%d/%d)", health.ReadyReplicas, health.Replicas))
		score -= 30
	}

	if health.CurrentReplicas != health.Replicas {
		health.Issues = append(health.Issues, fmt.Sprintf("Scaling in progress (%d/%d)", health.CurrentReplicas, health.Replicas))
		score -= 20
	}

	if len(ss.Spec.VolumeClaimTemplates) == 0 {
		health.Issues = append(health.Issues, "No persistent storage configured")
		health.Recommendations = append(health.Recommendations, "Consider adding persistent volume claims")
		score -= 15
	}

	if score < 0 {
		score = 0
	}
	health.HealthScore = score

	if score >= 80 {
		health.Status = "Healthy"
	} else if score >= 60 {
		health.Status = "Warning"
	} else {
		health.Status = "Critical"
	}

	return health
}

func (c *Client) analyzeDaemonSetHealth(ds *appsv1.DaemonSet) DaemonSetHealth {
	health := DaemonSetHealth{
		Name:                   ds.Name,
		Namespace:              ds.Namespace,
		DesiredNumberScheduled: ds.Status.DesiredNumberScheduled,
		CurrentNumberScheduled: ds.Status.CurrentNumberScheduled,
		NumberReady:            ds.Status.NumberReady,
		NumberUnavailable:      ds.Status.NumberUnavailable,
		Age:                    time.Since(ds.CreationTimestamp.Time).Truncate(time.Second).String(),
		Issues:                 []string{},
		Recommendations:        []string{},
	}

	score := 100

	if health.NumberReady != health.DesiredNumberScheduled {
		health.Issues = append(health.Issues, fmt.Sprintf("Not all instances ready (%d/%d)", health.NumberReady, health.DesiredNumberScheduled))
		score -= 30
	}

	if health.NumberUnavailable > 0 {
		health.Issues = append(health.Issues, fmt.Sprintf("%d instances unavailable", health.NumberUnavailable))
		score -= 25
	}

	if health.CurrentNumberScheduled != health.DesiredNumberScheduled {
		health.Issues = append(health.Issues, "Scheduling issues detected")
		score -= 20
	}

	if score < 0 {
		score = 0
	}
	health.HealthScore = score

	if score >= 80 {
		health.Status = "Healthy"
	} else if score >= 60 {
		health.Status = "Warning"
	} else {
		health.Status = "Critical"
	}

	return health
}

func (c *Client) analyzePodHealth(pod *corev1.Pod) PodHealth {
	health := PodHealth{
		Name:         pod.Name,
		Namespace:    pod.Namespace,
		Status:       string(pod.Status.Phase),
		RestartCount: c.getTotalPodRestarts(pod),
		Age:          time.Since(pod.CreationTimestamp.Time).Truncate(time.Second).String(),
		Node:         pod.Spec.NodeName,
		Issues:       []string{},
	}

	score := 100

	if pod.Status.Phase != corev1.PodRunning {
		health.Issues = append(health.Issues, fmt.Sprintf("Pod not running (status: %s)", pod.Status.Phase))
		score -= 40
	}

	if health.RestartCount > 5 {
		health.Issues = append(health.Issues, fmt.Sprintf("High restart count (%d)", health.RestartCount))
		score -= 20
	} else if health.RestartCount > 0 {
		health.Issues = append(health.Issues, fmt.Sprintf("Has restarted %d times", health.RestartCount))
		score -= 10
	}

	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerStatus.Ready {
			health.Issues = append(health.Issues, fmt.Sprintf("Container %s not ready", containerStatus.Name))
			score -= 15
		}
		
		if containerStatus.LastTerminationState.Terminated != nil {
			health.LastRestartTime = containerStatus.LastTerminationState.Terminated.FinishedAt.Time
		}
	}

	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady && condition.Status != corev1.ConditionTrue {
			health.Issues = append(health.Issues, "Pod not ready")
			score -= 25
		}
	}

	if score < 0 {
		score = 0
	}
	health.HealthScore = score

	return health
}

func (c *Client) calculateWorkloadSummary(deployments []DeploymentHealth, statefulSets []StatefulSetHealth, daemonSets []DaemonSetHealth, pods []PodHealth) WorkloadSummary {
	summary := WorkloadSummary{
		TotalDeployments:  len(deployments),
		TotalStatefulSets: len(statefulSets),
		TotalDaemonSets:   len(daemonSets),
		TotalPods:         len(pods),
	}

	totalScore := 0
	totalWorkloads := 0

	for _, deploy := range deployments {
		if deploy.HealthScore >= 80 {
			summary.HealthyDeployments++
		}
		if deploy.Status == "Critical" {
			summary.CriticalIssues++
		}
		totalScore += deploy.HealthScore
		totalWorkloads++
	}

	for _, ss := range statefulSets {
		if ss.HealthScore >= 80 {
			summary.HealthyStatefulSets++
		}
		if ss.Status == "Critical" {
			summary.CriticalIssues++
		}
		totalScore += ss.HealthScore
		totalWorkloads++
	}

	for _, ds := range daemonSets {
		if ds.HealthScore >= 80 {
			summary.HealthyDaemonSets++
		}
		if ds.Status == "Critical" {
			summary.CriticalIssues++
		}
		totalScore += ds.HealthScore
		totalWorkloads++
	}

	for _, pod := range pods {
		if pod.HealthScore >= 80 {
			summary.HealthyPods++
		}
	}

	if totalWorkloads > 0 {
		summary.OverallHealthScore = totalScore / totalWorkloads
	} else {
		summary.OverallHealthScore = 100
	}

	return summary
}

func (c *Client) getTotalPodRestarts(pod *corev1.Pod) int32 {
	var totalRestarts int32
	for _, containerStatus := range pod.Status.ContainerStatuses {
		totalRestarts += containerStatus.RestartCount
	}
	return totalRestarts
}

func (c *Client) shouldSkipPod(pod *corev1.Pod) bool {
	if pod.Status.Phase == corev1.PodSucceeded {
		return true
	}
	
	for _, ownerRef := range pod.OwnerReferences {
		if ownerRef.Kind == "Job" && pod.Status.Phase == corev1.PodSucceeded {
			return true
		}
	}

	return false
}