<!--
Copyright 2019, 2022, Oracle Corporation and/or its affiliates.
All rights reserved.  Licensed under the Universal
Permissive License v 1.0 as shown at
http://oss.oracle.com/licenses/upl.

-->

-----
![logo](docs/images/logo-with-name.png)

![Operator CI](https://github.com/oracle/coherence-operator/workflows/Operator%20CI/badge.svg?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/oracle/coherence-operator)](https://goreportcard.com/report/github.com/oracle/coherence-operator)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=oracle_coherence-operator&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=oracle_coherence-operator)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=oracle_coherence-operator&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=oracle_coherence-operator)

![GitHub release (latest by date)](https://img.shields.io/github/v/release/oracle/coherence-operator)
[![License](http://img.shields.io/badge/license-UPL%201.0-blue.svg)](https://oss.oracle.com/licenses/upl/)

# Coherence Operator

Oracle enables organizations using 
[Coherence](https://oracle.github.io/coherence) 
to move their clusters into the cloud, by supporting industry standards, 
such as Docker and Kubernetes, Oracle facilitates running Coherence on cloud-neutral infrastructure. 
In addition, Oracle provides an open-source Coherence Operator, which implements features to assist with 
deploying and managing Coherence clusters in a Kubernetes environment. You can:

* Run Coherence within the industry standard Kubernetes container orchestration framework, using Docker 
containers for the members of a Coherence cluster.
* Flexibly override and customize cluster configuration using a `Coherence` custom resource definition.
* Safely scale the roles of a Coherence cluster using Kubernetes verbs or updates.
* Expose ports (e.g. Coherence*Extend) to access your cluster with a variety of clients.
* Deploy custom code for your server side classes.
* Use Kubernetes persistent volumes when using Coherence’s disk-based storage features Elastic Data or Persistence.
* Use Kubernetes Zone information to ensure data stored in Coherence is resilient to loss of a Zone. Coherence goes 
to great efforts to ensure data is safe across processes, machines, racks and sites. When Coherence is deployed to 
Kubernetes with the Coherence Operator, data will be spread across zones to ensure this underlying principle is 
supported; thus by default, loss of any zone is a tolerated failure mode. This is reflected in the StatusHA 
value (SITE-SAFE) for partitioned services, in addition to the member level site information that is equivalent 
to the kubernetes zone label on the associated pod.
* Use popular industry standard tools such as Grafana, ELK (or more specifically the Elasticsearch, Fluentd and 
Kibana (EFK) stack), and Prometheus to monitor the performance, logs, and health of your clusters.

-------
## Documentation

Documentation for the Coherence Operator is available [here](https://docs.coherence.community/coherence-operator/docs/latest)

The fastest way to experience the operator is to follow the 
[Quick Start guide](https://docs.coherence.community/coherence-operator/docs/latest/docs/about/03_quickstart).
-------

# Need more help? Have a suggestion? Come and say "Hello!"

We have a **public Slack channel** where you can get in touch with us to ask questions about using the operator 
or give us feedback or suggestions about what features and improvements you would like to see. We would love 
to hear from you. To join our channel, 
please [visit this site to get an invitation](https://join.slack.com/t/oraclecoherence/shared_invite/enQtNzcxNTQwMTAzNjE4LTJkZWI5ZDkzNGEzOTllZDgwZDU3NGM2YjY5YWYwMzM3ODdkNTU2NmNmNDFhOWIxMDZlNjg2MzE3NmMxZWMxMWE).  
The invitation email will include details of how to access our Slack
workspace.  After you are logged in, please come to `#operator` and say, "hello!"

## Contributing

This project welcomes contributions from the community. Before submitting a pull request, please [review our contribution guide](./CONTRIBUTING.md)

## Security

Please consult the [security guide](./SECURITY.md) for our responsible security vulnerability disclosure process

## License

Copyright (c) 2019, 2025 Oracle and/or its affiliates.

*Replace this statement if your project is not licensed under the UPL*

Released under the Universal Permissive License v1.0 as shown at
<https://oss.oracle.com/licenses/upl/>.
