<doc-view>

<h2 id="_secure_coherence_with_tls">Secure Coherence with TLS</h2>
<div class="section">
<p>This example is going to show how to use TLS (or SSL) to secure communication between different parts of a Coherence cluster and applications. This is quite a long guide as there are a number of things that can be secured wth TLS.</p>

<p>This example shows how to secure various parts of Coherence clusters using TLS.</p>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p><img src="./images/GitHub-Mark-32px.png" alt="GitHub Mark 32px" />
 The complete source code for this example is in the <a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/tree/master/examples/090_tls">Coherence Operator GitHub</a> repository.</p>
</p>
</div>
<p>In this example we are going to use <a id="" title="" target="_blank" href="https://cert-manager.io">Cert Manager</a> to manage the keys and certs for our Coherence server and clients. Cert Manage makes managing certificates in Kubernetes very simple, but it isn&#8217;t the only solution.</p>

<p>Although securing clusters with TLS is a common request, if running in a secure isolated Kubernetes cluster, you need to weigh up the pros and cons regarding the performance impact TLS will give over the additional security.</p>

<p>Using Cert Manager we will ultimately end up with four k8s <code>Secrets</code>:</p>

<ul class="ulist">
<li>
<p>A <code>Secret</code> containing the server keys, certs, keystore and truststore</p>

</li>
<li>
<p>A <code>Secret</code> containing a single file containing the server keystore, truststore and key password</p>

</li>
<li>
<p>A <code>Secret</code> containing the client keys, certs, keystore and truststore</p>

</li>
<li>
<p>A <code>Secret</code> containing a single file containing the client keystore, truststore and key password</p>

</li>
</ul>
<p>If you do not want to use Cert Manager to try this example then a long as you have a way to create the required <code>Secrets</code> containing the keys and passwords above then you can skip to the section on <router-link to="#coherence" @click.native="this.scrollFix('#coherence')">Securing Coherence</router-link>.</p>


<h3 id="_what_the_example_will_cover">What the Example will Cover</h3>
<div class="section">
<ul class="ulist">
<li>
<p><router-link to="#install_operator" @click.native="this.scrollFix('#install_operator')">Install the Operator</router-link></p>

</li>
<li>
<p><router-link to="#setup_cert_manager" @click.native="this.scrollFix('#setup_cert_manager')">Setting Up Cert-Manager</router-link></p>
<ul class="ulist">
<li>
<p><router-link to="#create_self_signed_issuer" @click.native="this.scrollFix('#create_self_signed_issuer')">Create the SelfSigned Issuer</router-link></p>

</li>
<li>
<p><router-link to="#create_ce_cert" @click.native="this.scrollFix('#create_ce_cert')">Create the CA Certificate</router-link></p>

</li>
<li>
<p><router-link to="#create_ca_issuer" @click.native="this.scrollFix('#create_ca_issuer')">Create the CA issuer</router-link></p>

</li>
<li>
<p><router-link to="#create_coherence_keystores" @click.native="this.scrollFix('#create_coherence_keystores')">Create the Coherence Keys, Certs and KeyStores</router-link></p>
<ul class="ulist">
<li>
<p><router-link to="#server_password_secret" @click.native="this.scrollFix('#server_password_secret')">Create the Server Keystore Password Secret</router-link></p>

</li>
<li>
<p><router-link to="#server_cert" @click.native="this.scrollFix('#server_cert')">Create the Server Certificate</router-link></p>

</li>
<li>
<p><router-link to="#client_certs" @click.native="this.scrollFix('#client_certs')">Create the Client Certificate</router-link></p>

</li>
</ul>
</li>
</ul>
</li>
<li>
<p><router-link to="#coherence" @click.native="this.scrollFix('#coherence')">Securing Coherence Clusters</router-link></p>
<ul class="ulist">
<li>
<p><router-link to="#images" @click.native="this.scrollFix('#images')">Build the Example Images</router-link></p>

</li>
<li>
<p><router-link to="#socket_provider" @click.native="this.scrollFix('#socket_provider')">Configure a Socket Provider</router-link></p>

</li>
</ul>
</li>
<li>
<p><router-link to="#tcmp" @click.native="this.scrollFix('#tcmp')">Secure Cluster Membership</router-link></p>

</li>
<li>
<p><router-link to="#extend" @click.native="this.scrollFix('#extend')">Secure Extend</router-link></p>

</li>
<li>
<p><router-link to="#grpc" @click.native="this.scrollFix('#grpc')">Secure gRPC</router-link></p>

</li>
</ul>
</div>

<h3 id="install_operator">Install the Operator</h3>
<div class="section">
<p>To run the examples below, you will need to have installed the Coherence Operator, do this using whatever method you prefer from the <a id="" title="" target="_blank" href="https://oracle.github.io/coherence-operator/docs/latest/#/docs/installation/01_installation">Installation Guide</a></p>

</div>

<h3 id="setup_cert_manager">Setting Up Cert-Manager</h3>
<div class="section">
<p>In this example we will use self-signed certs as this makes everything easy to get going.
Cert Manager has a number of ways to configure real certificates for production use.
Assuming that you&#8217;ve installed Cert Manager using one of the methods in its <a id="" title="" target="_blank" href="https://cert-manager.io/docs/installation/">Install Guide</a> we can proceed to created all of the required resources.</p>


<h4 id="create_self_signed_issuer">Create the SelfSigned Issuer</h4>
<div class="section">
<p>This is used to generate a root CA for use with the CA Issuer.
Here we are using a <code>ClusterIssuer</code> so that we can use a single self-signed issuer across all namespaces.
We could have instead created an <code>Issuer</code> in a single namespace.</p>

