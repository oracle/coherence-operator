# Storage-disabled client in cluster via interceptor

The steps detailed in samples such as [Storage-disabled client in cluster via interceptor](../../coherence-deployments/storage-disabled/interceptor)
call for the creation of a sidecar docker image that conveys the application classes 
to Kubernetes. This docker image is tagged with a version number, and the version number 
is how Kubernetes enables safe rolling upgrades. You can read more about safe rolling upgrades 
in the [Helm documentation](https://helm.sh/docs/helm/#helm-upgrade). 

Briefly, as with the scaling described in the preceding section, the safe rolling upgrade 
feature allows you to instruct Kubernetes, via the operator, to replace the currently deployed 
version of your application classes with a different one. Kubernetes does not care if the different 
version is "newer" or "older", as long as it has a docker tag and can be pulled by the cluster, 
that is all Kubernetes needs to know. The operator will ensure this is done without data loss or interruption of service.

In this sample we will initially deploy version 1.0.0 of our sidecar and then do a rolling upgrade to
version 2.0.0 of the sidecar which introduces a server side Interceptor to modify 
data to ensure its stored as uppercase. 

[Return to Storage-Disabled clients samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/docker-v1/Dockerfile](src/main/docker-v1/Dockerfile) - Dockerfile for creating side-car which only includes v1 storage-cache-config.xml

* [src/main/docker-v2/Dockerfile](src/main/docker-v2/Dockerfile) - Dockerfile for creating side-car which includes v2 storage-cache-config.xml

* [src/main/resources/conf/v1/storage-cache-config.xml](src/main/resources/conf/v1/storage-cache-config.xml) - version 1 cache config for storage-enabled tier - no interceptor 

* [src/main/resources/conf/v2/storage-cache-config.xml](src/main/resources/conf/v2/storage-cache-config.xml) - version 2 cache config for storage-enabled tier - includes uppercase interceptor  

* [src/main/resources/client-cache-config.xml](src/main/resources/client-cache-config.xml) - Client config for extend client

* [src/main/java/com/oracle/coherence/examples/UppsercaseInterceptor.java](src/main/java/com/oracle/coherence/examples/UppsercaseInterceptor.java) - interceptor that will change all entries to uppercase - version 2.

Note if you wish to enable Prometheus or log capture, change the following in the helm installs to `true`. Their default values are false, but they are set to `false` in the samples below for completeness.

* Prometheus: `--set prometheusoperator.enabled=true`

* Log capture: `--set logCaptureEnabled=true`

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](here).

## Installation Steps

1. Change to the `samples/operator/rolling-upgrade` directory and ensure you have your you have your maven build     
   environment set for JDK11 and build the project.

   ```bash
   mvn clean install -P docker-v1
   ```

   The above will build the v1 Docker image called `rolling-upgrade-sample:1.0.0`. This will be the initial image deployed 
   to the storage server.
   
   ```bash
   mvn clean install -P docker-v2
   ```

   The above will build the v2 Docker image called `rolling-upgrade-sample:2.0.0`. This will be the image we upgrade
   the deployment with.

   **Note:** If you are running against a remote Kubernetes cluster you will need to
   push the above image to your repository accessible to that cluster.

1. Install the Coherence cluster with rolling-upgrade-sample:1.0.0 image

   We are also setting the cluster-name to `interceptor-cluster` which will we use in later steps.

   ```bash
   helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=interceptor-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set store.cacheConfig=storage-cache-config.xml \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --set userArtifacts.image=rolling-upgrade-sample:1.0.0 \
      --version 1.0.0-SNAPSHOT coherence-community/coherence
   ```

   Because we use stateful sets, the coherence cluster will start one pod at a time.
   
   You can change this by using `--set store.podManagementPolicy=Parallel` in the above command.
    
   Use `kubectl get pods -n sample-coherence-ns` to ensure that all pods are running.
   All 3 storage-coherence-0/1/2 pods should be running and ready, as below:

   ```bash
   NAME                                                     READY   STATUS    RESTARTS   AGE
   storage-coherence-0                                      1/1     Running   0          4m
   storage-coherence-1                                      1/1     Running   0          2m
   storage-coherence-2                                      1/1     Running   0          1m
   ```
   
