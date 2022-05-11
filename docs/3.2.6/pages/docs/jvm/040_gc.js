<doc-view>

<h2 id="_garbage_collector_settings">Garbage Collector Settings</h2>
<div class="section">
<p>The <code>Coherence</code> CRD has fields in the <code>jvm.gc</code> section to allow certain garbage collection parameters to be set.
These include GC logging, setting the collector to use and arbitrary GC arguments.</p>


<h3 id="_enable_gc_logging">Enable GC Logging</h3>
<div class="section">
<p>To enable GC logging set the <code>jvm.gc.logging</code> field to <code>true</code>.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    gc:
      logging: true</markup>

<p>Setting the field to true adds the following JVM arguments to the JVM in the <code>coherence</code> container:</p>

<div class="listing">
<pre>-verbose:gc
-XX:+PrintGCDetails
-XX:+PrintGCTimeStamps
-XX:+PrintHeapAtGC
-XX:+PrintTenuringDistribution
-XX:+PrintGCApplicationStoppedTime
-XX:+PrintGCApplicationConcurrentTime</pre>
</div>

<p>If different GC logging arguments are required then the relevant JVM arguments can be added to either the
<code>jvm.args</code> field or the <code>jvm.gc.args</code> field.</p>

</div>

<h3 id="_set_the_garbage_collector">Set the Garbage Collector</h3>
<div class="section">
<p>The garbage collector to use can be set using the <code>jvm.gc.collector</code> field.
This field can be set to either <code>G1</code>, <code>CMS</code> or <code>Parallel</code>
(the field is case-insensitive, invalid values will be silently ignored).
The default collector set, if none has been specified, will be <code>G1</code>.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 50%;">
<col style="width: 50%;">
</colgroup>
<thead>
</thead>
<tbody>
<tr>
<td class="">Parameter</td>
<td class="">JVM Argument Set</td>
</tr>
<tr>
<td class=""><code>G1</code></td>
<td class=""><code>-XX:+UseG1GC</code></td>
</tr>
<tr>
<td class=""><code>CMS</code></td>
<td class=""><code>-XX:+UseConcMarkSweepGC</code></td>
</tr>
<tr>
<td class=""><code>Parallel</code></td>
<td class=""><code>-XX:+UseParallelGC</code></td>
</tr>
</tbody>
</table>
</div>
<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    gc:
      collector: "G1"</markup>

<p>The example above will add <code>-XX:+UseG1GC</code> to the command line.</p>

</div>

<h3 id="_adding_arbitrary_gc_args">Adding Arbitrary GC Args</h3>
<div class="section">
<p>Any arbitrary GC argument can be added to the <code>jvm.gc.args</code> field.
These arguments will be passed verbatim to the JVM command line.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    gc:
      args:
        - "-XX:MaxGCPauseMillis=200"</markup>

<p>In the example above the <code>-XX:MaxGCPauseMillis=200</code> JVM argument will be added to the command line.</p>

<div class="admonition note">
<p class="admonition-inline">The <code>jvm.gc.args</code> field will add the provided arguments to the end of the command line exactly as they
are in the args list. This field provides the same functionality as <router-link to="#jvm/030_jvm_args.adoc" @click.native="this.scrollFix('#jvm/030_jvm_args.adoc')">JVM Args</router-link>
but sometimes it might be useful to be able to separate the two gorups of arguments in the CRD spec.</p>
</div>
</div>
</div>
</doc-view>
