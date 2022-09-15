<doc-view>

<h2 id="_coherence_operator_deployment_example">Coherence Operator Deployment Example</h2>
<div class="section">
<p>This example showcases how to deploy Coherence applications using the Coherence Operator.</p>

<p>This example shows how to use the Kubernetes Horizontal Pod Autoscaler to scale Coherence clusters.</p>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p><img src="./images/GitHub-Mark-32px.png" alt="GitHub Mark 32px" />
 The complete source code for this example is in the <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/tree/master/examples/021_deployment">Coherence Operator GitHub</a> repository.</p>
</p>
</div>
<p>The following scenarios are covered:</p>

<ul class="ulist">
<li>
<p><router-link to="#pre" @click.native="this.scrollFix('#pre')">Prerequisites</router-link></p>
<ul class="ulist">
<li>
<p><router-link to="#create-the-example-namespace" @click.native="this.scrollFix('#create-the-example-namespace')">Create the example namespace</router-link></p>

</li>
<li>
<p><router-link to="#clone-the-github-repository" @click.native="this.scrollFix('#clone-the-github-repository')">Clone the GitHub repository</router-link></p>

</li>
<li>
<p><router-link to="#install-operator" @click.native="this.scrollFix('#install-operator')">Install the Coherence Operator</router-link></p>

</li>
</ul>
</li>
<li>
<p><router-link to="#examples" @click.native="this.scrollFix('#examples')">Run the Examples</router-link></p>
<ul class="ulist">
<li>
<p><router-link to="#ex1" @click.native="this.scrollFix('#ex1')">Example 1 - Coherence cluster only</router-link></p>

</li>
<li>
<p><router-link to="#ex2" @click.native="this.scrollFix('#ex2')">Example 2 - Adding a Proxy tier</router-link></p>

</li>
<li>
<p><router-link to="#ex3" @click.native="this.scrollFix('#ex3')">Example 3 - Adding a User application tier</router-link></p>

</li>
<li>
<p><router-link to="#ex4" @click.native="this.scrollFix('#ex4')">Example 4 - Enabling Persistence</router-link></p>

</li>
<li>
<p><router-link to="#metrics" @click.native="this.scrollFix('#metrics')">View Cluster Metrics using Prometheus and Grafana</router-link></p>

</li>
</ul>
</li>
<li>
<p><router-link to="#cleaning-up" @click.native="this.scrollFix('#cleaning-up')">Cleaning Up</router-link></p>

</li>
</ul>
<p>After the initial installation of the Coherence cluster, the following examples
build on the previous ones by issuing a <code>kubectl apply</code> to modify
the installation adding additional tiers.</p>

<p>You can use <code>kubectl create</code> for any of the examples to install that one directly.</p>

</div>

<h2 id="pre">Prerequisites</h2>
<div class="section">
<p>Ensure you have the following software installed:</p>

