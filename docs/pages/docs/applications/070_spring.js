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

<p><strong>Spring Boot 2.x or 3.x</strong>
This type of image can be run by the Coherence Operator by specifying an application type of <code>spring</code>
for Spring Boot 2.x applications or <code>spring3</code> for SpringBoot 3.x applications.
The application type is set in the <code>spec.application.type</code> field and by setting the working directory
to the exploded directory, for example:</p>

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
<li data-value="1">The <code>type</code> field set to <code>spring</code> tells the Operator that this is a Spring Boot 2.x application.</li>
<li data-value="2">The working directory has been set to the directory containing the exploded Spring Boot application.</li>
</ul>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  image: my-spring-app:1.0.0
  application:
    type: spring3        <span class="conum" data-value="1" />
    workingDir: /spring  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>type</code> field set to <code>spring3</code> tells the Operator that this is a Spring Boot 3.x application.</li>
<li data-value="2">The working directory has been set to the directory containing the exploded Spring Boot application.</li>
</ul>
<p>When the Operator starts the application it will then run a command equivalent to:</p>

<p><strong>Spring Boot 2.x</strong></p>

<markup
lang="bash"

>cd /spring &amp;&amp; java org.springframework.boot.loader.PropertiesLauncher</markup>

<p><strong>Spring Boot 3.x</strong></p>

<markup
lang="bash"

>cd /spring &amp;&amp; java org.springframework.boot.loader.launch.PropertiesLauncher</markup>

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

<p><strong>Spring Boot 2.x</strong></p>

<markup
lang="bash"

>java --class-path /app/libs/catalogue-1.0.0.jar org.springframework.boot.loader.PropertiesLauncher</markup>

<p><strong>Spring Boot 3.x</strong></p>

<markup
lang="bash"

>java --class-path /app/libs/catalogue-1.0.0.jar org.springframework.boot.loader.launch.PropertiesLauncher</markup>

<div class="admonition note">
<p class="admonition-inline">The Operator does not run the fat jar using the <code>java -jar</code> command because it needs to add various other
JVM arguments and append to the classpath, so it has to run the <code>org.springframework.boot.loader.PropertiesLauncher</code>
class as opposed to the <code>org.springframework.boot.loader.JarLauncher</code> that <code>java -jar</code> would run.</p>
</div>
</div>

<h3 id="_using_could_native_buildpacks">Using Could Native Buildpacks</h3>
<div class="section">
<p>If the Spring Boot Maven or Gradle plugin has been used to produce an image using
<a id="" title="" target="_blank" href="https://docs.spring.io/spring-boot/reference/packaging/container-images/cloud-native-buildpacks.html">Cloud Native Buildpacks</a>
these images can work with the Coherence Operator.</p>

<p>Images using Cloud Native Buildpacks contain a special launcher executable the runs the Java application. This makes it more complex than normal for the Operator to provide a custom Java command.
For images built using Cloud Native Buildpacks to work the <code>Coherence</code> resource must be configured to execute the images entry point instead of the Operator injecting a command line.</p>

<div class="admonition important">
<p class="admonition-textlabel">Important</p>
<p ><p>Due to the way that the Coherence Operator configures JVM arguments
when configured to use an image entry point, the image must be running
Java 11 or higher.</p>
</p>
</div>
<p>Instead of building a custom command line, the Operator uses the <code>JDK_JAVA_OPTIONS</code> environment variable to pass and
configured JVM options and system properties to the Spring application.
This is a standard environment variable that the JVM will effectively use to pre-pend JVM arguments to its command line.</p>

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
    type: spring <span class="conum" data-value="1" />
    useImageEntryPoint: true <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The application type has been set to <code>spring</code> (for Spring Boot 2.x) or <code>spring3</code> (for Spring Boot 3.x) so that the
operator knows that this is a Spring Boot application, and the fact that the image is a Buildpacks image will be auto-discovered.</li>
<li data-value="2">The Operator will run the image&#8217;s entry point and set the <code>JDK_JAVA_OPTIONS</code> environment variable
to pass arguments to the JVM</li>
</ul>
<p>For more information on using image entry points with the Coherence operator see the
<router-link to="/docs/applications/080_entrypoint">Run an Image Entry Point</router-link> documentation.</p>


<h4 id="_buildpacks_jvm_arguments">Buildpacks JVM Arguments</h4>
<div class="section">
<p>A typical Spring Boot buildpack launcher will attempt to configure options such as heap size based on the container
resource limits configured, so this must be taken into account if using any of the memory options available in the
<code>Coherence</code> CRD as there may be conflicting configurations.</p>

</div>
</div>
</div>
</doc-view>
