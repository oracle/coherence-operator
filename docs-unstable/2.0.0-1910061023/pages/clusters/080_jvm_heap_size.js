<doc-view>

<h2 id="_setting_the_jvm_heap_size">Setting the JVM Heap Size</h2>
<div class="section">
<p>It is good practice to fix the Coherence JVM heap size and to set both the JVM <code>-Xmx</code> and <code>-Xms</code> options to the same value.
The heap size of the JVM can be configured for roles in the <code>jvm.heapSize</code> field of a role spec. If the <code>heapSize</code> value
is configured then that value is applied to bot the JVMs minimum and maximum heap sizes (i.e. used to set both
<code>-Xms</code> and -<code>Xmx</code>).</p>

<p>The format of the value of the <code>heapSize</code> field is any valid value that can be used when setting the <code>-Xmx</code> JVM option,
for example <code>10G</code> would set a 10 GB heap.</p>

<div class="admonition note">
<p class="admonition-inline">If the <code>Pod</code> resource limits are being set to limit memory usage of a <code>Pod</code> it is important that the <code>heapSize</code> value
is also set if using a Coherence image with a JVM version lower than Java 10. Prior to Java 10 the JVM could see all of
the memory available to a machine regardless of any Pod limits. The JVM could then easily attempt to consume more memory
that the <code>Pod</code> or <code>Container</code> was allowed and consequently crashing the <code>Pod</code>. With Coherence images that use a version
of Java above 10 this issue is less of a problem. Even so if using the <code>resources</code> section of the configuration to
limit a <code>Pod</code> or <code>Containers</code> memory it is a good idea to limit the JVM heap.</p>
</div>

<h3 id="_setting_the_jvm_heap_size_for_the_implicit_role">Setting the JVM Heap Size for the Implicit Role</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with a single implicit role the <code>heapSize</code> is set in the <code>spec.jvm</code> section of
the configuration. For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    heapSize: 10g <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The Coherence JVM for the implicit role defined above will have a 10 GB heap.
Equivalent to passing <code>-Xms10g -Xmx10g</code> to the JVM.</li>
</ul>
</div>

<h3 id="_setting_the_jvm_heap_size_for_explicit_roles">Setting the JVM Heap Size for Explicit Roles</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles the <code>heapSize</code> is set in the <code>jvm</code> section of
the configuration for each <code>role</code> in the <code>roles</code> list. For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  roles:
    - role: data
      jvm:
        heapSize: 10g <span class="conum" data-value="1" />
    - role: proxy
      jvm:
        heapSize: 500m <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The Coherence JVM for the <code>data</code> role defined above will have a 10 GB heap.
Equivalent to passing <code>-Xms10g -Xmx10g</code> to the JVM.</li>
<li data-value="2">The Coherence JVM for the <code>proxy</code> role defined above will have a 500 MB heap.
Equivalent to passing <code>-Xms500m -Xmx500m</code> to the JVM.</li>
</ul>
</div>

<h3 id="_setting_the_jvm_heap_size_for_explicit_roles_with_a_default">Setting the JVM Heap Size for Explicit Roles with a Default</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles a default <code>heapSize</code> value can be set in the
<code>CoherenceCluster</code> <code>spec</code> section that will apply t all of the roles in the <code>roles</code> list unless specifically
overridden by a role&#8217;s <code>jvm.heapSize</code> field. For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    heapSize: 500m    <span class="conum" data-value="1" />
  roles:
    - role: data
      jvm:
        heapSize: 10g <span class="conum" data-value="2" />
    - role: proxy    <span class="conum" data-value="3" />
    - role: web      <span class="conum" data-value="4" /></markup>

<ul class="colist">
<li data-value="1">The default max heap size of 500 MB will be applied to all of the roles in the cluster unless overridden for a
specific role.</li>
<li data-value="2">The <code>data</code> role overrides the default value to set the max heap for all JVMs in the <code>data</code> role to 10 GB.
Equivalent to passing <code>-Xms10g -Xmx10g</code> to the JVM.</li>
<li data-value="3">The <code>proxy</code> role does not specify a <code>heapSize</code> value so it will use the default value of 500 MB.</li>
<li data-value="4">The <code>web</code> role does not specify a <code>heapSize</code> value so it will use the default value of 500 MB.</li>
</ul>
</div>
</div>
</doc-view>
