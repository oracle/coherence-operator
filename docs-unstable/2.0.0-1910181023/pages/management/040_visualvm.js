<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Using VisualVM</dt>
<dd slot="desc"><p><a id="" title="" target="_blank" href="https://visualvm.github.io/">VisualVM</a> is a visual tool integrating commandline JDK tools and lightweight profiling capabilities.
Designed for both development and production time use.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_access_the_coherence_cluster_via_visualvm">Access the Coherence cluster via VisualVM</h2>
<div class="section">
<p>Coherence management is implemented using Java Management Extensions (JMX). JMX is a Java standard
for managing and monitoring Java applications and services. VisualVM and other JMX tools can be used to
manage and monitor Coherence Clusters via JMX.</p>

<p>This example shows how to connect to a cluster via VisualVM over JMXMP.</p>

<p>Please see <router-link to="#020_manegement_over_rest.adoc" @click.native="this.scrollFix('#020_manegement_over_rest.adoc')">Management over ReST</router-link> for how to connect
to a cluster via the VisualVM plugin using ReST.</p>

<div class="admonition note">
<p class="admonition-inline">See the <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/manage/introduction-oracle-coherence-management.html">Coherence Management Documentation</a>
for more information on JMX and Management.</p>
</div>

<h3 id="_create_a_coherencecluster_with_an_mbean_server_role">Create a CoherenceCluster With an MBean Server Role</h3>
<div class="section">
<p>The following <code>.yaml</code> file will create a <code>CoherenceCluster</code> with an additional <code>management</code> role that uses JMXMP as the
transport for JMX.</p>

<markup
lang="yaml"
title="cluster-with-jmx.yaml"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  readinessProbe:
    initialDelaySeconds: 10
    periodSeconds: 5
  roles:
    - role: data                                                  <span class="conum" data-value="1" />
        replicas: 3
    - role: management                                            <span class="conum" data-value="2" />
      replicas: 1                                                 <span class="conum" data-value="3" />
      coherence:
        storageEnabled: false                                     <span class="conum" data-value="4" />
      application:
        main: com.oracle.coherence.k8s.JmxmpServer                <span class="conum" data-value="5" />
      jvm:
        args:
          - -Dcoherence.distributed.localstorage=false            <span class="conum" data-value="6" />
          - -Dcoherence.management=all
          - -Dcoherence.management.remote=true
          - -Dcom.sun.management.jmxremote.ssl=false
          - -Dcom.sun.management.jmxremote.authenticate=false
        jmxmp:
          enabled: true                                           <span class="conum" data-value="7" />
      ports:
        - name: jmx                                               <span class="conum" data-value="8" />
          port: 9099</markup>

<ul class="colist">
<li data-value="1">This example uses a role named <code>data</code> as the storage enabled part of the cluster</li>
<li data-value="2">The <code>management</code> role will be configured to run the MBean server</li>
<li data-value="3">Only one replica is typically required for the MBean server role</li>
<li data-value="4">The MBean server should be storage disabled</li>
<li data-value="5">The main class that the JVM should run should be the custom MBean server class <code>com.oracle.coherence.k8s.JmxmpServer</code></li>
<li data-value="6">Additional system properties are added to enable Coherence management
See the <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/manage/introduction-oracle-coherence-management.html">Coherence Management Documentation</a></li>
<li data-value="7">JMXMP is enabled so that a reliable connection can be made to the MBean server from outside the <code>Pods</code></li>
<li data-value="8">The default port that the JMXMP server binds to is <code>9099</code> so this port is exposed as an additional port for the
<code>management</code> role</li>
</ul>
<p>Once the cluster is running a JMX connection can be made to the URL <code>service:jmx:jmxmp://&lt;host-name&gt;:9099</code> where the
 <code>&lt;host-name&gt;</code> to use will depend on how the container port is exposed.
If using <code>kubectl</code> port forwarding to expose
port <code>9099</code> on the management <code>Pod</code> then the URL would be <code>service:jmx:jmxmp://127.0.0.1:9099</code>.</p>

</div>
</div>
</doc-view>
