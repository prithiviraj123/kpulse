package alerts

import (
	"context"
	"fmt"

	"github.com/prithiviraj123/kpulse/pkg/cluster"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Show(k *cluster.KClient) {
	fmt.Println("\n⚠️  ALERTS")

	var warnings []string

	// Check nodes
	nodeList, _ := k.Client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if nodeList != nil {
		for _, node := range nodeList.Items {
			ready := false
			for _, c := range node.Status.Conditions {
				if c.Type == corev1.NodeReady && c.Status == corev1.ConditionTrue {
					ready = true
				}
			}
			if !ready {
				warnings = append(warnings, fmt.Sprintf("🔴 Node %s is NotReady", node.Name))
			}
		}
	}

	// Check pods
	podList, _ := k.Client.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	noLimits := 0
	noRequests := 0
	restarts := 0
	if podList != nil {
		for _, pod := range podList.Items {
			if pod.Status.Phase != corev1.PodRunning {
				continue
			}
			for _, c := range pod.Spec.Containers {
				if c.Resources.Limits.Cpu().IsZero() && c.Resources.Limits.Memory().IsZero() {
					noLimits++
				}
				if c.Resources.Requests.Cpu().IsZero() && c.Resources.Requests.Memory().IsZero() {
					noRequests++
				}
			}
			for _, cs := range pod.Status.ContainerStatuses {
				if cs.RestartCount > 5 {
					restarts++
					warnings = append(warnings, fmt.Sprintf("🟡 Pod %s/%s has %d restarts", pod.Namespace, pod.Name, cs.RestartCount))
				}
			}
		}
	}

	if noLimits > 0 {
		warnings = append(warnings, fmt.Sprintf("🟡 %d containers without resource limits", noLimits))
	}
	if noRequests > 0 {
		warnings = append(warnings, fmt.Sprintf("🟡 %d containers without resource requests", noRequests))
	}

	// Check metrics-server
	if !k.HasMetrics {
		warnings = append(warnings, "🔵 metrics-server not installed — install for live CPU/memory usage data")
	}

	if len(warnings) == 0 {
		fmt.Println("  ✅ No issues found — cluster looks healthy!")
	} else {
		for _, w := range warnings {
			fmt.Printf("  %s\n", w)
		}
	}
	fmt.Println()
}
