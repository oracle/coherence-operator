<doc-view>

<h2 id="_set_the_operational_configuration_file_name">Set the Operational Configuration File Name</h2>
<div class="section">
<p>The name of the Coherence operations configuration file (commonly called the overrides file) that the Coherence processes
in a <code>Coherence</code> resource will use can be set with the <code>spec.coherence.overrideConfig</code> field.
By setting this field the <code>coherence.override</code> system property will be set in the Coherence JVM.</p>

<p>When the <code>spec.coherence.overrideConfig</code> is blank or not specified, Coherence use its default behaviour to find the
operational configuration file to use. Typically, this is to use the first occurrence of <code>tangosol-coherence-override.xml</code>
that is found on the classpath
(consult the <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/develop-applications/understanding-configuration.html#GUID-360B798E-2120-44A9-8B09-1FDD9AB40EB5">Coherence documentation</a>
for an explanation of the default behaviour).</p>

<p>To set a specific operational configuration file to use set the <code>spec.coherence.overrideConfig</code> field, for example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  coherence:
    overrideConfig: test-override.xml <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>spec.coherence.overrideConfig</code> field has been set to <code>test-override.xml</code> which will effectively pass
<code>-Dcoherence.override=test-override.xml</code> to the JVM command line.</li>
</ul>
</div>
</doc-view>
