# Accessing the EFK (Elasticsearch, Fluentd and Kibana) Stack for Viewing Logs

The Coherence Operator (the "operator") manages Oracle Coherence on Kubernetes.
It manages monitoring data through Prometheus, and logging data through the EFK stack.

This use case is covered in the samples. Refer to the [samples documentation](samples/operator/logging/log-capture/README.md). To access and configure the Kibana user interface(UI),  follow the instructions:

## Install the Charts

When you install the `coherence-operator` and `coherence` charts, you must specify the following
option for `helm` for both charts. This ensures that the EFK stack is installed and correctly configured.

```bash
--set logCaptureEnabled=true
```

## Port Forward Kibana

Once you have installed both charts, use the following script to port forward the Kibana port ```5601```:

**Note**: If your chart is installed in a namespace other than the `default`,
then include the `--namespace` option for both `kubectl` commands.

```bash
#!/bin/bash
trap "exit" INT

while :
do
   kubectl port-forward $(kubectl get pods | grep kibana | awk '{print $1}') 5601:5601
done

```
## Default Dashboards

There are a number of dashboards created via the import process.

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

## Default Queries

There are many queries related to common Coherence messages, warnings and errors that are
loaded and can be accessed via the discover side-bar.

## Troubleshooting

### No Default Index Pattern

There are two index patterns created via the import process and the `coherence-cluster-*` pattern
will be set as the default. If for some reason this is not the case, then perform the following steps
to set `coherence-cluster-*` as the default:

1. Open the URL: `http://127.0.0.1:5601/`

2. Navigate to the `Management` side-bar.

3. Click `Index Patterns`, then select the `coherence-cluster-*` pattern.

4. Click `Star` on the top right to set as default.
