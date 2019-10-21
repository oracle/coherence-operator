<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Rolling Upgrades</dt>
<dd slot="desc"><p>The Coherence Operator facilitates safe rolling upgrade of either a application image or Coherence
image.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_rolling_upgrades">Rolling Upgrades</h2>
<div class="section">
<p>As described in <router-link to="/app-deployment/020_packaging">Packaging Applications</router-link> it is
usual to create sidecar Docker image which provides the application classes to Kubernetes.
The docker image is tagged with a version number and this version number is used by Kubernetes
to enable safe rolling upgrades.</p>

<p>The safe rolling upgrade feature allows you to instruct Kubernetes, through the operator,
to replace the currently installed version of your application classes with a different one.
Kubernetes does not verify whether the classes are new or old. The operator also ensures
that the replacement is done without data loss or interruption of service.</p>

<p>This is achieved simply with the Coherence Operator by using the <code>kubectl apply</code> command to against
an existing cluster to change the attached docker image.</p>

<p>This example shows how to issue a rolling upgrade to upgrade a cluster application image from <code>v1.0.0</code> to <code>v2.0.0</code> which introduces a second cache service while preserving the data in the first.</p>

<ul class="ulist">
<li>
<p>Version 1 - hr-* cache mapping maps to <code>HRPartitionedCache</code> service</p>

</li>
<li>
<p>Version 2 - additional fin-* cache mapping maps to <code>FINPartitionedCache</code> service.</p>

</li>
</ul>

<h3 id="_1_prerequisites">1. Prerequisites</h3>
<div class="section">
<ol style="margin-left: 15px;">
<li>
Install the Coherence Operator

</li>
<li>
Create any secrets required to pull Docker images

</li>
<li>
Create a new <code>working directory</code> and change to that directory

</li>
</ol>
</div>

<h3 id="_2_create_your_directory_structure">2. Create your directory structure</h3>
<div class="section">
<markup
lang="bash"

>mkdir -p files-v1/lib files-v1/conf files-v2/lib files-v2/conf</markup>

</div>

<h3 id="_3_create_the_dockerfiles">3. Create the Dockerfiles</h3>
<div class="section">
<p>In your working directory directory create a file called <code>Dockerfile-v1</code> with the following contents:</p>

<markup
lang="dockerfile"

>FROM scratch
COPY files-v1/lib/  /app/lib/
COPY files-v1/conf/ /app/conf/</markup>

<p>In your working directory directory create a file called <code>Dockerfile-v2</code> with the following contents:</p>

<markup
lang="dockerfile"

>FROM scratch
COPY files-v2/lib/  /app/lib/
COPY files-v2/conf/ /app/conf/</markup>

</div>

<h3 id="_4_add_the_required_config_files">4. Add the required config files</h3>
<div class="section">
<p>Add the following content to a file in <code>files-v1/conf</code> called <code>storage-cache-config.xml</code>.</p>

<div class="admonition note">
<p class="admonition-inline">This is the <code>VERSION 1</code> cache config which has a single service called <code>HRPartitionedCache</code>.</p>
</div>
<markup
lang="xml"

>&lt;?xml version='1.0'?&gt;
&lt;cache-config xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
              xmlns="http://xmlns.oracle.com/coherence/coherence-cache-config"
              xsi:schemaLocation="http://xmlns.oracle.com/coherence/coherence-cache-config coherence-cache-config.xsd"&gt;
  &lt;!-- v1 Cache Config --&gt;
  &lt;caching-scheme-mapping&gt;
    &lt;cache-mapping&gt;
      &lt;cache-name&gt;hr-*&lt;/cache-name&gt;
      &lt;scheme-name&gt;hr-scheme&lt;/scheme-name&gt;
    &lt;/cache-mapping&gt;
  &lt;/caching-scheme-mapping&gt;

  &lt;caching-schemes&gt;
    &lt;distributed-scheme&gt;
      &lt;scheme-name&gt;hr-scheme&lt;/scheme-name&gt;
      &lt;service-name&gt;HRPartitionedCache&lt;/service-name&gt;
      &lt;backing-map-scheme&gt;
        &lt;local-scheme&gt;
          &lt;high-units&gt;{back-limit-bytes 0B}&lt;/high-units&gt;
        &lt;/local-scheme&gt;
      &lt;/backing-map-scheme&gt;
      &lt;autostart&gt;true&lt;/autostart&gt;
    &lt;/distributed-scheme&gt;
  &lt;/caching-schemes&gt;
