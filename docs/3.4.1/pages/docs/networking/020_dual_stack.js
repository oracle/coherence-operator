<doc-view>

<h2 id="_dual_stack_networking">Dual Stack Networking</h2>
<div class="section">
<p>This section describes using Coherence and the Operator with a dual stack Kubernetes cluster,
where Pods and Services can have both IPv4 and IPv4 interfaces.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>This section only really applies to making Coherence bind to the correct local IP address for inter-cluster communication.
Normally for other Coherence endpoints, such as Extend, gRPC, management, metrics, etc. Coherence will bind to all
local addresses ubless specifically configured otherwise.
This means that in and environment such as dual-stack Kubernetes where a Pod has both an IPv4 and IPv6
address, those Coherence endpoints will be reachable using either the IPv4 or IPv6 address of the Pod.</p>
</p>
</div>
<p>Normally, using Coherence on a dual-stack server can cause issues due to the way that Coherence decides which local IP
address to use for inter-cluster communication. Similar problems can occur on any server that multiple IP addresses.
When Coherence is configured to use well known addressing for cluster discovery, a Coherence JVM will choose a local
address that is either in the WKA list, or is on an interface that can route to the WKA addresses.
In a dual stack environment the problem comes when an interface has both IPv4 and IPv6 addresses and Coherence is
inconsistent about which one to choose.</p>

<p>There are a few simple ways to fix this:</p>

<ul class="ulist">
<li>
<p>Set the JVM system property <code>java.net.preferIPv4Stack=true</code> or <code>java.net.preferIPv6Addresses=true</code> to set the Coherence
JVM to use the desired stack. If application code requires both stacks to be available though, this is not a good option.</p>

</li>
<li>
<p>Configure the WKA list to be only IPv4 addresses or IPv6 addresses. Coherence will then choose a matching local address.</p>

</li>
<li>
<p>Set the <code>coherence.localhost</code> system property (or <code>COHERENCE_LOCALHOST</code> environment variable) to the IP address
that Coherence should bind to. In a dual stack environment choose either the IPv4 address or IPv6 address and make sure
that the corresponding addresses are used in the WKA list.</p>

</li>
</ul>

<h3 id="_dual_stack_kubernetes_clusters">Dual Stack Kubernetes Clusters</h3>
<div class="section">
<p>In a dual-stack Kubernetes cluster, Pods will have both an IPv4 and IPv6 address.
These can be seen by looking at the status section of a Pod spec:</p>

<markup
lang="yaml"

>  podIP: 10.244.3.3
  podIPs:
  - ip: 10.244.3.3
  - ip: fd00:10:244:3::3</markup>

<p>The status section will have a <code>podIP</code> field, which is the Pods primary address.
There is also an array of the dual-stack addresses in the <code>podIPs</code> field.
The first address in <code>podIPs</code> is always the same as <code>podIP</code> and is usually the IPv4 address.</p>

<p>A Service in a dual-stack cluster can have a single IP family or multiple IP families configured in its spec.
The Operator will work out of the box if the default IP families configuration for Services is single stack, either IPv4 or IPv6.
When the WKA Service is created it will only be populated with one type of address, and Coherence will bind to the correct type.</p>

<p>In Kubernetes clusters where the WKA service has multiple IP families by default, there are a few options to fix this:</p>

<ul class="ulist">
<li>
<p>Set the JVM system property <code>java.net.preferIPv4Stack=true</code> or <code>java.net.preferIPv6Addresses=true</code> to set the Coherence
JVM to use the desired stack. If application code requires both stacks to be available though, this is not a good option.</p>

</li>
</ul>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  replicas: 3
  jvm:
    args:
      - "java.net.preferIPv4Stack=true"</markup>

<ul class="ulist">
<li>
<p>The <code>COHERENCE_LOCALHOST</code> environment variable can be configured to be the Pods IP address.
Typically, this will be the IPv4 address.</p>

</li>
</ul>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  replicas: 3
  env:
    - name: COHERENCE_LOCALHOST
      valueFrom:
        fieldRef:
          apiVersion: v1
          fieldPath: status.podIP</markup>

<ul class="ulist">
<li>
<p>Since Operator 3.4.1 it is possible to configure the IP family for the WKA Service. The <code>spec.coherence.wka.ipFamily</code>
field can be set to either "IPv4" or "IPv6". This will cause Coherence to bind to the relevant IP address type.</p>

</li>
</ul>
<p>For example, the yaml below will cause Coherence to bind to the IPv6 address.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  replicas: 3
  coherence:
    wka:
      ipFamily: IPv6</markup>

<p>Since Operator 3.4.1 it is also possible to configure the IP families used by the headless service created for the StatefulSet
if this is required.</p>

<p>The yaml below will configure WKA to use only IPv6, the headless Service created for the StatefulSet will be
a dual-stack, IPv4 and IPv6 service.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  replicas: 3
  headlessServiceIpFamilies:
    - IPv4
    - IPv6
  coherence:
    wka:
      ipFamily: IPv6</markup>

<p>The yaml below will configure both WKA and the headless Service created for the StatefulSet to use a single stack IPv6.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  replicas: 3
  headlessServiceIpFamilies:
    - IPv6
  coherence:
    wka:
      ipFamily: IPv6</markup>

</div>

<h3 id="_dual_stack_kubernetes_clusters_without_using_the_operator">Dual Stack Kubernetes Clusters Without Using the Operator</h3>
<div class="section">
<p>If not using the Coherence Operator to manage clusters the same techniques described above can be used to
manually configure Coherence to work correctly.</p>

<p>The simplest option is to ensure that the headless service used for well known addressing is configured to be single stack.
For example, the yaml below configures the service <code>storage-sts</code> to be a single stack IPv6 service.</p>

<markup
lang="yaml"

>apiVersion: v1
kind: Service
metadata:
  name: storage-sts
spec:
  clusterIP: None
  clusterIPs:
  - None
  ipFamilies:
  - IPv6
  ipFamilyPolicy: SingleStack</markup>

<p>If for some reason it is not possible to ise a dedicated single stack service for WKA, then the <code>COHERENCE_LOCALHOST</code>
environment variable can be set in the Pod to be the Pod IP address.</p>

<markup
lang="yaml"

>apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: storage
spec:
  template:
    spec:
      containers:
        - name: coherence
          env:
            - name: COHERENCE_LOCALHOST
              valueFrom:
                fieldRef:
                  apiVersion: v1
                  fieldPath: status.podIP</markup>

</div>
</div>
</doc-view>
