<doc-view>

<h2 id="_example_coherence_image_using_jib">Example Coherence Image using JIB</h2>
<div class="section">
<p>This example shows how to build a simple Coherence server image using
<a id="" title="" target="_blank" href="https://github.com/GoogleContainerTools/jib/blob/master/README.md">JIB</a> with either Maven or Gradle.
When building with Maven the project uses the <a id="" title="" target="_blank" href="https://github.com/GoogleContainerTools/jib/blob/master/jib-maven-plugin">JIB Maven Plugin</a>.
When building with Gradle the project uses the <a id="" title="" target="_blank" href="https://github.com/GoogleContainerTools/jib/tree/master/jib-gradle-plugin">JIB Gradle Plugin</a>.</p>

<p>The Coherence Operator has out of the box support for images built with JIB, for example it can automatically detect the class path to use and run the correct main class.</p>

<p>This simple application does not actually contain any code, a real application would obviously contain code and other resources.</p>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p><img src="./images/GitHub-Mark-32px.png" alt="GitHub Mark 32px" />
 The complete source code for this example is in the <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/tree/main/examples/015_simple_image">Coherence Operator GitHub</a> repository.</p>
</p>
</div>

<h3 id="_add_dependencies">Add Dependencies</h3>
<div class="section">
<p>To build a Coherence application there will obviously be at a minimum a dependency on <code>coherence.jar</code>.
Optionally we can also add dependencies on other Coherence modules.
In this example we&#8217;re going to add json support to the application by adding a dependency on <code>coherence-json</code>.</p>

<p>In the example we use the <code>coherence-bom</code> which ensures that we have consistent use of other Coherence modules.
In the <code>pom.xml</code> we have a <code>dependencyManagement</code> section.</p>

<markup
lang="xml"
title="pom.xml"
>    &lt;dependencyManagement&gt;
        &lt;dependencies&gt;
            &lt;dependency&gt;
                &lt;groupId&gt;com.oracle.coherence.ce&lt;/groupId&gt;
                &lt;artifactId&gt;coherence-bom&lt;/artifactId&gt;
                &lt;version&gt;${coherence.version}&lt;/version&gt;
                &lt;type&gt;pom&lt;/type&gt;
                &lt;scope&gt;import&lt;/scope&gt;
            &lt;/dependency&gt;
        &lt;/dependencies&gt;
    &lt;/dependencyManagement&gt;</markup>

<p>In the <code>build.gradle</code> file we add the bom as a platform dependency.</p>

