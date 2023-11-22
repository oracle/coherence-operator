<doc-view>

<h2 id="_example_coherence_image_using_a_dockerfile">Example Coherence Image using a Dockerfile</h2>
<div class="section">
<p>This example shows how to build a simple Coherence server image using a <code>Dockerfile</code>.
This image is built so that ot works out of the box with the Operator, with no additional configuration.
This is an alternative to the <router-link to="/examples/015_simple_image/README">Coherence Image using JIB</router-link> example.
There are many build tools and plugins for Maven and Gradle that are supposed to make building images easy.
Sometimes though, a simple <code>Dockerfile</code> approach is required.</p>

<p>A typical Coherence application image will still need to pull together various Coherence dependencies to add to the image.
This simple application does not actually contain any code, a real application would likely contain code and other resources.</p>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p><img src="./images/GitHub-Mark-32px.png" alt="GitHub Mark 32px" />
 The complete source code for this example is in the <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/tree/main/examples/016_simple_docker_image">Coherence Operator GitHub</a> repository.</p>
</p>
</div>
</div>

<h2 id="_the_dockerfile">The Dockerfile</h2>
<div class="section">
<p>The <code>Dockerfile</code> for the example is shown below:</p>

<markup

title="src/docker/Dockerfile"
>FROM gcr.io/distroless/java11-debian11

# Configure the image's health check command
# Health checks will only work with Coherence 22.06 and later
HEALTHCHECK  --start-period=10s --interval=30s \
    CMD ["java", \
    "-cp", "/app/libs/coherence.jar", \
    "com.tangosol.util.HealthCheckClient", \
    "http://127.0.0.1:6676/ready", \
    "||", "exit", "1"]

# Expose any default ports
# The default Coherence Extend port
EXPOSE 20000
# The default Coherence gRPC port
EXPOSE 1408
# The default Coherence metrics port
EXPOSE 9612
# The default Coherence health port
EXPOSE 6676

# Set the entry point to be the Java command to run
ENTRYPOINT ["java", "-cp", "/app/classes:/app/libs/*", "com.tangosol.net.Coherence"]

# Set any environment variables
# Set the health check port to a fixed value (corresponding to the command above)
ENV COHERENCE_HEALTH_HTTP_PORT=6676
# Fix the Extend Proxy to listen on port 20000
ENV COHERENCE_EXTEND_PORT=20000
# Enable Coherence metics
ENV COHERENCE_METRICS_HTTP_ENABLED=true
# Set the Coherence log level to debug logging
ENV COHERENCE_LOG_LEVEL=9
# Effectively disabled multicast cluster discovery, which does not work in containers
ENV COHERENCE_TTL=0

# Copy all the application files into the /app directory in the image
# This is the default structure supported by the Coherence Operator
COPY app app</markup>

<p><strong>Base Image</strong></p>

<p>The base image for this example is a distroless Java 11 image <code>gcr.io/distroless/java11-debian11</code></p>

<p><strong>Health Check</strong></p>

<p>The image is configured with a health check that uses the built-in Coherence health check on port 6676.</p>

<p><strong>Expose Ports</strong></p>

<p>A number of default Coherence ports are exposed.</p>

<p><strong>Entrypoint</strong></p>

