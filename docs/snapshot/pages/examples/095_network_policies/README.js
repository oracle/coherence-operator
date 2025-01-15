<doc-view>

<h2 id="_using_network_policies">Using Network Policies</h2>
<div class="section">
<p>This example covers running the Coherence Operator and Coherence clusters in Kubernetes with network policies.
In Kubernetes, a <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/services-networking/network-policies/">Network Policy</a>
is an application-centric construct which allow you to specify how a pod is allowed to communicate with various network
"entities" (we use the word "entity" here to avoid overloading the more common terms such as "endpoints" and "services",
which have specific Kubernetes connotations) over the network.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Network policies in Kubernetes are easy to get wrong if you are not careful.
In this case a policy will either block traffic it should not, in which case your application will not work,
or it will let traffic through it should block, which will be an invisible security hole.</p>

<p>It is obviously important to test your policies, but Kubernetes offers next to zero visibility into what the policies
are actually doing, as it is typically the network CNI extensions that are providing the policy implementation
and each of these may work in a different way.</p>
</p>
</div>

<h3 id="_introduction">Introduction</h3>
<div class="section">
<p>Kubernetes network policies specify the access permissions for groups of pods, similar to security groups in the
cloud are used to control access to VM instances and similar to firewalls.
The default behaviour of a Kubernetes cluster is to allow all Pods to freely talk to each other.
Whilst this sounds insecure, originally Kubernetes was designed to orchestrate services that communicated with each other,
it was only later that network policies were added.</p>

<p>A network policy is applied to a Kubernetes namespace and controls ingress into and egress out of Pods in that namespace.
The ports specified in a <code>NetworkPolicy</code> are the ports exposed by the <code>Pods</code>, they are not any ports that may be exposed by
any <code>Service</code> that exposes the <code>Pod</code> ports. For example, if a <code>Pod</code> exposed port 8080 and a <code>Service</code> exposing the <code>Pod</code>
mapped port 80 in the <code>Service</code> to port <code>8080</code> in the <code>Pod</code>, the <code>NetworkPolicy</code> ingress rule would be for the <code>Pod</code> port 8080.</p>

<p>Network polices would typically end up being dictated by corporate security standards where different companies may
apply stricter or looser rules than others.
The examples in this document start from the premise that everything will be blocked by a "deny all" policy and then opened up as needed.
This is the most secure use of network policies, and hence the examples can easily be tweaked if looser rules are applied.</p>

<p>This example has the following sections:</p>

<ul class="ulist">
<li>
<p><router-link to="#deny" @click.native="this.scrollFix('#deny')">Deny All Policy</router-link> - denying all ingress and egress</p>

</li>
<li>
<p><router-link to="#dns" @click.native="this.scrollFix('#dns')">Allow DNS</router-link> - almost every use case will require egress to DNS</p>

</li>
<li>
<p><router-link to="#operator" @click.native="this.scrollFix('#operator')">Coherence Operator Policies</router-link> - the network policies required to run the Coherence Operator</p>
<ul class="ulist">
<li>
<p><router-link to="#k8sapi" @click.native="this.scrollFix('#k8sapi')">Kubernetes API Server</router-link> - allow the Operator egress to the Kubernetes API server</p>

</li>
<li>
<p><router-link to="#cluster-access" @click.native="this.scrollFix('#cluster-access')">Coherence Clusters Pods</router-link> - allow the Operator egress to the Coherence cluster Pods</p>

</li>
<li>
<p><router-link to="#webhook" @click.native="this.scrollFix('#webhook')">Web Hooks</router-link> - allow ingress to the Operator&#8217;s web hook port</p>

</li>
</ul>
</li>
<li>
<p><router-link to="#coherence" @click.native="this.scrollFix('#coherence')">Coherence Cluster Policies</router-link> - the network policies required to run Coherence clusters</p>
<ul class="ulist">
<li>
<p><router-link to="#inter-cluster" @click.native="this.scrollFix('#inter-cluster')">Inter-Cluster Access</router-link> - allow Coherence cluster Pods to communicate</p>

</li>
<li>
<p><router-link to="#cluster-to-operator" @click.native="this.scrollFix('#cluster-to-operator')">Coherence Operator</router-link> - allow Coherence cluster Pods to communicate with the Operator</p>

</li>
<li>
<p><router-link to="#client" @click.native="this.scrollFix('#client')">Clients</router-link> - allows access by Extend and gRPC clients</p>

</li>
<li>
<p><router-link to="#metrics" @click.native="this.scrollFix('#metrics')">Metrics</router-link> - allow Coherence cluster member metrics to be scraped</p>

</li>
</ul>
</li>
<li>
<p><router-link to="#testing" @click.native="this.scrollFix('#testing')">Testing Connectivity</router-link> - using the Operator&#8217;s network connectivity test utility to test policies</p>

</li>
</ul>

<h4 id="deny">Deny All Policy</h4>
<div class="section">
<p>Kubernetes does not have a “deny all” policy, but this can be achieved with a regular network policy that specifies
a <code>policyTypes</code> of both 'Ingress` and <code>Egress</code> but omits any definitions.
A wild-card <code>podSelector: {}</code> applies the policy to all Pods in the namespace.</p>

<markup
lang="yaml"
title="manifests/deny-all.yaml"
>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all
spec:
  podSelector: {}
  policyTypes:
    - Ingress
    - Egress
  ingress: []
  egress: []</markup>

<p>The policy above can be installed into the <code>coherence</code> namespace with the following command:</p>

<markup
lang="bash"

>kubectl -n coherence apply -f manifests/deny-all.yaml</markup>

<p>After installing the <code>deny-all</code> policy, any <code>Pod</code> in the <code>coherence</code> namespace will not be allowed either ingress,
nor egress. Very secure, but probably impractical for almost all use cases. After applying the <code>deny-all</code> policy more polices can be added to gradually open up the required access to run the Coherence Operator and Coherence clusters.</p>

</div>

<h4 id="dns">Allow DNS</h4>
<div class="section">
<p>When enforcing egress, such as with the <code>deny-all</code> policy above, it is important to remember that virtually every Pod needs
to communicate with other Pods or Services, and will therefore need to access DNS.</p>

<p>The policy below allows all Pods (using <code>podSelector: {}</code>) egress to both TCP and UDP on port 53 in all namespaces.</p>

<markup
lang="yaml"
title="manifests/allow-dns.yaml"
>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-dns
spec:
  podSelector: { }
  policyTypes:
    - Egress
  egress:
    - to:
        - namespaceSelector: { }
      ports:
        - protocol: UDP
          port: 53
#        - protocol: TCP
#          port: 53</markup>