<markup
lang="groovy"
title="build.gradle"
>dependencies {
    implementation platform("com.oracle.coherence.ce:coherence-bom:22.06.7")</markup>

<p>We can then add the <code>coherence</code> and <code>coherence-json</code> modules as dependencies</p>

<markup
lang="xml"
title="pom.xml"
>    &lt;dependencies&gt;
        &lt;dependency&gt;
            &lt;groupId&gt;com.oracle.coherence.ce&lt;/groupId&gt;
            &lt;artifactId&gt;coherence&lt;/artifactId&gt;
        &lt;/dependency&gt;
        &lt;dependency&gt;
            &lt;groupId&gt;com.oracle.coherence.ce&lt;/groupId&gt;
            &lt;artifactId&gt;coherence-json&lt;/artifactId&gt;
        &lt;/dependency&gt;
    &lt;/dependencies&gt;</markup>

<p>In the <code>build.gradle</code> file we add the bom as a platform dependency.</p>

<markup
lang="groovy"
title="build.gradle"
>dependencies {
    implementation platform("com.oracle.coherence.ce:coherence-bom:22.06.7")

    implementation "com.oracle.coherence.ce:coherence"
    implementation "com.oracle.coherence.ce:coherence-json"
}</markup>

</div>

<h3 id="_add_the_jib_plugin">Add the JIB Plugin</h3>
<div class="section">
<p>To build the image using JIB we need to add the JIB plugin to the project.</p>

<p>In the <code>pom.xml</code> file we add JIB to the <code>plugins</code> section.</p>

<markup
lang="xml"
title="pom.xml"
>    &lt;build&gt;
        &lt;plugins&gt;
            &lt;plugin&gt;
                &lt;groupId&gt;com.google.cloud.tools&lt;/groupId&gt;
                &lt;artifactId&gt;jib-maven-plugin&lt;/artifactId&gt;
                &lt;version&gt;3.4.0&lt;/version&gt;
            &lt;/plugin&gt;
        &lt;/plugins&gt;
    &lt;/build&gt;</markup>

<p>In the <code>build.gradle</code> file we add JIB to the <code>plugins</code> section.</p>

<markup
lang="groovy"
title="build.gradle"
>plugins {
    id 'java'
    id 'com.google.cloud.tools.jib' version '3.3.2'
}</markup>

</div>

<h3 id="_configure_the_jib_plugin">Configure the JIB Plugin</h3>
<div class="section">
<p>Now we can configure the JIB plugin with the properties specific to our image.
In this example the configuration is very simple, the JIB plugin documentation shows many more options.</p>

<p>We are going to set the following options:
* The name and tags for the image we will build.
* The main class that we will run as the entry point to the image - in this case <code>com.tangosol.net.Coherence</code>.
* The base image. In this example we will us a distroless Java 11 image. A distroless image is more secure as it contains nothing more than core linux and a JRE. There is no shell or other tools to introduce CVEs. The downside of this is that there is no shell, so you cannot exec into the running container, or use a shell script as an entry point. If you don;t need those things a distroless image is a great choice.</p>


<h4 id="_maven_configuration">Maven Configuration</h4>
<div class="section">
<p>In the <code>pom.xml</code> file we configure the plugin where it is declared in the <code>plugins</code> section:</p>

<markup
lang="xml"
title="pom.xml"
>&lt;plugin&gt;
    &lt;groupId&gt;com.google.cloud.tools&lt;/groupId&gt;
    &lt;artifactId&gt;jib-maven-plugin&lt;/artifactId&gt;
    &lt;version&gt;${version.plugin.jib}&lt;/version&gt;
    &lt;configuration&gt;
        &lt;from&gt;
            &lt;image&gt;gcr.io/distroless/java11-debian11&lt;/image&gt;    <span class="conum" data-value="1" />
        &lt;/from&gt;
        &lt;to&gt;
            &lt;image&gt;${project.artifactId}&lt;/image&gt;        <span class="conum" data-value="2" />
            &lt;tags&gt;
                &lt;tag&gt;${project.version}&lt;/tag&gt;           <span class="conum" data-value="3" />
                &lt;tag&gt;latest&lt;/tag&gt;
            &lt;/tags&gt;
        &lt;/to&gt;
        &lt;container&gt;
            &lt;mainClass&gt;com.tangosol.net.Coherence&lt;/mainClass&gt;  <span class="conum" data-value="4" />
            &lt;format&gt;OCI&lt;/format&gt;                               <span class="conum" data-value="5" />
        &lt;/container&gt;
    &lt;/configuration&gt;
&lt;/plugin&gt;</markup>

<ul class="colist">
<li data-value="1">The base image will be <code>gcr.io/distroless/java11-debian11</code></li>
<li data-value="2">The image name is set to the Maven module name using the property <code>${project.artifactId}</code></li>
<li data-value="3">There will be two tags for the image, <code>latest</code> and the project version taken from the <code>${project.version}</code> property.</li>
<li data-value="4">The main class to use when the image is run is set to <code>com.tangosol.net.Coherence</code></li>
<li data-value="5">The image type is set to <code>OCI</code></li>
</ul>
</div>

<h4 id="_gradle_configuration">Gradle Configuration</h4>
<div class="section">
<p>In the <code>build.gradle</code> file we configure JIB in the <code>jib</code> section:</p>

<markup
lang="groovy"
title="build.gradle"
>jib {
  from {
    image = 'gcr.io/distroless/java11-debian11'    <span class="conum" data-value="1" />
  }
  to {
    image = "${project.name}"              <span class="conum" data-value="2" />
    tags = ["${version}", 'latest']        <span class="conum" data-value="3" />
  }
  container {
    mainClass = 'com.tangosol.net.Coherence'  <span class="conum" data-value="4" />
    format = 'OCI'                            <span class="conum" data-value="5" />
  }
}</markup>

<ul class="colist">
<li data-value="1">The base image will be <code>gcr.io/distroless/java11-debian11</code></li>
<li data-value="2">The image name is set to the Maven module name using the property <code>${project.artifactId}</code></li>
<li data-value="3">There will be two tags for the image, <code>latest</code> and the project version taken from the <code>${project.version}</code> property.</li>
<li data-value="4">The main class to use when the image is run is set to <code>com.tangosol.net.Coherence</code></li>
<li data-value="5">The image type is set to <code>OCI</code></li>
</ul>
</div>
</div>

<h3 id="_build_the_image">Build the Image</h3>
<div class="section">
<p>To create the server image run the relevant commands as documented in the JIB plugin documentation.
In this case we&#8217;re going to build the image using Docker, although JIB offers other alternatives.</p>

<p>Using Maven we run:</p>

<markup
lang="bash"

>./mvnw compile jib:dockerBuild</markup>

<p>Using Gradle we run:</p>

<markup
lang="bash"

>./gradlew compileJava jibDockerBuild</markup>

<p>The command above will create an image named <code>simple-coherence</code> with two tags, <code>latest</code> and <code>1.0.0</code>.
Listing the local images should show the new images.</p>

<markup
lang="bash"

>$ docker images | grep simple
simple-coherence   1.0.0   1613cd3b894e   51 years ago  227MB
simple-coherence   latest  1613cd3b894e   51 years ago  227MB</markup>

</div>

<h3 id="_run_the_image">Run the Image</h3>
<div class="section">
<p>The image just built can be run using Docker (or your chosen container tool).
In this example we&#8217;ll run it interactively, just to prove it runs and starts Coherence.</p>

<markup
lang="bash"

>docker run -it --rm simple-coherence:latest</markup>

<p>The console output should display Coherence starting and finally show the Coherence service list, which will look something like this:</p>

<markup
lang="bash"

>Services
  (
  ClusterService{Name=Cluster, State=(SERVICE_STARTED, STATE_JOINED), Id=0, OldestMemberId=1}
  TransportService{Name=TransportService, State=(SERVICE_STARTED), Id=1, OldestMemberId=1}
  InvocationService{Name=Management, State=(SERVICE_STARTED), Id=2, OldestMemberId=1}
  PartitionedCache{Name=$SYS:Config, State=(SERVICE_STARTED), Id=3, OldestMemberId=1, LocalStorage=enabled, PartitionCount=257, BackupCount=1, AssignedPartitions=257, BackupPartitions=0, CoordinatorId=1}
  PartitionedCache{Name=PartitionedCache, State=(SERVICE_STARTED), Id=4, OldestMemberId=1, LocalStorage=enabled, PartitionCount=257, BackupCount=1, AssignedPartitions=257, BackupPartitions=0, CoordinatorId=1}
  PartitionedCache{Name=PartitionedTopic, State=(SERVICE_STARTED), Id=5, OldestMemberId=1, LocalStorage=enabled, PartitionCount=257, BackupCount=1, AssignedPartitions=257, BackupPartitions=0, CoordinatorId=1}
  ProxyService{Name=Proxy, State=(SERVICE_STARTED), Id=6, OldestMemberId=1}
  )</markup>

<p>Press <code>ctrl-C</code> to exit the container, the <code>--rm</code> option we used above wil automatically delete the stopped container.</p>

<p>We now have a simple Coherence image we can use in other examples and when trying out the Coherence Operator.</p>

</div>

<h3 id="_configuring_the_image_at_runtime">Configuring the Image at Runtime</h3>
<div class="section">
<p>With recent Coherence versions, Coherence configuration items that can be set using system properties prefixed with <code>coherence.</code> can also be set using environment variables. This makes it simple to set those properties when running containers because environment variables can be set from the commandline.</p>

<p>To set a property the system property name needs to be converted to an environment variable name.
This is done by converting the name to uppercase and replacing dots ('.') with underscores ('_').</p>

<p>For example, to set the cluster name we would set the <code>coherence.cluster</code> system property.
To run the image and set cluster name with an environment variable we convert <code>coherence.cluster</code> to <code>COHERENCE_CLUSTER</code> and run:</p>

<markup
lang="bash"

>docker run -it --rm -e COHERENCE_CLUSTER=my-cluster simple-coherence:latest</markup>

<p>This is much simpler than trying to change the Java commandline the image entrypoint uses.</p>

</div>
</div>
</doc-view>
