<doc-view>

<h2 id="_set_the_cache_configuration_file_name">Set the Cache Configuration File Name</h2>
<div class="section">
<p>The name of the Coherence cache configuration file that the Coherence processes in a <code>Coherence</code> resource will
use can be set with the <code>spec.coherence.cacheConfig</code> field. By setting this field the <code>coherence.cacheconfig</code> system
property will be set in the Coherence JVM.</p>

<p>When the <code>spec.coherence.cacheConfig</code> is blank or not specified, Coherence use its default behaviour to find the
cache configuration file to use. Typically, this is to use the first occurrence of <code>coherence-cache-config.xml</code> that is
found on the classpath
(consult the <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/develop-applications/understanding-configuration.html#GUID-360B798E-2120-44A9-8B09-1FDD9AB40EB5">Coherence documentation</a>
for an explanation of the default behaviour).</p>

<p>To set a specific cache configuration file to use set the <code>spec.coherence.cacheConfig</code> field, for example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  coherence:
    cacheConfig: storage-cache-config.xml <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>spec.coherence.cacheConfig</code> field has been set to <code>storage-cache-config.xml</code> which will effectively pass
<code>-Dcoherence.cacheconfig=storage-cache-config.xml</code> to the JVM command line.</li>
</ul>
</div>
</doc-view>
