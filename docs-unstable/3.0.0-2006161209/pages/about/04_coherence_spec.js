<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Coherence Operator API Docs</dt>
<dd slot="desc"><p>A reference guide to the Coherence Operator CRD types.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_coherence_operator_api_docs">Coherence Operator API Docs</h2>
<div class="section">
<p>This is a reference for the Coherence Operator API types.
These are all the types and fields that are used in the Coherence CRD.</p>

<div class="admonition tip">
<p class="admonition-inline">This document was generated from comments in the Go structs in the pkg/api/ directory.</p>
</div>

<h3 id="_table_of_contents">Table of Contents</h3>
<div class="section">
<ul class="ulist">
<li>
<p><router-link to="#_coherenceresourcespec" @click.native="this.scrollFix('#_coherenceresourcespec')">CoherenceResourceSpec</router-link></p>

</li>
<li>
<p><router-link to="#_applicationspec" @click.native="this.scrollFix('#_applicationspec')">ApplicationSpec</router-link></p>

</li>
<li>
<p><router-link to="#_coherencespec" @click.native="this.scrollFix('#_coherencespec')">CoherenceSpec</router-link></p>

</li>
<li>
<p><router-link to="#_coherencetracingspec" @click.native="this.scrollFix('#_coherencetracingspec')">CoherenceTracingSpec</router-link></p>

</li>
<li>
<p><router-link to="#_coherencewkaspec" @click.native="this.scrollFix('#_coherencewkaspec')">CoherenceWKASpec</router-link></p>

</li>
<li>
<p><router-link to="#_configmapvolumespec" @click.native="this.scrollFix('#_configmapvolumespec')">ConfigMapVolumeSpec</router-link></p>

</li>
<li>
<p><router-link to="#_imagespec" @click.native="this.scrollFix('#_imagespec')">ImageSpec</router-link></p>

</li>
<li>
<p><router-link to="#_jvmspec" @click.native="this.scrollFix('#_jvmspec')">JVMSpec</router-link></p>

</li>
<li>
<p><router-link to="#_jvmdebugspec" @click.native="this.scrollFix('#_jvmdebugspec')">JvmDebugSpec</router-link></p>

</li>
<li>
<p><router-link to="#_jvmgarbagecollectorspec" @click.native="this.scrollFix('#_jvmgarbagecollectorspec')">JvmGarbageCollectorSpec</router-link></p>

</li>
<li>
<p><router-link to="#_jvmjmxmpspec" @click.native="this.scrollFix('#_jvmjmxmpspec')">JvmJmxmpSpec</router-link></p>

</li>
<li>
<p><router-link to="#_jvmmemoryspec" @click.native="this.scrollFix('#_jvmmemoryspec')">JvmMemorySpec</router-link></p>

</li>
<li>
<p><router-link to="#_jvmoutofmemoryspec" @click.native="this.scrollFix('#_jvmoutofmemoryspec')">JvmOutOfMemorySpec</router-link></p>

</li>
<li>
<p><router-link to="#_localobjectreference" @click.native="this.scrollFix('#_localobjectreference')">LocalObjectReference</router-link></p>

</li>
<li>
<p><router-link to="#_namedportspec" @click.native="this.scrollFix('#_namedportspec')">NamedPortSpec</router-link></p>

</li>
<li>
<p><router-link to="#_networkspec" @click.native="this.scrollFix('#_networkspec')">NetworkSpec</router-link></p>

</li>
<li>
<p><router-link to="#_persistencespec" @click.native="this.scrollFix('#_persistencespec')">PersistenceSpec</router-link></p>

</li>
<li>
<p><router-link to="#_persistentstoragespec" @click.native="this.scrollFix('#_persistentstoragespec')">PersistentStorageSpec</router-link></p>

</li>
<li>
<p><router-link to="#_poddnsconfig" @click.native="this.scrollFix('#_poddnsconfig')">PodDNSConfig</router-link></p>

</li>
<li>
<p><router-link to="#_portspecwithssl" @click.native="this.scrollFix('#_portspecwithssl')">PortSpecWithSSL</router-link></p>

</li>
<li>
<p><router-link to="#_probehandler" @click.native="this.scrollFix('#_probehandler')">ProbeHandler</router-link></p>

</li>
<li>
<p><router-link to="#_readinessprobespec" @click.native="this.scrollFix('#_readinessprobespec')">ReadinessProbeSpec</router-link></p>

</li>
<li>
<p><router-link to="#_resource" @click.native="this.scrollFix('#_resource')">Resource</router-link></p>

</li>
<li>
<p><router-link to="#_resources" @click.native="this.scrollFix('#_resources')">Resources</router-link></p>

</li>
<li>
<p><router-link to="#_sslspec" @click.native="this.scrollFix('#_sslspec')">SSLSpec</router-link></p>

</li>
<li>
<p><router-link to="#_scalingprobe" @click.native="this.scrollFix('#_scalingprobe')">ScalingProbe</router-link></p>

</li>
<li>
<p><router-link to="#_scalingspec" @click.native="this.scrollFix('#_scalingspec')">ScalingSpec</router-link></p>

</li>
<li>
<p><router-link to="#_secretvolumespec" @click.native="this.scrollFix('#_secretvolumespec')">SecretVolumeSpec</router-link></p>

</li>
<li>
<p><router-link to="#_servicemonitorspec" @click.native="this.scrollFix('#_servicemonitorspec')">ServiceMonitorSpec</router-link></p>

</li>
<li>
<p><router-link to="#_servicespec" @click.native="this.scrollFix('#_servicespec')">ServiceSpec</router-link></p>

</li>
<li>
<p><router-link to="#_startquorum" @click.native="this.scrollFix('#_startquorum')">StartQuorum</router-link></p>

</li>
<li>
<p><router-link to="#_startquorumstatus" @click.native="this.scrollFix('#_startquorumstatus')">StartQuorumStatus</router-link></p>

</li>
<li>
<p><router-link to="#_coherence" @click.native="this.scrollFix('#_coherence')">Coherence</router-link></p>

</li>
<li>
<p><router-link to="#_coherencelist" @click.native="this.scrollFix('#_coherencelist')">CoherenceList</router-link></p>

</li>
<li>
<p><router-link to="#_coherenceresourcestatus" @click.native="this.scrollFix('#_coherenceresourcestatus')">CoherenceResourceStatus</router-link></p>

</li>
</ul>
</div>

