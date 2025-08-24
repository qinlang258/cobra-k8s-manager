#!/bin/bash

# K8s Manager 节点过滤功能使用示例

echo "=== K8s Manager 节点过滤功能使用示例 ==="
echo

# 设置可执行文件路径
K8S_MANAGER="./k8s-manager"

echo "1. 获取所有节点信息（原始功能）"
echo "命令: $K8S_MANAGER node"
echo

echo "2. 使用新的过滤功能获取所有节点信息"
echo "命令: $K8S_MANAGER node-filter"
echo

echo "3. 按节点组名称过滤节点"
echo "命令: $K8S_MANAGER node-filter --nodegroup my-nodegroup"
echo

echo "4. 按节点名称过滤"
echo "命令: $K8S_MANAGER node-filter --node-name worker-node-1"
echo

echo "5. 按CPU使用率排序（降序）"
echo "命令: $K8S_MANAGER node-filter --sort-by cpu"
echo

echo "6. 按内存使用率排序（降序）"
echo "命令: $K8S_MANAGER node-filter --sort-by memory"
echo

echo "7. 按节点名称排序（升序）"
echo "命令: $K8S_MANAGER node-filter --sort-by name"
echo

echo "8. 按节点组名称排序（升序）"
echo "命令: $K8S_MANAGER node-filter --sort-by nodegroup"
echo

echo "9. 输出为JSON格式"
echo "命令: $K8S_MANAGER node-filter --output json"
echo

echo "10. 输出为YAML格式"
echo "命令: $K8S_MANAGER node-filter --output yaml"
echo

echo "11. 组合使用：过滤特定节点组并按CPU排序，输出JSON"
echo "命令: $K8S_MANAGER node-filter --nodegroup my-nodegroup --sort-by cpu --output json"
echo

echo "12. 显示节点组统计信息"
echo "命令: $K8S_MANAGER nodegroup stats"
echo

echo "13. 以JSON格式显示节点组统计"
echo "命令: $K8S_MANAGER nodegroup stats --output json"
echo

echo "14. 列出特定节点组中的节点"
echo "命令: $K8S_MANAGER nodegroup list my-nodegroup"
echo

echo "15. 导出Excel文件"
echo "命令: $K8S_MANAGER node-filter --nodegroup my-nodegroup --export"
echo

echo "=== 高级用法示例 ==="
echo

echo "16. 查找CPU使用率最高的节点"
echo "命令: $K8S_MANAGER node-filter --sort-by cpu --output json | jq '.[0]'"
echo

echo "17. 查找特定节点组中内存使用率最高的节点"
echo "命令: $K8S_MANAGER node-filter --nodegroup my-nodegroup --sort-by memory --output json | jq '.[0]'"
echo

echo "18. 统计各节点组的节点数量"
echo "命令: $K8S_MANAGER nodegroup stats --output json | jq 'to_entries | map({nodegroup: .key, count: .value.node_count})'"
echo

echo "=== 实际执行示例（需要有效的kubeconfig）==="
echo

# 检查是否存在可执行文件
if [ ! -f "$K8S_MANAGER" ]; then
    echo "注意: 可执行文件 $K8S_MANAGER 不存在"
    echo "请先编译项目: go build -o k8s-manager main.go"
    echo
fi

# 检查是否有kubeconfig
if [ ! -f "/root/.kube/config" ] && [ -z "$KUBECONFIG" ]; then
    echo "注意: 未找到kubeconfig文件"
    echo "请确保 /root/.kube/config 存在或设置 KUBECONFIG 环境变量"
    echo
fi

echo "如果环境配置正确，可以运行以下命令进行测试："
echo

# 安全的测试命令（不会失败）
echo "# 测试帮助信息"
echo "$K8S_MANAGER node-filter --help"
echo

echo "# 测试节点组命令帮助"
echo "$K8S_MANAGER nodegroup --help"
echo

echo "脚本执行完成！"
