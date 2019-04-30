# Produce and extract a Java Flight Recorder (JFR) file

Java Flight Recorder (JFR) is a tool for collecting diagnostic and profiling data 
about a running Java application. It is integrated into the Java Virtual Machine (JVM) 
and causes almost no performance overhead, so it can be used even in heavily loaded production environments.

In Coherence 12.2.1.4.0 and above, the [Management over REST](../../rest) functionality provides 
the ability to create and managed JFR recordings.

In this sample, we will execute an operation across all nodes of a cluster.

Valid JFR commands are jfrStart, jfrStop, jfrDump, and jfrCheck.
     
See [documentation](https://docs.oracle.com/javacomponents/jmc-5-4/jfr-runtime-guide/run.htm#JFRUH176) for more details on commands.  
      

[Return to Diagnostics Tools](../) / [Return to Management samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Install the Coherence cluster

   Issue the following to install a Coherence cluster if you don't already have one running.

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=coherence-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --version 1.0.0-SNAPSHOT coherence-community/coherence
   ```
   
1. Ensure the pods are running:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-66f9bb7b75-hqk4l   1/1     Running   0          13m
   storage-coherence-0                   1/1     Running   0          3m
   storage-coherence-1                   1/1     Running   0          2m
   storage-coherence-2                   1/1     Running   0          44s
   ```
   
1. Port-Forward the Management over REST port

   ```bash
   $ kubectl port-forward storage-coherence-0 -n sample-coherence-ns 30000:30000
   Forwarding from [::1]:30000 -> 30000
   Forwarding from 127.0.0.1:30000 -> 30000
   ```   
   
   The base URL for JFR commands is `http://127.0.0.1:30000/management/coherence/cluster/diagnostic-cmd/{command}`.
   
   Valid JFR commands are jfrStart, jfrStop, jfrDump, and jfrCheck.
   
1. Start Recording Across all nodes

   Issue the following to start a JFR recording across all nodes for a period of 60 seconds.
   
   ```bash
   $ curl http://127.0.0.1:30000/management/coherence/cluster/diagnostic-cmd/jfrStart?name=myJfr,duration=30s,filename=/tmp/myRecording.jfr
   ```

   **TBC**
  
## Uninstalling the Charts

Carry out the following commands to delete the chart installed in this sample.

```bash
$ helm delete storage --purge  
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.
