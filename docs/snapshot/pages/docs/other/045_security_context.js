<doc-view>

<h2 id="_pod_container_securitycontext">Pod &amp; Container SecurityContext</h2>
<div class="section">
<p>Kubernetes allows you to configure a <a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/configure-pod-container/security-context/">Security Context</a> for both Pods and Containers. The Coherence CRD exposes both of these to allow you to set the security context configuration for the Coherence Pods and for the Coherence containers withing the Pods.</p>

<p>For more details see the Kubernetes <a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/configure-pod-container/security-context/">Security Context</a> documentation.</p>


<h3 id="_setting_the_pod_security_context">Setting the Pod Security Context</h3>
<div class="section">
<p>To specify security settings for a Pod, include the <code>securityContext</code> field in the Coherence resource specification.
The securityContext field is a <a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#podsecuritycontext-v1-core">PodSecurityContext</a> object. The security settings that you specify for a Pod apply to all Containers in the Pod. Here is a configuration file for a Pod that has a securityContext:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  securityContext:
    runAsUser: 1000
    runAsGroup: 3000
    fsGroup: 2000</markup>

</div>

<h3 id="_setting_the_coherence_container_security_context">Setting the Coherence Container Security Context</h3>
<div class="section">
<p>To specify security settings for the Coherence container within the Pods, include the <code>containerSecurityContext</code> field in the Container manifest. The <code>containerSecurityContext</code> field is a <a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#securitycontext-v1-core">SecurityContext</a> object.
Security settings that you specify in the <code>containerSecurityContext</code> field apply only to the individual Coherence container and the Operator init-container, and they override settings made at the Pod level in the <code>securityContext</code> field when there is overlap. Container settings do not affect the Pod&#8217;s Volumes.</p>

<p>Here is the configuration file for a Coherence resource that has both the Pod and the container security context:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  securityContext:
    runAsUser: 1000
    runAsGroup: 3000
    fsGroup: 2000
  containerSecurityContext:
      runAsUser: 2000
      allowPrivilegeEscalation: false
      capabilities:
        add: ["NET_ADMIN", "SYS_TIME"]</markup>

</div>
</div>
</doc-view>
