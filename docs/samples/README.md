# Coherence Operator Samples

These samples provide simple demonstrations of how to accomplish common
tasks.  The samples are **not** intended to be used in production
deployments or to be depended upon to create production environments.
They are provided for educational and demonstration purposes only.

While these samples may be useful and usable as is, it is intended that
you would read through all of the sample code in detail, understand how
the given sample works, and then customize it to suit your needs.

# Table of Contents

1. [Start Here](#start-here)
   1. [Confirm Quickstart Runtime Prerequisites](#confirm-quickstart-runtime-prerequisites)
   1. [Ensure you have JDK11 and Maven Installed](#ensure-you-have-jdk11-and-maven-installed)
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
1. [Ensure you have JDK11 and Maven Installed](#ensure-you-have-jdk11-and-maven-installed)
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
  
  * [Obtain the Coherence Docker Image](../quickstart.md#obtain-the-coherence-docker-image)
  
  
> Note: For all helm install commands you can leave the --version option off 
> and the latest version of the chart will be retrieved. If you wanted 
> to use a specific version, such as 0.9.3, add --version 0.9.3 to all installs for the coherence-operator and coherence charts.  

## Ensure you have JDK11 and Maven Installed

Ensure you have the following installed:

* JDK 11+
* Maven 3.5.4+

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
   
   If you are using Coherence 12.2.1.4.0 or above you also need to install Coherence REST.
   
   ```bash
   $ mvn install:install-file -Dfile=$COHERENCE_HOME/lib/coherence-rest.jar \
                              -DpomFile=$COHERENCE_HOME/plugins/maven/com/oracle/coherence/coherence-rest/12.2.1/coherence-rest.12.2.1.pom
   ```   

1. Ensure that the [samples top level pom.xml](pom.xml) has the Coherence version set to the version you
   are using:  E.g. if you have Coherence 12.2.1.3.2 then `coherence.version` should be:

   ```xml
    <coherence.version>12.2.1-3-2</coherence.version>
   ```

## Create the sample namespace

You should only need to carry out the following the first time you
run any of the samples.

* Create your target namespace:

  ```bash
  $ kubectl create namespace sample-coherence-ns

  namespace/sample-coherence-ns created
  ```

* Create a secret for pulling images from private repositories

  If you are pulling images from private repositories, you must create a secret
  which will be used for this. In these samples we are assuming you have created a secret called `sample-coherence-secret` in your namespace `sample-coherence-ns`.

  See [https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/) for more information.

  Confirm your secret exists.

  ```bash
  $ kubectl get secret sample-coherence-secret -n sample-coherence-ns
  NAME                      TYPE                             DATA   AGE
  sample-coherence-secret   kubernetes.io/dockerconfigjson   1      18s
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

Issue the following to ensure all the projects with source code build ok. 

> Note: Any compilation errors will most likely indicate that the Coherence JAR's 
> are not properly installed or you have not set your JDK.

```bash
$ mvn clean install
```

## Install the Coherence Operator

Before you attempt any of the samples below, you should install the `coherence-operator` chart.  
This can be done once and can keep running for all the samples.

When you install the `coherence-operator` you can optionally enable the following:

1. Prometheus integration to capture metrics and display in Grafana. (Only available from Coherence 12.2.1.4.0 and above)

1. Log capture which will use Fluentd to send logs to Elasticsearch where
   they key be vied in Kiabana.

Enabling both Prometheus and log capture will require considerable extra system resources.

> Note: when running locally, (e.g. on Docker for Mac), you should allocate sufficient memory
> to you Docker for Mac process. The minimum recommended to run is 8G.

### Install the Coherence Operator (no Prometheus or log capture)

```bash
$ helm install \
   --namespace sample-coherence-ns \
   --set imagePullSecrets=sample-coherence-secret \
   --name coherence-operator \
   --set "targetNamespaces={sample-coherence-ns}" \
   coherence-community/coherence-operator
```

The above will install the `coherence-operator` without Prometheus or log capture enabled.

### Enabling Prometheus

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
     coherence-community/coherence-operator
```

> Note: The first time you install prometheusOperator, you should set the above `createCustomResource=true`. All subsequent `coherence-operator` installs should set this to `false`.

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
    coherence-community/coherence-operator
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
      1. [Include custom user metrics for scraping by Prometheus](operator/metrics/custom-metrics)  &#x2726;
      1. [Enable SSL for Metrics](operator/metrics/ssl) &#x2718;
      1. [Scrape metrics from your own Prometheus instance](operator/metrics/own-prometheus) &#x2726;
   1. [Scaling a Coherence deployment via kubectl](operator/scaling) &#x2714;
   1. [Change image version for Coherence or application container using rolling upgrade](operator/rolling-upgrade) &#x2714;
1. [Coherence Deployments](coherence-deployments)
   1. [Add application jars/Config to a Coherence deployment](coherence-deployments/sidecar) &#x2714;
   1. [Accessing Coherence via Coherence*Extend](coherence-deployments/extend)
      1. [Access Coherence via default proxy port](coherence-deployments/extend/default) &#x2714;
      1. [Access Coherence via separate proxy tier](coherence-deployments/extend/proxy-tier) &#x2714;
      1. [Enabling SSL for Proxy Servers](coherence-deployments/extend/ssl) 
         1. [Enable SSL in Coherence 12.2.1.3.X](coherence-deployments/extend/ssl/12213) &#x2714;
         1. [Enable SSL in Coherence 12.2.1.4.X and above](coherence-deployments/extend/ssl/12214) &#x2718;
      1. [Using multiple Coherence*Extend proxies](coherence-deployments/extend/multiple) &#x2714;
   1. [Accessing Coherence via storage-disabled clients](coherence-deployments/storage-disabled)
      1. [Storage-disabled client in cluster via interceptor](coherence-deployments/storage-disabled/interceptor) &#x2714;
      1. [Storage-disabled client in cluster as separate user image](coherence-deployments/storage-disabled/other) &#x2714;
   1. [Federation  (12.2.1.4.X only)](coherence-deployments/federation)
      1. [Within a single Kubernetes cluster](coherence-deployments/federation/within-cluster) &#x2718;
      1. [Across across separate Kubernets clusters](coherence-deployments/federation/across-clusters) &#x2718;
   1. [Persistence](coherence-deployments/persistence)
      1. [Use default persistent volume claim](coherence-deployments/persistence/default) &#x2714;
      1. [Use a specific persistent volume](coherence-deployments/persistence/pvc) &#x2714;
      1. [Specifying an archiver](coherence-deployments/persistence/archiver) &#x2718;
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
   1. [Manage and use the Reporter](management/reporter) &#x2714; 
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

# Accessing UI endpoints

## Access Grafana

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

If you have enabled log capture then you can use the `port-forward-prometheus.sh` script in the
[common](common) directory to view metrics.

1. Start the port-forward

   ```bash
   $ ./port-forward-prometheus.sh sample-coherence-ns

   Forwarding from 127.0.0.1:9090 -> 9090
   Forwarding from [::1]:9090 -> 9090
   ```
1. Access Prometheus using the following URL:

   [http://127.0.0.1:9090/](http://127.0.0.1:9090/)
