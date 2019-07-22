# Produce and Extract Heap Dump

Some of the debugging techniques described in [Debugging in
Coherence](https://docs.oracle.com/middleware/12213/coherence/develop-applications/debugging-coherence.htm)
require the creation of files, such as  log files
and JVM heap dumps, for analysis. You can also create and extract these files in the Coherence Operator.  

This sample shows how to collect a `.hprof` file for a
heap dump. A single-command technique is also included at the end of this sample.

[Return to Diagnostics Tools](../) / [Return to Management samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Prerequisites

Ensure you have installed the Coherence Operator using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Install the Coherence cluster

   Install a Coherence cluster if you don't have one running:

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=coherence-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      coherence/coherence
   ```
   
1. Ensure the pods are running:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   ```console
   NAME                                  READY   STATUS    RESTARTS   AGE
   coherence-operator-66f9bb7b75-hqk4l   1/1     Running   0          13m
   storage-coherence-0                   1/1     Running   0          3m
   storage-coherence-1                   1/1     Running   0          2m
   storage-coherence-2                   1/1     Running   0          44s
   ```
   
1. Open a terminal window in one of the storage nodes:

   ```bash
   $ kubectl exec -it storage-coherence-0 -n sample-coherence-ns -- bash
   ```

   Obtain the PID of the Coherence process. Generally, the PID is `1`. You can also use `jps` to get the actual PID.

   ```bash
   # /usr/java/default/bin/jps
   1 DefaultCacheServer
   4230 Jps
   ```

1. Use the `jcmd` command to extract the heap dump:

   ```bash
   $ rm /tmp/heap.hprof
   $ /usr/java/default/bin/jcmd 1 GC.heap_dump /tmp/heap.hprof
   $ exit
   ```
   
1. Copy the heap dump to local machine:

   ```bash
   $ kubectl cp sample-coherence-ns/storage-coherence-0:/tmp/heap.hprof heap.hprof 
   ```  
   
   Deepending upon whether the Kubernetes cluster is local or remote, this might take some time.
   
1. Single command usage

   Assuming that the Coherence PID is `1`, you can use this repeatable single-command technique to extract the heap dump:

   ```bash
   $ (kubectl exec storage-coherence-0 -n sample-coherence-ns  -- /bin/bash -c "rm -f /tmp/heap.hprof; /usr/java/default/bin/jcmd 1 GC.heap_dump /tmp/heap.hprof; cat /tmp/heap.hprof > /dev/stderr" ) 2> heap.hprof
   ```
    Note that we redirect the heap dump output to `stderr` to prevent the unsuppressable.

   ```bash
   1:
   Heap dump file created
   ```
  
## Uninstall the Charts

Use the following command to delete the chart installed in this sample:

```bash
$ helm delete storage --purge  
```

Before starting another sample, ensure that all the pods are removed from previous sample.

If you want to remove the `coherence-operator`, then include it in the `helm delete` command.
