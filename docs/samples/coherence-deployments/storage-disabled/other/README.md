# Storage-disabled client in cluster as separate user image

To ensure that a custom made container and Helm Chart can be installed and join an existing Coherence
cluster, the following needs be done for the custom container:

* Ensure the correct version of Coherence is included in the packaging

* Set the `Cluster Name` to match the cluster of your running Coherence cluster

* Set the `Well Known Address` to match the headless service name of your running Coherence cluster

* Ensure local storage is disabled via `-Dcoherence.distributed.localstorage=false`

Below is an example of a (Helidon)[https://helidon.io/#/] web application which exposes a `/query` endpoint
allowing for CohQL commands to be passed in and executed against a Coherence cluster.


[Return to Storage-Disabled clients samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/assembly/helm-assembly.xml](src/assembly/helm-assembly.xml) - Assembly file for helm

* [src/main/java/com/oracle/coherence/examples/Main.java](src/main/java/com/oracle/coherence/examples/Main.java) - Entry point for Helidon web application

* [src/main/helm/](src/main/helm) - Helm chart files

Note if you wish to enable Prometheus or log capture, change the following in the helm installs to `true`. Their default values are false, but they are set to `false` in the samples below for completeness.

* Prometheus: `--set prometheusoperator.enabled=true`

* Log capture: `--set logCaptureEnabled=true`

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/storage-disabled/other` directory and ensure you have your maven build     
   environment set for JDK11 and build the project.

   ```bash
   $ mvn clean install
   ```
   
1. The result of the above is the docker image will be built with the cache configuration files
   with the name in the format helidon-sample:${version}. This image
   will be used by the chart.

   > Note: If you are running against a remote Kubernetes cluster you will need to
   > push the above image to your repository accessible to that cluster. You will also need to 
   > prefix the image name in your `helm` command below.
   
1. Install the Coherence cluster

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=2 \
      --set cluster=helidon-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=true \
      --set logCaptureEnabled=false \
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
   
1. Install the Helidon web application

   When we install the Helidon web application, we must ensure we set the following to allow
   the Helidon application to join the Coherence cluster.

   * `--set cluster=helidon-cluster` - same cluster name

   * `--set store.wka=storage-coherence-headless` - ensures it can contact the cluster
  
   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --set wka=storage-coherence-headless \
      --set clusterName=helidon-cluster \
      --name helidon-web-app \
      target/webserver-1.0.0-SNAPSHOT-helm/webserver/
   ```
   
   As per the instructions output above, port forward port 8080.
   
   ```bash
   $ export POD_NAME=$(kubectl get pods --namespace sample-coherence-ns -l "app=webserver,release=helidon-web-app" -o jsonpath="{.items[0].metadata.name}")
   $ kubectl --namespace sample-coherence-ns port-forward $POD_NAME 8080:8080
   Forwarding from 127.0.0.1:8080 -> 8080
   Forwarding from [::1]:8080 -> 8080
   ```
   
1. Issue CohQL commands

   Now issue various CohQL commands to create and mutate data in the Coherence cluster.
   
   ```bash
   $ curl -i -w '\n' -X PUT http://127.0.0.1:8080/query -d '{"query":"create cache foo"}'

   HTTP/1.1 200 OK
   Date: Thu, 18 Apr 2019 06:48:15 GMT
   transfer-encoding: chunked
   connection: keep-alive

   $ curl -i -w '\n' -X PUT http://127.0.0.1:8080/query -d '{"query":"insert into foo key(\"foo\") value(\"bar\")"}'

   HTTP/1.1 200 OK
   Date: Thu, 18 Apr 2019 06:48:40 GMT
   transfer-encoding: chunked
   connection: keep-alive

   $ curl -i -w '\n' -X PUT http://127.0.0.1:8080/query -d '{"query":"select key(),value() from foo"}'
   HTTP/1.1 200 OK
   Content-Type: application/json
   Date: Thu, 18 Apr 2019 06:49:15 GMT
   transfer-encoding: chunked
   connection: keep-alive

   {"result":"{foo=[foo, bar]}"}
   ```

## Uninstalling the Charts

Carry out the following commands to delete the two charts installed in this sample.

```bash
$ helm delete storage helidon-web-app --purge
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.