<markup
lang="yaml"
title="manifests/selfsigned-issuer.yaml"
>apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned-issuer
spec:
  selfSigned: {}</markup>

<p>Create the <code>ClusterIssuer</code> with the following command. As this is a <code>ClusterIssuer</code>, is does not require a namespace.</p>

<markup
lang="bash"

>kubectl apply -f manifests/selfsigned-issuer.yaml</markup>

<p>We can list the <code>ClusterIssuers</code> in the cluster:</p>

<markup
lang="bash"

>kubectl get clusterissuer</markup>

<p>We should see that the <code>selfsigned-issuer</code> is present and is ready.</p>

<markup
lang="bash"

>NAME                READY   AGE
selfsigned-issuer   True    14m</markup>

</div>
</div>

<h3 id="create_ce_cert">Create the CA Certificate</h3>
<div class="section">
<p>We’re going to create an internal CA that will be used to sign our certificate requests for the Coherence server and clients that we will run later. Both the server and client will use the CA to validate a connection.</p>

<p>To create the CA issuer, first create a self-signed CA certificate.</p>

<markup
lang="yaml"
title="manifests/ca-cert.yaml"
>apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: ca-certificate
spec:
  issuerRef:
    name: selfsigned-issuer   <span class="conum" data-value="1" />
    kind: ClusterIssuer
    group: cert-manager.io
  secretName: ca-cert        <span class="conum" data-value="2" />
  duration: 2880h # 120d
  renewBefore: 360h # 15d
  commonName: Cert Admin
  isCA: true                 <span class="conum" data-value="3" />
  privateKey:
    size: 2048
  usages:
    - digital signature
    - key encipherment</markup>

<ul class="colist">
<li data-value="1">The certificate will use the <code>selfsigned-issuer</code> cluster issuer we created above.</li>
<li data-value="2">There will be a secret named <code>ca-cert</code> created containing the key and certificate</li>
<li data-value="3">Note that the <code>isCA</code> field is set to <code>true</code> in the body of the spec.</li>
</ul>
<p>The CA issuer that we will create later will also be a <code>ClusterIssuer</code>, so in order for the issuer to find the <code>Certificate</code> above we will create the certificate in the <code>cert-manager</code> namespace, which is where Cert Manager is running.</p>

<markup
lang="bash"

>kubectl -n cert-manager apply -f manifests/ca-cert.yaml</markup>

<p>We can see that the certificate was created and should be ready:</p>

<markup
lang="bash"

>kubectl -n cert-manager get certificate</markup>

<markup
lang="bash"

>NAME             READY   SECRET    AGE
ca-certificate   True    ca-cert   12m</markup>

<p>There will also be a secret named <code>ca-secret</code> created in the <code>cert-manager</code> namespace.
The Secret will contain the certificate and signing key, this will be created when the CA certificate is deployed, and the CA issuer will reference that secret.</p>

</div>

<h3 id="create_ca_issuer">Create the CA issuer.</h3>
<div class="section">
<p>As with the self-signed issuer above, we will create a <code>ClusterIssuer</code> for the CA issuer.</p>

<markup
lang="bash"
title="manifests/ca-cert.yaml"
>apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: ca-issuer
spec:
  ca:
    secretName: ca-cert  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The <code>ca-issuer</code> will use the <code>ca-cert</code> secret created by the <code>ca-certificate</code> <code>Certificate</code> we created above.</li>
</ul>
<p>Create the CA issuer with the following command. As this is a <code>ClusterIssuer</code>, is does not require a namespace.</p>

<markup
lang="bash"

>kubectl apply -f manifests/ca-issuer.yaml</markup>

<p>You can then check that the issuer have been successfully configured by checking the status.</p>

<markup
lang="bash"

>kubectl get clusterissuer</markup>

<p>We should see that both <code>ClusterIssuers</code> we created are present and is ready.</p>

<markup
lang="bash"

>NAME                READY   AGE
ca-issuer           True    22m
selfsigned-issuer   True    31m</markup>

</div>

<h3 id="create_coherence_keystores">Create the Coherence Keys, Certs and KeyStores</h3>
<div class="section">
<p>As the Coherence server, and client in this example, are Java applications they will require Java keystores to hold the certificates. We can use Cert-Manager to create these for us.</p>


<h4 id="_create_a_namespace">Create a Namespace</h4>
<div class="section">
<p>We will run the Coherence cluster in a namespace called <code>coherence-test</code>, so we will first create this:</p>

<markup
lang="bash"

>kubectl create ns coherence-test</markup>

</div>

<h4 id="server_password_secret">Create the Server Keystore Password Secret</h4>
<div class="section">
<p>The keystore will be secured with a password. We will create this password in a <code>Secret</code> so that Cert-Manager can find and use it.
The simplest way to create this secret is with kubectl:</p>

<markup
lang="bash"

>kubectl -n coherence-test create secret generic \
    server-keystore-secret --from-literal=password-key=[your-password]</markup>

<p>&#8230;&#8203;replacing <code>[your-password]</code> with the actual password you want to use.
Resulting in a <code>Secret</code> similar to this:</p>

<markup
lang="bash"
title="manifests/ca-cert.yaml"
>apiVersion: v1
kind: Secret
metadata:
  name: server-keystore-secret
data:
  password-key: "cGFzc3dvcmQ=" <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">In this example the password used is <code>password</code></li>
</ul>
</div>

<h4 id="server_cert">Create the Server Certificate</h4>
<div class="section">
<p>We can now create the server certificate and keystore.</p>

