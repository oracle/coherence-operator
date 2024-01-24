<doc-view>

<h2 id="_running_coherence_with_istio">Running Coherence with Istio</h2>
<div class="section">
<p>This example shows how to deploy a simple Coherence cluster in Kubernetes with Istio.</p>

<p>Coherence can be configured to work with <a id="" title="" target="_blank" href="https://istio.io">Istio</a>, even if Istio is configured in Strict Mode.
Coherence caches can be accessed from inside or outside the Kubernetes cluster via Coherence*Extend, REST,
and other supported Coherence clients.
Although Coherence itself can be configured to use TLS, when using Istio Coherence cluster members and clients can
just use the default socket configurations and Istio will control and route all the traffic over mTLS.</p>

</div>

<h2 id="_how_does_coherence_work_with_istio">How Does Coherence Work with Istio?</h2>
<div class="section">
<p>Istio is a "Service Mesh" so the clue to how Istio works in Kubernetes is in the name, it relies on the configuration
of Kubernetes Services.
This means that any ports than need to be accessed in Pods, including those using in "Pod to Pod" communication
must be exposed via a Service. Usually a Pod can reach any port on another Pod even if it is not exposed in the
container spec, but this is not the case when using Istio as only ports exposed by the Envoy proxy are allowed.</p>

<p>For Coherence cluster membership, this means the cluster port and the local port must be exposed on a Service.
To do this the local port must be configured to be a fixed port instead of the default ephemeral port.
The default cluster port is <code>7574</code> and there is no reason to ever change this when running in containers.
A fixed local port has to be configured for Coherence to work with Istio out of the box.
Additional ports, management port, metrics port, etc. also need to be exposed if they are being used.</p>

<p>Ideally, Coherence clusters are run as a StatefulSet in Kubernetes.
This means that the Pods are configured with a host name and a subdomain based on the name of the StatefulSet
headless service name, and it is this name that should be used to access Pods.</p>


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
</div>

<h2 id="_create_a_coherence_cluster">Create a Coherence Cluster</h2>
<div class="section">
<p>The best way to run Coherence cluster members is to use a StatefulSet. Multiple StatefulSets can be created that
are all part of the same Coherence cluster.</p>

<p>In this example we will run a Coherence cluster using the CE image. This image starts Coherence with health
checks enabled on port 6676,
an Extend proxy listening on port 20000, a gRPC proxy on port 1408, the cluster port set to 7574.
We will also enable Coherence Management over REST on port 30000, and metrics on port 9612.
We will set the Coherence local port to a fixed value of 7575.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Istio has a few requirements for how Kubernetes resources are configured.
One of those is labels, where an <code>app</code> and <code>version</code> label are required to specify the application name
that the resource is part of and the version of that application.
All the resources in this example contains those labels.</p>
</p>
</div>

<h3 id="_cluster_discovery_service">Cluster Discovery Service</h3>
<div class="section">
<p>For Coherence cluster discovery to work in Kubernetes we have to configure Coherence well-known-addresses which
requires a headless service. We cannot use the same headless service then we will create for the StatefulSet because
the WKA service must have the <code>publishNotReadyAddresses</code> field set to <code>true</code>, whereas the StatefulSet service does not.
We would not want the ports accessed via the StatefulSet service to route to unready Pods, but for cluster discovery
we must allow unready Pods to be part of the Service.</p>

<p>The discovery service can be created with yaml like that shown below.</p>

<markup
lang="yaml"
title="wka-service.yaml"
>apiVersion: v1
kind: Service
metadata:
  name: storage-wka    <span class="conum" data-value="1" />
spec:
  clusterIP: None
  publishNotReadyAddresses: true  <span class="conum" data-value="2" />
  selector:                       <span class="conum" data-value="3" />
    app: my-coherence-app
    version: 1.0.0
  ports:
    - name: coherence    <span class="conum" data-value="4" />
      port: 7574
      targetPort: coherence
      appProtocol: tcp</markup>

