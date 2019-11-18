<doc-view>

<h2 id="_configure_pod_scheduling">Configure Pod Scheduling</h2>
<div class="section">
<p>In Kubernetes <code>Pods</code> can be configured to control how and onto which nodes Kubernetes will schedule those <code>Pods</code>; the
Coherence Operator allows the same control for <code>Pods</code> in roles in a <code>CoherenceCluster</code> resource.</p>

<p>The following settings can be configured:</p>


<div class="table__overflow elevation-1 ">
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
<td><code>nodeSelector</code></td>
<td><code>nodeSelector</code> is the simplest recommended form of node selection constraint.
<code>nodeSelector</code> is a field of role spec, it specifies a map of key-value pairs.
For the <code>Pod</code> to be eligible to run on a node, the node must have each of the indicated key-value pairs as labels
(it can have additional labels as well).
See <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/configuration/assign-pod-node/">Assigning Pods to Nodes</a> in the
Kubernetes documentation</td>
</tr>
<tr>
<td><code>affinity</code></td>
<td>The affinity/anti-affinity feature, greatly expands the types of constraints you can express over just using labels
in a <code>nodeSelector</code>.
See <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/configuration/assign-pod-node/">Assigning Pods to Nodes</a> in the
Kubernetes documentation</td>
</tr>
<tr>
<td><code>tolerations</code></td>
<td><code>nodeSelector</code> and <code>affinity</code> are properties of <code>Pods</code> that attracts them to a set of nodes (either as a preference or
a hard requirement). Taints are the opposite – they allow a node to repel a set of <code>Pods</code>.
Taints and tolerations work together to ensure that <code>Pods</code> are not scheduled onto inappropriate nodes.
One or more taints are applied to a node; this marks that the node should not accept any <code>Pods</code> that do not tolerate
the taints. Tolerations are applied to <code>Pods</code>, and allow (but do not require) the <code>Pods</code> to schedule onto nodes with
matching taints.
See <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/">Taints and Tolerations</a> in the Kubernetes
documentation.</td><td>&#8230;&#8203;</td>
</tr>
</tbody>
</table>
</div>
<p>The <code>nodeSelector</code>, <code>affinity</code> and <code>tolerations</code> fields are all part of the role spec and like any other role spec
field can be configured at different levels depending on whether the <code>CoherenceCluster</code> has implicit or explicit roles.
The format of the fields is that same as documented in the Kubernetes documentation
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/configuration/assign-pod-node/">Assigning Pods to Nodes</a> and
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/">Taints and Tolerations</a></p>

</div>

<h2 id="_pod_scheduling_for_a_single_implicit_role">Pod Scheduling for a Single Implicit Role</h2>
<div class="section">
<p>When configuring a <code>CoherenceCluster</code> with a single implicit role the scheduling fields are configured directly in
the <code>CoherenceCluster</code> <code>spec</code> section.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  tolerations:                                             <span class="conum" data-value="1" />
   - key: "example-key"
     operator: "Exists"
     effect: "NoSchedule"
   nodeSelector:                                           <span class="conum" data-value="2" />
   - disktype: ssd
   affinity:                                               <span class="conum" data-value="3" />
     nodeAffinity:
       requiredDuringSchedulingIgnoredDuringExecution:
         nodeSelectorTerms:
         - matchExpressions:
           - key: kubernetes.io/e2e-az-name
             operator: In
             values:
             - e2e-az1
             - e2e-az2</markup>

<ul class="colist">
<li data-value="1">The <code>tolerations</code> are set for the implicit <code>storage</code> role</li>
<li data-value="2">A <code>nodeSelector</code> is set for the implicit <code>storage</code> role</li>
<li data-value="3"><code>affinity</code> is set for the implicit <code>storage</code> role</li>
</ul>
</div>

<h2 id="_pod_scheduling_for_explicit_roles">Pod Scheduling for Explicit Roles</h2>
<div class="section">
<p>When configuring one or more explicit roles in a <code>CoherenceCluster</code> the scheduling fields are configured for each role
in the <code>roles</code> list.
For example</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data
      nodeSelector:                                        <span class="conum" data-value="1" />
      - disktype: ssd
    - role: proxy
      tolerations:                                         <span class="conum" data-value="2" />
       - key: "example-key"
         operator: "Exists"
         effect: "NoSchedule"
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

<ul class="colist">
<li data-value="1">The <code>data</code> role has a <code>nodeSelector</code> configured</li>
<li data-value="2">The <code>proxy</code> role has <code>tolerations</code> and <code>affinity</code> configured</li>
</ul>
</div>

<h2 id="_pod_scheduling_for_explicit_roles_with_defaults">Pod Scheduling for Explicit Roles with Defaults</h2>
<div class="section">
<p>When configuring one or more explicit roles in a <code>CoherenceCluster</code> default values for the scheduling fields may be
configured directly in the <code>spec</code> section of the <code>CoherenceCluster</code> that will apply to all roles in the <code>roles</code> list
unless specifically overridden for a role.
Values specified for a role fully override the default values, so even though <code>nodeSelector</code> is a map the default and
role values are <strong>not</strong> merged.</p>

<p>For example</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  nodeSelector:                                        <span class="conum" data-value="1" />
    - disktype: ssd
  roles:
    - role: data
      nodeSelector:                                    <span class="conum" data-value="2" />
        - shape: massive
    - role: proxy                                      <span class="conum" data-value="3" />
    - role: web</markup>

<ul class="colist">
<li data-value="1">The default scheduling configuration specified a node selector label of <code>disktype=ssd</code></li>
<li data-value="2">The <code>data</code> role overrides the <code>nodeSelector</code> to be <code>shape=massive</code></li>
<li data-value="3">The <code>proxy</code> and <code>web</code> roles do not specify any scheduling fields so they will just ue the default node selector
label of <code>disktype=ssd</code></li>
</ul>
<p>The <code>tolerations</code> and <code>affinity</code> fields may be used in the same way.</p>

</div>
</doc-view>
