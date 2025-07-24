package kube

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"k8s-manager/pkg/excel"
	"k8s-manager/pkg/mtable"
	"k8s-manager/pkg/prometheusplugin"
	"k8s-manager/tools"

	prov1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

const (
	//RSS内存
	RssMemoryUsageTemplate = "avg(avg_over_time(container_memory_rss {pod=\"%s\",  container=\"%s\", namespace=\"%s\"}[7d])) / 1024 / 1024"
	podCpuUsageTemplate    = "sum(irate(container_cpu_usage_seconds_total{pod=\"%s\",  container=\"%s\", namespace=\"%s\"}[7d])) * 1000"
)

func FormatData(result model.Value, warnings prov1.Warnings, err error) string {
	var num_data string

	if err != nil {
		fmt.Println("prometheus没有获取到数据,请检查Prometheus是否能正常访问?")
		return ""
	}

	if result.String() == "" {
		return "0"
	}

	data := result.String()

	//提取 => 0.0031342189920170885 @[1711701880.602
	s1 := strings.Split(data, "=>")
	s2 := strings.Split(s1[1], "@")
	num_data = strings.ReplaceAll(s2[0], " ", "")

	return num_data
}

func GetWorkloadLimitRequests(ctx context.Context, kubeconfig, workload, namespace, name string, export bool) {
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
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
					deployMap["CPU限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
					deployMap["内存限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
					deployMap["CPU所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
					deployMap["内存所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

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
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
					deployMap["CPU限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
					deployMap["内存限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
					deployMap["CPU所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
					deployMap["内存所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

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
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
					deployMap["CPU限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
					deployMap["内存限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
					deployMap["CPU所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
					deployMap["内存所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

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
						deployMap["资源类型"] = "deployment"
						deployMap["资源名"] = values.Name
						deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
						deployMap["CPU限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
						deployMap["内存限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
						deployMap["CPU所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
						deployMap["内存所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

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
						deployMap["资源类型"] = "statefulsets"
						deployMap["资源名"] = values.Name
						deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
						deployMap["CPU限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
						deployMap["内存限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
						deployMap["CPU所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
						deployMap["内存所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

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
						deployMap["资源类型"] = "daemonsets"
						deployMap["资源名"] = values.Name
						deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
						deployMap["CPU限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
						deployMap["内存限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
						deployMap["CPU所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
						deployMap["内存所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

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
					deployMap["资源类型"] = "deployment"
					deployMap["资源名"] = values.Name
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
					deployMap["CPU限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
					deployMap["内存限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
					deployMap["CPU所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
					deployMap["内存所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()
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
						deployMap["资源类型"] = "deployment"
						deployMap["资源名"] = values.Name
						deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
						deployMap["CPU限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
						deployMap["内存限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
						deployMap["CPU所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
						deployMap["内存所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

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
					deployMap["资源类型"] = "statefulsets"
					deployMap["资源名"] = values.Name
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
					deployMap["CPU限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
					deployMap["内存限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
					deployMap["CPU所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
					deployMap["内存所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

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
						deployMap["资源类型"] = "statefulsets"
						deployMap["资源名"] = values.Name
						deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
						deployMap["CPU限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
						deployMap["内存限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
						deployMap["CPU所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
						deployMap["内存所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

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
					deployMap["资源类型"] = "daemonsets"
					deployMap["资源名"] = values.Name
					deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
					deployMap["CPU限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
					deployMap["内存限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
					deployMap["CPU所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
					deployMap["内存所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

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
						deployMap["资源类型"] = "daemonsets"
						deployMap["资源名"] = values.Name
						deployMap["容器名"] = values.Spec.Template.Spec.Containers[i].Name
						deployMap["CPU限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().String()
						deployMap["内存限制"] = values.Spec.Template.Spec.Containers[i].Resources.Limits.Memory().String()
						deployMap["CPU所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().String()
						deployMap["内存所需"] = values.Spec.Template.Spec.Containers[i].Resources.Requests.Memory().String()

						ItemList = append(ItemList, deployMap)
					}
				}
			}
		}
	}

	mtable.TablePrint("resource", ItemList)

	if export {
		excel.ExportXlsx(ctx, "resource", ItemList, kubeconfig)
	}
}

func AnalysisResourceAndLimitWithNamespace(ctx context.Context, kubeconfig, workload, namespace, prometheusUrl string, export bool) {
	client, err := NewClientset(kubeconfig)
	if err != nil {
		klog.Error(ctx, err.Error())
		return
	}

	prometheus_client := prometheusplugin.NewProme(prometheusUrl, 10)
	ItemList := make([]map[string]string, 0)

	switch namespace {
	case "all":
		nsList, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			klog.Error(ctx, err.Error())
		}

		for _, ns := range nsList.Items {
			podList, err := client.CoreV1().Pods(ns.Name).List(ctx, metav1.ListOptions{})
			if err != nil {
				klog.Error(ctx, err.Error())
			}
			//获取所有pod的资源信息
			for _, pod := range podList.Items {
				for i := 0; i < len(pod.Spec.Containers); i++ {
					deployMap := make(map[string]string)
					cpuLimit := pod.Spec.Containers[i].Resources.Limits.Cpu().String()
					cpuRequets := pod.Spec.Containers[i].Resources.Requests.Cpu().String()
					memoryLimit := pod.Spec.Containers[i].Resources.Limits.Memory().String()
					memoryRequests := pod.Spec.Containers[i].Resources.Requests.Memory().String()

					nodename := pod.Spec.NodeName

					memorySql := fmt.Sprintf(RssMemoryUsageTemplate, pod.Name, pod.Spec.Containers[i].Name, pod.Namespace)
					cpuSql := fmt.Sprintf(podCpuUsageTemplate, pod.Name, pod.Spec.Containers[i].Name, pod.Namespace)

					strmemorySize := FormatData(prometheus_client.Client.Query(ctx, memorySql, time.Now()))
					memorySize, err1 := strconv.ParseFloat(strmemorySize, 64)
					if err1 != nil {
						klog.Error(ctx, err1.Error())
					}

					strcpuSize := FormatData(prometheus_client.Client.Query(ctx, cpuSql, time.Now()))
					cpuSize, err := strconv.ParseFloat(strcpuSize, 64)
					if err != nil {
						klog.Error(ctx, err.Error())
					}

					deployMap["节点名"] = nodename
					deployMap["节点组名称"], err = tools.GetNodeGroupNameFromPod(ctx, client, pod)
					if err != nil {
						klog.Error(ctx, fmt.Sprintf("Error getting node group name for pod %s: %v", pod.Name, err))
					}
					deployMap["NAMESPACE"] = pod.Namespace
					deployMap["POD_NAME"] = pod.Name
					deployMap["容器名"] = pod.Spec.Containers[i].Name
					deployMap["CPU限制"] = cpuLimit
					deployMap["CPU所需"] = cpuRequets
					deployMap["最近7天已使用的CPU"] = fmt.Sprintf("%.2fm", cpuSize)
					deployMap["内存限制"] = memoryLimit
					deployMap["内存所需"] = memoryRequests
					deployMap["最近7天已使用的内存"] = fmt.Sprintf("%.2fMi", memorySize)

					ItemList = append(ItemList, deployMap)

				}

			}
		}
	default:
		podList, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			klog.Error(ctx, err.Error())
		}
		//获取所有pod的资源信息
		for _, pod := range podList.Items {
			for i := 0; i < len(pod.Spec.Containers); i++ {
				deployMap := make(map[string]string)
				cpuLimit := pod.Spec.Containers[i].Resources.Limits.Cpu().String()
				cpuRequets := pod.Spec.Containers[i].Resources.Requests.Cpu().String()
				memoryLimit := pod.Spec.Containers[i].Resources.Limits.Memory().String()
				memoryRequests := pod.Spec.Containers[i].Resources.Requests.Memory().String()
				nodename := pod.Spec.NodeName

				memorySql := fmt.Sprintf(RssMemoryUsageTemplate, pod.Name, pod.Spec.Containers[i].Name, pod.Namespace)
				cpuSql := fmt.Sprintf(podCpuUsageTemplate, pod.Name, pod.Spec.Containers[i].Name, pod.Namespace)

				strmemorySize := FormatData(prometheus_client.Client.Query(ctx, memorySql, time.Now()))
				memorySize, err1 := strconv.ParseFloat(strmemorySize, 64)
				if err1 != nil {
					klog.Error(ctx, err1.Error())
				}

				strcpuSize := FormatData(prometheus_client.Client.Query(ctx, cpuSql, time.Now()))
				cpuSize, err := strconv.ParseFloat(strcpuSize, 64)
				if err != nil {
					// Handle error if conversion fails
					klog.Error(ctx, err.Error())
				}

				var xmx, xms string

				item, _ := client.CoreV1().Pods(pod.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
				for _, container := range item.Spec.Containers {
					for _, env := range container.Env {
						if env.Name == "JAVA_OPTS" {
							for _, v := range strings.Fields(env.Value) {
								if strings.Contains(v, "Xmx") {
									xmx = v
								} else if strings.Contains(v, "Xms") {
									xms = v
								}

							}
						}
					}
				}

				deployMap["节点名"] = nodename
				deployMap["节点组名称"], err = tools.GetNodeGroupNameFromPod(ctx, client, pod)
				if err != nil {
					klog.Error(ctx, fmt.Sprintf("Error getting node group name for pod %s: %v", pod.Name, err))
				}
				deployMap["NAMESPACE"] = namespace
				deployMap["POD_NAME"] = pod.Name
				deployMap["容器名"] = pod.Spec.Containers[i].Name
				deployMap["CPU限制"] = cpuLimit
				deployMap["CPU所需"] = cpuRequets
				deployMap["最近7天已使用的CPU"] = fmt.Sprintf("%.2fm", cpuSize)
				deployMap["内存限制"] = memoryLimit
				deployMap["内存所需"] = memoryRequests
				deployMap["JAVA-XMX"] = xmx
				deployMap["JAVA-XMS"] = xms
				deployMap["最近7天已使用的内存"] = fmt.Sprintf("%.2fMi", memorySize)

				ItemList = append(ItemList, deployMap)

			}
		}
	}

	mtable.TablePrint("analysis-cpu-memory", ItemList)

	if export {
		excel.ExportXlsx(ctx, "analysis-cpu-memory", ItemList, kubeconfig)
	}

}

func TestPrometheus(ctx context.Context, pod_name, container_name, namespace, prometheusUrl string) {
	prometheus_client := prometheusplugin.NewProme(prometheusUrl, 10)

	//cpuSql := fmt.Sprintf(podCpuUsageTemplate, pod_name, container_name, namespace)
	cpuSql := "rate(container_cpu_usage_seconds_total{container!='',pod='prometheus-k8s-0'}[7d])"

	result, _, err := prometheus_client.Client.Query(ctx, cpuSql, time.Now())
	if err != nil {
		klog.Error(ctx, err.Error())
	}

	fmt.Println(result)

}

func AnalysisResourceAndLimitWithNode(ctx context.Context, kubeconfig, workload, namespace, node, prometheusUrl string, export bool) {
	client, err := NewClientset(kubeconfig)
	if err != nil {
		klog.Error(ctx, err.Error())
		return
	}

	prometheus_client := prometheusplugin.NewProme(prometheusUrl, 10)
	ItemList := make([]map[string]string, 0)

	if namespace != "all" {

		podsList, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
			FieldSelector: fmt.Sprintf("spec.nodeName=%s", node),
		})

		if err != nil {
			klog.Error(ctx, err.Error())
		}

		//获取所有pod的资源信息
		for _, pod := range podsList.Items {
			for i := 0; i < len(pod.Spec.Containers); i++ {
				deployMap := make(map[string]string)
				cpuLimit := pod.Spec.Containers[i].Resources.Limits.Cpu().String()
				cpuRequets := pod.Spec.Containers[i].Resources.Requests.Cpu().String()
				memoryLimit := pod.Spec.Containers[i].Resources.Limits.Memory().String()
				memoryRequests := pod.Spec.Containers[i].Resources.Requests.Memory().String()
				nodename := pod.Spec.NodeName

				var xmx, xms string

				item, _ := client.CoreV1().Pods(pod.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
				for _, container := range item.Spec.Containers {
					for _, env := range container.Env {
						if env.Name == "JAVA_OPTS" {
							for _, v := range strings.Fields(env.Value) {
								if strings.Contains(v, "Xmx") {
									xmx = v
								} else if strings.Contains(v, "Xms") {
									xms = v
								}

							}
						}
					}
				}

				memorySql := fmt.Sprintf(RssMemoryUsageTemplate, pod.Name, pod.Spec.Containers[i].Name, pod.Namespace)
				cpuSql := fmt.Sprintf(podCpuUsageTemplate, pod.Name, pod.Spec.Containers[i].Name, pod.Namespace)

				strmemorySize := FormatData(prometheus_client.Client.Query(ctx, memorySql, time.Now()))
				memorySize, err1 := strconv.ParseFloat(strmemorySize, 64)
				if err1 != nil {
					klog.Error(ctx, err1.Error())
				}

				strcpuSize := FormatData(prometheus_client.Client.Query(ctx, cpuSql, time.Now()))
				cpuSize, err := strconv.ParseFloat(strcpuSize, 64)
				if err != nil {
					// Handle error if conversion fails
					klog.Error(ctx, err.Error())
				}

				deployMap["节点名"] = nodename
				deployMap["节点组名称"], err = tools.GetNodeGroupNameFromPod(ctx, client, pod)
				if err != nil {
					klog.Error(ctx, fmt.Sprintf("Error getting node group name for pod %s: %v", pod.Name, err))
				}
				deployMap["NAMESPACE"] = pod.Namespace
				deployMap["POD_NAME"] = pod.Name
				deployMap["容器名"] = pod.Spec.Containers[i].Name
				deployMap["CPU限制"] = cpuLimit
				deployMap["CPU所需"] = cpuRequets
				deployMap["最近7天已使用的CPU"] = fmt.Sprintf("%.2fm", cpuSize)
				deployMap["内存限制"] = memoryLimit
				deployMap["内存所需"] = memoryRequests
				deployMap["JAVA-XMX"] = xmx
				deployMap["JAVA-XMS"] = xms
				deployMap["最近7天已使用的内存"] = fmt.Sprintf("%.2fMi", memorySize)

				ItemList = append(ItemList, deployMap)

			}
		}
	} else {
		nsList, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			klog.Error(ctx, err.Error())
		}

		for _, ns := range nsList.Items {
			podsList, err := client.CoreV1().Pods(ns.Name).List(ctx, metav1.ListOptions{
				FieldSelector: fmt.Sprintf("spec.nodeName=%s", node),
			})

			if err != nil {
				klog.Error(ctx, err.Error())
			}

			//获取所有pod的资源信息
			for _, pod := range podsList.Items {
				for i := 0; i < len(pod.Spec.Containers); i++ {
					deployMap := make(map[string]string)
					cpuLimit := pod.Spec.Containers[i].Resources.Limits.Cpu().String()
					cpuRequets := pod.Spec.Containers[i].Resources.Requests.Cpu().String()
					memoryLimit := pod.Spec.Containers[i].Resources.Limits.Memory().String()
					memoryRequests := pod.Spec.Containers[i].Resources.Requests.Memory().String()
					nodename := pod.Spec.NodeName
					var xmx, xms string

					item, _ := client.CoreV1().Pods(pod.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
					for _, container := range item.Spec.Containers {
						for _, env := range container.Env {
							if env.Name == "JAVA_OPTS" {
								for _, v := range strings.Fields(env.Value) {
									if strings.Contains(v, "Xmx") {
										xmx = v
									} else if strings.Contains(v, "Xms") {
										xms = v
									}

								}
							}
						}
					}

					memorySql := fmt.Sprintf(RssMemoryUsageTemplate, pod.Name, pod.Spec.Containers[i].Name, pod.Namespace)
					cpuSql := fmt.Sprintf(podCpuUsageTemplate, pod.Name, pod.Spec.Containers[i].Name, pod.Namespace)

					strmemorySize := FormatData(prometheus_client.Client.Query(ctx, memorySql, time.Now()))
					memorySize, err1 := strconv.ParseFloat(strmemorySize, 64)
					if err1 != nil {
						klog.Error(ctx, err1.Error())
					}

					strcpuSize := FormatData(prometheus_client.Client.Query(ctx, cpuSql, time.Now()))
					cpuSize, err := strconv.ParseFloat(strcpuSize, 64)
					if err != nil {
						// Handle error if conversion fails
						klog.Error(ctx, err.Error())
					}

					deployMap["节点名"] = nodename
					deployMap["节点组名称"], err = tools.GetNodeGroupNameFromPod(ctx, client, pod)
					if err != nil {
						klog.Error(ctx, fmt.Sprintf("Error getting node group name for pod %s: %v", pod.Name, err))
					}
					deployMap["NAMESPACE"] = ns.Name
					deployMap["POD_NAME"] = pod.Name
					deployMap["容器名"] = pod.Spec.Containers[i].Name
					deployMap["CPU限制"] = cpuLimit
					deployMap["CPU所需"] = cpuRequets
					deployMap["最近7天已使用的CPU"] = fmt.Sprintf("%.2fm", cpuSize)
					deployMap["内存限制"] = memoryLimit
					deployMap["内存所需"] = memoryRequests
					deployMap["JAVA-XMX"] = xmx
					deployMap["JAVA-XMS"] = xms
					deployMap["最近7天已使用的内存"] = fmt.Sprintf("%.2fMi", memorySize)

					ItemList = append(ItemList, deployMap)

				}
			}
		}
	}

	mtable.TablePrint("analysis-cpu-memory", ItemList)

	if export {
		excel.ExportXlsx(ctx, "analysis-cpu-memory", ItemList, kubeconfig)
	}
}