<h3 id="_coherenceresourcespec">CoherenceResourceSpec</h3>
<div class="section">
<p>CoherenceResourceSpec defines the specification of a Coherence resource. A Coherence resource is typically one or more Pods that perform the same functionality, for example storage members.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>image</code></td>
<td class="">The name of the image. More info: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/containers/images">https://kubernetes.io/docs/concepts/containers/images</a></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>imagePullPolicy</code></td>
<td class="">Image pull policy. One of Always, Never, IfNotPresent. More info: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/containers/images#updating-images">https://kubernetes.io/docs/concepts/containers/images#updating-images</a></td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#pullpolicy-v1-core">corev1.PullPolicy</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>imagePullSecrets</code></td>
<td class="">ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec. If specified, these secrets will be passed to individual puller implementations for them to use. For example, in the case of docker, only DockerConfig type secrets are honored. More info: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod">https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod</a></td>
<td class=""><code>[]<router-link to="#_localobjectreference" @click.native="this.scrollFix('#_localobjectreference')">LocalObjectReference</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>replicas</code></td>
<td class="">The desired number of cluster members of this deployment. This is a pointer to distinguish between explicit zero and not specified. If not specified a default value of 3 will be used. This field cannot be negative.</td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>cluster</code></td>
<td class="">The optional name of the Coherence cluster that this Coherence resource belongs to. If this value is set the Pods controlled by this Coherence resource will form a cluster with other Pods controlled by Coherence resources with the same cluster name. If not set the Coherence resource&#8217;s name will be used as the cluster name.</td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>role</code></td>
<td class="">The name of the role that this deployment represents in a Coherence cluster. This value will be used to set the Coherence role property for all members of this role</td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>coherence</code></td>
<td class="">The optional settings specific to Coherence functionality.</td>
<td class=""><code>&#42;<router-link to="#_coherencespec" @click.native="this.scrollFix('#_coherencespec')">CoherenceSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>application</code></td>
<td class="">The optional application specific settings.</td>
<td class=""><code>&#42;<router-link to="#_applicationspec" @click.native="this.scrollFix('#_applicationspec')">ApplicationSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>jvm</code></td>
<td class="">The JVM specific options</td>
<td class=""><code>&#42;<router-link to="#_jvmspec" @click.native="this.scrollFix('#_jvmspec')">JVMSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>ports</code></td>
<td class="">Ports specifies additional port mappings for the Pod and additional Services for those ports.</td>
<td class=""><code>[]<router-link to="#_namedportspec" @click.native="this.scrollFix('#_namedportspec')">NamedPortSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>scaling</code></td>
<td class="">The configuration to control safe scaling.</td>
<td class=""><code>&#42;<router-link to="#_scalingspec" @click.native="this.scrollFix('#_scalingspec')">ScalingSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>startQuorum</code></td>
<td class="">StartQuorum controls the start-up order of this Coherence resource in relation to other Coherence resources.</td>
<td class=""><code>[]<router-link to="#_startquorum" @click.native="this.scrollFix('#_startquorum')">StartQuorum</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>env</code></td>
<td class="">Env is additional environment variable mappings that will be passed to the Coherence container in the Pod. To specify extra variables add them as name value pairs the same as they would be added to a Pod containers spec.</td>
<td class=""><code>[]<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#envvar-v1-core">corev1.EnvVar</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>labels</code></td>
<td class="">The extra labels to add to the all of the Pods in this deployments. Labels here will add to or override those defined for the cluster. More info: <a id="" title="" target="_blank" href="http://kubernetes.io/docs/user-guide/labels">http://kubernetes.io/docs/user-guide/labels</a></td>
<td class=""><code>map[string]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>annotations</code></td>
<td class="">Annotations are free-form yaml that will be added to the store release as annotations Any annotations should be placed BELOW this annotations: key. For example if we wanted to include annotations for Prometheus it would look like this:<br>
<br>
annotations:<br>
  prometheus.io/scrape: \"true\" +<br>
  prometheus.io/port: \"2408\"<br></td>
<td class=""><code>map[string]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>initContainers</code></td>
<td class="">List of additional initialization containers to add to the deployment&#8217;s Pod. More info: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/workloads/pods/init-containers/">https://kubernetes.io/docs/concepts/workloads/pods/init-containers/</a></td>
<td class=""><code>[]<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#container-v1-core">corev1.Container</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>sideCars</code></td>
<td class="">List of additional side-car containers to add to the deployment&#8217;s Pod.</td>
<td class=""><code>[]<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#container-v1-core">corev1.Container</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>configMapVolumes</code></td>
<td class="">A list of ConfigMaps to add as volumes. Each entry in the list will be added as a ConfigMap Volume to the deployment&#8217;s Pods and as a VolumeMount to all of the containers and init-containers in the Pod.<br>
see: <router-link to="/other/050_configmap_volumes">Add ConfigMap Volumes</router-link></td>
<td class=""><code>[]<router-link to="#_configmapvolumespec" @click.native="this.scrollFix('#_configmapvolumespec')">ConfigMapVolumeSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>secretVolumes</code></td>
<td class="">A list of Secrets to add as volumes. Each entry in the list will be added as a Secret Volume to the deployment&#8217;s Pods and as a VolumeMount to all of the containers and init-containers in the Pod.<br>
see: <router-link to="#other/020_secret_volumes.adoc" @click.native="this.scrollFix('#other/020_secret_volumes.adoc')">Add Secret Volumes</router-link></td>
<td class=""><code>[]<router-link to="#_secretvolumespec" @click.native="this.scrollFix('#_secretvolumespec')">SecretVolumeSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>volumes</code></td>
<td class="">Volumes defines extra volume mappings that will be added to the Coherence Pod.<br>
  The content of this yaml should match the normal k8s volumes section of a Pod definition +<br>
  as described in <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/storage/volumes/">https://kubernetes.io/docs/concepts/storage/volumes/</a><br></td>
<td class=""><code>[]<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#volume-v1-core">corev1.Volume</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>volumeClaimTemplates</code></td>
<td class="">VolumeClaimTemplates defines extra PVC mappings that will be added to the Coherence Pod.<br>
  The content of this yaml should match the normal k8s volumeClaimTemplates section of a Pod definition +<br>
  as described in <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/storage/persistent-volumes/">https://kubernetes.io/docs/concepts/storage/persistent-volumes/</a><br></td>
<td class=""><code>[]<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#persistentvolumeclaim-v1-core">corev1.PersistentVolumeClaim</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>volumeMounts</code></td>
<td class="">VolumeMounts defines extra volume mounts to map to the additional volumes or PVCs declared above<br>
  in store.volumes and store.volumeClaimTemplates<br></td>
<td class=""><code>[]<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#volumemount-v1-core">corev1.VolumeMount</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>healthPort</code></td>
<td class="">The port that the health check endpoint will bind to.</td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>readinessProbe</code></td>
<td class="">The readiness probe config to be used for the Pods in this deployment. ref: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/">https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/</a></td>
<td class=""><code>&#42;<router-link to="#_readinessprobespec" @click.native="this.scrollFix('#_readinessprobespec')">ReadinessProbeSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>livenessProbe</code></td>
<td class="">The liveness probe config to be used for the Pods in this deployment. ref: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/">https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/</a></td>
<td class=""><code>&#42;<router-link to="#_readinessprobespec" @click.native="this.scrollFix('#_readinessprobespec')">ReadinessProbeSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>resources</code></td>
<td class="">Resources is the optional resource requests and limits for the containers<br>
 ref: <a id="" title="" target="_blank" href="http://kubernetes.io/docs/user-guide/compute-resources/">http://kubernetes.io/docs/user-guide/compute-resources/</a> +<br>
<br>
By default the cpu requests is set to zero and the cpu limit set to 32. This is because it appears that K8s defaults cpu to one and since Java 10 the JVM now correctly picks up cgroup cpu limits then the JVM will only see one cpu. By setting resources.requests.cpu=0 and resources.limits.cpu=32 it ensures that the JVM will see the either the number of cpus on the host if this is &#8656; 32 or the JVM will see 32 cpus if the host has &gt; 32 cpus. The limit is set to zero so that there is no hard-limit applied.<br>
<br>
No default memory limits are applied.</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core">corev1.ResourceRequirements</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>affinity</code></td>
<td class="">Affinity controls Pod scheduling preferences.<br>
  ref: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity">https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity</a><br></td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#affinity-v1-core">corev1.Affinity</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>nodeSelector</code></td>
<td class="">NodeSelector is the Node labels for pod assignment<br>
  ref: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector">https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector</a><br></td>
<td class=""><code>map[string]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>tolerations</code></td>
<td class="">Tolerations is for nodes that have taints on them.<br>
  Useful if you want to dedicate nodes to just run the coherence container +<br>
For example:<br>
  tolerations: +<br>
  - key: \"key\" +<br>
    operator: \"Equal\" +<br>
    value: \"value\" +<br>
    effect: \"NoSchedule\" +<br>
<br>
  ref: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/">https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/</a><br></td>
