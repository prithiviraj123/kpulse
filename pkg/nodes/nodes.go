package nodes

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/prithiviraj123/kpulse/pkg/cluster"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeInfo struct {
	Name       string
	Status     string
	CPUReq     int64
	CPUCap     int64
	CPUUsage   int64
	MemReq     int64
	MemCap     int64
	MemUsage   int64
	Pods       int
	PodCap     int64
	Instance   string
}

func Gather(k *cluster.KClient) []NodeInfo {
	nodeList, err := k.Client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil || nodeList == nil {
		return nil
	}

	podList, _ := k.Client.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})

	// Aggregate pod requests per node
	nodeReqs := map[string][2]int64{} // [cpu_milli, mem_bytes]
	nodePodCount := map[string]int{}
	if podList != nil {
		for _, pod := range podList.Items {
			if pod.Status.Phase != corev1.PodRunning || pod.Spec.NodeName == "" {
				continue
			}
			nodePodCount[pod.Spec.NodeName]++
			var cpuReq, memReq int64
			for _, c := range pod.Spec.Containers {
				if r, ok := c.Resources.Requests[corev1.ResourceCPU]; ok {
					cpuReq += r.MilliValue()
				}
				if r, ok := c.Resources.Requests[corev1.ResourceMemory]; ok {
					memReq += r.Value()
				}
			}
			cur := nodeReqs[pod.Spec.NodeName]
			cur[0] += cpuReq
			cur[1] += memReq
			nodeReqs[pod.Spec.NodeName] = cur
		}
	}

	// Get actual usage from metrics if available
	nodeUsage := map[string][2]int64{}
	if k.HasMetrics {
		ml, err := k.Metrics.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
		if err == nil && ml != nil {
			for _, m := range ml.Items {
				cpu := m.Usage[corev1.ResourceCPU]
				mem := m.Usage[corev1.ResourceMemory]
				nodeUsage[m.Name] = [2]int64{cpu.MilliValue(), mem.Value()}
			}
		}
	}

	var result []NodeInfo
	for _, node := range nodeList.Items {
		status := "NotReady"
		for _, c := range node.Status.Conditions {
			if c.Type == corev1.NodeReady && c.Status == corev1.ConditionTrue {
				status = "Ready"
			}
		}

		cpuCap := node.Status.Allocatable[corev1.ResourceCPU]
		memCap := node.Status.Allocatable[corev1.ResourceMemory]
		podCap := node.Status.Allocatable[corev1.ResourcePods]

		reqs := nodeReqs[node.Name]
		usage := nodeUsage[node.Name]

		instanceType := node.Labels["node.kubernetes.io/instance-type"]
		if instanceType == "" {
			instanceType = node.Labels["beta.kubernetes.io/instance-type"]
		}

		result = append(result, NodeInfo{
			Name:     shortenName(node.Name),
			Status:   status,
			CPUReq:   reqs[0],
			CPUCap:   cpuCap.MilliValue(),
			CPUUsage: usage[0],
			MemReq:   reqs[1],
			MemCap:   memCap.Value(),
			MemUsage: usage[1],
			Pods:     nodePodCount[node.Name],
			PodCap:   podCap.Value(),
			Instance: instanceType,
		})
	}

	sort.Slice(result, func(i, j int) bool { return result[i].CPUReq > result[j].CPUReq })
	return result
}

func Show(k *cluster.KClient) {
	nodeInfos := Gather(k)
	if len(nodeInfos) == 0 {
		fmt.Println("\n⚠️  No nodes found")
		return
	}

	fmt.Println("\n📦 NODES")
	for i, n := range nodeInfos {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("  Node %d\n", i+1)
		fmt.Println("  ┌──────────────────┬──────────────────────────┐")
		fmt.Printf("  │ Name             │ %-24s │\n", n.Name)
		fmt.Printf("  │ Instance         │ %-24s │\n", n.Instance)
		fmt.Printf("  │ CPU (req/cap)    │ %-24s │\n", fmt.Sprintf("%dm / %dm", n.CPUReq, n.CPUCap))
		fmt.Printf("  │ Memory (req/cap) │ %-24s │\n", fmt.Sprintf("%s / %s", fmtMem(n.MemReq), fmtMem(n.MemCap)))
		if k.HasMetrics {
			fmt.Printf("  │ CPU Usage        │ %-24s │\n", fmt.Sprintf("%dm (%d%%)", n.CPUUsage, pct(n.CPUUsage, n.CPUCap)))
			fmt.Printf("  │ Memory Usage     │ %-24s │\n", fmt.Sprintf("%s (%d%%)", fmtMem(n.MemUsage), pct(n.MemUsage, n.MemCap)))
		}
		fmt.Printf("  │ Pods             │ %-24s │\n", fmt.Sprintf("%d / %d", n.Pods, n.PodCap))
		fmt.Println("  └──────────────────┴──────────────────────────┘")
	}
}

func shortenName(name string) string {
	if strings.Contains(name, ".ec2.internal") {
		name = strings.TrimSuffix(name, ".ec2.internal")
	}
	if strings.HasPrefix(name, "ip-") {
		name = strings.ReplaceAll(name[3:], "-", ".")
	}
	return name
}

func fmtMem(bytes int64) string {
	gi := float64(bytes) / (1024 * 1024 * 1024)
	if gi >= 1 {
		return fmt.Sprintf("%.1fGi", gi)
	}
	mi := float64(bytes) / (1024 * 1024)
	return fmt.Sprintf("%.0fMi", mi)
}

func pct(used, total int64) int64 {
	if total == 0 {
		return 0
	}
	return (used * 100) / total
}

// Ensure resource import is used
var _ = resource.DecimalSI
