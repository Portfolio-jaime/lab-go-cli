package cmd

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestVersionFlag(t *testing.T) {
	// Test that --version flag works by running the built binary
	cmd := exec.Command("./bin/k8s-cli", "--version")
	cmd.Dir = "../"

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Expected no error running --version, got %v", err)
	}

	outputStr := string(output)

	expectedStrings := []string{
		"k8s-cli version",
		"Git commit:",
		"Built:",
		"Go version:",
		"OS/Arch:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Expected output to contain '%s', but got:\n%s", expected, outputStr)
		}
	}
}

func TestVersionShortFlag(t *testing.T) {
	// Test that -v flag works
	cmd := exec.Command("./bin/k8s-cli", "-v")
	cmd.Dir = "../"

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Expected no error running -v, got %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "k8s-cli version") {
		t.Errorf("Expected output to contain version info, but got:\n%s", outputStr)
	}
}

func TestVersionVsVersionCommand(t *testing.T) {
	// Skip this test if no kubeconfig is available
	if _, err := os.Stat(os.Getenv("HOME") + "/.kube/config"); os.IsNotExist(err) {
		t.Skip("Skipping test: no kubeconfig found")
	}

	// Test that 'version' command (without --) still shows Kubernetes info
	cmd := exec.Command("./bin/k8s-cli", "version", "--help")
	cmd.Dir = "../"

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Expected no error running version --help, got %v", err)
	}

	outputStr := string(output)
	// Should show help for the version command (Kubernetes cluster version)
	if !strings.Contains(outputStr, "Kubernetes cluster version") {
		t.Errorf("Expected 'version' command to show Kubernetes cluster version help, but got:\n%s", outputStr)
	}
}
