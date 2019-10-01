<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Configure Applications</dt>
<dd slot="desc"><p>Whilst the Coherence Operator can manage plain Coherence clusters typically custom application code and configuration
files would be added to a Coherence JVM&#8217;s classpath.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_configure_applications">Configure Applications</h2>
<div class="section">
<p>Different application code and configuration can be added to different roles in a <code>CoherenceCluster</code> by specifying
the application&#8217;s configuration in the <code>application</code> section of a role spec. There are a number of different fields
that can be configured for an application described below. All of the configuration described below is optional.</p>

<markup
lang="yaml"
title="Application Spec"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  application:
    image: acme/orders-data:1.0.0
    imagePullPolicy: Always
    appDir: "/app"
    libDir: "/lib"
    configDir: "/conf"
    main: io.acme.Server
    args:
      - "foo"
      - "bar"</markup>

<ul class="ulist">
<li>
<p><router-link to="#app-image" @click.native="this.scrollFix('#app-image')">Setting the Application Image</router-link> - The application&#8217;s image that is used to provide application
<code>.jar</code> files and configuration files to add to the JVM classpath.</p>

</li>
<li>
<p><router-link to="#pull-policy" @click.native="this.scrollFix('#pull-policy')">Setting the Application Image Pull Policy</router-link> - The pull policy that Kubernetes will use to pull
the application&#8217;s image</p>

</li>
<li>
<p><router-link to="#app-dir" @click.native="this.scrollFix('#app-dir')">Setting the Application Directory</router-link> - The name of the folder in the application&#8217;s image containing
application artifacts. This will be the working directory for the Coherence container.</p>

</li>
<li>
<p><router-link to="#app-lib" @click.native="this.scrollFix('#app-lib')">Setting the Application Lib Directory</router-link> - The name of the folder in the application&#8217;s image containing
<code>.jar</code> files to add to the JVM class path</p>

</li>
<li>
<p><router-link to="#app-conf" @click.native="this.scrollFix('#app-conf')">Setting the Application Config Directory</router-link> - The name of the folder in the application&#8217;s image containing
configuration files that will be add to the JVM classpath</p>

</li>
<li>
<p><router-link to="#app-main" @click.native="this.scrollFix('#app-main')">Setting the Application Main Class</router-link> - The application&#8217;s custom main main Class to use if running a
class other than Coherence <code>DefaultCacheServer</code></p>

</li>
<li>
<p><router-link to="#app-args" @click.native="this.scrollFix('#app-args')">Setting the Application Main Class Arguments</router-link> - The arguments to pass to the application&#8217;s main <code>Class</code></p>

</li>
</ul>
</div>

<h2 id="app-image">Setting the Application Image</h2>
<div class="section">
<p>The application image is the Docker image containing the <code>.jar</code> files and configuration files of the Coherence application
that should be deployed in the Coherence cluster. For more information see the
<router-link to="/guides/030_applications">Deploying Coherence Applications</router-link> guide.</p>

<p>Whilst the Coherence Operator makes it simple to deploy and manage a Coherence cluster in Kubernetes in the majority of
use cases there will be a requirement for application code to be deployed and run in the Coherence JVMs. This application
code and any application configuration files are supplied as a separate image. This image is loaded as an init-container
by the Coherence <code>Pods</code> and the relevant <code>.jar</code> files and configuration files from this image are added to the Coherence
JVM&#8217;s classpath.</p>

<p>As well as setting the image name it is also sometimes useful to set the application image&#8217;s  <router-link to="#pull-policy" @click.native="this.scrollFix('#pull-policy')">image pull policy</router-link>.</p>

<p>If the image being configured is from a registry requiring authentication see the section
on <router-link to="/clusters/200_private_repos">pulling from private registries</router-link>.</p>


<h3 id="_setting_the_application_image_for_the_implicit_role">Setting the Application Image for the Implicit Role</h3>
<div class="section">
<p>When using the implicit role configuration the application image to use is set directly in the <code>CoherenceCluster</code> <code>spec</code>
<code>application.image</code> field.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  application:
    image: acme/orders-data:1.0.0  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>acme/orders-data:1.0.0</code> will be used to add additional <code>.jar</code> files and configuration files to the classpath of
the Coherence container in the implicit <code>storage</code> role&#8217;s <code>Pods</code></li>
</ul>
</div>

