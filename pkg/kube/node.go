package kube

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s-manager/pkg/excel"
	"k8s-manager/pkg/mtable"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"sigs.k8s.io/yaml"
)

// NodeFilterOptions 定义节点过滤选项
type NodeFilterOptions struct {
	NodeGroupName string // 节点组名称过滤
	NodeName      string // 节点名称过滤
	OutputFormat  string // 输出格式 (table|json|yaml)
	SortBy        string // 排序字段 (cpu|memory|name|nodegroup)
}

// NodeInfo 定义节点信息结构
type NodeInfo struct {
	NodeName                string  `json:"node_name" yaml:"node_name"`
	NodeGroupName          string  `json:"node_group_name" yaml:"node_group_name"`
	NodeIP                 string  `json:"node_ip" yaml:"node_ip"`
	OSImage                string  `json:"os_image" yaml:"os_image"`
	KubeletVersion         string  `json:"kubelet_version" yaml:"kubelet_version"`
	ContainerRuntimeVersion string  `json:"container_runtime_version" yaml:"container_runtime_version"`
	UsedCPU                string  `json:"used_cpu" yaml:"used_cpu"`
	TotalCPU               string  `json:"total_cpu" yaml:"total_cpu"`
	CPUUsagePercent        string  `json:"cpu_usage_percent" yaml:"cpu_usage_percent"`
	UsedMemory             string  `json:"used_memory" yaml:"used_memory"`
	TotalMemory            string  `json:"total_memory" yaml:"total_memory"`
	MemoryUsagePercent     string  `json:"memory_usage_percent" yaml:"memory_usage_percent"`
	// 用于排序的数值字段
	CPUUsageValue    float64 `json:"-" yaml:"-"`
	MemoryUsageValue float64 `json:"-" yaml:"-"`
}

// NodeGroupStat 节点组统计信息
type NodeGroupStat struct {
	NodeGroupName string  `json:"node_group_name" yaml:"node_group_name"`
	NodeCount     int     `json:"node_count" yaml:"node_count"`
	TotalCPU      float64 `json:"total_cpu_millicores" yaml:"total_cpu_millicores"`
	TotalMemory   float64 `json:"total_memory_mi" yaml:"total_memory_mi"`
}

// GetNodeInfo 保持原有接口兼容性
func GetNodeInfo(ctx context.Context, nodeName, kubeconfig string, export bool) {
	options := &NodeFilterOptions{
		NodeName:     nodeName,
		OutputFormat: "table",
	}
	GetNodeInfoWithFilter(ctx, kubeconfig, export, options)
}

// GetNodeInfoWithFilter 获取节点信息并支持过滤
func GetNodeInfoWithFilter(ctx context.Context, kubeconfig string, export bool, options *NodeFilterOptions) {
	client, err := NewClientset(kubeconfig)
	if err != nil {
		klog.Error(ctx, err.Error())
		return
	}

	metricsClient, err := NewMetricsClient(kubeconfig)
	if err != nil {
		klog.Warningf("无法创建 metrics 客户端: %v", err)
		klog.Info("将显示节点基本信息，但无法显示资源使用情况")
	}

	// 构建 ListOptions
	listOptions := metav1.ListOptions{}
	
	// 如果指定了节点名称，使用 FieldSelector
	if options.NodeName != "" {
		listOptions.FieldSelector = fmt.Sprintf("metadata.name=%s", options.NodeName)
	}

	nodeList, err := client.CoreV1().Nodes().List(ctx, listOptions)
	if err != nil {
		klog.Error(ctx, err.Error())
		return
	}

	if len(nodeList.Items) == 0 {
		fmt.Println("没有找到匹配的节点")
		return
	}

	var nodeInfos []NodeInfo
	ItemList := make([]map[string]string, 0)

	for _, node := range nodeList.Items {
		nodeInfo := processNodeInfo(ctx, &node, metricsClient)
		
		// 应用节点组过滤器
		if options.NodeGroupName != "" && !strings.Contains(nodeInfo.NodeGroupName, options.NodeGroupName) {
			continue
		}
		
		nodeInfos = append(nodeInfos, nodeInfo)
		
		// 为了兼容现有的表格输出，也构建 map
		deployMap := nodeInfoToMap(nodeInfo)
		ItemList = append(ItemList, deployMap)
	}

	if len(nodeInfos) == 0 {
		fmt.Printf("没有找到匹配过滤条件的节点")
		if options.NodeGroupName != "" {
			fmt.Printf("（节点组: %s）", options.NodeGroupName)
		}
		fmt.Println()
		return
	}

	// 应用排序
	if options.SortBy != "" {
		sortNodeInfos(nodeInfos, options.SortBy)
		// 重新构建 ItemList 以保持排序
		ItemList = make([]map[string]string, 0)
		for _, nodeInfo := range nodeInfos {
			ItemList = append(ItemList, nodeInfoToMap(nodeInfo))
		}
	}

	// 根据输出格式显示结果
	switch options.OutputFormat {
	case "json":
		outputJSON(nodeInfos)
	case "yaml":
		outputYAML(nodeInfos)
	default:
		mtable.TablePrint("node", ItemList)
		
		// 如果有 metrics 问题，显示提示信息
		hasMetricsIssue := false
		for _, nodeInfo := range nodeInfos {
			if nodeInfo.UsedCPU == "N/A" {
				hasMetricsIssue = true
				break
			}
		}
		
		if hasMetricsIssue {
			fmt.Println("\n注意: 部分资源使用数据显示为 N/A，这通常是因为:")
			fmt.Println("1. 集群中没有安装 metrics-server")
			fmt.Println("2. metrics-server 没有正常运行")
			fmt.Println("3. 节点上的 kubelet 没有启用 metrics 功能")
			fmt.Println("\n要安装 metrics-server，请运行:")
			fmt.Println("kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml")
		}
	}

	// 导出 Excel
	if export {
		excel.ExportXlsx(ctx, "node", ItemList, kubeconfig)
	}
}

