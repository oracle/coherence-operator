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

<p>As an alternative to JMX see <router-link to="#020_manegement_over_rest.adoc" @click.native="this.scrollFix('#020_manegement_over_rest.adoc')">Management over ReST</router-link> for how to connect to a cluster via
the VisualVM plugin using ReST.</p>

<div class="admonition note">
<p class="admonition-inline">See the <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/manage/introduction-oracle-coherence-management.html">Coherence Management Documentation</a>
for more information on JMX and Management.</p>
</div>

<h3 id="_prerequisites">Prerequisites</h3>
<div class="section">
<ol style="margin-left: 15px;">
<li>
Install the Coherence Operator
<p>Ensure you have installed the Coherence Operator using the <router-link to="/install/01_installation">Install Guide</router-link>.</p>

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

<h3 id="_install_the_coherencecluster">Install the <code>CoherenceCluster</code></h3>
<div class="section">
<p>In this example a simple cluster with two roles will be created. The first role,
named <code>data</code> will be three storage enabled members. The second role named <code>management</code> will have a single replica
and will run the MBean server.</p>

<markup
lang="yaml"
title="cluster-with-jmx.yaml"
>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
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
<li data-value="1">This example uses a role named <code>data</code> as the storage enabled part of the cluster with three replicas</li>
<li data-value="2">The <code>management</code> role will be configured to run the MBean server</li>
<li data-value="3">Only one replica is typically required for the MBean server role</li>
<li data-value="4">The MBean server should be storage disabled</li>
<li data-value="5">The main class that the JVM should run should be the custom MBean server class <code>com.oracle.coherence.k8s.JmxmpServer</code></li>
<li data-value="6">Additional system properties are added to enable Coherence management
See the <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/12.2.1.4/manage/introduction-oracle-coherence-management.html">Coherence Management Documentation</a></li>
<li data-value="7">JMXMP is enabled so that a reliable JMX connection can be made to the MBean server from outside the <code>Pods</code></li>
<li data-value="8">The default port that the JMXMP server binds to is <code>9099</code> so this port is exposed as an additional port for the
<code>management</code> role</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">The default Coherence image used comes from <a id="" title="" target="_blank" href="https://container-registry.oracle.com">Oracle Container Registry</a>
so unless that image has already been pulled onto the Kubernetes nodes an <code>imagePullSecrets</code> field will be required
to pull the image.
See <router-link to="/about/04_obtain_coherence_images">Obtain Coherence Images</router-link></p>
</div>
<p>The example <code>cluster-with-jmx.yaml</code> can be installed into Kubernetes with the following command:</p>

<markup
lang="bash"

>kubectl -n sample-coherence-ns apply -f cluster-with-jmx.yaml</markup>

<p>This should install the cluster with two roles resulting in two <code>CoherenceRole</code> resources, two <code>StatefulSets</code> and four
<code>Pods</code> being created.</p>

</div>

<h3 id="_check_whether_the_mbean_server_pod_is_running">Check Whether the MBean server Pod is Running:</h3>
<div class="section">
<p>The following <code>kubectl</code> command can be used to check that the cluster is running:</p>

<markup
lang="bash"

>kubectl -n sample-coherence-ns get coherenceclusters</markup>

<p>&#8230;&#8203;which should display something like:</p>

<markup
lang="bash"

>NAME           ROLES
test-cluster   2             <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>test-cluster</code> was created and as expected has two roles.</li>
</ul>
<p>The following <code>kubectl</code> command can be used to check that the roles is running:</p>

<markup
lang="bash"

>kubectl -n sample-coherence-ns get coherenceroles</markup>

<p>&#8230;&#8203;which should display something like:</p>

<markup
lang="bash"

>NAME                      ROLE         CLUSTER        REPLICAS   READY   STATUS
test-cluster-data         data         test-cluster   3          3       Ready    <span class="conum" data-value="1" />
test-cluster-management   management   test-cluster   1          1       Ready    <span class="conum" data-value="2" /></markup>

