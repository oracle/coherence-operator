<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Obtain Coherence Images</dt>
<dd slot="desc"><p>Coherence images are not available from public registries such as Docker Hub and must be pulled from one of two
private registries.</p>
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
Coherence image.
See <router-link to="/clusters/200_private_repos">Using Private Image Registries</router-link></p>

</div>

<h2 id="_coherence_images_from_docker_store">Coherence Images from Docker Store</h2>
<div class="section">
<ul class="ulist">
<li>
<p>In a <a id="" title="" target="_blank" href="https://hub.docker.com/_/oracle-coherence-12c">https://hub.docker.com/_/oracle-coherence-12c</a></p>

</li>
<li>
<p>In a web browser, navigate to <a id="" title="" target="_blank" href="https://hub.docker.com/">Docker Hub</a> and click Sign In.</p>

</li>
<li>
<p>Search for the official Oracle Coherence images</p>

</li>
<li>
<p>Click on the <code>Proceed to Checkout</code> button</p>

</li>
<li>
<p>Accept the license agreements by clicking the check boxes.</p>

</li>
<li>
<p>Click the <code>Get Content</code> button</p>

</li>
</ul>
<p>Once this is done the Docker Hub credentials can be used to create Kubernetes secret to pull the Coherence image.</p>

</div>

<h2 id="_use_imagepullsecrets">Use ImagePullSecrets</h2>
<div class="section">
<p>Kubernetes supports configuring pods to use <code>imagePullSecrets</code> for pulling images. If possible, this is the preferable
and most portable route.
See the <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/containers/images/#specifying-imagepullsecrets-on-a-pod">kubernetes docs</a>
for this.</p>

<p>Once secrets have been created in the namespace that the <code>CoherenceCluster</code> is to be installed in then the secret name
can be specified in the <code>CoherenceCluster</code> <code>spec</code>. It is possible to specify multiple secrets in the case where the different
images being used are pulled from different registries.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  imagePullSecrets:
    - name: coherence-secret  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>coherence-secret</code> will be used for pulling images from the registry associated to the secret</li>
</ul>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  imagePullSecrets:           <span class="conum" data-value="1" />
    - name: coherence-secret
    - name: application-secret</markup>

<ul class="colist">
<li data-value="1">In this case two secrets have been specified, <code>coherence-secret</code> and <code>application-secret</code></li>
</ul>
<p>Also see <router-link to="/clusters/200_private_repos">Using Private Image Registries</router-link></p>

</div>
</doc-view>