<ul class="colist">
<li data-value="1">The service name is <code>storeage-wka</code> and this will be used to configure the Coherence WKA address in the cluster.</li>
<li data-value="2">The <code>publishNotReadyAddresses</code> field must be set to <code>true</code></li>
<li data-value="3">The <code>selector</code> is configured to match a sub-set of the Pod labels configured in the StatefulSet</li>
<li data-value="4">We do not really need or care about the port for the cluster discovery service, but all Kubernetes services must have
at least one port, so here we use the cluster port. We could use any random port, even one that nothing is listening on</li>
</ul>
</div>

<h3 id="_statefulset_headless_service">StatefulSet Headless Service</h3>
<div class="section">
<p>All StatefulSets require a headless Service creating and the name of this Service is specified in the StatefulSet spec.
All the ports mentioned above will be exposed on this service.
The yaml for the service could look like this:</p>

<markup
lang="yaml"
title="storage-service.yaml"
>apiVersion: v1
kind: Service
metadata:
  name: storage-headless
spec:
  clusterIP: None
  selector:
    app: my-coherence-app  <span class="conum" data-value="1" />
    version: 1.0.0
  ports:
    - name: coherence              <span class="conum" data-value="2" />
      port: 7574
      targetPort: coherence
      appProtocol: tcp
    - name: coh-local              <span class="conum" data-value="3" />
      port: 7575
      targetPort: coh-local
      appProtocol: tcp
    - name: extend-proxy           <span class="conum" data-value="4" />
      port: 20000
      targetPort: extend-proxy
      appProtocol: tcp
    - name: grpc-proxy             <span class="conum" data-value="5" />
      port: 1408
      targetPort: grpc-proxy
      appProtocol: grpc
    - name: management             <span class="conum" data-value="6" />
      port: 30000
      targetPort: management
      appProtocol: http
    - name: metrics                <span class="conum" data-value="7" />
      port: 9612
      targetPort: metrics
      appProtocol: http</markup>

<ul class="colist">
<li data-value="1">The selector labels will match a sub-set of the labels specified for the Pods in the StatefulSet</li>
<li data-value="2">The Coherence cluster port 7574 is exposed with the name <code>coherence</code> mapping to the container port in the StatefulSet named <code>coherence</code>.
This port has an <code>appProtocol</code> of <code>tcp</code> to tell Istio that the port traffic is raw TCP traffic.</li>
<li data-value="3">The Coherence local port 7575 is exposed with the name <code>coh-local</code> mapping to the container port in the StatefulSet named <code>coh-local</code>
This port has an <code>appProtocol</code> of <code>tcp</code> to tell Istio that the port traffic is raw TCP traffic.</li>
<li data-value="4">The Coherence Extend proxy port 20000 is exposed with the name <code>extend-proxy</code> mapping to the container port in the StatefulSet named <code>extend-proxy</code>
This port has an <code>appProtocol</code> of <code>tcp</code> to tell Istio that the port traffic is raw TCP traffic.</li>
<li data-value="5">The Coherence gRPC proxy port 1408 is exposed with the name <code>grpc-proxy</code> mapping to the container port in the StatefulSet named <code>grpc-proxy</code>
This port has an <code>appProtocol</code> of <code>grpc</code> to tell Istio that the port traffic is gRPC traffic.</li>
<li data-value="6">The Coherence Management over REST port 30000 is exposed with the name <code>management</code> mapping to the container port in the StatefulSet named <code>management</code>
This port has an <code>appProtocol</code> of <code>http</code> to tell Istio that the port traffic is http traffic.</li>
<li data-value="7">The Coherence Metrics port 9612 is exposed with the name <code>metrics</code> mapping to the container port in the StatefulSet named <code>metrics</code>
This port has an <code>appProtocol</code> of <code>http</code> to tell Istio that the port traffic is http traffic.</li>
</ul>
<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Istio requires ports to specify the protocol used for their traffic, and this can be done in two ways.
Either using the <code>appProtocol</code> field for the ports, as shown above.
Or, prefix the port name with the protocol, so instead of <code>management</code> the port name would be <code>http-management</code></p>
</p>
</div>
</div>

<h3 id="_the_statefulset">The StatefulSet</h3>
<div class="section">
<p>With the two Services defined, the StatefulSet can now be configured.
Istio</p>

<markup
lang="yaml"
title="storage.yaml"
>apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: storage
  labels:
    app: my-coherence-app
    version: 1.0.0
