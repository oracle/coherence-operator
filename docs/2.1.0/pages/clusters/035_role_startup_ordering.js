<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Define Coherence Role Start Order</dt>
<dd slot="desc"><p>The roles in a <code>CoherenceCluster</code> can be configured to start in a specific order.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_define_coherence_role_start_order">Define Coherence Role Start Order</h2>
<div class="section">
<p>The default behaviour of the operator is to create the <code>StatefulSets</code> for all of the roles in parallel so that they all start at
the same time. Sometimes this behaviour is not suitable if, for example, application code running in one role depends on the
availability of another role. The <code>CoherenceCluster</code> CRD allows roles to be configured with a <code>startQuorum</code> that defines a role&#8217;s
dependency on other roles in the cluster.</p>

<div class="admonition note">
<p class="admonition-inline">The <code>startQuorum</code> only applies when a cluster is being created by the operator, it does not apply in other functions such as
upgrades, scaling, shut down etc.</p>
</div>
<p>An individual role can depend on one or more other roles. The dependency can be such that the role will not be created until all
of the <code>Pods</code> of the dependent role are ready, or it can be configured so that just a single <code>Pod</code> of the dependent role must be
ready.</p>

<p>For example:
In the yaml snippet below there are two roles, <code>data</code> and <code>proxy</code></p>

<markup
lang="yaml"

>- role: data
  replicas: 3      <span class="conum" data-value="1" />
- role: proxy
  startQuorum:     <span class="conum" data-value="2" />
    - role: data
      podCount: 1</markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role does not specify a <code>startQuorum</code> so this role will be created immediately by the operator.</li>
<li data-value="2">The <code>proxy</code> role has a start quorum that means that the <code>proxy</code> role depends on the <code>data</code> role.
The <code>podCount</code> field is set to <code>1</code> meaning that the <code>proxy</code> role will not be created until at least <code>1</code> of the <code>data</code> role <code>Pods</code>
in in the <code>Ready</code> state.</li>
</ul>
<p>Omitting the <code>podCount</code> from the quorum means that the role will not start until all of the configured replicas of the dependent
role are ready; for example:</p>

<markup
lang="yaml"

>- role: data
  replicas: 3
- role: proxy
  startQuorum:  <span class="conum" data-value="1" />
    - role: data</markup>

<ul class="colist">
<li data-value="1">The <code>proxy</code> role&#8217;s <code>startQuorum</code> just specifies a dependency on the <code>data</code> role with no <code>podCount</code> so all <code>3</code> of the <code>data</code>
role&#8217;s <code>Pods</code> must be <code>Ready</code> before the <code>proxy</code> role is created by the operator.</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">Setting a <code>podCount</code> less than or equal to zero is the same as not specifying a count.</p>
</div>

<h3 id="_multiple_dependencies">Multiple Dependencies</h3>
<div class="section">
<p>The <code>startQuorum</code> can specify a dependency on more than on role; for example:</p>

<markup
lang="yaml"

>- role: data      <span class="conum" data-value="1" />
  replicas: 5
- role: proxy
  replicas: 3
- role: web
  startQuorum:    <span class="conum" data-value="2" />
    - role: data
    - role: proxy
      podCount: 1</markup>

<ul class="colist">
<li data-value="1">The <code>data</code> and <code>proxy</code> roles do not specify a <code>startQuorum</code> so these roles will be created immediately by the operator.</li>
<li data-value="2">The <code>web</code> role has a <code>startQuorum</code> the defines a dependency on both the <code>data</code> role and the <code>proxy</code> role. The <code>proxy</code>
dependency also specifies a <code>podCount</code> of <code>1</code>. This means that the operator wil not create the <code>web</code> role until all <code>5</code> replicas
of the <code>data</code> role are <code>Ready</code> and at least <code>1</code> of the <code>proxy</code> role&#8217;s <code>Pods</code> is <code>Ready</code>.</li>
</ul>
</div>

<h3 id="_chained_dependencies">Chained Dependencies</h3>
<div class="section">
<p>It is also possible to chain dependencies, for example:</p>

<markup
lang="yaml"

>- role: data      <span class="conum" data-value="1" />
  replicas: 5
- role: proxy
  replicas: 3
  startQuorum:    <span class="conum" data-value="2" />
    - role: data
- role: web
  startQuorum:    <span class="conum" data-value="3" />
    - role: proxy
      podCount: 1</markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role does not specify a <code>startQuorum</code> so this role will be created immediately by the operator.</li>
<li data-value="2">The <code>proxy</code> role defines a dependency on the <code>data</code> role without a <code>podCount</code> so all three <code>Pods</code> of the <code>data</code> role must be
in a <code>Ready</code> state before the operator will create the <code>proxy</code> role.</li>
<li data-value="3">The <code>web</code> role depends on the <code>proxy</code> role with a <code>podCount</code> of one, so the operator will not create the <code>web</code> role until
at least one <code>proxy</code> role <code>Pod</code> is in a <code>Ready</code> state.</li>
</ul>
<div class="admonition warning">
<p class="admonition-inline">The operator does not validate that a <code>startQuorum</code> makes sense. It is possible to declare a quorum with circular
dependencies, in which case the roles will never start. It would also be possible to create a quorum with a <code>podCount</code> greater
than the <code>replicas</code> value of the dependent role, in which case the quorum would never be met and the role would not start.</p>
</div>
<div class="admonition note">
<p class="admonition-inline">If creating a cluster with multiple explicit roles a <code>startQuorum</code> declared in the cluster&#8217;s default section will be
ignored. A <code>startQuorum</code> can only be specified at the individual role level.</p>
</div>
</div>
</div>
</doc-view>
