<doc-view>

<h2 id="_spring_boot_applications">Spring Boot Applications</h2>
<div class="section">
<p>The Coherence Operator supports running images that contain Spring Boot applications.
Exactly how easy this is depends on how the image has been built.</p>

<p>When the operator runs an image it overrides the default image entrypoint and uses its own launcher.
This allows the operator to properly configure various Coherence properties that the launcher then uses to build the
command line to actually run your application. With some types of image this is not a straight forward Java command line
so the Operator requires a bit more information adding to the <code>Coherence</code> deployment yaml.</p>


<h3 id="_using_jib_images">Using JIB Images</h3>
<div class="section">
<p>The simplest way to build an application image to run with the Coherence Operator (including Spring Boot applications)
is to use the <a id="" title="" target="_blank" href="https://github.com/GoogleContainerTools/jib/blob/master/README.md">JIB</a> tool.
JIB images will work out of the box with the operator, even for a Spring Boot application, as described in
<router-link to="/docs/applications/020_build_application">Building Applications</router-link> and
<router-link to="/docs/applications/030_deploy_application">Deploying Applications</router-link>.</p>

<p>If you have used the Spring Maven or Gradle plugins to build the application into a fat jar, but you then build the image
using the <a id="" title="" target="_blank" href="https://github.com/GoogleContainerTools/jib/blob/master/README.md">JIB</a> plugin then JIB will detect the fat
jar and package the image in an exploded form that will run out of the box with the operator.</p>

</div>

<h3 id="_using_an_exploded_spring_boot_image">Using an Exploded Spring Boot Image</h3>
<div class="section">
<p>Another way to build a Spring Boot image is to explode the Spring Boot jar into a directory structure in the image.</p>

<p>For example, if a Spring Boot jar has been exploded into a directory called <code>/spring</code>, the image contents might look
like the diagram below; where you can see the <code>/spring</code> directory contains the Spring Boot application.</p>

<markup


>├── bin
├── boot
├── dev
├─⊕ etc
├─⊕ home
├─⊕ lib
├─⊕ lib64
├── proc
├── root
├── run
├── sbin
├── spring
│   ├── BOOT-INF
│   │   ├─⊕ classes
│   │   ├── classpath.idx
│   │   └─⊕ lib
│   ├── META-INF
│   │   ├── MANIFEST.MF
│   │   └─⊕ maven
│   └── org
│       └── springframework
│           └─⊕ boot
├── sys
├── tmp
├─⊕ usr
└─⊕ var</markup>

<p>This type of image can be run by the Coherence Operator by specifying an application type of <code>spring</code> in the
<code>spec.application.type</code> field and by setting the working directory to the exploded directory, for example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  image: my-spring-app:1.0.0
  application:
    type: spring         <span class="conum" data-value="1" />
    workingDir: /spring  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>type</code> field set to <code>spring</code> tells the Operator that this is a Spring Boot application.</li>
<li data-value="2">The working directory has been set to the directory containing the exploded Spring Boot application.</li>
</ul>
<p>When the Operator starts the application it will then run a command equivalent to:</p>

<markup
lang="bash"

>cd /spring &amp;&amp; java org.springframework.boot.loader.PropertiesLauncher</markup>

</div>

<h3 id="_using_a_spring_boot_fat_jar">Using a Spring Boot Fat Jar</h3>
<div class="section">
<p>It is not recommended to build images containing fat jars for various reasons which can easily be found on the internet.
If you feel that you must build your application as a Spring Boot fat jar then this can still work with the Coherence Operator.</p>

<p>The Java command line to run a Spring Boot fat jar needs to be something like <code>java -jar my-app.jar</code>
where <code>my-app.jar</code> is the fat jar.
This means that the Operator&#8217;s launcher needs to know the location of the fat jar in the image, so this must
be provided in the <code>Coherence</code> deployment yaml.</p>

<p>For example, suppose that an application has been built into a fat jar names <code>catalogue-1.0.0.jar</code> which is in the
<code>/app/libs</code> directory in the image, so the full path to the jar is <code>/app/libs/catalogue-1.0.0.jar</code>.
This needs to be set in the <code>spec.applicaton.springBootFatJar</code> field of the <code>Coherence</code> yaml.</p>

<p>The <code>spec.application.type</code> field also needs to be set to <code>spring</code> so that the Operator knows that this is a
Spring Boot application</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  image: catalogue:1.0.0
  application:
    type: spring                                      <span class="conum" data-value="1" />
    springBootFatJar: /app/libs/catalogue-1.0.0.jar   <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>type</code> field set to <code>spring</code> tells the Operator that this is a Spring Boot application.</li>
<li data-value="2">The location of the Spring Boot jar has been set.</li>
</ul>
<p>When the Operator starts the application it will then run a command equivalent to:</p>

<markup
lang="bash"

>java -cp /app/libs/catalogue-1.0.0.jar org.springframework.boot.loader.PropertiesLauncher</markup>

<div class="admonition note">
<p class="admonition-inline">The Operator does not run the fat jar using the <code>java -jar</code> command because it needs to add various other
JVM arguments and append to the classpath, so it has to run the <code>org.springframework.boot.loader.PropertiesLauncher</code>
class as opposed to the <code>org.springframework.boot.loader.JarLauncher</code> that <code>java -jar</code> would run.</p>
</div>
</div>

