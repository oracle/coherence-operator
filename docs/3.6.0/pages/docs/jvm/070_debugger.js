<doc-view>

<h2 id="_debugger_configuration">Debugger Configuration</h2>
<div class="section">
<p>Occasionally it is useful to be able to connect a debugger to a JVM, and the <code>Coherence</code> CRD spec has fields to
configure the Coherence container&#8217;s JVM to work with a debugger. The fields in the CRD will ultimately result in
arguments being passed to the JVM and could have been added as plain JVM arguments, but having specific fields in the
CRD makes it simpler to configure and the intention more obvious.</p>

<p>The fields to control debug settings of the JVM are in the <code>jvm.debug</code> section of the CRD spec.</p>


<h3 id="_listening_for_a_debugger_connection">Listening for a Debugger Connection</h3>
<div class="section">
<p>One scenario for debugging is for the Coherence JVM to open a port and listen for a debugger connection request.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    debug:
      enabled: true   <span class="conum" data-value="1" />
      port: 5005      <span class="conum" data-value="2" />
      suspend: false  <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The <code>jvm.debug.enabled</code> flag is set to <code>true</code> to enable debug mode.</li>
<li data-value="2">The <code>jvm.debug.port</code> field specifies the port the JVM will listen on for a debugger connection.</li>
<li data-value="3">The <code>jvm.debug.suspend</code> flag is set to <code>false</code> so that the JVM will start without waiting for a debugger to connect.</li>
</ul>
<p>The example above results in the following arguments being passed to the JVM:</p>

<markup


>-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5005</markup>

<ul class="ulist">
<li>
<p>The <code>address=*:5005</code> value comes from the <code>jvm.debug.port</code> field</p>

</li>
<li>
<p>The <code>suspend=n</code> value comes from the <code>jvm.debug.suspend</code> field</p>

</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">If the <code>jvm.debug.port</code> is not specified the default value used by the Operator will be <code>5005</code>.</p>
</div>
</div>

<h3 id="_attaching_to_a_debugger_connection">Attaching to a Debugger Connection</h3>
<div class="section">
<p>Another scenario for debugging is for the Coherence JVM to connect out to a listening debugger.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  jvm:
    debug:
      enabled: true               <span class="conum" data-value="1" />
      attach:  "10.10.100.2:5000" <span class="conum" data-value="2" />
      suspend: false              <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The <code>jvm.debug.enabled</code> flag is set to <code>true</code> to enable debug mode.</li>
<li data-value="2">The <code>jvm.debug.attach</code> field specifies the address of the debugger that the JVM will connect to.</li>
<li data-value="3">The <code>jvm.debug.suspend</code> flag is set to <code>false</code> so that the JVM will start without waiting for a debugger to connect.</li>
</ul>
<p>The example above results in the following arguments being passed to the JVM:</p>

<markup


>-agentlib:jdwp=transport=dt_socket,server=n,address=10.10.100.2:5000,suspend=n,timeout=10000</markup>

</div>
</div>
</doc-view>
