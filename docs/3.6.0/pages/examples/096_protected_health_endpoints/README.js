<doc-view>

<v-layout row wrap>
<v-flex xs12 sm10 lg10>
<v-card class="section-def" v-bind:color="$store.state.currentColor">
<v-card-text class="pa-3">
<v-card class="section-def__card">
<v-card-text>
<dl>
<dt slot=title>Protected Coherence Health Endpoints Example</dt>
<dd slot="desc"><p>This example shows how to configure the Coherence Operator to call a protected Coherence Health mutator during safe
shutdown and scale-to-zero. It covers the recommended Operator-managed configuration and the application-owned
alternatives.</p>

<p>See the <router-link to="/docs/health/010_protected_health_endpoints">Protected Coherence Health Endpoints</router-link> documentation for
compatible versions, configuration ownership, upgrade guidance, security constraints, and troubleshooting.</p>

<div class="admonition tip">
<p class="admonition-textlabel">Tip</p>
<p ><p><img src="./images/GitHub-Mark-32px" alt="GitHub Mark 32px" />
 The complete source for this example is in the
<a id="" title="" target="_blank" href="https://github.com/oracle/coherence-operator/tree/main/examples/096_protected_health_endpoints">Coherence Operator GitHub</a> repository.</p>
</p>
</div></dd>
</dl>
</v-card-text>
</v-card>
</v-card-text>
</v-card>
</v-flex>
</v-layout>

<h2 id="_configuration_choices">Configuration choices</h2>
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
<th>Configuration</th>
<th>Use when</th>
<th>Example</th>
</tr>
</thead>
<tbody>
<tr>
<td class="">Managed security with packaged JAAS (recommended)</td>
<td class="">The application can use the Coherence Health server and Basic authentication. The Operator owns the member and client
configuration.</td>
<td class=""><code>manifests/managed-packaged-jaas.yaml</code></td>
</tr>
<tr>
<td class="">Managed security with an existing JAAS file</td>
<td class="">The application already owns the JVM&#8217;s global JAAS configuration and supplies the required <code>CoherenceREST</code> context.</td>
<td class=""><code>manifests/managed-existing-jaas.yaml</code></td>
</tr>
<tr>
<td class="">Managed security with application-mounted TLS files</td>
<td class="">The application already mounts the member keystore and its password files.</td>
<td class=""><code>manifests/managed-mounted-store-paths.yaml</code></td>
</tr>
<tr>
<td class="">Application-owned Basic-over-TLS server</td>
<td class="">The application owns the Coherence Health system properties, TLS, JAAS, and kubelet probes. The Operator owns only
the authenticated suspension request.</td>
<td class=""><code>manifests/application-owned-basic.yaml</code></td>
</tr>
<tr>
<td class="">Application-owned server with an HTTP action</td>
<td class="">The application already owns a Basic-authenticated HTTPS server and needs only the Operator-client fragment.</td>
<td class=""><code>manifests/application-owned-http-action.yaml</code></td>
</tr>
<tr>
<td class="">Application-owned server with a custom action</td>
<td class="">The application needs a client flow that the Operator&#8217;s HTTP action cannot express.</td>
<td class=""><code>manifests/application-owned-exec-action.yaml</code></td>
</tr>
</tbody>
</table>
</div>
<p>Use the managed configuration unless the application already owns the protected Health server. Do not combine
<code>healthMutatorConnection.secure: true</code> with a custom <code>suspendProbe</code> or application-supplied <code>coherence.health.*</code>
properties.</p>

</div>

<h2 id="_prerequisites">Prerequisites</h2>
<div class="section">
<p>Before using this example:</p>

<ul class="ulist">
<li>
<p>install Coherence Operator 3.6.0 or later and its matching CRD;</p>

</li>
<li>
<p>select a Coherence release that supports protected <code>PUT</code> Health mutators;</p>

</li>
<li>
<p>make the application image available to every Kubernetes node, preferably by immutable digest;</p>