<markup
lang="yaml"
title="manifests/server-keystore.yaml"
>apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: server-keystore
spec:
  issuerRef:
    name: ca-issuer                   <span class="conum" data-value="1" />
    kind: ClusterIssuer
    group: cert-manager.io
  secretName: coherence-server-certs  <span class="conum" data-value="2" />
  keystores:
    jks:
      create: true
      passwordSecretRef:
        key: password-key
        name: server-keystore-secret  <span class="conum" data-value="3" />
  duration: 2160h # 90d
  renewBefore: 360h # 15d
  privateKey:
    size: 2048
    algorithm: RSA
    encoding: PKCS1
  usages:
    - digital signature
    - key encipherment
    - client auth
    - server auth
  commonName: Coherence Certs</markup>

<ul class="colist">
<li data-value="1">The issuer will the <code>ClusterIssuer</code> named <code>ca-issuer</code> that we created above.</li>
<li data-value="2">The keys, certs and keystores will be created in a secret named <code>coherence-server-certs</code></li>
<li data-value="3">The keystore password secret is the <code>Secret</code> named <code>server-keystore-secret</code> we created above</li>
</ul>
<p>We can create the certificate in the <code>coherence-test</code> namespace with the following command:</p>

<markup
lang="bash"

>kubectl -n coherence-test apply -f manifests/server-keystore.yaml</markup>

<p>If we list the certificate in the <code>coherence-test</code> namespace we should see the new certificate and that it is ready.</p>

<markup
lang="bash"

>kubectl -n coherence-test get certificate</markup>

<markup
lang="bash"

>NAME              READY   SECRET                   AGE
server-keystore   True    coherence-server-certs   4s</markup>

<p>If we list the secrets in the <code>coherence-test</code> namespace we should see both the password secret and the keystore secret:</p>

<markup
lang="bash"

>kubectl -n coherence-test get secret</markup>

<markup
lang="bash"

>NAME                     TYPE                 DATA   AGE
coherence-server-certs   kubernetes.io/tls    5      117s
server-keystore-secret   Opaque               1      2m9s</markup>

</div>

<h4 id="client_certs">Create the Client Certificate</h4>
<div class="section">
<p>We can create the certificates and keystores for the client in exactly the same way we did for the server.</p>

<p>Create a password secret for the client keystore:</p>

<markup
lang="bash"

>kubectl -n coherence-test create secret generic \
    client-keystore-secret --from-literal=password-key=[your-password]</markup>

<p>Create the client certificate and keystore.</p>

<markup
lang="yaml"
title="manifests/client-keystore.yaml"
>apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: client-keystore
spec:
  issuerRef:
    name: ca-issuer                   <span class="conum" data-value="1" />
    kind: ClusterIssuer
    group: cert-manager.io
  secretName: coherence-client-certs  <span class="conum" data-value="2" />
  keystores:
    jks:
      create: true
      passwordSecretRef:
        key: password-key
        name: client-keystore-secret  <span class="conum" data-value="3" />
  duration: 2160h # 90d
  renewBefore: 360h # 15d
  privateKey:
    size: 2048
    algorithm: RSA
    encoding: PKCS1
  usages:
    - digital signature
    - key encipherment
    - client auth
  commonName: Coherence Certs</markup>

<ul class="colist">
<li data-value="1">The issuer is the same cluster-wide <code>ca-issuer</code> that we used for the server.</li>
<li data-value="2">The keys, certs and keystores will be created in a secret named <code>coherence-client-certs</code></li>
<li data-value="3">The keystore password secret is the <code>Secret</code> named <code>client-keystore-secret</code> we created above</li>
</ul>
<markup
lang="bash"

>kubectl -n coherence-test apply -f manifests/client-keystore.yaml</markup>

<p>If we list the certificate in the <code>coherence-test</code> namespace we should see the new client certificate and that it is ready.</p>

<markup
lang="bash"

>kubectl -n coherence-test get certificate</markup>

<markup


>NAME              READY   SECRET                   AGE
client-keystore   True    coherence-client-certs   12s
server-keystore   True    coherence-server-certs   2m13s</markup>

</div>
</div>
</div>

<h2 id="coherence">Securing Coherence</h2>
<div class="section">
<p>By this point, you should have installed the Operator and have the four <code>Secrets</code> required, either created by Cert Manager, or manually. Now we can secure Coherence clusters.</p>


<h3 id="images">Build the Test Images</h3>
<div class="section">
<p>This example includes a Maven project that will build a Coherence server and client images with configuration files that allow us to easily demonstrate TLS. To build the images run the following command:</p>

<markup
lang="bash"

>./mvnw clean package jib:dockerBuild</markup>

<p>This will produce two images:</p>

<ul class="ulist">
<li>
<p><code>tls-example-server:1.0.0</code></p>

</li>
<li>
<p><code>tls-example-client:1.0.0</code></p>

</li>
</ul>
<p>These images can run secure or insecure depending on various system properties passed in at runtime.</p>

</div>

<h3 id="socket_provider">Configure a Socket Provider</h3>
<div class="section">
<p>When configuring Coherence to use TLS, we need to configure a socket provider that Coherence can use to create secure socket. We then tell Coherence to use this provider in various places, such as Extend connections, cluster member TCMP connections etc.
This configuration is typically done by adding the provider configuration to the Coherence operational configuration override file.</p>

<p>The Coherence documentation has a lot of details on configuring socket providers in the section on <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/secure/using-ssl-secure-communication.html#GUID-21CBAF48-BA78-4373-AC90-BF668CF31776">Using SSL Secure Communication</a></p>

<p>Below is an example that we will use on the server cluster members</p>

