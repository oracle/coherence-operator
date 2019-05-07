# Accessing the EFK stack for viewing logs

The Oracle Coherence Kubernetes Operator manages Oracle Coherence on Kubernetes.
It manages monitoring data through Prometheus, and logging data through the EFK
(ElasticSearch, Fluentd and Kibana) stack.

This use-case is covered [in the samples](docs/samples/operator/logging/log-capture/).

To access and configure the Kibana UI, please follow the instructions below.

## Installing the charts

When you install the `coherence-operator` and `coherence` charts, you must specify the following
option for `helm` , for both charts, to ensure the EFK stack (Elasitcsearch, Fluentd and Kibana) 
is installed and correctly configured.

```bash
--set logCaptureEnabled=true 
```

## Port forward Kibana

Once you have installed both charts, use the following script to port forward the Kibana port 5601.

Note: If your chart is installed in a namespace other than `default`
then include `--namespace` option for both `kubectl` commands.

```bash
#!/bin/bash
trap "exit" INT
  
while :
do
   kubectl port-forward $(kubectl get pods | grep kibana | awk '{print $1}') 5601:5601
done

```

## Default dashboards

There are a number of dashboard created via the import process.

**Coherence Operator**

* *Coherence Operator - All Messages* - Shows all Coherence Operator messages

**Coherence Cluster**

* *Coherence Cluster - All Messages* - Shows all messages

* *Coherence Cluster - Errors and Warnings* - Shows only errors and warnings

* *Coherence Cluster - Persistence* - Shows persistence related messages 

* *Coherence Cluster - Partitions* - Shows Ppartition related messages 

* *Coherence Cluster - Message Sources* - Allows visualization of messages via the message source (Thread)

* *Coherence Cluster - Configuration Messages* - Shows configuration related messages

* *Coherence Cluster - Network* - Shows network related messages such as communication delays and TCP ring disconnects 

## Default queries

There are many queries related to common Coherence messages, warnings and errors that are 
loaded and can be accessed via the discover side-bar.

## Troubleshooting

### No default index pattern

There are two index patterns created via the import process and the `coherence-cluster-*` pattern
will be set as the default. If for some reason this is not the case, then carry out the following
to set `coherence-cluster-*` as the default:

* Open a the following URL `http://127.0.0.1:5601/`

* Click on `Management` side-bar.

* Click on `Index Patterns`, then select `coherence-cluster-*' pattern.

* Click the `Star` on the top right to set as default.
