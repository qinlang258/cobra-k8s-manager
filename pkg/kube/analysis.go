package kube

import (
	"context"
	"fmt"
	"k8s-manager/pkg/mtable"
	"sort"

	"strconv"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func AnalysisNode(ctx context.Context, kubeconfig, nodeName string) {
	client, err := NewClientset(kubeconfig)
	if err != nil {
		klog.Error(ctx, err.Error())
		return
	}

	metricsClient, err := NewMetricsClient(kubeconfig)
	if err != nil {
		klog.Error(ctx, err.Error())
		return
	}

	ItemList := make([]map[string]string, 0)

	nodeData, err := client.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		klog.Error(ctx, err.Error())
	}

	totalMemoryMi := float64(nodeData.Status.Capacity.Memory().Value()) / 1024 / 1024
	totalCpuCores := nodeData.Status.Capacity.Cpu().MilliValue()

	podList, _ := client.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	for _, pod := range podList.Items {
		if pod.Spec.NodeName == nodeName {
			deployMap := make(map[string]string)
			deployMap["NODE_NAME"] = nodeName
			deployMap["NAMESPACE"] = pod.Namespace
			deployMap["POD_NAME"] = pod.Name

			// 获取 Pod Metrics
			podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(pod.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
			if err != nil {
				klog.Error(ctx, fmt.Sprintf("Error fetching metrics for pod %s: %v", pod.Name, err))
				// 继续处理其他 Pod，即使当前 Pod 的 metrics 获取失败
				continue
			}

			// 检查 Pod Metrics 是否存在容器数据
			if len(podMetrics.Containers) > 0 {
				// 获取 CPU 和内存使用数据
				for i := 0; i < len(podMetrics.Containers); i++ {
					cpuUsage := podMetrics.Containers[i].Usage.Cpu()
					memoryUsage := podMetrics.Containers[i].Usage.Memory()
					usedMemoryMi := float64(memoryUsage.Value()) / 1024 / 1024
					usedCpuCores := float64(cpuUsage.MilliValue())

					deployMap["CPU_USED"] = fmt.Sprintf("%.2fm", usedCpuCores)
					deployMap["CPU_PERCENT"] = fmt.Sprintf("%.2f%%", (usedCpuCores/float64(totalCpuCores))*100)
					deployMap["MEMORY_USED"] = fmt.Sprintf("%.2fm", usedMemoryMi)
					deployMap["MEMORY_PERCENT"] = fmt.Sprintf("%.2f%%", (usedMemoryMi/totalMemoryMi)*100)
				}
			} else {
				klog.Warning(ctx, fmt.Sprintf("Pod %s has no container metrics", pod.Name))
				deployMap["CPU_USED"] = "N/A"
				deployMap["MEMORY_USED"] = "N/A"
			}

			ItemList = append(ItemList, deployMap)
		}
	}

	// 按照 MEMORY_USED 倒序排列
	sort.Slice(ItemList, func(i, j int) bool {
		// 获取 MEMORY_USED 的数值部分
		memI, errI := parseMemory(ItemList[i]["MEMORY_USED"])
		memJ, errJ := parseMemory(ItemList[j]["MEMORY_USED"])

		// 如果其中一个值是 "N/A" 或无法解析，则视为最小值
		if errI != nil || errJ != nil {
			return memI > memJ // "N/A" 处理为最小值
		}
		return memI > memJ
	})

	// 打印排序后的结果
	mtable.TablePrint("analysis", ItemList)
}

// 辅助函数：解析 MEMORY_USED 字符串，返回数值
func parseMemory(memStr string) (float64, error) {
	// 去除 "m" 后缀并解析为浮动值
	if strings.HasSuffix(memStr, "m") {
		memStr = strings.TrimSuffix(memStr, "m")
	}
	return strconv.ParseFloat(memStr, 64)
}