<p>If allowing DNS egress to all namespaces is overly permissive, DNS could be further restricted to just the <code>kube-system</code>
namespace, therefore restricting DNS lookups to only Kubernetes internal DNS.
Kubernetes applies the <code>kubernetes.io/metadata.name</code> label to namespaces, and sets its value to the namespace name,
so this can be used in label matchers.</p>

<p>With the policy below, Pods will be able to use internal Kubernetes DNS only.</p>

<markup
lang="yaml"
title="manifests/allow-dns-kube-system.yaml"
>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-dns
spec:
  podSelector: { }
  policyTypes:
    - Egress
  egress:
    - to:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: kube-system
      ports:
        - protocol: UDP
          port: 53
#        - protocol: TCP
#          port: 53</markup>

<p>The policy above can be installed into the <code>coherence</code> namespace with the following command:</p>

<markup
lang="bash"

>kubectl -n coherence apply -f manifests/allow-dns-kube-system.yaml</markup>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p>Some documentation regarding allowing DNS with Kubernetes network policies only shows opening up UDP connections.
During our testing with network policies, we discovered that with only UDP allowed any lookup for a fully qualified
name would fail. For example <code>nslookup my-service.my-namespace.svc</code> would work, but the fully qualified
<code>nslookup my-service.my-namespace.svc.cluster.local</code> would not. Adding TCP to the DNS policy allowed DNS lookups with
<code>.cluster.local</code> to also work.</p>

<p>Neither the Coherence Operator, nor Coherence itself use a fully qualified service name for a DNS lookup.
It appears that Java&#8217;s <code>InetAddress.findAllByName()</code> method still works only with UDP, albeit extremely slowly.
By default, the service name used for the Coherence WKA setting uses just the <code>.svc</code> suffix.</p>
</p>
</div>
</div>
</div>

<h3 id="operator">Coherence Operator Policies</h3>
<div class="section">
<p>Assuming the <code>coherence</code> namespace exists, and the <code>deny-all</code> and <code>allow-dns</code> policies described above have been applied,
if the Coherence Operator is installed, it wil fail to start as it has no access to endpoints it needs to operate.
The following sections will add network polices to allow the Coherence Operator to access Kubernetes services and
Pods it requires.</p>


<h4 id="k8sapi">Access the Kubernetes API Server</h4>
<div class="section">
<p>The Coherence Operator uses Kubernetes APIs to manage various resources in the Kubernetes cluster.
For this to work, the Operator Pod must be allowed egress to the Kubernetes API server.
Configuring access to the API server is not as straight forward as other network policies.
The reason for this is that there is no Pod available with labels that can be used in the configuration,
instead, the IP address of the API server itself must be used.</p>

<p>There are various methods to find the IP address of the API server.
The exact method required may vary depending on the type of Kubernetes cluster being used, for example a simple
development cluster running in KinD on a laptop may differ from a cluster running in a cloud provider&#8217;s infrastructure.</p>

<p>The common way to find the API server&#8217;s IP address is to use <code>kubectl cluster-info</code> as follows:</p>

<markup
lang="bash"

>$ kubectl cluster-info
Kubernetes master is running at https://192.168.99.100:8443</markup>

<p>In the above case the IP address of the API server would be <code>192.168.99.100</code> and the port is <code>8443</code>.</p>

<p>In a simple KinD development cluster, the API server IP address can be obtained using <code>kubectl</code> as shown below:</p>

<markup
lang="bash"

>$ kubectl -n default get endpoints kubernetes -o json
{
    "apiVersion": "v1",
    "kind": "Endpoints",
    "metadata": {
        "creationTimestamp": "2023-02-08T10:31:26Z",
        "labels": {
            "endpointslice.kubernetes.io/skip-mirror": "true"
        },
        "name": "kubernetes",
        "namespace": "default",
        "resourceVersion": "196",
        "uid": "68b0a7de-c0db-4524-a1a2-9d29eb137f28"
    },
    "subsets": [
        {
            "addresses": [
                {
                    "ip": "192.168.49.2"
                }
            ],
            "ports": [
                {
                    "name": "https",
                    "port": 8443,
                    "protocol": "TCP"
                }
            ]
        }
    ]
}</markup>

<p>In the above case the IP address of the API server would be <code>192.168.49.2</code> and the port is <code>8443</code>.</p>

<p>The IP address displayed for the API server can then be used in the network policy.
The policy shown below allows Pods with the <code>app.kubernetes.io/name: coherence-operator</code> label (which the Operator has)
egress access to the API server.</p>

<markup
lang="yaml"
title="manifests/allow-k8s-api-server.yaml"
>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: operator-to-apiserver-egress
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: coherence-operator
  policyTypes:
    - Egress
    - Ingress
  egress:
    - to:
        - ipBlock:
            cidr: 172.18.0.2/24
        - ipBlock:
            cidr: 10.96.0.1/24
      ports:
        - port: 6443
          protocol: TCP
        - port: 443
          protocol: TCP</markup>

<p>The <code>allow-k8s-api-server.yaml</code> policy can be installed into the <code>coherence</code> namespace to allow the Operator to communicate with the API server.</p>

<markup
lang="bash"

>kubectl -n coherence apply -f manifests/allow-k8s-api-server.yaml</markup>

<p>With the <code>allow-k8s-api-server.yaml</code> policy applied, the Coherence Operator should now start correctly and its Pods should reach the "ready" state.</p>

</div>

<h4 id="cluster-access">Ingress From and Egress Into Coherence Cluster Member Pods</h4>
<div class="section">
<p>When a Coherence cluster is deployed, on start-up of a Pod the cluster member will connect to the Operator&#8217;s REST endpoint to query the site name and rack name, based on the Node the Coherence member is running on. To allow this to happen the Operator needs to be configured with the relevant ingress policy.</p>

<p>The <code>coherence-operator-rest-ingress</code> policy applies to the Operator Pod, as it has a <code>podSelector</code> label of <code>app.kubernetes.io/name: coherence-operator</code>, which is a label applied to the Operator Pod. The policy allows any Pod with the label <code>coherenceComponent: coherencePod</code> ingress
into the <code>operator</code> REST port. When the Operator creates a Coherence cluster, it applies the label <code>coherenceComponent: coherencePod</code> to all the Coherence cluster Pods.
The policy below allows access from all namespaces using <code>namespaceSelector: { }</code> but it could be tightened up to specific namespaces if required.</p>

<markup
lang="yaml"
title="manifests/allow-operator-rest-ingress.yaml"
>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: coherence-operator-rest-ingress
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: coherence-operator
  policyTypes:
    - Ingress
  ingress:
    - from:
        - namespaceSelector: { }
          podSelector:
            matchLabels:
              coherenceComponent: coherencePod
      ports:
        - port: operator
          protocol: TCP</markup>