</li>
<li>
<p>obtain a server keystore containing the member private key and certificate;</p>

</li>
<li>
<p>obtain the keystore and key password files; and</p>

</li>
<li>
<p>obtain the PEM CA bundle that the Operator will use to verify the member certificate.</p>

</li>
</ul>
<p>The member certificate must contain <code>health.example</code>, or the value selected for <code>tls.serverName</code>, in its subject
alternative names. The managed provider uses one-way TLS and does not require a member truststore.</p>

<p>Run the commands below from the root of the Operator repository. Select the context and namespace explicitly:</p>

<markup
lang="bash"

>export CONTEXT='YOUR_KUBERNETES_CONTEXT'
export NAMESPACE='YOUR_NAMESPACE'
export COHERENCE_IMAGE='registry.example/application@sha256:REPLACE_WITH_DIGEST'

kubectl --context "$CONTEXT" create namespace "$NAMESPACE" \
  --dry-run=client -o yaml | kubectl --context "$CONTEXT" apply -f -</markup>

</div>

<h2 id="_create_the_secrets">Create the Secrets</h2>
<div class="section">
<p>Create immutable, versioned Secrets in the same namespace as the <code>Coherence</code> resource. Do not put credentials or
private keys directly in the resource.</p>

<p>The following example creates a high-entropy password without a trailing newline and copies externally issued TLS
material into a private temporary directory. Replace the paths under <code>/secure/material</code> with the locations supplied by
the approved certificate and secret-management process.</p>

<markup
lang="bash"

>umask 077
export MATERIAL="$(mktemp -d "${TMPDIR:-/tmp}/coherence-health.XXXXXX")"
trap 'rm -f "$MATERIAL/username" "$MATERIAL/password" "$MATERIAL/server.jks" "$MATERIAL/store-password" "$MATERIAL/ca.pem"; rmdir "$MATERIAL" 2&gt;/dev/null || true' EXIT

printf '%s' 'health-operator' &gt; "$MATERIAL/username"
openssl rand -hex 32 | tr -d '\n' &gt; "$MATERIAL/password"
install -m 600 /secure/material/server.jks "$MATERIAL/server.jks"
install -m 600 /secure/material/store-password "$MATERIAL/store-password"
install -m 600 /secure/material/ca.pem "$MATERIAL/ca.pem"

kubectl --context "$CONTEXT" -n "$NAMESPACE" create secret generic health-auth-v1 \
  --from-file=username="$MATERIAL/username" \
  --from-file=password="$MATERIAL/password"

kubectl --context "$CONTEXT" -n "$NAMESPACE" create secret generic health-stores-v1 \
  --from-file=server.jks="$MATERIAL/server.jks" \
  --from-file=store-password="$MATERIAL/store-password"

kubectl --context "$CONTEXT" -n "$NAMESPACE" create secret generic health-ca-v1 \
  --from-file=ca.pem="$MATERIAL/ca.pem"

for secret in health-auth-v1 health-stores-v1 health-ca-v1; do
  kubectl --context "$CONTEXT" -n "$NAMESPACE" patch secret "$secret" \
    --type=merge -p '{"immutable":true}'
done</markup>

<p>Credential and password files are consumed as exact bytes. An unintended trailing newline changes the value. Versioned
Secret names make the selected generation explicit; they do not by themselves prove that a mixed-generation live
rotation is safe.</p>

</div>

<h2 id="_use_managed_security">Use managed security</h2>
<div class="section">
<p>The recommended configuration uses <code>healthMutatorConnection.secure: true</code>. The Operator configures the member&#8217;s
<code>HealthSSLProvider</code>, authenticated Coherence mutators, packaged <code>CoherenceREST</code> JAAS entry, Secret mounts, HTTPS kubelet
probes, and authenticated <code>PUT /suspend</code> client.</p>

