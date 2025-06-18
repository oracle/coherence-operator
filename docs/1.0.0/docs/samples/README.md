# Samples

Redirecting, please update your bookmarks...

<script>
  window.location.href = "https://docs.coherence.community/coherence-operator/docs/latest/examples/000_overview";
</script>


The samples provide demonstrations of how to accomplish common tasks. These samples are provided for educational and demonstration purposes only; they are not intended to be used in production deployments or to be depended upon to create production environments.

Read through the samples, understand how they work, and then customize them to your requirements.

> See the [Oracle Coherence Demonstration](https://github.com/coherence-community/coherence-demo) for an example of an application using Coherence and running via the Coherence Operator.

# Table of Contents
- [Get Started](#get-started)
  * [Check Runtime Prerequisites](#check-runtime-prerequisites)
  * [Check JDK 8 and Maven Installation](#check-jdk-8-and-maven-installation)
  * [Install Coherence into Local Maven Repository](#install-coherence-into-local-maven-repository)
  * [Create a Secret](#create-a-secret)
    + [Use Docker Hub for Images](#use-docker-hub-for-images)
    + [Provide a Secret to Kubernetes for a Repository on Docker Hub](#provide-a-secret-to-kubernetes-for-a-repository-on-docker-hub)
  * [Create the Sample Namespace](#create-the-sample-namespace)
  * [Clone the GitHub Repository](#clone-the-github-repository)
  * [Install the Coherence Operator](#install-the-coherence-operator)
    + [Install the Coherence Operator Without Prometheus and Log Capture](#install-the-coherence-operator-without-prometheus-and-log-capture)
    + [Enable Prometheus](#enable-prometheus)
    + [Enable Log Capture](#enable-log-capture)
  * [Enable Prometheus and Log Capture](#enable-prometheus-and-log-capture)
  * [Check Operator Status](#check-operator-status)
- [List Of Samples](#list-of-samples)
- [Troubleshooting Tips](#troubleshooting-tips)
  * [Coherence Cluster pods never reach ready "1/1"](#coherence-cluster-pods-never-reach-ready--1-1-)
  * [Error: ImagePullBackOff after installing Operator or coherence](#error--imagepullbackoff-after-installing-operator-or-coherence)
  * [Error: configmaps "coherence-internal-config" not found](#error--configmaps--coherence-internal-config--not-found)
  * [Unable to delete pods when using log capture](#unable-to-delete-pods-when-using-log-capture)
  * [Receive Error 'no matches for kind "Prometheus"'](#receive-error--no-matches-for-kind--prometheus--)
- [Accessing UI endpoints](#accessing-ui-endpoints)
  * [Access Grafana](#access-grafana)
  * [Access Kibana](#access-kibana)
  * [Access Prometheus](#access-prometheus)
- [Run Samples Integration Tests](#run-samples-integration-tests)

# Get Started

To setup Coherence Operator, follow these steps:

1. [Check Runtime Prerequisites](#check-runtime-prerequisites)
1. [Check JDK 8 and Maven Installation](#check-jdk-8-and-maven-installation)
1. [Install Coherence into Local Maven Repository](#install-coherence-into-local-maven-repository)
1. [Create the Sample Namespace](#create-the-sample-namespace)
5. [Create a Secret](#create-a-secret)
1. [Clone the GitHub Repository](#clone-the-github-repository)
1. [Install the Coherence Operator](#install-the-coherence-operator)

If you have already run samples before, you can go to the [List of Samples](#list-of-samples).

## Check Runtime Prerequisites

Refer to the following sections in the Quick Start Guide for software versions and runtime environment configuration:

* Runtime Environment Prerequisites

  * [Software Requirements](../quickstart.md#software-requirements) - Helm and Kubernetes versions

  * [Runtime Environment Requirements](../quickstart.md#runtime-environment-requirements) - Helm and Kubernetes configuration

* Environment Configuration

  * [Add the Helm repository for Coherence](../quickstart.md#add-the-helm-repository-for-coherence)

  * [Obtain Images from Oracle Container Registry](../quickstart.md#obtain-images-from-oracle-container-registry)


> **Note:** For all the `helm install` commands, you can leave the --version option off and the latest version of the chart is retrieved. If you wanted to use a specific version, such as 0.9.8, add `--version 0.9.8` to all installs for the `coherence-operator` and `coherence` charts.  

## Check JDK 8 and Maven Installation

Ensure that you have the following installed:

* JDK 8+

* Maven 3.5.4+

> **Note:** You can use a later version of Java, for example, JDK11, as the
> `maven.compiler.source` and `target` are set to JDK 8 in the sample
> `pom.xml` files.

## Install Coherence into Local Maven Repository

If you are running samples that have a Maven project, then follow these steps:

1. Download and install Oracle Coherence 12.2.1.3 from [Oracle Technology Network](https://www.oracle.com/technetwork/middleware/coherence/downloads/index.html).

2. Ensure that the COHERENCE_HOME environment variable is set to point to the `coherence` directory under your install location containing the bin, lib, and doc directories. This is required only for the Maven `install-file` commands.

3. Install Coherence into your local Maven repository:

   ```bash
   $ mvn install:install-file -Dfile=$COHERENCE_HOME/lib/coherence.jar   \
    -DpomFile=$COHERENCE_HOME/plugins/maven/com/oracle/coherence/coherence/12.2.1/coherence.12.2.1.pom
   ```

   If you are running Coherence 12.2.1.4, you need to install `coherence-metrics`.

   ```bash
   $ mvn install:install-file -Dfile=$COHERENCE_HOME/lib/coherence-metrics.jar -DpomFile=$COHERENCE_HOME/plugins/maven/com/oracle/coherence/coherence-metrics/12.2.1/coherence-metrics.12.2.1.pom
   ``` 

## Create the Sample Namespace

You need to create the namespace for the first time to run any of the samples. Create your target namespace:

  ```bash
  $ kubectl create namespace sample-coherence-ns

  namespace/sample-coherence-ns created
  ```
In the samples, a Kubernetes namespace called `sample-coherence-ns` is used. If you want to change this namespace, ensure that you change any references to this namespace to match your selected namespace.

## Create a Secret

If all of your images can be pulled from public repositories, this step is not required. Otherwise, you need to enable your Kubernetes cluster to pull images from private repositories. You must create a secret to convey the docker credentials to Kubernetes. In the samples, the secret named `sample-coherence-secret` is used in the namespace `sample-coherence-ns`.

See [https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/) for more information.

  ```bash
  $ kubectl get secret sample-coherence-secret -n sample-coherence-ns
  ```
  ```console
  NAME                      TYPE                             DATA   AGE
  sample-coherence-secret   kubernetes.io/dockerconfigjson   1      18s
  ```

### Use Docker Hub for Images

You can pull Docker images from one repository, modify them, re-tag, and push the images to the repositories that you own. For example, repositories that you can create on [Docker Hub](http://hub.docker.com). The repositories are public by default, but can also be private.

### Provide a Secret to Kubernetes for a Repository on Docker Hub

In this example, the Coherence 12.2.1.3 Docker image is pulled from the
Docker Store, re-tagged, and pushed in to a freshly created docker
repository within your Docker Hub account. Then, create a secret to allow Kubernetes to pull that image. This approach can also be used when pushing sidecar images in other samples.

1. Sign in to [Docker Hub](https://hub.docker.com/).

2. Click **Create a Repository** on the Docker Welcome page.
3. Enter `Coherence` as the Name of the repository and in the Description type `Re-tags of Official Coherence Image`.
4. Select Public or Private for the repository and click **Create**.
5. Note the docker tag displayed.
   ```bash
   docker push <dockerid>/coherence:tagname
   ```
   `<dockerid>/coherence:tagname` is the docker tag for the repository.
6. Follow the steps in the section [Obtain Images from Oracle Container Registry](../quickstart.md#obtain-images-from-oracle-container-registry) to get the Coherence 12.2.1.3.x Docker image.
7. Re-tag the Coherence 12.2.1.3.x Docker image with your repository and docker tag.
   ```bash
   docker tag store/oracle/coherence:12.2.1.3 <dockerid>/coherence:12.2.1.3
    ```
8. Log in to your Docker Hub account:
   ```bash
   docker login
   ```
   Enter your Docker ID and password.
9. Push the re-tagged image to your Docker repository:
   ```bash
   docker push <dockerid>/coherence:12.2.1.3
    ```
   Note the value of the `The push refers to repository` statement. The first part of that is necessary to create the secret for Kubernetes. In this case, it should be something like `docker.io/mydockerid/coherence`.
10. Create the secret within your namespace:
    ```bash
    kubectl create secret docker-registry sample-coherence-secret 
      --namespace sample-coherence-ns --docker-server=hub.docker.com 
      --docker-username=docker.io/<dockerid> --docker-password="your docker password" 
      --docker-email="the email address of your docker account"
      ```
      When invoking `helm`, you can specify one or more secrets using the `--set imagePullSecrets` option.

      ```bash
       --set "imagePullSecrets{sample-coherence-secret}"
       ```

## Clone the GitHub Repository

The samples exist in the `gh-pages` branch of the Coherence Operator GitHub repository - https://github.com/oracle/coherence-operator.

Clone the repository and switch to the `gh-pages` branch:

```bash
$ git clone https://github.com/oracle/coherence-operator

$ cd coherence-operator

$ git checkout gh-pages

$ cd docs/samples
```
In the `samples` root directory, check the [`pom.xml`](pom.xml) and verify that the value of the `coherence.version` property matches the version of Coherence that you are actually using. For example, if you have Coherence 12.2.1.3.0, then the value of `coherence.version` must be `12.2.1-3-0`.  If this value needs ajustment, use the `-Dcoherence.version=` argument for all invocations of `mvn`.

Use the following command to ensure that all the projects with source code build correcly:

```bash
$ mvn clean install
```
> **Note**: Any compilation errors indicates that the Coherence JARs are not properly installed or you have not set the JDK correctly in your system.

## Install the Coherence Operator

Install the operator first to try out the samples. When you install the operator, you can optionally enable the following:

1. Prometheus integration: Captures metrics and displays in
   Grafana. (Available only for Coherence 12.2.1.4.0 or later)

2. Log capture: Uses Fluentd to send logs to Elasticsearch which can be then viewed in Kibana.

When you enable both Prometheus and log capture, you require extra system resources.

> **Note:** When you are running the operator locally, for example, Docker on Mac, you should allocate sufficient memory to your Docker for Mac process. The minimum recommended memory to run is 8G.

### Install the Coherence Operator Without Prometheus and Log Capture

The following command installs the operator without Prometheus or log capture enabled:

```bash
$ helm install \
   --namespace sample-coherence-ns \
   --set imagePullSecrets=sample-coherence-secret \
   --name coherence-operator \
   --set "targetNamespaces={sample-coherence-ns}" \
   coherence/coherence-operator
```
### Enable Prometheus

> **Note:** Use of Prometheus and Grafana is available only when using the
operator with Coherence 12.2.1.4 or later version.

To enable Prometheus, add the following options to the operator installation command:

```bash
   --set prometheusoperator.enabled=true \
   --set prometheusoperator.prometheusOperator.createCustomResource=false
```

The complete `helm install` example to enable Prometheus is as follows:

```bash
$ helm install \
     --namespace sample-coherence-ns \
     --set imagePullSecrets=sample-coherence-secret \
     --name coherence-operator \
     --set prometheusoperator.enabled=true \
     --set prometheusoperator.prometheusOperator.createCustomResource=false \
     --set "targetNamespaces={sample-coherence-ns}" \
     coherence/coherence-operator
```

> **Note:** When you install `prometheusOperator` for the first time, you must set `createCustomResource=true`. For subsequent installation of the operator, it must be set to `false`.

### Enable Log Capture

 To enable log capture, which includes Fluentd, Elasticsearch and Kibana, add the following options to the `helm install` command:

 ```bash
 --set logCaptureEnabled=true
 ```

 The complete `helm install` example to enable log capture is as follows:

 ```bash
$ helm install \
    --namespace sample-coherence-ns \
    --set imagePullSecrets=sample-coherence-secret \
    --name coherence-operator \
    --set logCaptureEnabled=true \
    --set "targetNamespaces={sample-coherence-ns}" \
    coherence/coherence-operator
  ```

## Enable Prometheus and Log Capture

You can enable both Prometheus and log capture by setting both of the options to `true`.

## Check Operator Status

Use `kubectl get pods -n sample-coherence-ns` to ensure that all the pods are in running status. When you enable Prometheus, the following output is displayed:

```console
NAME                                                     READY   STATUS    RESTARTS   AGE
coherence-operator-66f9bb7b75-nxwdc                      1/1     Running   0          13m
coherence-operator-grafana-898fc8bbd-nls45               3/3     Running   0          13m
coherence-operator-kube-state-metrics-5d5f6855bd-klzj5   1/1     Running   0          13m
coherence-operator-prometh-operator-58bd58ddfd-dhd9q     1/1     Running   0          13m
coherence-operator-prometheus-node-exporter-5hxwh        1/1     Running   0          13m
prometheus-coherence-operator-prometh-prometheus-0       3/3     Running   1          12m
```
Depending upon the number of CPU cores, you can see multiple node-exporter processes.

When you enable log capture, the following output is displayed:

```console
NAME                                  READY   STATUS    RESTARTS   AGE
coherence-operator-64b4f8f95d-fmz2x   2/2     Running   0          2m
elasticsearch-5b5474865c-tlr44        1/1     Running   0          2m
kibana-f6955c4b9-n8krf                1/1     Running   0          2m
```

# List Of Samples

Samples legend:
* &#x2714; - Available for Coherence 12.2.1.3.x and above

* &#x2726; - Available for Coherence 12.2.1.4.x and above


1. [Coherence Operator](operator/)
   1. [Logging](operator/logging)
      1. [Enable Log capture to View Logs in Kiabana ](operator/logging/log-capture) &#x2714;
      1. [Configure Custom Logger and View in Kibana ](operator/logging/custom-logs) &#x2714;
      1. [Push Logs to Your Elasticsearch Instance](operator/logging/own-elasticsearch) &#x2714;
   1. [Metrics (12.2.1.4.X only)](operator/metrics)
      1. [Deploy the Operator with Prometheus Enabled and View in Grafana](operator/metrics/enable-metrics)  &#x2726;
      1. [Enable SSL for Metrics](operator/metrics/ssl) &#x2726;
      1. [Scrape Metrics from Your Prometheus Instance](operator/metrics/own-prometheus) &#x2726;
   1. [Scaling a Coherence Deployment](operator/scaling) &#x2714;
   1. [Change Image Version for Coherence or aApplication Container Using Rolling Upgrade](operator/rolling-upgrade) &#x2714;
1. [Coherence Deployments](coherence-deployments)
   1. [Add Application JARs/Config to a Coherence Deployment](coherence-deployments/sidecar) &#x2714;
   1. [Accessing Coherence via Coherence*Extend](coherence-deployments/extend)
      1. [Access Coherence via the Default Proxy Port](coherence-deployments/extend/default) &#x2714;
      1. [Access Coherence via the Separate Proxy Tier](coherence-deployments/extend/proxy-tier) &#x2714;
      1. [Enabling SSL for Proxy Servers](coherence-deployments/extend/ssl) &#x2714;     
      1. [Using multiple Coherence*Extend Proxies](coherence-deployments/extend/multiple) &#x2714;
   1. [Accessing Coherence via storage-disabled Clients](coherence-deployments/storage-disabled)
      1. [Storage-disabled Client in Cluster via Interceptor](coherence-deployments/storage-disabled/interceptor) &#x2714;
      1. [Storage-disabled Client in Cluster as Separate User image](coherence-deployments/storage-disabled/other) &#x2714;
   1. [Federation  (12.2.1.4.X only)](coherence-deployments/federation)
      1. [Within a Single Kubernetes Cluster](coherence-deployments/federation/within-cluster) &#x2726;
      1. [Across Separate Kubernets Clusters](coherence-deployments/federation/across-clusters) &#x2726;
   1. [Persistence](coherence-deployments/persistence)
      1. [Use Default Persistent Volume Claim](coherence-deployments/persistence/default) &#x2714;
      1. [Use a Specific Persistent Volume](coherence-deployments/persistence/pvc) &#x2714;
   1. [Elastic Data](coherence-deployments/elastic-data)
      1. [Deploy Using Default FlashJournal Locations](coherence-deployments/elastic-data/default) &#x2714;
      1. [Deploy Using External Volume Mapped to the Host](coherence-deployments/elastic-data/external) &#x2714;
   1. [Installing Multiple Coherence Clusters with One Operator](coherence-deployments/multiple-clusters)    
1. [Management](management)
   1. [Using Management over REST (12.2.1.4.X only)](management/rest)
      1. [Access Management over REST](management/rest/standard) &#x2726;
      1. [Access Management over REST Using JVisualVM plugin](management/rest/jvisualvm) &#x2726;
      1. [Enable SSL with Management over REST](management/rest/ssl) &#x2726;
      1. [Modify Writable MBeans](management/rest/mbeans) &#x2726;
   1. [Access JMX in the Coherence Cluster via JConsole and JVisualVM](management/jmx) &#x2714;
   1. [Access Coherence Console and CohQL on a Cluster Node](management/console-cohql) &#x2714;
   1. [Diagnostic Tools](management/diagnostics)
      1. [Produce and Extract a Heap Dump](management/diagnostics/heap-dump) &#x2714;
      1. [Produce and Extract a Java Flight Recorder (JFR) file](management/diagnostics/jfr) &#x2726;
   1. [Manage and Use the Reporter](management/reporter) &#x2726;
   1. [Provide Arguments to the JVM that Runs Coherence](management/jvmarguments) &#x2714;      

# Troubleshooting Tips

## Coherence Cluster pods never reach ready "1/1"

Use the following `kubectl` command to see the message from the pod:

```bash
$ kubectl describe pod pod-name -n sample-coherence-ns
```

## Error: ImagePullBackOff after installing Operator or coherence

When you see `Error: ImagePullBackOff` for one of the pod status,
examine the pod using `kubectl describe pod -n sample-coherence-ns pod-name` to determine the image causing the issue.

Ensure that you have set the following to a valid secret:
```bash
--set imagePullSecrets=sample-coherence-secret
```

## Error: configmaps "coherence-internal-config" not found

If your pods don't start and the `kubectl describe` command shows the error, ensure that you have included the `--targetNamespaces` option when installing the `coherence-operator`.

```bash
Error: configmaps "coherence-internal-config" not found
```

## Unable to delete pods when using log capture

If you are using Kubernetes version older than 1.13.0, you cannot delete pods when you have enabled log capture feature. This is a known issue (fluentd) and you need to add the options `--force --grace-period=0` to force delete the pods.

Refer to [https://github.com/kubernetes/kubernetes/issues/45688](https://github.com/kubernetes/kubernetes/issues/45688).

## Receive Error 'no matches for kind "Prometheus"'

If you have enabled metrics, you can get the following error when you try to install the operator:
```console
Error: validation failed: [unable to recognize "": no matches for kind "Prometheus" in
version "monitoring.coreos.com/v1", unable to recognize "": no matches for kind "PrometheusRule" in
version "monitoring.coreos.com/v1", unable to recognize "": no
...
```
Ensure that you have set the following option in the operator installation:

```bash
   --set prometheusoperator.prometheusOperator.createCustomResource=true
```
This is required only when you install Prometheus for the first time in the namespace.

# Accessing UI endpoints

## Access Grafana

> **Note:** Use of Prometheus and Grafana is available only when using the
operator with Oracle Coherence 12.2.1.4 version.

If you have enabled Prometheus, then you can use the `port-forward-grafana.sh` script in the [common](common) directory to view metrics.

1. Start the port-forward

   ```bash
   $ ./port-forward-grafana.sh sample-coherence-ns
   ```
   ```console
   Forwarding from 127.0.0.1:3000 -> 3000
   Forwarding from [::1]:3000 -> 3000
   ```

2. Access Grafana using the following URL:

   [http://127.0.0.1:3000/d/coh-main/coherence-dashboard-main](http://127.0.0.1:3000/d/coh-main/coherence-dashboard-main)

   * Username: admin  

   * Password: prom-operator

## Access Kibana

If you have enabled log capture, then you can use the `port-forward-kibana.sh` script in the [common](common) directory to view metrics.

1. Start the port-forward

   ```bash
   $ ./port-forward-kibana.sh sample-coherence-ns
   ```
   ```console
   Forwarding from 127.0.0.1:5601 -> 5601
   Forwarding from [::1]:5601 -> 5601
   ```
2. Access Kibana using the following URL:

   [http://127.0.0.1:5601/](http://127.0.0.1:5601/)

## Access Prometheus

> **Note:** Use of Prometheus and Grafana is available only when using the
operator with Oracle Coherence 12.2.1.4 version.

If you have enabled Prometheus, then you can use the `port-forward-prometheus.sh` script in the [common](common) directory to view metrics directly.

1. Start the port-forward

   ```bash
   $ ./port-forward-prometheus.sh sample-coherence-ns
    ```

   ```console
   Forwarding from 127.0.0.1:9090 -> 9090
   Forwarding from [::1]:9090 -> 9090
   ```
2. Access Prometheus using the following URL:

   [http://127.0.0.1:9090/](http://127.0.0.1:9090/)

# Run Samples Integration Tests

Refer to [Developer Guide](developer.md) for more information about how to run the samples integration tests.
