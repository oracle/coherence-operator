<doc-view>

<h2 id="_install_using_helm">Install Using Helm</h2>
<div class="section">
<p>For more flexibility the simplest way to install the Coherence Operator is to use the Helm chart.
This ensures that all the correct resources will be created in Kubernetes.</p>

<div class="admonition warning">
<p class="admonition-textlabel">Warning</p>
<p ><p><strong>Helm Upgrades</strong></p>

<p>Now that the Coherence Operator no longer installs the CRDs when it starts, the CRDs are
installed as part of the Helm chart. This works ok when installing the operator for the first
time into a Kubernetes cluster. If the Helm chart is being used to upgrade an existing install
of an earlier Coherence Operator version, or the CRDs already exist, then the installation
can fail with an error message similar to this:</p>

<p><code>Error: INSTALLATION FAILED: Unable to continue with install: CustomResourceDefinition "coherence.coherence.oracle.com" in namespace "" exists and cannot be imported into the current release: invalid ownership metadata; label validation error: missing key "app.kubernetes.io/managed-by": must be set to "Helm"; annotation validation error: missing key "meta.helm.sh/release-name": must be set to "operator"; annotation validation error: missing key "meta.helm.sh/release-namespace": must be set to "default"</code></p>

<p>This is because Helm will refuse to overwrite any resources that it did not originally install.</p>

<p>In this case the CRDs have to be installed manually from the CRD manifest files before the
Helm install or upgrade can be performed.
The Helm install or upgrade then needs to set the <code>installCrd</code> value to <code>false</code> so that the CRDs
are not installed as part of the Helm chart install.</p>
</p>
</div>

<h3 id="_add_the_coherence_helm_repository">Add the Coherence Helm Repository</h3>
<div class="section">
<p>Add the <code>coherence</code> helm repository using the following commands:</p>

<markup
lang="bash"

>helm repo add coherence https://oracle.github.io/coherence-operator/charts

helm repo update</markup>

<div class="admonition note">
<p class="admonition-inline">To avoid confusion, the URL <code><a id="" title="" target="_blank" href="https://oracle.github.io/coherence-operator/charts">https://oracle.github.io/coherence-operator/charts</a></code> is a Helm repo, it is not
a website you open in a browser. You may think we shouldn&#8217;t have to say this, but you&#8217;d be surprised.</p>
</div>
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
</div>

<h3 id="helm-operator-image">Set the Operator Image</h3>
<div class="section">
<p>The Helm chart uses a default Operator image from
<code>container-registry.oracle.com/middleware/coherence-operator:3.5.0</code>.
If the image needs to be pulled from a different location (for example an internal registry) then there are two ways to override the default.
Either set the individual <code>image.registry</code>, <code>image.name</code> and <code>image.tag</code> values, or set the whole image name by setting the <code>image</code> value.</p>

<p>For example, if the Operator image has been deployed into a private registry named <code>foo.com</code> but
with the same image name <code>coherence-operator</code> and tag <code>3.5.0</code> as the default image,
then just the <code>image.registry</code> needs to be specified.</p>

<p>In the example below, the image used to run the Operator will be <code>foo.com/coherence-operator:3.5.0</code>.</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set image.registry=foo.com \
    coherence-operator \
    coherence/coherence-operator</markup>

<p>All three of the image parts can be specified individually using <code>--set</code> options.
In the example below, the image used to run the Operator will
be <code>foo.com/operator:1.2.3</code>.</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set image.registry=foo.com \
    --set image.name=operator \
    --set image.tag=1.2.3
    coherence-operator \
    coherence/coherence-operator</markup>

<p>Alternatively, the image can be set using a single <code>image</code> value.
For example, the command below will set the Operator image to <code>images.com/coherence-operator:0.1.2</code>.</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set image=images.com/coherence-operator:0.1.2 \
    coherence-operator \
    coherence/coherence-operator</markup>

</div>

<h3 id="helm-pull-secrets">Image Pull Secrets</h3>
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

<h4 id="_add_pull_secrets_using_set">Add Pull Secrets Using --set</h4>
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

<h3 id="_change_the_operator_replica_count">Change the Operator Replica Count</h3>
<div class="section">
<p>To change the replica count when installing the Operator using Helm, the <code>replicas</code> value can be set.</p>

<p>For example, to change the replica count from 3 to 1, the <code>--set replicas=1</code> option can be used.</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set replicas=1
    coherence \
    coherence/coherence-operator</markup>

</div>

