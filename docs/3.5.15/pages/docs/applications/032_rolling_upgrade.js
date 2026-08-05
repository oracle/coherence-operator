<doc-view>

<h2 id="_rolling_upgrades_of_coherence_applications">Rolling Upgrades of Coherence Applications</h2>
<div class="section">
<p>The Coherence Operator supports safe rolling upgrades of Coherence clusters.</p>

</div>

<h2 id="_default_behaviour">Default Behaviour</h2>
<div class="section">
<p>The Coherence Operator uses a StatefulSet to manage the application Pods.
The StatefulSet is configured to perform its default rolling upgrade behaviour.
This means that when a Coherence resource is updated the StatefulSet will control the rolling upgrade.
First a Pod is killed and rescheduled with the updated specification.
When this Pod is "ready" the next Pod is killed and rescheduled, and so on until all the Pods are updated.
Because the default readiness probe configured by the Operator will wait for Coherence members to be "safe"
(i.e. no endangered partitions) and redistribution to be complete when the new Pod is ready, it is safe
to kill the next Pod.</p>

</div>

<h2 id="_custom_rolling_upgrades">Custom Rolling Upgrades</h2>
<div class="section">
<p>The Coherence resource yaml has a field named <code>RollingUpdateStrategy</code> which can be used to override the default
rolling upgrade strategy. The field can be set to one of the following values:</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 50%;">
<col style="width: 50%;">
</colgroup>
<thead>
<tr>
<th>RollingUpdateStrategy</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>Pod</code></td>
<td class="">This is the same as the default behaviour, one Pod at a time is upgraded.</td>
</tr>
<tr>
<td class=""><code>Node</code></td>
<td class="">This strategy will upgrade all Pods on a Node at the same time.</td>
</tr>
<tr>
<td class="">`NodeLabel `</td>
<td class="">This strategy will upgrade all Pods on all Nodes that have a matching value for a give Node label.</td>
</tr>
<tr>
<td class=""><code>Manual</code></td>
<td class="">This strategy is the same as the <code>Manual</code> rolling upgrade configuration for a StatefulSet.</td>
</tr>
</tbody>
</table>
</div>
<p>The default "by Pod" strategy is the slowest but safest strategy.
For a very large cluster upgrading by Pod may take a long time. For example, if each Pod takes two minutes to be
rescheduled and become ready, and a cluster has 100 Pods, that will be 200 minutes, (over three hours) to upgrade.
In a lot of use cases the time taken to upgrade is not an issue, Coherence continues to serve requests while the
upgrade is in progress. But, sometimes a faster upgrade is required, which is where the other strategies can be used.</p>


<h3 id="_upgrade_by_pod">Upgrade By Pod</h3>
<div class="section">
<p>Upgrading by Pod is the default strategy described above.</p>

<p>The <code>Pod</code> strategy is configured by omitting the <code>rollingUpdateStrategy</code> field,
or by setting the <code>rollingUpdateStrategy</code> field to <code>Pod</code> as shown below:</p>

<markup
lang="yaml"
title="cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  rollingUpdateStrategy: Pod
  image: my-app:1.0.0</markup>

</div>

<h3 id="_upgrade_by_node">Upgrade By Node</h3>
<div class="section">
<p>By default, the Operator configures Coherence to be at least "machine safe",
using the Node as the machine identifier. This means that it should be safe to
lose all Pods on a Node. By upgrading multiple Pods at the same time the overall time to perform a
rolling upgrade is less than using the default one Pod at a time behaviour.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>When using the <code>Node</code> strategy where multiple Pods will be killed, the remaining cluster must have enough
capacity to recover the data and backups from the Pods that are killed.</p>

<p>For example, if a cluster of 18 Pods is distributed over three Nodes, each Node will be running six Pods.
When upgrading by Node, six Pods will be killed at the same time, so there must be enough capacity in the
remaining 12 Pods to hold all the data that was in the original 18 Pods.</p>
</p>
</div>
<p>The <code>Node</code> strategy is configured by setting the <code>rollingUpdateStrategy</code> field to <code>Node</code> as shown below:</p>

