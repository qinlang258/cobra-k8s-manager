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
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		kube.GetWorkloadLimitRequests(ctx, Kubeconfig, Workload, Namespace, Name)
	},
}

func init() {
	//rootCmd.AddCommand(resourceCmd)

}
