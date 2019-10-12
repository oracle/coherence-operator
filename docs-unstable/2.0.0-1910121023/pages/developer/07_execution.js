<doc-view>

<h2 id="_running_coherence_operator_development">Running Coherence Operator Development</h2>
<div class="section">
<p>There are two ways to run the Coherence Operator, either deployed into a k8s cluster or by using the Operator SDK
to run it locally on your dev machine (assuming your dev machine has access to a k8s cluster such as Docker Desktop
on MacOS).</p>


<h3 id="_namespaces">Namespaces</h3>
<div class="section">
<p><strong>NOTE:</strong> The Coherence Operator by default runs in and monitors a <strong>single</strong> namespace.
This is different behaviour to v1.0 of the Coherence Operator.
For more details see the Operator SDK document on
<a id="" title="" target="_blank" href="https://github.com/operator-framework/operator-sdk/blob/v0.9.0/doc/operator-scope.md">Operator Scope</a>.</p>

</div>

<h3 id="_running_locally">Running Locally</h3>
<div class="section">
<p>During development running the Coherence Operator locally is by far the simplest option as it is faster and
it also allows remote debugging if you are using a suitable IDE.</p>

<p>To run a local copy of the operator that will connect to whatever you local kubernetes config is pointing to:</p>

<markup
lang="bash"

>make run</markup>


<h4 id="_stopping_the_local_operator">Stopping the Local Operator</h4>
<div class="section">
<p>To stop the local operator just use CTRL-Z or CTRL-C. Sometimes processes can be left around even after exiting in
this way. To make sure all of the processes are dead you can run the kill script:</p>

<markup
lang="bash"

>./hack/kill-local.sh</markup>

</div>
</div>

<h3 id="_clean_up">Clean-up</h3>
<div class="section">
<p>After running the operator the CRDs can be removed from the k8s cluster by running the make target:</p>

<markup
lang="bash"

>make uninstall-crds</markup>

</div>

<h3 id="_deploying_to_kubernetes">Deploying to Kubernetes</h3>
<div class="section">
<p>The simplest and most reliable way to deploy the operator to K8s is to use the Helm chart.
After building the operator the chart is created in the <code>build/_output/helm-charts/coherence-operator</code> directory.
Using the Helm chart will ensure that all of the required RBAC rules are created when deploying to an environment
where RBAC is enabled.
The chart can be installed in the usual way with Helm</p>

<markup
lang="bash"

>helm install --name operator \
  --namespace operator-test \
  build/_output/helm-charts/coherence-operator</markup>

</div>
</div>
</doc-view>