<markup
lang="xml"
title="src/main/resources/tls-coherence-override.xml"
>&lt;coherence xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xmlns="http://xmlns.oracle.com/coherence/coherence-operational-config"
    xsi:schemaLocation="http://xmlns.oracle.com/coherence/coherence-operational-config coherence-operational-config.xsd"&gt;
  &lt;cluster-config&gt;
    &lt;socket-providers&gt;
      &lt;socket-provider id="tls"&gt;
        &lt;ssl&gt;
          &lt;protocol&gt;TLS&lt;/protocol&gt;
          &lt;identity-manager&gt;
            &lt;key-store&gt;
              &lt;url system-property="coherence.tls.keystore"/&gt;
              &lt;password-provider&gt;
                &lt;class-name&gt;com.oracle.coherence.k8s.FileBasedPasswordProvider&lt;/class-name&gt;
                  &lt;init-params&gt;
                    &lt;init-param&gt;
                      &lt;param-type&gt;String&lt;/param-type&gt;
                      &lt;param-value system-property="coherence.tls.keystore.password"&gt;/empty.txt&lt;/param-value&gt;
                    &lt;/init-param&gt;
                &lt;/init-params&gt;
              &lt;/password-provider&gt;
            &lt;/key-store&gt;
            &lt;password-provider&gt;
              &lt;class-name&gt;com.oracle.coherence.k8s.FileBasedPasswordProvider&lt;/class-name&gt;
              &lt;init-params&gt;
                &lt;init-param&gt;
                  &lt;param-type&gt;String&lt;/param-type&gt;
                  &lt;param-value system-property="coherence.tls.key.password"&gt;/empty.txt&lt;/param-value&gt;
              &lt;/init-param&gt;
            &lt;/init-params&gt;
          &lt;/password-provider&gt;
          &lt;/identity-manager&gt;
          &lt;trust-manager&gt;
            &lt;key-store&gt;
              &lt;url system-property="coherence.tls.truststore"/&gt;
              &lt;password-provider&gt;
                &lt;class-name&gt;com.oracle.coherence.k8s.FileBasedPasswordProvider&lt;/class-name&gt;
                &lt;init-params&gt;
                  &lt;init-param&gt;
                    &lt;param-type&gt;String&lt;/param-type&gt;
                    &lt;param-value system-property="coherence.tls.truststore.password"&gt;/empty.txt&lt;/param-value&gt;
                  &lt;/init-param&gt;
                &lt;/init-params&gt;
              &lt;/password-provider&gt;
            &lt;/key-store&gt;
          &lt;/trust-manager&gt;
        &lt;/ssl&gt;
      &lt;/socket-provider&gt;
    &lt;/socket-providers&gt;
  &lt;/cluster-config&gt;
&lt;/coherence&gt;</markup>

<p>The file above has a number of key parts.</p>

<p>We must give the provider a name so that we can refer to it in other configuration.
This is done by setting the <code>id</code> attribute of the <code>&lt;socket-provider&gt;</code> element. In this case we name the provider "tls" in <code>&lt;socket-provider id="tls"&gt;</code>.</p>

<p>We set the <code>&lt;protocol&gt;</code> element to TLS to tell Coherence that this is a TLS socket.</p>

<p>We need to set the keystore URL. If we always used a common location, we could hard code it in the configuration. In this case we will configure the <code>&lt;keystore&gt;&lt;url&gt;</code> element to be injected from a system property which we will configure at runtime <code>&lt;url system-property="coherence.tls.keystore"/&gt;</code>.</p>

<p>We obviously do not want hard-coded passwords in our configuration.
In this example we will use a password provider, which is a class implementing the <code>com.tangosol.net.PasswordProvider</code> interface, that can provide the password by reading file.
In this case the file will be the one from the password secret created above that we will mount into the container.</p>

<markup
lang="xml"
title="src/main/resources/server-cache-config.xml"
>&lt;password-provider&gt;
  &lt;class-name&gt;com.oracle.coherence.k8s.FileBasedPasswordProvider&lt;/class-name&gt;
    &lt;init-params&gt;
      &lt;init-param&gt;
        &lt;param-type&gt;String&lt;/param-type&gt;
        &lt;param-value system-property="coherence.tls.keystore.password"/&gt;
      &lt;/init-param&gt;
  &lt;/init-params&gt;
&lt;/password-provider&gt;</markup>

<p>In the snippet above the password file location will be passed in using the
<code>coherence.tls.keystore.password</code> system property.</p>

<p>We declare another password provider for the private key password.</p>

<p>We then declare the configuration for the truststore, which follows the same pattern as the keystore.</p>

<p>The configuration above is included in both of the example images that we built above.</p>

</div>
</div>

<h2 id="tcmp">Secure Cluster Membership</h2>
<div class="section">
<p>Now we have a "tls" socket provider we can use it to secure Coherence. The Coherence documentation has a section on <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/secure/using-ssl-secure-communication.html#GUID-21CBAF48-BA78-4373-AC90-BF668CF31776">Securing Coherence TCMP with TLS</a>.
Securing communication between cluster members is very simple, we just set the <code>coherence.socketprovider</code> system property to the name of the socket provider we want to use. In our case this will be the "tls" provider we configured above, so we would use <code>-Dcoherence.socketprovider=tls</code></p>

<p>The yaml below is a <code>Coherence</code> resource that will cause the Operator to create a three member Coherence cluster.</p>