&lt;/cache-config&gt;</markup>

<p>Add the following content to a file in <code>files-v2/conf</code> called <code>storage-cache-config.xml</code>.</p>

<div class="admonition note">
<p class="admonition-inline">This is the <code>VERSION 2</code> cache config which adds an additional cahe mapping and cache service called <code>FINPartitionedCache</code>.</p>
</div>
<markup
lang="xml"

>&lt;?xml version='1.0'?&gt;
&lt;cache-config xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
              xmlns="http://xmlns.oracle.com/coherence/coherence-cache-config"
              xsi:schemaLocation="http://xmlns.oracle.com/coherence/coherence-cache-config coherence-cache-config.xsd"&gt;
  &lt;!-- v2 Cache Config --&gt;
  &lt;caching-scheme-mapping&gt;
    &lt;cache-mapping&gt;
      &lt;cache-name&gt;hr-*&lt;/cache-name&gt;
      &lt;scheme-name&gt;hr-scheme&lt;/scheme-name&gt;
    &lt;/cache-mapping&gt;
    &lt;cache-mapping&gt;
      &lt;cache-name&gt;fin-*&lt;/cache-name&gt;
      &lt;scheme-name&gt;fin-scheme&lt;/scheme-name&gt;
    &lt;/cache-mapping&gt;
  &lt;/caching-scheme-mapping&gt;

  &lt;caching-schemes&gt;
    &lt;distributed-scheme&gt;
      &lt;scheme-name&gt;hr-scheme&lt;/scheme-name&gt;
      &lt;service-name&gt;HRPartitionedCache&lt;/service-name&gt;
      &lt;backing-map-scheme&gt;
        &lt;local-scheme&gt;
          &lt;high-units&gt;{back-limit-bytes 0B}&lt;/high-units&gt;
        &lt;/local-scheme&gt;
      &lt;/backing-map-scheme&gt;
      &lt;autostart&gt;true&lt;/autostart&gt;
    &lt;/distributed-scheme&gt;

    &lt;distributed-scheme&gt;
      &lt;scheme-name&gt;fin-scheme&lt;/scheme-name&gt;
      &lt;service-name&gt;FINPartitionedCache&lt;/service-name&gt;
      &lt;backing-map-scheme&gt;
        &lt;local-scheme&gt;
          &lt;high-units&gt;{back-limit-bytes 0B}&lt;/high-units&gt;
        &lt;/local-scheme&gt;
      &lt;/backing-map-scheme&gt;
      &lt;autostart&gt;true&lt;/autostart&gt;
    &lt;/distributed-scheme&gt;
  &lt;/caching-schemes&gt;
&lt;/cache-config&gt;</markup>

</div>

<h3 id="_5_build_the_docker_images">5. Build the Docker images</h3>
<div class="section">
<p>In your <code>working directory</code>, issue the following:</p>

<markup
lang="bash"

>docker build -t rolling-example:1.0.0 -f Dockerfile-v1 .

docker build -t rolling-example:2.0.0 -f Dockerfile-v2 .

docker images | grep rolling-example
REPOSITORY              TAG     IMAGE ID            CREATED             SIZE
rolling-example         2.0.0   3e195af6d5e1        8 seconds ago       1.36kB
rolling-example         1.0.0   5ce9152dd12c        26 seconds ago      890B</markup>

</div>

<h3 id="_6_create_the_coherence_cluster_yaml">6. Create the Coherence cluster yaml</h3>
<div class="section">
<p>Create the file <code>rolling-cluster.yaml</code> with the following contents.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: rolling-cluster
spec:
  replicas: 3
  coherence:
    cacheConfig: storage-cache-config.xml
  application:
    image: rolling-example:1.0.0</markup>

<div class="admonition note">
<p class="admonition-inline">Add an <code>imagePullSecrets</code> entry if required to pull images from a private repository.</p>
</div>
</div>

<h3 id="_7_install_the_coherence_cluster">7. Install the Coherence Cluster</h3>
<div class="section">
<p>Issue the following to install the cluster:</p>

<markup
lang="bash"

>kubectl create -n &lt;namespace&gt; -f rolling-cluster.yaml

coherencecluster.coherence.oracle.com/rolling-cluster created

kubectl -n &lt;namespace&gt; get pod -l coherenceCluster=rolling-cluster

