/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"k8s-manager/pkg/kube"

	"github.com/spf13/cobra"
)

// nodeCmd represents the node command
var nodeCmd = &cobra.Command{
	Use:     "node",
	Short:   "获取节点的资源信息",
	Long:    "获取节点的资源信息",
	Example: "./k8s-manager node --kubeconfig <可选，配置文件地址，默认/root/.kube/config> -n <查询的namespace空间>",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		kube.GetNodeInfo(ctx, Node, Kubeconfig)
	},
}

func init() {
}
