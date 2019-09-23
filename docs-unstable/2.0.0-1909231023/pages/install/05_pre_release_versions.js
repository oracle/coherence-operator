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
<p class="admonition-inline">Pre-release versions of the Coherence Operator are not guaranteed to be bug free and should not be used for
production use. Pre-release versions of the Helm chart and Docker images may be removed and hence made unavailable
without notice. APIs and CRD specifications are not guaranteed to remain stable or backwards compatible  between
pre-release versions.</p>
</div>
<p>To access pre-release versions of the Helm chart add the unstable chart repository.</p>

<markup
lang="bash"

>helm repo add coherence-unstable https://oracle.github.io/coherence-operator/charts-unstable

helm repo update</markup>

<p>To list all of the available Coherence Operator chart versions:</p>

<markup
lang="bash"

>helm search coherence-operator -l</markup>

<p>The <code>-l</code> parameter shows all versions as opposed to just the latest versions if it was omitted.</p>

<p>A specific pre-release version of the Helm chart can be installed using the <code>--version</code> argument,
for example to use version <code>2.0.0-alpha1</code>:</p>

<markup
lang="bash"

>helm install coherence-unstable/coherence-operator \
    --version 2.0.0-alpha1 \    <span class="conum" data-value="1" />
    --namespace &lt;namespace&gt; \       <span class="conum" data-value="2" />
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

<h3 id="_accessing_pre_release_coherence_operator_docker_images">Accessing Pre-Release Coherence Operator Docker Images</h3>
<div class="section">
<div class="admonition note">
<p class="admonition-inline">Not all pre-release images are pushed to public repositories such as Docker Hub.
Consequently when installing those versions of the Coherence Operator credentials and Kubernetes pull secrets will be required.</p>
</div>
<p>For example to access an image in the <code>iad.ocir.io/odx-stateservice</code> repository you would need to have your own credentials
for that repository so that a secret can be created.</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; \                                     <span class="conum" data-value="1" />
  create secret docker-registry coherence-operator-secret \  <span class="conum" data-value="2" />
  --docker-server=$DOCKER_REPO \                             <span class="conum" data-value="3" />
  --docker-username=$DOCKER_USERNAME \                       <span class="conum" data-value="4" />
  --docker-password=$DOCKER_PASSWORD \                       <span class="conum" data-value="5" />
  --docker-email=$DOCKER_EMAIL                               <span class="conum" data-value="6" /></markup>

<ul class="colist">
<li data-value="1">Replace &lt;namespace&gt; with the Kubernetes namespace that the Coherence Operator will be installed into.</li>
<li data-value="2">In this example the name of the secret to be created is <code>coherence-operator-secret</code></li>
<li data-value="3">Replace <code>$DOCKER_REPO</code> with the name of the Docker repository that the images are to be pulled from.</li>
<li data-value="4">Replace <code>$DOCKER_USERNAME</code> with your username for that repository.</li>
<li data-value="5">Replace <code>$DOCKER_PASSWORD</code> with your password for that repository.</li>
<li data-value="6">Replace <code>$DOCKER_EMAIL</code> with your email (or even a fake email).</li>
</ul>
<p>See the <a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/">Kubernetes documentation</a>
on pull secrets for more details.</p>

<p>Once a secret has been created in the namespace the Coherence Operator can be installed with an extra value parameter
to specify the secret to use:</p>

<markup
lang="bash"

>helm install coherence-unstable/coherence-operator \
    --version 2.0.0-1909130555 \
    --namespace &lt;namespace&gt; \
    --set imagePullSecrets=coherence-operator-secret \  <span class="conum" data-value="1" />
    --name coherence-operator</markup>

<ul class="colist">
<li data-value="1">Set the pull secret to use to the same name that was created above.</li>
</ul>
</div>
</div>
</doc-view>
