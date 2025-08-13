package kubernetes

import (
	"fmt"
	"sort"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterEvent struct {
	Type      string
	Reason    string
	Message   string
	Object    string
	Namespace string
	FirstTime time.Time
	LastTime  time.Time
	Count     int32
	Severity  string
	Component string
}

type LogAnalysis struct {
	CriticalEvents []ClusterEvent
	WarningEvents  []ClusterEvent
	ErrorPatterns  []ErrorPattern
	ResourceEvents []ResourceEvent
	SecurityEvents []SecurityEvent
}

type ErrorPattern struct {
	Pattern        string
	Count          int
	LastSeen       time.Time
	Severity       string
	Description    string
	Recommendation string
}

type ResourceEvent struct {
	Type         string
	ResourceName string
	Namespace    string
	Event        string
	Timestamp    time.Time
	Impact       string
}

type SecurityEvent struct {
	Type        string
	Description string
	Object      string
	Namespace   string
	Timestamp   time.Time
	RiskLevel   string
	Action      string
}

type PodLogSummary struct {
	PodName        string
	Namespace      string
	ErrorCount     int
	WarningCount   int
	CriticalIssues []string
	Status         string
	LastRestart    time.Time
}

func (c *Client) GetClusterEvents(namespace string, hours int) ([]ClusterEvent, error) {
	timeWindow := time.Now().Add(-time.Duration(hours) * time.Hour)

	listOptions := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("firstTimestamp>%s", timeWindow.Format(time.RFC3339)),
	}

	events, err := c.Clientset.CoreV1().Events(namespace).List(c.Context, listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	var clusterEvents []ClusterEvent
	for _, event := range events.Items {
		severity := categorizeSeverity(&event)
		component := extractComponent(&event)

		clusterEvent := ClusterEvent{
			Type:      event.Type,
			Reason:    event.Reason,
			Message:   event.Message,
			Object:    fmt.Sprintf("%s/%s", event.InvolvedObject.Kind, event.InvolvedObject.Name),
			Namespace: event.Namespace,
			FirstTime: event.FirstTimestamp.Time,
			LastTime:  event.LastTimestamp.Time,
			Count:     event.Count,
			Severity:  severity,
			Component: component,
		}
		clusterEvents = append(clusterEvents, clusterEvent)
	}

	sort.Slice(clusterEvents, func(i, j int) bool {
		return clusterEvents[i].LastTime.After(clusterEvents[j].LastTime)
	})

	return clusterEvents, nil
}

func (c *Client) GetLogAnalysis(namespace string, hours int) (*LogAnalysis, error) {
	events, err := c.GetClusterEvents(namespace, hours)
	if err != nil {
		return nil, err
	}

	analysis := &LogAnalysis{
		CriticalEvents: []ClusterEvent{},
		WarningEvents:  []ClusterEvent{},
		ErrorPatterns:  []ErrorPattern{},
		ResourceEvents: []ResourceEvent{},
		SecurityEvents: []SecurityEvent{},
	}

	for _, event := range events {
		switch event.Severity {
		case "Critical":
			analysis.CriticalEvents = append(analysis.CriticalEvents, event)
		case "Warning":
			analysis.WarningEvents = append(analysis.WarningEvents, event)
		}

		if resourceEvent := analyzeResourceEvent(&event); resourceEvent != nil {
			analysis.ResourceEvents = append(analysis.ResourceEvents, *resourceEvent)
		}

		if securityEvent := analyzeSecurityEvent(&event); securityEvent != nil {
			analysis.SecurityEvents = append(analysis.SecurityEvents, *securityEvent)
		}
	}

	analysis.ErrorPatterns = findErrorPatterns(events)

	return analysis, nil
}

func (c *Client) GetPodLogsAnalysis(namespace string) ([]PodLogSummary, error) {
	pods, err := c.Clientset.CoreV1().Pods(namespace).List(c.Context, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	var summaries []PodLogSummary
	for _, pod := range pods.Items {
		summary := PodLogSummary{
			PodName:        pod.Name,
			Namespace:      pod.Namespace,
			Status:         string(pod.Status.Phase),
			CriticalIssues: []string{},
		}

		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.RestartCount > 0 && containerStatus.LastTerminationState.Terminated != nil {
				summary.LastRestart = containerStatus.LastTerminationState.Terminated.FinishedAt.Time
			}
		}

		events, err := c.getEventsForPod(&pod)
		if err == nil {
			summary.ErrorCount, summary.WarningCount, summary.CriticalIssues = analyzePodEvents(events)
		}

		if summary.ErrorCount > 0 || summary.WarningCount > 0 || len(summary.CriticalIssues) > 0 {
			summaries = append(summaries, summary)
		}
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].ErrorCount+summaries[i].WarningCount > summaries[j].ErrorCount+summaries[j].WarningCount
	})

	return summaries, nil
}

