<doc-view>

<h2 id="_coherence_extend_clients">Coherence Extend Clients</h2>
<div class="section">
<p>Coherence*Extend is the mechanism used by remote Coherence client applications to connect to a Coherence cluster.
Coherence*Extend includes support for native Coherence clients (Java, C++, and .NET) and non-native Coherence clients (REST and Memcached).
Coherence*Extend can be used to connect clients to Coherence clusters running in Kubernetes.
There are two scenarios, the client could also be in kubernetes, or the client could be external connecting via a service or some other form of ingress.
There are different ways to configure the client in these scenarios.</p>

<p>These examples are not going to cover all the possible use-cases for Extend, the examples are specifically about different ways to connect a client to a Coherence cluster running inside kubernetes.
Extend is extensively documented in the <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.2206/develop-remote-clients/getting-started-coherenceextend.html">official Coherence documentation</a>.</p>


<h3 id="_prerequisites">Prerequisites</h3>
<div class="section">

<h4 id="_server_image">Server Image</h4>
<div class="section">
<p>To show Extend working the example will require a Coherence cluster to connect to.
For the server the example will use the image built in the <router-link to="/examples/015_simple_image/README">Build a Coherence Server Image using JIB</router-link> example (or could also use the <router-link to="/examples/016_simple_docker_image/README">Build a Coherence Server Image using a Dockerfile</router-link> example.
If you have not already done so, you should build the image from that example, so it is available to deploy to your Kubernetes cluster.</p>

</div>

<h4 id="_install_the_operator">Install the Operator</h4>
<div class="section">
<p>If you have not already done so, you need to install the Coherence Operator.
There are a few simple ways to do this as described in the <router-link to="/docs/installation/01_installation">Installation Guide</router-link></p>

</div>
</div>
</div>

<h2 id="_the_client_application">The Client Application</h2>
<div class="section">
<p>To demonstrate different configurations and connectivity we need a simple Extend client application.</p>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p><img src="./images/GitHub-Mark-32px.png" alt="GitHub Mark 32px" />
 The complete source code for this example is in the <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/tree/main/examples/025_extend_client">Coherence Operator GitHub</a> repository.</p>
</p>
</div>
<p>As the client application only needs to demonstrate connectivity to Coherence using different configurations it is not going to do very much.
There is a single class with a <code>main</code> method. In the <code>main</code> method the code obtains a <code>NamedMap</code> from Coherence via Extend and does some simple put and get operations. If these all function correctly the application exits with an exit code of zero. If there is an exception, the stack trace is printed and the application exits with an exit code of 1.</p>

<p>There are also some different cache configuration files for the different ways to configure Extend, these are covered in the relevant examples below.</p>

</div>

<h2 id="_building_the_client_image">Building the Client Image</h2>
<div class="section">
<p>The client application is both a Maven and Gradle project, so you can use whichever you are comfortable with.
The only dependency the client application needs is <code>coherence.jar</code>.</p>


<h3 id="_using_the_maven_or_gradle_jib_plugin">Using the Maven or Gradle JIB Plugin</h3>
<div class="section">
<p>The image can be built using the JIB plugin with either Maven or Gradle, as described below.</p>

<p>Using Maven we run:</p>

<markup
lang="bash"

>./mvnw compile jib:dockerBuild</markup>

<p>Using Gradle we run:</p>

<markup
lang="bash"

>./gradlew compileJava jibDockerBuild</markup>

<p>The command above will create an image named <code>simple-extend-client</code> with two tags, <code>latest</code> and <code>1.0.0</code>.
Listing the local images should show the new images.</p>

<markup
lang="bash"

>$ docker images | grep simple
simple-extend-client   1.0.0   1613cd3b894e   51 years ago  220MB
simple-extend-client   latest  1613cd3b894e   51 years ago  220MB</markup>

</div>

<h3 id="_using_a_dockerfile">Using a Dockerfile</h3>
<div class="section">
<p>Alternatively, if you cannot use the JIB plugin in your environment, the client image can be built using a simple Dockerfile and <code>docker build</code> command. We will still use Maven or Gradle to pull all the required dependencies together.</p>

<p>Using Maven we run:</p>

<markup
lang="bash"

>./mvnw package
docker build -t simple-extend-client:1.0.0 -t simple-extend-client:latest target/docker</markup>

<p>Using Gradle we run:</p>

<markup
lang="bash"

>./gradlew assembleImage
docker build -t simple-extend-client:1.0.0 -t simple-extend-client:latest build/docker</markup>

<p>Again, the build should result in the Extend client images</p>

<p>The command above will create an image named <code>simple-extend-client</code> with two tags, <code>latest</code> and <code>1.0.0</code>.
Listing the local images should show the new images.</p>

<markup
lang="bash"

>$ docker images | grep simple
simple-extend-client   1.0.0   1613cd3b894e   2 minutes ago  220MB
simple-extend-client   latest  1613cd3b894e   2 minutes ago  220MB</markup>

<p>If we tried to run the application or image at this point it would fail with an exception as there is no cluster to connect to.</p>

</div>
</div>

<h2 id="_extend_inside_kubernetes_using_the_coherence_nameservice">Extend Inside Kubernetes Using the Coherence NameService</h2>
<div class="section">
<p>If the Extend client is going to run inside Kubernetes then we have a number of choices for configuration.
In this section we are going to use the simplest way to configure Extend in Coherence, which is to use the Coherence NameService.
In this configuration we do not need to specify any ports, the Extend proxy in the server cluster will bind to an ephemeral port.
The Extend client will then use the Coherence NameService to find the addresses and ports that the Extend proxy is listening on.</p>


<h3 id="_proxy_server_configuration">Proxy Server Configuration</h3>
<div class="section">
<p>The default cache configuration file, built into <code>coherence.jar</code> configures an Extend proxy that binds to an ephemeral port.
The proxy-scheme configuration looks like this:</p>

<markup
lang="xml"
title="coherence-cache-config.xml"
>    &lt;proxy-scheme&gt;
      &lt;service-name&gt;Proxy&lt;/service-name&gt;
      &lt;autostart system-property="coherence.proxy.enabled"&gt;true&lt;/autostart&gt;
    &lt;/proxy-scheme&gt;</markup>

<p>That is all that is required in a cache configuration file to create a proxy service that will bind to an ephemeral port.
The proxy is enabled by default, but could be disabled by setting the system property <code>coherence.proxy.enabled</code> to false.</p>

</div>

<h3 id="_deploy_the_server">Deploy the Server</h3>
<div class="section">
<p>To run the NameService examples below the server needs to be deployed.
The example includes a <code>manifests/</code> directory containing Kubernetes yaml files used by the example.</p>

<p>For the NameService examples below the server will use the default cache configuration file from <code>coherence.jar</code> which has the <code>Proxy</code> service configured above. The yaml to deploy the server cluster is in the <code>manifests/default-server.yaml</code> file.</p>

<markup
lang="yaml"
title="manifests/default-server.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  image: simple-coherence:1.0.0
  replicas: 3
  coherence:
    cacheConfig: coherence-cache-config.xml</markup>

<p>The yaml above will deploy a three member cluster configured to use the default <code>coherence-cache-config.xml</code> configuration file.</p>

<p>There are no additional ports exposed in the configuration. The Extend proxy will be listening on an ephemeral port, so we have no idea what that port will be.</p>

<p>We can deploy the server into the default namespace in kubernetes with the following command:</p>

<markup
lang="bash"

>kubectl apply -f manifests/default-server.yaml</markup>

<p>We can list the resources created by the Operator.</p>

<markup
lang="bash"

>kubectl get all</markup>

<p>Which should display something like this:</p>

<markup
lang="bash"

>NAME            READY   STATUS    RESTARTS   AGE
pod/storage-0   1/1     Running   0          81s
pod/storage-1   1/1     Running   0          81s
pod/storage-2   1/1     Running   0          81s

NAME                    TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
service/storage-sts     ClusterIP   None         &lt;none&gt;        7/TCP     81s
service/storage-wka     ClusterIP   None         &lt;none&gt;        7/TCP     81s

NAME                       READY   AGE
statefulset.apps/storage   3/3     81s</markup>

<ul class="ulist">
<li>
<p>We can see that the Operator has created a <code>StatefulSet</code>, with three <code>Pods</code> and there are two <code>Services</code>.</p>

</li>
<li>
<p>The <code>storage-sts</code> service is the headless service required for the <code>StatefulSet</code>.</p>

</li>
<li>
<p>The <code>storage-wka</code> service is the headless service that Coherence will use for well known address cluster discovery.</p>

</li>
</ul>
</div>

<h3 id="_minimal_extend_client_configuration">Minimal Extend Client Configuration</h3>
<div class="section">
<p>The configuration required for the Extend client is equally minimal.
The example source code includes a configuration file named <code>src/main/resources/minimal-client-cache-config.xml</code> that can be used to connect to the proxy configured above.</p>

<markup
lang="xml"
title="src/main/resources/minimal-client-cache-config.xml"
>&lt;?xml version="1.0"?&gt;
&lt;cache-config xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
              xmlns="http://xmlns.oracle.com/coherence/coherence-cache-config"
              xsi:schemaLocation="http://xmlns.oracle.com/coherence/coherence-cache-config coherence-cache-config.xsd"&gt;
  &lt;caching-scheme-mapping&gt;
    &lt;cache-mapping&gt;                        <span class="conum" data-value="1" />
      &lt;cache-name&gt;*&lt;/cache-name&gt;
      &lt;scheme-name&gt;remote&lt;/scheme-name&gt;
    &lt;/cache-mapping&gt;
  &lt;/caching-scheme-mapping&gt;

  &lt;caching-schemes&gt;
    &lt;remote-cache-scheme&gt;
      &lt;scheme-name&gt;remote&lt;/scheme-name&gt;                 <span class="conum" data-value="2" />
      &lt;service-name&gt;RemoteService&lt;/service-name&gt;        <span class="conum" data-value="3" />
      &lt;proxy-service-name&gt;Proxy&lt;/proxy-service-name&gt;    <span class="conum" data-value="4" />
    &lt;/remote-cache-scheme&gt;
  &lt;/caching-schemes&gt;
&lt;/cache-config&gt;</markup>

<ul class="colist">
<li data-value="1">There is a single <code>cache-mapping</code> that maps all cache names to the scheme named <code>remote</code>.</li>
<li data-value="2">The <code>remote-scheme</code> is named <code>remote</code>.</li>
<li data-value="3">The <code>remote-scheme</code> has a service name of <code>RemoteService</code>.</li>
<li data-value="4">The remote service will connect to a proxy service on the server that is named <code>Proxy</code>, this must correspond to the name of the proxy service in our server cache configuration file.</li>
</ul>

<h4 id="_deploy_the_client">Deploy the Client</h4>
<div class="section">
<p>The simplest way to run the Extend client in Kubernetes is as a <code>Job</code>. The client just connects to a cache and does a <code>put</code>, then exits, so a <code>Job</code> is ideal for this type of process. The example contains yaml to create a Job <code>manifests/minimal-job.yaml</code> that looks like this:</p>

<markup
lang="yaml"
title="manifests/minimal-job.yaml"
>apiVersion: batch/v1
kind: Job
metadata:
  name: extend-client
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: client
          image: simple-extend-client:1.0.0
          env:
            - name: COHERENCE_CACHE_CONFIG
              value: minimal-client-cache-config.xml
            - name: COHERENCE_WKA
              value: storage-wka
            - name: COHERENCE_CLUSTER
              value: storage</markup>

<p>To be able to run the client we need to set in three pieces of information.</p>

<ul class="ulist">
<li>
<p>The name of the cache configuration file. We set this using the <code>COHERENCE_CACHE_CONFIG</code> environment variable, and set the value to <code>minimal-client-cache-config.xml</code>, which is the configuration file we&#8217;re using in this example.</p>

</li>
<li>
<p>The client needs to be able to discover the storage Pods to connect to. Just like the server cluster uses well known addresses to discover a cluster, the client can do the same. We set the <code>COHERENCE_WKA</code> environment variable to the name of the WKA service created for the server when we deployed it above, in this case it is <code>storage-wka</code>.</p>

</li>
<li>
<p>Finally, we set the name of the Coherence cluster the client will connect to. When we deployed the server we did not specify a name, so the default cluster name will be the same as the <code>Coherence</code> resource name, in this case <code>storage</code>. So we set the <code>COHERENCE_CLUSTER</code> environment variable to <code>storage</code>.</p>

</li>
</ul>
<p>The client <code>Job</code> can be deployed into the default namespace in Kubernetes with the following command:</p>

<markup
lang="bash"

>kubectl apply -f manifests/minimal-job.yaml</markup>

<p>The <code>Jobs</code> deployed can then be listed</p>

<markup
lang="bash"

>kubectl get job</markup>

<p>Which should display something like this:</p>

<markup
lang="bash"

>NAME            COMPLETIONS   DURATION   AGE
extend-client   1/1           4s         5s</markup>

<p>The <code>Job</code> above completed very quickly, which we would expect as it is just doing a trivial put to a cache.</p>

<p>We can list the <code>Pods</code> created for the <code>Job</code> and then look at the log from the client.
All <code>Pods</code> associated to a <code>Job</code> have a label in the form <code>job-name: &lt;name-of-job&gt;</code>, so in our case the label will be <code>job-name: extend-client</code>.
We can use this with <code>kubectl</code> to list <code>Pods</code> associated to the <code>Job</code>. If the <code>Job</code> ran successfully there should be only one <code>Pod</code>. If the <code>Job</code> failed and has a restart policy, or was restarted by Kubernetes for other reasons there could be multiple <code>Pods</code>. In this case we expect a single successful <code>Pod</code>.</p>

<markup
lang="bash"

>kubectl get pod -l job-name=extend-client</markup>

<markup
lang="bash"

>NAME                  READY   STATUS      RESTARTS   AGE
extend-client-k7wfq   0/1     Completed   0          4m24s</markup>

<p>If we look at the log for the <code>Pod</code> we should see the last line printed to <code>System.out</code> by the client:</p>

<markup
lang="bash"

>kubectl logs extend-client-k7wfq</markup>

<p>The last line of the log will be something like this:</p>

<markup
lang="bash"

>Put key=key-1 value=0.9332279895860512 previous=null</markup>

<p>The values will be different as we put different random values each time the client runs.
The previous value was <code>null</code> in this case as we have not run any other client with this cluster. If we re-ran the client <code>Job</code> the previous value would be displayed as the cache on the server now has data in it.</p>

<p><strong>Clean-Up</strong></p>

<p>We have shown a simple Extend client running in Kubernetes, connecting to a Coherence cluster using the NameService.
We can now delete the <code>Job</code> using <code>kubectl</code>.</p>

<markup
lang="bash"

>kubectl delete job extend-client</markup>

<p>We can also delete the server.</p>

<markup
lang="bash"

>kubectl delete -f manifests/default-server.yaml</markup>

</div>
</div>

<h3 id="_deploy_the_client_to_a_different_namespace">Deploy the Client to a Different Namespace</h3>
<div class="section">
<p>In the first example we deployed the client to the same namespace as the server.
If we wanted to deploy the client to a different namespace we would need to ensure the fully qualified name of the WKA service is used when setting the <code>COHERENCE_WKA</code> environment variable. The Coherence cluster is deployed into the <code>default</code> namespace so the fully qualified WKA service name is <code>storage-wka.default.svc</code>, or we could also use <code>storage-wka.default.svc.cluster.local</code>.</p>

<markup
lang="yaml"
title="manifests/minimal-job.yaml"
>apiVersion: batch/v1
kind: Job
metadata:
  name: extend-client
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: client
          image: simple-extend-client:1.0.0
          env:
            - name: COHERENCE_CACHE_CONFIG
              value: minimal-client-cache-config.xml
            - name: COHERENCE_WKA
              value: storage-wka.default.svc.cluster.local
            - name: COHERENCE_CLUSTER
              value: storage</markup>

<p>We can deploy this client <code>Job</code> into a different namespace than the cluster is deployed into:</p>

<markup
lang="bash"

>kubectl create ns coherence-test
kubectl apply -f manifests/minimal-other-namespace-job.yaml -n coherence-test</markup>

<p>We should see the <code>Job</code> complete successfully.</p>

</div>
</div>

<h2 id="_extend_clients_external_to_kubernetes">Extend Clients External to Kubernetes</h2>
<div class="section">
<p>The NameService example above will only work if the client is running inside the same Kubernetes cluster as the server.
When the client uses the Coherence NameService to look up the addresses of the Extend proxy service, the cluster only knows its internal IP addresses. If a client external to Kubernetes tried to use the NameService the addresses returned would be unreachable, as they are internal to the Kubernetes cluster.</p>

<p>To connect external Extend clients, the proxy must be bound to known ports and those ports exposed to the client via some form of service or ingress.</p>


<h3 id="_proxy_server_configuration_2">Proxy Server Configuration</h3>
<div class="section">
<p>The Extend proxy service on the server must be configured to have a fixed port, so there is a little more configuration than previously.</p>

<p>The example server image contains a Coherence configuration file named <code>test-cache-config.xml</code>, which contains an Extend proxy configured to bind to all host addresses (<code>0.0.0.0</code>) on port 20000.</p>

<markup
lang="xml"
title="test-cache-config.xml"
>&lt;proxy-scheme&gt;
  &lt;service-name&gt;Proxy&lt;/service-name&gt;
  &lt;acceptor-config&gt;
    &lt;tcp-acceptor&gt;
      &lt;local-address&gt;
        &lt;!-- The proxy will listen on all local addresses --&gt;
        &lt;address&gt;0.0.0.0&lt;/address&gt;
        &lt;port&gt;20000&lt;/port&gt;
      &lt;/local-address&gt;
    &lt;/tcp-acceptor&gt;
  &lt;/acceptor-config&gt;
  &lt;autostart&gt;true&lt;/autostart&gt;
&lt;/proxy-scheme&gt;</markup>

</div>

<h3 id="_deploy_the_server_2">Deploy the Server</h3>
<div class="section">
<p>The example contains a yaml file that can be used to deploy a Coherence server with the fixed proxy address, as shown above.</p>

<markup
lang="yaml"
title="manifests/fixed-port-server.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  image: simple-coherence:1.0.0
  replicas: 3
  coherence:
    cacheConfig: test-cache-config.xml
  ports:
    - name: extend
      port: 20000</markup>

<p>The yaml above will deploy a three member cluster configured to use the default <code>test-cache-config.xml</code> configuration file and expose the Extend port  via a service.</p>

<p>The server can be deployed with the following command.</p>

<markup
lang="bash"

>kubectl apply -f manifests/fixed-port-server.yaml</markup>

<p>The resources created by the Coherence Operator can be listed:</p>

<markup
lang="bash"

>kubectl get all</markup>

<markup
lang="bash"

>NAME            READY   STATUS    RESTARTS   AGE
pod/storage-0   1/1     Running   0          61s
pod/storage-1   1/1     Running   0          61s
pod/storage-2   1/1     Running   0          61s

NAME                     TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)     AGE
service/storage-extend   ClusterIP   10.101.99.24   &lt;none&gt;        20000/TCP   61s
service/storage-sts      ClusterIP   None           &lt;none&gt;        7/TCP       61s
service/storage-wka      ClusterIP   None           &lt;none&gt;        7/TCP       61s

