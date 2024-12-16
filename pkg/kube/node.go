package kube

import (
	"context"
	"fmt"
	"k8s-manager/pkg/mtable"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func GetNodeInfo(ctx context.Context, nodeName, kubeconfig string) {
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

	nodeList, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Error(ctx, err.Error())
	}

	ItemList := make([]map[string]string, 0)

	for _, values := range nodeList.Items {
		deployMap := make(map[string]string)
		deployMap["NODE_NAME"] = values.Name
		deployMap["OS_IMAGE"] = values.Status.NodeInfo.OSImage

		deployMap["KUBELET_VERSION"] = values.Status.NodeInfo.KubeletVersion
		deployMap["CONTAINER_RUNTIME_VERSION"] = values.Status.NodeInfo.ContainerRuntimeVersion
		deployMap["NODE_ADDRESS"] = values.Status.Addresses[0].Address

		nodeMetrics, err := metricsClient.MetricsV1beta1().NodeMetricses().Get(ctx, values.Name, metav1.GetOptions{})
		if err != nil {
			klog.Error(ctx, err.Error())
		}

		// 获取 CPU 和内存使用数据
		cpuUsage := nodeMetrics.Usage[corev1.ResourceCPU]
		memoryUsage := nodeMetrics.Usage[corev1.ResourceMemory]

		// 将内存从 Ki 转换为 Mi
		usedMemoryMi := float64(memoryUsage.Value()) / 1024 / 1024
		totalMemoryMi := float64(values.Status.Capacity.Memory().Value()) / 1024 / 1024

		// 转换 CPU 使用量为毫核心数
		usedCpuCores := float64(cpuUsage.MilliValue())
		totalCpuCores := values.Status.Capacity.Cpu().MilliValue()

		deployMap["CPU_USED"] = fmt.Sprintf("%.2fm", usedCpuCores)
		deployMap["CPU_TOTAL"] = fmt.Sprintf("%dm", totalCpuCores)
		deployMap["CPU_PERCENT"] = fmt.Sprintf("%.2f%%", (usedCpuCores/float64(totalCpuCores))*100)
		deployMap["MEMORY_USED"] = fmt.Sprintf("%.2fMi", usedMemoryMi)
		deployMap["MEMORY_TOTAL"] = fmt.Sprintf("%.2fMi", totalMemoryMi)
		deployMap["MEMORY_PERCENT"] = fmt.Sprintf("%.2f%%", (usedMemoryMi/totalMemoryMi)*100)
		ItemList = append(ItemList, deployMap)
	}

	mtable.TablePrint("node", ItemList)

}