<p>During operations such as scaling and shutting down of a Coherence cluster, the Operator needs to connect to the health endpoint of the Coherence cluster Pods.</p>

<p>The <code>coherence-operator-cluster-member-egress</code> policy below applies to the Operator Pod, as it has a <code>podSelector</code> label of <code>app.kubernetes.io/name: coherence-operator</code>, which is a label applied to the Operator Pod. The policy allows egress to the <code>health</code> port in any Pod with the label <code>coherenceComponent: coherencePod</code>. When the Operator creates a Coherence cluster, it applies the label <code>coherenceComponent: coherencePod</code> to all the Coherence cluster Pods.
The policy below allows egress to Coherence Pods in all namespaces using <code>namespaceSelector: { }</code> but it could be tightened up to specific namespaces if required.</p>

<markup
lang="yaml"
title="manifests/allow-operator-cluster-member-egress.yaml"
>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: coherence-operator-cluster-member-egress
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: coherence-operator
  policyTypes:
    - Egress
  egress:
    - to:
        - namespaceSelector: { }
          podSelector:
            matchLabels:
              coherenceComponent: coherencePod
      ports:
        - port: health
          protocol: TCP</markup>

<p>The two policies can be applied to the <code>coherence</code> namespace.</p>

<markup
lang="bash"

>kubectl -n coherence apply -f manifests/allow-operator-rest-ingress.yaml
kubectl -n coherence apply -f manifests/allow-operator-cluster-member-egress.yaml</markup>

</div>

<h4 id="webhook">Webhook Ingress</h4>
<div class="section">
<p>With all the above policies in place, the Operator is able to work correctly, but if a <code>Coherence</code> resource is now created
Kubernetes will be unable to call the Operator&#8217;s webhook without the correct ingress policy.</p>

<p>The following example demonstrates this. Assume there is a minimal`Coherence` yaml file named <code>minimal.yaml</code>
that will create a single member Coherence cluster.</p>

<markup
lang="yaml"
title="minimal.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  replicas: 1</markup>

<p>If <code>minimal.yaml</code> is applied using <code>kubectl</code> with a small timeout of 10 seconds, the creation of the resource will
fail due to Kubernetes not having access to the Coherence Operator webhook.</p>

<markup
lang="bash"

>$ kubectl apply --timeout=10s -f minimal.yaml
Error from server (InternalError): error when creating "minimal.yaml": Internal error occurred: failed calling webhook "coherence.oracle.com": failed to call webhook: Post "https://coherence-operator-webhook.operator-test.svc:443/mutate-coherence-oracle-com-v1-coherence?timeout=10s": context deadline exceeded</markup>

<p>The simplest solution is to allow ingress from any IP address to the webhook on port, with a policy like that shown below.
This policy uses and empty <code>from: []</code> attribute, which allows access from anywhere to the <code>webhook-server</code> port in the Pod.</p>

<markup
lang="yaml"
title="manifests/allow-webhook-ingress-from-all.yaml"
>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: apiserver-to-operator-webhook-ingress
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: coherence-operator
  policyTypes:
    - Ingress
  ingress:
    - from: []
      ports:
        - port: webhook-server
          protocol: TCP</markup>

<p>Allowing access to the webhook from anywhere is not very secure, so a more restrictive <code>from</code> attribute could be used
to limit access to the IP address (or addresses) of the Kubernetes API server.
As with the API server policy above, the trick here is knowing the API server addresses to use.</p>

<p>The policy below only allows access from specific addresses:</p>

<markup
lang="yaml"
title="manifests/allow-webhook-ingress-from-all.yaml"
>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: apiserver-to-operator-webhook-ingress
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: coherence-operator
  policyTypes:
    - Ingress
  ingress:
    - from:
        - ipBlock:
            cidr: 172.18.0.2/24
        - ipBlock:
            cidr: 10.96.0.1/24
      ports:
        - port: webhook-server
          protocol: TCP
        - port: 443
          protocol: TCP</markup>

</div>
</div>

<h3 id="coherence">Coherence Cluster Member Policies</h3>
<div class="section">
<p>Once the policies are in place to allow the Coherence Operator to work, the policies to allow Coherence clusters to run can be put in place.
The exact set of policies requires will vary depending on the Coherence functionality being used.
If Coherence is embedded in another application, such as a web-server, then additional policies may also be needed to allow ingress to other endpoints.
Conversely, if the Coherence application needs access to other services, for example a database, then additional egress policies may need to be created.</p>

<p>This example is only going to cover Coherence use cases, but it should be simple enough to apply the same techniques to policies for other applications.</p>


