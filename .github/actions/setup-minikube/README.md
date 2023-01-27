Setup Minikube GitHub Action
===============================

[<img src="https://github.com/manusa/actions-setup-minikube/workflows/Perform checks/badge.svg"/>](https://github.com/manusa/actions-setup-minikube/actions)
[<img src="https://github.com/manusa/actions-setup-minikube/workflows/Run action and validate environment/badge.svg"/>](https://github.com/manusa/actions-setup-minikube/actions)

Set up your GitHub Actions workflow with a specific version of
[Minikube](https://github.com/kubernetes/minikube)
and [Kubernetes](https://github.com/kubernetes/kubernetes).

_Currently only Linux Ubuntu 18.04, 20.04, or 22.04
[CI environment](https://help.github.com/en/github/automating-your-workflow-with-github-actions/virtual-environments-for-github-actions)
is supported._

## Usage

### Basic

```yaml
name: Example workflow

on: [push]

jobs:
  example:
    name: Example Minikube-Kubernetes Cluster interaction
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3
      - name: Setup Minikube
        uses: manusa/actions-setup-minikube@v2.7.2
        with:
          minikube version: 'v1.28.0'
          kubernetes version: 'v1.25.4'
          github token: ${{ secrets.GITHUB_TOKEN }}
      - name: Interact with the cluster
        run: kubectl get nodes
```

### Required input parameters

| Parameter            | Description                                                                       |
|----------------------|-----------------------------------------------------------------------------------|
| `minikube version`   | Minikube [version](https://github.com/kubernetes/minikube/releases) to deploy     |
| `kubernetes version` | Kubernetes [version](https://github.com/kubernetes/kubernetes/releases) to deploy |

### Optional input parameters

| Parameter           | Description                                                                                                                              |
|---------------------|------------------------------------------------------------------------------------------------------------------------------------------|
| `github token`      | GITHUB_TOKEN secret value to access GitHub REST API with an unlimited number of requests (optional but recommended)                      |
| `driver`            | Minikube [driver](https://minikube.sigs.k8s.io/docs/drivers/) to use. This action supports `none` (default if not specified) or `docker` |
| `container runtime` | The container runtime to be used (valid options: docker, cri-o, containerd)                                                              |
| `start args`        | Additional arguments to append to [`minikube start`](https://minikube.sigs.k8s.io/docs/commands/start/) command                          |

## License

The scripts and documentation in this project are released under the [Apache 2.0](./LICENSE) license.