spec:
  selector:
    matchLabels:
        app: my-coherence-app
        version: 1.0.0
  serviceName: storage-headless  <span class="conum" data-value="1" />
  replicas: 3
  podManagementPolicy: Parallel
  updateStrategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1
  template:
    metadata:
      labels:
        app: my-coherence-app
        version: 1.0.0
    spec:
      containers:
        - name: coherence
          image: ghcr.io/oracle/coherence-ce:22.06.7   <span class="conum" data-value="2" />
          env:
            - name: COHERENCE_CLUSTER          <span class="conum" data-value="3" />
              value: "test-cluster"
            - name: NAMESPACE                  <span class="conum" data-value="4" />
              valueFrom:
                fieldRef:
                  fieldPath: "metadata.namespace"
            - name: COHERENCE_WKA                   <span class="conum" data-value="5" />
              value: "storage-wka.${NAMESPACE}.svc"
            - name: COHERENCE_LOCALPORT        <span class="conum" data-value="6" />
              value: "7575"
            - name: COHERENCE_LOCALHOST        <span class="conum" data-value="7" />
              valueFrom:
                fieldRef:
                  fieldPath: "metadata.name"
            - name: COHERENCE_MACHINE          <span class="conum" data-value="8" />
              valueFrom:
                fieldRef:
                  fieldPath: "spec.nodeName"
            - name: COHERENCE_MEMBER           <span class="conum" data-value="9" />
              valueFrom:
                fieldRef:
                  fieldPath: "metadata.name"
            - name: COHERENCE_EXTEND_PORT
              value: "20000"
            - name: COHERENCE_GRPC_SERVER_PORT
              value: "1408"
          ports:
           - name: coherence         <span class="conum" data-value="10" />
             containerPort: 7574
           - name: coh-local
             containerPort: 7575
           - name: extend-proxy
             containerPort: 20000
           - name: grpc-proxy
             containerPort: 1408
           - name: management
             containerPort: 30000
           - name: metrics
             containerPort: 9162
          readinessProbe:            <span class="conum" data-value="11" />
            httpGet:
              path: "/ready"
              port: 6676
              scheme: "HTTP"
          livenessProbe:
            httpGet:
              path: "/healthz"
              port: 6676
              scheme: "HTTP"</markup>

<ul class="colist">
<li data-value="1">All StatefulSets require a headless service, in this case the service will be named <code>storage-headless</code> to match the
service above</li>
<li data-value="2">This example is using the CE 22.06 image</li>
<li data-value="3">The <code>COHERENCE_CLUSTER</code> environment variable sets the Coherence cluster name to <code>test-cluster</code></li>
<li data-value="4">The <code>NAMESPACE</code> environment variable contains the namespace the StatefulSet is deployed into.
The value is taken from the <code>matadata.namespace</code> field of the Pod. This is then used to create a fully qualified
well known address value</li>
<li data-value="5">The <code>COHERENCE_WKA</code> environment variable sets address Coherence uses to perform a DNS lookup for cluster member IP
addresses. In this case we use the name of the WKA service created above combined with the <code>NAMESPACE</code> environment
variable to give a fully qualified service name.</li>
<li data-value="6">The <code>COHERENCE_LOCALPORT</code> environment variable sets the Coherence localport to 7575, which matches what was exposed
in the Service ports and container ports</li>
<li data-value="7">The <code>COHERENCE_LOCAHOST</code> environment variable sets the hostname that Coherence binds to, in this case it will be
the same as the Pod name by using the "valueFrom" setting to get the value from the Pod&#8217;s <code>metadata.name</code> field</li>
<li data-value="8">It is best practice to use the <code>COHERENCE_MACHINE</code> environment variable to set the Coherence machine label to the
Kubernetes Node name. The machine name is used by Coherence when assigning backup partitions, so a backup of a partition will
not be on the same Node as the primary owner of the partition.
the same as the Pod name by using the "valueFrom" setting to get the value from the Pod&#8217;s <code>metadata.name</code> field</li>
<li data-value="9">It is best practice to use the <code>COHERENCE_MEMBER</code> environment variable to set the Coherence member name to the
Pod name.</li>
<li data-value="10">All the ports required are exposed as container ports. The names must correspond to the names used for the container ports in the Service spec.</li>
<li data-value="11">As we are using Coherence CE 22.06 we can use Coherence built in health check endpoints for the readiness and liveness probes.</li>
</ul>
<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>The example above is using Coherence 22.06 which has built in health checks and health check endpoints which can be used
as readiness and liveness probes in Kubernetes. These endpoints are only available if you start Coherence correctly
using the Bootstrap API introduced in 22.06.</p>