<h4 id="inter-cluster">Access Other Cluster Members</h4>
<div class="section">
<p>All Pods in a Coherence cluster must be able to talk to each other (otherwise they wouldn&#8217;t be a cluster).
This means that there needs to be ingress and egress policies to allow this.</p>

<p><strong>Cluster port</strong>: The default cluster port is 7574, and there is almost never any need to change this, especially in a containerised environment where there is little chance of port conflicts.</p>

<p><strong>Unicast ports</strong>: Unicast uses TMB (default) and UDP. Each cluster member listens on one UDP and one TCP port and both ports need to be opened in the network policy. The default behaviour of Coherence is for the unicast ports to be automatically assigned from the operating system&#8217;s available ephemeral port range. When securing Coherence with network policies, the use of ephemeral ports will not work, so a range of ports can be specified for coherence to operate within. The Coherence Operator sets values for both unicast ports so that ephemeral ports will not be used. The default values are <code>7575</code> and <code>7576</code>.</p>

<p>The two unicast ports can be changed in the <code>Coherence</code> spec by setting the <code>spec.coherence.localPort</code> field,
and the <code>spec.coherence.localPortAdjust</code> field for example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  coherence:
    localPort: 9000
    localPortAdjust: 9001</markup>

<p>Alternatively the values can also be configured using environment variables</p>

<markup
lang="yaml"

>env:
  - name: COHERENCE_LOCALPORT
    value: "9000"
  - name: COHERENCE_LOCALPORT_ADJUST
    value: "9001"</markup>

<p><strong>Echo port <code>7</code></strong>: The default TCP port of the IpMonitor component that is used for detecting hardware failure of cluster members. Coherence doesn&#8217;t bind to this port, it only tries to connect to it as a means of pinging remote machines, or in this case Pods.</p>

<p>The Coherence Operator applies the <code>coherenceComponent: coherencePod</code> label to all Coherence Pods, so this can be used in the network policy <code>podSelector</code>, to apply the policy to only the Coherence Pods.</p>

<p>The policy below works with the default ports configured by the Operator.</p>

<markup
lang="yaml"
title="manifests/allow-cluster-member-access.yaml"
>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-coherence-cluster
spec:
  podSelector:
    matchLabels:
      coherenceComponent: coherencePod
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              coherenceComponent: coherencePod
      ports:
        - port: 7574
          endPort: 7576
          protocol: TCP
        - port: 7574
          endPort: 7576
          protocol: UDP
        - port: 7
          protocol: TCP
  egress:
    - to:
        - podSelector:
            matchLabels:
              coherenceComponent: coherencePod
      ports:
        - port: 7574
          endPort: 7576
          protocol: TCP
        - port: 7574
          endPort: 7576
          protocol: UDP
        - port: 7
          protocol: TCP</markup>

<p>If the Coherence local port and local port adjust values are changed, then the policy would need to be amended.
For example, if <code>COHERENCE_LOCALPORT=9000</code> and <code>COHERENCE_LOCALPORT_ADJUST=9100</code></p>

<markup
lang="yaml"
title="manifests/allow-cluster-member-access-non-default.yaml"
>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-coherence-cluster
spec:
  podSelector:
    matchLabels:
      coherenceComponent: coherencePod
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              coherenceComponent: coherencePod
      ports:
        - port: 7574
          protocol: TCP
        - port: 7574
          protocol: UDP
        - port: 9000
          endPort: 9100
          protocol: TCP
        - port: 9000
          endPort: 9100
          protocol: UDP
        - port: 7
          protocol: TCP
  egress:
    - to:
        - podSelector:
            matchLabels:
              coherenceComponent: coherencePod
      ports:
        - port: 7574
          protocol: TCP
        - port: 7574
          protocol: UDP
        - port: 9000
          endPort: 9100
          protocol: TCP
        - port: 9000
          endPort: 9100
          protocol: UDP
        - port: 7
          protocol: TCP</markup>

<p>Both of the policies above should be applied to the namespace where the Coherence cluster will be deployed.
With the two policies above in place, the Coherence Pods will be able to communicate.</p>

</div>

<h4 id="cluster-to-operator">Egress to and Ingress From the Coherence Operator</h4>
<div class="section">
<p>When a Coherence Pod starts Coherence calls back to the Operator to obtain the site name and rack name based on the Node the Pod is scheduled onto.
For this to work, there needs to be an egress policy to allow Coherence Pods to access the Operator.</p>

<p>During certain operations the Operator needs to call the Coherence members health endpoint to check health and status.
For this to work there needs to be an ingress policy to allow the Operator access to the health endpoint in the Coherence Pods</p>

<p>The policy below applies to Pods with the <code>coherenceComponent: coherencePod</code> label, which will match Coherence cluster member Pods.
The policy allows ingress from the Operator to the Coherence Pod health port from namespace <code>coherence</code>
using the namespace selector label <code>kubernetes.io/metadata.name: coherence</code>
and Pod selector label <code>app.kubernetes.io/name: coherence-operator</code>
The policy allows egress from the Coherence pods to the Operator&#8217;s REST server <code>operator</code> port.</p>

<markup
lang="yaml"
title="manifests/allow-cluster-member-operator-access.yaml"
>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: coherence-operator-cluster-member-access
spec:
  podSelector:
    matchLabels:
      coherenceComponent: coherencePod
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: coherence
          podSelector:
            matchLabels:
              app.kubernetes.io/name: coherence-operator
      ports:
        - port: health
          protocol: TCP
  egress:
    - to:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: coherence
          podSelector:
            matchLabels:
              app.kubernetes.io/name: coherence-operator
      ports:
        - port: operator
          protocol: TCP</markup>

<p>If the Operator is not running in the <code>coherence</code> namespace then the namespace match label can be changed to the required value.
The policy above should be applied to the namespace where the Coherence cluster will be deployed.</p>

</div>

<h4 id="client">Client Access (Coherence*Extend and gRPC)</h4>
<div class="section">
<p>A typical Coherence cluster does not run in isolation but as part of a larger application.
If the application has other Pods that are Coherence clients, then they will need access to the Coherence cluster.
This would usually mean creating ingress and egress policies for the Coherence Extend port and gRPC port,
depending on which Coherence APIs are being used.</p>

<p>Instead of using actual port numbers, a <code>NetworkPolicy</code> can be made more flexible by using port names.
When ports are defined in a container spec of a Pod, they are usually named.
By using the names of the ports in the <code>NetworkPolicy</code> instead of port numbers, the real port numbers can be changed without
affecting the network policy.</p>


<h5 id="_coherence_extend_access">Coherence Extend Access</h5>
<div class="section">
<p>If Coherence Extend is being used, then first the Extend Proxy must be configured to use a fixed port.
The default behaviour of Coherence is to bind the Extend proxy to an ephemeral port and clients use the Coherence
NameService to look up the port to use.</p>

<p>When using the default Coherence images, for example <code>ghcr.io/oracle/coherence-ce:22.06</code> the Extend proxy is already
configured to run on a fixed port <code>20000</code>. When using this image, or any image that uses the default Coherence cache
configuration file, this port can be changed by setting the <code>COHERENCE_EXTEND_PORT</code> environment variable.</p>

<p>When using the Coherence Concurrent extensions over Extend, the Concurrent Extend proxy also needs to be configured with a fixed port.
When using the default Coherence images, for example <code>ghcr.io/oracle/coherence-ce:22.06</code> the Concurrent Extend proxy is already
configured to run on a fixed port <code>20001</code>. When using this image, or any image that uses the default Coherence cache
configuration file, this port can be changed by setting the <code>COHERENCE_CONCURRENT_EXTEND_PORT</code> environment variable.</p>

<p>For the examples below, a <code>Coherence</code> deployment has the following configuration.
This will expose Extend on a port named <code>extend</code> with a port number of <code>20000</code>, and a port named <code>extend-atomics</code>
with a port number of <code>20001</code>. The polices described below will then use the port names,
so if required the port number could be changed and the policies would still work.</p>

<markup
lang="yaml"
title="coherence-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  ports:
    - name: extend
      port: 20000
    - name: extend-atomics
      port: 20001</markup>

<p>The ingress policy below will work with the default Coherence image and allow ingress into the Coherence Pods
to both the default Extend port and Coherence Concurrent Extend port.
The policy allows ingress from Pods that have the <code>coherence.oracle.com/extendClient: true</code> label, from any namespace.
It could be tightened further by using a more specific namespace selector.</p>

<markup
lang="yaml"
title="manifests/allow-extend-ingress.yaml"
>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: coherence-extend-ingress
spec:
  podSelector:
    matchLabels:
      coherenceComponent: coherencePod
  policyTypes:
    - Ingress
  ingress:
    - from:
        - namespaceSelector: {}
          podSelector:
            matchLabels:
              coherence.oracle.com/extendClient: "true"
      ports:
        - port: extend
          protocol: TCP
        - port: extend-atomics
          protocol: TCP</markup>

<p>The policy above should be applied to the namespace where the Coherence cluster is running.</p>

<p>Instead of using fixed port numbers in the</p>

<p>The egress policy below will work with the default Coherence image and allow egress from Pods with the
<code>coherence.oracle.com/extendClient: true</code> label to Coherence Pods with the label <code>coherenceComponent: coherencePod</code>.
on both the default Extend port and Coherence Concurrent Extend port.</p>

<markup
lang="yaml"
title="manifests/allow-extend-egress.yaml"
>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: coherence-extend-egress
spec:
  podSelector:
    matchLabels:
      coherence.oracle.com/extendClient: "true"
  policyTypes:
    - Ingress
  egress:
    - to:
      - namespaceSelector: { }
        podSelector:
          matchLabels:
            coherenceComponent: coherencePod
      ports:
        - port: extend
          protocol: TCP
        - port: extend-atomics
          protocol: TCP</markup>

<p>The policy above allows egress to Coherence Pods in any namespace. This would ideally be tightened up to the specific
namespace that the Coherence cluster is deployed in.
For example, if the Coherence cluster is deployed in the <code>datastore</code> namespace, then the <code>to</code> section of policy could
be changed as follows:</p>

<markup
lang="yaml"
title="manifests/allow-extend-egress.yaml"
>- to:
  - namespaceSelector:
      matchLabels:
        kubernetes.io/metadata.name: datastore
    podSelector:
      matchLabels:
        coherenceComponent: coherencePod</markup>

<p>This policy must be applied to the namespace <em>where the client Pods will be deployed</em>.</p>

</div>

<h5 id="_coherence_grpc_access">Coherence gRPC Access</h5>
<div class="section">
<p>If Coherence gRPC is being used, then first the gRPC Proxy must be configured to use a fixed port.</p>

<p>When using the default Coherence images, for example <code>ghcr.io/oracle/coherence-ce:22.06</code> the gRPC proxy is already
configured to run on a fixed port <code>1408</code>. The gRPC proxy port can be changed by setting the <code>COHERENCE_GRPC_PORT</code> environment variable.</p>

<p>The ingress policy below will allow ingress into the Coherence Pods gRPC port.
The policy allows ingress from Pods that have the <code>coherence.oracle.com/grpcClient: true</code> label, from any namespace.
It could be tightened further by using a more specific namespace selector.</p>

<markup
lang="yaml"
title="manifests/allow-grpc-ingress.yaml"
>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: coherence-grpc-ingress
spec:
  podSelector:
    matchLabels:
      coherenceComponent: coherencePod
  policyTypes:
    - Ingress
  ingress:
    - from:
        - namespaceSelector: {}
          podSelector:
            matchLabels:
              coherence.oracle.com/grpcClient: "true"
      ports:
        - port: grpc
          protocol: TCP</markup>

<p>The policy above should be applied to the namespace where the Coherence cluster is running.</p>

<p>The egress policy below will allow egress to the gRPC port from Pods with the <code>coherence.oracle.com/grpcClient: true</code> label to
Coherence Pods with the label <code>coherenceComponent: coherencePod</code>.</p>

<markup
lang="yaml"
title="manifests/allow-extend-egress.yaml"
>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: coherence-extend-egress
spec:
  podSelector:
    matchLabels:
      coherence.oracle.com/extendClient: "true"
  policyTypes:
    - Ingress
  egress:
    - to:
      - namespaceSelector: { }
        podSelector:
          matchLabels:
            coherenceComponent: coherencePod
      ports:
        - port: extend
          protocol: TCP
        - port: extend-atomics
          protocol: TCP</markup>

<p>The policy above allows egress to Coherence Pods in any namespace. This would ideally be tightened up to the specific
namespace that the Coherence cluster is deployed in.
For example, if the Coherence cluster is deployed in the <code>datastore</code> namespace, then the <code>to</code> section of policy could
be changed as follows:</p>

<markup
lang="yaml"
title="manifests/allow-extend-egress.yaml"
>- to:
  - namespaceSelector:
      matchLabels:
        kubernetes.io/metadata.name: datastore
    podSelector:
      matchLabels:
        coherenceComponent: coherencePod</markup>

<p>This policy must be applied to the namespace <em>where the client Pods will be deployed</em>.</p>

</div>
</div>

<h4 id="metrics">Coherence Metrics</h4>
<div class="section">
<p>If Coherence metrics is enabled there will need to be an ingress policy to allow connections from metrics clients.
There would also need to be a similar egress policy in the metrics client&#8217;s namespace to allow it to access the Coherence metrics endpoints.</p>

<p>A simple <code>Coherence</code> resource that will create a cluster with metrics enabled is shown below.
This yaml will create a Coherence cluster with a port names <code>metrics</code> that maps to the default metrics port of '9612`.</p>

<markup
lang="yaml"
title="manifests/coherence-cluster-with-metrics.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  coherence:
    metrics:
      enabled: true
  ports:
    - name: metrics
      serviceMonitor:
        enabled: true</markup>

<p>The example below will assume that metrics will be scraped by Prometheus, and that Prometheus is installed into a namespace called <code>monitoring</code>.
An ingress policy must be created in the namespace where the Coherence cluster is deployed allowing ingress to the metrics port from the Prometheus Pods.
The Pods running Prometheus have a label <code>app.kubernetes.io/name: prometheus</code> so this can be used in the policy&#8217;s Pod selector.
This policy should be applied to the namespace where the Coherence cluster is running.</p>

<markup
lang="yaml"
title="manifests/allow-metrics-ingress.yaml"
>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: coherence-metrics-ingress
spec:
  podSelector:
    matchLabels:
      coherenceComponent: coherencePod
  policyTypes:
    - Ingress
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              kubernetes.io/metadata.name: monitoring
          podSelector:
            matchLabels:
              app.kubernetes.io/name: prometheus
      ports:
        - port: metrics
          protocol: TCP</markup>

<p>If the <code>monitoring</code> namespace also has a "deny-all" policy and needs egress opening up for Prometheus to scrape metrics then an egress policy will need to be added to the <code>monitoring</code> namespace.</p>

<p>The policy below will allow Pods with the label <code>app.kubernetes.io/name: prometheus</code> egress to Pods with the <code>coherenceComponent: coherencePod</code> label in any namespace. The policy could be further tightened up by adding a namespace selector to restrict egress to the specific namespace where the Coherence cluster is running.</p>

<markup
lang="yaml"
title="manifests/allow-metrics-egress.yaml"
>apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: coherence-metrics-egress
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: prometheus
  policyTypes:
    - Egress
  egress:
    - to:
        - namespaceSelector: { }
          podSelector:
            matchLabels:
              coherenceComponent: coherencePod
      ports:
        - port: metrics
          protocol: TCP</markup>

</div>
</div>

<h3 id="testing">Testing Network Policies</h3>
<div class="section">
<p>At the time of writing this documentation, Kubernetes provides no way to verify the correctness of network policies.
It is easy to mess up a policy, in which case policies will either block too much traffic, in which case your application
will work, or worse they will not be blocking access and leave a security hole.</p>

<p>As we have had various requests for help from customers who cannot get Coherence to work with network policies enabled,
the Operator has a simple utility to test connectivity outside of Coherence. This will allow testing pf policies without
the complications of having to start a Coherence server.</p>

<p>This example includes some simple yaml files that will create simulator Pods that listen on all the ports used by the Operator
and by a Coherence cluster member. These simulator Pods are configured with the same labels that the real Operator and
Coherence Pods would have and the same labels used by the network policies in this example. Also included are some yaml files
that start a test client, that simulates either the Operator connecting to Coherence Pods or a Coherence Pod connecting to
the Operator and to other Coherence Pods.</p>

<p>To run these tests, the Operator does not have to be installed.</p>


<h4 id="_create_the_test_namespaces">Create the Test Namespaces</h4>
<div class="section">
<p>In this example we will assume the Operator will eventually be running in a namespace called <code>coherence</code> and the Coherence
cluster will run in a namespace called <code>coh-test</code>. We can create the namespaces using <code>kubectl</code></p>

<markup
lang="bash"

>kubectl create ns coherence</markup>

<markup
lang="bash"

>kubectl create ns coh-test</markup>

<p>At this point there are no network policies installed, this will allow us to confirm the connectivity tests work.</p>

</div>

<h4 id="_start_the_operator_simulator">Start the Operator Simulator</h4>
<div class="section">
<p>The Operator simulator server should run in the <code>coherence</code> namespace.
It can be created using the following command:</p>

<markup
lang="bash"

>kubectl -n coherence apply -f examples/095_network_policies/manifests/net-test-operator-server.yaml</markup>

</div>

<h4 id="_start_the_coherence_cluster_simulator">Start the Coherence Cluster Simulator</h4>
<div class="section">
<p>The Coherence cluster member simulator server should run in the <code>coh-test</code> namespace.
It can be created using the following command:</p>

<markup
lang="bash"

>kubectl -n coh-test apply -f examples/095_network_policies/manifests/net-test-coherence-server.yaml</markup>

</div>

<h4 id="_run_the_operator_test">Run the Operator Test</h4>
<div class="section">
<p>We can now run the Operator test Job. This wil run a Kubernetes Job that simulates the Operator connecting to the
Kubernetes API server and to the Operator Pods.</p>

<markup
lang="bash"

>kubectl -n coherence apply -f examples/095_network_policies/manifests/net-test-operator.yaml</markup>

<p>The test Job should complete very quickly as it is only testing connectivity to various ports.
The results of the test can be seen by looking at the Pod log. The command below will display the log:</p>

<markup
lang="bash"

>kubectl -n coherence logs $(kubectl -n coherence get pod -l 'coherenceNetTest=operator-client' -o name)</markup>

<p>The output from a successful test will look like this:</p>

<markup


>1.6727606592497227e+09	INFO	runner	Operator Version: 3.3.2
1.6727606592497835e+09	INFO	runner	Operator Build Date: 2023-01-03T12:25:58Z
1.6727606592500978e+09	INFO	runner	Operator Built By: jonathanknight
1.6727606592501197e+09	INFO	runner	Operator Git Commit: c8118585b8f3d72b083ab1209211bcea364c85c5
1.6727606592501485e+09	INFO	runner	Go Version: go1.19.2
1.6727606592501757e+09	INFO	runner	Go OS/Arch: linux/amd64
1.6727606592504115e+09	INFO	net-test	Starting test	{"Name": "Operator Simulator"}
1.6727606592504556e+09	INFO	net-test	Testing connectivity	{"PortName": "K8s API Server"}
1.6727606592664087e+09	INFO	net-test	Testing connectivity PASSED	{"PortName": "K8s API Server", "Version": "v1.24.7"}
1.6727606592674055e+09	INFO	net-test	Testing connectivity	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "Health", "Port": 6676}
1.6727606592770455e+09	INFO	net-test	Testing connectivity PASSED	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "Health", "Port": 6676}</markup>

