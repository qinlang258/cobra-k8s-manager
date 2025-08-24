# cobra-k8s-manager
cobra-k8s-manager

该命令有以下几个功能：analysis，image，node，resource，top, config, nodegroup

## 通用配置  
+ 所有命令均可附带 --kubeconfig指定配置文件
./k8s-manager --kubeconfig <指定使用的k8s配置文件>

+ 监听因切换集群导致的/root/.kube/config变化。  prometheus.yaml的配置路径是 当前家目录下的.jcrose-prometheus/prometheus.yaml
执行了该命令之后切换集群导致的config变化，使家目录下的~/jcrose-prometheus/prometheus.yaml 的/root/.kube/config字段对应的prometheus地址产生同样的变化
nohup pkg/fsnotify/fsnotify &  

## 1 node 分析所有node的资源情况（已优化，支持过滤和排序）
示例代码
```powershell
# 获取所有节点的资源信息
./k8s-manager node 

# 按节点组名称过滤
./k8s-manager node --nodegroup my-nodegroup

# 按节点名称过滤  
./k8s-manager node --node-name worker-node-1

# 按CPU使用率排序（降序）
./k8s-manager node --sort-by cpu

# 按内存使用率排序（降序）
./k8s-manager node --sort-by memory

# 按节点名称排序（升序）
./k8s-manager node --sort-by name

# 按节点组名称排序（升序）
./k8s-manager node --sort-by nodegroup

# 输出JSON格式
./k8s-manager node --output json

# 输出YAML格式
./k8s-manager node --output yaml

# 组合使用：过滤特定节点组并按CPU排序，输出JSON
./k8s-manager node --nodegroup my-nodegroup --sort-by cpu --output json

# 导出Excel文件
./k8s-manager node --nodegroup my-nodegroup --export
```

## 2 nodegroup 节点组管理和统计（新增功能）
示例代码
```powershell
# 显示所有节点组统计信息
./k8s-manager nodegroup stats

# 以JSON格式显示节点组统计
./k8s-manager nodegroup stats --output json

# 以YAML格式显示节点组统计
./k8s-manager nodegroup stats --output yaml

# 列出特定节点组中的节点
./k8s-manager nodegroup list my-nodegroup

# 以JSON格式列出节点组中的节点
./k8s-manager nodegroup list my-nodegroup --output json
```

## 3 analysis 分析Node节点上的资源使用构成

示例代码
```powershell
# 分析指定节点上的所有容器的资源开销
./k8s-manager analysis --node <节点名>  
```

## 4 image 获取指定namespace的所有镜像地址

示例代码
```powershell
# 获取所有namespace的镜像地址  
./k8s-manager image  
# 获取指定namespace的镜像地址
./k8s-manager image -n <namespace>
```

## 5 resource 获取指定namespace的所有limit 与 Requests大小
示例代码
```powershell
# 获取所有namespace的limit 与 Requests大小  
./k8s-manager resource  
# 获取指定namespace的limit 与 Requests大小
./k8s-manager resource -n <namespace>

# 在prometheus查询最近七天的内存CPU使用情况
go run main.go resource prometheus  -p
```

## 6 top 获取指定namespace的资源使用情况
示例代码
```powershell
# 获取所有namespace的资源开销
./k8s-manager top
# 获取指定namespace的资源开销
./k8s-manager top -n <namespace> 
```

## 7 config 自动获取K8S配置文件路径下所有集群的prometheus的ingress地址
示例代码
```powershell
# 默认获取 /root/.kube/文件下所有yaml文件的prometheus地址
go run main.go config
# 获取指定文件下所有yaml文件的prometheus地址
go run main.go config -p /data/k8s
```

## 新增功能特性

### 节点过滤和排序
- **节点组过滤**: 支持按节点组名称过滤节点
- **节点名称过滤**: 支持按具体节点名称过滤
- **多种排序**: 支持按CPU使用率、内存使用率、节点名称、节点组名称排序
- **多种输出格式**: 支持table、JSON、YAML三种输出格式

### 节点组管理
- **统计信息**: 显示各节点组的节点数量、总CPU和内存资源
- **节点列表**: 列出指定节点组中的所有节点
- **多格式输出**: 支持表格、JSON、YAML格式输出

### 兼容性
- 保持原有命令的完全兼容性
- 新增功能通过可选参数提供
- 支持原有的Excel导出功能

### 使用场景
1. **资源监控**: 快速找到CPU或内存使用率最高的节点
2. **节点组管理**: 统计各节点组的资源分布情况
3. **故障排查**: 按节点组过滤，快速定位问题节点
4. **自动化集成**: JSON/YAML输出便于脚本处理和自动化工具集成
