<doc-view>

<h2 id="_configure_the_jvm">Configure the JVM</h2>
<div class="section">
<p>There are a number of fields in the <code>CoherenceCluster</code> CRD that can be used to configure the JVM.
These fields are all in the <code>jvm</code> section when configuring a <code>role</code> in the CRD.</p>

<ul class="ulist">
<li>
<p><router-link to="#args" @click.native="this.scrollFix('#args')">JVM Arguments</router-link></p>

</li>
<li>
<p><router-link to="#gc" @click.native="this.scrollFix('#gc')">Garbage Collector Configuration</router-link></p>
<ul class="ulist">
<li>
<p><router-link to="#gc-collector" @click.native="this.scrollFix('#gc-collector')">Configuring the Garbage Collector to Use</router-link></p>

</li>
<li>
<p><router-link to="#gc-args" @click.native="this.scrollFix('#gc-args')">Configuring the Garbage Collector Arguments</router-link></p>

</li>
<li>
<p><router-link to="#gc-logging" @click.native="this.scrollFix('#gc-logging')">Configuring Garbage Collector Logging</router-link></p>

</li>
</ul>
</li>
<li>
<p><router-link to="#memory" @click.native="this.scrollFix('#memory')">Memory Configuration</router-link></p>
<ul class="ulist">
<li>
<p><router-link to="#heap-size" @click.native="this.scrollFix('#heap-size')">Heap Size</router-link></p>

</li>
<li>
<p><router-link to="#metaspace-size" @click.native="this.scrollFix('#metaspace-size')">Metaspace size</router-link></p>

</li>
<li>
<p><router-link to="#stack-size" @click.native="this.scrollFix('#stack-size')">Stack size</router-link></p>

</li>
<li>
<p><router-link to="#nio-size" @click.native="this.scrollFix('#nio-size')">Native Memory Size</router-link></p>

</li>
<li>
<p><router-link to="#nmt" @click.native="this.scrollFix('#nmt')">Native Memory Tracking</router-link></p>

</li>
</ul>
</li>
<li>
<p><router-link to="#oom" @click.native="this.scrollFix('#oom')">Behaviour on Out Of Memory Error</router-link></p>

</li>
<li>
<p><router-link to="#useContainerLimits" @click.native="this.scrollFix('#useContainerLimits')">Container Resource Limits</router-link></p>

</li>
<li>
<p><router-link to="#flightRecorder" @click.native="this.scrollFix('#flightRecorder')">Flight Recorder</router-link></p>

</li>
<li>
<p><router-link to="#diagnosticsVolume" @click.native="this.scrollFix('#diagnosticsVolume')">Diagnostic Volume</router-link></p>

</li>
<li>
<p><router-link to="">JVM Debug Arguments</router-link></p>

</li>
</ul>
<p>The following sections describe the different JVM configuration options available in the <code>CoherenceCluster</code> CRD.
These CRD fields all result in the addition or omission of various JVM arguments.
A number of arguments are always passed to the Coherence container&#8217;s JVM:</p>

<markup


>-XX:HeapDumpPath=/jvm/${POD_NAME}/${POD_UID}/heap-dumps/${POD_NAME}-${POD_UID}.hprof  <span class="conum" data-value="1" />
-XX:ErrorFile=/jvm/${POD_NAME}/${POD_UID}/hs-err-${POD_NAME}-${POD_UID}.log           <span class="conum" data-value="2" />
-Dcoherence.ttl=0                                                                     <span class="conum" data-value="3" />
-XshowSettings:all
-XX:+PrintCommandLineFlags
-XX:+PrintFlagsFinal
-XX:+UnlockDiagnosticVMOptions
-XX:+UnlockCommercialFeatures
-XX:+UnlockExperimentalVMOptions</markup>

<ul class="colist">
<li data-value="1">Any heap dumps created by the JVM when an out of memory error occurs will be written to a file called
<code>/jvm/${POD_NAME}/${POD_UID}/heap-dumps/${POD_NAME}-${POD_UID}.hprof</code></li>
<li data-value="2">Any error files created by a JVM crash will be written to a file called
<code>/jvm/${POD_NAME}/${POD_UID}/hs-err-${POD_NAME}-${POD_UID}.log</code></li>
<li data-value="3">Coherence multicast discovery is disabled as multicast cannot be relied on in containers</li>
</ul>
<p>The <code>/jvm</code> root directory used for heap dumps and error files can be
<router-link to="#diagnosticsVolume" @click.native="this.scrollFix('#diagnosticsVolume')">mounted to an external volume</router-link> to allow easier access to these files.</p>

<v-divider class="my-5"/>
</div>

<h2 id="args">JVM Arguments</h2>
<div class="section">
<p>The <code>jvm.args</code> field is a string array of arbitrary JVM options. Any valid JVM option or system property argument may be
passed to the JVM in the Coherence container by setting the value in this field.</p>


<h3 id="_setting_the_jvm_arguments_for_the_implicit_role">Setting the JVM Arguments for the Implicit Role</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with a single implicit role the <code>args</code> is set in the <code>spec.jvm</code> section of
the configuration. For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    args:
      - "-XX:G1HeapRegionSize=16m"  <span class="conum" data-value="1" />
      - "-Dfoo=bar"</markup>

<ul class="colist">
<li data-value="1">The <code>-XX:G1HeapRegionSize=16m</code> JVM option and the <code>-Dfoo=bar</code> system property will be passed as arguments to the
JVM for the implicit storage role.</li>
</ul>
</div>

<h3 id="_setting_the_jvm_arguments_for_explicit_roles">Setting the JVM Arguments for Explicit Roles</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles the <code>args</code> are set in the <code>jvm</code> section of
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
        args:
          - "-XX:G1HeapRegionSize=16m"  <span class="conum" data-value="1" />
          - "-Dcoherence.pof.config=storage-pof-config.xml"
    - role: proxy
      jvm:
        args:
          - "-XX:MaxGCPauseMillis=500"  <span class="conum" data-value="2" />
          - "-Dcoherence.pof.config=proxy-pof-config.xml"</markup>

