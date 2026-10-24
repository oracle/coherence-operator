<doc-view>

<h2 id="_set_the_classpath">Set the Classpath</h2>
<div class="section">
<p>The Coherence container in the <code>Pods</code> in a <code>Coherence</code> resource deployment runs a Java application and as such requires a classpath
with at a minimum <code>coherence.jar</code>. There are certain defaults that the Operator will use to work out the classpath to use
but additional classpath elements can be provided to the configuration.</p>


<h3 id="_the_classpath_environment_variable">The <code>CLASSPATH</code> Environment Variable</h3>
<div class="section">
<p>If the image to be run has the <code>CLASSPATH</code> environment variable set this will be used as part of the classpath.</p>

</div>

<h3 id="_the_coherence_home_environment_variable">The <code>COHERENCE_HOME</code> Environment Variable</h3>
<div class="section">
<p>If the image to be run has the <code>COHERENCE_HOME</code> environment variable set this will be used to add the following elements
to the classpath:</p>

<ul class="ulist">
<li>
<p><code>$COHERENCE_HOME/lib/coherence.jar</code></p>

</li>
<li>
<p><code>$COHERENCE_HOME/conf</code></p>

</li>
</ul>
<p>These will be added to the end of the classpath. For example in an image that has <code>CLASSPATH=/home/root/lib/*</code>
and <code>COHERENCE_HOME</code> set to <code>/oracle/coherence</code> the effective classpath used will be:</p>

<pre>/home/root/lib/*:/oracle/coherence/lib/coherence.jar:/oracle/coherence/conf</pre>
</div>

<h3 id="_jib_image_classpath">JIB Image Classpath</h3>
<div class="section">
<p>A simple way to build Java images is using <a id="" title="" target="_blank" href="https://github.com/GoogleContainerTools/jib/blob/master/README.md">JIB</a>.
When JIB was with its Maven or Gradle plugin to produce an image it packages the application&#8217;s dependencies, classes
and resources into a set of well-known locations:</p>

<ul class="ulist">
<li>
<p><code>/app/libs/</code> - the jar files that the application depends on</p>

</li>
<li>
<p><code>/app/classes</code> - the application&#8217;s class files</p>

</li>
<li>
<p><code>/app/resources</code> - the application&#8217;s other resources</p>

</li>
</ul>
<p>By default, the Operator will add these locations to the classpath. These classpath elements will be added before any
value set by the <code>CLASSPATH</code> or <code>COHERENCE_HOME</code> environment variables.</p>

<p>For example in an image that has <code>CLASSPATH=/home/root/lib/\*</code>
and <code>COHERENCE_HOME</code> set to <code>/oracle/coherence</code> the effective classpath used will be:</p>

<pre>/app/libs/*:/app/classes:/app/resources:/home/root/lib/*:/oracle/coherence/lib/coherence.jar:/oracle/coherence/conf</pre>

<h4 id="_exclude_the_jib_classpath">Exclude the JIB Classpath</h4>
<div class="section">
<p>If the image is not a JIB image there could be occasions when automatically adding <code>/app/libs/*:/app/classes:/app/resources</code>
to the classpath causes issues, for example one or more of those locations exists with files in that should not be on the
classpath. In this case the <code>Coherence</code> CRD spec has a field to specify that the JIB classpath should not be used.</p>

<p>The <code>spec.jvm.useJibClasspath</code> field can be set to <code>false</code> to exclude the JIB directories from the classpath
(the default value is <code>true</code>).</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    useJibClasspath: false  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>useJibClasspath</code> is set to <code>false</code>. Even if any of the <code>/app/resources</code>, <code>/app/classes</code> or <code>/app/libs/</code>
directories exist in the image they will not be added to the classpath.</li>
</ul>
</div>
</div>

<h3 id="_additional_classpath_elements">Additional Classpath Elements</h3>
<div class="section">
<p>If an image will be used that has artifacts in locations other than the defaults discussed above then it is possible
to specify additional elements to be added to the classpath. The <code>jvm.classpath</code> field in the <code>Coherence</code> CRD spec
allows a list of extra classpath values to be provided. These elements will be added <em>after</em> the JIB classpath
described above.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    classpath:          <span class="conum" data-value="1" />
      - "/data/lib/*"
      - "/data/config"</markup>

<ul class="colist">
<li data-value="1">The <code>classpath</code> field adds <code>/data/lib/*</code> and <code>/data/config</code> to the classpath.
In an image without the <code>CLASSPATH</code> or <code>COHERENCE_HOME</code> environment variables the effective classpath would be:</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">There is no validation of the elements of the classpath. The elements will not be verified to ensure that the locations
exist. As long as they are valid values to be used in a JVM classpath they will be accepted.</p>
</div>
</div>
</div>

<h2 id="_environment_variable_expansion">Environment Variable Expansion</h2>
<div class="section">
<p>The Operator supports environment variable expansion in classpath entries.
The runner in the Coherence container will replace <code>${var}</code> or <code>$var</code> in classpath entries with the corresponding environment variable name.</p>

<p>For example if a container has an environment variable of <code>APP_HOME</code> set to <code>/myapp</code> then it could be used in the classpath like this:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    classpath:
      - "${APP_HOME}/lib/*"  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The actual classpath entry at runtime will resolve to <code>/myapp/lib/*</code></li>
</ul>
<p>Any environment variable that is present when the Coherence container starts can be used, this would include variables created as part of the image and variables specified in the Coherence yaml.</p>

</div>
</doc-view>
