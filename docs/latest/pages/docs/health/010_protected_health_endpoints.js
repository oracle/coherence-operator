<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Protected Coherence Health Endpoints</dt>
<dd slot="desc"><p>Recent Coherence releases protect the state-changing Health endpoints used to suspend and resume cache services.
Coherence Operator 3.6.0 is the first Operator release that can configure and call these protected endpoints during
safe shutdown and scale-to-zero.</p>

<p>The recommended configuration uses <code>spec.healthMutatorConnection</code> to configure both sides of the connection:</p>

<ul class="ulist">
<li>
<p>the Coherence member exposes <code>PUT /suspend</code> using Basic authentication over TLS; and</p>

</li>
<li>
<p>the Operator reads credentials and a CA certificate from Kubernetes Secrets and sends an authenticated, verified
HTTPS request.</p>

</li>
</ul>
<p>Authentication or TLS failure blocks the transition. The Operator does not retry anonymously, use <code>GET</code>, or disable
certificate verification.</p>
</dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_hardened_health_mutator_behavior">Hardened Health mutator behavior</h2>
<div class="section">
<p>The hardened Coherence Health server changes the mutator endpoints as follows:</p>

<ul class="ulist">
<li>
<p><code>/suspend</code> and <code>/resume</code> accept <code>PUT</code> only. The old <code>GET</code> routes are not available.</p>

</li>
<li>
<p>Mutators are disabled by default. A request returns <code>403</code> until
<code>coherence.health.http.mutators.enabled=true</code> is configured.</p>

</li>
<li>
<p>An enabled mutator must use <code>basic</code>, <code>cert</code>, or <code>cert+basic</code> authentication. There is no unauthenticated mode.</p>

</li>
<li>
<p>Read-only endpoints such as <code>/ready</code>, <code>/started</code>, <code>/live</code>, <code>/healthz</code>, <code>/safe</code>, and <code>/ha</code> remain available for
Kubernetes and Operator health checks without mutator credentials.</p>

</li>
</ul>
<p>An older Operator sends an unauthenticated <code>GET /suspend</code>, so it cannot safely scale a deployment using a hardened
Coherence Health server to zero. Operator 3.6.0 sends <code>PUT</code> and adds the managed and application-owned configuration
options described on this page.</p>

</div>

<h2 id="_compatible_versions">Compatible versions</h2>
<div class="section">
<p>The following are the first releases on each maintained Coherence line that contain the Health mutator hardening. Later
patches on the same line also contain it.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 50%;">
<col style="width: 50%;">
</colgroup>
<thead>
<tr>
<th>Coherence line</th>
<th>First hardened release</th>
</tr>
</thead>
<tbody>
<tr>
<td class="">Oracle Coherence 14.1.1.2206</td>
<td class="">14.1.1.2206.18</td>
</tr>
<tr>
<td class="">Oracle Coherence 14.1.2</td>
<td class="">14.1.2.0.8</td>
</tr>
<tr>
<td class="">Oracle Coherence 15.1.1</td>
<td class="">15.1.1.0.4</td>
</tr>
<tr>
<td class="">Coherence Community Edition 14.1.2</td>
<td class="">14.1.2-0-8</td>
</tr>
<tr>
<td class="">Coherence Community Edition 15.1.1</td>
<td class="">15.1.1-0-4</td>
</tr>
<tr>
<td class="">Coherence Community Edition 22.06</td>
<td class="">22.06.18</td>
</tr>
<tr>
<td class="">Coherence Community Edition feature releases</td>
<td class="">26.07</td>
</tr>
</tbody>
</table>
</div>
<p>Use Coherence Operator 3.6.0 or later with a hardened Coherence release. The Operator image and CRD must be upgraded
together so that the <code>healthMutatorConnection</code> and method-aware <code>suspendProbe.http</code> fields are both available.</p>

</div>

<h2 id="_choose_a_configuration_model">Choose a configuration model</h2>
<div class="section">

