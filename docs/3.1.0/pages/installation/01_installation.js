<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Coherence Operator Installation</dt>
<dd slot="desc"><p>The Coherence Operator is available as a Docker image <code>oracle/coherence-operator:3.1.0</code> that can
easily be installed into a Kubernetes cluster.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_prerequisites">Prerequisites</h2>
<div class="section">
<ul class="ulist">
<li>
<p>Access to Oracle Coherence Docker images, or self-built Coherence images.</p>

</li>
<li>
<p>Access to a Kubernetes v1.12.0+ cluster. The Operator test pipeline runs on Kubernetes tested on v1.12 upto v1.19</p>

</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">OpenShift - the Coherence Operator works without modification on OpenShift, but some versions
of the Coherence images will not work out of the box.
See the <router-link to="/installation/06_openshift">OpensShift</router-link> section of the documentation that explains how to
run Coherence clusters with the Operator on OpenShift.</p>
</div>
</div>

<h2 id="_installing_with_helm">Installing With Helm</h2>
<div class="section">
<p>The simplest way to install the Coherence Operator is to use the Helm chart.
This ensures that all the correct resources will be created in Kubernetes.</p>


<h3 id="_add_the_coherence_helm_repository">Add the Coherence Helm Repository</h3>
<div class="section">
<p>Add the <code>coherence</code> helm repository using the following commands:</p>

<markup
lang="bash"

>helm repo add coherence https://oracle.github.io/coherence-operator/charts

helm repo update</markup>

</div>

<h3 id="_install_the_coherence_operator_helm_chart">Install the Coherence Operator Helm chart</h3>
<div class="section">
<p>Once the Coherence Helm repo has been configured the Coherence Operator can be installed using a normal Helm 3
install command:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \      <span class="conum" data-value="1" />
    coherence \                    <span class="conum" data-value="2" />
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">where <code>&lt;namespace&gt;</code> is the namespace that the Coherence Operator will be installed into.</li>
<li data-value="2"><code>coherence</code> is the name of this Helm installation.</li>
</ul>

<h4 id="_uninstall_the_coherence_operator_helm_chart">Uninstall the Coherence Operator Helm chart</h4>
<div class="section">
<p>To uninstall the operator:</p>

<markup
lang="bash"

>helm delete coherence-operator --namespace &lt;namespace&gt;</markup>

</div>
</div>
</div>

<h2 id="_operator_scope">Operator Scope</h2>
<div class="section">
<p>The recommended way to install the Coherence Operator is to install a single instance of the operator into a namespace
and where it will then control <code>Coherence</code> resources in all namespaces across the Kubernetes cluster.
Alternatively it may be configured to watch a sub-set of namespaces by setting the <code>WATCH_NAMESPACE</code> environment variable.
The watch namespace(s) does not have to include the installation namespace.</p>

<p>In theory, it is possible to install multiple instances of the Coherence Operator into different namespaces, where
each instances monitors a different set of namespaces. There are a number of potential issues with this approach, so
it is not recommended.</p>

<ul class="ulist">
<li>
<p>Only one CRD can be installed - Different releases of the Operator may use slightly different CRD versions, for example
v3.1.1 may introduce extra fields not in v3.1.0. As the CRD version is <code>v1</code> there is no guarantee which CRD version has
actually installed, which could lead to subtle issues.</p>

</li>
<li>
<p>The operator creates and installs defaulting and validating web-hooks. A web-hook is associated to a CRD resource so
installing multiple web-hooks for the same resource may lead to issues. If an operator is uninstalled, but the web-hook
configuration remains, then Kubernetes will not accept modifications to resources of that type as it will be
unable to contact the web-hook.</p>

</li>
</ul>
<p>To set the watch namespaces when installing with helm set the <code>watchNamespaces</code> value, for example:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set watchNamespaces=payments,catalog,customers <span class="conum" data-value="1" />
    coherence-operator \
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">The <code>payments</code>, <code>catalog</code> and <code>customers</code> namespaces will be watched by the Operator.</li>
</ul>
</div>

<h2 id="_operator_image">Operator Image</h2>
<div class="section">
<p>The Helm chart uses a default registry to pull the Operator image from.
If the image needs to be pulled from a different location (for example an internal registry) then the <code>image</code> field
in the values file can be set, for example:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set image=images.com/coherence/coherence-operator:3.1.0 <span class="conum" data-value="1" />
    coherence-operator \
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">The image used to run the Operator will be <code>images.com/coherence/coherence-operator:3.1.0</code>.</li>
</ul>

<h3 id="_image_pull_secrets">Image Pull Secrets</h3>
<div class="section">
<p>If the image is to be pulled from a secure repository that requires credentials then the image pull secrets
can be specified.
See the Kubernetes documentation on <a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/">Pulling from a Private Registry</a>.</p>


<h4 id="_add_pull_secrets_using_a_values_file">Add Pull Secrets Using a Values File</h4>
<div class="section">
<p>Create a values file that specifies the secrets, for example the <code>private-repo-values.yaml</code> file below:</p>

<markup
lang="yaml"
title="private-repo-values.yaml"
>imagePullSecrets:
- name: registry-secrets</markup>

<p>Now use that file in the Helm install command:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    -f private-repo-values.yaml <span class="conum" data-value="1" />
    coherence-operator \
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">the <code>private-repo-values.yaml</code> values fle will be used by Helm to inject the settings into the Operator deployment</li>
</ul>
</div>

<h4 id="_add_pull_secrets_using_set">Add Pull Secrets Using --Set</h4>
<div class="section">
<p>Although the <code>imagePullSecrets</code> field in the values file is an array of <code>name</code> to value pairs it is possible to set
these values with the normal Helm <code>--set</code> parameter.</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set imagePullSecrets[0].name=registry-secrets <span class="conum" data-value="1" />
    coherence-operator \
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">this creates the same imagePullSecrets as the values file above.</li>
</ul>
</div>
</div>
</div>
</doc-view>
