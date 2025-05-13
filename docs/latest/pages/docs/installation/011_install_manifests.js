<doc-view>

<h2 id="_install_using_manifests">Install Using Manifests</h2>
<div class="section">
<p>If you want the default Coherence Operator installation then the simplest solution is use <code>kubectl</code> to
apply the manifests from the Operator release.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>As of v3.5.0 of the Operator the manifest yaml also installs the two CRDs that the Operator uses.
In previous releases the Operator would install the CRDs when it started but this behaviour is disabled by default
when installing with the manifest yaml.</p>
</p>
</div>
<p>The following command will install the Operator. This assumes that the Kubernetes account being used to perform
the installation has all the RBAC permissions required to install all the resource types in the yaml file.</p>

<markup
lang="bash"

>kubectl apply -f https://github.com/oracle/coherence-operator/releases/download/v3.5.0/coherence-operator.yaml</markup>

<p>This will create a namespace called <code>coherence</code> and install the CRDs and the Operator into the namespace,
along with all the required <code>ClusterRole</code> and <code>RoleBinding</code> resources. The <code>coherence</code> namespace can be changed by
downloading and editing the yaml file.</p>

<p>In some restricted environments, a Kubernetes user might not have RBAC permissions to install CRDs.
In this case the <code>coherence-operator.yaml</code> file will need to be edited to remove the two CRDs from the
beginning of the file. The CRDs <strong><em>must be manually installed before the Operator is installed</em></strong>, as described
below in <router-link to="#manual-crd" @click.native="this.scrollFix('#manual-crd')">Manually Install the CRDs</router-link>.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Because the <code>coherence-operator.yaml</code> manifest also creates the namespace, the corresponding <code>kubectl delete</code>
command will <em>remove the namespace and everything deployed to it</em>! If you do not want this behaviour you should edit
the <code>coherence-operator.yaml</code> to remove the namespace section from the start of the file.</p>
</p>
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

</div>

<h2 id="manifest-restrict">Installing Without Cluster Roles</h2>
<div class="section">
<p>The default install for the Operator is to have one Operator deployment that manages all Coherence resources across
all the namespaces in a Kubernetes cluster. This requires the Operator to have cluster role RBAC permissions
to manage and monitor all the resources.</p>

<p>Sometimes, for security reasons or for example in a shared Kubernetes cluster this is not desirable.
The Operator can therefore be installed with plain namespaced scoped roles and role bindings.
The Operator release includes a single yaml file named <code>coherence-operator-restricted.yaml</code> that may be used to install
the Operator into a single namespace without any cluster roles.</p>

<p>The Operator installed with this yaml</p>

<ul class="ulist">
<li>
<p>will not use WebHooks</p>

</li>
<li>
<p>will not look-up Node labels for Coherence site and rack configurations</p>

</li>
</ul>
<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>As of v3.5.0 of the Operator the <code>coherence-operator-restricted.yaml</code> also installs the two CRDs that the Operator uses.
In previous releases the Operator would install the CRDs when it started but this behaviour is disabled by default
when installing with the manifest yaml.</p>
</p>
</div>
<p>The following command will install the Operator. This assumes that the Kubernetes account being used to perform
the installation has all the RBAC permissions required to install all the resource types in the yaml file.</p>

<markup
lang="bash"

>kubectl apply -f https://github.com/oracle/coherence-operator/releases/download/v3.5.0/coherence-operator-restricted.yaml</markup>

<div class="admonition important">
<p class="admonition-textlabel">Important</p>
<p ><p>In some restricted environments, a Kubernetes user might not have RBAC permissions to install CRDs.
In this case the <code>coherence-operator.yaml</code> file will need to be edited to remove the two CRDs from the
beginning of the file. The CRDs <strong><em>must be manually installed before the Operator is installed</em></strong>, as described
below in <router-link to="#manual-crd" @click.native="this.scrollFix('#manual-crd')">Manually Install the CRDs</router-link>.</p>
</p>
</div>
</div>

<h2 id="manual-crd">Manually Install the CRDs</h2>
<div class="section">
<p>Although by default the Operator will install its CRDs, they can be manually installed into Kubernetes.
This may be required where the Operator is running with restricted permissions as described above.</p>

<p>The Operator release artifacts include small versions of the two CRDs which can be installed with the following commands:</p>

<markup
lang="bash"

>kubectl apply -f https://github.com/oracle/coherence-operator/releases/download/v3.5.0/coherence.oracle.com_coherence_small.yaml
kubectl apply -f https://github.com/oracle/coherence-operator/releases/download/v3.5.0/coherencejob.oracle.com_coherence_small.yaml</markup>

<p>The small versions of the CRDs are identical to the full versions but hav a cut down OpenAPI spec with a lot of comments
removed so that the CRDs are small enough to be installed with <code>kubectl apply</code></p>

</div>

<h2 id="_change_the_operator_replica_count">Change the Operator Replica Count</h2>
<div class="section">
<p>When installing with single manifest yaml file, the replica count can be changed by editing the yaml file itself
to change the occurrence of <code>replicas: 3</code> in the manifest yaml to <code>replicas: 1</code></p>

<p>For example, this could be done using <code>sed</code></p>

<markup
lang="bash"

>sed -i -e 's/replicas: 3/replicas: 1/g' coherence-operator.yaml</markup>

<p>Or on MacOS, where <code>sed</code> is slightly different:</p>

<markup
lang="bash"

>sed -i '' -e 's/replicas: 3/replicas: 1/g' coherence-operator.yaml</markup>

</div>
</doc-view>
