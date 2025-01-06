package kube

import (
	"context"

	"k8s-manager/pkg/mtable"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func GetWorkloadImage(ctx context.Context, kubeconfig, workload, namespace string) {
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
					deployMap["资源类型"] = "deployment"
					deployMap["资源名"] = values.Name
					deployMap["镜像地址"] = values.Spec.Template.Spec.Containers[i].Image
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
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
					deployMap["资源类型"] = "statefulsets"
					deployMap["资源名"] = values.Name
					deployMap["镜像地址"] = values.Spec.Template.Spec.Containers[i].Image
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
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
					deployMap["资源类型"] = "daemonsets"
					deployMap["资源名"] = values.Name
					deployMap["镜像地址"] = values.Spec.Template.Spec.Containers[i].Image
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
					ItemList = append(ItemList, deployMap)
				}
			}

		} else if namespace == "all" {
			deployItems, err := client.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, values := range deployItems.Items {
				for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
					deployMap := make(map[string]string)
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["资源类型"] = "deployment"
					deployMap["资源名"] = values.Name
					deployMap["镜像地址"] = values.Spec.Template.Spec.Containers[i].Image
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
					ItemList = append(ItemList, deployMap)
				}
			}

			stsItems, err := client.AppsV1().StatefulSets("").List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, values := range stsItems.Items {
				for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
					deployMap := make(map[string]string)
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["资源类型"] = "statefulsets"
					deployMap["资源名"] = values.Name
					deployMap["镜像地址"] = values.Spec.Template.Spec.Containers[i].Image
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
					ItemList = append(ItemList, deployMap)
				}
			}

			dsItems, err := client.AppsV1().DaemonSets("").List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, values := range dsItems.Items {
				for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
					deployMap := make(map[string]string)
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["资源类型"] = "daemonsets"
					deployMap["资源名"] = values.Name
					deployMap["镜像地址"] = values.Spec.Template.Spec.Containers[i].Image
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
					ItemList = append(ItemList, deployMap)
				}
			}

		}

	case "deployment", "deploy":
		if namespace != "all" {
			items, err := client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, values := range items.Items {
				for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
					deployMap := make(map[string]string)
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["资源类型"] = "deployment"
					deployMap["资源名"] = values.Name
					deployMap["镜像地址"] = values.Spec.Template.Spec.Containers[i].Image
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
					ItemList = append(ItemList, deployMap)
				}
			}
		} else if namespace == "all" {
			items, err := client.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, values := range items.Items {
				for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
					deployMap := make(map[string]string)
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["资源类型"] = "deployment"
					deployMap["资源名"] = values.Name
					deployMap["镜像地址"] = values.Spec.Template.Spec.Containers[i].Image
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
					ItemList = append(ItemList, deployMap)
				}
			}

		}
	case "sts", "statefulsets":
		if namespace != "all" {
			items, err := client.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, values := range items.Items {
				for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
					deployMap := make(map[string]string)
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["资源类型"] = "statefulsets"
					deployMap["资源名"] = values.Name
					deployMap["镜像地址"] = values.Spec.Template.Spec.Containers[i].Image
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
					ItemList = append(ItemList, deployMap)
				}
			}
		} else if namespace == "all" {
			items, err := client.AppsV1().StatefulSets("").List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, values := range items.Items {
				for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
					deployMap := make(map[string]string)
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["资源类型"] = "statefulsets"
					deployMap["资源名"] = values.Name
					deployMap["镜像地址"] = values.Spec.Template.Spec.Containers[i].Image
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
					ItemList = append(ItemList, deployMap)
				}
			}

		}
	case "ds", "daemonsets":
		if namespace != "all" {
			items, err := client.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, values := range items.Items {
				for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
					deployMap := make(map[string]string)
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["资源类型"] = "daemonsets"
					deployMap["资源名"] = values.Name
					deployMap["镜像地址"] = values.Spec.Template.Spec.Containers[i].Image
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
					ItemList = append(ItemList, deployMap)
				}
			}
		} else if namespace == "all" {
			items, err := client.AppsV1().DaemonSets("").List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			for _, values := range items.Items {
				for i := 0; i < len(values.Spec.Template.Spec.Containers); i++ {
					deployMap := make(map[string]string)
					deployMap["NAMESPACE"] = values.Namespace
					deployMap["资源类型"] = "daemonsets"
					deployMap["资源名"] = values.Name
					deployMap["镜像地址"] = values.Spec.Template.Spec.Containers[i].Image
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
					ItemList = append(ItemList, deployMap)
				}
			}

		}
	}

	mtable.TablePrint("image", ItemList)
}
