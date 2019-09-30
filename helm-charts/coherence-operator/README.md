# coherence-operator
Install coherence-operator to work with Coherence clusters on Kubernetes.

## Introduction

This chart install a coherence-operator deployment on a 
[Kubernetes](https://kubernetes.io) cluster using the [Helm](https://helm.sh)
package manager.

## Prerequisites
* Kubernetes 1.10.3 or above
* Helm 2.11.0 or above

## Installing the Chart
To install the chart with the release name `sample-coherence-operator`:

```
@ helm install --name sample-coherence-operator coherence-operator
```

The command deploys coherence-operator on the Kubernetes cluster in the
default configuration. The [configuration](#configuration) section list
parameters that can be configured during installation.

## Uninstalling the Chart
To uninstall the `sample-coherence-operator` deployment:

```
$ helm delete sample-coherence-operator
```

The command removes all the Kubernetes components associated with the chart
and deletes the release.

```
$ kubectl delete secret coherence-monitoring-config
```