<td class=""><code>[]<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#toleration-v1-core">corev1.Toleration</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>securityContext</code></td>
<td class="">SecurityContext is the PodSecurityContext that will be added to all of the Pods in this deployment. See: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/configure-pod-container/security-context/">https://kubernetes.io/docs/tasks/configure-pod-container/security-context/</a></td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#podsecuritycontext-v1-core">corev1.PodSecurityContext</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>shareProcessNamespace</code></td>
<td class="">Share a single process namespace between all of the containers in a pod. When this is set containers will be able to view and signal processes from other containers in the same pod, and the first process in each container will not be assigned PID 1. HostPID and ShareProcessNamespace cannot both be set. Optional: Default to false.</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>hostIPC</code></td>
<td class="">Use the host&#8217;s ipc namespace. Optional: Default to false.</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>network</code></td>
<td class="">Configure various networks and DNS settings for Pods in this rolw.</td>
<td class=""><code>&#42;<router-link to="#_networkspec" @click.native="this.scrollFix('#_networkspec')">NetworkSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>coherenceUtils</code></td>
<td class="">The configuration for the Coherence utils image</td>
<td class=""><code>&#42;<router-link to="#_imagespec" @click.native="this.scrollFix('#_imagespec')">ImageSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>serviceAccountName</code></td>
<td class="">The name to use for the service account to use when RBAC is enabled The role bindings must already have been created as this chart does not create them it just sets the serviceAccountName value in the Pod spec.</td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>automountServiceAccountToken</code></td>
<td class="">Whether or not to auto-mount the Kubernetes API credentials for a service account</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>operatorRequestTimeout</code></td>
<td class="">The timeout to apply to REST requests made back to the Operator from Coherence Pods. These requests are typically to obtain site and rack information for the Pod.</td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_applicationspec">ApplicationSpec</h3>
<div class="section">
<p>The specification of the application deployed into the Coherence.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>type</code></td>
<td class="">The application type to execute. This field would be set if using the Coherence Graal image and running a none-Java application. For example if the application was a Node application this field would be set to \"node\". The default is to run a plain Java application.</td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>main</code></td>
<td class="">Class is the Coherence container main class.  The default value is com.tangosol.net.DefaultCacheServer. If the application type is non-Java this would be the name of the corresponding language specific runnable, for example if the application type is \"node\" the main may be a Javascript file.</td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>args</code></td>
<td class="">Args is the optional arguments to pass to the main class.</td>
<td class=""><code>[]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>workingDir</code></td>
<td class="">The application folder in the custom artifacts Docker image containing application artifacts. This will effectively become the working directory of the Coherence container. If not set the application directory default value is \"/app\".</td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_coherencespec">CoherenceSpec</h3>
<div class="section">
<p>This section of the CRD configures settings specific to Coherence.<br>
see: <router-link to="#coherence_settings/010_overview.adoc" @click.native="this.scrollFix('#coherence_settings/010_overview.adoc')">Coherence Configuration</router-link></p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>cacheConfig</code></td>
<td class="">CacheConfig is the name of the cache configuration file to use<br>
see: <router-link to="#coherence_settings/030_cache_config.adoc" @click.native="this.scrollFix('#coherence_settings/030_cache_config.adoc')">Configure Cache Config File</router-link></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>overrideConfig</code></td>
<td class="">OverrideConfig is name of the Coherence operational configuration override file, the default is tangosol-coherence-override.xml<br>
see: <router-link to="#coherence_settings/040_override_file.adoc" @click.native="this.scrollFix('#coherence_settings/040_override_file.adoc')">Configure Operational Config File</router-link></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>storageEnabled</code></td>
<td class="">A boolean flag indicating whether members of this deployment are storage enabled. This value will set the corresponding coherence.distributed.localstorage System property. If not specified the default value is true. This flag is also used to configure the ScalingPolicy value if a value is not specified. If the StorageEnabled field is not specified or is true the scaling will be safe, if StorageEnabled is set to false scaling will be parallel.<br>
see: <router-link to="#coherence_settings/050_storage_enabled.adoc" @click.native="this.scrollFix('#coherence_settings/050_storage_enabled.adoc')">Configure Storage Enabled</router-link></td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>persistence</code></td>
<td class="">Persistence values configure the on-disc data persistence settings. The bool Enabled enables or disabled on disc persistence of data.<br>
see: <router-link to="#coherence_settings/080_persistence.adoc" @click.native="this.scrollFix('#coherence_settings/080_persistence.adoc')">Configure Persistence</router-link></td>
<td class=""><code>&#42;<router-link to="#_persistencespec" @click.native="this.scrollFix('#_persistencespec')">PersistenceSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>logLevel</code></td>
<td class="">The Coherence log level, default being 5 (info level).<br>
see: <router-link to="#coherence_settings/060_log_level.adoc" @click.native="this.scrollFix('#coherence_settings/060_log_level.adoc')">Configure Coherence log level</router-link></td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>management</code></td>
<td class="">Management configures Coherence management over REST Note: Coherence management over REST will is available in Coherence version &gt;= 12.2.1.4.<br>
see: <router-link to="/management/010_overview">Management &amp; Diagnostics</router-link></td>
<td class=""><code>&#42;<router-link to="#_portspecwithssl" @click.native="this.scrollFix('#_portspecwithssl')">PortSpecWithSSL</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>metrics</code></td>
<td class="">Metrics configures Coherence metrics publishing Note: Coherence metrics publishing will is available in Coherence version &gt;= 12.2.1.4.<br>
see: <router-link to="/metrics/010_overview">Metrics</router-link></td>
<td class=""><code>&#42;<router-link to="#_portspecwithssl" @click.native="this.scrollFix('#_portspecwithssl')">PortSpecWithSSL</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>tracing</code></td>
<td class="">Tracing is used to configure Coherence distributed tracing functionality.</td>
<td class=""><code>&#42;<router-link to="#_coherencetracingspec" @click.native="this.scrollFix('#_coherencetracingspec')">CoherenceTracingSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>allowEndangeredForStatusHA</code></td>
<td class="">AllowEndangeredForStatusHA is a list of Coherence partitioned cache service names that are allowed to be in an endangered state when testing for StatusHA. Instances where a StatusHA check is performed include the readiness probe and when scaling a deployment. This field would not typically be used except in cases where a cache service is configured with a backup count greater than zero but it does not matter if caches in those services loose data due to member departure. Normally, such cache services would have a backup count of zero, which would automatically excluded them from the StatusHA check.</td>
<td class=""><code>[]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>excludeFromWKA</code></td>
<td class="">Exclude members of this deployment from being part of the cluster&#8217;s WKA list.<br>
see: <router-link to="#coherence_settings/070_wka.adoc" @click.native="this.scrollFix('#coherence_settings/070_wka.adoc')">Well Known Addressing</router-link></td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>wka</code></td>
<td class="">Specify an existing Coherence deployment to be used for WKA. If an existing deployment is to be used for WKA the ExcludeFromWKA is implicitly set to true.<br>
see: <router-link to="#coherence_settings/070_wka.adoc" @click.native="this.scrollFix('#coherence_settings/070_wka.adoc')">Well Known Addressing</router-link></td>
<td class=""><code>&#42;<router-link to="#_coherencewkaspec" @click.native="this.scrollFix('#_coherencewkaspec')">CoherenceWKASpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>skipVersionCheck</code></td>
<td class="">Certain features rely on a version check prior to starting the server, e.g. metrics requires &gt;= 12.2.1.4. The version check relies on the ability of the start script to find coherence.jar but if due to how the image has been built this check is failing then setting this flag to true will skip version checking and assume that the latest coherence.jar is being used.</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_coherencetracingspec">CoherenceTracingSpec</h3>
<div class="section">
<p>CoherenceTracingSpec configures Coherence tracing.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>ratio</code></td>
<td class="">Ratio is the tracing sampling-ratio, which controls the likelihood of a tracing span being collected. For instance, a value of 1.0 will result in all tracing spans being collected, while a value of 0.1 will result in roughly 1 out of every 10 tracing spans being collected.<br>
<br>
A value of 0 indicates that tracing spans should only be collected if they are already in the context of another tracing span.  With such a configuration, Coherence will not initiate tracing on its own, and it is up to the application to start an outer tracing span, from which Coherence will then collect inner tracing spans.<br>
<br>
A value of -1 disables tracing completely.<br>
<br>
The Coherence default is -1 if not overridden. For values other than -1, numbers between 0 and 1 are acceptable.<br>
<br>
Due to decimal values not being allowed in a CRD field the ratio value is held as a resource Quantity type and may be entered in yaml or json as a valid Quantity string as well as a decimal string. For example: A value of 0.5 is represented in json as a Quantity of \"500m\" for 500 millis. A value of 0.0005 is represented in json as a Quantity of \"500u\" for 500 micros.</td>
<td class=""><code>&#42;resource.Quantity</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_coherencewkaspec">CoherenceWKASpec</h3>
<div class="section">

