# Deploy Elastic Data using default FlashJournal locations

In this sample we will enable Elastic Data using the default FlashJournal directory
which is `/tmp/`.  

Please see [Oracle Elastic Data Documentation](https://docs.oracle.com/middleware/12213/coherence/COHDG/implementing-storage-and-backing-maps.htm#COHDG5496) 
for more information on Elastic Data.

[Return to Elastic Data samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/docker/Dockerfile](src/main/docker/Dockerfile) - Dockerfile for creating side-car image from which configuration
  will be read from at pod startup

* [src/main/resources/conf/elastic-data-cache-config.xml](src/main/resources/conf/elastic-data-cache-config.xml) - cache config for storage-tier

Note if you wish to enable Prometheus or log capture, change the following in the helm installs to `true`. Their default values are false, but they are set to `false` in the instructions below for completeness.

* Prometheus: `--set prometheusoperator.enabled=true`

* Log capture: `--set logCaptureEnabled=true`

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/elastic-data/default` directory and ensure you have your maven build     
   environment set for JDK11 and build the project.

   ```bash
   mvn clean install -P docker
   ```

1. The result of the above is the docker image will be built with the cache configuration files
   with the name in the format proxy-tier-sample:${version}.

   For Example:

   ```bash
   elastic-data-sample-default:1.0.0-SNAPSHOT
   ```

   **Note:** If you are running against a remote Kubernetes cluster you will need to
   push the above image to your repository accessible to that cluster.
   
1. Install the Coherence cluster

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=elastic-data-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set store.cacheConfig=elastic-data-cache-config.xml \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --set userArtifacts.image=elastic-data-sample-default:1.0.0-SNAPSHOT \
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
   storage-coherence-2                                      1/1     Running   0          2m
   ```   
   
1. Add data to Elastic Data

   Connect to the Coherence `console` using the following to create a cache using `FlashJournal`.

   ```bash
   kubectl exec -it --namespace sample-coherence-ns storage-coherence-0 bash /scripts/startCoherence.sh console
   ```   
   
   At the `Map (?):` prompt, type `cache flash-01`.  This will create a cache in the service `DistributedSchemeFlash`
   which is a FlashJournal scheme.
   
   Use the following to add 100,000 objects of size 1024 bytes, starting at index 0 and using batches of 100.
   
   ```bash
   bulkput 100000 1024 0 100

   Mon Apr 15 07:37:03 GMT 2019: adding 100000 items (starting with #0) each 1024 bytes ...
   Mon Apr 15 07:37:15 GMT 2019: done putting (11578ms, 8878KB/sec, 8637 items/sec)
   ```
   
   At the prompt, type `size` and it should show 100000.
   
   Then type `bye` to exit the `console`.
   
1. Inspect the Elastic Data FlashJournal files  

   Issue the following to exec into one of the Coherence pods and list the files used by Elastic Data.
   
   ```bash
   # ls -l /tmp/
   total 84744
   -rw-r--r-- 1 root root 86769664 Apr 15 07:37 coh1781907747204398478.tmp
   drwxr-xr-x 2 root root     4096 Apr 15 07:49 hsperfdata_root

   ```
  
   The contents of the file(s) are not important, just the fact that they are using
   the default location.
   
   Type `exit` to leave the `exec` session.
   

## Verifying Grafana Data (If you enabled Prometheus)

Access Grafana using the instructions [here](../../../README.md#access-grafana).

## Verifying Kibana Logs (if you enabled log capture)

Access Kibana using the instructions [here](../../../README.md#access-kibana).

## Uninstalling the Charts

Carry out the following commands to delete the chart installed in this sample.

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.
