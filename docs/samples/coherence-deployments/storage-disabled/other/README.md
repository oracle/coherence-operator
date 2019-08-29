# Storage-Disabled Client in Cluster as Separate User Image

To ensure that a custom made container and Helm Chart can be installed and joins an existing Coherence
cluster, perform the following steps for the custom container:

* Ensure that the correct version of Coherence is included in the package.

* Set `Cluster Name` to match the name of the Coherence cluster you are running.

* Set  `Well Known Address` to match the headless service name of the Coherence cluster you are running.

* Disable local storage using the setting, `-Dcoherence.distributed.localstorage=false`

This sample shows how a [Helidon](https://helidon.io/#/) web application exposes a `/query` endpoint,
allowing `CohQL` commands to be passed and executed against a Coherence cluster.


[Return to Storage-Disabled clients samples](../) / [Return to Coherence Deployments samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/assembly/helm-assembly.xml](src/assembly/helm-assembly.xml) - Assembly file for Helm

* [src/main/java/com/oracle/coherence/examples/Main.java](src/main/java/com/oracle/coherence/examples/Main.java) - Entry point for Helidon web application

* [src/main/helm/](src/main/helm) - Helm chart files

> **Note:** If you want to enable Prometheus or log capture, ensure you set the appropriate properties for the Coherence Operator install:

* Prometheus: `--set prometheusoperator.enabled=true`

* Log capture: `--set logCaptureEnabled=true`

## Prerequisites

Ensure that you have installed the Coherence Operator by following the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/coherence-deployments/storage-disabled/other` directory. Ensure that you have your maven build environment set for JDK8, and build the project.

   ```bash
   $ mvn clean install -P docker
   ```

   This builds the Docker image with the cache configuration files, with the name in the format, `helidon-sample:${version}`. This image
   is subsequently used by the chart.

   > **Note:** If you are running against a remote Kubernetes cluster, you must
   > push the above image to your repository accessible to that cluster. You must also
   > prefix the image name in the `helm` command, as shown below.

2. Install the Coherence cluster.

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=helidon-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set logCaptureEnabled=false \
      coherence/coherence
   ```

   After the installation is complete, get the list of pods by running the following command:

   ```bash
   $ kubectl get pods -n sample-coherence-ns
   ```
   All the three storage-coherence pods should be running and ready, as shown in the output:

   ```console
   NAME                    READY   STATUS    RESTARTS   AGE
   storage-coherence-0     1/1     Running   0          4m
   storage-coherence-1     1/1     Running   0          2m   
   storage-coherence-2     1/1     Running   0          2m
   ```


3. Install the Helidon web application.

   When you install the Helidon web application, ensure that  you set the following properties to allow the Helidon application to join the Coherence cluster:

   * `--set cluster=helidon-cluster` - Uses the same cluster name

   * `--set store.wka=storage-coherence-headless` - Ensures that it can contact the cluster

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

4. Run `CohQL` commands.

   Use the various `CohQL` commands to create and mutate data in the Coherence cluster.

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

## Uninstall the Charts

Run the following command to delete both the charts installed in this sample.

```bash
$ helm delete storage helidon-web-app --purge
```

Before starting another sample, ensure that all  pods are removed from the previous sample. To remove `coherence-operator`, use the `helm delete` command.
