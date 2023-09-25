<doc-view>

<h2 id="_using_coherence_with_istio">Using Coherence with Istio</h2>
<div class="section">
<p>You can run the Coherence cluster and manage then using the Coherence Operator alongside <a id="" title="" target="_blank" href="https://istio.io">Istio</a>.
Coherence clusters managed with the Coherence Operator 3.2.0 and later work with Istio 1.9.1 and later.
Coherence caches can be accessed from outside the Coherence cluster via Coherence*Extend, REST, and other supported Coherence clients.
Using Coherence clusters with Istio does not require the Coherence Operator to also be using Istio (and vice-versa) .
The Coherence Operator can manage Coherence clusters independent of whether those clusters are using Istio or not.</p>


<h3 id="_why_doesnt_coherence_work_with_istio">Why Doesn&#8217;t Coherence Work with Istio?</h3>
<div class="section">
<p>Coherence uses a custom TCP message protocol for inter-cluster member communication.
When a cluster member sends a message to another member, the "reply to" address of the sending member is in the message. This address is the socket address the member is listening on (i.e. it is the IP address and port Coherence has bound to).
When Istio is intercepting traffic the message ends up being sent via the Envoy proxy and the actual port Coherence is listening on is blocked by Istio. When the member that receives the message tries to send a response to the reply to address, that port is not visible to it.</p>

<p>Coherence clients will work with Istio, so Extend, gRPC and http clients for things like REST, metrics and management will work when routed through the Envoy proxy.</p>

</div>

<h3 id="_prerequisites">Prerequisites</h3>
<div class="section">
<p>The instructions assume that you are using a Kubernetes cluster with Istio installed and configured already.</p>


<h4 id="_enable_istio_strict_mode">Enable Istio Strict Mode</h4>
<div class="section">
<p>For this example we make Istio run in "strict" mode so that it will not allow any traffic between Pods outside the Envoy proxy. If other modes are used, such as permissive, then Coherence will work as normal as its ports will not be blocked.</p>

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
<p>To use Coherence operator with Istio, you can deploy the operator into a namespace which has Istio automatic sidecar injection enabled.  Before installing the operator, create the namespace in which you want to run the Coherence operator and label it for automatic injection.</p>

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

<p>There are a number of ways to exclude the web-hook port, the simplest is to add a <code>PeerAuthentication</code> resource to the Operator namespace.</p>

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

<p>You can then install the operator using your preferred method in the Operator <router-link to="/docs/installation/01_installation">Installation Guide</router-link>.</p>

<p>After installed operator, use the following command to confirm the operator is running:</p>

<markup
lang="bash"

>kubectl get pods -n coherence

NAME                                                     READY   STATUS    RESTARTS   AGE
coherence-operator-controller-manager-7d76f9f475-q2vwv   2/2     Running   1          17h</markup>

<p>The output should show 2/2 in READY column, meaning there are 2 containers running in the Operator pod. One is Coherence Operator and the other is Envoy Proxy.</p>

</div>
</div>

<h3 id="_creating_a_coherence_cluster_with_istio">Creating a Coherence cluster with Istio</h3>
<div class="section">
<p>You can configure your cluster to run with Istio automatic sidecar injection enabled. Before creating your cluster, create the namespace in which you want to run the cluster and label it for automatic injection.</p>

<markup
lang="bash"

>kubectl create namespace coherence-example
kubectl label namespace coherence-example istio-injection=enabled</markup>


<h4 id="_exclude_the_coherence_cluster_ports">Exclude the Coherence Cluster Ports</h4>
<div class="section">
<p>As explained above, Coherence cluster traffic must be excluded from the Envoy proxy, there are various ways to do this.</p>

<p>There are three ports that must be excluded:</p>

<ul class="ulist">
<li>
<p>The cluster port - defaults to 7574, there is no need to set this to any other value.</p>

</li>
<li>
<p>The TCP first local port - the Operator will default this to 7575 using its web-hook (if the web-hook is disabled this needs to be manually set).</p>

</li>
<li>
<p>The TCP second local port - the Operator will default this to 7576 using its web-hook (if the web-hook is disabled this needs to be manually set).</p>

</li>
</ul>
<p><strong>1 Use an Annotation in the Coherence Resource</strong></p>

<p>The Istio exclusion annotation <code>traffic.sidecar.istio.io/excludeInboundPorts</code> can be added to the Coherence yaml to list the ports to be excluded,</p>

<p>For example, using the default ports the following annotation will exclude those ports from Istio:</p>

<markup
lang="yaml"
title="coherence-storage.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  annotations:
    traffic.sidecar.istio.io/excludeInboundPorts: "7574,7575,7576"</markup>

<p>If the Coherence Operator&#8217;s web-hook has been disabled, the local ports must be set in the yaml too:</p>

<markup
lang="yaml"
title="coherence-storage.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  annotations:
    traffic.sidecar.istio.io/excludeInboundPorts: "7574,7575,7576"
  coherence:
    localPort: 7575
    localPortAdjust: 7576</markup>

<p><strong>2 Use a PeerAuthentication resource</strong></p>

<p>A <code>PeerAuthentication</code> resource can be added to the Coherence cluster&#8217;s namespace <strong>before the cluster is deployed</strong>.</p>

<markup
lang="yaml"
title="istio-coherence.yaml"
>apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: "coherence"
spec:
  selector:
    matchLabels:
      coherenceComponent: coherencePod
  mtls:
    mode: STRICT
  portLevelMtls:
    7574:
      mode: PERMISSIVE
    7575:
      mode: PERMISSIVE
    7576:
      mode: PERMISSIVE</markup>