<markup
lang="yaml"
title="manifests/managed-packaged-jaas.yaml"
># Copyright (c) 2026, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
# Complete recommended managed configuration. Render NAMESPACE and
# COHERENCE_IMAGE with envsubst before applying. The image must contain the
# application and a Coherence runtime that supports protected PUT mutators.
apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: protected-health
  namespace: ${NAMESPACE}
spec:
  replicas: 3
  image: ${COHERENCE_IMAGE}
  imagePullPolicy: IfNotPresent
  healthPort: 8089
  healthMutatorConnection:
    secure: true
    basicAuth:
      username:
        name: health-auth-v1
        key: username
      password:
        name: health-auth-v1
        key: password
    tls:
      caSecret:
        name: health-ca-v1
        key: ca.pem
      serverName: health.example
    serverTLS:
      enabled: true
      requireClientCert: false
      secrets: health-stores-v1
      keyStore: server.jks
      keyStorePasswordFile: store-password
      keyPasswordFile: store-password
      keyStoreType: JKS
    jaas:
      mode: packaged
  coherence:
    persistence:
      mode: active
      persistentVolumeClaim:
        accessModes:
          - ReadWriteOnce
        resources:
          requests:
            storage: 1Gi</markup>

<p>Render and apply the manifest:</p>

<markup
lang="bash"

>envsubst '${NAMESPACE} ${COHERENCE_IMAGE}' \
  &lt; examples/096_protected_health_endpoints/manifests/managed-packaged-jaas.yaml \
  &gt; "$MATERIAL/protected-health.yaml"

kubectl --context "$CONTEXT" apply -f "$MATERIAL/protected-health.yaml"
kubectl --context "$CONTEXT" -n "$NAMESPACE" rollout status \
  statefulset/protected-health --timeout=10m</markup>

<p>The example enables active persistence and therefore requires a default <code>ReadWriteOnce</code> StorageClass. Replace the image,
application, cache configuration, persistence settings, and resource sizing with the deployment&#8217;s actual values.</p>


<h3 id="_use_an_existing_jaas_file">Use an existing JAAS file</h3>
<div class="section">
<p>When the application already sets <code>java.security.auth.login.config</code>, use <code>jaas.mode: existing</code>. The existing file must
contain a login context named exactly <code>CoherenceREST</code> and must validate the same credentials used by the Operator client.</p>

<markup
lang="yaml"
title="manifests/managed-existing-jaas.yaml"
># Copyright (c) 2026, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
# Managed Health security with an application-owned JAAS file. The file must
# retain the exact CoherenceREST context.
apiVersion: v1
kind: ConfigMap
metadata:
  name: health-jaas-v1
  namespace: ${NAMESPACE}
data:
  jaas.conf: |
    CoherenceREST {
      com.oracle.coherence.k8s.SecretLoginModule required
        usernameFile="/coherence/health-auth/username/value"
        passwordFile="/coherence/health-auth/password/value";
    };
---
# Merge this spec fragment into a StatefulSet-backed Coherence resource.
spec:
  healthPort: 8089
  jvm:
    args:
      - -Djava.security.auth.login.config=/application/security/jaas.conf
  volumes:
    - name: health-jaas
      configMap:
        name: health-jaas-v1
  volumeMounts:
    - name: health-jaas
      mountPath: /application/security
      readOnly: true
  healthMutatorConnection:
    secure: true
    basicAuth:
      username: {name: health-auth-v1, key: username}
      password: {name: health-auth-v1, key: password}
    tls:
      caSecret: {name: health-ca-v1, key: ca.pem}
      serverName: health.example
    serverTLS:
      enabled: true
      requireClientCert: false
      secrets: health-stores-v1
      keyStore: server.jks
      keyStorePasswordFile: store-password
      keyPasswordFile: store-password
      keyStoreType: JKS
    jaas:
      mode: existing</markup>

