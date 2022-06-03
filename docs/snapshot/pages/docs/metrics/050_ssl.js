<doc-view>

<h2 id="_ssl_with_metrics">SSL with Metrics</h2>
<div class="section">
<p>It is possible to configure metrics endpoint to use SSL to secure the communication between server and
client. The SSL configuration is in the <code>coherence.metrics.ssl</code> section of the CRD spec.</p>

<p>For example:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test-cluster
spec:
  coherence:
    metrics:
      enabled: true
      ssl:
        enabled: true                            <span class="conum" data-value="1" />
        keyStore: metrics-keys.jks               <span class="conum" data-value="2" />
        keyStoreType: JKS                        <span class="conum" data-value="3" />
        keyStorePasswordFile: store-pass.txt     <span class="conum" data-value="4" />
        keyPasswordFile: key-pass.txt            <span class="conum" data-value="5" />
        keyStoreProvider:                        <span class="conum" data-value="6" />
        keyStoreAlgorithm: SunX509               <span class="conum" data-value="7" />
        trustStore: metrics-trust.jks            <span class="conum" data-value="8" />
        trustStoreType: JKS                      <span class="conum" data-value="9" />
        trustStorePasswordFile: trust-pass.txt   <span class="conum" data-value="10" />
        trustStoreProvider:                      <span class="conum" data-value="11" />
        trustStoreAlgorithm: SunX509             <span class="conum" data-value="12" />
        requireClientCert: true                  <span class="conum" data-value="13" />
        secrets: metrics-secret                  <span class="conum" data-value="14" /></markup>

<ul class="colist">
<li data-value="1">The <code>enabled</code> field when set to <code>true</code> enables SSL for metrics or when set to <code>false</code> disables SSL</li>
<li data-value="2">The <code>keyStore</code> field sets the name of the Java key store file that should be used to obtain the server&#8217;s key</li>
<li data-value="3">The optional <code>keyStoreType</code> field sets the type of the key store file, the default value is <code>JKS</code></li>
<li data-value="4">The optional <code>keyStorePasswordFile</code> sets the name of the text file containing the key store password</li>
<li data-value="5">The optional <code>keyPasswordFile</code> sets the name of the text file containing the password of the key in the key store</li>
<li data-value="6">The optional <code>keyStoreProvider</code> sets the provider name for the key store</li>
<li data-value="7">The optional <code>keyStoreAlgorithm</code> sets the algorithm name for the key store, the default value is <code>SunX509</code></li>
<li data-value="8">The <code>trustStore</code> field sets the name of the Java trust store file that should be used to obtain the server&#8217;s key</li>
<li data-value="9">The optional <code>trustStoreType</code> field sets the type of the trust store file, the default value is <code>JKS</code></li>
<li data-value="10">The optional <code>trustStorePasswordFile</code> sets the name of the text file containing the trust store password</li>
<li data-value="11">The optional <code>trustStoreProvider</code> sets the provider name for the trust store</li>
<li data-value="12">The optional <code>trustStoreAlgorithm</code> sets the algorithm name for the trust store, the default value is <code>SunX509</code></li>
<li data-value="13">The optional <code>requireClientCert</code> field if set to <code>true</code> enables two-way SSL where the client must also provide
a valid certificate</li>
<li data-value="14">The optional <code>secrets</code> field sets the name of the Kubernetes <code>Secret</code> to use to obtain the key store, truct store
and password files from.</li>
</ul>
<p>The various files and keystores referred to in the configuration above can be any location accessible in the image
used by the <code>coherence</code> container in the deployment&#8217;s <code>Pods</code>. Typically, for things such as SSL keys and certs,
these would be provided by obtained from <code>Secrets</code> loaded as additional <code>Pod</code> <code>Volumes</code>.
See <router-link to="#other/060_secret_volumes.adoc" @click.native="this.scrollFix('#other/060_secret_volumes.adoc')">Add Secrets Volumes</router-link> for the documentation on how to specify
secrets as additional volumes.</p>

</div>
</doc-view>