<ul class="colist">
<li data-value="1">The <code>-XX:G1HeapRegionSize=16m -Dcoherence.pof.config=storage-pof-config.xml</code> arguments will be passed to the JVM for
the explicit <code>data</code> role.</li>
<li data-value="2">The <code>-XX:MaxGCPauseMillis=500 coherence.pof.config=proxy-pof-config.xml</code> argument will be passed to the JVM for the
explicit <code>proxy</code> role.</li>
</ul>
</div>

<h3 id="_setting_the_jvm_arguments_for_explicit_roles_with_a_default">Setting the JVM Arguments for Explicit Roles with a Default</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles a default <code>args</code> value can be set in the
<code>CoherenceCluster</code> <code>spec</code> section that will apply to all of the roles in the <code>roles</code> list.</p>

<div class="admonition note">
<p class="admonition-inline">Any <code>args</code> set explicitly in the <code>jvm.args</code> field for a <code>role</code> will be <strong>merged</strong> with those in the defaults
section.</p>
</div>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    args:
      - "-XX:MaxGCPauseMillis=500"  <span class="conum" data-value="1" />
      - "-XX:G1HeapRegionSize=16m"
  roles:
    - role: data                    <span class="conum" data-value="2" />
      jvm:
        args:
          - "-XX:+AggressiveHeap"
    - role: proxy                   <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The default JVM <code>args</code> of <code>-XX:MaxGCPauseMillis=500</code> and <code>-XX:G1HeapRegionSize=16m</code> will be passed to the JVM
for <strong>all</strong> roles.</li>
<li data-value="2">The <code>data</code> role adds an additional argument <code>-XX:+AggressiveHeap</code> so the JVM will be passed three arguments:
<code>-XX:MaxGCPauseMillis=500 -XX:G1HeapRegionSize=16m -XX:+AggressiveHeap</code></li>
<li data-value="3">The <code>proxy</code> role does not specify any additional args so will just use the two default JVM arguments
<code>-XX:MaxGCPauseMillis=500 -XX:G1HeapRegionSize=16m</code></li>
</ul>
<v-divider class="my-5"/>
</div>
</div>

<h2 id="gc">Garbage Collector Configuration</h2>
<div class="section">
<p>The <code>CoherenceCluster</code> CRD allows garbage collector settings to be applied to the Coherence JVMs. Whilst any GC
parameters could actually be applied using the <code>jvm.args</code> field these GC specific fields allow options to be set
without having to look up and remember specific GC options. The garbage collector configuration is set in the
<code>jvm.gc</code> section of the CRD.</p>

<ul class="ulist">
<li>
<p><router-link to="#gc-collector" @click.native="this.scrollFix('#gc-collector')">Configuring the Garbage Collector to Use</router-link></p>

</li>
<li>
<p><router-link to="#gc-args" @click.native="this.scrollFix('#gc-args')">Configuring the Garbage Collector Arguments</router-link></p>

</li>
<li>
<p><router-link to="#gc-logging" @click.native="this.scrollFix('#gc-logging')">Configuring Garbage Collector Logging</router-link></p>

</li>
</ul>

<h3 id="gc-collector">Configuring the Garbage Collector to Use</h3>
<div class="section">
<p>The <code>CoherenceCluster</code> CRD supports setting the garbage collectors to use automatically. The supported collectors are
<code>G1</code>, <code>CMS</code>, <code>Parallel</code> or the JVM default.
The garbage collector to use is set using the <code>jvm.gc.collector</code> field.
The value sould be one of:</p>


<div class="table__overflow elevation-1 ">
<table class="datatable table">
<colgroup>
<col style="width: 50%;">
<col style="width: 50%;">
</colgroup>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td><code>G1</code></td>
<td>Enables the G1 garbage collector by adding the <code>-XX:+UseG1GC</code> JVM option</td>
</tr>
<tr>
<td><code>CMS</code></td>
<td>Enables the CMS garbage collector by adding the <code>-XX:+UseConcMarkSweepGC</code> JVM option</td>
</tr>
<tr>
<td><code>Parallel</code></td>
<td>Enables the parallel garbage collector by adding the <code>-XX:+UseParallelGC</code> JVM option</td>
</tr>
<tr>
<td><code>Default</code></td>
<td>Deos not add any extra GC parameter; the JVM will use its default garbage collector</td><td>&#8230;&#8203;</td>
</tr>
</tbody>
</table>
</div>
<p>The <code>jvm.gc.collector</code> value is not case sensitive so for example <code>CMS</code>, <code>cms</code> and <code>CmS</code> will all enable the <code>CMS</code>
collector.
The contents of the <code>jvm.gc.collector</code> are not validated, any value other than those described above will be treated
as <code>Default</code> enabling the JVMs default garbage collector.</p>

<div class="admonition note">
<p class="admonition-inline">The default value for <code>jvm.gc.collector</code> is <code>G1</code> which will enable the recommended G1 garbage collector.</p>
</div>

<h4 id="_setting_the_garbage_collector_for_the_implicit_role">Setting the Garbage Collector for the Implicit Role</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with a single implicit role the garbage collector to use is set in the <code>spec</code>
section of the yaml. For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    gc:
      collector: CMS  <span class="conum" data-value="1" /></markup>

<p>The implicit storage role will use the <code>CMS</code> garbage collector.</p>

</div>

