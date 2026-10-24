<doc-view>

<h2 id="_tls_cipher_suites">TLS Cipher Suites</h2>
<div class="section">
<p>The Coherence Operator uses TLS for various client connections and server sockets.
TLS can support a number of cipher suites, some of which are deemed legacy and insecure.
These insecure ciphers are usually only present for backwards compatability.</p>

<p>The Coherence Operator is written in Go, and the ciphers supported are determined by the version og Go
used to build the operator.
Go splits ciphers into two lists a secure list and an insecure list, the insecure ciphers are disabled by default.</p>

<p>Oracle Global Security has stricter requirements than the default Go cipher list.
By default, the Coherence Operator enables only ciphers in Go&#8217;s secure list, except for
<code>TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA</code> and <code>TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA</code>, which are disabled.</p>

<p>It is possible to enable or disable cipher suites when installing the Coherence Operator.
The Coherence Operator has two command line flags which can be used to specify ciphers to be allowed or denied.</p>

<ul class="ulist">
<li>
<p>The <code>--cipher-allow-list</code> command line flag is used to specify cipher names to add to the allowed list.</p>

</li>
<li>
<p>The <code>--cipher-deny-list</code> command line flag is used to specify cipher names to add to the disabled list.</p>

</li>
</ul>
<p>Multiple ciphers can be enabled and disabled by specifying the relevant command line flag multiple times.</p>

<p>If a cipher name is added to both the allow list and to the deny list, it will be disabled.</p>

<div class="admonition note">
<p class="admonition-textlabel">Note</p>
<p ><p>If either the <code>--cipher-allow-list</code> or <code>--cipher-deny-list</code> is set to a name that does not match any of the
supported Go cipher names, the Operator will display an error in its log and will not start.
See the <a id="" title="" target="_blank" href="https://pkg.go.dev/crypto/tls#pkg-constants">Go TLS package documentation</a> for a lost of valid names.</p>
</p>
</div>
<p><strong>Only Allow FIPS Ciphers</strong></p>

<p>The Coherence Operator can be installed in FIPS mode to only support FIPS compliant ciphers,
see the <router-link to="/docs/installation/100_fips">FIPS modes</router-link> documentation for details.</p>

<p>How the command line flags are set depends on how the Coherence Operator is installed.</p>


<h3 id="_install_using_yaml_manifests">Install Using Yaml Manifests</h3>
<div class="section">
<p>If <router-link to="/docs/installation/011_install_manifests">installing using the yaml manifests</router-link>,
the yaml must be edited to add the required flags:</p>

<p>Find the <code>args:</code> section of the operator <code>Deployment</code> in the yaml file, it looks like this:</p>

<markup
lang="yaml"

>        args:
          - operator
          - --enable-leader-election</markup>

<p>then add the required allow or disallow flags. For example to allow <code>TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA</code>
the args can be edited as shown below:</p>

<markup
lang="yaml"

>        args:
          - operator
          - --enable-leader-election
          - --cipher-allow-list=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA</markup>

<p>To enable both <code>TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA</code> and <code>TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA</code> ciphers:</p>

<markup
lang="yaml"

>        args:
          - operator
          - --enable-leader-election
          - --cipher-allow-list=TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA
          - --cipher-allow-list=TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA</markup>

</div>

<h3 id="_install_using_helm">Install Using Helm</h3>
<div class="section">
<p>If <router-link to="/docs/installation/012_install_helm">installing the operator using Helm</router-link>
The Coherence Operator Helm chart has a <code>cipherAllowList</code> field and <code>cipherDenyList</code> field in its values file.
These values are Helm arrays and can be set to a list of ciphers to be enabled or disabled.</p>

<p>The simplest way to set lists on the Helm command line is using the <code>--set-json</code> command line flag.
For example to allow <code>TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA</code></p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set-json='cipherAllowList=["TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA"]'
    coherence-operator \
    coherence/coherence-operator</markup>

<p>To enable both <code>TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA</code> and <code>TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA</code> ciphers:</p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set-json='cipherAllowList=["TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA", "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA"]'
    coherence-operator \
    coherence/coherence-operator</markup>

<p>To disable <code>TLS_CHACHA20_POLY1305_SHA256</code></p>

<markup
lang="bash"

>helm install  \
    --namespace &lt;namespace&gt; \
    --set-json='cipherDenyList=["TLS_CHACHA20_POLY1305_SHA256"]'
    coherence-operator \
    coherence/coherence-operator</markup>

</div>
</div>
</doc-view>
