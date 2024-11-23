## Introduction
Saver is a simple kubectl plugin that helps you save your deployed workloads from your Kubernetes cluster to you localhost.
It helps you simply save your workloads with 'YAML' format.
So you can easily backup your workloads and try to deploy them somewhere.
## Features
Supported workloads for now: 
* Deployment
* Statefulset
* Daemonset

## Installation
1. Build the deployment
    ```bash
    go build -o kubectl-save
    ```
2. Copy the binary plugin to `/usr/local/bin`
    ```bash
    sudo cp kubectl-save /usr/local/bin/
    ```
## Usage
After installation, you can easily save all deployments, daemonsets, and statefulsets in a specific namespace using simple commands like these:
```bash
kubectl save deployment --namespace controller -o controller-deployments.yaml
kubectl save daemonset --namespace controller -o controller-daemonsets.yaml
kubectl save statefulset --namespace controller -o controller-statefulsets.yaml
```