<h3 id="helm-watch-ns">Set the Watch Namespaces</h3>
<div class="section">
<p>To set the watch namespaces when installing with helm set the <code>watchNamespaces</code> value, for example:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set watchNamespaces=payments,catalog,customers \
    coherence-operator \
    coherence/coherence-operator</markup>

<p>The <code>payments</code>, <code>catalog</code> and <code>customers</code> namespaces will be watched by the Operator.</p>


<h4 id="_set_the_watch_namespace_to_the_operators_install_namespace">Set the Watch Namespace to the Operator&#8217;s Install Namespace</h4>
<div class="section">
<p>When installing the Operator using the Helm chart, there is a convenience value that can be set if the
Operator should only monitor the same namespace that it is installed into.
By setting the <code>onlySameNamespace</code> value to <code>true</code> the watch namespace will be set to the installation namespace.
If the <code>onlySameNamespace</code> value is set to <code>true</code> then any value set for the <code>watchNamespaces</code> value will be ignored.</p>

<p>For example, the command below will set <code>onlySameNamespace</code> to true, and the Operator will be installed into,
and only monitor the <code>coh-testing</code> namespace.</p>

<markup
lang="bash"

>helm install  \
    --namespace coh-testing \
    --set onlySameNamespace=true \
    coherence-operator \
    coherence/coherence-operator</markup>

<p>In the example below, the <code>onlySameNamespace</code> is set to true, so the Operator will be installed into,
and only monitor the <code>coh-testing</code> namespace. Even though the <code>watchNamespaces</code> value is set, it will be ignored.</p>

<markup
lang="bash"

>helm install  \
    --namespace coh-testing \
    --set watchNamespaces=payments,catalog,customers \
    --set onlySameNamespace=true \
    coherence-operator \
    coherence/coherence-operator</markup>

</div>
</div>

<h3 id="helm-sec-context">Install the Operator with a Security Context</h3>
<div class="section">
<p>The Operator container can be configured with a Pod <code>securityContext</code> or a container <code>securityContext</code>,
so that it runs as a non-root user.</p>

<p>This can be done using a values file:</p>

<p><strong>Set the Pod securityContext</strong></p>

<markup
lang="yaml"
title="security-values.yaml"
>podSecurityContext:
  runAsNonRoot: true
  runAsUser: 1000</markup>

<p><strong>Set the Container securityContext</strong></p>

<markup
lang="yaml"
title="security-values.yaml"
>securityContext:
  runAsNonRoot: true
  runAsUser: 1000</markup>

<p>Then the <code>security-values.yaml</code> values file above can be used in the Helm install command.</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --values security-values.yaml \
    coherence \
    coherence/coherence-operator</markup>

<p>Alternatively, the Pod or container <code>securityContext</code> values can be set on the command line as <code>--set</code> parameters:</p>

<p><strong>Set the Pod securityContext</strong></p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set podSecurityContext.runAsNonRoot=true \
    --set podSecurityContext.runAsUser=1000 \
    coherence \
    coherence/coherence-operator</markup>

<p><strong>Set the Container securityContext</strong></p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set securityContext.runAsNonRoot=true \
    --set securityContext.runAsUser=1000 \
    coherence \
    coherence/coherence-operator</markup>

</div>

<h3 id="helm-labels">Set Additional Labels</h3>
<div class="section">
<p>When installing the Operator with Helm, it is possible to set additional labels to be applied to the Operator Pods
and to the Operator Deployment.</p>


<h4 id="_adding_pod_labels">Adding Pod Labels</h4>
<div class="section">
<p>To add labels to the Operator Pods set the <code>labels</code> value, either on the command line using <code>--set</code> or in the values file.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Setting <code>labels</code> will only apply the additional labels to the Operator Pods, they will not be applied to any other resource created by the Helm chart.</p>
</p>
</div>
<p>For example, using the command line:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set labels.one=value-one \
    --set labels.two=value-two \
    coherence \
    coherence/coherence-operator</markup>

<p>The command above would add the following additional labels <code>one</code> and <code>two</code> to the Operator Pod as shown below:</p>

<markup
lang="yaml"

>apiVersion: v1
kind: Pod
metadata:
  name: coherence-operator
  labels:
    one: value-one
    two: value-two</markup>

<p>The same labels could also be specified in a values file:</p>

<markup

title="add-labels-values.yaml"
>labels:
  one: value-one
  two: value-two</markup>

</div>

