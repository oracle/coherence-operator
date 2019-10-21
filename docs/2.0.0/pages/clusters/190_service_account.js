<doc-view>

<h2 id="_configure_the_kubernetes_service_account">Configure the Kubernetes Service Account</h2>
<div class="section">
<p>In Kubernetes clusters that have RBAC enabled it may be a requirement to set the
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/">service account</a>
that will be used by the <code>Pods</code> created for a <code>CoherenceCluster</code></p>

<p>The service account name is set for the <code>CoherenceCluster</code> as a whole and will be applied to all <code>Pods</code>.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  serviceAccountName: foo    <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">All <code>Pods</code> in the <code>test-cluster</code> will use the service account <code>foo</code></li>
</ul>
</div>
</doc-view>