<p>The Coherence Operator labels Coherence Pods with the label <code>coherenceComponent: coherencePod</code> so this can be used in the <code>PeerAuthentication</code>. Then each port to be excluded is listed in the <code>portLevelMtls</code> and set to be <code>PERMISSIVE</code>.</p>

<p>This yaml can then be installed into the namespace that the Coherence cluster will be deployed into.</p>

</div>
</div>

<h3 id="_tls">TLS</h3>
<div class="section">
<p>Coherence clusters work with mTLS and Coherence clients can also support TLS through the Istio Gateway with TLS termination to connect to Coherence cluster running inside kubernetes. For example, you can apply the following Istio Gateway and Virtual Service in the namespace of the Coherence cluster.  Before applying the gateway, create a secret for the credential from the certificate and key (e.g. server.crt and server.key) to be used by the Gateway:</p>

<markup
lang="bash"

>kubectl create -n istio-system secret tls extend-credential --key=server.key --cert=server.crt</markup>

<p>Then, create a keystore (server.jks) to be used by the Coherence Extend client, e.g.:</p>

<markup
lang="bash"

>openssl pkcs12 -export -in server.crt -inkey server.key -chain -CAfile ca.crt -name "server" -out server.p12

keytool -importkeystore -deststorepass password -destkeystore server.jks -srckeystore server.p12 -srcstoretype PKCS12</markup>

<p>tlsGateway.yaml</p>

<markup
lang="bash"

>apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: tlsgateway
spec:
  selector:
    istio: ingressgateway # use istio default ingress gateway
  servers:
  - port:
      number: 8043
      name: tls
      protocol: TLS
    tls:
      mode: SIMPLE
      credentialName: "extend-credential" # the secret created in the previous step
      maxProtocolVersion: TLSV1_3
    hosts:
    - "*"</markup>

<p>tlsVS.yaml</p>

<markup
lang="bash"

>apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: extend
spec:
  hosts:
  - "*"
  gateways:
  - tlsgateway
  tcp:
  - match:
    route:
    - destination:
        host: example-cluster-proxy-proxy  # the service name used to expose the Extend proxy port</markup>

<p>Apply the Gateway and VirtualService:</p>

<markup
lang="bash"

>kubectl apply -f tlsGateway.yaml -n coherence-example
kubectl apply -f tlsVS.yaml -n coherence-example</markup>

<p>Then configure a Coherence*Extend client to connect to the proxy server via TLS protocol.  Below is an example of a &lt;remote-cache-scheme&gt; configuration of an Extend client using TLS port 8043 configured in the Gateway and server.jks created earlier in the example.</p>

<p>client-cache-config.xml</p>

<div class="listing">
<pre>...
    &lt;remote-cache-scheme&gt;
        &lt;scheme-name&gt;extend-direct&lt;/scheme-name&gt;
        &lt;service-name&gt;ExtendTcpProxyService&lt;/service-name&gt;
        &lt;initiator-config&gt;
            &lt;tcp-initiator&gt;
                &lt;socket-provider&gt;
                    &lt;ssl&gt;
                        &lt;protocol&gt;TLS&lt;/protocol&gt;
                        &lt;trust-manager&gt;
                            &lt;algorithm&gt;PeerX509&lt;/algorithm&gt;
                            &lt;key-store&gt;
                                &lt;url&gt;file:server.jks&lt;/url&gt;
                                &lt;password&gt;password&lt;/password&gt;
                            &lt;/key-store&gt;
                        &lt;/trust-manager&gt;
                    &lt;/ssl&gt;
                &lt;/socket-provider&gt;
                &lt;remote-addresses&gt;
                    &lt;socket-address&gt;
                        &lt;address&gt;$INGRESS_HOST&lt;/address&gt;
                        &lt;port&gt;8043&lt;/port&gt;
                    &lt;/socket-address&gt;
                &lt;/remote-addresses&gt;
            &lt;/tcp-initiator&gt;
        &lt;/initiator-config&gt;
    &lt;/remote-cache-scheme&gt;
...</pre>
</div>

<p>If you are using Docker for Desktop, <code>$INGRESS_HOST</code> is <code>127.0.0.1</code>, and you can use the Kubectl port-forward to allow the Extend client to access the Coherence cluster from your localhost:</p>

<markup
lang="bash"

>kubectl port-forward -n istio-system &lt;istio-ingressgateway-pod&gt; 8043:8043</markup>

</div>

<h3 id="_prometheus">Prometheus</h3>
<div class="section">
<p>The coherence metrics that record and track the health of Coherence cluster using Prometheus are also available in Istio environment and can be viewed through Grafana. However, Coherence cluster traffic is not visible by Istio.</p>

</div>

<h3 id="_traffic_visualization">Traffic Visualization</h3>
<div class="section">
<p>Istio provides traffic management capabilities, including the ability to visualize traffic in <a id="" title="" target="_blank" href="https://kiali.io">Kiali</a>. You do not need to change your applications to use this feature. The Istio proxy (envoy) sidecar that is injected into your pods provides it. The image below shows an example with traffic flow. In this example, you can see how the traffic flows in from the Istio gateway on the left, to the cluster services, and then to the individual cluster members.  This example has storage members (example-cluster-storage), a proxy member running proxy service (example-cluster-proxy), and a REST member running http server (example-cluster-rest).  However, Coherence cluster traffic between members is not visible.</p>



<v-card>
<v-card-text class="overflow-y-hidden" >
<img src="./images/images/istioKiali.png" alt="istioKiali"width="1024" />
</v-card-text>
</v-card>

<p>To learn more, see <a id="" title="" target="_blank" href="https://istio.io/latest/docs/concepts/traffic-management/">Istio traffic management</a>.</p>

</div>
</div>
</doc-view>
