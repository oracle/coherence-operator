<doc-view>

<h2 id="_coherence_clusters_on_openshift">Coherence Clusters on OpenShift</h2>
<div class="section">
<p>Whilst the Coherence Operator will run out of the box on OpenShift some earlier versions of the Coherence Docker
image will not work without configuration changes.</p>

<p>These earlier versions of the Coherence Docker images that Oracle publishes default the container user
as <code>oracle</code>. When running the Oracle images or layered images that retain the default user as <code>oracle</code>
with OpenShift, the <code>anyuid</code> security context constraint is required to ensure proper access to the file
system within the Docker image. Later versions of the Coherence images have been modified to work without
needing <code>anyuid</code>.</p>

<p>To work with older image versions , the administrator must:</p>

<ul class="ulist">
<li>
<p>Ensure the <code>anyuid</code> security content is granted</p>

</li>
<li>
<p>Ensure that Coherence containers are annotated with <code>openshift.io/scc: anyuid</code></p>

</li>
</ul>
<p>For example, to update the OpenShift policy, use:</p>

<markup
lang="bash"

>oc adm policy add-scc-to-user anyuid -z default</markup>

<p>and to annotate the Coherence containers, update the <code>Coherence</code> resource to include annotations</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test-cluster
spec:
  annotations:
    openshift.io/scc: anyuid  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>openshift.io/scc: anyuid</code> annotation will be applied to all of the Coherence Pods.</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">For additional information about OpenShift requirements see the
<a id="" title="" target="_blank" href="https://docs.openshift.com/container-platform/3.3/creating_images/guidelines.html">OpenShift documentation</a></p>
</div>
</div>
</doc-view>
