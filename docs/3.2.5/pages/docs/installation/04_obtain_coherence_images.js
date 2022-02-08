<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Obtain Coherence Images</dt>
<dd slot="desc"><p>For most use-cases we expect the developer to provide a suitable Coherence application image to be
run by the operator. For POCs, demos and experimentation the Coherence Operator uses the OSS Coherence CE image
when no image has been specified for a <code>Coherence</code> resource.
Commercial Coherence images are not available from public image registries and must be pulled from the
middleware section of <a id="" title="" target="_blank" href="https://container-registry.oracle.com">Oracle Container Registry.</a></p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_coherence_images_from_oracle_container_registry">Coherence Images from Oracle Container Registry</h2>
<div class="section">
<p>Get the Coherence Docker image from the Oracle Container Registry:</p>

<ul class="ulist">
<li>
<p>In a web browser, navigate to <a id="" title="" target="_blank" href="https://container-registry.oracle.com/">Oracle Container Registry</a> and click Sign In.</p>

</li>
<li>
<p>Enter your Oracle credentials or create an account if you don&#8217;t have one.</p>

</li>
<li>
<p>Search for coherence in the Search Oracle Container Registry field.</p>

</li>
<li>
<p>Click coherence in the search result list.</p>

</li>
<li>
<p>On the Oracle Coherence page, select the language from the drop-down list and click Continue.</p>

</li>
<li>
<p>Click Accept on the Oracle Standard Terms and Conditions page.</p>

</li>
</ul>
<p>Once this is done the Oracle Container Registry credentials can be used to create Kubernetes secret to pull the
Coherence image.</p>

</div>

<h2 id="_use_imagepullsecrets">Use ImagePullSecrets</h2>
<div class="section">
<p>Kubernetes supports configuring pods to use <code>imagePullSecrets</code> for pulling images. If possible, this is the preferable
and most portable route.
See the <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/containers/images/#specifying-imagepullsecrets-on-a-pod">kubernetes docs</a>
for this.</p>

<p>Once secrets have been created in the namespace that the <code>Coherence</code> resource is to be installed in then the secret name
can be specified in the <code>Coherence</code> CRD <code>spec</code>. It is possible to specify multiple secrets in the case where the different
images being used are pulled from different registries.</p>

<p>For example to use the commercial Coherence 14.1.1.0.0 image from OCR specify the image and image pull secrets in
the <code>Coherence</code> resource yaml</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test-cluster
spec:
  image: container-registry.oracle.com/middleware/coherence:14.1.1.0.0
  imagePullSecrets:
    - name: coherence-secret  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>coherence-secret</code> will be used for pulling images from the registry associated to the secret</li>
</ul>
<p>Also see <router-link to="/docs/installation/05_private_repos">Using Private Image Registries</router-link></p>

</div>
</doc-view>
