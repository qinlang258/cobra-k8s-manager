package main

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, _ := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)

	client, _ := kubernetes.NewForConfig(config)

	data, _ := client.CoreV1().Pods("freshx").Get(context.Background(), "openapi-biz-747648fd95-bnd7p", metav1.GetOptions{})

	envs := data.Spec.Containers[0].Env

	for _, values := range envs {
		if values.Name == "JAVA_OPTS" {
			labels := strings.Fields(values.Value)
			for _, v := range labels {
				fmt.Println(v)
			}
		}
	}

}
