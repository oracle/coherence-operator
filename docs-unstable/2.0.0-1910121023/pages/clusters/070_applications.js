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

<p>See the in-depth guide on <router-link to="#app-deployments/010_overview.adoc" @click.native="this.scrollFix('#app-deployments/010_overview.adoc')">Coherence application deployments</router-link> for more details on
creating and deploying <code>CoherenceClusters</code> with custom application code.</p>

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
    libDir: "/app/lib"
    configDir: "/app/conf"
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
<v-divider class="my-5"/>
</div>

<h2 id="app-image">Setting the Application Image</h2>
<div class="section">
<p>The application image is the Docker image containing the <code>.jar</code> files and configuration files of the Coherence application
that should be deployed in the Coherence cluster. For more information see the
<router-link to="#app-deployments/010_overview.adoc" @click.native="this.scrollFix('#app-deployments/010_overview.adoc')">Coherence application deployments</router-link> guide.</p>

<p>Whilst the Coherence Operator makes it simple to deploy and manage a Coherence cluster in Kubernetes in the majority of
use cases there will be a requirement for application code to be deployed and run in the Coherence JVMs. This application
code and any application configuration files are supplied as a separate image. This image is loaded as an init-container
by the Coherence <code>Pods</code> and the relevant <code>.jar</code> files and configuration files from this image are added to the Coherence
JVM&#8217;s classpath.</p>

<p>As well as setting the image name it is also sometimes useful to set the application image&#8217;s <router-link to="#pull-policy" @click.native="this.scrollFix('#pull-policy')">image pull policy</router-link>.</p>

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
        image: acme/orders-data:1.0.0   <span class="conum" data-value="1" />
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
    image: acme/orders-data:1.0.0           <span class="conum" data-value="1" />
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
<v-divider class="my-5"/>
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
        imagePullPolicy: Always        <span class="conum" data-value="1" />
    - role: proxy
      application:
        image: acme/orders-proxy/1.0.0
        imagePullPolicy: IfNotPresent  <span class="conum" data-value="2" /></markup>

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
    imagePullPolicy: Always                 <span class="conum" data-value="1" />
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
        imagePullPolicy: IfNotPresent       <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The default image pull policy is set to <code>Always</code>. The <code>data</code> and <code>proxy</code> roles will use the default value because
they do not specifically set the value in their specs.</li>
<li data-value="2">The image pull policy for the <code>web</code> role above has been set to <code>IfNotPresent</code></li>
</ul>
<v-divider class="my-5"/>
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

<p>The lib directory is set in the <code>application.libDir</code> field. This field is optional and if not specified the default
directory name used will be <code>/app/lib</code>.</p>


<h3 id="_setting_the_application_lib_directory_for_the_implicit_role">Setting the Application Lib Directory for the Implicit Role</h3>
<div class="section">
<p>When configuring a <code>CoherenceCluster</code> with a single implicit role the application&#8217;s lib directory is specified in the
<code>application.libDir</code> field.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  application:
    libDir: /lib  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The application image contains a directory named <code>/app-lib</code> that contains the <code>.jar</code> files to add to the JVM
classpath.</li>
</ul>
</div>

<h3 id="_setting_the_application_lib_directory_for_explicit_roles">Setting the Application Lib Directory for Explicit Roles</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with explicit roles in the <code>roles</code> list the <code>application.libDir</code> field can be set
specifically for each role:</p>

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
        libDir: app-lib    <span class="conum" data-value="1" />
    - role: proxy
      application:
        libDir: proxy-lib  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The application image contains a directory named <code>/app-lib</code> that contains the <code>.jar</code> files to add to the JVM
classpath in all of the <code>Pods</code> for the <code>data</code> role.</li>
<li data-value="2">The application image contains a directory named <code>/proxy-lib</code> that contains the <code>.jar</code> files to add to the JVM
classpath in all of the <code>Pods</code> for the <code>proxy</code> role.</li>
</ul>
</div>