<h4 id="_setting_the_garbage_collector_for_explicit_roles">Setting the Garbage Collector for Explicit Roles</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles the garbage collector to use is set in the
<code>jvm.gc.collector</code> section for each role.</p>

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
      jvm:
        gc:
          collector: G1   <span class="conum" data-value="1" />
    - role: proxy
      jvm:
        gc:
          collector: CMS  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The JVMs for the <code>data</code> role will use the G1 garbage collector</li>
<li data-value="2">The JVMs for the <code>proxy</code> role will use the CMS garbage collector</li>
</ul>
</div>

<h4 id="_setting_the_garbage_collector_for_explicit_roles_with_a_default">Setting the Garbage Collector for Explicit Roles with a Default</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles a default garbage collector can be set in the
<code>spec.jvm.gc.collector</code> field of the CRD. This value can then be overridden for specific roles in the <code>jvm.gc.collector</code>
field for each role in the <code>roles</code> list.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    gc:
      collector: CMS     <span class="conum" data-value="1" />
  roles:
    - role: data         <span class="conum" data-value="2" />
      jvm:
        gc:
          collector: G1
    - role: proxy        <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The default garbage collector us set to <code>CMS</code> which will be used by all roles in the <code>roles</code> list that do not
set a specific collector to use.</li>
<li data-value="2">The <code>data</code> role overrides the default collector so that the JVMs for the <code>data</code> role will use the G1 garbage
collector</li>
<li data-value="3">The <code>proxy</code> role does not specify a collector to use so that JVMs for the <code>proxy</code> role will use the CMS garbage
collector</li>
</ul>
</div>
</div>

<h3 id="gc-args">Configuring Garbage Collector Arguments</h3>
<div class="section">
<p>Arbitrary GC arguments can be passed to the JVM in the <code>jvm.gc.args</code> field. This field is a string array where each
argument to be passed to the JVM is a separate string value.</p>


<h4 id="_setting_garbage_collector_arguments_for_the_implicit_role">Setting Garbage Collector Arguments for the Implicit Role</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with a single implicit role the GC arguments are set in the <code>spec.jvm.gc.args</code> field.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    gc:
      args:                           <span class="conum" data-value="1" />
        - "-XX:MaxGCPauseMillis=500"
        - "-XX:G1ReservePercent=20"</markup>

<ul class="colist">
<li data-value="1">The implicit storage role will have the additional GC arguments <code>-XX:MaxGCPauseMillis=500</code> and
<code>-XX:G1ReservePercent=20</code> passed to the JVM.</li>
</ul>
</div>

<h4 id="_setting_garbage_collector_arguments_for_explicit_roles">Setting Garbage Collector Arguments for Explicit Roles</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles the GC arguments are set in the <code>jvm.gc.args</code> field
for each role in the <code>roles</code> list.
For example:</p>

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
        gc:
          args:                            <span class="conum" data-value="1" />
            - "-XX:MaxGCPauseMillis=500"
            - "-XX:G1ReservePercent=20"
    - role: proxy
      jvm:
        gc:
          args:                            <span class="conum" data-value="2" />
            - "-XX:MaxGCPauseMillis=1000"</markup>

<ul class="colist">
<li data-value="1">The explicit <code>data</code> role will have the additional GC arguments <code>-XX:MaxGCPauseMillis=500</code> and
<code>-XX:G1ReservePercent=20</code> passed to the JVM.</li>
<li data-value="2">The explicit <code>proxy</code> role will have the additional GC argument <code>-XX:MaxGCPauseMillis=1000</code> passed to the JVM.</li>
</ul>
</div>

<h4 id="_setting_garbage_collector_arguments_for_explicit_roles_with_a_default">Setting Garbage Collector Arguments for Explicit Roles with a Default</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles a default GC arguments are set in the
<code>spec.jvm.gc.args</code> field and will be applied to all roles in the roles list that do not set specific GC arguments.</p>

<div class="admonition note">
<p class="admonition-inline">GC arguments set for explicit roles override the defaults. The role&#8217;s GC arguments are <strong>not merged</strong> with the
default GC arguments.</p>
</div>
<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    gc:
      args:                                <span class="conum" data-value="1" />
        - "-XX:MaxGCPauseMillis=500"
        - "-XX:G1ReservePercent=20"
  roles:
    - role: data                           <span class="conum" data-value="2" />
    - role: proxy
      jvm:
        gc:
          args:                            <span class="conum" data-value="3" />
            - "-XX:MaxGCPauseMillis=1000"</markup>

<ul class="colist">
<li data-value="1">The default GC arguments are <code>-XX:MaxGCPauseMillis=500</code> and <code>-XX:G1ReservePercent=20</code></li>
<li data-value="2">The <code>data</code> role does not specify any GC arguments so the default arguments of <code>-XX:MaxGCPauseMillis=500</code> and
<code>-XX:G1ReservePercent=20</code> will be passed to the <code>data</code> role JVMs.</li>
<li data-value="3">The <code>proxy</code> role specifies the GC arguments <code>-XX:MaxGCPauseMillis=1000</code> which will <strong>override</strong> the defaults so only
<code>-XX:MaxGCPauseMillis=1000</code> will be passed to the <code>proxy</code> role JVMs.</li>
</ul>
</div>
</div>

<h3 id="gc-logging">Configuring Garbage Collector Logging</h3>
<div class="section">
<p>The Coherence documentation recommends enabling GC logging for Coherence JVMs. To this end the <code>CoherenceCluster</code> CRD
has a boolean field <code>jvm.gc.logging</code> to enable or disable default GC logging JVM arguments. By default the value of this
field is set to <code>true</code> if it is not specified for a <code>CoherenceCluster</code>.</p>

<p>The following GC logging JVM arguments are added if the <code>jvm.gc.logging</code> field is omitted or explicitly set to <code>true</code>:</p>

<markup


