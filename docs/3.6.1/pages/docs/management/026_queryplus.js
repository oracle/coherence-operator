<doc-view>

<h2 id="_the_coherence_query_plus">The Coherence Query Plus</h2>
<div class="section">
<p>The Coherence Query Plus utility is a console application that allows simple SQL like queries
to be made against caches, see the
<a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/14.1.2/develop-applications/using-coherence-query-language.html">Using Coherence Query Language</a>
section of the Coherence documentation.</p>


<h3 id="_using_query_plus_in_pods">Using Query Plus in Pods</h3>
<div class="section">
<p>Most official Coherence images are distroless images, so they do not have a shell that can be used to
create a command line session and execute commands.
The Operator works around this to support a few selected commands by injecting its <code>runner</code> utility.
The Operator installs the <code>runner</code> at the location <code>/coherence-operator/utils/runner</code>.</p>

<p>The <code>runner</code> utility is a simple CLI that executes commands, one of those is <code>queryplus</code> which will
start a Java process running Query Plus.</p>

<div class="admonition caution">
<p class="admonition-textlabel">Caution</p>
<p ><p>The Query Plus JVM will join the cluster as a storage disabled member alongside the JVM running in the
Coherence container in the Pod.
The Query Plus session will have all the same configuration parameters as the Coherence container.</p>

<p>For this reason, great care must be taken with the commands that are executed so that the cluster does not become unstable.</p>
</p>
</div>
</div>

<h3 id="_start_a_query_plus_session">Start a Query Plus Session</h3>
<div class="section">
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

<p>A Query Plus session can be run by exec&#8217;ing into one of the Pods to execute the runner with the argument <code>queryplus</code>.</p>

<markup
lang="bash"

>kubectl exec -it storage-0 -c coherence -- /coherence-operator/runner queryplus</markup>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>The <code>kubectl exec</code> command must include the <code>-it</code> options so that <code>kubectl</code> creates an interactive terminal session.</p>
</p>
</div>
<p>After executing the above command, the <code>CohQL&gt;</code> prompt will be displayed ready to accept input.
Using the Query Plus utility is documented in the
<a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/14.1.2/develop-applications/using-coherence-query-language.html#GUID-1CBE48A8-1009-4656-868D-663AA85CB021">Using the CohQL Command-Line Tool</a>
section of the Coherence documentation</p>

</div>

<h3 id="_run_query_plus_with_command_line_arguments">Run Query Plus With Command Line Arguments</h3>
<div class="section">
<p>Instead of running an interactive Query Plus session, arguments can be passed into Query Plus as part of the exec command.
Query Plus will execute the commands and exit.</p>

<p>The command line for this is slightly complicated because there are two CLI programs involved in the full command line,
first <code>kubectl</code> and second the Operator&#8217;s runner.
In each case the <code>--</code> command line separator needs to be used so that each CLI knows the everything after a <code>--</code>
is to be passed to the next process.</p>

<p>For example a simple string key and value could be inserted into a cache named "test" with the following
CohQL statement <code>insert into "test" key "one" value "value-one"</code>.
This statement can be executed in a Pod with the following command</p>

<markup
lang="bash"

>kubectl exec storage-0 -c coherence -- /coherence-operator/runner queryplus -- -c -l 'insert into test key "one" value "value-one"'</markup>

<p>In the above example the first <code>--</code> tels <code>kubectl</code> that all the remaining arguments are to be passed
as arguments to the exec session. The second <code>--</code> tells the Operator runner that all the remaining arguments
are to be passed to Query Plus.</p>

<p>After running the above command the cache <code>test</code> will contain an entry with the key <code>"one"</code> and value <code>"value-one"</code>.
If the statement <code>select * from test</code> is executed the value in the cache will be displayed.</p>

<markup
lang="bash"

>kubectl exec storage-0 -c coherence -- /coherence-operator/runner queryplus -- -c -l 'select * from test'</markup>

<p>The last few lines of the console output will display the results of executing the statement:</p>

<markup


>Results
"value-one"</markup>

</div>
</div>
</doc-view>
