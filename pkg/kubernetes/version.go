package kubernetes

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterInfo struct {
	ServerVersion   string
	GitVersion      string
	Major           string
	Minor           string
	Platform        string
	BuildDate       string
	GoVersion       string
	Compiler        string
	GitCommit       string
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
}

func (c *Client) GetInstalledComponents() ([]ComponentInfo, error) {
	var components []ComponentInfo

	commonComponents := []struct {
		name      string
		namespace string
	}{
		{"metrics-server", "kube-system"},
		{"argocd-server", "argocd"},
		{"argocd-server", "argo-cd"},
		{"kuma-control-plane", "kuma-system"},
		{"istio-proxy", "istio-system"},
		{"traefik", "traefik-system"},
		{"nginx-ingress", "ingress-nginx"},
		{"cert-manager", "cert-manager"},
		{"prometheus", "monitoring"},
		{"grafana", "monitoring"},
		{"jaeger", "jaeger"},
		{"kiali", "istio-system"},
	}

	for _, comp := range commonComponents {
		deployments, err := c.Clientset.AppsV1().Deployments(comp.namespace).List(c.Context, metav1.ListOptions{})
		if err != nil {
			continue
		}

		for _, dep := range deployments.Items {
			if dep.Name == comp.name || containsString(dep.Name, comp.name) {
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
				})
			}
		}
	}

	daemonsets, err := c.Clientset.AppsV1().DaemonSets("").List(c.Context, metav1.ListOptions{})
	if err == nil {
		for _, ds := range daemonsets.Items {
			for _, comp := range commonComponents {
				if ds.Name == comp.name || containsString(ds.Name, comp.name) {
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
					})
				}
			}
		}
	}

	return components, nil
}