package kube

import (
	"context"
	"fmt"
	"k8s-manager/pkg/mtable"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func GetPodTopInfo(ctx context.Context, kubeconfig, workload, namespace, name string) {
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
	switch workload {
	case "all":
		if namespace != "all" {
			deploymentLtems, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, values := range deploymentLtems.Items {
				deployMap := make(map[string]string)
				deployMap["NAMESPACE"] = values.Namespace
				deployMap["TYPE"] = values.OwnerReferences[0].Kind
				deployMap["RESOURCE_NAME"] = values.OwnerReferences[0].Name
				deployMap["POD_NAME"] = values.Name

				podMetrics, _ := metricsClient.MetricsV1beta1().PodMetricses(namespace).Get(ctx, values.Name, metav1.GetOptions{})
				// 获取 CPU 和内存使用数据
				for i := 0; i < len(podMetrics.Containers); i++ {
					cpuUsage := podMetrics.Containers[i].Usage.Cpu()
					memoryUsage := podMetrics.Containers[i].Usage.Memory()
					usedMemoryMi := float64(memoryUsage.Value()) / 1024 / 1024
					usedCpuCores := float64(cpuUsage.MilliValue())

					deployMap["CPU_USED"] = fmt.Sprintf("%.2fm", usedCpuCores)
					deployMap["MEMORY_USED"] = fmt.Sprintf("%.2fm", usedMemoryMi)
				}

				ItemList = append(ItemList, deployMap)

			}

		}
	}

	mtable.TablePrint("top", ItemList)
}
