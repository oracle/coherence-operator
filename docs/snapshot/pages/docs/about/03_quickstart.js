<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Quick Start</dt>
<dd slot="desc"><p>This guide is a simple set of steps to install the Coherence Operator and then use that
to install a simple Coherence cluster.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_prerequisites">Prerequisites</h2>
<div class="section">
<p>Ensure that the <router-link to="/docs/installation/01_installation">Coherence Operator prerequisites</router-link> are available.</p>

</div>

<h2 id="_1_install_the_coherence_operator">1. Install the Coherence Operator</h2>
<div class="section">
<p>If you want the default Coherence Operator installation then the simplest solution is use <code>kubectl</code> to apply the manifests from the Operator release.</p>

<markup
lang="bash"

>kubectl apply -f https://github.com/oracle/coherence-operator/releases/download/v3.4.2/coherence-operator.yaml</markup>

<p>This will create a namespace called <code>coherence</code> and install the Operator into it along with all the required <code>ClusterRole</code> and <code>RoleBinding</code> resources. The <code>coherence</code> namespace can be changed by downloading and editing the yaml file.</p>

<div class="admonition note">
<p class="admonition-inline">Because the <code>coherence-operator.yaml</code> manifest also creates the namespace, the corresponding <code>kubectl delete</code> command will <em>remove the namespace and everything deployed to it</em>! If you do not want this behaviour you should edit the <code>coherence-operator.yaml</code> to remove the namespace section from the start of the file.</p>
</div>

<h3 id="_alternatively_install_using_helm">Alternatively Install Using Helm</h3>
<div class="section">
<p>Alternatively you can install the Operator using the Helm chart.</p>


<h4 id="_add_the_coherence_operator_helm_repository">Add the Coherence Operator Helm repository</h4>
<div class="section">
<p>Add the Coherence Operator Helm repo to your local Helm.</p>

<markup
lang="bash"

>helm repo add coherence https://oracle.github.io/coherence-operator/charts

helm repo update</markup>

<div class="admonition note">
<p class="admonition-inline">To avoid confusion, the URL <code><a id="" title="" target="_blank" href="https://oracle.github.io/coherence-operator/charts">https://oracle.github.io/coherence-operator/charts</a></code> is a Helm repo, it is not a web site you open in a browser. You may think we shouldn&#8217;t have to say this, but you&#8217;d be surprised.</p>
</div>
</div>

<h4 id="_install_the_coherence_operator_helm_chart">Install the Coherence Operator Helm chart</h4>
<div class="section">
<markup
lang="bash"
title="helm v3 install command"
>helm install  \
    --namespace &lt;namespace&gt; \
    &lt;release-name&gt; \
    coherence/coherence-operator</markup>

<p>e.g. if the Kubernetes namespace is <code>coherence-test</code> the command would be:</p>

<markup
lang="bash"
title="helm v3 install command"
>helm install --namespace coherence-test  operator coherence/coherence-operator</markup>

<p>or with Helm v2</p>

<markup
lang="bash"

>helm install --namespace coherence-test  --name operator coherence/coherence-operator</markup>

<p>See the <router-link to="/docs/installation/01_installation">full install guide</router-link> for more details.</p>

</div>
</div>
</div>

<h2 id="_2_install_a_coherence_deployment">2. Install a Coherence Deployment</h2>
<div class="section">
<p>Ensure that the Coherence images can be pulled by the Kubernetes cluster,
see <router-link to="/docs/installation/04_obtain_coherence_images">Obtain Coherence Images</router-link>.
By default, a <code>Coherence</code> resource will use the OSS Coherence CE image from Docker Hub.
If a different image is to be used the image name will need to be specified in the <code>Coherence</code> yaml,
see <router-link to="/docs/applications/010_overview">Setting the Application Image</router-link> for documentation on how to
specify a different images to use.</p>


<h3 id="_2_1_install_a_coherence_resource_using_the_minimal_required_configuration">2.1 Install a Coherence resource using the minimal required configuration.</h3>
<div class="section">
<p>The minimal required yaml to create a <code>Coherence</code> resource is shown below.</p>

<markup
lang="yaml"
title="my-deployment.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: my-deployment <span class="conum" data-value="1" /></markup>

