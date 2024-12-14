/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// var cfgFile string
var command string

var (
	Workload, Namespace, Name, Kubeconfig string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-manager",
	Short: "A brief description of your application",
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

	rootCmd.PersistentFlags().StringVarP(&Kubeconfig, "kubeconfig", "", "/root/.kube/config", "请输入 kubeconfig的文件路径")
	rootCmd.PersistentFlags().StringVarP(&Namespace, "namespace", "n", "all", "请输入 namespace空间，如果不填写则输出所有空间下的镜像")
	rootCmd.PersistentFlags().StringVarP(&Workload, "workload", "", "all", "请输入 workload的种类，如果不填写输出所有类型的镜像")
	rootCmd.PersistentFlags().StringVarP(&Name, "name", "", "", "请输入资源的name信息")
}
