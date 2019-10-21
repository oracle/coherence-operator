# Change Image Version for Coherence or Application Container Using Rolling Upgrade

Samples, such as [Storage-disabled client in cluster via interceptor](../../coherence-deployments/storage-disabled/interceptor) call for the creation of a sidecar docker image.
Sidecar docker image provides the application classes to Kubernetes. The docker image is tagged with a version number and this version number is used by Kubernetes to enable safe rolling upgrades. You can read more about safe rolling upgrades in the [Helm documentation](https://helm.sh/docs/helm/#helm-upgrade).

The safe rolling upgrade feature allows you to instruct Kubernetes, through the operator, to replace the currently installed version of your application classes with a different one. Kubernetes does not verify whether the classes are new or old. It checks whether the image can be pulled by the cluster and image has a docker tag. The operator also ensures that the replacement is done without data loss or interruption of service.

This sample initially deploys version 1.0.0 of the sidecar Docker image and then does a rolling upgrade to
version 2.0.0 of the sidecar image which introduces a server side Interceptor to modify data to ensure it is stored as uppercase.

[Return to Storage-Disabled clients samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../README.md#list-of-samples)

## Sample files

* [src/main/docker-v1/Dockerfile](src/main/docker-v1/Dockerfile) - Dockerfile for creating sidecar which includes  only v1 storage-cache-config.xml

* [src/main/docker-v2/Dockerfile](src/main/docker-v2/Dockerfile) - Dockerfile for creating sidecar which includes v2 storage-cache-config.xml

* [src/main/resources/conf/v1/storage-cache-config.xml](src/main/resources/conf/v1/storage-cache-config.xml) - Cache configuration version 1 for storage-enabled tier without interceptor

* [src/main/resources/conf/v2/storage-cache-config.xml](src/main/resources/conf/v2/storage-cache-config.xml) - Cache configuration version 2 for storage-enabled tier that includes uppercase interceptor

* [src/main/resources/client-cache-config.xml](src/main/resources/client-cache-config.xml) - Client configuration for extend client

* [src/main/java/com/oracle/coherence/examples/UppercaseInterceptor.java](src/main/java/com/oracle/coherence/examples/UppercaseInterceptor.java) - Interceptor that changes all entries to uppercase - version 2.0.0

## Prerequisites

Ensure you have already installed the Coherence Operator using the instructions [here](../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/operator/rolling-upgrade` directory and ensure you have your Maven build environment set for JDK 8 and build the project:

   ```bash
   $ mvn clean install -P docker-v1,docker-v2
   ```

   The version 1 and version 2  Docker images are created:

   * `rolling-upgrade-sample:1.0.0`

   * `rolling-upgrade-sample:2.0.0`

   `rolling-upgrade-sample:1.0.0` is the initial image installed in the chart.

   > **Note**: If you are using a remote Kubernetes cluster, you need to push the created images to your repository accessible to that cluster. You need to prefix the image name in the `helm` command.

1. Install the Coherence cluster with rolling-upgrade-sample:1.0.0 image as a sidecar:

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=rolling-upgrade-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set store.cacheConfig=storage-cache-config.xml \
      --set logCaptureEnabled=false \
      --set userArtifacts.image=rolling-upgrade-sample:1.0.0 \
      coherence/coherence
   ```

   After the installation completes, list the pods:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   ```console
   NAME                   READY   STATUS    RESTARTS   AGE
   storage-coherence-0    1/1     Running   0          4m
   storage-coherence-1    1/1     Running   0          2m
   storage-coherence-2    1/1     Running   0          1m
   ```

   All the three storage-coherence-0/1/2 pods are in running state.

1. Port forward the proxy port on the `storage-coherence-0` pod:

   ```bash
   $ kubectl port-forward -n sample-coherence-ns storage-coherence-0 20000:20000
   ```

1. Connect via CohQL commands and execute the following command:

   ```bash
   $ mvn exec:java
   ```

   Run the following CohQL commands to insert data into the cluster:

   ```sql
   insert into 'test' key('key-1') value('value-1');
   insert into 'test' key('key-2') value('value-2');

   select key(), value() from 'test';
   Results
   ["key-1", "value-1"]
   ["key-2", "value-2"]
   ```

1. Upgrade the helm release to use the `rolling-upgrade-sample:2.0.0` image.

   Use the following arguments to upgrade to version 2.0.0 of the image:

   * `--reuse-values` - specifies to reuse all previous values associated with the release

   * `--set userArtifacts.image=rolling-upgrade-sample:2.0.0` - the new artifact version

   ```bash
   $ helm upgrade storage coherence/coherence \
      --namespace sample-coherence-ns \
      --reuse-values \
      --set imagePullSecrets=sample-coherence-secret \
      --set userArtifacts.image=rolling-upgrade-sample:2.0.0
   ```

1. Check the status of the upgrade.

   Use the following command to check the status of the rolling upgrade of all pods.

   > **Note**: The command below will not return until upgrade of all pods is complete.

   ```bash
   $ kubectl rollout status sts/storage-coherence --namespace sample-coherence-ns
   Waiting for 1 pods to be ready...
   Waiting for 1 pods to be ready...
   waiting for statefulset rolling update to complete 1 pods at revision storage-coherence-67b75785f6...
   Waiting for 1 pods to be ready...
   Waiting for 1 pods to be ready...
   waiting for statefulset rolling update to complete 2 pods at revision storage-coherence-67b75785f6...
   Waiting for 1 pods to be ready...
   Waiting for 1 pods to be ready...
   statefulset rolling update complete 3 pods at revision storage-coherence-67b75785f6...
   ```

1. Verify the data through CohQL commands.

   When the upgrade is running, you can execute the following commands in the CohQL session:

   ```sql
   select key(), value() from 'test';
   ```

   You can note that the data always remains the same.

   > **Note**: Your port-forward fails when the storage-coherence-0` pod restarts. You have to stop and restart it.

   In an environment where you have configured a load balancer, then the Coherence*Extend session automatically reconnects when it detects a disconnect.

1. Add new data to confirm the interceptor is now active.  

   ```sql
   insert into 'test' key('key-3') value('value-3');

   select key(), value() from 'test';
   Results
   ["key-1", "value-1"]
   ["key-3", "VALUE-3"]
   ["key-2", "value-2"]
   ```

   You can note that the value for `key-3` has been converted to uppercase which shows that the server-side interceptor is now active.

1. Verify that the 2.0.0 image on one of the pods.

   Use the following command to verify that the 2.0.0 image is active:

   ```bash
   $ kubectl describe pod storage-coherence-0  -n sample-coherence-ns | grep rolling-upgrade
   ```
   ```console
   Image:         rolling-upgrade-sample:2.0.0
   Normal  Pulled                 4m59s  kubelet, docker-for-desktop  Container image "rolling-upgrade-sample:2.0.0" already present on machine
   ```

   The output shows that the version 2.0.0 image is now present.

## Uninstall the Chart

Use the following command: to delete the chart installed in this sample:

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are removed from previous samples.

If you want to remove the `coherence-operator`, then use `helm delete` command.
