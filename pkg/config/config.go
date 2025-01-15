package config

import (
	"context"
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
}

func getPrometheusUrl(ctx context.Context, kubeconfigPath string) string {
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
	}

	return "http://" + data.Spec.Rules[0].Host
}

// 定义Cluster结构体
type Cluster struct {
	Cluster struct {
		CertificateAuthorityData string `yaml:"certificate-authority-data"`
		Server                   string `yaml:"server"`
	} `yaml:"cluster"`
	Name string `yaml:"name"`
}

type KubeConfig struct {
	Clusters []Cluster `yaml:"clusters"`
}

func findYamlFiles(root string) ([]string, error) {
	var yamlFiles []string

	// Walk the directory
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// If there's an error walking, return it
		if err != nil {
			return err
		}

		// Check if the file has a .yaml extension
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") || info.Name() == "config" {
			yamlFiles = append(yamlFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
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
		prometheus.Url = getPrometheusUrl(ctx, file)

		pcs = append(pcs, prometheus)

	}

	// 输出到文件 config/prometheus.yaml
	prometheusConfig := PrometheusConfig{Prometheus: pcs}
	outputFile := "config/prometheus.yaml"

	// 创建并打开文件，写入数据
	file, err := os.Create(outputFile)
	if err != nil {
		klog.Error(ctx, "Failed to create file: "+err.Error())
		return false
	}
	defer file.Close()

	// 将 pcs 数据转换为 YAML 格式
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
	klog.Info(ctx, "Prometheus configuration written to ", outputFile)
	return true
}
