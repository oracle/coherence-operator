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
<li data-value="1">The <code>coherence</code> container in the <code>Pods</code> has a request of 0.25 cpu and 64MiB of memory.
The <code>coherence</code> container has a limit of 0.5 cpu and 128MiB of memory.</li>
</ul>
</div>
</doc-view>
