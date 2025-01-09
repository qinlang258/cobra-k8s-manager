package kube

import (
	"context"
	"io/ioutil"
	"regexp"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog/v2"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// 定义 kubeconfig 文件结构体
type KubeConfig struct {
	APIVersion     string        `yaml:"apiVersion"`
	Kind           string        `yaml:"kind"`
	Contexts       []api.Context `yaml:"contexts"`
	CurrentContext string        `yaml:"currentContext"`
}

func NewClientset(configPath string) (*kubernetes.Clientset, error) {

	if configPath != "" {
		config, err := clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			klog.Error(err.Error())
		}
		return kubernetes.NewForConfig(config)
	} else {
		config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			klog.Error(err.Error())
		}
		return kubernetes.NewForConfig(config)
	}
}

func GetClientgoNamespace(configPath string) string {
	if configPath != "" {
		bytes, err := ioutil.ReadFile(configPath)
		if err != nil {
			klog.Error("没有找到该配置文件", configPath)
		}

		// 正则表达式查找 namespace 字段
		re := regexp.MustCompile(`namespace:\s*(\S+)`)
		matches := re.FindStringSubmatch(string(bytes))

		// 如果找到了匹配项
		if len(matches) > 1 {
			return matches[1]
		} else {
			return "default"
		}
	} else {
		bytes, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
		if err != nil {
			klog.Error("没有找到该配置文件", clientcmd.RecommendedHomeFile)
		}

		// 正则表达式查找 namespace 字段
		re := regexp.MustCompile(`namespace:\s*(\S+)`)
		matches := re.FindStringSubmatch(string(bytes))

		// 如果找到了匹配项
		if len(matches) > 1 {
			return matches[1]
		} else {
			return "default"
		}
	}

}

func NewMetricsClient(configPath string) (*versioned.Clientset, error) {

	if configPath != "" {
		config, err := clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			klog.Error(context.Background(), err)
		}
		metricsClient, err := versioned.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		return metricsClient, err
	} else {
		config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			klog.Error(context.Background(), err)
		}
		metricsClient, err := versioned.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		return metricsClient, err
	}

}