<p>The only required field is <code>metadata.name</code> which will be used as the Coherence cluster name, in this case <code>my-deployment</code></p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; apply -f my-deployment.yaml</markup>

<div class="admonition note">
<p class="admonition-inline">Use the same namespace that the operator was installed into,
e.g. if the namespace is <code>coherence-test</code> the command would be
<code>kubectl -n coherence-test create -f my-deployment.yaml</code></p>
</div>
</div>

<h3 id="_2_2_list_the_coherence_resources">2.2 List the Coherence Resources</h3>
<div class="section">
<p>After installing the <code>my-deployment.yaml</code> above here should be a single <code>Coherence</code> resource  named <code>my-deployment</code> in the Coherence Operator namespace.</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; get coherence</markup>

<p>or alternatively using the <code>Coherence</code> CRD a short name of <code>coh</code></p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; get coh</markup>

<p>e.g. if the namespace is <code>coherence-test</code> the command would be <code>kubectl -n coherence-test get coherence</code></p>

<markup
lang="bash"

>NAME                                                  AGE
coherence.coherence.oracle.com/my-deployment   19s</markup>

</div>

<h3 id="_2_3_list_all_of_the_pods_for_the_coherence_resource">2.3 List all of the <code>Pods</code> for the Coherence resource.</h3>
<div class="section">
<p>The Coherence Operator applies a <code>coherenceDeployment</code> label to all <code>Pods</code> so this label can be used with the <code>kubectl</code> command to find <code>Pods</code> for a <code>CoherenceCoherence</code> resource.</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; get pod -l coherenceDeployment=my-deployment</markup>

<p>e.g. if the namespace is <code>coherence</code> the command would be:
<code>kubectl -n coherence get pod -l coherenceDeployment=my-deployment</code></p>

<markup
lang="bash"

>NAME              READY   STATUS    RESTARTS   AGE
my-deployment-0   1/1     Running   0          2m58s
my-deployment-1   1/1     Running   0          2m58s
my-deployment-2   1/1     Running   0          2m58s</markup>

</div>

<h3 id="_2_3_list_all_the_pods_for_the_coherence_cluster">2.3 List all the <code>Pods</code> for the Coherence cluster.</h3>
<div class="section">
<p>The Coherence Operator applies a <code>coherenceCluster</code> label to all <code>Pods</code>, so this label can be used with the <code>kubectl</code>
command to find all <code>Pods</code> for a Coherence cluster, which will be made up of multiple <code>Coherence</code> resources.</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; get pod -l coherenceCluster=my-cluster</markup>

<p>e.g. If there is a cluster named <code>my-cluster</code> made up of two <code>Coherence</code> resources in the namespace
<code>coherence-test</code>, one named <code>storage</code> and one named <code>front-end</code>
then the <code>kubectl</code> command to list all Pods for the cluster would be:</p>

<markup
lang="bash"

>kubectl -n coherence-test get pod -l coherenceCluster=my-cluster</markup>

<p>The result of which might look something like this</p>

<markup
lang="bash"

>NAME          READY   STATUS    RESTARTS   AGE
storage-0     1/1     Running   0          2m58s
storage-1     1/1     Running   0          2m58s
storage-2     1/1     Running   0          2m58s
front-end-0   1/1     Running   0          2m58s
front-end-1   1/1     Running   0          2m58s
front-end-2   1/1     Running   0          2m58s</markup>

</div>
</div>

<h2 id="_3_scale_the_coherence_cluster">3. Scale the Coherence Cluster</h2>
<div class="section">

<h3 id="_3_1_use_kubectl_to_scale_up">3.1 Use kubectl to Scale Up</h3>
<div class="section">
<p>Using the <code>kubectl scale</code> command a specific <code>Coherence</code> resource can be scaled up or down.</p>

<markup
lang="bash"

>kubectl -n &lt;namespace&gt; scale coherence/my-deployment --replicas=6</markup>

<p>e.g. if the namespace is <code>coherence-test</code> the command would be:
<code>kubectl -n coherence scale coherence/my-deployment --replicas=6</code></p>

</div>
</div>
</doc-view>
