<doc-view>

<h2 id="_using_jshell">Using JShell</h2>
<div class="section">
<p>JShell is a Java utility that allows Java code to be executed in a console.
Whilst it is simple to exec into a Pod and run JShell, the Coherence Operator will run JShell
configured with the same class path and system properties as the running Coherence container.
This makes it simpler to invoke JShell commands knowing that everything required to
access the running Coherence JVM is present.</p>


<h3 id="_using_jshell_in_pods">Using JShell in Pods</h3>
<div class="section">
<p>The Operator installs a simple CLI named <code>runner</code> at the location <code>/coherence-operator/utils/runner</code>.
One of the commands the runner can execute is <code>jshell</code> which will start a JShell process.</p>

<div class="admonition caution">
<p class="admonition-textlabel">Caution</p>
<p ><p>JShell can be a useful debugging tool, but running JShell in a production cluster is not recommended.</p>

<p>The JShell JVM will join the cluster as a storage disabled member alongside the JVM running in the
Coherence container in the Pod.
The JShell session will have all the same configuration parameters as the Coherence container.</p>

<p>For this reason, great care must be taken with the commands that are executed so that the cluster does not become unstable.</p>
</p>
</div>
</div>

<h3 id="_start_a_jshell_session">Start a JShell Session</h3>
<div class="section">
<p>The <code>kubectl exec</code> command can be used to create an interactive session in a Pod using the Coherence Operator runner
to start a JShell session.</p>

<p><strong>Example</strong></p>

<p>The yaml below will create a simple three member cluster.</p>

<markup

title="minimal.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
    replicas: 3</markup>

<p>The cluster name is <code>storage</code> and there will be three Pods created, <code>storage-0</code>, <code>storage-1</code> and <code>storage-2</code>.</p>

<p>A Query Plus session can be run by exec&#8217;ing into one of the Pods to execute the runner with the argument <code>jshell</code>.</p>

<markup
lang="bash"

>kubectl exec -it storage-0 -c coherence -- /coherence-operator/runner jshell</markup>

<p>After executing the above command, the <code>jshell&gt;</code> prompt will be displayed ready to accept input.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>The <code>kubectl exec</code> command must include the <code>-it</code> options so that <code>kubectl</code> creates an interactive terminal session.</p>
</p>
</div>
</div>

<h3 id="_starting_coherence">Starting Coherence</h3>
<div class="section">
<p>The JShell session only starts the JShell REPL and Coherence is not started in the JShell process.
As the JShell process has all the same configuration except it is configured to be storage disabled.
As the Coherence container in the Pod any of the normal ways to bootstrap Coherence can be used.
Any configuration changes, for example setting system properties, can be done before Coherence is started.</p>

<p>For example:</p>

<markup
lang="java"

>jshell&gt; import com.tangosol.net.*;

jshell&gt; Coherence c = Coherence.clusterMember().start().join();

jshell&gt; Session s = c.getSession();
s ==&gt; com.tangosol.internal.net.ConfigurableCacheFactorySession@3d0f8e03

jshell&gt; NamedCache&lt;String, String&gt; cache = s.getCache("test");
cache ==&gt; com.tangosol.internal.net.SessionNamedCache@91213130

jshell&gt; cache.size();
$5 ==&gt; 0

jshell&gt;</markup>

</div>
</div>
</doc-view>
