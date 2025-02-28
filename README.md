# cobra-k8s-manager
cobra-k8s-manager

该命令有以下几个功能：analysis，image，node，resource，top, config

## 通用配置  
+ 所有命令均可附带 --kubeconfig指定配置文件
./k8s-manager --kubeconfig <指定使用的k8s配置文件>

+ 监听因切换集群导致的/root/.kube/config变化。  prometheus.yaml的配置路径是 当前家目录下的.jcrose-prometheus/prometheus.yaml
执行了该命令之后切换集群导致的config变化，使家目录下的~/jcrose-prometheus/prometheus.yaml 的/root/.kube/config字段对应的prometheus地址产生同样的变化
nohup pkg/fsnotify/fsnotify &  


## 1 node 分析所有node的资源情况
示例代码
```powershell
1 获取所有节点的资源信息
./k8s-manager node 
```

## 2 analysis 分析Node节点上的资源使用构成

示例代码
```powershell
1 分析指定节点上的所有容器的资源开销
./k8s-manager analysis --node <节点名>  
```

## 3 image 获取指定namespace的所有镜像地址

示例代码
```powershell
1 获取所有namespace的镜像地址  
./k8s-manager image  
2 获取指定namespace的镜像地址
./k8s-manager image -n <namespace>
```

## 4 resource 获取指定namespace的所有limit 与 Requests大小
示例代码
```powershell
1 获取所有namespace的limit 与 Requests大小  
./k8s-manager resource  
2 获取指定namespace的limit 与 Requests大小
./k8s-manager resource -n <namespace>

3 在prometheus查询最近七天的内存CPU使用情况
go run main.go resource prometheus  -p
```

## 5 top 获取指定namespace的资源使用情况
示例代码
```powershell
1 获取所有namespace的资源开销
./k8s-manager top
2 获取指定namespace的资源开销
./k8s-manager top -n <namespace> 
```

## 6 config 自动获取K8S配置文件路径下所有集群的prometheus的ingress地址
示例代码
```powershell
1 默认获取 /root/.kube/文件下所有yaml文件的prometheus地址
go run main.go config
2 获取指定文件下所有yaml文件的prometheus地址
go run main.go config -p /data/k8s
```