<markup
lang="yaml"
title="cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  rollingUpdateStrategy: Node
  image: my-app:1.0.0</markup>

</div>

<h3 id="_upgrade_by_node_label">Upgrade By Node Label</h3>
<div class="section">
<p>The <code>NodeLabel</code> strategy will perform a rolling upgrade by using a label on Nodes to group Nodes together.
Then all Pods on all the Nodes in a group (i.e. with the same label value) will be upgraded together.</p>

<p>In many production Kubernetes clusters, there is a concept of zones and fault domains, with each Node belonging to
one of these zones and domains. Typically, Nodes are labelled to indicate which zone and domain they are in.
For example the <code>topology.kubernetes.io/zone</code> is a standard label for the zone name.</p>

<p>These labels are used by the Coherence Operator to configure the site and rack names for a Coherence cluster.
(see the documentation on <router-link to="/docs/coherence/021_member_identity">Configuring Site and Rack</router-link>).
In a properly configured cluster that is site or rack safe, it is possible to upgrade all Pods in a site or rack
at the same time. In a typical Cloud Kubernetes Cluster there may be three zones, so a rolling upgrade by zon (or site)
would upgrade the cluster in three parts, which would be much faster than Pod by Pod.
This is a more extreme version of the <code>Node</code> strategy.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>When using the <code>NodeLabel</code> strategy where multiple Pods will be killed, the remaining cluster must have enough
capacity to recover the data and backups from the Pods that are killed.</p>

<p>For example, if a cluster of 18 Pods is distributed over Nodes in three zones, each zone will be running six Pods.
When upgrading by Node label, six Pods will be killed at the same time, so there must be enough capacity in the
remaining 12 Pods to hold all the data that was in the original 18 Pods.</p>
</p>
</div>
<p>The <code>Node</code> strategy is configured by setting the <code>rollingUpdateStrategy</code> field to <code>NodeLabel</code>
and also setting the <code>rollingUpdateLabel</code> field to the name of the label to use.</p>

<p>For example, to perform a rolling upgrade of all Pods by zone the yaml below could be used:</p>

<markup
lang="yaml"
title="cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  rollingUpdateStrategy: NodeLabel
  rollingUpdateLabel: "topology.kubernetes.io/zone"
  image: my-app:1.0.0</markup>

<div class="admonition caution">
<p class="admonition-textlabel">Caution</p>
<p ><p>It is up to the customer to verify that the label used is appropriate, i.e. is one of the labels used for the
Coherence site or rack configuration. It is also important to ensure that all Nodes in the cluster actually have
the label.</p>

<p>It is also up to the customer to verify that the Coherence cluster to be upgraded is site or rack safe before the
upgrade begins. The Coherence Operator can determine that no services are endangered, but it cannot determine site
or rack safety.</p>
</p>
</div>
</div>

<h3 id="_manual_upgrade">Manual Upgrade</h3>
<div class="section">
<p>If the <code>rollingUpdateStrategy</code> is set to <code>Manual</code> then neither the Coherence Operator, nor the StatefulSet controller in
Kubernetes will upgrade the Pods.
When the manual strategy is used the StatefulSet&#8217;s <code>spec.</code> field is set to <code>OnDelete</code>.
After updating a Coherence resource, the StatefulSet will be updated with the new state, but none of the Pods will be upgraded.
Pods must then be manually deleted so that they are rescheduled with the new configuration.
Pods can be deleted in any order and any number at a time.
(see <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#update-strategies">StatefulSet Update Strategies</a>
in the Kubernetes documentation).</p>

<p>The <code>Manual</code> strategy is configured by setting the <code>rollingUpdateStrategy</code> field to <code>Manual</code> as shown below:</p>

<markup
lang="yaml"
title="cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  rollingUpdateStrategy: Manual
  image: my-app:1.0.0</markup>

<div class="admonition caution">
<p class="admonition-textlabel">Caution</p>
<p ><p>When using the manual upgrade strategy, the customer is in full control of the upgrade process.
The Operator will not do anything. It is important that the customer understands how to perform
a safe rolling upgrade if no data loss is desired.</p>
</p>
</div>
</div>
</div>
</doc-view>
