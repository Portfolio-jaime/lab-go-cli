package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// getBinaryPath returns the correct binary path for the current OS
func getBinaryPath() string {
	binaryName := "k8s-cli"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	return filepath.Join("..", "bin", binaryName)
}

func TestVersionFlag(t *testing.T) {
	// Skip integration tests on Windows in CI due to make/build complexity
	if runtime.GOOS == "windows" && os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("Skipping integration test on Windows in CI")
	}
	
	// Test that --version flag works by running the built binary
	binaryPath := getBinaryPath()
	
	// Build binary if it doesn't exist
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		buildCmd := exec.Command("make", "-f", "Makefile.dev", "build")
		buildCmd.Dir = "../"
		if err := buildCmd.Run(); err != nil {
			t.Fatalf("Failed to build binary: %v", err)
		}
	}
	
	// Use relative path from the binary
	cmd := exec.Command(binaryPath, "--version")

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
	// Skip integration tests on Windows in CI due to make/build complexity
	if runtime.GOOS == "windows" && os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("Skipping integration test on Windows in CI")
	}
	
	// Test that -v flag works
	binaryPath := getBinaryPath()
	
	// Build binary if it doesn't exist
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		buildCmd := exec.Command("make", "-f", "Makefile.dev", "build")
		buildCmd.Dir = "../"
		if err := buildCmd.Run(); err != nil {
			t.Fatalf("Failed to build binary: %v", err)
		}
	}
	
	cmd := exec.Command(binaryPath, "-v")

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
	// Skip integration tests on Windows in CI due to make/build complexity
	if runtime.GOOS == "windows" && os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("Skipping integration test on Windows in CI")
	}
	
	// Skip this test if no kubeconfig is available
	if _, err := os.Stat(os.Getenv("HOME") + "/.kube/config"); os.IsNotExist(err) {
		t.Skip("Skipping test: no kubeconfig found")
	}

	binaryPath := getBinaryPath()
	
	// Build binary if it doesn't exist
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		buildCmd := exec.Command("make", "-f", "Makefile.dev", "build")
		buildCmd.Dir = "../"
		if err := buildCmd.Run(); err != nil {
			t.Fatalf("Failed to build binary: %v", err)
		}
	}

	// Test that 'version' command (without --) still shows Kubernetes info
	cmd := exec.Command(binaryPath, "version", "--help")

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