>-verbose:gc
-XX:+PrintGCDetails
-XX:+PrintGCTimeStamps
-XX:+PrintHeapAtGC
-XX:+PrintTenuringDistribution
-XX:+PrintGCApplicationStoppedTime
-XX:+PrintGCApplicationConcurrentTime</markup>


<h4 id="_configuring_garbage_collector_logging_for_the_implicit_role">Configuring Garbage Collector Logging for the Implicit Role</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with a single implicit role GC logging can be enabled or disabled in the <code>spec</code>
section of the yaml.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    gc:
      logging: true  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The implicit storage role has GC logging explicitly enabled so that the JVM arguments listed above will
be added to the JVM&#8217;s command line.</li>
</ul>
</div>

<h4 id="_configuring_garbage_collector_logging_for_explicit_roles">Configuring Garbage Collector Logging for Explicit Roles</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles GC logging can be enabled or disabled in the
<code>jvm.gc.logging</code> field of each role in the <code>roles</code> list.
section of the yaml
For example:</p>

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
        gc:
          logging: true   <span class="conum" data-value="1" />
    - role: proxy
      jvm:
        gc:
          logging: false  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role has GC logging explicitly enabled so that the JVM arguments listed above will be added to the
JVM&#8217;s command line</li>
<li data-value="2">The <code>proxy</code> role has GC logging explicitly disabled so that the JVM arguments listed above will not be added to
the JVM&#8217;s command line</li>
</ul>
</div>

<h4 id="_configuring_garbage_collector_logging_for_explicit_roles_with_a_default">Configuring Garbage Collector Logging for Explicit Roles with a Default</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles a default GC logging setting can be specified in the
<code>spec</code> section of the CRD which can then be overridden for individual roles in the <code>roles</code> list.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    gc:
      logging: false      <span class="conum" data-value="1" />
  roles:
    - role: data          <span class="conum" data-value="2" />
      jvm:
        gc:
          logging: true
    - role: proxy         <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The default value for <code>jvm.gc.logging</code> is false, which will disable GC logging.</li>
<li data-value="2">The <code>data</code> role overrides the default and sets GC logging to <code>true</code></li>
<li data-value="3">The <code>proxy</code> role does not specify a value for <code>jvm.gc.logging</code> so it will use the default, which will disable GC
logging.</li>
</ul>
<v-divider class="my-5"/>
</div>
</div>
</div>

<h2 id="memory">Memory Configuration</h2>
<div class="section">
<p>The JVM has a number of options that can be set to fix the size of different memory regions. The <code>CoherenceCluster</code> CRD
provides fields to set that most common values. None of these fields have default values so if they are not specified
the JVMs default behaviour will apply.</p>

<p>The memory options that can be configured are:</p>

<ul class="ulist">
<li>
<p><router-link to="#heap-size" @click.native="this.scrollFix('#heap-size')">Heap Size</router-link></p>

</li>
<li>
<p><router-link to="#metaspace-size" @click.native="this.scrollFix('#metaspace-size')">Metaspace size</router-link></p>

</li>
<li>
<p><router-link to="#stack-size" @click.native="this.scrollFix('#stack-size')">Stack size</router-link></p>

</li>
<li>
<p><router-link to="#nio-size" @click.native="this.scrollFix('#nio-size')">Max Native Memory</router-link></p>

</li>
<li>
<p><router-link to="#nmt" @click.native="this.scrollFix('#nmt')">Native Memory Tracking</router-link></p>

</li>
<li>
<p><router-link to="#oom" @click.native="this.scrollFix('#oom')">Behaviour on Out Of Memory Error</router-link></p>

</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">If the <code>Pod</code> resource limits are being set to limit memory usage of a <code>Pod</code> it is recommended that some of the JVM
memory regions are fixed to ensure that the JVM does not exceed the container&#8217;s resource limits in a JVM before Java 10.
Prior to Java 10 the JVM could see all of the memory available to a machine regardless of any Pod limits.
The JVM could then easily attempt to consume more memory that the <code>Pod</code> or <code>Container</code> was allowed and consequently
crashing the <code>Pod</code>. With Coherence images that use a version of Java above 10 this issue is less of a problem.
Even so if using the <code>resources</code> section of the configuration to limit a <code>Pod</code> or <code>Containers</code> memory it is a good idea
to limit the JVM heap. Also see <router-link to="#useContainerLimits" @click.native="this.scrollFix('#useContainerLimits')">the useContainerLimits setting</router-link>.</p>
</div>

<h3 id="heap-size">JVM Heap Size</h3>
<div class="section">
<p>It is good practice to fix the Coherence JVM heap size and to set both the JVM <code>-Xmx</code> and <code>-Xms</code> options to the same
value.
The heap size of the JVM can be configured for roles in the <code>jvm.heapSize</code> field of a role spec. If the <code>heapSize</code> value
is configured then that value is applied to bot the JVMs minimum and maximum heap sizes (i.e. used to set both
<code>-Xms</code> and -<code>Xmx</code>).</p>

<p>The format of the value of the <code>heapSize</code> field is any valid value that can be used when setting the <code>-Xmx</code> JVM option,
for example <code>10G</code> would set a 10 GB heap.</p>


<h4 id="_setting_the_jvm_heap_size_for_the_implicit_role">Setting the JVM Heap Size for the Implicit Role</h4>
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
    memory:
      heapSize: 10g <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The Coherence JVM for the implicit role defined above will have a 10 GB heap.
Equivalent to passing <code>-Xms10g -Xmx10g</code> to the JVM.</li>
</ul>
</div>

<h4 id="_setting_the_jvm_heap_size_for_explicit_roles">Setting the JVM Heap Size for Explicit Roles</h4>
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
        memory:
          heapSize: 10g   <span class="conum" data-value="1" />
    - role: proxy
      jvm:
        memory:
          heapSize: 500m  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The Coherence JVM for the <code>data</code> role defined above will have a 10 GB heap.
