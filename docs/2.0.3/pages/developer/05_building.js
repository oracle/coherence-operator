<doc-view>

<h2 id="_how_to_build_the_coherence_operator">How to Build the Coherence Operator</h2>
<div class="section">
<p>The Operator SDK generates Go projects that use Go Modules and hence the Coherence Operator uses Go Modules too.
The Coherence Operator can be checked out from Git to any location, it does not have to be under your <code>$GOPATH</code>.
The first time that the project is built may require Go to fetch a number of dependencies and may take longer than
usual to complete.</p>

<p>The easiest way to build the whole project is using <code>make</code>.
To build the Coherence Operator, package the Helm charts and create the various Docker images run the following
command:</p>

<markup
lang="bash"

>make all</markup>

<p>The <code>all</code> make target will build the Go and Java parts of the Operator and create all of the images required.</p>

<div class="admonition note">
<p class="admonition-inline">There have been issues with Go not being able to resolve all of the module dependencies required to build the
Coherence Operator. This can be resolved by setting the <code>GOPROXY</code> environment variable <code>GOPROXY=https://proxy.golang.org</code></p>
</div>
</div>

<h2 id="_build_versions">Build Versions</h2>
<div class="section">
<p>By default the version number used to tag the Docker images and Helm charts is set in the <code>VERSION</code> property
in the <code>Makefile</code> and in the <code>pom.xml</code> files in the <code>java/</code> directory.</p>

<p>The <code>Makefile</code> also contains a <code>VERSION_SUFFIX</code> variable that is used to add a suffix to the build. By default
this suffix is <code>ci</code> so the default version of the build artifacts is <code>2.0.3-ci</code>. Change this suffix, for
example when building a release candidate or a full release.</p>

<p>For example, if building a release called <code>alpha2</code> the following command can be used:</p>

<markup
lang="bash"

>make build-all-images VERSION_SUFFIX=alpha2</markup>

<p>If building a full release without a suffix the following command can be used</p>

<markup
lang="bash"

>make build-all-images VERSION_SUFFIX=""</markup>


<h3 id="_testing">Testing</h3>
<div class="section">

<h4 id="_unit_tests">Unit Tests</h4>
<div class="section">
<p>The Coherence Operator contains tests that can be executed using <code>make</code>. The tests are plain Go tests and
also <a id="" title="" target="_blank" href="https://github.com/onsi/ginkgo">Ginkgo</a> test suites.</p>

<p>To execute the unit and functional tests that do not require a k8s cluster you can execute the following command:</p>

<markup
lang="bash"

>make test-all</markup>

<p>This will build and execute all of the Go and Java tests, you do not need to have run a <code>make build</code> first.</p>

</div>

<h4 id="_go_unit_tests">Go Unit Tests</h4>
<div class="section">
<p>To only tun the Go tests use:</p>

<markup
lang="bash"

>make test-operator</markup>

</div>

<h4 id="_java_unit_tests">Java Unit Tests</h4>
<div class="section">
<p>To only tun the Java tests use:</p>

<markup
lang="bash"

>make test-mvn</markup>

</div>

<h4 id="_end_to_end_tests">End-to-End Tests</h4>
<div class="section">
<p>End to end tests require the Operator to be running. There are three types of end-to-end tests, Helm tests, local
tests and remote tests.</p>

<ul class="ulist">
<li>
<p>Helm tests are tests that install the Coherence Operator Helm chart and then make assertions about the state fo the
resulting install. These tests do not test functionality of the Operator itself.
The Helm tests suite is run using make:</p>

</li>
</ul>
<div class="listing">
<pre>make helm-test</pre>
</div>

<ul class="ulist">
<li>
<p>Local tests, which is the majority ot the tests, can be executed with a locally running operator (i.e. the operator
does not need to be deployed in a container in k8s). This makes the tests faster to run and also makes it possible
to run the operator in a debugger while the test is executing
The local end-to-end test suite is run using make:</p>

</li>
</ul>
<div class="listing">
<pre>make e2e-local-test</pre>
</div>

<p>It is possible to run a sub-set of the tests or an individual test by using the <code>GO_TEST_FLAGS=&lt;regex&gt;</code> parameter.
For example, to just run the <code>TestMinimalCoherenceCluster</code> clustering test in the <code>test/e2e/local/clustering_test.go</code>
file:</p>

<markup
lang="bash"

>make e2e-local-test GO_TEST_FLAGS='-run=^TestMinimalCoherenceCluster$$'</markup>

<p>The reg-ex above matches exactly the <code>TestMinimalCoherenceCluster</code> test name because it uses the reg-ex start <code>^</code> and
end <code>$</code> characters.</p>

<p>For example, to run all of the clustering tests where the test name starts with <code>TestOneRole</code> we can use
the reg-ex <code>^TestOneRole.*'</code></p>

<markup
lang="bash"

>make e2e-local-test  GO_TEST_FLAGS='-run=^TestOneRole.*'</markup>

<p><strong>Note</strong> Any <code>$</code> signs in the reg-ex need to be escaped by using a double dollar sign <code>$$</code>.</p>

<p>The <code>GO_TEST_FLAGS</code> parameter can actually consist of any valid argument to be passed to the <code>go test</code> command. There is plenty of
documentation on <a id="" title="" target="_blank" href="https://tip.golang.org/cmd/go/#hdr-Test_packages">go test</a></p>

<ul class="ulist">
<li>
<p>Remote tests require the operator to actually be installed in a container in k8s. An example of this is the scaling
tests because the operator needs to be able to directly reach the Pods. Very few end-to-end tests fall into this categrory.
The local end-to-end test suite is run using make:</p>

</li>
</ul>
<div class="listing">
<pre>make e2e-test</pre>
</div>

<p>As with local tests the <code>GO_TEST_FLAGS</code> parameter can be used to execute a sub-set of tests or a single test.</p>

</div>
</div>
</div>
</doc-view>