<h3 id="_setting_the_application_lib_directory_for_explicit_roles_with_a_default">Setting the Application Lib Directory for Explicit Roles with a Default</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with explicit roles in the <code>roles</code> list the <code>application.libDir</code> field can be set
at the <code>spec</code> level as a default that applies to all of the roles in the <code>roles</code> list unless specifically overridden
for an individual role:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  application:
    libDir: app-lib        <span class="conum" data-value="1" />
  roles:
    - role: data           <span class="conum" data-value="2" />
    - role: proxy
      application:
        libDir: proxy-lib  <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The default value for the <code>libDir</code> for all roles will be <code>/app-lib</code> unless overridden for a specific role.</li>
<li data-value="2">The <code>data</code> role does not specify a value for <code>libDir</code> so it will use the default <code>app-lib</code>. The application image
should contain a directory named <code>/app-lib</code> that contains the <code>.jar</code> files to add to the JVM classpath in all of the
<code>Pods</code> for the <code>data</code> role.</li>
<li data-value="3">The <code>proxy</code> role has an explicit value set for the <code>libDir</code> field. The application image should a directory named
<code>/proxy-lib</code> that contains the <code>.jar</code> files to add to the JVM classpath in all of the <code>Pods</code> for the <code>proxy</code> role.</li>
</ul>
<v-divider class="my-5"/>
</div>
</div>

<h2 id="app-conf">Setting the Application Config Directory</h2>
<div class="section">
<p>A Coherence application may require additional files added to the classpath such as configuration files and other
resources. These additional files can be placed into the config directory of the application&#8217;s image and this directory
added to the classpath of the Coherence JVM. Just the directory is added to the classpath (e.g. <code>-cp /conf</code>) the contents
themselves are not added.</p>

<p>The configuration directory is set in the <code>application.configDir</code> field. This field is optional and if not specified
the default directory name used will be <code>/app/conf</code>.</p>


<h3 id="_setting_the_application_config_directory_for_the_implicit_role">Setting the Application Config Directory for the Implicit Role</h3>
<div class="section">
<p>When configuring a <code>CoherenceCluster</code> with a single implicit role the application&#8217;s configuration directory is specified
in the <code>application.configDir</code> field.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  application:
    configDir: app-conf  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The application image contains a directory named <code>/app-conf</code> that contains any configuration files to add to the JVM
classpath.</li>
</ul>
</div>

<h3 id="_setting_the_application_config_directory_for_explicit_roles">Setting the Application Config Directory for Explicit Roles</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with explicit roles in the <code>roles</code> list the <code>application.configDir</code> field can be set
specifically for each role:</p>

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
        configDir: app-conf    <span class="conum" data-value="1" />
    - role: proxy
      application:
        configDir: proxy-conf  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The application image contains a directory named <code>/app-conf</code> that contains the configuration files to add to the JVM
classpath in all of the <code>Pods</code> for the <code>data</code> role.</li>
<li data-value="2">The application image contains a directory named <code>/proxy-conf</code> that contains the configuration files to add to the
JVM classpath in all of the <code>Pods</code> for the <code>proxy</code> role.</li>
</ul>
</div>

<h3 id="_setting_the_application_config_directory_for_explicit_roles_with_a_default">Setting the Application Config Directory for Explicit Roles with a Default</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with explicit roles in the <code>roles</code> list the <code>application.configDir</code> field can be set
at the <code>spec</code> level as a default that applies to all of the roles in the <code>roles</code> list unless specifically overridden
for an individual role:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  application:
    configDir: app-conf       <span class="conum" data-value="1" />
  roles:
    - role: data              <span class="conum" data-value="2" />
    - role: proxy
      application:
        configDir: proxy-conf <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The default value for the <code>configDir</code> field is <code>app-conf/</code> which will be used for all roles unless specifically