<p>We can see that the test has connected to the Kubernetes API server and has connected to the health port on the
Coherence cluster test server in the <code>coh-test</code> namespace.</p>

<p>The test Job can then be deleted:</p>

<markup
lang="bash"

>kubectl -n coherence delete -f examples/095_network_policies/manifests/net-test-operator.yaml</markup>

</div>

<h4 id="_run_the_cluster_member_test">Run the Cluster Member Test</h4>
<div class="section">
<p>The cluster member test simulates a Coherence cluster member connecting to other cluster members in the same namespace
and also making calls to the Operator&#8217;s REST endpoint.</p>

<markup
lang="bash"

>kubectl -n coh-test apply -f examples/095_network_policies/manifests/net-test-coherence.yaml</markup>

<p>Again, the test should complete quickly as it is just connecting to various ports.
The results of the test can be seen by looking at the Pod log. The command below will display the log:</p>

<markup
lang="bash"

>kubectl -n coh-test logs $(kubectl -n coh-test get pod -l 'coherenceNetTest=coherence-client' -o name)</markup>

<p>The output from a successful test will look like this:</p>

<markup


>1.6727631152848177e+09	INFO	runner	Operator Version: 3.3.2
1.6727631152849226e+09	INFO	runner	Operator Build Date: 2023-01-03T12:25:58Z
1.6727631152849536e+09	INFO	runner	Operator Built By: jonathanknight
1.6727631152849755e+09	INFO	runner	Operator Git Commit: c8118585b8f3d72b083ab1209211bcea364c85c5
1.6727631152849965e+09	INFO	runner	Go Version: go1.19.2
1.6727631152850187e+09	INFO	runner	Go OS/Arch: linux/amd64
1.6727631152852216e+09	INFO	net-test	Starting test	{"Name": "Cluster Member Simulator"}
1.6727631152852666e+09	INFO	net-test	Testing connectivity	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "UnicastPort1", "Port": 7575}
1.6727631152997334e+09	INFO	net-test	Testing connectivity PASSED	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "UnicastPort1", "Port": 7575}
1.6727631152998908e+09	INFO	net-test	Testing connectivity	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "UnicastPort2", "Port": 7576}
1.6727631153059115e+09	INFO	net-test	Testing connectivity PASSED	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "UnicastPort2", "Port": 7576}
1.6727631153063197e+09	INFO	net-test	Testing connectivity	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "Management", "Port": 30000}
1.6727631153116117e+09	INFO	net-test	Testing connectivity PASSED	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "Management", "Port": 30000}
1.6727631153119817e+09	INFO	net-test	Testing connectivity	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "Metrics", "Port": 9612}
1.6727631153187876e+09	INFO	net-test	Testing connectivity PASSED	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "Metrics", "Port": 9612}
1.6727631153189638e+09	INFO	net-test	Testing connectivity	{"Host": "net-test-operator-server.coherence.svc", "PortName": "OperatorRest", "Port": 8000}
1.6727631153265746e+09	INFO	net-test	Testing connectivity PASSED	{"Host": "net-test-operator-server.coherence.svc", "PortName": "OperatorRest", "Port": 8000}
1.6727631153267298e+09	INFO	net-test	Testing connectivity	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "Echo", "Port": 7}
1.6727631153340726e+09	INFO	net-test	Testing connectivity PASSED	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "Echo", "Port": 7}
1.6727631153342876e+09	INFO	net-test	Testing connectivity	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "ClusterPort", "Port": 7574}
1.6727631153406997e+09	INFO	net-test	Testing connectivity PASSED	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "ClusterPort", "Port": 7574}</markup>