<h3 id="_setting_the_application_image_for_explicit_roles">Setting the Application Image for Explicit Roles</h3>
<div class="section">
<p>When using the explicit roles in a <code>CoherenceCluster</code> <code>roles</code> list the application image to use is set for each role.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data
      application:
        image: acme/orders-data:1.0.0  <span class="conum" data-value="1" />
    - role: proxy
      application:
        image: acme/orders-proxy/1.0.0  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role <code>Pods</code> will use the application image <code>acme/orders-data:1.0.0</code></li>
<li data-value="2">The <code>proxy</code> role <code>Pods</code> will use the application image <code>acme/orders-proxy/1.0.0</code></li>
</ul>
</div>

<h3 id="_setting_the_application_image_for_explicit_roles_with_a_default">Setting the Application Image for Explicit Roles with a Default</h3>
<div class="section">
<p>When using the explicit roles in a <code>CoherenceCluster</code> <code>roles</code> list the application image to use can be set in the
<code>CoherenceCluster</code> <code>spec</code> section and will apply to all roles unless specifically overridden for a <code>role</code> in the
<code>roles</code> list.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  application:
    image: acme/orders-data:1.0.0  <span class="conum" data-value="1" />
  roles:
    - role: data
    - role: proxy</markup>

<ul class="colist">
<li data-value="1">The <code>data</code> and the <code>proxy</code> roles will both use the application image <code>acme/orders-data:1.0.0</code></li>
</ul>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  application:
    image: acme/orders-data:1.0.0  <span class="conum" data-value="1" />
  roles:
    - role: data
    - role: proxy
    - role: web
      application:
        image: acme/orders-front-end/1.0.0  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> and the <code>proxy</code> roles will both use the application image <code>acme/orders-data:1.0.0</code></li>
<li data-value="2">The <code>web</code> role will use the application image <code>acme/orders-web/1.0.0</code></li>
</ul>
</div>
</div>

<h2 id="pull-policy">Setting the Application Image Pull Policy</h2>
<div class="section">
<p>The image pull policy controls when (and if) Kubernetes will pull the application image onto the node where the Coherence
<code>Pods</code> are being schedules.
See <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/containers/images/#updating-images">Kubernetes imagePullPolicy</a> for more information.</p>

<div class="admonition note">
<p class="admonition-inline">The Kubernetes default pull policy is <code>IfNotPresent</code> unless the image tag is <code>:latest</code> in which case the default
policy is <code>Always</code>. The <code>IfNotPresent</code> policy causes the Kubelet to skip pulling an image if it already exists.
Note that you should avoid using the <code>:latest</code> tag, see
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/configuration/overview/#container-images">Kubernetes Best Practices for Configuration</a>
for more information.</p>
</div>
<p>The application image&#8217;s pull policy is set using the <code>imagePullPolicy</code> field in the <code>spec.application</code> section.</p>


<h3 id="_setting_the_image_pull_policy_for_the_implicit_role">Setting the Image Pull Policy for the Implicit Role</h3>
<div class="section">
<p>To set the <code>imagePullPolicy</code> for the implicit role:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  application:
    image: acme/orders-data:1.0.0
    imagePullPolicy: Always <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The image pull policy for the implicit role above has been set to <code>Always</code></li>
</ul>
</div>

<h3 id="_setting_the_image_pull_policy_for_explicit_roles">Setting the Image Pull Policy for Explicit Roles</h3>
<div class="section">
<p>To set the <code>imagePullPolicy</code> for the explicit roles in the <code>roles</code> list:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data
      application:
        image: acme/orders-data:1.0.0
        imagePullPolicy: Always <span class="conum" data-value="1" />
    - role: proxy
      application:
        image: acme/orders-proxy/1.0.0
        imagePullPolicy: IfNotPresent <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The image pull policy for the <code>data</code> role has been set to <code>Always</code></li>
<li data-value="2">The image pull policy for the <code>proxy</code> role above has been set to <code>IfNotPresent</code></li>
</ul>
</div>

<h3 id="_setting_the_image_pull_policy_for_explicit_roles_with_default">Setting the Image Pull Policy for Explicit Roles with Default</h3>
<div class="section">
<p>To set the <code>imagePullPolicy</code> for the explicit roles with a default value:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  application:
    imagePullPolicy: Always <span class="conum" data-value="1" />
  roles:
    - role: data
      application:
        image: acme/orders-data:1.0.0
    - role: proxy
      application:
        image: acme/orders-proxy/1.0.1
    - role: web
      application:
        image: acme/orders-front-end/1.0.1
        imagePullPolicy: IfNotPresent <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The default image pull policy is set to <code>Always</code>. The <code>data</code> and <code>proxy</code> roles will use the default value because
