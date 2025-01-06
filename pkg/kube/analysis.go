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

func AnalysisNodeWithNode(ctx context.Context, kubeconfig, nodeName string) {
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
					deployMap := make(map[string]string)
					deployMap["节点名"] = nodeName
					deployMap["NAMESPACE"] = pod.Namespace
					deployMap["POD_NAME"] = pod.Name
					deployMap["容器名"] = podMetrics.Containers[i].Name
					cpuUsage := podMetrics.Containers[i].Usage.Cpu()
					memoryUsage := podMetrics.Containers[i].Usage.Memory()
					usedMemoryMi := float64(memoryUsage.Value()) / 1024 / 1024
					usedCpuCores := float64(cpuUsage.MilliValue())

					deployMap["当前已使用的CPU"] = fmt.Sprintf("%.2fm", usedCpuCores)
					deployMap["CPU使用占服务器的百分比"] = fmt.Sprintf("%.2f%%", (usedCpuCores/float64(totalCpuCores))*100)
					deployMap["当前已使用的内存"] = fmt.Sprintf("%.2fm", usedMemoryMi)
					deployMap["内存使用占服务器的百分比"] = fmt.Sprintf("%.2f%%", (usedMemoryMi/totalMemoryMi)*100)
					ItemList = append(ItemList, deployMap)
				}
			}

		}
	}

	// 按照 当前已使用的内存内存 倒序排列
	sort.Slice(ItemList, func(i, j int) bool {
		// 获取 当前已使用的内存内存 的数值部分
		memI, errI := parseMemory(ItemList[i]["当前已使用的内存"])
		memJ, errJ := parseMemory(ItemList[j]["当前已使用的内存"])

		// 如果其中一个值是 "N/A" 或无法解析，则视为最小值
		if errI != nil || errJ != nil {
			return memI > memJ // "N/A" 处理为最小值
		}
		return memI > memJ
	})

	// 打印排序后的结果
	mtable.TablePrint("analysis", ItemList)
}

// 辅助函数：解析 当前已使用的内存内存 字符串，返回数值
func parseMemory(memStr string) (float64, error) {
	// 去除 "m" 后缀并解析为浮动值
	if strings.HasSuffix(memStr, "m") {
		memStr = strings.TrimSuffix(memStr, "m")
	}
	return strconv.ParseFloat(memStr, 64)
}
