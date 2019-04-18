<!--
Copyright 2018, Oracle Corporation and/or its affiliates.
All rights reserved.  Licensed under the Universal
Permissive License v 1.0 as shown at
http://oss.oracle.com/licenses/upl.

-->

Oracle is finding ways for organizations using Oracle Coherence to move
their important workloads into the cloud. By certifying on industry
standards, such as Docker and Kubernetes, Coherence now runs in a cloud
neutral infrastructure. In addition, we've provided an open-source
Oracle Coherence Kubernetes Operator (the “operator”) which has several
key features to assist you with deploying and managing Coherence
clusters in a Kubernetes environment. You can:

* Run the Coherence you know and love on the industry standard
  Kubernetes container orchestration framework, using Docker containers
  for the system elements of Coherence.

* Use popular industry standard tools such as
  [Grafana](https://grafana.com/),
  [ELK](https://www.elastic.co/elk-stack), and
  [Prometheus](https://prometheus.io/) to monitor the performance,
  logs and and health from your clusters.

* Flexibly override and customize cluster configuration.

* Scale the Coherence deployment.

* Use
  [Coherence*Extend](https://docs.oracle.com/middleware/12213/coherence/develop-remote-clients/building-your-first-extend-application.htm#COHCG5033)
  to access your cluster with a variety of clients.

* Start clusters based on declarative startup parameters and desired
  states.

* Use Kubernetes persistent volumes to provide the storage for "storage
enabled" cluster nodes.

* Deploy custom code for your `EntryProcessor` classes and other
server-side Coherence constructs.

The fastest way to experience the operator is to follow the [Quick Start
guide](https://oracle.github.io/coherence-operator/docs/quickstart.html).
More information is available on the [Documentation
Site](https://oracle.github.io/coherence-operator/).
