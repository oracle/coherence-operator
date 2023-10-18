<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Coherence Operator Installation</dt>
<dd slot="desc"><p>The Coherence Operator is available as an image from the GitHub container registry <code>ghcr.io/oracle/coherence-operator:3.3.1</code> that can
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
<p>Installation Options</p>
<ul class="ulist">
<li>
<p><router-link to="#manifest" @click.native="this.scrollFix('#manifest')">Simple installation using Kubectl</router-link></p>

</li>
<li>
<p><router-link to="#helm" @click.native="this.scrollFix('#helm')">Install the Helm chart</router-link></p>
<ul class="ulist">
<li>
<p><router-link to="#helm-operator-image" @click.native="this.scrollFix('#helm-operator-image')">Set the Operator Image</router-link></p>

</li>
<li>
<p><router-link to="#helm-pull-secrets" @click.native="this.scrollFix('#helm-pull-secrets')">Image Pull Secrets</router-link></p>

</li>
<li>
<p><router-link to="#helm-watch-ns" @click.native="this.scrollFix('#helm-watch-ns')">Set the Watch Namespaces</router-link></p>

</li>
<li>
<p><router-link to="#helm-sec-context" @click.native="this.scrollFix('#helm-sec-context')">Install the Operator with a Security Context</router-link></p>

</li>
<li>
<p><router-link to="#helm-labels" @click.native="this.scrollFix('#helm-labels')">Set Additional Labels</router-link></p>

</li>
<li>
<p><router-link to="#helm-annotations" @click.native="this.scrollFix('#helm-annotations')">Set Additional Annotations</router-link></p>

</li>
<li>
<p><router-link to="#helm-uninstall" @click.native="this.scrollFix('#helm-uninstall')">Uninstall the Coherence Operator Helm chart</router-link></p>

</li>
</ul>
</li>
<li>
<p><router-link to="#kubectl" @click.native="this.scrollFix('#kubectl')">Kubectl with Kustomize</router-link></p>

</li>
<li>
<p><router-link to="#tanzu" @click.native="this.scrollFix('#tanzu')">VMWare Tanzu Package (kapp-controller)</router-link></p>

</li>
</ul>
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
<p>Access to a Kubernetes v1.19.0+ cluster. The Operator test pipeline is run using Kubernetes versions v1.19 upto v1.26</p>

</li>
<li>
<p>A Coherence application image using Coherence version 12.2.1.3 or later. Note that some functionality (e.g. metrics) is only
available in Coherence 12.2.1.4 and later.</p>

</li>
</ul>
<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>ARM Support: As of version 3.2.0, the Coherence Operator is build as a multi-architecture image that supports running in Kubernetes on both Linux/amd64 and Linux/arm64. The prerequisite is that the Coherence application image used has been built to support ARM.</p>
</p>
</div>
<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Istio (or similar service meshes)</p>

<p>When installing the Operator and Coherence into Kubernetes cluster that use Istio or similar meshes there are a
number of pre-requisites that must be understood.
See the <router-link to="/examples/400_Istio/README">Istio example</router-link> for more details.</p>
</p>
</div>
<p>There are a number of ways to install the Coherence Operator documented below:</p>

<ul class="ulist">
<li>
<p><router-link to="#manifest" @click.native="this.scrollFix('#manifest')">Simple installation using Kubectl</router-link></p>

</li>
<li>
<p><router-link to="#helm" @click.native="this.scrollFix('#helm')">Install the Helm chart</router-link></p>

</li>
<li>
<p><router-link to="#kubectl" @click.native="this.scrollFix('#kubectl')">Kubectl with Kustomize</router-link></p>

</li>
<li>
<p><router-link to="#tanzu" @click.native="this.scrollFix('#tanzu')">VMWare Tanzu Package (kapp-controller)</router-link></p>

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
(see the <router-link to="/docs/installation/09_RBAC">RBAC</router-link> documentation) for more details.</p>
</div>
<div class="admonition note">
<p class="admonition-inline">OpenShift - the Coherence Operator works without modification on OpenShift, but some versions
of the Coherence images will not work out of the box.
See the <router-link to="/docs/installation/06_openshift">OpensShift</router-link> section of the documentation that explains how to
run Coherence clusters with the Operator on OpenShift.</p>
</div>
<div class="admonition note">
<p class="admonition-inline">Whilst Coherence works out of the box on many Kubernetes installations, some Kubernetes
installations may configure iptables in a way that causes Coherence to fail to create clusters.
See the <router-link to="/docs/installation/08_networking">O/S Network Configuration</router-link> section of the documentation
for more details if you have well-known-address issues when Pods attempt to form a cluster.</p>
</div>
</div>
</div>

