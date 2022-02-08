<doc-view>

<h2 id="_troubleshooting_guide">Troubleshooting Guide</h2>
<div class="section">
<p>The purpose of this page is to list troubleshooting guides and work-arounds for issues that you may run into when using the Coherence Operator.
This page will be updated and maintained over time to include common issues we see from customers</p>

</div>

<h2 id="_contents">Contents</h2>
<div class="section">
<ul class="ulist">
<li>
<p><router-link to="#restart" @click.native="this.scrollFix('#restart')">Why Does the Operator Pod Restart</router-link></p>

</li>
<li>
<p><router-link to="#ready" @click.native="this.scrollFix('#ready')">Why are the Coherence Pods not reaching ready</router-link></p>

</li>
<li>
<p><router-link to="#stuck-pending" @click.native="this.scrollFix('#stuck-pending')">My Coherence cluster is stuck with some running Pods and some pending Pods, I want to scale down</router-link></p>

</li>
<li>
<p><router-link to="#stuck-delete" @click.native="this.scrollFix('#stuck-delete')">My Coherence cluster is stuck with all pending/crashing Pods, I cannot delete the deployment</router-link></p>

</li>
<li>
<p><router-link to="#site-safe" @click.native="this.scrollFix('#site-safe')">My cache services will not reach Site Safe</router-link></p>

</li>
<li>
<p><router-link to="#dashboards" @click.native="this.scrollFix('#dashboards')">My Grafana Dashboards do not display any metrics</router-link></p>

</li>
</ul>
</div>

<h2 id="_issues">Issues</h2>
<div class="section">

<h3 id="restart">Why Does the Operator Pod Restart After Installation</h3>
<div class="section">
<p>You might notice that when the Operator is installed that the Operator Pod starts, dies and is then restarted.
This is expected behaviour. The Operator uses a K8s web-hook for defaulting and validation of the Coherence resources.
A web-hook requires certificates to be present in a <code>Secret</code> mounted to the Operator Pod as a volume.
If this Secret is not present the Operator creates it and populates it with self-signed certs.
In order for the Secret to then be mounted correctly the Pod must be restarted.
See the <router-link to="/docs/webhooks/01_introduction">the WebHooks</router-link> documentation.</p>

</div>

<h3 id="ready">Why are the Coherence Pods not reaching ready</h3>
<div class="section">
<p>The readiness/liveness probe used by the Operator in the Coherence Pods checks a number of things to determine whether the Pods is ready, one of these is whether the JVM is a cluster member.
If your application uses a custom main class and is not properly bootstrapping Coherence then the Pod will not be ready until your application code actually touches a Coherence resource causing Coherence to start and join the cluster.</p>

<p>When running in clusters with the Operator using custom main classes it is advisable to properly bootstrap Coherence
from within your <code>main</code> method. This can be done using the new Coherence bootstrap API available from CE release 20.12
or by calling <code>com.tangosol.net.DefaultCacheServer.startServerDaemon().waitForServiceStart();</code></p>

</div>

<h3 id="stuck-pending">My Coherence cluster is stuck with some running Pods and some pending Pods, I want to scale down</h3>
<div class="section">
<p>If you try to create a Coherence deployment that has a replica count that is greater than your k8s cluster can actually
provision then one or more Pods will fail to be created or can be left in a pending state.
The obvious solution to this is to just scale down your Coherence deployment to a smaller size that can be provisioned.
The issue here is that the safe scaling functionality built into the operator will not allow scaling down to take place
because it cannot guarantee no parition/data loss. The Coherence deployment is now stuck in this state.</p>

<p>The simplest solution would be to completely delete the the Coherence deployment and redeploy with a lower replica count.</p>

<p>If this is not possible then the following steps will allow the deployment to be scaled down.</p>

<p>1 Update the stuck Coherence deployment&#8217;s scaling policy to be <code>Parallel</code></p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  scaling:
    policy: Parallel</markup>

