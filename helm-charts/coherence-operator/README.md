<!--
  Copyright 2020, 2025, Oracle Corporation and/or its affiliates.
  Licensed under the Universal Permissive License v 1.0 as shown at
  http://oss.oracle.com/licenses/upl.
-->

# coherence-operator
Install coherence-operator to work with Coherence clusters on Kubernetes.

## Introduction

This chart install a coherence-operator deployment on a 
[Kubernetes](https://kubernetes.io) cluster using the [Helm](https://helm.sh)
package manager.

## Prerequisites
* A non end of life version of Kubernetes (the operator is tested on all [currently supported Kubernetes versions](https://kubernetes.io/releases/))
* Helm 3 or above

## Installing the Chart
To install the chart with the release name `sample-coherence-operator`:

```
@ helm install sample-coherence-operator coherence-operator
```

The command deploys coherence-operator on the Kubernetes cluster using the
default configuration. 
See the [Coherence Operator installation guide](https://docs.coherence.community/coherence-operator/docs/latest/docs/installation/012_install_helm)
for full details of how to 
install the operator using Helm

## Uninstalling the Chart
To uninstall the `sample-coherence-operator` deployment:

```
$ helm delete sample-coherence-operator
```