<h2 id="images">Coherence Operator Images</h2>
<div class="section">
<p>The Coherence Operator uses a single image, the Operator also runs as an init-container in the Coherence cluster Pods.</p>

<ul class="ulist">
<li>
<p><code>ghcr.io/oracle/coherence-operator:3.3.1</code> - The Operator image.</p>

</li>
</ul>
<p>If no image is specified in the <code>Coherence</code> yaml, then the default Coherence image will also be used,</p>

<ul class="ulist">
<li>
<p><code>ghcr.io/oracle/coherence-ce:22.06.6</code> - The default Coherence image.</p>

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
caveats see the <router-link to="/docs/installation/07_webhooks">Web Hooks</router-link> documentation for how to do this.</p>
</p>
</div>
<div class="admonition important">
<p class="admonition-textlabel">Important</p>
<p ><p>If multiple instance of the Operator are installed, where they are monitoring the same namespaces, this can cause issues.
For example, when a <code>Coherence</code> resource is then changed, all the Operator deployments will receive the same events
from Etcd and try to apply the same changes. Sometimes this may work, sometimes there may be errors, for example multiple
Operators trying to remove finalizers and delete a Coherence cluster.</p>
</p>
</div>
</div>

<h2 id="manifest">Default Install with Kubectl</h2>
<div class="section">
<p>If you want the default Coherence Operator installation then the simplest solution is use <code>kubectl</code> to apply the manifests from the Operator release.</p>

<markup
lang="bash"

>kubectl apply -f https://github.com/oracle/coherence-operator/releases/download/v3.3.1/coherence-operator.yaml</markup>

<p>This will create a namespace called <code>coherence</code> and install the Operator into it along with all the required <code>ClusterRole</code> and <code>RoleBinding</code> resources. The <code>coherence</code> namespace can be changed by downloading and editing the yaml file.</p>

<div class="admonition note">
<p class="admonition-inline">Because the <code>coherence-operator.yaml</code> manifest also creates the namespace, the corresponding <code>kubectl delete</code> command will <em>remove the namespace and everything deployed to it</em>! If you do not want this behaviour you should edit the <code>coherence-operator.yaml</code> to remove the namespace section from the start of the file.</p>
</div>
<p>Instead of using a hard coded version in the command above you can find the latest Operator version using <code>curl</code>:</p>

<markup
lang="bash"

>export VERSION=$(curl -s \
  https://api.github.com/repos/oracle/coherence-operator/releases/latest \
  | grep '"name": "v' \
  | cut -d '"' -f 4 \
  | cut -b 2-10)</markup>

<p>Then download with:</p>

<markup
lang="bash"

>kubectl apply -f https://github.com/oracle/coherence-operator/releases/download/${VERSION}/coherence-operator.yaml</markup>


<h3 id="_change_the_operator_replica_count">Change the Operator Replica Count</h3>
<div class="section">
<p>When installing with single manifest yaml file, the replica count can be changed by editing the yaml file itself to change the occurrence of <code>replicas: 3</code> in the manifest yaml to <code>replicas: 1</code></p>

<p>For example, this could be done using <code>sed</code></p>

<markup
lang="bash"

>sed -i -e 's/replicas: 3/replicas: 1/g' coherence-operator.yaml</markup>

<p>Or on MacOS, where <code>sed</code> is slightly different:</p>

<markup
lang="bash"

>sed -i '' -e 's/replicas: 3/replicas: 1/g' coherence-operator.yaml</markup>

</div>
</div>

<h2 id="helm">Installing With Helm</h2>
<div class="section">
<p>For more flexibility but the simplest way to install the Coherence Operator is to use the Helm chart.
This ensures that all the correct resources will be created in Kubernetes.</p>


<h3 id="_add_the_coherence_helm_repository">Add the Coherence Helm Repository</h3>
<div class="section">
<p>Add the <code>coherence</code> helm repository using the following commands:</p>

