<doc-view>

<h2 id="_upgrading_from_operator_v2">Upgrading from Operator v2</h2>
<div class="section">
<p>Version 3 of the Coherence Operator is very different to version 2.
There is only a single CRD named <code>Coherence</code> instead of the three CRDs used by v2,
and the operator no longer uses Helm internally to install the Kubernetes resources.</p>

<p>In terms of usage and concepts, the biggest change is that there are no longer clusters and roles.
The <code>Coherence</code> CRD represents what would previously in v2 have been a role. A Coherence cluster that is made up
of multiple roles will just require multiple <code>Coherence</code> resources deploying to Kubernetes.
The simplification of the operator, and consequently the better reliability, far outweigh any advantage of being able
to put multiple roles in a single yaml file. If this is desire just put multiple <code>Coherence</code> resource definitions in
a single yaml file with the <code>---</code> separator.</p>

<p>For example:</p>

<p>In Operator v2 a cluster may have been defined with two roles, <code>storage</code> and <code>proxy</code> like this:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: CoherenceCluster
metadata:
  name: my-cluster
spec:
  roles:
    - role: storage
      replicas: 3
    - role: proxy
      replicas: 2</markup>

<p>In Operator v3 this needs to be two separate`Coherence` resources.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: my-cluster-storage
spec:
  - role: storage
    replicas: 3
    cluster: my-cluster
---
apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: my-cluster-proxy
spec:
  - role: proxy
    replicas: 2
    cluster: my-cluster</markup>

<div class="admonition note">
<p class="admonition-inline">To make both <code>Coherence</code> resources part of the same cluster the <code>cluster</code> field must now be set in both
resources to the same value, in this case <code>my-cluster</code>.</p>
</div>
</div>

<h2 id="_applications">Applications</h2>
<div class="section">
<p>Coherence applications in Operator v2 worked by application resources (jar files etc) being provided in an image
that was loaded as an init-container in the <code>Pod</code>, and the application artifacts copied to the classpath of the Coherence
container.
In version 3 of the Operator there is only one image required that should contain all of the resources required for the
application, including Coherence jar. This gives the application developer much more control over how the image is built
and what resources it contains, as well as making it more obvious what is going to be run when the container starts.</p>


<h3 id="_images">Images</h3>
<div class="section">
<p>In Operator v2 there were multiple images defined, one for Coherence and one used to provide application artifacts.
Because of the application changes described only a single image now needs to be specified in the <code>image</code> field
of the <code>CRD</code> spec.</p>

<p>See the <router-link to="#applications/010_overview.adoc" @click.native="this.scrollFix('#applications/010_overview.adoc')">Applications</router-link> section of the doecumentation for more details.</p>

</div>
</div>

<h2 id="_crd_differences">CRD Differences</h2>
<div class="section">
<p>A lot of the fields in the <code>Coherence</code> CRD are the same as when defining a role in version 2.
Whilst a number of new fields and features have been added in version 3, a handful of fields have moved,
and a small number, that no longer made sense, have been removed.
The <router-link to="#about/04_coherence_spec.adoc" @click.native="this.scrollFix('#about/04_coherence_spec.adoc')">Coherence Spec</router-link> page documents the full <code>Coherence</code> CRD, so it is
simple to locate where a field might have moved to.</p>

</div>

<h2 id="_logging_and_fluentd">Logging and Fluentd</h2>
<div class="section">
<p>Version 3 of the operator no longer has fields to configure a Fluentd side-car container.
There are a lot of different ways to configure Fluentd and making the Operator accomodate all of these was becoming
too much of a head-ache to do in a backwards compatible way.
If a Fluentd side-car is required it can just be added to the <code>Coherence</code> resource spec as an additional container,
so there is no limitation on the Fluentd configuration.
See the <router-link to="#logging/010_overview.adoc" @click.native="this.scrollFix('#logging/010_overview.adoc')">Logging documentation</router-link> for more examples.</p>

</div>

<h2 id="_prometheus_and_elasticsearch">Prometheus and Elasticsearch</h2>
<div class="section">
<p>Version 3 of the Operator no  longer comes with the option to install Prometheus and/or Elasticsearch.
This feature was only ever intended to make it easier to demo features that required Prometheus and Elasticsearch and
keeping this up to date was a headache nobody needed.
Both Prometheus and Elasticsearch have operators of their own which make installing them simple and importing the
dashboards provided by the Coherence Operator simple too.</p>

</div>
</doc-view>
