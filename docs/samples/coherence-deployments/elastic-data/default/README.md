# Deploy Elastic Data Using Default FlashJournal Locations

This sample shows how to enable Elastic Data using the default FlashJournal directory, `/tmp/`.  

Refer to the [Oracle Elastic Data documentation](https://docs.oracle.com/middleware/12213/coherence/COHDG/implementing-storage-and-backing-maps.htm#COHDG5496)
for more information about Elastic Data.

[Return to Elastic Data samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/docker/Dockerfile](src/main/docker/Dockerfile) - Dockerfile for creating sidecar image from which the configuration will be read at pod startup

* [src/main/resources/conf/elastic-data-cache-config.xml](src/main/resources/conf/elastic-data-cache-config.xml) - Cache configuration for storage-tier

## Prerequisites

Ensure that you have installed Coherence Operator by following the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/elastic-data/default` directory. Ensure that you have your maven build environment set for JDK8, and build the project.

   ```bash
   $ mvn clean install -P docker
   ```
 This builds the Docker image with the cache configuration files, with the name in the format `elastic-data-sample-default:${version}`. For example,

   ```bash
   elastic-data-sample-default:1.0.0-SNAPSHOT
   ```

   > **Note:** If you are running against a remote Kubernetes cluster, then you must
   > push the above image to your repository accessible to that cluster. You must also
   > prefix the image name in the `helm` command, as shown  below.

2. Install the Coherence cluster.

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
      coherence/coherence
   ```

   After the installation is complete, list the pods by running the command:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   All the three storage-coherence pods should be running and ready, as shown in the output:
   ```console
   NAME                   READY   STATUS    RESTARTS   AGE
   storage-coherence-0    1/1     Running   0          4m
   storage-coherence-1    1/1     Running   0          2m   
   storage-coherence-2    1/1     Running   0          2m
   ```   

3. Add data to Elastic Data.

   a. Connect to the Coherence console to create a cache using `FlashJournal`:

   ```bash
   $ kubectl exec -it --namespace sample-coherence-ns storage-coherence-0 -- bash /scripts/startCoherence.sh console
   ```   

   At the `Map (?):` prompt, type `cache flash-01`.  This creates a cache in the service, `DistributedSchemeFlash`
   which is a FlashJournal scheme.

   b. Use the following to add 100,000 objects of size 1024 bytes, starting at index 0, and using batches of 100.

   ```bash
   bulkput 100000 1024 0 100

   Mon Apr 15 07:37:03 GMT 2019: adding 100000 items (starting with #0) each 1024 bytes ...
   Mon Apr 15 07:37:15 GMT 2019: done putting (11578ms, 8878KB/sec, 8637 items/sec)
   ```

   At the prompt, type `size` and it should return 100000.

   Then, type `bye` to exit the `console`.

4. Ensure that the Elastic Data FlashJournal files exist.

   Run the following command against one of the Coherence pods to list the files used by Elastic Data:

   ```bash
   $ kubectl exec -it -n sample-coherence-ns storage-coherence-0 -- bash -c 'ls -l /tmp/'
   ```
   ```console
   total 84744
   -rw-r--r-- 1 root root 86769664 Apr 15 07:37 coh1781907747204398478.tmp
   drwxr-xr-x 2 root root     4096 Apr 15 07:49 hsperfdata_root
   ```

   Type `exit` to leave the `exec` session.

## Uninstall the Charts

Run the following commands to delete the chart installed in this sample:

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all  pods are removed from the previous sample. To remove `coherence-operator`, use the `helm delete` command.
