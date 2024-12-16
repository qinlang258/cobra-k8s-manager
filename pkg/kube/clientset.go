package kube

import (
	"context"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

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