<p>The test client successfully connected to the Coherence cluster port (7475), the two unicast ports (7575 and 7576),
the Coherence management port (30000), the Coherence metrics port (9612), the Operator REST port (8000), and the echo port (7).</p>

<p>The test Job can then be deleted:</p>

<markup
lang="bash"

>kubectl -n coh-test delete -f examples/095_network_policies/manifests/net-test-coherence.yaml</markup>

</div>

<h4 id="_testing_the_operator_web_hook">Testing the Operator Web Hook</h4>
<div class="section">
<p>The Operator has a web-hook that k8s calls to validate Coherence resource configurations and to provide default values.
Web hooks in Kubernetes use TLS by default and listen on port 443. The Operator server simulator also listens on port 443
to allow this connectivity to be tested.</p>

<p>The network policy in this example that allows ingress to the web-hook allows any client to connect.
This is because it is not always simple to work out the IP address that the API server will connect to the web-hook from.</p>

<p>We can use the network tester to simulate this by running a Job that will connect to the web hook port.
The web-hook test job in this example does not label the Pod and can be run from the default namespace to simulate a random
external connection.</p>

<markup
lang="bash"

>kubectl -n default apply -f examples/095_network_policies/manifests/net-test-webhook.yaml</markup>

<p>We can then check the results of the Job by looking at the Pod log.</p>

