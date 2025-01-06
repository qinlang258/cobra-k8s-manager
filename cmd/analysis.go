/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"k8s-manager/pkg/kube"

	"github.com/spf13/cobra"
)

// analysisCmd represents the analysis command
var analysisCmd = &cobra.Command{
	Use:     "analysis",
	Short:   "分析某一节点的资源使用情况",
	Long:    "分析某一节点的资源使用情况",
	Example: "./k8s-manager analysis <节点名>",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		kube.AnalysisNode(ctx, Kubeconfig, Node)
	},
}

func init() {
	analysisCmd.PersistentFlags().StringVarP(&Node, "node", "", "", "请输入想要查询的Node名字")

}
