<doc-view>

<h2 id="_adding_arbitrary_jvm_arguments">Adding Arbitrary JVM Arguments</h2>
<div class="section">
<p>The <code>Coherence</code> CRD allows any arbitrary JVM arguments to be passed to the JVM in the <code>coherence</code> container
by using the <code>jvm.args</code> field of rthe CRD spec.
Any valid system property or JVM argument can be added to the <code>jvm.args</code> list.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    args:
      - "-Dcoherence.pof.config=storage-pof-config.xml"
      - "-Dcoherence.tracing.ratio=0.1"
      - "-agentpath:/yourkit/bin/linux-x86-64/libyjpagent.so"</markup>

<p>In this example the <code>args</code> list adds two System properties <code>coherence.pof.config=storage-pof-config.xml</code>
and <code>coherence.tracing.ratio=0.1</code> and also adds the YourKit profiling agent.</p>

<div class="admonition note">
<p class="admonition-inline">When the Operator builds the command line to use when starting Coherence Pods, any arguments added to
the <code>jvm.args</code> field will be added after all the arguments added by the Operator from other configuration fields.
This means that arguments such as system properties added to <code>jvm.args</code> will override any added by the Operator.</p>
</div>
<p>For example</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  coherence:
    cacheConfig: storage-config.xml                   <span class="conum" data-value="1" />
  jvm:
    args:
      - "-Dcoherence.cache.config=test-config.xml"    <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">Setting the <code>coherence.cacheConfig</code> field will make the operator add
<code>-Dcoherence.cache.config=storage-config.xml</code> to the command line.</li>
<li data-value="2">Adding <code>-Dcoherence.cache.config=test-config.xml</code> to the <code>jvm.args</code> field will make the Operator add
<code>-Dcoherence.cache.config=test-config.xml</code> to the end of the JVM arguments in the command line.</li>
</ul>
<p>When duplicate system properties are present on a command line the last one wins so in the example above the cache
configuration used would be  <code>test-config.xml</code>.</p>

</div>
</doc-view>
