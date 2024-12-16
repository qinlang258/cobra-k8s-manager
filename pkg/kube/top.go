package kube

import (
	"context"
	"fmt"
	"k8s-manager/pkg/mtable"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func GetPodTopInfo(ctx context.Context, kubeconfig, workload, namespace string) {
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
			deploymentItems, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
				return
			}
			for _, values := range deploymentItems.Items {
				if values.OwnerReferences[0].Kind == "StatefulSet" || values.OwnerReferences[0].Kind == "DaemonSet" || values.OwnerReferences[0].Kind == "ReplicaSet" {
					deployMap := make(map[string]string)
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["TYPE"] = values.OwnerReferences[0].Kind
					deployMap["RESOURCE_NAME"] = values.OwnerReferences[0].Name
					deployMap["POD_NAME"] = values.Name

					podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).Get(ctx, values.Name, metav1.GetOptions{})
					if err != nil {
						klog.Error(ctx, "Error fetching metrics for pod ", values.Name, ": ", err.Error())
						continue
					}

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

		} else if namespace == "all" {
			nsList, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, "Error listing namespaces: ", err.Error())
				return
			}
			for _, ns := range nsList.Items {
				deploymentItems, err := client.CoreV1().Pods(ns.Name).List(ctx, metav1.ListOptions{})
				if err != nil {
					klog.Error(ctx, "Error listing pods in namespace ", ns.Name, ": ", err.Error())
					continue
				}
				for _, values := range deploymentItems.Items {
					if values.OwnerReferences[0].Kind == "StatefulSet" || values.OwnerReferences[0].Kind == "DaemonSet" || values.OwnerReferences[0].Kind == "ReplicaSet" {
						deployMap := make(map[string]string)
						deployMap["NAMESPACE"] = values.Namespace
						deployMap["TYPE"] = values.OwnerReferences[0].Kind
						deployMap["RESOURCE_NAME"] = values.OwnerReferences[0].Name
						deployMap["POD_NAME"] = values.Name

						// 使用正确的命名空间
						podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(ns.Name).Get(ctx, values.Name, metav1.GetOptions{})
						if err != nil {
							klog.Error(ctx, "Error fetching metrics for pod ", values.Name, ": ", err.Error())
							continue
						}

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
		}
	}

	mtable.TablePrint("top", ItemList)
}