<ul class="colist">
<li data-value="1">The <code>data</code> role was created with three replicas and three <code>Pods</code> are ready</li>
<li data-value="2">The <code>management</code> role was created with one replica and one <code>Pod</code> is ready</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">The output above may not all of the <code>Pods</code> are ready depending on how quickly the command is entered after
creating the <code>CoherenceCluster</code>, eventually all of the <code>Pods</code> should reach a ready state.</p>
</div>
<p>The following <code>kubectl</code> command can be used to list the <code>Pods</code></p>

<markup
lang="bash"

>kubectl -n sample-coherence-ns get pods</markup>

<p>&#8230;&#8203;which should display something like:</p>

<markup
lang="bash"

>NAME                                          READY   STATUS    RESTARTS   AGE
operator-coherence-operator-5d779ffc7-6pnfk   1/1     Running   0          4m33s  <span class="conum" data-value="1" />
test-cluster-data-0                           1/1     Running   0          2m39s  <span class="conum" data-value="2" />
test-cluster-data-1                           1/1     Running   0          2m39s
test-cluster-data-2                           1/1     Running   0          2m39s
test-cluster-management-0                     1/1     Running   0          2m36s  <span class="conum" data-value="3" /></markup>

<ul class="colist">
<li data-value="1">The Coherence Operator <code>Pod</code> is running in the namespace</li>
<li data-value="2">There are three pods prefixed <code>test-cluster-data-</code> that are the <code>Pods</code> for the <code>data</code> role</li>
<li data-value="3">There is one pod <code>test-cluster-management-0</code> that is the <code>Pod</code> for the <code>management</code> role</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">The output above may not all of the <code>Pods</code> are ready depending on how quickly the command is entered after
creating the <code>CoherenceCluster</code>, eventually all of the <code>Pods</code> should reach a ready state.</p>
</div>
</div>

<h3 id="_optional_add_data_to_a_cache">(Optional) Add Data to a Cache</h3>
<div class="section">
<div class="admonition note">
<p class="admonition-inline">If you do not carry out this step, then you will not see any <code>CacheMBeans</code>.</p>
</div>
<ol style="margin-left: 15px;">
<li>
The following command will run <code>kubectl</code> to exec into the first <code>data</code> role <code>Pod</code> and start an interactive
Coherence console session.
<markup
lang="bash"

>kubectl exec -it --namespace sample-coherence-ns \
    test-cluster-data-0 bash /scripts/startCoherence.sh console</markup>

</li>
<li>
At the <code>Map (?):</code> prompt, enter the command:
<markup


>cache test</markup>

<p>This will create a cache names <code>test</code> in the cache service <code>PartitionedCache</code>.</p>

</li>
<li>
Enter the following command to add 100,000 objects of size 1024 bytes, starting at index 0 and using batches of 100.
<markup
lang="bash"

>bulkput 100000 1024 0 100</markup>

</li>
<li>
When the <code>Map (?):</code> prompt returns, enter the <code>size</code> command and the console should display <code>100000</code>.

</li>
<li>
Finally type the command <code>bye</code> and press <code>&lt;enter&gt;</code> to exit the <code>console</code>.

</li>
</ol>
</div>

<h3 id="_port_forward_the_mbean_server_pod">Port Forward the MBean Server Pod:</h3>
<div class="section">
<p>The simplest way to connect from a dev machine into the management node is to just use <code>kubectl</code> to forward a local
port to the management <code>Pod</code>, which is named <code>test-cluster-management-0</code>.</p>

<markup
lang="bash"

>kubectl --namespace sample-coherence-ns port-forward \
  test-cluster-management-0 9099:9099</markup>

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
<p>Refer to the [Coherence MBean Reference](<a id="" title="" target="_blank" href="https://docs.oracle.com/middleware/12213/coherence/COHMG/oracle-coherence-mbeans-reference.htm#COHMG5442">https://docs.oracle.com/middleware/12213/coherence/COHMG/oracle-coherence-mbeans-reference.htm#COHMG5442</a>) for detailed information about Coherence MBeans.</p>

</div>

<h3 id="_clean_up">Clean Up</h3>
<div class="section">
<p>Finally to clean up the cluster run the <code>kubectl</code> command:</p>

<markup
lang="bash"

>kubectl -n sample-coherence-ns delete -f cluster-with-jmx.yaml</markup>

<p>And finally, if required, uninstall the Coherence Operator.</p>

</div>
</div>
</doc-view>
