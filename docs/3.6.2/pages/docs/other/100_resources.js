<doc-view>

<h2 id="_container_resource_limits">Container Resource Limits</h2>
<div class="section">
<p>When creating a <code>Coherence</code> resource you can optionally specify how much CPU and memory (RAM) each Coherence Container
is allowed to consume. The container resources are specified in the <code>resources</code> section of the <code>Coherence</code> spec;
the format is exactly the same as documented in the Kubernetes documentation
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/">Managing Compute Resources for Containers</a>.</p>

<div class="admonition warning">
<p class="admonition-inline">When setting resource limits, in particular memory limits, for a container it is important to ensure that the
Coherence JVM is properly configured so that it does not consume more memory than the limits. If the JVM attempts to
consume more memory than the resource limits allow the <code>Pod</code> can be killed by Kubernetes.
See <router-link to="/docs/jvm/050_memory">Configuring the JVM Memory</router-link> for details on the different memory settings.</p>
</div>
<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test-cluster
spec:
  resources:           <span class="conum" data-value="1" />
    requests:
      memory: "64Mi"
      cpu: "250m"
    limits:
      memory: "128Mi"
      cpu: "500m"</markup>

<ul class="colist">
<li data-value="1">The <code>coherence</code> container in the <code>Pods</code> will have requests of 0.25 cpu and 64MiB of memory,
and limits of 0.5 cpu and 128MiB of memory.</li>
</ul>
</div>

<h2 id="_initcontainer_resource_limits">InitContainer Resource Limits</h2>
<div class="section">
<p>The Coherence Operator adds an init-container to the Pods that it manages. This init container does nothing more
than copy some files and ensure some directories exist. In terms of resource use it is extremely light.
Some customers have expressed a desire to still be able to set limits fo this init container, so this is possible
using the <code>spec.initResources</code> field.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test-cluster
spec:
  initResources:           <span class="conum" data-value="1" />
    requests:
      memory: "64Mi"
      cpu: "250m"
    limits:
      memory: "128Mi"
      cpu: "500m"</markup>

<ul class="colist">
<li data-value="1">The <code>coherence-k8s-utils</code> init-container in the <code>Pods</code> will have requests of 0.25 cpu and 64MiB of memory,
and limits of 0.5 cpu and 128MiB of memory.</li>
</ul>
<p>These resources only applies to the init-container that the Operator creates, any other init-containers added in the
<code>spec.initContainers</code> section should have their own resources configured.</p>

</div>
</doc-view>