<h3 id="_using_could_native_buildpacks">Using Could Native Buildpacks</h3>
<div class="section">
<p>If the Spring Boot Maven or Gradle plugin has been used to produce an image using
<a id="" title="" target="_blank" href="https://spring.io/blog/2020/01/27/creating-docker-images-with-spring-boot-2-3-0-m1">Cloud Native Buildpacks</a>
these images can work with the Coherence Operator.</p>

<div class="admonition warning">
<p class="admonition-textlabel">Warning</p>
<p ><p>Due to limitation on the way that arguments can be passed to the JVM when using Buildpacks images the Coherence
operator will only work with images containing a JVM greater than Java 11.
Although the Buildpacks launcher will honour the <code>JAVA_OPTS</code> or <code>JAVA_TOOL_OPTIONS</code> environment variables there appear
to be size limitations for the values of these variables that make it impractical for the Operator to use them.
The Operator therefore creates a JVM arguments file to pass the arguments to the JVM.
At the time of writing these docs, Java 8 (which is the default version of Java used by the Spring Boot plugin) does not
support the use of argument files for the JVM.</p>

<p>It is simple to configure the version of the JVM used by the Spring Boot plugin, for example in Maven:</p>

<markup
lang="xml"

>&lt;plugin&gt;
  &lt;groupId&gt;org.springframework.boot&lt;/groupId&gt;
  &lt;artifactId&gt;spring-boot-maven-plugin&lt;/artifactId&gt;
  &lt;version&gt;2.3.4.RELEASE&lt;/version&gt;
  &lt;configuration&gt;
    &lt;image&gt;
      &lt;env&gt;
        &lt;BP_JVM_VERSION&gt;11.*&lt;/BP_JVM_VERSION&gt;
      &lt;/env&gt;
    &lt;/image&gt;
  &lt;/configuration&gt;
&lt;/plugin&gt;</markup>
</p>
</div>
<p>When creating a <code>Coherence</code> deployment for a Spring Boot Buildpacks image The application type must be set to <code>spring</code>.
The Operator&#8217;s launcher will automatically detect that the image is a Buildpacks image and launch the application using
the Buildpacks launcher.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  image: catalogue:1.0.0
  application:
    type: spring <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The application type has been set to <code>spring</code> so that the operator knows that this is a Spring Boot application,
and the fact that the image is a Buildpacks image will be auto-discovered.</li>
</ul>
<p>When the Operator starts the application it will then run the buildpacks launcher with a command equivalent
to this:</p>

<markup
lang="bash"

>/cnb/lifecycle/launcher java @jvm-args-file org.springframework.boot.loader.PropertiesLauncher</markup>


<h4 id="_buildpacks_detection">Buildpacks Detection</h4>
<div class="section">
<p>If for some reason buildpacks auto-detection does not work properly the <code>Coherence</code>
CRD contains a filed to force buildpacks to be enabled or disabled.</p>

<p>The <code>boolean</code> field <code>spec.application.cloudNativeBuildPack.enabled</code> can be set to <code>true</code> to enable buildpacks or false
to disable buildpack.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  image: catalogue:1.0.0
  application:
    type: spring            <span class="conum" data-value="1" />
    cloudNativeBuildPack:
      enabled: true         <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The application type has been set to <code>spring</code> so that the operator knows that this is a Spring Boot application</li>
<li data-value="2">The <code>cloudNativeBuildPack.enabled</code> field has been set to <code>true</code> to force the Operator to use the Buildpacks launcher.</li>
</ul>
</div>

<h4 id="_specify_the_buildpacks_launcher">Specify the Buildpacks Launcher</h4>
<div class="section">
<p>A Cloud Native Buildpacks image uses a launcher mechanism to run the executable(s) in the image. The Coherence Operator
launcher will configure the application and then invoke the same buildpacks launcher.
The Coherence Operator assumes that the buildpacks launcher is in the image in the location <code>/cnb/lifecycle/launcher</code>.
If a buildpacks image has been built with the launcher in a different location then the <code>Coherence</code> CRD contains
a field to set the new location.</p>

<p>The <code>spec.application.cloudNativeBuildPack.enabled</code> field.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  image: catalogue:1.0.0
  application:
    type: spring                    <span class="conum" data-value="1" />
    cloudNativeBuildPack:
      launcher: /buildpack/launcher <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The application type has been set to <code>spring</code> so that the operator knows that this is a Spring Boot application</li>
<li data-value="2">The buildpacks launcher that the Operator will invoke is located at <code>/buildpack/launcher</code>.</li>
</ul>
</div>

<h4 id="_buildpack_jvm_arguments">Buildpack JVM Arguments</h4>
<div class="section">
<p>A typical Spring Boot buildpack launcher will attempt to configure options such as heap size based on the container
resource limits configured, so this must be taken into account if using any of the memory options available in the
<code>Coherence</code> CRD as there may be conflicting configurations.</p>

</div>
</div>
</div>
</doc-view>
