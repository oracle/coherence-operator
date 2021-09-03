<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Add ConfigMap Volumes</dt>
<dd slot="desc"><p>Additional <code>Volumes</code> and <code>VolumeMounts</code> from <code>ConfigMaps</code> can easily be added to a <code>Coherence</code> resource.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_add_configmap_volumes">Add ConfigMap Volumes</h2>
<div class="section">
<p>To add a <code>ConfigMap</code> as an additional volume to the <code>Pods</code> of a Coherence deployment add entries to the
<code>configMapVolumes</code> list in the CRD spec.
Each entry in the list has a mandatory <code>name</code> and <code>mountPath</code> field, all other fields are optional.
The <code>name</code> field is the name of the <code>ConfigMap</code> to mount and is also used as the volume name.
The <code>mountPath</code> field is the path in the container to mount the volume to.</p>

<div class="admonition note">
<p class="admonition-inline">Additional volumes added in this way (either <code>ConfigMaps</code> shown here, or <code>Secrets</code> or plain <code>Volumes</code>) will be
added to all containers in the <code>Pod</code>.</p>
</div>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  configMapVolumes:
    - name: storage-config               <span class="conum" data-value="1" />
      mountPath: /home/coherence/config  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>ConfigMap</code> named <code>storage-config</code> will be mounted to the <code>Pod</code> as an additional <code>Volume</code> named <code>storage-config</code></li>
<li data-value="2">The <code>ConfigMap</code> will be mounted at <code>/home/coherence/config</code> in the containers.</li>
</ul>
<p>The yaml above would result in a <code>Pod</code> spec similar to the following (a lot of the <code>Pod</code> spec has been omitted to just
show the relevant volume information):</p>

<markup
lang="yaml"

>apiVersion: v1
kind: Pod
metadata:
  name: storage-0
spec:
  containers:
    - name: coherence
      volumeMounts:
        - name: storage-config
          mountPath: /home/coherence/config
  volumes:
    - name: storage-config
      configMap:
        name: storage-config</markup>

<p>As already stated, if the <code>Coherence</code> resource has additional containers the <code>ConfigMap</code> will be mounted in all of them.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  sideCars:
    - name: fluentd
      image: "fluent/fluentd:v1.3.3"
  configMapVolumes:
    - name: storage-config
      mountPath: /home/coherence/config</markup>

<p>In this example the <code>storage-config</code> <code>ConfigMap</code> will be mounted as a <code>Volume</code> and mounted to both the <code>coherence</code>
container and the <code>fluentd</code> container.
The yaml above would result in a <code>Pod</code> spec similar to the following (a lot of the <code>Pod</code> spec has been omitted to just
show the relevant volume information):</p>

<markup
lang="yaml"

>apiVersion: v1
kind: Pod
metadata:
  name: storage-0
spec:
  containers:
    - name: coherence
      volumeMounts:
        - name: storage-config
          mountPath: /home/coherence/config
    - name: fluentd
      image: "fluent/fluentd-kubernetes-daemonset:v1.3.3-debian-elasticsearch-1.3"
      volumeMounts:
        - name: storage-config
          mountPath: /home/coherence/config
  volumes:
    - name: storage-config
      configMap:
        name: storage-config</markup>

</div>
</doc-view>
