package pods

import (
	"context"
	"fmt"
	"sort"

	"github.com/prithiviraj123/kpulse/pkg/cluster"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodInfo struct {
	Name      string
	Namespace string
	CPUReq    int64
	CPUUsage  int64
	MemReq    int64
	MemUsage  int64
	Status    string
}

func Show(k *cluster.KClient) {
	podList, err := k.Client.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil || podList == nil {
		fmt.Println("\n⚠️  Could not list pods")
		return
	}

	// Get pod metrics if available
	podUsage := map[string][2]int64{} // key: ns/name
	if k.HasMetrics {
		ml, err := k.Metrics.MetricsV1beta1().PodMetricses(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
		if err == nil && ml != nil {
			for _, m := range ml.Items {
				var cpu, mem int64
				for _, c := range m.Containers {
					cpu += c.Usage.Cpu().MilliValue()
					mem += c.Usage.Memory().Value()
				}
				podUsage[m.Namespace+"/"+m.Name] = [2]int64{cpu, mem}
			}
		}
	}

	var infos []PodInfo
	for _, pod := range podList.Items {
		if pod.Status.Phase != corev1.PodRunning {
			continue
		}
		var cpuReq, memReq int64
		for _, c := range pod.Spec.Containers {
			if r, ok := c.Resources.Requests[corev1.ResourceCPU]; ok {
				cpuReq += r.MilliValue()
			}
			if r, ok := c.Resources.Requests[corev1.ResourceMemory]; ok {
				memReq += r.Value()
			}
		}
		usage := podUsage[pod.Namespace+"/"+pod.Name]
		infos = append(infos, PodInfo{
			Name:      shortenPod(pod.Name),
			Namespace: pod.Namespace,
			CPUReq:    cpuReq,
			CPUUsage:  usage[0],
			MemReq:    memReq,
			MemUsage:  usage[1],
			Status:    string(pod.Status.Phase),
		})
	}

	// Sort by CPU (usage if available, else requests)
	sort.Slice(infos, func(i, j int) bool {
		if infos[i].CPUUsage > 0 || infos[j].CPUUsage > 0 {
			return infos[i].CPUUsage > infos[j].CPUUsage
		}
		return infos[i].CPUReq > infos[j].CPUReq
	})

	// Show top 10
	limit := 10
	if len(infos) < limit {
		limit = len(infos)
	}

	fmt.Printf("\n🔥 TOP PODS (by %s) — showing %d of %d running\n", metricLabel(k), limit, len(infos))
	fmt.Println("┌──────────────────────────┬──────────────┬──────────┬──────────┐")
	fmt.Println("│ Pod                      │ Namespace    │ CPU      │ Memory   │")
	fmt.Println("├──────────────────────────┼──────────────┼──────────┼──────────┤")
	for _, p := range infos[:limit] {
		cpu := fmt.Sprintf("%dm req", p.CPUReq)
		mem := fmtMem(p.MemReq)
		if k.HasMetrics && p.CPUUsage > 0 {
			cpu = fmt.Sprintf("%dm", p.CPUUsage)
			mem = fmtMem(p.MemUsage)
		}
		fmt.Printf("│ %-24s │ %-12s │ %-8s │ %-8s │\n", p.Name, p.Namespace, cpu, mem)
	}
	fmt.Println("└──────────────────────────┴──────────────┴──────────┴──────────┘")
}

func metricLabel(k *cluster.KClient) string {
	if k.HasMetrics {
		return "usage"
	}
	return "requests"
}

func shortenPod(name string) string {
	if len(name) > 24 {
		return name[:24]
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
