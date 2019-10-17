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
<p>Once this is done the Docker Hub credentials can be used to create Kubernetes secret to pull the Coherence image.
See <router-link to="/clusters/200_private_repos">Using Private Image Registries</router-link></p>

</div>
</doc-view>