<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 25%;">
<col style="width: 37.5%;">
<col style="width: 37.5%;">
</colgroup>
<thead>
<tr>
<th>Model</th>
<th>Use when</th>
<th>Configuration owner</th>
</tr>
</thead>
<tbody>
<tr>
<td class="">Operator-managed Basic authentication and one-way TLS (recommended)</td>
<td class="">The application can use the Coherence Health server and Basic authentication.</td>
<td class="">The Operator configures the Coherence listener, TLS identity, JAAS, kubelet probes, and Operator client from
<code>healthMutatorConnection</code>.</td>
</tr>
<tr>
<td class="">Operator-managed security with an existing JAAS file</td>
<td class="">The application already owns the JVM&#8217;s global JAAS configuration.</td>
<td class="">The application provides the <code>CoherenceREST</code> JAAS entry; the Operator configures the remaining server and client
settings.</td>
</tr>
<tr>
<td class="">Operator-managed security with application-mounted TLS files</td>
<td class="">The application already mounts the member keystore and password files.</td>
<td class="">The application owns the volume; the Operator uses explicit file paths and configures the remaining settings.</td>
</tr>
<tr>
<td class="">Application-owned server with an Operator HTTP action</td>
<td class="">The application already configures a Basic-authenticated HTTPS Health server.</td>
<td class="">The application configures Coherence, TLS, JAAS, and kubelet probes. The Operator configures only the authenticated
<code>PUT</code> request.</td>
</tr>
<tr>
<td class="">Application-owned server with an exec action</td>
<td class="">The server uses client certificates, <code>cert+basic</code>, or another client flow that the Operator HTTP action cannot express.</td>
<td class="">The application image supplies and owns the client command and all server configuration.</td>
</tr>
<tr>
<td class="">No managed security</td>
<td class="">The deployment uses an older, unhardened Health endpoint or a completely application-owned endpoint.</td>
<td class="">Omit <code>healthMutatorConnection</code>. The application is responsible for any endpoint configuration and protection.</td>
</tr>
</tbody>
</table>
</div>
<p>For complete manifests and application instructions for each model, see the
<router-link to="/examples/096_protected_health_endpoints/README">Protected Coherence Health Endpoints Example</router-link>.</p>

</div>

<h2 id="_operator_managed_basic_authentication_over_one_way_tls">Operator-managed Basic authentication over one-way TLS</h2>
<div class="section">
<p>Set <code>spec.healthMutatorConnection.secure</code> to <code>true</code> to enable managed security. The Operator then:</p>

<ul class="ulist">
<li>
<p>enables Coherence Health mutators;</p>

</li>
<li>
<p>selects Basic authentication and the <code>CoherenceREST</code> JAAS context;</p>

</li>
<li>
<p>configures the member&#8217;s <code>HealthSSLProvider</code> and server keystore;</p>

</li>
<li>
<p>reads the Basic username, password, and client CA from namespace-local Secrets;</p>

</li>
<li>
<p>sends <code>PUT /suspend</code> over verified HTTPS; and</p>

</li>
<li>
<p>changes generated readiness, liveness, and startup probes on the Health port to HTTPS.</p>

</li>
</ul>
<p>Managed mode uses one-way TLS. The member presents a certificate to the Operator and kubelet, but does not request a
client certificate. The member therefore needs a keystore, but no truststore. The separate <code>tls.caSecret</code> is the CA
used by the Operator to verify the member certificate.</p>


<h3 id="_required_secrets">Required Secrets</h3>
<div class="section">
<p>Create the following Secrets in the same namespace as the <code>Coherence</code> resource. Use your normal credential and
certificate management process; do not put credentials or private keys directly in the <code>Coherence</code> resource.</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 28.571%;">
<col style="width: 28.571%;">
<col style="width: 42.857%;">
</colgroup>
<thead>
<tr>
<th>Secret</th>
<th>Required keys</th>
<th>Purpose</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>health-auth-v1</code></td>
<td class=""><code>username</code>, <code>password</code></td>
<td class="">Credentials used by the Operator and the <code>CoherenceREST</code> JAAS login module.</td>
</tr>
<tr>
<td class=""><code>health-stores-v1</code></td>
<td class=""><code>server.jks</code>, <code>store-password</code></td>
<td class="">Member certificate and password files. The key and store passwords may use different keys when required.</td>
</tr>
<tr>
<td class=""><code>health-ca-v1</code></td>
<td class=""><code>ca.pem</code></td>
<td class="">PEM CA bundle used by the Operator to verify the member certificate.</td>
</tr>
</tbody>
</table>
</div>
<p>Secret key values are consumed as exact bytes. Avoid an unintended trailing newline in usernames and passwords. Prefer
immutable, versioned Secrets so that a resource update clearly selects a new credential or certificate generation.</p>

