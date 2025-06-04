<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Coherence Operator Installation</dt>
<dd slot="desc"><p>The Coherence Operator is available as an image from the GitHub container registry
<code>container-registry.oracle.com/middleware/coherence-operator:3.5.1</code> that can
easily be installed into a Kubernetes cluster.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_coherence_operator_installation">Coherence Operator Installation</h2>
<div class="section">
<p><strong>Contents</strong></p>

<ul class="ulist">
<li>
<p><router-link to="#prereq" @click.native="this.scrollFix('#prereq')">Prerequisites before installation</router-link></p>

</li>
<li>
<p><router-link to="#ha" @click.native="this.scrollFix('#ha')">Operator High Availability</router-link></p>

</li>
<li>
<p><router-link to="#images" @click.native="this.scrollFix('#images')">Coherence Operator Images</router-link></p>

</li>
<li>
<p><router-link to="#scope" @click.native="this.scrollFix('#scope')">Operator Scope - monitoring all or a fixed set of namespaces</router-link></p>

</li>
<li>
<p><router-link to="/docs/installation/090_tls_cipher">Configure TLS Cipher Suites</router-link></p>

</li>
<li>
<p><router-link to="/docs/installation/100_fips">FIPS Compliance</router-link></p>

</li>
</ul>

<h3 id="prereq">Prerequisites</h3>
<div class="section">
<p>The prerequisites apply to all installation methods.</p>

<ul class="ulist">
<li>
<p>Access to Oracle Coherence Operator images.</p>

</li>
<li>
<p>Access to a Kubernetes cluster. The Operator test pipeline is run using against all the currently supported Kubernetes versions.</p>

</li>
<li>
<p>A Coherence application image using Coherence version 12.2.1.3 or later. Note that some functionality (e.g. metrics) is only
available in Coherence 12.2.1.4 and later.</p>

</li>
</ul>
<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Istio (or similar service meshes)</p>

<p>When installing the Operator and Coherence into Kubernetes cluster that use Istio or similar meshes there are a
number of pre-requisites that must be understood.
See the <router-link to="/examples/400_Istio/README">Istio example</router-link> for more details.</p>
</p>
</div>
</div>

<h3 id="_installation_options">Installation Options</h3>
<div class="section">
<p>There are a number of ways to install the Coherence Operator.</p>

<ul class="ulist">
<li>
<p><router-link to="/docs/installation/011_install_manifests">Install using the yaml manifest file</router-link></p>

</li>
<li>
<p><router-link to="/docs/installation/012_install_helm">Install using Helm</router-link></p>

</li>
<li>
<p><router-link to="/docs/installation/013_install_kustomize">Install using Kustomize</router-link></p>

</li>
<li>
<p><router-link to="/docs/installation/014_install_openshift">Install on OpenShift</router-link></p>

</li>
<li>
<p><router-link to="/docs/installation/015_install_olm">Install using the Operator Lifecycle Manager (OLM)</router-link></p>

</li>
<li>
<p><router-link to="/docs/installation/016_install_tanzu">Install on VMWare Tanzu</router-link></p>

</li>
</ul>
</div>

<h3 id="ha">High Availability</h3>
<div class="section">
<p>The Coherence Operator runs in HA mode by default. The <code>Deployment</code> created by the installation will have a replica count of 3.
In reduced capacity Kubernetes clusters, for example, local laptop development and test, the replica count can be reduced. It is recommended to leave the default of 3 for production environments.
Instructions on how to change the replica count for the different install methods are included below.</p>

<p>The Coherence Operator runs a REST server that the Coherence cluster members will query to discover the site and rack names that should be used by Coherence. If the Coherence Operator is not running when a Coherence Pod starts, then the Coherence member in that Pod will be unable to properly configure its site and rack names, possibly leading to data distribution that is not safely distributed over sites. In production, and in Kubernetes clusters that are spread over multiple availability zones and failure domains, it is important to run the Operator in HA mode.</p>

<p>The Operator yaml files and Helm chart include a default Pod scheduling configuration that uses anti-affinity to distribute the three replicas onto nodes that have different <code>topology.kubernetes.io/zone</code> labels. This label is a standard Kubernetes label used to describe the zone the node is running in, and is typically applied by Kubernetes cloud vendors.</p>

</div>

