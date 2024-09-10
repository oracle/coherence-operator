<doc-view>

<h2 id="_scale_coherence_deployments">Scale Coherence Deployments</h2>
<div class="section">
<p>The Coherence Operator provides the ability to safely scale up and down a <code>Coherence</code> deployment.
A <code>Coherence</code> deployment is backed by a <code>StatefulSet</code>, which can easily be scaled using existing Kubernetes features.
The problem with directly scaling down the <code>StatefulSet</code> is that Kubernetes will immediately kill the required number
of <code>Pods</code>. This is obviously very bad for Coherence as killing multiple storage enabled members would almost certainly
cause data loss.</p>

<p>The Coherence Operator supports scaling by applying the scaling update directly to <code>Coherence</code> deployment rather than
to the underlying <code>StatefulSet</code>. There are two methods to scale a <code>Coherence</code> deployment:</p>

<ul class="ulist">
<li>
<p>Update the <code>replicas</code> field in the <code>Coherence</code> CRD spec.</p>

</li>
<li>
<p>Use the <code>kubectl scale</code> command</p>

</li>
</ul>
<p>When either of these methods is used the Operator will detect that a change to the size of the deployment is required
and ensure that the change will be applied safely. The logical steps the Operator will perform are:</p>

<ol style="margin-left: 15px;">
<li>
Detect desired replicas is different to current replicas

</li>
<li>
Check the cluster is StatusHA - i.e. no cache services are endangered. If any service is not StatusHA requeue the
scale request  (go back to step one).

</li>
<li>
If scaling up, add the required number of members.

</li>
<li>
If scaling down, scale down by one member and requeue the request (go back to step one).

</li>
</ol>
<p>What these steps ensure is that the deployment will not be resized unless the cluster is in a safe state.
When scaling down only a single member will be removed at a time, ensuring that the cluster is in a safe state before
removing the next member.</p>

<div class="admonition note">
<p class="admonition-inline">The Operator will only apply safe scaling functionality to deployments that are storage enabled.
If a deployment is storage disabled then it can be scaled up or down by the required number of members
in one step as there is no fear of data loss in a storage disabled member.</p>
</div>
</div>

<h2 id="_controlling_safe_scaling">Controlling Safe Scaling</h2>
<div class="section">
<p>The <code>Coherence</code> CRD has a number of fields that control the behaviour of scaling.</p>


<h3 id="_scaling_policy">Scaling Policy</h3>
<div class="section">
<p>The <code>Coherence</code> CRD spec has a field <code>scaling.policy</code> that can be used to override the default scaling
behaviour. The scaling policy has three possible values:</p>


<div class="table__overflow elevation-1  ">
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
<td class=""><code>ParallelUpSafeDown</code></td>
<td class="">This is the default scaling policy.
With this policy when scaling up <code>Pods</code> are added in parallel (the same as using the <code>Parallel</code> <code>podManagementPolicy</code>
in a <a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetspec-v1-apps">StatefulSet</a>) and
when scaling down <code>Pods</code> are removed one at a time (the same as the <code>OrderedReady</code> <code>podManagementPolicy</code> for a
StatefulSet). When scaling down a check is done to ensure that the members of the cluster have a safe StatusHA value
before a <code>Pod</code> is removed (i.e. none of the Coherence cache services have an endangered status).
This policy offers faster scaling up and start-up because pods are added in parallel as data should not be lost when
adding members, but offers safe, albeit slower,  scaling down as <code>Pods</code> are removed one by one.</td>
</tr>
<tr>
<td class=""><code>Parallel</code></td>
<td class="">With this policy when scaling up <code>Pods</code> are added in parallel (the same as using the <code>Parallel</code> <code>podManagementPolicy</code>
in a <a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetspec-v1-apps">StatefulSet</a>).
With this policy no StatusHA check is performed either when scaling up or when scaling down.
This policy allows faster start and scaling times but at the cost of no data safety; it is ideal for deployments that are
storage disabled.</td>
</tr>
<tr>
<td class=""><code>Safe</code></td>
<td class="">With this policy when scaling up and down <code>Pods</code> are removed one at a time (the same as the <code>OrderedReady</code>
<code>podManagementPolicy</code> for a StatefulSet). When scaling down a check is done to ensure that the members of the deployment
have a safe StatusHA value before a <code>Pod</code> is removed (i.e. none of the Coherence cache services have an endangered status).
This policy is slower to start, scale up and scale down.</td>
</tr>
</tbody>
</table>
</div>
<p>Both the <code>ParallelUpSafeDown</code> and <code>Safe</code> policies will ensure no data loss when scaling a deployment.</p>

<p>The policy can be set as shown below:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  scaling:
    policy: Safe <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">This deployment will scale both up and down with StatusHA checks.</li>
</ul>
</div>

<h3 id="_scaling_statusha_probe">Scaling StatusHA Probe</h3>
<div class="section">
<p>The StatusHA check performed by the Operator uses a http endpoint that the Operator runs on a well-known port in the
Coherence JVM. This endpoint performs a simple check to verify that none of the partitioned cache services known
about by Coherence have an endangered status. If an application has a different concept of what "safe" means it can
implement a different method to check the status during scaling.</p>

<p>The operator supports different types of safety check probes, these are exactly the same as those supported by
Kubernetes for readiness and liveness probes. The <code>scaling.probe</code> section of the <code>Coherence</code> CRD allows different
types of probe to be configured.</p>


<h4 id="_using_a_http_get_probe">Using a HTTP Get Probe</h4>
<div class="section">
<p>An HTTP get probe works the same way as a
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-a-liveness-http-request">Kubernetes liveness http request</a></p>

<p>The probe can be configured as follows</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  scaling:
    probe:
      httpGet:             <span class="conum" data-value="1" />
        port: 8080
        path: /statusha</markup>

<ul class="colist">
<li data-value="1">This deployment will check the status of the services by performing a http GET on <code><a id="" title="" target="_blank" href="http://&lt;pod-ip&gt;:8080/statusha">http://&lt;pod-ip&gt;:8080/statusha</a></code>.
If the response is <code>200</code> the check will pass, any other response the check is assumed to be false.</li>
</ul>
</div>

<h4 id="_using_a_tcp_probe">Using a TCP Probe</h4>
<div class="section">
<p>A TCP probe works the same way as a
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-a-tcp-liveness-probe">Kubernetes TCP liveness probe</a></p>

<p>The probe can be configured as follows</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  scaling:
    probe:
      tcpSocket:    <span class="conum" data-value="1" />
        port: 7000</markup>

<ul class="colist">
<li data-value="1">This deployment will check the status of the services by connecting to the socket on port <code>7000</code>.</li>
</ul>
</div>

<h4 id="_using_an_exec_command_probe">Using an Exec Command Probe</h4>
<div class="section">
<p>An exec probe works the same way as a
<a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-a-liveness-command">Kubernetes Exec liveness probe</a></p>

<p>The probe can be configured as follows</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  scaling:
    probe:
      exec:
        command:      <span class="conum" data-value="1" />
          - /bin/ah
          - safe.sh</markup>

<ul class="colist">
<li data-value="1">This deployment will check the status of the services by running the <code>sh safe.sh</code> command in the <code>Pod</code>.</li>
</ul>
</div>
</div>
</div>
</doc-view>
