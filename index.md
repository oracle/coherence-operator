<!--
Copyright 2018, 2019, Oracle Corporation and/or its affiliates.
All rights reserved.  Licensed under the Universal
Permissive License v 1.0 as shown at
http://oss.oracle.com/licenses/upl.

-->

# Coherence Operator

Oracle enables organizations using [Coherence](https://www.oracle.com/middleware/technologies/coherence.html) to move their 
clusters into the cloud. By supporting industry standards, such as Docker and Kubernetes, Oracle facilitates running 
Coherence on cloud-neutral infrastructure. In addition, Oracle provides an open-source Coherence Operator ("the operator"), 
which implements features to assist with deploying and managing Coherence clusters in a Kubernetes environment. You can:

The current version of the [Coherence Operator is 2.0.0](https://oracle.github.io/coherence-operator/docs/2.0.0)

## Features

* The Coherence Operator is a Go based application built using the 
[Operator SDK](https://github.com/operator-framework/operator-sdk)

* Deploy Coherence clusters into Kubernetes using a custom resource definition

* Flexibly override and customize cluster configuration.

* Safely scale the Coherence deployments using standard Kubernetes tools such as the `kubectl scale` command.


* Use Kubernetes Zone information to ensure data stored in Coherence is resilient to loss of a Zone.
  Coherence goes to great efforts to ensure data is safe across processes, machines, racks and sites. 
  When Coherence is deployed to Kubernetes with the Coherence Operator, data will be spread across zones
  to ensure this underlying principle is supported; thus by default, loss of any zone is a tolerated failure mode. 
  This is reflected in the StatusHA value (SITE-SAFE) for partitioned services, in addition to the member level 
  site information that is equivalent to the kubernetes zone label on the associated pod.

* Use Kubernetes persistent volumes when using Coherence's disk-based storage features Elastic Data or Persistence.

* Deploy custom code for your server side classes.

* Use
  [Coherence*Extend](https://docs.oracle.com/middleware/12213/coherence/develop-remote-clients/building-your-first-extend-application.htm#COHCG5033)
  to access your cluster with a variety of clients.

* Use popular industry standard tools such as
  [Grafana](https://grafana.com/),
  [ELK](https://www.elastic.co/elk-stack) (or more specifically the EFK stack including Fluentd), and
  [Prometheus](https://prometheus.io/)
  to monitor the performance, logs and and health of your clusters.


## Getting Started

The fastest way to experience the operator is to follow the 
[Quick Start guide](https://oracle.github.io/coherence-operator/docs/2.0.0/#/about/03_quickstart)


## Need more help? Have a suggestion? Come and say "Hello!"

### Join Our Public Slack Channel

We have a **public Slack channel** where you can get in touch with us to ask questions about using the operator or give us
feedback or suggestions about what features and improvements you would like to see.
We would love to hear from you. To join our channel, please 
[visit this site to get an invitation](https://join.slack.com/t/oraclecoherence/shared_invite/enQtNzcxNTQwMTAzNjE4LTJkZWI5ZDkzNGEzOTllZDgwZDU3NGM2YjY5YWYwMzM3ODdkNTU2NmNmNDFhOWIxMDZlNjg2MzE3NmMxZWMxMWE).
The invitation email will include details of how to access our Slack workspace. 
After you are logged in, please come to `#operator` and say, "hello!"


#### Previous Versions

Looking for a previous version's documentation?

[Version 1.0.0](https://oracle.github.io/coherence-operator/docs/1.0.0/)