<p>Start Coherence using <code>com.tangosol.net.Coherence</code> as the main class.</p>

<markup
lang="bash"

>java -cp coherence.jar com.tangosol.net.Coherence</markup>

<p>Start Coherence in code:</p>

<markup
lang="java"

>Coherence coherence = Coherence.clusterMember().start().join();</markup>

<p>See the Coherence <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.2206/develop-applications/starting-and-stopping-cluster-members.html#GUID-A74280C1-AB23-48DB-8C0D-E494230ABFCA">Bootstrap API</a>
documentation for more details</p>
</p>
</div>
</div>

<h3 id="_deploy_the_cluster">Deploy the Cluster</h3>
<div class="section">
<p>We will deploy the cluster into a Kubernetes namespace names <code>coherence</code>.
Before deploying the cluster we need to ensure it has been labeled so that Istio will inject the
Envoy proxy sidecar into the Pods.</p>

<markup
lang="bash"

>kubectl create namespace coherence
kubectl label namespace coherence istio-injection=enabled</markup>

<p>To deploy the cluster we just apply all three yaml files to Kubernetes.
We could combine them into  a single yaml file if we wanted to.</p>

<markup
lang="bash"

>kubectl -n coherence apply -f wka-service.yaml
kubectl -n coherence apply -f storage-service.yaml
kubectl -n coherence apply -f storage.yaml</markup>

<p>If we list the services, we see the two services we created</p>

<markup
lang="bash"

>$ kubectl get svc
NAME               TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)                                                   AGE
storage-headless   ClusterIP   None         &lt;none&gt;        7574/TCP,7575/TCP,20000/TCP,1408/TCP,30000/TCP,9612/TCP   37m
storage-wka        ClusterIP   None         &lt;none&gt;        7574/TCP                                                  16m</markup>

<p>If we list the Pods, we see three Pods, as the StatefulSet replicas field is set to three.</p>

<markup
lang="bash"

>$ kubectl get pod
NAME        READY   STATUS    RESTARTS   AGE
storage-0   2/2     Running   0          7m47s
storage-1   2/2     Running   0          7m47s
storage-2   2/2     Running   0          7m47s</markup>

<p>We can use Istio&#8217;s Kiali dashboard to visualize the cluster we created.</p>

<p>We labelled the resources with the <code>app</code> label with a value of <code>my-coherence-app</code> and we can see this application
in the Kiali dashboard. The graph shows the cluster member Pods communicating with each other via the <code>storage-headless</code>
service. The padlock icons show that this traffic is using mTLS even though Coherence has not been configured with TLS,
this is being provided by Istio.</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/images/kiali-cluster-start.png" alt="kiali cluster start"width="1024" />
</v-card-text>
</v-card>

</div>
</div>

<h2 id="_coherence_clients">Coherence Clients</h2>
<div class="section">
<p>Coherence clients (Extend or gRPC) can be configured to connect to the Coherence cluster.</p>


<h3 id="_extend_proxy_configuration">Extend Proxy Configuration</h3>
<div class="section">
<p>To work correctly with Istio a Coherence Extend proxy in the server&#8217;s cache configuration file must be configured
to use a fixed port.
For example, the XML snippet below configures the proxy to bind to all interfaces (<code>0.0.0.0</code>) on port 20000.</p>

<markup
lang="xml"

>    &lt;proxy-scheme&gt;
      &lt;service-name&gt;Proxy&lt;/service-name&gt;
      &lt;acceptor-config&gt;
        &lt;tcp-acceptor&gt;
          &lt;local-address&gt;
            &lt;address system-property="coherence.extend.address"&gt;0.0.0.0&lt;/address&gt;
            &lt;port system-property="coherence.extend.port"&gt;20000&lt;/port&gt;
          &lt;/local-address&gt;
        &lt;/tcp-acceptor&gt;
      &lt;/acceptor-config&gt;
      &lt;autostart&gt;true&lt;/autostart&gt;
    &lt;/proxy-scheme&gt;</markup>

