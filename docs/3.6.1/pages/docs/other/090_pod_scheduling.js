<doc-view>

<h2 id="_configure_pod_scheduling">Configure Pod Scheduling</h2>
<div class="section">
<p>In Kubernetes <code>Pods</code> can be configured to control how, and onto which nodes, Kubernetes will schedule those <code>Pods</code>; the
Coherence Operator allows the same control for <code>Pods</code> owned by a <code>Coherence</code> resource.</p>

<p>The following settings can be configured:</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 50%;">
<col style="width: 50%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>nodeSelector</code></td>
<td class=""><code>nodeSelector</code> is the simplest recommended form of node selection constraint.
<code>nodeSelector</code> is a field of role spec, it specifies a map of key-value pairs.
For the <code>Pod</code> to be eligible to run on a node, the node must have each of the indicated key-value pairs as labels
(it can have additional labels as well).
See <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/configuration/assign-pod-node/">Assigning Pods to Nodes</a> in the
Kubernetes documentation</td>
</tr>
<tr>
<td class=""><code>affinity</code></td>
<td class="">The affinity/anti-affinity feature, greatly expands the types of constraints you can express over just using labels
in a <code>nodeSelector</code>.
See <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/configuration/assign-pod-node/">Assigning Pods to Nodes</a> in the
Kubernetes documentation</td>
</tr>
<tr>
<td class=""><code>tolerations</code></td>
<td class=""><code>nodeSelector</code> and <code>affinity</code> are properties of <code>Pods</code> that attracts them to a set of nodes (either as a preference or
a hard requirement). Taints are the opposite – they allow a node to repel a set of <code>Pods</code>.
Taints and tolerations work together to ensure that <code>Pods</code> are not scheduled onto inappropriate nodes.
One or more taints are applied to a node; this marks that the node should not accept any <code>Pods</code> that do not tolerate
the taints. Tolerations are applied to <code>Pods</code>, and allow (but do not require) the <code>Pods</code> to schedule onto nodes with
matching taints.
See <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/">Taints and Tolerations</a> in the Kubernetes
documentation.</td>
</tr>
</tbody>
</table>
</div>
<p>The <code>nodeSelector</code>, <code>affinity</code> and <code>tolerations</code> fields are all part of the <code>Coherence</code> CRD spec.
The format of the fields is that same as documented in the Kubernetes documentation
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/configuration/assign-pod-node/">Assigning Pods to Nodes</a> and
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/">Taints and Tolerations</a></p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test-cluster
spec:
  tolerations:
    - key: "example-key"
      operator: "Exists"
      effect: "NoSchedule"
  nodeSelector:
    disktype: ssd
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
          - matchExpressions:
             - key: kubernetes.io/e2e-az-name
               operator: In
               values:
                 - e2e-az1
                 - e2e-az2</markup>

</div>
</doc-view>
