package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var cfgFile string

var (
	cliVersion string
	gitCommit  string
	buildTime  string
	goVersion  string
)

var rootCmd = &cobra.Command{
	Use:   "k8s-cli",
	Short: "A CLI tool for Kubernetes cluster analysis and monitoring",
	Long: `k8s-cli is a comprehensive command-line tool that helps you analyze 
your Kubernetes cluster by providing information about:
- Kubernetes version and installed components
- Resource consumption (cluster, nodes, pods)
- Recommendations for optimization`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if version, _ := cmd.Flags().GetBool("version"); version {
			fmt.Printf("k8s-cli version %s\n", cliVersion)
			fmt.Printf("Git commit: %s\n", gitCommit)
			fmt.Printf("Built: %s\n", buildTime)
			fmt.Printf("Go version: %s\n", goVersion)
			fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
			return nil
		}
		return cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func SetVersionInfo(version, commit, buildT, goVer string) {
	cliVersion = version
	gitCommit = commit
	buildTime = buildT
	goVersion = goVer
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.k8s-cli.yaml)")
	rootCmd.PersistentFlags().String("kubeconfig", "", "path to kubeconfig file (default is $HOME/.kube/config)")
	rootCmd.Flags().BoolP("version", "v", false, "Show CLI version")
}

func initConfig() {
	if cfgFile != "" {
		fmt.Printf("Using config file: %s\n", cfgFile)
	}
}