<p>The port could be changed by setting the <code>COHERENCE_EXTEND_PORT</code> environment variable in the server yaml.</p>

<markup
lang="yaml"

>    spec:
      containers:
        - name: coherence
          image: ghcr.io/oracle/coherence-ce:22.06.7
          env:
            - name: COHERENCE_EXTEND_PORT
              value: "20001"</markup>

<p>The Extend port should be exposed on the StatefulSet (as shown in the <router-link to="#sts" @click.native="this.scrollFix('#sts')">StatefulSet yaml</router-link> above) and
on the StatefulSet headless service so that clients can discover it and connect
to it (as shown in the <router-link to="#headless-svc" @click.native="this.scrollFix('#headless-svc')">Service yaml</router-link> above).</p>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p>The default cache configuration file used by Coherence, and used in the Coherence images published on GitHub,
contains an Extend Proxy service that uses the <code>COHERENCE_EXTEND_PORT</code> environment variable to set the port.</p>
</p>
</div>
</div>

<h3 id="_grpc_proxy_configuration">gRPC Proxy Configuration</h3>
<div class="section">
<p>The Coherence gRPC proxy binds to an ephemeral port by default.
This port can be changed by using the <code>COHERENCE_GRPC_SERVER_PORT</code> environment variable;</p>

<markup
lang="yaml"

>    spec:
      containers:
        - name: coherence
          image: ghcr.io/oracle/coherence-ce:22.06.7
          env:
            - name: COHERENCE_GRPC_SERVER_PORT
              value: "1408"</markup>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p>The default configuration used by Coherence images published on GitHub sets the gRPC port to 1408.</p>
</p>
</div>
<p>Once the server StatefulSet and Service have been properly configured the clients can be configured.
The options available for this will depend on where the client will run.</p>

</div>

<h3 id="_clients_inside_kubernetes">Clients Inside Kubernetes</h3>
<div class="section">
<p>If the clients are also inside the cluster they can be configured to connect using the StatefulSet headless service
as the hostname for the proxy endpoints.
There are two options for configuring Extend and</p>

<ul class="ulist">
<li>
<p>Clients inside Kubernetes can also use the minimal Coherence NameService configuration where the
StatefulSet service name is used as the client&#8217;s WKA address and the same cluster name is configured.</p>

</li>
<li>
<p>Clients external to the Kubernetes cluster can be configured using any of the ingress or gateway features of Istio and Kubernetes.
All the different ways to do this are beyond the scope of this simple example as there are many, and they
depend on the versions of Istio and Kubernetes being used.</p>

</li>
</ul>

<h4 id="_build_a_client_image">Build a Client Image</h4>
<div class="section">
<p>For this example we need a simple client image that can be run with different configurations.
Instead of building an application we will use a Coherence Image from GitHub combined with the
utilities from the Coherence Operator.</p>

<p>The simple Dockerfile below is a multistage build file.
It uses the Operator image as a "builder" and then the Coherence image as the base.
Various utilities are copied from the Operator image into the base.</p>

<markup

title="Dockerfile"
>FROM ghcr.io/oracle/coherence-operator:3.3.3 AS Builder

FROM ghcr.io/oracle/coherence-ce:22.06.7
COPY --from=Builder /files /files
COPY --from=Builder /files/lib/coherence-operator.jar /app/libs/coherence-operator.jar
COPY coherence-java-client-22.06.7.jar /app/libs/coherence-java-client-22.06.7.jar

ENTRYPOINT ["files/runner"]
CMD ["-h"]</markup>

<p>As we are going to show both the Coherence Extend client and gRPC client we need to add the
Coherence gRPC client jar. We can download this with <code>curl</code> to the same directory as the Dockerfile.</p>

<markup
lang="bash"

>curl -s https://repo1.maven.org/maven2/com/oracle/coherence/ce/coherence-java-client/22.06.7/coherence-java-client-22.06.7.jar \
  -o coherence-java-client-22.06.7.jar</markup>

<p>Build the image with the following command:</p>

<markup
lang="bash"

