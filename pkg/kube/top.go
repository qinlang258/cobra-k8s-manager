package kube

import (
	"context"
	"fmt"
	"k8s-manager/pkg/excel"
	"k8s-manager/pkg/mtable"
	"k8s-manager/tools"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func GetPodTopInfoWithNamespaceAndNode(ctx context.Context, kubeconfig, workload, node, namespace string, export bool) {
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
			deployMap["节点组名称"], err = tools.GetNodeGroupNameFromPod(ctx, client, values)
			if err != nil {
				klog.Error(ctx, fmt.Sprintf("Error getting node group name for pod %s: %v", values.Name, err))
			}
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

	mtable.TablePrint("top", ItemList)
	if export {
		if export {
			excel.ExportXlsx(ctx, "top", ItemList, kubeconfig)
		}
	}

}

func GetPodTopInfoWithNode(ctx context.Context, kubeconfig, workload, node string, export bool) {
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
					deployMap["节点组名称"], err = tools.GetNodeGroupNameFromPod(ctx, client, values)
					if err != nil {
						klog.Error(ctx, fmt.Sprintf("Error getting node group name for pod %s: %v", values.Name, err))
					}
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
					deployMap["节点组名称"], err = tools.GetNodeGroupNameFromPod(ctx, client, values)
					if err != nil {
						klog.Error(ctx, fmt.Sprintf("Error getting node group name for pod %s: %v", values.Name, err))
					}

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
	if export {
		if export {
			excel.ExportXlsx(ctx, "top", ItemList, kubeconfig)
		}
	}
}

func GetPodTopInfoWithCurrentNamespace(ctx context.Context, kubeconfig string, export bool) {
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
	namespace := GetClientgoNamespace(kubeconfig)

	podListItem, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Error(ctx, err.Error())
	}

	for _, pod := range podListItem.Items {
		deployMap := make(map[string]string)
		deployMap["节点名"] = pod.Spec.NodeName
		deployMap["节点组名称"], err = tools.GetNodeGroupNameFromPod(ctx, client, pod)
		if err != nil {
			klog.Error(ctx, fmt.Sprintf("Error getting node group name for pod %s: %v", pod.Name, err))
		}
		deployMap["NAMESPACE"] = pod.Namespace
		deployMap["资源类型"] = pod.OwnerReferences[0].Kind
		deployMap["资源名"] = pod.OwnerReferences[0].Name
		deployMap["POD_NAME"] = pod.Name

		podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(pod.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
		if err != nil {
			klog.Error(ctx, "Error fetching metrics for pod ", pod.Name, ": ", err.Error())
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

func GetPodAllTopInfo(ctx context.Context, kubeconfig string) {
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

	podListItem, err := client.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Error(ctx, err.Error())
	}

	nodeListItem, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Error(ctx, err.Error())
	}

	for _, node := range nodeListItem.Items {
		for _, pod := range podListItem.Items {
			if node.Name == pod.Spec.NodeName {
				deployMap := make(map[string]string)
				deployMap["节点名"] = pod.Spec.NodeName
				deployMap["节点组名称"], err = tools.GetNodeGroupNameFromPod(ctx, client, pod)
				if err != nil {
					klog.Error(ctx, fmt.Sprintf("Error getting node group name for pod %s: %v", pod.Name, err))
				}
				deployMap["NAMESPACE"] = pod.Namespace
				deployMap["资源类型"] = pod.OwnerReferences[0].Kind
				deployMap["资源名"] = pod.OwnerReferences[0].Name
				deployMap["POD_NAME"] = pod.Name

				podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(pod.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
				if err != nil {
					klog.Error(ctx, "Error fetching metrics for pod ", pod.Name, ": ", err.Error())
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
	}

	mtable.TablePrint("top", ItemList)
}

func GetPodTopInfoWithNamespace(ctx context.Context, kubeconfig, workload, namespace string, export bool) {
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
					deployMap["节点组名称"], err = tools.GetNodeGroupNameFromPod(ctx, client, values)
					if err != nil {
						klog.Error(ctx, fmt.Sprintf("Error getting node group name for pod %s: %v", values.Name, err))
					}
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
					deployMap["节点组名称"], err = tools.GetNodeGroupNameFromPod(ctx, client, values)
					if err != nil {
						klog.Error(ctx, fmt.Sprintf("Error getting node group name for pod %s: %v", values.Name, err))
					}
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
	if export {
		excel.ExportXlsx(ctx, "top", ItemList, kubeconfig)
	}
}