// processNodeInfo 处理单个节点信息
func processNodeInfo(ctx context.Context, node *corev1.Node, metricsClient *versioned.Clientset) NodeInfo {
	nodeInfo := NodeInfo{
		NodeName:                node.Name,
		OSImage:                 node.Status.NodeInfo.OSImage,
		KubeletVersion:          node.Status.NodeInfo.KubeletVersion,
		ContainerRuntimeVersion: node.Status.NodeInfo.ContainerRuntimeVersion,
	}

	// 获取节点 IP
	if len(node.Status.Addresses) > 0 {
		nodeInfo.NodeIP = node.Status.Addresses[0].Address
	}

	// 获取节点组名称
	nodeInfo.NodeGroupName = getNodeGroupName(node)

	// 获取节点容量信息（这些信息总是可用的）
	totalCpuCores := node.Status.Allocatable.Cpu().MilliValue()
	totalMemoryMi := float64(node.Status.Allocatable.Memory().Value()) / 1024 / 1024

	nodeInfo.TotalCPU = fmt.Sprintf("%dm", totalCpuCores)
	nodeInfo.TotalMemory = fmt.Sprintf("%.2fMi", totalMemoryMi)

	// 尝试获取资源使用情况
	nodeMetrics, err := metricsClient.MetricsV1beta1().NodeMetricses().Get(ctx, node.Name, metav1.GetOptions{})
	if err != nil {
		// 如果无法获取 metrics，设置默认值并记录警告
		klog.Warningf("无法获取节点 %s 的 metrics 数据: %v", node.Name, err)
		nodeInfo.UsedCPU = "N/A"
		nodeInfo.CPUUsagePercent = "N/A"
		nodeInfo.UsedMemory = "N/A"
		nodeInfo.MemoryUsagePercent = "N/A"
		nodeInfo.CPUUsageValue = 0
		nodeInfo.MemoryUsageValue = 0
		return nodeInfo
	}

	// 计算 CPU 和内存使用情况
	cpuUsage := nodeMetrics.Usage[corev1.ResourceCPU]
	memoryUsage := nodeMetrics.Usage[corev1.ResourceMemory]

	// 内存转换
	usedMemoryMi := float64(memoryUsage.Value()) / 1024 / 1024

	// CPU 转换
	usedCpuCores := float64(cpuUsage.MilliValue())

	// 计算使用百分比
	cpuPercent := (usedCpuCores / float64(totalCpuCores)) * 100
	memoryPercent := (usedMemoryMi / totalMemoryMi) * 100

	nodeInfo.UsedCPU = fmt.Sprintf("%.2fm", usedCpuCores)
	nodeInfo.CPUUsagePercent = fmt.Sprintf("%.2f%%", cpuPercent)
	nodeInfo.UsedMemory = fmt.Sprintf("%.2fMi", usedMemoryMi)
	nodeInfo.MemoryUsagePercent = fmt.Sprintf("%.2f%%", memoryPercent)

	// 保存数值用于排序
	nodeInfo.CPUUsageValue = cpuPercent
	nodeInfo.MemoryUsageValue = memoryPercent

	return nodeInfo
}

// getNodeGroupName 获取节点组名称
func getNodeGroupName(node *corev1.Node) string {
	// 检查 EKS 节点组标签
	if nodeGroup, exists := node.Labels["eks.amazonaws.com/nodegroup"]; exists {
		return nodeGroup
	}
	
	// 检查其他可能的节点组标签
	possibleLabels := []string{
		"kops.k8s.io/instancegroup",
		"node.kubernetes.io/instance-type",
		"beta.kubernetes.io/instance-type",
		"kubernetes.io/hostname",
	}
	
	for _, label := range possibleLabels {
		if value, exists := node.Labels[label]; exists {
			return value
		}
	}
	
	return "未知"
}

