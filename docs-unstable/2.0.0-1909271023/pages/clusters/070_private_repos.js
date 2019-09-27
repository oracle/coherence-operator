<doc-view>

<h2 id="_using_private_image_registries">Using Private Image Registries</h2>
<div class="section">
<p>Sometimes the images used by a Coherence cluster need to be pulled from a private image registry that requires credentials.
The Coherence Operator supports supplying credentials in the <code>CoherenceCluster</code> configuration.
The Kubernetes documentation on <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/containers/images/#using-a-private-registry"> using a private registries</a>
gives a number of options for supplying credentials.</p>

</div>

<h2 id="_use_imagepullsecrets">Use ImagePullSecrets</h2>
<div class="section">
<p>Kubernetes supports configuring pods to use <code>imagePullSecrets</code> for pulling images. If possible, this is the preferable
and most portable route.
See the <a id="" title="" target="_blank" href="https://kubernetes.io/docs/concepts/containers/images/#specifying-imagepullsecrets-on-a-pod">kubernetes docs</a>
for this.
Once secrets have been created in the namespace that the <code>CoherenceCluster</code> is to be installed in then the secret name
can be specified in the <code>CoherenceCluster</code> <code>spec</code>. It is possible to specify multiple secrets in the case where the different
images being used are pulled from different registries.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  imagePullSecrets:
    - coherence-secret  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>coherence-secret</code> will be used for pulling images from the registry associated to the secret</li>
</ul>
<p>The <code>imagePullSecrets</code> field is a list of string values so multiple secrets can be specified for different authenticated
registries in the case where the Coherence cluster will use images from different authenticated registries..</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: test-cluster
spec:
  imagePullSecrets:     <span class="conum" data-value="1" />
    - coherence-secret
    - ocr-secret</markup>

<ul class="colist">
<li data-value="1">The <code>imagePullSecrets</code> list specifies two secrets to use <code>coherence-secret</code> and <code>ocr-secret</code></li>
</ul>
<p>Image pull secrets are only specified for the <code>CoherenceCluster</code> as a whole as there is no benefit to being able to
specify different secrets for different roles within a cluster.</p>

</div>
</doc-view>
