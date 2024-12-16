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
		kube.GetPodTopInfo(ctx, Kubeconfig, Workload, Namespace)
	},
}

func init() {
	//rootCmd.AddCommand(topCmd)

}
