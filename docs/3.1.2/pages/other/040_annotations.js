<doc-view>

<h2 id="_pod_annotations">Pod Annotations</h2>
<div class="section">
<p>Additional annotations can be added to the <code>Pods</code> managed by the Operator.
Annotations should be added to the <code>annotations</code> map in the <code>Coherence</code> CRD spec.
The entries in the <code>annotations</code> map should confirm to the recommendations and rules in the Kubernetes
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/">Annotations</a> documentation.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  annotations:                        <span class="conum" data-value="1" />
    prometheus.io/path: /metrics
    prometheus.io/port: "9612"
    prometheus.io/scheme: http
    prometheus.io/scrape: "true"</markup>

<ul class="colist">
<li data-value="1">A number of Prometheus annotations will be added to this <code>Coherence</code> deployment&#8217;s <code>Pods</code></li>
</ul>
</div>
</doc-view>
