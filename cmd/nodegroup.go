/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s-manager/pkg/kube"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

var (
	NodeGroupOutputFormat string
)

// nodegroupCmd represents the nodegroup command
var nodegroupCmd = &cobra.Command{
	Use:   "nodegroup",
	Short: "节点组管理和统计",
	Long: `节点组管理和统计功能，支持：
- 列出所有节点组
- 显示节点组统计信息
- 按节点组过滤节点`,
	Example: `# 显示所有节点组统计
./k8s-manager nodegroup stats

# 以 JSON 格式输出节点组信息
./k8s-manager nodegroup stats --output json

# 列出指定节点组的节点
./k8s-manager nodegroup list my-nodegroup`,
}

// nodegroupStatsCmd represents the nodegroup stats command
var nodegroupStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "显示节点组统计信息",
	Long:  "显示所有节点组的统计信息，包括节点数量、总CPU和内存资源",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		showNodeGroupStats(ctx)
	},
}

// nodegroupListCmd represents the nodegroup list command
var nodegroupListCmd = &cobra.Command{
	Use:   "list [nodegroup-name]",
	Short: "列出节点组中的节点",
	Long:  "列出指定节点组中的所有节点信息",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		
		if len(args) == 0 {
			fmt.Println("请指定节点组名称")
			return
		}
		
		nodeGroupName := args[0]
		listNodesInNodeGroup(ctx, nodeGroupName)
	},
}

func showNodeGroupStats(ctx context.Context) {
	stats, err := kube.GetNodeGroupStats(ctx, Kubeconfig)
	if err != nil {
		fmt.Printf("获取节点组统计信息失败: %v\n", err)
		return
	}

	switch NodeGroupOutputFormat {
	case "json":
		jsonData, err := json.MarshalIndent(stats, "", "  ")
		if err != nil {
			fmt.Printf("JSON 序列化失败: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	case "yaml":
		yamlData, err := yaml.Marshal(stats)
		if err != nil {
			fmt.Printf("YAML 序列化失败: %v\n", err)
			return
		}
		fmt.Println(string(yamlData))
	default:
		// 表格输出
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"节点组名称", "节点数量", "总CPU(毫核)", "总内存(Mi)"})
		
		for _, stat := range stats {
			table.Append([]string{
				stat.NodeGroupName,
				fmt.Sprintf("%d", stat.NodeCount),
				fmt.Sprintf("%.0f", stat.TotalCPU),
				fmt.Sprintf("%.0f", stat.TotalMemory),
			})
		}
		
		table.Render()
	}
}

func listNodesInNodeGroup(ctx context.Context, nodeGroupName string) {
	nodeInfos, err := kube.GetNodesByNodeGroup(ctx, Kubeconfig, nodeGroupName)
	if err != nil {
		fmt.Printf("获取节点组 %s 的节点信息失败: %v\n", nodeGroupName, err)
		return
	}

	if len(nodeInfos) == 0 {
		fmt.Printf("节点组 %s 中没有找到节点\n", nodeGroupName)
		return
	}

	switch NodeGroupOutputFormat {
	case "json":
		jsonData, err := json.MarshalIndent(nodeInfos, "", "  ")
		if err != nil {
			fmt.Printf("JSON 序列化失败: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	case "yaml":
		yamlData, err := yaml.Marshal(nodeInfos)
		if err != nil {
			fmt.Printf("YAML 序列化失败: %v\n", err)
			return
		}
		fmt.Println(string(yamlData))
	default:
		// 表格输出
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"节点名", "节点IP", "CPU使用率", "内存使用率", "Kubelet版本",
		})
		
		for _, nodeInfo := range nodeInfos {
			table.Append([]string{
				nodeInfo.NodeName,
				nodeInfo.NodeIP,
				nodeInfo.CPUUsagePercent,
				nodeInfo.MemoryUsagePercent,
				nodeInfo.KubeletVersion,
			})
		}
		
		table.Render()
	}
}

func init() {
	// 添加子命令
	nodegroupCmd.AddCommand(nodegroupStatsCmd)
	nodegroupCmd.AddCommand(nodegroupListCmd)
	
	// 添加标志
	nodegroupCmd.PersistentFlags().StringVarP(&NodeGroupOutputFormat, "output", "o", "table", "输出格式 (table|json|yaml)")
	
	// 将命令添加到根命令
	rootCmd.AddCommand(nodegroupCmd)
}
