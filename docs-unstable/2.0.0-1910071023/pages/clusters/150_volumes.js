<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Configure Additional Volumes</dt>
<dd slot="desc"><p>Although a Coherence cluster member may not need access to specific volumes custom applications deployed into a cluster
may require them. For this reason it is possible to configure roles in a <code>CoherenceCluster</code> with arbitrary <code>VolumeMounts</code>.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_configure_additional_volumes">Configure Additional Volumes</h2>
<div class="section">
<p>There are two parts to configuring a <code>Volume</code> that will be accessible to an application running in th Coherence <code>Pods</code>.
First a <code>Volume</code> must be defined for the <code>Pod</code> itself and then a corresponding <code>VolumeMount</code> must be configured that
will be added to the Coherence container in the <code>Pods</code>.
Additional <code>Volumes</code> and <code>VolumeMounts</code> can be added to a <code>CoherenceCluster</code> by defining each additional <code>Volume</code>
and <code>VolumeMount</code> using exactly the same yaml that would be used if adding <code>Volumes</code> to Kubernetes <code>Pods</code> as described
in the <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/storage/volumes/">Kubernetes Volumes documentation</a></p>


<h3 id="_adding_volumes_to_the_implicit_role">Adding Volumes to the Implicit Role</h3>
<div class="section">
<p>When defining a <code>CoherenceCluster</code> with a single implicit role the <code>Volumes</code> and <code>VolumeMounts</code> are added directly to
the <code>CoherenceCluster</code> <code>spec</code></p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  volumes:
     - name: test-volume      <span class="conum" data-value="1" />
       hostPath:
         path: /data
         type: Directory
  volumeMounts:
     - name: test-volume      <span class="conum" data-value="2" />
       mountPath: /test-data</markup>

<ul class="colist">
<li data-value="1">An additional <code>Volume</code> named <code>test-volume</code> will be added to the <code>Pod</code>. In this case the <code>Volume</code> is
a <code>hostPath</code> volume type.</li>
<li data-value="2">A corresponding <code>VolumeMount</code> is added so that the <code>test-volume</code> will be mounted into the Coherence container
with the path <code>/test-data</code></li>
</ul>
<p>Multiple <code>Volumes</code> and <code>VolumeMappings</code> can be added:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data
      volumes:
         - name: test-volume      <span class="conum" data-value="1" />
           hostPath:
             path: /data
             type: Directory
      volumeMounts:
         - name: test-volume      <span class="conum" data-value="2" />
           mountPath: /test-data
    - role: proxy
      volumes:
         - name: proxy-volume     <span class="conum" data-value="3" />
           hostPath:
             path: /proxy-data
             type: Directory
      volumeMounts:
         - name: test-volume      <span class="conum" data-value="4" />
           mountPath: /data</markup>

<ul class="colist">
<li data-value="1">An additional Host Path <code>Volume</code> named <code>test-volume</code> will be added to the containers in the <code>Pods</code> for the <code>data</code> role.</li>
<li data-value="2">An additional <code>VolumeMount</code> to mount the <code>test-volume</code> to the <code>/test-data</code> path will be added to the containers in
the <code>Pods</code> for the <code>data</code> role.</li>
<li data-value="3">An additional Host Path <code>Volume</code> named <code>proxy-volume</code> will be added to the containers in the <code>Pods</code> for the <code>proxy</code> role.</li>
<li data-value="4">An additional <code>VolumeMount</code> to mount the <code>proxy-volume</code> to the <code>/proxy-data</code> path will be added to the containers in
the <code>Pods</code> for the <code>proxy</code> role.</li>
</ul>
</div>

<h3 id="_adding_volumes_to_explicit_roles">Adding Volumes to Explicit Roles</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles additional <code>Volumes</code> can be added to the configuration
of each role.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  volumes:                   <span class="conum" data-value="1" />
     - name: data-volume
       hostPath:
         path: /test-data
         type: Directory
     - name: config-volume
       hostPath:
         path: /test-config
         type: Directory
  volumeMounts:              <span class="conum" data-value="2" />
     - name: test-volume
       mountPath: /data
     - name: config-volume
       mountPath: /config</markup>

<ul class="colist">
<li data-value="1">The <code>volumes</code> list has two additional <code>Volumes</code>, <code>data-volume</code> and <code>config-volume</code></li>
<li data-value="2">The <code>volumeMounts</code> list has two corresponding <code>VolumeMounts</code>.</li>
</ul>
</div>

<h3 id="_adding_volumes_to_explicit_roles_with_defaults">Adding Volumes to Explicit Roles with Defaults</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles additional <code>Volumes</code> and <code>VolumeMounts</code> can be added
as defaults that will apply to all roles in the <code>roles</code> list. The additional <code>Volumes</code> and <code>VolumeMounts</code> in the defaults
section will be merged with any additional <code>Volumes</code> and <code>VolumeMounts</code> specified for the role.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  volumes:
     - name: data-volume          <span class="conum" data-value="1" />
       hostPath:
         path: /test-data
         type: Directory
  volumeMounts:
     - name: test-volume          <span class="conum" data-value="2" />
       mountPath: /data
  roles:
    - role: data                  <span class="conum" data-value="3" />
    - role: proxy
      volumes:
         - name: config-volume    <span class="conum" data-value="4" />
           hostPath:
             path: /proxy-config
             type: Directory
      volumeMounts:
         - name: config-volume    <span class="conum" data-value="5" />
           mountPath: /config</markup>

<ul class="colist">
<li data-value="1">The default <code>volumes</code> list has one additional <code>Volumes</code>, <code>data-volume</code></li>
<li data-value="2">The default <code>volumeMounts</code> list has one corresponding <code>VolumeMounts</code></li>
<li data-value="3">The <code>data</code> role does not have any additional <code>Volumes</code> or <code>VolumeMounts</code> so it will just inherit the default
<code>Volume</code> named <code>data-volume</code> and <code>VolumeMount</code> named <code>test-volume</code></li>
<li data-value="4">The <code>proxy</code> role has an additional <code>Volume</code> named <code>config-volume</code> so when the <code>Volume</code> lists are merged it will
have two additional <code>Volumes</code> <code>config-volume</code> and <code>test-volume</code></li>
<li data-value="5">The <code>proxy</code> role has an additional <code>VolumeMount</code> named <code>config-volume</code> so when the <code>VolumeMount</code> lists are merged
it will have two additional <code>VolumeMounts</code> <code>config-volume</code> and <code>test-volume</code></li>
</ul>
<p>When configuring explicit roles with default <code>Volumes</code> and <code>VolumeMounts</code> if the <code>role</code> defines a <code>Volume</code>
or <code>VolumeMount</code> with the same name as one defined in the defaults then the role&#8217;s definition overrides the
default definition. For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  volumes:
     - name: data-volume
       hostPath:
         path: /test-data
         type: Directory
  volumeMounts:
     - name: test-volume
       mountPath: /data
  roles:
    - role: data
    - role: proxy
      volumes:
         - name: data-volume      <span class="conum" data-value="1" />
           hostPath:
             path: /proxy-data
             type: Directory
         - name: config-volume
           hostPath:
             path: /proxy-config
             type: Directory
      volumeMounts:
         - name: config-volume
           mountPath: /config</markup>

<ul class="colist">
<li data-value="1">The <code>proxy</code> role overrides the default <code>data-volume</code> <code>Volume</code> with a different configuration.</li>
</ul>
</div>
</div>
</doc-view>