<markup
lang="yaml"
title="manifests/coherence-cluster.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: tls-cluster
spec:
  replicas: 3
  image: tls-example-server:1.0.0    <span class="conum" data-value="1" />
  cluster: test-cluster              <span class="conum" data-value="2" />
  coherence:
    overrideConfig: tls-coherence-override.xml  <span class="conum" data-value="3" />
    cacheConfig: server-cache-config.xml        <span class="conum" data-value="4" />
  jvm:
    args:
      - -Dcoherence.socketprovider=tls  <span class="conum" data-value="5" />
      - -Dcoherence.tls.keystore=file:/coherence/certs/keystore.jks
      - -Dcoherence.tls.keystore.password=file:/coherence/certs/credentials/password-key
      - -Dcoherence.tls.key.password=file:/coherence/certs/credentials/password-key
      - -Dcoherence.tls.truststore=file:/coherence/certs/truststore.jks
      - -Dcoherence.tls.truststore.password=file:/coherence/certs/credentials/password-key
  secretVolumes:
    - mountPath: coherence/certs             <span class="conum" data-value="6" />
      name: coherence-server-certs
    - mountPath: coherence/certs/credentials
      name: server-keystore-secret
  ports:
    - name: extend  <span class="conum" data-value="7" />
      port: 20000
    - name: grpc
      port: 1408
    - name: management
      port: 30000
    - name: metrics
      port: 9612</markup>

<ul class="colist">
<li data-value="1">The image name is the server image built from this example project</li>
<li data-value="2">We specify a cluster name because we want to be able to demonstrate other Coherence deployments can or cannot join this cluster, so their yaml files will use this same cluster name.</li>
<li data-value="3">We set the Coherence override file to the file containing the "tls" socket provider configuration.</li>
<li data-value="4">We use a custom cache configuration file that has an Extend proxy that we can secure later.</li>
<li data-value="5">We set the <code>coherence.socketprovider</code> system property to use the "tls" provider, we also set a number of other properties that will set the locations of the keystores and password files to map to the secret volume mounts.</li>
<li data-value="6">We mount the certificate and password secrets to volumes</li>
<li data-value="7">We expose some ports for clients which we will use later, and for management, so we can enquire on the cluster state using REST.</li>
</ul>
<p>Install the yaml above into the <code>coherence-test</code> namespace:</p>

<markup
lang="bash"

>kubectl -n coherence-test apply -f manifests/coherence-cluster.yaml</markup>

<p>If we list the Pods in the <code>coherence-test</code> namespace then after a minute or so there should be three ready Pods.</p>

<markup
lang="bash"

>kubectl -n coherence-test get pods</markup>

<markup
lang="bash"

>NAME             READY   STATUS    RESTARTS   AGE
tls-cluster-0    1/1     Running   0          88s
tls-cluster-1    1/1     Running   0          88s
tls-cluster-2    1/1     Running   0          88s</markup>


<h3 id="_port_forward_to_the_rest_management_port">Port Forward to the REST Management Port</h3>
<div class="section">
<p>Remember that we exposed a number of ports in our Coherence cluster, one of these was REST management on port <code>30000</code>.
We can use this along with <code>curl</code> to enquire about the cluster state.
We need to use <code>kubectl</code> to forward a local port to one of the Coherence Pods.</p>

<p>Open another terminal session and run the following command:</p>

<markup
lang="bash"

>kubectl -n coherence-test port-forward tls-cluster-0 30000:30000</markup>

<p>This will forward port <code>30000</code> on the local machine (e.g. your dev laptop) to the <code>tls-cluster-0</code> Pod.</p>

<p>We can now obtain the cluster state from the REST endpoint with the following command:</p>

<markup
lang="bash"

>curl -X GET http://127.0.0.1:30000/management/coherence/cluster</markup>

<p>or if you have the <a id="" title="" target="_blank" href="https://stedolan.github.io/jq/">jq</a> utility we can pretty print the json output:</p>

<markup
lang="bash"

>curl -X GET http://127.0.0.1:30000/management/coherence/cluster | jq</markup>

<p>We will see json something like this:</p>

<markup
lang="json"

>{
  "links": [
  ],
  "clusterSize": 3,      <span class="conum" data-value="1" />
  "membersDeparted": [],
  "memberIds": [
    1,
    2,
    3
  ],
  "oldestMemberId": 1,
  "refreshTime": "2021-03-07T12:27:20.193Z",
  "licenseMode": "Development",
  "localMemberId": 1,
  "version": "22.06",
  "running": true,
  "clusterName": "test-cluster",
  "membersDepartureCount": 0,
  "members": [                     <span class="conum" data-value="2" />
    "Member(Id=1, Timestamp=2021-03-07 12:24:32.982, Address=10.244.1.6:38271, MachineId=17483, Location=site:zone-two,rack:two,machine:operator-worker2,process:33,member:tls-cluster-1, Role=tls-cluster)",
    "Member(Id=2, Timestamp=2021-03-07 12:24:36.572, Address=10.244.2.5:36139, MachineId=21703, Location=site:zone-one,rack:one,machine:operator-worker,process:35,member:tls-cluster-0, Role=tls-cluster)",
    "Member(Id=3, Timestamp=2021-03-07 12:24:36.822, Address=10.244.1.7:40357, MachineId=17483, Location=site:zone-two,rack:two,machine:operator-worker2,process:34,member:tls-cluster-2, Role=tls-cluster)"
  ],
  "type": "Cluster"
}</markup>

<ul class="colist">
<li data-value="1">We can see that the cluster size is three.</li>
<li data-value="2">The member list shows details of the three Pods in the cluster</li>
</ul>
</div>

<h3 id="_start_non_tls_cluster_members">Start Non-TLS Cluster Members</h3>
<div class="section">
<p>To demonstrate that the cluster is secure we can start another cluster with yaml that does not enable TLS.</p>