>docker build -t coherence-client:1.0.0 -f Dockerfile .</markup>

<p>There will now be an imaged named <code>coherence-client:1.0.0</code> which can be pushed somewhere Kubernetes can see it.
We will use this example below.</p>

</div>

<h4 id="_using_the_coherence_nameservice_configuration">Using the Coherence NameService Configuration</h4>
<div class="section">
<p>The minimal configuration in a client&#8217;s cache configuration file is shown below.
This configuration will use the Coherence NameService to look up the endpoints for the Extend Proxy services
running in the Coherence cluster.</p>

<markup
lang="xml"

>&lt;remote-cache-scheme&gt;
  &lt;scheme-name&gt;thin-remote&lt;/scheme-name&gt;
  &lt;service-name&gt;RemoteCache&lt;/service-name&gt;
  &lt;proxy-service-name&gt;Proxy&lt;/proxy-service-name&gt;
&lt;/remote-cache-scheme&gt;</markup>

<p>For the NameService to work in Kubernetes, the client must be configured with the same cluster name, the same
well known addresses and same cluster port as the server. When using Istio the server&#8217;s cluster port, local port
and Extend port should be exposed on the StatefulSet headless service.
The client&#8217;s well known address is then set to the qualified Kubernetes DNS name for the server&#8217;s StatefulSet headless service.</p>

<p>These can all be set using environment variables in the yaml for the client.</p>

<p>For example, assuming the client will connect to the Coherence cluster configured in the <router-link to="#sts" @click.native="this.scrollFix('#sts')">StatefulSet</router-link> above:</p>

<markup
lang="yaml"

>  env:
    - name: COHERENCE_CLUSTER             <span class="conum" data-value="1" />
      value: "test-cluster"
    - name: COHERENCE_WKA                 <span class="conum" data-value="2" />
      value: "storage-headless.coherence.svc"</markup>

<ul class="colist">
<li data-value="1">The cluster name is set to <code>test-cluster</code> the same as the StatefulSet</li>
<li data-value="2">The <code>COHERENCE_WKA</code> value is set to the DNS name of the StatefulSet headless service, which has the format
<code>&lt;service-name&gt;.&lt;namespace&gt;.svc</code> so in this case <code>storage-headless.coherence.svc</code></li>
</ul>
</div>

<h4 id="_run_an_extend_client_pod">Run an Extend Client Pod</h4>
<div class="section">
<p>Using the <code>coherence-client:1.0.0</code> image created above, we can run a simple Coherence client Pod.</p>

<markup
lang="yaml"
title="extend-client-pod.yaml"
>apiVersion: v1
kind: Pod
metadata:
  name: client
  labels:
    app: coherence-client
    version: 1.0.0
spec:
  containers:
    - name: coherence
      image: coherence-client:1.0.0  <span class="conum" data-value="1" />
      command:
        - /files/runner   <span class="conum" data-value="2" />
        - sleep
        - "15m"
      env:
        - name: COHERENCE_CLUSTER              <span class="conum" data-value="3" />
          value: "test-cluster"
        - name: COHERENCE_WKA                  <span class="conum" data-value="4" />
          value: "storage-headless.coherence.svc"
        - name: COHERENCE_CLIENT               <span class="conum" data-value="5" />
          value: "remote"</markup>

<ul class="colist">
<li data-value="1">The container image is set to the client image built above <code>coherence-client:1.0.0</code></li>
<li data-value="2">The command line the container will run is <code>/files/runner sleep 15m</code> which will just sleep for 15 minutes</li>
<li data-value="3">The Coherence cluster name is set to the same name as the server deployed above in the <router-link to="#sts" @click.native="this.scrollFix('#sts')">StatefulSet yaml</router-link></li>
<li data-value="4">The WKA address is set to the StatefulSet&#8217;s headless service name <code>storage-headless.coherence.svc</code></li>
<li data-value="5">For this example the <code>COHERENCE_CLIENT</code> which sets the default cache configuration file to run as an Extend
client, using the NameService to look up the proxies.</li>
</ul>
<p>We can deploy the client into Kubernetes</p>

<markup
lang="bash"

>kubectl -n coherence apply -f extend-client-pod.yaml</markup>