overridden for a role.</li>
<li data-value="2">The <code>data</code> role does not specify a value for <code>configDir</code> so it will use the default. The application image should
contain a directory named <code>/app-conf</code> that contains the configuration files to add to the JVM classpath in all of the
<code>Pods</code> for the <code>data</code> role.</li>
<li data-value="3">The <code>proxy</code> role has an explicit value set for the <code>configDir</code> field. The application image should a directory named
<code>/proxy-conf</code> that contains the configuration files to add to the JVM classpath in all of the <code>Pods</code> for the <code>proxy</code>
role.</li>
</ul>
<v-divider class="my-5"/>
</div>
</div>

<h2 id="app-dir">Setting the Application Directory</h2>
<div class="section">
<p>Sometimes an application may have more than just <code>.jar</code> files or configuration files in the <code>conf</code> folder.
An application may have a number of artifacts that it needs to access from a working directory so for this use case
an application directory can be specified that will effectively become the working directory for the Coherence JVM
in the <code>Pods</code>. The application directory may be a parent directory of the lib or configuration directory or they may
be separate directory trees.</p>

<p>The application directory is set in the <code>application.appDir</code> field. This field is optional and if not specified
the default directory name used will be <code>/app</code>.</p>


<h3 id="_setting_the_application_directory_for_the_implicit_role">Setting the Application Directory for the Implicit Role</h3>
<div class="section">
<p>When configuring a <code>CoherenceCluster</code> with a single implicit role the application&#8217;s working directory is specified
in the <code>spec.application.appDir</code> field.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  application:
    appDir: app  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The application image contains a directory named <code>/app</code> that will effectively become the working directory for
the Coherence JVM in the <code>Pods</code> for the role.</li>
</ul>
</div>

<h3 id="_setting_the_application_directory_for_explicit_roles">Setting the Application Directory for Explicit Roles</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with explicit roles in the <code>roles</code> list the <code>application.appDir</code> field can be set
specifically for each role:</p>

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
        appDir: data-app   <span class="conum" data-value="1" />
    - role: proxy
      application:
        appDir: proxy-app  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The application image contains a directory named <code>/data-app</code> that will effectively become the working directory for
the Coherence JVM in the <code>Pods</code> for the <code>data</code> role.</li>
<li data-value="2">The application image contains a directory named <code>/proxy-app</code> that will effectively become the working directory for
the Coherence JVM in the <code>Pods</code> for the <code>proxy</code> role.</li>
</ul>
</div>

<h3 id="_setting_the_application_directory_for_explicit_roles_with_a_default">Setting the Application Directory for Explicit Roles with a Default</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with explicit roles in the <code>roles</code> list the <code>application.appDir</code> field can be set
at the <code>spec</code> level as a default that applies to all of the roles in the <code>roles</code> list unless specifically overridden
for an individual role:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  application:
    appDir: app           <span class="conum" data-value="1" />
  roles:
    - role: data          <span class="conum" data-value="2" />
    - role: proxy
      application:
        appDir: proxy-app <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The default value for the <code>appDir</code> field is <code>/app</code> which will be used for all roles unless specifically
overridden for a role.</li>
<li data-value="2">The <code>data</code> role does not specify a value for <code>appDir</code> so it will use the default. The application image should
contain a directory named <code>/app</code> will effectively become the working directory for the Coherence JVM in the <code>Pods</code> for
the <code>data</code> role.</li>
<li data-value="3">The <code>proxy</code> role has an explicit value set for the <code>appDir</code> field. The application image should a directory named
<code>/proxy-app</code> will effectively become the working directory for the Coherence JVM in the <code>Pods</code> for the <code>proxy</code> role</li>
</ul>
<v-divider class="my-5"/>
</div>
</div>

<h2 id="app-main">Setting the Application Main</h2>
<div class="section">
<p>By default Coherence containers will run the <code>main</code> method in the <code>com.tangosol.coherence.DefaultCacheServer</code>
class. Sometimes an application requires a different class as the main class (or even a main that is not a class at all,
for example when running a Node JS application on top of the Graal VM the <code>main</code> could be a Javascript file).
The main to be used can be configured for each role in a <code>CoherenceCluster</code>.</p>


