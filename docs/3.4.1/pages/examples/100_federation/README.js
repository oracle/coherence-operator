<doc-view>

<h2 id="_coherence_federation">Coherence Federation</h2>
<div class="section">
<p>This simple example demonstrates the Coherence federation feature.  It shows how to deploy two Coherence clusters that federating data between them using the Coherence Operator. The Coherence federation feature requires Coherence Grid Edition. See <a id="" title="" target="_blank" href="https://oracle.github.io/coherence-operator/docs/latest/#/docs/installation/04_obtain_coherence_images">Obtain Coherence Images</a> on how to get a commercial Coherence image.</p>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p><img src="./images/GitHub-Mark-32px.png" alt="GitHub Mark 32px" />
 The complete source code for this example is in the <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/tree/main/examples/100_federation">Coherence Operator GitHub</a> repository.</p>
</p>
</div>

<h3 id="_what_the_example_will_cover">What the Example will Cover</h3>
<div class="section">
<ul class="ulist">
<li>
<p><router-link to="#install-operator" @click.native="this.scrollFix('#install-operator')">Install the Coherence Operator</router-link></p>

</li>
<li>
<p><router-link to="#create-the-example-namespace" @click.native="this.scrollFix('#create-the-example-namespace')">Create the example namespace</router-link></p>

</li>
<li>
<p><router-link to="#create-secret" @click.native="this.scrollFix('#create-secret')">Create image pull and config store secrets</router-link></p>

</li>
<li>
<p><router-link to="#example" @click.native="this.scrollFix('#example')">Run the Example</router-link></p>

</li>
<li>
<p><router-link to="#cleanup" @click.native="this.scrollFix('#cleanup')">Cleaning Up</router-link></p>

</li>
</ul>
</div>

<h3 id="install-operator">Install the Coherence Operator</h3>
<div class="section">
<p>To run the examples below, you will need to have installed the Coherence Operator, do this using whatever method you prefer from the <a id="" title="" target="_blank" href="https://oracle.github.io/coherence-operator/docs/latest/#/docs/installation/01_installation">Installation Guide</a>.</p>

<p>Once you complete, confirm the operator is running, for example:</p>

<markup
lang="bash"

>kubectl get pods -n coherence

NAME                                                     READY   STATUS    RESTARTS   AGE
coherence-operator-controller-manager-74d49cd9f9-sgzjr   1/1     Running   1          27s</markup>

</div>
</div>

<h2 id="create-the-example-namespace">Create the example namespace</h2>
<div class="section">
<p>First, run the following command to create the namespace, coherence-example, for the example:</p>

<markup
lang="bash"

>kubectl create namespace coherence-example

namespace/coherence-example created</markup>


<h3 id="create-secret">Create image pull and configure store secrets</h3>
<div class="section">
<p>This example reqires two secrets:</p>

<ul class="ulist">
<li>
<p>An image pull secret named ocr-pull-secret containing your OCR credentials to be used by the example.</p>

</li>
</ul>
<p>Use a command similar to the following to create the image pull secret:</p>

<markup
lang="bash"

>kubectl create secret docker-registry ocr-pull-secret \
    --docker-server=container-registry.oracle.com \
    --docker-username="&lt;username&gt;" --docker-password="&lt;password&gt;" \
    --docker-email="&lt;email&gt;" -n coherence-example</markup>

<ul class="ulist">
<li>
<p>A configure store secret named storage-config to store the Coherence configuration files.</p>

</li>
</ul>
<p>Run the following command to create the configure store secret:</p>

<markup
lang="bash"

>kubectl create secret generic storage-config -n coherence-example \
    --from-file=src/main/resources/tangosol-coherence-override.xml \
    --from-file=src/main/resources/storage-cache-config.xml</markup>

</div>

<h3 id="example">Run the Example</h3>
<div class="section">
<p>Ensure you are in the <code>examples/federation</code> directory to run the example. This example uses the yaml files <code>src/main/yaml/primary-cluster.yaml</code> and <code>src/main/yaml/secondary-cluster.yaml</code>, which
define a primary cluster and a secondary cluster.</p>


<h4 id="_1_install_the_coherence_clusters">1. Install the Coherence clusters</h4>
<div class="section">
<p>Run the following commands to create the primary and secondary clusters:</p>

<markup
lang="bash"

>kubectl -n coherence-example create -f src/main/yaml/primary-cluster.yaml

coherence.coherence.oracle.com/primary-cluster created</markup>

<markup
lang="bash"

>kubectl -n coherence-example create -f src/main/yaml/secondary-cluster.yaml

coherence.coherence.oracle.com/secondary-cluster created</markup>

