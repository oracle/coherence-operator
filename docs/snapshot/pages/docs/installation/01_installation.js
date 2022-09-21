<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Coherence Operator Installation</dt>
<dd slot="desc"><p>The Coherence Operator is available as a Docker image <code>oracle/coherence-operator:3.2.9</code> that can
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

<h3 id="_prerequisites">Prerequisites</h3>
<div class="section">
<p>The prerequisites apply to all installation methods.</p>

<ul class="ulist">
<li>
<p>Access to Oracle Coherence Operator images.</p>

</li>
<li>
<p>Access to a Kubernetes v1.18.0+ cluster. The Operator test pipeline is run using Kubernetes versions v1.18 upto v1.24</p>

</li>
<li>
<p>A Coherence application image using Coherence version 12.2.1.3 or later. Note that some functionality (e.g. metrics) is only
available in Coherence 12.2.1.4 and later.</p>

</li>
</ul>
<div class="admonition note">
<p class="admonition-inline">ARM Support: As of version 3.2.0, the Coherence Operator is build as a multi-architecture image that supports running in Kubernetes on both Linux/amd64 and Linux/arm64. The prerequisite is that the Coherence application image used has been built to support ARM.</p>
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

<h3 id="_high_availability">High Availability</h3>
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

<h2 id="_coherence_operator_images">Coherence Operator Images</h2>
<div class="section">
<p>The Coherence Operator uses a single images, the Operator also runs as an init-container in the Coherence cluster Pods.</p>

<ul class="ulist">
<li>
<p><code>ghcr.io/oracle/coherence-operator:3.2.9</code> - The Operator image.</p>

</li>
</ul>
<p>If no image is specified in the <code>Coherence</code> yaml, then the default Coherence image will also be used,</p>

<ul class="ulist">
<li>
<p><code>ghcr.io/oracle/coherence-ce:22.06.1</code> - The default Coherence image.</p>

</li>
</ul>
<p>If using a private image registry then these images will all need to be pushed to that registry for the Operator to work. The default Coherence image may be omitted if all Coherence applications will use custom Coherence images.</p>

</div>

<h2 id="manifest">Default Install with Kubectl</h2>
<div class="section">
<p>If you want the default Coherence Operator installation then the simplest solution is use <code>kubectl</code> to apply the manifests from the Operator release.</p>

<markup
lang="bash"

>kubectl apply -f https://github.com/oracle/coherence-operator/releases/download/v3.2.9/coherence-operator.yaml</markup>

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

<h2 id="_installing_with_helm">Installing With Helm</h2>
<div class="section">
<p>For more flexibility but the simplest way to install the Coherence Operator is to use the Helm chart.
This ensures that all the correct resources will be created in Kubernetes.</p>


<h3 id="helm">Add the Coherence Helm Repository</h3>
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

<h3 id="_run_the_operator_as_a_non_root_user">Run the Operator as a Non-Root User</h3>
<div class="section">
<p>The Operator container can be configured with a <code>securityContext</code> so that it runs as a non-root user.</p>

<p>This can be done using a values file:</p>

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

<p>Alternatively, the <code>securityContext</code> values can be set on the command line as <code>--set</code> parameters:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set securityContext.runAsNonRoot=true \
    --set securityContext.runAsUser=1000 \
    coherence \
    coherence/coherence-operator</markup>

</div>

<h3 id="_uninstall_the_coherence_operator_helm_chart">Uninstall the Coherence Operator Helm chart</h3>
<div class="section">
<p>To uninstall the operator:</p>

<markup
lang="bash"

>helm delete coherence-operator --namespace &lt;namespace&gt;</markup>

</div>
</div>

<h2 id="_operator_scope">Operator Scope</h2>
<div class="section">
<p>The recommended way to install the Coherence Operator is to install a single instance of the operator into a namespace
and where it will then control <code>Coherence</code> resources in all namespaces across the Kubernetes cluster.
Alternatively it may be configured to watch a sub-set of namespaces by setting the <code>WATCH_NAMESPACE</code> environment variable.
The watch namespace(s) does not have to include the installation namespace.</p>

<p>In theory, it is possible to install multiple instances of the Coherence Operator into different namespaces, where
each instances monitors a different set of namespaces. There are a number of potential issues with this approach, so
it is not recommended.</p>

