<doc-view>

<h2 id="_debugging_the_coherence_operator">Debugging the Coherence Operator</h2>
<div class="section">
<p>Assuming that you have an IDE capable of debugging Go and have
<a id="" title="" target="_blank" href="https://github.com/go-delve/delve/tree/master/Documentation/installation">delve</a> installed you can debug the operator.
When debugging an instance of the operator is run locally so functionality that will only work when the operator is
deployed into k8s cannot be properly debugged.</p>

<p>To start an instance of the operator that can be debugged use the make target <code>run-debug</code>, for example:</p>

<markup
lang="bash"

>make run-debug</markup>

<p>This will start the operator and listen for a debugger to connect on the default delve port <code>2345</code>.
The operator will connect to whichever k8s cluster the current environment is configured to point to.</p>


<h3 id="_stopping_the_debug_session">Stopping the Debug Session</h3>
<div class="section">
<p>To stop the local operator just use CTRL-Z or CTRL-C. Sometimes processes can be left around even after exiting in
this way. To make sure all of the processes are dead you can run the kill script:</p>

<markup
lang="bash"

>make debug-stop</markup>

</div>

<h3 id="_debugging_tests">Debugging Tests</h3>
<div class="section">
<p>To debug the operator while running a particular tests first start the debugger as described above.
Then use the debug make test target to execute the test.</p>

<p>For example to debug the <code>TestMinimalCoherenceCluster</code> test first start the debug session:</p>

<markup
lang="bash"

>make run-debug</markup>

<p>Then execute the test with the <code>debug-e2e-local-test</code> make target:</p>

<markup
lang="bash"

>make debug-e2e-local-test GO_TEST_FLAGS='-run=^TestMinimalCoherenceCluster$$'</markup>

</div>
</div>
</doc-view>
