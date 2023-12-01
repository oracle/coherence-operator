<doc-view>

<h2 id="_using_coherence_with_istio">Using Coherence with Istio</h2>
<div class="section">
<p>You can run the Coherence cluster and manage them using the Coherence Operator alongside <a id="" title="" target="_blank" href="https://istio.io">Istio</a>.
Coherence clusters managed with the Coherence Operator 3.3.2 and later work with Istio 1.9.1 and later out of the box.
Coherence caches can be accessed from outside the Coherence cluster via Coherence*Extend, REST, and other supported
Coherence clients.
Using Coherence clusters with Istio does not require the Coherence Operator to also be using Istio (and vice-versa) .
The Coherence Operator can manage Coherence clusters independent of whether those clusters are using Istio or not.</p>

<p>Although Coherence itself can be configured to use TLS, when using Istio Coherence cluster members and clients can
just use the default socket configurations and Istio will control and route all the traffic over mTLS.</p>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p>Coherence clusters can be manually configured to work with Istio, even if not using the Operator.
See the Istio example in the <router-link to="/examples/no-operator/04_istio/README">No Operator Examples</router-link></p>
</p>
</div>

<h3 id="_how_does_coherence_work_with_istio">How Does Coherence Work with Istio?</h3>
<div class="section">
<p>Istio is a "Service Mesh" so the clue to how Istio works in Kubernetes is in the name, it relies on the configuration
of Kubernetes Services.
This means that any ports than need to be accessed in Pods, including those using in "Pod to Pod" communication
must be exposed via a Service. Usually a Pod can reach any port on another Pod even if it is not exposed in the
container spec, but this is not the case when using Istio as only ports exposed by the Envoy proxy are allowed.</p>

<p>For Coherence cluster membership, this means the cluster port and the local port must be exposed on a Service.
To do this the local port must be configured to be a fixed port instead of the default ephemeral port.
The Coherence Operator uses the default cluster port of <code>7574</code> and there is no reason to ever change this.
The Coherence Operator always configures a fixed port for the local port so this works with Istio out of the box.
In addition, the Operator uses the health check port to determine the status of a cluster, so this needs to be
exposed so that the Operator can reach Coherence Pods.</p>

<p>The Coherence localhost property can be set to the name of the Pod.
This is easily done using the container environment variables, which the Operator does automatically.</p>

