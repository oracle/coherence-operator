<doc-view>

<h2 id="_useful_info">Useful Info</h2>
<div class="section">

<h3 id="_labeling_your_k8s_node">Labeling Your K8s Node</h3>
<div class="section">
<p>For local testing, for example in Docker Desktop it is useful to add the zone label to your local K8s node with
the fault domain that is then used by the Coherence Pods to set their <code>site</code> property.</p>

<p>For example, if your local node is called <code>docker-desktop</code> you can use the following command to set
the zone name to <code>twilight-zone</code>:</p>

<markup
lang="bash"

>kubectl label node docker-desktop failure-domain.beta.kubernetes.io/zone=twilight-zone</markup>

<p>With this label set all Coherence Pods installed by the Coherence Operator on that node will be
running in the <code>twilight-zone</code>.</p>

</div>

<h3 id="_kubernetes_dashboard">Kubernetes Dashboard</h3>
<div class="section">
<p>Assuming that you have the <a id="" title="" target="_blank" href="https://github.com/kubernetes/dashboard">Kubernetes Dashboard</a> then you can easily
start the local proxy and display the required login token by running:</p>

<markup
lang="bash"

>./hack/kube-dash.sh</markup>

<p>This will display the authentication token, the local k8s dashboard URL and then start <code>kubectl proxy</code>.</p>

</div>

<h3 id="_stuck_coherenceinternal_resources">Stuck CoherenceInternal Resources</h3>
<div class="section">
<p>Sometimes a CoherenceInternal resource becomes stuck in k8s. This is because the operator adds finalizers to the
resources causing k8s to be unable to delete them. The simplest way to delete them is to use the <code>kubectl patch</code>
command to remove the finalizer.</p>

<p>For example, if there was a CoherenceInternal resource called <code>test-role</code> in namespace <code>testing</code> then
the following command could be used.</p>

<markup
lang="bash"

>kubectl -n testing patch coherenceinternal/test-role \
  -p '{"metadata":{"finalizers": []}}' \
  --type=merge;</markup>

<p>Alternatively there is a make target that wil clean up and remove all CoherenceCLuster, CoherenceRole and CoherenceInternal
resources from the test namespace.</p>

<markup
lang="bash"

>make delete-coherence-clusters</markup>

</div>
</div>
</doc-view>
