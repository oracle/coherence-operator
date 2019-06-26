# Quick Start Guide

The Coherence Kubernetes Operator manages Coherence through Kubernetes,
monitoring MBean attributes through Prometheus and server logs through
Elasticsearch and Kibana.

> **Note**: Use of Prometheus and Grafana is available only when using the
> operator with Oracle Coherence 12.2.1.4.

Use this quick start guide to deploy Coherence applications in a
Kubernetes cluster managed by the Coherence Kubernetes Operator. Please
note that this quick start guide is for illustrative purposes only, and
not sufficiently prescriptive or thorough for a production environment.
These instructions assume that you are already familiar with Kubernetes
and Helm.  If you need to learn more about these two important and
complimentary technologies, please refer to the
[Kubernetes](https://kubernetes.io/docs/home/?path=users&persona=app-developer&level=foundational)
and [Helm](https://helm.sh/docs/) documentation.

## More Advanced Actions

For more advanced actions, such as accessing Kibana for viewing server
logs, see the [User Guide](user-guide.md).

> **Note**: If you have an old version of the operator installed on your cluster
> you must remove it before installing any of the charts by using the
> `helm delete --purge` command.

## Prerequisites

### Software and Version Prerequisites

* Kubernetes 1.11.5+, 1.12.3+, 1.13.0+ (check with `kubectl version`)
* Docker 18.03.1-ce (check with `docker version`)
* Flannel networking v0.10.0-amd64 (check with `docker images | grep flannel`)
* Helm 2.12.3 or above (and all of its prerequisites)
* Oracle Coherence 12.2.1.3.2

### Runtime Environment Prerequisites

* Kubernetes must be able to pull the docker images required by the
  Coherence Operator.

* Some of the Helm charts in this project require, configuration on each
  Kubernetes pod that will be running the workloads related to the
  chart.  This configuration currently includes:

    * setting the value of the `max_map_count` kernel parameter to at
      least `262144`.  The manner for doing this is OS specific and is
      described
      [in the docker documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html#docker-cli-run-prod-mode).

* If running the operator with a local Kubernetes cluster on a developer
  workstation, ensure the workstation meets the reasonable hardware
  requirements and the Docker preferences have been tuned to run
  optimally for the hardware.  In particular, ensure memory and disk
  limits have been correctly tuned.

## 1. Environment Configuration

### Add the Helm repository for Coherence

Issue the following to create a `coherence` helm repository:

```bash
$ helm repo add coherence https://oracle.github.io/coherence-operator/charts

$ helm repo update

Hang tight while we grab the latest from your chart repositories...
...Skip local chart repository
...Successfully got an update from the "coherence" chart repository
```

> **Note**: For all helm install commands you can leave the `--version`
> option off and the latest version of the chart will be retrieved.  If
> you wanted to use a specific version, such as `0.9.8`, add `--version
> 0.9.8` to all installs for the `coherence-operator` and `coherence`
> charts.

If you wish to build the Coherence Operator from source, please refer to the
[Developer Guide](developer.md) and ensure you replace `coherence`
Helm repository prefix in all samples with the full qualified directory as described at the end of the guide.

## 2. Use Helm to install the Coherence Operator

You may like to customize the value of the of the `--name` and
(optional) `--namespace` arguments to `helm` for the operator.  You
may also like to customize `targetNamespaces` which the operator manages
and `imagePullSecrets` (if it is necessary).

```bash
$ helm --debug install coherence/coherence-operator \
    --name sample-coherence-operator \
    --set "targetNamespaces={}" \
    --set imagePullSecrets=sample-coherence-secret
```

> **Note**: Remove the `--debug` option if you do not want very verbose
> output.  Please consult the `values.yaml` in the chart for important
> information regarding the `--set targetNamespaces` argument.

If the operation completes successfully, you should see output similar to the following.

```bash
NOTES:
1. Get the application URLs by running these commands:

  export POD_NAME=$(kubectl get pods -l "app=coherence-operator,release=sample-coherence-operator" -o jsonpath="{.items[0].metadata.name}")

  To forward a local port to the Pod http port run:

      kubectl port-forward $POD_NAME 8000:8000

  then access the http endpoint at http://127.0.0.1:8000
```

Use `helm ls` to view the installed releases.

```bash
$ helm ls

NAME                     	REVISION	UPDATED                 	STATUS  	CHART                   	APP VERSION	NAMESPACE
sample-coherence-operator	1       	Thu May  9 13:59:22 2019	DEPLOYED	coherence-operator-0.9.8	0.9.8      	default  
```

You can also query the status with `helm status`:

```bash
$ helm status sample-coherence-operator
```

If the deployment was successful, the output should include output
similar to the following (abbreviated):

```bash
LAST DEPLOYED: Thu Feb  7 14:11:17 2019
STATUS: DEPLOYED

[...]
```

## 3. Use Helm to install Coherence

By default the Oracle Coherence Docker image pulled by the Coherence Helm
chart is from the Oracle Container Registry.

To be able to pull Coherence Docker Images from the Oracle Container Registry:

a) Login to [Oracle Container Registry](https://container-registry.oracle.com)
   and accept the terms and conditions to download Coherence images:

   > 1. Go to to [Oracle Container Registry](https://container-registry.oracle.com)
   > 2. Search for "Coherence".
   > 3. Select `coherence` from the list.
   > 4. Click on `Sign-in` on the right and enter your credentials, or create and account if you don't already have one.
   > 5. On the right, select the language for the  `Oracle Standard Terms and Restrictions`.
   > 6. Click `Continue` and scroll down to accept the terms and conditions.

b) Create Kubernetes docker-registry secret with the same credentials that is
   used in step (a) to login into Oracle Container Registry and tell Kubernetes
   to use that secret when pulling the image.

```bash
$ kubectl create secret docker-registry oracle-container-registry-secret \
    --docker-server=container-registry.oracle.com \
    --docker-username='<USERNAME>' --docker-password='<PASSWORD>'
```

In the above command, replace &lt;USERNAME&gt; and &lt;PASSWORD&gt; with the actual
username and password to authenticate with container-registry.oracle.com

Install the `coherence` helm chart using the imagePullSecrets value
'oracle-container-registry-secret' which was just created in the previous
mentioned `kubectl create secret` command.

You may want to customize the value for the `--name` option that you want
to use for your coherence cluster.

```bash
$ helm --debug install coherence/coherence \
    --name sample-coherence \
    --set imagePullSecrets=oracle-container-registry-secret
```

> **Note**: If you want to use a different version of Coherence than the
> one specified in the `coherence` helm chart, supply a `--set` argument
> for the `coherence.image` value, as shown next.
> `--set coherence.image="<prefix>/coherence:<version>"`
> Use the command `helm inspect readme <chart name>` to print out the
> `README.md` of the chart.  For example `helm inspect readme
> coherence/coherence` will print out the `README.md` for the operator
> chart.  This includes documentation on all the possible values that
> can be configured with `--set` options to `helm`.  In particular, look
> at the *Configuration* section of the `README.md`.

If the operation completes successfully, you should see output similar
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

You can also query the status with `helm status`:

```bash
$ helm status sample-coherence
LAST DEPLOYED: Wed Feb 13 14:51:38 2019
STATUS: DEPLOYED
[...]
```

Running `helm install` creates what is called a "helm release".  See the
[glossary](https://docs.helm.sh/glossary/) for the formal definition of
"release".

## 4. Access the Coherence running within Kubernetes using the default Coherence*Extend feature

When starting Coherence with no options, as in the preceding section,
the Coherence cluster created has three nodes, and exposes a
Coherence*Extend proxy server on port 20000 in the cluster.  As the
chart notes explain, you must "forward" this port so that it is
available outside the cluster.  The `kubectl` command can do this, but
you must supply the name of the Kubernetes "pod" that is running
Coherence.  Thankfully, the chart notes say exactly how to do this:

```bash
$ export POD_NAME=$(kubectl get pods -l "app=coherence,release=sample-coherence" -o jsonpath="{.items[0].metadata.name}")
$ kubectl port-forward $POD_NAME 20000:20000
```

The first command queries the Kubernetes cluster to get the name of the
first of the three nodes in the Coherence cluster.  It may be something
like `sample-coherence-65f558c987-5bdxr`.  The second command tells
Kubernetes to map port 20000 inside the cluster to port 20000 outside
the cluster.  With this information, it is possible to write a
Coherence*Extend client configuration to access the cluster.  Save the
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

Finally, we can write a small Java program to interact with the
Coherence cluster.  Save this file next to the preceding XML file, as
`HelloCoherence.java`

```
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

Assuming you are in the same directory as the XML and Java source files,
and that the correct `coherence.jar` is available at
`${COHERENCE_HOME}/lib/coherence.jar`, compile and run the program as shown
next:

```bash
$ javac -cp .:${COHERENCE_HOME}/lib/coherence.jar HelloCoherence.java
$ java -cp .:${COHERENCE_HOME}/lib/coherence.jar \
       -Dcoherence.cacheconfig=$PWD/example-client-config.xml HelloCoherence
```

This should produce output similar to the following:

```bash
The value of the key is 1
```

Running the program again should produce:

```bash
The value of the key is 2
```

> **Note**: If you are using JDK 11 or newer, you can omit the `javac`
> step and simply run the program as shown next.

```bash
$ java -cp $${COHERENCE_HOME}/lib/coherence.jar \
  -Dcoherence.cacheconfig=$PWD/example-client-config.xml  HelloCoherence.java
```

## 5. Use Helm to delete Coherence and the Operator

Remove the `coherence` release:

```bash
$ helm delete --purge sample-coherence sample-coherence-operator
```
