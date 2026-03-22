# KPulse

Kubernetes cluster health & cost snapshot — one command, full visibility.

## What it shows

```
$ kpulse

╔══════════════════════════════════════════════════════════════╗
║  🔵 KPULSE — Cluster: eks-demo                              ║
║     Version: v1.34   Region: us-east-1   Nodes: 2           ║
╚══════════════════════════════════════════════════════════════╝

📦 NODES          — CPU/MEM requests, usage, pod count per node
🔥 TOP PODS       — Top 10 pods by CPU usage or requests
💰 COST ESTIMATE  — Monthly cost per node (AWS on-demand pricing)
⚠️  ALERTS         — Unhealthy nodes, missing limits, high restarts
```

## Install

### CLI (standalone binary)

```bash
go install github.com/kpulse/kpulse/cmd/kpulse@latest
```

### Helm

```bash
helm install kpulse ./charts/kpulse -n kpulse --create-namespace
```

## Usage

```bash
kpulse              # Full cluster overview
kpulse nodes        # Node resource usage only
kpulse pods         # Top pods by CPU/memory
kpulse cost         # Node cost estimates
kpulse alerts       # Warnings and recommendations
kpulse version      # Show version
```

## Features

- Works **with or without metrics-server** — shows requests/limits always, adds live usage when metrics-server is available
- **Cost estimation** with hardcoded AWS pricing for common instance types
- **Alerts** for NotReady nodes, containers without resource limits, high restart counts
- **Read-only** — only needs `get` and `list` permissions
- **Single binary** — no dependencies, runs anywhere with a kubeconfig

## Requirements

- `kubectl` configured with cluster access (kubeconfig)
- Optional: metrics-server for live CPU/memory usage data
