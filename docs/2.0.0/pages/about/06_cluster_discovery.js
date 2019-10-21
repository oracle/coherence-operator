<doc-view>

<h2 id="_coherence_cluster_discovery">Coherence Cluster Discovery</h2>
<div class="section">
<p>A Coherence cluster is made up of one or more JVMs. In order for these JVMs to form a cluster they need to be able to
discover other cluster members. The default mechanism for discovery is multicast broadcast but this does not work in
most container environments. Coherence provides an alternative mechanism where the addresses of the hosts where the
members of the cluster will run is provided in the form of a "well known address" (or WKA) list. This address list is
then used by Coherence when it starts in a JVM to discover other cluster members running on the hosts in the WKA list.</p>

<p>When running in containers each container is effectively a host and has its own host name and IP address (or addresses)
and in Kubernetes it is the <code>Pod</code> that is effectively a host. When starting a container it is usually not possible to
know in advance what the host names of the containers or <code>Pods</code> will be so there needs to be another solution to
providing the WKA list.</p>

<p>When Coherence processes a WKA list it will perform a DNS lookup for each host name in the list. If a host name resolves
to more than one IP address then <em>all</em> of those IP addresses are used in cluster discovery. This feature of Coherence
when combined with Kubernetes <code>Services</code> allows discovery of cluster members without resorting to a custom discovery
mechanism.</p>

<p>A Kubernetes <code>Service</code> has a DNS name and that name will resolve to all of the IP addresses for the <code>Pods</code> that match
that <code>Service</code> selector. This means that a Coherence JVM only needs to be given the DNS name of a <code>Service</code> as the
single host name in its WKA list and it will form a cluster with any other JVM using the same host name for WKA and the
same cluster name.</p>

<p>When the Coherence Operator creates resolves a <code>CoherenceCluster</code> configuration into a running set of <code>Pods</code> if creates
a headless service specifically for the purposes of WKA for that cluster.</p>

<p>For example, if a <code>CoherenceCluster</code> is created with the following yaml:</p>

<markup
lang="yaml"
title="test-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">A Coherence cluster will br created with a cluster name <code>test-cluster</code></li>
</ul>
<p>The yaml for the WKA <code>Service</code> would look like the following:</p>

<markup
lang="yaml"
title="wka-service.yaml"
>apiVersion: v1
kind: Service
metadata:
  name: test-cluster-wka                                             <span class="conum" data-value="1" />
  annotations:
    service.alpha.kubernetes.io/tolerate-unready-endpoints: "true"   <span class="conum" data-value="2" />
  labels:
    coherenceCluster: test-cluster
    component: coherenceWkaService
spec:
  clusterIP: None                                                    <span class="conum" data-value="3" />
  ports:
    - name: coherence                                                <span class="conum" data-value="4" />
      protocol: TCP
      port: 7
      targetPort: 7
  selector:
    coherenceCluster: test-cluster                                   <span class="conum" data-value="5" />
    component: coherencePod</markup>

<ul class="colist">
<li data-value="1">The <code>Service</code> name is made up of the cluster name with the suffix <code>-wka</code> so in this case <code>test-cluster-wka</code></li>
<li data-value="2">The <code>Service</code> is configured to allow unready <code>Pods</code> so that all <code>Pods</code> matching the selector will be resolved as
members of this service regardless of their ready state. This is important so that Coherence JVMs can discover other
members before they are fully ready.</li>
<li data-value="3">The service has a <code>clusterIP</code> of <code>Nonde</code> so it is headless</li>
<li data-value="4">A single port is exposed, in this case the echo port (7) even though nothing in the Coherence <code>Pods</code> binds to this
port. Ideally no port would be included but the service has to have at least one port defined.</li>
<li data-value="5">The selector will match all <code>Pods</code> with the labels <code>coherenceCluster=test-cluster</code> and <code>component=coherencePod</code>
which are labels that the Coherence Operator will assign to all <code>Pods</code> in this cluster</li>
</ul>
<p>Because this <code>Service</code> is created in the same <code>Namespace</code> as the rest of the Coherence cluster <code>Pods</code> the JVMs can use
the raw <code>Service</code> name as the WKA list, in the example above the WKA list would just be <code>test-cluster-wka</code>.</p>

</div>
</doc-view>
