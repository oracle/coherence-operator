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

The fastest way to experience the operator is to follow the
[Quick Start guide](docs/quickstart.md), or you can read our
[blogs](https://blogs.oracle.com/weblogicserver/how-to-weblogic-server-on-kubernetes),
or try out the [samples](docs/samples/README.md).

This documentation is for the current release of the operator.  For
documentation on previous releases, use the GitHub `Branch` pulldown,
select the `Tags` tab, and click on the desired release.

## Known issues

| Issue | Description |
|-------|-------------|
| TODO | MVP Docs: a collection of hyperlinked markdown files |

<!--
Operator version 0.1.0
Documentation for the 0.1.0 release of the operator is
available [here](docs/0.1.0/README.md).

Backward compatibility guidelines
PENDING
-->

# About this documentation

This documentation includes sections targeted to different audiences.
To help you find what you are looking for more easily, please consult
this table of contents:

* The [Quick Start guide](docs/quickstart.md) explains how to
  quickly get the Coherence running on Kubernetes, using the defaults,
  nothing special.

* The [User guide](docs/user-guide.md) contains detailed usage
  information on the Coherence Operator, including how to install and
  configure the operator, and how to use it to create and manage
  WebLogic domains.

* The [Samples](docs/samples/README.md) provide detailed example
  code and instructions that show you how to perform various tasks
  related to the operator.

<!--
* The [Developer guide](docs/developer.md) provides details for people
  who want to understand how the operator is built, tested, and so
  on. Those who wish to contribute to the operator code will find useful
  information here.  This section also includes API documentation
  (Javadoc) and Swagger/OpenAPI documentation for the REST APIs.

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

<!--
Need more help? Have a suggestion? Come and say "Hello!"

We have a **public Slack channel** where you can get in
touch with us to ask questions about using the operator or
give us feedback or suggestions about what features and
improvements you would like to see.  We would love to hear
from you. To join our channel, please
[visit this site to get an invitation](https://weblogic-slack-inviter.herokuapp.com/).
The invitation email will include details of how to access
our Slack workspace.  After you are logged in, please come
to `#operator` and say, "hello!"

-->
