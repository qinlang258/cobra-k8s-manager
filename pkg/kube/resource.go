package kube

import (
	"context"

	"k8s-manager/pkg/mtable"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func GetWorkloadLimitRequests(ctx context.Context, kubeconfig, workload, namespace, name string) {
	client, err := NewClientset(kubeconfig)
	if err != nil {
		klog.Error(err)
		return
	}

	ItemList := make([]map[string]string, 0)
	switch workload {
	case "all":
		if namespace != "all" {
			deploymentLtems, err := client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, values := range deploymentLtems.Items {
				for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
					deployMap := make(map[string]string)
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["TYPE"] = "deployment"
					deployMap["RESOURCE_NAME"] = values.Name
					deployMap["CONTAINER_NAME"] = values.Spec.Template.Spec.Containers[i].Name
					deployMap["CPU_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
					deployMap["MEMORY_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
					deployMap["CPU_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
					deployMap["MEMORY_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

					ItemList = append(ItemList, deployMap)
				}
			}

			stsItems, err := client.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, values := range stsItems.Items {
				for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
					deployMap := make(map[string]string)
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["TYPE"] = "statefulsets"
					deployMap["RESOURCE_NAME"] = values.Name
					deployMap["CONTAINER_NAME"] = values.Spec.Template.Spec.Containers[i].Name
					deployMap["CPU_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
					deployMap["MEMORY_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
					deployMap["CPU_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
					deployMap["MEMORY_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

					ItemList = append(ItemList, deployMap)
				}
			}

			dsItems, err := client.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, values := range dsItems.Items {
				for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
					deployMap := make(map[string]string)
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["TYPE"] = "daemonsets"
					deployMap["RESOURCE_NAME"] = values.Name
					deployMap["CONTAINER_NAME"] = values.Spec.Template.Spec.Containers[i].Name
					deployMap["CPU_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
					deployMap["MEMORY_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
					deployMap["CPU_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
					deployMap["MEMORY_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

					ItemList = append(ItemList, deployMap)
				}
			}

		} else if namespace == "all" {
			nsList, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, ns := range nsList.Items {
				deployItems, err := client.AppsV1().Deployments(ns.Name).List(ctx, metav1.ListOptions{})
				if err != nil {
					klog.Error(ctx, err.Error())
				}
				for _, values := range deployItems.Items {
					for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
						deployMap := make(map[string]string)
						deployMap["NAMESPACE"] = values.Namespace
						deployMap["TYPE"] = "deployment"
						deployMap["RESOURCE_NAME"] = values.Name
						deployMap["CONTAINER_NAME"] = values.Spec.Template.Spec.Containers[i].Name
						deployMap["CPU_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
						deployMap["MEMORY_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
						deployMap["CPU_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
						deployMap["MEMORY_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

						ItemList = append(ItemList, deployMap)
					}
				}

				stsItems, err := client.AppsV1().StatefulSets(ns.Name).List(ctx, metav1.ListOptions{})
				if err != nil {
					klog.Error(ctx, err.Error())
				}
				for _, values := range stsItems.Items {
					for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
						deployMap := make(map[string]string)
						deployMap["NAMESPACE"] = values.Namespace
						deployMap["TYPE"] = "statefulsets"
						deployMap["RESOURCE_NAME"] = values.Name
						deployMap["CONTAINER_NAME"] = values.Spec.Template.Spec.Containers[i].Name
						deployMap["CPU_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
						deployMap["MEMORY_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
						deployMap["CPU_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
						deployMap["MEMORY_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

						ItemList = append(ItemList, deployMap)
					}
				}

				dsItems, err := client.AppsV1().DaemonSets(ns.Name).List(ctx, metav1.ListOptions{})
				if err != nil {
					klog.Error(ctx, err.Error())
				}
				for _, values := range dsItems.Items {
					for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
						deployMap := make(map[string]string)
						deployMap["NAMESPACE"] = values.Namespace
						deployMap["TYPE"] = "daemonsets"
						deployMap["RESOURCE_NAME"] = values.Name
						deployMap["CONTAINER_NAME"] = values.Spec.Template.Spec.Containers[i].Name
						deployMap["CPU_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
						deployMap["MEMORY_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
						deployMap["CPU_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
						deployMap["MEMORY_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

						ItemList = append(ItemList, deployMap)
					}
				}
			}

		}

	case "deployment":
		if namespace != "all" {
			items, err := client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, values := range items.Items {
				for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
					deployMap := make(map[string]string)
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["TYPE"] = "deployment"
					deployMap["RESOURCE_NAME"] = values.Name
					deployMap["CONTAINER_NAME"] = values.Spec.Template.Spec.Containers[i].Name
					deployMap["CPU_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
					deployMap["MEMORY_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
					deployMap["CPU_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
					deployMap["MEMORY_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()
					ItemList = append(ItemList, deployMap)
				}
			}
		} else if namespace == "all" {
			nsList, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, ns := range nsList.Items {
				items, err := client.AppsV1().Deployments(ns.Name).List(ctx, metav1.ListOptions{})
				if err != nil {
					klog.Error(ctx, err.Error())
				}
				for _, values := range items.Items {
					for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
						deployMap := make(map[string]string)
						deployMap["NAMESPACE"] = values.Namespace
						deployMap["TYPE"] = "deployment"
						deployMap["RESOURCE_NAME"] = values.Name
						deployMap["CONTAINER_NAME"] = values.Spec.Template.Spec.Containers[i].Name
						deployMap["CPU_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
						deployMap["MEMORY_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
						deployMap["CPU_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
						deployMap["MEMORY_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

						ItemList = append(ItemList, deployMap)
					}
				}
			}
		}
	case "sts":
		if namespace != "all" {
			items, err := client.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, values := range items.Items {
				for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
					deployMap := make(map[string]string)
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["TYPE"] = "statefulsets"
					deployMap["RESOURCE_NAME"] = values.Name
					deployMap["CONTAINER_NAME"] = values.Spec.Template.Spec.Containers[i].Name
					deployMap["CPU_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
					deployMap["MEMORY_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
					deployMap["CPU_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
					deployMap["MEMORY_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

					ItemList = append(ItemList, deployMap)
				}
			}
		} else if namespace == "all" {
			nsList, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, ns := range nsList.Items {
				items, err := client.AppsV1().StatefulSets(ns.Name).List(ctx, metav1.ListOptions{})
				if err != nil {
					klog.Error(ctx, err.Error())
				}
				for _, values := range items.Items {
					for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
						deployMap := make(map[string]string)
						deployMap["NAMESPACE"] = values.Namespace
						deployMap["TYPE"] = "statefulsets"
						deployMap["RESOURCE_NAME"] = values.Name
						deployMap["CONTAINER_NAME"] = values.Spec.Template.Spec.Containers[i].Name
						deployMap["CPU_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
						deployMap["MEMORY_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
						deployMap["CPU_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
						deployMap["MEMORY_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

						ItemList = append(ItemList, deployMap)
					}
				}
			}
		}
	case "ds":
		if namespace != "all" {
			items, err := client.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, values := range items.Items {
				for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
					deployMap := make(map[string]string)
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["TYPE"] = "daemonsets"
					deployMap["RESOURCE_NAME"] = values.Name
					deployMap["CONTAINER_NAME"] = values.Spec.Template.Spec.Containers[i].Name
					deployMap["CPU_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
					deployMap["MEMORY_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
					deployMap["CPU_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
					deployMap["MEMORY_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

					ItemList = append(ItemList, deployMap)
				}
			}
		} else if namespace == "all" {
			nsList, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, ns := range nsList.Items {
				items, err := client.AppsV1().DaemonSets(ns.Name).List(ctx, metav1.ListOptions{})
				if err != nil {
					klog.Error(ctx, err.Error())
				}
				for _, values := range items.Items {
					for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
						deployMap := make(map[string]string)
						deployMap["NAMESPACE"] = values.Namespace
						deployMap["TYPE"] = "daemonsets"
						deployMap["RESOURCE_NAME"] = values.Name
						deployMap["CONTAINER_NAME"] = values.Spec.Template.Spec.Containers[i].Name
						deployMap["CPU_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
						deployMap["MEMORY_LIMIT"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
						deployMap["CPU_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
						deployMap["MEMORY_REQUESTS"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

						ItemList = append(ItemList, deployMap)
					}
				}
			}
		}
	}

	mtable.TablePrint("resource", ItemList)
}