<h3 id="_setting_the_application_main_class_for_the_implicit_role">Setting the Application Main Class for the Implicit Role</h3>
<div class="section">
<p>When configuring a <code>CoherenceCluster</code> with a single implicit role the application&#8217;s working directory is specified
in the <code>application.main</code> field.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  application:
    libDir: lib          <span class="conum" data-value="1" />
    main: com.acme.Main  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The application image should contain a directory named <code>/lib</code> that will contain the <code>.jar</code> files containing the
application classes and dependencies.</li>
<li data-value="2">One of those classes will be <code>com.acme.Main</code> which will be executed as the main class when starting the JVMs for
the <code>data</code> role.</li>
</ul>
</div>

<h3 id="_setting_the_application_main_class_for_explicit_roles">Setting the Application Main Class for Explicit Roles</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with explicit roles in the <code>roles</code> list the <code>application.main</code> field can be set
specifically for each role:</p>

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
        libDir: lib           <span class="conum" data-value="1" />
        main: com.acme.Main   <span class="conum" data-value="2" />
    - role: proxy
      application:
        libDir: lib
        main: com.acme.Proxy  <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The application image should contain a directory named <code>/lib</code> that will contain the <code>.jar</code> files containing the
application classes and dependencies.</li>
<li data-value="2">One of those classes will be <code>com.acme.Main</code> which will be executed as the main class when starting the JVMs for
the <code>data</code> role.</li>
<li data-value="3">The <code>proxy</code> role will use the <code>com.acme.Proxy</code> class as the main class</li>
</ul>
</div>

<h3 id="_setting_the_application_main_class_for_explicit_roles_with_a_default">Setting the Application Main Class for Explicit Roles with a Default</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with explicit roles in the <code>roles</code> list the <code>application.main</code> field can be set
at the <code>spec</code> level as a default that applies to all of the roles in the <code>roles</code> list unless specifically overridden
for an individual role:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  application:
    libDir: lib               <span class="conum" data-value="1" />
    main: com.acme.Main       <span class="conum" data-value="2" />
  roles:
    - role: data              <span class="conum" data-value="3" />
    - role: proxy
      application:
        main: com.acme.Proxy  <span class="conum" data-value="4" /></markup>

<ul class="colist">
<li data-value="1">The application image should contain a directory named <code>/lib</code> that will contain the <code>.jar</code> files containing the
application classes and dependencies.</li>
<li data-value="2">One of those classes will be <code>com.acme.Main</code> which will be executed as the main class for all roles that do not
specifically specify a <code>main</code>.</li>
<li data-value="3">The <code>data</code> role does not specify a <code>main</code> field so the Coherence JVM in the <code>Pods</code> for the <code>data</code> role will all use
the <code>com.acme.Main</code> class as the main class.</li>
<li data-value="4">The <code>proxy</code> role will specifies a <code>main</code> class to use so all Coherence JVMs in the <code>Pods</code> for the <code>proxy</code> role
will use the <code>com.acme.Proxy</code> class as the main class.</li>
</ul>
<v-divider class="my-5"/>
</div>
</div>

<h2 id="app-args">Setting the Application Main Arguments</h2>
<div class="section">
<p>Some applications that specify a custom <code>main</code> may also require command line arguments to be passed to the <code>main</code>,
These additional arguments can also be configured for the roles in a <code>CoherenceCluster</code>. Application arguments are
specified as a string array.</p>


<h3 id="_setting_the_application_main_arguments_for_the_implicit_role">Setting the Application Main Arguments for the Implicit Role</h3>
<div class="section">
<p>When configuring a <code>CoherenceCluster</code> with a single implicit role the application&#8217;s working directory is specified
in the <code>application.main</code> field.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  application:
    libDir: lib          <span class="conum" data-value="1" />
    main: com.acme.Main  <span class="conum" data-value="2" />
    args:                <span class="conum" data-value="3" />
      - "argOne"
      - "argTwo"</markup>