<markup
lang="yaml"
title="manifests/coherence-cluster-no-tls.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: no-tls-cluster
spec:
  replicas: 3
  image: tls-example-server:1.0.0     <span class="conum" data-value="1" />
  cluster: test-cluster               <span class="conum" data-value="2" />
  coherence:
    cacheConfig: server-cache-config.xml
  ports:
    - name: extend
      port: 20000
    - name: grpc
      port: 1408
    - name: management
      port: 30000
    - name: metrics
      port: 9612</markup>

<ul class="colist">
<li data-value="1">This <code>Coherence</code> resource uses the same server image as the secure cluster</li>
<li data-value="2">This <code>Coherence</code> resource also uses the same cluster name as the secure cluster, <code>test-cluster</code>, so it should attempt to join with the secure cluster.
If the existing cluster is not secure, we will end up with a cluster of six members.</li>
</ul>
<p>Install the yaml above into the <code>coherence-test</code> namespace:</p>

<markup
lang="bash"

>kubectl -n coherence-test apply -f manifests/coherence-cluster-no-tls.yaml</markup>

<p>If we list the Pods in the <code>coherence-test</code> namespace then after a minute or so there should be three ready Pods.</p>

<markup
lang="bash"

>kubectl -n coherence-test get pods</markup>

<markup
lang="bash"

>NAME                READY   STATUS    RESTARTS   AGE
tls-cluster-0       1/1     Running   0          15m
tls-cluster-1       1/1     Running   0          15m
tls-cluster-2       1/1     Running   0          15m
no-tls-cluster-0    1/1     Running   0          78s
no-tls-cluster-1    1/1     Running   0          78s
no-tls-cluster-2    1/1     Running   0          78s</markup>

<p>There are six pods running, but they have not formed a six member cluster.
If we re-run the curl command to query the REST management endpoint of the secure cluster we will see that the cluster size is still three:</p>

<markup
lang="bash"

>curl -X GET http://127.0.0.1:30000/management/coherence/cluster -s | jq '.clusterSize'</markup>

<p>What happens is that the non-TLS members have effectively formed their own cluster of three members, but have not been able to form a cluster with the TLS enabled members.</p>

</div>

<h3 id="_cleanup">Cleanup</h3>
<div class="section">
<p>After trying the example, remove both clusters with the corresponding <code>kubectl delete</code> commands so that they do not interfere with the next example.</p>

<markup
lang="bash"

>kubectl -n coherence-test delete -f manifests/coherence-cluster-no-tls.yaml

kubectl -n coherence-test delete -f manifests/coherence-cluster.yaml</markup>

</div>

<h3 id="extend">Secure Extend Connections</h3>
<div class="section">
<p>A common connection type to secure are client connections into the cluster from Coherence Extend clients. The Coherence documentation contains details on <a id="" title="" target="_blank" href="https://docs.oracle.com/en/middleware/standalone/coherence/14.1.1.0/secure/using-ssl-secure-communication.html#GUID-0F636928-8731-4228-909C-8B8AB09613DB">Using SSL to Secure Extend Client Communication</a> for more in-depth details.</p>

<p>As with securing TCMP, we can specify a socket provider in the Extend proxy configuration in the server&#8217;s cache configuration file and also in the remote scheme in the client&#8217;s cache configuration. In this example we will use exactly the same TLS socket provider configuration that we created above. The only difference being the name of the <code>PasswordProvider</code> class used by the client. At the time of writing this, Coherence does not include an implementation of <code>PasswordProvider</code> that reads from a file. The Coherence Operator injects one into the classpath of the server, but our simple client is not managed by the Operator. We have added a simple <code>FileBasedPasswordProvider</code> class to the client code in this example.</p>


<h4 id="_secure_the_proxy">Secure the Proxy</h4>
<div class="section">
<p>To enable TLS for an Extend proxy, we can just specify the name of the socket provider that we want to use in the <code>&lt;proxy-scheme&gt;</code> in the server&#8217;s cache configuration file.</p>

<p>The snippet of configuration below is taken from the <code>server-cache-config.xml</code> file in the example source.</p>

<markup
lang="xml"
title="src/main/resources/server-cache-config.xml"
>&lt;proxy-scheme&gt;
    &lt;service-name&gt;Proxy&lt;/service-name&gt;
    &lt;acceptor-config&gt;
        &lt;tcp-acceptor&gt;
            &lt;socket-provider system-property="coherence.extend.socket.provider"/&gt;       <span class="conum" data-value="1" />
            &lt;local-address&gt;
                &lt;address system-property="coherence.extend.address"&gt;0.0.0.0&lt;/address&gt;   <span class="conum" data-value="2" />
                &lt;port system-property="coherence.extend.port"&gt;20000&lt;/port&gt;              <span class="conum" data-value="3" />
            &lt;/local-address&gt;
        &lt;/tcp-acceptor&gt;
    &lt;/acceptor-config&gt;
    &lt;load-balancer&gt;client&lt;/load-balancer&gt;
    &lt;autostart&gt;true&lt;/autostart&gt;
&lt;/proxy-scheme&gt;</markup>

<ul class="colist">
<li data-value="1">The <code>&lt;socket-provider&gt;</code> element is empty by default, but is configured to be set from the system property named <code>coherence.extend.socket.provider</code>. This means that by default, Extend will run without TLS. If we start the server with the system property set to "tls", the name of our socket provider, then the proxy will use TLS.</li>
<li data-value="2">The Extend proxy will bind to all local addresses.</li>
<li data-value="3">The Extend proxy service will bind to port 20000.</li>
</ul>
<p>We add the additional <code>coherence.extend.socket.provider</code> system property to the <code>spec.jvm.args</code> section of the Coherence resource yaml we will use to deploy the server. The yaml below is identical to the yaml we used above to secure TCMP, but with the addition of the <code>coherence.extend.socket.provider</code> property.</p>

