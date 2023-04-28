<doc-view>

<h2 id="_build_custom_application_images">Build Custom Application Images</h2>
<div class="section">
<p>To deploy a Coherence application using the operator the application code must be packaged into an image that the
Coherence container in the Pods will run. This image can be any image that contains a JVM as well as the application&#8217;s
jar files, including obviously <code>coherence.jar</code>.</p>

<p>There are many ways to build an image for a Java application so it would be of little value to document the exact steps
for one of them here that might turn out to be used by very few people. One of the simplest ways to build a Java image
is to use <a id="" title="" target="_blank" href="https://github.com/GoogleContainerTools/jib/blob/master/README.md">JIB</a>.
The Operator supports JIB images automatically but any image that meets the requirements of having a JVM and <code>coherence.jar</code>
will be supported. Any version of Java which works with the version of <code>coherence.jar</code> in the image will be suitable.
This can be a JRE, it does not need to be a full JDK.</p>

<p>At a bare minimum the directories in an image might look like this example
(obviously there would be more O/S related files and more JVM files, but they are not relevant for the example):</p>

<markup


>/
|-- app
|    |-- libs                      <span class="conum" data-value="1" />
|         |-- application.jar
|         |-- coherence.jar
|-- usr
     |-- bin
     |    |-- java                 <span class="conum" data-value="2" />
     |
     |-- lib
          |-- jvm
               |-- java-11-openjdk <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The <code>/app/libs</code> directory contains the application jar files. This will be the classpath used to run the application.</li>
