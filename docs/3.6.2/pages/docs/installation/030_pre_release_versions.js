<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Accessing Pre-Release Versions</dt>
<dd slot="desc"><p>Pre-release version of the Coherence Operator are made available from time to time.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_accessing_pre_release_versions">Accessing Pre-Release Versions</h2>
<div class="section">
<div class="admonition warning">
<p class="admonition-inline">We cannot guarantee that pre-release versions of the Coherence Operator are bug free and hence they should
not be used in production.
We reserve the right to remove pre-release versions of the Helm chart and Docker images ant any time and without notice.
We cannot guarantee that APIs and CRD specifications will remain stable or backwards compatible between pre-release versions.</p>
</div>
<p>To access pre-release versions of the Helm chart add the unstable chart repository.</p>

<markup
lang="bash"

>helm repo add coherence-unstable https://oracle.github.io/coherence-operator/charts-unstable

helm repo update</markup>

<p>To list all the available Coherence Operator chart versions:</p>

<markup
lang="bash"

>helm search coherence-operator -l</markup>

<p>The <code>-l</code> parameter shows all versions as opposed to just the latest versions if it was omitted.</p>

<p>A specific pre-release version of the Helm chart can be installed using the <code>--version</code> argument,
for example to use version <code>3.0.0-2005140315</code>:</p>

<markup
lang="bash"

>helm install coherence-unstable/coherence-operator \
    --version 3.0.0-2005140315 \   <span class="conum" data-value="1" />
    --namespace &lt;namespace&gt; \      <span class="conum" data-value="2" />
    --name coherence-operator</markup>

<ul class="colist">
<li data-value="1">The <code>--version</code> argument is used to specify the exact version of the chart</li>
<li data-value="2">The optional <code>--namespace</code> parameter to specify which namespace to install the operator into, if omitted then
Helm will install into whichever is currently the default namespace for your Kubernetes configuration.</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">When using pre-release versions of the Helm chart it is always advisable to install a specific version otherwise
Helm will try to work out the latest version in the pre-release repo and as pre-release version numbers are not strictly
sem-ver compliant this may be unreliable.</p>
</div>
</div>
</doc-view>