NAME                       READY   AGE
statefulset.apps/storage   3/3     61s</markup>

<p>As well as the <code>Pods</code> and <code>Services</code> created in the previous example, there is now a <code>Service</code> named <code>storage-extend</code>, which exposes the Extend port.</p>

</div>

<h3 id="_configure_the_extend_client">Configure the Extend Client</h3>
<div class="section">
<p>An external client needs to be configured with a remote scheme that connects to a known address and port.
The example contains a cache configuration file named <code>src/main/resources/fixed-address-cache-config.xml</code> that has this configuration.</p>

<markup
lang="xml"
title="src/main/resources/fixed-address-cache-config.xml"
>&lt;remote-cache-scheme&gt;
  &lt;scheme-name&gt;remote&lt;/scheme-name&gt;
  &lt;service-name&gt;RemoteCache&lt;/service-name&gt;
  &lt;proxy-service-name&gt;Proxy&lt;/proxy-service-name&gt;
  &lt;initiator-config&gt;
    &lt;tcp-initiator&gt;
      &lt;remote-addresses&gt;
        &lt;socket-address&gt;
            &lt;!-- the 127.0.0.1 loop back address will only work in local dev testing --&gt;
            &lt;address system-property="coherence.extend.address"&gt;127.0.0.1&lt;/address&gt;
            &lt;port system-property="coherence.extend.port"&gt;20000&lt;/port&gt;
        &lt;/socket-address&gt;
      &lt;/remote-addresses&gt;
    &lt;/tcp-initiator&gt;
  &lt;/initiator-config&gt;
