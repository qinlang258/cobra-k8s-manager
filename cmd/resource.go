/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"k8s-manager/pkg/kube"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

var (
	Prometheus bool
)

func getPrometheusUrl(ctx context.Context, file string) string {
	// 设置 Viper 配置文件路径和类型
	viper.SetConfigType("yaml")
	viper.SetConfigFile("config/prometheus.yaml")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		klog.Error(ctx, "Failed to read config: "+err.Error())
		return ""
	}

	// 获取 prometheus 配置部分
	var prometheusConfig []map[string]interface{}
	err = viper.UnmarshalKey("prometheus", &prometheusConfig)
	if err != nil {
		klog.Error(ctx, "Failed to unmarshal prometheus config: "+err.Error())
		return ""
	}

	// 遍历配置并查找匹配的 kubeconfig
	for _, item := range prometheusConfig {
		kubeconfig := item["kubeconfig"].(string)
		if kubeconfig == file {
			// 获取对应的 URL
			url := item["url"].(string)
			return url
		}
	}

	// 如果没有找到匹配的 kubeconfig
	return ""
}

// resourceCmd represents the resource command
var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "获取pod资源的相关 Limit与Resource信息",
	Long:  "获取pod资源的相关 Limit与Resource信息",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		var url string
		if Prometheus {
			url = getPrometheusUrl(ctx, Kubeconfig)
		}
		//读取 Prometheus地址

		//如果需要输出prometheus实际开销
		if Prometheus != false && Node == "" {
			kube.AnalysisResourceAndLimitWithNamespace(ctx, Kubeconfig, Workload, Namespace, url)
		} else if Prometheus != false && Node != "" {
			kube.AnalysisResourceAndLimitWithNode(ctx, Kubeconfig, Workload, Namespace, Node, url)

		} else if Prometheus == false && Node == "" {
			kube.GetWorkloadLimitRequests(ctx, Kubeconfig, Workload, Namespace, Name)
		}
	},
}

/*
可以供执行命令主机访问的 prometheus的地址 http://192.168.44.134:20248/
*/

func init() {
	//rootCmd.AddCommand(resourceCmd)
	resourceCmd.PersistentFlags().StringVarP(&Node, "node", "", "", "请输入想要查询的Node名字")
	//resourceCmd.PersistentFlags().StringVarP(&Prometheus, "url", "u", "", "需要分析的话，填写prometheus的地址，仅支持当前集群的状态查询")
	resourceCmd.PersistentFlags().BoolVarP(&Prometheus, "prometheus", "p", false, "需要分析的话，在配置文件填写prometheus的地址，仅支持当前集群的状态查询")
	resourceCmd.PersistentFlags().StringVarP(&Namespace, "namespace", "n", "all", "请输入 namespace空间，如果不填写则输出所有空间下的镜像")
	resourceCmd.PersistentFlags().StringVarP(&Workload, "workload", "", "all", "请输入 workload的种类，如果不填写输出所有类型的镜像")
}
