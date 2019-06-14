# Produce and extract a Java Flight Recorder (JFR) file

Java Flight Recorder (JFR) is a tool for collecting diagnostic and profiling 
data about a running Java application. It is integrated into the Java Virtual Machine (JVM) 
and causes almost no performance overhead, so it can be used even in heavily loaded production environments.

By default when the Coherence chart is installed the Management over REST endpoint will be exposed
as port 30000 on each of the Pods. 

This sample shows how you can create and manage Flight Recordings by using the Management over REST endpoint
which is exposed via the following endpoint:

* `http://host:300000`/management/coherence/cluster/diagnostic-cmd`.

The Swagger document is available via the following URL:  

* `http://host:300000`/management/coherence/cluster/metadata-catalog`.

The endpoint makes use of `jcmd` which is described in the [Oracle documentation](https://docs.oracle.com/javacomponents/jmc-5-4/jfr-runtime-guide/comline.htm).

> **Note**: Use of Management over REST is only available when using the
> operator with Coherence 12.2.1.4.

[Return to Diagnostics Tools](../) / [Return to Management samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Install the Coherence cluster

   Issue the following to install the cluster:

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=coherence-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --set coherence.image=your-12.2.1.4.0-Coherence-image \
      coherence/coherence
   ```
   
   > *Note:* If your version of the Coherence Operator does not default to using Coherence
   > 12.2.1.4.0, then you will need to replace `your-12.2.1.4.0-Coherence-image` with an
   > appropriate 12.2.1.4.0 image.
   
   Once the install has completed, issue the following command to list the pods:
   
   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                   READY   STATUS    RESTARTS   AGE
   storage-coherence-0    1/1     Running   0          4m
   storage-coherence-1    1/1     Running   0          2m   
   storage-coherence-2    1/1     Running   0          2m
   ```
   
1. Port-Forward the Management over REST port

   ```bash
   $ kubectl port-forward storage-coherence-0 -n sample-coherence-ns 30000:30000
   Forwarding from [::1]:30000 -> 30000
   Forwarding from 127.0.0.1:30000 -> 30000
   ```   
   
1. Start a Flight Recording

   Using `curl`, issue the following command to start a recording with a name of `test1` on member 1.
   
   ```bash
   $ curl -X POST -H 'Content-type: application/json' -v \
       "http://127.0.0.1:30000/management/coherence/cluster/members/1/diagnostic-cmd/jfrStart?options=name%3Dtest1"
   ```
   
   The above should return a HTTP 200 OK status.
   
   Check the status of the recording by using:
   
   ```bash
   $ curl -X POST -H 'Content-type: application/json' -v \
        "http://127.0.0.1:30000/management/coherence/cluster/members/1/diagnostic-cmd/jfrCheck?options=name%3Dtest1"
   ```
   
   The following status line should be shown in the output:
   
   ```json
   "status":"Recording: recording=1 name=\"test1\" (running)\n"}
   ```
   
1. Dump the Flight Recording to a file
       
   Using `curl`, issue the following command to dump the currently running recording with a 
   name of `test1` on member 1 to a file called `/tmp/test1.jfr`.
   
   ```bash
   $ curl -X POST -H 'Content-type: application/json' -v \
       "http://127.0.0.1:30000/management/coherence/cluster/members/1/diagnostic-cmd/jfrDump?options=name%3Dtest1,filename%3D/tmp/test1.jfr"
   ```
   
   The following status line should be shown in the output:
   
   ```json
   "status":"Dumped recording \"test1\", 717.0 kB written to:\n\n/tmp/test1.jfr\n"}
   ```
   
1. Copy the JFR recording to local machine

   ```bash
   $ kubectl cp sample-coherence-ns/storage-coherence-0:/tmp/test1.jfr test1.jfr
   ```  
   
   This may take a while depending upon if you Kubernetes cluster is local or remote.
  
1. Stop the Flight Recording

   Using `curl`, issue the following command to stop the currently running recording with a 
   name of `test1` on member 1.
   
   ```bash
   $ curl -X POST -H 'Content-type: application/json' -v \
       "http://127.0.0.1:30000/management/coherence/cluster/members/1/diagnostic-cmd/jfrStop?options=name%3Dtest1"
   ```
   
   The following status line should be shown in the output:
   
   ```json
   "status":"Stopped recording \"test1\".\n"}
   ```

> *Note:* All of the above commands can be run against all nodes by leaving out `/members/1` path. E.g.
> `http://127.0.0.1:30000/management/coherence/cluster/diagnostic-cmd/jfrStart` 

## Uninstalling the Charts

Carry out the following commands to delete the chart installed in this sample.

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are gone from previous samples.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above. 
  

    



   