func (c *Client) getEventsForPod(pod *corev1.Pod) ([]corev1.Event, error) {
	fieldSelector := fmt.Sprintf("involvedObject.name=%s,involvedObject.uid=%s", pod.Name, pod.UID)
	listOptions := metav1.ListOptions{
		FieldSelector: fieldSelector,
	}

	events, err := c.Clientset.CoreV1().Events(pod.Namespace).List(c.Context, listOptions)
	if err != nil {
		return nil, err
	}

	return events.Items, nil
}

func categorizeSeverity(event *corev1.Event) string {
	criticalReasons := []string{
		"Failed", "FailedScheduling", "FailedMount", "FailedAttachVolume",
		"FailedCreatePodSandBox", "FailedPodSandBoxStatus", "NetworkNotReady",
		"FailedKillPod", "FailedCreatePodContainer", "InspectFailed",
	}

	warningReasons := []string{
		"Unhealthy", "ProbeWarning", "BackOff", "ImagePullBackOff",
		"ErrImagePull", "NodeNotReady", "SystemOOM", "FreeDiskSpaceFailed",
		"DeadlineExceeded", "EvictionThresholdMet",
	}

	reason := event.Reason

	for _, criticalReason := range criticalReasons {
		if strings.Contains(reason, criticalReason) {
			return "Critical"
		}
	}

	for _, warningReason := range warningReasons {
		if strings.Contains(reason, warningReason) {
			return "Warning"
		}
	}

	if event.Type == "Warning" {
		return "Warning"
	}

	return "Info"
}

func extractComponent(event *corev1.Event) string {
	if event.Source.Component != "" {
		return event.Source.Component
	}
	if event.Source.Host != "" {
		return event.Source.Host
	}
	if event.ReportingController != "" {
		return event.ReportingController
	}
	return "Unknown"
}

