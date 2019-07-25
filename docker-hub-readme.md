# Docker Images for Coherence Operator

> **Note:** The images from this repository must be used only through Coherence Operator Helm Chart as documented in the [Quick Start Guide](https://oracle.github.io/coherence-operator/docs/quickstart.html), and are not intended to be used directly.
>
> If you want to pull them to your local system directly using the `docker pull` command, note that we do not publish an image with the `latest` tag. Instead, you have to pull the image tagged with specific version. The list of available versions or tags is available under **Tags**.

Oracle enables organizations using [Coherence](https://www.oracle.com/technetwork/middleware/coherence/overview/index.html) to move their clusters
into the cloud. By supporting industry standards such as Docker and Kubernetes, Oracle facilitates running Coherence on cloud-neutral infrastructure. In particular, Oracle provides an open source Coherence Operator, which implements features to assist with deploying and managing Coherence clusters in a Kubernetes environment. You can:

* Run Coherence within the de facto standard Kubernetes container
  orchestration framework, using Docker containers for the members of a
  Coherence cluster.

* Use popular industry standard tools such as
  [Grafana](https://grafana.com/),
  [ELK](https://www.elastic.co/elk-stack) (or more specifically the EFK stack including Fluentd), and
  [Prometheus](https://prometheus.io/)
  to monitor the performance, logs, and health of your clusters.

* Flexibly override and customize cluster configuration.

* Scale the Coherence deployment using Kubernetes verbs or updates.

* Use
  [Coherence*Extend](https://docs.oracle.com/middleware/12213/coherence/develop-remote-clients/building-your-first-extend-application.htm#COHCG5033)
  to access your cluster with a variety of clients.

* Start clusters based on declarative startup parameters and desired
  states.

* Use Kubernetes persistent volumes when using Coherence's disk-based
  storage features Elastic Data or Persistence.

* Deploy custom code for `EntryProcessor` classes and other
server-side Coherence constructs.

## Software and Version Prerequisites

* Kubernetes 1.11.5+, 1.12.3+, 1.13.0+ (check with `kubectl version`)
* Docker 18.03.1-ce (check with `docker version`)
* Flannel networking v0.10.0-amd64 (check with `docker images | grep flannel`)
* Helm 2.12.3 or above (and all of its prerequisites)
* Oracle Coherence 12.2.1.3

## Getting Started

  The following documentation includes sections targeted to different audiences. To help you find what you are looking for more easily, consult this table of contents:

  * The [Quick Start guide](https://oracle.github.io/coherence-operator/docs/quickstart.html) explains how to
    quickly get Coherence running on Kubernetes, using the defaults, nothing special.

  * The [User guide](https://oracle.github.io/coherence-operator/docs/user-guide.html) contains detailed usage information on the Coherence Operator, including how to install and configure the operator and several common use cases.

  * The [Samples](https://oracle.github.io/coherence-operator/docs/samples/) provide detailed example code and instructions that show you how to perform various tasks
    related to the operator.

  * The [Developer guide](https://oracle.github.io/coherence-operator/docs/developer.html) provides details for people
    who want to understand how the operator is built and tested. Those who want to contribute to the operator code can find useful information here.

## Licenses

Coherence Operator images in this repository folder are licensed under the [Universal Permissive License 1.0](http://oss.oracle.com/licenses/upl).

Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
