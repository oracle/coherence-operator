<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Installing Manually</dt>
<dd slot="desc"><p>It is possible to install the Coherence Operator by crating the required yaml manually and installing it into Kubernetes.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_rbac">RBAC</h2>
<div class="section">
<p>In Kubernetes clusters with RBAC enabled the Coherence Operator requires certain RBAC resources to be created.</p>

<markup
lang="bash"

>sh examples/create_role.sh</markup>

<p>The example RBAC script creates a <code>ServiceAccount</code> with the name <code>coherence-operator</code> if a different name is required
then the name in the <code>example/example-rbac.yaml</code> should be modified.</p>

<p>The shell script above will use default values for the role name, role binding name and install into the <code>default</code>
namespace. If the default values need to be changed or the operator is to be installed into a namespace other than
<code>default</code> the script can be run with the following environment variables:</p>

<markup
lang="bash"

>export ROLE_NAME=&lt;role-name&gt;
export ROLE_BINDING_NAME=&lt;role-binding-name&gt;
export NAMESPACE=&lt;namespace&gt;
sh example/create_role.sh</markup>

</div>

<h2 id="_deployment">Deployment</h2>
<div class="section">
<p>Create a <code>Deployment</code> for the Coherence Operator using the example script.</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; create -f example/example-deployment.yaml</markup>

<p>where <code>&lt;namespace&gt;</code> is the name of the namespace that the operator is to be deployed into.
If RBAC is being used this will also be the same namespace used to create the RBAC roles.
If the <code>default</code> namespace is being used the <code>-n &lt;namespace&gt;</code> argument can be omitted.</p>

</div>
</doc-view>
