# User guide

This document provides detailed user information for the Oracle
Coherence Kubernetes Operator.  It provides instructions on how to
install the operator in your Kubernetes cluster and how to use it to
manage Coherence Clusters.

For convenience, unless otherwise stated, all the examples in this guide
will be installed in a Kubernetes namespace called
*sample-coherence-ns*.  To set this namespace as the active namespace
execute this command:

```
$ kubectl config set-context $(kubectl config current-context) --namespace=sample-coherence-ns
```

The steps listed in [Environment Configuration in the
quickstart](quickstart.md#1-environment-configuration) must be performed
before any of the steps in this guide.

### Table of Use-Cases

#### Common Coherence Tasks

* [Supply Configuration Files And Application Classes to the Coherence Cluster within Kubernetes](#supply-configuration-files-and-application-classes-to-the-coherence-cluster-within-kubernetes)

   * [Supply a Jar File Containing Application Classes](#first-lets-show-the-simple-example-of-including-a-jar-file)
   
   * [Supply a Configuration File Outside of a Jar File](#now-lets-modify-the-preceding-example-to-deploy-a-config-file)
   
   * [Supply a Configuration File and/or Application Classes In a Jar File](#finally-lets-combine-the-preceding-two-use-cases-and-deploy-a-jar-containing-both-application-classes-and-configuration-files)
   
   * [Extract Reporter Files from Kubernetes](#extract-reporter-files-from-a-kubernetes-coherence-pod)
   
#### Kubernetes Specific Use Cases
   
* [Scale a Coherence Cluster With Helm](#using-helm-to-scale-the-coherence-deployment)

* [Perform a Safe Rolling Upgrade](#perform-a-safe-rolling-upgrade)

* [Deploy Multiple Coherence Clusters Managed by the Operator](#deploy-multiple-coherence-clusters-managed-by-the-operator)

-------------


## Common Coherence Tasks

Most of the administrative tasks one must do with Coherence still apply
when running within Kubernetes.  As such, the [official
documentation](https://docs.oracle.com/middleware/12213/coherence/)
remains a very useful resource.  This section covers a few common
scenarios that require special treatment regarding Kubernetes.

### Use-Cases

#### Supply Configuration Files And Application Classes to the Coherence Cluster within Kubernetes

Two of the most common administrative tasks with Coherence are
[Overriding the Default Cache Configuration
File](https://docs.oracle.com/middleware/12213/coherence/develop-applications/understanding-configuration.htm#COHDG-GUID-C5335E66-6D7F-4C15-B7EC-F6D7D1494066)
and deploying jars for [Processing Data in a
Cache](https://docs.oracle.com/middleware/12213/coherence/develop-applications/processing-data-cache.htm).

This section explains how to make custom configuration and jar files
available to your Coherence Cluster running in Kubernetes.  The same
approach can be used for any administrative task that requires making JAR
files, or XML or other configuration files available to the Coherence
Cluster.

The Oracle Coherence Kubernetes Operator uses the "sidecar pattern", [as
recommended by
Kubernetes](https://kubernetes.io/docs/concepts/cluster-administration/logging/#sidecar-container-with-a-logging-agent),
to make resources available to Coherence within the Kubernetes cluster.
Because Docker containers are the most flexible way to allow interaction
with the Coherence Cluster running in Kubernetes, the general steps for
this usage of the sidecar pattern include:

* Discern what configuration files or jars you want to make available to
  the server.
  
* Package those items in a Docker image and deploy that image to a
  Docker registery accessible to the Kubernetes cluster.

* When installing the Helm chart, refer to the image by name.
  
Here are the concrete steps for performing this action.  These steps
adapt the general idea from *Building Your First Extend Application* in
[Oracle Fusion Middleware Developing Remote Clients for Oracle
Coherence](https://docs.oracle.com/middleware/12213/coherence/develop-remote-clients/building-your-first-extend-application.htm),
for use within Kubernetes.

##### First let's show the simple example of including a jar file. 

##### 1. Create a directory for the files.

```
$ mkdir -p hello-example/files/lib
$ cd hello-example
```

##### 2. Create the Java program that will access the cluster.

In the same directory, create this simple java program, in the file `HelloExample.java`.

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

    System.out.print("The value of the key is " + ts.toString());
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

This program uses a static inner class, `Timestamp` as the values to be
stored in Coherence.  Any Java object that is stored in Coherence must
be accessible by Coherence in compiled form.  This usually means the
Java objects are compiled classes in jar files on the Coherence
classpath.  Therefore, we must compile and jar the program, as shown
here.

```
$ javac -cp .:${COHERENCE_HOME}/lib/coherence.jar HelloExample.java
$ jar -cf files/lib/hello-example.jar *.class
```

##### 3. Create a Docker image for the sidecar

We must now package the jar file within the sidecar Docker image.

1. Create a `Dockerfile` next to the XML and JAR files, with the
   following contents.

    ```
    FROM oraclelinux:7-slim
    RUN mkdir -p /files/lib
    COPY files/lib/hello-example.jar files/lib
    ```

    Note that the jar file is placed in the `files/lib` directory
    relative to the root of the docker image.  This is the default
    location where Coherence will look for jar files to add to the
    classpath.  Any jar files in `files/lib` will be added to the
    classpath.  It is possible to change the location where Coherence
    looks for jars to add to the classpath as shown in the following
    step.

2. Ensure docker is running on current host.  If not, run through [the
   Docker getting-started](https://docs.docker.com/get-started/).
   
3. Build and tag a docker image for your *coherence-operator-hello-example*:

    ```
    $ docker build -t "coherence-operator-hello-example:1.0.0-SNAPSHOT" .
    ```

    Note that the trailing dot "." is very significant.  It means, "run
    the build relative to the current directory."

4. Push your image to the docker registry which the Kubernetes cluster
   can reach.  See [the
   quickstart](quickstart.md#prepare-the-namespace-and-docker-registry-access) to learn how
   to make the Kubernetes cluster aware of the Docker credentials so it
   can pull down images.
   
    If you are using a local Kubernetes, you can omit this step, since
    the Kubernetes pulls from the same Docker server as the one to which
    the local build command built the image.
   
##### 4. Install the Helm chart, passing the arguments to make the chart aware of the sidecar image
   
```
$ helm --debug install --version 1.0.0-SNAPSHOT \
     ./coherence --name hello-example \
     --set userArtifacts.image=coherence-operator-hello-example:1.0.0-SNAPSHOT \
     --set imagePullSecrets=sample-coherence-secret
```

> If your jar files are in a different location within the sidecar Docker
> image, use the `--set userArtifacts.libDir=<absolute path within
> docker image>` argument to `helm install` to configure the correct location.

In a separate shell, set up a Kubernetes "port forward" to expose the
Extend port so that your local client can use it.  The instructions for
doing this are output from the above `helm install` command, but they
are repeated hear for your convenience.

```
$ export POD_NAME=$(kubectl get pods --namespace default -l "app=coherence,release=hello-example" -o jsonpath="{.items[0].metadata.name}")
$ kubectl --namespace default port-forward $POD_NAME 20000:20000
```

This prints the following output and blocks the shell:

```
Forwarding from 127.0.0.1:20000 -> 20000
Forwarding from [::1]:20000 -> 20000
```

##### 5. Create the local extend client config and run the client.

Because the local client will connect to the service via
Coherence*Extend, a local client config is necessary.  Create a file
called `hello-client-config.xml` with the following contents.

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

Run the client with the following command.

```
$ java -cp files/lib/hello-example.jar:${COHERENCE_HOME}/lib/coherence.jar \
  -Dcoherence.log.level=-1 \
  -Dcoherence.cacheconfig=$PWD/hello-client-config.xml HelloExample
```

This should show output similar to the following:

```
Oracle Coherence Version 12.2.1.4.0 Build 73407
 Grid Edition: Development mode
Copyright (c) 2000, 2019, Oracle and/or its affiliates. All rights reserved.

The value of the key is Timestamp: previousTime: 11:47:04 currentTime: 16:09:30
```

Running the command again will show the `Timestamp` being updated:

```
Oracle Coherence Version 12.2.1.4.0 Build 73407
 Grid Edition: Development mode
Copyright (c) 2000, 2019, Oracle and/or its affiliates. All rights reserved.

The value of the key is Timestamp: previousTime: 16:09:30 currentTime: 16:10:20
```

##### 6. Delete the Helm relese

```
$ helm delete --purge hello-example
```

  3. Change the cache configuration that is used to one in the application
     jar.

##### Now let's modify the preceding example to deploy a config file

The same sidecar approach used in the preceding example is also used to
deploy configuration files to Coherence inside Kubernetes.  Coherence
already has the necessary configuration built-in, but for the sake of
the illustration, we will use a subset of that configuration.

##### 1. Create a directory for the files.

```
$ mkdir -p hello-config-example/files/conf
$ cd hello-config-example
```

##### 2. Create the Java program that will access the cluster.

In the same directory, create this simple java program, in the file
`HelloConfigXml.java`.  This is exactly the same java code as in the
[quickstart](quickstart.md).


```
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

Create the following XML file, next to the java file, called
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

This is exactly the same XML as in the [quickstart](quickstart.md).

##### 3. Create the server side cache configuration

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

##### 4. Create a Docker image for the sidecar

We must now package the XML file within the sidecar Docker image.

1. Create a `Dockerfile` next to the java file, with the
   following contents.
   
    ```
    FROM oraclelinux:7-slim
    RUN mkdir -p /files/conf
    COPY files/conf/hello-server-config.xml files/conf/hello-server-config.xml
    ```

    Note that the XML file is placed in the `files/conf` directory
    relative to the root of the docker image.  This is the default
    location where Coherence will look for config files apply to
    Coherence.  It is possible to change the location where Coherence
    looks for jars to add to the classpath as shown in the following
    step.

2. Ensure docker is running on current host.  If not, run through [the
   Docker getting-started](https://docs.docker.com/get-started/).
   
3. Build and tag a docker image for your *coherence-operator-hello-example*:

    ```
    $ docker build -t "coherence-operator-hello-server-config:1.0.0-SNAPSHOT" .
    ```

    Note that the trailing dot "." is very significant.  It means, "run
    the build relative to the current directory."

4. Push your image to the docker registry which the Kubernetes cluster
   can reach.  See [the
   quickstart](quickstart.md#prepare-the-namespace-and-docker-registry-access) to learn how
   to make the Kubernetes cluster aware of the Docker credentials so it
   can pull down images.
   
    If you are using a local Kubernetes, you can omit this step, since
    the Kubernetes pulls from the same Docker server as the one to which
    the local build command built the image.


##### 5. Install the Helm chart, passing the arguments to make the chart aware of the sidecar image
   
```
$ helm --debug install --version 1.0.0-SNAPSHOT \
     ./coherence --name hello-server-config \
     --set userArtifacts.image=coherence-operator-hello-server-config:1.0.0-SNAPSHOT \
     --set imagePullSecrets=sample-coherence-secret \
     --set store.cacheConfig=hello-server-config.xml
```

> If your XML files are in a different location within the sidecar Docker
> image, use the `--set userArtifacts.configDir=<absolute path within
> docker image>` argument to `helm install` to configure the correct
> location.

In a separate shell, set up a Kubernetes "port forward" to expose the
Extend port so that your local client can use it.  The instructions for
doing this are output from the above `helm install` command, but they
are repeated hear for your convenience.

```
$ export POD_NAME=$(kubectl get pods --namespace default -l "app=coherence,release=hello-server-config" -o jsonpath="{.items[0].metadata.name}")
$ kubectl --namespace default port-forward $POD_NAME 20000:20000
```

This prints the following output and blocks the shell:

```
Forwarding from 127.0.0.1:20000 -> 20000
Forwarding from [::1]:20000 -> 20000
```

##### 6. Run the client program

Assuming you are in the same directory as the XML and Java source files,
and that the correct `coherence.jar` is available at
`${COHERENCE_HOME}/lib/coherence.jar`, compile and run the program as
shown next:

```
$ javac -cp .:${COH_JAR} HelloConfigXml.java
$ java -cp .:${COH_JAR} -Dcoherence.cacheconfig=$PWD/hello-client-config.xml -Dcoherence.log.level=-1 HelloConfigXml
```

This should produce output similar to the following:

```
Oracle Coherence Version 12.2.1.4.0 Build 73407
 Grid Edition: Development mode
Copyright (c) 2000, 2019, Oracle and/or its affiliates. All rights reserved.

The value of the key is 1
```

Running the program again shows the value has been incremented.

```
Oracle Coherence Version 12.2.1.4.0 Build 73407
 Grid Edition: Development mode
Copyright (c) 2000, 2019, Oracle and/or its affiliates. All rights reserved.

The value of the key is 2.
```

##### Finally, let's combine the preceding two use-cases and deploy a jar containing both application classes and configuration files

Frequently, the sidecar image contains one or more jar files, each of
which may contain application classes, configuration files, or both.
Any jar files included in the sidecare image using the approach detailed
above will end up on the Coherence Classpath.  Any Java classes in those
jar files, will therefore be available for Classloading by the entire
Coherence cluster.  Any configuration files *must be included in the top
level of a jar file* in order to be referenced by the Coherence helm
helm chart.  Consider the following sidecar image layout.

```
files/
   lib/
      coherence-operator-hello-server-config-1.0.0-SNAPSHOT.jar
```

Within the jar, consider the following excerpt from the file layout.

```
META-INF/
META-INF/LICENSE
META-INF/beans.xml
META-INF/maven/org.javassist/javassist/pom.xml
META-INF/services/org.glassfish.jersey.server.spi.ContainerProvider
com/foo/demo/model/Price.class
cache-config.xml
pof-config.xml
```

The following invocation will cause the coherence helm chart to be
installed such that the `cache-config.xml` and `pof-config.xml` are
supplied to coherence, as well as all the java classes in the jar file
being in the Coherence classpath.

```
$ helm --debug install --version 1.0.0-SNAPSHOT \
     TODO/coherence --name hello-server-config \
     --set userArtifacts.image=coherence-operator-hello-server-config:1.0.0-SNAPSHOT \
     --set imagePullSecrets=sample-coherence-secret \
     --set store.cacheConfig=cache-config.xml \
     --set store.pof.config=pof-config.xml
```

#### Extract Reporter Files from a Kubernetes Coherence Pod

Any of the debugging techniques described in [Debugging in
Coherence](https://docs.oracle.com/middleware/12213/coherence/develop-applications/debugging-coherence.htm)
that call for the creation of files to be examined, such as log files
and JVM heap dumps, can also be accomplished with the Coherence
Operator.  Let's take the example of collecting a `.hprof` file for a
heap dump.  A single-command technique is included at the end of this
use case.

Assuming you have the operator and Coherence as the only apps running in
the Kubernetes cluster, the following command lists the pods of the
operator and Coherence.

```
$ kubectl get pods
NAME                                 READY     STATUS    RESTARTS   AGE
coherence-demo-storage-0             1/1       Running   0          45m
coherence-demo-storage-1             1/1       Running   0          44m
coherence-operator-7bc94cfb4-g4kz2   1/1       Running   0          47m
```

Get a shell into the storage node:

```
$ kubectl exec -it coherence-demo-storage-0 -- /bin/bash
```

Obtain the PID of the Coherence process.  Usually this is PID `1`, but
it is a good idea to use `jps` to get the actual PID.

```
bash-4.2# /usr/java/default/bin/jps
1 DefaultCacheServer
4230 Jps
```

Now use the `jcmd` command to extract the heap dump and exit the shell.

```
bash-4.2# /usr/java/default/bin/jcmd 1 GC.heap_dump /DefaultCache.hprof
bash-4.2# exit
```

Finally, use `kubectl exec` to extract the heap dump.

```
$ (kubectl exec coherence-demo-storage-0 -it -- cat /DefaultCache.hprof ) > DefaultCache.hprof
```

Assuming the Coherence PID is `1`, a repeatable single-command version of this technique is:

```
$ (kubectl exec coherence-demo-storage-0 -- /bin/bash -c "rm -f /tmp/heap.hprof; /usr/java/default/bin/jcmd 1 GC.heap_dump /tmp/heap.hprof; cat /tmp/heap.hprof > /dev/stderr" ) 2> heap.hprof
```

Note that we redirect the heap dump output to `stderr` to prevent the unsuppressable 

```
1:
Heap dump file created
```

output from `jcmd` from showing up in the heap dump file.


## Kubernetes Specific Use Cases

### Using Helm to Scale the Coherence Deployment

The Coherence Operator leverages Kubernetes Statefulsets to ensure that
scale up and scale down operations allow the underlying Coherence
cluster nodes sufficient time to rebalance the cluster data.

Assume the Coherence helm chart has been run with the default options
and a [Helm release](https://helm.sh/docs/glossary/#release) by the name of
`coherence-deploy` has been created and is successfully running, the
following command will increase the number of Coherence cluster nodes
from the default to the new value of `4`.

```
kubectl scale statefulsets coherence-deploy --replicas=4
```

Monitoring the progress of the cluster as Kubernetes adjusts to the new
intent will show the number of pods being adjusted and the status of
each pod progressing through the various stages to end up at `Running`
status.

### Perform a Safe Rolling Upgrade

The steps detailed in section [Supply a Jar File Containing Application
Classes](#table-of-use-cases) call for the creation of a sidecar docker
image that conveys the application classes to Kubernetes.  This docker
image is tagged with a version number, and the version number is how
Kubernetes enables safe rolling upgrades.  You can read more about safe
rolling upgrades in [the Helm
documentation](https://helm.sh/docs/helm/#helm-upgrade).  Briefly, as
with the scaling described in the preceding section, the safe rolling
upgrade feature allows you to instruct Kubernetes, via the operator, to
replace the currently deployed version of your application classes with
a different one.  Kubernetes does not care if the different version is
"newer" or "older", as long as it has a docker tag and can be pulled by
the cluster, that is all Kubernetes needs to know.  The operator will
ensure this is done without data loss or interruption of service.

Assuming the sidecar has been installed using the steps detailed in
[Supply a Jar File Containing Application Classes](#table-of-use-cases),
and the upgrade destination is available and has been tagged with
`coherence-operator-hello-example:1.0.1`, the following command will
instruct the operator to upgrade from
`coherence-operator-hello-example:1.0.0-SNAPSHOT` to
`coherence-operator-hello-example:1.0.1`.

```
$ helm --debug upgrade \
     ./coherence --name hello-example --reuse-values \
     --set userArtifacts.image=coherence-operator-hello-example:1.0.1 --wait \
     --set imagePullSecrets=sample-coherence-secret
```

### Deploy Multiple Coherence Clusters Managed by the Operator

The Coherence Operator is designed to be installed once on a given
Kubernetes cluster. This one [Helm
release](https://helm.sh/docs/glossary/#release) of the Coherence
Operator is able to monitor and manage all of the Coherence clusters
installed on the given Kubernetes cluster.  The following commands
install the Coherence operator, then install multiple, independent
Coherence clusters, on the same Kubernetes cluster, managed by that one
operator.

```
$ helm --debug install --version 0.9.1 TODO/coherence-operator \
    --name sample-coherence-operator \
    --set "targetNamespaces={}" \
    --set imagePullSecrets=sample-coherence-secret

$ helm --debug install --version 1.0.0-SNAPSHOT \
     --set cluster=revenue-management \
     --set imagePullSecrets=sample-coherence-secret \
     --set userArtifacts.image=revenue-app:2.0.1 \
     --name revenue-management \
     TODO/coherence

$ helm --debug install --version 1.0.0-SNAPSHOT \
     --set cluster=charging \
     --set imagePullSecrets=sample-coherence-secret \
     --set userArtifacts.image=charging-app:2.0.1 \
     --name charging \
     TODO/coherence
```

The first `helm install` installs the operator with an empty list for
the `targetNamespaces` parameter.  This causes the operator to manage
all namespaces for Coherence clusters.

The second and third `helm install` invocations differ in the values
passed to the `cluster` and `userArtifacts.image` parameters, and
`--name` option.  These values must be unique to ensure that the two
Coherence clusters to not merge and form one cluster.


## Monitoring Performance and Logging

### Use-Cases


  1. Deploy Operator w/ Prometheus enabled showing coherence metrics in
     prometheus.

  2. Deploy Operator and Grafana showing OOTB grafana dashboard and
     visualized display of the metrics in prometheus.
     
  3. Deploy Operator w/ ELK enabled, and deploy a Coherence Deployment.
If you wish to access Kibana to view logs, then please see [Accessing
the EFK stack for viewing logs](logcapture.md).

     i) Show logs going to elasticsearch and output in kibana in addition
        to highlighting what is extracted from the log record, i.e. thread
        name, member id, service name, etc.

     ii) Describe how customers put their own logs into ES.

  4. How does a customer add their own metrics? (a research item that we
     should work out how and then doc)

  5. Configuring SSL endpoints for management over REST and metrics publishing

``   Note: Management over REST and metrics publishing will be available in Coherence 12.2.1.4.``

This section describes how to configure SSL for management over REST and Prometheus metrics through two examples.

i) Configuring SSL for management over REST <p />
The following is an example on how to configure a two-way SSL for 
Coherence management over REST:

 a) Create k8s secrets for your key store and trust store files <p />
 Coherence SSL requires Java key store and trust store files. These files
 are usually password protected.
 Let's say our key store and trust store are password protected.  Below are 
 the required files:
 
```
keyStore - name of the Java keystore file: myKeystore.jks
keyStorePasswordFile - name of the keystore password file: storepassword.txt
keyPasswordFile - name of the key password file: keypassword.txt
trustStore - name of the Java trust store file: myTruststore.jks
trustStorePasswordFile - name of the trust store password file: trustpassword.txt
```

The following command creates a k8s secret, ssl-secret, to contain these files:

```
kubectl create secret generic ssl-secret \
   --namespace myNamespace \
   --from-file=./myKeystore.jks \
   --from-file=./myTruststore.jks \
   --from-file=./storepassword.txt \
   --from-file=./keypassword.txt \
   --from-file=./trustpassword.txt
```

 b) Create a YAML file, helm-values-ssl-management.yaml, to enable SSL for 
 Coherence management over REST 
 using the keystore, trust store, and password files in the ssl-secret
 we created in a):

     
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
  
 c) Install the Coherence helm chart using the YAML file created in step b): <p />

```
  helm install --version 1.0.0-SNAPSHOT TODO/coherence \
    --name coherence \
    --namespace myNamespace \
    --set imagePullSecrets=my-imagePull-secret \
    -f helm-values-ssl-management.yaml
```

To verify that Coherence management over REST is running
with https, you could forward the management listen port
to your local machine and access the REST endpoint
using the following command and URL respectively:

```
kubectl port-forward <pod name> 30000:30000
 
https://localhost:30000/management/coherence/cluster
```

If you have self-signed certificate, you may get "Your connection is not secure" from the browser.
You can click "Advanced" button, then "Add Exception..." to allow the request go through.

You could also look for the following message in the log file of the Coherence pod: <br />
`Started: HttpAcceptor{Name=Proxy:ManagementHttpProxy:HttpAcceptor, State=(SERVICE_STARTED), HttpServer=NettyHttpServer{Protocol=HTTPS, AuthMethod=cert}`

  
ii) Configuring SSL for metrics publishing for Prometheus <p />
You can either create a different k8s secret with a different set of keystore,
trust store, etc. or use the same secret used by management over rest. For our example,
we will just use the same secret, ssl-secret.  Here is an example on how to configure 
a SSL endpont for Coherence metrics:

 a) Create k8s secret for your key store and trust store files <p />
 We can skip this step since we already have a k8s secret created 
 from the management over REST example. <p />

 b) Create a YAML file, helm-values-ssl-metrics.yaml, using the keystore, trust store,
 and password file stored in ssl-secret we created in example i):

     
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
  
 c) Install the Coherence helm chart using the YAML created in step b): <p />

```
  helm install --version 1.0.0-SNAPSHOT TODO/coherence \
    --name coherence \
    --namespace myNamespace \
    --set imagePullSecrets=my-imagePull-secret \
    -f helm-values-ssl-metrics.yaml
```

To verify that Coherence metrics for Prometheus is running
with https, you could forward the Coherence metrics port and access the metrics
 from your local machine use the following commands:

```
kubectl port-forward <Coherence pod> 9095:9095
 
curl -X GET https://localhost:9095/metrics \
--cacert <caCert> --cert <certificate>
 
add "--insecure" if you use self-signed certificate.
```

You could also look for the following message in the log file of the Coherence pod:
`Started: HttpAcceptor{Name=Proxy:MetricsHttpProxy:HttpAcceptor, State=(SERVICE_STARTED), HttpServer=NettyHttpServer{Protocol=HTTPS, AuthMethod=cert}`

To configure Prometheus SSL (TLS) connections with the Coherence metrics SSL endpoints,
see: https://github.com/helm/charts/blob/master/stable/prometheus-operator/README.md
on how to specify k8s secrets that contain the certificates required for two-way SSL in Prometheus; <br />
see: https://prometheus.io/docs/prometheus/latest/configuration/configuration/#tls_config
on how to configure Prometheus to use SSL (TLS) connections.

Once you configured Prometheus to use SSL, you can verify that Prometheus is scraping Coherence
metrics over https by forwarding the Prometheus service port to your local machine
and access the following URL:

```
kubectl port-forward <Prometheus pod> 9090:9090
 
http://localhost:9090/graph
```

You shoud see many coherence_* metrics    

To enable SSL for both management over REST and metrics publishing for Prometheus, install the
Coherence chart with both YAML files:

```
  helm install --version 1.0.0-SNAPSHOT TODO/coherence \
    --name coherence \
    --namespace myNamespace \
    --set imagePullSecrets=my-imagePull-secret \
    -f helm-values-ssl-management.yaml,helm-values-ssl-metrics.yaml
```
    
## Kubernetes Specific Tasks

### Use-Cases

  2. Access the management over REST json info.
  
  3. Illustrate that changing of an image version (Coherence of
     application container) results in a safe rolling restart.

  4. Enable a PVC and therefore illustrate Coherence persisting to the
     specified location.
     
     Ensure documenting the case of using local SSD for active
     persistence and Block Storage for snapshots.

  5. Highlight how a Coherence deployment uses the zone information; zone
     should be viewable in the Member information.

<!-- 

    Moved from quickstart, work into userguide.

    ## 4. ELK Stack
    By default, the ELK is installed with Coherence Operator. Kibana is deployed
    with a NodePort type service. The Kibana host and port can be found by looking
    at service information:
    ```
    $ kubectl get svc -n sample-coherence-ns
    ```
    Note that in Docker/Kubernetes environment over Mac. The host is `localhost`.

    ## 5. Remove Coherence.
    Remove Coherence:
    ```
    $ helm delete --purge sample-coherence
    ```

-->
