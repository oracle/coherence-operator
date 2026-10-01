<doc-view>

<h2 id="_install_using_olm">Install Using OLM</h2>
<div class="section">
<p>The <a id="" title="" target="_blank" href="https://olm.operatorframework.io">Operator Lifecycle Manager</a> (OLM) can be used to install the Coherence Operator.</p>

<p>As part of the Coherence Operator release bundle and catalog images are pushed to the release image registry.
These images can be used to deploy the operator on Kubernetes clusters that are running OLM.</p>

<p>The Coherence Operator is not currently available on Operator Hub, but the required resource files can be created
manually to install the operator into Kubernetes.</p>


<h3 id="_install_the_coherence_operator_catalogsource">Install The Coherence Operator CatalogSource</h3>
<div class="section">
<p>Create a yaml manifest that will install the Coherence Operator CatalogSource as shown below.</p>

<markup
lang="yaml"
title="operator-catalog-source.yaml"
>apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: coherence-operator-catalog
  namespace: olm
spec:
  displayName: Oracle Coherence Operators
  image: ghcr.io/oracle/coherence-operator-catalog:latest
  publisher: Oracle Corporation
  sourceType: grpc
  updateStrategy:
    registryPoll:
      interval: 60m</markup>

<p>Install the CatalogSource into the <code>olm</code> namespace using the following command:</p>

<markup
lang="bash"

>kubectl apply -f operator-catalog-source.yaml</markup>

<p>Running the following command should list the catalog sources installed in the <code>olm</code> namespace, including the Coherence
catalog source.</p>

<markup
lang="bash"

>kubectl -n olm get catalogsource</markup>

<p>The Coherence catalog source Pod should eventually be ready, which can be verified with the following command:</p>

<markup
lang="bash"

>POD=$(kubectl -n olm get pod -l olm.catalogSource=coherence-operator-catalog)
kubectl -n olm wait --for condition=ready --timeout 480s $(POD)</markup>

</div>
</div>
</doc-view>
