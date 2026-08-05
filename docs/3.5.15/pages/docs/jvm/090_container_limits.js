<doc-view>

<h2 id="_respect_container_resource_limits">Respect Container Resource Limits</h2>
<div class="section">
<p>The JVM can be configured to respect container limits set, for example cpu and memory limits.
This can be important if container limits have been set for the container in the <code>resources</code> section as a JVM that
does not respect these limits can cause the <code>Pod</code> to be killed.
This is done by adding the <code>-XX:+UseContainerSupport</code> JVM option.
It is possible to control this using the <code>jvm.useContainerLimits</code> field in the <code>Coherence</code> CRD spec.
If the field is not set, the operator adds the <code>-XX:+UseContainerSupport</code> option by default.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    useContainerLimits: false   <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>useContainerLimits</code> field is set to false, so the <code>-XX:+UseContainerSupport</code> will not be passed to the JVM.</li>
</ul>
<p>See the <router-link to="/docs/other/100_resources">Resource Limits</router-link> documentation on how to specify resource limits
for the Coherence container.</p>

</div>
</doc-view>