<p>The file is a <code>ConfigMap</code> plus a <code>spec</code> fragment. Merge that fragment into the StatefulSet-backed <code>Coherence</code> resource.
Managed startup fails instead of overwriting an existing global JAAS selection.</p>

</div>

<h3 id="_use_application_mounted_keystore_files">Use application-mounted keystore files</h3>
<div class="section">
<p>When the application already mounts the server keystore and password files, omit <code>serverTLS.secrets</code> and supply absolute
container paths:</p>

<markup
lang="yaml"
title="manifests/managed-mounted-store-paths.yaml"
># Copyright (c) 2026, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
# Use explicit paths when the application owns the TLS volume mount. The
# credentials and client CA are still namespace-local Secret references.
spec:
  healthPort: 8089
  volumes:
    - name: application-health-tls
      secret:
        secretName: health-stores-v1
  volumeMounts:
    - name: application-health-tls
      mountPath: /application/health-tls
      readOnly: true
  healthMutatorConnection:
    secure: true
    basicAuth:
      username: {name: health-auth-v1, key: username}
      password: {name: health-auth-v1, key: password}
    tls:
      caSecret: {name: health-ca-v1, key: ca.pem}
      serverName: health.example
    serverTLS:
      enabled: true
      requireClientCert: false
      # No 'secrets' field: every value below is an absolute container path.
      keyStore: /application/health-tls/server.jks
      keyStorePasswordFile: /application/health-tls/store-password
      keyPasswordFile: /application/health-tls/store-password
      keyStoreType: JKS
    jaas:
      mode: packaged</markup>

<p>This choice changes only ownership of the member TLS mount. The Operator still configures the Health server, Basic
credentials, client CA, JAAS mode, and authenticated suspension request.</p>

</div>
</div>

<h2 id="_use_application_owned_health_security">Use application-owned Health security</h2>
<div class="section">
<p>Use application-owned mode when the application must configure the Health listener itself or cannot use the managed
launcher arrangement. Omit <code>healthMutatorConnection</code>; the application owns:</p>

<ul class="ulist">
<li>
<p><code>coherence.health.http.mutators.enabled=true</code>;</p>

</li>
<li>
<p>the <code>basic</code>, <code>cert</code>, or <code>cert+basic</code> authentication selection;</p>

</li>
<li>
<p>the HTTPS socket provider and server identity;</p>

</li>
<li>
<p>the <code>CoherenceREST</code> JAAS context when Basic authentication is used;</p>

</li>
<li>
<p>the required Secret volumes and JVM properties; and</p>

</li>
<li>
<p>HTTPS readiness, liveness, and any startup probes.</p>

</li>
</ul>
<p>The following complete Basic-over-TLS example uses the Operator support library&#8217;s <code>SecretLoginModule</code>. The selected
application launcher must place both Coherence and that LoginModule on the member JVM classpath.</p>

<markup
lang="yaml"
title="manifests/application-owned-basic.yaml"
># Copyright (c) 2026, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
# Complete application-owned Basic-over-TLS configuration. Render NAMESPACE
# and COHERENCE_IMAGE with envsubst before applying. The application image must
# place Coherence and the selected JAAS LoginModule on the member classpath.
apiVersion: v1
kind: ConfigMap
metadata:
  name: application-health-jaas-v1
  namespace: ${NAMESPACE}
data:
  jaas.conf: |
    CoherenceREST {
      com.oracle.coherence.k8s.SecretLoginModule required
        usernameFile="/application/health-auth/username"
        passwordFile="/application/health-auth/password";
    };
---
apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: application-owned-health
  namespace: ${NAMESPACE}