// nodeInfoToMap 将 NodeInfo 转换为 map（兼容现有代码）
func nodeInfoToMap(nodeInfo NodeInfo) map[string]string {
	return map[string]string{
		"节点名":                      nodeInfo.NodeName,
		"节点组名称":                    nodeInfo.NodeGroupName,
		"节点IP":                     nodeInfo.NodeIP,
		"OS镜像":                     nodeInfo.OSImage,
		"Kubelet版本":                nodeInfo.KubeletVersion,
		"CONTAINER_RUNTIME_VERSION": nodeInfo.ContainerRuntimeVersion,
		"当前已使用的CPU":                nodeInfo.UsedCPU,
		"CPU总大小":                   nodeInfo.TotalCPU,
		"CPU使用占服务器的百分比":            nodeInfo.CPUUsagePercent,
		"当前已使用的内存":                nodeInfo.UsedMemory,
		"内存总大小":                    nodeInfo.TotalMemory,
		"内存使用占服务器的百分比":            nodeInfo.MemoryUsagePercent,
	}
}

// sortNodeInfos 对节点信息进行排序
func sortNodeInfos(nodeInfos []NodeInfo, sortBy string) {
	switch sortBy {
	case "cpu":
		sort.Slice(nodeInfos, func(i, j int) bool {
			return nodeInfos[i].CPUUsageValue > nodeInfos[j].CPUUsageValue
		})
	case "memory":
		sort.Slice(nodeInfos, func(i, j int) bool {
			return nodeInfos[i].MemoryUsageValue > nodeInfos[j].MemoryUsageValue
		})
	case "name":
		sort.Slice(nodeInfos, func(i, j int) bool {
			return nodeInfos[i].NodeName < nodeInfos[j].NodeName
		})
	case "nodegroup":
		sort.Slice(nodeInfos, func(i, j int) bool {
			return nodeInfos[i].NodeGroupName < nodeInfos[j].NodeGroupName
		})
	}
}

// outputJSON 输出 JSON 格式
func outputJSON(nodeInfos []NodeInfo) {
	jsonData, err := json.MarshalIndent(nodeInfos, "", "  ")
	if err != nil {
		klog.Errorf("Failed to marshal JSON: %v", err)
		return
	}
	fmt.Println(string(jsonData))
}

// outputYAML 输出 YAML 格式
func outputYAML(nodeInfos []NodeInfo) {
	yamlData, err := yaml.Marshal(nodeInfos)
	if err != nil {
		klog.Errorf("Failed to marshal YAML: %v", err)
		return
	}
	fmt.Println(string(yamlData))
}

// GetNodesByNodeGroup 根据节点组名称获取节点列表
func GetNodesByNodeGroup(ctx context.Context, kubeconfig, nodeGroupName string) ([]NodeInfo, error) {
	client, err := NewClientset(kubeconfig)
	if err != nil {
		return nil, err
	}

	metricsClient, err := NewMetricsClient(kubeconfig)
	if err != nil {
		return nil, err
	}

	// 使用标签选择器过滤节点组
	labelSelector := fmt.Sprintf("eks.amazonaws.com/nodegroup=%s", nodeGroupName)
	nodeList, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, err
	}

	var nodeInfos []NodeInfo
	for _, node := range nodeList.Items {
		nodeInfo := processNodeInfo(ctx, &node, metricsClient)
		nodeInfos = append(nodeInfos, nodeInfo)
	}

	return nodeInfos, nil
}

// GetNodeGroupStats 获取节点组统计信息
func GetNodeGroupStats(ctx context.Context, kubeconfig string) (map[string]NodeGroupStat, error) {
	client, err := NewClientset(kubeconfig)
	if err != nil {
		return nil, err
	}

	nodeList, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	stats := make(map[string]NodeGroupStat)
	
	for _, node := range nodeList.Items {
		nodeGroupName := getNodeGroupName(&node)
		
		stat := stats[nodeGroupName]
		stat.NodeGroupName = nodeGroupName
		stat.NodeCount++
		
		// 累加资源
		cpuCapacity := node.Status.Allocatable.Cpu().MilliValue()
		memoryCapacity := node.Status.Allocatable.Memory().Value() / 1024 / 1024 // 转换为 Mi
		
		stat.TotalCPU += float64(cpuCapacity)
		stat.TotalMemory += float64(memoryCapacity)
		
		stats[nodeGroupName] = stat
	}

	return stats, nil
}