<h3 id="_notes">Notes</h3>
<div class="section">
<div class="admonition note">
<p class="admonition-inline">Installing the Coherence Operator using the methods below will create a number of <code>ClusterRole</code> RBAC resources.
Some corporate security policies do not like to give cluster wide roles to third-party products.
To help in this situation the operator can be installed without cluster roles, but with caveats
(see the <router-link to="#docs/installation/09_RBAC.adoc" @click.native="this.scrollFix('#docs/installation/09_RBAC.adoc')">RBAC</router-link> documentation) for more details.</p>
</div>
<div class="admonition note">
<p class="admonition-inline">OpenShift - the Coherence Operator works without modification on OpenShift, but some versions
of the Coherence images will not work out of the box.
See the <router-link to="#docs/installation/06_openshift.adoc" @click.native="this.scrollFix('#docs/installation/06_openshift.adoc')">OpensShift</router-link> section of the documentation that explains how to
run Coherence clusters with the Operator on OpenShift.</p>
</div>
<div class="admonition note">
<p class="admonition-inline">Whilst Coherence works out of the box on many Kubernetes installations, some Kubernetes
installations may configure iptables in a way that causes Coherence to fail to create clusters.
See the <router-link to="#docs/installation/08_networking.adoc" @click.native="this.scrollFix('#docs/installation/08_networking.adoc')">O/S Network Configuration</router-link> section of the documentation
for more details if you have well-known-address issues when Pods attempt to form a cluster.</p>
</div>
</div>
</div>

<h2 id="images">Coherence Operator Images</h2>
<div class="section">
<p>The Coherence Operator uses a single image, the Operator also runs as an init-container in the Coherence cluster Pods.</p>

<ul class="ulist">
<li>
<p><code>container-registry.oracle.com/middleware/coherence-operator:3.5.1</code> - The Operator image.</p>

</li>
</ul>
<p>If no image is specified in the <code>Coherence</code> yaml, then the default Coherence image will also be used,</p>

<ul class="ulist">
<li>
<p><code>container-registry.oracle.com/middleware/coherence-ce:14.1.2-0-2</code> - The default Coherence image.</p>

</li>
</ul>
<p>If using a private image registry then these images will all need to be pushed to that registry for the Operator to work. The default Coherence image may be omitted if all Coherence applications will use custom Coherence images.</p>

</div>

<h2 id="scope">Operator Scope</h2>
<div class="section">
<p>The recommended way to install the Coherence Operator is to install a single instance of the operator into a namespace
and where it will then control <code>Coherence</code> resources in all namespaces across the Kubernetes cluster.
Alternatively it may be configured to watch a sub-set of namespaces by setting the <code>WATCH_NAMESPACE</code> environment variable.
The watch namespace(s) does not have to include the installation namespace.</p>

<div class="admonition caution">
<p class="admonition-textlabel">Caution</p>
<p ><p>In theory, it is possible to install multiple instances of the Coherence Operator into different namespaces, where
each instance monitors a different set of namespaces. There are a number of potential issues with this approach, so
it is not recommended.</p>

<ul class="ulist">
<li>
<p>Only one version of a CRD can be installed - There is currently only a single version of the CRD, but different
releases of the Operator may use slightly different specs of this CRD version, for example
a new Operator release may introduce extra fields not in the previous releases.
As the CRD version is fixed at <code>v1</code> there is no guarantee which CRD version has actually installed, which could lead to
subtle issues.</p>

</li>
<li>
<p>The operator creates and installs defaulting and validating web-hooks. A web-hook is associated to a CRD resource so
installing multiple web-hooks for the same resource may lead to issues. If an operator is uninstalled, but the web-hook
configuration remains, then Kubernetes will not accept modifications to resources of that type as it will be
unable to contact the web-hook.</p>

</li>
</ul>
<p>It is possible to run the Operator without web-hooks, but this has its own
caveats see the <router-link to="#docs/installation/07_webhooks.adoc" @click.native="this.scrollFix('#docs/installation/07_webhooks.adoc')">Web Hooks</router-link> documentation for how to do this.</p>
</p>
</div>
<div class="admonition important">
<p class="admonition-textlabel">Important</p>
<p ><p>If multiple instances of the Operator are installed, where they are monitoring the same namespaces, this can cause issues.
For example, when a <code>Coherence</code> resource is then changed, all the Operator deployments will receive the same events
from Etcd and try to apply the same changes. Sometimes this may work, sometimes there may be errors, for example multiple
Operators trying to remove finalizers and delete a Coherence cluster.</p>
</p>
</div>
</div>
</doc-view>
