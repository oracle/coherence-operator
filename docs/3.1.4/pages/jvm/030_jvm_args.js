<doc-view>

<h2 id="_adding_arbitrary_jvm_arguments">Adding Arbitrary JVM Arguments</h2>
<div class="section">
<p>The <code>Coherence</code> CRD allows any arbitrary JVM arguments to be passed to the JVM in the <code>coherence</code> container
by using the <code>jvm.args</code> field of the CRD spec.
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
<p>When duplicate system properties are present on a command line, the last one wins so in the example above the cache
configuration used would be  <code>test-config.xml</code>.</p>


<h3 id="_default_arguments">Default Arguments</h3>
<div class="section">
<p>The Coherence Operator will add the following JVM arguments by default:</p>

<markup


>-Dcoherence.cluster=&lt;cluster-name&gt;
-Dcoherence.role=&lt;role&gt;
-Dcoherence.wka=&lt;deployment-name&gt;-wka.svc.cluster.local
-Dcoherence.cacheconfig=coherence-cache-config.xml
-Dcoherence.k8s.operator.health.port=6676
-Dcoherence.management.http.port=30000
-Dcoherence.metrics.http.port=9612
-Dcoherence.distributed.persistence-mode=on-demand
-Dcoherence.override=k8s-coherence-override.xml
-Dcoherence.ttl=0

-XX:+UseG1GC
-XX:+PrintCommandLineFlags
-XX:+PrintFlagsFinal
-XshowSettings:all
-XX:+UseContainerSupport
-XX:+HeapDumpOnOutOfMemoryError
-XX:+ExitOnOutOfMemoryError
-XX:HeapDumpPath=/jvm/&lt;member&gt;/&lt;pod-uid&gt;/heap-dumps/&lt;member&gt;-&lt;pod-uid&gt;.hprof
-XX:ErrorFile=/jvm/&lt;member&gt;/&lt;pod-uid&gt;/hs-err-&lt;member&gt;-&lt;pod-uid&gt;.log
-XX:+UnlockDiagnosticVMOptions
-XX:NativeMemoryTracking=summary
-XX:+PrintNMTStatistics</markup>

<p>Some arguments and system properties above can be overridden or changed by setting values in the <code>Coherence</code> CDR spec.</p>

</div>
</div>

<h2 id="_environment_variable_expansion">Environment Variable Expansion</h2>
<div class="section">
<p>The Operator supports environment variable expansion in JVM arguments.
The runner in the Coherence container will replace <code>${var}</code> or <code>$var</code> in the JVM arguments with the corresponding environment variable name.</p>

<p>For example a JVM argument of <code>"-Dmy.host.name=${HOSTNAME}"</code> when run on a Pod with a host name of <code>COH-1</code> will resolve to <code>"-Dmy.host.name=COH-1"</code>.</p>

<p>Any environment variable that is present when the Coherence container starts can be used, this would include variables created as part of the image and variables specified in the Coherence yaml.</p>

</div>
</doc-view>
