<doc-view>

<h2 id="_fips_compatibility">FIPS Compatibility</h2>
<div class="section">
<p>The Coherence Operator image uses an empty scratch image for its base image.
This means that the Coherence Operator image is FIPS compatible and can be run in a FIPS compliant Kubernetes cluster.</p>

<p>As the Coherence Operator is written in Go, it can use Go&#8217;s built in FIPS support.
To run the Coherence Operator in a FIPS compliant mode, it needs to be installed with the environment variable <code>GODEBUG</code>
set to either <code>fips140=on</code> or <code>fips140=only</code>. This is explained in the Golang <a id="" title="" target="_blank" href="https://go.dev/doc/security/fips140">FIPS-140 documentation</a>.</p>

<p>How the <code>GODEBUG</code> environment variable is set depends on how the operator is installed.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>Although the Coherence Operator image can easily be installed in a FIPS compliant mode, none of the default
Oracle Coherence images used by the operator are FIPS complaint.
The Oracle Coherence team does not currently publish FIPS compliant Coherence images.
Coherence is FIPS compatible and correctly configured applications running in an image that has a FIPS
compliant JDK and FIPS compliant base O/S will be FIPS complaint.
Customers must build their own FIPS complaint Java and Coherence images, which the operator will then manage.</p>
</p>
</div>

<h3 id="_install_using_yaml_manifests">Install Using Yaml Manifests</h3>
<div class="section">
<p>If <router-link to="/docs/installation/011_install_manifests">installing using the yaml manifests</router-link>,
the yaml must be edited to add the <code>GODEBUG</code> environment variable to
the operator deployments environment variables:</p>

<p>Find the <code>env:</code> section of the operator <code>Deployment</code> in the yaml file, it looks like this:</p>

<markup
lang="yaml"

>        env:
        - name: OPERATOR_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace</markup>

<p>then add the required <code>GODEBUG</code> value, for example</p>

<markup
lang="yaml"

>        env:
        - name: OPERATOR_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: GODEBUG
          value: fips140=on</markup>

</div>

<h3 id="_install_using_helm">Install Using Helm</h3>
<div class="section">
<p>If <router-link to="/docs/installation/012_install_helm">installing the operator using Helm</router-link>
The Coherence Operator Helm chart has a <code>fips</code> field in its values file.
This value is used to set the <code>GODEBUG</code> environment variables.
The <code>fips</code> value is unset by default, if set it must be one of the values, "off", "on" or "only".
If <code>fips</code> is set to any other value the chart will fail to install.</p>

<p>For example, to install the operator with fips set to "on"</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set fips=on
    coherence-operator \
    coherence/coherence-operator</markup>

</div>
</div>
</doc-view>