<markup
lang="bash"

>helm repo add coherence https://oracle.github.io/coherence-operator/charts

helm repo update</markup>

<div class="admonition note">
<p class="admonition-inline">To avoid confusion, the URL <code><a id="" title="" target="_blank" href="https://oracle.github.io/coherence-operator/charts">https://oracle.github.io/coherence-operator/charts</a></code> is a Helm repo, it is not a website you open in a browser. You may think we shouldn&#8217;t have to say this, but you&#8217;d be surprised.</p>
</div>
</div>

<h3 id="_install_the_coherence_operator_helm_chart">Install the Coherence Operator Helm chart</h3>
<div class="section">
<p>Once the Coherence Helm repo has been configured the Coherence Operator can be installed using a normal Helm 3
install command:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \      <span class="conum" data-value="1" />
    coherence \                    <span class="conum" data-value="2" />
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">where <code>&lt;namespace&gt;</code> is the namespace that the Coherence Operator will be installed into.</li>
<li data-value="2"><code>coherence</code> is the name of this Helm installation.</li>
</ul>
</div>

<h3 id="helm-operator-image">Set the Operator Image</h3>
<div class="section">
<p>The Helm chart uses a default Operator image from <code>ghcr.io/oracle/coherence-operator:3.3.1</code>.
If the image needs to be pulled from a different location (for example an internal registry) then there are two ways to override the default.
Either set the individual <code>image.registry</code>, <code>image.name</code> and <code>image.tag</code> values, or set the whole image name by setting the <code>image</code> value.</p>

<p>For example, if the Operator image has been deployed into a private registry named <code>foo.com</code> but
with the same image name <code>coherence-operator</code> and tag <code>3.3.1</code> as the default image,
then just the <code>image.registry</code> needs to be specified.</p>

<p>In the example below, the image used to run the Operator will be <code>foo.com/coherence-operator:3.3.1</code>.</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set image.registry=foo.com \
    coherence-operator \
    coherence/coherence-operator</markup>

<p>All three of the image parts can be specified individually using <code>--set</code> options.
In the example below, the image used to run the Operator will
be <code>foo.com/operator:1.2.3</code>.</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set image.registry=foo.com \
    --set image.name=operator \
    --set image.tag=1.2.3
    coherence-operator \
    coherence/coherence-operator</markup>

<p>Alternatively, the image can be set using a single <code>image</code> value.
For example, the command below will set the Operator image to <code>images.com/coherence-operator:0.1.2</code>.</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set image=images.com/coherence-operator:0.1.2 \
    coherence-operator \
    coherence/coherence-operator</markup>

</div>

<h3 id="helm-pull-secrets">Image Pull Secrets</h3>
<div class="section">
<p>If the image is to be pulled from a secure repository that requires credentials then the image pull secrets
can be specified.
See the Kubernetes documentation on <a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/">Pulling from a Private Registry</a>.</p>


<h4 id="_add_pull_secrets_using_a_values_file">Add Pull Secrets Using a Values File</h4>
<div class="section">
<p>Create a values file that specifies the secrets, for example the <code>private-repo-values.yaml</code> file below:</p>

<markup
lang="yaml"
title="private-repo-values.yaml"
>imagePullSecrets:
- name: registry-secrets</markup>

<p>Now use that file in the Helm install command:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    -f private-repo-values.yaml <span class="conum" data-value="1" />
    coherence-operator \
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">the <code>private-repo-values.yaml</code> values fle will be used by Helm to inject the settings into the Operator deployment</li>
</ul>
</div>

<h4 id="_add_pull_secrets_using_set">Add Pull Secrets Using --set</h4>
<div class="section">
<p>Although the <code>imagePullSecrets</code> field in the values file is an array of <code>name</code> to value pairs it is possible to set
these values with the normal Helm <code>--set</code> parameter.</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set imagePullSecrets[0].name=registry-secrets <span class="conum" data-value="1" />
    coherence-operator \
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">this creates the same imagePullSecrets as the values file above.</li>
</ul>
</div>
</div>

<h3 id="_change_the_operator_replica_count_2">Change the Operator Replica Count</h3>
<div class="section">
<p>To change the replica count when installing the Operator using Helm, the <code>replicas</code> value can be set.</p>