<ul class="colist">
<li data-value="1">The application image should contain a directory named <code>/lib</code> that will contain the <code>.jar</code> files containing the
application classes and dependencies.</li>
<li data-value="2">One of those classes will be <code>com.acme.Main</code> which will be executed as the main class when starting the JVMs for
the <code>data</code> role.</li>
<li data-value="3">The arguments <code>"argOne"</code> and <code>"argTwo"</code> will be passed to the <code>com.acme.Main</code> class <code>main()</code> method.</li>
</ul>
</div>

<h3 id="_setting_the_application_main_arguments_for_explicit_roles">Setting the Application Main Arguments for Explicit Roles</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with explicit roles in the <code>roles</code> list the <code>application.args</code> field can be set
specifically for each role:</p>

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
        libDir: lib           <span class="conum" data-value="1" />
        main: com.acme.Main   <span class="conum" data-value="2" />
        args:                 <span class="conum" data-value="3" />
          - "argOne"
          - "argTwo"
    - role: proxy
      application:
        libDir: lib
        main: com.acme.Main
        args:                 <span class="conum" data-value="4" />
          - "argThree"
          - "argFour"</markup>

<ul class="colist">
<li data-value="1">The application image should contain a directory named <code>/lib</code> that will contain the <code>.jar</code> files containing the
application classes and dependencies.</li>
<li data-value="2">One of those classes will be <code>com.acme.Main</code> which will be executed as the main class when starting the JVMs for
the <code>data</code> role.</li>
<li data-value="3">The arguments <code>"argOne"</code> and <code>"argTwo"</code> will be passed to the <code>com.acme.Main</code> class <code>main()</code> method in <code>Pods</code> for
the <code>data</code> role.</li>
<li data-value="4">The <code>proxy</code> role specifies different arguments. The arguments <code>"argThree"</code> and <code>"argFour"</code> will be passed to the
<code>com.acme.Main</code> class <code>main()</code> method in <code>Pods</code> for the <code>proxy</code> role.</li>
</ul>
</div>

<h3 id="_setting_the_application_main_arguments_for_explicit_roles_with_a_default">Setting the Application Main Arguments for Explicit Roles with a Default</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with explicit roles in the <code>roles</code> list the <code>application.main</code> field can be set
at the <code>spec</code> level as a default that applies to all of the roles in the <code>roles</code> list unless specifically overridden
for an individual role:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  application:
    libDir: lib               <span class="conum" data-value="1" />
    main: com.acme.Main       <span class="conum" data-value="2" />
    args:                     <span class="conum" data-value="3" />
      - "argOne"
      - "argTwo"
  roles:
    - role: data              <span class="conum" data-value="4" />
    - role: proxy
      application:
        args:                 <span class="conum" data-value="5" />
          - "argThree"
          - "argFour"
    - role: web
      application:
        args: []              <span class="conum" data-value="6" /></markup>

<ul class="colist">
<li data-value="1">The application image should contain a directory named <code>/lib</code> that will contain the <code>.jar</code> files containing the
application classes and dependencies.</li>
<li data-value="2">One of those classes will be <code>com.acme.Main</code> which will be executed as the main class for all roles that do not
specifically specify a <code>main</code>.</li>
<li data-value="3">The default args are <code>"argOne"</code> and <code>"argTwo"</code></li>
<li data-value="4">The <code>data</code> role does not specify an <code>args</code> field so the Coherence JVM in the <code>Pods</code> for the <code>data</code> role will all use
the default arguments of <code>"argOne"</code> and <code>"argTwo"</code></li>
<li data-value="5">The <code>proxy</code> role specifies different arguments. The arguments <code>"argThree"</code> and <code>"argFour"</code> will be passed to the
<code>com.acme.Main</code> class <code>main()</code> method in <code>Pods</code> for the <code>proxy</code> role.</li>
<li data-value="6">The <code>web</code> role specifies an empty array for the <code>args</code> field so no arguments will be passed to its main class.</li>
</ul>
</div>
</div>
</doc-view>
