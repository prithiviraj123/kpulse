# KPulse

**Kubernetes cluster health & cost snapshot — one command, full visibility.**

KPulse is a lightweight CLI tool that gives you an instant overview of your Kubernetes cluster — node resources, top pods, cost estimates, and health alerts — all in a single command.

---

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Usage](#usage)
- [Sample Output](#sample-output)
- [In-Cluster Deployment](#in-cluster-deployment)
- [How It Works](#how-it-works)
- [Configuration](#configuration)
- [Supported Platforms](#supported-platforms)
- [Helm Chart](#helm-chart)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

---

## Features

| Feature | Description |
|---------|-------------|
| **Cluster Overview** | Cluster name, Kubernetes version, region, node count |
| **Node Resources** | CPU & memory requests vs capacity, pod count, instance type per node |
| **Top Pods** | Top 10 pods sorted by CPU usage or requests |
| **Cost Estimates** | Monthly cost per node based on AWS on-demand pricing |
| **Health Alerts** | NotReady nodes, containers without limits, high restart counts |
| **Metrics Server Support** | Shows live CPU/memory usage when metrics-server is installed |
| **In-Cluster Mode** | Runs as a pod inside the cluster with ServiceAccount auth |
| **Zero Dependencies** | Single binary, no external tools required |

---

## Quick Start

```bash
# Install
brew tap prithiviraj123/tap
brew install kpulse

# Run
kpulse
```

That's it. If `kubectl` works, `kpulse` works.

---

## Installation

### Homebrew (macOS / Linux)

```bash
brew tap prithiviraj123/tap
brew install kpulse
```

### Linux (amd64)

```bash
curl -sL https://github.com/prithiviraj123/kpulse/releases/download/v0.1.3/kpulse_0.1.3_linux_amd64.tar.gz | tar xz
sudo mv kpulse /usr/local/bin/
```

### Linux (ARM64 — Graviton, Raspberry Pi)

```bash
curl -sL https://github.com/prithiviraj123/kpulse/releases/download/v0.1.3/kpulse_0.1.3_linux_arm64.tar.gz | tar xz
sudo mv kpulse /usr/local/bin/
```

### macOS (Apple Silicon)

```bash
curl -sL https://github.com/prithiviraj123/kpulse/releases/download/v0.1.3/kpulse_0.1.3_darwin_arm64.tar.gz | tar xz
sudo mv kpulse /usr/local/bin/
```

### macOS (Intel)

```bash
curl -sL https://github.com/prithiviraj123/kpulse/releases/download/v0.1.3/kpulse_0.1.3_darwin_amd64.tar.gz | tar xz
sudo mv kpulse /usr/local/bin/
```

### Windows

Download from [GitHub Releases](https://github.com/prithiviraj123/kpulse/releases/download/v0.1.3/kpulse_0.1.3_windows_amd64.zip), extract, and add to PATH.

### Go Install

```bash
go install github.com/prithiviraj123/kpulse/cmd/kpulse@latest
```

### Docker

```bash
docker run -v ~/.kube:/root/.kube prithiviraj123/kpulse:v0.1.3
```

---

## Usage

```bash
kpulse              # Full cluster overview (nodes + pods + cost + alerts)
kpulse nodes        # Node resource usage only
kpulse pods         # Top pods by CPU/memory
kpulse cost         # Node cost estimates
kpulse alerts       # Warnings and recommendations
kpulse version      # Show version
kpulse help         # Show help
```

---

## Sample Output

### Full Overview (`kpulse`)

```
🔵 KPULSE
┌──────────────┬──────────────────┐
│ Cluster      │ eks-demo         │
│ Version      │ v1.34.4          │
│ Region       │ us-east-1        │
│ Nodes        │ 1                │
└──────────────┴──────────────────┘

📦 NODES
  Node 1
  ┌──────────────────┬──────────────────────────┐
  │ Name             │ 172.31.6.7               │
  │ Instance         │ t3.medium                │
  │ CPU (req/cap)    │ 370m / 1930m             │
  │ Memory (req/cap) │ 204Mi / 3.2Gi            │
  │ Pods             │ 6 / 17                   │
  └──────────────────┴──────────────────────────┘

🔥 TOP PODS (by requests) — showing 6 of 6 running
┌──────────────────────────┬──────────────┬──────────┬──────────┐
│ Pod                      │ Namespace    │ CPU      │ Memory   │
├──────────────────────────┼──────────────┼──────────┼──────────┤
│ coredns-7d58d485c9-jbh4l │ kube-system  │ 100m req │ 70Mi     │
│ coredns-7d58d485c9-tktwt │ kube-system  │ 100m req │ 70Mi     │
│ kube-proxy-m6zpv         │ kube-system  │ 100m req │ 0Mi      │
│ aws-node-pj5pb           │ kube-system  │ 50m req  │ 0Mi      │
│ kpulse-df54cd555-7dqnm   │ kpulse       │ 10m req  │ 32Mi     │
└──────────────────────────┴──────────────┴──────────┴──────────┘

💰 COST ESTIMATE (monthly, on-demand, us-east-1)
┌────────────────────────┬──────────────┬──────────┬──────────────┐
│ Node                   │ Instance     │ $/hour   │ $/month      │
├────────────────────────┼──────────────┼──────────┼──────────────┤
│ 172.31.6.7             │ t3.medium    │ $0.0416  │ $30.37       │
├────────────────────────┼──────────────┼──────────┼──────────────┤
│ TOTAL                  │              │          │ $30.37       │
└────────────────────────┴──────────────┴──────────┴──────────────┘

⚠️  ALERTS
  🟡 3 containers without resource limits
  🔵 metrics-server not installed — install for live CPU/memory usage data
```

### Nodes Only (`kpulse nodes`)

```
🔵 KPULSE
┌──────────────┬──────────────────┐
│ Cluster      │ eks-demo         │
│ Version      │ v1.34.4          │
│ Region       │ us-east-1        │
│ Nodes        │ 1                │
└──────────────┴──────────────────┘

📦 NODES
  Node 1
  ┌──────────────────┬──────────────────────────┐
  │ Name             │ 172.31.6.7               │
  │ Instance         │ t3.medium                │
  │ CPU (req/cap)    │ 360m / 1930m             │
  │ Memory (req/cap) │ 172Mi / 3.2Gi            │
  │ Pods             │ 5 / 17                   │
  └──────────────────┴──────────────────────────┘
```

---

## In-Cluster Deployment

Deploy kpulse inside any Kubernetes cluster with a single command:

### Install

```bash
kubectl apply -f https://raw.githubusercontent.com/prithiviraj123/kpulse/master/deploy/components.yaml
```

This creates:
- `kpulse` namespace
- ServiceAccount with read-only ClusterRole
- Deployment running the kpulse container

### Run

```bash
# Full overview
kubectl exec -n kpulse deploy/kpulse -- kpulse

# Nodes only
kubectl exec -n kpulse deploy/kpulse -- kpulse nodes

# Pods only
kubectl exec -n kpulse deploy/kpulse -- kpulse pods

# Cost only
kubectl exec -n kpulse deploy/kpulse -- kpulse cost

# Alerts only
kubectl exec -n kpulse deploy/kpulse -- kpulse alerts
```

### Uninstall

```bash
kubectl delete -f https://raw.githubusercontent.com/prithiviraj123/kpulse/master/deploy/components.yaml
```

---

## How It Works

```
kpulse binary
    │
    ├── Detects environment (kubeconfig or in-cluster ServiceAccount)
    │
    ├── Connects to Kubernetes API server
    │   ├── GET /api/v1/nodes           → Node capacity, labels, conditions
    │   ├── GET /api/v1/pods            → Pod requests, limits, status, restarts
    │   └── GET /apis/metrics.k8s.io    → Live usage (optional, if metrics-server exists)
    │
    ├── Aggregates data
    │   ├── CPU/memory requests per node
    │   ├── Pod count per node
    │   ├── Instance type → cost lookup
    │   └── Health checks (NotReady, no limits, restarts)
    │
    └── Prints formatted tables to stdout
```

### With vs Without Metrics Server

| Data | Without metrics-server | With metrics-server |
|------|----------------------|-------------------|
| Node CPU/Memory | Requests vs capacity | Actual usage + requests |
| Pod CPU/Memory | Requested amounts | Real-time usage |
| Cost estimates | ✅ Works | ✅ Works |
| Alerts | ✅ Works | ✅ Works + usage-based |

KPulse works without metrics-server. When metrics-server is available, it automatically shows live usage data alongside requests.

---

## Configuration

KPulse requires no configuration files. It uses:

| Setting | Source |
|---------|--------|
| Cluster connection | `~/.kube/config` (local) or ServiceAccount (in-cluster) |
| Current context | Whatever `kubectl` is pointing to |
| Cost pricing | Built-in AWS on-demand pricing table (us-east-1) |

### Supported Instance Types for Cost Estimation

KPulse includes pricing for common AWS instance types:

- **T3 family:** t3.micro, t3.small, t3.medium, t3.large, t3.xlarge, t3.2xlarge
- **T3a family:** t3a.micro, t3a.small, t3a.medium, t3a.large
- **M5 family:** m5.large, m5.xlarge, m5.2xlarge, m5.4xlarge
- **M5a family:** m5a.large, m5a.xlarge
- **M6i family:** m6i.large, m6i.xlarge, m6i.2xlarge
- **M7i family:** m7i.large, m7i.xlarge
- **C5 family:** c5.large, c5.xlarge, c5.2xlarge
- **C6i family:** c6i.large, c6i.xlarge
- **R5 family:** r5.large, r5.xlarge
- **R6i family:** r6i.large, r6i.xlarge

For unlisted instance types, cost shows as "N/A".

---

## Supported Platforms

### Binaries

| OS | Architecture | Download |
|----|-------------|----------|
| macOS | ARM64 (Apple Silicon) | `kpulse_0.1.3_darwin_arm64.tar.gz` |
| macOS | amd64 (Intel) | `kpulse_0.1.3_darwin_amd64.tar.gz` |
| Linux | amd64 | `kpulse_0.1.3_linux_amd64.tar.gz` |
| Linux | ARM64 | `kpulse_0.1.3_linux_arm64.tar.gz` |
| Windows | amd64 | `kpulse_0.1.3_windows_amd64.zip` |
| Windows | ARM64 | `kpulse_0.1.3_windows_arm64.zip` |

### Container Image

```bash
docker pull prithiviraj123/kpulse:v0.1.3
docker pull prithiviraj123/kpulse:latest
```

### Kubernetes Compatibility

Tested with:
- EKS (v1.28 — v1.34)
- GKE
- AKS
- Self-managed clusters
- k3s / kind / minikube

---

## Helm Chart

For teams that prefer Helm:

```bash
helm install kpulse ./charts/kpulse -n kpulse --create-namespace
```

The Helm chart deploys:
- ServiceAccount (`kpulse`)
- ClusterRole (read-only: nodes, pods, metrics)
- ClusterRoleBinding
- Deployment

### Values

| Key | Default | Description |
|-----|---------|-------------|
| `image.repository` | `prithiviraj123/kpulse` | Container image |
| `image.tag` | `0.1.3` | Image tag |
| `image.pullPolicy` | `IfNotPresent` | Pull policy |
| `serviceAccount.name` | `kpulse` | ServiceAccount name |

---

## RBAC Permissions

KPulse requires minimal read-only permissions:

```yaml
rules:
  - apiGroups: [""]
    resources: ["nodes", "pods", "namespaces"]
    verbs: ["get", "list"]
  - apiGroups: ["metrics.k8s.io"]
    resources: ["nodes", "pods"]
    verbs: ["get", "list"]
```

No write access. No secrets access. Read-only.

---

## Roadmap

- [ ] Namespace breakdown (resource usage per namespace)
- [ ] Deployment/StatefulSet status overview
- [ ] GCP and Azure cost estimation
- [ ] `--watch` mode (auto-refresh every N seconds)
- [ ] `--output json` for scripting
- [ ] Custom pricing file support
- [ ] Kubectl plugin (`kubectl kpulse`)
- [ ] Web dashboard mode

---

## Contributing

1. Fork the repo
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Commit changes (`git commit -m 'Add my feature'`)
4. Push (`git push origin feature/my-feature`)
5. Open a Pull Request

### Build from source

```bash
git clone https://github.com/prithiviraj123/kpulse.git
cd kpulse
go build -o kpulse ./cmd/kpulse/
./kpulse version
```

---

## License

MIT License — see [LICENSE](LICENSE) for details.

---

## Links

- **GitHub:** https://github.com/prithiviraj123/kpulse
- **Releases:** https://github.com/prithiviraj123/kpulse/releases
- **Docker Hub:** https://hub.docker.com/r/prithiviraj123/kpulse
- **Homebrew:** `brew tap prithiviraj123/tap && brew install kpulse`