<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>deployment</code></td>
<td class="">The name of the existing Coherence deployment to use for WKA.</td>
<td class=""><code>string</code></td>
<td class="">true</td>
</tr>
<tr>
<td class=""><code>namespace</code></td>
<td class="">The optional namespace of the existing Coherence deployment to use for WKA if different from this deployment&#8217;s namespace.</td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_configmapvolumespec">ConfigMapVolumeSpec</h3>
<div class="section">
<p>Represents a ConfigMap that will be added to the deployment&#8217;s Pods as an additional Volume and as a VolumeMount in the containers.<br>
see: <router-link to="/other/050_configmap_volumes">Add ConfigMap Volumes</router-link></p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>name</code></td>
<td class="">The name of the ConfigMap to mount. This will also be used as the name of the Volume added to the Pod if the VolumeName field is not set.</td>
<td class=""><code>string</code></td>
<td class="">true</td>
</tr>
<tr>
<td class=""><code>mountPath</code></td>
<td class="">Path within the container at which the volume should be mounted.  Must not contain ':'.</td>
<td class=""><code>string</code></td>
<td class="">true</td>
</tr>
<tr>
<td class=""><code>volumeName</code></td>
<td class="">The optional name to use for the Volume added to the Pod. If not set, the ConfigMap name will be used as the VolumeName.</td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>readOnly</code></td>
<td class="">Mounted read-only if true, read-write otherwise (false or unspecified). Defaults to false.</td>
<td class=""><code>bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>subPath</code></td>
<td class="">Path within the volume from which the container&#8217;s volume should be mounted. Defaults to \"\" (volume&#8217;s root).</td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>mountPropagation</code></td>
<td class="">mountPropagation determines how mounts are propagated from the host to container and the other way around. When not set, MountPropagationNone is used.</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#mountpropagationmode-v1-core">corev1.MountPropagationMode</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>subPathExpr</code></td>
<td class="">Expanded path within the volume from which the container&#8217;s volume should be mounted. Behaves similarly to SubPath but environment variable references $(VAR_NAME) are expanded using the container&#8217;s environment. Defaults to \"\" (volume&#8217;s root). SubPathExpr and SubPath are mutually exclusive.</td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>items</code></td>
<td class="">If unspecified, each key-value pair in the Data field of the referenced ConfigMap will be projected into the volume as a file whose name is the key and content is the value. If specified, the listed keys will be projected into the specified paths, and unlisted keys will not be present. If a key is specified which is not present in the ConfigMap, the volume setup will error unless it is marked optional. Paths must be relative and may not contain the '..' path or start with '..'.</td>
<td class=""><code>[]<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#keytopath-v1-core">corev1.KeyToPath</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>defaultMode</code></td>
<td class="">Optional: mode bits to use on created files by default. Must be a value between 0 and 0777. Defaults to 0644. Directories within the path are not affected by this setting. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.</td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>optional</code></td>
<td class="">Specify whether the ConfigMap or its keys must be defined</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_imagespec">ImageSpec</h3>
<div class="section">
<p>ImageSpec defines the settings for a Docker image</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>image</code></td>
<td class="">The image name. More info: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/containers/images">https://kubernetes.io/docs/concepts/containers/images</a></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>imagePullPolicy</code></td>
<td class="">Image pull policy. One of Always, Never, IfNotPresent. More info: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/containers/images#updating-images">https://kubernetes.io/docs/concepts/containers/images#updating-images</a></td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#pullpolicy-v1-core">corev1.PullPolicy</a></code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_jvmspec">JVMSpec</h3>
<div class="section">
<p>The JVM configuration.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>classpath</code></td>
<td class="">Classpath specifies additional items to add to the classpath of the JVM.</td>
<td class=""><code>[]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>args</code></td>
<td class="">Args specifies the options (System properties, -XX: args etc) to pass to the JVM.</td>
<td class=""><code>[]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>debug</code></td>
<td class="">The settings for enabling debug mode in the JVM.</td>
<td class=""><code>&#42;<router-link to="#_jvmdebugspec" @click.native="this.scrollFix('#_jvmdebugspec')">JvmDebugSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>useContainerLimits</code></td>
<td class="">If set to true Adds the  -XX:+UseContainerSupport JVM option to ensure that the JVM respects any container resource limits. The default value is true</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>gc</code></td>
<td class="">Set JVM garbage collector options.</td>
<td class=""><code>&#42;<router-link to="#_jvmgarbagecollectorspec" @click.native="this.scrollFix('#_jvmgarbagecollectorspec')">JvmGarbageCollectorSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>diagnosticsVolume</code></td>
<td class="">&#160;</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#volumesource-v1-core">corev1.VolumeSource</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>memory</code></td>
<td class="">Configure the JVM memory options.</td>
<td class=""><code>&#42;<router-link to="#_jvmmemoryspec" @click.native="this.scrollFix('#_jvmmemoryspec')">JvmMemorySpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>jmxmp</code></td>
<td class="">Configure JMX using JMXMP.</td>
<td class=""><code>&#42;<router-link to="#_jvmjmxmpspec" @click.native="this.scrollFix('#_jvmjmxmpspec')">JvmJmxmpSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>useJibClasspath</code></td>
<td class="">A flag indicating whether to automatically add the default classpath for images created by the JIB tool <a id="" title="" target="_blank" href="https://github.com/GoogleContainerTools/jib">https://github.com/GoogleContainerTools/jib</a> If true then the /app/lib/* /app/classes and /app/resources entries are added to the JVM classpath. The default value fif not specified is true.</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_jvmdebugspec">JvmDebugSpec</h3>
<div class="section">
<p>The JVM Debug specific configuration.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>enabled</code></td>
<td class="">Enabled is a flag to enable or disable running the JVM in debug mode. Default is disabled.</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>suspend</code></td>
<td class="">A boolean true if the target VM is to be suspended immediately before the main class is loaded; false otherwise. The default value is false.</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>attach</code></td>
<td class="">Attach specifies the address of the debugger that the JVM should attempt to connect back to instead of listening on a port.</td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>port</code></td>
<td class="">The port that the debugger will listen on; the default is 5005.</td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_jvmgarbagecollectorspec">JvmGarbageCollectorSpec</h3>
<div class="section">
<p>Options for managing the JVM garbage collector.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>collector</code></td>
<td class="">The name of the JVM garbage collector to use. G1 - adds the -XX:+UseG1GC option CMS - adds the -XX:+UseConcMarkSweepGC option Parallel - adds the -XX:+UseParallelGC Default - use the JVMs default collector The field value is case insensitive If not set G1 is used. If set to a value other than those above then the default collector for the JVM will be used.</td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>args</code></td>
<td class="">Args specifies the GC options to pass to the JVM.</td>
<td class=""><code>[]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>logging</code></td>
<td class="">Enable the following GC logging args  -verbose:gc -XX:+PrintGCDetails -XX:+PrintGCTimeStamps -XX:+PrintHeapAtGC -XX:+PrintTenuringDistribution -XX:+PrintGCApplicationStoppedTime -XX:+PrintGCApplicationConcurrentTime Default is true</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_jvmjmxmpspec">JvmJmxmpSpec</h3>
<div class="section">
<p>Options for configuring JMX using JMXMP.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>enabled</code></td>
<td class="">If set to true the JMXMP support will be enabled. Default is false</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>port</code></td>
<td class="">The port tht the JMXMP MBeanServer should bind to. If not set the default port is 9099</td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_jvmmemoryspec">JvmMemorySpec</h3>
<div class="section">
<p>Options for managing the JVM memory.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>heapSize</code></td>
<td class="">HeapSize is the min/max heap value to pass to the JVM. The format should be the same as that used for Java&#8217;s -Xms and -Xmx JVM options. If not set the JVM defaults are used.</td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>stackSize</code></td>
<td class="">StackSize is the stack size value to pass to the JVM. The format should be the same as that used for Java&#8217;s -Xss JVM option. If not set the JVM defaults are used.</td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>metaspaceSize</code></td>
<td class="">MetaspaceSize is the min/max metaspace size to pass to the JVM. This sets the -XX:MetaspaceSize and -XX:MaxMetaspaceSize=size JVM options. If not set the JVM defaults are used.</td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>directMemorySize</code></td>
<td class="">DirectMemorySize sets the maximum total size (in bytes) of the New I/O (the java.nio package) direct-buffer allocations. This value sets the -XX:MaxDirectMemorySize JVM option. If not set the JVM defaults are used.</td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>nativeMemoryTracking</code></td>
<td class="">Adds the -XX:NativeMemoryTracking=mode  JVM options where mode is on of \"off\", \"summary\" or \"detail\", the default is \"summary\" If not set to \"off\" also add -XX:+PrintNMTStatistics</td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>onOutOfMemory</code></td>
<td class="">Configure the JVM behaviour when an OutOfMemoryError occurs.</td>
<td class=""><code>&#42;<router-link to="#_jvmoutofmemoryspec" @click.native="this.scrollFix('#_jvmoutofmemoryspec')">JvmOutOfMemorySpec</router-link></code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_jvmoutofmemoryspec">JvmOutOfMemorySpec</h3>
<div class="section">
<p>Options for managing the JVM behaviour when an OutOfMemoryError occurs.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>exit</code></td>
<td class="">If set to true the JVM will exit when an OOM error occurs. Default is true</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>heapDump</code></td>
<td class="">If set to true adds the -XX:+HeapDumpOnOutOfMemoryError JVM option to cause a heap dump to be created when an OOM error occurs. Default is true</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_localobjectreference">LocalObjectReference</h3>
<div class="section">
<p>LocalObjectReference contains enough information to let you locate the referenced object inside the same namespace.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>name</code></td>
<td class="">Name of the referent. More info: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names">https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names</a></td>
<td class=""><code>string</code></td>
<td class="">true</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_namedportspec">NamedPortSpec</h3>
<div class="section">

