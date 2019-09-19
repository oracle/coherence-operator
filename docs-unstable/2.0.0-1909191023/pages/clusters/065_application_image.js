<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Setting the Application Image</dt>
<dd slot="desc"><p>The application image is the Docker image containing the <code>.jar</code> files and configuration files of the Coherence application
that should be deployed in the Coherence cluster. For more information see the
<router-link to="/guides/030_applications">Deploying Coherence Applications</router-link> guide.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_setting_the_application_image">Setting the Application Image</h2>
<div class="section">
<p>Whilst the Coherence Operator makes it simple to deploy and manage a Coherence cluster in Kubernetes in the majority of
use cases there will be a requirement for application code to be deployed and run in the Coherence JVMs. This application
code and any application configuration files are supplied as a separate image. This image is loaded as an init-container
by the Coherence <code>Pods</code> and the relevant <code>.jar</code> files and configuration files from this image are added to the Coherence
JVM&#8217;s classpath.</p>

<p>As well as setting the image name it is also sometimes useful to set the application image&#8217;s  <router-link to="#pull-policy" @click.native="this.scrollFix('#pull-policy')">image pull policy</router-link>.</p>


<h3 id="_setting_the_application_image_for_the_implicit_role">Setting the Application Image for the Implicit Role</h3>
<div class="section">
<p>When using the implicit role configuration the application image to use is set directly in the <code>CoherenceCluster</code> <code>spec</code>
<code>images.userArtifacts.image</code> section.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  images:
    userArtifacts:
      image: acme/orders-data/1.0.0  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>acme/orders-data/1.0.0</code> will be used to add additional <code>.jar</code> files and configuration files to the classpath of
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
      images:
        userArtifacts:
          image: acme/orders-data/1.0.0  <span class="conum" data-value="1" />
    - role: proxy
      images:
        userArtifacts:
          image: acme/orders-proxy/1.0.0  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role <code>Pods</code> will use the application image <code>acme/orders-data/1.0.0</code></li>
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
  images:
    userArtifacts:
      image: acme/orders-data/1.0.0  <span class="conum" data-value="1" />
  roles:
    - role: data
    - role: proxy</markup>

<ul class="colist">
<li data-value="1">The <code>data</code> and the <code>proxy</code> roles will both use the application image <code>acme/orders-data/1.0.0</code></li>
</ul>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  images:
    userArtifacts:
      image: acme/orders-data/1.0.0  <span class="conum" data-value="1" />
  roles:
    - role: data
    - role: proxy
    - role: web
      images:
        userArtifacts:
          image: acme/orders-front-end/1.0.0  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> and the <code>proxy</code> roles will both use the application image <code>acme/orders-data/1.0.0</code></li>
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
<p>The application image&#8217;s pull policy is set using the <code>imagePullPolicy</code> field in the <code>spec.images.coherence</code> section.</p>

<p>For example:</p>

<ol style="margin-left: 15px;">
<li>
To set the <code>imagePullPolicy</code> for the implicit role.

</li>
</ol>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  images:
    userArtifacts:
      image: acme/orders-data/1.0.0
      imagePullPolicy: Always <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The image pull policy for the implicit role above has been set to <code>Always</code></li>
</ul>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data
      images:
        userArtifacts:
          image: acme/orders-data/1.0.0
          imagePullPolicy: Always <span class="conum" data-value="1" />
    - role: proxy
      images:
        userArtifacts:
          image: acme/orders-proxy/1.0.0
          imagePullPolicy: IfNotPresent <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The image pull policy for the <code>data</code> role has been set to <code>Always</code></li>
<li data-value="2">The image pull policy for the <code>proxy</code> role above has been set to <code>IfNotPresent</code></li>
</ul>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  images:
    userArtifacts:
      imagePullPolicy: Always <span class="conum" data-value="1" />
  roles:
    - role: data
      images:
        userArtifacts:
          image: acme/orders-data/1.0.0
    - role: proxy
      images:
        userArtifacts:
          image: acme/orders-proxy/1.0.1
    - role: web
      images:
        userArtifacts:
          image: acme/orders-front-end/1.0.1
          imagePullPolicy: IfNotPresent <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The default image pull policy is set to <code>Always</code>. The <code>data</code> and <code>proxy</code> roles will use the default value because
they do not specifically set the value in their specs.</li>
<li data-value="2">The image pull policy for the <code>web</code> role above has been set to <code>IfNotPresent</code></li>
</ul>
</div>
</doc-view>
