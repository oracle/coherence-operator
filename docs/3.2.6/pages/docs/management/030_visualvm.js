<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Using VisualVM</dt>
<dd slot="desc"><p><a id="" title="" target="_blank" href="https://visualvm.github.io/">VisualVM</a> is a visual tool integrating commandline JDK tools and lightweight profiling
capabilities, designed for both development and production time use.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_access_a_coherence_cluster_via_visualvm">Access A Coherence Cluster via VisualVM</h2>
<div class="section">
<p>Coherence management is implemented using Java Management Extensions (JMX). JMX is a Java standard
for managing and monitoring Java applications and services. VisualVM and other JMX tools can be used to
manage and monitor Coherence Clusters via JMX.</p>

<p>The default transport used by JMX is RMI but RMI can be difficult to set-up reliably in containers and Kubernetes so
that it can be accessed externally due to its use of multiple TCP ports that are difficult to configure and it does
not work well with the NAT&#8217;ed type of networking typically found in these environments. JMXMP on the other hand is an
alternative to RMI that does work well in containers and only requires a single TCP port.</p>

<p>This example shows how to connect to a cluster via JMX over JMXMP.</p>

<p>As an alternative to JMX see <router-link to="/docs/management/020_management_over_rest">Management over REST</router-link>
for how to connect to a cluster via the VisualVM plugin using REST.</p>

<div class="admonition note">
<p class="admonition-inline">See the <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/manage/introduction-oracle-coherence-management.html">Coherence Management Documentation</a>
for more information on JMX and Management.</p>
</div>

<h3 id="_prerequisites">Prerequisites</h3>
<div class="section">
<ol style="margin-left: 15px;">
<li>
Install the Coherence Operator
<p>Ensure you have installed the Coherence Operator using the <router-link to="/docs/installation/01_installation">Install Guide</router-link>.</p>

</li>
<li>
Download the JMXMP connector JAR
<p>The JMX endpoint does not use RMI, instead it uses JMXMP. This requires an additional JAR on the classpath
of the Java JMX client (VisualVM and JConsole). You can use curl to download the required JAR.</p>

<markup
lang="bash"

>curl http://central.maven.org/maven2/org/glassfish/external/opendmk_jmxremote_optional_jar/1.0-b01-ea/opendmk_jmxremote_optional_jar-1.0-b01-ea.jar \
    -o opendmk_jmxremote_optional_jar-1.0-b01-ea.jar</markup>

<p>This jar can also be downloaded as a Maven dependency if you are connecting through a Maven project.</p>

<markup
lang="xml"

>&lt;dependency&gt;
  &lt;groupId&gt;org.glassfish.external&lt;/groupId&gt;
  &lt;artifactId&gt;opendmk_jmxremote_optional_jar&lt;/artifactId&gt;
  &lt;version&gt;1.0-b01-ea&lt;/version&gt;
&lt;/dependency&gt;</markup>

</li>
</ol>
</div>

<h3 id="_install_a_jmx_enabled_coherence_cluster">Install a JMX Enabled Coherence Cluster</h3>
<div class="section">
<p>In this example a simple minimal cluster will be created running the MBean server.</p>

<markup
lang="yaml"
title="cluster-with-jmx.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test-cluster
spec:
  jvm:
    args:                                                  <span class="conum" data-value="1" />
      - -Dcoherence.management=all
      - -Dcoherence.management.remote=true
      - -Dcom.sun.management.jmxremote.ssl=false
      - -Dcom.sun.management.jmxremote.authenticate=false
    jmxmp:
      enabled: true                                        <span class="conum" data-value="2" />
      port: 9099
  ports:
    - name: jmx                                            <span class="conum" data-value="3" />
      port: 9099</markup>

<ul class="colist">
<li data-value="1">Additional system properties are added to enable Coherence management
See the <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/manage/introduction-oracle-coherence-management.html">Coherence Management Documentation</a></li>
<li data-value="2">JMXMP is enabled on port <code>9099</code> so that a reliable JMX connection can be made to the MBean server from outside the <code>Pods</code></li>
<li data-value="3">The JMXMP port is exposed as an additional port</li>
</ul>
<p>The example <code>cluster-with-jmx.yaml</code> can be installed into Kubernetes with the following command:</p>

