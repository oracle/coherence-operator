<doc-view>

<h2 id="_using_network_policies">Using Network Policies</h2>
<div class="section">
<p>This example covers running the Coherence Operator and Coherence clusters in Kubernetes with network policies.
In Kubernetes, a <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/services-networking/network-policies/">Network Policy</a>
is an application-centric construct which allow you to specify how a pod is allowed to communicate with various network
"entities" (we use the word "entity" here to avoid overloading the more common terms such as "endpoints" and "services",
which have specific Kubernetes connotations) over the network.</p>


<h3 id="_introduction">Introduction</h3>
<div class="section">
<p>Kubernetes network policies specify the access permissions for groups of pods, similar to security groups in the
cloud are used to control access to VM instances and similar to firewalls.
The default behaviour of a Kubernetes cluster is to allow all Pods to freely talk to each other.
Whilst this sounds insecure, originally Kubernetes was designed to orchestrate services that communicated with each other,
it was only later that network policies were added.</p>

<p>A network policy is applied to a Kubernetes namespace and controls ingress into and egress out of Pods in that namespace.</p>

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
  podSelector: { }
  policyTypes:
    - Ingress
    - Egress</markup>

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

<p>The policy below allows all Pods (using <code>podSelector: {}</code>) egress to UDP port 53 in all namespaces.</p>

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
          port: 53</markup>

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
          port: 53</markup>

<p>The policy above can be installed into the <code>coherence</code> namespace with the following command:</p>

<markup
lang="bash"

>kubectl -n coherence apply -f manifests/allow-dns-kube-system.yaml</markup>

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

<p>In the above case the IP address of the API server would be <code>192.168.99.100</code>.</p>

<p>In a simple KinD development cluster, the API server IP address can be obtained using <code>kubectl</code> as shown below:</p>

<markup
lang="bash"

>$ kubectl get pod -n kube-system kube-apiserver-operator-control-plane -o wide
NAME                                    READY   STATUS    RESTARTS   AGE     IP           NODE
kube-apiserver-operator-control-plane   1/1     Running   0          7h43m   172.18.0.5   operator-control-plane</markup>

<p>In the above case the IP address of the API server would be <code>172.18.0.5</code>.</p>

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
  egress:
    - to:
        - ipBlock:
            cidr: 172.18.0.5/32
      ports:
        - port: 6443
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

<p>The trick here is to know where the webhook call is coming from so that a network policy can be sufficiently secure.</p>

<p>The simplest solution is to allow ingress from any IP address to the webhook with a policy like that shown below.
This policy uses and empty <code>from: []</code> attribute, which allows access from anywhere.</p>

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

<p>Allowing all access to the webhook is not very secure, so a more restrictive <code>from</code> attribute could be used to limit
access to the IP address of the Kubernetes API server.</p>

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
  name: coherence-operator-cluster-member-egress
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

<p>When using the default Coherence images, for example <code>ghcr.io/oracle/coherence-ce:21.12.4</code> the Extend proxy is already
configured to run on a fixed port <code>20000</code>. When using this image, or any image that uses the default Coherence cache
configuration file, this port can be changed by setting the <code>COHERENCE_EXTEND_PORT</code> environment variable.</p>

<p>When using the Coherence Concurrent extensions over Extend, the Concurrent Extend proxy also needs to be configured with a fixed port.
When using the default Coherence images, for example <code>ghcr.io/oracle/coherence-ce:21.12.4</code> the Concurrent Extend proxy is already
configured to run on a fixed port <code>20001</code>. When using this image, or any image that uses the default Coherence cache
configuration file, this port can be changed by setting the <code>COHERENCE_CONCURRENT_EXTEND_PORT</code> environment variable.</p>

<p>For the examples below, a <code>Coherence</code> deployment has the following configuration.
This will expose Extend on a port named <code>extend</code> with a port number of <code>20000</code>, and a port named <code>extend-concurrent</code>
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
    - name: extend-concurrent
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
        - port: extend-concurrent
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
        - port: extend-concurrent
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

<p>When using the default Coherence images, for example <code>ghcr.io/oracle/coherence-ce:21.12.4</code> the gRPC proxy is already
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
        - port: extend-concurrent
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
</div>
</doc-view>