<p>The image entry point will run <code>com.tangosol.net.Coherence</code> to run a Coherence storage enabled server.
The classpath is set to <code>/app/classes:/app/libs/*</code>. This is the same classpath that the JIB plugin would add artifacts to and is also supported out of the box by the Coherence operator.</p>

<p><strong>Environment Variables</strong></p>

<p>A number of environment variables are set to configure Coherence.
These values could have been set as system properties in the entry point, but using environment variables is a simpler option when running containers as they can easily be overridden at deploy time.</p>

<p><strong>Copy the Image Artifacts</strong></p>

<p>The Maven and Gradle build will copy all the classes and dependencies into a directory named <code>app/</code> in the same directory as the <code>Dockerfile</code>.
Using <code>COPY app app</code> will copy all the files into the image.</p>

</div>

<h2 id="_assemble_the_image_directory">Assemble the Image Directory</h2>
<div class="section">
<p>The next step is to assemble all the artifacts required to build the image.
Looking at the <code>Dockerfile</code> above, this means copying any dependencies and other files into a directory named <code>app/</code> in the same directory that the <code>Dockerfile</code> is in.
This example contains both a Maven <code>pom.xml</code> file and Gradle build files, that show how to use these tools to gather all the files required for the image.</p>

<p>There are other build tools such as <code>make</code> or <code>ant</code> or just plain scripts, but as the task involves pulling together all the Coherence jar files from Maven central, it is simplest to use Maven or Gradle.</p>

<p>To build a Coherence application there will obviously be at a minimum a dependency on <code>coherence.jar</code>.
Optionally we can also add dependencies on other Coherence modules and other dependencies, for example Coherence coul dbe configured to use SLF4J for logging.
In this example we&#8217;re going to add json support to the application by adding a dependency on <code>coherence-json</code> and <code>coherence-grpc-proxy</code>.</p>

<p>Jump to the relevant section, depending on the build tool being used:</p>

<ul class="ulist">
<li>
<p><router-link to="#maven" @click.native="this.scrollFix('#maven')">Using Maven</router-link></p>

</li>
<li>
<p><router-link to="#gradle" @click.native="this.scrollFix('#gradle')">Using Gradle</router-link></p>

</li>
</ul>

<h3 id="maven">Using Maven</h3>
<div class="section">
<p>To assemble the image artifacts using Maven, everything is configured in the Maven <code>pom.xml</code> file.
The Maven build will pull all the artifacts required in the image, including the <code>Dockerfile</code> into a directory under <code>target\docker</code>.</p>


<h4 id="_adding_dependencies">Adding Dependencies</h4>
<div class="section">
<p>In the example the <code>coherence-bom</code> is added to the <code>&lt;dependencyManagement&gt;</code> section as an import, to ensure consistent versioning of other Coherence modules.</p>

<p>In the <code>pom.xml</code> we have a <code>dependencyManagement</code> section.</p>

<markup
lang="xml"
title="pom.xml"
>&lt;dependencyManagement&gt;
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

<p>We can then add the <code>coherence</code> <code>coherence-json</code> and <code>coherence-grpc-proxy</code> modules as dependencies</p>

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
        &lt;dependency&gt;
            &lt;groupId&gt;com.oracle.coherence.ce&lt;/groupId&gt;
            &lt;artifactId&gt;coherence-grpc-proxy&lt;/artifactId&gt;
        &lt;/dependency&gt;
    &lt;/dependencies&gt;</markup>

</div>

<h4 id="_assembling_the_image_artifacts">Assembling the Image Artifacts</h4>
<div class="section">
<p>This example will use the Maven Assembly Plugin to gather all the dependencies and other files together into the <code>target/docker</code> directory. The assembly plugin is configured in the <code>pom.xml</code> file.</p>

<p>The assembly plugin is configured to use the <code>src/assembly/image-assembly.xml</code> descriptor file to determine what to assemble. The <code>&lt;finalName&gt;</code> configuration element is set to <code>docker</code> so all the files will be assembled into a directory named <code>docker/</code> under the <code>target/</code> directory.
The assembly plugin execution is bound to the <code>package</code> build phase.</p>

<markup
lang="xml"

>&lt;plugin&gt;
    &lt;groupId&gt;org.apache.maven.plugins&lt;/groupId&gt;
    &lt;artifactId&gt;maven-assembly-plugin&lt;/artifactId&gt;
    &lt;version&gt;${maven.assembly.plugin.version}&lt;/version&gt;
    &lt;executions&gt;
        &lt;execution&gt;
            &lt;id&gt;prepare-image&lt;/id&gt;
            &lt;phase&gt;package&lt;/phase&gt;
            &lt;goals&gt;
                &lt;goal&gt;single&lt;/goal&gt;
            &lt;/goals&gt;
            &lt;configuration&gt;
                &lt;finalName&gt;docker&lt;/finalName&gt;
                &lt;appendAssemblyId&gt;false&lt;/appendAssemblyId&gt;
                &lt;descriptors&gt;
                    &lt;descriptor&gt;${project.basedir}/src/assembly/image-assembly.xml&lt;/descriptor&gt;
                &lt;/descriptors&gt;
                &lt;attach&gt;false&lt;/attach&gt;
            &lt;/configuration&gt;
        &lt;/execution&gt;
    &lt;/executions&gt;
&lt;/plugin&gt;</markup>

<p>The <code>image-assembly.xml</code> descriptor file is shown below, and configures the following:</p>

<ul class="ulist">
<li>
<p>The <code>&lt;format&gt;dir&lt;/format&gt;</code> element tells the assembly plugin to assemble all the artifacts into a directory.</p>

</li>
<li>
<p>There are two <code>&lt;fileSets&gt;</code> configured:</p>
<ul class="ulist">
<li>
<p>The first copies any class files in <code>target/classes</code> to <code>app/classes</code> (which will actually be <code>target/docker/app/classes</code>)</p>

</li>
<li>
<p>The second copies all files under <code>src/docker</code> (i.e. the <code>Dockerfile</code>) into <code>target/docker</code></p>

</li>
</ul>
</li>
<li>
<p>The <code>&lt;dependencySets&gt;</code> configuration copies all the project dependencies (including transitive dependencies) to the <code>app/libs</code> directory (actually the <code>target/docker/app/libs</code> directory). Any version information will be stripped from the files, so <code>coherence-22.06.6.jar</code> would become <code>coherence.jar</code>.</p>

</li>
</ul>
<markup
lang="xml"
title="src/assembly/image-assembly.xml"
>&lt;assembly xmlns="http://maven.apache.org/ASSEMBLY/2.1.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xsi:schemaLocation="http://maven.apache.org/ASSEMBLY/2.1.0 http://maven.apache.org/xsd/assembly-2.1.0.xsd"&gt;
  &lt;id&gt;image&lt;/id&gt;
  &lt;formats&gt;
    &lt;format&gt;dir&lt;/format&gt;
  &lt;/formats&gt;

  &lt;includeBaseDirectory&gt;false&lt;/includeBaseDirectory&gt;

  &lt;fileSets&gt;
    &lt;!-- copy the module's compiled classes --&gt;
    &lt;fileSet&gt;
      &lt;directory&gt;target/classes&lt;/directory&gt;
      &lt;outputDirectory&gt;app/classes&lt;/outputDirectory&gt;
      &lt;fileMode&gt;755&lt;/fileMode&gt;
      &lt;filtered&gt;false&lt;/filtered&gt;
    &lt;/fileSet&gt;
    &lt;!-- copy the Dockerfile --&gt;
    &lt;fileSet&gt;
      &lt;directory&gt;${project.basedir}/src/docker&lt;/directory&gt;
      &lt;outputDirectory/&gt;
      &lt;fileMode&gt;755&lt;/fileMode&gt;
    &lt;/fileSet&gt;
  &lt;/fileSets&gt;

  &lt;!-- copy the application dependencies --&gt;
  &lt;dependencySets&gt;
    &lt;dependencySet&gt;
      &lt;outputDirectory&gt;app/libs&lt;/outputDirectory&gt;
      &lt;directoryMode&gt;755&lt;/directoryMode&gt;
      &lt;fileMode&gt;755&lt;/fileMode&gt;
      &lt;unpack&gt;false&lt;/unpack&gt;
      &lt;useProjectArtifact&gt;false&lt;/useProjectArtifact&gt;
      &lt;!-- strip the version from the jar files --&gt;
      &lt;outputFileNameMapping&gt;${artifact.artifactId}${dashClassifier?}.${artifact.extension}&lt;/outputFileNameMapping&gt;
    &lt;/dependencySet&gt;
  &lt;/dependencySets&gt;
&lt;/assembly&gt;</markup>

<p>Running the following command will pull all the required image artifacts and <code>Dockerfile</code> into the <code>target/docker</code> directory:</p>

<markup
lang="bash"

>./mvnw package</markup>

</div>
</div>

<h3 id="gradle">Using Gradle</h3>
<div class="section">
<p>To assemble the image artifacts using Maven, everything is configured in the Maven <code>build.gradle</code> file.
The Gradle build will pull all the artifacts required in the image, including the <code>Dockerfile</code> into a directory under <code>build\docker</code>.</p>


<h4 id="_adding_dependencies_2">Adding Dependencies</h4>
<div class="section">
<p>In the example the <code>coherence-bom</code> is added to the <code>&lt;dependencyManagement&gt;</code> section as an import, to ensure consistent versioning of other Coherence modules.</p>

<p>In the <code>build.gradle</code> file we add the bom as a platform dependency and then add dependencies on <code>coherence</code> and <code>coherence-json</code>.</p>

<markup
lang="groovy"
title="build.gradle"
>dependencies {
    implementation platform("com.oracle.coherence.ce:coherence-bom:22.06.6")

    implementation "com.oracle.coherence.ce:coherence"
    implementation "com.oracle.coherence.ce:coherence-json"
    implementation "com.oracle.coherence.ce:coherence-grpc-proxy"
}</markup>

</div>

<h4 id="_assembling_the_image_artifacts_2">Assembling the Image Artifacts</h4>
<div class="section">
<p>To assemble all the image artifacts into the <code>build/docker</code> directory, the Gradle copy task can be used.
There will be multiple copy tasks to copy each type of artifact, the dependencies, any compile classes, and the <code>Dockerfile</code>.</p>

<p>The following task named <code>copyDependencies</code> is added to <code>build.gradle</code> to copy the dependencies. This task has additional configuration to rename the jar files to strip off any version.</p>

<markup
lang="groovy"
title="build.gradle"
>task copyDependencies(type: Copy) {
    from configurations.runtimeClasspath
    into "$buildDir/docker/app/libs"
    configurations.runtimeClasspath.resolvedConfiguration.resolvedArtifacts.each {
        rename "${it.artifact.name}-${it.artifactId.componentIdentifier.version}", "${it.artifact.name}"
    }
}</markup>

<p>The following task named <code>copyClasses</code> copies any compiled classes (although this example does not actually have any).</p>

<markup
lang="groovy"
title="build.gradle"
>task copyClasses(type: Copy) {
    dependsOn classes
    from "$buildDir/classes/java/main"
    into "$buildDir/docker/app/classes"
}</markup>

<p>The final copy task named <code>copyDocker</code> copies the contents of the <code>src/docker</code> directory:</p>

<markup
lang="groovy"
title="build.gradle"
>task copyDocker(type: Copy) {
    from "src/docker"
    into "$buildDir/docker"
}</markup>

<p>To be able to run the image assembly as a single command, an empty task named `` is created that depends on all the copy tasks.</p>

<p>Running the following command will pull all the required image artifacts and <code>Dockerfile</code> into the <code>build/docker</code> directory:</p>

<markup
lang="bash"

>./gradlew assembleImage</markup>

</div>
</div>
</div>

<h2 id="_build_the_image">Build the Image</h2>
<div class="section">
<p>After running the Maven or Gradle commands to assemble the image artifacts, Docker can be used to actually build the image from the relevant <code>docker/</code> directory.</p>

<p>Using Maven:</p>

<markup
lang="bash"

>cd target/docker
docker build -t simple-coherence-server:1.0.0 .</markup>

<p>Using Gradle:</p>

<markup
lang="bash"

>cd build/docker
docker build -t simple-coherence-server:1.0.0 .</markup>

<p>The command above will create an image named <code>simple-coherence-server:1.0.0</code>.
Listing the local images should show the new images, similar to the output below:</p>

<markup
lang="bash"

>$ docker images | grep simple
simple-coherence-server   1.0.0   1613cd3b894e   51 years ago  227MB</markup>

</div>
</doc-view>