spec:
  replicas: 3
  image: ${COHERENCE_IMAGE}
  imagePullPolicy: IfNotPresent
  healthPort: 8089

  # The application owns every Coherence Health security property. There is
  # deliberately no healthMutatorConnection block.
  jvm:
    args:
      - -Dcoherence.health.http.mutators.enabled=true
      - -Dcoherence.health.http.auth=basic
      - -Dcoherence.health.http.provider=HealthSSLProvider
      - -Dcoherence.health.security.keystore=file:/application/health-tls/server.jks
      - -Dcoherence.health.security.keystore.password=/application/health-tls/store-password
      - -Dcoherence.health.security.key.password=/application/health-tls/store-password
      - -Dcoherence.health.security.keystore.type=JKS
      - -Djava.security.auth.login.config=/application/health-jaas/jaas.conf

  volumes:
    - name: application-health-auth
      secret:
        secretName: health-auth-v1
    - name: application-health-tls
      secret:
        secretName: health-stores-v1
    - name: application-health-jaas
      configMap:
        name: application-health-jaas-v1
  volumeMounts:
    - name: application-health-auth
      mountPath: /application/health-auth
      readOnly: true
    - name: application-health-tls
      mountPath: /application/health-tls
      readOnly: true
    - name: application-health-jaas
      mountPath: /application/health-jaas
      readOnly: true

  # This block configures only the Operator client. The JVM arguments, JAAS
  # file, and volumes above configure the member server.
  suspendProbe:
    timeoutSeconds: 60
    http:
      method: PUT
      scheme: HTTPS
      path: /suspend
      port: 8089
      basicAuth:
        username: {name: health-auth-v1, key: username}
        password: {name: health-auth-v1, key: password}
      tls:
        caSecret: {name: health-ca-v1, key: ca.pem}
        serverName: health.example

  # Application-owned mode does not rewrite these probes to HTTPS.
  readinessProbe:
    httpGet:
      scheme: HTTPS
      path: /ready
      port: 8089
  livenessProbe:
    httpGet:
      scheme: HTTPS
      path: /live
      port: 8089
  startupProbe:
    httpGet:
      scheme: HTTPS
      path: /started
      port: 8089</markup>

<p>Render and apply it in the same way as the managed example:</p>

<markup
lang="bash"

>envsubst '${NAMESPACE} ${COHERENCE_IMAGE}' \
  &lt; examples/096_protected_health_endpoints/manifests/application-owned-basic.yaml \
  &gt; "$MATERIAL/application-owned-health.yaml"

kubectl --context "$CONTEXT" apply -f "$MATERIAL/application-owned-health.yaml"
kubectl --context "$CONTEXT" -n "$NAMESPACE" rollout status \
  statefulset/application-owned-health --timeout=10m</markup>

<p>The <code>suspendProbe.http</code> block configures only the Operator client. Unlike managed mode, it does not configure the member
JVM or rewrite kubelet probes to HTTPS.</p>


<h3 id="_use_only_the_custom_http_action_fragment">Use only the custom HTTP-action fragment</h3>
<div class="section">
<p>If the application already supplies the protected server configuration, merge the following Operator-client and probe
fragment into its <code>Coherence</code> resource:</p>

<markup
lang="yaml"
title="manifests/application-owned-http-action.yaml"
># Copyright (c) 2026, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
# Operator-client fragment for an application-owned protected Health server.
# Use only when the application independently configures the protected
# Coherence health server. Do not add healthMutatorConnection with this probe.
spec:
  healthPort: 8089
  suspendProbe:
    timeoutSeconds: 60
    http:
      method: PUT
      scheme: HTTPS
      path: /suspend
      port: 8089
      basicAuth:
        username: {name: health-auth-v1, key: username}
        password: {name: health-auth-v1, key: password}
      tls:
        caSecret: {name: health-ca-v1, key: ca.pem}
        serverName: health.example

  # Custom probes are required because unmanaged mode does not rewrite the
  # generated probes to HTTPS. Read-only endpoints must remain unauthenticated.
  readinessProbe:
    httpGet:
      scheme: HTTPS
      path: /ready
      port: 8089
  livenessProbe:
    httpGet:
      scheme: HTTPS
      path: /live
      port: 8089
  startupProbe:
    httpGet:
      scheme: HTTPS
      path: /started
      port: 8089</markup>

