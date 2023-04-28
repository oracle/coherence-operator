<doc-view>

<h2 id="_operator_web_hooks">Operator Web-Hooks</h2>
<div class="section">
<p>The Coherence Operator uses Kubernetes
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/">dynamic admission control</a>
commonly known as defaulting and validating web-hooks. As the name implies, these are used to provide default values
for some fields in a <code>Coherence</code> resource and to also validate <code>Coherence</code> resources on creation and update.
The operator creates and configures the two web-hooks when it starts.</p>


<h3 id="_webhook_scope">Webhook Scope</h3>
<div class="section">
<p>Webhooks in Kubernetes are a cluster resource, not a namespaced scoped resource, so consequently there is typically only
a single webhook installed for a given resource type. If the Coherence Operator has been installed as a cluster scoped
operator then this is not a problem but if multiple Coherence Operators have been deployed then they could all attempt
to install the webhooks and update or overwrite a previous configuration.
This might not be an issue if all the operators deployed in a Kubernetes cluster are the same version but different
versions could cause issues.
This is one of the reasons that it is recommended to install a single cluster scoped Coherence Operator.</p>

</div>
</div>

<h2 id="_manage_web_hook_certificates">Manage Web-Hook Certificates</h2>
<div class="section">
<p>A web-hook requires certificates to be able to work in Kubernetes.
By default, the operator will create and manage self-signed certificates for this purpose.
These certificates are created using the Kubernetes certificate
It is possible to use other certificates, either managed by the
<a id="" title="" target="_blank" href="https://cert-manager.io/docs/installation/kubernetes/">Kubernetes cert-manager</a> or managed manually.</p>

<p>The certificates should be stored in a <code>Secret</code> named <code>coherence-webhook-server-cert</code> in the same namespace that
the operator has installed in. (although this name can be changed if required). This <code>Secret</code> must exist, or the operator
wil fail to start. The Operator Helm chart will create this <code>Secret</code> when the Operator is managing its own self-signed
certs, otherwise the <code>Secret</code> must be created manually or by an external certificate manager.</p>


<h3 id="_self_signed_certificates">Self-Signed Certificates</h3>
<div class="section">
<p>This is the default option, the operator will create and manage a set of self-signed certificates.
The Operator will update the <code>Secret</code> with its certificates and create the <code>MutatingWebhookConfiguration</code> and
<code>ValidatingWebhookConfiguration</code> resources configured to use those certificates.</p>

</div>

<h3 id="_cert_manager_self_signed">Cert Manager (Self-Signed)</h3>
<div class="section">
<p>Assuming <a id="" title="" target="_blank" href="https://cert-manager.io/docs/installation/kubernetes/">Kubernetes Cert Manager</a> has been installed in the
Kubernetes cluster then to use it for managing the web-hook certificates,
the Operator needs to be installed with the <code>CERT_TYPE</code> environment variable set to <code>cert-manager</code>.</p>

<p>The Operator will then detect the version of Cert Manager and automatically create the required self-signed <code>Issuer</code>
and <code>Certificate</code> resources. Cert Manager will detect these and create the <code>Secret</code>. This may cause the operator Pod to
re-start until the <code>Secret</code> has been created.</p>


<h4 id="_install_using_manifest_file">Install Using Manifest File</h4>
<div class="section">
<p>If installing the operator using the manifest yaml file first replace the occurrences of <code>self-signed</code> in the yaml file with <code>cert-manager</code>.</p>

<p>For example:</p>

<markup
lang="bash"

>curl -L https://github.com/oracle/coherence-operator/releases/download/v3.2.11/coherence-operator.yaml \
    -o coherence-operator.yaml
sed -i s/self-signed/cert-manager/g coherence-operator.yaml
kubectl apply -f coherence-operator.yaml</markup>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>On MacOS the <code>sed</code> command is slightly different for in-place replacement
and requires an empty string after the <code>-i</code> parameter:</p>

<markup
lang="bash"

>sed -i '' s/self-signed/cert-manager/g coherence-operator.yaml</markup>
</p>
</div>
</div>

<h4 id="_install_using_helm">Install Using Helm</h4>
<div class="section">
<p>To set the certificate manager to use when installing the Helm chart, set the <code>webhookCertType</code> value:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set webhookCertType=cert-manager <span class="conum" data-value="1" />
    coherence-operator \
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">The certificate manager will be set to <code>cert-manager</code></li>
</ul>
</div>
</div>

<h3 id="_manual_certificates">Manual Certificates</h3>
<div class="section">
<p>If certificates will be managed some other way (for example by Cert Manager managing real certificates) then the
<code>CERT_TYPE</code> environment variable should be set to <code>manual</code>.</p>

<p>A <code>Secret</code> must exist in the namespace the operator will be installed into containing the CA certificate, certificate
and key files that the operator will use to configure the web-hook. The files must exist with the names expected by the operator.
The default name of the <code>Secret</code> expected by the operator is <code>coherence-webhook-server-cert</code> but this can be changed.</p>

