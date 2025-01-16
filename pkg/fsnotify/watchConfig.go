package main

import (
	"fmt"
	"io/ioutil"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type Cluster struct {
	Name string `yaml:"name"`
}

type KubeConfig struct {
	Clusters []Cluster `yaml:"clusters"`
}

type PrometheusConfig struct {
	Prometheus []struct {
		Kubeconfig string `yaml:"kubeconfig"`
		URL        string `yaml:"url"`
	} `yaml:"prometheus"`
}

func watchKubeConfig(filePath string, watcher *fsnotify.Watcher) error {
	err := watcher.Add(filePath)
	if err != nil {
		return err
	}
	return nil
}

func getClusterNameFromKubeConfig(filePath string) (string, error) {
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

	var prometheusConfig PrometheusConfig
	err = yaml.Unmarshal(data, &prometheusConfig)
	if err != nil {
		return err
	}

	// 更新 prometheus.yaml 中的 URL
	for i := range prometheusConfig.Prometheus {
		if prometheusConfig.Prometheus[i].Kubeconfig == "/root/.kube/config" {
			prometheusConfig.Prometheus[i].URL = fmt.Sprintf("http://prometheus.%s.yunlizhi.net", clusterName)
		}
	}

	// 将更新后的内容写回文件
	updatedData, err := yaml.Marshal(&prometheusConfig)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(prometheusConfigPath, updatedData, 0644)
	if err != nil {
		return err
	}

	fmt.Println("prometheus.yaml updated with new URL:", prometheusConfig.Prometheus[0].URL)
	return nil
}

func main() {
	// 监听 kubeconfig 和 prometheus.yaml 的路径
	kubeConfigPath := "/root/.kube/config"
	prometheusConfigPath := "../../config/prometheus.yaml"

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
	clusterName, err := getClusterNameFromKubeConfig(kubeConfigPath)
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
				clusterName, err := getClusterNameFromKubeConfig(kubeConfigPath)
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
