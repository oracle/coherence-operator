# Manage the Reporter

Oracle Coherence reports show key management information over time. The reports identify trends that are valuable for troubleshooting and planning. Reporting is disabled by default and must be explicitly enabled by setting enabled in an operational override file or by using system properties.  This approach is valid for all versions of Coherence.

This sample shows how to access and manage the reporter over REST, which is available only in Oraccle Coherence 12.2.1.4.0 or later version.

> **Note**: To enable the Coherence Reporter using system properties, refer to the section
> [Tune JVM Runtime Settings](../../jvmarguments/) to use `store.javaOpts`.

The [Oracle Reporter documentation](
https://docs.oracle.com/middleware/1221/coherence/manage/reporter.htm#COHMG4885) explains how to set system propertites.

[Return to Management over REST samples](../)  [Return to Management samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Prerequisites

Ensure you have installed the Coherence Operator using the instructions [here](../../README.md#install-the-coherence-operator).

## Installation Steps

1. Install the Coherence cluster

   Install the cluster with 3 nodes:

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
   
   > *Note:* If your version of the Coherence Operator does not default to using Oracle Coherence 12.2.1.4.0, then you need to replace `your-12.2.1.4.0-Coherence-image` with an appropriate 12.2.1.4.0 image.
   
   Use `kubectl get pods -n sample-coherence-ns` and wait until all the three pods are running.
   
1. Port forward the Management over REST port:

   ```bash
   $ kubectl port-forward storage-coherence-0 -n sample-coherence-ns 30000:30000
   ```
   ```console
   Forwarding from [::1]:30000 -> 30000
   Forwarding from 127.0.0.1:30000 -> 30000
   ```   
   
1. Access the Reporter endpoint

   The Reporter is available at the following endpoint: `http://127.0.0.1:30000/management/coherence/cluster/reporters`
    
   The reporter needs to be started only on one node, but is available to be started on all nodes.
    
   View the reporter state on node 1:
    
    ```bash
    $ curl http://127.0.0.1:30000/management/coherence/cluster/reporters/1 2> /dev/null| json_pp
    ```   
    
    ```json
    {
    "outputPath" : "/.",
    "type" : "Reporter",
    "runMaxMillis" : 0,
    "runAverageMillis" : 0,
    "intervalSeconds" : 60,
    "runLastMillis" : 0,
    "state" : "Stopped",
    "refreshTime" : "2019-04-26T09:02:10.146Z",
    "links" : [
      {
         "href" : "http://127.0.0.1:30000/management/coherence/cluster/reporters",
         "rel" : "parent"
      },
      {
         "href" : "http://127.0.0.1:30000/management/coherence/cluster/reporters/1",
         "rel" : "self"
      },
      {
         "rel" : "canonical",
         "href" : "http://127.0.0.1:30000/management/coherence/cluster/reporters/1"
      }
    ],
    "nodeId" : "1",
    "reports" : [
       "reports/report-node.xml",
       "reports/report-network-health.xml",
       "reports/report-network-health-detail.xml",
       "reports/report-memory-status.xml",
       "reports/report-service.xml",
       "reports/report-cache-effectiveness.xml",
       "reports/report-proxy.xml",
       "reports/report-proxy-http.xml",
       "reports/report-management.xml",
       "reports/report-flashjournal.xml",
       "reports/report-ramjournal.xml",
       "reports/report-persistence.xml",
       "reports/report-persistence-detail.xml",
       "reports/report-federation-destination.xml",
       "reports/report-federation-origin.xml"
    ],
    "autoStart" : false,
    "lastReport" : null,
    "lastExecuteTime" : "1970-01-01T00:00:00.000Z",
    "currentBatch" : 0,
    "configFile" : "reports/report-group.xml"
    }
    ```
    
1.  Set the reporter output directory:

    ```bash
    $ curl -v -X POST -H 'Content-type: application/json' \
        http://127.0.0.1:30000/management/coherence/cluster/reporters/1 -d '{"outputPath": "/tmp/"}'
    ```   
    
    Validate the set output path:
    
    ```bash
    $ curl http://127.0.0.1:30000/management/coherence/cluster/reporters/1?fields=outputPath \
       | json_pp | grep outputPath

    "outputPath" : "/tmp"
    ```
    
1. Start the reporter:

   ```bash
   $ curl -X POST http://127.0.0.1:30000/management/coherence/cluster/reporters/1/start
   ```   
   
   Confim that the reporter has started:
   
   ```bash
   $ curl http://127.0.0.1:30000/management/coherence/cluster/reporters/1?fields=currentBatch,lastReport,lastExecuteTime,state,runLastMillis |json_pp
   ```
   
   The `State` must be either `Started` or `Sleeping`. `Sleeping` means that the reporter has run reports and is sleeping until the next execution.
   
   ```json
   {
   "lastExecuteTime" : "2019-04-26T09:11:53.172Z",
   "currentBatch" : 2,
   "lastReport" : "reports/report-federation-origin.xml",
   "runLastMillis" : 66,
   "links" : [
      {
         "rel" : "parent",
         "href" : "http://127.0.0.1:30000/management/coherence/cluster/reporters"
      },
      {
         "rel" : "self",
         "href" : "http://127.0.0.1:30000/management/coherence/cluster/reporters/1"
      },
      {
         "href" : "http://127.0.0.1:30000/management/coherence/cluster/reporters/1",
         "rel" : "canonical"
      }
   ],
   "state" : "Sleeping"
   }
   ```
   
1. View the reporter files

   Execute the following command to `exec` into the pod and view the Reporter files:

   ```bash
   $ kubectl exec -it -n sample-coherence-ns storage-coherence-0 bash

   $ ls -l /tmp/*.txt
   -rw-r--r-- 1 root root  618 Apr 26 09:15 /tmp/2019042609-Management.txt
   -rw-r--r-- 1 root root 1653 Apr 26 09:15 /tmp/2019042609-memory-status.txt
   -rw-r--r-- 1 root root 1089 Apr 26 09:15 /tmp/2019042609-network-health-detail.txt
   -rw-r--r-- 1 root root  395 Apr 26 09:15 /tmp/2019042609-network-health.txt
   -rw-r--r-- 1 root root 1377 Apr 26 09:15 /tmp/2019042609-nodes.txt
   -rw-r--r-- 1 root root  711 Apr 26 09:15 /tmp/2019042609-persistence-detail.txt
   -rw-r--r-- 1 root root  559 Apr 26 09:15 /tmp/2019042609-persistence.txt
   -rw-r--r-- 1 root root 2472 Apr 26 09:15 /tmp/2019042609-report-proxy-http.txt
   -rw-r--r-- 1 root root  798 Apr 26 09:15 /tmp/2019042609-report-proxy.txt
   -rw-r--r-- 1 root root 3366 Apr 26 09:15 /tmp/2019042609-service.txt
   ```  

   To copy the files to your current directory, use the following command:
   
   ```bash
   $ kubectl exec -it -n sample-coherence-ns storage-coherence-0 -- bash -c 'cd /tmp && tar cf /tmp/reports.tar  *.txt'
   $ kubectl cp sample-coherence-ns/storage-coherence-0:/tmp/reports.tar reports.tar
   ```
   
## Uninstall the Charts

Use the following commands to delete the chart installed in this sample:

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are removed from previous sample.

If you want to remove the `coherence-operator`, then include it in the `helm delete` command.