</div>

<h3 id="_managed_configuration_example">Managed configuration example</h3>
<div class="section">
<p>The recommended example uses the Operator-provided JAAS login module and Secret-mounted JKS material. The
<router-link to="/examples/096_protected_health_endpoints/README">Protected Coherence Health Endpoints Example</router-link> contains the
complete resource, Secret creation, rendering, application, and verification steps.</p>

<p>The important fields are:</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 33.333%;">
<col style="width: 66.667%;">
</colgroup>
<thead>
<tr>
<th>Field</th>
<th>Meaning</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>secure: true</code></td>
<td class="">Enables the complete managed configuration. Security fields are rejected when this value is <code>false</code>.</td>
</tr>
<tr>
<td class=""><code>basicAuth</code></td>
<td class="">References non-optional username and password keys in namespace-local Secrets.</td>
</tr>
<tr>
<td class=""><code>tls.caSecret</code></td>
<td class="">References the PEM CA bundle used by the Operator client. It is not a member truststore.</td>
</tr>
<tr>
<td class=""><code>tls.serverName</code></td>
<td class="">DNS name verified against the member certificate&#8217;s subject alternative names.</td>
</tr>
<tr>
<td class=""><code>serverTLS</code></td>
<td class="">Configures the member&#8217;s TLS identity. <code>enabled</code> must be <code>true</code> and <code>requireClientCert</code> must be <code>false</code>.</td>
</tr>
<tr>
<td class=""><code>serverTLS.keyStoreProvider</code></td>
<td class="">Optionally selects the Java keystore provider propagated to the managed <code>HealthSSLProvider</code>.</td>
</tr>
<tr>
<td class=""><code>jaas.mode</code></td>
<td class="">Selects <code>packaged</code> (the default) or <code>existing</code> JAAS configuration.</td>
</tr>
</tbody>
</table>
</div>
<p>The managed <code>HealthSSLProvider</code> does not use the <code>serverTLS.trustStore</code>, <code>trustStorePasswordFile</code>, <code>trustStoreType</code>,
<code>trustStoreAlgorithm</code>, or <code>trustStoreProvider</code> fields. They remain accepted for compatibility with the shared TLS API,
but are not used for one-way TLS.</p>

</div>

<h3 id="_use_an_existing_jaas_configuration">Use an existing JAAS configuration</h3>
<div class="section">
<p>Set <code>jaas.mode: existing</code> when the application already sets <code>java.security.auth.login.config</code>. The existing file must
contain a login context named exactly <code>CoherenceREST</code> and must validate the same Secret-backed credentials used by the
Operator client.</p>

<p>Do not use the default <code>packaged</code> mode when another global JAAS file is already configured. Managed mode refuses to
overwrite it. The published example includes the required <code>ConfigMap</code>, volume mount, JVM argument, and
<code>healthMutatorConnection</code> fragment.</p>

</div>

<h3 id="_use_application_mounted_tls_files">Use application-mounted TLS files</h3>
<div class="section">
<p>If the application already mounts the server keystore and password files, omit <code>serverTLS.secrets</code> and specify absolute
container paths instead. The published example includes the required volume and explicit-path configuration.</p>

<p>The Operator still manages the Coherence Health properties, JAAS mode, Basic credentials, client CA, and authenticated
request.</p>

</div>

<h3 id="_supported_application_packaging">Supported application packaging</h3>
<div class="section">
<p>Managed mode supports normal Java/classpath applications and Spring Boot 2 or 3 applications whose executable JAR is
configured with <code>spec.application.springBootFatJar</code>. It does not support Cloud Native Buildpack launchers or a Spring
Boot application without a configured fat JAR.</p>

