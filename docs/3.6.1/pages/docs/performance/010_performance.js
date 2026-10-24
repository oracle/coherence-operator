<doc-view>

<h2 id="_performance_testing_in_kubernetes">Performance Testing in Kubernetes</h2>
<div class="section">
<p>Many customers use Coherence because they want access to data at in-memory speeds. Customers who want the best performance from their application typically embark on performance testing and load testing of their application. When doing this sort of testing on Kubernetes, it is useful to understand the ways that your Kubernetes environment can impact your test results.</p>

</div>

<h2 id="_where_are_your_nodes">Where are your Nodes?</h2>
<div class="section">
<p>When an application has been deployed into Kubernetes, pods will typically be distributed over many nodes in the Kubernetes cluster.
When deploying into Kubernetes cluster in the cloud, for example on Oracle OKE, the nodes can be distributed across different availability zones. These zones are effectively different data centers, meaning that the network speed can differ considerable between nodes in different zones.
Performance testing in this sort of environment can be difficult if you use default Pod scheduling. Different test runs could distribute Pods to different nodes, in different zones, and skew results depending on how "far" test clients and servers are from each other.
For example, when testing a simple Coherence <code>EntryProcessor</code> invocation in a Kubernetes cluster spread across zones, we saw the 95% response time when the client and server were in the same zone was 0.1 milli-seconds. When the client and server were in different zones, the 95% response time could be as high as 0.8 milli-seconds. This difference is purely down to the network distance between nodes. Depending on the actual use-cases being tested, this difference might not have much impact on overall response times, but for simple operations it can be a significant enough overhead to impact test results.</p>

<p>The solution to the issue described above is to use Pod scheduling to fix the location of the Pods to be used for tests. In a cluster like Oracle OKE, this would ensure all the Pods will be scheduled into the same availability zone.</p>


<h3 id="_finding_node_zones">Finding Node Zones</h3>
<div class="section">
<p>This example is going to talks about scheduling Pods to a single availability zone in a Kubernetes cluster in the cloud. Pod scheduling in this way uses Node labels, and in fact any label on the Nodes in your cluster could be used to fix the location of the Pods.</p>

<p>To schedule all the Coherence Pods into a single zone we first need to know what zones we have and what labels have used.
The standard Kubernetes Node label for the availability zone is <code>topology.kubernetes.io/zone</code> (as documented in the <a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/labels-annotations-taints/">Kubernetes Labels Annotations and Taints</a> documentation). To slightly confuse the situation, prior to Kubernetes 1.17 the label was <code>failure-domain.beta.kubernetes.io/zone</code>, which has now been deprecated. Some Kubernetes clusters, even after 1.17, still use the deprecated label, so you need to know what labels your Nodes have.</p>

<p>Run the following command so list the nodes in a Kubernetes cluster with the value of the two zone labels for each node:</p>

<markup
lang="bash"

>kubectl get nodes -L topology.kubernetes.io/zone,failure-domain.beta.kubernetes.io/zone</markup>

<p>The output will be something like this:</p>

<markup


>NAME      STATUS   ROLES   AGE   VERSION   ZONE             ZONE
node-1    Ready    node    66d   v1.19.7   US-ASHBURN-AD-1
node-2    Ready    node    66d   v1.19.7   US-ASHBURN-AD-2
node-3    Ready    node    66d   v1.19.7   US-ASHBURN-AD-3
node-4    Ready    node    66d   v1.19.7   US-ASHBURN-AD-2
node-5    Ready    node    66d   v1.19.7   US-ASHBURN-AD-3
node-6    Ready    node    66d   v1.19.7   US-ASHBURN-AD-1</markup>

<p>In the output above the first <code>Zone</code> column has values, and the second does not. This means that the zone label used is the first in the label list in our <code>kubectl</code> command, i.e., <code>topology.kubernetes.io/zone</code>.</p>

<p>If the nodes had been labeled with the second, deprecated, label in the <code>kubectl</code> command list <code>failure-domain.beta.kubernetes.io/zone</code> the output would look like this:</p>

<markup


>NAME      STATUS   ROLES   AGE   VERSION   ZONE   ZONE
node-1    Ready    node    66d   v1.19.7          US-ASHBURN-AD-1
node-2    Ready    node    66d   v1.19.7          US-ASHBURN-AD-2
node-3    Ready    node    66d   v1.19.7          US-ASHBURN-AD-3
node-4    Ready    node    66d   v1.19.7          US-ASHBURN-AD-2
node-5    Ready    node    66d   v1.19.7          US-ASHBURN-AD-3
node-6    Ready    node    66d   v1.19.7          US-ASHBURN-AD-1</markup>

<p>From the list of nodes above we can see that there are three zones, <code>US-ASHBURN-AD-1</code>, <code>US-ASHBURN-AD-2</code> and <code>US-ASHBURN-AD-3</code>.
In this example we will schedule all the Pods to zome <code>US-ASHBURN-AD-1</code>.</p>

</div>

<h3 id="_scheduling_pods_of_a_coherence_cluster">Scheduling Pods of a Coherence Cluster</h3>
<div class="section">
<p>The <code>Coherence</code> CRD supports a number of ways to schedule Pods, as described in the <router-link to="/docs/other/090_pod_scheduling">Configure Pod Scheduling</router-link> documentation. Using node labels is the simplest of the scheduling methods.
In this case we need to schedule Pods onto nodes that have the label <code>topology.kubernetes.io/zone=US-ASHBURN-AD-1</code>.
In the <code>Coherence</code> yaml we use the <code>nodeSelector</code> field.</p>

<markup
lang="yaml"
title="coherence-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  replicas: 3
  nodeSelector:
    topology.kubernetes.io/zone: US-ASHBURN-AD-1</markup>

<p>When the yaml above is applied, a cluster of three Pods will be created, all scheduled onto nodes in the <code>US-ASHBURN-AD-1</code> availability zone.</p>

</div>

<h3 id="_other_performance_factors">Other Performance Factors</h3>
<div class="section">
<p>Depending on the Kubernetes cluster you are using there could be various other factors to bear in mind. Many Kubernetes clusters run on virtual machines, which can be poor for repeated performance comparisons unless you know what else might be running on the underlying hardware that the VM is on. If a test run happens at the same time as another VM is consuming a lot of the underlying hardware resource this can obviously impact the results. Unfortunately bear-metal hardware, the best for repeated performance tests, is not always available, so it is useful to bear this in mind if there is suddenly a strange outlier in the tests.</p>

</div>
</div>
</doc-view>