<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>name</code></td>
<td class="">Name specifies the name of the port.</td>
<td class=""><code>string</code></td>
<td class="">true</td>
</tr>
<tr>
<td class=""><code>port</code></td>
<td class="">Port specifies the port used.</td>
<td class=""><code>int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>protocol</code></td>
<td class="">Protocol for container port. Must be UDP or TCP. Defaults to \"TCP\"</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#protocol-v1-core">corev1.Protocol</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>nodePort</code></td>
<td class="">The port on each node on which this service is exposed when type=NodePort or LoadBalancer. Usually assigned by the system. If specified, it will be allocated to the service if unused or else creation of the service will fail. Default is to auto-allocate a port if the ServiceType of this Service requires one. More info: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport">https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport</a></td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>hostPort</code></td>
<td class="">Number of port to expose on the host. If specified, this must be a valid port number, 0 &lt; x &lt; 65536. If HostNetwork is specified, this must match ContainerPort. Most containers do not need this.</td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>hostIP</code></td>
<td class="">What host IP to bind the external port to.</td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>service</code></td>
<td class="">Service configures the Kubernetes Service used to expose the port.</td>
<td class=""><code>&#42;<router-link to="#_servicespec" @click.native="this.scrollFix('#_servicespec')">ServiceSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>serviceMonitor</code></td>
<td class="">The specification of a Prometheus ServiceMonitor resource that will be created for the Service being exposed for this port.</td>
<td class=""><code>&#42;<router-link to="#_servicemonitorspec" @click.native="this.scrollFix('#_servicemonitorspec')">ServiceMonitorSpec</router-link></code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_networkspec">NetworkSpec</h3>
<div class="section">
<p>NetworkSpec configures various networking and DNS settings for Pods in a deployment.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>dnsConfig</code></td>
<td class="">Specifies the DNS parameters of a pod. Parameters specified here will be merged to the generated DNS configuration based on DNSPolicy.</td>
<td class=""><code>&#42;<router-link to="#_poddnsconfig" @click.native="this.scrollFix('#_poddnsconfig')">PodDNSConfig</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>dnsPolicy</code></td>
<td class="">Set DNS policy for the pod. Defaults to \"ClusterFirst\". Valid values are 'ClusterFirstWithHostNet', 'ClusterFirst', 'Default' or 'None'. DNS parameters given in DNSConfig will be merged with the policy selected with DNSPolicy. To have DNS options set along with hostNetwork, you have to specify DNS policy explicitly to 'ClusterFirstWithHostNet'.</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#dnspolicy-v1-core">corev1.DNSPolicy</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>hostAliases</code></td>
<td class="">HostAliases is an optional list of hosts and IPs that will be injected into the pod&#8217;s hosts file if specified. This is only valid for non-hostNetwork pods.</td>
<td class=""><code>[]<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#hostalias-v1-core">corev1.HostAlias</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>hostNetwork</code></td>
<td class="">Host networking requested for this pod. Use the host&#8217;s network namespace. If this option is set, the ports that will be used must be specified. Default to false.</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>hostname</code></td>
<td class="">Specifies the hostname of the Pod If not specified, the pod&#8217;s hostname will be set to a system-defined value.</td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_persistencespec">PersistenceSpec</h3>
<div class="section">
<p>The spec for Coherence persistence.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>mode</code></td>
<td class="">The persistence mode to use. Valid choices are \"on-demand\", \"active\", \"active-async\". This field will set the coherence.distributed.persistence-mode System property to \"default-\" + Mode.</td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>persistentVolumeClaim</code></td>
<td class="">PersistentVolumeClaim allows the configuration of a normal k8s persistent volume claim for persistence data.</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#persistentvolumeclaimspec-v1-core">corev1.PersistentVolumeClaimSpec</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>volume</code></td>
<td class="">Volume allows the configuration of a normal k8s volume mapping for persistence data instead of a persistent volume claim. If a value is defined for store.persistence.volume then no PVC will be created and persistence data will instead be written to this volume. It is up to the deployer to understand the consequences of this and how the guarantees given when using PVCs differ to the storage guarantees for the particular volume type configured here.</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#volumesource-v1-core">corev1.VolumeSource</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>snapshots</code></td>
<td class="">Snapshot values configure the on-disc persistence data snapshot (backup) settings. These settings enable a different location for persistence snapshot data. If not set then snapshot files will be written to the same volume configured for persistence data in the Persistence section.</td>
<td class=""><code>&#42;<router-link to="#_persistentstoragespec" @click.native="this.scrollFix('#_persistentstoragespec')">PersistentStorageSpec</router-link></code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_persistentstoragespec">PersistentStorageSpec</h3>
<div class="section">
<p>PersistenceStorageSpec defines the persistence settings for the Coherence</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>persistentVolumeClaim</code></td>
<td class="">PersistentVolumeClaim allows the configuration of a normal k8s persistent volume claim for persistence data.</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#persistentvolumeclaimspec-v1-core">corev1.PersistentVolumeClaimSpec</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>volume</code></td>
<td class="">Volume allows the configuration of a normal k8s volume mapping for persistence data instead of a persistent volume claim. If a value is defined for store.persistence.volume then no PVC will be created and persistence data will instead be written to this volume. It is up to the deployer to understand the consequences of this and how the guarantees given when using PVCs differ to the storage guarantees for the particular volume type configured here.</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#volumesource-v1-core">corev1.VolumeSource</a></code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_poddnsconfig">PodDNSConfig</h3>
<div class="section">
<p>PodDNSConfig defines the DNS parameters of a pod in addition to those generated from DNSPolicy.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>nameservers</code></td>
<td class="">A list of DNS name server IP addresses. This will be appended to the base nameservers generated from DNSPolicy. Duplicated nameservers will be removed.</td>
<td class=""><code>[]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>searches</code></td>
<td class="">A list of DNS search domains for host-name lookup. This will be appended to the base search paths generated from DNSPolicy. Duplicated search paths will be removed.</td>
<td class=""><code>[]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>options</code></td>
<td class="">A list of DNS resolver options. This will be merged with the base options generated from DNSPolicy. Duplicated entries will be removed. Resolution options given in Options will override those that appear in the base DNSPolicy.</td>
<td class=""><code>[]<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#poddnsconfigoption-v1-core">corev1.PodDNSConfigOption</a></code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_portspecwithssl">PortSpecWithSSL</h3>
<div class="section">
<p>PortSpecWithSSL defines a port with SSL settings for a Coherence component</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>enabled</code></td>
<td class="">Enable or disable flag.</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>port</code></td>
<td class="">The port to bind to.</td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>ssl</code></td>
<td class="">SSL configures SSL settings for a Coherence component</td>
<td class=""><code>&#42;<router-link to="#_sslspec" @click.native="this.scrollFix('#_sslspec')">SSLSpec</router-link></code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_probehandler">ProbeHandler</h3>
<div class="section">
<p>The definition of a probe handler.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>exec</code></td>
<td class="">One and only one of the following should be specified. Exec specifies the action to take.</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#execaction-v1-core">corev1.ExecAction</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>httpGet</code></td>
<td class="">HTTPGet specifies the http request to perform.</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#httpgetaction-v1-core">corev1.HTTPGetAction</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>tcpSocket</code></td>
<td class="">TCPSocket specifies an action involving a TCP port. TCP hooks not yet supported</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#tcpsocketaction-v1-core">corev1.TCPSocketAction</a></code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_readinessprobespec">ReadinessProbeSpec</h3>
<div class="section">
<p>ReadinessProbeSpec defines the settings for the Coherence Pod readiness probe</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>exec</code></td>
<td class="">One and only one of the following should be specified. Exec specifies the action to take.</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#execaction-v1-core">corev1.ExecAction</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>httpGet</code></td>
<td class="">HTTPGet specifies the http request to perform.</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#httpgetaction-v1-core">corev1.HTTPGetAction</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>tcpSocket</code></td>
<td class="">TCPSocket specifies an action involving a TCP port. TCP hooks not yet supported</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#tcpsocketaction-v1-core">corev1.TCPSocketAction</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>initialDelaySeconds</code></td>
<td class="">Number of seconds after the container has started before liveness probes are initiated. More info: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes">https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes</a></td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>timeoutSeconds</code></td>
<td class="">Number of seconds after which the probe times out. More info: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes">https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes</a></td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>periodSeconds</code></td>
<td class="">How often (in seconds) to perform the probe.</td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>successThreshold</code></td>
<td class="">Minimum consecutive successes for the probe to be considered successful after having failed.</td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>failureThreshold</code></td>
<td class="">Minimum consecutive failures for the probe to be considered failed after having succeeded.</td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_resource">Resource</h3>
<div class="section">