Equivalent to passing <code>-Xms10g -Xmx10g</code> to the JVM.</li>
<li data-value="2">The Coherence JVM for the <code>proxy</code> role defined above will have a 500 MB heap.
Equivalent to passing <code>-Xms500m -Xmx500m</code> to the JVM.</li>
</ul>
</div>

<h4 id="_setting_the_jvm_heap_size_for_explicit_roles_with_a_default">Setting the JVM Heap Size for Explicit Roles with a Default</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles a default <code>heapSize</code> value can be set in the
<code>CoherenceCluster</code> <code>spec</code> section that will apply to all of the roles in the <code>roles</code> list unless specifically
overridden by a role&#8217;s <code>jvm.heapSize</code> field. For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    memory:
      heapSize: 500m     <span class="conum" data-value="1" />
  roles:
    - role: data
      jvm:
        memory:
          heapSize: 10g  <span class="conum" data-value="2" />
    - role: proxy        <span class="conum" data-value="3" />
    - role: web          <span class="conum" data-value="4" /></markup>

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

<h3 id="metaspace-size">JVM Metaspace Size</h3>
<div class="section">
<p>The metaspace size is the amount of native memory that can be allocated for class metadata. By default the JVM does not
limit this size. When running in size limited containers this size may be set to ensure that the JVM does not cause the
container to exceed its configured memory limits. The metaspace size is set using the <code>jvm.memory.metaspaceSize</code> field.
Setting this field causes the <code>-XX:MetaspaceSize</code> and <code>-XX:MaxMetaspaceSize</code> JVM arguments to be set.
There is no default value for the <code>metaspaceSize</code> field so if it is omitted the JVMs default behaviour will control the
metaspace size.</p>


<h4 id="_configuring_the_jvm_metaspace_size_for_the_implicit_role">Configuring the JVM Metaspace Size for the Implicit Role</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with a single implicit role the metaspace size can be set in the <code>spec</code> section of
the CRD.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    memory:
      metaspaceSize: 256m  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The metaspace size will for the implicit storage role will be set to <code>256m</code> by setting the JVM arguments
<code>-XX:MetaspaceSize=256m -XX:MaxMetaspaceSize=256m</code></li>
</ul>
</div>

<h4 id="_configuring_the_jvm_metaspace_size_for_explicit_roles">Configuring the JVM Metaspace Size for Explicit Roles</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles the metaspace size can be set in the
<code>jvm.memory.metaspaceSize</code> field for each role in the <code>roles</code> list.
For example:</p>

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
        memory:
          metaspaceSize: 256m  <span class="conum" data-value="1" />
    - role: proxy
      jvm:
        memory:
          metaspaceSize: 512m  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The metaspace size will for the <code>data</code> role will be set to <code>256m</code> by setting the JVM arguments
<code>-XX:MetaspaceSize=256m -XX:MaxMetaspaceSize=256m</code></li>
<li data-value="2">The metaspace size will for the <code>proxy</code> role will be set to <code>512m</code> by setting the JVM arguments
<code>-XX:MetaspaceSize=512m -XX:MaxMetaspaceSize=512m</code></li>
</ul>
</div>

<h4 id="_configuring_the_jvm_metaspace_size_for_explicit_roles_with_a_default">Configuring the JVM Metaspace Size for Explicit Roles with a Default</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles a default
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    memory:
      metaspaceSize: 512m      <span class="conum" data-value="1" />
  roles:
    - role: data               <span class="conum" data-value="2" />
      jvm:
        memory:
          metaspaceSize: 256m
    - role: proxy              <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The metaspace size will for the <code>data</code> role will be set to <code>256m</code> by setting the JVM arguments
<code>-XX:MetaspaceSize=256m -XX:MaxMetaspaceSize=256m</code></li>
<li data-value="2">The metaspace size will for the <code>proxy</code> role will be set to <code>512m</code> by setting the JVM arguments
<code>-XX:MetaspaceSize=512m -XX:MaxMetaspaceSize=512m</code></li>
</ul>
</div>
</div>