NAME                        READY   STATUS    RESTARTS   AGE
rolling-cluster-storage-0   1/1     Running   0          58s
rolling-cluster-storage-1   1/1     Running   0          58s
rolling-cluster-storage-2   1/1     Running   0          58s</markup>

<div class="admonition note">
<p class="admonition-inline">Ensure all pods are running and ready before you continue.</p>
</div>
</div>

<h3 id="_8_add_data_to_a_cache_in_the_hr_service">8. Add Data to a cache in the HR service</h3>
<div class="section">
<markup
lang="bash"

>kubectl exec -it -n &lt;namespace&gt; rolling-cluster-storage-0 bash /scripts/startCoherence.sh console</markup>

<p>At the prompt, type <code>cache hr-test</code> and you will notice the following indicating your
cache configuration file with the service name of <code>HRPartitionedCache</code> is being loaded.</p>

<markup
lang="bash"

>...
Cache Configuration: hr-test
  SchemeName: server
  AutoStart: true
  ServiceName: HRPartitionedCache
..</markup>

<p>Use the following to create 10,000 entries of 100 bytes:</p>

<markup
lang="bash"

>bulkput 10000 100 0 100</markup>

<p>Lastly issue the command <code>size</code> to verify the cache entry count.</p>

<p>Issue the following to confirm there is no cache mapping and service for <code>fin-*</code> as yet.</p>

<markup
lang="bash"

>cache fin-test

java.lang.IllegalArgumentException: ensureCache cannot find a mapping for cache fin-test</markup>

<p>Type <code>bye</code> to exit the console.</p>

</div>

<h3 id="_9_update_the_application_image_version_to_2_0_0">9. Update the application image version to 2.0.0</h3>
<div class="section">
<p>Edit the <code>rolling-cluster.yaml</code> file and change the <code>image:</code> version from <code>1.0.0</code> to <code>2.0.0</code>.</p>

<markup
lang="yaml"

>image: rolling-example:2.0.0</markup>

<p>Issue the following to apply the new yaml:</p>

<markup
lang="bash"

>kubectl apply -n &lt;namespace&gt; -f rolling-cluster.yaml

coherencecluster.coherence.oracle.com/rolling-cluster configured</markup>

<p>Use the following command to check the status of the rolling upgrade of all pods.</p>

<div class="admonition note">
<p class="admonition-inline">The command below will not return until upgrade of all pods is complete.</p>
</div>
<markup
lang="bash"

>kubectl -n &lt;namespace&gt; rollout status sts/rolling-cluster-storage

Waiting for 1 pods to be ready...
statefulset rolling update complete 3 pods at revision rolling-cluster-storage-67f5cfdcb...</markup>

</div>

<h3 id="_10_validate_the_hr_cache_data_still_exists">10. Validate the HR cache data still exists</h3>
<div class="section">
<markup
lang="bash"

>kubectl exec -it -n &lt;namespace&gt; rolling-cluster-storage-0 bash /scripts/startCoherence.sh console</markup>

<p>At the prompt, type <code>cache hr-test</code> and then <code>size</code> and you will see the 10,000 entries are still present
because the upgrade was done is a safe manner.</p>

</div>

<h3 id="_11_add_data_to_a_cache_in_the_new_hr_service">11. Add Data to a cache in the new HR service</h3>
<div class="section">
<p>At the prompt, type <code>cache find-test</code> and you will notice the following indicating your
cache configuration file with the service name of <code>FINPartitionedCache</code> is now being loaded.</p>

<markup
lang="bash"

>...
Cache Configuration: fin-test
  SchemeName: server
  AutoStart: true
  ServiceName: FINPartitionedCache
..</markup>

<p>Use the following to create 10,000 entries of 100 bytes:</p>

<markup
lang="bash"

>bulkput 10000 100 0 100</markup>

<p>Lastly issue the command <code>size</code> to verify the cache entry count.</p>

<p>Type <code>bye</code> to exit the console.</p>

</div>

<h3 id="_12_uninstall_the_coherence_cluster">12. Uninstall the Coherence Cluster</h3>
<div class="section">
<markup
lang="bash"

>kubectl delete -n &lt;namespace&gt; -f rolling-cluster.yaml

coherencecluster.coherence.oracle.com "rolling-cluster" deleted</markup>

</div>
</div>
</doc-view>
