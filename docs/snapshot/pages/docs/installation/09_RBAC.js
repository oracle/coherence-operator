<doc-view>

<h2 id="_rbac_roles">RBAC Roles</h2>
<div class="section">
<p>When installing the Coherence Operator into Kubernetes clusters with RBAC enabled, the Operator will require certain roles to work properly. Both the Operator Helm chart, and the Operator k8s manifest files will install all the required roles, role bindings and create a service account.</p>

</div>

<h2 id="_cluster_roles">Cluster Roles</h2>
<div class="section">
<p>By default, both install methods will create ClusterRole resources and ClusterRoleBinding resources to bind those roles to the Operator ServiceAccount. Some Kubernetes administrators are wary of letting arbitrary installations have ClusterRole permissions and try to discourage it. The Coherence Operator can run without ClusterRole permissions, but it is important to understand what this means from an operational point of view.</p>

<p>Cluster roles are used for a number of operator features:</p>

<ul class="ulist">
<li>
<p>Installing the CRDs - the Operator automatically ensures that the CRDs it requires are installed when it starts.</p>

</li>
<li>
<p>Installing the Web-Hook - the Operator automatically installs the defaulting and validating web-hooks for the <code>Coherence</code> resource when it starts. Without the validating web-hook a lot more care must be taken to ensure that only valid <code>Coherence</code> resource yaml is added to k8s. In the worst case, invalid yaml may ultimately cause the Operator to panic where invalid yaml would normally have been disallowed by the web-hook.</p>

</li>
<li>
<p>Coherence CLuster site and rack information - the Operator is used to supply site and rack values for the Coherence clusters that it manages. These values come from <code>Node</code> labels that the Operator must be able to look up. Without this information a Coherence cluster will have empty values for the <code>coherence.site</code> and <code>coherence.rack</code> system properties, meaning that Coherence will be unable to make data site-safe in k8s clusters that have multiple availability zones.</p>

</li>
<li>
<p>Monitoring multiple namespaces - if the Operator is to monitor multiple namespaces it must have cluster wide roles to do this</p>

</li>
</ul>
<p>Assuming that all the above reductions in features are acceptable then the Operator can be installed without creating cluster roles.</p>

</div>

<h2 id="_install_the_operator_without_clusterroles">Install the Operator Without ClusterRoles</h2>
<div class="section">
<p>The two methods of installing the Operator discussed in the <router-link to="/docs/installation/01_installation">Install Guide</router-link> can be used to install the Operator without ClusterRoles.</p>


<h3 id="_manually_install_crds">Manually Install CRDs</h3>
<div class="section">
<div class="admonition important">
<p class="admonition-textlabel">Important</p>
<p ><p>Before installing the Operator, with either method described below, the CRDs MUST be manually installed from the Operator manifest files.</p>

<p>The manifest files are published with the GitHub release at this link:
<a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/releases/download/v3.2.8/coherence-operator-manifests.tar.gz">3.2.8 Manifests</a></p>

<p>You MUST ensure that the CRD manifests match the version of the Operator being installed.</p>

<ul class="ulist">
<li>
<p>Download the manifests and unpack them.</p>

</li>
<li>
<p>In the directory that the .tar.gz file the was unpacked the <code>crd/</code> directory will the Coherence CRD.
The CRD can be installed with kubectl</p>

</li>
</ul>
<markup
lang="bash"

>kubectl create -f crd/coherence.oracle.com_coherence.yaml</markup>
</p>
</div>
</div>

<h3 id="_install_using_helm">Install Using Helm</h3>
<div class="section">
<p>The Operator can be installed from the Helm chart, as described in the <router-link to="/docs/installation/01_installation">Install Guide</router-link>.
The Helm chart contains values that control whether cluster roles are created when installing the chart. To install the chart without any cluster roles set the <code>clusterRoles</code> value to <code>false</code>.</p>

<markup
lang="bash"

>helm install  \
    --set clusterRoles=false       <span class="conum" data-value="1" />
    --namespace &lt;namespace&gt; \      <span class="conum" data-value="2" />
    coherence \
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">The <code>clusterRoles</code> value is set to false.</li>
<li data-value="2">The <code>&lt;namespace&gt;</code> value is the namespace that the Coherence Operator will be installed into
and without cluster roles will be the <em>only</em> namespace that the Operator monitors.</li>
</ul>

<h4 id="_allow_node_lookup">Allow Node Lookup</h4>
<div class="section">
<p>The Helm chart allows the Operator to be installed with a single <code>ClusterRole</code> allowing it to read k8s <code>Node</code> information. This is used to provide site, and rack labels, for Coherence cluster members. In environments where Kubernetes administrators are happy to allow the Operator read-only access to <code>Node</code> information the <code>nodeRoles</code> value can be set to <code>true</code>.</p>

<markup
lang="bash"

>helm install  \
    --set clusterRoles=false       <span class="conum" data-value="1" />
    --set nodeRoles=true           <span class="conum" data-value="2" />
    --namespace &lt;namespace&gt; \
    coherence \
    coherence/coherence-operator</markup>