<markup
lang="bash"

>kubectl -n default logs $(kubectl -n default get pod -l 'coherenceNetTest=webhook-client' -o name)</markup>

<p>The output from a successful test will look like this:</p>

<markup


>1.6727639834559627e+09	INFO	runner	Operator Version: 3.3.2
1.6727639834562948e+09	INFO	runner	Operator Build Date: 2023-01-03T12:25:58Z
1.6727639834563956e+09	INFO	runner	Operator Built By: jonathanknight
1.6727639834565024e+09	INFO	runner	Operator Git Commit: c8118585b8f3d72b083ab1209211bcea364c85c5
1.6727639834566057e+09	INFO	runner	Go Version: go1.19.2
1.6727639834567096e+09	INFO	runner	Go OS/Arch: linux/amd64
1.6727639834570327e+09	INFO	net-test	Starting test	{"Name": "Web-Hook Client"}
1.6727639834571698e+09	INFO	net-test	Testing connectivity	{"Host": "net-test-operator-server.coherence.svc", "PortName": "WebHook", "Port": 443}
1.6727639834791095e+09	INFO	net-test	Testing connectivity PASSED	{"Host": "net-test-operator-server.coherence.svc", "PortName": "WebHook", "Port": 443}</markup>

<p>We can see that the client successfully connected to port 443.</p>

<p>The test Job can then be deleted:</p>

<markup
lang="bash"

>kubectl -n default delete -f examples/095_network_policies/manifests/net-test-webhook.yaml</markup>

</div>

<h4 id="_testing_ad_hoc_ports">Testing Ad-Hoc Ports</h4>
<div class="section">
<p>The test client is able to test connectivity to any host and port. For example suppose we want to simulate a Prometheus Pod
connecting to the metrics port of a Coherence cluster. The server simulator is listening on port 9612, so we need to run
the client to connect to that port.</p>

<p>We can create a Job yaml file to run the test client. As the test will simulate a Prometheus client we add the labels
that a standard Prometheus Pod would have and that we also use in the network policies in this example.</p>

<p>In the Job yaml, we need to set the <code>HOST</code>, <code>PORT</code> and optionally the <code>PROTOCOL</code> environment variables.
In this test, the host is the DNS name for the Service created for the Coherence server simulator <code>net-test-coherence-server.coh-test.svc</code>, the port is the metrics port <code>9612</code> and the protocol is <code>tcp</code>.</p>

<markup
lang="yaml"
title="manifests/net-test-client.yaml"
>apiVersion: batch/v1
kind: Job
metadata:
  name: test-client
  labels:
    app.kubernetes.io/name: prometheus
    coherenceNetTest: client
spec:
  template:
    metadata:
      labels:
        app.kubernetes.io/name: prometheus
        coherenceNetTest: client
    spec:
      containers:
      - name: net-test
        image: ghcr.io/oracle/coherence-operator:3.4.2
        env:
          - name: HOST
            value: net-test-coherence-server.coh-test.svc
          - name: PORT
            value: "9612"
          - name: PROTOCOL
            value: tcp
        command:
          - /files/runner
        args:
          - net-test
          - client
      restartPolicy: Never
  backoffLimit: 4</markup>

<p>We need to run the test Job in the <code>monitoring</code> namespace, which is the same namespace that Prometheus is
usually deployed into.</p>

<markup
lang="bash"

>kubectl -n monitoring apply -f examples/095_network_policies/manifests/net-test-client.yaml</markup>

<p>We can then check the results of the Job by looking at the Pod log.</p>

<markup
lang="bash"

>kubectl -n monitoring logs $(kubectl -n monitoring get pod -l 'coherenceNetTest=client' -o name)</markup>

<p>The output from a successful test will look like this:</p>

<markup


