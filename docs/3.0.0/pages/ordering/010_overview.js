<doc-view>

<h2 id="_coherence_deployment_dependencies_and_start_order">Coherence Deployment Dependencies and Start Order</h2>
<div class="section">
<p>The default behaviour of the operator is to create the <code>StatefulSet</code> for a <code>Coherence</code> deployment immediately.
Sometimes this behaviour is not suitable if, for example, when application code running in one deployment depends on the
availability of another deployment.
Typically, this might be storage disabled members having functionality that relies on the storage members being ready first.
The <code>Coherence</code> CRD allows can be configured with a <code>startQuorum</code> that defines a deployment&#8217;s dependency on other
deployments in the cluster.</p>

<div class="admonition note">
<p class="admonition-inline">The <code>startQuorum</code> only applies when a cluster is initially being started by the operator, it does not apply in other
functions such as upgrades, scaling, shut down etc.</p>
</div>
<p>An individual deployment can depend on one or more other deployment. The dependency can be such that the deployment will
not be created until all of the <code>Pods</code> of the dependent deployment are ready, or it can be configured so that just a
single <code>Pod</code> of the dependent deployment must be ready.</p>

<p>For example:
In the yaml snippet below there are two <code>Coherence</code> deployments, <code>data</code> and <code>proxy</code></p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: data
spec:
  replicas: 3           <span class="conum" data-value="1" />
---
apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: proxy
spec:
  startQuorum:          <span class="conum" data-value="2" />
    - deployment: data
      podCount: 1</markup>

<ul class="colist">
<li data-value="1">The <code>data</code> deployment does not specify a <code>startQuorum</code> so this role will be created immediately by the operator.</li>
<li data-value="2">The <code>proxy</code> deployment has a start quorum that means that the <code>proxy</code> deployment depends on the <code>data</code> deployment.
The <code>podCount</code> field has been set to <code>1</code> meaning the <code>proxy</code> deployment&#8217;s <code>StatefulSet</code> will not be created until at
least <code>1</code> of the <code>data</code> deployment&#8217;s <code>Pods</code> is in the <code>Ready</code> state.</li>
</ul>
<p>Omitting the <code>podCount</code> from the quorum means that the role will not start until all the configured replicas of the
dependent deployment are ready; for example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: data
spec:
  replicas: 3
---
apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: proxy
spec:
  startQuorum:          <span class="conum" data-value="1" />
    - deployment: data</markup>

<ul class="colist">
<li data-value="1">The <code>proxy</code> deployment&#8217;s <code>startQuorum</code> just specifies a dependency on the <code>data</code> deployment with no <code>podCount</code> so
all <code>3</code> of the <code>data</code> deployment&#8217;s <code>Pods</code> must be <code>Ready</code> before the <code>proxy</code> deployment&#8217;s <code>StatefulSet</code> is created by
the operator.</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">Setting a <code>podCount</code> less than or equal to zero is the same as not specifying a count.</p>
</div>

<h3 id="_multiple_dependencies">Multiple Dependencies</h3>
<div class="section">
<p>The <code>startQuorum</code> can specify a dependency on more than on deployment; for example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: data      <span class="conum" data-value="1" />
spec:
  replicas: 5
---
apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: proxy        <span class="conum" data-value="1" />
spec:
  replicas: 3
---
apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: web
spec:
  startQuorum:          <span class="conum" data-value="2" />
    - deployment: data
    - deployment: proxy
      podCount: 1</markup>

<ul class="colist">
<li data-value="1">The <code>data</code> and <code>proxy</code> deployments do not specify a <code>startQuorum</code>, so the <code>StatefulSets</code> for these deployments will
be created immediately by the operator.</li>
<li data-value="2">The <code>web</code> deployment has a <code>startQuorum</code> the defines a dependency on both the <code>data</code> deployment and the <code>proxy</code>
deployment. The <code>proxy</code> dependency also specifies a <code>podCount</code> of <code>1</code>.
This means that the operator wil not create the <code>web</code> role&#8217;s <code>StatefulSet</code> until all <code>5</code> replicas of the <code>data</code>
deployment are <code>Ready</code> and at least <code>1</code> of the <code>proxy</code> deployment&#8217;s <code>Pods</code> is <code>Ready</code>.</li>
</ul>
</div>

<h3 id="_chained_dependencies">Chained Dependencies</h3>
<div class="section">
<p>It is also possible to chain dependencies, for example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: data            <span class="conum" data-value="1" />
spec:
  replicas: 5
---
apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: proxy
spec:
  replicas: 3
  startQuorum:          <span class="conum" data-value="2" />
    - deployment: data
---
apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: web
spec:
  startQuorum:          <span class="conum" data-value="3" />
    - deployment: proxy
      podCount: 1</markup>

<ul class="colist">
<li data-value="1">The <code>data</code> deployment does not specify a <code>startQuorum</code> so this deployment&#8217;s <code>StatefulSet</code> will be created immediately
by the operator.</li>
<li data-value="2">The <code>proxy</code> deployment defines a dependency on the <code>data</code> deployment without a <code>podCount</code> so all five <code>Pods</code> of the
<code>data</code> role must be in a <code>Ready</code> state before the operator will create the <code>proxy</code> deployment&#8217;s <code>StatefulSet</code>.</li>
<li data-value="3">The <code>web</code> deployment depends on the <code>proxy</code> deployment with a <code>podCount</code> of one, so the operator will not create the
<code>web</code> deployment&#8217;s <code>StatefulSet</code> until at least one <code>proxy</code> deployment <code>Pod</code> is in a <code>Ready</code> state.</li>
</ul>
<div class="admonition warning">
<p class="admonition-inline">The operator does not validate that a <code>startQuorum</code> makes sense. It is possible to declare a quorum with circular
dependencies, in which case the roles will never start. It would also be possible to create a quorum with a <code>podCount</code> greater
than the <code>replicas</code> value of the dependent deployment, in which case the quorum would never be met, and the role would not start.</p>
</div>
</div>
</div>
</doc-view>