<p>For example, to change the replica count from 3 to 1, the <code>--set replicas=1</code> option can be used.</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set replicas=1
    coherence \
    coherence/coherence-operator</markup>

</div>

<h3 id="helm-watch-ns">Set the Watch Namespaces</h3>
<div class="section">
<p>To set the watch namespaces when installing with helm set the <code>watchNamespaces</code> value, for example:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set watchNamespaces=payments,catalog,customers \
    coherence-operator \
    coherence/coherence-operator</markup>

<p>The <code>payments</code>, <code>catalog</code> and <code>customers</code> namespaces will be watched by the Operator.</p>


<h4 id="_set_the_watch_namespace_to_the_operators_install_namespace">Set the Watch Namespace to the Operator&#8217;s Install Namespace</h4>
<div class="section">
<p>When installing the Operator using the Helm chart, there is a convenience value that can be set if the
Operator should only monitor the same namespace that it is installed into.
By setting the <code>onlySameNamespace</code> value to <code>true</code> the watch namespace will be set to the installation namespace.
If the <code>onlySameNamespace</code> value is set to <code>true</code> then any value set for the <code>watchNamespaces</code> value will be ignored.</p>

<p>For example, the command below will set <code>onlySameNamespace</code> to true, and the Operator will be installed into,
and only monitor the <code>coh-testing</code> namespace.</p>

<markup
lang="bash"

>helm install  \
    --namespace coh-testing \
    --set onlySameNamespace=true \
    coherence-operator \
    coherence/coherence-operator</markup>

<p>In the example below, the <code>onlySameNamespace</code> is set to true, so the Operator will be installed into,
and only monitor the <code>coh-testing</code> namespace. Even though the <code>watchNamespaces</code> value is set, it will be ignored.</p>

<markup
lang="bash"

>helm install  \
    --namespace coh-testing \
    --set watchNamespaces=payments,catalog,customers \
    --set onlySameNamespace=true \
    coherence-operator \
    coherence/coherence-operator</markup>

</div>
</div>

<h3 id="helm-sec-context">Install the Operator with a Security Context</h3>
<div class="section">
<p>The Operator container can be configured with a Pod <code>securityContext</code> or a container <code>securityContext</code>,
so that it runs as a non-root user.</p>

<p>This can be done using a values file:</p>

<p><strong>Set the Pod securityContext</strong></p>

<markup
lang="yaml"
title="security-values.yaml"
>podSecurityContext:
  runAsNonRoot: true
  runAsUser: 1000</markup>

<p><strong>Set the Container securityContext</strong></p>

<markup
lang="yaml"
title="security-values.yaml"
>securityContext:
  runAsNonRoot: true
  runAsUser: 1000</markup>

<p>Then the <code>security-values.yaml</code> values file above can be used in the Helm install command.</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --values security-values.yaml \
    coherence \
    coherence/coherence-operator</markup>

<p>Alternatively, the Pod or container <code>securityContext</code> values can be set on the command line as <code>--set</code> parameters:</p>

<p><strong>Set the Pod securityContext</strong></p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set podSecurityContext.runAsNonRoot=true \
    --set podSecurityContext.runAsUser=1000 \
    coherence \
    coherence/coherence-operator</markup>

<p><strong>Set the Container securityContext</strong></p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set securityContext.runAsNonRoot=true \
    --set securityContext.runAsUser=1000 \
    coherence \
    coherence/coherence-operator</markup>

</div>

<h3 id="helm-labels">Set Additional Labels</h3>
<div class="section">
<p>When installing the Operator with Helm, it is possible to set additional labels to be applied to the Operator Pods
and to the Operator Deployment.</p>


<h4 id="_adding_pod_labels">Adding Pod Labels</h4>
<div class="section">
<p>To add labels to the Operator Pods set the <code>labels</code> value, either on the command line using <code>--set</code> or in the values file.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Setting <code>labels</code> will only apply the additional labels to the Operator Pods, they will not be applied to any other resource created by the Helm chart.</p>
</p>
</div>
<p>For example, using the command line:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set labels.one=value-one \
    --set labels.two=value-two \
    coherence \
    coherence/coherence-operator</markup>

<p>The command above would add the following additional labels <code>one</code> and <code>two</code> to the Operator Pod as shown below:</p>

<markup
lang="yaml"

