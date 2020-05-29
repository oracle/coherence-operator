<doc-view>

<h2 id="_well_known_addressing_and_cluster_discovery">Well Known Addressing and Cluster Discovery</h2>
<div class="section">
<p>A Coherence cluster is made up of one or more JVMs. In order for these JVMs to form a cluster they need to be able to
discover other cluster members. The default mechanism for discovery is multicast broadcast but this does not work in
most container environments. Coherence provides an alternative mechanism where the addresses of the hosts where the
members of the cluster will run is provided in the form of a
<a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/develop-applications/setting-cluster.html#GUID-E8CC7C9A-5739-4D12-B88E-A3575F20D63B">"well known address" (or WKA) list</a>.
This address list is then used by Coherence when it starts in a JVM to discover other cluster members running on the
hosts in the WKA list.</p>

<p>When running in containers each container is effectively a host and has its own host name and IP address (or addresses)
and in Kubernetes it is the <code>Pod</code> that is effectively a host. When starting a container it is usually not possible to
know in advance what the host names of the containers or <code>Pods</code> will be so there needs to be another solution to
providing the WKA list.</p>

<p>Coherence processes a WKA list it by performing a DNS lookup for each host name in the list. If a host name resolves
to more than one IP address then <em>all</em> of those IP addresses will be used in cluster discovery. This feature of Coherence
when combined with Kubernetes <code>Services</code> allows discovery of cluster members without resorting to a custom discovery
mechanism.</p>

<p>A Kubernetes <code>Service</code> has a DNS name and that name will resolve to all the IP addresses of the <code>Pods</code> that match
that <code>Service</code> selector. This means that a Coherence JVM only needs to be given the DNS name of a <code>Service</code> as the
single host name in its WKA list so that it will form a cluster with any other JVM using in a Pod matching the selector.</p>

<p>When the Coherence Operator creates reconciles a <code>Coherence</code> CRD configuration to create a running set of <code>Pods</code>
it creates a headless service specifically for the purposes of WKA for that <code>Coherence</code> resource with a selector that
matches any Pod with the same cluster name.</p>

<p>For example, if a <code>Coherence</code> resource is created with the following yaml:</p>

<markup
lang="yaml"
title="test-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  cluster: test-cluster <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">In this yaml the <code>Coherence</code> resource has a cluster name of <code>test-cluster</code></li>
</ul>
<p>The Operator will create a <code>Service</code> for the <code>Coherence</code> resource using the same name as the deployment
with a <code>-wka</code> suffix.
So in the example above the Operator would create a <code>Service</code> with the name <code>storage-wka</code>.</p>

<p>The yaml for the WKA <code>Service</code> would look like the following:</p>

<markup
lang="yaml"
title="wka-service.yaml"
>apiVersion: v1
kind: Service
metadata:
  name: storage-wka                                                  <span class="conum" data-value="1" />
  labels:
    coherenceCluster: test-cluster
    component: coherenceWkaService
spec:
  clusterIP: None                                                    <span class="conum" data-value="2" />
  publishNotReadyAddresses: true                                     <span class="conum" data-value="3" />
  ports:
    - name: coherence                                                <span class="conum" data-value="4" />
      protocol: TCP
      port: 7
      targetPort: 7
  selector:
    coherenceCluster: test-cluster                                   <span class="conum" data-value="5" />
    component: coherencePod</markup>

<ul class="colist">
<li data-value="1">The <code>Service</code> name is made up of the cluster name with the suffix <code>-wka</code> so in this case <code>storage-wka</code></li>
<li data-value="2">The service has a <code>clusterIP</code> of <code>None</code> so it is headless</li>
<li data-value="3">The <code>Service</code> is configured to allow unready <code>Pods</code> so that all <code>Pods</code> matching the selector will be resolved as
members of this service regardless of their ready state. This is important so that Coherence JVMs can discover other
members before they are fully ready.</li>
<li data-value="4">A single port is exposed, in this case the echo port (7), even though nothing in the Coherence <code>Pods</code> binds to this
port. Ideally no port would be included, but a Kubernetes service has to have at least one port defined.</li>
<li data-value="5">The selector will match all <code>Pods</code> with the labels <code>coherenceCluster=test-cluster</code> and <code>component=coherencePod</code>
which are labels that the Coherence Operator will assign to all <code>Pods</code> in this cluster</li>
</ul>
<p>Because this <code>Service</code> is created in the same <code>Namespace</code> as the deployment&#8217;s <code>Pods</code> the JVMs can use
the raw <code>Service</code> name as the WKA list, in the example above the WKA list would just be <code>test-cluster-wka</code>.</p>


<h3 id="_exclude_a_deployment_from_wka">Exclude a Deployment From WKA</h3>
<div class="section">
<p>In some situations it may be desirable to exclude the Pods belonging to certain deployments in the cluster from being
members of the well known address list. For example certain K8s network configurations such as host networking can
cause issues with WKA if other deployments in the cluster are using host networking.</p>

<p>A role can be excluded from the WKA list by setting the <code>excludeFromWKA</code> field of the <code>coherence</code> section of the
deployment&#8217;s spec to <code>true</code>.</p>

<markup
lang="yaml"
title="test-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test-client
spec:
  cluster: `my-cluster`    <span class="conum" data-value="1" />
  coherence:
    excludeFromWKA: true   <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>cluster</code> field is set to the name of the Coherence cluster that this deployment wil be part of (there is no
point in excluding a deployment from WKA unless it is part of a wider cluster).</li>
<li data-value="2">The <code>excludeFromWKA</code> field is <code>true</code> so that <code>Pods</code> in the <code>test-client</code> deployment will not form part of the WKA
list for the Coherence cluster.</li>
</ul>
<div class="admonition warning">
<p class="admonition-inline">The operator does not validate the <code>excludeFromWKA</code> field for a deployment so it is possible to try to create
a cluster where all of the deployment have <code>excludeFromWKA</code> set to <code>true</code> which will cause the cluster fail to start.</p>
</div>
<div class="admonition warning">
<p class="admonition-inline">When excluding a deployment from WKA it is important that at least one deployment that is part of the WKA list
has been started first otherwise the non-WKA role members cannot start.Eventually the K8s readiness probe for these Pods
would time-out causing K8s to restart them but this would not be a desirable way to start a cluster.
The start-up order can be controlled by configuring the deployment&#8217;s <code>startQuorum</code> list, as described in the documentation
section on <router-link to="/ordering/010_overview">deployment start-up ordering</router-link>.</p>
</div>
</div>
</div>
</doc-view>