<markup
lang="bash"

>kubectl -n coherence-test apply -f cluster-with-jmx.yaml</markup>

<p>This should install the cluster into the namespace <code>coherence-test</code> with a default replica count of three, resulting in
a <code>StatefulSet</code> with three <code>Pods</code>.</p>

</div>

<h3 id="_port_forward_the_mbean_server_pod">Port Forward the MBean Server Pod:</h3>
<div class="section">
<p>After installing the basic <code>cluster-with-jmx.yaml</code> from the example above there would be a three member
Coherence cluster installed into Kubernetes.</p>

<p>The <code>kubectl</code> CLI can be used to list <code>Pods</code> for the cluster:</p>

<markup
lang="bash"

>kubectl -n coherence-test get pod -l coherenceCluster=test-cluster

NAME             READY   STATUS    RESTARTS   AGE
test-cluster-0   1/1     Running   0          36s
test-cluster-1   1/1     Running   0          36s
test-cluster-2   1/1     Running   0          36s</markup>

<p>In a test or development environment the simplest way to reach an exposed port is to use the <code>kubectl port-forward</code> command.
For example to connect to the first <code>Pod</code> in the deployment:</p>

<markup
lang="bash"

>kubectl -n coherence-test port-forward test-cluster-0 9099:9099

Forwarding from [::1]:9099 -&gt; 9099
Forwarding from 127.0.0.1:9099 -&gt; 9099</markup>

<p>JMX can now be access using the URL <code>service:jmx:jmxmp://127.0.0.1:9099</code></p>

</div>

<h3 id="_access_mbeans_through_jconsole">Access MBeans Through JConsole</h3>
<div class="section">
<ol style="margin-left: 15px;">
<li>
Run JConsole with the JMXMP connector on the classpath:
<markup
lang="bash"

>jconsole -J-Djava.class.path="$JAVA_HOME/lib/jconsole.jar:$JAVA_HOME/lib/tools.jar:opendmk_jmxremote_optional_jar-1.0-b01-ea.jar" service:jmx:jmxmp://127.0.0.1:9099</markup>

</li>
<li>
In the console UI, select the <code>MBeans</code> tab and then <code>Coherence Cluster</code> attributes.
You should see the Coherence MBeans as shown below:
<p><img src="./images/jconsole.png" alt="VisualVM"width="513" />
</p>

</li>
</ol>
</div>

<h3 id="_access_mbeans_through_visualvm">Access MBeans Through VisualVM</h3>
<div class="section">
<ol style="margin-left: 15px;">
<li>
Ensure you run VisualVM with the JMXMP connector on the classpath:
<markup
lang="bash"

>jvisualvm -cp "$JAVA_HOME/lib/tools.jar:opendmk_jmxremote_optional_jar-1.0-b01-ea.jar"</markup>

<div class="admonition note">
<p class="admonition-inline">If you have downloaded VisualVM separately (as VisualVM has not been part of the JDK from Java 9 onwards),
then the executable is <code>visualvm</code> (or on MacOS it is <code>/Applications/VisualVM.app/Contents/MacOS/visualvm</code>).</p>
</div>
</li>
<li>
From the VisualVM menu select <code>File</code> / <code>Add JMX Connection</code>

</li>
<li>
Enter <code>service:jmx:jmxmp://127.0.0.1:9099</code> for the <code>Connection</code> value and click <code>OK</code>.
<p>A JMX connection should be added under the <code>Local</code> section of the left hand panel.</p>

</li>
<li>
Double-click the new local connection to connect to the management <code>Pod</code>.
You can see the <code>Coherence</code> MBeans under the <code>MBeans</code> tab.
If you have installed the Coherence VisualVM plugin, you can also see a <code>Coherence</code> tab.
<p><img src="./images/jvisualvm.png" alt="VisualVM"width="735" />
</p>

</li>
</ol>
<p>Refer to the <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/manage/oracle-coherence-mbeans-reference.html#GUID-5E57FA4D-9CF8-4069-A8FD-B50E4FAB2687">Coherence MBean Reference</a>
for detailed information about Coherence MBeans.</p>

</div>
</div>
</doc-view>