<markup
lang="yaml"
title="coherence-cluster-extend.yaml"
>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: tls-cluster
spec:
  replicas: 3
  image: tls-example-server:1.0.0
  cluster: test-cluster
  coherence:
    cacheConfig: server-cache-config.xml
    overrideConfig: tls-coherence-override.xml
  jvm:
    args:
      - -Dcoherence.socketprovider=tls
      - -Dcoherence.extend.socket.provider=tls    <span class="conum" data-value="1" />
      - -Dcoherence.tls.keystore=file:/coherence/certs/keystore.jks
      - -Dcoherence.tls.keystore.password=file:/coherence/certs/credentials/password-key
      - -Dcoherence.tls.key.password=file:/coherence/certs/credentials/password-key
      - -Dcoherence.tls.truststore=file:/coherence/certs/truststore.jks
      - -Dcoherence.tls.truststore.password=file:/coherence/certs/credentials/password-key
  secretVolumes:
    - mountPath: coherence/certs
      name: coherence-server-certs
    - mountPath: coherence/certs/credentials
      name: server-keystore-secret
  ports:
    - name: extend
      port: 20000
    - name: grpc
      port: 1408</markup>

<ul class="colist">
<li data-value="1">The <code>-Dcoherence.extend.socket.provider=tls</code> has been added to enable TLS for the Extend proxy.</li>
</ul>
<p>Installing the yaml above will give us a Coherence cluster that uses TLS for both TCMP inter-cluster communication and for Extend connections.</p>

</div>

<h4 id="_install_the_cluster">Install the Cluster</h4>
<div class="section">
<p>We can install the Coherence cluster defined in the yaml above using <code>kubectl</code>:</p>

<markup
lang="bash"

>kubectl -n coherence-test apply -f manifests/coherence-cluster-extend.yaml</markup>

<p>After a minute or two the three Pods should be ready, which can be confirmed with <code>kubectl</code>.
Because the yaml above declares a port named <code>extend</code> on port <code>20000</code>, the Coherence Operator will create a k8s <code>Service</code> to expose this port. The service name will be the Coherence resource name suffixed with the port name, so in this case <code>tls-cluster-extend</code>. As a <code>Service</code> in k8s can be looked up by DNS, we can use this service name as the host name for the client to connect to.</p>

</div>

<h4 id="_configure_the_extend_client">Configure the Extend Client</h4>
<div class="section">
<p>Just like the server, we can include a socket provider configuration in the override file and configure the name of the socket provider that the client should use in the client&#8217;s cache configuration file. The socket provider configuration is identical to that shown already above (with the different <code>FileBasedPasswordProvider</code> class name).</p>

<p>The Extend client code used in the <code>src/main/java/com/oracle/coherence/examples/k8s/client/Main.java</code> file in this example just starts a Coherence client, then obtains a <code>NamedMap</code>, and in a very long loop just puts data into the map, logging out the keys added. This is very trivial but allows us to see that the client is connected and working (or not).</p>

<p>The snippet of xml below is from the client&#8217;s cache configuration file.</p>

<markup
lang="xml"
title="src/main/resources/client-cache-config.xml"
>&lt;remote-cache-scheme&gt;
    &lt;scheme-name&gt;remote&lt;/scheme-name&gt;
    &lt;service-name&gt;Proxy&lt;/service-name&gt;
    &lt;initiator-config&gt;
        &lt;tcp-initiator&gt;
            &lt;socket-provider system-property="coherence.extend.socket.provider"/&gt;           <span class="conum" data-value="1" />
            &lt;remote-addresses&gt;
                &lt;socket-address&gt;
                    &lt;address system-property="coherence.extend.address"&gt;127.0.0.1&lt;/address&gt; <span class="conum" data-value="2" />
                    &lt;port system-property="coherence.extend.port"&gt;20000&lt;/port&gt;              <span class="conum" data-value="3" />
                &lt;/socket-address&gt;
            &lt;/remote-addresses&gt;
        &lt;/tcp-initiator&gt;
    &lt;/initiator-config&gt;
&lt;/remote-cache-scheme&gt;</markup>

<ul class="colist">
<li data-value="1">The <code>&lt;socket-provider&gt;</code> element is empty by default, but is configured to be set from the system property named <code>coherence.extend.socket.provider</code>. This means that by default, the Extend client will connect without TLS. If we start the client with the system property set to "tls", the name of our socket provider, then the client will use TLS.</li>
<li data-value="2">By default, the Extend client will connect loopback, on <code>127.0.0.1</code> but this can be overridden by setting the <code>coherence.extend.address</code> system property. We will use this when we deploy the client to specify the name of the <code>Service</code> that is used to expose the server&#8217;s Extend port.</li>
<li data-value="3">The Extend client will connect to port 20000. Although this can be overridden with a system property, port 20000 is also the default port used by the server, so there is no need to override it.</li>
</ul>
</div>

<h4 id="_start_an_insecure_client">Start an Insecure Client</h4>
<div class="section">
<p>As a demonstration we can first start a non-TLS client and see what happens. We can create a simple <code>Pod</code> that will run the client image using the yaml below.</p>

<p>One of the features of newer Coherence CE versions is that configuration set via system properties prefixed with <code>coherence.</code> can also be set with corresponding environment variable names. The convention used for the environment variable name is to convert the system property name to uppercase and convert "." characters to "_", so setting the cache configuration file with the <code>coherence.cacheconfig</code> system property can be done using the <code>COHERENCE_CACHECONFIG</code> environment variable.
This makes it simple to set Coherence configuration properties in a Pod yaml using environment variables instead of having to build a custom Java command line.</p>

