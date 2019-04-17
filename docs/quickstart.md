# Quick Start guide

The Oracle Coherence Kubernetes Operator manages Oracle Coherence on Kubernetes.
It manages monitoring data through Prometheus, and logging data through
ElasticSearch and Kibana.

> Note, use of Prometheus and Grafana is only available when using the
> operator with Coherence 12.2.1.4.

Use this quick start guide to deploy Coherence applications in a
Kubernetes cluster managed by the Coherence Operator. Please note that
this walk-through is for demonstration purposes only, not for use in
production.  These instructions assume that you are already familiar with
Kubernetes and Helm.  If you need to learn more about these two
important and complimentary technologies, please refer to the
[Kubernetes](https://kubernetes.io/docs/home/?path=users&persona=app-developer&level=foundational) and 
[Helm](https://helm.sh/docs/) documentation.

## More Advanced Actions

For more advanced actions, such as accessing Kibana for viewing logs,
see the [User Guide](user-guide.md).

> If you have an old version of the operator installed on your cluster
> you must remove it before installing any of the charts by using the
> `helm delete --purge` command.

## Prerequisites

### Software and Version Prerequisites

* Kubernetes 1.10.3 or above
* Helm 2.12.3 or above

### Runtime Environment Prerequisites

* Kubernetes must be able to pull the docker images required by the
  Coherence Kubernetes Operator.
  
* Some of the Helm charts in this project require additional,
  non-standard, configuration on the each of the Kubernetes PODs that
  will be running the workloads related to the chart.  This
  configuration currently includes:
  
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

### Set up Helm repository for Coherence

In the absence of a Helm repository, execute the steps [Developer
Guide](developer.md) and derive the correct value for `HELM_PREFIX`.

Whenever `HELM_PREFIX` appears in this document or the
[user-guide](user-guide.md), replace it with the value from the
[Developer Guide](developer.md).

Whenever `OPERATOR_VERSION` appears in this document or the
[user-guide](user-guide.md), replace it with the value from the
[Developer Guide](developer.md).


## 2. Install the Coherence Operator

##### 1. Use Helm to install the Coherence Operator

You may like to customize the value of the of the `--name` and
(optional) `--namespace` arguments to `kubectl` for the operator.  You
may also like to customize `targetNamespaces` which the operator manages
and `imagePullSecrets` (if it is necessary).

```
$ helm --debug install --version OPERATOR_VERSION HELM_PREFIX/coherence-operator \
    --name sample-coherence-operator \
    --set "targetNamespaces={}" \
    --set imagePullSecrets=sample-coherence-secret
``` 

**Note**: Remove the `--debug` option if you do not want very verbose
output.  Please consult the `values.yaml` in the chart for important
information regarding the `--set targetNamespaces` argument.

If the operation completes successfully, you should see output similar
to the following.

```
NOTES:
1. Get the application URLs by running these commands:

  export POD_NAME=$(kubectl get pods -l "app=coherence-operator,release=sample-coherence-operator" -o jsonpath="{.items[0].metadata.name}")

  To forward a local port to the Pod http port run:

      kubectl port-forward $POD_NAME 8000:8000

  then access the http endpoint at http://127.0.0.1:8000
```

##### 2. Verify that the Helm install is running by checking Helm status:

```
$ helm status sample-coherence-operator
```

If the deployment was succesfull, the output should include output
similar to the following (abbreviated):

```
LAST DEPLOYED: Thu Feb  7 14:11:17 2019
STATUS: DEPLOYED

[...]
```

## 3. Use Helm to install Coherence

Install the `coherence` helm chart.  You may want to customize the values
for the `--name`, `--namespace` and `imagePullSecrets` options.

```
$ helm --debug install --version OPERATOR_VERSION HELM_PREFIX/coherence \
    --name sample-coherence \
    --set imagePullSecrets=sample-coherence-secret
``` 

If the operation completes successfully, you should see output similar
to the following.

```
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

```
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
the cluster created has three nodes and exposes a Coherence*Extend on
port 20000 in the cluster.  As the chart notes explain, you must
"forward" this port so that it is available outside the cluster.  The
`kubectl` command can do this, but you must supply the name of the
Kubernetes "pod" that is running Coherence.  Thankfully, the chart notes
say exactly how to do this:

```
export POD_NAME=$(kubectl get pods -l "app=coherence,release=sample-coherence" -o jsonpath="{.items[0].metadata.name}")
kubectl port-forward $POD_NAME 20000:20000
```

The first command queries the Kubernetes cluster to get the name of the
first of the three nodes in the Coherence cluster.  It may be something
like `sample-coherence-65f558c987-5bdxr`.  The second command tells
Kubernetes to map port 20000 inside the cluster to port 20000 outside
the cluster.  With this information, it is possible to write a
Coherence*Extend client config that accesses the cluster.  Save the
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

```
$ javac -cp .:${COHERENCE_HOME}/lib/coherence.jar HelloCoherence.java
$ java -cp .:${COHERENCE_HOME}/lib/coherence.jar -Dcoherence.cacheconfig=$PWD/example-client-config.xml -Dcoherence.log.level=-1 HelloCoherence
```

This should produce output similar to the following:

```
Oracle Coherence Version 12.2.1.4.0 Build 73407
 Grid Edition: Development mode
Copyright (c) 2000, 2019, Oracle and/or its affiliates. All rights reserved.

The value of the key is 1
```

Running the program again should produce:

```
Oracle Coherence Version 12.2.1.4.0 Build 73407
 Grid Edition: Development mode
Copyright (c) 2000, 2019, Oracle and/or its affiliates. All rights reserved.

The value of the key is 2
```

## 5. Use Helm to delete Coherence and the Operator

Remove the `coherence` release:

```
$ helm delete --purge sample-coherence sample-coherence-operator
```