&lt;/remote-cache-scheme&gt;</markup>

<p>When the client runs using the configuration above it will attempt to connect to an Extend proxy on <code>127.0.0.1:20000</code>.
The address to connect to can be overridden by setting the <code>coherence.extend.address</code> system property.
The port to connect to can be overridden by setting the <code>coherence.extend.port</code> system property.</p>

</div>

<h3 id="_run_the_extend_client">Run the Extend Client</h3>
<div class="section">
<p>This example assumes that you are running Kubernetes on a development machine, for example with <code>KinD</code>, of <code>Minikube</code> or in Docker, etc.
In this case the <code>Service</code> created is of type <code>ClusterIP</code>, so it is not actually exposed outside of Kubernetes as most development Kubernetes clusters do not support services of type <code>LoadBalancer</code>.</p>

<p>This means that to test the external client we will need to use port forwarding.
In a console start the port forwarder using <code>kubectl</code> as follows</p>

<markup
lang="bash"

>kubectl port-forward svc/storage-extend 20000:20000</markup>

<p>The example client can not connect to the Extend proxy via the host machine on port <code>20000</code>.</p>

<p>The simplest way to run the Extend client locally is to use either Maven or Gradle.
The Maven <code>pom.xml</code> file uses the Maven Exec plugin to run the client.
The Gradle <code>build.gradle</code> file configures a run task to execute the client.</p>

