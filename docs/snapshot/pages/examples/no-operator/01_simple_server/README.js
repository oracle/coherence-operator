<doc-view>

<h2 id="_a_simple_coherence_cluster_in_kubernetes">A Simple Coherence Cluster in Kubernetes</h2>
<div class="section">
<p>This example shows how to deploy a simple Coherence cluster in Kubernetes manually, without using the Coherence Operator.</p>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p><img src="./images/GitHub-Mark-32px.png" alt="GitHub Mark 32px" />
 The complete source code for this example is in the <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/tree/main/examples/no-operator/01_simple_server">Coherence Operator GitHub</a> repository.</p>
</p>
</div>
<p><strong>Prerequisites</strong></p>

<p>This example assumes that you have already built the example server image.</p>

</div>

<h2 id="_create_the_kubernetes_resources">Create the Kubernetes Resources</h2>
<div class="section">
<p>Now we have an image we can create the yaml files required to run the Coherence cluster in Kubernetes.</p>


<h3 id="_statefulset_and_services">StatefulSet and Services</h3>
<div class="section">
<p>We will run Coherence using a <code>StatefulSet</code> and in Kubernetes all <code>StatefulSet</code> resources also require a headless <code>Service</code>.</p>


<h4 id="_statefulset_headless_service">StatefulSet Headless Service</h4>
<div class="section">
<markup
lang="yaml"
title="coherence.yaml"
>apiVersion: v1
kind: Service
metadata:
  name: storage-sts
  labels:
    coherence.oracle.com/cluster: test-cluster
    coherence.oracle.com/deployment: storage
    coherence.oracle.com/component: statefulset-service
spec:
  type: ClusterIP
  clusterIP: None
  ports:
  - name: tcp-coherence
    port: 7
    protocol: TCP
    targetPort: 7
  publishNotReadyAddresses: true
  selector:
    coherence.oracle.com/cluster: test-cluster
    coherence.oracle.com/deployment: storage</markup>

<p>The <code>Service</code> above named <code>storage-sts</code> has a selector that must match labels on the <code>Pods</code> in the <code>StatefulSet</code>.
We use port 7 in this <code>Service</code> because all services must define at least one port, but we never use this port and nothing in the Coherence <code>Pods</code> will bind to port 7.</p>

</div>

<h4 id="_coherence_well_known_address_headless_service">Coherence Well Known Address Headless Service</h4>
<div class="section">
<p>When running Coherence clusters in Kubernetes we need to use well-known-addressing for Coherence cluster discovery.
For this to work we create a <code>Service</code> that we can use for discovery of <code>Pods</code> that are in the cluster.
In this example we only have a single <code>StatefulSet</code>, so we could just use the headless service above for WKA too.
But in Coherence clusters where there are multiple <code>StatefulSets</code> in the cluster we would have to use a separate <code>Service</code>.</p>

<markup
lang="yaml"
title="coherence.yaml"
>apiVersion: v1
kind: Service
metadata:
  name: storage-wka
  labels:
    coherence.oracle.com/cluster: test-cluster
    coherence.oracle.com/deployment: storage
    coherence.oracle.com/component: wka-service
spec:
  type: ClusterIP
  clusterIP: None
  ports:
  - name: tcp-coherence
    port: 7
    protocol: TCP
    targetPort: 7
  publishNotReadyAddresses: true
  selector:
    coherence.oracle.com/cluster: test-cluster</markup>

<p>The <code>Service</code> above named <code>storage-wka</code> is almost identical to the <code>StatefulSet</code> service.
It only has a single selector label, so will match all <code>Pods</code> with the label <code>coherence.oracle.com/cluster: test-cluster</code> regardless of which <code>StatefulSet</code> they belong to.</p>

<p>The other important property of the WKA <code>Service</code> is that it must have the field <code>publishNotReadyAddresses: true</code> so that <code>Pods</code> with matching labels are assigned to the <code>Service</code> even when those <code>Pods</code> are not ready.</p>