<ul class="ulist">
<li>
<p>Only one CRD can be installed - Different releases of the Operator may use slightly different CRD versions, for example
a new version may introduce extra fields not in the previous version. As the CRD version is <code>v1</code> there is no guarantee
which CRD version has actually installed, which could lead to subtle issues.</p>

</li>
<li>
<p>The operator creates and installs defaulting and validating web-hooks. A web-hook is associated to a CRD resource so
installing multiple web-hooks for the same resource may lead to issues. If an operator is uninstalled, but the web-hook
configuration remains, then Kubernetes will not accept modifications to resources of that type as it will be
unable to contact the web-hook.</p>

</li>
</ul>
<p>To set the watch namespaces when installing with helm set the <code>watchNamespaces</code> value, for example:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set watchNamespaces=payments,catalog,customers <span class="conum" data-value="1" />
    coherence-operator \
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">The <code>payments</code>, <code>catalog</code> and <code>customers</code> namespaces will be watched by the Operator.</li>
</ul>
</div>

<h2 id="_operator_image">Operator Image</h2>
<div class="section">
<p>The Helm chart uses a default registry to pull the Operator image from.
If the image needs to be pulled from a different location (for example an internal registry) then the <code>image</code> field
in the values file can be set, for example:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set image=images.com/coherence-operator:0.1.2 <span class="conum" data-value="1" />
    coherence-operator \
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">The image used to run the Operator will be <code>images.com/coherence-operator:0.1.2</code>.</li>
</ul>

<h3 id="_image_pull_secrets">Image Pull Secrets</h3>
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
</div>

<h2 id="kubectl">Install with Kubectl and Kustomize</h2>
<div class="section">
<p>If you want to use yaml directly to install the operator, with something like <code>kubectl</code>, you can use the manifest files
published with the GitHub release at this link:
<a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/releases/download/v3.2.9/coherence-operator-manifests.tar.gz">3.2.9 Manifests</a></p>

<p>These manifest files are for use with a tool called Kustomize, which is built into <code>kubectl</code>
see the documentation here: <a id="" title="" target="_blank" href="https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/">https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/</a></p>

<p>Download the
<a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/releases/download/v3.2.9/coherence-operator-manifests.tar.gz">3.2.9 Manifests</a>
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

>cd ./manager &amp;&amp; kustomize edit set image controller=myregistry/coherence-operator:3.2.9</markup>

<p>Change the name of the Operator image by running the command below, changing the image name to the registry and image name
that you are using for the Operator utilities image</p>

<markup
lang="bash"

>cd ./manager &amp;&amp; kustomize edit add configmap env-vars --from-literal OPERATOR_IMAGE=myregistry/coherence-operator:3.2.9</markup>

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
<p><code>ghcr.io/oracle/coherence-operator-package:3.2.9</code> - the Coherence Operator package</p>

</li>
<li>
<p><code>ghcr.io/oracle/coherence-operator-repo:3.2.9</code> - the Coherence Operator repository</p>

</li>
</ul>

<h3 id="_install_the_coherence_repository">Install the Coherence Repository</h3>
<div class="section">
<p>The first step to deploy the Coherence Operator package in Tanzu is to add the repository.
This can be done using the Tanzu CLI.</p>

<markup
lang="bash"

>tanzu package repository add coherence-repo \
    --url ghcr.io/oracle/coherence-operator-repo:3.2.9 \
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
coherence-operator.oracle.github.com  Oracle Coherence Operator  A Kubernetes operator for managing Oracle Coherence clusters  3.2.9</markup>

</div>

<h3 id="_install_the_coherence_operator_package">Install the Coherence Operator Package</h3>
<div class="section">
<p>Once the Coherence Operator repository has been installed, the <code>coherence-operator.oracle.github.com</code> package can be installed, which will install the Coherence Operator itself.</p>

<markup
lang="bash"

>tanzu package install coherence \
    --package-name coherence-operator.oracle.github.com \
    --version 3.2.9 \
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
coherence  coherence-operator.oracle.github.com  3.2.9            Reconcile succeeded</markup>

<p>The Operator is now installed and ready to mage Coherence clusters.</p>

</div>
</div>
</doc-view>
