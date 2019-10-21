<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>CoherenceCluster CRD Overview</dt>
<dd slot="desc"><p>Creating a Coherence cluster using the Coherence Operator is as simple as creating any other Kubernetes resource.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_coherencecluster_crd_overview">CoherenceCluster CRD Overview</h2>
<div class="section">
<p>The Coherence Operator uses a Kubernetes custom resource definition, (CRD) named <code>CoherenceCluster</code> to define the
configuration for a Coherence cluster.
All of the fields in the <code>CoherenceCluster</code> CRD are optional and a Coherence cluster can be created with a simple yaml
file:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: my-cluster  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>metadata.name</code> field will be used as the Coherence cluster name.</li>
</ul>
<p>The yaml above will create a Coherence cluster with three storage enabled members.
There is not much that can actually be achived with this cluster because no ports are exposed outside of Kubernetes
so the cluster is inaccessible. It could be possibly be accessed by other <code>Pods</code> in the same Kubernetes cluster but
in most use cases additional configuration would be required.</p>

</div>

<h2 id="_coherence_roles">Coherence Roles</h2>
<div class="section">
<p>A role is what is actually configured in the <code>CoherenceCluster</code> spec. In a traditional Coherence application that may have
had a number of storage enabled members and a number of storage disable Coherence*Extend proxy members this cluster would
have effectively had two roles, "storage" and "proxy".
Some clusters may simply have just a storage role and some complex Coherence applications and clusters may have many roles
and even different roles storage enabled for different caches/services within the same cluster.</p>

<p>A role in a <code>CoherenceCluster</code> is either configured as a single implicit <code>role</code> or one or more explicit <code>roles</code>.</p>

<markup
lang="yaml"
title="Single Implicit Role"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: my-cluster
spec:               <span class="conum" data-value="1" />
  replicas: 6</markup>

<ul class="colist">
<li data-value="1">The configuration for the <code>role</code> (in this case just the <code>replicas</code> field) is added directly to the <code>spec</code> section
of the <code>CoherenceCluster</code>.</li>
</ul>
<markup
lang="yaml"
title="Single Explicit Role"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: my-cluster
spec:
  roles:
  - role: data <span class="conum" data-value="1" />
    replicas: 6</markup>

<ul class="colist">
<li data-value="1">The configuration for a single explicit <code>role</code> named <code>data</code> is added to the <code>roles</code> list.
of the <code>CoherenceCluster</code>.</li>
</ul>
<markup
lang="yaml"
title="Multiple Explicit Roles"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: my-cluster
spec:
  roles:
  - role: data   <span class="conum" data-value="1" />
    replicas: 6
  - role: data   <span class="conum" data-value="2" />
    replicas: 3</markup>

<ul class="colist">
<li data-value="1">The first role in the <code>roles</code> list is named <code>data</code> with a <code>replicas</code> value of <code>6</code></li>
<li data-value="2">The second role in the <code>roles</code> list is named <code>proxy</code> with a <code>replicas</code> value of <code>3</code></li>
</ul>
</div>

<h2 id="_the_coherence_role_specification">The Coherence Role Specification</h2>
<div class="section">
<p>The specification for a <code>role</code> in the <code>CoherenceCluster</code> CRD (both implicit or expilict) has the following top level
fields that may be configured:</p>

<markup
lang="yaml"

>  role:                      <span class="conum" data-value="1" />
  replicas:                  <span class="conum" data-value="2" />
  application: {}            <span class="conum" data-value="3" />
  coherence: {}              <span class="conum" data-value="4" />
  jvm: {}                    <span class="conum" data-value="5" />
  scaling: {}                <span class="conum" data-value="6" />
  ports: []                  <span class="conum" data-value="7" />
  logging: {}                <span class="conum" data-value="8" />
  volumes: []                <span class="conum" data-value="9" />
  volumeClaimTemplates: []   <span class="conum" data-value="10" />
  volumeMounts: []           <span class="conum" data-value="11" />
  env: []                    <span class="conum" data-value="12" />
  annotations: {}            <span class="conum" data-value="13" />
  labels: []                 <span class="conum" data-value="14" />
  nodeSelector: {}           <span class="conum" data-value="15" />
  tolerations: []            <span class="conum" data-value="16" />
  affinity: {}               <span class="conum" data-value="17" />
  resources: {}              <span class="conum" data-value="18" />
  readinessProbe: {}         <span class="conum" data-value="19" />
  livenessProbe: {}          <span class="conum" data-value="20" /></markup>

