# Quick Start Guide

The Coherence Kubernetes Operator manages Coherence through Kubernetes,
monitoring MBean attributes through Prometheus and server logs through
Elasticsearch and Kibana.

> **Note**: Use of Prometheus and Grafana is available only when using the
> operator with Oracle Coherence 12.2.1.4.

Use this Quick Start guide to deploy Coherence applications in a
Kubernetes cluster managed by the Coherence Kubernetes Operator. Note that this guide is for illustrative purposes only, and
not sufficiently prescriptive or thorough for a production environment.
These instructions assume that you are already familiar with Kubernetes
and Helm.  If you want to learn more about these two  technologies, refer to the
[Kubernetes](https://kubernetes.io/docs/home/?path=users&persona=app-developer&level=foundational)
and [Helm](https://helm.sh/docs/) documentation.

For more advanced actions, such as accessing Kibana for viewing server
logs, see the [User Guide](user-guide.md).

> **Note**: If you have an older version of the operator installed on your cluster,
> you must remove it before installing the current version.

## Before You Begin

### Software Requirements
The Coherence Kubernetes Operator has the following requirements:

* Kubernetes 1.11.5+, 1.12.3+, 1.13.0+ (`kubectl version`)
* Docker 18.03.1-ce (`docker version`)
* Flannel networking v0.10.0-amd64 (`docker images | grep flannel`)
* Helm 2.12.3 or above, and all of its prerequisites
* Oracle Coherence 12.2.1.3.2

### Runtime Environment Requirements

* You will need a Kubernetes cluster that can pull the docker images required by the operator.
* Some of the Helm charts in this operator, require configuration on each Kubernetes pod that will be running the workloads related to the chart. This configuration currently includes:

    * setting the value of the `max_map_count` kernel parameter to at least `262144`. It is OS specific and is described in the [docker documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html#docker-cli-run-prod-mode).

* If you are running the operator with a local Kubernetes cluster on a developer workstation, ensure that the workstation meets the hardware requirements and the Docker preferences have been tuned to run optimally for the hardware.  In particular, ensure that the memory and limits are correctly tuned.

## Configure the environment

### Add the Helm Repository for Coherence

Create a `coherence` helm repository:

```bash
$ helm repo add coherence https://oracle.github.io/coherence-operator/charts

$ helm repo update
```

If you want to build the Coherence Operator from source, refer to the [Developer Guide](developer.md) and ensure that you replace `coherence`
Helm repository prefix in all samples with the full qualified directory as described at the end of the guide.

## Install the Coherence Kubernetes Operator

1. Use `helm` to install the operator:

```bash
$ helm --debug install coherence/coherence-operator \
    --name sample-coherence-operator \
    --set "targetNamespaces={}" \
    --set imagePullSecrets=sample-coherence-secret
```
> **Note**: For all `helm install` commands, you can leave the `--version` option off and the latest version of the chart is retrieved. If you wanted to use a specific version, such as `0.9.8`, add `--version 0.9.8` to all installs for the `coherence-operator` and `coherence` charts.

> **Note**: Remove the `--debug` option if you do not want verbose output. Refer to the `values.yaml` in the chart for more information about the `--set targetNamespaces` argument.

When the operation completes successfully, output similar to the following displays:

```bash
NOTES:
1. Get the application URLs by running these commands:

  export POD_NAME=$(kubectl get pods -l "app=coherence-operator,release=sample-coherence-operator" -o jsonpath="{.items[0].metadata.name}")

  To forward a local port to the Pod http port run:

      kubectl port-forward $POD_NAME 8000:8000

  then access the http endpoint at http://127.0.0.1:8000
```
2. Use `helm ls` to view the installed releases.

```bash
$ helm ls

NAME                     	REVISION	UPDATED                 	STATUS  	CHART                   	APP VERSION	NAMESPACE
sample-coherence-operator	1       	Thu May  9 13:59:22 2019	DEPLOYED	coherence-operator-0.9.8	0.9.8      	default  
```

3. Verify the status of the operator using `helm status`:

```bash
$ helm status sample-coherence-operator
```

A successful deployment displays an output similar to the following output:

```bash
LAST DEPLOYED: Thu Feb  7 14:11:17 2019
STATUS: DEPLOYED

[...]
```

## Install Coherence

### Obtain Images from Oracle Container Registry

By default, the Helm chart pulls the Oracle Coherence Docker image from the Oracle Container Registry.

To pull Coherence Docker images from the Oracle Container Registry:

1. In a web browser, navigate to [Oracle Container Registry](https://container-registry.oracle.com) and click **Sign In**.
2. Enter your Oracle credentials or create an account if you don't have one.
3. Search for coherence in the Search Oracle Container Registry field.
4. Click `coherence` in the search result list.
5. In the Oracle Coherence page, select the language from the drop-down list and click Continue.
6. Click **Accept** in the Oracle Standard Terms and Conditions page.
7. In a terminal window, log in to the Oracle Container Registry:

    `docker login container-registry.com`
8. Pull the coherence image with the command:

  `docker pull container-registry.oracle.com/middleware/coherence:12.2.1.3.2`

### Set Up Secrets to Access the Oracle Container Registry

 Create a Kubernetes secret with the Oracle Container Registry credentials:
 ```bash
 $ kubectl create secret docker-registry oracle-container-registry-secret \
     --docker-server=container-registry.oracle.com \
     --docker-username='<username>' --docker-password='<password>' \
     --docker-email=`<emailid>`
 ```

 When installing Coherence, refer to the secret and Kubernetes will use that secret when pulling the image.

### Use Helm to Install Coherence

Install the `coherence` helm chart using the secret 'oracle-container-registry-secret' which was created using the `kubectl create secret` command.

```bash
$ helm --debug install coherence/coherence \
    --name sample-coherence \
    --set imagePullSecrets=oracle-container-registry-secret
```

>**Note**: If you want to use a different version of Coherence than the one specified in the `coherence` Helm chart, supply a `--set` argument
> for the `coherence.image` value:
>
> `--set coherence.image="<prefix>/coherence:<version>"`
>
> Use the command `helm inspect readme <chart name>` to print out the
> `README.md` of the chart. For example, `helm inspect readme
> coherence/coherence` prints out the `README.md` for the operator
> chart. This includes documentation on all the possible values that
> can be configured with `--set` options to `helm`. Refer to the Configuration section of [README](README.md).

When the operation completes successfully, you see output similar
to the following.

```bash
NOTES:
1. Get the application URLs by running these commands:

  export POD_NAME=$(kubectl get pods -l "app=coherence,release=sample-coherence" -o jsonpath="{.items[0].metadata.name}")

  To forward a local port to the Pod Coherence*Extend port run:

      kubectl port-forward $POD_NAME 20000:20000

  then access the Coherence*Extend endpoint at 127.0.0.1:20000

  To forward a local port to the Pod Http port run:

      kubectl port-forward $POD_NAME 30000:30000

  then access the http endpoint at http://127.0.0.1:30000
```

You can also query the status of the installed Coherence with `helm status`:

```bash
$ helm status sample-coherence
LAST DEPLOYED: Wed Feb 13 14:51:38 2019
STATUS: DEPLOYED
[...]
```
Running `helm install` creates a **helm release**.  See the
[Helm Glossary](https://docs.helm.sh/glossary/) for the  definition of Release.

## Access the Installed Coherence
You can access the installed Coherence running within your Kubernetes using the default Coherence*Extend feature.

When starting Coherence with no options, the created Coherence cluster has three nodes, and exposes a Coherence*Extend proxy server on port 20000 in the cluster. You must forward this port so that it is available outside the cluster. You can use the `kubectl` command by supplying the name of the Kubernetes pod that is running the Coherence.

Query the Kubernetes cluster to get the name of the first of the three nodes in the Coherence cluster.

```bash
$ export POD_NAME=$(kubectl get pods -l "app=coherence,release=sample-coherence" -o jsonpath="{.items[0].metadata.name}")
```
The name of the cluster will be similar to `sample-coherence-65f558c987-5bdxr`. Then, map the port 20000 inside the cluster to port 20000 outside the cluster.

```bash
$ kubectl port-forward $POD_NAME 20000:20000
```

With this information, you can write a Coherence*Extend client configuration to access the cluster. Save the
following XML in a file called `example-client-config.xml`:

```xml
<cache-config xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
   xmlns="http://xmlns.oracle.com/coherence/coherence-cache-config"
   xsi:schemaLocation="http://xmlns.oracle.com/coherence/coherence-cache-config
   coherence-cache-config.xsd">
   <caching-scheme-mapping>
      <cache-mapping>
         <cache-name>*</cache-name>
         <scheme-name>thin-remote</scheme-name>
      </cache-mapping>
   </caching-scheme-mapping>

   <caching-schemes>
       <remote-cache-scheme>
           <scheme-name>thin-remote</scheme-name>
           <service-name>Proxy</service-name>
           <initiator-config>
               <tcp-initiator>
                  <remote-addresses>
                      <socket-address>
                          <address>127.0.0.1</address>
                          <port>20000</port>
                      </socket-address>
                  </remote-addresses>
               </tcp-initiator>
           </initiator-config>
       </remote-cache-scheme>
   </caching-schemes>
</cache-config>
```

Write a simple Java program to interact with the Coherence cluster. Save this file as `HelloCoherence.java`in the same directory as the XML file.

```bash
import com.tangosol.net.CacheFactory;
import com.tangosol.net.NamedCache;

public class HelloCoherence {
  public static void main(String[] asArgs) throws Throwable {
    NamedCache<String, Integer> cache = CacheFactory.getCache("HelloCoherence");
    Integer IValue = (Integer) cache.get("key");
    IValue = (null == IValue) ? Integer.valueOf(1) : Integer.valueOf(IValue + 1);
    cache.put("key", IValue);
    System.out.println("The value of the key is " + IValue);
  }
}
```
With the XML and Java source files in the same directory, and the `coherence.jar` at `${COHERENCE_HOME}/lib/coherence.jar`, compile and run the program:

```bash
$ javac -cp .:${COHERENCE_HOME}/lib/coherence.jar HelloCoherence.java
$ java -cp .:${COHERENCE_HOME}/lib/coherence.jar \
       -Dcoherence.cacheconfig=$PWD/example-client-config.xml HelloCoherence
```
This produces an output similar to the following:

```bash
The value of the key is 1
```

Run the program again and you get an output similar to the following:

```bash
The value of the key is 2
```

> **Note**: If you are using JDK 11 or later, you can omit the javac step and simply run the program as:

```bash
$ java -cp $${COHERENCE_HOME}/lib/coherence.jar \
  -Dcoherence.cacheconfig=$PWD/example-client-config.xml  HelloCoherence.java
```

## Remove Coherence and Coherence Kubernetes Operator

Remove the `coherence` release:

```bash
$ helm delete --purge sample-coherence sample-coherence-operator
```