</div>

<h4 id="_2_list_the_created_coherence_clusters">2. List the created Coherence clusters</h4>
<div class="section">
<p>Run the following command to list the clusters:</p>

<markup
lang="bash"

>kubectl -n coherence-example get coherence

NAME                CLUSTER             ROLE                REPLICAS   READY   PHASE
primary-cluster     primary-cluster     primary-cluster     2          2       Ready
secondary-cluster   secondary-cluster   secondary-cluster   2          2       Ready</markup>

<p>To see the Coherence cache configuration file loaded from the secret volumn we defined, run the following command:</p>

<markup
lang="bash"

>kubectl logs -n coherence-example primary-cluster-0 | grep "Loaded cache"

... Oracle Coherence GE 14.1.1.0.0 &lt;Info&gt; (thread=main, member=n/a): Loaded cache configuration from "file:/config/storage-cache-config.xml"</markup>

</div>

<h4 id="_3_view_the_running_pods">3. View the running pods</h4>
<div class="section">
<p>Run the following command to view the Pods:</p>

<markup
lang="bash"

>kubectl -n coherence-example get pods</markup>

<markup
lang="bash"

>NAME                  READY   STATUS    RESTARTS   AGE
primary-cluster-0     1/1     Running   0          83s
primary-cluster-1     1/1     Running   0          83s
secondary-cluster-0   1/1     Running   0          74s
secondary-cluster-1   1/1     Running   0          73s</markup>

</div>

<h4 id="_4_connect_to_the_coherence_console_inside_the_primary_cluster_to_add_data">4. Connect to the Coherence Console inside the primary cluster to add data</h4>
<div class="section">
<p>We will connect via Coherence console to add some data using the following commands:</p>

<markup
lang="bash"

>kubectl exec -it -n coherence-example primary-cluster-0 /coherence-operator/utils/runner console</markup>

<p>At the prompt type the following to create a cache called <code>test</code>:</p>

<markup
lang="bash"

>cache test</markup>

<p>Use the following to add an entry with "primarykey" and "primaryvalue":</p>

<markup
lang="bash"

>put "primarykey" "primaryvalue"</markup>

<p>Use the following to create 10,000 entries of 100 bytes:</p>

<markup
lang="bash"

>bulkput 10000 100 0 100</markup>

<p>Lastly issue the command <code>size</code> to verify the cache entry count. It should be 10001.</p>

<p>Type <code>bye</code> to exit the console.</p>

</div>

<h4 id="_6_connect_to_the_coherence_console_inside_the_secondary_cluster_to_verify_that_data_is_federated_from_primary_cluster">6. Connect to the Coherence Console inside the secondary cluster to verify that data is federated from primary cluster</h4>
<div class="section">
<p>We will connect via Coherence console to confirm that the data we added to the primary cluster is federated to the secondary cluster.</p>

<markup
lang="bash"

>kubectl exec -it -n coherence-example secondary-cluster-0 /coherence-operator/utils/runner console</markup>

<p>At the prompt type the following to set the cache to <code>test</code>:</p>

<markup
lang="bash"

>cache test</markup>

<p>Use the following to get entry with "primarykey":</p>

<markup
lang="bash"

>get "primarykey"
primaryvalue</markup>

<p>Issue the command <code>size</code> to verify the cache entry count. It should be 10001.</p>

<p>Our federation has Active/Active topology. So, the data changes in both primary and secondary clusters are federated between the clusters. Use the following to add an entry with "secondarykey" and "secondaryvalue":</p>

<markup
lang="bash"

>put "secondarykey" "secondaryvalue"</markup>

</div>

<h4 id="_7_confirm_the_primary_cluster_also_received_secondarykey_secondaryvalue_entry">7. Confirm the primary cluster also received "secondarykey", "secondaryvalue" entry</h4>
<div class="section">
<p>Follow the command in the previous section to connect to the Coherence Console inside the primary cluster.</p>

<p>Use the following command to confirm that entry with "secondarykey" is federated to primary cluster:</p>

<markup
lang="bash"

>get "secondarykey"
secondaryvalue</markup>

</div>
</div>

<h3 id="cleanup">Cleaning up</h3>
<div class="section">
<p>Use the following commands to delete the primary and secondary clusters:</p>

<markup
lang="bash"

>kubectl -n coherence-example delete -f src/main/yaml/primary-cluster.yaml

kubectl -n coherence-example delete -f src/main/yaml/secondary-cluster.yaml</markup>

<p>Uninstall the Coherence operator using the undeploy commands for whichever method you chose to install it.</p>

</div>
</div>
</doc-view>
