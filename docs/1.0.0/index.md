<!--
Copyright 2018, 2019, Oracle Corporation and/or its affiliates.
All rights reserved.  Licensed under the Universal
Permissive License v 1.0 as shown at
http://oss.oracle.com/licenses/upl.

-->

# Coherence Operator Documentation

Oracle enables organizations using Coherence to move their clusters into the cloud. By supporting industry standards, such as Docker and Kubernetes, Oracle facilitates running Coherence on cloud-neutral infrastructure. In addition, Oracle provides an open-source Coherence Operator ("the operator"), which implements features to assist with deploying and managing Coherence clusters in a Kubernetes environment. You can:

* Run Coherence within the de facto standard Kubernetes container
  orchestration framework, using Docker containers for the members of a
  Coherence cluster.

* Use popular industry standard tools such as
  [Grafana](https://grafana.com/),
  [ELK](https://www.elastic.co/elk-stack) (or more specifically the EFK stack including Fluentd), and
  [Prometheus](https://prometheus.io/)
  to monitor the performance, logs and and health of your clusters.

* Flexibly override and customize cluster configuration.

* Scale the Coherence deployment.

* Use
  [Coherence*Extend](https://docs.oracle.com/middleware/12213/coherence/develop-remote-clients/building-your-first-extend-application.htm#COHCG5033)
  to access your cluster with a variety of clients.

* Use Kubernetes Zone information to ensure data stored in Coherence is resilient to loss of a Zone.
  Coherence goes to great efforts to ensure data is safe across processes, machines, racks and sites. When Coherence is deployed to Kubernetes with the Coherence Operator, data will be spread across zones to ensure this underlying principle is supported; thus by default, loss of any zone is a tolerated failure mode. This is reflected in the StatusHA value (SITE-SAFE) for partitioned services, in addition to the member level site information that is equivalent to the kubernetes zone label on the associated pod.

* Start clusters based on declarative startup parameters and desired
  states.

* Use Kubernetes persistent volumes when using Coherence's disk-based
  storage features Elastic Data or Persistence.

* Deploy custom code for your `EntryProcessor` classes and other
server-side Coherence constructs.

The fastest way to experience the operator is to follow the [Quick Start guide](docs/quickstart.md), or try out the
[samples](docs/samples/README.md).

## About This Documentation

This documentation includes sections targeted to different audiences. To help you find what you are looking for more easily, consult this table of contents:

* The [Quick Start Guide](docs/quickstart.md) explains how to
  quickly get Coherence running on Kubernetes, using the defaults, nothing special.

* The [User Guide](docs/user-guide.md) contains detailed usage
  information on the Coherence Operator, including how to install and  configure the operator and several common use cases.

* The [Samples](docs/samples/README.md) provide detailed example code and instructions that show you how to perform various tasks related to the operator.

* The [Developer Guide](docs/developer.md) provides details for users who want to understand how the operator is built, tested, and so on. Those who wish to contribute to the operator code will find useful information here.

* The [Access the EFK (Elasticsearch, Fluentd and Kibana) Stack to View Logs](docs/logcapture.md) page describes how to enable log capture, and manage data logging through the EFK stack to view logs.

* The [Monitor Coherence Services via Grafana Dashboards](docs/prometheusoperator.md) page explains how to configure the Prometheus Operator and monitor Coherence services through Grafana dashboards.
  
<!--
* The [Contributing](#contributing-to-the-operator) section provides information about contribution requirements.
-->

## User Guide

The [User Guide](docs/user-guide.md) provides detailed information on all aspects of using the operator including:

* Installing and configuring the operator.

* Using the operator to create and manage Coherence clusters.

* Manually creating Coherence clusters to be managed by the operator.

* Configuring Elasticsearch and Kibana to access the operator's log files.

* Shutting down clusters.

* And much more!

## Samples

Refer to our [samples](docs/samples/README.md) for
information about the available sample code.

## Things to Keep In Mind for Existing Coherence Users

Software running in Kubernetes must provide "health checks" so that Kubernetes can make informed decisions about starting, stopping, or even killing, the containers running the software.  The operator provides everything required to do this for Coherence.  Keep in mind  that these health checks cause frequent `MemberJoined` and  `MemberLeft` events to happen.  If these events refer to something  like `OracleCoherenceK8sPodChecker`, they are normal and be safely ignored.

## Need more help? Have a suggestion? Come and say "Hello!"

### Join Our Public Slack Channel

We have a **public Slack channel** where you can get in touch with us to
ask questions about using the operator or give us feedback or
suggestions about what features and improvements you would like to see.
We would love to hear from you. To join our channel, please [visit this
site to get an
invitation](https://join.slack.com/t/oraclecoherence/shared_invite/enQtNzcxNTQwMTAzNjE4LTJkZWI5ZDkzNGEzOTllZDgwZDU3NGM2YjY5YWYwMzM3ODdkNTU2NmNmNDFhOWIxMDZlNjg2MzE3NmMxZWMxMWE). The
invitation email will include details of how to access our Slack workspace.  After you are logged in, please come to `#operator` and say, "hello!"
