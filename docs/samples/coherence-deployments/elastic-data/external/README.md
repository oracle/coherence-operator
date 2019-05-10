# Deploy using external volume mapped to the host

This sample shows how to create Persistent Volumes (PV) and then map the Elastic Data 
to be stored on these PV's.  This would allow for a specific size to be used for storing Elastic Data
rather than just relying on the size of the underlying default "/tmp/ directory.

Please see [Oracle Elastic Data Documentation](https://docs.oracle.com/middleware/12213/coherence/COHDG/implementing-storage-and-backing-maps.htm#COHDG5496) 
for more information on Elastic Data.

[Return to Elastic Data samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/docker/Dockerfile](src/main/docker/Dockerfile) - Dockerfile for creating side-car image from which configuration
  will be read from at pod startup

* [src/main/resources/conf/elastic-data-cache-config.xml](src/main/resources/conf/elastic-data-cache-config.xml) - cache config for storage-tier

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/elastic-data/external` directory and ensure you have your maven build     
   environment set for JDK8 and build the project.

   ```bash
   $ mvn clean install -P docker
   ```

1. The result of the above is the docker image will be built with the cache configuration files
   with the name in the format elastic-data-sample-external:${version}.

   For Example:

   ```bash
   elastic-data-sample-external:1.0.0-SNAPSHOT
   ```

   > Note: If you are running against a remote Kubernetes cluster you will need to
   > push the above image to your repository accessible to that cluster. You will also need to 
   > prefix the image name in your `helm` command below.
    
1. Install the Coherence cluster

   You can use the following three options to specify volumes, volumeMounts & volumeClaimTemplates
   
   * `--set store.volumes` - defines extra volume mappings that will be added to the Coherence Pod
   
   * `--set store.volumeClaimTemplates` - defines extra PVC mappings that will be added to the Coherence Pod
   
   * `--set store.volumeMounts` - defines extra volume mounts to map to the additional volumes or PVCs declared above
      in store.volumes and store.volumeClaimTemplates
   
   In our example we are going to use a `yaml` file ([volumes.yaml](src/main/yaml/volumes.yaml)) to specify
   hostPath volumes. 
   
   > Note: You should set the values appropriately for your Kubernetes environment and needs.
   
   We will also set `--set store.javaOpts="-Dcoherence.flashjournal.dir=/elastic-data" ` - to point Elastic data to the mount path
   
   > Note: The `coherence.flashjournal.dir` option was only added in Coherence 12.2.1.4, so we must include
   > an override file to define this so it works in 12.2.1.3.X as well. 
   
   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set store.javaOpts="-Dcoherence.flashjournal.dir=/elastic-data" \
      --set clusterSize=3 \
      --set cluster=elastic-data-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set store.cacheConfig=elastic-data-cache-config.xml \
      --set store.overrideConfig=elastic-data-override.xml \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --set userArtifacts.image=elastic-data-sample-external:1.0.0-SNAPSHOT \
      -f src/main/yaml/volumes.yaml \
      coherence/coherence
   ```
   
   Use `kubectl get pods -n sample-coherence-ns` to ensure that all pods are running.
   All 3 storage-coherence-0/1/2 pods should be running and ready, as below:

   ```bash
   NAME                  READY   STATUS    RESTARTS   AGE
   storage-coherence-0   1/1     Running   0          4m
   storage-coherence-1   1/1     Running   0          2m   
   storage-coherence-2   1/1     Running   0          2m
   ```   
   
1. Confirm the mounted volume

   ```bash
   $ kubectl exec -it storage-coherence-0   -n sample-coherence-ns -- bash -c df
   Filesystem     1K-blocks     Used Available Use% Mounted on
   overlay         61252420 16809112  41302140  29% /
   tmpfs              65536        0     65536   0% /dev
   tmpfs            4334408        0   4334408   0% /sys/fs/cgroup
   /dev/sda1       61252420 16809112  41302140  29% /logs
   overlay          4334408      356   4334052   1% /elastic-data
   shm                65536        0     65536   0% /dev/shm
   tmpfs            4334408       12   4334396   1% /run/secrets/kubernetes.io/serviceaccount
   tmpfs            4334408        0   4334408   0% /proc/acpi
   tmpfs            4334408        0   4334408   0% /sys/firmware
   ```   
   
   You should see your `/elastic-data` volume mounted.
   
1. Add data to Elastic Data

   Connect to the Coherence `console` using the following to create a cache using `FlashJournal`.

   ```bash
   $ kubectl exec -it --namespace sample-coherence-ns storage-coherence-0 bash /scripts/startCoherence.sh console
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
   
1. Ensure the Elastic Data FlashJournal files exist 

   Issue the following to exec into one of the Coherence pods and list the files used by Elastic Data.
   
   ```bash
   $ kubectl exec -it -n sample-coherence-ns storage-coherence-0 -- bash -c 'ls -l /elastic-data'

   total 202496
   -rw-r--r-- 1 root root 69468160 May  1 08:44 coh1207347469383108692.tmp
   -rw-r--r-- 1 root root 68943872 May  1 08:44 coh5447980795344195354.tmp
   -rw-r--r-- 1 root root 68943872 May  1 08:44 coh6857569664465911116.tmp
   ```
  
   Type `exit` to leave the `exec` session.

## Uninstalling the Charts

Carry out the following commands to delete the chart installed in this sample.

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.
