<doc-view>

<h2 id="_pod_labels">Pod Labels</h2>
<div class="section">
<p>Additional labels can be added to the <code>Pods</code> managed by the Operator.
Additional labels should be added to the <code>labels</code> map in the <code>Coherence</code> CRD spec.
The entries in the <code>labels</code> map should confirm to the recommendations and rules in the Kubernetes
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/">Labels</a> documentation.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  labels:             <span class="conum" data-value="1" />
    tier: backend
    environment: dev</markup>

<ul class="colist">
<li data-value="1">Two labels will be added to the <code>Pods</code>, <code>tier=backend</code> and <code>environment=dev</code></li>
</ul>
</div>
</doc-view>
