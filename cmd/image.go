/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"k8s-manager/pkg/kube"

	"github.com/spf13/cobra"
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "获取镜像信息",
	Long:  "获取容器镜像的信息",
	Example: `
	./k8s-manager image   默认获取所有空间的image信息
	./k8s-manager image -n	monitoring 获取monitoring空间下的所有镜像
	./k8s-manager image --kubeconfig /root/.kube/k8s-ops-zjk-aliyun.yaml 使用指定的配置文件链接对应的集群
	./k8s-manager image --workload deploy || sts || ds获取deployment sts ds的镜像
	`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		// if name != "" && namespace == "" && workload == "" {
		// 	cmd.Help()
		// 	return
		// }
		kube.GetWorkloadImage(ctx, Kubeconfig, Workload, Namespace)
	},
}

func init() {
	//rootCmd.AddCommand(imageCmd)
}
