<doc-view>

<h2 id="_member_identity">Member Identity</h2>
<div class="section">
<p>Each JVM in a Coherence cluster has an identity. This is made up of a number of values for <code>site</code>, <code>rack</code>, <code>member</code>,
<code>machine</code> and <code>node-id</code>.
The <code>node-id</code> is assigned by Coherence when a node joins a cluster.
The other values can be assigned using system properties, or will have defaults assigned by Coherence if not set.
The Coherence Operator will configure properties for these values.</p>

<ul class="ulist">
<li>
<p>The member name is set to the Pod name.</p>

</li>
<li>
<p>The machine name is set to the name of the Node that the Pod has been scheduled onto.</p>

</li>
<li>
<p>The rack name is taken from the <code>topology.kubernetes.io/region</code> label on the Node that the Pod has been scheduled onto.
If the <code>topology.kubernetes.io/region</code> label is not set then the deprecated <code>failure-domain.beta.kubernetes.io/region</code>
label will be tried.
If neither of these labels are set then the site will be set to the same value as the site name.</p>

</li>
<li>
<p>The site name is taken from the <code>topology.kubernetes.io/zone</code> label on the Node that the Pod has been scheduled onto.
If the <code>topology.kubernetes.io/zone</code> label is not set then the deprecated <code>failure-domain.beta.kubernetes.io/zone</code> label
will be tried.
If neither of these labels are set then the site and rack will be unset, and the cache services may not reach site safe.</p>

</li>
</ul>
</div>

<h2 id="_status_ha_values">Status HA Values</h2>
<div class="section">
<p>As well as identifying cluster members, these values are also used by the partitioned cache service to distribute data
as widely (safely) as possible in the cluster. The backup owner will be as far away as possible from the primary owner.
Ideally this would be on a member with a different site; failing that, a different rack, machine and finally member.</p>

</div>

<h2 id="_changing_site_and_rack_values">Changing Site and Rack Values</h2>
<div class="section">
<p>You should not usually need to change the default values applied for the <code>member</code> and <code>machine</code> names, but you may need
to change the values used for the site, or rack. The labels used for the <code>site</code> and <code>rack</code> are standard k8s labels but
the k8s cluster being used may not have these labels set</p>


<h3 id="_apply_node_labels">Apply Node Labels</h3>
<div class="section">
<p>One solution to missing site and rack values is to apply the required labels to the Nodes in the k8s cluster.</p>

<p>For example the command below labels the node in Docker dDesktop on MacOS to "twighlight-zone".</p>

<markup
lang="bash"

>kubectl label node docker-desktop topology.kubernetes.io/zone=twighlight-zone</markup>

</div>

<h3 id="_specify_site_and_rack_using_system_properties">Specify Site and Rack Using System Properties</h3>
<div class="section">
<p>The site and rack values can be specified as system properties as part of the Coherence deployment yaml.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: my-cluster
spec:
  jvm:
    args:
      - "-Dcoherence.site=foo"
      - "-Dcoherence.rack=fbar"</markup>

<p>In the deployment above the site name is set to "foo" using the <code>coherence.site</code> system property.
The rack name is set to "bar" using the <code>coherence.rack</code> system property.</p>

</div>

<h3 id="_configure_the_operator_to_use_different_labels">Configure the Operator to Use Different Labels</h3>
<div class="section">
<p>The Operator can be configured to use different labels to obtain values for the site and rack names.
This will obviously apply to all Coherence deployments managed by the Operator, but is useful if the Nodes in the
k8s cluster do not have the normal k8s labels.
The <code>SITE_LABEL</code> and <code>RACK_LABEL</code> environment variables are used to specify different labels to use.
How these environment variables are set depends on how you are installing the Operator.</p>


<h4 id="_using_helm">Using Helm</h4>
<div class="section">
<p>If the Operator is installed using the Helm chart then the site and rack labels can be set using the
<code>siteLabel</code> and <code>rackLabel</code> values;
for example:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set siteLabel=identity/site \
    --set siteLabel=identity/rack \
    coherence-operator \
    coherence/coherence-operator</markup>

<p>In the example above the Node label used by the Operator to get the value for the site will be <code>identity/site</code>,
and the Node label used to get the value for the rack will be <code>identity/rack</code>.</p>

</div>

<h4 id="_using_kubectl_or_kustomize">Using Kubectl or Kustomize</h4>
<div class="section">
<p>If using <code>kubectl</code> or <code>kustomize</code> as described in the <router-link to="/docs/installation/01_installation">Installation Guide</router-link>
the additional environment variables can be applied using <code>kustomize</code> commands.</p>

<markup
lang="bash"

>cd ./manager &amp;&amp; $(GOBIN)/kustomize edit add configmap env-vars --from-literal SITE_LABEL='identity/site'</markup>

<markup
lang="bash"

>cd ./manager &amp;&amp; $(GOBIN)/kustomize edit add configmap env-vars --from-literal RACK_LABEL='identity/rack'</markup>

</div>
</div>
</div>
</doc-view>
