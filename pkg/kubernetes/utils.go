package kubernetes

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
)

func containsString(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func getMainContainerImage(deployment *appsv1.Deployment) string {
	if len(deployment.Spec.Template.Spec.Containers) > 0 {
		return deployment.Spec.Template.Spec.Containers[0].Image
	}
	return ""
}

func getMainContainerImageDS(daemonset *appsv1.DaemonSet) string {
	if len(daemonset.Spec.Template.Spec.Containers) > 0 {
		return daemonset.Spec.Template.Spec.Containers[0].Image
	}
	return ""
}

func extractVersionFromImage(image string) string {
	parts := strings.Split(image, ":")
	if len(parts) > 1 {
		version := parts[len(parts)-1]
		if version != "latest" {
			return version
		}
	}
	return "latest"
}

func isInterestingComponent(name string, commonComponents []string) bool {
	name = strings.ToLower(name)
	for _, comp := range commonComponents {
		if strings.Contains(name, strings.ToLower(comp)) {
			return true
		}
	}
	return false
}

func removeDuplicateComponents(components []ComponentInfo) []ComponentInfo {
	seen := make(map[string]ComponentInfo)

	// First pass - collect all components, prioritizing Helm sources
	for _, comp := range components {
		key := comp.Namespace + "/" + comp.Name
		existing, exists := seen[key]

		if !exists {
			seen[key] = comp
		} else {
			// Prefer Helm source over others
			if comp.Source == "Helm" && existing.Source != "Helm" {
				seen[key] = comp
			}
		}
	}

	// Convert back to slice
	result := make([]ComponentInfo, 0, len(seen))
	for _, comp := range seen {
		result = append(result, comp)
	}

	return result
}