<ul class="colist">
<li data-value="1">The <code>clusterRoles</code> value is set to <code>false</code>.</li>
<li data-value="2">The <code>nodeRoles</code> value is set to <code>true</code>, so a single ClusterRole will be applied to the Operator&#8217;s service account</li>
</ul>
</div>
</div>

<h3 id="_install_using_kustomize">Install Using Kustomize</h3>
<div class="section">
<p>The Operator can be installed using Kustomize with the manifest files, as described in the <router-link to="/docs/installation/01_installation">Install Guide</router-link>.</p>


<h4 id="_exclude_the_clusterrole_manifests">Exclude the ClusterRole Manifests</h4>
<div class="section">
<p>To install without cluster roles, after unpacking the manifests <code>.tar.gz</code> edit the <code>config/kustomization.yaml</code> file to comment out the inclusion of the cluster role bindings.</p>

<p>For example:</p>

<markup
lang="yaml"
title="kustomization.yaml"
>resources:
- service_account.yaml
- role.yaml
- role_binding.yaml
#- node_viewer_role.yaml
#- node_viewer_role_binding.yaml
#- cluster_role.yaml
#- cluster_role_binding.yaml</markup>

</div>

<h4 id="_disable_web_hooks_and_crd_installation">Disable Web-Hooks and CRD Installation</h4>
<div class="section">
<p>The Operator would normally install validating and defaulting web-hooks as well as ensuring that the Coherence CRDs are installed. Without cluster roles this must be disabled by editing the <code>manager/manager.yaml</code> file in the manifests.</p>

<p>Edit the Operator container <code>args</code> section of the deployment yaml to add command line arguments to <code>--enable-webhook=false</code> to disable web-hook creation and <code>--install-crd=false</code> to disable CRD installation.</p>

<p>For example, change the section of the <code>manager/manager.yaml</code> file that looks like this:</p>

<markup
lang="yaml"
title="manager/manager.yaml"
>        command:
          - /manager
        args:
          - --enable-leader-election
        envFrom:</markup>

<p>to be:</p>

<markup
lang="yaml"
title="manager/manager.yaml"
>        command:
          - /manager
        args:
          - --enable-leader-election
          - --enable-webhook=false
          - --install-crd=false
        envFrom:</markup>

</div>

<h4 id="_edit_the_operator_clusterrole_clusterrolebinding">Edit the Operator ClusterRole &amp; ClusterRoleBinding</h4>
<div class="section">
<p>The Operator will require a role and role binding to work in a single namespace.
Edit the <code>config/role.yaml</code> to change its type from <code>ClusterRole</code> to <code>Role</code>.</p>

<p>For example, change:</p>

<markup
lang="yaml"
title="role.yaml"
>apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: manager-role</markup>

<p>to be:</p>

<markup
lang="yaml"
title="role.yaml"
>apiVersion: rbac.authorization.k8s.io/v1
kind: Role  <span class="conum" data-value="1" />
metadata:
  creationTimestamp: null
  name: manager-role</markup>

<ul class="colist">
<li data-value="1"><code>ClusterRole</code> has been changed to <code>Role</code></li>
</ul>
<p>Edit the <code>config/role_binding.yaml</code> to change its type from <code>ClusterRoleBinding</code> to <code>RoleBinding</code>.</p>

<p>For example change:</p>

<markup
lang="yaml"
title="role_binding.yaml"
>apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: manager-rolebinding
  labels:
    control-plane: coherence
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: manager-role
subjects:
- kind: ServiceAccount
  name: coherence-operator
  namespace: default</markup>

<p>to be:</p>

<markup
lang="yaml"
title="role_binding.yaml"
>apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding  <span class="conum" data-value="1" />
metadata:
  name: manager-rolebinding
  labels:
    control-plane: coherence
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role <span class="conum" data-value="2" />
  name: manager-role
subjects:
- kind: ServiceAccount
  name: coherence-operator
  namespace: default</markup>

<ul class="colist">
<li data-value="1">The type has been changed from <code>ClusterRoleBinding</code> to <code>RoleBinding</code></li>
<li data-value="2">The role being bound has been changed from <code>ClusterRole</code> to <code>Role</code>.</li>
</ul>
</div>

<h4 id="_allow_node_lookup_2">Allow Node Lookup</h4>
<div class="section">
<p>In environments where Kubernetes administrators are happy to allow the Operator read-only access to <code>Node</code> information, the required <code>ClusterRole</code> can be created by leaving the relevant lines uncommented in the <code>config/kustomization.yaml</code> file.</p>

<p>For example:</p>

<markup
lang="yaml"
title="kustomization.yaml"
>resources:
- service_account.yaml
- role.yaml
- role_binding.yaml
- node_viewer_role.yaml         <span class="conum" data-value="1" />
- node_viewer_role_binding.yaml
#- cluster_role.yaml
#- cluster_role_binding.yaml</markup>

<ul class="colist">
<li data-value="1">The <code>node_viewer_role.yaml</code> and <code>node_viewer_role_binding.yaml</code> will now be left in the installation.</li>
</ul>
</div>
</div>
</div>
</doc-view>
