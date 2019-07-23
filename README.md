<!--
Copyright 2019, Oracle Corporation and/or its affiliates.
All rights reserved.  Licensed under the Universal
Permissive License v 1.0 as shown at
http://oss.oracle.com/licenses/upl.

-->

-----

# Oracle Coherence Operator

Oracle enables organizations using Coherence to move their clusters into the cloud. By supporting industry standards, such as Docker and Kubernetes, Oracle facilitates running Coherence on cloud-neutral infrastructure. In addition, Oracle provides an open-source Coherence Operator ("the operator"), which implements features to assist with deploying and managing Coherence clusters in a Kubernetes environment. You can:

* Run Coherence within the industry standard Kubernetes container orchestration framework, using Docker containers for the members of a Coherence cluster.
* Use popular industry standard tools such as Grafana, ELK (or more specifically the Elasticsearch, Fluentd and Kibana (EFK) stack), and Prometheus to monitor the performance, logs, and health of your clusters.
* Flexibly override and customize cluster configuration.
* Scale the Coherence deployment.
* Use Coherence*Extend to access your cluster with a variety of clients.
* Use Kubernetes Zone information to ensure data stored in Coherence is resilient to loss of a Zone. Coherence goes to great efforts to ensure data is safe across processes, machines, racks and sites. When Coherence is deployed to Kubernetes with the Coherence Operator, data will be spread across zones to ensure this underlying principle is supported; thus by default, loss of any zone is a tolerated failure mode. This is reflected in the StatusHA value (SITE-SAFE) for partitioned services, in addition to the member level site information that is equivalent to the kubernetes zone label on the associated pod.
* Start clusters based on declarative startup parameters and desired states.
* Use Kubernetes persistent volumes when using Coherence’s disk-based storage features Elastic Data or Persistence.
* Deploy custom code for your EntryProcessor classes and other server-side Coherence constructs.

The fastest way to experience the operator is to follow the [Quick Start guide](https://oracle.github.io/coherence-operator/docs/quickstart.html), or you can look through our  [documentation](https://oracle.github.io/coherence-operator/), or try out the [samples](https://oracle.github.io/coherence-operator/docs/samples/).

-------
The current release of the operator is 1.0. This release was published on .

-------

# Need more help? Have a suggestion? Come and say "Hello!"

We have a **public Slack channel** where you can get in touch with us to ask questions about using the operator or give us feedback or suggestions about what features and improvements you would like to see. We would love to hear from you. To join our channel, please [visit this site to get an invitation](https://join.slack.com/t/oraclecoherence/shared_invite/enQtNjA3MTU3MTk0MTE3LWZhMTdhM2E0ZDY2Y2FmZDhiOThlYzJjYTc5NzdkYWVlMzUzODZiNTI4ZWU3ZTlmNDQ4MmE1OTRhOWI1MmIxZjQ).  The
invitation email will include details of how to access our Slack
workspace.  After you are logged in, please come to `#operator` and say, "hello!"



# Documentation

Documentation for the operator is available [here](https://oracle.github.io/coherence-operator/) and includes information for users and for developers. It provides [Samples](https://oracle.github.io/coherence-operator/docs/samples/), [User Guide](https://oracle.github.io/coherence-operator/docs/user-guide.html), [Developer Guide](https://oracle.github.io/coherence-operator/docs/developer.html) and a [Quick Start guide](https://oracle.github.io/coherence-operator/docs/quickstart.html) if you just want to get up and running quickly.