>1.6727665901488597e+09	INFO	runner	Operator Version: 3.3.2
1.6727665901497366e+09	INFO	runner	Operator Build Date: 2023-01-03T12:25:58Z
1.6727665901498337e+09	INFO	runner	Operator Built By: jonathanknight
1.6727665901498716e+09	INFO	runner	Operator Git Commit: c8118585b8f3d72b083ab1209211bcea364c85c5
1.6727665901498966e+09	INFO	runner	Go Version: go1.19.2
1.6727665901499205e+09	INFO	runner	Go OS/Arch: linux/amd64
1.6727665901501486e+09	INFO	net-test	Starting test	{"Name": "Simple Client"}
1.6727665901501985e+09	INFO	net-test	Testing connectivity	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "net-test-coherence-server.coh-test.svc-9612", "Port": 9612}
1.6727665901573336e+09	INFO	net-test	Testing connectivity PASSED	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "net-test-coherence-server.coh-test.svc-9612", "Port": 9612}</markup>

<p>We can see that the test client successfully connected to the Coherence cluster member simulator on port 9612.</p>

<p>The test Job can then be deleted:</p>

<markup
lang="bash"

>kubectl -n monitoring delete -f examples/095_network_policies/manifests/net-test-client.yaml</markup>

</div>

<h4 id="_test_with_network_policies">Test with Network Policies</h4>
<div class="section">
<p>All the above tests ran successfully without any network policies. We can now start to apply policies and re-run the
tests to see what happens.</p>

<p>In a secure environment we would start with a policy that blocks all access and then gradually open up required ports.
We can apply the <code>deny-all.yaml</code> policy and then re-run the tests. We should apply the policy to both of the namespaces we are using in this example:</p>

<markup
lang="bash"

>kubectl -n coherence apply -f examples/095_network_policies/manifests/deny-all.yaml
kubectl -n coh-test apply -f examples/095_network_policies/manifests/deny-all.yaml</markup>

<p>Now, re-run the Operator test client:</p>

<markup
lang="bash"

>kubectl -n coherence apply -f examples/095_network_policies/manifests/net-test-operator.yaml</markup>

<p>and check the result:</p>

<markup
lang="bash"

>kubectl -n coherence logs $(kubectl -n coherence get pod -l 'coherenceNetTest=operator-client' -o name)</markup>

<markup


>1.6727671834237397e+09	INFO	runner	Operator Version: 3.3.2
1.6727671834238796e+09	INFO	runner	Operator Build Date: 2023-01-03T12:25:58Z
1.6727671834239576e+09	INFO	runner	Operator Built By: jonathanknight
1.6727671834240365e+09	INFO	runner	Operator Git Commit: c8118585b8f3d72b083ab1209211bcea364c85c5
1.6727671834240875e+09	INFO	runner	Go Version: go1.19.2
1.6727671834241736e+09	INFO	runner	Go OS/Arch: linux/amd64
1.6727671834244306e+09	INFO	net-test	Starting test	{"Name": "Operator Simulator"}
1.6727671834245417e+09	INFO	net-test	Testing connectivity	{"PortName": "K8s API Server"}
1.6727672134268515e+09	INFO	net-test	Testing connectivity FAILED	{"PortName": "K8s API Server", "Error": "Get \"https://10.96.0.1:443/version?timeout=32s\": dial tcp 10.96.0.1:443: i/o timeout"}
1.6727672134269848e+09	INFO	net-test	Testing connectivity	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "Health", "Port": 6676}
1.6727672234281697e+09	INFO	net-test	Testing connectivity FAILED	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "Health", "Port": 6676, "Error": "dial tcp: lookup net-test-coherence-server.coh-test.svc: i/o timeout"}</markup>

<p>We can see that the test client failed to connect to the Kubernetes API server and failed to connect
to the Coherence cluster health port. This means the deny-all policy is working.</p>

<p>We can now apply the various polices to fix the test</p>

<markup
lang="bash"

>kubectl -n coherence apply -f examples/095_network_policies/manifests/allow-dns.yaml
kubectl -n coherence apply -f examples/095_network_policies/manifests/allow-k8s-api-server.yaml
kubectl -n coherence apply -f examples/095_network_policies/manifests/allow-operator-cluster-member-egress.yaml
kubectl -n coherence apply -f examples/095_network_policies/manifests/allow-operator-rest-ingress.yaml
kubectl -n coherence apply -f examples/095_network_policies/manifests/allow-webhook-ingress-from-all.yaml

kubectl -n coh-test apply -f examples/095_network_policies/manifests/allow-dns.yaml
kubectl -n coh-test apply -f examples/095_network_policies/manifests/allow-cluster-member-access.yaml
kubectl -n coh-test apply -f examples/095_network_policies/manifests/allow-cluster-member-operator-access.yaml
kubectl -n coh-test apply -f examples/095_network_policies/manifests/allow-metrics-ingress.yaml</markup>

<p>Now, delete and re-run the Operator test client:</p>

<markup
lang="bash"

>kubectl -n coherence delete -f examples/095_network_policies/manifests/net-test-operator.yaml
kubectl -n coherence apply -f examples/095_network_policies/manifests/net-test-operator.yaml</markup>

<p>and check the result:</p>

<markup
lang="bash"

>kubectl -n coherence logs $(kubectl -n coherence get pod -l 'coherenceNetTest=operator-client' -o name)</markup>

<p>Now with the policies applied the test should have passed.</p>

<markup


>1.6727691273634596e+09	INFO	runner	Operator Version: 3.3.2
1.6727691273635025e+09	INFO	runner	Operator Build Date: 2023-01-03T12:25:58Z
1.6727691273635256e+09	INFO	runner	Operator Built By: jonathanknight
1.6727691273635616e+09	INFO	runner	Operator Git Commit: c8118585b8f3d72b083ab1209211bcea364c85c5
1.6727691273637156e+09	INFO	runner	Go Version: go1.19.2
1.6727691273637407e+09	INFO	runner	Go OS/Arch: linux/amd64
1.6727691273639407e+09	INFO	net-test	Starting test	{"Name": "Operator Simulator"}
1.6727691273639877e+09	INFO	net-test	Testing connectivity	{"PortName": "K8s API Server"}
1.6727691273857167e+09	INFO	net-test	Testing connectivity PASSED	{"PortName": "K8s API Server", "Version": "v1.24.7"}
1.6727691273858056e+09	INFO	net-test	Testing connectivity	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "Health", "Port": 6676}
1.6727691273933685e+09	INFO	net-test	Testing connectivity PASSED	{"Host": "net-test-coherence-server.coh-test.svc", "PortName": "Health", "Port": 6676}</markup>

<p>The other tests can also be re-run and should also pass.</p>

</div>

<h4 id="_clean_up">Clean-Up</h4>
<div class="section">
<p>Once the tests are completed, the test servers and Jobs can be deleted.</p>

<markup
lang="bash"

>kubectl -n coherence delete -f examples/095_network_policies/manifests/net-test-operator-server.yaml
kubectl -n coh-test delete -f examples/095_network_policies/manifests/net-test-coherence-server.yaml</markup>

</div>
</div>
</div>
</doc-view>