</div>

<h4 id="_the_statefulset">The StatefulSet</h4>
<div class="section">
<p>We can now create the <code>StatefulSet</code> yaml.</p>

<markup
lang="yaml"
title="coherence.yaml"
>apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: storage
spec:
  selector:
    matchLabels:
      coherence.oracle.com/cluster: test-cluster
      coherence.oracle.com/deployment: storage
  serviceName: storage-sts
  replicas: 3
  template:
    metadata:
      labels:
        coherence.oracle.com/cluster: test-cluster
        coherence.oracle.com/deployment: storage
    spec:
      containers:
        - name: coherence
          image: simple-coherence:1.0.0
          command:
            - java
          args:
            - -cp
            - "@/app/jib-classpath-file"
            - -Xms1800m
            - -Xmx1800m
            - "@/app/jib-main-class-file"
          env:
            - name: COHERENCE_CLUSTER
              value: storage
            - name: COHERENCE_WKA
              value: storage-wka.svc.cluster.local
            - name: COHERENCE_CACHECONFIG
              value: "test-cache-config.xml"
          ports:
            - name: extend
              containerPort: 20000</markup>

<ul class="ulist">
<li>
<p>The <code>StatefulSet</code> above will create a Coherence cluster with three replicas (or <code>Pods</code>).</p>

</li>
<li>
<p>There is a single <code>container</code> in the <code>Pod</code> named <code>coherence</code> that will run the image <code>simple-coherence:1.0.0</code> we created above.</p>

</li>
<li>
<p>The command line used to run the container will be <code>java -cp @/app/jib-classpath-file -Xms1800m -Xmx1800m @/app/jib-main-class-file</code></p>

</li>
<li>
<p>Because we used JIB to create the image, there will be a file named <code>/app/jib-classpath-file</code> that contains the classpath for the application. We can use this to set the classpath on the JVM command line using <code>-cp @/app/jib-classpath-file</code> so in our yaml we know we will have the correct classpath for the image we built. If we change the classpath by changing project dependencies in the <code>pom.xml</code> file for our project and rebuild the image the container in Kubernetes will automatically use the changed classpath.</p>

</li>
<li>
<p>JIB also creates a file in the image named <code>/app/jib-main-class-file</code> which contains the name of the main class we specified in the JIB Maven plugin. We can use <code>@/app/jib-main-class-file</code> in place of the main class in our command line so that we run the correct main class in our container. If we change the main class in the JIB settings when we build the image our container in Kubernetes will automatically run the correct main class.</p>

</li>
<li>
<p>We set both the min and max heap to 1.8 GB (it is a Coherence recommendation to set both min and max heap to the same value rather than set a smaller -Xms).</p>

</li>
<li>
<p>The main class that will run will be <code>com.tangosol.net.Coherence</code>.</p>

</li>
<li>
<p>The cache configuration file configures a Coherence Extend proxy service, which will listen on port <code>20000</code>. We need to expose this port in the container&#8217;s ports section.</p>

</li>
<li>
<p>We set a number of environment variables for the container:</p>

</li>
</ul>

<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 33.333%;">
<col style="width: 33.333%;">
<col style="width: 33.333%;">
</colgroup>
<thead>
<tr>
<th>Name</th>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td class="">COHERENCE_CLUSTER</td>
<td class="">storage</td>
<td class="">This sets the cluster name in Coherence (the same as setting <code>-Dcoherence.cluster=storage</code>)</td>
</tr>
<tr>
<td class="">COHERENCE_WKA</td>
<td class="">storage-wka</td>
<td class="">This sets the DNS name Coherence will use for discovery of other Pods in cluster. It is set to the name of the WKA <code>Service</code> created above.</td>
</tr>
<tr>
<td class="">COHERENCE_CACHECONFIG</td>
<td class="">"test-cache-config.xml"</td>
<td class="">This tells Coherence the name of the cache configuration file to use (the same as setting <code>-Dcoherence.cacheconfig=test-cache-config.xml</code>);</td>
</tr>
</tbody>
</table>
</div>
</div>

