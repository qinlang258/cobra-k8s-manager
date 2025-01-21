/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"k8s-manager/pkg/kube"
	"k8s-manager/pkg/prometheusplugin"

	"github.com/spf13/cobra"
)

var (
	Prometheus bool
)

// resourceCmd represents the resource command
var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "获取pod资源的相关 Limit与Resource信息",
	Long:  "获取pod资源的相关 Limit与Resource信息",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		var url string
		if Prometheus {
			url = prometheusplugin.GetPrometheusUrl(ctx, Kubeconfig)
			fmt.Println("所使用的prometheus地址是： ", url)
		}
		//读取 Prometheus地址

		//如果需要输出prometheus实际开销
		if Prometheus != false && Node == "" {
			kube.AnalysisResourceAndLimitWithNamespace(ctx, Kubeconfig, Workload, Namespace, url, Export)
		} else if Prometheus != false && Node != "" {
			kube.AnalysisResourceAndLimitWithNode(ctx, Kubeconfig, Workload, Namespace, Node, url, Export)

		} else if Prometheus == false && Node == "" {
			kube.GetWorkloadLimitRequests(ctx, Kubeconfig, Workload, Namespace, Name, Export)
		}
	},
}

func init() {
	//rootCmd.AddCommand(resourceCmd)
	resourceCmd.PersistentFlags().StringVarP(&Node, "node", "", "", "请输入想要查询的Node名字")
	//resourceCmd.PersistentFlags().StringVarP(&Prometheus, "url", "u", "", "需要分析的话，填写prometheus的地址，仅支持当前集群的状态查询")
	resourceCmd.PersistentFlags().BoolVarP(&Prometheus, "prometheus", "p", false, "需要分析的话，在配置文件填写prometheus的地址，仅支持当前集群的状态查询")
	resourceCmd.PersistentFlags().StringVarP(&Namespace, "namespace", "n", "all", "请输入 namespace空间，如果不填写则输出所有空间下的镜像")
	resourceCmd.PersistentFlags().StringVarP(&Workload, "workload", "", "all", "请输入 workload的种类，如果不填写输出所有类型的镜像")
}
