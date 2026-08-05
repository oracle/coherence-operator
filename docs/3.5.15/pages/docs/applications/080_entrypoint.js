<doc-view>

<h2 id="_run_an_image_entry_point">Run an Image Entry Point</h2>
<div class="section">
<p>The default behaviour of the Coherence operator is to configure the entry point and arguments to
use to run the Coherence container. This command line is created from the various configuration
elements in the <code>Coherence</code> resource yaml. Any entry point and arguments actually configured in
the image being used will be ignored.
The behaviour can be changed so that the images own entry point is used for the container.
This could be useful for example when an image contains a shell script that performs initialisation
before running the Java Coherence application.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Using an image entry point is only supported in images that use Java 11 or higher.</p>
</p>
</div>
<p>To use an image&#8217;s entry point set the <code>spec.application.useImageEntryPoint</code> field in the <code>Coherence</code>
resource to <code>true</code>.</p>

<p>For example:</p>

<markup

title="storage.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  replicas: 3
  image: container-registry.oracle.com/middleware/coherence-ce:14.1.2-0-1
  application:
    useImageEntryPoint: true</markup>


<h3 id="_how_are_the_jvm_and_coherence_configured">How are the JVM and Coherence Configured</h3>
<div class="section">
<p>When the operator builds the container command line it can pass all the required JVM options
and system properties to configure the application on the command line.
But, when the image entry point it being used the operator needs to pass configuration another way.</p>

<p>All the Coherence configuration system properties can also be passed as environment variables, so
the operator configures the container to have all the required environment variables to configure
Coherence. For example, the <code>coherence.role</code> system property is used to configure the role name
of a Coherence process, but Coherence will also use the <code>COHERENCE_ROLE</code> environment variable for this.
If <code>spec.role</code> value is set in the <code>Coherence</code> resource, this is be used to set <code>COHERENCE_ROLE</code>
environment variable in the Coherence container configuration in the Pod.</p>

<p>The operator then uses a combination of Java arguments files and the <code>JDK_JAVA_OPTIONS</code> environment
variable to configure the JVM. This means that most of the features of the <code>Coherence</code> CRD can be
used, even when running an image entry point.</p>


<h4 id="_java_argument_files">Java Argument Files</h4>
<div class="section">
<p>Various other environment variables are set by the Coherence operator to configure the container.
When the Pod starts an init-container that the operator has configured uses these environment
variables to produce a number of Java command line argument files.
These files contain all the JVM command line options that the operator would have used in its
custom command line if it was running the container.
For more information on argument files see the
<a id="" title="" target="_blank" href="https://docs.oracle.com/en/java/javase/17/docs/specs/man/java.html#java-command-line-argument-files">
Java Arguments Files</a> documentation.</p>

<p>The operator creates multiple arguments files for different purposes.
The Java argument files are always created by the init-container as these are used in the command line
that the operator normally configures for container.
There will be a file for the class path, a file for JVM options, a file for Spring Boot options
if the application is Spring Boot, etc.</p>

</div>
</div>

<h3 id="_the_jdk_java_options_environment_variable">The <code>JDK_JAVA_OPTIONS</code> Environment Variable</h3>
<div class="section">
<p>The <code>JDK_JAVA_OPTIONS</code> is a special environment variable recognised by the JVM.
Any values in the <code>JDK_JAVA_OPTIONS</code> environment variable are effectively prepended to the JVM
command line.
It is described fully in the
<a id="" title="" target="_blank" href="https://docs.oracle.com/en/java/javase/21/docs/specs/man/java.html#using-the-jdk_java_options-launcher-environment-variable">https://docs.oracle.com/en/java/javase/21/docs/specs/man/java.html#using-the-jdk_java_options-launcher-environment-variable</a>
[Java Command] documentation.</p>

<p>There are limitations on the size of the value for an environment variable, so the operator could
not specify all the options it needs in the <code>JDK_JAVA_OPTIONS</code> environment variable.
This is why the operator uses argument files instead, so all it needs to set into the <code>JDK_JAVA_OPTIONS</code> environment
variable are the names of the argument files to load.</p>


<h4 id="_what_if_the_application_already_sets_jdk_java_options">What If The Application Already Sets <code>JDK_JAVA_OPTIONS</code></h4>
<div class="section">
<p>If the <code>JDK_JAVA_OPTIONS</code> environment variable is set in the <code>Coherence</code> resource then the operator
will append its additional configuration onto the existing value.</p>

</div>

<h4 id="_disabling_use_of_jdk_java_options">Disabling Use of <code>JDK_JAVA_OPTIONS</code></h4>
<div class="section">
<p>There may be occasions that the operator should not configure the <code>JDK_JAVA_OPTIONS</code> environment variable.
For example, an image may run a shell script that runs various other Java commands before starting the
main Coherence application. If the <code>JDK_JAVA_OPTIONS</code> environment variable was set it would be applied
to all these Java processes too.</p>

<p>Setting the <code>spec.application.useJdkJavaOptions</code> field to <code>false</code> in the Coherence resource will
disable the use of the <code>JDK_JAVA_OPTIONS</code> environment variable and the operator will not set it.</p>

<p>For example,</p>

<markup

title="storage.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  replicas: 3
  image: container-registry.oracle.com/middleware/coherence-ce:14.1.2-0-1
  application:
    useImageEntryPoint: true
    useJdkJavaOptions: false</markup>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>When the  <code>spec.application.useJdkJavaOptions</code> field is set to false the operator has no way to pass
a number of configuration options to the JVM. Coherence configurations that are passed as environment
variables will still work. Anything passed as JVM options, such as memory configurations, system
properties, etc cannot be configured.</p>

<p>As long as the application that the image runs is a Coherence application correctly configured
to run in Kubernetes with the options required by the operator then it should still work.</p>
</p>
</div>
</div>

<h4 id="_using_an_alternative_to_jdk_java_options">Using An Alternative To <code>JDK_JAVA_OPTIONS</code></h4>
<div class="section">
<p>In use cases where the <code>JDK_JAVA_OPTIONS</code> environment variable cannot be used and is disabled as
described above, an alternative environment variable name can be specified that the operator will
configure instead. This allows an application to use this alternative environment variable at runtime
to obtain all the configurations that the Operator would have applied to the JVM.</p>

<p>The name of the alternative environment variable is set in the <code>spec.application.alternateJdkJavaOptions</code>
field of the <code>Coherence</code> resource.</p>

<p>For example, using the the yaml below will cause the operator to set the Java options values
into the <code>ALT_JAVA_OPTS</code> environment variable.</p>

<markup

title="storage.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  replicas: 3
  image: container-registry.oracle.com/middleware/coherence-ce:14.1.2-0-1
  application:
    useImageEntryPoint: true
    useJdkJavaOptions: false
    alternateJdkJavaOptions: "ALT_JAVA_OPTS"</markup>

<p>In the Coherence container the application code can then access the The <code>ALT_JAVA_OPTS</code> environment variable
to obtain the JVM options the Operator configured.</p>

</div>

<h4 id="_use_java_argument_files_directly">Use Java Argument Files Directly</h4>
<div class="section">
<p>In use cases where the <code>JDK_JAVA_OPTIONS</code> environment variable has been disabled application code
could also directly access the Java argument files the operator configured and use those to
configure the Coherence JVM.</p>

</div>
</div>
</div>
</doc-view>