<h4 id="_coherence_extend_service">Coherence Extend Service</h4>
<div class="section">
<p>In the cache configuration used in the image Coherence will run a Coherence Extend proxy service, listening on port 20000.
This port has been exposed in the Coherence container in the <code>StatefulSet</code> and we can also expose it via a <code>Service</code>.</p>

<markup
lang="yaml"
title="coherence.yaml"
>apiVersion: v1
kind: Service
metadata:
  name: storage-extend
  labels:
    coherence.oracle.com/cluster: test-cluster
    coherence.oracle.com/deployment: storage
    coherence.oracle.com/component: wka-service
spec:
  type: ClusterIP
  ports:
  - name: extend
    port: 20000
    protocol: TCP
    targetPort: extend
  selector:
    coherence.oracle.com/cluster: test-cluster
    coherence.oracle.com/deployment: storage</markup>

<p>The type of the <code>Service</code> above is <code>ClusterIP</code>, but we could just as easily use a different type depending on how the service will be used. For example, we might use ingress, or Istio, or a load balancer if the Extend clients were connecting from outside the Kubernetes cluster. In local development we can just port forward to the service above.</p>

</div>
</div>
</div>

<h2 id="_deploy_to_kubernetes">Deploy to Kubernetes</h2>
<div class="section">
<p>We can combine all the snippets of yaml above into a single file and deploy it to Kubernetes.
The source code for this example contains a file named <code>coherence.yaml</code> containing all the configuration above.
We can deploy it with the following command:</p>

<markup
lang="bash"

>kubectl apply -f coherence.yaml</markup>

<p>We can see all the resources created in Kubernetes by running the following command:</p>

<markup
lang="bash"

>kubectl get all</markup>

<p>Which will display something like the following:</p>

<markup


>NAME            READY   STATUS    RESTARTS   AGE
pod/storage-0   1/1     Running   0          19s
pod/storage-1   1/1     Running   0          17s
pod/storage-2   1/1     Running   0          16s

NAME                     TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)     AGE
service/storage-extend   ClusterIP   10.105.78.34   &lt;none&gt;        20000/TCP   19s
service/storage-sts      ClusterIP   None           &lt;none&gt;        7/TCP       19s
service/storage-wka      ClusterIP   None           &lt;none&gt;        7/TCP       19s

NAME                       READY   AGE
statefulset.apps/storage   3/3     19s</markup>

<p>We can see there are three <code>Pods</code> as we specified three replicas.
The three <code>Services</code> we specified have been created.
Finally, the <code>StatefulSet</code> exists and has three ready replicas.</p>

</div>

<h2 id="_connect_an_extend_client">Connect an Extend Client</h2>
<div class="section">
<p>Now we have a Coherence cluster running in Kubernetes we can try connecting a simple Extend client.
For this example we will use the test client Maven project to run the client.</p>

<p>To connect from our local dev machine into the server we will use port-forward in this example.
We could have configured ingress and load balancing, etc. but for local dev and test port-forward is simple and easy.</p>

<p>The client is configured to connect to an Extend proxy listening on <code>127.0.0.1:20000</code>. The server we have deployed into Kubernetes is listening also listening on port 20000 via the <code>storage-extend</code> service. If we run a port-forward process that forwards port 20000 on our local machine to port 20000 of the service we can connect the client without needing any other configuration.</p>

<markup
lang="bash"

>kubectl port-forward service/storage-extend 20000:20000</markup>

<p>Now in another terminal window, we can run the test client from the <code>test-client/</code> directory execute the following command:</p>

<markup
lang="bash"

>mvn exec:java</markup>

