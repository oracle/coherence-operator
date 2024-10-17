<doc-view>

<h2 id="_what_is_the_coherence_operator">What is the Coherence Operator?</h2>
<div class="section">
<p>The Coherence Operator is a <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/extend-kubernetes/operator/">Kubernetes Operator</a> that
is used to manage <a id="" title="" target="_blank" href="https://oracle.github.io/coherence">Oracle Coherence</a> clusters in Kubernetes.
The Coherence Operator takes on the tasks of that human DevOps resource might carry out when managing Coherence clusters,
such as configuration, installation, safe scaling, management and metrics.</p>

<p>The Coherence Operator is a Go based application built using the <a id="" title="" target="_blank" href="https://github.com/operator-framework/operator-sdk">Operator SDK</a>.
It is distributed as a Docker image and Helm chart for easy installation and configuration.</p>

</div>

<h2 id="_why_use_the_coherence_kubernetes_operator">Why use the Coherence Kubernetes Operator</h2>
<div class="section">
<p>Using the Coherence Operator to manage Coherence clusters running in Kubernetes has many advantages over just deploying
and running clusters with the resources provided by Kubernetes.
Coherence can be treated as just another library that your application depends on and uses and hence, a Coherence
application can run in Kubernetes without requiring the Operator, but in this case there are
a number of things that the DevOps team for an application would need to build or do manually.</p>


<h3 id="_cluster_discovery">Cluster Discovery</h3>
<div class="section">
<p>JVMs that run as Coherence cluster members need to discover the other members of the cluster.
This is discussed in the <router-link to="/docs/coherence/070_wka">Coherence Well Known Addressing</router-link> section of the documentation.
When using the Operator the well known addressing configuration for clusters is managed automatically to allow a Coherence
deployment to create its own cluster or to join with other deployments to form larger clusters.</p>

</div>

<h3 id="_better_fault_tolerant_data_distribution">Better Fault Tolerant Data Distribution</h3>
<div class="section">
<p>The Operator configures the Coherence site and rack properties for cluster members based on Kubernetes Node topology
labels. This allows Coherence to better distribute data across sites when a Kubernetes cluster spans availability domains.</p>

</div>

<h3 id="_safe_scaling">Safe Scaling</h3>
<div class="section">
<p>When scaling down a Coherence cluster, care must be taken to ensure that there will be no data loss.
This typically means scaling down by a single Pod at a time and waiting for the cluster to become "safe" before scaling
down the next Pod.
The Operator has built in functionality to do this, so scaling a Coherence cluster is as simple as scaling any other
Kubernetes Deployment or StatefulSet.</p>

</div>

<h3 id="_autoscaling">Autoscaling</h3>
<div class="section">
<p>Alongside safe scaling, because the Coherence CRD supports the Kubernetes scale sub-resource it is possible to configure
the Kubernetes Horizontal Pod Autoscaler to scale Coherence
clusters based on metrics.</p>

</div>

<h3 id="_readiness_probes">Readiness Probes</h3>
<div class="section">
<p>The Operator has an understanding of when a Coherence JVM is "ready", so it configures a readiness probe that k8s will
use to signal whether a Pod is ready or not.</p>

</div>

<h3 id="_persistence">Persistence</h3>
<div class="section">
<p>Using the Operator makes it simple to configure and use Coherence Persistence, storing data on Kubernetes Persistent
Volumes to allow state to be maintained between cluster restarts.</p>

</div>

<h3 id="_graceful_shutdown">Graceful Shutdown</h3>
<div class="section">
<p>When a Coherence cluster is deployed with persistence enabled, the Operator will gracefully shutdown a cluster by suspending
services before stopping all the Pods.
This ensures that all persistence files are properly closed and allows for quicker recovery and restart of the cluster.
Without the Operator, if a cluster is shutdown, typically by removing the controlling StatefulSet from Kubernetes then
the Pods will be shutdown but not all at the same time.
It is obviously impossible for k8s to kill all the Pods at the exact same instant in time. As some Pods die the remaining
storage enabled Pods will be trying to recover data for the lost Pods, this can cause a lot of needles work and moving of
data over the network. It is much cleaner to suspend all the services before shutdown.</p>

</div>

