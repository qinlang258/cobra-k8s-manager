package prometheusplugin

import (
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"

	"context"
	"fmt"
	"time"
)

const (
	nodeMeasureQueryTemplate = "sum_over_time(node_network_receive_bytes_total{device=\"%s\"}[%ss]) * on(instance) group_left(nodename) (node_uname_info{nodename=\"%s\"})"
	podMemoryUsageTemplate   = "container_memory_working_set_bytes{container=\"%s\"}"
	podCpuUsageTemplate      = "container_cpu_usage_seconds_total{container=\"%s\"}"
)

type PrometheusHandle struct {
	timeRange time.Duration
	ip        string
	Client    v1.API
}

func NewProme(ip string, timeRace time.Duration) *PrometheusHandle {
	client, err := api.NewClient(api.Config{Address: ip})
	if err != nil {
		klog.Fatalf("[NetworkTraffic Plugin] FatalError creating prometheus client: %s", err.Error())
	}
	return &PrometheusHandle{
		ip:        ip,
		timeRange: timeRace,
		Client:    v1.NewAPI(client),
	}
}

func GetPrometheusUrl(ctx context.Context, file string) string {
	// 获取当前用户的家目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		klog.Error(ctx, "Failed to get user home directory: "+err.Error())
	}

	// 定义目标文件路径
	targetDir := filepath.Join(homeDir, ".kube", "jcrose-prometheus")
	targetFile := filepath.Join(targetDir, "prometheus.yaml")

	// 设置 Viper 配置文件路径和类型
	viper.SetConfigType("yaml")
	viper.SetConfigFile(targetFile)

	// 读取配置文件
	err = viper.ReadInConfig()
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

func (p *PrometheusHandle) GetCpuUsage(container string) (*model.Sample, error) {
	value, err := p.query(fmt.Sprintf(podCpuUsageTemplate, container))
	//fmt.Println(fmt.Sprintf(podCpuUsageTemplate, container))
	if err != nil {
		return nil, fmt.Errorf("[NetworkTraffic Plugin] Error querying prometheus: %w", err)
	}

	cpuMeasure := value.(model.Vector)
	if len(cpuMeasure) != 1 {
		return nil, fmt.Errorf("[NetworkTraffic Plugin] Invalid response, expected 1 value, got %d", len(cpuMeasure))
	}
	return cpuMeasure[0], err
}

func (p *PrometheusHandle) GetMemoryUsage(container string) (*model.Sample, error) {
	value, err := p.query(fmt.Sprintf(podMemoryUsageTemplate, container))
	//fmt.Println(fmt.Sprintf(podMemoryUsageTemplate, container))
	if err != nil {
		return nil, fmt.Errorf("[NetworkTraffic Plugin] Error querying prometheus: %w", err)
	}

	memoryMeasure := value.(model.Vector)
	if len(memoryMeasure) != 1 {
		return nil, fmt.Errorf("[NetworkTraffic Plugin] Invalid response, expected 1 value, got %d", len(memoryMeasure))
	}
	return memoryMeasure[0], err
}

func (p *PrometheusHandle) query(promQL string) (model.Value, error) {
	results, warnings, err := p.Client.Query(context.Background(), promQL, time.Now())
	if len(warnings) > 0 {
		klog.Warningf("[prometheus] Warnings: %v\n", warnings)
	}

	return results, err
}
