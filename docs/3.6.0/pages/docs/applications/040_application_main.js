<doc-view>

<h2 id="_set_the_application_main">Set the Application Main</h2>
<div class="section">
<p>The Coherence container in the deployment&#8217;s Pods will, by default, run <code>com.tangosol.net.DefaultCacheServer</code> as the Java main class.
It is possible to change this when running a custom application that requires a different main.</p>

<p>The name of the main is set in the <code>application.main</code> field in the <code>Coherence</code> spec.</p>

<p>For example, if the deployment is using a custom image <code>catalogue:1.0.0</code> that requires a custom main class
called <code>com.acme.Catalogue</code> the <code>Coherence</code> resource would look like this:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  image: catalogue:1.0.0
  application:
    main: com.acme.Catalogue <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>com.acme.Catalogue</code> will be run as the main class.</li>
</ul>
<p>The example would be equivalent to the Coherence container running:</p>

<markup
lang="bash"

>$ java com.acme.Catalogue</markup>

</div>
</doc-view>