<ul class="ulist">
<li>
<p>Java 11+ JDK either [OpenJDK](<a id="" title="" target="_blank" href="https://adoptopenjdk.net/">https://adoptopenjdk.net/</a>) or [Oracle JDK](<a id="" title="" target="_blank" href="https://www.oracle.com/java/technologies/javase-downloads.html">https://www.oracle.com/java/technologies/javase-downloads.html</a>)</p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://docs.docker.com/install/">Docker</a> version 17.03+.</p>

</li>
<li>
<p>Access to a Kubernetes v1.14.0+ cluster.</p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/tools/install-kubectl/">kubectl</a> version matching your Kubernetes cluster.</p>

</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">This example requires Java 11+ because it creates a Helidon web application and Helidon requires Java 11+. Coherence and running Coherence in Kubernetes only requires Java 8+.</p>
</div>
</div>

<h2 id="create-the-example-namespace">Create the example namespace</h2>
<div class="section">
<p>You need to create the namespace for the first time to run any of the examples. Create your target namespace:</p>

<markup
lang="bash"

>kubectl create namespace coherence-example

namespace/coherence-example created</markup>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>In the examples, a Kubernetes namespace called <code>coherence-example</code> is used.
If you want to change this namespace, ensure that you change any references to this namespace
to match your selected namespace when running the examples.</p>
</p>
</div>
</div>

<h2 id="clone-the-github-repository">Clone the GitHub repository</h2>
<div class="section">
<p>These examples exist in the <code>examples/021_deployment</code> directory in the
<a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator">Coherence Operator GitHub repository</a>.</p>

<p>Clone the repository:</p>

<markup
lang="bash"

>git clone https://github.com/oracle/coherence-operator

cd coherence-operator/examples/021_deployment</markup>

<p>Ensure you have Docker running and JDK 11+ build environment set and use the
following command from the deployment example directory to build the project and associated Docker image:</p>

<markup
lang="bash"

>./mvnw package jib:dockerBuild</markup>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>If you are running behind a corporate proxy and receive the following message building the Docker image:
<code>Connect to gcr.io:443 [gcr.io/172.217.212.82] failed: connect timed out</code> you must modify the build command
to add the proxy hosts and ports to be used by the <code>jib-maven-plugin</code> as shown below:</p>

<markup
lang="bash"

>mvn package jib:dockerBuild -Dhttps.proxyHost=host \
    -Dhttps.proxyPort=80 -Dhttp.proxyHost=host -Dhttp.proxyPort=80</markup>
</p>
</div>
<p>This will result in the following Docker image being created which contains the configuration and server-side
artifacts to be use by all deployments.</p>

<markup


>deployment-example:1.0.0</markup>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>If you are running against a remote Kubernetes cluster, you need to tag and
push the Docker image to your repository accessible to that cluster.
You also need to prefix the image name in the <code>yaml</code> files below.</p>
</p>
</div>
</div>

<h2 id="install-operator">Install the Coherence Operator</h2>
<div class="section">
<p>Install the Coherence Operator using your preferred method in the Operator
<a id="" title="" target="_blank" href="https://oracle.github.io/coherence-operator/docs/latest/#/docs/installation/01_installation">Installation Guide</a></p>

<p>Confirm the operator is running, for example if the operator is installed into the <code>coherence-example</code> namespace:</p>

<markup
lang="bash"

>kubectl get pods -n coherence-example

NAME                                                     READY   STATUS    RESTARTS   AGE
coherence-operator-controller-manager-74d49cd9f9-sgzjr   1/1     Running   1          27s</markup>

</div>

<h2 id="examples">Run the Examples</h2>
<div class="section">
<p>Ensure you are in the <code>examples/021_deployment</code> directory to run the following commands.</p>


<h3 id="ex1">Example 1 - Coherence cluster only</h3>
<div class="section">
<p>The first example uses the yaml file <code>src/main/yaml/example-cluster.yaml</code>, which
defines a single tier <code>storage</code> which will store cluster data.</p>

<div class="admonition note">
<p class="admonition-inline">If you have pushed your Docker image to a remote repository, ensure you update the above file to prefix the image.</p>
</div>

<h4 id="_1_install_the_coherence_cluster_storage_tier">1. Install the Coherence cluster <code>storage</code> tier</h4>
<div class="section">
<markup
lang="bash"

>kubectl -n coherence-example create -f src/main/yaml/example-cluster.yaml

coherence.coherence.oracle.com/example-cluster-storage created</markup>

</div>

<h4 id="_2_list_the_created_coherence_cluster">2. List the created Coherence cluster</h4>
<div class="section">
<markup
lang="bash"

>kubectl -n coherence-example get coherence

NAME                      CLUSTER           ROLE                      REPLICAS   READY   PHASE
example-cluster-storage   example-cluster   example-cluster-storage   3                  Created

NAME                                                         AGE
coherencerole.coherence.oracle.com/example-cluster-storage   18s</markup>

</div>

<h4 id="_3_view_the_running_pods">3. View the running pods</h4>
<div class="section">
<p>Run the following command to view the Pods:</p>

<markup
lang="bash"

>kubectl -n coherence-example get pods</markup>

<markup
lang="bash"

>NAME                                                     READY   STATUS    RESTARTS   AGE
coherence-operator-controller-manager-74d49cd9f9-sgzjr   1/1     Running   1          6m46s
example-cluster-storage-0                                0/1     Running   0          119s
example-cluster-storage-1                                1/1     Running   0          119s
example-cluster-storage-2                                0/1     Running   0          118s</markup>

</div>

<h4 id="_connect_to_the_coherence_console_inside_the_cluster_to_add_data">Connect to the Coherence Console inside the cluster to add data</h4>
<div class="section">
<p>Since we cannot yet access the cluster via Coherence*Extend, we will connect via Coherence console to add data.</p>

<markup
lang="bash"

>kubectl exec -it -n coherence-example example-cluster-storage-0 /coherence-operator/utils/runner console</markup>

<p>At the prompt type the following to create a cache called <code>test</code>:</p>

<markup
lang="bash"

>cache test</markup>

<p>Use the following to create 10,000 entries of 100 bytes:</p>

<markup
lang="bash"

>bulkput 10000 100 0 100</markup>

<p>Lastly issue the command <code>size</code> to verify the cache entry count.</p>

<p>Type <code>bye</code> to exit the console.</p>

</div>

<h4 id="_scale_the_storage_tier_to_6_members">Scale the <code>storage</code> tier to 6 members</h4>
<div class="section">
<p>To scale up the cluster the <code>kubectl scale</code> command can be used:</p>

<markup
lang="bash"

>kubectl -n coherence-example scale coherence/example-cluster-storage --replicas=6</markup>

<p>Use the following to verify all 6 nodes are Running and READY before continuing.</p>

<markup
lang="bash"

>kubectl -n coherence-example get pods</markup>

<markup
lang="bash"

>NAME                                                     READY   STATUS    RESTARTS   AGE
coherence-operator-controller-manager-74d49cd9f9-sgzjr   1/1     Running   1          53m
example-cluster-storage-0                                1/1     Running   0          49m
example-cluster-storage-1                                1/1     Running   0          49m
example-cluster-storage-2                                1/1     Running   0          49m
example-cluster-storage-3                                1/1     Running   0          54s
example-cluster-storage-4                                1/1     Running   0          54s
example-cluster-storage-5                                1/1     Running   0          54s</markup>

</div>

<h4 id="_confirm_the_cache_count">Confirm the cache count</h4>
<div class="section">
<p>Re-run step 3 above and just use the <code>cache test</code> and <code>size</code> commands to confirm the number of entries is still 10,000.</p>

<p>This confirms that the scale-out was done in a <code>safe</code> manner ensuring no data loss.</p>

</div>
</div>

<h3 id="_scale_the_storage_tier_back_to_3_members">Scale the <code>storage</code> tier back to 3 members</h3>
<div class="section">
<p>To scale back doewn to three members run the following command:</p>

<markup
lang="bash"

>kubectl -n coherence-example scale coherence/example-cluster-storage --replicas=3</markup>

<p>By using the following, you will see that the number of members will gradually scale back to
3 during which the is done in a <code>safe</code> manner ensuring no data loss.</p>

<markup
lang="bash"

>kubectl -n coherence-example get pods</markup>

<markup
lang="bash"

>NAME                        READY   STATUS        RESTARTS   AGE
example-cluster-storage-0   1/1     Running       0          19m
example-cluster-storage-1   1/1     Running       0          19m
example-cluster-storage-2   1/1     Running       0          19m
example-cluster-storage-3   1/1     Running       0          3m41s
example-cluster-storage-4   0/1     Terminating   0          3m41s</markup>

</div>

<h3 id="ex2">Example 2 - Adding a Proxy tier</h3>
<div class="section">
<p>The second example uses the yaml file <code>src/main/yaml/example-cluster-proxy.yaml</code>, which
adds a proxy server <code>example-cluster-proxy</code> to allow for Coherence*Extend connections via a Proxy server.</p>

<p>The additional yaml added below shows:</p>

<ul class="ulist">
<li>
<p>A port called <code>proxy</code> being exposed on 20000</p>

</li>
<li>
<p>The tier being set as storage-disabled</p>

</li>
<li>
<p>A different cache config being used which will start a Proxy Server. See [here](src/main/resources/proxy-cache-config.xml) for details</p>

</li>
</ul>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: example-cluster-proxy
spec:
  cluster: example-cluster
  jvm:
    memory:
      heapSize: 512m
  ports:
    - name: metrics
      port: 9612
      serviceMonitor:
        enabled: true
    - name: proxy
      port: 20000
  coherence:
    cacheConfig: proxy-cache-config.xml
    storageEnabled: false
    metrics:
      enabled: true
  image: deployment-example:1.0.0
  imagePullPolicy: Always
  replicas: 1</markup>


<h4 id="_install_the_proxy_tier">Install the <code>proxy</code> tier</h4>
<div class="section">
<markup
lang="bash"

>  kubectl -n coherence-example apply -f src/main/yaml/example-cluster-proxy.yaml

  kubectl get coherence -n coherence-example

  NAME                      CLUSTER           ROLE                      REPLICAS   READY   PHASE
  example-cluster-proxy     example-cluster   example-cluster-proxy     1          1       Ready
  example-cluster-storage   example-cluster   example-cluster-storage   3          3       Ready</markup>

</div>

<h4 id="_view_the_running_pods">View the running pods</h4>
<div class="section">
<markup
lang="bash"

>kubectl -n coherence-example get pods

NAME                                  READY   STATUS    RESTARTS   AGE
coherence-operator-578497bb5b-w89kt   1/1     Running   0          68m
example-cluster-proxy-0               1/1     Running   0          2m41s
example-cluster-storage-0             1/1     Running   0          29m
example-cluster-storage-1             1/1     Running   0          29m
example-cluster-storage-2             1/1     Running   0          2m43s</markup>

<p>Ensure the <code>example-cluster-proxy-0</code> pod is Running and READY before continuing.</p>

</div>

<h4 id="_port_forward_the_proxy_port">Port forward the proxy port</h4>
<div class="section">
<pre>In a separate terminal, run the following:</pre>
<markup
lang="bash"

>    kubectl port-forward -n coherence-example example-cluster-proxy-0 20000:20000</markup>

</div>

<h4 id="_connect_via_cohql_and_add_data">Connect via CohQL and add data</h4>
<div class="section">
<p>In a separate terminal, change to the <code>examples/021_deployments</code> directory and run the following to
start Coherence Query Language (CohQL):</p>

<markup
lang="bash"

>    mvn exec:java

    Coherence Command Line Tool

    CohQL&gt;</markup>

<p>Run the following <code>CohQL</code> commands to view and insert data into the cluster.</p>

<markup


>CohQL&gt; select count() from 'test';

Results
10000

CohQL&gt; insert into 'test' key('key-1') value('value-1');

CohQL&gt; select key(), value() from 'test' where key() = 'key-1';
Results
["key-1", "value-1"]

CohQL&gt; select count() from 'test';
Results
10001

CohQL&gt; quit</markup>

<p>The above results will show that you can see the data previously inserted and
can add new data into the cluster using Coherence*Extend.</p>

</div>
</div>

<h3 id="ex3">Example 3 - Adding a User application tier</h3>
<div class="section">
<p>The third example uses the yaml file <code>src/main/yaml/example-cluster-app.yaml</code>, which
adds a new tier <code>rest</code>. This tier defines a user application which uses <a id="" title="" target="_blank" href="https://helidon.io/">Helidon</a> to create a <code>/query</code> endpoint allowing the user to send CohQL commands via this endpoint.</p>

<p>The additional yaml added below shows:</p>

<ul class="ulist">
<li>
<p>A port called <code>http</code> being exposed on 8080 for the application</p>

</li>
<li>
<p>The tier being set as storage-disabled</p>

</li>
<li>
<p>Using the storage-cache-config.xml but as storage-disabled</p>

</li>
<li>
<p>An alternate main class to run - <code>com.oracle.coherence.examples.Main</code></p>

</li>
</ul>
<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: example-cluster-rest
spec:
  cluster: example-cluster
  jvm:
    memory:
      heapSize: 512m
  ports:
    - name: metrics
      port: 9612
      serviceMonitor:
        enabled: true
    - name: http
      port: 8080
  coherence:
    cacheConfig: storage-cache-config.xml
    storageEnabled: false
    metrics:
      enabled: true
  image: deployment-example:1.0.0
  imagePullPolicy: Always
  application:
    main: com.oracle.coherence.examples.Main</markup>


<h4 id="_install_the_rest_tier">Install the <code>rest</code> tier</h4>
<div class="section">
<p>Install the yaml with the following command:</p>

<markup
lang="bash"

>kubectl -n coherence-example apply -f src/main/yaml/example-cluster-app.yaml

kubectl get coherence -n coherence-example

NAME                      CLUSTER           ROLE                      REPLICAS   READY   PHASE
example-cluster-proxy     example-cluster   example-cluster-proxy     1          1       Ready
example-cluster-rest      example-cluster   example-cluster-rest      1          1       Ready
example-cluster-storage   example-cluster   example-cluster-storage   3          3       Ready</markup>

</div>

<h4 id="_view_the_running_pods_2">View the running pods</h4>
<div class="section">
<markup
lang="bash"

>kubectl -n coherence-example get pods

NAME                              READY   STATUS    RESTARTS   AGE
coherence-operator-578497bb5b-w89kt   1/1     Running   0          90m
example-cluster-proxy-0               1/1     Running   0          3m57s
example-cluster-rest-0                1/1     Running   0          3m57s
example-cluster-storage-0             1/1     Running   0          3m59s
example-cluster-storage-1             1/1     Running   0          3m58s
example-cluster-storage-2             1/1     Running   0          3m58s</markup>

</div>

<h4 id="_port_forward_the_application_port">Port forward the application port</h4>
<div class="section">
<p>In a separate terminal, run the following:</p>

<markup
lang="bash"

>kubectl port-forward -n coherence-example example-cluster-rest-0 8080:8080</markup>

</div>

<h4 id="_access_the_custom_query_endpoint">Access the custom <code>/query</code> endpoint</h4>
<div class="section">
<p>Use the various <code>CohQL</code> commands via the <code>/query</code> endpoint to access, and mutate data in the Coherence cluster.</p>

<markup
lang="bash"

>curl -i -w '\n' -X PUT http://127.0.0.1:8080/query -d '{"query":"create cache foo"}'</markup>

<markup
lang="bash"

>HTTP/1.1 200 OK
Date: Fri, 19 Jun 2020 06:29:40 GMT
transfer-encoding: chunked
connection: keep-alive</markup>

<markup
lang="bash"

>curl -i -w '\n' -X PUT http://127.0.0.1:8080/query -d '{"query":"insert into foo key(\"foo\") value(\"bar\")"}'</markup>

<markup
lang="bash"

>HTTP/1.1 200 OK
Date: Fri, 19 Jun 2020 06:29:44 GMT
transfer-encoding: chunked
connection: keep-alive</markup>

<markup
lang="bash"

>curl -i -w '\n' -X PUT http://127.0.0.1:8080/query -d '{"query":"select key(),value() from foo"}'</markup>

<markup
lang="bash"

>HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 19 Jun 2020 06:29:55 GMT
transfer-encoding: chunked
connection: keep-alive

{"result":"{foo=[foo, bar]}"}</markup>

<markup
lang="bash"

>curl -i -w '\n' -X PUT http://127.0.0.1:8080/query -d '{"query":"create cache test"}'</markup>

<markup
lang="bash"

>HTTP/1.1 200 OK
Date: Fri, 19 Jun 2020 06:30:00 GMT
transfer-encoding: chunked
connection: keep-alive</markup>

<markup
lang="bash"

>curl -i -w '\n' -X PUT http://127.0.0.1:8080/query -d '{"query":"select count() from test"}'</markup>

<markup
lang="bash"

>HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 19 Jun 2020 06:30:20 GMT
transfer-encoding: chunked
connection: keep-alive

{"result":"10001"}</markup>

</div>
</div>

<h3 id="ex4">Example 4 - Enabling Persistence</h3>
<div class="section">
<p>The fourth example uses the yaml file <code>src/main/yaml/example-cluster-persistence.yaml</code>, which
enabled Active Persistence for the <code>storage</code> tier by adding a <code>persistence:</code> element.</p>

<p>The additional yaml added to the storage tier below shows:</p>

<ul class="ulist">
<li>
<p>Active Persistence being enabled via <code>persistence.enabled=true</code></p>

</li>
<li>
<p>Various Persistence Volume Claim (PVC) values being set under <code>persistentVolumeClaim</code></p>

</li>
</ul>
<markup
lang="yaml"

>  coherence:
    cacheConfig: storage-cache-config.xml
    metrics:
      enabled: true
    persistence:
      enabled: true
      persistentVolumeClaim:
        accessModes:
          - ReadWriteOnce
        resources:
          requests:
            storage: 1Gi</markup>

<p>NOTE:By default, when you enable Coherence Persistence, the required infrastructure in terms of persistent volumes (PV) and persistent volume claims (PVC) is set up automatically. Also, the persistence-mode is set to <code>active</code>. This allows the Coherence cluster to be restarted, and the data to be retained.</p>


<h4 id="_delete_the_existing_deployment">Delete the existing deployment</h4>
<div class="section">
<p>We must first delete the existing deployment as we need to redeploy to enable Active Persistence.</p>

<markup
lang="bash"

>kubectl -n coherence-example delete -f src/main/yaml/example-cluster-app.yaml</markup>

<p>Ensure all the pods have terminated before you continue.</p>

</div>

<h4 id="_install_the_cluster_with_persistence_enabled">Install the cluster with Persistence enabled</h4>
<div class="section">
<markup
lang="bash"

>kubectl -n coherence-example create -f src/main/yaml/example-cluster-persistence.yaml</markup>

</div>

<h4 id="_view_the_running_pods_and_pvcs">View the running pods and PVC&#8217;s</h4>
<div class="section">
<markup
lang="bash"

>kubectl -n coherence-example get pods</markup>

<markup
lang="bash"

>NAME                            READY   STATUS    RESTARTS   AGE
example-cluster-rest-0          1/1     Running   0          5s
example-cluster-proxy-0         1/1     Running   0          5m1s
example-cluster-storage-0       1/1     Running   0          5m3s
example-cluster-storage-1       1/1     Running   0          5m3s
example-cluster-storage-2       1/1     Running   0          5m3s</markup>

<p>Check the Persistent Volumes and PVC are automatically created.</p>

<markup
lang="bash"

>kubectl get pvc -n coherence-example</markup>

<markup
lang="bash"

>NAME                                           STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
persistence-volume-example-cluster-storage-0   Bound    pvc-15b46996-eb35-11e9-9b4b-025000000001   1Gi        RWO            hostpath       55s
persistence-volume-example-cluster-storage-1   Bound    pvc-15bd99e9-eb35-11e9-9b4b-025000000001   1Gi        RWO            hostpath       55s
persistence-volume-example-cluster-storage-2   Bound    pvc-15e55b6b-eb35-11e9-9b4b-025000000001   1Gi        RWO            hostpath       55s</markup>

<p>Wait until all  nodes are Running and READY before continuing.</p>

</div>

<h4 id="_check_active_persistence_is_enabled">Check Active Persistence is enabled</h4>
<div class="section">
<p>Use the following to view the logs of the <code>example-cluster-storage-0</code> pod and validate that Active Persistence is enabled.</p>

<markup
lang="bash"

>kubectl logs example-cluster-storage-0 -c coherence -n coherence-example | grep 'Created persistent'</markup>

<markup
lang="bash"

>...
019-10-10 04:52:00.179/77.023 Oracle Coherence GE 12.2.1.4.0 &lt;Info&gt; (thread=DistributedCache:PartitionedCache, member=4): Created persistent store /persistence/active/example-cluster/PartitionedCache/126-2-16db40199bc-4
2019-10-10 04:52:00.247/77.091 Oracle Coherence GE 12.2.1.4.0 &lt;Info&gt; (thread=DistributedCache:PartitionedCache, member=4): Created persistent store /persistence/active/example-cluster/PartitionedCache/127-2-16db40199bc-4
...</markup>

<p>If you see output similar to above then Active Persistence is enabled.</p>

</div>

<h4 id="_connect_to_the_coherence_console_to_add_data">Connect to the Coherence Console to add data</h4>
<div class="section">
<markup
lang="bash"

>kubectl exec -it -n coherence-example example-cluster-storage-0 /coherence-operator/utils/runner console</markup>

<p>At the prompt type the following to create a cache called <code>test</code>:</p>

<markup
lang="bash"

>cache test</markup>

<p>Use the following to create 10,000 entries of 100 bytes:</p>

<markup
lang="bash"

>bulkput 10000 100 0 100</markup>

<p>Lastly issue the command <code>size</code> to verify the cache entry count.</p>

<p>Type <code>bye</code> to exit the console.</p>

</div>

<h4 id="_delete_the_cluster">Delete the cluster</h4>
<div class="section">
<div class="admonition note">
<p class="admonition-inline">This will not delete the PVC&#8217;s.</p>
</div>
<markup
lang="bash"

>kubectl -n coherence-example delete -f src/main/yaml/example-cluster-persistence.yaml</markup>

<p>Use <code>kubectl get pods -n coherence-example</code> to confirm the pods have terminated.</p>

</div>

<h4 id="_confirm_the_pvcs_are_still_present">Confirm the PVC&#8217;s are still present</h4>
<div class="section">
<markup
lang="bash"

>kubectl get pvc -n coherence-example</markup>

<markup
lang="bash"

>NAME                                           STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
persistence-volume-example-cluster-storage-0   Bound    pvc-730f86fe-eb19-11e9-9b4b-025000000001   1Gi        RWO            hostpath       116s
persistence-volume-example-cluster-storage-1   Bound    pvc-73191751-eb19-11e9-9b4b-025000000001   1Gi        RWO            hostpath       116s
persistence-volume-example-cluster-storage-2   Bound    pvc-73230889-eb19-11e9-9b4b-025000000001   1Gi        RWO            hostpath       116s</markup>

</div>

<h4 id="_re_install_the_cluster">Re-install the cluster</h4>
<div class="section">
<markup
lang="bash"

>kubectl -n coherence-example create -f src/main/yaml/example-cluster-persistence.yaml</markup>

</div>

<h4 id="_follow_the_logs_for_persistence_messages">Follow the logs for Persistence messages</h4>
<div class="section">
<markup
lang="bash"

>kubectl logs example-cluster-storage-0 -c coherence -n coherence-example -f</markup>

<p>You should see a message regarding recovering partitions, similar to the following:</p>

<markup
lang="bash"

>2019-10-10 05:00:14.255/32.206 Oracle Coherence GE 12.2.1.4.0 &lt;D5&gt; (thread=DistributedCache:PartitionedCache, member=1): Recovering 86 partitions
...
2019-10-10 05:00:17.417/35.368 Oracle Coherence GE 12.2.1.4.0 &lt;Info&gt; (thread=DistributedCache:PartitionedCache, member=1): Created persistent store /persistence/active/example-cluster/PartitionedCache/50-3-16db409d035-1 from SafeBerkeleyDBStore(50-2-16db40199bc-4, /persistence/active/example-cluster/PartitionedCache/50-2-16db40199bc-4)
...</markup>

<p>Finally, you should see the following indicating active recovery has completed.</p>

<markup
lang="bash"

>2019-10-10 08:18:04.870/59.565 Oracle Coherence GE 12.2.1.4.0 &lt;Info&gt; (thread=DistributedCache:PartitionedCache, member=1):
   Recovered PartitionSet{172..256} from active persistent store</markup>

</div>

<h4 id="_confirm_the_data_has_been_recovered">Confirm the data has been recovered</h4>
<div class="section">
<markup
lang="bash"

>kubectl exec -it -n coherence-example example-cluster-storage-0 /coherence-operator/utils/runner console</markup>

<p>At the prompt type the following to create a cache called <code>test</code>:</p>

<markup
lang="bash"

>cache test</markup>

<p>Lastly issue the command <code>size</code> to verify the cache entry count is 10,000 meaning the data has been recovered.</p>

<p>Type <code>bye</code> to exit the console.</p>

</div>
</div>

<h3 id="metrics">View Cluster Metrics Using Prometheus and Grafana</h3>
<div class="section">
<p>If you wish to view metrics via Grafana, you must carry out the following steps <strong>before</strong> you
install any of the examples above.</p>


<h4 id="_install_prometheus_operator">Install Prometheus Operator</h4>
<div class="section">
<p>Install the Prometheus Operator, as documented in the Prometheus Operator <a id="" title="" target="_blank" href="https://prometheus-operator.dev/docs/prologue/quick-start/">Quick Start</a> page. Prometheus can then be accessed as documented in the
<a id="" title="" target="_blank" href="https://prometheus-operator.dev/docs/prologue/quick-start/#access-prometheus">Access Prometheus section of the Quick Start</a> page.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p><strong>Using RBAC</strong></p>

<p>If installing Prometheus into RBAC enabled k8s clusters, you may need to create the required RBAC resources
as described in the <a id="" title="" target="_blank" href="https://prometheus-operator.dev/docs/operator/rbac/">Prometheus RBAC</a> documentation.
The Coherence Operator contains an example that works with the out-of-the-box Prometheus Operator install
that we use for testing <a id="" title="" target="_blank" href="https://raw.githubusercontent.com/oracle/coherence-operator/master/hack/prometheus-rbac.yaml">prometheus-rbac.yaml</a>
This yaml creates a <code>ClusterRole</code> with the required permissions and a <code>ClusterRoleBinding</code> that binds the role to the
<code>prometheus-k8s</code> service account (which is the name of the account created, and used by the Prometheus Operator).
This yaml file can be installed into k8s before installing the Prometheus Operator.</p>
</p>
</div>
</div>

<h4 id="_access_grafana">Access Grafana</h4>
<div class="section">
<p>The Prometheus Operator also installs Grafana. Grafana can be accessed as documented in the
<a id="" title="" target="_blank" href="https://prometheus-operator.dev/docs/prologue/quick-start/#access-grafana">Access Grafana section of the Quick Start</a> page.
Note that the default credentials are specified in that section of the documentation.</p>

</div>

<h4 id="_import_the_grafana_dashboards">Import the Grafana Dashboards</h4>
<div class="section">
<p>To import the Coherence Grafana dashboards follow the instructions in the Operator documentation section
<router-link to="#metrics/030_importing.adoc" @click.native="this.scrollFix('#metrics/030_importing.adoc')">Importing Grafana Dashboards</router-link>.</p>

<p>After importing the dashboards into Grafana and with the port-forward still running the Coherence dashboards can be
accessed at <a id="" title="" target="_blank" href="http://localhost:3000/d/coh-main/coherence-dashboard-main">localhost:3000/d/coh-main/coherence-dashboard-main</a></p>

</div>

<h4 id="_troubleshooting">Troubleshooting</h4>
<div class="section">
<ul class="ulist">
<li>
<p>It may take up to 5 minutes for data to start appearing in Grafana.</p>

</li>
<li>
<p>If you are not seeing data after 5 minutes, access the Prometheus endpoint as described above.
Ensure that the endpoints named <code>coherence-example/example-cluster-storage-metrics/0 (3/3 up)</code> are up.
If the endpoints are not up then wait 60 seconds and refresh the browser.</p>

</li>
<li>
<p>If you do not see any values in the <code>Cluster Name</code> dropdown in Grafana, ensure the endpoints are up as  described above and click on <code>Manage Alerts</code> and then <code>Back to Main Dashboard</code>. This will re-query the data and load the list of clusters.</p>

</li>
</ul>
</div>
</div>

<h3 id="cleaning-up">Cleaning Up</h3>
<div class="section">

<h4 id="_delete_the_cluster_2">Delete the cluster</h4>
<div class="section">
<markup
lang="bash"

>kubectl -n coherence-example delete -f src/main/yaml/example-cluster-persistence.yaml</markup>

</div>

<h4 id="_delete_the_pvcs">Delete the PVC&#8217;s</h4>
<div class="section">
<p>Ensure all the pods have all terminated before you delete the PVC&#8217;s.</p>

<markup
lang="bash"

>kubectl get pvc -n coherence-example | sed 1d | awk '{print $1}' | xargs kubectl delete pvc -n coherence-example</markup>

</div>

<h4 id="_remove_the_coherence_operator">Remove the Coherence Operator</h4>
<div class="section">
<p>Uninstall the Coherence operator using the undeploy commands for whichever method you chose to install it.</p>

</div>

<h4 id="_delete_prometheus_operator">Delete Prometheus Operator</h4>
<div class="section">
<p>Uninstall the Prometheus Operator as documented in the
<a id="" title="" target="_blank" href="https://prometheus-operator.dev/docs/prologue/quick-start/#remove-kube-prometheus">Remove kube-prometheus section of the Quick Start</a> page.</p>

</div>
</div>
</div>
</doc-view>
