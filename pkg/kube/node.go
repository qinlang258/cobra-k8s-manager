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
		deployMap["节点名"] = values.Name
		deployMap["OS镜像"] = values.Status.NodeInfo.OSImage

		deployMap["Kubelet版本"] = values.Status.NodeInfo.KubeletVersion
		deployMap["CONTAINER_RUNTIME_VERSION"] = values.Status.NodeInfo.ContainerRuntimeVersion
		deployMap["节点IP"] = values.Status.Addresses[0].Address

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

		deployMap["当前已使用的CPU"] = fmt.Sprintf("%.2fm", usedCpuCores)
		deployMap["CPU总大小"] = fmt.Sprintf("%dm", totalCpuCores)
		deployMap["CPU_PERCENT"] = fmt.Sprintf("%.2f%%", (usedCpuCores/float64(totalCpuCores))*100)
		deployMap["当前已使用的内存"] = fmt.Sprintf("%.2fMi", usedMemoryMi)
		deployMap["内存总大小"] = fmt.Sprintf("%.2fMi", totalMemoryMi)
		deployMap["内存使用占服务器的百分比"] = fmt.Sprintf("%.2f%%", (usedMemoryMi/totalMemoryMi)*100)
		ItemList = append(ItemList, deployMap)
	}

	mtable.TablePrint("node", ItemList)

}
