<doc-view>

<h2 id="_readiness_liveness_probes">Readiness &amp; Liveness Probes</h2>
<div class="section">
<p>The Coherence Operator injects a Readiness/Liveness endpoint into the Coherence container that is used as the default
readiness and liveness check for the <code>Pods</code> deployed by the operator.
This endpoint is suitable for most use-cases, but it is possible to configure a different readiness and liveness probe,
or just change the timings of the probes if required.</p>

<p>The readiness/liveness probe used by the Operator in the Coherence Pods checks a number of things to determine whether the Pods is ready, one of these is whether the JVM is a cluster member.
If your application uses a custom main class and is not properly bootstrapping Coherence then the Pod will not be ready until your application code actually touches a Coherence resource causing Coherence to start and join the cluster.</p>

<p>When running in clusters with the Operator using custom main classes it is advisable to properly bootstrap Coherence
from within your <code>main</code> method. This can be done using the new Coherence bootstrap API available from CE release 20.12
or by calling <code>com.tangosol.net.DefaultCacheServer.startServerDaemon().waitForServiceStart();</code></p>


<h3 id="_coherence_readiness">Coherence Readiness</h3>
<div class="section">
<p>The default endpoint used by the Operator for readiness checks that the <code>Pod</code> is a member of the Coherence cluster and
that none of the partitioned cache services have a StatusHA value of <code>endangered</code>.
If the <code>Pod</code> is the only cluster member at the time of the ready check the StatusHA check will be skipped.
If a partitioned service has a backup count of zero the StatusHA check will be skipped for that service.</p>

<p>There are scenarios where the StatusHA check can fail but should be ignored because the application does not care
about data loss for caches on that particular cache service. Normally in this case the backup count for the cache
service would be zero, and the service would automatically be skipped in the StatusHA test.</p>

<p>The ready check used by the Operator can be configured to skip the StatusHA test for certain services.
In the <code>Coherence</code> CRD the <code>coherence.allowEndangeredForStatusHA</code> is a list of string values that can be
set to the names of partitioned cache services that should not be included in the StatusHA check.
For a service to be skipped its name must exactly match one of the names in the <code>allowEndangeredForStatusHA</code> list.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test-cluster
spec:
  coherence:
    allowEndangeredForStatusHA:   <span class="conum" data-value="1" />
      - TempService</markup>

<ul class="colist">
<li data-value="1">The <code>allowEndangeredForStatusHA</code> field is a list of string values. In this case the <code>TempService</code> will not
be checked for StatusHA in the ready check.</li>
</ul>
</div>

<h3 id="_configure_readiness">Configure Readiness</h3>
<div class="section">
<p>The <code>Coherence</code> CRD <code>spec.readinessProbe</code> field is identical to configuring a readiness probe for a <code>Pod</code>
in Kubernetes; see <a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/">Configure Liveness &amp; Readiness</a></p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test-cluster
spec:
  readinessProbe:
    httpGet:
      port: 8080
      path: "/ready"
    timeoutSeconds: 60
    initialDelaySeconds: 300
    periodSeconds: 120
    failureThreshold: 10
    successThreshold: 1</markup>

<p>The example above configures a http probe for readiness and sets different timings for the probe.
The <code>Coherence</code> CRD supports the other types of readiness probe too, <code>exec</code> and <code>tcpSocket</code>.</p>

</div>

<h3 id="_configure_liveness">Configure Liveness</h3>
<div class="section">
<p>The <code>Coherence</code> CRD <code>spec.livenessProbe</code> field is identical to configuring a liveness probe for a <code>Pod</code>
in Kubernetes; see <a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/">Configure Liveness &amp; Readiness</a></p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test-cluster
spec:
  livenessProbe:
    httpGet:
      port: 8080
      path: "/live"
    timeoutSeconds: 60
    initialDelaySeconds: 300
    periodSeconds: 120
    failureThreshold: 10
    successThreshold: 1</markup>

<p>The example above configures a http probe for liveness and sets different timings for the probe.
The <code>Coherence</code> CRD supports the other types of readiness probe too, <code>exec</code> and <code>tcpSocket</code>.</p>

</div>
</div>
</doc-view>
