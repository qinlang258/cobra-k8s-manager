/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"k8s-manager/pkg/kube"

	"github.com/spf13/cobra"
)

var (
	Prometheus string
)

// resourceCmd represents the resource command
var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "获取pod资源的相关 Limit与Resource信息",
	Long:  "获取pod资源的相关 Limit与Resource信息",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		//如果需要输出prometheus实际开销
		if Prometheus != "" && Node == "" {
			kube.AnalysisResourceAndLimitWithNamespace(ctx, Kubeconfig, Workload, Namespace, Prometheus)
		} else if Prometheus != "" && Node != "" {
			kube.AnalysisResourceAndLimitWithNode(ctx, Kubeconfig, Workload, Namespace, Node, Prometheus)

		} else if Prometheus == "" && Node == "" {
			kube.GetWorkloadLimitRequests(ctx, Kubeconfig, Workload, Namespace, Name)
		}
	},
}

/*
可以供执行命令主机访问的 prometheus的地址 http://192.168.44.134:20248/
*/

func init() {
	//rootCmd.AddCommand(resourceCmd)
	resourceCmd.PersistentFlags().StringVarP(&Prometheus, "url", "u", "", "需要分析的话，填写prometheus的地址，仅支持当前集群的状态查询")
}
