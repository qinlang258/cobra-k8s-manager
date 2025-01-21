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
	Workload, Namespace, Name, Kubeconfig, Node, Analysis, KubeconfigPath string
	Export                                                                bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-manager",
	Short: "获取K8S的资源使用情况",
	Long: `jcrose的k8s管理工具，功能涵盖如下
	1 analysis 分析某一节点的资源使用情况，类似于describe node xxx
	2 config 生成 /root/.kube/jcrose-prometheus/prometheus.yaml 配置文件，获取默认路径/root/.kube/下的yaml文件导入prometheus的地址
	3 image 获取镜像地址
	4 node 获取所有node节点的信息
	5 resource 获取pod所使用的 limit与requests 与java_opts的环境变量，可以加上-p获取对应集群的最近7天平均内存和CPU数值
	6 获取容器的实际使用资源开销`,
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
	rootCmd.AddCommand(configCmd)

	rootCmd.PersistentFlags().StringVarP(&Kubeconfig, "kubeconfig", "", "/root/.kube/config", "请输入 kubeconfig的文件路径")
	rootCmd.PersistentFlags().BoolVarP(&Export, "export", "", false, "是否输出Excel?默认不输出")
}
