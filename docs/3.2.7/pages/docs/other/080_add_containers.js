<doc-view>

<h2 id="_configure_additional_containers">Configure Additional Containers</h2>
<div class="section">
<p>Additional containers and init-containers can easily be added to a <code>Coherence</code> resource Pod.
There are two types of container that can be added, init-containers and normal containers.
An example use case for this would be to add something like a Fluentd side-car container to ship logs to Elasticsearch.</p>

<div class="admonition note">
<p class="admonition-inline">A note about Volumes:<br>
The Operator created a number of volumes and volume mounts by default. These default volume mounts will be added
to <strong>all</strong> containers in the <code>Pod</code> including containers added as described here.<br>
Any additional volumes and volume mounts added to the <code>Coherence</code> resource spec will also be added <strong>all</strong> containers.</p>
</div>

<h3 id="_add_a_container">Add a Container</h3>
<div class="section">
<p>To add a container to the <code>Pods</code> specify the container in the <code>sideCars</code> list in the <code>Coherence</code> CRD spec.</p>

<p>See the <router-link to="/docs/logging/020_logging">Logging Documentation</router-link> for a bigger example of adding a side-car container.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  sideCars:
    - name: fluentd                   <span class="conum" data-value="1" />
      image: "fluent/fluentd:v1.3.3"</markup>

<ul class="colist">
<li data-value="1">An additional container named <code>fluentd</code> has been added to the CRD spec.</li>
</ul>
<p>The containers will added to the <code>sideCars</code> will be added to the <code>Pods</code> exactly as configured.
Any configuration that is valid in a Kubernetes
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#container-v1-core">Container Spec</a>
may be added to an entry in <code>sideCars</code></p>

</div>

<h3 id="_add_an_init_container">Add an Init-Container</h3>
<div class="section">
<p>Just like normal containers above, additional init-containers can also be added to the <code>Pods</code>.
To add an init-container to the <code>Pods</code> specify the container in the <code>initContainers</code> list in the <code>Coherence</code> CRD spec.
As with containers, for init-containers any configuration that is valid in a Kubernetes
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#container-v1-core">Container Spec</a>
may be added to an entry in <code>initContainers</code></p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  initContainers:
    - name: setup                   <span class="conum" data-value="1" />
      image: "app-setup:1.0.0"</markup>

<ul class="colist">
<li data-value="1">An additional init-container named <code>setup</code> has been added to the CRD spec with an image named <code>app-setup:1.0.0</code>.</li>
</ul>
</div>
</div>
</doc-view>
