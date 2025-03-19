<doc-view>

<h2 id="_coherence_federation">Coherence Federation</h2>
<div class="section">
<p>This example demonstrates the Coherence federation feature which allows you to federate cache data asynchronously across multiple geographically dispersed clusters.
Cached data is federated across clusters to provide redundancy, off-site backup, and multiple points of access for application users in different geographical locations.</p>

<div class="admonition note">
<p class="admonition-inline"><strong>Coherence federation feature requires Coherence Grid Edition.</strong></p>
</div>
<p>To demonstrate this feature, we will deploy two Coherence clusters, using the <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator">Coherence Operator</a>, located on separate OKE (Kubernetes Engine) clusters on Oracle Cloud (OCI).
It is assumed that you have two OCI regions, with these OKE clusters already configured, and have configured Dynamic Routing Gateways (DRGs)
to connect the two regions together.</p>

<div class="admonition tip">
<p class="admonition-inline">Although this example uses OCI, the concepts can be applied to other cloud providers to achieve the same result.</p>
</div>
<p>They key (cloud platform-agnostic) aspects of the example are:</p>

<ul class="ulist">
<li>
<p>The two Kubernetes clusters are located in separate cloud regions</p>

</li>
<li>
<p>Each region must be able to be connected with or communicate with the other region</p>

</li>
<li>
<p>Each region must have a network load balancer (LB), created either via deployment yaml or configured, that can forward traffic on port 40000, in our case, to the federation service in the OKE cluster</p>

</li>
<li>
<p>Routing rules must be setup to allow LBs to send traffic over the specified ports between the regions</p>

</li>
</ul>
<p>The diagram below outlines the setup for this example and uses the following OCI regions that have been configured with the relevant Kubernetes contexts.</p>

<ul class="ulist">
<li>
<p>Melbourne region has a Kubernetes context of <code>c1</code></p>

</li>
<li>
<p>Sydney region has a Kubernetes context of <code>c3</code></p>

</li>
</ul>


<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/images/federated-coherence.png" alt="Federated Setup"width="100%" />
</v-card-text>
</v-card>

