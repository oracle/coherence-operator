# Deploy Elastic Data Using External Volume Mapped to the Host

This sample shows how to create persistent volumes (PV), and then map the Elastic Data to be stored on these PV. This would allow for a specific size to be used for storing Elastic Data, rather than only relying on the size of the underlying default `/tmp/` directory.

Refer to the [Oracle Elastic Data Documentation](https://docs.oracle.com/middleware/12213/coherence/COHDG/implementing-storage-and-backing-maps.htm#COHDG5496)
for more information about Elastic Data.

[Return to Elastic Data samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/docker/Dockerfile](src/main/docker/Dockerfile) - Dockerfile for creating sidecar image from which the configuration will be read at pod startup

* [src/main/resources/conf/elastic-data-cache-config.xml](src/main/resources/conf/elastic-data-cache-config.xml) - Cache configuration for storage-tier

## Prerequisites

Ensure that you have installed Coherence Operator by following the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/elastic-data/external` directory. Ensure that you have your maven build environment set for JDK8, and build the project.

   ```bash
   $ mvn clean install -P docker
   ```

  This builds the docker image with the cache configuration files, with the name in the format, `elastic-data-sample-external:${version}`. For example,

   ```bash
   elastic-data-sample-external:1.0.0-SNAPSHOT
   ```

   > **Note:** If you are running against a remote Kubernetes cluster, then you must
   > push the above image to your repository accessible to that cluster. You must also
   > prefix the image name in the `helm` command, as shown below.

2. Install the Coherence cluster.

   You can use the following three options to specify `volumes`, `volumeMounts` & `volumeClaimTemplates`:

   * `--set store.volumes` - Defines extra volume mappings that will be added to the Coherence Pod

   * `--set store.volumeClaimTemplates` - Defines extra PVC mappings that will be added to the Coherence Pod

   * `--set store.volumeMounts` - Defines extra volume mounts to map to the additional volumes or PVC declared above in `store.volumes` and `store.volumeClaimTemplates`

   For this sample, we are going to use the YAML file, [volumes.yaml](src/main/yaml/volumes.yaml) to specify hostPath volumes.

   > **Note:** You should set the values appropriately for your Kubernetes environment and needs.

   Also, set `--set store.javaOpts="-Dcoherence.flashjournal.dir=/elastic-data" ` - to point Elastic data to the mount path.

   > **Note:** The `coherence.flashjournal.dir` option was only added in Coherence 12.2.1.4. Therefore, we must include
   > an override file to define this so that it works in 12.2.1.3.x as well.

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set store.javaOpts="-Dcoherence.flashjournal.dir=/elastic-data" \
      --set clusterSize=1 \
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

   After installing, list the pods by running the following command:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   All the three storage-coherence pods must be running and ready, as shown in the output:
   ```console
   NAME                  READY   STATUS    RESTARTS   AGE
   storage-coherence-0   1/1     Running   0          4m
   storage-coherence-1   1/1     Running   0          2m
   storage-coherence-2   1/1     Running   0          1m
   ```
3. Confirm the mounted volume.

   ```bash
   $ kubectl exec -it storage-coherence-0   -n sample-coherence-ns -- bash -c df
   ```
   You should see your `/elastic-data` volume mounted.
   ```console
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
4. Add data to Elastic Data.

   a. Connect to the Coherence console to create a cache using `FlashJournal`:

   ```bash
   $ kubectl exec -it --namespace sample-coherence-ns storage-coherence-0 bash /scripts/startCoherence.sh console
   ```   

   At the `Map (?):` prompt, type `cache flash-01`.  This creates a cache in the service, `DistributedSchemeFlash`
   which is a FlashJournal scheme.

   b. Use the following to add 100,000 objects of size 1024 bytes, starting at index 0, and using batches of 100.

   ```bash
   bulkput 100000 1024 0 100

   Mon Apr 15 07:37:03 GMT 2019: adding 100000 items (starting with #0) each 1024 bytes ...
   Mon Apr 15 07:37:15 GMT 2019: done putting (11578ms, 8878KB/sec, 8637 items/sec)
   ```

   At the prompt, type `size` and it should show 100000.

   Then, type `bye` to exit the `console`.

5. Ensure that the Elastic Data FlashJournal files exist.

   Run the following command against one of the Coherence pods to list the files used by Elastic Data:

   ```bash
   $ kubectl exec -it -n sample-coherence-ns storage-coherence-0 -- bash -c 'ls -l /elastic-data'
   ```
   ```console
   total 202496
   -rw-r--r-- 1 root root 69468160 May  1 08:44 coh1207347469383108692.tmp
   -rw-r--r-- 1 root root 68943872 May  1 08:44 coh5447980795344195354.tmp
   -rw-r--r-- 1 root root 68943872 May  1 08:44 coh6857569664465911116.tmp
   ```

   Type `exit` to leave the `exec` session.

## Uninstall the Charts

Run the following commands to delete the chart installed in this sample:

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all  pods are removed from the previous sample. To remove `coherence-operator`, use the `helm delete` command.
