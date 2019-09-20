<doc-view>

<h2 id="_setting_the_coherence_image">Setting the Coherence Image</h2>
<div class="section">
<p>The Coherence Operator has a default setting for the Coherence image that will be used when by <code>Pods</code> in a <code>CoherenceCluster</code>.
This default value can be overridden to enable roles in the cluster to use a different image.
If the image being configured is from a registry requiring authentication see the section on <router-link to="/clusters/070_private_repos">pulling from private registries</router-link>.</p>

<p>As well as setting the image name it is also sometimes useful to set the Coherence image&#8217;s  <router-link to="#pull-policy" @click.native="this.scrollFix('#pull-policy')">image pull policy</router-link>.</p>


<h3 id="_setting_the_coherence_image_for_the_implicit_role">Setting the Coherence Image for the Implicit Role</h3>
<div class="section">
<p>When using the implicit role configuration the Coherence image to use is set directly in the <code>CoherenceCluster</code> <code>spec</code>
<code>images.coherence.image`</code> section.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  images:
    coherence:
      image: container-registry.oracle.com/middleware/coherence/12.2.1.4  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>coherence</code> container in the implicit role&#8217;s <code>Pod</code> will use the Coherence image <code>container-registry.oracle.com/middleware/coherence/12.2.1.4</code></li>
</ul>
</div>

<h3 id="_setting_the_coherence_image_for_explicit_roles">Setting the Coherence Image for Explicit Roles</h3>
<div class="section">
<p>When using the explicit roles in a <code>CoherenceCluster</code> <code>roles</code> list the Coherence image to use is set for each role.</p>

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
        coherence:
          image: container-registry.oracle.com/middleware/coherence/12.2.1.4.1  <span class="conum" data-value="1" />
    - role: proxy
      images:
        coherence:
          image: container-registry.oracle.com/middleware/coherence/12.2.1.4.0  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>coherence</code> container in the  <code>data</code> role <code>Pods</code> will use the Coherence
image <code>container-registry.oracle.com/middleware/coherence/12.2.1.4.1</code></li>
<li data-value="2">The <code>coherence</code> container in the  <code>proxy</code> role <code>Pods</code> will use the Coherence
image <code>container-registry.oracle.com/middleware/coherence/12.2.1.4.1</code></li>
</ul>
</div>

<h3 id="_setting_the_coherence_image_for_explicit_roles_with_a_default">Setting the Coherence Image for Explicit Roles with a Default</h3>
<div class="section">
<p>When using the explicit roles in a <code>CoherenceCluster</code> <code>roles</code> list the Coherence image to use can be set in the
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
    coherence:
      image: container-registry.oracle.com/middleware/coherence/12.2.1.4.0  <span class="conum" data-value="1" />
  roles:
    - role: data
    - role: proxy</markup>

<ul class="colist">
<li data-value="1">The image <code>container-registry.oracle.com/middleware/coherence/12.2.1.4.0</code> set in the <code>spec</code> section will be used by
both the <code>data</code> and the <code>proxy</code> roles. The <code>coherence</code> container in all of the <code>Pods</code> will use this image.</li>
</ul>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  images:
    coherence:
      image: container-registry.oracle.com/middleware/coherence/12.2.1.4.0  <span class="conum" data-value="1" />
  roles:
    - role: data
      images:
        coherence:
          image: container-registry.oracle.com/middleware/coherence/12.2.1.4.1  <span class="conum" data-value="2" />
    - role: proxy
    - role: web</markup>

<ul class="colist">
<li data-value="1">The image <code>container-registry.oracle.com/middleware/coherence/12.2.1.4.0</code> set in the <code>spec</code> section will be used by
both the <code>proxy</code> and the <code>web</code> roles. The <code>coherence</code> container in all of the <code>Pods</code> will use this image.</li>
<li data-value="2">The <code>container-registry.oracle.com/middleware/coherence/12.2.1.4.1</code> image is specifically set for the <code>data</code> role
so the <code>coherence</code> container in the <code>Pods</code> for the <code>data</code> role will use this image.</li>
</ul>
</div>
</div>

<h2 id="pull-policy">Setting the Coherence Image Pull Policy</h2>
<div class="section">
<p>The image pull policy controls when (and if) Kubernetes will pull the Coherence image onto the node where the Coherence
<code>Pods</code> are being schedules.
See <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/containers/images/#updating-images">Kubernetes imagePullPolicy</a> for more information.</p>

<div class="admonition note">
<p class="admonition-inline">The Kubernetes default pull policy is <code>IfNotPresent</code> unless the image tag is <code>:latest</code> in which case the default
policy is <code>Always</code>. The <code>IfNotPresent</code> policy causes the Kubelet to skip pulling an image if it already exists.
Note that you should avoid using the <code>:latest</code> tag, see
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/configuration/overview/#container-images">Kubernetes Best Practices for Configuration</a>
for more information.</p>
</div>
<p>The Coherence image&#8217;s pull policy is set using the <code>imagePullPolicy</code> field in the <code>spec.images.coherence</code> section.</p>

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
    coherence:
      image: container-registry.oracle.com/middleware/coherence/12.2.1.4.0
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
        coherence:
          image: container-registry.oracle.com/middleware/coherence/12.2.1.4.1
          imagePullPolicy: Always <span class="conum" data-value="1" />
    - role: proxy
      images:
        coherence:
          image: container-registry.oracle.com/middleware/coherence/12.2.1.4.0
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
    coherence:
      imagePullPolicy: Always <span class="conum" data-value="1" />
  roles:
    - role: data
      images:
        coherence:
          image: container-registry.oracle.com/middleware/coherence/12.2.1.4.1
    - role: proxy
      images:
        coherence:
          image: container-registry.oracle.com/middleware/coherence/12.2.1.4.1
    - role: web
      images:
        coherence:
          image: container-registry.oracle.com/middleware/coherence/12.2.1.4.0
          imagePullPolicy: IfNotPresent <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The default image pull policy is set to <code>Always</code>. The <code>data</code> and <code>proxy</code> roles will use the default value because
they do not specifically set the value in their specs.</li>
<li data-value="2">The image pull policy for the <code>web</code> role above has been set to <code>IfNotPresent</code></li>
</ul>
</div>
</doc-view>
