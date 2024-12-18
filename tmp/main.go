package main

import (
	"context"
	"k8s-manager/pkg/kube"
)

func main() {
	ctx := context.Background()
	kube.TestPrometheus(ctx, "prometheus-k8s-0", "prometheus", "monitoring", "http://192.168.44.134:20248/")
}
