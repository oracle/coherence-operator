<doc-view>

<h2 id="_install_using_kustomize">Install Using Kustomize</h2>
<div class="section">
<p>If you want to use yaml directly to install the operator, with something like <code>kubectl</code>, you can use the manifest files
published with the GitHub release at this link:
<a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/releases/download/v3.5.0/coherence-operator-manifests.tar.gz">3.5.0 Manifests</a></p>

<p>These manifest files are for use with a tool called Kustomize, which is built into <code>kubectl</code>
see the documentation here: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/">https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/</a></p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>As of v3.5.0 of the Operator the manifest yaml also installs the two CRDs that the Operator uses.
In previous releases the Operator would install the CRDs when it started but this behaviour is disabled by default
when installing with the manifest yaml.</p>
</p>
</div>
<p>Download the
<a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/releases/download/v3.5.0/coherence-operator-manifests.tar.gz">3.5.0 Manifests</a>
from the release page and unpack the file, which should produce a directory called <code>manifests</code> with a structure like this:</p>

<markup


>manifests
    default
        config.yaml
        kustomization.yaml
    manager
        kustomization.yaml
        manager.yaml
        service.yaml
    rbac
        coherence_editor_role.yaml
        coherence_viewer_role.yaml
        kustomization.yaml
        leader_election_role.yaml
        leader_election_role_binding.yaml
        role.yaml
        role_binding.yaml</markup>

<p>There are two ways to use these manifest files, either install using <code>kustomize</code> or generate the yaml and manually
install with <code>kubectl</code>.</p>

<div class="admonition note">
<p class="admonition-inline">All the commands below are run from a console in the <code>manifests/</code> directory from the extracted file above.</p>
</div>

<h3 id="_install_with_kustomize">Install with Kustomize</h3>
<div class="section">
<p>If you have Kustomize installed (or can install it from <a id="" title="" target="_blank" href="https://github.com/kubernetes-sigs/kustomize">https://github.com/kubernetes-sigs/kustomize</a>) you can use
Kustomize to configure the yaml and install.</p>


<h4 id="_change_the_operator_replica_count">Change the Operator Replica Count</h4>
<div class="section">
<p>To change the replica count using Kustomize a patch file needs to be applied.
The Operator manifests include a patch file, named <code>manager/single-replica-patch.yaml</code>, that changes the replica count from 3 to 1. This patch can be applied with the following Kustomize command.</p>

<markup
lang="bash"

>cd ./manager &amp;&amp; kustomize edit add patch \
  --kind Deployment --name controller-manager \
  --path single-replica-patch.yaml</markup>

</div>

<h4 id="_set_image_names">Set Image Names</h4>
<div class="section">
<p>If you need to use different iamge names from the defaults <code>kustomize</code> can be used to specify different names:</p>

<p>Change the name of the Operator image by running the command below, changing the image name to the registry and image name
that you are using for the Operator, for example if you have the images in a custom registry</p>

<markup
lang="bash"

>cd ./manager &amp;&amp; kustomize edit set image controller=myregistry/coherence-operator:3.5.0</markup>

<p>Change the name of the Operator image by running the command below, changing the image name to the registry and image name
that you are using for the Operator utilities image</p>

<markup
lang="bash"

>cd ./manager &amp;&amp; kustomize edit add configmap env-vars --from-literal OPERATOR_IMAGE=myregistry/coherence-operator:3.5.0</markup>

<p>Change the name of the default Coherence image. If you are always going to be deploying your own application images then this
does not need to change.</p>

<markup
lang="bash"

>cd ./manager &amp;&amp; $(GOBIN)/kustomize edit add configmap env-vars --from-literal COHERENCE_IMAGE=$(COHERENCE_IMAGE)</markup>

<p>Set the namespace to install into, the example below sets the namespace to <code>coherence-test</code>:</p>

<markup
lang="bash"

>cd ./default &amp;&amp; /kustomize edit set namespace coherence-test</markup>

</div>

<h4 id="_install">Install</h4>
<div class="section">
<p>The Operator requires a <code>Secret</code> for its web-hook certificates. This <code>Secret</code> needs to exist but can be empty.
The <code>Secret</code> must be in the same namespace that the Operator will be deployed to.
For example, if the Operator namespace is <code>coherence-test</code>, then the <code>Secret</code> can be created with this command:</p>

<markup
lang="bash"

>kubectl -n coherence-test create secret generic coherence-webhook-server-cert</markup>

<p>The Operator can now be installed by running the following command from the <code>manifests</code> directory:</p>

<markup
lang="bash"

>kustomize build ./default | kubectl apply -f -</markup>

</div>
</div>

<h3 id="_generate_yaml_install_with_kubectl">Generate Yaml - Install with Kubectl</h3>
<div class="section">
<p>Instead of using Kustomize to modify and install the Operator we can use <code>kubectl</code> to generate the yaml from the manifests.
You can then edit this yaml and manually deploy it with <code>kubectl</code>.</p>

<p>Run the following command from the <code>manifests</code> directory:</p>

<markup
lang="bash"

>kubectl create --dry-run -k default/ -o yaml &gt; operator.yaml</markup>

<p>This will create a file in the <code>manifests</code> directory called <code>operator.yaml</code> that contains all the yaml required
to install the Operator. You can then edit this yaml to change image names or add other settings.</p>

<p>The Operator can be installed using the generated yaml.</p>

<p>For example if the Operator is to be deployed to the <code>coherence-test</code> namespace:</p>

<markup
lang="bash"

>kubectl -n coherence-test create secret generic coherence-webhook-server-cert
kubectl -n coherence-test create -f operator.yaml</markup>

</div>
</div>
</doc-view>