>apiVersion: v1
kind: Pod
metadata:
  name: coherence-operator
  labels:
    one: value-one
    two: value-two</markup>

<p>The same labels could also be specified in a values file:</p>

<markup

title="add-labels-values.yaml"
>labels:
  one: value-one
  two: value-two</markup>

</div>

<h4 id="_adding_deployment_labels">Adding Deployment Labels</h4>
<div class="section">
<p>To add labels to the Operator Deployment set the <code>deploymentLabels</code> value, either on the command line using <code>--set</code> or in the values file.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Setting <code>deploymentLabels</code> will only apply the additional labels to the Deployment, they will not be applied to any other resource created by the Helm chart.</p>
</p>
</div>
<p>For example, using the command line:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set deploymentLabels.one=value-one \
    --set deploymentLabels.two=value-two \
    coherence \
    coherence/coherence-operator</markup>

<p>The command above would add the following additional labels <code>one</code> and <code>two</code> to the Operator Pod as shown below:</p>

<markup
lang="yaml"

>apiVersion: apps/v1
kind: Deployment
metadata:
  name: coherence-operator
  labels:
    one: value-one
    two: value-two</markup>

<p>The same labels could also be specified in a values file:</p>

<markup

title="add-labels-values.yaml"
>deploymentLabels:
  one: value-one
  two: value-two</markup>

</div>
</div>

<h3 id="helm-annotations">Set Additional Annotations</h3>
<div class="section">
<p>When installing the Operator with Helm, it is possible to set additional annotations to be applied to the Operator Pods
and to the Operator Deployment.</p>


<h4 id="_adding_pod_annotations">Adding Pod Annotations</h4>
<div class="section">
<p>To add annotations to the Operator Pods set the <code>annotations</code> value, either on the command line using <code>--set</code> or in the values file.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Setting <code>annotations</code> will only apply the additional annotations to the Operator Pods, they will not be applied to any other resource created by the Helm chart.</p>
</p>
</div>
<p>For example, using the command line:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set annotations.one=value-one \
    --set annotations.two=value-two \
    coherence \
    coherence/coherence-operator</markup>

<p>The command above would add the following additional annotations <code>one</code> and <code>two</code> to the Operator Pod as shown below:</p>

<markup
lang="yaml"

>apiVersion: v1
kind: Pod
metadata:
  name: coherence-operator
  annotations:
    one: value-one
    two: value-two</markup>

<p>The same annotations could also be specified in a values file:</p>

<markup

title="add-annotations-values.yaml"
>annotations:
  one: value-one
  two: value-two</markup>

</div>

<h4 id="_adding_deployment_annotations">Adding Deployment Annotations</h4>
<div class="section">
<p>To add annotations to the Operator Deployment set the <code>deploymentAnnotations</code> value, either on the command line using <code>--set</code> or in the values file.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Setting <code>deploymentAnnotations</code> will only apply the additional annotations to the Deployment, they will not be applied to any other resource created by the Helm chart.</p>
</p>
</div>
<p>For example, using the command line:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set deploymentAnnotations.one=value-one \
    --set deploymentAnnotations.two=value-two \
    coherence \
    coherence/coherence-operator</markup>

<p>The command above would add the following additional annotations <code>one</code> and <code>two</code> to the Operator Pod as shown below:</p>

<markup
lang="yaml"

>apiVersion: apps/v1
kind: Deployment
metadata:
  name: coherence-operator
  annotations:
    one: value-one
    two: value-two</markup>

<p>The same annotations could also be specified in a values file:</p>

<markup

title="add-annotations-values.yaml"
>deploymentAnnotations:
  one: value-one
  two: value-two</markup>

</div>
</div>

<h3 id="helm-uninstall">Uninstall the Coherence Operator Helm chart</h3>
<div class="section">
<p>To uninstall the operator:</p>

<markup
lang="bash"

>helm delete coherence-operator --namespace &lt;namespace&gt;</markup>

</div>
</div>

<h2 id="kubectl">Install with Kubectl and Kustomize</h2>
<div class="section">
<p>If you want to use yaml directly to install the operator, with something like <code>kubectl</code>, you can use the manifest files
published with the GitHub release at this link:
<a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/releases/download/v3.3.1/coherence-operator-manifests.tar.gz">3.3.1 Manifests</a></p>

