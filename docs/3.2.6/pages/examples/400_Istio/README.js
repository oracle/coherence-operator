<doc-view>

<h2 id="_using_coherence_with_istio">Using Coherence with Istio</h2>
<div class="section">
<p>You can run the Coherence cluster and manage then using the Coherence Operator alongside <a id="" title="" target="_blank" href="https://istio.io">Istio</a>.
Coherence clusters managed with the Coherence Operator 3.2.0 and later work with Istio 1.9.1 and later.
Coherence caches can be accessed from outside the Coherence cluster via Coherence*Extend, REST, and other supported Coherence clients.
Using Coherence clusters with Istio does not require the Coherence Operator to also be using Istio (and vice-versa) .
The Coherence Operator can manage Coherence clusters independent of whether those clusters are using Istio or not.</p>

<div class="admonition important">
<p class="admonition-textlabel">Important</p>
<p ><p>The current support for Istio has the following limitation:</p>

<p>Ports that are exposed in the ports list of the container spec in a Pod will be intercepted by the Envoy proxy in the Istio side-car container. Coherence cluster traffic must not pass through Envoy proxies as this will break Coherence, so the Coherence cluster port must never be exposed as a container port if using Istio. There is no real reason to expose the Coherence cluster port in a container because there is no requirement to have this port externally visible.</p>
</p>
</div>

<h3 id="_prerequisites">Prerequisites</h3>
<div class="section">
<p>The instructions assume that you are using a Kubernetes cluster with Istio installed and configured already.</p>

</div>

<h3 id="_using_the_coherence_operator_with_istio">Using the Coherence operator with Istio</h3>
<div class="section">
<p>To use Coherence operator with Istio, you can deploy the operator into a namespace which has Istio automatic sidecar injection enabled.  Before installing the operator, create the namespace in which you want to run the Coherence operator and label it for automatic injection.</p>

<markup
lang="bash"

>kubectl create namespace coherence
kubectl label namespace coherence istio-injection=enabled</markup>

<p>Istio Sidecar AutoInjection is done automatically when you label the coherence namespace with istio-injection.</p>

<p>After the namespace is labeled, you can install the operator using your preferred method in the Operator <router-link to="/docs/installation/01_installation">Installation Guide</router-link>.</p>

<p>After installed operator, use the following command to confirm the operator is running:</p>

<markup
lang="bash"

>kubectl get pods -n coherence

NAME                                                     READY   STATUS    RESTARTS   AGE
coherence-operator-controller-manager-7d76f9f475-q2vwv   2/2     Running   1          17h</markup>

<p>2/2 in READY column means that there are 2 containers running in the operator Pod. One is Coherence operator and the other is Envoy Proxy.</p>

</div>

<h3 id="_creating_a_coherence_cluster_with_istio">Creating a Coherence cluster with Istio</h3>
<div class="section">
<p>You can configure your cluster to run with Istio automatic sidecar injection enabled. Before creating your cluster, create the namespace in which you want to run the cluster and label it for automatic injection.</p>

<markup
lang="bash"

>kubectl create namespace coherence-example
kubectl label namespace coherence-example istio-injection=enabled</markup>

<p>There is no other requirements to run Coherence in Istio environment.</p>

<p>The following is an example that creates a cluster named example-cluster-storage:</p>

<p>example.yaml</p>

<markup
lang="bash"

># Example
apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: example-cluster-storage</markup>

<markup
lang="bash"

>$ kubectl -n coherence-example apply -f example.yaml</markup>

<p>After you installed the Coherence cluster, run the following command to view the pods:</p>

<markup
lang="bash"

>$ kubectl -n coherence-example get pods

NAME                                             READY   STATUS    RESTARTS   AGE
example-cluster-storage-0                        2/2     Running   0          45m
example-cluster-storage-1                        2/2     Running   0          45m
example-cluster-storage-2                        2/2     Running   0          45m</markup>

<p>You can see that 3 members in the cluster are running with 3 pods. 2/2 in READY column means that there are 2 containers running in each Pod. One is Coherence member and the other is Envoy Proxy.</p>

</div>

<h3 id="_tls">TLS</h3>
<div class="section">
<p>Coherence cluster works with mTLS. Coherence client can also support TLS through Istio Gateway with TLS termination to connect to Coherence cluster running inside kubernetes.  For example, you can apply the following Istio Gateway and Virtual Service in the namespace of the Coherence cluster.  Before applying the gateway, create a secret for the credential from the certificate and key (e.g. server.crt and server.key) to be used by the Gateway:</p>

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
