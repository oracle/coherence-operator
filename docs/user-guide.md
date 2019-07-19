# User Guide

The User Guide provides detailed information about how to install
Coherence and the Coherence Operator in your Kubernetes cluster, and how to use the operator to manage Coherence clusters.

- [Guide to this Document](#guide-to-this-document)
- [Before You Begin](#before-you-begin)
- [Common Coherence Tasks](#common-coherence-tasks)
  * [Provide Configuration Files and Application Classes to the Coherence Cluster within Kubernetes](#provide-configuration-files-and-application-classes-to-the-coherence-cluster-within-kubernetes)
    + [Deploy JAR Files](#deploy-jar-files)
      - [Create a JAR File](#create-a-jar-file)
      - [Create a Docker image for the Sidecar](#create-a-docker-image-for-the-sidecar)
      - [Use Helm to Install Coherence](#use-helm-to-install-coherence)
      - [Create the Local Extend Client Configuration](#create-the-local-extend-client-configuration)
      - [Remove Coherence](#remove-coherence)
    + [Deploy Configuration Files](#deploy-configuration-files)
      - [Create the Server Side Cache Configuration](#create-the-server-side-cache-configuration)
      - [Create a Docker Image for the Sidecar](#create-a-docker-image-for-the-sidecar)
      - [Run the client program](#run-the-client-program)
  * [Deploy JAR Containing Application Classes and Configuration files](#deploy-jar-containing-application-classes-and-configuration-files)
  * [Extract Heap Dump Files from a Kubernetes Coherence Pod](#extract-heap-dump-files-from-a-kubernetes-coherence-pod)
  * [Extract Coherence Log Files from Kubernetes](#extract-coherence-log-files-from-kubernetes)
    + [Query Elasticsearch](#query-elasticsearch)
  * [Use Java Management Extensions (JMX) to Inspect and Manage Coherence](#use-java-management-extensions--jmx--to-inspect-and-manage-coherence)
    + [Download the `opendmk_jmxremote_optional_jar` JAR](#download-the--opendmk-jmxremote-optional-jar--jar)
    + [Run VisualVM with the Additional JAR](#run-visualvm-with-the-additional-jar)
    + [Manipulate the VisualVM UI to View the Coherence MBeans](#manipulate-the-visualvm-ui-to-view-the-coherence-mbeans)
  * [Provide Arguments to the JVM that Runs Coherence](#provide-arguments-to-the-jvm-that-runs-coherence)
- [Kubernetes Specific Tasks](#kubernetes-specific-tasks)
  * [Using Helm to Scale the Coherence Deployment](#using-helm-to-scale-the-coherence-deployment)
  * [Perform a Safe Rolling Upgrade](#perform-a-safe-rolling-upgrade)
  * [Deploy Multiple Coherence Clusters](#deploy-multiple-coherence-clusters)
- [Monitoring Performance and Logging](#monitoring-performance-and-logging)
  * [Configuring SSL Endpoints for Management over REST and Metrics Publishing](#configuring-ssl-endpoints-for-management-over-rest-and-metrics-publishing)
  * [Configuring SSL Endpoints for Management over REST](#configuring-ssl-endpoints-for-management-over-rest)
  * [Configuring SSL for Metrics Publishing for Prometheus](#configuring-ssl-for-metrics-publishing-for-prometheus)

# Guide to this Document

The User Guide provides exclusive steps for managing Coherence within Kubernetes. For most of the administrative tasks for managing Kubernetes, refer to [Kubernetes](https://kubernetes.io/docs/home/) Documentation.

The information in this guide is organized into sections that are common and Kubernetes specific tasks. Refer to these sections accordingly for managing Coherence:
* [Common Coherence Tasks](#common-coherence-tasks)
* [Kubernetes Specific Tasks](#kubernetes-specific-tasks)

# Before You Begin

See [Before You Begin](quickstart.md#before-you-begin) section in the Quick Start guide.

All the examples in this guide are installed in a Kubernetes namespace called *sample-coherence-ns*.  To set this namespace as the active namespace, execute the command:

```bash
$ kubectl config set-context $(kubectl config current-context) --namespace=sample-coherence-ns
```

# Common Coherence Tasks

The most common administrative tasks with Coherence are [Overriding the Default Cache Configuration
File](https://docs.oracle.com/middleware/12213/coherence/develop-applications/understanding-configuration.htm#COHDG-GUID-C5335E66-6D7F-4C15-B7EC-F6D7D1494066)
and deploying JARs for [Processing Data in a
Cache](https://docs.oracle.com/middleware/12213/coherence/develop-applications/processing-data-cache.htm).

Most of the administrative tasks to do with Coherence apply when running within Kubernetes.  The [official documentation](https://docs.oracle.com/middleware/12213/coherence/) remains a very useful resource. This section covers a few common scenarios that require special treatment regarding Kubernetes.

## Provide Configuration Files and Application Classes to the Coherence Cluster within Kubernetes

This section explains how to make custom configuration and JAR files
available to your Coherence cluster running in Kubernetes. This approach can be used for any administrative task that requires to make JAR, XML, or other configuration files available to the Coherence cluster.

The Oracle Coherence Operator uses the *sidecar pattern*, as
recommended by [Kubernetes](https://kubernetes.io/docs/concepts/cluster-administration/logging/#sidecar-container-with-a-logging-agent),
to make resources available to Coherence within the Kubernetes cluster.
Docker containers are the most flexible way to allow interaction
with the Coherence cluster running in Kubernetes. The steps to use the sidecar pattern include:

1. Determine the JARs or configuration files that you want to make them available to the cluster.

2. Package the files in a Docker image and deploy that image to a
  Docker registry accessible to the Kubernetes cluster.

3. Install the docker image using the Helm chart.

### Deploy JAR Files

The concept to create a JAR file is derived from [Building Your First Extend Application](https://docs.oracle.com/middleware/12213/coherence/develop-remote-clients/building-your-first-extend-application.htm) in
*Oracle Fusion Middleware Developing Remote Clients for Oracle
Coherence* for use within Kubernetes.

#### Create a JAR File

To create a JAR file:

1. Create a directory for the files.

  ```bash
  $ mkdir -p hello-example/files/lib
  $ cd hello-example
  ```
2. Create a Java program to access the cluster. Save the java file as HelloExample.java in the `hello-example` directory.

  ```java
  import java.io.Serializable;
  import java.text.SimpleDateFormat;
  import java.util.Date;

  import com.tangosol.net.NamedCache;
  import com.tangosol.net.CacheFactory;  

  public class HelloExample {
  public static void main(String[] asArgs) throws Throwable {
    NamedCache<String, Timestamp> cache = CacheFactory.getCache("hello-example");
    Timestamp ts = cache.get("ts1");
    cache.put("ts1",
              ts = new Timestamp((null == ts) ? Long.MIN_VALUE : ts.currentTime));

    System.out.println("The value of the key is " + ts.toString());
  }

  public static class Timestamp implements Serializable {
    public long currentTime;
    public long previousTime;

    public Timestamp(long previousTime) {
      this.currentTime = System.currentTimeMillis();
      this.previousTime = previousTime;
    }

    public String toString() {
      SimpleDateFormat f = new SimpleDateFormat("HH:mm:ss");
      return "Timestamp: previousTime: " + f.format(new Date(previousTime)) +
             " currentTime: " + f.format(new Date(currentTime));
    }
  }
}
```

This program uses a static inner class, `Timestamp`, to store the values in Coherence. Any Java object that is stored in Coherence must
be accessible by Coherence in compiled form. The Java objects are compiled classes in JAR files on the Coherence classpath. Therefore, compile and archive the file:

  ```
  $ javac -cp .:${COHERENCE_HOME}/lib/coherence.jar HelloExample.java
  $ jar -cf files/lib/hello-example.jar *.class
  ```
#### Create a Docker image for the Sidecar

Package the created JAR file within the sidecar Docker image:

1. Create a `Dockerfile` with the following contents.

    ```bash
    FROM oraclelinux:7-slim
    RUN mkdir -p /files/lib
    COPY files/lib/hello-example.jar files/lib
    ```
    Note that the JAR file is placed in the `files/lib` directory
    relative to the root of the Docker image. This is the default
    location where Coherence will look for JAR files to add to the
    classpath. Any JAR files in `files/lib` will be added to the
    classpath. You can change the location where Coherence looks for JARs to add to the classpath.

2. Ensure that the Docker is running on the current host. If not, see [ Get Started](https://docs.docker.com/get-started/) in Docker documentation.

3. Build and tag a Docker image for `hello-example-sidecar`:

    ```bash
    $ docker build -t "hello-example-sidecar:1.0.0-SNAPSHOT" .
    ```

    The trailing dot "." in the command refers to run the build relative to the current directory.
4. Push the created Docker image to the Docker registry which the Kubernetes cluster can reach:

   See [Quick Start](quickstart.md#obtain-images-from-oracle-container-registry) guide to learn how to make the Kubernetes cluster aware of the Docker credentials so it
   can pull down images.

   > **Note:** If you are using a local Kubernetes, you can omit this step, since the Kubernetes pulls from the same Docker server as the one to which the local build command built the image.

#### Use Helm to Install Coherence

Install Coherence using the Helm chart with the details of the sidecar image arguments:

```bash
$ helm --debug install coherence/coherence --name hello-example \
     --set userArtifacts.image=hello-example-sidecar:1.0.0-SNAPSHOT \
     --set imagePullSecrets=sample-coherence-secret
```

> **Note:** If the JAR files are in a different location within the sidecar Docker image, use the `--set userArtifacts.libDir=<absolute path within docker image>` argument to `helm install` to configure the correct location.

In a new terminal window, set up a Kubernetes port forward to expose the Extend port so that your local client can use it.

```bash
$ export POD_NAME=$(kubectl get pods --namespace default -l \
   "app=coherence,release=hello-example" -o jsonpath="{.items[0].metadata.name}")
$ kubectl --namespace default port-forward $POD_NAME 20000:20000
```

This prints the following output and blocks the shell:

```bash
Forwarding from 127.0.0.1:20000 -> 20000
Forwarding from [::1]:20000 -> 20000
```

#### Create the Local Extend Client Configuration

A local client configuration is necessary because the local client connects to the service through Coherence*Extend. Create a file
named `hello-client-config.xml` with the following contents:

```
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

Run the client using the following command:

```bash
$ java -cp files/lib/hello-example.jar:${COHERENCE_HOME}/lib/coherence.jar \
       -Dcoherence.log.level=1 -Dcoherence.distributed.localstorage=false \
       -Dcoherence.cacheconfig=$PWD/hello-client-config.xml HelloExample
```

An output similar to the following is displayed:

```
The value of the key is Timestamp: previousTime: 11:47:04 currentTime: 16:09:30
```

Run the command again and it will show the updated `Timestamp`:

```
The value of the key is Timestamp: previousTime: 16:09:30 currentTime: 16:10:20
```

#### Remove Coherence

```
$ helm delete --purge hello-example
```

### Deploy Configuration Files

The similar sidecar approach is used to deploy configuration files to Coherence inside Kubernetes. Though, Coherence has the necessary built-in configuration, a subset of that configuration is used in this example.

1. Create a directory for the files.

  ```bash
  $ cd ..
  $ mkdir -p hello-config-example/files/conf
  $ cd hello-config-example
  ```
2. Create the Java program to access the cluster. In the same directory, create a simple java program `HelloConfigXml.java`.

  ```java
  import com.tangosol.net.CacheFactory;
  import com.tangosol.net.NamedCache;

  public class HelloConfigXml {
    public static void main(String[] asArgs) throws Throwable {
    NamedCache<String, Integer> cache = CacheFactory.getCache("hello-config-xml");
    Integer IValue = (Integer) cache.get("key");
    IValue = (null == IValue) ? Integer.valueOf(1) : Integer.valueOf(IValue + 1);
    cache.put("key", IValue);
    System.out.println("The value of the key is " + IValue);
  }
}
  ```
3. Create the following XML file, next to the java file, called
`hello-client-config.xml`:

  ```
  <cache-config xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
   xmlns="http://xmlns.oracle.com/coherence/coherence-cache-config"
   xsi:schemaLocation="http://xmlns.oracle.com/coherence/coherence-cache-config
   coherence-cache-config.xsd">
   <caching-scheme-mapping>
      <cache-mapping>
         <cache-name>hello-config-xml</cache-name>
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

#### Create the Server Side Cache Configuration

Create the file `files/conf/hello-server-config.xml` with the following
content:

```
<cache-config xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
              xmlns="http://xmlns.oracle.com/coherence/coherence-cache-config"
              xsi:schemaLocation="http://xmlns.oracle.com/coherence/coherence-cache-config coherence-cache-config.xsd">
  <caching-scheme-mapping>
    <cache-mapping>
      <cache-name>hello-config-xml</cache-name>
      <scheme-name>${coherence.profile near}-${coherence.client direct}</scheme-name>
    </cache-mapping>
  </caching-scheme-mapping>

  <caching-schemes>
    <!-- near caching scheme for clustered clients -->
    <near-scheme>
      <scheme-name>near-direct</scheme-name>
      <front-scheme>
        <local-scheme>
          <high-units>{front-limit-entries 10000}</high-units>
        </local-scheme>
      </front-scheme>
      <back-scheme>
        <distributed-scheme>
          <scheme-ref>thin-direct</scheme-ref>
        </distributed-scheme>
      </back-scheme>
    </near-scheme>

    <!-- near caching scheme for extend clients -->
    <near-scheme>
      <scheme-name>near-remote</scheme-name>
      <scheme-ref>near-direct</scheme-ref>
      <back-scheme>
        <remote-cache-scheme>
          <scheme-ref>thin-remote</scheme-ref>
        </remote-cache-scheme>
      </back-scheme>
    </near-scheme>

    <!-- remote caching scheme for accessing the proxy from extend clients -->
    <remote-cache-scheme>
      <scheme-name>thin-remote</scheme-name>
      <service-name>RemoteCache</service-name>
      <proxy-service-name>Proxy</proxy-service-name>
    </remote-cache-scheme>

    <!-- partitioned caching scheme for clustered clients -->
    <distributed-scheme>
      <scheme-name>thin-direct</scheme-name>
      <scheme-ref>server</scheme-ref>
      <local-storage system-property="coherence.distributed.localstorage">false</local-storage>
      <autostart>false</autostart>
    </distributed-scheme>

    <!-- partitioned caching scheme for servers -->
    <distributed-scheme>
      <scheme-name>server</scheme-name>
      <service-name>PartitionedCache</service-name>
      <local-storage system-property="coherence.distributed.localstorage">true</local-storage>
      <backing-map-scheme>
        <local-scheme>
          <high-units>{back-limit-bytes 0B}</high-units>
        </local-scheme>
      </backing-map-scheme>
      <autostart>true</autostart>
    </distributed-scheme>

    <!-- proxy scheme that allows extend clients to connect to the cluster over TCP/IP -->
    <proxy-scheme>
      <service-name>Proxy</service-name>
      <acceptor-config>
        <tcp-acceptor>
          <local-address>
            <address system-property="coherence.extend.address"/>
            <port system-property="coherence.extend.port"/>
          </local-address>
        </tcp-acceptor>
      </acceptor-config>
      <load-balancer>client</load-balancer>
      <autostart>true</autostart>
    </proxy-scheme>
  </caching-schemes>
</cache-config>
```
#### Create a Docker Image for the Sidecar

Package the XML file within the sidecar Docker image:

1. Create a `Dockerfile` next to the java file, with the following contents.

   ```bash
   FROM oraclelinux:7-slim
   RUN mkdir -p /files/conf
   COPY files/conf/hello-server-config.xml files/conf/hello-server-config.xml
   ```
   Note that the XML file is placed in the `files/conf` directory
   relative to the root of the Docker image.  This is the default
   location where Coherence will look for configuration files that apply to Coherence.  You can change the location where Coherence
   looks for configuration files to add to the classpath.

2. Ensure Docker is running on current host. If not, refer to [Get Started with Docker](https://docs.docker.com/get-started/).

3. Build and tag a Docker image for `hello-server-config-sidecar`:

    ```bash
    $ docker build -t "hello-server-config-sidecar:1.0.0-SNAPSHOT" .
    ```
  Note that the trailing dot "." is very significant.  It means, "run the build relative to the current directory."

4. Push your image to the Docker registry which the Kubernetes cluster
   can reach.  See [Quick Start](./quickstart.md) guide to learn how
   to make the Kubernetes cluster aware of the Docker credentials so it
   can pull down images.

   > **Note:** If you are using a local Kubernetes, you can omit this step, since the Kubernetes pulls from the same Docker server as the one to which the local build command built the image.

5. Install the Helm chart, passing the arguments with the sidecar image:

```bash
$ helm --debug install coherence/coherence --name hello-server-config \
     --set userArtifacts.image=hello-server-config-sidecar:1.0.0-SNAPSHOT \
     --set imagePullSecrets=sample-coherence-secret \
     --set store.cacheConfig=hello-server-config.xml
```

> **Note:** If your XML files are in a different location within the sidecar Docker image, use the `--set userArtifacts.configDir=<absolute path within docker image>` argument to `helm install` to configure the correct location.

In a separate shell, set up a Kubernetes port forward to expose the
Extend port so that your local client can use it.  The instructions for
doing this are output from the above `helm install` command, but they
are repeated here for your convenience.

```bash
$ export POD_NAME=$(kubectl get pods --namespace default -l "app=coherence,release=hello-server-config" -o jsonpath="{.items[0].metadata.name}")
$ kubectl --namespace default port-forward $POD_NAME 20000:20000
```

This prints the following output and blocks the shell:

```bash
Forwarding from 127.0.0.1:20000 -> 20000
Forwarding from [::1]:20000 -> 20000
```

#### Run the client program

In the same directory as the XML and Java source files, run the client:

```bash
$ javac  -cp .:${COHERENCE_HOME}/lib/coherence.jar HelloConfigXml.java
$ java -cp .:${COHERENCE_HOME}/lib/coherence.jar \
       -Dcoherence.distributed.localstorage=false \
       -Dcoherence.cacheconfig=$PWD/hello-client-config.xml -Dcoherence.log.level=1 HelloConfigXml
```
The the correct `coherence.jar` must be available at`${COHERENCE_HOME}/lib/coherence.jar`. An output similar to the following s displayed:

```bash
The value of the key is 1
```

Run the program again and it shows that the value has been incremented.

```bash
The value of the key is 2.
```

## Deploy JAR Containing Application Classes and Configuration files

You can deploy a JAR that contains both application classes and configuration files. The sidecar image contains one or more JAR files, each of which can contain application classes, configuration files, or both.
JAR files included in the sidecar image will be available on the Coherence classpath and all Java classes in those JAR files will be available for Classloading by the entire Coherence cluster. The configuration files must be included in the top
level of a JAR file so that it can be referenced by the Coherence helm
chart. An example of the sidecar image layout:

```bash
files/
   lib/
      coherence-operator-hello-server-config-1.0.0-SNAPSHOT.jar
```

The file layout in the JAR file must be:

```bash
META-INF/
META-INF/LICENSE
META-INF/beans.xml
META-INF/maven/org.javassist/javassist/pom.xml
META-INF/services/org.glassfish.jersey.server.spi.ContainerProvider
com/foo/demo/model/Price.class
cache-config.xml
pof-config.xml
```

In the following example, Coherence is installed with the configuration files `cache-config.xml` and `pof-config.xml`, and the JAR file with all the Java classes in the Coherence classpath.

```
$ helm --debug install coherence/coherence --name hello-server-config \
     --set userArtifacts.image=coherence-operator-hello-server-config:1.0.0-SNAPSHOT \
     --set imagePullSecrets=sample-coherence-secret \
     --set store.cacheConfig=cache-config.xml \
     --set store.pof.config=pof-config.xml
```

## Extract Heap Dump Files from a Kubernetes Coherence Pod

You can use the operator to create log files and JVM heap dump files to debug issues in Coherence applications.

See more about [Debugging in
Coherence](https://docs.oracle.com/middleware/12213/coherence/develop-applications/debugging-coherence.htm) in *Oracle Fusion Middleware Developing Applications with Oracle Coherence*.

In this example, a `.hprof` file is collected for a heap dump:

1. Execute the following command to list the pods of the installed operator and Coherence:

  ```bash
$ kubectl get pods
NAME                                 READY     STATUS    RESTARTS   AGE
coherence-demo-storage-0             1/1       Running   0          45m
coherence-demo-storage-1             1/1       Running   0          44m
coherence-operator-7bc94cfb4-g4kz2   1/1       Running   0          47m
```
2. Get a shell to the storage node:

  ```bash
$ kubectl exec -it coherence-demo-storage-0 -- /bin/bash
```
3. Obtain the process ID (PID) of the Coherence process. Use `jps` to get the actual PID:

  ```bash
bash-4.2# /usr/java/default/bin/jps
1 DefaultCacheServer
4230 Jps
```
4. Use the `jcmd` command to extract the heap dump and exit the shell:

  ```bash
bash-4.2# /usr/java/default/bin/jcmd 1 GC.heap_dump /DefaultCache.hprof
bash-4.2# exit
```
5. Use `kubectl exec` to extract the heap dump:

  ```bash
$ (kubectl exec coherence-demo-storage-0 -it -- cat /DefaultCache.hprof ) > DefaultCache.hprof
```

You can also extract the heap dump using the following single command which can be used repeatedly:

```bash
$ (kubectl exec coherence-demo-storage-0 -- /bin/bash -c "rm -f /tmp/heap.hprof; /usr/java/default/bin/jcmd 1 GC.heap_dump /tmp/heap.hprof; cat /tmp/heap.hprof > /dev/stderr" ) 2> heap.hprof
```

```bash
1:
Heap dump file created
```

In this command, the PID of the Coherence is assumed to be `1`. Also, the heap dump output is redirected to `stderr` to prevent the unsuppressable output from `jcmd` from showing up in the heap dump file.


## Extract Coherence Log Files from Kubernetes

When you install the operator and Coherence with the feature log capture enabled, all the log messages from each Coherence cluster are captured to Elasticsearch, stored with Fluentd, and analyzed in Kibana.
The common practice is to capture all the log messages from all of the Coherence clusters into a log aggregator than examining individual Coherence cluster node log files for errors.
Refer to the [sample] for installing the operator and Coherence with log capture enabled.

With log capture feature enabled, all the log messages from every Coherence cluster member are captured including the cluster members that are not running. The persistence of the stored log messages depends on how Fluentd is configured and is beyond the scope of this documentation.  

You can reconstruct the log message of each Coherence cluster member by querying Elasticsearch using `curl`, and manipulating the result using [jq](https://stedolan.github.io/jq/) to produce output equivalent to a regular Coherence log file.

### Query Elasticsearch

To reach the Elasticsearch endpoint using `curl`, you can use port forwarding in Kubernetes to access the Elasticsearch pod. Use the `kubectl port-forward` with the Elasticsearch name and the port number (default 9200).

Use the following `curl` command to capture the log message for each Coherence cluster member:

```bash
curl -s --output coherence-0.json http://ES_HOST:ES_PORT/coherence-cluster-*/_search?size=9999&q=host%3A%22my-20190514-storage-coherence-1%22&sort=@timestamp
```
In the command, `ES_HOST:ES_PORT` is the Elasticsearch pod name and port number in which it can reached. This command extracts the log message from the Coherence cluster member named `my-storage-coherence-0` and stores in the `coherence-0.json` file.

Reformat the `coherence-0.json` file using `jq`:

```
jq -j '.["hits"] | .["hits"] | .[] | .["_source"] | .["@timestamp"]," ", .["product"]," <", .["level"], "> (thread=", .["thread"], ", member=", .["member"], "):", .["log"], "newline"' coherence-0.json | sed -e $'s/newline/\\\n/g' > coherence-0.log
```

## Use Java Management Extensions (JMX) to Inspect and Manage Coherence

JMX is the standard way to inspect and manage enterprise Java applications. Applications that expose themselves through JMX does not incur runtime performance penalty unless a tool is actively connected to the JMX connection. The [Java Tutorials](https://docs.oracle.com/javase/tutorial/) provide [Introduction to JMX](https://docs.oracle.com/javase/tutorial/jmx/index.html). The section [JMX with Coherence](https://docs.oracle.com/middleware/12213/coherence/COHMG/using-jmx-manage-oracle-coherence.htm#COHMG239) in *Oracle Fusion Middleware Developing Applications with Oracle Coherence* describes how to use JMX with Coherence.

The Coherence Helm chart must be installed with additional arguments so that you can use the JMX feature in the operator. This section covers how to install Coherence in a Kubernetes cluster with JMX enabled.

Note that to fully appreciate this use-case, deploy an application that
uses Coherence and creates some caches.  Such an application can be
installed using the steps detailed in [Deploy JAR Files](#deploy-jar-files).

See the [Quick Start](quickstart.md#install-the-operator) guide to install the operator. Install Coherence using Helm chart with the following additional argument for JMX `--set store.jmx.enabled=true*`:

```bash
$ helm --debug install coherence/coherence --name hello-example \
     --set userArtifacts.image=coherence-demo-app:1.0 \
     --set store.jmx.enabled=true \
     --set imagePullSecrets=sample-coherence-secret
```

After Coherence installation, you must expose the network port for JMX using the `kubectl port-forward` command.

The instructions will also include suggestions on how to use JConsole or
[VisualVM](https://visualvm.github.io/).  For the sake of completeness,
this use-case documents how to use VisualVM to access and manipulate
Coherence MBeans when running within Kubernetes.

### Download the `opendmk_jmxremote_optional_jar` JAR

The JMX endpoint does not use RMI, instead it uses JMXMP. This requires an
additional JAR on the classpath of the Java JMX client (VisualVM, or
JConsole). This can be downloaded as a Maven dependency:

```xml
<dependency>
    <groupId>org.glassfish.external</groupId>
    <artifactId>opendmk_jmxremote_optional_jar</artifactId>
    <version>1.0-b01-ea</version>
</dependency>
```
or directly from:

    http://central.maven.org/maven2/org/glassfish/external/opendmk_jmxremote_optional_jar/1.0-b01-ea/opendmk_jmxremote_optional_jar-1.0-b01-ea.jar

### Run VisualVM with the Additional JAR

Download and start VisualVM 1.4.2 version to enable connection to Coherence in Kubernetes:

```bash
visualvm --jdkhome ${JAVA_HOME} --cp:a PATH_TO_DOWNLOADED.jar
```

### Manipulate the VisualVM UI to View the Coherence MBeans

1. In the **File** menu, choose **Add JMX Connection**. In the **Connection** field, enter the value that was output by the Helm chart instructions. For example, `service:jmx:jmxmp://127.0.0.1:9099`. Click **OK**.

2. In the left navigation pane, click **Applications**.  Double click the **service:jmx:jmxmp...** link.

3. Click **MBeans** to open the MBeans browser.

4. In the tree view on the left, open Coherence > Cache >
  DistributedCache and keep drilling down , until you can find a cache
  created by your application.

  Here you can see the MBeans in https://docs.oracle.com/middleware/12213/coherence/COHMG/oracle-coherence-mbeans-reference.htm#GUID-A443DF50-F151-4E9B-AFC9-DFEDF4B149E7__CHDFJDAC

  In particular `HighUnits`, which defaults to 0.  This can be
  interactively changed in `visualvm`.

5. Expand the tree view to Coherence > Node and pick one of the nodes.

  Here you can see the MBeans in https://docs.oracle.com/middleware/12213/coherence/COHMG/oracle-coherence-mbeans-reference.htm#GUID-0AB8710B-2A1D-432D-AFBF-8E73B8230D51__CHDBIJFA

  In particular `LoggingLevel`, which defaults to 5. This can be
  also be interactively changed.

  Note that any changes to MBean attributes done in this way will not
  persist when the cluster restarts.  To make persistent changes, you
  must modify the Coherence configuration files.

## Provide Arguments to the JVM that Runs Coherence

Any production enterprise Java application must carefully tune the JVM
arguments for maximum performance, and Coherence is no exception.  This
use-case explains how to convey JVM arguments to Coherence running
inside Kubernetes.

This use-case is covered [in the samples](samples/management/jvmarguments/).

Please see [the Coherence Performance Tuning
documentation](https://docs.oracle.com/middleware/12213/coherence/administer/performance-tuning.htm#GUID-2A0BC9E6-C3AA-4012-B3D8-EC51963B0CEB)
for authoritative information on this topic.

There are several values in the
[values.yaml](https://github.com/oracle/coherence-operator/blob/master/operator/src/main/helm/coherence/values.yaml)
file of the Coherence Helm chart that convey JVM arguments to the JVM
that runs Coherence within Kubernetes. Please see the source code for
the authoritative documentation on these values.  Such values include
the following.

| `--set` left hand side | Meaning |
|------------------------|---------|
| `store.maxHeap`        | Heap size arguments to the JVM. The format should be the same as that used for Java's -Xms and -Xmx JVM options. If not set the JVM defaults are used. |
| `store.jmx.maxHeap` | Heap size arguments passed to the MBean server JVM.  Same format and meaning as the preceding row. |
| `store.jvmArgs` | Options passed directly to the JVM running Coherence within Kubernetes |
| `store.javaOpts` | Miscellaneous JVM options to pass to the Coherence store container |

The following invocation installs and starts Coherence with specific
values to be passed to the JVM.

```bash
$ helm --debug install coherence/coherence --name hello-example \
     --set store.maxHeap="8g" \
     --set store.jvmArgs="-Xloggc:/tmp/gc-log -server -Xcomp" \
     --set store.javaOpts="-Dcoherence.log.level=6 -Djava.net.preferIPv4Stack=true" \
     --set userArtifacts.image=hello-example-sidecar:1.0.0-SNAPSHOT \
     --set imagePullSecrets=sample-coherence-secret
```

The JVM arguments will include the `store.` arguments specified above,
in addition to many others required by the operator and Coherence.

```bash
-Xloggc:/tmp/gc-log -server -Xcomp -Xms8g -Xmx8g -Dcoherence.log.level=6 -Djava.net.preferIPv4Stack=true
```

To inspect the full JVM arguments, you can use `kubectl get logs -f <pod-name>` and search for one of the arguments you specified.

# Kubernetes Specific Tasks


## Using Helm to Scale the Coherence Deployment

The Coherence Operator leverages Kubernetes `Statefulsets` to ensure that the scale up and scale down operations allow the underlying Coherence
cluster nodes sufficient time to rebalance the cluster data.

Use the following command to scale up the number of Coherence cluster nodes:
```bash
$ kubectl scale statefulsets <helm_release_name> --replicas=<number_of_nodes>

```
For example, to increase the number of Coherence cluster nodes from 2 to 4 for a Coherence Operator installed using Helm chart and the Helm release named `coherence_deploy`:
```bash
$ kubectl scale statefulsets coherence-deploy --replicas=4
```
You can monitor the progress of the cluster as Kubernetes adjusts to the new configuration. Kubernetes shows the number of pods being adjusted and the status of each pod. The pods state change to `Running` status. The coherence cluster also completes rebalancing before the increment.

To scale down the number of Coherence cluster nodes:

```bash
$ kubectl scale statefulsets coherence-deploy --replicas=3
```
The Coherence cluster rebalances the decrease in the number of nodes and the number of pods are adjusted accordingly.

## Perform a Safe Rolling Upgrade

The sidecar docker image created using the JAR file containing application classes is tagged with a version number. The version number enables safe rolling upgrades. See Helm documentation for more information about safe rolling upgrades. The safe rolling upgrade allow you to instruct Kubernetes to replace the currently deployed version of the application classes with a different one.  Coherence and Kubernetes ensures that this upgrade is done without any data loss or interruption of service.

Assuming that you have installed the sidecar docker image using the steps detailed in the <procedure link> and the new sidecar docker image is available for the upgrade, you can use the following command to upgrade:

```bash
$ helm --debug upgrade coherence/coherence --name hello-example --reuse-values \
     --set userArtifacts.image=hello-example-sidecar:1.0.1 --wait \
     --set imagePullSecrets=sample-coherence-secret
```
In this example, the `hello-example-sidecar:1.0.1` is the upgrade destination tagged image. The operator upgrades the application from `hello-example-sidecar:1.0.0-SNAPSHOT` to
`hello-example-sidecar:1.0.1`.

## Deploy Multiple Coherence Clusters

The operator is designed to be installed once on a given
Kubernetes cluster. This one [Helm release](https://helm.sh/docs/glossary/#release) of the Coherence Operator can monitor and manage all of the Coherence clusters installed on the given Kubernetes cluster.

The following commands install the Coherence operator, then install multiple independent Coherence clusters on the same Kubernetes cluster. All the clusters are managed by one operator.

First, install the Coherence Operator with an empty list for the `targetNamespaces` parameter. This causes the operator to manage
all namespaces for Coherence clusters.

```bash
$ helm --debug install coherence/coherence-operator \
    --name sample-coherence-operator \
    --set "targetNamespaces={}" \
    --set imagePullSecrets=sample-coherence-secret
```
Then, install two independent clusters which differ in the values
passed to the `cluster` and `userArtifacts.image` parameters, and
`--name` option.

```bash
$ helm --debug install coherence/coherence \
     --set cluster=revenue-management \
     --set imagePullSecrets=sample-coherence-secret \
     --set userArtifacts.image=revenue-app:2.0.1 \
     --name revenue-management

$ helm --debug install coherence/coherence \
     --set cluster=charging \
     --set imagePullSecrets=sample-coherence-secret \
     --set userArtifacts.image=charging-app:2.0.1 \
     --name charging
```
The values must be unique to ensure that the two Coherence clusters to not merge and form one cluster.

>**Note**: Use the command `helm inspect readme <chart name>` to print the `README.md` of the chart. For example `helm inspect readme coherence/coherence-operator` prints the `README.md` for the operator chart. This includes documentation on all the possible values that can be configured with `--set` options to `helm`.

# Monitoring Performance and Logging

See the following guides for monitoring services and viewing logs:

* [Monitoring Coherence services via Grafana dashboards](prometheusoperator.md)
* [Accessing the EFK stack for viewing logs](logcapture.md)

> **Note**: Use of Prometheus and Grafana is available only when using the operator with Oracle Coherence 12.2.1.4.

## Configuring SSL Endpoints for Management over REST and Metrics Publishing

This section describes how to configure SSL for management over REST and Prometheus metrics:
* Configurating SSL Endpoints for Management over REST
* Configuring SSL for Metrics Publishing for Prometheus

> **Note:** SSL and Management over REST and metrics publishing are available in Oracle Coherence 12.2.1.4.

## Configuring SSL Endpoints for Management over REST

This section describes how to configure a two way SSL for Coherence management over REST with an example:

1. Create Kubernetes secrets for your key store and trust store files. Coherence SSL requires Java key store and trust store files. These files are password protected. Let's define the password protected key store and trust store files:

  ```
  keyStore - name of the Java keystore file: myKeystore.jks
  keyStorePasswordFile - name of the keystore password file: storepassword.txt
  keyPasswordFile - name of the key password file: keypassword.txt
  trustStore - name of the Java trust store file: myTruststore.jks
  trustStorePasswordFile - name of the trust store password file: trustpassword.txt
  ```
  The following command creates a Kubernetes secret named `ssl-secret` which contains the Java key store and trust store files:

  ```bash
  kubectl create secret generic ssl-secret \
     --namespace myNamespace \
     --from-file=./myKeystore.jks \
     --from-file=./myTruststore.jks \
     --from-file=./storepassword.txt \
     --from-file=./keypassword.txt \
     --from-file=./trustpassword.txt
  ```
2. Create a YAML file, `helm-values-ssl-management.yaml`, to enable SSL for Coherence management over REST using the keystore, trust store, and password files in the `ssl-secret`:

  ```yaml     
       store:
         management:
           ssl:
             enabled: true
             secrets: ssl-secret
             keyStore: myKeystore.jks
             keyStorePasswordFile: storepassword.txt
             keyPasswordFile: keypassword.txt
             keyStoreType: JKS
             trustStore: myTruststore.jks
             trustStorePasswordFile: trustpassword.txt
             trustStoreType: JKS
             requireClientCert: true

         readinessProbe:
           initialDelaySeconds: 10  
  ```
3. Install the Coherence Helm chart using the YAML file `helm-values-ssl-management.yaml`:

  ```bash
  helm install coherence/coherence \
    --name coherence \
    --namespace myNamespace \
    --set imagePullSecrets=my-imagePull-secret \
    -f helm-values-ssl-management.yaml
  ```
To verify that the Coherence management over REST is running with HTTPS, forward the management listen port to your local machine: and

  ```bash
  $ kubectl port-forward <pod name> 30000:30000
  ```
Access the REST endpoint using the URL `https://localhost:30000/management/coherence/cluster`.

SSL certificates are required to access sites using the HTTPS protocol. If you have self-signed certificate, you will get the message that your connection is not secure from the browser. Click **Advanced** and then **Add Exception...** to allow the request to access the URL.

Also, look for the following message in the log file of the Coherence pod: <br />
  `Started: HttpAcceptor{Name=Proxy:ManagementHttpProxy:HttpAcceptor, State=(SERVICE_STARTED), HttpServer=NettyHttpServer{Protocol=HTTPS, AuthMethod=cert}`

## Configuring SSL for Metrics Publishing for Prometheus

To configure a SSL endpoint for Coherence metrics:

  1. Create Kubernetes secrets for your key store and trust store files. Coherence SSL requires Java key store and trust store files. These files are password protected. Let's define the password protected key store and trust store files:

    ```bash
    keyStore - name of the Java keystore file: myKeystore.jks
    keyStorePasswordFile - name of the keystore password file: storepassword.txt
    keyPasswordFile - name of the key password file: keypassword.txt
    trustStore - name of the Java trust store file: myTruststore.jks
    trustStorePasswordFile - name of the trust store password file: trustpassword.txt
    ```
    The following command creates a Kubernetes secret named `ssl-secret` which contains the Java key store and trust store files:

    ```bash
    kubectl create secret generic ssl-secret \
       --namespace myNamespace \
       --from-file=./myKeystore.jks \
       --from-file=./myTruststore.jks \
       --from-file=./storepassword.txt \
       --from-file=./keypassword.txt \
       --from-file=./trustpassword.txt
    ```
2. Create a YAML file, `helm-values-ssl-metrics.yaml`, using the keystore, trust store, and password file stored in `ssl-secret`:

  ```yaml   
     store:
       metrics:
         ssl:
           enabled: true
           secrets: ssl-secret
           keyStore: myKeystore.jks
           keyStorePasswordFile: storepassword.txt
           keyPasswordFile: keypassword.txt
           keyStoreType: JKS
           trustStore: myTruststore.jks
           trustStorePasswordFile: trustpassword.txt
           trustStoreType: JKS
           requireClientCert: true
           readinessProbe:
           initialDelaySeconds: 10
  ```        
3. Install the Coherence helm chart using the YAML file, `helm-values-ssl-metrics.yaml`:

  ```bash
  helm install coherence/coherence \
    --name coherence \
    --namespace myNamespace \
    --set imagePullSecrets=my-imagePull-secret \
    -f helm-values-ssl-metrics.yaml
  ```

To verify that the Coherence metrics for Prometheus is running with HTTPS, forward the Coherence metrics port and access the metrics from your local machine use the following commands:

  ```bash
  $ kubectl port-forward <Coherence pod> 9612:9612

  $ curl -X GET https://localhost:9612/metrics --cacert <caCert> --cert <certificate>
  ```

You can add `--insecure` if you use self-signed certificate. Also, look for the following message in the log file of the Coherence pod:</br>
  `Started: HttpAcceptor{Name=Proxy:MetricsHttpProxy:HttpAcceptor, State=(SERVICE_STARTED), HttpServer=NettyHttpServer{Protocol=HTTPS, AuthMethod=cert}`

To configure Prometheus SSL (TLS) connections with the Coherence metrics SSL endpoints, see https://github.com/helm/charts/blob/master/stable/prometheus-operator/README.md for more information about how to specify Kubernetes secrets that contain the certificates required for two-way SSL in Prometheus.
To configure Prometheus to use SSL (TLS) connections, see https://prometheus.io/docs/prometheus/latest/configuration/configuration/#tls_config.

After configuring Prometheus to use SSL, verify that the Prometheus is scraping Coherence metrics over HTTPS by forwarding the Prometheus service port to your local machine and access the following URL:

```bash
  $ kubectl port-forward <Prometheus pod> 9090:9090

  http://localhost:9090/graph
  ```
You should see many vendor:coherence_* metrics.   

To enable SSL for both management over REST and metrics publishing for Prometheus, install the Coherence chart with both YAML files:

  ```bash
    helm --debug install coherence/coherence \
      --name coherence \
      --namespace myNamespace \
      --set imagePullSecrets=my-imagePull-secret \
      -f helm-values-ssl-management.yaml,helm-values-ssl-metrics.yaml
```
