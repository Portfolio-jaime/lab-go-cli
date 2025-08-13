package main

import (
	"fmt"
	"os"
	"runtime"

	"k8s-cli/cmd"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
	GoVersion = runtime.Version()
)

func main() {
	cmd.SetVersionInfo(Version, GitCommit, BuildTime, GoVersion)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