<p>With Maven:</p>

<markup
lang="bash"

>./mvnw compile exec:java</markup>

<p>With Gradle:</p>

<markup
lang="bash"

>./gradlew runClient</markup>

<p>Both of the above commands run successfully and the final line of output should be the line printed by the client showing the result of the put.</p>

<markup
lang="bash"

>Put key=key-1 value=0.5274436018741687 previous=null</markup>

<p><strong>Clean-up</strong></p>

<p>We can now delete the server.</p>

<markup
lang="bash"

>kubectl delete -f manifests/fixed-port-server.yaml</markup>

</div>
</div>

<h2 id="_mixing_internal_and_external_extend_clients">Mixing Internal and External Extend Clients</h2>
<div class="section">
<p>The example server configuration used for connecting external clients can also be used for internal Extend clients, which is useful for use-cases where some clients are inside Kubernetes and some outside.
An Extend client running inside Kubernetes then has the choice of using the NameService configuration from the first example, or using the fixed address and port configuration of the second example.</p>

<p>If an internal Extend client is configured to use a fixed address then the host name of the proxy can be set to the service used to expose the server&#8217;s extend port.</p>

<p>For example, if the client&#8217;s cache configuration file contains a remote scheme like the external example above:</p>

<markup
lang="xml"
title="src/main/resources/fixed-address-cache-config.xml"
>&lt;remote-cache-scheme&gt;
  &lt;scheme-name&gt;remote&lt;/scheme-name&gt;
  &lt;service-name&gt;RemoteCache&lt;/service-name&gt;
  &lt;proxy-service-name&gt;Proxy&lt;/proxy-service-name&gt;
  &lt;initiator-config&gt;
    &lt;tcp-initiator&gt;
      &lt;remote-addresses&gt;
        &lt;socket-address&gt;
            &lt;!-- the 127.0.0.1 loopback address will only work in local dev testing --&gt;
            &lt;address system-property="coherence.extend.address"&gt;127.0.0.1&lt;/address&gt;
            &lt;port system-property="coherence.extend.port"&gt;20000&lt;/port&gt;
        &lt;/socket-address&gt;
      &lt;/remote-addresses&gt;
    &lt;/tcp-initiator&gt;
  &lt;/initiator-config&gt;