<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>kind</code></td>
<td class="">&#160;</td>
<td class=""><code>ResourceType</code></td>
<td class="">true</td>
</tr>
<tr>
<td class=""><code>name</code></td>
<td class="">&#160;</td>
<td class=""><code>string</code></td>
<td class="">true</td>
</tr>
<tr>
<td class=""><code>spec</code></td>
<td class="">&#160;</td>
<td class=""><code>runtime.Object</code></td>
<td class="">true</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_resources">Resources</h3>
<div class="section">

<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>version</code></td>
<td class="">&#160;</td>
<td class=""><code>int32</code></td>
<td class="">true</td>
</tr>
<tr>
<td class=""><code>items</code></td>
<td class="">&#160;</td>
<td class=""><code>[]<router-link to="#_resource" @click.native="this.scrollFix('#_resource')">Resource</router-link></code></td>
<td class="">true</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_sslspec">SSLSpec</h3>
<div class="section">
<p>SSLSpec defines the SSL settings for a Coherence component over REST endpoint.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>enabled</code></td>
<td class="">Enabled is a boolean flag indicating whether enables or disables SSL on the Coherence management over REST endpoint, the default is false (disabled).</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>secrets</code></td>
<td class="">Secrets is the name of the k8s secrets containing the Java key stores and password files.<br>
  This value MUST be provided if SSL is enabled on the Coherence management over ReST endpoint.<br></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>keyStore</code></td>
<td class="">Keystore is the name of the Java key store file in the k8s secret to use as the SSL keystore<br>
  when configuring component over REST to use SSL.<br></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>keyStorePasswordFile</code></td>
<td class="">KeyStorePasswordFile is the name of the file in the k8s secret containing the keystore<br>
  password when configuring component over REST to use SSL.<br></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>keyPasswordFile</code></td>
<td class="">KeyStorePasswordFile is the name of the file in the k8s secret containing the key<br>
  password when configuring component over REST to use SSL.<br></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>keyStoreAlgorithm</code></td>
<td class="">KeyStoreAlgorithm is the name of the keystore algorithm for the keystore in the k8s secret<br>
  used when configuring component over REST to use SSL. If not set the default is SunX509<br></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>keyStoreProvider</code></td>
<td class="">KeyStoreProvider is the name of the keystore provider for the keystore in the k8s secret<br>
  used when configuring component over REST to use SSL.<br></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>keyStoreType</code></td>
<td class="">KeyStoreType is the name of the Java keystore type for the keystore in the k8s secret used<br>
  when configuring component over REST to use SSL. If not set the default is JKS.<br></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>trustStore</code></td>
<td class="">TrustStore is the name of the Java trust store file in the k8s secret to use as the SSL<br>
  trust store when configuring component over REST to use SSL.<br></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>trustStorePasswordFile</code></td>
<td class="">TrustStorePasswordFile is the name of the file in the k8s secret containing the trust store<br>
  password when configuring component over REST to use SSL.<br></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>trustStoreAlgorithm</code></td>
<td class="">TrustStoreAlgorithm is the name of the keystore algorithm for the trust store in the k8s<br>
  secret used when configuring component over REST to use SSL.  If not set the default is SunX509.<br></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>trustStoreProvider</code></td>
<td class="">TrustStoreProvider is the name of the keystore provider for the trust store in the k8s<br>
  secret used when configuring component over REST to use SSL.<br></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>trustStoreType</code></td>
<td class="">TrustStoreType is the name of the Java keystore type for the trust store in the k8s secret<br>
  used when configuring component over REST to use SSL. If not set the default is JKS.<br></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>requireClientCert</code></td>
<td class="">RequireClientCert is a boolean flag indicating whether the client certificate will be<br>
  authenticated by the server (two-way SSL) when configuring component over REST to use SSL. +<br>
  If not set the default is false<br></td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_scalingprobe">ScalingProbe</h3>
