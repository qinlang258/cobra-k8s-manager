/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

// var cfgFile string
var command string

var (
	Workload, Namespace, Name, Kubeconfig, Node, Analysis string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-manager",
	Short: "获取K8S的资源使用情况",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if len(command) == 0 {
			cmd.Help()
			klog.Error("必须输入想要查询的内容")
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(imageCmd)
	rootCmd.AddCommand(resourceCmd)
	rootCmd.AddCommand(topCmd)
	rootCmd.AddCommand(nodeCmd)
	rootCmd.AddCommand(analysisCmd)

	rootCmd.PersistentFlags().StringVarP(&Kubeconfig, "kubeconfig", "", "/root/.kube/config", "请输入 kubeconfig的文件路径")
	rootCmd.PersistentFlags().StringVarP(&Name, "name", "", "", "请输入资源的name信息")
	rootCmd.PersistentFlags().StringVarP(&Analysis, "analysis", "", "", "请输入想分析的Node名字")

}