<p>This will start a Coherence interactive console which will eventually print the <code>Map (?):</code> prompt.
The console is now waiting for commands, so we can go ahead and create a cache.</p>

<p>At the <code>Map (?):</code> prompt type the command <code>cache test</code> and press enter. This will create a cache named <code>test</code></p>

<markup


>Map (?): cache test</markup>

<p>We should see output something like this:</p>

<markup


>2021-09-17 12:25:12.143/14.600 Oracle Coherence CE 21.12.1 &lt;Info&gt; (thread=com.tangosol.net.CacheFactory.main(), member=1): Loaded cache configuration from "file:/Users/jonathanknight/dev/Projects/GitOracle/coherence-operator-3.0/examples/no-operator/test-client/target/classes/client-cache-config.xml"
2021-09-17 12:25:12.207/14.664 Oracle Coherence CE 21.12.1 &lt;D5&gt; (thread=com.tangosol.net.CacheFactory.main(), member=1): Created cache factory com.tangosol.net.ExtensibleConfigurableCacheFactory

Cache Configuration: test
  SchemeName: remote
  ServiceName: RemoteCache
  ServiceDependencies: DefaultRemoteCacheServiceDependencies{RemoteCluster=null, RemoteService=Proxy, InitiatorDependencies=DefaultTcpInitiatorDependencies{EventDispatcherThreadPriority=10, RequestTimeoutMillis=30000, SerializerFactory=null, TaskHungThresholdMillis=0, TaskTimeoutMillis=0, ThreadPriority=10, WorkerThreadCount=0, WorkerThreadCountMax=2147483647, WorkerThreadCountMin=0, WorkerThreadPriority=5}{Codec=null, FilterList=[], PingIntervalMillis=0, PingTimeoutMillis=30000, MaxIncomingMessageSize=0, MaxOutgoingMessageSize=0}{ConnectTimeoutMillis=30000, RequestSendTimeoutMillis=30000}{LocalAddress=null, RemoteAddressProviderBldr=com.tangosol.coherence.config.builder.WrapperSocketAddressProviderBuilder@35f8cdc1, SocketOptions=SocketOptions{LingerTimeout=0, KeepAlive=true, TcpNoDelay=true}, SocketProvideBuilderr=com.tangosol.coherence.config.builder.SocketProviderBuilder@1e4cf40, isNameServiceAddressProvider=false}}{DeferKeyAssociationCheck=false}

Map (test):</markup>

<p>The cache named <code>test</code> has been created and prompt has changed to <code>Map (test):</code>, so this confirms that we have connected to the Extend proxy in the server running in Kubernetes.</p>

<p>We can not put data into the cache using the <code>put</code> command</p>

<markup


>Map (test): put key-1 value-1</markup>

<p>The command above puts an entry into the <code>test</code> cache with a key of <code>"key-1"</code> and a value of <code>"value-1"</code> and will print the previous value mapped to the <code>"key-1"</code> key, which in this case is <code>null</code>.</p>

<markup


>Map (test): put key-1 value-1
null

Map (test):</markup>

<p>We can now do a <code>get</code> command to fetch the entry we just put, which should print <code>value-1</code> and re-display the command prompt.</p>

<markup


>Map (test): get key-1
value-1

Map (test):</markup>

<p>To confirm we really have connected to the server we can kill the console wil ctrl-C, restart it and execute the <code>cache</code> and <code>get</code> commands again.</p>

<markup


>Map (?): cache test

... output removed for brevity ...

Map (test): get key-1
value-1

Map (test):</markup>

<p>We can see above that the get command returned <code>value-1</code> which we previously inserted.</p>

</div>

<h2 id="_clean_up">Clean-UP</h2>
<div class="section">
<p>We can now exit the test client by pressing ctrl-C, stop the port-forward process with crtl-C and undeploy the server:</p>

<markup
lang="bash"

>kubectl delete -f coherence.yaml</markup>

</div>
</doc-view>