they do not specifically set the value in their specs.</li>
<li data-value="2">The image pull policy for the <code>web</code> role above has been set to <code>IfNotPresent</code></li>
</ul>
</div>
</div>

<h2 id="app-dir">Setting the Application Directory</h2>
<div class="section">

<h3 id="_setting_the_application_directory_for_the_implicit_role">Setting the Application Directory for the Implicit Role</h3>
<div class="section">

</div>

<h3 id="_setting_the_application_directory_for_explicit_roles">Setting the Application Directory for Explicit Roles</h3>
<div class="section">

</div>

<h3 id="_setting_the_application_directory_for_explicit_roles_with_a_default">Setting the Application Directory for Explicit Roles with a Default</h3>
<div class="section">

</div>
</div>

<h2 id="app-lib">Setting the Application Lib Directory</h2>
<div class="section">
<p>A typical Coherence application may also require additional dependencies (usually <code>.jar</code> files) that need to be added
to the classpath.
The applications&#8217;s lib directory is a directory in the application&#8217;s image that contains these additional <code>.jar</code> files.
The Coherence Operator will add the files to the classpath with the wildcard setting (e.g. <code>-cp /lib/*</code>) it does not add
each file in the lib directory individually to the classpath. This means that the contents of the lib directory are
added to the classpath using the rules that the JVM uses to process wild card classpath entries.</p>


<h3 id="_setting_the_application_lib_directory_for_the_implicit_role">Setting the Application Lib Directory for the Implicit Role</h3>
<div class="section">

</div>

<h3 id="_setting_the_application_lib_directory_for_explicit_roles">Setting the Application Lib Directory for Explicit Roles</h3>
<div class="section">

</div>

<h3 id="_setting_the_application_lib_directory_for_explicit_roles_with_a_default">Setting the Application Lib Directory for Explicit Roles with a Default</h3>
<div class="section">

</div>
</div>

<h2 id="app-conf">Setting the Application Config Directory</h2>
<div class="section">
<p>A Coherence application may require additional files added to the classpath such as configuration files and other
resources. These additional files can be placed into the config directory of the application&#8217;s image and this directory
added to the classpath of the Coherence JVM. Just the directory is added to the classpath (e.g. <code>-cp /conf</code>) the contents
themselves are not added.</p>


<h3 id="_setting_the_application_config_directory_for_the_implicit_role">Setting the Application Config Directory for the Implicit Role</h3>
<div class="section">

</div>

<h3 id="_setting_the_application_config_directory_for_explicit_roles">Setting the Application Config Directory for Explicit Roles</h3>
<div class="section">

</div>

<h3 id="_setting_the_application_config_directory_for_explicit_roles_with_a_default">Setting the Application Config Directory for Explicit Roles with a Default</h3>
<div class="section">

</div>
</div>

<h2 id="app-main">Setting the Application Main Class</h2>
<div class="section">
<p>By default the Coherence container will run the <code>main</code> method in the <code>com.tangosol.coherence.DefaultCacheServer</code>
class. Sometimes an application requires a different class as the main class and this can be configured for roles
in a <code>CoherenceCluster</code>.</p>


<h3 id="_setting_the_application_main_class_for_the_implicit_role">Setting the Application Main Class for the Implicit Role</h3>
<div class="section">

</div>

<h3 id="_setting_the_application_main_class_for_explicit_roles">Setting the Application Main Class for Explicit Roles</h3>
<div class="section">

</div>

<h3 id="_setting_the_application_main_class_for_explicit_roles_with_a_default">Setting the Application Main Class for Explicit Roles with a Default</h3>
<div class="section">

</div>
</div>

<h2 id="app-args">Setting the Application Main Class Arguments</h2>
<div class="section">
<p>If a custom main class is being used it is sometimes required to pass arguments to that class&#8217;s <code>main</code> method.
These additional arguments can be configured alongside the main class.</p>


<h3 id="_setting_the_application_main_class_arguments_for_the_implicit_role">Setting the Application Main Class Arguments for the Implicit Role</h3>
<div class="section">

</div>

<h3 id="_setting_the_application_main_class_arguments_for_explicit_roles">Setting the Application Main Class Arguments for Explicit Roles</h3>
<div class="section">

</div>

<h3 id="_setting_the_application_main_class_arguments_for_explicit_roles_with_a_default">Setting the Application Main Class Arguments for Explicit Roles with a Default</h3>
<div class="section">

</div>
</div>
</doc-view>
