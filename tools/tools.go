package tools

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetNodeGroupNameFromPod(ctx context.Context, client *kubernetes.Clientset, pod v1.Pod) (string, error) {
	// 获取Pod所在节点的信息
	node, err := client.CoreV1().Nodes().Get(ctx, pod.Spec.NodeName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get node: %v", err)
	}

	// 检查节点标签
	if nodeGroupName, ok := node.Labels["eks.amazonaws.com/nodegroup"]; ok {
		return nodeGroupName, nil
	}

	// 对于自管理节点组，可能使用不同的标签
	if nodeGroupName, ok := node.Labels["eks.amazonaws.com/nodegroup"]; ok {
		return nodeGroupName, nil
	}

	return "", fmt.Errorf("node group name not found in node labels")
}
