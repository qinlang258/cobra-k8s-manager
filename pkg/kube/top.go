package kube

import (
	"context"
	"fmt"
	"k8s-manager/pkg/mtable"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func GetPodTopInfoWithNamespaceAndNode(ctx context.Context, kubeconfig, workload, node, namespace string) {
	ItemList := make([]map[string]string, 0)

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

	deploymentItems, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", node),
	})
	if err != nil {
		klog.Error(ctx, err.Error())
		return
	}
	for _, values := range deploymentItems.Items {
		if values.OwnerReferences[0].Kind == "StatefulSet" || values.OwnerReferences[0].Kind == "DaemonSet" || values.OwnerReferences[0].Kind == "ReplicaSet" {
			deployMap := make(map[string]string)
			deployMap["节点名"] = values.Spec.NodeName
			deployMap["NAMESPACE"] = values.Namespace
			deployMap["资源类型"] = values.OwnerReferences[0].Kind
			deployMap["资源名"] = values.OwnerReferences[0].Name
			deployMap["POD_NAME"] = values.Name

			podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(values.Namespace).Get(ctx, values.Name, metav1.GetOptions{})
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

				deployMap["当前已使用的CPU"] = fmt.Sprintf("%.2fm", usedCpuCores)
				deployMap["当前已使用的内存"] = fmt.Sprintf("%.2fm", usedMemoryMi)
			}

			ItemList = append(ItemList, deployMap)

		}
		mtable.TablePrint("top", ItemList)
	}
}

func GetPodTopInfoWithNode(ctx context.Context, kubeconfig, workload, node string) {
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
		if node != "all" {
			deploymentItems, err := client.CoreV1().Pods("").List(ctx, metav1.ListOptions{
				FieldSelector: fmt.Sprintf("spec.nodeName=%s", node),
			})
			if err != nil {
				klog.Error(ctx, err.Error())
				return
			}
			for _, values := range deploymentItems.Items {
				if values.OwnerReferences[0].Kind == "StatefulSet" || values.OwnerReferences[0].Kind == "DaemonSet" || values.OwnerReferences[0].Kind == "ReplicaSet" {
					deployMap := make(map[string]string)
					deployMap["节点名"] = values.Spec.NodeName
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["资源类型"] = values.OwnerReferences[0].Kind
					deployMap["资源名"] = values.OwnerReferences[0].Name
					deployMap["POD_NAME"] = values.Name

					podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(values.Namespace).Get(ctx, values.Name, metav1.GetOptions{})
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

						deployMap["当前已使用的CPU"] = fmt.Sprintf("%.2fm", usedCpuCores)
						deployMap["当前已使用的内存"] = fmt.Sprintf("%.2fm", usedMemoryMi)
					}

					ItemList = append(ItemList, deployMap)
				}
			}

		} else if node == "all" {
			deploymentItems, err := client.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, "Error listing pods in namespace ", err.Error())
			}
			for _, values := range deploymentItems.Items {
				if values.OwnerReferences[0].Kind == "StatefulSet" || values.OwnerReferences[0].Kind == "DaemonSet" || values.OwnerReferences[0].Kind == "ReplicaSet" {
					deployMap := make(map[string]string)
					deployMap["节点名"] = values.Spec.NodeName
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["资源类型"] = values.OwnerReferences[0].Kind
					deployMap["资源名"] = values.OwnerReferences[0].Name
					deployMap["POD_NAME"] = values.Name

					// 使用正确的命名空间
					podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(values.Namespace).Get(ctx, values.Name, metav1.GetOptions{})
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

						deployMap["当前已使用的CPU"] = fmt.Sprintf("%.2fm", usedCpuCores)
						deployMap["当前已使用的内存"] = fmt.Sprintf("%.2fMi", usedMemoryMi)
					}

					ItemList = append(ItemList, deployMap)
				}
			}

		}
	}

	mtable.TablePrint("top", ItemList)
}

func GetPodTopInfoWithNamespace(ctx context.Context, kubeconfig, workload, namespace string) {
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
					deployMap["节点名"] = values.Spec.NodeName
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["资源类型"] = values.OwnerReferences[0].Kind
					deployMap["资源名"] = values.OwnerReferences[0].Name
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

						deployMap["当前已使用的CPU"] = fmt.Sprintf("%.2fm", usedCpuCores)
						deployMap["当前已使用的内存"] = fmt.Sprintf("%.2fm", usedMemoryMi)
					}

					ItemList = append(ItemList, deployMap)
				}
			}

		} else if namespace == "all" {
			deploymentItems, err := client.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, "Error listing pods in namespace ", err.Error())
			}
			for _, values := range deploymentItems.Items {
				if values.OwnerReferences[0].Kind == "StatefulSet" || values.OwnerReferences[0].Kind == "DaemonSet" || values.OwnerReferences[0].Kind == "ReplicaSet" {
					deployMap := make(map[string]string)
					deployMap["节点名"] = values.Spec.NodeName
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["资源类型"] = values.OwnerReferences[0].Kind
					deployMap["资源名"] = values.OwnerReferences[0].Name
					deployMap["POD_NAME"] = values.Name

					// 使用正确的命名空间
					podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(values.Namespace).Get(ctx, values.Name, metav1.GetOptions{})
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

						deployMap["当前已使用的CPU"] = fmt.Sprintf("%.2fm", usedCpuCores)
						deployMap["当前已使用的内存"] = fmt.Sprintf("%.2fMi", usedMemoryMi)
					}

					ItemList = append(ItemList, deployMap)
				}
			}

		}
	}

	mtable.TablePrint("top", ItemList)
}
