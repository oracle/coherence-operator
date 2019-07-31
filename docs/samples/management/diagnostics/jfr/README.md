# Produce and Extract a Java Flight Recorder (JFR) File

Java Flight Recorder (JFR) is a tool for collecting diagnostic and profiling data about a running Java application. It is integrated into the Java Virtual Machine (JVM) 
and does cause any performance overhead and it can be used in heavily loaded production environments.

By default, when the Coherence chart is installed, the Management over REST endpoint is exposed at port 30000 on each of the pods.

This sample shows how you can create and manage Flight Recordings using the Management over REST endpoint, which is exposed via the following endpoint:

* `http://host:30000/management/coherence/cluster/diagnostic-cmd`

The Swagger document is available at:  

* `http://host:30000/management/coherence/cluster/metadata-catalog`

The endpoint makes use of `jcmd` which is described in the [Oracle documentation](https://docs.oracle.com/javacomponents/jmc-5-4/jfr-runtime-guide/comline.htm).

> **Note**: Use of Management over REST is available only when using the operator with Oracle Coherence 12.2.1.4.0.

[Return to Diagnostics Tools](../) / [Return to Management samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Prerequisites

Ensure you have installed the Coherence Operator using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Install the Coherence cluster

   Use the following command to install the cluster:

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=coherence-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set logCaptureEnabled=false \
      --set coherence.image=your-12.2.1.4.0-Coherence-image \
      coherence/coherence
   ```
   
   > *Note:* If your version of the Coherence Operator does not default to using Coherence 12.2.1.4.0, then you need to replace `your-12.2.1.4.0-Coherence-image` with an appropriate 12.2.1.4.0 image.
   
   Ensure the pods are running:
   
   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   ```console
   NAME                   READY   STATUS    RESTARTS   AGE
   storage-coherence-0    1/1     Running   0          4m
   storage-coherence-1    1/1     Running   0          2m   
   storage-coherence-2    1/1     Running   0          2m
   ```
   
1. Port forward the Management over REST port:

   ```bash
   $ kubectl port-forward storage-coherence-0 -n sample-coherence-ns 30000:30000
   ```
   ```console
   Forwarding from [::1]:30000 -> 30000
   Forwarding from 127.0.0.1:30000 -> 30000
   ```   
   
1. Start a Flight Recording

   Use curl to start a recording with a name `test1` on member 1:
   
   ```bash
   $ curl -X POST -H 'Content-type: application/json' -v \
       "http://127.0.0.1:30000/management/coherence/cluster/members/1/diagnostic-cmd/jfrStart?options=name%3Dtest1"
   ```
   
   It returns a HTTP 200 OK status.
   
   Check the status of the recording:
   
   ```bash
   $ curl -X POST -H 'Content-type: application/json' -v \
        "http://127.0.0.1:30000/management/coherence/cluster/members/1/diagnostic-cmd/jfrCheck?options=name%3Dtest1"
   ```
   
   The following status line is displayed in the output:
   
   ```json
   "status":"Recording: recording=1 name=\"test1\" (running)\n"}
   ```
   
1. Dump the Flight Recording to a file
       
   Use curl to dump the currently running recording with the name `test1` on member 1 to a file called `/tmp/test1.jfr`:
   
   ```bash
   $ curl -X POST -H 'Content-type: application/json' -v \
       "http://127.0.0.1:30000/management/coherence/cluster/members/1/diagnostic-cmd/jfrDump?options=name%3Dtest1,filename%3D/tmp/test1.jfr"
   ```
   
   The following status line is displayed in the output:
   
   ```json
   "status":"Dumped recording \"test1\", 717.0 kB written to:\n\n/tmp/test1.jfr\n"}
   ```
   
1. Copy the JFR recording to a local machine:

   ```bash
   $ kubectl cp sample-coherence-ns/storage-coherence-0:/tmp/test1.jfr test1.jfr
   ```  
   
   Depending upon whether your Kubernetes cluster is local or remote, this might take some time.
  
1. Stop the Flight Recording

   Use curl to stop the currently running recording with the name `test1` on member 1:
   
   ```bash
   $ curl -X POST -H 'Content-type: application/json' -v \
       "http://127.0.0.1:30000/management/coherence/cluster/members/1/diagnostic-cmd/jfrStop?options=name%3Dtest1"
   ```
   
   The following status line is displayed in the output:
   
   ```json
   "status":"Stopped recording \"test1\".\n"}
   ```

> *Note:* The commands in this procedure can be run on all nodes by leaving out `/members/1` in the path. For example,
> `http://127.0.0.1:30000/management/coherence/cluster/diagnostic-cmd/jfrStart`

## Uninstall the Charts

Use the following command to delete the chart installed in this sample:

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are removed from previous samples.

If you want to remove the `coherence-operator`, then include it in the `helm delete` command.
