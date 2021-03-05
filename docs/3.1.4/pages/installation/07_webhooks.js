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
<p>Assuming Cert Manager has been installed in the Kubernetes cluster then to use it for managing the web-hook certificates,
the Operator needs to be installed with the <code>CERT_TYPE</code> environment variable set to <code>cert-manager</code>.</p>

<p>The Operator will then detect the version of Cert Manager and automatically create the required self-signed <code>Issuer</code>
and <code>Certificate</code> resources. Cert Manager will detect these and create the <code>Secret</code>. This may cause the operator Pod to
re-start until the <code>Secret</code> has been created.</p>

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

<h3 id="_manual_certificates">Manual Certificates</h3>
<div class="section">
<p>If certificates will managed some other way (for example by Cert Manager managing real certificates) then the
<code>CERT_TYPE</code> environment variable should be set to <code>manual</code>.</p>

<p>Before the Operator starts the <code>Secret</code> must exist containing the valid certificates.
The Operator will use the certificates that it finds in the <code>Secret</code> to create the web-hook resources.</p>

<p>To set the certificate manager to use when installing the Helm chart, set the <code>webhookCertType</code> value:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set webhookCertType=manual <span class="conum" data-value="1" />
    coherence-operator \
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">The certificate manager will be set to <code>manual</code></li>
</ul>
</div>
</div>
</doc-view>