<p>Managed mode also requires an explicit path for the Operator&#8217;s generated JVM arguments to reach the member JVM. The
supported launcher configurations are:</p>

<ul class="ulist">
<li>
<p>the default Operator-controlled launcher, with no custom <code>spec.application.entryPoint</code> and
<code>spec.application.useImageEntryPoint</code> omitted or <code>false</code>;</p>

</li>
<li>
<p><code>spec.application.useImageEntryPoint: true</code> with <code>useJdkJavaOptions</code> omitted or <code>true</code>, which supplies the generated
argument-file reference through <code>JDK_JAVA_OPTIONS</code>; or</p>

</li>
<li>
<p><code>spec.application.useImageEntryPoint: true</code> with <code>useJdkJavaOptions: false</code> and a non-empty
<code>alternateJdkJavaOptions</code> environment variable name.</p>

</li>
</ul>
<p>The last configuration is an explicit application contract: the image entry point must read the named environment
variable and pass its generated argument-file reference to the actual member JVM. The Operator cannot verify that the
application honors this contract. A custom <code>application.entryPoint</code> is supported only when
<code>useImageEntryPoint: true</code> also selects one of the two environment-variable delivery paths above.</p>

<p>The Operator rejects an unsupported launcher combination before reconciling the workload. This validation applies only
to <code>healthMutatorConnection.secure: true</code>; omitted or insecure managed Health configuration retains the existing
launcher behavior.</p>

<p>Managed mode requires the Coherence Health server. Do not set <code>COHERENCE_OPERATOR_HEALTH_CHECK=true</code>; that selects the
legacy Operator Health server and the Pod configuration fails closed.</p>

</div>
</div>

<h2 id="application_owned">Application-owned Health security</h2>
<div class="section">
<p>Use an application-owned configuration when the application already configures its Health listener, needs certificate
authentication, or cannot use the managed launcher arrangement. Omit <code>healthMutatorConnection</code>; do not combine managed
and application-owned server properties.</p>

<p>The application must configure all of the following:</p>

<ul class="ulist">
<li>
<p><code>coherence.health.http.mutators.enabled=true</code>;</p>

</li>
<li>
<p><code>coherence.health.http.auth=basic</code>, <code>cert</code>, or <code>cert+basic</code>;</p>

</li>
<li>
<p>an HTTPS socket provider selected by <code>coherence.health.http.provider</code>;</p>

</li>
<li>
<p>the socket provider&#8217;s identity and, for client-certificate authentication, trust configuration;</p>

</li>
<li>
<p>a JAAS <code>CoherenceREST</code> context when Basic authentication is used; and</p>

</li>
<li>
<p>explicit readiness, liveness, and startup probes that can use the application&#8217;s TLS contract, because the Operator
does not rewrite probes in application-owned mode.</p>

</li>
</ul>
<p>The read-only endpoints do not require mutator authentication. With one-way TLS, Kubernetes HTTP probes can call them
over HTTPS. If the listener requires a client certificate, use an exec or other probe capable of presenting it.</p>


<h3 id="_basic_authentication_with_an_operator_http_action">Basic authentication with an Operator HTTP action</h3>
<div class="section">
<p>For a Basic-authenticated endpoint, configure a method-aware <code>suspendProbe.http</code>. The Operator reads credentials and
the CA from Secrets and performs a verified HTTPS request. The published example includes both a complete
application-owned Basic-over-TLS resource and a reusable Operator-client fragment.</p>

<p>The application must independently configure matching Basic credentials, TLS material, and Coherence Health system
properties. The Operator HTTP action supports <code>GET</code> and <code>PUT</code>, verified server TLS, Secret-backed Basic authentication,
and non-secret headers. Use <code>PUT</code> for hardened <code>/suspend</code>.</p>

</div>

<h3 id="_certificate_authentication_and_other_custom_clients">Certificate authentication and other custom clients</h3>
<div class="section">
<p>Coherence also supports <code>cert</code> and <code>cert+basic</code> mutator authentication. These modes require a TLS client certificate.
The Operator HTTP action does not present one, so configure an application-owned exec action or another
application-owned client that can satisfy the server&#8217;s authentication contract.</p>