<p>We deployed the client into the same namespace as the cluster, we could easily have deployed it to another namespace.
If we list the Pods we will see the cluster and the client. All Pods has two containers, one being the Istio side-car.</p>

<markup
lang="bash"

>$ k -n coherence get pod
NAME        READY   STATUS    RESTARTS   AGE
storage-0   2/2     Running   0          105m
storage-1   2/2     Running   0          105m
storage-2   2/2     Running   0          105m
client      2/2     Running   0          8m27s</markup>

<p>Now we can exec into the Pod and start a Coherence QueryPlus console session using the following command:</p>

<markup
lang="bash"

>kubectl -n coherence exec -it client -- /files/runner queryplus</markup>

<p>The QueryPlus session will start and eventually display the <code>CohQL&gt;</code> prompt:</p>

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

<p>If we now look at the Kiali dashboard we can see that the client application has communicated with the storage cluster.
All of this communication was using mTLS but without configuring Coherence to use TLS.</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/images/kiali-client-remote-app.png" alt="kiali client remote app" />
</v-card-text>
</v-card>

<p>If we look at the Kiali dashboard traffic tab for the client application we can see the traffic was TCP over mTLS.</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/images/kiali-client-remote-traffic.png" alt="kiali client remote traffic" />
</v-card-text>
</v-card>

<p>To exit from the <code>CohQL&gt;</code> prompt type the <code>bye</code> command. The delete the client Pod</p>

<markup
lang="bash"

>kubectl -n coherence delete -f extend-client-pod.yaml</markup>

</div>

<h4 id="_run_a_grpc_client_pod">Run a gRPC Client Pod</h4>
<div class="section">
<p>We can run the same image as a gRPC client.
For this example, instead of the NameService we will configure Coherence to</p>

<markup
lang="yaml"
title="grpc-client-pod.yaml"
>apiVersion: v1
kind: Pod
metadata:
  name: client
  labels:
    app: coherence-client
    version: 1.0.0
spec:
  containers:
    - name: coherence
      image: coherence-client:1.0.0
      command:
        - /files/runner
        - sleep
        - "15m"
      env:
        - name: COHERENCE_CLIENT
          value: "grpc-fixed"
        - name: COHERENCE_GRPC_ADDRESS
          value: "storage-headless.coherence.svc"
        - name: COHERENCE_GRPC_PORT
          value: "1408"</markup>

<p>We can now deploy the gRPC client Pod</p>

<markup
lang="bash"

>kubectl -n coherence delete -f grpc-client-pod.yaml</markup>

<p>And exec into the Pod to create a QueryPlus session.</p>

<markup
lang="bash"

>kubectl -n coherence exec -it client -- /files/runner queryplus</markup>

<p>We can run the same <code>create cache test</code> and <code>select * from test</code> command that we ran above to connect the client
to the cluster. This time the client should be connecting over gRPC.</p>

<p>If we now look at the Kiali dashboard we can see again that the client application has communicated with the storage cluster.
All of this communication was using mTLS but without configuring Coherence to use TLS.</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/images/kiali-client-grpc-app.png" alt="kiali client grpc app" />
</v-card-text>
</v-card>

<p>If we look at the Kiali dashboard traffic tab for the client application we can see that this time the
traffic was gRPC over mTLS.</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/images/kiali-client-grpc-traffic.png" alt="kiali client grpc traffic" />
</v-card-text>
</v-card>

<p>To exit from the <code>CohQL&gt;</code> prompt type the <code>bye</code> command. The delete the client Pod</p>

<markup
lang="bash"

>kubectl -n coherence delete -f extend-client-pod.yaml</markup>

</div>
</div>

<h3 id="_clients_outside_kubernetes">Clients Outside Kubernetes</h3>
<div class="section">
<p>When connecting Coherence Extend or gRPC clients from outside Kubernetes, the Coherence NameService cannot be used
by clients to look up the endpoints. The clients must be configured with fixed endpoints using the hostnames and ports
of the configured ingress or gateway services. Exactly how this is done will depend on the versions of Istio and Kubernetes
being used and whether Ingress or the Kubernetes Gateway API is used. The different options available make it impossible
to build an example that can cover all these scenarios.</p>

</div>
</div>
</doc-view>
