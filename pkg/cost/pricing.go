package cost

import (
	"context"
	"fmt"

	"github.com/prithiviraj123/kpulse/pkg/cluster"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// On-demand pricing per hour (us-east-1, Linux)
var pricing = map[string]float64{
	"t3.micro":    0.0104,
	"t3.small":    0.0208,
	"t3.medium":   0.0416,
	"t3.large":    0.0832,
	"t3.xlarge":   0.1664,
	"t3.2xlarge":  0.3328,
	"t3a.micro":   0.0094,
	"t3a.small":   0.0188,
	"t3a.medium":  0.0376,
	"t3a.large":   0.0752,
	"m5.large":    0.0960,
	"m5.xlarge":   0.1920,
	"m5.2xlarge":  0.3840,
	"m5.4xlarge":  0.7680,
	"m5a.large":   0.0860,
	"m5a.xlarge":  0.1720,
	"m6i.large":   0.0960,
	"m6i.xlarge":  0.1920,
	"m6i.2xlarge": 0.3840,
	"m7i.large":   0.1008,
	"m7i.xlarge":  0.2016,
	"c5.large":    0.0850,
	"c5.xlarge":   0.1700,
	"c5.2xlarge":  0.3400,
	"c6i.large":   0.0850,
	"c6i.xlarge":  0.1700,
	"r5.large":    0.1260,
	"r5.xlarge":   0.2520,
	"r6i.large":   0.1260,
	"r6i.xlarge":  0.2520,
}

func Show(k *cluster.KClient) {
	nodeList, err := k.Client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil || nodeList == nil {
		return
	}

	fmt.Println("\n💰 COST ESTIMATE (monthly, on-demand, us-east-1)")
	fmt.Println("┌────────────────────────┬──────────────┬──────────┬──────────────┐")
	fmt.Println("│ Node                   │ Instance     │ $/hour   │ $/month      │")
	fmt.Println("├────────────────────────┼──────────────┼──────────┼──────────────┤")

	var totalMonthly float64
	for _, node := range nodeList.Items {
		instanceType := node.Labels["node.kubernetes.io/instance-type"]
		if instanceType == "" {
			instanceType = node.Labels["beta.kubernetes.io/instance-type"]
		}

		name := node.Name
		if len(name) > 22 {
			name = name[:22]
		}

		hourly, ok := pricing[instanceType]
		if !ok {
			hourly = 0
		}
		monthly := hourly * 730 // avg hours per month
		totalMonthly += monthly

		if hourly > 0 {
			fmt.Printf("│ %-22s │ %-12s │ $%-7.4f │ $%-11.2f │\n", name, instanceType, hourly, monthly)
		} else {
			fmt.Printf("│ %-22s │ %-12s │ %-8s │ %-12s │\n", name, instanceType, "N/A", "N/A")
		}
	}

	fmt.Println("├────────────────────────┼──────────────┼──────────┼──────────────┤")
	fmt.Printf("│ %-22s │              │          │ $%-11.2f │\n", "TOTAL", totalMonthly)
	fmt.Println("└────────────────────────┴──────────────┴──────────┴──────────────┘")
}
