<doc-view>

<h2 id="_scale_coherence_deployments">Scale Coherence Deployments</h2>
<div class="section">
<p>The Coherence Operator provides the ability to safely scale up and down a <code>Coherence</code> deployment.
A <code>Coherence</code> deployment is backed by a <code>StatefulSet</code>, which can easily be scaled using existing Kubernetes features.
The problem with directly scaling down the <code>StatefulSet</code> is that Kubernetes will immediately kill the required number
of <code>Pods</code>. This is obviously very bad for Coherence as killing multiple storage enabled members would almost certainly
cause data loss.</p>

<p>The Coherence Operator supports scaling by applying the scaling update directly to <code>Coherence</code> deployment rather than
to the underlying <code>StatefulSet</code>. There are two methods to scale a <code>Coherence</code> deployment:</p>

<ul class="ulist">
<li>
<p>Update the <code>replicas</code> field in the <code>Coherence</code> CRD spec.</p>

</li>
<li>
<p>Use the <code>kubectl scale</code> command</p>

</li>
</ul>
<p>When either of these methods is used the Operator will detect that a change to the size of the deployment is required
and ensure that the change will be applied safely. The logical steps the Operator will perform are:</p>

<ol style="margin-left: 15px;">
<li>
Detect desired replicas is different to current replicas

</li>
<li>
Check the cluster is StatusHA - i.e. no cache services are endangered. If any service is not StatusHA requeue the
scale request  (go back to step one).

</li>
<li>
If scaling up, add the required number of members.

</li>
<li>
If scaling down, scale down by one member and requeue the request (go back to step one).

</li>
</ol>
<p>What these steps ensure is that the deployment will not be resized unless the cluster is in a safe state.
When scaling down only a single member will be removed at a time, ensuring that the cluster is in a safe state before
removing the next member.</p>

<div class="admonition note">
<p class="admonition-inline">The Operator will only apply safe scaling functionality to deployments that are storage enabled.
If a deployment is storage disabled then it can be scaled up or down by the required number of members
in one step as there is no fear of data loss in a storage disabled member.</p>
</div>
</div>

<h2 id="_controlling_safe_scaling">Controlling Safe Scaling</h2>
<div class="section">
<p>The <code>Coherence</code> CRD has a number of fields that control the behaviour of scaling.</p>


<h3 id="_scaling_mode">Scaling Mode</h3>
<div class="section">

</div>

<h3 id="_scaling_status_ha_probe">Scaling Status HA Probe</h3>
<div class="section">

</div>
</div>
</doc-view>
