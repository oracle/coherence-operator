<doc-view>

<h2 id="_deploy_coherence_applications">Deploy Coherence Applications</h2>
<div class="section">
<p>Once a custom application image has been built (as described in <router-link to="/docs/applications/020_build_application">Build Custom Application Images</router-link>)
a <code>Coherence</code> resource can be configured to use that image.</p>


<h3 id="_specify_the_image_to_use">Specify the Image to Use</h3>
<div class="section">
<p>To specify the image to use set the <code>image</code> field in the <code>Coherence</code> spec to the name of the image.</p>

<p>For example if there was an application image called <code>catalogue:1.0.0</code> it can be specified like this:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  image: catalogue:1.0.0  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>catalogue:1.0.0</code> will be used in the <code>coherence</code> container in the deployment&#8217;s Pods.</li>
</ul>
<p>The example above would assume that the <code>catalogue:1.0.0</code> has a JVM on the <code>PATH</code> and all the required <code>.jar</code> files,
or Java classes, in the default classpath locations used by the Operator.</p>

</div>

<h3 id="_image_pull_secrets">Image Pull Secrets</h3>
<div class="section">
<p>If your image needs to be pulled from a private registry you may need to provide
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/">image pull secrets</a> for this.</p>

<p>For example, supposing the application image is <code>repo.acme.com/catalogue:1.0.0</code> and that <code>repo.acme.com</code> is a private registry; we might a <code>Secret</code> to the k8s namespace named <code>repo-acme-com-secrets</code>. We can then specify that these secrets are used in the <code>Coherence</code> resource by setting the <code>imagePullSecrets</code> fields. The <code>imagePullSecrets</code> field is a list of secret names, the same format as that used when specifying secrets for a Pod spec.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  image: repo.acme.com/catalogue:1.0.0  <span class="conum" data-value="1" />
  imagePullSecrets:
    - name: repo-acme-com-secrets       <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>repo.acme.com/catalogue:1.0.0</code> image will be used for the application image</li>
<li data-value="2">The <code>Secret</code> named <code>repo-acme-com-secrets</code> will be used to pull images.</li>
</ul>
<p>Multiple secrets can be specified in the case where different images used by different containers are pulled from different registries.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  image: repo.acme.com/catalogue:1.0.0
  imagePullSecrets:
    - name: repo-acme-com-secrets               <span class="conum" data-value="1" />
    - name: oracle-container-registry-secrets</markup>

<ul class="colist">
<li data-value="1">The example above has two image pull secrets, <code>repo-acme-com-secrets</code> and <code>oracle-container-registry-secrets</code></li>
</ul>
</div>

<h3 id="_more_application_configuration">More Application Configuration</h3>
<div class="section">
<p>Additional configuration can be added to specify other application settings, these include:</p>

<ul class="ulist">
<li>
<p>setting the <router-link to="/docs/jvm/020_classpath">classpath</router-link></p>

</li>
<li>
<p>specifying the <router-link to="/docs/applications/040_application_main">application main</router-link></p>

</li>
<li>
<p>specifying <router-link to="/docs/applications/050_application_args">application arguments</router-link></p>

</li>
<li>
<p>specifying the <router-link to="/docs/applications/060_application_working_dir">working directory</router-link></p>

</li>
</ul>
</div>
</div>
</doc-view>