<li data-value="2">The <code>/usr/bin/java</code> file is the Java executable and on the <code>PATH</code> in the image (this would be a link to the actual
Java executable location, in this example to <code>/usr/lib/jvm/java-11-openjdk/bin/java</code>.</li>
<li data-value="3">The <code>/usr/lib/jvm/java-11-openjdk/</code> is the actual JVM install location.</li>
</ul>

<h3 id="_image_entrypoint_what_does_the_operator_run">Image <code>EntryPoint</code> - What Does the Operator Run?</h3>
<div class="section">
<p>The image does not need to have an <code>EntryPoint</code> or command specified, it does not need to actually be executable.
If the image does have an <code>EntryPoint</code>, it will just be ignored.</p>

<p>The Coherence Operator actually injects its own <code>runner</code> executable into the container which the container runs and which
in turn builds the Java command line to execute. The <code>runner</code> process looks at arguments and environment variables configured
for the Coherence container and from these constructs a Java command line that it then executes.</p>

<p>The default command might look something like this:</p>

<markup
lang="bash"

>java -cp `/app/resources:/app/classes:/app/libs/*` \
    &lt;JVM args&gt; \
    &lt;System Properties&gt; \
    com.tangosol.net.DefaultCacheServer</markup>

<p>The <code>runner</code> will work out the JVM&#8217;s classpath, args and system properties to add to the command line
and execute the main class <code>com.tangosol.net.DefaultCacheServer</code>.
All these are configurable in the <code>Coherence</code> resource spec.</p>

</div>

<h3 id="_optional_classpath_environment_variable">Optional <code>CLASSPATH</code> Environment Variable</h3>
<div class="section">
<p>If the <code>CLASSPATH</code> environment variable has been set in an image that classpath will be used when running the Coherence
container. Other elements may also be added to the classpath depending on the configuration of the <code>Coherence</code> resource.</p>

</div>

<h3 id="_setting_the_classpath">Setting the Classpath</h3>
<div class="section">
<p>An application image contains <code>.jar</code> files (at least <code>coherence.jar</code>), possibly Java class files, also possibly
other ad-hoc files, all of which need to be on the application&#8217;s classpath.
There are certain classpath values that the operator supports out of the box without needing any extra configuration,
but for occasions where the location of files in the image does not match the defaults a classpath can be specified.</p>

<p>Images built with <a id="" title="" target="_blank" href="https://github.com/GoogleContainerTools/jib/blob/master/README.md">JIB</a>
have a default classpath of <code>/app/resources:/app/classes:/app/libs/*</code>.
When the Coherence container starts if the directories <code>/app/resources</code>, <code>/app/classes</code> or <code>/app/libs/</code> exist in the
image they will automatically be added to the classpath of the JVM. In this way the Operator supports standard JIB
images without requiring additional configuration.</p>

<p>If the image is not a JIB image, or is a JIB image without the standard classpath but one or more of the
<code>/app/resources</code>, <code>/app/classes</code> or <code>/app/libs/</code> directories exist they will still be added to the classpath.
This may be desired or in some cases it may cause issues. It is possible to disable automatically adding these
directories in the <code>Coherence</code> resource spec by setting the <code>jvm.useJibClasspath</code> field to <code>false</code> (the default
value of the field is <code>true</code>).</p>

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
<li data-value="1">The <code>useJibClasspath</code> is set to <code>false</code>. Even if any of the the <code>/app/resources</code>, <code>/app/classes</code> or <code>/app/libs/</code>
directories exist in the image they will not be added to the classpath.</li>
</ul>
<p>If the image is not a JIB image, or is a JIB image without the standard classpath, then additional classpath entries
can be configured as described in the <router-link to="/docs/jvm/020_classpath">setting the classpath</router-link> documentation.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    classpath:          <span class="conum" data-value="1" />
      - "/data/libs/*"  <span class="conum" data-value="2" />
      - "/data/config"</markup>

<ul class="colist">
<li data-value="1">The <code>jvm.classpath</code> field will be used to add additional items to the classpath, the field is a list of strings.</li>
<li data-value="2">Each entry in the <code>jvm.classpath</code> will be appended to the classpath exactly as it is declared, so in this case
the classpath will be <code>/data/libs/*:/data/config</code></li>
</ul>
</div>

<h3 id="_optional_java_home_environment_variable">Optional <code>JAVA_HOME</code> Environment Variable</h3>
<div class="section">
<p>The <code>JAVA_HOME</code> environment variable does not have to be set in the image. If it is set the JVM at that location will
be used to run the application. If it is not set then the <code>java</code> executable <strong>must</strong> be on the <code>PATH</code> in the image.</p>

</div>

<h3 id="_optional_coherence_home_environment_variable">Optional <code>COHERENCE_HOME</code> Environment Variable</h3>
<div class="section">
<p>The <code>COHERENCE_HOME</code> environment variable does not have to be set in an image.
Typically, all the jar files, including <code>coherence.jar</code> would be packaged into a single directory which is then used as
the classpath.
It is possible to run official Coherence images published by Oracle, which have <code>COHERENCE_HOME</code> set, which is then used
by the Operator to set the classpath.</p>

<p>If the <code>COHERENCE_HOME</code> environment variable is set in an image the following entries will be added to the end of the
classpath:</p>

<ul class="ulist">
<li>
<p><code>$COHERENCE_HOME/lib/coherence.jar</code></p>

</li>
<li>
<p><code>$COHERENCE_HOME/conf</code></p>

</li>
</ul>
</div>

<h3 id="_additional_data_volumes">Additional Data Volumes</h3>
<div class="section">
<p>If the application requires access to external storage volumes in Kubernetes it is possible to add additional <code>Volumes</code>
and <code>VolumeMappings</code> to the Pod and containers.</p>

<p>There are three ways to add additional volumes:</p>

<ul class="ulist">
<li>
<p>ConfigMaps - easily add a <code>ConfigMap</code> volume and volume mapping see: <router-link to="/docs/other/050_configmap_volumes">Add ConfigMap Volumes</router-link></p>

</li>
<li>
<p>Secrets - easily add a <code>Secret</code> volume and volume mapping see: <router-link to="/docs/other/060_secret_volumes">Add Secret Volumes</router-link></p>

</li>
<li>
<p>Volumes - easily add any additional volume and volume mapping see: <router-link to="/docs/other/070_add_volumes">Add Volumes</router-link></p>

</li>
</ul>
<p>Both of <code>ConfigMaps</code> and <code>Secrets</code> have been treated as a special case because they are quite commonly used to provide
configurations to Pods, so the <code>Coherence</code> spec provides a simpler way to declare them than for ad-hoc <code>Volumes</code>.</p>

</div>
</div>
</doc-view>
