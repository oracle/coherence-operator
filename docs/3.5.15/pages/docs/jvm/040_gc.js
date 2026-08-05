<doc-view>

<h2 id="_garbage_collector_settings">Garbage Collector Settings</h2>
<div class="section">
<p>The <code>Coherence</code> CRD has fields in the <code>jvm.gc</code> section to allow certain garbage collection parameters to be set.
These setting the collector to use and arbitrary GC arguments.</p>

<div class="admonition important">
<p class="admonition-textlabel">Important</p>
<p ><p>If running Kubernetes on ARM processors and using Coherence images built on Java 8 for ARM,
note that the G1 garbage collector in that version of Java on ARM is marked as experimental.</p>

<p>By default, the Operator configures the Coherence JVM to use G1.
This will cause errors on Arm64 Java 8 JMS unless the JVM option <code>-XX:+UnlockExperimentalVMOptions</code> is
added in the Coherence resource spec.
Alternatively specify a different garbage collector, ideally on a version of Java this old, use CMS.</p>
</p>
</div>

<h3 id="_set_the_garbage_collector">Set the Garbage Collector</h3>
<div class="section">
<p>The garbage collector to use can be set using the <code>jvm.gc.collector</code> field.
This field can be set to either <code>G1</code>, <code>CMS</code> or <code>Parallel</code>
(the field is case-insensitive, invalid values will be silently ignored).</p>

<p>The default collector set, if none has been specified, will be whatever is the default for the JVM being used.</p>


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
<tr>
<td class=""><code>Serial</code></td>
<td class=""><code>-XX:+UseSerialGC</code></td>
</tr>
<tr>
<td class=""><code>ZGC</code></td>
<td class=""><code>-XX:+UseZGC</code></td>
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
      collector: "ZGC"</markup>

<p>The example above will add <code>-XX:+UseZGC</code> to the command line.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>The JVM only allows a single <code>-XX:Use*<strong>*</strong></code> option that sets the garbage collector to use, so the collector should not be
specified in both the <code>spec.jvm.gc.collector</code> field, and in the <code>spec.jvm.args</code> field.</p>
</p>
</div>
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
are in the args list. This field provides the same functionality as <router-link to="/docs/jvm/030_jvm_args">JVM Args</router-link>
but sometimes it might be useful to be able to separate the two gorups of arguments in the CRD spec.</p>
</div>
</div>
</div>
</doc-view>
