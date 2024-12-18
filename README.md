# cobra-k8s-manager
cobra-k8s-manager

该命令有以下几个功能：analysis，image，node，resource，top

## 通用配置  
+ 所有命令均可附带 --kubeconfig指定配置文件
./k8s-manager --kubeconfig <指定使用的k8s配置文件>


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
go run main.go resource prometheus -u <prometheus访问地址>
```

## 5 top 获取指定namespace的资源使用情况
示例代码
```powershell
1 获取所有namespace的资源开销
./k8s-manager top
2 获取指定namespace的资源开销
./k8s-manager top -n <namespace> 
```


