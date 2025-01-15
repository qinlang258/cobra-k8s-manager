/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"k8s-manager/pkg/config"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "初始化prometheus配置文件",
	Long:  `初始化prometheus配置文件`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		config.InitPrometheus(ctx, KubeconfigPath)
	},
}

func init() {
	configCmd.PersistentFlags().StringVarP(&KubeconfigPath, "kubeconfig_path", "p", "/root/.kube/", "默认位置在当前家目录的.kube文件夹下")

}
