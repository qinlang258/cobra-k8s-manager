/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"k8s-manager/pkg/kube"

	"github.com/spf13/cobra"
)

// resourceCmd represents the resource command
var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "获取pod资源的相关 Limit与Resource信息",
	Long:  "获取pod资源的相关 Limit与Resource信息",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		kube.GetWorkloadLimitRequests(ctx, Kubeconfig, Workload, Namespace, Name)
	},
}

func init() {
	//rootCmd.AddCommand(resourceCmd)

}
