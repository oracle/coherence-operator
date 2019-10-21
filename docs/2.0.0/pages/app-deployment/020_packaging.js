<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Packaging Applications</dt>
<dd slot="desc"><p>Whilst it is simple to deploy a Coherence cluster into Kubernetes in most cases there is also a requirement to add
application code and configuration to the Coherence JVMs class path.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_introduction">Introduction</h2>
<div class="section">
<p>A common scenario for Coherence deployments is to include specific user artefacts such as cache and
operational configuration files as well as user classes.</p>

<p>This can be achieved with Coherence Operator by specifying the application configuration
in the <code>application</code> section of the spec or for an individual role.
The <code>image</code> field specifies a Docker image from which the configuration and classes
are copied and added to the JVM classpath at runtime.</p>

<p>The <code>libDir</code> and <code>configDir</code> are optional fields below <code>application</code> and are described below:</p>

<ul class="ulist">
<li>
<p><code>libDir</code> - contains application classes, default value is <code>/app/lib</code></p>

</li>
<li>
<p><code>configDir</code>  - contains cache and operational configuration files, default value is <code>/app/conf</code></p>

</li>
</ul>
<p>The example yaml below instructs the Coherence Operator to attach a Docker image called <code>acme/orders-data:1.0.0</code>
at Pod startup and copy the artefacts in the <code>libDir</code> and <code>configDir</code> to the Pod and add
to the JVM classpath.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  application:
    image: acme/orders-data:1.0.0</markup>

<p>Since we are using the default directories, we would expect that the Docker image referenced above
would include two directories <code>/app/lib</code> and <code>/app/conf</code> containing the appropriate files.</p>

<p>The following Dockerfile could be used to create such an image,
assuming the the directories <code>files/lib</code> and <code>files/conf</code> contain the files to copy.</p>

<markup
lang="dockerfile"

>FROM scratch
COPY files/lib/  /app/lib/
COPY files/conf/ /app/conf/</markup>

<div class="admonition note">
<p class="admonition-inline">Is is recommended to use the <code>scratch</code> image in the <code>FROM</code> clause to minimize the size of the resultant image.</p>
</div>
<p>See the <router-link to="/clusters/070_applications">Coherence Applications</router-link> section for
full details on each of the fields in the application section.</p>


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

<h4 id="_2_create_the_dockerfile">2. Create the Dockerfile</h4>
<div class="section">
<p>In your working directory directory create a file called <code>Dockerfile</code> with the following contents:</p>

<markup
lang="dockerfile"

>FROM scratch
COPY files/lib/  /app/lib/
COPY files/conf/ /app/conf/</markup>

</div>

<h4 id="_3_add_the_required_config_files">3. Add the required config files</h4>
<div class="section">
<markup
lang="bash"

>mkdir -p files/lib files/conf</markup>

<p>Add the following content to a file in <code>/files/conf</code> called <code>storage-cache-config.xml</code>.</p>

<markup
lang="xml"

>&lt;?xml version='1.0'?&gt;
&lt;cache-config xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
              xmlns="http://xmlns.oracle.com/coherence/coherence-cache-config"
              xsi:schemaLocation="http://xmlns.oracle.com/coherence/coherence-cache-config coherence-cache-config.xsd"&gt;

  &lt;caching-scheme-mapping&gt;
    &lt;cache-mapping&gt;
      &lt;cache-name&gt;*&lt;/cache-name&gt;
      &lt;scheme-name&gt;server&lt;/scheme-name&gt;
    &lt;/cache-mapping&gt;
  &lt;/caching-scheme-mapping&gt;

  &lt;caching-schemes&gt;
    &lt;distributed-scheme&gt;
      &lt;scheme-name&gt;server&lt;/scheme-name&gt;
      &lt;service-name&gt;ExamplePartitionedCache&lt;/service-name&gt;
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

<h4 id="_4_build_the_docker_image">4. Build the Docker image</h4>
<div class="section">
<p>In your <code>working directory</code>, issue the following:</p>

<markup
lang="bash"

>docker build -t packaging-example:1.0.0 .

Step 1/3 : FROM scratch
 ---&gt;
Step 2/3 : COPY files/lib/  /app/lib/
 ---&gt; c91db5a34f5c
Step 3/3 : COPY files/conf/ /app/conf/
 ---&gt; 7dd0b5f3e37a
Successfully built 7dd0b5f3e37a
Successfully tagged packaging-example:1.0.0</markup>

<div class="admonition note">
<p class="admonition-inline">In this example we have created but not populated the <code>lib</code> directory which would be used for application classes.</p>
</div>
</div>

<h4 id="_5_create_the_coherence_cluster_yaml">5. Create the Coherence cluster yaml</h4>
<div class="section">
<p>Create the file <code>packaging-cluster.yaml</code> with the following contents.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: packaging-cluster
spec:
  replicas: 3
  coherence:
    cacheConfig: storage-cache-config.xml
  application:
    image: packaging-example:1.0.0</markup>

<div class="admonition note">
<p class="admonition-inline">The default Coherence image used comes from <a id="" title="" target="_blank" href="https://container-registry.oracle.com">Oracle Container Registry</a>
so unless that image has already been pulled onto the Kubernetes nodes an <code>imagePullSecrets</code> field will be required
to pull the image.
See <router-link to="/about/04_obtain_coherence_images">Obtain Coherence Images</router-link></p>
</div>
</div>

<h4 id="_6_install_the_coherence_cluster">6. Install the Coherence Cluster</h4>
<div class="section">
<p>Issue the following to install the cluster:</p>

<markup
lang="bash"

>kubectl create -n &lt;namespace&gt; -f packaging-cluster.yaml

coherencecluster.coherence.oracle.com/packaging-cluster created

kubectl -n &lt;namespace&gt; get pod -l coherenceCluster=packaging-cluster

NAME                          READY   STATUS    RESTARTS   AGE
packaging-cluster-storage-0   1/1     Running   0          58s
packaging-cluster-storage-1   1/1     Running   0          58s
packaging-cluster-storage-2   1/1     Running   0          58s</markup>

</div>

<h4 id="_7_add_data_to_the_coherence_cluster_via_the_console">7. Add Data to the Coherence Cluster via the Console</h4>
<div class="section">
<markup
lang="bash"

>kubectl exec -it -n &lt;namespace&gt; packaging-cluster-storage-0 bash /scripts/startCoherence.sh console</markup>

<p>At the prompt, type <code>cache test</code> and you will notice the following indicating your
cache configuration file with the service name of <code>ExamplePartitionedCache</code> is being loaded.</p>

<markup
lang="bash"

>...
Cache Configuration: test
  SchemeName: server
  AutoStart: true
  ServiceName: ExamplePartitionedCache
..</markup>

</div>

<h4 id="_8_uninstall_the_coherence_cluster">8. Uninstall the Coherence Cluster</h4>
<div class="section">
<markup
lang="bash"

>kubectl delete -n &lt;namespace&gt; -f packaging-cluster.yaml

coherencecluster.coherence.oracle.com "packaging-cluster" deleted</markup>

</div>
</div>
</div>
</doc-view>