1. Port forward the proxy port on the storage-coherence-0 pod.

   ```bash
   $ kubectl port-forward -n sample-coherence-ns storage-coherence-0 20000:20000
   ```

1. Connect via QueryPlus and issue CohQL commands

   Issue the following command to run QueryPlus:

   ```bash
   $ mvn exec:java
   ```

   Run the following CohQL commands to insert data into the cluster.

   ```sql
   CohQL> insert into 'test' key('key-1') value('value-1');
   CohQL> insert into 'test' key('key-2') value('value-2');

   CohQL> select key(), value() from 'test';
   Results
   ["key-1", "value-1"]   
   ["key-2", "value-2"]
   ```
   
   Keep CohQL open while you continue.
   
1. Upgrade the server to rolling-upgrade-sample:2.0.0 image  

   Issue the following to upgrade to version 2.0.0 of the image.
   
   * `--reuse-values` - specifies to re-use all previous values associated with the release
   
   * `--set userArtifacts.image=rolling-upgrade-sample:2.0.0` - the new artifact version

   ```bash
   helm upgrade storage coherence-community/coherence \
      --namespace sample-coherence-ns \
      --reuse-values \
      --set imagePullSecrets=sample-coherence-secret \
      --version 1.0.0-SNAPSHOT \
      --set userArtifacts.image=rolling-upgrade-sample:2.0.0   
     
   ```
   
   While the upgrade is running you can issue the following in your `CohQL` session:
   
   ```sql
   select key(), value() from 'test';
   ```
   
   You will notice that the data always remains the same.
   
   *Note*: Your port-forward will fail once the `storage-coherence-0` pod restarts, so you will have 
   stop and restart it.  In an environment where you have configured a load balancer, then the 
   Coherence*Extend session will automatically reconnect for you.
   
1. Check the status of the upgrade

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   NAME                                  READY   STATUS        RESTARTS   AGE
   coherence-operator-66f9bb7b75-pprgg   1/1     Running       0          30m
   storage-coherence-0                   1/1     Running       0          19m
   storage-coherence-1                   0/1     Terminating   0          18m
   storage-coherence-2                   1/1     Running       0          1m 
   ```
   
   When all of the pods have status of `Running` and `1/1` for Ready, you can continue.
   
   *Note*: The above shows not all pods are finished restarting.
   
1. Add more data via CohQL commands

   ```sql
   CohQL> insert into 'test' key('key-3') value('value-3');

   CohQL> select key(), value() from 'test';
   Results
   ["key-1", "value-1"]
   ["key-3", "VALUE-3"]
   ["key-2", "value-2"]
   ```    

   You will notice that the value for `key-3` has been converted to uppercase which shows that the
   server-side interceptor is now active.
 
1. Verify the 2.0.0 image is active

   Use the following to verify the 2.0.0 image is active:
   
   ```bash
   $ kubectl describe pod storage-coherence-0  -n sample-coherence-ns | grep rolling-upgrade
   Image:         rolling-upgrade-sample:2.0.0
   Normal  Pulled                 4m59s  kubelet, docker-for-desktop  Container image "rolling-upgrade-sample:2.0.0" already present on machine
   ```
   
   The above shows that the version 2.0.0 image is now present.
   
## Verifying Grafana Data (If you enabled Prometheus)

Access Grafana using the instructions [here](../../../README.md#access-grafana).

## Verifying Kibana Logs (if you enabled log capture)

Access Kibana using the instructions [here](../../../README.md#access-kibana).

## Uninstalling the Chart

Carry out the following commands to delete the two charts created in this sample.

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are gone from previous samples.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.
