package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

type PrometheusConfig struct {
	Prometheus []Prometheus `yaml:"prometheus"`
}

type Prometheus struct {
	KubeConfig string `yaml:"kubeconfig"`
	Url        string `yaml:"url"`
	Port       int    `yaml:"port"`
}

type Cluster struct {
	Name string `yaml:"name"`
}

type KubeConfig struct {
	Clusters []Cluster `yaml:"clusters"`
}

type PrometheusConfigs struct {
	Prometheus []Prometheus
}

func getPrometheusUrl(ctx context.Context, kubeconfigPath string) (string, error) {
	fmt.Println(kubeconfigPath)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		klog.Error(ctx, err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Error(ctx, err.Error())
	}

	data, err := clientset.NetworkingV1().Ingresses("monitoring").Get(ctx, "prometheus", metav1.GetOptions{})
	if err != nil {
		klog.Error(ctx, err.Error())
		return "", err
	}

	return "http://" + data.Spec.Rules[0].Host, nil
}

func findYamlFiles(root string) ([]string, error) {
	var yamlFiles []string

	dir := filepath.Dir(root)
	// 读取目录内容
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// 遍历目录中的文件和文件夹
	for _, file := range files {
		// 如果是文件并且后缀是 .yaml 或 .yml
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml") || file.Name() == "config") {
			// 拼接完整路径并添加到结果中
			yamlFiles = append(yamlFiles, dir+"/"+file.Name())
		}
	}

	return yamlFiles, nil
}

func InitPrometheus(ctx context.Context, kubeconfigPath string) bool {
	yamlFiles, err := findYamlFiles(kubeconfigPath)
	if err != nil {
		klog.Error(ctx, err.Error())
		return false
	}

	if yamlFiles == nil {
		klog.Error(ctx, "没有找到以.yaml为后缀或者 config相关的配置文件")
	}

	pcs := []Prometheus{}

	for _, file := range yamlFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			klog.Error(ctx, err.Error())
			continue
		}

		kubeConfig := KubeConfig{}
		err = yaml.Unmarshal(data, &kubeConfig)
		if err != nil {
			klog.Error(ctx, "Failed to unmarshal kubeconfig: "+err.Error())
			continue
		}

		var prometheus Prometheus
		prometheus.KubeConfig = file
		url, err := getPrometheusUrl(ctx, file)
		if err != nil {
			klog.Error(ctx, err.Error())
			continue
		}
		prometheus.Url = url

		pcs = append(pcs, prometheus)
	}

	// 获取当前用户的家目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		klog.Error(ctx, "Failed to get user home directory: "+err.Error())
		return false
	}

	// 定义目标文件路径
	targetDir := filepath.Join(homeDir, ".kube", "jcrose-prometheus")
	targetFile := filepath.Join(targetDir, "prometheus.yaml")

	// 确保目标目录存在
	err = os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		klog.Error(ctx, "Failed to create directories: "+err.Error())
		return false
	}

	// 创建并打开文件，写入数据
	file, err := os.Create(targetFile)
	if err != nil {
		klog.Error(ctx, "Failed to create file: "+err.Error())
		return false
	}
	defer file.Close()

	// 将 pcs 数据转换为 YAML 格式
	prometheusConfig := PrometheusConfig{Prometheus: pcs}
	dataToWrite, err := yaml.Marshal(prometheusConfig)
	if err != nil {
		klog.Error(ctx, "Failed to marshal prometheus config: "+err.Error())
		return false
	}

	// 写入 YAML 数据到文件
	_, err = file.Write(dataToWrite)
	if err != nil {
		klog.Error(ctx, "Failed to write to file: "+err.Error())
		return false
	}

	// 输出成功日志
	klog.Info(ctx, "Prometheus configuration written to ", targetFile)
	return true
}

func GetClusterNameFromPrometheusUrl(kubeconfigPath string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		klog.Error(context.Background(), "Failed to get user home directory: "+err.Error())
	}
	prometheusConfigPath := filepath.Join(homeDir, ".kube", "jcrose-prometheus", "prometheus.yaml")
	// 读取 prometheus.yaml 配置文件

	data, err := ioutil.ReadFile(prometheusConfigPath)
	if err != nil {
		return "", err
	}

	var name string
	var prometheusUrl string
	var prometheusConfigs PrometheusConfigs
	err = yaml.Unmarshal(data, &prometheusConfigs)
	if err != nil {
		return "", err
	}

	// 查找匹配的 kubeconfig
	for i := range prometheusConfigs.Prometheus {
		if strings.Contains(prometheusConfigs.Prometheus[i].KubeConfig, kubeconfigPath) {
			// 获取目标URL
			prometheusUrl = prometheusConfigs.Prometheus[i].Url
			break
		}
	}

	// 查找匹配的 kubeconfig
	for j := range prometheusConfigs.Prometheus {
		if strings.Contains(prometheusConfigs.Prometheus[j].Url, prometheusUrl) && prometheusConfigs.Prometheus[j].KubeConfig != "/root/.kube/config" {
			// 获取目标URL
			name = prometheusConfigs.Prometheus[j].KubeConfig
			break
		}
	}

	fieldsName := strings.Split(name, "/")
	yamlName := fieldsName[len(fieldsName)-1]
	return strings.Split(yamlName, ".")[0], err

}
