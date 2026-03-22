package cluster

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

type KClient struct {
	Client         *kubernetes.Clientset
	Metrics        *metricsv.Clientset
	HasMetrics     bool
	ClusterName    string
	ClusterVersion string
}

func NewClient() (*KClient, error) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("kubeconfig: %w", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("k8s client: %w", err)
	}

	// Extract cluster name from kubeconfig context
	rawConfig, _ := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	clusterName := "unknown"
	if rawConfig != nil && rawConfig.CurrentContext != "" {
		clusterName = rawConfig.CurrentContext
	}

	// Get server version
	sv, err := client.Discovery().ServerVersion()
	clusterVersion := "unknown"
	if err == nil {
		clusterVersion = sv.GitVersion
	}

	// Try metrics client
	mc, merr := metricsv.NewForConfig(config)
	hasMetrics := false
	if merr == nil {
		_, err := mc.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
		hasMetrics = err == nil
	}

	return &KClient{
		Client:         client,
		Metrics:        mc,
		HasMetrics:     hasMetrics,
		ClusterName:    clusterName,
		ClusterVersion: clusterVersion,
	}, nil
}

func Show(k *KClient) {
	nodeList, _ := k.Client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	nodeCount := 0
	region := ""
	if nodeList != nil {
		nodeCount = len(nodeList.Items)
		if nodeCount > 0 {
			for key, val := range nodeList.Items[0].Labels {
				if strings.Contains(key, "topology.kubernetes.io/region") {
					region = val
					break
				}
			}
		}
	}

	metricsStatus := "✅ Available"
	if !k.HasMetrics {
		metricsStatus = "⚠️  Not installed (showing requests/limits only)"
	}

	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Printf("║  🔵 KPULSE — Cluster: %-37s ║\n", k.ClusterName)
	fmt.Printf("║     Version: %-10s  Region: %-8s  Nodes: %-5d    ║\n", k.ClusterVersion, region, nodeCount)
	fmt.Printf("║     Metrics Server: %-39s ║\n", metricsStatus)
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}
