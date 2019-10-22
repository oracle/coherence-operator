# Docker Images for Oracle Coherence Operator

> **Note:** The images from this repo are not intended to be used directly, but only via Coherence Operator Helm Chart,
> as documented in the [Quick Start Guide](https://oracle.github.io/coherence-operator/docs/2.0.0/#/about/03_quickstart).
>
> If, for some reason, you want to pull them to your local system directly, using `docker pull` command, please note
> that we do not publish an image with the `latest` tag, so you will have to pull the specific, version-tagged image instead.
> The list of available versions/tags is available on the `Tags` tab above.

Oracle enables organizations using [Coherence](https://www.oracle.com/technetwork/middleware/coherence/overview/index.html) to move their clusters
into the cloud. By supporting de facto standards such as Docker and Kubernetes, Oracle facilitates running Coherence on cloud-neutral
infrastructure. In particular, Oracle provides an open-source Coherence Operator, which implements features to assist
with deploying and managing Coherence clusters in a Kubernetes environment. You can:


* Run Coherence within the industry standard Kubernetes container orchestration framework, using Docker containers for the members of a Coherence cluster.
* Flexibly override and customize cluster configuration using a `CoherenceCluster` custom resource definition.
* Safely scale the roles of a Coherence cluster using Kubernetes verbs or updates.
* Use
  [Coherence*Extend](https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/develop-remote-clients/building-your-first-extend-application.html#GUID-2E360BF7-1541-4997-97F2-9D3739AE3B48)
  to access your cluster with a variety of clients.
* Start clusters based on declarative startup parameters and desired states.
* Deploy custom code for `EntryProcessor` classes and other server-side Coherence constructs.
* Use Kubernetes persistent volumes when using Coherence’s disk-based storage features Elastic Data or Persistence.
* Use Kubernetes Zone information to ensure data stored in Coherence is resilient to loss of a Zone. Coherence goes to great efforts to ensure data is safe across processes, machines, racks and sites. When Coherence is deployed to Kubernetes with the Coherence Operator, data will be spread across zones to ensure this underlying principle is supported; thus by default, loss of any zone is a tolerated failure mode. This is reflected in the StatusHA value (SITE-SAFE) for partitioned services, in addition to the member level site information that is equivalent to the kubernetes zone label on the associated pod.
* Use popular industry standard tools such as
  [Grafana](https://grafana.com/),
  [ELK](https://www.elastic.co/elk-stack) (or more specifically the EFK stack including Fluentd), and
  [Prometheus](https://prometheus.io/)
  to monitor the performance, logs and and health of your clusters.


## Software and Version Prerequisites

* Kubernetes 1.11.3+ cluster
* Access to Oracle Coherence 12.2.1.3.2 or 12.2.1.4.0 Docker images

## Getting Started

  The quickest way to use and experience the Coherence Operator is to follow the [Quick Start guide](https://oracle.github.io/coherence-operator/docs/2.0.0/#/about/03_quickstart)

## Licenses

Coherence Kubernetes Operator images in this repository folder are licensed under the (Universal Permissive License 1.0) [http://oss.oracle.com/licenses/upl].

Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