<p>The exec fragment in the published example demonstrates the resource shape but uses Basic authentication. Replace its
command and mounted client material for the selected certificate contract. The application image must contain the
required client program. Exec actions are less portable than the managed HTTP client and provide fewer sanitized
diagnostics.</p>

<p>An empty custom probe or a gRPC-only probe is retained as a successful no-op for compatibility with existing
resources. The Operator emits a <code>ProbeActionNoOp</code> Warning event when it encounters either form. No request or command
is sent, so the result does not prove StatusHA or that services were suspended. Configure an Operator-supported
<code>http</code>, <code>httpGet</code>, <code>exec</code>, or <code>tcpSocket</code> action instead.</p>

</div>
</div>

<h2 id="_disable_or_bypass_managed_security">Disable or bypass managed security</h2>
<div class="section">
<p>Omitting <code>healthMutatorConnection</code>, or setting only <code>healthMutatorConnection.secure: false</code>, disables the Operator&#8217;s
managed Health security configuration. The Operator then sends the default <code>PUT /suspend</code> over HTTP unless a custom
<code>suspendProbe</code> is configured. These two disabled forms generate equivalent Pod templates, so switching between them
does not by itself replace Pods. Switching from <code>secure: true</code> to either disabled form does change the template and
requires a rollout to remove the secure-generation annotation, credentials, mounts, JVM configuration, and HTTPS probe
wiring. Enabling managed security makes the corresponding template changes and also requires a rollout.</p>

<p>This does not disable the security layer in a hardened Coherence release. Hardened Coherence keeps mutators disabled
until they are explicitly enabled and always requires one of its supported authentication modes. Consequently,
<code>secure: false</code> by itself cannot suspend a hardened member.</p>

<p>Setting <code>spec.suspendServicesOnShutdown: false</code> skips service suspension entirely. This is an operational bypass, not
a security setting. It can make scale-to-zero or shutdown unsafe, particularly with persistence, and should not be used
to work around an authentication or TLS configuration error.</p>

</div>

<h2 id="_upgrade_an_existing_deployment">Upgrade an existing deployment</h2>
<div class="section">
<p>Upgrade the Operator deployment and CRD to 3.6.0 or later before rolling a hardened Coherence image into an existing
deployment. Add the selected Health security configuration as part of the Coherence image rollout. Upgrading only the
Operator does not enable security inside existing member Pods, and upgrading only Coherence leaves an older Operator
unable to authenticate or use the required <code>PUT</code> method.</p>

<div class="admonition warning">
<p class="admonition-textlabel">Warning</p>
<p ><p>The 3.6.0 CRD limits <code>spec.actions</code> to 128 entries. Installing the CRD does not remove or retroactively reject an
existing <code>Coherence</code> resource with more than 128 actions, but Kubernetes will reject subsequent updates until the
submitted resource contains no more than 128 actions. Check existing resources and reduce longer action lists before
the upgrade or as the first change after upgrading the CRD.</p>
</p>
</div>
<p>For Operator-managed mode, confirm before the rollout that:</p>

<ul class="ulist">
<li>
<p>all referenced Secrets and keys exist in the <code>Coherence</code> resource namespace;</p>

</li>
<li>
<p>the member certificate is valid for <code>tls.serverName</code>;</p>

</li>
<li>
<p>the member keystore and password files are readable by the Coherence container;</p>

</li>
<li>
<p>an existing JAAS file contains the exact <code>CoherenceREST</code> context when <code>jaas.mode: existing</code> is selected; and</p>

</li>
<li>
<p>custom <code>coherence.health.*</code> JVM properties are removed because they conflict with managed mode.</p>

</li>
</ul>
<p>Keep the previous image, resource definition, and Secret generation available until the new Pods are ready and safe
shutdown has been tested. Rotate credentials and certificates by creating a new immutable Secret generation and
updating all related references together.</p>

</div>