<h4 id="_adding_deployment_labels">Adding Deployment Labels</h4>
<div class="section">
<p>To add labels to the Operator Deployment set the <code>deploymentLabels</code> value, either on the command line using <code>--set</code> or in the values file.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Setting <code>deploymentLabels</code> will only apply the additional labels to the Deployment, they will not be applied to any other resource created by the Helm chart.</p>
</p>
</div>
<p>For example, using the command line:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set deploymentLabels.one=value-one \
    --set deploymentLabels.two=value-two \
    coherence \
    coherence/coherence-operator</markup>

<p>The command above would add the following additional labels <code>one</code> and <code>two</code> to the Operator Pod as shown below:</p>

<markup
lang="yaml"

>apiVersion: apps/v1
kind: Deployment
metadata:
  name: coherence-operator
  labels:
    one: value-one
    two: value-two</markup>

<p>The same labels could also be specified in a values file:</p>

<markup

title="add-labels-values.yaml"
>deploymentLabels:
  one: value-one
  two: value-two</markup>

</div>
</div>

<h3 id="helm-annotations">Set Additional Annotations</h3>
<div class="section">
<p>When installing the Operator with Helm, it is possible to set additional annotations to be applied to the Operator Pods
and to the Operator Deployment.</p>


<h4 id="_adding_pod_annotations">Adding Pod Annotations</h4>
<div class="section">
<p>To add annotations to the Operator Pods set the <code>annotations</code> value, either on the command line using <code>--set</code> or in the values file.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Setting <code>annotations</code> will only apply the additional annotations to the Operator Pods, they will not be applied to any other resource created by the Helm chart.</p>
</p>
</div>
<p>For example, using the command line:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set annotations.one=value-one \
    --set annotations.two=value-two \
    coherence \
    coherence/coherence-operator</markup>

<p>The command above would add the following additional annotations <code>one</code> and <code>two</code> to the Operator Pod as shown below:</p>

<markup
lang="yaml"

>apiVersion: v1
kind: Pod
metadata:
  name: coherence-operator
  annotations:
    one: value-one
    two: value-two</markup>

<p>The same annotations could also be specified in a values file:</p>

<markup

title="add-annotations-values.yaml"
>annotations:
  one: value-one
  two: value-two</markup>

</div>

<h4 id="_adding_deployment_annotations">Adding Deployment Annotations</h4>
<div class="section">
<p>To add annotations to the Operator Deployment set the <code>deploymentAnnotations</code> value, either on the command line using <code>--set</code> or in the values file.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Setting <code>deploymentAnnotations</code> will only apply the additional annotations to the Deployment, they will not be applied to any other resource created by the Helm chart.</p>
</p>
</div>
<p>For example, using the command line:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set deploymentAnnotations.one=value-one \
    --set deploymentAnnotations.two=value-two \
    coherence \
    coherence/coherence-operator</markup>

<p>The command above would add the following additional annotations <code>one</code> and <code>two</code> to the Operator Pod as shown below:</p>

<markup
lang="yaml"

>apiVersion: apps/v1
kind: Deployment
metadata:
  name: coherence-operator
  annotations:
    one: value-one
    two: value-two</markup>

<p>The same annotations could also be specified in a values file:</p>

<markup

title="add-annotations-values.yaml"
>deploymentAnnotations:
  one: value-one
  two: value-two</markup>

</div>
</div>

<h3 id="helm-job">CoherenceJob CRD Support</h3>
<div class="section">
<p>By default, the Operator will install both CRDs, <code>Coherence</code> and <code>CoherenceJob</code>.
If support for <code>CoherenceJob</code> is not required then it can be excluded from being installed setting the
Operator command line parameter <code>--install-job-crd</code> to <code>false</code>.</p>

<p>When installing with Helm, the <code>allowCoherenceJobs</code> value can be set to <code>false</code> to disable support for <code>CoherenceJob</code>
resources (the default value is <code>true</code>).</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set allowCoherenceJobs=false \
    coherence \
    coherence/coherence-operator</markup>

</div>
</div>

<h2 id="helm-upgrade">Upgrade the Coherence Operator Using Helm</h2>
<div class="section">
<p>If the Coherence operator was originally installed using Helm then it can be upgraded to a new
version using a newer Helm chart.</p>

<p>To upgrade to the latest version of the Coherence operator simply use the Helm upgrade command as
shown below.</p>

<markup
lang="bash"

>helm upgrade  \
    --namespace &lt;namespace&gt; \
    coherence \
    coherence/coherence-operator</markup>

<p>The command above will use all the default configurations, but the usual methods of applying
values to the install can be used.</p>


<h3 id="helm-upgrade-350">Upgrading From pre-3.5.0 Versions</h3>
<div class="section">
<p>Before version 3.5.0 of the Coherence operator, the operator used to install the CRDs
when it started. In 3.5.0 this behaviour was changed and the operator no longer installs
the CRDs, these must be installed along with the operator. The 3.5.0 and above Helm chart
includes the CRDs.</p>