<p>The HTTP action supports verified server TLS and Secret-backed Basic credentials, but it does not present a TLS client
certificate.</p>

</div>

<h3 id="_use_a_custom_in_pod_client">Use a custom in-Pod client</h3>
<div class="section">
<p>For certificate authentication, <code>cert+basic</code>, or another client flow that the HTTP action cannot express, use an
application-owned exec action. The application image must contain the required client program and must mount its client
material.</p>

<markup
lang="yaml"
title="manifests/application-owned-exec-action.yaml"
># Copyright (c) 2026, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
# Last-resort application-owned action. The image must contain /bin/sh and curl.
# The application independently configures its protected health server.
spec:
  healthPort: 8089
  volumes:
    - name: health-client-material
      projected:
        sources:
          - secret:
              name: health-auth-v1
              items:
                - {key: username, path: username}
                - {key: password, path: password}
          - secret:
              name: health-ca-v1
              items:
                - {key: ca.pem, path: ca.pem}
  volumeMounts:
    - name: health-client-material
      mountPath: /application/health-client
      readOnly: true
  suspendProbe:
    timeoutSeconds: 60
    exec:
      command:
        - /bin/sh
        - -ec
        - &gt;-
          curl --silent --show-error --fail
          --cacert /application/health-client/ca.pem
          --resolve health.example:8089:127.0.0.1
          --user "$(cat /application/health-client/username):$(cat /application/health-client/password)"
          --request PUT https://health.example:8089/suspend
  readinessProbe:
    httpGet:
      scheme: HTTPS
      path: /ready
      port: 8089
  livenessProbe:
    httpGet:
      scheme: HTTPS
      path: /live
      port: 8089
  startupProbe:
    httpGet:
      scheme: HTTPS
      path: /started
      port: 8089</markup>

<p>The supplied fragment demonstrates the resource shape with <code>curl</code> and Basic authentication. Replace the command and
mounted material for the selected certificate contract. An exec action is less portable and provides fewer sanitized
diagnostics than the managed HTTP client.</p>

</div>
</div>

<h2 id="_verify_the_deployment">Verify the deployment</h2>
<div class="section">
<p>For managed mode, confirm that generated probes targeting the Health port use HTTPS:</p>

<markup
lang="bash"

>kubectl --context "$CONTEXT" -n "$NAMESPACE" get statefulset protected-health \
  -o jsonpath='{range .spec.template.spec.containers[*]}{.name}{" ready="}{.readinessProbe.httpGet.scheme}{":"}{.readinessProbe.httpGet.port}{" live="}{.livenessProbe.httpGet.scheme}{":"}{.livenessProbe.httpGet.port}{"\n"}{end}'</markup>

<p>Then verify the deployment&#8217;s application data before and after scale-to-zero. A failed credential, CA, or certificate
name must block the transition while the existing members remain running. Ready Pods and unchanged PVC names do not by
themselves prove that application data was recovered.</p>

<p>Do not add an anonymous or <code>GET</code> fallback, disable certificate verification, or set
<code>spec.suspendServicesOnShutdown: false</code> to hide a security configuration failure.</p>

<div class="admonition warning">
<p class="admonition-inline">If scale-to-zero is blocked, correct the configuration or restore <code>spec.replicas</code> to a non-zero value before
deleting the <code>Coherence</code> resource or its namespace. Desired replicas can be zero while members are still running, and
deletion in that state can terminate those members without the expected suspension.</p>
</div>
</div>

<h2 id="_clean_up_an_example_namespace">Clean up an example namespace</h2>
<div class="section">
<p>Delete only a namespace that was created specifically for this example:</p>

<markup
lang="bash"

>kubectl --context "$CONTEXT" delete namespace "$NAMESPACE"</markup>

<p>Namespace deletion can also delete dynamically provisioned volumes according to the StorageClass reclaim policy.</p>

</div>
</doc-view>