<h2 id="_verify_the_configuration">Verify the configuration</h2>
<div class="section">
<p>After the rollout, verify these observable results:</p>


<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 40%;">
<col style="width: 60%;">
</colgroup>
<thead>
<tr>
<th>Check</th>
<th>Expected result</th>
</tr>
</thead>
<tbody>
<tr>
<td class="">Generated kubelet probes in managed mode</td>
<td class="">Readiness, liveness, and startup probes that target the Health port use <code>HTTPS</code>.</td>
</tr>
<tr>
<td class="">Unauthenticated <code>GET /ready</code></td>
<td class=""><code>200</code> when the member is ready. Read-only probes do not require mutator credentials.</td>
</tr>
<tr>
<td class=""><code>PUT /suspend</code> while mutators are disabled</td>
<td class=""><code>403</code>.</td>
</tr>
<tr>
<td class="">Unauthenticated or invalid-credential <code>PUT /suspend</code></td>
<td class=""><code>401</code> when Basic authentication is configured.</td>
</tr>
<tr>
<td class="">Authenticated <code>PUT /suspend</code></td>
<td class=""><code>200</code>. Send <code>PUT /resume</code> after a direct manual test.</td>
</tr>
<tr>
<td class="">Invalid CA, certificate name, credential, or Secret reference during scale-to-zero</td>
<td class="">The transition remains blocked and existing members remain running.</td>
</tr>
</tbody>
</table>
</div>
<p>Use application-level write and read checks when validating persistence recovery. Ready Pods and unchanged PVC names do
not by themselves prove that application data was recovered.</p>

</div>

<h2 id="_troubleshooting">Troubleshooting</h2>
<div class="section">

<div class="table__overflow elevation-1  ">
<table class="datatable table">
<colgroup>
<col style="width: 40%;">
<col style="width: 60%;">
</colgroup>
<thead>
<tr>
<th>Observation</th>
<th>Action</th>
</tr>
</thead>
<tbody>
<tr>
<td class=""><code>403</code> from <code>PUT /suspend</code></td>
<td class="">Mutators are disabled. In managed mode, confirm that the member Pod contains the current
<code>healthMutatorConnection</code> generation and started successfully. In application-owned mode, set
<code>coherence.health.http.mutators.enabled=true</code>.</td>
</tr>
<tr>
<td class=""><code>401</code> from <code>PUT /suspend</code></td>
<td class="">Confirm that the Operator client and <code>CoherenceREST</code> JAAS context use the same Secret keys and exact credential bytes.</td>
</tr>
<tr>
<td class="">TLS name or trust failure</td>
<td class="">Confirm <code>tls.serverName</code>, the certificate subject alternative names, and the PEM bundle in <code>tls.caSecret</code>. Do not
disable certificate verification. A member-side truststore does not repair Operator-client trust.</td>
</tr>
<tr>
<td class="">Pod configuration rejects the selected Health server</td>
<td class="">Remove <code>COHERENCE_OPERATOR_HEALTH_CHECK=true</code>. Managed mode requires the Coherence Health server.</td>
</tr>
<tr>
<td class="">Pod configuration reports a direct-property conflict</td>
<td class="">Remove application-supplied <code>coherence.health.<strong></code> and <code>coherence.health.security.</strong></code> JVM properties, or omit
<code>healthMutatorConnection</code> and use the application-owned model.</td>
</tr>
<tr>
<td class="">Scale-to-zero remains blocked</td>
<td class="">Treat the block as fail-closed. Correct the authentication, TLS, Secret, or runtime configuration while members remain
running, or restore <code>spec.replicas</code> to a non-zero value. Do not delete the <code>Coherence</code> resource or its namespace while
the transition is blocked: desired replicas can be zero while members are still running, and deletion can terminate
those members without the expected suspension. Wait until scale-to-zero completes or the restored deployment is healthy
before deleting it.</td>
</tr>
</tbody>
</table>
</div>
<p>Do not add an anonymous or <code>GET</code> fallback, combine a custom <code>suspendProbe</code> with
<code>healthMutatorConnection.secure: true</code>, or disable safe suspension to hide a security configuration failure.</p>

</div>
</doc-view>