<p>2 Scale down the cluster to the required size using whatever scaling commands you want, i.e <code>kubectl scale</code>
or just update the replica value of the Coherence deployment yaml. Note: If updating the Coherence yaml, this
should not be done as part of step 1, above.</p>

<p>3 Once the Coherence deployment has scaled to the required size then change the scaling policy value back to the
default by updating the Coherence yaml to have no scaling policy value in it.</p>

<div class="admonition warning">
<p class="admonition-inline">When using this work around to scale down a stuck deployment that contains data it is important that
only the missing or pending Pods are removed. For example if a Coherence deployment is deployed with a replica count
of 100 and 90 Pods are ready, but the other 10 are either missing or stuck pending then the replica value used in
step 2 above must be 90. Because the scaling policy has been set to <code>Parallel</code> the operator will not check any
Status HA values before scaling down Pods, so removing "ready" Pods that contain data will almost certainly result
in data loss. To safely scale down lower, then first follow the three steps above then after changing the scaling policy
back to the default further scaling down can be done as normal.</p>
</div>
</div>

<h3 id="stuck-delete">My Coherence cluster is stuck with all pending/crashing Pods, I cannot delete the deployment</h3>
<div class="section">
<p>A Coherence deployment can become stuck where none of the Pods can start, for example the image used is incorrect
and all Pods are stuck in ImagePullBackoff. It can then become impossible to delete the broken deployment.
This is because the Operator has installed a finalizer but this finalizer cannot execute.</p>

<p>For example, suppose we have deployed a Coherence deployment named <code>my-cluster</code> into namespace <code>coherence-test</code>.</p>

<p>First try to delete the deployment as normal:</p>

<markup
lang="console"

>kubectl -n coherence-test delete coherence/my-cluster</markup>

<p>If this command hangs, then press <code>ctrl-c</code> to exit and then run the following patch command.</p>

<markup
lang="console"

>kubectl -n coherence-test patch coherence/my-cluster -p '{"metadata":{"finalizers":[]}}' --type=merge</markup>

<p>This will remove the Operator&#8217;s finalizer from the Coherence deployment.</p>

<p>At this point the <code>my-cluster</code> Coherence deployment might already have been removed,
if not try the delete command again.</p>

</div>

<h3 id="site-safe">My cache services will not reach Site-Safe</h3>
<div class="section">
<p>Coherence distributes data in a cluster to achieve the highest status HA value that it can, the best being site-safe.
This is done using the various values configured for the site, rack, machine, and member names.
The Coherence Operator configures these values for the Pods in a Coherence deployment.
By default, the values for the site and rack names are taken from standard k8s labels applied to the Nodes in the k8s cluster.
If the Nodes in the cluster do not have these labels set then the site and rack names will be unset and Coherence
will not be able to reach rack or site safe.</p>

<p>There are a few possible solutions to this, see the explanation in the
documentation explaining <router-link to="/docs/coherence/021_member_identity">Member Identity</router-link></p>

</div>

<h3 id="dashboards">My Grafana Dashboards do not display any metrics</h3>
<div class="section">
<p>If you have imported the Grafana dashboards provided by the Operator into Grafana, but they are not displaying any metric
values, it may be that you have imported the wrong format dashboards. The Operator has multiple sets of dashboards,
one for the default Coherence metric name format, one for Microprofile metric name format, and one for
<a id="" title="" target="_blank" href="https://micrometer.io">Micrometer</a> metric name format.</p>

<p>The simplest way to find out which version corresponds to your Coherence cluster
is to query the metrics endpoint with something like <code>curl</code>.
If the metric names are in the format <code>vendor:coherence_cluster_size</code>, i.e. prefixed with <code>vendor:</code> then this is
the default Coherence format.
If metric names are in the format <code>vendor_Coherence_Cluster_Size</code>, i.e. prefixed with <code>vendor_</code> then this is
Microprofile format.
If the metric name has no <code>vendor</code> prefix then it is using Micrometer metrics.</p>

<p>See: the <router-link to="/docs/metrics/030_importing">Importing Grafana Dashboards</router-link> documentation.</p>

</div>
</div>
</doc-view>
