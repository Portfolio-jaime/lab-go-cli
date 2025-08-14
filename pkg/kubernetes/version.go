package kubernetes

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ClusterInfo struct {
	ServerVersion string
	GitVersion    string
	Major         string
	Minor         string
	Platform      string
	BuildDate     string
	GoVersion     string
	Compiler      string
	GitCommit     string
}

func (c *Client) GetClusterVersion() (*ClusterInfo, error) {
	version, err := c.Clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get server version: %w", err)
	}

	return &ClusterInfo{
		ServerVersion: version.String(),
		GitVersion:    version.GitVersion,
		Major:         version.Major,
		Minor:         version.Minor,
		Platform:      version.Platform,
		BuildDate:     version.BuildDate,
		GoVersion:     version.GoVersion,
		Compiler:      version.Compiler,
		GitCommit:     version.GitCommit,
	}, nil
}

type ComponentInfo struct {
	Name      string
	Namespace string
	Status    string
	Version   string
	Ready     string
	Source    string // "Kubernetes", "Helm", "StatefulSet", etc.
}

func (c *Client) GetInstalledComponents() ([]ComponentInfo, error) {
	var components []ComponentInfo

	// Get components from Helm releases
	helmComponents, err := c.getHelmComponents()
	if err == nil {
		components = append(components, helmComponents...)
	}

	// Get components from standard Kubernetes resources
	k8sComponents, err := c.getKubernetesComponents()
	if err == nil {
		components = append(components, k8sComponents...)
	}

	// Remove duplicates - prefer Helm info when available
	components = removeDuplicateComponents(components)

	return components, nil
}

func (c *Client) getHelmComponents() ([]ComponentInfo, error) {
	var components []ComponentInfo

	// Define the Helm secret resource
	gvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "secrets",
	}

	// List all secrets in all namespaces that are Helm releases
	secrets, err := c.DynamicClient.Resource(gvr).List(c.Context, metav1.ListOptions{
		LabelSelector: "owner=helm",
	})
	if err != nil {
		return components, fmt.Errorf("failed to list helm secrets: %w", err)
	}

	for _, secret := range secrets.Items {
		labels := secret.GetLabels()
		if labels == nil {
			continue
		}

		name, hasName := labels["name"]
		status, hasStatus := labels["status"]
		if !hasName {
			continue
		}

		version := "Unknown"
		if appVersion, ok := labels["app.kubernetes.io/version"]; ok {
			version = appVersion
		} else if chartVersion, ok := labels["version"]; ok {
			version = chartVersion
		}

		statusStr := "Unknown"
		if hasStatus {
			// Capitalize first letter (replacement for deprecated strings.Title)
			if len(status) > 0 {
				statusStr = strings.ToUpper(string(status[0])) + strings.ToLower(status[1:])
			} else {
				statusStr = status
			}
		}

		components = append(components, ComponentInfo{
			Name:      name,
			Namespace: secret.GetNamespace(),
			Status:    statusStr,
			Version:   version,
			Ready:     "Helm",
			Source:    "Helm",
		})
	}

	return components, nil
}

func (c *Client) getKubernetesComponents() ([]ComponentInfo, error) {
	var components []ComponentInfo

	// Get all namespaces to search comprehensively
	namespaces, err := c.Clientset.CoreV1().Namespaces().List(c.Context, metav1.ListOptions{})
	if err != nil {
		return components, fmt.Errorf("failed to list namespaces: %w", err)
	}

	// Common components to look for (expanded list)
	commonComponents := []string{
		"metrics-server", "argocd", "argo", "kuma", "istio", "traefik",
		"nginx", "cert-manager", "prometheus", "grafana", "jaeger",
		"kiali", "fluentd", "elasticsearch", "kibana", "vault",
		"consul", "etcd", "redis", "postgres", "mysql", "mongodb",
		"kafka", "zookeeper", "rabbitmq", "jenkins", "sonarqube",
		"nexus", "harbor", "docker-registry", "ingress", "gateway",
	}

	for _, ns := range namespaces.Items {
		nsName := ns.Name

		// Skip certain system namespaces that are unlikely to have interesting components
		if nsName == "kube-node-lease" || nsName == "kube-public" {
			continue
		}

		// Check Deployments
		deployments, err := c.Clientset.AppsV1().Deployments(nsName).List(c.Context, metav1.ListOptions{})
		if err == nil {
			for _, dep := range deployments.Items {
				if isInterestingComponent(dep.Name, commonComponents) {
					status := "Running"
					if dep.Status.ReadyReplicas == 0 {
						status = "Not Ready"
					}

					version := "Unknown"
					if image := getMainContainerImage(&dep); image != "" {
						version = extractVersionFromImage(image)
					}

					components = append(components, ComponentInfo{
						Name:      dep.Name,
						Namespace: dep.Namespace,
						Status:    status,
						Version:   version,
						Ready:     fmt.Sprintf("%d/%d", dep.Status.ReadyReplicas, dep.Status.Replicas),
						Source:    "Deployment",
					})
				}
			}
		}

		// Check StatefulSets
		statefulsets, err := c.Clientset.AppsV1().StatefulSets(nsName).List(c.Context, metav1.ListOptions{})
		if err == nil {
			for _, sts := range statefulsets.Items {
				if isInterestingComponent(sts.Name, commonComponents) {
					status := "Running"
					if sts.Status.ReadyReplicas == 0 {
						status = "Not Ready"
					}

					version := "Unknown"
					if len(sts.Spec.Template.Spec.Containers) > 0 {
						image := sts.Spec.Template.Spec.Containers[0].Image
						version = extractVersionFromImage(image)
					}

					components = append(components, ComponentInfo{
						Name:      sts.Name,
						Namespace: sts.Namespace,
						Status:    status,
						Version:   version,
						Ready:     fmt.Sprintf("%d/%d", sts.Status.ReadyReplicas, sts.Status.Replicas),
						Source:    "StatefulSet",
					})
				}
			}
		}

		// Check DaemonSets
		daemonsets, err := c.Clientset.AppsV1().DaemonSets(nsName).List(c.Context, metav1.ListOptions{})
		if err == nil {
			for _, ds := range daemonsets.Items {
				if isInterestingComponent(ds.Name, commonComponents) {
					status := "Running"
					if ds.Status.NumberReady == 0 {
						status = "Not Ready"
					}

					version := "Unknown"
					if image := getMainContainerImageDS(&ds); image != "" {
						version = extractVersionFromImage(image)
					}

					components = append(components, ComponentInfo{
						Name:      ds.Name,
						Namespace: ds.Namespace,
						Status:    status,
						Version:   version,
						Ready:     fmt.Sprintf("%d/%d", ds.Status.NumberReady, ds.Status.DesiredNumberScheduled),
						Source:    "DaemonSet",
					})
				}
			}
		}
	}

	return components, nil
}
