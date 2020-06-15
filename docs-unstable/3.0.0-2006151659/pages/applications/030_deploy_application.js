<doc-view>

<h2 id="_deploy_coherence_application_images">Deploy Coherence Application Images</h2>
<div class="section">
<p>Once a custom application image has been built (as described in <router-link to="/applications/020_build_application">Build Application Images</router-link>)
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
or Java classes, in the default classpath locations used by the Operator.
When this is not the case additional configuration can be added to specify additional application settings, these include:</p>

<ul class="ulist">
<li>
<p>setting the <router-link to="/jvm_settings/020_classpath">classpath</router-link></p>

</li>
<li>
<p>specifying the <router-link to="/applications/040_application_main">application main</router-link></p>

</li>
<li>
<p>specifying <router-link to="/applications/050_application_args">application arguments</router-link></p>

</li>
<li>
<p>specifying the <router-link to="/applications/060_application_working_dir">working directory</router-link></p>

</li>
</ul>
</div>
</div>
</doc-view>
