<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Environment Variables</dt>
<dd slot="desc"><p>It is possible to pass arbitrary environment variables to the <code>Pods</code> that are created for a Coherence cluster.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_environment_variables">Environment Variables</h2>
<div class="section">
<p>Environment variables can be configured in a <code>CoherenceCluster</code> and will be passed through to the Coherence <code>Pods</code>
created for the roles in the cluster. Environment variables are configured in the <code>env</code> field of the spec. The format
for setting environment variables is exactly the same as when configuring them in a
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/inject-data-application/define-environment-variable-container/">Kubernetes <code>Container</code></a>.</p>


<h3 id="_environment_variables_in_the_implicit_role">Environment Variables in the Implicit Role</h3>
<div class="section">
<p>If configuring a single implicit role environment variables are set in the <code>spec.env</code> section; for example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  env:
    - name: FOO        <span class="conum" data-value="1" />
      value: "foo-val"
    - name: BAR        <span class="conum" data-value="2" />
      value: "bar-val"</markup>

<ul class="colist">
<li data-value="1">The <code>FOO</code> environment variable with a value of <code>foo-val</code> will be passed to the <code>coherence</code> container in the <code>Pods</code>
created for the implicit role.</li>
<li data-value="2">The <code>BAR</code> environment variable with a value of <code>bar-val</code> will be passed to the <code>coherence</code> container in the <code>Pods</code>
created for the implicit role.</li>
</ul>
</div>

<h3 id="_environment_variables_in_explicit_roles">Environment Variables in Explicit Roles</h3>
<div class="section">
<p>When configuring one or more explicit roles in the <code>roles</code> section of the spec environment variables can be configured
for each role.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data
      env:
        - name: FOO        <span class="conum" data-value="1" />
          value: "foo-val"
    - role: proxy
      env:
        - name: BAR        <span class="conum" data-value="2" />
          value: "bar-val"</markup>

<ul class="colist">
<li data-value="1">All <code>Pods</code> created for the <code>data</code> role will have the <code>FOO</code> environment variable set to <code>foo-val</code></li>
<li data-value="2">All <code>Pods</code> created for the <code>proxy</code> role will have the <code>BAR</code> environment variable set to <code>bar-val</code></li>
</ul>
</div>

<h3 id="_environment_variables_in_explicit_roles_with_defaults">Environment Variables in Explicit Roles With Defaults</h3>
<div class="section">
<p>When configuring one or more explicit roles it is also possible to configure environment variables at the
defaults level. These environment variables will be shared by all <code>Pods</code> in all roles unless specifically
overridden for a role. An environment variable is only overridden in a role by declaring a role level
environment variable with the same name. When creating the <code>Pods</code> configuration the Coherence Operator will
merge the list of default environment variables with the role&#8217;s list of environment variables.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  env:
    - name: FOO
      value: "foo-val"
  roles:
    - role: data           <span class="conum" data-value="1" />
    - role: proxy          <span class="conum" data-value="2" />
      env:
        - name: BAR
          value: "bar-val"
    - role: web            <span class="conum" data-value="3" />
      env:
        - name: FOO
          value: "foo-web"
        - name: BAR
          value: "bar-web"</markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role does not have any environment variables configured so it will just inherit the <code>FOO=foo-val</code>
environment variable from the defaults.</li>
<li data-value="2">The <code>proxy</code> role has the <code>BAR=bar-val</code> environment variables configured and will also inherit the <code>FOO=foo-val</code>
environment variable from the defaults.</li>
<li data-value="3">The <code>web</code> role has will override the <code>FOO</code> environment variable from the default with <code>FOO=foo-web</code>. It also
has its own <code>BAR=bar-web</code> environment variable.</li>
</ul>
</div>
</div>
</doc-view>