<div class="section">
<p>ScalingProbe is the handler that will be used to determine how to check for StatusHA in a Coherence. StatusHA checking is primarily used during scaling of a deployment, a deployment must be in a safe Phase HA state before scaling takes place. If StatusHA handler is disabled for a deployment (by specifically setting Enabled to false then no check will take place and a deployment will be assumed to be safe).</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>timeoutSeconds</code></td>
<td class="">Number of seconds after which the handler times out (only applies to http and tcp handlers). Defaults to 1 second. Minimum value is 1.</td>
<td class=""><code>&#42;int</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_scalingspec">ScalingSpec</h3>
<div class="section">
<p>The configuration to control safe scaling.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>policy</code></td>
<td class="">ScalingPolicy describes how the replicas of the deployment will be scaled. The default if not specified is based upon the value of the StorageEnabled field. If StorageEnabled field is not specified or is true the default scaling will be safe, if StorageEnabled is set to false the default scaling will be parallel.</td>
<td class=""><code>&#42;ScalingPolicy</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>probe</code></td>
<td class="">The probe to use to determine whether a deployment is Phase HA. If not set the default handler will be used. In most use-cases the default handler would suffice but in advanced use-cases where the application code has a different concept of Phase HA to just checking Coherence services then a different handler may be specified.</td>
<td class=""><code>&#42;<router-link to="#_scalingprobe" @click.native="this.scrollFix('#_scalingprobe')">ScalingProbe</router-link></code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_secretvolumespec">SecretVolumeSpec</h3>
<div class="section">
<p>Represents a Secret that will be added to the deployment&#8217;s Pods as an additional Volume and as a VolumeMount in the containers.<br>
see: <router-link to="#other/020_secret_volumes.adoc" @click.native="this.scrollFix('#other/020_secret_volumes.adoc')">Add Secret Volumes</router-link></p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>name</code></td>
<td class="">The name of the Secret to mount. This will also be used as the name of the Volume added to the Pod if the VolumeName field is not set.</td>
<td class=""><code>string</code></td>
<td class="">true</td>
</tr>
<tr>
<td class=""><code>mountPath</code></td>
<td class="">Path within the container at which the volume should be mounted.  Must not contain ':'.</td>
<td class=""><code>string</code></td>
<td class="">true</td>
</tr>
<tr>
<td class=""><code>volumeName</code></td>
<td class="">The optional name to use for the Volume added to the Pod. If not set, the Secret name will be used as the VolumeName.</td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>readOnly</code></td>
<td class="">Mounted read-only if true, read-write otherwise (false or unspecified). Defaults to false.</td>
<td class=""><code>bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>subPath</code></td>
<td class="">Path within the volume from which the container&#8217;s volume should be mounted. Defaults to \"\" (volume&#8217;s root).</td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>mountPropagation</code></td>
<td class="">mountPropagation determines how mounts are propagated from the host to container and the other way around. When not set, MountPropagationNone is used.</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#mountpropagationmode-v1-core">corev1.MountPropagationMode</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>subPathExpr</code></td>
<td class="">Expanded path within the volume from which the container&#8217;s volume should be mounted. Behaves similarly to SubPath but environment variable references $(VAR_NAME) are expanded using the container&#8217;s environment. Defaults to \"\" (volume&#8217;s root). SubPathExpr and SubPath are mutually exclusive.</td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>items</code></td>
<td class="">If unspecified, each key-value pair in the Data field of the referenced Secret will be projected into the volume as a file whose name is the key and content is the value. If specified, the listed keys will be projected into the specified paths, and unlisted keys will not be present. If a key is specified which is not present in the Secret, the volume setup will error unless it is marked optional. Paths must be relative and may not contain the '..' path or start with '..'.</td>
<td class=""><code>[]<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#keytopath-v1-core">corev1.KeyToPath</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>defaultMode</code></td>
<td class="">Optional: mode bits to use on created files by default. Must be a value between 0 and 0777. Defaults to 0644. Directories within the path are not affected by this setting. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.</td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>optional</code></td>
<td class="">Specify whether the Secret or its keys must be defined</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_servicemonitorspec">ServiceMonitorSpec</h3>
<div class="section">
<p>The ServiceMonitor spec for a port service.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>enabled</code></td>
<td class="">Enabled is a flag to enable or disable creation of a Prometheus ServiceMonitor for a port. If Prometheus ServiceMonitor CR is not installed no ServiceMonitor then even if this flag is true no ServiceMonitor will be created.</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>labels</code></td>
<td class="">Additional labels to add to the ServiceMonitor. More info: <a id="" title="" target="_blank" href="http://kubernetes.io/docs/user-guide/labels">http://kubernetes.io/docs/user-guide/labels</a></td>
<td class=""><code>map[string]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>jobLabel</code></td>
<td class="">The label to use to retrieve the job name from. See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#servicemonitorspec">https://coreos.com/operators/prometheus/docs/latest/api.html#servicemonitorspec</a></td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>targetLabels</code></td>
<td class="">TargetLabels transfers labels on the Kubernetes Service onto the target. See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#servicemonitorspec">https://coreos.com/operators/prometheus/docs/latest/api.html#servicemonitorspec</a></td>
<td class=""><code>[]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>podTargetLabels</code></td>
<td class="">PodTargetLabels transfers labels on the Kubernetes Pod onto the target. See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#servicemonitorspec">https://coreos.com/operators/prometheus/docs/latest/api.html#servicemonitorspec</a></td>
<td class=""><code>[]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>sampleLimit</code></td>
<td class="">SampleLimit defines per-scrape limit on number of scraped samples that will be accepted. See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#servicemonitorspec">https://coreos.com/operators/prometheus/docs/latest/api.html#servicemonitorspec</a></td>
<td class=""><code>uint64</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>path</code></td>
<td class="">HTTP path to scrape for metrics. See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint">https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint</a></td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>scheme</code></td>
<td class="">HTTP scheme to use for scraping. See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint">https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint</a></td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>params</code></td>
<td class="">Optional HTTP URL parameters See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint">https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint</a></td>
<td class=""><code>map[string][]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>interval</code></td>
<td class="">Interval at which metrics should be scraped See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint">https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint</a></td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>scrapeTimeout</code></td>
<td class="">Timeout after which the scrape is ended See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint">https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint</a></td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>tlsConfig</code></td>
<td class="">TLS configuration to use when scraping the endpoint See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint">https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint</a></td>
<td class=""><code>&#42;monitoringv1.TLSConfig</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>bearerTokenFile</code></td>
<td class="">File to read bearer token for scraping targets. See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint">https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint</a></td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>bearerTokenSecret</code></td>
<td class="">Secret to mount to read bearer token for scraping targets. The secret needs to be in the same namespace as the service monitor and accessible by the Prometheus Operator. See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint">https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint</a></td>
<td class=""><code><a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#secretkeyselector-v1-core">corev1.SecretKeySelector</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>honorLabels</code></td>
<td class="">HonorLabels chooses the metric&#8217;s labels on collisions with target labels. See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint">https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint</a></td>
<td class=""><code>bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>honorTimestamps</code></td>
<td class="">HonorTimestamps controls whether Prometheus respects the timestamps present in scraped data. See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint">https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint</a></td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>basicAuth</code></td>
<td class="">BasicAuth allow an endpoint to authenticate over basic authentication More info: <a id="" title="" target="_blank" href="https://prometheus.io/docs/operating/configuration/#endpoints">https://prometheus.io/docs/operating/configuration/#endpoints</a> See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint">https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint</a></td>
<td class=""><code>&#42;monitoringv1.BasicAuth</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>metricRelabelings</code></td>
<td class="">MetricRelabelings to apply to samples before ingestion. See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint">https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint</a></td>
<td class=""><code>[]&#42;monitoringv1.RelabelConfig</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>relabelings</code></td>
<td class="">Relabelings to apply to samples before scraping. More info: <a id="" title="" target="_blank" href="https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config">https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config</a> See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint">https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint</a></td>
<td class=""><code>[]&#42;monitoringv1.RelabelConfig</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>proxyURL</code></td>
<td class="">ProxyURL eg <a id="" title="" target="_blank" href="http://proxyserver:2195">http://proxyserver:2195</a> Directs scrapes to proxy through this endpoint. See <a id="" title="" target="_blank" href="https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint">https://coreos.com/operators/prometheus/docs/latest/api.html#endpoint</a></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_servicespec">ServiceSpec</h3>
<div class="section">