<ul class="colist">
<li data-value="1">The <code>role</code> field sets the name of the role, if omitted the default name of <code>storage</code> will be used. If configuring
multiple roles in a <code>CoherenceCluster</code> each role must have a unique name.
See <router-link to="/clusters/030_roles">Define Coherence Roles</router-link> for more details.</li>
<li data-value="2">The <code>replicas</code> field sets the number of replicas (<code>Pods</code>) that will be vreated for the role. If not specified the
default value is <code>3</code>.
See configuring a <router-link to="/clusters/040_replicas">Role&#8217;s Replica Count</router-link> for more details.</li>
<li data-value="3">The <code>application</code> section contains fields for configuring custom application code.
See <router-link to="/clusters/070_applications">Configure Applications</router-link> for more details.</li>
<li data-value="4">The <code>coherence</code> section contains fields for configuring Coherence specific settings.
See <router-link to="/clusters/050_coherence">Configure Coherence</router-link> for more details.</li>
<li data-value="5">The <code>jvm</code> section contains fields for configuring how the JVM behaves.
See <router-link to="/clusters/080_jvm">Configure the JVM</router-link> for more details.</li>
<li data-value="6">The <code>scaling</code> section contains fields for configuring how the number of replicas in a role is safely scaled up
and down.
See <router-link to="/clusters/085_safe_scaling">Configure Safe Scaling</router-link> for more details.</li>
<li data-value="7">The <code>ports</code> section contains fields for configuring how ports are exposed via services.
See <router-link to="/clusters/090_ports_and_services">Expose Ports and Services</router-link> for more details.</li>
<li data-value="8">The <code>logging</code> section contains fields for configuring logging.
See <router-link to="/clusters/100_logging">Configure Logging</router-link> for more details.</li>
<li data-value="9">The <code>volumes</code> section contains fields for configuring additional volumes to add to the <code>Pods</code> for a role.
See <router-link to="/clusters/110_volumes">Configure Additional Volumes</router-link> for more details.</li>
<li data-value="10">The <code>volumeClaimTemplates</code> section contains fields for configuring additional PVCs to add to the <code>Pods</code> for a role.
See <router-link to="/clusters/110_volumes">Configure Additional Volumes</router-link> for more details.</li>
<li data-value="11">The <code>volumeMounts</code> section contains fields for configuring additional volume mounts to add to the <code>Pods</code> for a role.
See <router-link to="/clusters/110_volumes">Configure Additional Volumes</router-link> for more details.</li>
<li data-value="12">The <code>env</code> section contains extra environment variables to add to the Coherence container.
See <router-link to="/clusters/115_environment_variables">Environment Variables</router-link> for more details.</li>
<li data-value="13">The <code>annotations</code> map contains extra annotations to add to the <code>Pods</code> for the role.
See <router-link to="/clusters/120_annotations">Configure Pod Annotations</router-link> for more details.</li>
<li data-value="14">The <code>labels</code> map contains extra labels to add to the <code>Pods</code> for the role.
See <router-link to="/clusters/125_labels">Configure Pod Labels</router-link> for more details.</li>
<li data-value="15">The <code>nodeSelector</code> map contains node selectors to determine how Kubernetes schedules the <code>Pods</code> in the role.
See <router-link to="/clusters/130_pod_scheduling">Configure Pod Scheduling</router-link> for more details.</li>
<li data-value="16">The <code>tolerations</code> array contains taints and tolerations to determine how Kubernetes schedules the <code>Pods</code> in the role.
See <router-link to="/clusters/130_pod_scheduling">Configure Pod Scheduling</router-link> for more details.</li>
<li data-value="17">The <code>affinity</code> contains <code>Pod</code> affinity fields to determine how Kubernetes schedules the <code>Pods</code> in the role.
See <router-link to="/clusters/130_pod_scheduling">Configure Pod Scheduling</router-link> for more details.</li>
<li data-value="18">The <code>resources</code> contains configures resource limits for the Coherence containers.
See <router-link to="/clusters/140_resource_constraints">Configure Container Resource Limits</router-link> for more details.</li>
<li data-value="19">The <code>readinessProbe</code> section configures the readiness probe for the Coherence containers.
See <router-link to="/clusters/150_readiness_liveness">Readiness &amp; Liveness Probes</router-link> for more details.</li>
<li data-value="20">The <code>livenessProbe</code> section configures the liveness probe for the Coherence containers.
See <router-link to="/clusters/150_readiness_liveness">Readiness &amp; Liveness Probes</router-link> for more details.</li>
</ul>
</div>
</doc-view>
