<doc-view>

<h2 id="_coherence_operator_kubernetes_web_hooks">Coherence Operator Kubernetes Web-Hooks</h2>
<div class="section">
<p>The Coherence Operator uses Kubernetes admission control webhooks to validate and provide default values for
Coherence resources
(see the <a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/">Kubernetes documentation</a>
for more details on web-hooks).</p>

<p>The Coherence Operator webhooks will validate a <code>Coherence</code> resources when it is created or updated contain.
For example, the <code>replicas</code> count is not negative. If a <code>Coherence</code> resource is invalid it will be rejected before it
gets stored into Kubernetes.</p>


<h3 id="_webhook_scope">Webhook Scope</h3>
<div class="section">
<p>Webhooks in Kubernetes are a cluster resource, not a namespaced scoped resource, so consequently there is typically only
a single webhook installed for a given resource type. If the Coherence Operator is installed as a cluster scoped operator
then this is not a problem but if multiple Coherence Operators are deployed then they could all attempt to install the
webhooks and update or overwrite a previous configuration. This might not be an issue if all of the operators deployed
in a Kubernetes cluster are the same version but different versions could cause issues.</p>

</div>
</div>

<h2 id="_webhook_certificates">Webhook Certificates</h2>
<div class="section">
<p>Kubernetes requires webhooks to expose an API over https and consequently this requires certificates to be created.
By default, the Coherence Operator will create a self-signed CA certificate and key for use with its webhooks.
Alternatively it is possible to use an external certificate manager such as the commonly used
<a id="" title="" target="_blank" href="https://github.com/jetstack/cert-manager">Cert Manager</a>.
Configuring and using Cert Manager is beyond the scope of this documentation.</p>

</div>
</doc-view>