&lt;/remote-cache-scheme&gt;</markup>

<p>The client would be run with the <code>coherence.extend.address</code> system property, (or <code>COHERENCE_EXTEND_ADDRESS</code> environment variable) set to the fully qualified name of the Extend service, in the case of our example server running in the default namespace, this would be <code>-Dcoherence.extend.address=storage-extend.default.svc.cluster.local</code></p>

</div>

<h2 id="_external_client_in_the_real_world">External Client in the Real World</h2>
<div class="section">
<p>The example above used port-forward to connect the external Extend client to the cluster.
This showed how to configure the client and server but is not how a real world application would work.
In a real deployment the server would typically be deployed with the Extend service behind a load balancer or some other form of ingress, such as Istio. The Extend client would then be configured to connect to the external ingress address and port.
Some ingress, such as Istio, can also be configured to add TLS security, which Extend will work with.</p>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p>There is an open source project named <a id="" title="" target="_blank" href="https://metallb.universe.tf">MetalLB</a> that can easily be deployed into development environment Kubernetes clusters and provides support for load balancer services. This is a simple way to test and try out load balancers in development Kubernetes.</p>

<p>If MetalLB was installed (or your cluster supports LoadBalancer services) the yaml for deploying the cluster can be altered to make the Extend service a load balancer.</p>

<markup
lang="yaml"
title="manifests/fixed-port-lb-server.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  image: simple-coherence:1.0.0
  replicas: 3
  coherence:
    cacheConfig: test-cache-config.xml
  ports:
    - name: extend
      port: 20000
      service:
        type: LoadBalancer</markup>

<p>This can be deployed using:</p>

<markup
lang="bash"

>kubectl apply -f manifests/fixed-port-lb-server.yaml</markup>

<p>Now if we look at the Extend service, we see it is a load balancer</p>

<markup
lang="bash"

>kubectl get svc storage-extend</markup>

<markup
lang="bash"

>NAME             TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)           AGE
storage-extend   LoadBalancer   10.110.84.229   127.0.0.240   20000:30710/TCP   2m20s</markup>

<p>Exactly how you connect to the MetalLB load balancer, and on which address, varies depending on where your Kubernetes cluster is running.</p>
</p>
</div>
</div>
</doc-view>
