package main

import (
	"fmt"
	"os"

	"github.com/prithiviraj123/kpulse/pkg/alerts"
	"github.com/prithiviraj123/kpulse/pkg/cluster"
	"github.com/prithiviraj123/kpulse/pkg/cost"
	"github.com/prithiviraj123/kpulse/pkg/nodes"
	"github.com/prithiviraj123/kpulse/pkg/pods"
)

const version = "0.1.0"

func main() {
	cmd := "all"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "version", "--version", "-v":
		fmt.Printf("kpulse v%s\n", version)
		return
	case "help", "--help", "-h":
		printHelp()
		return
	}

	k, err := cluster.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to connect to cluster: %v\n", err)
		os.Exit(1)
	}

	switch cmd {
	case "all":
		cluster.Show(k)
		nodes.Show(k)
		pods.Show(k)
		cost.Show(k)
		alerts.Show(k)
	case "nodes":
		cluster.Show(k)
		nodes.Show(k)
	case "pods":
		cluster.Show(k)
		pods.Show(k)
	case "cost":
		cluster.Show(k)
		cost.Show(k)
	case "alerts":
		cluster.Show(k)
		alerts.Show(k)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`kpulse - Kubernetes cluster health & cost snapshot

Usage:
  kpulse              Show full cluster overview
  kpulse nodes        Show node resource usage
  kpulse pods         Show top pods by resource usage
  kpulse cost         Show node cost estimates
  kpulse alerts       Show warnings and recommendations
  kpulse version      Show version`)
}
