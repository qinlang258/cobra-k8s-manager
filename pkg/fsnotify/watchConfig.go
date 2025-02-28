package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

type Cluster struct {
	Name string `yaml:"name"`
}

type KubeConfig struct {
	Clusters []Cluster `yaml:"clusters"`
}

type PrometheusConfigs struct {
	Prometheus []PrometheusConfig
}

type PrometheusConfig struct {
	Kubeconfig string `yaml:"kubeconfig"`
	URL        string `yaml:"url"`
	Port       int    `yaml:"port"`
}

func watchKubeConfig(filePath string, watcher *fsnotify.Watcher) error {
	err := watcher.Add(filePath)
	if err != nil {
		return err
	}
	return nil
}

func GetClusterNameFromKubeConfig(filePath string) (string, error) {
	// 读取 kubeconfig 文件
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	var kubeConfig KubeConfig
	err = yaml.Unmarshal(data, &kubeConfig)
	if err != nil {
		return "", err
	}

	// 获取第一个 cluster 的 name
	if len(kubeConfig.Clusters) > 0 {
		return kubeConfig.Clusters[0].Name, nil
	}

	return "", fmt.Errorf("no clusters found in the kubeconfig")
}

func updatePrometheusConfig(prometheusConfigPath, clusterName string) error {
	// 读取 prometheus.yaml 配置文件
	data, err := ioutil.ReadFile(prometheusConfigPath)
	if err != nil {
		return err
	}

	var targetUrl string
	var targetPort int
	var prometheusConfigs PrometheusConfigs
	err = yaml.Unmarshal(data, &prometheusConfigs)
	if err != nil {
		return err
	}

	// 查找匹配的 kubeconfig
	for i := range prometheusConfigs.Prometheus {
		if strings.Contains(prometheusConfigs.Prometheus[i].Kubeconfig, clusterName) {
			// 获取目标URL
			targetUrl = prometheusConfigs.Prometheus[i].URL
			targetPort = prometheusConfigs.Prometheus[i].Port
			break
		}
	}

	// 查找匹配的 kubeconfig
	for j := range prometheusConfigs.Prometheus {
		if strings.Contains(prometheusConfigs.Prometheus[j].Kubeconfig, "/root/.kube/config") {
			// 获取目标URL
			prometheusConfigs.Prometheus[j].URL = targetUrl
			prometheusConfigs.Prometheus[j].Port = targetPort
			break
		}
	}

	// 将更新后的内容写回文件
	updatedData, err := yaml.Marshal(&prometheusConfigs)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(prometheusConfigPath, updatedData, 0644)
	if err != nil {
		return err
	}

	fmt.Println("prometheus.yaml updated with new URL:", prometheusConfigs.Prometheus[0].URL)
	return nil
}

func main() {
	// 监听 kubeconfig 和 prometheus.yaml 的路径
	// 获取当前用户的家目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		klog.Error(context.Background(), "Failed to get user home directory: "+err.Error())
	}
	kubeConfigPath := filepath.Join(homeDir, ".kube", "config")
	prometheusConfigPath := filepath.Join(homeDir, ".kube", "jcrose-prometheus", "prometheus.yaml")

	// 创建 fsnotify watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error creating watcher:", err)
		return
	}
	defer watcher.Close()

	// 开始监听 kubeconfig 文件的变动
	err = watchKubeConfig(kubeConfigPath, watcher)
	if err != nil {
		fmt.Println("Error watching kubeconfig:", err)
		return
	}

	fmt.Println("Watching for changes in", kubeConfigPath)

	// 获取初始的 cluster 名称并更新 prometheus.yaml
	clusterName, err := GetClusterNameFromKubeConfig(kubeConfigPath)
	if err != nil {
		fmt.Println("Error getting cluster name:", err)
		return
	}
	err = updatePrometheusConfig(prometheusConfigPath, clusterName)
	if err != nil {
		fmt.Println("Error updating prometheus config:", err)
		return
	}

	// 监听文件变化
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("Detected change in kubeconfig:", event.Name)
				clusterName, err := GetClusterNameFromKubeConfig(kubeConfigPath)
				if err != nil {
					fmt.Println("Error getting cluster name:", err)
					continue
				}
				err = updatePrometheusConfig(prometheusConfigPath, clusterName)
				if err != nil {
					fmt.Println("Error updating prometheus config:", err)
				}
			}
		case err := <-watcher.Errors:
			fmt.Println("Watcher error:", err)
		}
	}
}
