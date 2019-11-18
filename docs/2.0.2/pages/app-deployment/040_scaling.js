<doc-view>

<h2 id="_scaling_roles">Scaling Roles</h2>
<div class="section">
<p>For various reasons it is sometimes desirable to scale up or down the number of members of a Coherence cluster.
Whilst it has always been simple to add new members to a cluster (scale up) it needs care when removing members
(scale down) so that data is not lost. The Coherence Operator makes both of these operations simple by properly
managing safe scale down.</p>

<p>As already described, a cluster managed by the Coherence Operator is made up of roles, whether that is a single implicit
role or one or more explicit roles. The member count of a role is controlled by the role&#8217;s <code>replicas</code> field and the
Coherence Operator manages scaling at the role level by monitoring changes to the <code>replicas</code> field for a role.
Individual roles in a cluster can be scaled up or down independently (and without affecting) other roles in the cluster.</p>


<h3 id="_scale_up">Scale Up</h3>
<div class="section">
<p>By default operations that scale up a cluster will add members in parallel. Adding members to a role is a safe operation
and will not result in data loss so adding in parallel will scale up faster. This can be important if a cluster is under
heavy load and members need to be added quickly to keep the cluster healthy.</p>

</div>

<h3 id="_scale_down">Scale Down</h3>
<div class="section">
<p>By default when scaling down the Coherence Operator will remove members safely to ensure no data loss. Before removing
a member from a role the operator will check that the cluster is STatus HA (this is, no partitions are endangered) and
only then will a member be removed. In this way members are removed one at a time, which can be slow if scaling down by
a large number but the slowness is outweighed by the fact that there will be no data loss.</p>


<h4 id="_storage_disabled_roles">Storage Disabled Roles</h4>
<div class="section">
<p>When scaling down a storage disabled role the default will be to remove members in parallel. If a role is storage
disabled (i.e. the role&#8217;s <code>coherence.storageEnabled</code> field is set to <code>false</code>) then scaling down is parallel is safe as
those members are not managing data that might be lost.</p>

</div>
</div>

<h3 id="_scale_down_to_zero">Scale Down to Zero</h3>
<div class="section">
<p>Scaling down a role to have a replica count of zero is a special case which basically tells the Coherence Operator to
effectively un-deploy that role from the cluster. Scaling to zero will terminate all of the <code>Pods</code> of a role at the
same time. Obviously if the members of the role are storage enabled and persistence is not used then data will be lost.</p>

<p>Scaling down to zero is a way to remove all of the members of a role from a cluster without actually deleting the role
or cluster yaml from the Kubernetes cluster. This could be useful for example in cases where a role is used for a
one-off purpose such as data loading where it can run and then after completion be scaled back to zero until it is
required again when it can be scaled back up.</p>

</div>

<h3 id="_scaling_policy">Scaling Policy</h3>
<div class="section">
<p>Whether a role is scaled up or down in parallel or safely is controlled by the <code>scaling.policy</code> field of the role&#8217;s
spec. This is described in detail in the
<router-link to="/clusters/085_safe_scaling">Scaling section of the CoherenceCluster CRD documentation</router-link></p>

</div>
</div>

<h2 id="_the_mechanics_of_scaling">The Mechanics of Scaling</h2>
<div class="section">
<p>There are two ways to scale a role within a cluster; update the <code>replicas</code> field in the cluster <code>yaml</code>/<code>json</code> or
use the <code>kubectl scale</code> command to scale a specific role.</p>

<p>For example:</p>

<p>A Coherence cluster can be defined in a <code>.yaml</code> file called <code>test-cluster.yaml</code> like this:</p>

<markup
lang="yaml"
title="test-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: storage
      replicas: 6
    - role: http
      replicas: 2</markup>

<p>The cluster can be created in Kubernetes using <code>kubectl</code>:</p>

<markup


>kubectl create -f test-cluster.yaml</markup>

<p>After the cluster has started the roles in the cluster can be listed:</p>

<markup


>kubectl get coherenceroles</markup>

<p>&#8230;&#8203;which might display something like the following:</p>

<markup


>NAME                   ROLE      CLUSTER        REPLICAS   READY   STATUS
test-cluster-http      http      test-cluster   2          2       Ready
test-cluster-storage   storage   test-cluster   6          6       Ready</markup>

<p>As defined in the <code>yaml</code> the <code>test-cluster</code> has two roles <code>test-cluster-http</code> and <code>test-cluster-storage</code>.</p>


<h3 id="_update_the_coherencecluster_yaml">Update the CoherenceCluster YAML</h3>
<div class="section">
<p>To scale up the <code>storage</code> role to nine members one option would be to update the yaml:</p>

<markup
lang="yaml"
title="test-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: storage
      replicas: 9        <span class="conum" data-value="1" />
    - role: http
      replicas: 2</markup>

<p>The <code>storage</code> role now has <code>replicas</code> set to <code>9</code> so re-apply the <code>yaml</code> Kubernetes using <code>kubectl</code>:</p>

<markup


>kubectl apply -f test-cluster.yaml</markup>

<p>&#8230;&#8203;after the new <code>Pods</code> have started listing the roles might look like this:</p>

<markup


>kubectl get coherenceroles</markup>

<markup


>NAME                   ROLE      CLUSTER        REPLICAS   READY   STATUS
test-cluster-http      http      test-cluster   2          2       Ready
test-cluster-storage   storage   test-cluster   9          9       Ready</markup>

</div>

<h3 id="_use_the_kubectl_scale_command">Use the kubectl scale Command</h3>
<div class="section">
<p>The <code>kubectl</code> CLI offers a simple way to scale a Kubernetes resource providing that the resource is properly configured
to allow this (which the Coherence CRDs are).</p>

<p>Continuing the example if the <code>storage</code> role is to now be scale down from nine back to six then <code>kubectl</code> can be used as follows:</p>

<markup


>kubectl scale coherencerole test-cluster-storage --replicas=6</markup>

<p>The Coherence Operator will now scale the <code>storage</code> role down by removing one member at a time until the desired replica
count is reached. Eventually listing the roles will show the desired state:</p>

<markup


>kubectl get coherenceroles</markup>

<markup


>NAME                   ROLE      CLUSTER        REPLICAS   READY   STATUS
test-cluster-http      http      test-cluster   2          2       Ready
test-cluster-storage   storage   test-cluster   6          6       Ready</markup>

</div>
</div>
</doc-view>