<markup
lang="yaml"
title="manifests/client-no-tls.yaml"
>apiVersion: v1
kind: Pod
metadata:
  name: client
spec:
  containers:
    - name: client
      image: tls-example-client:1.0.0
      env:
        - name: COHERENCE_CACHECONFIG       <span class="conum" data-value="1" />
          value: client-cache-config.xml
        - name: COHERENCE_EXTEND_ADDRESS    <span class="conum" data-value="2" />
          value: tls-cluster-extend</markup>

<ul class="colist">
<li data-value="1">The client will use the <code>client-cache-config.xml</code> cache configuration file.</li>
<li data-value="2">The <code>COHERENCE_EXTEND_ADDRESS</code> is set to <code>tls-cluster-extend</code>, which is the name of the service exposing the server&#8217;s Extend port and which will be injected into the client&#8217;s cache configuration file, as explained above.</li>
</ul>
<p>We can run the client Pod with the following command:</p>

<markup
lang="bash"

>kubectl -n coherence-test apply -f manifests/client-no-tls.yaml</markup>

<p>If we look at the Pods now in the <code>coherence-test</code> namespace we will see the client running:</p>

<markup
lang="bash"

>$ kubectl -n coherence-test get pod</markup>

<markup
lang="bash"

>NAME            READY   STATUS    RESTARTS   AGE
client          1/1     Running   0          3s
tls-cluster-0   1/1     Running   0          2m8s
tls-cluster-1   1/1     Running   0          2m8s
tls-cluster-2   1/1     Running   0          2m8s</markup>

<p>If we look at the log of the client Pod though we will see a stack trace with the cause:</p>

<markup
lang="bash"

>kubectl -n coherence-test logs client</markup>

<markup


>2021-03-07 12:53:13.481/1.992 Oracle Coherence CE 22.06 &lt;Error&gt; (thread=main, member=n/a): Error while starting service "Proxy": com.tangosol.net.messaging.ConnectionException: could not establish a connection to one of the following addresses: []</markup>

<p>This tells us that the client failed to connect to the cluster, because the client is not using TLS.</p>

<p>We can remove the non-TLS client:</p>

<markup


>kubectl -n coherence-test delete -f manifests/client-no-tls.yaml</markup>

</div>

<h4 id="_start_a_tls_enabled_client">Start a TLS Enabled Client</h4>
<div class="section">
<p>We can now modify the client yaml to run the client with TLS enabled.
The client image already contains the <code>tls-coherence-override.xml</code> file with the configuration for the TLS socket provider.
We need to set the relevant environment variables to inject the location of the keystores and tell Coherence to use the "tls" socket provider for the Extend connection.</p>

<markup
lang="yaml"
title="manifests/client.yaml"
>apiVersion: v1
kind: Pod
metadata:
  name: client
spec:
  containers:
    - name: client
      image: tls-example-client:1.0.0
      env:
        - name: COHERENCE_CACHECONFIG
          value: client-cache-config.xml
        - name: COHERENCE_EXTEND_ADDRESS
          value: tls-cluster-extend
        - name: COHERENCE_OVERRIDE
          value: tls-coherence-override.xml                 <span class="conum" data-value="1" />
        - name: COHERENCE_EXTEND_SOCKET_PROVIDER
          value: tls
        - name: COHERENCE_TLS_KEYSTORE
          value: file:/coherence/certs/keystore.jks
        - name: COHERENCE_TLS_KEYSTORE_PASSWORD
          value: /coherence/certs/credentials/password-key
        - name: COHERENCE_TLS_KEY_PASSWORD
          value: /coherence/certs/credentials/password-key
        - name: COHERENCE_TLS_TRUSTSTORE
          value: file:/coherence/certs/truststore.jks
        - name: COHERENCE_TLS_TRUSTSTORE_PASSWORD
          value: /coherence/certs/credentials/password-key
      volumeMounts:                                         <span class="conum" data-value="2" />
        - name: coherence-client-certs
          mountPath: coherence/certs
        - name: keystore-credentials
          mountPath: coherence/certs/credentials
  volumes:                                                  <span class="conum" data-value="3" />
    - name: coherence-client-certs
      secret:
        defaultMode: 420
        secretName: coherence-client-certs
    - name: keystore-credentials
      secret:
        defaultMode: 420
        secretName: client-keystore-secret</markup>

<ul class="colist">
<li data-value="1">The yaml is identical to the non-TLS client with the addition of the environment variables to configure TLS.</li>
<li data-value="2">We create volume mount points to map the Secret volumes containing the keystores and password to directories in the container</li>
<li data-value="3">We mount the Secrets as volumes</li>
</ul>
<p>We can run the client Pod with the following command:</p>

<markup
lang="bash"

>kubectl -n coherence-test apply -f manifests/client.yaml</markup>

<p>If we now look at the client&#8217;s logs:</p>

<markup
lang="bash"

>kubectl -n coherence-test logs client</markup>

<p>The end of the log should show the messages from the client as it puts each entry into a <code>NamedMap</code>.</p>

<markup


>Put 0
Put 1
Put 2
Put 3
Put 4
Put 5</markup>

<p>So now we have a TLS secured Extend proxy and client.
We can remove the client and test cluster:</p>

<markup
lang="bash"

>kubectl -n coherence-test delete -f manifests/client.yaml

kubectl -n coherence-test delete -f manifests/coherence-cluster-extend.yaml</markup>

</div>
</div>
</div>
</doc-view>
