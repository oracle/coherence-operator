# Coherence Operator Samples

These samples provide simple demonstrations of how to accomplish common
tasks.  The samples are **not** intended to be used in production
deployments or to be depended upon to create production environments.
They are provided for educational and demonstration purposes only.

While these samples may be useful and usable as is, it is intended that
you would read through all of the sample code in detail, understand how
the given sample works, and then customize it to suit your needs.

> Please see the [Oracle Coherence Demonstration](https://github.com/coherence-community/coherence-demo)
> for an example of an application using Coherence and running via the Coherence Operator.

# Table of Contents

1. [Start Here](#start-here)
   1. [Confirm Quickstart Runtime Prerequisites](#confirm-quickstart-runtime-prerequisites)
   1. [Ensure you have JDK8 and Maven Installed](#ensure-you-have-jdk8-and-maven-installed)
   1. [Ensure you install Coherence into your local Maven repository](#ensure-you-install-coherence-into-your-local-maven-repository)
   1. [Create the sample namespace](#create-the-sample-namespace)
   1. [Clone the GitHub Repository](#clone-the-github-repository)
   1. [Install the Coherence Operator](#install-the-coherence-operator)
1. [List of Samples](#list-of-samples)
1. [Troubleshooting Tips](#troubleshooting-tips)
1. [Accessing UI endpoints](#accessing-ui-endpoints)
   1. [Access Grafana](#access-grafana)
   1. [Access Kibana](#access-kibana)
   1. [Access Prometheus](#access-prometheus)

# Start Here

If you have never setup Coherence Operator before, please carry out the following:

1. [Confirm Quickstart Runtime Prerequisites](#confirm-quickstart-runtime-prerequisites)
1. [Ensure you have JDK8 and Maven Installed](#ensure-you-have-jdk8-and-maven-installed)
1. [Ensure you install Coherence into your local Maven repository](#ensure-you-install-coherence-into-your-local-maven-repository)
1. [Create the sample namespace](#create-the-sample-namespace)
1. [Clone the GitHub Repository](#clone-the-github-repository)
1. [Install the Coherence Operator](#install-the-coherence-operator)

If you have already run samples before, please go to the [List of Samples](#list-of-samples).

Throughout all these samples we are using a Kubernetes namespace called `sample-coherence-ns`.
If you wish to change this namespace,
please ensure you change any references to this namespace to match your selected namespace.

## Confirm Quickstart Runtime Prerequisites

Confirm you have completed the following sections from the `Quick Start Guide` before continuing:

* Runtime Environment Prerequisites
  
  * [Software and Version Prerequisites](../quickstart.md#software-and-version-prerequisites) - Helm and Kubernetes versions
  
  * [Runtime Environment Prerequisites](../quickstart.md#runtime-environment-prerequisites) - Helm & Kubernetes configuration
  
* Environment Configuration

  * [Add the Helm repository for Coherence](../quickstart.md#add-the-helm-repository-for-coherence)
  
  * [Obtain Images from Oracle Container Registry](../quickstart.md#obtain-images-from-oracle-container-registry)
  
  
> **Note:** For all helm install commands you can leave the --version option off 
> and the latest version of the chart will be retrieved. If you wanted 
> to use a specific version, such as 0.9.8, add --version 0.9.8 to all installs for the coherence-operator and coherence charts.  

## Ensure you have JDK8 and Maven Installed

Ensure you have the following installed:

* JDK 8+

* Maven 3.5.4+

> **Note:** You may use a later version of Java, e.g. JDK11 as the
> `maven.compiler.source` and `target` are set to 8 in the sample
> `pom.xml` files.

## Ensure you install Coherence into your local Maven repository

If you are not running samples that have a Maven project, then you can skip this step, otherwise continue on.

1. Download and install Coherence 12.2.1.3 from [Oracle Technology Network](https://www.oracle.com/technetwork/middleware/coherence/downloads/index.html).

1. Make sure COHERENCE_HOME environment variable is set to point to your `coherence` directory under your install location.
   I.e. the directory containing the bin, lib, doc directories. This is only required for the Maven install-file commands.

1. Use the following to install Coherence into your local repository:

   ```bash
   $ mvn install:install-file -Dfile=$COHERENCE_HOME/lib/coherence.jar   \
                              -DpomFile=$COHERENCE_HOME/plugins/maven/com/oracle/coherence/coherence/12.2.1/coherence.12.2.1.pom
   ```
   
   If you are running Coherence 12.2.1.4, you will also need to install `coherence-metrics`.
   
   ```bash
   $ mvn install:install-file -Dfile=$COHERENCE_HOME/lib/coherence-metrics.jar \
                              -DpomFile=$COHERENCE_HOME/plugins/maven/com/oracle/coherence/coherence-metrics/12.2.1/coherence-metrics.12.2.1.pom
   ```   

## Create the sample namespace

You should only need to carry out the following the first time you
run any of the samples.

* Create your target namespace:

  ```bash
  $ kubectl create namespace sample-coherence-ns

  namespace/sample-coherence-ns created
  ```

## Docker Images: Create an imagePullSecret

* Enable Kubernetes to pull docker images from private docker
  repositories by creating a "secret"

  If you must enable your Kubernetes cluster to pull images from private
  repositories, you must create a "secret" to convey the docker
  credentials to Kubernetes. In these samples we are assuming you have
  created a secret called `sample-coherence-secret` in your namespace
  `sample-coherence-ns`.  If all of your images can be pulled from
  public repositories, this step is not required.

  See [https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/) for more information.

  Confirm your secret exists.

  ```bash
  $ kubectl get secret sample-coherence-secret -n sample-coherence-ns
  NAME                      TYPE                             DATA   AGE
  sample-coherence-secret   kubernetes.io/dockerconfigjson   1      18s
  ```

### Using DockerHub For Your Images

It is frequently convenient to pull docker images from one repository,
modify, re-tag, and then push the images to repositories under your own
control, such as repositories you can create on `hub.docker.com`.  Such
repositories are public by default, but can also be private.

#### Providing an imagePullSecret to Kubernetes for a repository on Docker Hub

In this example, we pull the Coherence 12.2.1.3 docker image from the
Docker Store, re-tag it, and push it to a freshly created docker
repository within your Docker Hub account.  We then create an
imagePullSecret to allow Kubernetes to pull that image.  This approach
can also be used when pushing "sidecar images", as several samples
require.

1. Sign in to your accont on [Docker Hub](https://hub.docker.com/).
   Upon login, you are taken to your list of repositories.  For
   discussion, let's say your docker ID is `mydockerid`.

1. Create a new repository by clicking the `Create Repository +` button.
   This takes you to the `Create Repository` page.  You will see a
   dropdown with `mydockerid` followed by a blank field where you are to
   type the "Name".  Type `coherence` for the name.  Type "Re-tags of
   Official Coherence Image" as the description.
   
1. Choose whether to make the repository public or private, then click
   create.  This example assumes private, but either will work.
   
1. Click the `Create` button at the bottom of the page.  Note the value
   of the "To push a new tag to this repository," box.  The value after
   `docker push` will be `mydockerid/coherence:tagname`.  This whole
   value is called the "docker tag"
   
1. Follow the steps in [the developer guide](/docs/developer/#how-to-build-the-operator-without-running-any-testssamples)
   in the section "Obtain a Coherence 12.2.1.3.* Docker image and tag it
   correctly", up to and including the `docker pull` command.
   
1. Re-tag the Coherence 12.2.1.3 docker image with your repository and tag name: 

   `docker tag store/oracle/coherence:12.2.1.3 mydockerid/coherence:12.2.1.3`
   
1. At the command line, login to your docker hub account.

   `docker login`
   
   Enter your docker userid and password as requested.
   
1. Push the re-tagged image to your docker repository.

   `docker push mydockerid/coherence:12.2.1.3`
   
   Note the value of the `The push refers to repository` statement.  The
   first part of that is necessary to create the secret for Kubernetes.
   In this case, it should be something like `docker.io/mydockerid/coherence`.
   
1. Create the secret within your namespace.

   ```
   kubectl create secret docker-registry sample-coherence-secret \
    --namespace sample-coherence-ns --docker-server=hub.docker.com \
    --docker-username=docker.io/mydockerid --docker-password="your docker password" \
    --docker-email="the email address of your docker account"
   ```
   
When invoking `helm` you may specify one or more imagePullSecrets by
enclosing them within a comma separated list inside the curly braces of
a `--set imagePullSecrets` option.

   ```
   --set "imagePullSecrets{sample-coherence-secret}"
   ```

## Clone the GitHub Repository

The samples exist in the `gh-pages` branch of the Coherence Operator GitHub repository - https://github.com/oracle/coherence-operator.

Issue the following to clone the repository and switch to the `gh-pages` branch.

```bash
$ git clone https://github.com/oracle/coherence-operator

$ cd coherence-operator

$ git checkout gh-pages

$ cd docs/samples
```

1. Inspect the [samples top level pom.xml](pom.xml) and verify that the
   value of the `coherence.version` property matches the version of
   Coherence you are actaully using. For example if you have Coherence
   12.2.1.3.0 then the value of `coherence.version` must be
   `12.2.1-3-0`.  If this value needs ajustment, use the
   `-Dcoherence.version=` argument to all invocations of `mvn`.

Issue the following to ensure all the projects with source code build ok. 

> **Note**: Any compilation errors will most likely indicate that the
> Coherence JAR's are not properly installed or you have not set your
> JDK.

```bash
$ mvn clean install
```

## Install the Coherence Operator

Before you attempt any of the samples below, you should install the
`coherence-operator` chart.This can be done once and can keep running
for all the samples.

When you install the `coherence-operator` you can optionally enable the following:

1. Prometheus integration: to capture metrics and display in
   Grafana. (Only available from Coherence 12.2.1.4.0 and above)

1. Log capture: to use Fluentd to send logs to Elasticsearch where they
   can be viewed in Kiabana.

Enabling both Prometheus and log capture will require considerable extra system resources.

> **Note:** when running locally, (e.g. on Docker for Mac), you should allocate sufficient memory
> to you Docker for Mac process. The minimum recommended to run is 8G.

### Install the Coherence Operator (no Prometheus or log capture)

```bash
$ helm install \
   --namespace sample-coherence-ns \
   --set imagePullSecrets=sample-coherence-secret \
   --name coherence-operator \
   --set "targetNamespaces={sample-coherence-ns}" \
   coherence/coherence-operator
```

The above will install the `coherence-operator` without Prometheus or log capture enabled.

### Enabling Prometheus

> **Note:** Use of Prometheus and Grafana is only available when using the
> operator with Coherence 12.2.1.4.

To enable Prometheus, add the following options to the above command:

```bash
   --set prometheusoperator.enabled=true \
   --set prometheusoperator.prometheusOperator.createCustomResource=false
```

Full helm install example for enabling Prometheus is:

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

> **Note:** The first time you install prometheusOperator, you should set the above `createCustomResource=true`. All subsequent `coherence-operator` installs should set this to `false`.

### Enabling log Capture

 To enable log capture, which includes Fluentd, Elasticsearch and Kibana, add the following options to your helm commands:

 ```bash
 --set logCaptureEnabled=true
 ```

 Full helm install example for enabling log capture is:

 ```bash
$ helm install \
    --namespace sample-coherence-ns \
    --set imagePullSecrets=sample-coherence-secret \
    --name coherence-operator \
    --set logCaptureEnabled=true \
    --set "targetNamespaces={sample-coherence-ns}" \
    coherence/coherence-operator
  ```

## Enabling Prometheus and log Capture

You can enable both Prometheus and log capture by setting 
both of the options above to `true`.

## Checking that the Operator is Running

Use `kubectl get pods -n sample-coherence-ns` to ensure that all pods are running.

When enabling Prometheus the following will be displayed.  Depending upon the
number of CPU cores, you may see multiple node-exporter processes.

```bash
NAME                                                     READY   STATUS    RESTARTS   AGE
coherence-operator-66f9bb7b75-nxwdc                      1/1     Running   0          13m
coherence-operator-grafana-898fc8bbd-nls45               3/3     Running   0          13m
coherence-operator-kube-state-metrics-5d5f6855bd-klzj5   1/1     Running   0          13m
coherence-operator-prometh-operator-58bd58ddfd-dhd9q     1/1     Running   0          13m
coherence-operator-prometheus-node-exporter-5hxwh        1/1     Running   0          13m
prometheus-coherence-operator-prometh-prometheus-0       3/3     Running   1          12m
```

When enabling log capture the following will be displayed:

```bash
NAME                                  READY   STATUS    RESTARTS   AGE
coherence-operator-64b4f8f95d-fmz2x   2/2     Running   0          2m
elasticsearch-5b5474865c-tlr44        1/1     Running   0          2m
kibana-f6955c4b9-n8krf                1/1     Running   0          2m
```

# List Of Samples

Samples legend:
* &#x2714; - Available for Coherence 12.2.1.3.x and above

* &#x2726; - Available for Coherence 12.2.1.4.x and above

* &#x2718; - Not yet available

1. [Coherence Operator](operator/)
   1. [Logging](operator/logging)
      1. [Enable log capture to view logs in Kiabana ](operator/logging/log-capture) &#x2714;
      1. [Configure custom logger and view in Kibana ](operator/logging/custom-logs) &#x2714;
      1. [Push logs to your own Elasticsearch Instance](operator/logging/own-elasticsearch) &#x2714;
   1. [Metrics (12.2.1.4.X only)](operator/metrics)
      1. [Deploy the operator with Prometheus enabled and view in Grafana](operator/metrics/enable-metrics)  &#x2726;
      1. [Enable SSL for Metrics](operator/metrics/ssl) &#x2726;
      1. [Scrape metrics from your own Prometheus instance](operator/metrics/own-prometheus) &#x2726;
   1. [Scaling a Coherence deployment via kubectl](operator/scaling) &#x2714;
   1. [Change image version for Coherence or application container using rolling upgrade](operator/rolling-upgrade) &#x2714;
1. [Coherence Deployments](coherence-deployments)
   1. [Add application jars/Config to a Coherence deployment](coherence-deployments/sidecar) &#x2714;
   1. [Accessing Coherence via Coherence*Extend](coherence-deployments/extend)
      1. [Access Coherence via default proxy port](coherence-deployments/extend/default) &#x2714;
      1. [Access Coherence via separate proxy tier](coherence-deployments/extend/proxy-tier) &#x2714;
      1. [Enabling SSL for Proxy Servers](coherence-deployments/extend/ssl) &#x2714;     
      1. [Using multiple Coherence*Extend proxies](coherence-deployments/extend/multiple) &#x2714;
   1. [Accessing Coherence via storage-disabled clients](coherence-deployments/storage-disabled)
      1. [Storage-disabled client in cluster via interceptor](coherence-deployments/storage-disabled/interceptor) &#x2714;
      1. [Storage-disabled client in cluster as separate user image](coherence-deployments/storage-disabled/other) &#x2714;
   1. [Federation  (12.2.1.4.X only)](coherence-deployments/federation)
      1. [Within a single Kubernetes cluster](coherence-deployments/federation/within-cluster) &#x2726;
      1. [Across across separate Kubernets clusters](coherence-deployments/federation/across-clusters) &#x2726;
   1. [Persistence](coherence-deployments/persistence)
      1. [Use default persistent volume claim](coherence-deployments/persistence/default) &#x2714;
      1. [Use a specific persistent volume](coherence-deployments/persistence/pvc) &#x2714;
   1. [Elastic Data](coherence-deployments/elastic-data)
      1. [Deploy using default FlashJournal locations](coherence-deployments/elastic-data/default) &#x2714;
      1. [Deploy using external volume mapped to the host](coherence-deployments/elastic-data/external) &#x2714;
   1. [Installing Multiple Coherence clusters with one Operator](coherence-deployments/multiple-clusters)    
1. [Management](management)
   1. [Using Management over REST (12.2.1.4.X only)](management/rest)
      1. [Access management over REST](management/rest/standard) &#x2726;
      1. [Access management over REST using JVisualVM plugin](management/rest/jvisualvm) &#x2726;
      1. [Enable SSL with management over REST](management/rest/ssl) &#x2726;
      1. [Modify Writable MBeans](management/rest/mbeans) &#x2726;
   1. [Access JMX in the Coherence cluster via JConsole and JVisualVM](management/jmx) &#x2714;
   1. [Access Coherence Console and CohQL on a cluster node](management/console-cohql) &#x2714;
   1. [Diagnostic Tools](management/diagnostics)
      1. [Produce and extract a heap dump](management/diagnostics/heap-dump) &#x2714; 
      1. [Produce and extract a Java Flight Recorder (JFR) file](management/diagnostics/jfr) &#x2726; 
   1. [Manage and use the Reporter](management/reporter) &#x2726;
   1. [Provide arguments to the JVM that runs Coherence](management/jvmarguments) &#x2714;      
 
# Troubleshooting Tips

## Coherence Cluster pods never reach ready "1/1"

Use the following `kubectl` command to see what the message from the pod is:

```bash
$ kubectl describe pod pod-name -n sample-coherence-ns
```

## Error: ImagePullBackOff after installing Operator or coherence

If when you list pods, you see `Error: ImagePullBackOff` for one of the pod status,
examine the pod via `kubectl describe pod -n sample-coherence-ns pod-name` to determine the
image causing the problem.

Ensure you have set the following to a valid secret:
```bash
--set imagePullSecrets=sample-coherence-secret
```

## Error: configmaps "coherence-internal-config" not found

If your pods don't start and the `kubectl describe` command shows the error, then ensure you have included
the `--targetNamespaces` option when installing the `coherence-operator`.

```bash
Error: configmaps "coherence-internal-config" not found
```

## Unable to delete pods when using log capture

If you are using a Kubernetes version below 1.13.0 then you may hit an issue where
you cannot delete pods when you have enabled log capture. This is an known issue (specifically with fluentd)
and you will need to add the options `--force --grace-period=0` to force deletion of the pods.

Refer to [https://github.com/kubernetes/kubernetes/issues/45688](https://github.com/kubernetes/kubernetes/issues/45688).

## Receive Error 'no matches for kind "Prometheus"'

If you are enabling metrics, and you see the following error after attempting to install the
Coherence Operator:

```bash
Error: validation failed: [unable to recognize "": no matches for kind "Prometheus" in 
version "monitoring.coreos.com/v1", unable to recognize "": no matches for kind "PrometheusRule" in 
version "monitoring.coreos.com/v1", unable to recognize "": no
...
```

Ensure that you have set the following, as you may not have installed with prometheus operator enabled before:

```bash
   --set prometheusoperator.prometheusOperator.createCustomResource=true

```

This is only required when you first install prometheus operator into a namespace. 

# Accessing UI endpoints

## Access Grafana

> **Note:** Use of Prometheus and Grafana is only available when using the
> operator with Coherence 12.2.1.4.

If you have enabled Prometheus then you can use the `port-forward-grafana.sh` script in the
[common](common) directory to view metrics.

1. Start the port-forward

   ```bash
   $ ./port-forward-grafana.sh sample-coherence-ns

   Forwarding from 127.0.0.1:3000 -> 3000
   Forwarding from [::1]:3000 -> 3000
   ```

1. Access Grafana using the following URL:

   [http://127.0.0.1:3000/d/coh-main/coherence-dashboard-main](http://127.0.0.1:3000/d/coh-main/coherence-dashboard-main)

   * Username: admin  

   * Password: prom-operator

## Access Kibana

If you have enabled log capture then you can use the `port-forward-kibana.sh` script in the
[common](common) directory to view metrics.

1. Start the port-forward

   ```bash
   $ ./port-forward-kibana.sh sample-coherence-ns

   Forwarding from 127.0.0.1:5601 -> 5601
   Forwarding from [::1]:5601 -> 5601
   ```
1. Access Kibana using the following URL:

   [http://127.0.0.1:5601/](http://127.0.0.1:5601/)

## Access Prometheus


> **Note:** Use of Prometheus and Grafana is only available when using the
> operator with Coherence 12.2.1.4.

If you have enabled Prometheus then you can use the `port-forward-prometheus.sh` script in the
[common](common) directory to view metrics directly.

1. Start the port-forward

   ```bash
   $ ./port-forward-prometheus.sh sample-coherence-ns

   Forwarding from 127.0.0.1:9090 -> 9090
   Forwarding from [::1]:9090 -> 9090
   ```
1. Access Prometheus using the following URL:

   [http://127.0.0.1:9090/](http://127.0.0.1:9090/)

# Running the samples integration tests

Please see [developer.md](developer.md) for instructions how to run the 
samples integration tests.