<p>The certificates in the <code>Secret</code> must be valid for the <code>Service</code> name that exposes the Coherence web-hook.
The default format of the DNS used for the certificate CN (common name) is <code>coherence-operator-webhook.&lt;namespace&gt;.svc</code>
where <code>&lt;namespace&gt;</code> is the namespace the operator is installed into.
Additional names may also be configured using the different formats of Kubernetes <code>Service</code> DNS names.</p>

<p>For example, if the Operator is installed into a namespace named <code>coherence</code> the <code>Service</code> DNS names would be:</p>

<markup


>  - coherence-operator-webhook.coherence
  - coherence-operator-webhook.coherence.svc
  - coherence-operator-webhook.coherence.svc.cluster.local</markup>

<p>An example of the format of the <code>Secret</code> is shown below:</p>

<markup
lang="yaml"
title="sh"
>apiVersion: v1
kind: Secret
metadata:
  name: coherence-webhook-server-cert
type: Opaque
data:
  ca.crt: ... # &lt;base64 endocde CA certificate file&gt;
  tls.crt: ... # &lt;base64 endocde certificate file&gt;
  tls.key: ... # &lt;base64 endocde private key file&gt;</markup>

<div class="admonition warning">
<p class="admonition-textlabel">Warning</p>
<p ><p>If a <code>Secret</code> with the name specified in <code>webhookCertSecret</code> does not exist in the namespace the operator
is being installed into then the operator Pod will not start as the <code>Secret</code> will be mounted as a volume
in the operator Pod.</p>
</p>
</div>

<h4 id="_install_using_manifest_file_2">Install Using Manifest File</h4>
<div class="section">
<p>If installing the operator using the manifest yaml file first replace the occurrences of <code>self-signed</code> in the yaml file with <code>cert-manager</code>.</p>

<p>For example:</p>

<markup
lang="bash"

>curl -L https://github.com/oracle/coherence-operator/releases/download/v3.2.11/coherence-operator.yaml \
    -o coherence-operator.yaml
sed -i s/self-signed/manual/g coherence-operator.yaml
kubectl apply -f coherence-operator.yaml</markup>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>On MacOS the <code>sed</code> command is slightly different for in-place replacement
and requires an empty string after the <code>-i</code> parameter:</p>

<markup
lang="bash"

>sed -i '' s/self-signed/cert-manager/g coherence-operator.yaml</markup>
</p>
</div>
</div>

<h4 id="_install_using_helm_2">Install Using Helm</h4>
<div class="section">
<p>To configure the operator to use manually managed certificates when installing the Helm chart,
set the <code>webhookCertType</code> value.</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set webhookCertType=manual \ <span class="conum" data-value="1" />
    coherence-operator \
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">The certificate manager will be set to <code>manual</code> and the operator will expect to find a <code>Secret</code> named <code>coherence-webhook-server-cert</code></li>
</ul>
<p>To use manually managed certificates and store the keys and certs in a different secret, set the secret
name using the <code>webhookCertSecret</code> value.</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set webhookCertType=manual \ <span class="conum" data-value="1" />
    --set webhookCertSecret=operator-certs \ <span class="conum" data-value="2" />
    coherence-operator \
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">The certificate manager will be set to <code>manual</code></li>
<li data-value="2">The name of the secret is set to <code>operator-certs</code></li>
</ul>
<p>The Coherence Operator will now expect to find the keys and certs in a <code>Secret</code> named <code>operator-certs</code> in
the same namespace that the Operator is deployed into.</p>

</div>
</div>

<h3 id="no-hooks">Install the Operator Without Web-Hooks</h3>
<div class="section">
<p>It is possible to start the Operator without it registering any web-hooks with the API server.</p>

<div class="admonition caution">
<p class="admonition-textlabel">Caution</p>
<p ><p>Running the Operator without web-hooks is not recommended.
The admission web-hooks validate the <code>Coherence</code> resource yaml before it gets into the k8s cluster.
Without the web-hooks, invalid yaml will be accepted by k8s and the Operator will then log errors
when it tries to reconcile invalid yaml. Or worse, the Operator will create an invalid <code>StatefulSet</code>
which will then fail to start.</p>
</p>
</div>

<h4 id="_install_using_manifest_file_3">Install Using Manifest File</h4>
<div class="section">
<p>If installing using the manifest yaml files, then you need to edit the <code>coherence-operator.yaml</code> manifest to add a
command line argument to the Operator.</p>

<p>Update the <code>controller-manager</code> deployment and add an argument, edit the section that looks like this:</p>

<markup
lang="yaml"

>        args:
          - operator
          - --enable-leader-election</markup>

<p>and add the additional <code>--enable-webhook=false</code> argument like this:</p>

<markup
lang="yaml"

>        args:
          - operator
          - --enable-leader-election
          - --enable-webhook=false</markup>

<p>apiVersion: apps/v1
kind: Deployment
metadata:
  name: controller-manager</p>

</div>

<h4 id="_installing_using_helm">Installing Using Helm</h4>
<div class="section">
<p>If installing the Operator using Helm, the <code>webhooks</code> value can be set to false in the values file or
on the command line.</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set webhooks=false \
    coherence-operator \
    coherence/coherence-operator</markup>

</div>
</div>
</div>
</doc-view>