<h3 id="_simpler_configuration">Simpler Configuration</h3>
<div class="section">
<p>The Coherence CRD is designed to make the more commonly used configuration parameters for Coherence, and the JVM simpler
to configure. The Coherence CRD is simple to use, in fact none of its fields are mandatory, so an application can be
deployed with nothing more than a name, and a container image.</p>

</div>

<h3 id="_dual_stack_kubernetes_clusters">Dual-Stack Kubernetes Clusters</h3>
<div class="section">
<p>The Operator supports running Coherence on dual-stack IPv4 and IPv6 Kubernetes clusters.</p>

</div>

<h3 id="_consistency">Consistency</h3>
<div class="section">
<p>By using the Operator to manage Coherence clusters all clusters are configured and managed the same way making it easier
for DevOps to manage multiple clusters and applications.</p>

</div>

<h3 id="_expertise">Expertise</h3>
<div class="section">
<p>The Operator has been built and tested by the Coherence engineering team, who understand Coherence and the various scenarios
and edge cases that can occur when managing Coherence clusters at scale in Kubernetes.</p>

</div>
</div>

<h2 id="_coherence_clusters">Coherence Clusters</h2>
<div class="section">
<p>A Coherence cluster is a number of distributed Java Virtual Machines (JVMs) that communicate to form a single coherent cluster.
In Kubernetes, this concept can be related to a number of Pods that form a single cluster.
In each <code>Pod</code> is a JVM running a Coherence <code>DefaultCacheServer</code>, or a custom application using Coherence.</p>

<p>The operator uses a Kubernetes Custom Resource Definition (CRD) to represent a group of members in a Coherence cluster.
Typically, a deployment would be used to configure one or more members of a specific role in a cluster.
Every field in the <code>Coherence</code> CRD <code>Spec</code> is optional, so a simple cluster can be defined in  yaml as:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: my-cluster <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">In this case the <code>metadata.name</code> field in the <code>Coherence</code> resource yaml will be used as the Coherence cluster name.</li>
</ul>
<p>The operator will use default values for fields that have not been entered, so the above yaml will create
a Coherence deployment using a <code>StatefulSet</code> with a replica count of three, which means that will be three storage
enabled Coherence <code>Pods</code>.</p>

<p>See the <router-link to="/docs/about/04_coherence_spec">Coherence CRD spec</router-link> page for details of all the fields in the CRD.</p>

<p>In the above example no <code>spec.image</code> field has been set, so the Operator will use a publicly available Coherence CE
image as its default. These images are meant for demos, POCs and experimentation, but for a production application you
should build your own image.</p>

</div>

<h2 id="_using_commercial_coherence_versions">Using Commercial Coherence Versions</h2>
<div class="section">
<div class="admonition note">
<p class="admonition-inline">Whilst the Coherence CE version can be freely deployed anywhere, if your application image uses a commercial
version of Oracle Coherence then you are responsible for making sure your deployment has been properly licensed.</p>
</div>
<p>Oracle&#8217;s current policy is that a license will be required for each Kubernetes Node that images are to be pulled to.
While an image exists on a node it is effectively the same as having installed the software on that node.</p>

<p>One way to ensure that the Pods of a Coherence deployment only get scheduled onto nodes that meet the
license requirement is to configure Pod scheduling, for example a node selector. Node selectors, and other scheduling,
is simple to configure in the <code>Coherence</code> CRD, see the <router-link to="/docs/other/090_pod_scheduling">scheduling documentation</router-link></p>

<p>For example, if a commercial Coherence license exists such that a sub-set of nodes in a Kubernetes cluster
have been covered by the license then those nodes could all be given a label, e.g. <code>coherenceLicense=true</code></p>

<p>When creating a <code>Coherence</code> deployment specify a node selector to match the label:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  image: my-app:1.0.0         <span class="conum" data-value="1" />
  nodeSelector:
    coherenceLicense: 'true'  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>my-app:1.0.0</code> image contains a commercial Coherence version.</li>
<li data-value="2">The <code>nodeSelector</code> will ensure Pods only get scheduled to nodes with the <code>coherenceLicense=true</code> label.</li>
</ul>
<p>There are other ways to configure Pod scheduling supported by the Coherence Operator (such as taints and tolerations)
and there are alternative ways to restrict nodes that Pods can be schedule to, for example a namespace in kubernetes
can be restricted to a sub-set of the cluster&#8217;s nodes. Using a node selector as described above is probably the
simplest approach.</p>

</div>
</doc-view>
