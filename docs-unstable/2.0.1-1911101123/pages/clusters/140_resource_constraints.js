<doc-view>

<h2 id="_container_resource_limits">Container Resource Limits</h2>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> you can optionally specify how much CPU and memory (RAM) each Coherence Container
needs. The container resources are specified in the <code>resources</code> section of a role in a <code>CoherenceCluster</code>, the format
is exactly the same as documented in the Kubernetes documentation
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/">Managing Compute Resources for Containers</a>.</p>

<div class="admonition warning">
<p class="admonition-inline">When setting resource limits, in particular memory limits, for a container it is important to ensure that the
Coherence JVM is properly configured so that it does not consume more memory than the limits. If the JVM attempts to
consume more memory than the resource limits allow the <code>Pod</code> can be killed by Kubernetes.
See <router-link to="/clusters/080_jvm">Configuring the JVM</router-link> for details on the different memory settings.</p>
</div>

<h3 id="_configure_resource_limits_for_the_single_implicit_role">Configure Resource Limits for the Single Implicit Role</h3>
<div class="section">
<p>When using the implicit role configuration of the resource limits is set directly in the <code>CoherenceCluster</code> <code>spec</code>
<code>resources</code> section.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  resources:               <span class="conum" data-value="1" />
    requests:
      memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"</markup>

<ul class="colist">
<li data-value="1">The <code>coherence</code> container in the implicit role&#8217;s <code>Pods</code> has a request of 0.25 cpu and 64MiB (226 bytes) of memory.
The <code>coherence</code> container has a limit of 0.5 cpu and 128MiB of memory.</li>
</ul>
</div>

<h3 id="_configure_resource_limits_for_explicit_roles">Configure Resource Limits for Explicit Roles</h3>
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
      resources:               <span class="conum" data-value="1" />
        requests:
          memory: "10Gi"
            cpu: "4"
          limits:
            memory: "15Gi"
            cpu: "4"
    - role: proxy
      resources:               <span class="conum" data-value="2" />
        requests:
          memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"</markup>

<ul class="colist">
<li data-value="1">The <code>coherence</code> container in the <code>data</code> role&#8217;s <code>Pods</code> has a request of 4 cpu and 10GiB of memory.
The <code>coherence</code> container has a limit of 4 cpu and 15GiB of memory.</li>
<li data-value="2">The <code>coherence</code> container in the <code>proxy</code> role&#8217;s <code>Pods</code> has a request of 0.25 cpu and 64MiB of memory.
The <code>coherence</code> container has a limit of 0.5 cpu and 128MiB of memory.</li>
</ul>
</div>

<h3 id="_configure_resource_limits_for_explicit_roles_with_a_default">Configure Resource Limits for Explicit Roles with a Default</h3>
<div class="section">
<p>When using the explicit roles in a <code>CoherenceCluster</code> <code>roles</code> list the resource limits to use can be set in the
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
  resources:                   <span class="conum" data-value="1" />
    requests:
      memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"
  roles:
    - role: data
      resources:               <span class="conum" data-value="2" />
        requests:
          memory: "10Gi"
            cpu: "4"
          limits:
            memory: "15Gi"
            cpu: "4"
    - role: proxy              <span class="conum" data-value="3" />
    - role: web</markup>

<ul class="colist">
<li data-value="1">The default resource limits has a request of 0.25 cpu and 64MiB (226 bytes) of memory and has a limit of 0.5 cpu
and 128MiB of memory.</li>
<li data-value="2">The <code>data</code> role overrides the defaults and specifies a request of 4 cpu and 10GiB of memory.
The <code>coherence</code> container has a limit of 4 cpu and 15GiB of memory.</li>
<li data-value="3">The <code>proxy</code> role and the <code>web</code> role do not specify resource limits so the defaults will apply so that <code>Pods</code> in the
<code>proxy</code> and <code>web</code> roles have a request of 0.25 cpu and 64MiB (226 bytes) of memory and has a limit of 0.5 cpu and 128MiB
of memory.</li>
</ul>
</div>
</div>
</doc-view>
