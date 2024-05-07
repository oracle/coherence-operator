<doc-view>

<h2 id="_global_labels_and_annotations">Global Labels and Annotations</h2>
<div class="section">
<p>It is possible to specify a global set of labels and annotations that will be applied to all resources.
Global labels and annotations can be specified in two ways:</p>

<ul class="ulist">
<li>
<p>For an individual <code>Coherence</code> deployment, in which case they will be applied to all the Kubernetes resources
created for that deployment</p>

</li>
<li>
<p>As part of the Operator install, in which case they will be applied to all Kubernetes resources managed by the
Operator, including all Coherence clusters and related resources</p>

</li>
</ul>
</div>

<h2 id="_specify_global_labels_for_a_coherence_resource">Specify Global Labels for a Coherence Resource</h2>
<div class="section">
<p>The <code>Coherence</code> CRD contains a <code>global</code> field that allows global labels and annotations to be specified.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  replicas: 3
  global:
    labels:
      one: "label-one"
      two: "label-two"</markup>

<p>If the yaml above is applied to Kubernetes, then every resource the Operator creates for the <code>storage</code> Coherence
deployment, it will add the two labels, <code>one=label-one</code> and <code>two=label-two</code>. This includes the <code>StatefulSet</code>,
the <code>Pods</code>, any <code>Service</code> such as the stateful set service, the WKA service, etc.</p>

<p>If any of the labels in the <code>global</code> section are also in the Pod labels section or for the Services for exposed ports,
those labels will take precedence.</p>

<p>For example</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: storage
spec:
  replicas: 3
  labels:
    one: "pod-label-one"
  global:
    labels:
      one: "label-one"
      two: "label-one"</markup>

<p>In the yaml above, the global label <code>one=label-one</code> and <code>two=labl-two</code> will be applied to every resource created for
the <code>Coherence</code> deployment except for the Pods. The Operator uses the <code>spec.labels</code> field to define Pods specific labels,
so in this case the Pod labels will be <code>one=pod-label-one</code> from the <code>spec.labels</code> field and <code>two=labl-two</code> from the global
labels.</p>

</div>

<h2 id="_specify_global_labels_when_installing_the_operator">Specify Global Labels when Installing the Operator</h2>
<div class="section">
<p>The Operator <code>runner</code> binary has various command line flags that can be specified on its command line.
Two of these flags when starting the Operator are:</p>

<ul class="ulist">
<li>
<p><code>--global-label</code> to specify a global label key and value</p>

</li>
<li>
<p><code>--global-annotation</code> to specify a global annotation key and value</p>

</li>
</ul>
<p>Both of these command line flags can be specified multiple times if required.</p>

<p>For example:</p>

<markup
lang="bash"

>runner operator --global-label one=label-one --global-annoataion foo=bar --global-label two=label-two</markup>

<p>The command above will start the Operator with two global labels,<code>one=label-one</code> and <code>two=labl-two</code> and with
one global annotation <code>foo=bar</code>.</p>

<p>The Operator will then apply these labels and annotations to every Kubernetes resource that it creates.</p>


<h3 id="_installing_using_the_manifest_files">Installing Using the Manifest Files</h3>
<div class="section">
<p>When installing the Operator using the manifest yaml files, additional command line flags can be configured
by manually editing the yaml file before installing.</p>

<p>Download the yaml manifest file from the GitHub repo
<a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/releases/download/v3.3.5/coherence-operator.yaml">https://github.com/oracle/coherence-operator/releases/download/v3.3.5/coherence-operator.yaml</a></p>

<p>Find the section of the yaml file the defines the Operator container args, the default looks like this</p>

<markup
lang="yaml"
title="coherence-operator.yaml"
>      - args:
        - operator
        - --enable-leader-election</markup>

<p>Then edit the argument list to add the required <code>--global-label</code> and <code>--global-annotation</code> flags.</p>

<p>For example, to add the same <code>--global-label one=label-one --global-annotation foo=bar --global-label two=label-two</code>
flags, the file would look like this:</p>

<markup
lang="yaml"
title="coherence-operator.yaml"
>      - args:
        - operator
        - --enable-leader-election
        - --global-label
        - one=label-one
        - --global-annotation
        - foo=bar
        - --global-label
        - two=label-two`</markup>

<div class="admonition important">
<p class="admonition-textlabel">Important</p>
<p ><p>Container arguments must each be a separate entry in the arg list.
This is valid</p>

<markup
lang="yaml"
title="coherence-operator.yaml"
>      - args:
        - operator
        - --enable-leader-election
        - --global-label
        - one=label-one</markup>

<p>This is not valid</p>

<markup
lang="yaml"
title="coherence-operator.yaml"
>      - args:
        - operator
        - --enable-leader-election
        - --global-label  one=label-one</markup>
</p>
</div>
</div>

<h3 id="_installing_using_the_helm_chart">Installing Using the Helm Chart</h3>
<div class="section">
<p>If installing the Operator using the Helm chart, the global labels and annotations can be specified as values
as part of the Helm command or in a values file.</p>

<p>For example, to add the same <code>--global-label one=label-one --global-annotation foo=bar --global-label two=label-two</code>
flags, create a simple values file:</p>

<markup

title="global-values.yaml"
>globalLabels:
  one: "label-one"
  two: "label-two"

globalAnnotations:
  foo: "bar"</markup>

<p>Use the values file when installing the Helm chart</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --values global-values.yaml
    coherence \
    coherence/coherence-operator</markup>

</div>
</div>
</doc-view>
