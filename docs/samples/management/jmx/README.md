# Access JMX in the Coherence Cluster via JConsole and JVisualVM

Java Management Extensions (JMX) is the standard way to inspect and
manage enterprise Java applications.  Applications that expose
themselves via JMX incur no runtime performance penalty for doing so,
unless a tool is actively connected to the JMX connection, and only then
in certain cases.  [The Java Tutorials](https://docs.oracle.com/javase/tutorial/) provide
[introduction to JMX](https://docs.oracle.com/javase/tutorial/jmx/index.html).  

Once familiar with JMX, [the Coherence documentation](https://docs.oracle.com/middleware/12213/coherence/COHMG/toc.htm)
has complete coverage of [how to use JMX with Coherence](https://docs.oracle.com/middleware/12213/coherence/COHMG/using-jmx-manage-oracle-coherence.htm#COHMG239).

All of the capabilities of JMX with Coherence are also present with the
operator. 

This sample shows how to connecting to a Coherence JMX MBean Server when using the Coherence Operator. You use the `--set store.jmx.enabled=true`
option which will create an MBean Server Pod from which you can connect to.

By default there will be one replica for the MBean Server. You can create more MBean server pods by setting 
the `store.jmx.replicas` value, e.g. `--set store.jmx.replicas=2`.

See [Here](../../management/rest/) for information on connecting to Management over REST endpoint.

[Return to Management samples](../) / [Return to samples](../../README.md#list-of-samples)

## Prerequisites

1. Install the Coherence Operator

   Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

1. Download the JMXMP connector jar

   The JMX endpoint does not use RMI, it uses JMXMP. This requires an additional jar on the classpath
   of the Java JMX client (i.e. VisualVM, JConsole, etc). You can use curl to download the required JAR.

   ```bash
   curl http://central.maven.org/maven2/org/glassfish/external/opendmk_jmxremote_optional_jar/1.0-b01-ea/opendmk_jmxremote_optional_jar-1.0-b01-ea.jar \
        -o opendmk_jmxremote_optional_jar-1.0-b01-ea.jar
   ```     
  
   This also can be downloaded as a Maven dependency if you are connecting via a Maven built project.
  
   ```xml
   <dependency>
      <groupId>org.glassfish.external</groupId>
      <artifactId>opendmk_jmxremote_optional_jar</artifactId>
      <version>1.0-b01-ea</version>
   </dependency>
   ```

## Installation Steps

1. Install the Coherence cluster

   Issue the following to install the cluster with 1 MBean Server Pod:

   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=coherence-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=false \
      --set logCaptureEnabled=false \
      --set store.jmx.enabled=true \
      --set store.jmx.replicas=1 \
      coherence-community/coherence
   ```
   
   *Note*: There are many other `store.jmx.*` options which control other aspects of the MBean Server node.
   Refer to the [Coherence Operator values.yaml file](https://github.com/oracle/coherence-operator/blob/master/operator/src/main/helm/coherence/values.yaml)
   for more information.
   
   After the chart is installed, instructions are displayed to help you utilize this feature.
   You can follow these instructions or use the commands below:
   
1. Check the MBean Server Pod is running

    ```bash
    $ kubectl get pods -n sample-coherence-ns
    NAME                                    READY   STATUS    RESTARTS   AGE
    coherence-operator-5899f6444b-tckm4     1/1     Running   0          1h
    storage-coherence-0                     1/1     Running   0          29m
    storage-coherence-1                     1/1     Running   0          28m
    storage-coherence-2                     1/1     Running   0          27m
    storage-coherence-jmx-54f5d779d-svh29   1/1     Running   0          29m
    ```   
    
    You should see a pod prefixed with `storage-coherence-jmx` in the above output.
   
1. Port-forward the MBean Server Pod   
   
   ```bash
   $ export POD_NAME=$(kubectl get pods --namespace sample-coherence-ns -l "app=coherence,release=storage,component=coherenceJMXPod" -o jsonpath="{.items[0].metadata.name}")

   $ kubectl --namespace sample-coherence-ns port-forward $POD_NAME 9099:9099
   ```
   
   Access the JMX endpoint at the URL `service:jmx:jmxmp://127.0.0.1:9099` 
   
1. (Optionally) Add data to a cache

   Connect to the Coherence `console` using the following to create a cache.
   
   *Note*: If you do not carry out this step, then you will not see any `CacheMBeans` as described below.

   ```bash
   $ kubectl exec -it --namespace sample-coherence-ns storage-coherence-0 bash /scripts/startCoherence.sh console
   ```   
   
   At the `Map (?):` prompt, type `cache test`.  This will create a cache in the service `PartitionedCache`.
   
   Use the following to add 100,000 objects of size 1024 bytes, starting at index 0 and using batches of 100.
   
   ```bash
   bulkput 100000 1024 0 100

   Mon Apr 15 07:37:03 GMT 2019: adding 100000 items (starting with #0) each 1024 bytes ...
   Mon Apr 15 07:37:15 GMT 2019: done putting (11578ms, 8878KB/sec, 8637 items/sec)
   ```
   
   At the prompt, type `size` and it should show 100000.
   
   Then type `bye` to exit the `console`.   

1. Access via JConsole
    
   Ensure you run JConsole with the JMXMP connector on the classpath:

   ```bash
   $ jconsole -J-Djava.class.path="$JAVA_HOME/lib/jconsole.jar:$JAVA_HOME/lib/tools.jar:opendmk_jmxremote_optional_jar-1.0-b01-ea.jar" service:jmx:jmxmp://127.0.0.1:9099
   ```
   
   Select the `MBeans` tab and then `Coherence Cluster` attributes. You should see the Coherence MBeans as below:
   
   ![JConsole](img/jconsole.png)
   
1. Access via JVisualVM 
   
   Ensure you run JVisualVM with the JMXMP connector on the classpath:

   ```bash
   $ jvisualvm -cp "$JAVA_HOME/lib/tools.jar:opendmk_jmxremote_optional_jar-1.0-b01-ea.jar" 
   ```
    
   *Note*: If you have downloaded JVisualVM separatley, then the executable will be `visualvm`.
   
   From the `Applications` tab, right-click on `Local` and then `Add JMX Connection`.
   
   Under `Connection` enter `service:jmx:jmxmp://127.0.0.1:9099` and click `OK`.
   
   This will add the JMX connection. Double-click the new connection and when this opens
   you should be able to see the `Coherence` MBeans under the `Mbeans` tab. If you have the Coherence
   JVisualVM Plugin installed you should also see a `Coherence` tab as below.
   
  ![JVisualVM](img/jvisualvm.png)
  
  Please refer to the [Coherence MBean Reference](https://docs.oracle.com/middleware/12213/coherence/COHMG/oracle-coherence-mbeans-reference.htm#COHMG5442)
  for detailed information on Coherence MBeans.

## Uninstalling the Charts

Carry out the following commands to delete the chart installed in this sample.

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are gone from previous samples.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.
 
   
   
   
   