<h3 id="stack-size">JVM Stack Size</h3>
<div class="section">
<p>Setting the stack size sets the thread stack size (in bytes) used by the JVM. The stack size is configured in for roles
in a <code>CoherenceCluster</code> by setitng the <code>jvm.memory`stackSize</code> field. Setting this fields sets the <code>-Xss</code> JVM argument.
Omitting this fields does not set the <code>-Xss</code> argument leaving the JVM to its default configuration which sets the stack
size based on the O/S being used.</p>


<h4 id="_configuring_the_jvm_stack_size_for_the_implicit_role">Configuring the JVM Stack Size for the Implicit Role</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with a single implicit role the stack size can be set in the <code>spec</code> section of CRD.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    memory:
      stackSize: 1024k  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The stack size for the implicit storage role is set to <code>1024k</code> which will cause the <code>-Xss1024k</code> argument to be
passed to the JVM.</li>
</ul>
</div>

<h4 id="_configuring_the_jvm_stack_size_for_explicit_roles">Configuring the JVM Stack Size for Explicit Roles</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles
For example:</p>

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
        memory:
          stackSize: 512k   <span class="conum" data-value="1" />
    - role: proxy
      jvm:
        memory:
          stackSize: 1024k  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The stack size for the <code>data</code> role is set to <code>512k</code> which will cause the <code>-Xss512k</code> argument to be passed to the JVM.</li>
<li data-value="2">The stack size for the <code>proxy</code> role is set to <code>1024k</code> which will cause the <code>-Xss1024k</code> argument to be passed to the JVM.</li>
</ul>
</div>

<h4 id="_configuring_the_jvm_stack_size_for_explicit_roles_with_a_default">Configuring the JVM Stack Size for Explicit Roles with a Default</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles a default stack size can be set in the <code>spec</code> section
of the yaml that will apply to all roles in the <code>roles</code> list unless overridden for a specific role.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    memory:
      stackSize: 1024k      <span class="conum" data-value="1" />
  roles:
    - role: data            <span class="conum" data-value="2" />
      jvm:
        memory:
          stackSize: 512k
    - role: proxy           <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The default stack size is set to <code>1024k</code> which will cause the <code>-Xss1024k</code> argument to be passed to the JVM for all
roles in the <code>roles</code> list unless overridden.</li>
<li data-value="2">The stack size for the <code>data</code> role is specifically set to <code>512k</code> which will cause the <code>-Xss512k</code> argument to be passed
to the JVMs for the <code>data</code> role.</li>
<li data-value="3">The stack size for the <code>proxy</code> role is not configured so the default value will be used which will cause the
<code>-Xss1024k</code> argument to be passed to the JVMs for the <code>proxy</code> role.</li>
</ul>
</div>
</div>

<h3 id="nio-size">JVM Native Memory Size</h3>
<div class="section">
<p>Native memory is used by the JVM and by Coherence for a number of reasons. In a resource limited container it may be
useful to limit the amount of nio memory available to the JVM to stop the JVM exceeding the containers memory limits.
The nio size is set using the <code>jvm.directMemorySize</code> field which will cause the <code>-XX:MaxDirectMemorySize</code> JVM argument
to be set. There is no default value for the <code>jvm.directMemorySize</code> field so if it is omitted the JVM&#8217;s default size
will be used.</p>


<h4 id="_configuring_the_jvm_native_memory_size_for_the_implicit_role">Configuring the JVM Native Memory Size for the Implicit Role</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with a single implicit role</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    memory:
      directMemorySize: 2g  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">the maximum direct memory size for the implicit storage role is set to <code>2g</code> causing the <code>-XX:MaxDirectMemorySize=2g</code>
argument to be passed to the JVM.</li>
</ul>
</div>

<h4 id="_configuring_the_jvm_native_memory_size_for_explicit_roles">Configuring the JVM Native Memory Size for Explicit Roles</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles
For example:</p>

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
        memory:
          directMemorySize: 2g  <span class="conum" data-value="1" />
    - role: proxy
      jvm:
        memory:
          directMemorySize: 1g  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">the maximum direct memory size for the <code>data</code> role is set to <code>2g</code> causing the <code>-XX:MaxDirectMemorySize=2g</code>
argument to be passed to the JVM.</li>
<li data-value="2">the maximum direct memory size for the <code>proxy</code> role is set to <code>1g</code> causing the <code>-XX:MaxDirectMemorySize=1g</code>
argument to be passed to the JVM.</li>
</ul>
</div>

<h4 id="_configuring_the_jvm_native_memory_size_for_explicit_roles_with_a_default">Configuring the JVM Native Memory Size for Explicit Roles with a Default</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles a default
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    memory:
      directMemorySize: 1g      <span class="conum" data-value="1" />
  roles:
    - role: data                <span class="conum" data-value="2" />
      jvm:
        memory:
          directMemorySize: 2g
    - role: proxy               <span class="conum" data-value="3" /></markup>

</div>
</div>

<h3 id="nmt">Native Memory Tracking</h3>
<div class="section">
<p>The Native memory tracking mode can be configured for JVMs using the <code>jvm.memory.nativeMemoryTracking</code> field to track
JVM nio memory usage, which can be useful when  debugging nio memory issues. Setting the <code>nativeMemoryTracking</code> value
causes the <code>-XX:NativeMemoryTracking</code> JVM argument to be set.
If the <code>jvm.memory.nativeMemoryTracking</code> field is not specified a value of <code>summary</code> is used passing
<code>-XX:NativeMemoryTracking=summary</code> to the JVM.
See the <a id="" title="" target="_blank" href="https://docs.oracle.com/javase/8/docs/technotes/guides/troubleshoot/tooldescr007.html">native memory tracking</a>
documentation.</p>


<h4 id="_configuring_native_memory_tracking_for_the_implicit_role">Configuring Native Memory Tracking for the Implicit Role</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with a single implicit role native memory tracking can be configured in the <code>spec</code>
section of the yaml.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    memory:
      nativeMemoryTracking: detail  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The native memory tracking mode for the JVMs in the implicit storage role will be set to <code>detail</code> causing the
<code>-XX:NativeMemoryTracking=detail</code> to be passed to the JVMs.</li>
</ul>
</div>

<h4 id="_configuring_native_memory_tracking_for_explicit_roles">Configuring Native Memory Tracking for Explicit Roles</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles native memory tracking can br configured specifically
for each role in the <code>roles</code> list.
For example:</p>

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
        memory:
          nativeMemoryTracking: detail   <span class="conum" data-value="1" />
    - role: proxy
      jvm:
        memory:
          nativeMemoryTracking: summary  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The native memory tracking mode for the JVMs in the <code>data</code> role will be set to <code>detail</code> causing the
<code>-XX:NativeMemoryTracking=detail</code> to be passed to the JVMs.</li>
<li data-value="2">The native memory tracking mode for the JVMs in the <code>proxy</code> role will be set to <code>summary</code> causing the
<code>-XX:NativeMemoryTracking=summary</code> to be passed to the JVMs.</li>
</ul>
</div>

<h4 id="_configuring_native_memory_tracking_for_explicit_roles_with_a_default">Configuring Native Memory Tracking for Explicit Roles with a Default</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles a default native memory tracking mode can be set in
the <code>spec</code> section which will apply to all roles in the <code>roles</code> list unless specifically overridden for a role.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    memory:
      nativeMemoryTracking: off         <span class="conum" data-value="1" />
  roles:
    - role: data                        <span class="conum" data-value="2" />
      jvm:
        memory:
          nativeMemoryTracking: detail
    - role: proxy                       <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The default native memory tracking mode is set to <code>off</code> for all roles in the <code>roles</code> list unless specifically
overridden. This will cause the <code>-XX:NativeMemoryTracking=off</code> to be passed to the JVMs.</li>
<li data-value="2">The native memory tracking mode is specifically set to <code>detail</code> for the <code>data</code> role causing the
<code>-XX:NativeMemoryTracking=detail</code> to be passed to the JVMs in the <code>data</code> role.</li>
<li data-value="3">The native memory tracking mode is not set for the <code>proxy</code> role so it will use the default value of <code>off</code> causing
the <code>-XX:NativeMemoryTracking=off</code> to be passed to the JVMs in the <code>proxy</code> role.</li>
</ul>
</div>
</div>

<h3 id="oom">JVM Behaviour on Out Of Memory</h3>
<div class="section">
<p>It is an important recommendation in the Coherence documentation to specifically set the behaviour of a JVM when it
encounters an out of memory error. The JVM should be set to exit and generate a heap dump. A JVM that encounters an OOM
error is left in an undefined state and this can cause a Coherence cluster to become unstable if the JVM does not exit.
Generating a heap dump is useful to diagnose why the JVM had the OOM error.</p>

<p>There are two boolean fields in the <code>CoherenceCluster</code> CRD that control this behaviour:</p>

<ul class="ulist">
<li>
<p><code>jvm.memory.onOutOfMemory.exit</code> which determines whether the JVM will exit on an OOM error; the default value if
the field is not specified is <code>true</code>.
A value of <code>true</code> causes the <code>-XX:+ExitOnOutOfMemoryError</code> argument to be passed to the JVM.</p>

</li>
<li>
<p><code>jvm.memory.onOutOfMemory.heapDump</code> which determines whether the JVM will generate a heap dump on an OOM error; the
default value if the field is not specified is <code>true</code>.</p>

</li>
</ul>
<p>Heap dumps will be written to a file <code>/jvm/${POD_NAME}/${POD_UID}/heap-dumps/${POD_NAME}-${POD_UID}.hprof</code>. The root
<code>/jvm</code> directory can be mapped to an external volume for easier access to the heap dumps
(see: <router-link to="#diagnosticsVolume" @click.native="this.scrollFix('#diagnosticsVolume')">setting the disgnostic volume</router-link>)</p>


<h4 id="_configuring_oom_behaviour_for_the_implicit_role">Configuring OOM Behaviour for the Implicit Role</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with a single implicit role native memory tracking can be configured in the <code>spec</code>
section of the yaml.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    memory:
      onOutOfMemory:
        exit: true      <span class="conum" data-value="1" />
        heapDump: true  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The implicit storage role will exit if an out of memory error occurs, the <code>-XX:+ExitOnOutOfMemoryError</code> argument
will be passed to the JVM</li>
<li data-value="2">The implicit storage role will generate a heap dump if an out of memory error occurs, the
<code>-XX:+HeapDumpOnOutOfMemoryError"</code> argument will be passed to the JVM</li>
</ul>
</div>

<h4 id="_configuring_oom_behaviour_for_explicit_roles">Configuring OOM Behaviour for Explicit Roles</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles native memory tracking can br configured specifically
for each role in the <code>roles</code> list.
For example:</p>

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
        memory:
          onOutOfMemory:
            exit: true       <span class="conum" data-value="1" />
            heapDump: true   <span class="conum" data-value="2" />
    - role: proxy
      jvm:
        memory:
          onOutOfMemory:
            exit: false      <span class="conum" data-value="3" />
            heapDump: false  <span class="conum" data-value="4" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role will exit if an out of memory error occurs, the <code>-XX:+ExitOnOutOfMemoryError</code> argument
will be passed to the JVM</li>
<li data-value="2">The <code>data</code> role will generate a heap dump if an out of memory error occurs, the
<code>-XX:+HeapDumpOnOutOfMemoryError"</code> argument will be passed to the JVM</li>
<li data-value="3">The <code>proxy</code> role will not exit if an out of memory error occurs, the <code>-XX:+ExitOnOutOfMemoryError</code> argument
will be not passed to the JVM</li>
<li data-value="4">The <code>proxy</code> role will not generate a heap dump if an out of memory error occurs, the
<code>-XX:+HeapDumpOnOutOfMemoryError"</code> argument will not be passed to the JVM</li>
</ul>
</div>

<h4 id="_configuring_oom_behaviour_for_explicit_roles_with_a_default">Configuring OOM Behaviour for Explicit Roles with a Default</h4>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles a default native memory tracking mode can be set in
the <code>spec</code> section which will apply to all roles in the <code>roles</code> list unless specifically overridden for a role.
For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    memory:
      onOutOfMemory:
        exit: false          <span class="conum" data-value="1" />
        heapDump: false      <span class="conum" data-value="2" />
  roles:
    - role: data
      jvm:
        memory:
          onOutOfMemory:
            exit: true       <span class="conum" data-value="3" />
            heapDump: true   <span class="conum" data-value="4" />
    - role: proxy            <span class="conum" data-value="5" /></markup>

<ul class="colist">
<li data-value="1">The default setting for exit on out of memory error is <code>false</code></li>
<li data-value="2">The default setting for generating a heap dump on out of memory error is <code>false</code></li>
<li data-value="3">The <code>data</code> role overrides the default <code>jvm.memory.onOutOfMemory.exit</code> value to <code>true</code> and will exit if an out of
memory error occurs, the <code>-XX:+ExitOnOutOfMemoryError</code> argument will be passed to the JVM</li>
<li data-value="4">The <code>data</code> role overrides the default <code>jvm.memory.onOutOfMemory.heapDump</code> value to <code>true</code> and will generate a heap
dump if an out of memory error occurs, <code>-XX:+HeapDumpOnOutOfMemoryError"</code> argument will be passed to the JVM</li>
<li data-value="5">The <code>proxy</code> role does not specify any values for <code>jvm.memory.onOutOfMemory.exit</code> or
<code>jvm.memory.onOutOfMemory.heapDump</code> so it will use the default values of <code>false</code>, the <code>-XX:+ExitOnOutOfMemoryError</code> and
<code>-XX:+HeapDumpOnOutOfMemoryError"</code> arguments will not be passed to the JVM</li>
</ul>
<v-divider class="my-5"/>
</div>
</div>
</div>

<h2 id="useContainerLimits">Container Resource Limits</h2>
<div class="section">
<p>When running JVMs inside containers it is recommended to configure the JVM to respect the memory and CPU resource limits
that are configured for the container. This is especially important in Kubernetes where the <code>Pod</code> may be terminated if a
container exceeds the configured resource limits. The <code>jvm.useContainerLimits</code> field is used to either add or omit the
<code>-XX:+UseContainerSupport</code> JVM argument. If <code>useContainerLimits</code> is set to <code>true</code> then <code>-XX:+UseContainerSupport</code> is
added to the JVM arguments, if <code>useContainerLimits</code> is set to <code>false</code> then <code>-XX:+UseContainerSupport</code> is not
added to the JVM arguments.</p>

<p>The default value of <code>useContainerLimits</code> if not specified is <code>true</code> so <code>-XX:+UseContainerSupport</code> will always be added
to the JVM arguments unless <code>useContainerLimits</code> is explicitly set to <code>false</code>. It is recommended that this value be left
unspecified as the default <code>true</code> unless other arguments are being passed to the JVM to limit its resource usage.</p>


<h3 id="_setting_container_resource_limits_for_the_implicit_role">Setting Container Resource Limits for the Implicit Role</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with a single implicit role the <code>useContainerLimits</code> is set in the <code>spec.jvm</code>
section of the configuration. For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    useContainerLimits: true  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>-XX:+UseContainerSupport</code> JVM option will be passed as arguments to the JVM for the implicit storage role.</li>
</ul>
</div>

<h3 id="_setting_container_resource_limits_for_explicit_roles">Setting Container Resource Limits for Explicit Roles</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles the <code>useContainerLimits</code> are set in the <code>jvm</code>
section of the configuration for each <code>role</code> in the <code>roles</code> list. For example:</p>

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
        useContainerLimits: true  <span class="conum" data-value="1" />
    - role: proxy
      jvm:
        useContainerLimits: false  <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>-XX:+UseContainerSupport</code> JVM option will be passed as arguments to the JVM for the explicit <code>data</code> role.</li>
<li data-value="2">The <code>-XX:+UseContainerSupport</code> JVM option will not be passed as arguments to the JVM for the explicit <code>proxy</code> role.</li>
</ul>
</div>

<h3 id="_setting_container_resource_limits_for_explicit_roles_with_a_default">Setting Container Resource Limits for Explicit Roles with a Default</h3>
<div class="section">
<p>When creating a <code>CoherenceCluster</code> with one or more explicit roles a default <code>useContainerLimits</code> value can be set in
the <code>CoherenceCluster</code> <code>spec</code> section that will apply to all of the roles in the <code>roles</code> list unless explicitly
overridden for a role. For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  jvm:
    useContainerLimits: true        <span class="conum" data-value="1" />
  roles:
    - role: data                    <span class="conum" data-value="2" />
      jvm:
        useContainerLimits: false
    - role: proxy                   <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The default <code>useContainerLimits</code> is set to <code>true</code>.</li>
<li data-value="2">The <code>data</code> role overrides the default <code>useContainerLimits</code> and sets it to <code>false</code>.</li>
<li data-value="3">The <code>proxy</code> role does not specify any <code>useContainerLimits</code> value so will use the default of <code>true</code>.</li>
</ul>
<v-divider class="my-5"/>
</div>
</div>

<h2 id="flightRecorder">Flight Recorder</h2>
<div class="section">
<p>Flight Recorder is a useful tool to use when diagnosing issues with a Coherence application or as an aid to performance
and GC tuning. By default the JVMs in a <code>CoherenceCluster</code> are configured to produce a continual flight recording that
will be dumped to a file when the JVM exits.</p>

<p>The <code>/jvm</code> root directory used for <code>.jfr</code> files can be <router-link to="#diagnosticsVolume" @click.native="this.scrollFix('#diagnosticsVolume')">mounted to an external volume</router-link> to allow
easier access to these files.</p>

<v-divider class="my-5"/>
</div>

<h2 id="diagnosticsVolume">Diagnostic Volume</h2>
<div class="section">
<p>By default the Coherence JVMs are configured to write heap dumps, error logs and flight recordings to directories in the
container under the root <code>/jvm</code> directory. The <code>/jvm</code> directory is mapped to <code>volumeMount</code> named <code>jvm</code> which is in turn
mapped to a <code>volume</code> named <code>jvm</code>.</p>

<p>The default configuration for the <code>jvm</code> volume in the Coherence <code>Pods</code> is an empty directory.</p>

<markup
lang="yaml"

>volumeMounts:
  - name: jvm
    mountPath: /jvm
volumes:
  - name: jvm
    emptyDir: {}</markup>

<p>The default may be changed to map the <code>jvm</code> volume to any supported
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/storage/volumes/#types-of-volumes">Kubernetes <code>VolumeSource</code></a>.</p>

</div>

<h2 id="debug">JVM Debug Arguments</h2>
<div class="section">
<p>Sometimes attaching a debugger to a JVM is the best way to track down the cause of an issue. The <code>CoherenceCluster</code> CRD
has a number of fields that can be used to configure how the JVM can be started in debug mode.</p>

</div>
</doc-view>