<p>Coherence clusters are run as a StatefulSet in Kubernetes. This means that the Pods are configured with a host name
and a subdomain based on the name of the StatefulSet headless service name, and it is this name that should be used
to access Pods.
For example for a Coherence resource named <code>storage</code> the Operator will create a StatefulSet named <code>storgage</code> with a
headless service named <code>storage-sts</code>. Each Pod in a StatefulSet is numbered with a fixed identity, so the first Pod
in this cluster will be <code>storage-0</code>. The Pod has a number of DNS names that it is reachable with, but the fully
qualified name using the headless service will be <code>storage-0.storage-sts</code> or storage-0.storage-sts.&lt;namespace&gt;.svc`.</p>

<p>By default, the Operator will expose all the ports configured for the <code>Coherence</code> resource on the StatefulSet headless
service. This allows Coherence Extend and gRPC clients to use this service name as the WKA address when using the
Coherence NameService to lookup endpoints (see the client example below).</p>

</div>

<h3 id="_prerequisites">Prerequisites</h3>
<div class="section">
<p>The instructions assume that you are using a Kubernetes cluster with Istio installed and configured already.</p>


<h4 id="_enable_istio_strict_mode">Enable Istio Strict Mode</h4>
<div class="section">
<p>For this example we make Istio run in "strict" mode so that it will not allow any traffic between Pods outside the
Envoy proxy.
If other modes are used, such as permissive, then Istio allows Pod to Pod communication so a cluster may appear to work
in permissive mode, when it would not in strict mode.</p>

<p>To set Istio to strict mode create the following yaml file.</p>

<markup
lang="yaml"
title="istio-strict.yaml"
>apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: "default"
spec:
  mtls:
    mode: STRICT</markup>

<p>Install this yaml into the Istio system namespace with the following command:</p>

<markup
lang="bash"

>kubectl -n istio-system apply istio-strict.yaml</markup>

</div>
</div>

<h3 id="_using_the_coherence_operator_with_istio">Using the Coherence operator with Istio</h3>
<div class="section">
<p>To use Coherence operator with Istio, you can deploy the operator into a namespace which has Istio automatic sidecar
injection enabled.
Before installing the operator, create the namespace in which you want to run the Coherence operator and label it for
automatic injection.</p>

<markup
lang="bash"

>kubectl create namespace coherence
kubectl label namespace coherence istio-injection=enabled</markup>

<p>Istio Sidecar AutoInjection is done automatically when you label the coherence namespace with istio-injection.</p>


<h4 id="_exclude_the_operator_web_hook_from_the_envoy_proxy">Exclude the Operator Web-Hook from the Envoy Proxy</h4>
<div class="section">
<p>The Coherence Operator uses an admissions web-hook, which Kubernetes will call to validate Coherence resources.
This web-hook binds to port <code>9443</code> in the Operator Pods and is already configured to use TLS as is standard for
Kubernetes admissions web-hooks. If this port is routed through the Envoy proxy Kubernetes will be unable to
access the web-hook.</p>

<p>The Operator yaml manifests and Helm chart already add the <code>traffic.sidecar.istio.io/excludeInboundPorts</code> annotation
to the Operator Pods. This should exclude the web-hook port from being Istio.</p>

<p>Another way to do this is to add a <code>PeerAuthentication</code> resource to the Operator namespace.</p>

<p><strong>Before installing the Operator</strong>, create the following <code>PeerAuthentication</code> yaml.</p>

<markup
lang="yaml"
title="istio-operator.yaml"
>apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: "coherence-operator"
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: coherence-operator
      app.kubernetes.io/instance: coherence-operator-manager
      app.kubernetes.io/component: manager
  mtls:
    mode: STRICT
  portLevelMtls:
    9443:
      mode: PERMISSIVE</markup>

<p>Then install this <code>PeerAuthentication</code> resource into the same namespace that the Operator will be installed into.
For example, if the Operator will be in the <code>coherence</code> namespace:</p>

<markup
lang="bash"

>kubectl -n coherence apply istio-operator.yaml</markup>

<p>You can then install the operator using your preferred method in the
Operator <router-link to="/docs/installation/01_installation">Installation Guide</router-link>.</p>

<p>After installed operator, use the following command to confirm the operator is running:</p>

<markup
lang="bash"

>kubectl get pods -n coherence

NAME                                                     READY   STATUS    RESTARTS   AGE
coherence-operator-controller-manager-7d76f9f475-q2vwv   2/2     Running   1          17h</markup>

<p>The output should show 2/2 in READY column, meaning there are 2 containers running in the Operator pod.
One is Coherence Operator and the other is Envoy Proxy.</p>

<p>If we use the Istio Kiali addon to visualize Istio we can see the Operator in the list of applications</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/images/kiali-operator-app.png" alt="kiali operator app"width="1024" />
</v-card-text>
</v-card>

<p>We can also see on the detailed view, that the Operator talks to the Kubernetes API server</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/images/kiali-operator-app-graph.png" alt="kiali operator app graph"width="1024" />
</v-card-text>
</v-card>

</div>
</div>

<h3 id="_creating_a_coherence_cluster_with_istio">Creating a Coherence cluster with Istio</h3>
<div class="section">
<p>You can configure a cluster to run with Istio automatic sidecar injection enabled. Before creating the cluster,
create the namespace in which the cluster will run and label it for automatic injection.</p>

<markup
lang="bash"

>kubectl create namespace coherence-example
kubectl label namespace coherence-example istio-injection=enabled</markup>

<p>Now create a Coherence resource as normal, there is no additional configuration required to work in Istio.</p>

<p>For example using the yaml below to create a three member cluster with management and metrics enabled:</p>

<markup
lang="yaml"
title="storage.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  replicas: 3
  image: ghcr.io/oracle/coherence-ce:22.06.6
  labels:
    app: storage      <span class="conum" data-value="1" />
    version: 1.0.0    <span class="conum" data-value="2" />
  coherence:
    management:
      enabled: true
    metrics:
      enabled: true
  ports:
    - name: management  <span class="conum" data-value="3" />
    - name: metrics
    - name: extend
      port: 20000
      appProtocol: tcp  <span class="conum" data-value="4" />
    - name: grpc-proxy
      port: 1408
      appProtocol: grpc <span class="conum" data-value="5" /></markup>

<ul class="colist">
<li data-value="1">Istio prefers applications to have an <code>app</code> label</li>
<li data-value="2">Istio prefers applications to have a <code>version</code> label</li>
<li data-value="3">The Coherence Pods will expose ports for Management over REST, metrics, a Coherence*Extend proxy and a gRPC proxy</li>
<li data-value="4">The Operator will set the <code>appProtocol</code> for the management and metrics ports to <code>http</code>, but the Extend port must be
set manually to <code>tcp</code> so that Istio knows what sort of traffic is being used by that port</li>
<li data-value="5">The gRPC port&#8217;s <code>appProtocol</code> field is set to <code>grpc</code></li>
</ul>
<p>Using the Kiali console, we can now see two applications, the Coherence Operator in the "coherence" namespace
and the "storage" application in the "coherence-example" namespace.</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/images/kiali-storage-app.png" alt="kiali storage app"width="1024" />
</v-card-text>
</v-card>

<p>If we look at the graph view we can see all the traffic between the different parts of the system</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/images/kiali-post-deploy.png" alt="kiali post deploy"width="1024" />
</v-card-text>
</v-card>

<ul class="ulist">
<li>
<p>We can see the Kubernetes API server accessing the Operator web-hook to validate the yaml</p>

</li>
<li>
<p>We can see tge storage pods (the box marked "storage 1.0.0") communicate with each other via the storage-sts service to from a Coherence cluster</p>

</li>
<li>
<p>We can see the storage pods communicate with the Operator REST service to request their Coherence site and rack labels</p>

</li>
<li>
<p>We can see the Operator ping the storage pods health endpoints via the storage-sts service</p>

</li>
</ul>
<p>All of this traffic is using mTLS controlled by Istio</p>

</div>

<h3 id="_coherence_clients_running_in_kubernetes">Coherence Clients Running in Kubernetes</h3>
<div class="section">
<p>Coherence Extend clients and gRPC clients running inside the cluster will also work with Istio.</p>

<p>For this example the clients will be run in the <code>coherence-client</code> namespace, so it needs to be
created and labelled so that Istio injection works in that namespace.</p>

<markup
lang="bash"

>kubectl create namespace coherence-client
kubectl label namespace coherence-client istio-injection=enabled</markup>

<p>To simulate a client application a <code>CoherenceJob</code> resource will be used with different configurations
for the different types of client.</p>

<p>The simplest way to configure a Coherence extend client in a cache configuration file is a default configuration
similar to that shown below. No ports or addresses need to be configured. Coherence will use the JVM&#8217;s configured
cluster name and well know addresses to locate to look up the Extend endpoints using the Coherence NameService.</p>

<markup
lang="xml"

>&lt;remote-cache-scheme&gt;
  &lt;scheme-name&gt;thin-remote&lt;/scheme-name&gt;
  &lt;service-name&gt;RemoteCache&lt;/service-name&gt;
  &lt;proxy-service-name&gt;Proxy&lt;/proxy-service-name&gt;
&lt;/remote-cache-scheme&gt;</markup>

<p>We can configure a <code>CoherenceJob</code> to run an Extend client with this configuration as shown below:</p>

<markup
lang="yaml"
title="extend-client.yaml"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceJob
metadata:
  name: client
spec:
  image: ghcr.io/oracle/coherence-ce:22.06.6  <span class="conum" data-value="1" />
  restartPolicy: Never
  cluster: storage  <span class="conum" data-value="2" />
  coherence:
    wka:
      addresses:
        - "storage-sts.coherence-example.svc"  <span class="conum" data-value="3" />
  application:
    type: operator  <span class="conum" data-value="4" />
    args:
      - sleep
      - "300s"
  env:
    - name: COHERENCE_CLIENT    <span class="conum" data-value="5" />
      value: "remote"
    - name: COHERENCE_PROFILE   <span class="conum" data-value="6" />
      value: "thin"</markup>

<ul class="colist">
<li data-value="1">The client will use the CE image published on GitHub, which will use the default cache configuration file from Coherence jar.</li>
<li data-value="2">The cluster name must be set to the cluster name of the cluster started above, in this case <code>storage</code></li>
<li data-value="3">The WKA address needs to be set to the DNS name of the headless service for the storage cluster created above. As this
Job is running in a different name space this is the fully qualified name <code>&lt;service-name&gt;.&lt;namespace&gt;.svc</code> which is <code>storage-sts.coherence-example.svc</code></li>
<li data-value="4">Instead of running a normal command this Job will run the Operator&#8217;s <code>sleep</code> command and sleep for <code>300s</code> (300 seconds).</li>
<li data-value="5">The <code>COHERENCE_CLIENT</code> environment variable value of <code>remote</code> sets the Coherence cache configuration to be an Extend client using the NameService</li>
<li data-value="6">The <code>COHERENCE_PROFILE</code> environment variable value of <code>thin</code> sets the Coherence cache configuration not to use a Near Cache.</li>
</ul>
<p>The yaml above can be deployed into Kubernetes:</p>

<markup
lang="bash"

>kubectl -n coherence-client apply -f extend-client.yaml</markup>

<markup
lang="bash"

>$ kubectl -n coherence-client get pod
NAME           READY   STATUS    RESTARTS   AGE
client-qgnw5   2/2     Running   0          80s</markup>

<p>The Pod is now running but not doing anything, just sleeping.
If we look at the Kiali dashboard we can see the client application started and communicated wth the Operator.</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/images/kiali-client-started-graph.png" alt="kiali client started graph"width="1024" />
</v-card-text>
</v-card>

<p>We can use this sleeping Pod to exec into and run commands. In this case we will create a Coherence QueryPlus
client and run some CohQL commands. The command below will exec into the sleeping Pod.</p>

<markup
lang="bash"

>kubectl -n coherence-client exec -it client-qgnw5 -- /coherence-operator/utils/runner queryplus</markup>

<p>A QueryPlus client will be started and eventually display the <code>CohQL&gt;</code> prompt.</p>

<markup
lang="bash"

>Coherence Command Line Tool

CohQL&gt;</markup>

<p>A simple command to try is just creating a cache, so at the prompt type the command <code>create cache test</code> which will
create a cache named <code>test</code>. If all is configured correctly this client will connect to the cluster over Extend
and create the cache called <code>test</code> and return to the <code>CohQL</code> prompt.</p>

<markup
lang="bash"

>Coherence Command Line Tool

CohQL&gt; create cache test</markup>

<p>We can also try selecting data from the cache using the CohQL query <code>select * from test</code>
(which will return nothing as the cache is empty).</p>

<markup
lang="bash"

>CohQL&gt; select * from test
Results

CohQL&gt;</markup>

<p>If we now look at the Kiali dashboard we can see that the client has communicated with the storage cluster.
All of this communication was using mTLS but without configuring Coherence to use TLS.</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/images/kiali-client-storage.png" alt="kiali client storage"width="1024" />
</v-card-text>
</v-card>

<p>To exit from the <code>CohQL&gt;</code> prompt type the <code>bye</code> command.</p>

<p>Coherence Extend clients can connect to the cluster also using Istio to provide mTLS support.
Coherence clusters work with mTLS and Coherence clients can also support TLS through the Istio Gateway with TLS
termination to connect to Coherence cluster running inside kubernetes.
For example, you can apply the following Istio Gateway and Virtual Service in the namespace of the Coherence cluster.
Before applying the gateway, create a secret for the credential from the certificate and key
(e.g. server.crt and server.key) to be used by the Gateway:</p>

</div>

<h3 id="_coherence_clients_running_outside_kubernetes">Coherence Clients Running Outside Kubernetes</h3>
<div class="section">
<p>Coherence clients running outside the Kubernetes can be configured to connect to a Coherence cluster inside
Kubernetes using any of the ingress or gateway features of Istio and Kubernetes.
All the different ways to do this are beyond the scope of this simple example as there are many and they
depend on the versions of Istio and Kubernetes being used.</p>

<p>When connecting Coherence Extend or gRPC clients from outside Kubernetes, the Coherence NameService cannot be used
by clients to look up the endpoints. The clients must be configured with fixed endpoints using the hostnames and ports
of the configured ingress or gateway services.</p>

</div>
</div>
</doc-view>
