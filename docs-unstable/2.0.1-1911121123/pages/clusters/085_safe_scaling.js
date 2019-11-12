<doc-view>

<h2 id="_configure_safe_scaling">Configure Safe Scaling</h2>
<div class="section">
<p>The Coherence Operator contains functionality to allow it to safely scale a role within a Coherence cluster without
losing data. Scaling can be configured in the <code>scaling</code> section of the <code>CoherenceCluster</code> CRD.</p>

<p>A role in a <code>CoherenceCluster</code> can be scaled by changing the replica count in the role&#8217;s spec or by using the
<code>kubectl scale</code> command.</p>


<h3 id="_scaling_policy">Scaling Policy</h3>
<div class="section">
<p>The Coherence Operator uses a scaling policy to determine how the <code>StatefulSet</code> that makes up a role withing a
cluster is scaled.
Scaling policy has the following values:</p>


<div class="table__overflow elevation-1 ">
<table class="datatable table">
<colgroup>
<col style="width: 50%;">
<col style="width: 50%;">
</colgroup>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td><code>ParallelUpSafeDown</code></td>
<td>This is the default scaling policy.
With this policy when scaling up <code>Pods</code> are added in parallel (the same as using the <code>Parallel</code> <code>podManagementPolicy</code>
in a <a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.16/#statefulsetspec-v1-apps">StatefulSet</a>) and
when scaling down <code>Pods</code> are removed one at a time (the same as the <code>OrderedReady</code> <code>podManagementPolicy</code> for a
StatefulSet). When scaling down a check is done to ensure that the members of the role have a safe StatusHA value
before a <code>Pod</code> is removed (i.e. none of the Coherence cache services have an endangered status).
This policy offers faster scaling up and start-up because pods are added in parallel as data should not be lost when
adding members, but offers safe, albeit slower,  scaling down as <code>Pods</code> are removed one by one.</td>
</tr>
<tr>
<td><code>Parallel</code></td>
<td>With this policy when scaling up <code>Pods</code> are added in parallel (the same as using the <code>Parallel</code> <code>podManagementPolicy</code>
in a <a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.16/#statefulsetspec-v1-apps">StatefulSet</a>).
With this policy no StatusHA check is performed either when scaling up or when scaling down.
This policy allows faster start and scaling times but at the cost of no data safety; it is ideal for roles that are
storage disabled.</td>
</tr>
<tr>
<td><code>Safe</code></td>
<td>With this policy when scaling up and down <code>Pods</code> are removed one at a time (the same as the <code>OrderedReady</code>
<code>podManagementPolicy</code> for a StatefulSet). When scaling down a check is done to ensure that the members of the role have
a safe StatusHA value before a <code>Pod</code> is removed (i.e. none of the Coherence cache services have an endangered status).
This policy is slow to start, scale up and scale down.</td><td>&#8230;&#8203;</td>
</tr>
</tbody>
</table>
</div>
<p>The scaling policy is set in the <code>scaling.policy</code> section of the configuration of a role.</p>

</div>

<h3 id="_configure_scaling_policy_for_a_single_implicit_role">Configure Scaling Policy for a Single Implicit Role</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with a single implicit role the scaling policy can be defined at the <code>spec</code> level.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  scaling:
    policy: Safe <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The implicit role will have a scaling policy of <code>Safe</code></li>
</ul>
</div>

<h3 id="_configure_scaling_policy_for_explicit_roles">Configure Scaling Policy for Explicit Roles</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with explicit roles in the <code>roles</code> list scaling policy can be defined for each role,
for example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data
      scaling:
        policy: ParallelUpSafeDown <span class="conum" data-value="1" />
    - role: proxy
      scaling:
        policy: Parallel           <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role will have the scaling policy <code>ParallelUpSafeDown</code></li>
<li data-value="2">The <code>proxy</code> role will have the scaling policy <code>Parallel</code></li>
</ul>
</div>

<h3 id="_configure_pod_labels_for_explicit_roles_with_defaults">Configure Pod Labels for Explicit Roles With Defaults</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with explicit roles in the <code>roles</code> list scaling policy can be defined as defaults
applied to all roles unless specifically overridden for a role.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  scaling:
    policy: Parallel               <span class="conum" data-value="1" />
  roles:
    - role: data
      scaling:
        policy: ParallelUpSafeDown <span class="conum" data-value="2" />
    - role: proxy                  <span class="conum" data-value="3" />
    - role: web                    <span class="conum" data-value="4" /></markup>

<ul class="colist">
<li data-value="1">Thedefault scaling policy is <code>Parallel</code> that will apply to all roles unless specifically overridden.</li>
<li data-value="2">The <code>data</code> role overrides the default and specifies a scaling policy of <code>ParallelUpSafeDown</code></li>
<li data-value="3">The <code>proxy</code> role does not specify a scaling policy so will use the defautl of <code>Parallel</code></li>
<li data-value="4">The <code>web</code> role does not specify a scaling policy so will use the defautl of <code>Parallel</code></li>
</ul>
</div>
</div>
</doc-view>
