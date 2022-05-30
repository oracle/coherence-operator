<!--
  Copyright 2021, Oracle Corporation and/or its affiliates.
  Licensed under the Universal Permissive License v 1.0 as shown at
  http://oss.oracle.com/licenses/upl.
-->

# Coherence Helm Chart Example

An example of a Helm chart to manage Coherence resources in Kubernetes clusters.

## Introduction

This chart manages a `Coherence` resource on a 
[Kubernetes](https://kubernetes.io) cluster using the [Helm](https://helm.sh)
package manager.

## Prerequisites
* Kubernetes 1.18 or above
* Helm 3 or above

## Installing the Chart
To install the chart with the release name `coherence-example`:

```
@ helm install coherence-example my-coherence-demo
```

The command deploys a Coherence resource in the Kubernetes cluster.

## Uninstalling the Chart
To uninstall the `my-coherence-demo` deployment:

```
$ helm delete my-coherence-demo
```