<p>These manifest files are for use with a tool called Kustomize, which is built into <code>kubectl</code>
see the documentation here: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/">https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/</a></p>

<p>Download the
<a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/releases/download/v3.3.1/coherence-operator-manifests.tar.gz">3.3.1 Manifests</a>
from the release page and unpack the file, which should produce a directory called <code>manifests</code> with a structure like this:</p>

<markup


>manifests
    default
        config.yaml
        kustomization.yaml
    manager
        kustomization.yaml
        manager.yaml
        service.yaml
    rbac
        coherence_editor_role.yaml
        coherence_viewer_role.yaml
        kustomization.yaml
        leader_election_role.yaml
        leader_election_role_binding.yaml
        role.yaml
        role_binding.yaml</markup>

<p>There are two ways to use these manifest files, either install using <code>kustomize</code> or generate the yaml and manually
install with <code>kubectl</code>.</p>

<div class="admonition note">
<p class="admonition-inline">All the commands below are run from a console in the <code>manifests/</code> directory from the extracted file above.</p>
</div>

<h3 id="_install_with_kustomize">Install with Kustomize</h3>
<div class="section">
<p>If you have Kustomize installed (or can install it from <a id="" title="" target="_blank" href="https://github.com/kubernetes-sigs/kustomize">https://github.com/kubernetes-sigs/kustomize</a>) you can use
Kustomize to configure the yaml and install.</p>


<h4 id="_change_the_operator_replica_count_3">Change the Operator Replica Count</h4>
<div class="section">
<p>To change the replica count using Kustomize a patch file needs to be applied.
The Operator manifests include a patch file, named <code>manager/single-replica-patch.yaml</code>, that changes the replica count from 3 to 1. This patch can be applied with the following Kustomize command.</p>

<markup
lang="bash"

>cd ./manager &amp;&amp; kustomize edit add patch \
  --kind Deployment --name controller-manager \
  --path single-replica-patch.yaml</markup>

</div>

<h4 id="_set_image_names">Set Image Names</h4>
<div class="section">
<p>If you need to use different iamge names from the defaults <code>kustomize</code> can be used to specify different names:</p>

<p>Change the name of the Operator image by running the command below, changing the image name to the registry and image name
that you are using for the Operator, for example if you have the images in a custom registry</p>

<markup
lang="bash"

>cd ./manager &amp;&amp; kustomize edit set image controller=myregistry/coherence-operator:3.3.1</markup>

<p>Change the name of the Operator image by running the command below, changing the image name to the registry and image name
that you are using for the Operator utilities image</p>

<markup
lang="bash"

>cd ./manager &amp;&amp; kustomize edit add configmap env-vars --from-literal OPERATOR_IMAGE=myregistry/coherence-operator:3.3.1</markup>

<p>Change the name of the default Coherence image. If you are always going to be deploying your own application images then this
does not need to change.</p>

<markup
lang="bash"

>cd ./manager &amp;&amp; $(GOBIN)/kustomize edit add configmap env-vars --from-literal COHERENCE_IMAGE=$(COHERENCE_IMAGE)</markup>

<p>Set the namespace to install into, the example below sets the namespace to <code>coherence-test</code>:</p>

<markup
lang="bash"

>cd ./default &amp;&amp; /kustomize edit set namespace coherence-test</markup>

</div>

<h4 id="_install">Install</h4>
<div class="section">
<p>The Operator requires a <code>Secret</code> for its web-hook certificates. This <code>Secret</code> needs to exist but can be empty.
The <code>Secret</code> must be in the same namespace that the Operator will be deployed to.
For example, if the Operator namespace is <code>coherence-test</code>, then the <code>Secret</code> can be created with this command:</p>

<markup
lang="bash"

>kubectl -n coherence-test create secret generic coherence-webhook-server-cert</markup>

<p>The Operator can now be installed by running the following command from the <code>manifests</code> directory:</p>

<markup
lang="bash"

>kustomize build ./default | kubectl apply -f -</markup>

</div>
</div>

<h3 id="_generate_yaml_install_with_kubectl">Generate Yaml - Install with Kubectl</h3>
<div class="section">
<p>Instead of using Kustomize to modify and install the Operator we can use <code>kubectl</code> to generate the yaml from the manifests.
You can then edit this yaml and manually deploy it with <code>kubectl</code>.</p>

<p>Run the following command from the <code>manifests</code> directory:</p>

<markup
lang="bash"

>kubectl create --dry-run -k default/ -o yaml &gt; operator.yaml</markup>

<p>This will create a file in the <code>manifests</code> directory called <code>operator.yaml</code> that contains all the yaml required
to install the Operator. You can then edit this yaml to change image names or add other settings.</p>

<p>The Operator can be installed using the generated yaml.</p>

<p>For example if the Operator is to be deployed to the <code>coherence-test</code> namespace:</p>

<markup
lang="bash"

>kubectl -n coherence-test create secret generic coherence-webhook-server-cert
kubectl -n coherence-test create -f operator.yaml</markup>

</div>
</div>

<h2 id="tanzu">Install as a VMWare Tanzu Package (Carvel kapp-controller)</h2>
<div class="section">
<p>If using <a id="" title="" target="_blank" href="https://tanzucommunityedition.io">VMWare Tanzu</a> the Coherence Operator can be installed as a package.
Under the covers, Tanzu uses the <a id="" title="" target="_blank" href="https://carvel.dev">Carvel</a> tool set to deploy packages.
The Carvel tools can be used outside Tanzu, so the Coherence Operator repo and package images could also be deployed
using a standalone Carvel <a id="" title="" target="_blank" href="https://carvel.dev/kapp-controller/">kapp-controller</a>.</p>

<p>The Coherence Operator release published two images required to deploy the Operator as a Tanzu package.</p>

<ul class="ulist">
<li>
<p><code>ghcr.io/oracle/coherence-operator-package:3.3.1</code> - the Coherence Operator package</p>

</li>
<li>
<p><code>ghcr.io/oracle/coherence-operator-repo:3.3.1</code> - the Coherence Operator repository</p>

</li>
</ul>

<h3 id="_install_the_coherence_repository">Install the Coherence Repository</h3>
<div class="section">
<p>The first step to deploy the Coherence Operator package in Tanzu is to add the repository.
This can be done using the Tanzu CLI.</p>

<markup
lang="bash"

>tanzu package repository add coherence-repo \
    --url ghcr.io/oracle/coherence-operator-repo:3.3.1 \
    --namespace coherence \
    --create-namespace</markup>

<p>The installed repositories can be listed using the CLI:</p>

<markup
lang="bash"

>tanzu package repository list --namespace coherence</markup>

<p>which should display something like the following</p>

<markup
lang="bash"

>NAME            REPOSITORY                              TAG  STATUS               DETAILS
coherence-repo  ghcr.io/oracle/coherence-operator-repo  1h   Reconcile succeeded</markup>

<p>The available packages in the Coherence repository can also be displayed using the CLI</p>

<markup
lang="bash"

>tanzu package available list --namespace coherence</markup>

<p>which should include the Operator package, <code>coherence-operator.oracle.github.com</code> something like the following</p>

<markup
lang="bash"

>NAME                                  DISPLAY-NAME               SHORT-DESCRIPTION                                             LATEST-VERSION
coherence-operator.oracle.github.com  Oracle Coherence Operator  A Kubernetes operator for managing Oracle Coherence clusters  3.3.1</markup>

</div>

<h3 id="_install_the_coherence_operator_package">Install the Coherence Operator Package</h3>
<div class="section">
<p>Once the Coherence Operator repository has been installed, the <code>coherence-operator.oracle.github.com</code> package can be installed, which will install the Coherence Operator itself.</p>

<markup
lang="bash"

>tanzu package install coherence \
    --package-name coherence-operator.oracle.github.com \
    --version 3.3.1 \
    --namespace coherence</markup>

<p>The Tanzu CLI will display the various steps it is going through to install the package and if all goes well, finally display <code>Added installed package 'coherence'</code>
The packages installed in the <code>coherence</code> namespace can be displayed using the CLI.</p>

<markup
lang="bash"

>tanzu package installed list --namespace coherence</markup>

<p>which should display the Coherence Operator package.</p>

<markup
lang="bash"

>NAME       PACKAGE-NAME                          PACKAGE-VERSION  STATUS
coherence  coherence-operator.oracle.github.com  3.3.1            Reconcile succeeded</markup>

<p>The Operator is now installed and ready to mage Coherence clusters.</p>

</div>
</div>
</doc-view>