<div class="admonition note">
<p class="admonition-inline"><strong>Although some network setup information is outlined below, it is assumed you have knowledge of, or access to
Oracle Cloud Infrastructure (OCI) administrators who can set up cross region networking using Dynamic Routing Gateways (DRG&#8217;s) and
Remote Peering Connections. DRG Documentation below for more information.</strong></p>
</div>
<p>See the links below for more information on various topics described above:</p>

<ul class="ulist">
<li>
<p><a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator">Coherence Operator on GitHub</a></p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://docs.oracle.com/en-us/iaas/Content/ContEng/Tasks/create-cluster.htm">Creating OKE Clusters</a></p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://www.oracle.com/au/cloud/networking/dynamic-routing-gateway/">Dynamic Routing Gateway Documentation</a></p>

</li>
<li>
<p><a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/fusion-middleware/coherence/14.1.2/administer/federating-caches-clusters.html">Coherence Federation Documentation</a></p>

</li>
<li>
<p><router-link to="#../installation/04_obtain_coherence_images.adoc" @click.native="this.scrollFix('#../installation/04_obtain_coherence_images.adoc')">Obtaining Commercial Coherence Images</router-link></p>

</li>
</ul>
<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p>The complete source code for this example is in the <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/tree/main/examples/100_federation">Coherence Operator GitHub</a> repository.</p>
</p>
</div>

<h3 id="_what_the_example_will_cover">What the Example will Cover</h3>
<div class="section">
<ul class="ulist">
<li>
<p><router-link to="#prereqs" @click.native="this.scrollFix('#prereqs')">Prerequisites and Network Setup</router-link></p>
<ul class="ulist">
<li>
<p><router-link to="#container-registry" @click.native="this.scrollFix('#container-registry')">Obtain your container registry Auth token</router-link></p>

</li>
<li>
<p><router-link to="#create-oke" @click.native="this.scrollFix('#create-oke')">Create OKE clusters</router-link></p>

</li>
<li>
<p><router-link to="#network-setup" @click.native="this.scrollFix('#network-setup')">Complete Network Setup</router-link></p>

</li>
</ul>
</li>
<li>
<p><router-link to="#prepare" @click.native="this.scrollFix('#prepare')">Prepare the OKE clusters</router-link></p>
<ul class="ulist">
<li>
<p><router-link to="#install-operator" @click.native="this.scrollFix('#install-operator')">Install the Coherence Operator</router-link></p>

</li>
<li>
<p><router-link to="#create-the-example-namespace" @click.native="this.scrollFix('#create-the-example-namespace')">Create the example namespace</router-link></p>

</li>
<li>
<p><router-link to="#create-secret" @click.native="this.scrollFix('#create-secret')">Create image pull secrets and configmaps</router-link></p>

</li>
</ul>
</li>
<li>
<p><router-link to="#explore" @click.native="this.scrollFix('#explore')">Explore the Configuration</router-link></p>
<ul class="ulist">
<li>
<p><router-link to="#explore-coherence" @click.native="this.scrollFix('#explore-coherence')">Coherence Configuration</router-link></p>

</li>
<li>
<p><router-link to="#expore-yaml" @click.native="this.scrollFix('#expore-yaml')">Coherence Operator YAML</router-link></p>

</li>
</ul>
</li>
<li>
<p><router-link to="#deploy" @click.native="this.scrollFix('#deploy')">Deploy the Example</router-link></p>
<ul class="ulist">
<li>
<p><router-link to="#update-ocid" @click.native="this.scrollFix('#update-ocid')">Update the OCID&#8217;s for the federation service</router-link></p>

</li>
<li>
<p><router-link to="#install-primary" @click.native="this.scrollFix('#install-primary')">Install the Primary Coherence cluster</router-link></p>

</li>
<li>
<p><router-link to="#install-secondary" @click.native="this.scrollFix('#install-secondary')">Install the Secondary Coherence cluster</router-link></p>

</li>
<li>
<p><router-link to="#re-apply" @click.native="this.scrollFix('#re-apply')">Re-apply the yaml for both clusters</router-link></p>

</li>
<li>
<p><router-link to="#check-federation" @click.native="this.scrollFix('#check-federation')">Check federation status</router-link></p>

</li>
</ul>
</li>
<li>
<p><router-link to="#run" @click.native="this.scrollFix('#run')">Run the example</router-link></p>
<ul class="ulist">
<li>
<p><router-link to="#add-data-p" @click.native="this.scrollFix('#add-data-p')">Add Data in the primary cluster</router-link></p>

</li>
<li>
<p><router-link to="#validate-data-p" @click.native="this.scrollFix('#validate-data-p')">Validate data in the secondary cluster</router-link></p>

</li>
<li>
<p><router-link to="#add-data-s" @click.native="this.scrollFix('#add-data-s')">Add Data in the secondary cluster</router-link></p>

</li>
<li>
<p><router-link to="#validate-data-s" @click.native="this.scrollFix('#validate-data-s')">Validate data in the primary cluster</router-link></p>

</li>
</ul>
</li>
<li>
<p><router-link to="#cleanup" @click.native="this.scrollFix('#cleanup')">Cleaning Up</router-link></p>

</li>
</ul>
</div>

<h3 id="prereqs">Prerequisites and Network Setup</h3>
<div class="section">

<h4 id="container-registry">Obtain your container registry Auth token</h4>
<div class="section">
<p>In this example we are using the Coherence Grid Edition container from Oracle&#8217;s container registry.</p>

<ul class="ulist">
<li>
<p><code>container-registry.oracle.com/middleware/coherence:14.1.2.0.0</code></p>

</li>
</ul>
<p>To be able to pull the above image you need to do the following:</p>

<ol style="margin-left: 15px;">
<li>
Sign in to the Container Registry at <code><a id="" title="" target="_blank" href="https://container-registry.oracle.com/">https://container-registry.oracle.com/</a></code>. (If you don&#8217;t have an account, with Oracle you will need to create one)

</li>
<li>
Once singed in, search for <code>coherence</code> and select the link for <code>Oracle Coherence</code>

</li>
<li>
If you have not already, you will need to accept the 'Oracle Standard Terms and Conditions' on the right before you are able to pull the image

</li>
<li>
Once you have accepted the terms and condition click on the drop-down next to your name on the top right and select <code>Auth Token</code>.

</li>
<li>
Click on <code>Generate Token</code> and save this in a secure place for use further down.

</li>
</ol>
<div class="admonition note">
<p class="admonition-inline">See <a id="" title="" target="_blank" href="https://docs.oracle.com/en-us/iaas/Content/Registry/Tasks/registrygettingauthtoken.htm">the OCI documentation</a> for more information on creating your Auth token.</p>
</div>
</div>

<h4 id="create-oke">Create OKE clusters</h4>
<div class="section">
<p>This example assumes you have already created two OKE clusters in separate regions.</p>

<div class="admonition note">
<p class="admonition-inline">For more information on creating OKE clusters, see <a id="" title="" target="_blank" href="https://docs.oracle.com/en-us/iaas/Content/devops/using/create_oke_environment.htm">the OCI documentation</a>.</p>
</div>
</div>

<h4 id="network-setup">Complete Network Setup</h4>
<div class="section">
<p>You must ensure you have the following for each region:</p>

<ol style="margin-left: 15px;">
<li>
A Dynamic Routing Gateway must be setup and must have a remote peering connection to the other region

</li>
<li>
The routing table associated with the worker nodes (OKE cluster) subnet for each region needs to have a rule to route traffic to the other region VCN via the Dynamic Routing Gateway

</li>
<li>
In this example, we are exposing coherence on port 40000 for federation, we require the following Security rules:
<ol style="margin-left: 15px;">
<li>
Egress for worker nodes to get other cluster load balancer via the DRG on port 40000

</li>
<li>
Ingress from worker to receive from their own load balancer on port 40000

</li>
<li>
Ingress rule on own LBR to receive traffic from other region on port 40000

</li>
</ol>
</li>
<li>
Each region has a dedicated private subnet for the network load balancer

</li>
</ol>
<p>Item 3 above is available via the OCI console via VCN&#8594;Security&#8594;Network Security Group Information&#8594;Worker Node Subnet&#8594;Security Rules.</p>

<div class="admonition note">
<p class="admonition-inline">Refer to the <a id="" title="" target="_blank" href="https://www.oracle.com/au/cloud/networking/dynamic-routing-gateway/">DRG Documentation</a> for more information.</p>
</div>
</div>
</div>

<h3 id="prepare">Prepare the OKE clusters</h3>
<div class="section">
<p>For each of the OKE clusters, carry out the following</p>

<ol style="margin-left: 15px;">
<li>
Install the Coherence Operator

</li>
<li>
Create the example namespace

</li>
<li>
Create image pull secrets and configmaps

</li>
</ol>

<h4 id="install-operator">1. Install the Coherence Operator</h4>
<div class="section">
<p>To run the examples below, you will need to have installed the Coherence Operator,
do this using whatever method you prefer from the <a id="" title="" target="_blank" href="https://oracle.github.io/coherence-operator/docs/latest/#/docs/installation/01_installation">Installation Guide</a>.</p>

<markup
lang="bash"

>kubectl apply -f https://github.com/oracle/coherence-operator/releases/download/v3.4.3/coherence-operator.yaml</markup>

<p>Once you complete, confirm the operator is running, for example:</p>

<markup
lang="bash"

>kubectl get pods -n coherence

NAME                                                    READY   STATUS    RESTARTS   AGE
coherence-operator-controller-manager-846887895-d7wzg   1/1     Running   0          40s
coherence-operator-controller-manager-846887895-l7dcl   1/1     Running   0          40s
coherence-operator-controller-manager-846887895-rtsjv   1/1     Running   0          40s</markup>

</div>

<h4 id="create-the-example-namespace">2. Create the example namespace</h4>
<div class="section">
<p>First, run the following command to create the namespace, coherence-example, for the example:</p>

<markup
lang="bash"

>kubectl create namespace coherence-example

namespace/coherence-example created</markup>

</div>

<h4 id="create-secret">3. Create image pull secrets and configmaps</h4>
<div class="section">
<p>This example requires a number of secrets:</p>

<ul class="ulist">
<li>
<p>An image pull secret named <code>ocr-pull-secret</code> containing your OCR credentials to be used by the example.</p>

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
<p>A store secret named <code>storage-config</code> to store the Coherence configuration files.</p>

</li>
</ul>
<p>Run the following command to create the config secret:</p>

<markup
lang="bash"

>kubectl create secret generic storage-config -n coherence-example \
    --from-file=src/main/resources/tangosol-coherence-override.xml \
    --from-file=src/main/resources/storage-cache-config.xml</markup>

</div>
</div>

<h3 id="explore">Explore the Configuration</h3>
<div class="section">

<h4 id="explore-coherence">Coherence Configuration</h4>
<div class="section">
<p>In this example there are two Coherence Configuration files.</p>

<p><strong>Coherence Operational Override</strong></p>

<p>The Coherence override file typically contains information regarding the cluster configuration. In this case we are specifying the edition
and the federation configuration.</p>

<markup
lang="xml"

>&lt;?xml version="1.0"?&gt;
&lt;!--
  ~ Copyright (c) 2021, 2025 Oracle and/or its affiliates.
  ~ Licensed under the Universal Permissive License v 1.0 as shown at
  ~ http://oss.oracle.com/licenses/upl.
  --&gt;

&lt;!--
  Grid Edition version of the override file which includes Federation.
--&gt;
&lt;coherence xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
           xmlns="http://xmlns.oracle.com/coherence/coherence-operational-config"
           xsi:schemaLocation="http://xmlns.oracle.com/coherence/coherence-operational-config coherence-operational-config.xsd"&gt;

  &lt;license-config&gt;
    &lt;edition-name system-property="coherence.edition"&gt;GE&lt;/edition-name&gt;
  &lt;/license-config&gt;

  &lt;!--
    Define a federation configuration for PrimaryCluster and SecondaryCluster
    where the default topology is Active-Active.
    --&gt;
  &lt;federation-config&gt;
    &lt;participants&gt;
      &lt;participant&gt;
        &lt;name system-property="primary.cluster"/&gt; <span class="conum" data-value="1" />
        &lt;initial-action&gt;start&lt;/initial-action&gt;
        &lt;remote-addresses&gt;
          &lt;socket-address&gt;
            &lt;address system-property="primary.cluster.address"/&gt;
            &lt;port    system-property="primary.cluster.port"/&gt;
          &lt;/socket-address&gt;
        &lt;/remote-addresses&gt;
      &lt;/participant&gt;
      &lt;participant&gt;
        &lt;name system-property="secondary.cluster"/&gt; <span class="conum" data-value="2" />
        &lt;initial-action&gt;start&lt;/initial-action&gt;
        &lt;remote-addresses&gt;
          &lt;socket-address&gt;
            &lt;address system-property="secondary.cluster.address"/&gt;
            &lt;port    system-property="secondary.cluster.port"/&gt;
          &lt;/socket-address&gt;
        &lt;/remote-addresses&gt;
      &lt;/participant&gt;
    &lt;/participants&gt;
    &lt;topology-definitions&gt;
      &lt;active-active&gt;
        &lt;name&gt;Active&lt;/name&gt; <span class="conum" data-value="3" />
        &lt;active system-property="primary.cluster"/&gt;
        &lt;active system-property="secondary.cluster"/&gt;
      &lt;/active-active&gt;
    &lt;/topology-definitions&gt;
  &lt;/federation-config&gt;
&lt;/coherence&gt;</markup>

<ul class="colist">
<li data-value="1">Defines the primary-cluster name, address and port. These are overridden by environment variables in the deployment yaml.</li>
<li data-value="2">Defines the secondary-cluster name, address and port.</li>
<li data-value="3">Sets the topology to be active-active</li>
</ul>
<p><strong>Coherence Cache Configuration</strong></p>

<p>The Coherence cache configuration file contains the cache definitions. In this case we are specifying that all caches are federated.</p>

<markup
lang="xml"

>&lt;?xml version='1.0'?&gt;

&lt;!--
  ~ Copyright (c) 2021, 2025 Oracle and/or its affiliates.
  ~ Licensed under the Universal Permissive License v 1.0 as shown at
  ~ http://oss.oracle.com/licenses/upl.
  --&gt;

&lt;cache-config xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
              xmlns="http://xmlns.oracle.com/coherence/coherence-cache-config"
              xsi:schemaLocation="http://xmlns.oracle.com/coherence/coherence-cache-config coherence-cache-config.xsd"&gt;

  &lt;caching-scheme-mapping&gt;
    &lt;cache-mapping&gt; <span class="conum" data-value="1" />
      &lt;cache-name&gt;*&lt;/cache-name&gt;
      &lt;scheme-name&gt;server&lt;/scheme-name&gt;
    &lt;/cache-mapping&gt;
  &lt;/caching-scheme-mapping&gt;

  &lt;caching-schemes&gt;

    &lt;federated-scheme&gt; <span class="conum" data-value="2" />
      &lt;scheme-name&gt;server&lt;/scheme-name&gt;
      &lt;backing-map-scheme&gt;
        &lt;local-scheme/&gt;
      &lt;/backing-map-scheme&gt;
      &lt;autostart&gt;true&lt;/autostart&gt;
      &lt;address-provider&gt;
        &lt;local-address&gt;  <span class="conum" data-value="3" />
          &lt;address system-property="coherence.extend.address"/&gt;
          &lt;port system-property="coherence.federation.port"&gt;40000&lt;/port&gt;
        &lt;/local-address&gt;
      &lt;/address-provider&gt;
      &lt;topologies&gt;
        &lt;topology&gt;
          &lt;name&gt;Active&lt;/name&gt;  <span class="conum" data-value="4" />
        &lt;/topology&gt;
      &lt;/topologies&gt;
    &lt;/federated-scheme&gt;
  &lt;/caching-schemes&gt;
&lt;/cache-config&gt;</markup>

<ul class="colist">
<li data-value="1">Defines the cache mapping for all caches, '*', to map to the <code>server</code> scheme.</li>
<li data-value="2">The <code>server</code> scheme is a federated cache scheme</li>
<li data-value="3">Defines the local address and port on which to listen. Empty address will translate to <code>0.0.0.0</code></li>
<li data-value="4">Specifies the topology, referenced above, to use.</li>
</ul>
</div>

<h4 id="expore-yaml">Coherence Operator YAML</h4>
<div class="section">
<p><strong>Primary Cluster</strong></p>

<p>The following yaml file is used by the Coherence Operator to deploy the primary cluster.</p>

<markup
lang="yaml"

>#
# Copyright (c) 2021, 2025 Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
# Federation Example
# Primary cluster in an Active/Active topology
apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: primary-cluster <span class="conum" data-value="1" />
spec:
  jvm:
    classpath:
      - /config
      - /u01/oracle/oracle_home/coherence/lib/coherence.jar
  env: <span class="conum" data-value="2" />
    - name: "PRIMARY_CLUSTER"
      value: "primary-cluster"
    - name: "PRIMARY_CLUSTER_ADDRESS"
      value: ""
    - name: "PRIMARY_CLUSTER_PORT"
      value: "40000"
    - name: "SECONDARY_CLUSTER"
      value: "secondary-cluster"
    - name: "SECONDARY_CLUSTER_ADDRESS"
      value: ""
    - name: "SECONDARY_CLUSTER_PORT"
      value: "40000"
  secretVolumes: <span class="conum" data-value="3" />
    - mountPath: /config
      name: storage-config
  ports:
    - name: "federation" <span class="conum" data-value="4" />
      port: 40000
      protocol: TCP
      service:
        port: 40000
        type: LoadBalancer
        annotations: <span class="conum" data-value="5" />
          oci.oraclecloud.com/load-balancer-type: "nlb"
          oci-network-load-balancer.oraclecloud.com/internal: "true"
          oci-network-load-balancer.oraclecloud.com/subnet: "(Internal subnet OCID - REPLACE ME)"
          oci-network-load-balancer.oraclecloud.com/oci-network-security-groups: "(OCID of the NSG - REPLACE ME)"
    - name: management
  coherence:  <span class="conum" data-value="6" />
    cacheConfig: /config/storage-cache-config.xml
    overrideConfig: /config/tangosol-coherence-override.xml
    logLevel: 9
  image: container-registry.oracle.com/middleware/coherence:14.1.2.0.0 <span class="conum" data-value="7" />
  imagePullSecrets:
    - name: ocr-pull-secret
  replicas: 3</markup>

<ul class="colist">
<li data-value="1">The cluster name</li>
<li data-value="2">Environment variables to override the settings in the <code>tangosol-coherence-override.xml</code>.</li>
<li data-value="3">The <code>/config</code> path containing the xml files</li>
<li data-value="4">The federation port definitions</li>
<li data-value="5">The annotations to instruct OCI where to create the load balancer</li>
<li data-value="6">The cache config and override files</li>
<li data-value="7">The image and pull secrets</li>
</ul>
<p><strong>Secondary Cluster</strong></p>

<p>The following yaml is for the secondary cluster.</p>

<markup
lang="yaml"

>#
# Copyright (c) 2021, 2025 Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
# Federation Example
# Primary cluster in an Active/Active topology
apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: primary-cluster <span class="conum" data-value="1" />
spec:
  jvm:
    classpath:
      - /config
      - /u01/oracle/oracle_home/coherence/lib/coherence.jar
  env: <span class="conum" data-value="2" />
    - name: "PRIMARY_CLUSTER"
      value: "primary-cluster"
    - name: "PRIMARY_CLUSTER_ADDRESS"
      value: ""
    - name: "PRIMARY_CLUSTER_PORT"
      value: "40000"
    - name: "SECONDARY_CLUSTER"
      value: "secondary-cluster"
    - name: "SECONDARY_CLUSTER_ADDRESS"
      value: ""
    - name: "SECONDARY_CLUSTER_PORT"
      value: "40000"
  secretVolumes: <span class="conum" data-value="3" />
    - mountPath: /config
      name: storage-config
  ports:
    - name: "federation" <span class="conum" data-value="4" />
      port: 40000
      protocol: TCP
      service:
        port: 40000
        type: LoadBalancer
        annotations: <span class="conum" data-value="5" />
          oci.oraclecloud.com/load-balancer-type: "nlb"
          oci-network-load-balancer.oraclecloud.com/internal: "true"
          oci-network-load-balancer.oraclecloud.com/subnet: "(Internal subnet OCID - REPLACE ME)"
          oci-network-load-balancer.oraclecloud.com/oci-network-security-groups: "(OCID of the NSG - REPLACE ME)"
    - name: management
  coherence:  <span class="conum" data-value="6" />
    cacheConfig: /config/storage-cache-config.xml
    overrideConfig: /config/tangosol-coherence-override.xml
    logLevel: 9
  image: container-registry.oracle.com/middleware/coherence:14.1.2.0.0 <span class="conum" data-value="7" />
  imagePullSecrets:
    - name: ocr-pull-secret
  replicas: 3</markup>

</div>
</div>

<h3 id="deploy">Deploy the Example</h3>
<div class="section">

<h4 id="update-ocid">1. Update the OCID&#8217;s for the federation service</h4>
<div class="section">
<p>Before we apply the yaml for the primary <strong>AND</strong> secondary clusters, we must update the annotations in the
definition of the federation port to ensure we create the load balancer in the correct subnet.</p>

<p>The existing entry looks like the following:</p>

<markup
lang="yaml"

>  ports:
    - name: "federation"
      port: 40000
      protocol: TCP
      service:
        port: 40000
        type: LoadBalancer
        annotations:
          oci.oraclecloud.com/load-balancer-type: "nlb"
          oci-network-load-balancer.oraclecloud.com/internal: "true"
          oci-network-load-balancer.oraclecloud.com/subnet: "(Internal subnet OCID - REPLACE ME)"
          oci-network-load-balancer.oraclecloud.com/oci-network-security-groups: "(OCID of the NSG - REPLACE ME)"</markup>

<p>The following values should be changed:</p>

<p><code>oci-network-load-balancer.oraclecloud.com/subnet</code></p>

<p>This value should be set to the OCID of the subnet created for the network load balancer.
In our case in Sydney region (c1) the subnet <code>int_lb-xkhplh</code> is used for the load balancer. (Virtual Cloud Networks&#8594; Subnets)</p>



<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/images/subnet-ocid.png" alt="Subnet OCI" />
</v-card-text>
</v-card>

<p><code>oci-network-load-balancer.oraclecloud.com/oci-network-security-groups</code></p>

<p>This value should be set to the OCID of the network security group for the above subnet that has the security rules for ingress/egress defined previously.
In our case in Sydney region (c1) the network security group <code>int_lb-xkhplh</code> is used to apply the security rules we discussed previously. (Virtual Cloud Networks&#8594; Security &#8594; Network Security Groups)</p>



<v-card>
<v-card-text class="overflow-y-hidden" style="text-align:center">
<img src="./images/images/nsg-ocid.png" alt="NSG OCI" />
</v-card-text>
</v-card>

<div class="admonition note">
<p class="admonition-inline">You should switch to the secondary region, Melbourne (c3), and change the values in <code>secondary-cluster.yaml</code> to reflect the OCID values in this region.</p>
</div>
</div>

<h4 id="install-primary">2. Install the Primary Coherence cluster</h4>
<div class="section">
<p>Ensure you are in the <code>examples/federation</code> directory to run the example. This example uses the yaml files <code>src/main/yaml/primary-cluster.yaml</code> and <code>src/main/yaml/secondary-cluster.yaml</code>, which you modified in the previous step.</p>

<p>Ensure your Kubernetes context is set to the <strong>primary</strong> cluster, which is <code>c1</code> or Sydney in our case, and run the following
commands to create the primary cluster and load balancer for the federation port:</p>

<markup
lang="bash"

>kubectx c1
Switched to context "c1".</markup>

<markup
lang="bash"

>kubectl -n coherence-example apply -f src/main/yaml/primary-cluster.yaml

coherence.coherence.oracle.com/primary-cluster created</markup>

<p>Issue the following command to view the pods:</p>

<markup
lang="bash"

> kubectl -n coherence-example get pods

NAME                READY   STATUS    RESTARTS   AGE
primary-cluster-0   1/1     Running   0          2m
primary-cluster-1   1/1     Running   0          2m
primary-cluster-2   1/1     Running   0          2m</markup>

<p>Once all the pods are ready, view the services and wait for an external IP for the <code>primary-cluster-federation</code> service
as we have defined this as a load balancer.
Once this has been populated, update the <code>PRIMARY_CLUSTER_ADDRESS</code> in both <code>yaml</code> files to this value.
In our example this is <code>10.1.2.22</code>.</p>

<markup
lang="bash"

>kubectl get svc -n coherence-example

NAME                         TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)                                                AGE
primary-cluster-federation   LoadBalancer   10.101.77.23    10.1.2.22     40000:30876/TCP                                        60s
primary-cluster-management   ClusterIP      10.101.202.62   &lt;none&gt;        30000/TCP                                              60s
primary-cluster-sts          ClusterIP      None            &lt;none&gt;        7/TCP,7575/TCP,7574/TCP,6676/TCP,40000/TCP,30000/TCP   60s
primary-cluster-wka          ClusterIP      None            &lt;none&gt;        7/TCP,7575/TCP,7574/TCP,6676/TCP                       60s</markup>

<markup
lang="bash"

>  env:
    - name: "PRIMARY_CLUSTER"
      value: "primary-cluster"
    - name: "PRIMARY_CLUSTER_ADDRESS"
      value: "10.1.2.22"</markup>

</div>

<h4 id="install-secondary">3. Install the Secondary Coherence cluster</h4>
<div class="section">
<p>Ensure your Kubernetes context is set to the <strong>secondary</strong> cluster, which is <code>c3</code> or Melbourne in our case, and run the following
commands to create the secondary cluster and load balancer for federation port:</p>

<markup
lang="bash"

>kubectx c3
Switched to context "c3".</markup>

<markup
lang="bash"

>kubectl -n coherence-example apply -f src/main/yaml/secondary-cluster.yaml

coherence.coherence.oracle.com/primary-cluster created</markup>

<p>Issue the following command to view the pods:</p>

<markup
lang="bash"

>kubectl -n coherence-example get pods

NAME                  READY   STATUS    RESTARTS   AGE
secondary-cluster-0   1/1     Running   0          2m
secondary-cluster-1   1/1     Running   0          2m
secondary-cluster-2   1/1     Running   0          2m</markup>

<p>Wait for an external-ip for the <code>secondary-cluster-federation</code> service.</p>

<markup
lang="bash"

>kubectl get svc -n coherence-example

NAME                           TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                                                AGE
secondary-cluster-federation   LoadBalancer   10.103.205.128   10.3.2.14     40000:30186/TCP                                        43s
secondary-cluster-management   ClusterIP      10.103.96.216    &lt;none&gt;        30000/TCP                                              43s
secondary-cluster-sts          ClusterIP      None             &lt;none&gt;        7/TCP,7575/TCP,7574/TCP,6676/TCP,40000/TCP,30000/TCP   43s
secondary-cluster-wka          ClusterIP      None             &lt;none&gt;        7/TCP,7575/TCP,7574/TCP,6676/TCP                       43s</markup>

<p>Once this has been populated then update the <code>SECONDARY_CLUSTER_ADDRESS</code> in both <code>yaml</code> files to this value.
In our example this is <code>10.3.2.14</code>.</p>

<markup
lang="bash"

>    - name: "SECONDARY_CLUSTER"
      value: "secondary-cluster"
    - name: "SECONDARY_CLUSTER_ADDRESS"
      value: "10.3.2.14"</markup>

</div>

<h4 id="re-apply">4. Re-apply the yaml for both clusters</h4>
<div class="section">
<p>Run the following to re-apply the yaml in both regions to correctly set the load balancer addresses for each cluster.
Making a change to the yaml will be carried out in a safe manner via a rolling redeploy of the stateful set.</p>

<p><strong>Primary Cluster</strong></p>

<markup
lang="bash"

>kubectx c1
Switched to context "c1".

kubectl -n coherence-example apply -f src/main/yaml/primary-cluster.yaml</markup>

<p>Issue the following command and wait until the <code>PHASE</code> is ready which indicates the rolling upgrade has been completed. Once it is ready, you
can use <code>CTRL-C</code> to exit the command.</p>

<markup
lang="bash"

>kubectl get coh -n coherence-example -w

NAME              CLUSTER           ROLE              REPLICAS   READY   PHASE
primary-cluster   primary-cluster   primary-cluster   3          2       RollingUpgrade
primary-cluster   primary-cluster   primary-cluster   3          3       RollingUpgrade
primary-cluster   primary-cluster   primary-cluster   3          2       RollingUpgrade
primary-cluster   primary-cluster   primary-cluster   3          3       RollingUpgrade
primary-cluster   primary-cluster   primary-cluster   3          2       RollingUpgrade
primary-cluster   primary-cluster   primary-cluster   3          3       Ready</markup>

<p><strong>Secondary Cluster</strong></p>

<p>Change to the secondary cluster context and run the following to restart secondary cluster using a rolling restart.</p>

<markup
lang="bash"

>kubectx c3
Switched to context "c3".

kubectl -n coherence-example apply -f src/main/yaml/secondary-cluster.yaml</markup>

<p>Wait for the cluster to be ready.</p>

<markup
lang="bash"

>kubectl get coh -n coherence-example -w

NAME                CLUSTER             ROLE                REPLICAS   READY   PHASE
secondary-cluster   secondary-cluster   secondary-cluster   3          2       RollingUpgrade
secondary-cluster   secondary-cluster   secondary-cluster   3          3       RollingUpgrade
secondary-cluster   secondary-cluster   secondary-cluster   3          2       RollingUpgrade
secondary-cluster   secondary-cluster   secondary-cluster   3          3       RollingUpgrade
secondary-cluster   secondary-cluster   secondary-cluster   3          2       RollingUpgrade
secondary-cluster   secondary-cluster   secondary-cluster   3          3       Ready</markup>

</div>

<h4 id="check-federation">5. Check federation status</h4>
<div class="section">
<p>When both clusters have completed the rolling upgrade, open two terminals, to view the federation status of the primary and secondary clusters.</p>

<p>This command uses the Coherence CLI, which is embedded in the operator. For more information on the CLI, see <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-cli">here</a>.</p>

<p>Run the following command using <code>--context c1</code> to view the Federation status. Because we have not yet added data, we should see
it showing the <code>STATES</code> as <code>[INITIAL]</code>.  The <code>-W</code> option will continuously display the federation output until you use <code>CTRL-C</code> to exit.</p>

<div class="admonition note">
<p class="admonition-inline">From now on when we use <code>--context</code> parameter, you should replace this with the appropriate context for your environment.</p>
</div>
<markup
lang="bash"

>kubectl --context c1 exec -it -n coherence-example primary-cluster-0 -- /coherence-operator/utils/cohctl get federation all -o wide -W

2025-03-18 06:20:01.645984469 +0000 UTC m=+20.337976443
Using cluster connection 'default' from current context.

SERVICE         OUTGOING           MEMBERS  STATES     DATA SENT  MSG SENT  REC SENT  CURR AVG BWIDTH  AVG APPLY  AVG ROUND TRIP  AVG BACKLOG DELAY  REPLICATE  PARTITIONS  ERRORS  UNACKED
FederatedCache  secondary-cluster        3  [INITIAL]       0 MB         0         0          0.0Mbps        0ms             0ms                0ms      0.00%           0       0        0</markup>

<p>In a second window, change the context and run the second command:</p>

<markup
lang="bash"

>kubectl --context c3 exec -it -n coherence-example secondary-cluster-0 -- /coherence-operator/utils/cohctl get federation all -o wide -W

2025-03-18 06:20:46.142911529 +0000 UTC m=+10.337433424
Using cluster connection 'default' from current context.

SERVICE         OUTGOING         MEMBERS  STATES     DATA SENT  MSG SENT  REC SENT  CURR AVG BWIDTH  AVG APPLY  AVG ROUND TRIP  AVG BACKLOG DELAY  REPLICATE  PARTITIONS  ERRORS  UNACKED
FederatedCache  primary-cluster        3  [INITIAL]       0 MB         0         0          0.0Mbps        0ms             0ms                0ms      0.00%           0       0        0</markup>

</div>
</div>

<h3 id="run">Run the Example</h3>
<div class="section">

<h4 id="add-data-p">1. Add Data in the primary cluster</h4>
<div class="section">
<p>In a separate terminal, on the primary-cluster, run the following command to run the console application to add data.</p>

<markup
lang="bash"

>kubectl exec --context c1 -it -n coherence-example primary-cluster-0 -- /coherence-operator/utils/runner console</markup>

<p>At the <code>Map:</code> prompt, type the following to create a new cache called <code>test</code>.</p>

<markup
lang="bash"

>cache test</markup>

<p>Then, issue the command <code>put key1 primary</code> to add an entry with key=key1, and value=primary</p>

<markup
lang="bash"
title="Output"
>Map (test): put key1 primary
null

Map (test): list
key1 = primary</markup>

<div class="admonition tip">
<p class="admonition-inline">Type <code>bye</code> to exit the console.</p>
</div>
<p>Once you run the command, switch back to the terminals with the <code>federation all</code> commands running for the primary cluster, and you
should see data has being sent from the primary to secondary cluster.</p>

<markup
lang="bash"
title="Primary Cluster"
>2025-03-18 06:25:22.649699832 +0000 UTC m=+341.341691816
Using cluster connection 'default' from current context.

SERVICE         OUTGOING           MEMBERS  STATES          DATA SENT  MSG SENT  REC SENT  CURR AVG BWIDTH  AVG APPLY  AVG ROUND TRIP  AVG BACKLOG DELAY  REPLICATE  PARTITIONS  ERRORS  UNACKED
FederatedCache  secondary-cluster        3  [INITIAL IDLE]       0 MB         1         1          0.0Mbps       20ms             7ms               88ms      0.00%           0       0        0</markup>

<p>You should also see something similar on the secondary cluster, which shows data being received.</p>

<markup
lang="bash"
title="Secondary Cluster"
>2025-03-18 06:28:07.5999651 +0000 UTC m=+451.794486974
Using cluster connection 'default' from current context.

SERVICE         OUTGOING         MEMBERS  STATES     DATA SENT  MSG SENT  REC SENT  CURR AVG BWIDTH  AVG APPLY  AVG ROUND TRIP  AVG BACKLOG DELAY  REPLICATE  PARTITIONS  ERRORS  UNACKED
FederatedCache  primary-cluster        3  [INITIAL]       0 MB         0         0          0.0Mbps        0ms             0ms                0ms      0.00%           0       0        0

SERVICE         INCOMING         MEMBERS RECEIVING  DATA REC  MSG REC  REC REC  AVG APPLY  AVG BACKLOG DELAY
FederatedCache  primary-cluster                  1      0 MB        1        1       60ms              265ms</markup>

</div>

<h4 id="validate-data-s">2. Validate data in the secondary cluster</h4>
<div class="section">
<p>Run the following command against the secondary cluster to validate the data has been received:</p>

<markup
lang="bash"

>kubectl exec --context c3 -it -n coherence-example secondary-cluster-0 -- /coherence-operator/utils/runner console</markup>

<p>At the <code>Map:</code> prompt, type <code>cache test</code> and then <code>list</code> to see the data that was transfered.</p>

<markup
lang="bash"

>Map (test): list
key1 = primary</markup>

<div class="admonition note">
<p class="admonition-inline">If you do not see any data, then you should check the pod logs to see if there are any errors.</p>
</div>
<markup
lang="bash"

>size</markup>

</div>

<h4 id="add-data-s">3. Add Data in the secondary cluster</h4>
<div class="section">
<div class="admonition tip">
<p class="admonition-inline">Since we have configured federation to be active-active, any changes in either cluster will be replicated to the other.</p>
</div>
<p>Without leaving the console, run the command <code>put key2 secondary</code> and then <code>list</code>:</p>

<markup
lang="bash"

>Map (test): put key2 secondary
null

Map (test): list
key1 = primary
key2 = secondary</markup>

</div>

<h4 id="validate-data-p">4. Validate data in the primary cluster</h4>
<div class="section">
<p>Switch back to the primary cluster, and start the console to validate that the data has been receieved:</p>

<markup
lang="bash"

>kubectl exec --context c1 -it -n coherence-example primary-cluster-0 -- /coherence-operator/utils/runner console</markup>

<p>At the <code>Map:</code> prompt, type <code>cache test</code> and the <code>list</code></p>

<markup
lang="bash"

>Map (test): list
key2 = secondary
key1 = primary</markup>

<p>Use the following command to add 100,000 entries to the cache with random data:</p>

<markup
lang="bash"

>bulkput 100000 100 0 100</markup>

<p>If you view the <code>federation all</code> commands that are still running you should see the data being federated.</p>

<markup
lang="bash"
title="Primary Cluster"
>2025-03-18 06:34:04.272571513 +0000 UTC m=+862.964563476
Using cluster connection 'default' from current context.

SERVICE         OUTGOING           MEMBERS  STATES     DATA SENT  MSG SENT  REC SENT  CURR AVG BWIDTH  AVG APPLY  AVG ROUND TRIP  AVG BACKLOG DELAY  REPLICATE  PARTITIONS  ERRORS  UNACKED
FederatedCache  secondary-cluster        3  [SENDING]       3 MB    12,128    14,973          1.2Mbps      278ms            25ms              132ms      0.00%           0       0        0

SERVICE         INCOMING           MEMBERS RECEIVING  DATA REC  MSG REC  REC REC  AVG APPLY  AVG BACKLOG DELAY
FederatedCache  secondary-cluster                  1      0 MB        1        1       20ms              285ms</markup>

<div class="admonition tip">
<p class="admonition-inline">Run the console against the secondary cluster and validate that the cache size is 100,000.</p>
</div>
</div>
</div>

<h3 id="cleanup">Cleaning up</h3>
<div class="section">
<p>Use the following commands to delete the primary and secondary clusters:</p>

<markup
lang="bash"

>kubectx c1
Switched to context "c1".

kubectl -n coherence-example delete -f src/main/yaml/primary-cluster.yaml

kubectx c3
Switched to context "c3".

kubectl -n coherence-example delete -f src/main/yaml/secondary-cluster.yaml</markup>

<p>Uninstall the Coherence operator using the undeploy commands for whichever method you chose to install it.</p>

</div>
</div>
</doc-view>
