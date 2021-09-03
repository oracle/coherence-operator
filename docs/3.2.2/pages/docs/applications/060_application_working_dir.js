<doc-view>

<h2 id="_set_the_working_directory">Set the Working Directory</h2>
<div class="section">
<p>When running a custom application there may be a requirement to run in a specific working directory.
The working directory can be specified in the <code>application.workingDir</code> field in the <code>Coherence</code> spec.</p>

<p>For example, a deployment uses a custom image <code>catalogue:1.0.0</code> that requires a custom main class
called <code>com.acme.Catalogue</code>, and that class takes additional arguments.
In this example we&#8217;ll use two fictitious arguments such as a name and a language for the catalogue.
the <code>Coherence</code> resource would look like this:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  image: catalogue:1.0.0          <span class="conum" data-value="1" />
  application:
    workingDir: "/apps/catalogue" <span class="conum" data-value="2" />
    main: "com.acme.Catalogue"    <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The <code>catalogue:1.0.0</code> image will be used.</li>
<li data-value="2">The Java command will be executed in the <code>/apps/catalogue</code> working directory.</li>
<li data-value="3">The Java main class executed will be <code>com.acme.Catalogue</code></li>
</ul>
<p>The example would be equivalent to the Coherence container running:</p>

<markup
lang="bash"

>$ cd /apps/catalogue
$ java com.acme.Catalogue</markup>

</div>
</doc-view>
