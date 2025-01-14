<doc-view>

<h2 id="_environment_variables">Environment Variables</h2>
<div class="section">
<p>Environment variables can be added to the Coherence container in the <code>Pods</code> managed by the Operator.
Additional variables should be added to the <code>env</code> list in the <code>Coherence</code> CRD spec.
The entries in the <code>env</code> list are Kubernetes
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#envvar-v1-core">EnvVar</a>
values, exactly the same as when adding environment variables to a container spec.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  env:
    - name: VAR_ONE            <span class="conum" data-value="1" />
      value: VALUE_ONE
    - name: VAR_TWO            <span class="conum" data-value="2" />
      valueFrom:
        secretKeyRef:
          name: test-secret
          key: secret-key</markup>

<ul class="colist">
<li data-value="1">The <code>VAR_ONE</code> environment variable is a simple variable with a value of <code>VALUE_ONE</code></li>
<li data-value="2">The <code>VAR_TWO</code> environment variable is variable that is loaded from a secret.</li>
</ul>

<h3 id="_environment_variables_from">Environment Variables From</h3>
<div class="section">
<p>It is also possible to specify environment variables from a <code>ConfigMap</code> or <code>Secret</code> as you would for
a Kubernetes container.</p>

<p>For example, if there was a <code>ConfigMap</code> named <code>special-config</code> that contained environment variable values,
it can be added to the <code>Coherence</code> spec as shown below.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  envFrom:
    - configMapRef:
      name: special-config</markup>

</div>
</div>
</doc-view>
