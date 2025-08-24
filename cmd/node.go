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
	NodeGroupFilter string
	NodeNameFilter  string
	OutputFormat    string
	SortBy          string
)

// nodeCmd represents the node command
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "获取节点的资源信息（支持过滤和排序）",
	Long: `获取节点的资源信息，支持多种过滤和排序选项：

过滤选项：
- 按节点组名称过滤
- 按节点名称过滤

排序选项：
- 按CPU使用率排序
- 按内存使用率排序
- 按节点名称排序
- 按节点组名称排序

输出格式：
- table: 表格格式（默认）
- json: JSON格式
- yaml: YAML格式`,
	Example: `# 获取所有节点信息
./k8s-manager node

# 按节点组名称过滤
./k8s-manager node --nodegroup my-nodegroup

# 按节点名称过滤
./k8s-manager node --node-name worker-node-1

# 按CPU使用率排序
./k8s-manager node --sort-by cpu

# 按内存使用率排序
./k8s-manager node --sort-by memory

# 输出JSON格式
./k8s-manager node --output json

# 组合使用
./k8s-manager node --nodegroup my-nodegroup --sort-by cpu --output json`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		
		// 创建过滤选项
		filterOptions := &kube.NodeFilterOptions{
			NodeGroupName: NodeGroupFilter,
			NodeName:      NodeNameFilter,
			OutputFormat:  OutputFormat,
			SortBy:        SortBy,
		}
		
		kube.GetNodeInfoWithFilter(ctx, Kubeconfig, Export, filterOptions)
	},
}

func init() {
	// 添加过滤和排序相关的标志
	nodeCmd.Flags().StringVarP(&NodeGroupFilter, "nodegroup", "g", "", "按节点组名称过滤")
	nodeCmd.Flags().StringVarP(&NodeNameFilter, "node-name", "n", "", "按节点名称过滤")
	nodeCmd.Flags().StringVarP(&OutputFormat, "output", "o", "table", "输出格式 (table|json|yaml)")
	nodeCmd.Flags().StringVarP(&SortBy, "sort-by", "s", "", "排序字段 (cpu|memory|name|nodegroup)")
}
