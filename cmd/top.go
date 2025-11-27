/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"k8s-manager/pkg/kube"

	"github.com/spf13/cobra"
)

// topCmd represents the top command
var topCmd = &cobra.Command{
	Use:   "top",
	Short: "获取容器的实际使用资源开销",
	Long:  "获取容器的实际使用资源开销",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		if Namespace == "" && Node == "" {
			kube.GetPodTopInfoWithCurrentNamespace(ctx, Kubeconfig, LabelSelector, Export)
		}

		if Namespace != "" && Node == "" {
			kube.GetPodTopInfoWithNamespace(ctx, Kubeconfig, Workload, Namespace, LabelSelector, Export)
		}
		if Node != "" && Namespace == "" {
			kube.GetPodTopInfoWithNode(ctx, Kubeconfig, Workload, Node, LabelSelector, Export)
		}
		if Node != "" && Namespace != "" {
			kube.GetPodTopInfoWithNamespaceAndNode(ctx, Kubeconfig, Workload, Node, Namespace, LabelSelector, Export)
		}

	},
}

func init() {
	//rootCmd.AddCommand(topCmd)
	topCmd.PersistentFlags().StringVarP(&Namespace, "namespace", "n", "", "请输入 namespace空间，如果不填写则输出所有空间下的镜像")
	topCmd.PersistentFlags().StringVarP(&Workload, "workload", "", "all", "请输入 workload的种类，如果不填写输出所有类型的镜像")
	topCmd.PersistentFlags().StringVarP(&Node, "node", "", "", "请输入 节点名空间，如果不填写则输出该Node下的镜像")
	topCmd.PersistentFlags().StringVarP(&LabelSelector, "selector", "l", "", "标签选择器，例如: app=nginx,version=v1")
}