<p>This causes an issue when performing a Helm upgrade from a pre-3.5.0 version because Helm
did not install the CRDs. When attempting an upgrade Helm will display an error similar to
the one below:</p>

<markup


>Error: INSTALLATION FAILED: Unable to continue with install: CustomResourceDefinition
"coherence.coherence.oracle.com" in namespace "" exists and cannot be imported into the
current release: invalid ownership metadata; label validation error: missing key
"app.kubernetes.io/managed-by": must be set to "Helm"; annotation validation error:
missing key "meta.helm.sh/release-name": must be set to "operator"; annotation validation
error: missing key "meta.helm.sh/release-namespace": must be set to "default"</markup>

<p>This is because Helm will refuse to overwrite any resources that it did not originally install.
There are a few options to work around this.</p>

<div class="admonition warning">
<p class="admonition-textlabel">Warning</p>
<p ><p>As a work-around to the issue, you should not uninstall the existing CRDs.
Any running Coherence clusters being managed by the Operator will be deleted
if the CRDs are deleted.</p>
</p>
</div>

<h4 id="_continue_to_install_the_crds_manually">Continue to install the CRDs manually</h4>
<div class="section">
<p>The CRDs can be installed manually from the manifest yaml files as described
in the documentation section <router-link :to="{path: '/docs/installation/011_install_manifests', hash: '#manual-crd'}">Manually Install the CRDs</router-link>
The Helm install or upgrade then needs to set the <code>installCrd</code> value to <code>false</code> so that the CRDs
are not installed as part of the Helm chart install.</p>

<div class="admonition warning">
<p class="admonition-textlabel">Warning</p>
<p ><p>The CRDs for the new version <em>MUST</em> be installed <em>BEFORE</em> running the Helm upgrade.</p>
</p>
</div>
<markup
lang="bash"

>helm upgrade  \
    --namespace &lt;namespace&gt; \
    --set installCrd=false
    coherence \
    coherence/coherence-operator</markup>

</div>

<h4 id="_patch_the_crds_so_helm_manages_them">Patch the CRDs So Helm Manages Them</h4>
<div class="section">
<p>The CRDs can be patched with the required labels and annotations so that Helm thinks it
originally installed them and will then update them.</p>

<p>The commands below can be used to patch the CRDs:</p>

<markup
lang="bash"

>export HELM_RELEASE=operator
export HELM_NAMESPACE=coherence
kubectl patch customresourcedefinition coherence.coherence.oracle.com \
    --patch '{"metadata": {"labels": {"app.kubernetes.io/managed-by": "Helm"}}}'
kubectl patch customresourcedefinition coherence.coherence.oracle.com \
    --patch "{\"metadata\": {\"annotations\": {\"meta.helm.sh/release-name\": \"$HELM_RELEASE\"}}}"
kubectl patch customresourcedefinition coherence.coherence.oracle.com \
    --patch "{\"metadata\": {\"annotations\": {\"meta.helm.sh/release-namespace\": \"$HELM_NAMESPACE\"}}}"
kubectl patch customresourcedefinition coherencejob.coherence.oracle.com \
    --patch '{"metadata": {"labels": {"app.kubernetes.io/managed-by": "Helm"}}}'
kubectl patch customresourcedefinition coherencejob.coherence.oracle.com \
    --patch "{\"metadata\": {\"annotations\": {\"meta.helm.sh/release-name\": \"$HELM_RELEASE\"}}}"
kubectl patch customresourcedefinition coherencejob.coherence.oracle.com \
    --patch "{\"metadata\": {\"annotations\": {\"meta.helm.sh/release-namespace\": \"$HELM_NAMESPACE\"}}}"</markup>

<p>The first line exports the name of the Helm release being upgraded.
The second line exports the name of the Kubernetes namespace the operator was installed into.</p>

<p>After patching as described above the operator can be upgraded with a normal Helm upgrade command:</p>

<markup
lang="bash"

>helm upgrade  \
    --namespace $HELM_NAMESPACE \
    $HELM_RELEASE \
    coherence/coherence-operator</markup>

</div>
</div>
</div>

<h2 id="helm-uninstall">Uninstall the Coherence Operator Helm chart</h2>
<div class="section">
<p>To uninstall the operator:</p>

<markup
lang="bash"

>helm delete coherence-operator --namespace &lt;namespace&gt;</markup>

</div>
</doc-view>
