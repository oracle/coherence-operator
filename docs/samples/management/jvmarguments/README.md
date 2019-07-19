# Tune JVM Runtime Settings

All production enterprise Java application must carefully tune the JVM settings for maximum performance. This sample shows how to set JVM arguments to Coherence running in Kubernetes.

Refer to [Coherence Performance Tuning
documentation](https://docs.oracle.com/middleware/12213/coherence/administer/performance-tuning.htm#GUID-2A0BC9E6-C3AA-4012-B3D8-EC51963B0CEB) for more information about JVM tuning.

There are many values in the
[values.yaml](https://github.com/oracle/coherence-operator/blob/master/operator/src/main/helm/coherence/values.yaml)
file of the Coherence Helm chart that sets JVM arguments to the JVM that runs Coherence within Kubernetes. The values are described in the following table:

| `--set` Arguments | Descirption |
|------------------------|---------|
| `store.maxHeap`        | Heap size arguments to the JVM. The format should be the same as that used for Java's -Xms and -Xmx JVM options. If not set, the JVM defaults are used. |
| `store.jmx.maxHeap` | Heap size arguments passed to the MBean server JVM.  Same format and meaning as the preceding row. |
| `store.jvmArgs` | Options passed directly to the JVM running Coherence within Kubernetes |
| `store.javaOpts` | Miscellaneous JVM options to pass to the Coherence store container |

[Return to Management samples](../) / [Return to samples](../../README.md#list-of-samples)

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../README.md#install-the-coherence-operator).

## Installation Steps
                    
1. Install the Coherence cluster

   Install the cluster with the following settings:
   
   * `--set store.maxHeap=1G` - Set `maxHeap` to 1 GB
   
   * `--set store.jvmArgs="-Xloggc:/tmp/gc-log -server -Xcomp"` - Set generic options
   
   * `--set store.javaOpts="-Dcoherence.log.level=6 -Djava.net.preferIPv4Stack=true` - Set Coherence specific arguments

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=coherence-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --set store.jvmArgs="-Xloggc:/tmp/gc-log -server -Xcomp"  \
      --set store.javaOpts="-Dcoherence.log.level=6 -Djava.net.preferIPv4Stack=true" \
      --set store.maxHeap=1g \
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
   storage-coherence-2                   0/1     Running   0          44s
   ```

   The JVM arguments include the `store.` argument and other arguments required by the operator and Coherence.

   ```
   -Xmx1g -Xms1g -Xloggc:/tmp/gc-log -server -Xcomp -Xms8g -Xmx8g -Dcoherence.log.level=6 -Djava.net.preferIPv4Stack=true
   ```

1. Inspect the result

   To inspect the full JVM arguments, you can use `kubectl logs storage-coherence-0 -n sample-coherence-ns > /tmp/storage-coherence-0.log` and search for one of the arguments that you specified.

## Uninstall the Charts

Use the following command to delete the chart installed in this sample:

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are removed from previous sample.

If you want to remove the `coherence-operator`, then include it in the `helm delete` command.