func findErrorPatterns(events []ClusterEvent) []ErrorPattern {
	patterns := make(map[string]*ErrorPattern)

	for _, event := range events {
		if event.Severity == "Critical" || event.Severity == "Warning" {
			key := fmt.Sprintf("%s:%s", event.Reason, event.Type)

			if pattern, exists := patterns[key]; exists {
				pattern.Count += int(event.Count)
				if event.LastTime.After(pattern.LastSeen) {
					pattern.LastSeen = event.LastTime
				}
			} else {
				patterns[key] = &ErrorPattern{
					Pattern:        event.Reason,
					Count:          int(event.Count),
					LastSeen:       event.LastTime,
					Severity:       event.Severity,
					Description:    generatePatternDescription(event.Reason),
					Recommendation: generatePatternRecommendation(event.Reason),
				}
			}
		}
	}

	var result []ErrorPattern
	for _, pattern := range patterns {
		result = append(result, *pattern)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	return result
}

func analyzeResourceEvent(event *ClusterEvent) *ResourceEvent {
	resourceReasons := map[string]string{
		"FailedScheduling":     "Scheduling Issues",
		"FailedMount":          "Volume Issues",
		"SystemOOM":            "Memory Issues",
		"FreeDiskSpaceFailed":  "Disk Issues",
		"NodeNotReady":         "Node Issues",
		"EvictionThresholdMet": "Resource Pressure",
	}

	for reason, eventType := range resourceReasons {
		if strings.Contains(event.Reason, reason) {
			impact := determineImpact(reason)
			return &ResourceEvent{
				Type:         eventType,
				ResourceName: event.Object,
				Namespace:    event.Namespace,
				Event:        event.Reason,
				Timestamp:    event.LastTime,
				Impact:       impact,
			}
		}
	}

	return nil
}

func analyzeSecurityEvent(event *ClusterEvent) *SecurityEvent {
	securityPatterns := map[string]struct {
		description string
		riskLevel   string
		action      string
	}{
		"FailedMount": {
			"Volume mount failed - potential security misconfiguration",
			"Medium",
			"Check volume permissions and security contexts",
		},
		"ImagePullBackOff": {
			"Image pull failed - potential registry access issue",
			"Low",
			"Verify image registry credentials and policies",
		},
		"Forbidden": {
			"Access denied - RBAC or security policy violation",
			"High",
			"Review RBAC permissions and security policies",
		},
	}

	for pattern, info := range securityPatterns {
		if strings.Contains(event.Reason, pattern) || strings.Contains(event.Message, pattern) {
			return &SecurityEvent{
				Type:        "Security",
				Description: info.description,
				Object:      event.Object,
				Namespace:   event.Namespace,
				Timestamp:   event.LastTime,
				RiskLevel:   info.riskLevel,
				Action:      info.action,
			}
		}
	}

	return nil
}

func analyzePodEvents(events []corev1.Event) (int, int, []string) {
	errorCount := 0
	warningCount := 0
	var criticalIssues []string

	for _, event := range events {
		severity := categorizeSeverity(&event)

		switch severity {
		case "Critical":
			errorCount++
			criticalIssues = append(criticalIssues, fmt.Sprintf("%s: %s", event.Reason, event.Message))
		case "Warning":
			warningCount++
		}
	}

	return errorCount, warningCount, criticalIssues
}

func generatePatternDescription(reason string) string {
	descriptions := map[string]string{
		"FailedScheduling":       "Pod cannot be scheduled to any node",
		"FailedMount":            "Volume mounting failed",
		"ImagePullBackOff":       "Cannot pull container image",
		"SystemOOM":              "Out of memory condition",
		"FreeDiskSpaceFailed":    "Insufficient disk space",
		"NodeNotReady":           "Node is not in ready state",
		"EvictionThresholdMet":   "Node resource pressure detected",
		"FailedCreatePodSandBox": "Pod sandbox creation failed",
		"NetworkNotReady":        "Network not ready for pod",
		"DeadlineExceeded":       "Operation timed out",
	}

	if desc, exists := descriptions[reason]; exists {
		return desc
	}
	return "Unknown error pattern"
}

func generatePatternRecommendation(reason string) string {
	recommendations := map[string]string{
		"FailedScheduling":       "Check node capacity, taints, and pod resource requests",
		"FailedMount":            "Verify volume availability and mount permissions",
		"ImagePullBackOff":       "Check image name, registry credentials, and network connectivity",
		"SystemOOM":              "Increase memory limits or optimize application memory usage",
		"FreeDiskSpaceFailed":    "Clean up disk space or add more storage capacity",
		"NodeNotReady":           "Check node health and kubelet status",
		"EvictionThresholdMet":   "Add more resources or reduce workload density",
		"FailedCreatePodSandBox": "Check container runtime and network configuration",
		"NetworkNotReady":        "Verify CNI plugin and network policies",
		"DeadlineExceeded":       "Increase timeout values or optimize operation performance",
	}

	if rec, exists := recommendations[reason]; exists {
		return rec
	}
	return "Review logs and cluster configuration"
}

func determineImpact(reason string) string {
	highImpact := []string{"FailedScheduling", "SystemOOM", "NodeNotReady"}
	mediumImpact := []string{"FailedMount", "FreeDiskSpaceFailed", "EvictionThresholdMet"}

	for _, high := range highImpact {
		if strings.Contains(reason, high) {
			return "High"
		}
	}

	for _, medium := range mediumImpact {
		if strings.Contains(reason, medium) {
			return "Medium"
		}
	}

	return "Low"
}
