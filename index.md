<!--
Copyright 2018, 2019, Oracle Corporation and/or its affiliates.
All rights reserved.  Licensed under the Universal
Permissive License v 1.0 as shown at
http://oss.oracle.com/licenses/upl.

-->

---------

## Join our Public Slack Channel

We have a **public Slack channel** where you can get in touch with us to
ask questions about using the operator or give us feedback or
suggestions about what features and improvements you would like to see.
We would love to hear from you. To join our channel, please [visit this
site to get an
invitation](https://join.slack.com/t/oraclecoherence/shared_invite/enQtNjA3MTU3MTk0MTE3LWZhMTdhM2E0ZDY2Y2FmZDhiOThlYzJjYTc5NzdkYWVlMzUzODZiNTI4ZWU3ZTlmNDQ4MmE1OTRhOWI1MmIxZjQ).  The
invitation email will include details of how to access our Slack
workspace.  After you are logged in, please come to `#operator` and say,
"hello!"

--------


Oracle is finding ways for organizations using Oracle Coherence to move
their important workloads into the cloud. By certifying on industry
standards, such as Docker and Kubernetes, Coherence now runs in a cloud
neutral infrastructure. In addition, we've provided an open-source
Oracle Coherence Operator (the “operator”) which has several key
features to assist you with deploying and managing Coherence clusters in
a Kubernetes environment. You can:

* Run the Coherence you know and love on the industry standard
  Kubernetes container orchestration framework, using Docker containers
  for the system elements of Coherence.

* Use popular industry standard tools such as
  [Grafana](https://grafana.com/),
  [EFK](https://www.digitalocean.com/community/tutorials/how-to-set-up-an-elasticsearch-fluentd-and-kibana-efk-logging-stack-on-kubernetes), and
  [Prometheus](https://prometheus.io/) to monitor the performance,
  logs and health from your clusters.

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

The fastest way to experience the operator is to follow the [Quick Start guide](docs/quickstart.md), or try out the
[samples](docs/samples/README.md).

# About this documentation

This documentation includes sections targeted to different audiences.
To help you find what you are looking for more easily, please consult
this table of contents:

* The [Quick Start guide](docs/quickstart.md) explains how to
  quickly get Coherence running on Kubernetes, using the defaults,
  nothing special.

* The [User guide](docs/user-guide.md) contains detailed usage
  information on the Coherence Operator, including how to install and
  configure the operator and several common use-cases.

* The [Samples](docs/samples/README.md) provide detailed example
  code and instructions that show you how to perform various tasks
  related to the operator.

* The [Developer guide](docs/developer.md) provides details for people
  who want to understand how the operator is built, tested, and so
  on. Those who wish to contribute to the operator code will find useful
  information here.
  
<!--
* The [Contributing](#contributing-to-the-operator) section provides information about contribution requirements.
-->

# User guide

The [User guide](docs/user-guide.md) provides detailed information
about all aspects of using the operator including:

* Installing and configuring the operator.

* Using the operator to create and manage Coherence clusters.

* Manually creating Coherence clusters to be managed by the operator.

* Configuring Elasticsearch and Kibana to access the operator's log files.

* Shutting down clusters.

* And much more!

# Samples

Please refer to our [samples](docs/samples/README.md) for
information about the available sample code.

Need more help? Have a suggestion? Come and say "Hello!"

# Things to Keep In Mind for Existing Coherence Users

* Software running in Kubernetes must provide "health checks" so that
  Kubernetes can make informed decisions about starting, stopping, or
  even killing, the containers running the software.  The operator
  provides everything required to do this for Coherence.  Keep in mind
  that these health checks cause frequent `MemberJoined` and
  `MemberLeft` events to happen.  If these events refer to something
  like `OracleCoherenceK8sPodChecker`, they are normal and be safely
  ignored.