<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>enabled</code></td>
<td class="">Enabled controls whether to create the service yaml or not</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>name</code></td>
<td class="">An optional name to use to override the generated service name.</td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>port</code></td>
<td class="">The service port value</td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>type</code></td>
<td class="">Kind is the K8s service type (typically ClusterIP or LoadBalancer) The default is \"ClusterIP\".</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#servicetype-v1-core">corev1.ServiceType</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>externalIPs</code></td>
<td class="">externalIPs is a list of IP addresses for which nodes in the cluster will also accept traffic for this service.  These IPs are not managed by Kubernetes.  The user is responsible for ensuring that traffic arrives at a node with this IP.  A common example is external load-balancers that are not part of the Kubernetes system.</td>
<td class=""><code>[]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>clusterIP</code></td>
<td class="">clusterIP is the IP address of the service and is usually assigned randomly by the master. If an address is specified manually and is not in use by others, it will be allocated to the service; otherwise, creation of the service will fail. This field can not be changed through updates. Valid values are \"None\", empty string (\"\"), or a valid IP address. \"None\" can be specified for headless services when proxying is not required. Only applies to types ClusterIP, NodePort, and LoadBalancer. Ignored if type is ExternalName. More info: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies">https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies</a></td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>loadBalancerIP</code></td>
<td class="">LoadBalancerIP is the IP address of the load balancer</td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>labels</code></td>
<td class="">The extra labels to add to the service. More info: <a id="" title="" target="_blank" href="http://kubernetes.io/docs/user-guide/labels">http://kubernetes.io/docs/user-guide/labels</a></td>
<td class=""><code>map[string]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>annotations</code></td>
<td class="">Annotations is free form yaml that will be added to the service annotations</td>
<td class=""><code>map[string]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>sessionAffinity</code></td>
<td class="">Supports \"ClientIP\" and \"None\". Used to maintain session affinity. Enable client IP based session affinity. Must be ClientIP or None. Defaults to None. More info: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies">https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies</a></td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#serviceaffinity-v1-core">corev1.ServiceAffinity</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>loadBalancerSourceRanges</code></td>
<td class="">If specified and supported by the platform, this will restrict traffic through the cloud-provider load-balancer will be restricted to the specified client IPs. This field will be ignored if the cloud-provider does not support the feature.\" More info: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/access-application-cluster/configure-cloud-provider-firewall/">https://kubernetes.io/docs/tasks/access-application-cluster/configure-cloud-provider-firewall/</a></td>
<td class=""><code>[]string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>externalName</code></td>
<td class="">externalName is the external reference that kubedns or equivalent will return as a CNAME record for this service. No proxying will be involved. Must be a valid RFC-1123 hostname (<a id="" title="" target="_blank" href="https://tools.ietf.org/html/rfc1123">https://tools.ietf.org/html/rfc1123</a>) and requires Kind to be ExternalName.</td>
<td class=""><code>&#42;string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>externalTrafficPolicy</code></td>
<td class="">externalTrafficPolicy denotes if this Service desires to route external traffic to node-local or cluster-wide endpoints. \"Local\" preserves the client source IP and avoids a second hop for LoadBalancer and Nodeport type services, but risks potentially imbalanced traffic spreading. \"Cluster\" obscures the client source IP and may cause a second hop to another node, but should have good overall load-spreading.</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#serviceexternaltrafficpolicytype-v1-core">corev1.ServiceExternalTrafficPolicyType</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>healthCheckNodePort</code></td>
<td class="">healthCheckNodePort specifies the healthcheck nodePort for the service. If not specified, HealthCheckNodePort is created by the service api backend with the allocated nodePort. Will use user-specified nodePort value if specified by the client. Only effects when Kind is set to LoadBalancer and ExternalTrafficPolicy is set to Local.</td>
<td class=""><code>&#42;int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>publishNotReadyAddresses</code></td>
<td class="">publishNotReadyAddresses, when set to true, indicates that DNS implementations must publish the notReadyAddresses of subsets for the Endpoints associated with the Service. The default value is false. The primary use case for setting this field is to use a StatefulSet&#8217;s Headless Service to propagate SRV records for its Pods without respect to their readiness for purpose of peer discovery.</td>
<td class=""><code>&#42;bool</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>sessionAffinityConfig</code></td>
<td class="">sessionAffinityConfig contains the configurations of session affinity.</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#sessionaffinityconfig-v1-core">corev1.SessionAffinityConfig</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>ipFamily</code></td>
<td class="">ipFamily specifies whether this Service has a preference for a particular IP family (e.g. IPv4 vs. IPv6).  If a specific IP family is requested, the clusterIP field will be allocated from that family, if it is available in the cluster.  If no IP family is requested, the cluster&#8217;s primary IP family will be used. Other IP fields (loadBalancerIP, loadBalancerSourceRanges, externalIPs) and controllers which allocate external load-balancers should use the same IP family.  Endpoints for this Service will be of this family.  This field is immutable after creation. Assigning a ServiceIPFamily not available in the cluster (e.g. IPv6 in IPv4 only cluster) is an error condition and will fail during clusterIP assignment.</td>
<td class=""><code>&#42;<a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#ipfamily-v1-core">corev1.IPFamily</a></code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_startquorum">StartQuorum</h3>
<div class="section">
<p>StartQuorum defines the order that deployments will be started in a Coherence cluster made up of multiple deployments.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>deployment</code></td>
<td class="">The name of deployment that this deployment depends on.</td>
<td class=""><code>string</code></td>
<td class="">true</td>
</tr>
<tr>
<td class=""><code>namespace</code></td>
<td class="">The namespace that the deployment that this deployment depends on is installed into. Default to the same namespace as this deployment</td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>podCount</code></td>
<td class="">The number of the Pods that should have been started before this deployments will be started, defaults to all Pods for the deployment.</td>
<td class=""><code>int32</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_startquorumstatus">StartQuorumStatus</h3>
<div class="section">
<p>StartQuorumStatus tracks the state of a deployment&#8217;s start quorums.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>deployment</code></td>
<td class="">The name of deployment that this deployment depends on.</td>
<td class=""><code>string</code></td>
<td class="">true</td>
</tr>
<tr>
<td class=""><code>namespace</code></td>
<td class="">The namespace that the deployment that this deployment depends on is installed into. Default to the same namespace as this deployment</td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>podCount</code></td>
<td class="">The number of the Pods that should have been started before this deployments will be started, defaults to all Pods for the deployment.</td>
<td class=""><code>int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>ready</code></td>
<td class="">Whether this quorum&#8217;s condition has been met</td>
<td class=""><code>bool</code></td>
<td class="">true</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_coherence">Coherence</h3>
<div class="section">
<p>Coherence is the Schema for the Coherence API.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>metadata</code></td>
<td class="">&#160;</td>
<td class=""><code><a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta">metav1.ObjectMeta</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>spec</code></td>
<td class="">&#160;</td>
<td class=""><code><router-link to="#_coherenceresourcespec" @click.native="this.scrollFix('#_coherenceresourcespec')">CoherenceResourceSpec</router-link></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>status</code></td>
<td class="">&#160;</td>
<td class=""><code><router-link to="#_coherenceresourcestatus" @click.native="this.scrollFix('#_coherenceresourcestatus')">CoherenceResourceStatus</router-link></code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_coherencelist">CoherenceList</h3>
<div class="section">
<p>CoherenceList contains a list of Coherence resources.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>metadata</code></td>
<td class="">&#160;</td>
<td class=""><code><a id="" title="" target="_blank" href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta">metav1.ListMeta</a></code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>items</code></td>
<td class="">&#160;</td>
<td class=""><code>[]<router-link to="#_coherence" @click.native="this.scrollFix('#_coherence')">Coherence</router-link></code></td>
<td class="">true</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>

<h3 id="_coherenceresourcestatus">CoherenceResourceStatus</h3>
<div class="section">
<p>CoherenceResourceStatus defines the observed state of Coherence resource.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 7.692%;">
<col style="width: 76.923%;">
<col style="width: 7.692%;">
<col style="width: 7.692%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
<th>Required</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>phase</code></td>
<td class="">The phase of a Coherence resource is a simple, high-level summary of where the Coherence resource is in its lifecycle. The conditions array, the reason and message fields, and the individual container status arrays contain more detail about the pod&#8217;s status. There are eight possible phase values:<br>
<br>
Initialized:    The deployment has been accepted by the Kubernetes system. Created:        The deployments secondary resources, (e.g. the StatefulSet, Services etc) have been created. Ready:          The StatefulSet for the deployment has the correct number of replicas and ready replicas. Waiting:        The deployment&#8217;s start quorum conditions have not yet been met. Scaling:        The number of replicas in the deployment is being scaled up or down. RollingUpgrade: The StatefulSet is performing a rolling upgrade. Stopped:        The replica count has been set to zero. Failed:         An error occurred reconciling the deployment and its secondary resources.</td>
<td class=""><code>status.ConditionType</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>coherenceCluster</code></td>
<td class="">The name of the Coherence cluster that this deployment is part of.</td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>replicas</code></td>
<td class="">Replicas is the desired number of members in the Coherence deployment represented by the Coherence resource.</td>
<td class=""><code>int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>currentReplicas</code></td>
<td class="">CurrentReplicas is the current number of members in the Coherence deployment represented by the Coherence resource.</td>
<td class=""><code>int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>readyReplicas</code></td>
<td class="">ReadyReplicas is the number of number of members in the Coherence deployment represented by the Coherence resource that are in the ready state.</td>
<td class=""><code>int32</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>role</code></td>
<td class="">The effective role name for this deployment. This will come from the Spec.Role field if set otherwise the deployment name will be used for the role name</td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>selector</code></td>
<td class="">label query over deployments that should match the replicas count. This is same as the label selector but in the string format to avoid introspection by clients. The string will be in the same format as the query-param syntax. More info about label selectors: <a id="" title="" target="_blank" href="http://kubernetes.io/docs/user-guide/labels#label-selectors">http://kubernetes.io/docs/user-guide/labels#label-selectors</a></td>
<td class=""><code>string</code></td>
<td class="">false</td>
</tr>
<tr>
<td class=""><code>conditions</code></td>
<td class="">The status conditions.</td>
<td class=""><code>status.Conditions</code></td>
<td class="">false</td>
</tr>
</tbody>
</table>
</div>
<p><router-link to="#_table_of_contents" @click.native="this.scrollFix('#_table_of_contents')">Back to TOC</router-link></p>

</div>
</div>
</doc-view>