package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "k8s-cli",
	Short: "A CLI tool for Kubernetes cluster analysis and monitoring",
	Long: `k8s-cli is a comprehensive command-line tool that helps you analyze 
your Kubernetes cluster by providing information about:
- Kubernetes version and installed components
- Resource consumption (cluster, nodes, pods)
- Recommendations for optimization`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.k8s-cli.yaml)")
	rootCmd.PersistentFlags().String("kubeconfig", "", "path to kubeconfig file (default is $HOME/.kube/config)")
}

func initConfig() {
	if cfgFile != "" {
		fmt.Printf("Using config file: %s\n", cfgFile)
	}
}