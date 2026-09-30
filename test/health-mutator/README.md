<!--
Copyright (c) 2026, Oracle and/or its affiliates.
Licensed under the Universal Permissive License v 1.0 as shown at
https://oss.oracle.com/licenses/upl.
-->

# Protected Health mutator acceptance fixtures

This directory contains repository-only qualification fixtures for protected
Coherence Health mutators. The procedures use test applications, generated
credentials, disposable namespaces, failure injection, and destructive
cleanup. They are not production deployment instructions.

Customer documentation and reusable manifests are published in the
[Protected Coherence Health Endpoints example](../../examples/096_protected_health_endpoints/README.adoc).

## Qualification scope and provenance

Qualify the exact Operator image digest, Coherence artifact, Java version,
application image digest, and Kubernetes version used for a deployment.

Complete runtime qualification used Coherence CE 15.1.1-0-5 with normal
Java/classpath packaging and proved `3 -> blocked zero -> 0 -> 3` plus
persistent recovery. Coherence CE 22.06.18 passed startup and
HTTPS/authentication smoke checks only. Base CE 22.06 fails closed at the
secure-mutator capability boundary. Multi-worker Pod placement remains
unqualified.

## Managed automated fixture

The automated fixture validates the recommended
`spec.healthMutatorConnection.secure: true` path. Use a separate test namespace
and a test application image built with the affected hardened Coherence
artifact. The image must:

- contain `com.oracle.coherence.k8s.testing.RestServer` and
  `test-cache-config.xml`;
- be available to every node; and
- be selected by immutable digest.

The scripts generate 30-day test credentials and TLS material, create immutable
Secrets, deploy two storage-enabled members with active persistence, and prove:

- the unconfigured hardened mutator returns HTTP 403;
- read-only Health remains available over verified HTTPS;
- missing and invalid Basic credentials return HTTP 401;
- valid credentials can suspend and resume services;
- invalid Operator CA trust blocks scale-to-zero;
- repairing trust allows scale-to-zero; and
- the persistent canary is recovered after restoring two members.

From the repository root:

```bash
test/health-mutator/apply.sh \
  kind-coh40013873-qual \
  coh40013873-qual \
  registry.example/operator-test@sha256:REPLACE_WITH_VERIFIED_DIGEST \
  /tmp/coh40013873-health-v1

test/health-mutator/verify.sh \
  kind-coh40013873-qual \
  coh40013873-qual \
  /tmp/coh40013873-health-v1
```

The arguments are the explicit Kubernetes context, new disposable namespace,
pinned application-image digest, and new material directory. Preserve the
rendered manifests, Pod placement, events, exact image identities, Kubernetes
version, PVC identities, and canary results before deleting the namespace.

## Application-owned manual workflow

Use this workflow when the application, rather than
`spec.healthMutatorConnection`, owns the Coherence Health server's TLS, JAAS,
and mutator properties. The Operator owns only the request it uses during safe
shutdown or scale-to-zero.

The supplied qualification fixture is `application-owned.yaml`. It shows the exact
application-owned configuration qualified with Coherence 15.1.1-0-5. The
application main class, cache configuration, and persistence-canary endpoint
are qualification-image details; replace them with the corresponding parts of
the application being tested.

### Ownership boundary

The `Coherence` resource deliberately omits `healthMutatorConnection`.

The application owns:

- enabling the Coherence Health mutators;
- selecting the Health server authentication mode;
- selecting and configuring the TLS socket provider;
- mounting the server keystore and password files;
- providing the `CoherenceREST` JAAS context; and
- declaring HTTPS kubelet probes.

The Operator owns only `spec.suspendProbe`. For the Basic-over-TLS workflow it
reads the client username, password, and CA from namespace-local Secrets and
sends a verified HTTPS `PUT /suspend`.

### Preconditions

- Use a disposable namespace and a StatefulSet-backed `Coherence` resource.
- The application must use Coherence 15.1.1-0-5 or another version that exposes
  the hardened PUT mutator contract.
- The selected application launcher must place the Coherence and Operator
  support classes on the member JVM classpath.
- Do not select the legacy Operator health server with
  `COHERENCE_OPERATOR_HEALTH_CHECK=true`.
- Use a digest-pinned application image.
- The cluster needs a default `ReadWriteOnce` StorageClass for this demo.

### 1. Generate test credentials and TLS material

From the root of this Operator repository:

```bash
export NAMESPACE=coh-app-owned-health
export CONTEXT=YOUR_CONTEXT
export MATERIAL=/tmp/coh-app-owned-health-material

test/health-mutator/prepare.sh "$MATERIAL" health.example
```

The generated certificate contains `health.example` as a DNS subject
alternative name. The Operator will connect to a Pod IP but verify this name.

For production, replace this generated material with certificates and random
credentials issued through the normal security process.

### 2. Create immutable, versioned Secrets

```bash
kubectl --context "$CONTEXT" create namespace "$NAMESPACE"

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
done
```

The credential files contain exact bytes without trailing newlines. The JAAS
login module performs an exact comparison, so an accidental newline changes
the credential.

### 3. Understand the application-owned server properties

The full example supplies these member JVM properties:

```text
coherence.health.http.mutators.enabled=true
coherence.health.http.auth=basic
coherence.health.http.provider=HealthSSLProvider
coherence.health.security.keystore=file:/application/health-tls/server.jks
coherence.health.security.keystore.password=/application/health-tls/store-password
coherence.health.security.key.password=/application/health-tls/store-password
coherence.health.security.keystore.type=JKS
java.security.auth.login.config=/application/health-jaas/jaas.conf
```

Their responsibilities are:

- `mutators.enabled=true` permits PUT `/suspend` and `/resume`.
- `auth=basic` requires an HTTP Basic identity for mutator routes.
- `provider=HealthSSLProvider` selects the Operator-supplied socket-provider
  definition used by the Coherence Health listener.
- The `coherence.health.security.*` properties identify the application-mounted
  server keystore and password files consumed by that provider. One-way TLS
  does not use a member-side truststore.
- `java.security.auth.login.config` selects the application-owned JAAS file.

These direct properties are valid only because `healthMutatorConnection` is
absent. They conflict with managed secure mode and are rejected when both
approaches are configured.

The `CoherenceREST` JAAS entry uses the Operator support library's
`SecretLoginModule`, but its file paths are application-owned mounts:

```text
CoherenceREST {
  com.oracle.coherence.k8s.SecretLoginModule required
    usernameFile="/application/health-auth/username"
    passwordFile="/application/health-auth/password";
};
```

An application can substitute another compatible JAAS LoginModule, but it must
retain the exact `CoherenceREST` context and accept Coherence's
`UsernameAndPassword` identity assertion.

### 4. Render and apply the complete resource

Build the application image with the Coherence version being qualified. For
Kind, load it into every Kubernetes node. For a shared cluster, push it to a
registry reachable by every node. In either case, obtain and use its immutable
manifest digest, not a mutable tag. The image used by this exact example must
contain `com.oracle.coherence.k8s.testing.HealthMutatorDemo` and
`test-cache-config.xml`; when qualifying another application, update those two
fields in the resource.

```bash
export IMAGE_DIGEST='registry.example/demo@sha256:REPLACE_WITH_DIGEST'

envsubst '${NAMESPACE} ${IMAGE_DIGEST}' \
  < test/health-mutator/application-owned.yaml \
  > "$MATERIAL/application-owned-health.yaml"

kubectl --context "$CONTEXT" apply \
  -f "$MATERIAL/application-owned-health.yaml"

until kubectl --context "$CONTEXT" -n "$NAMESPACE" get \
  statefulset/application-owned-health >/dev/null 2>&1; do sleep 2; done

kubectl --context "$CONTEXT" -n "$NAMESPACE" rollout status \
  statefulset/application-owned-health --timeout=10m
```

The example starts three storage-enabled members with active persistence.

### 5. Verify the server before testing Operator suspension

Confirm that the running JVM contains the application-owned settings and that
the kubelet probes use HTTPS:

```bash
kubectl --context "$CONTEXT" -n "$NAMESPACE" get statefulset application-owned-health \
  -o jsonpath='{range .spec.template.spec.containers[*]}{.name}{" ready="}{.readinessProbe.httpGet.scheme}{":"}{.readinessProbe.httpGet.port}{" live="}{.livenessProbe.httpGet.scheme}{":"}{.livenessProbe.httpGet.port}{"\n"}{end}'
```

Expected Coherence-container output:

```text
coherence ready=HTTPS:8089 live=HTTPS:8089
```

Port-forward one member:

```bash
kubectl --context "$CONTEXT" -n "$NAMESPACE" port-forward \
  pod/application-owned-health-0 18089:8089 18080:8080
```

In another terminal, verify the protocol. Read-only health stays
unauthenticated; mutations require valid Basic credentials:

```bash
curl --fail --cacert "$MATERIAL/ca.pem" \
  --resolve health.example:18089:127.0.0.1 \
  https://health.example:18089/ready

curl --silent --output /dev/null --write-out '%{http_code}\n' \
  --request PUT --cacert "$MATERIAL/ca.pem" \
  --resolve health.example:18089:127.0.0.1 \
  https://health.example:18089/suspend
# Expected: 401

curl --fail --request PUT \
  --user "$(<"$MATERIAL/username"):$(<"$MATERIAL/password")" \
  --cacert "$MATERIAL/ca.pem" \
  --resolve health.example:18089:127.0.0.1 \
  https://health.example:18089/suspend

curl --fail --request PUT \
  --user "$(<"$MATERIAL/username"):$(<"$MATERIAL/password")" \
  --cacert "$MATERIAL/ca.pem" \
  --resolve health.example:18089:127.0.0.1 \
  https://health.example:18089/resume
```

Always resume after a direct suspension check before continuing.

### 6. Establish a persistent canary

The qualification application exposes `/canaryStart` and `/canaryCheck` on
port 8080. With the same port-forward running:

```bash
curl --fail http://127.0.0.1:18080/canaryStart
# Expected with the default partition count: stored 257 entries

kubectl --context "$CONTEXT" -n "$NAMESPACE" get pvc -o name | sort \
  > "$MATERIAL/pvcs-before-zero.txt"
```

For another application, replace these two endpoints with an application-level
write and read that prove data is present in every required persistent cache.
Record enough information to distinguish recovered data from newly created
data after restart.

### 7. Prove that bad Operator trust fails closed

Create a bad CA Secret and point only the custom Operator suspension action to
it. This does not change the server certificate or kubelet probes.

```bash
printf 'not a certificate\n' > "$MATERIAL/bad-ca.pem"

kubectl --context "$CONTEXT" -n "$NAMESPACE" create secret generic health-bad-ca-v1 \
  --from-file=ca.pem="$MATERIAL/bad-ca.pem"

kubectl --context "$CONTEXT" -n "$NAMESPACE" patch \
  coherence.coherence.oracle.com application-owned-health --type=merge \
  -p '{"spec":{"suspendProbe":{"timeoutSeconds":60,"http":{"method":"PUT","scheme":"HTTPS","path":"/suspend","port":8089,"basicAuth":{"username":{"name":"health-auth-v1","key":"username"},"password":{"name":"health-auth-v1","key":"password"}},"tls":{"caSecret":{"name":"health-bad-ca-v1","key":"ca.pem"},"serverName":"health.example"}}}}}'

kubectl --context "$CONTEXT" -n "$NAMESPACE" patch \
  coherence.coherence.oracle.com application-owned-health --type=merge \
  -p '{"spec":{"replicas":0}}'
```

The Operator must report a suspension failure and keep the StatefulSet and all
three Pods:

```bash
kubectl --context "$CONTEXT" -n "$NAMESPACE" get events \
  --field-selector involvedObject.kind=Coherence,involvedObject.name=application-owned-health,type=Warning

kubectl --context "$CONTEXT" -n "$NAMESPACE" get statefulset,pods
```

Do not proceed unless the transition was blocked.

### 8. Repair trust and reach zero

Restore the valid CA reference. The unchanged `replicas: 0` request should be
retried and complete:

```bash
kubectl --context "$CONTEXT" -n "$NAMESPACE" patch \
  coherence.coherence.oracle.com application-owned-health --type=merge \
  -p '{"spec":{"suspendProbe":{"timeoutSeconds":60,"http":{"method":"PUT","scheme":"HTTPS","path":"/suspend","port":8089,"basicAuth":{"username":{"name":"health-auth-v1","key":"username"},"password":{"name":"health-auth-v1","key":"password"}},"tls":{"caSecret":{"name":"health-ca-v1","key":"ca.pem"},"serverName":"health.example"}}}}}'

kubectl --context "$CONTEXT" -n "$NAMESPACE" wait --for=delete \
  statefulset/application-owned-health --timeout=10m

for _ in $(seq 1 120); do
  count=$(kubectl --context "$CONTEXT" -n "$NAMESPACE" get pods \
    -l coherenceDeployment=application-owned-health -o name | wc -l | tr -d ' ')
  [[ "$count" = 0 ]] && break
  sleep 2
done

test "$(kubectl --context "$CONTEXT" -n "$NAMESPACE" get pods \
  -l coherenceDeployment=application-owned-health -o name | wc -l | tr -d ' ')" = 0

kubectl --context "$CONTEXT" -n "$NAMESPACE" get pvc -o name | sort \
  > "$MATERIAL/pvcs-at-zero.txt"
diff -u "$MATERIAL/pvcs-before-zero.txt" "$MATERIAL/pvcs-at-zero.txt"
```

### 9. Restore three members and verify persistent recovery

```bash
kubectl --context "$CONTEXT" -n "$NAMESPACE" patch \
  coherence.coherence.oracle.com application-owned-health --type=merge \
  -p '{"spec":{"replicas":3}}'

until kubectl --context "$CONTEXT" -n "$NAMESPACE" get \
  statefulset/application-owned-health >/dev/null 2>&1; do sleep 2; done

kubectl --context "$CONTEXT" -n "$NAMESPACE" rollout status \
  statefulset/application-owned-health --timeout=10m

kubectl --context "$CONTEXT" -n "$NAMESPACE" port-forward \
  pod/application-owned-health-0 18089:8089 18080:8080
```

In another terminal:

```bash
curl --fail http://127.0.0.1:18080/canaryCheck
# Expected: recovered 257 entries
```

Record the final PVC names and require an empty diff:

```bash
kubectl --context "$CONTEXT" -n "$NAMESPACE" get pvc -o name | sort \
  > "$MATERIAL/pvcs-after-restart.txt"
diff -u "$MATERIAL/pvcs-before-zero.txt" "$MATERIAL/pvcs-after-restart.txt"
```

### 10. Preserve evidence and clean up

Retain the rendered resource, Pod placement, warning events, image digest,
Operator image/revision, Kubernetes version, PVC lists, and canary results.

After copying the evidence to durable storage:

```bash
kubectl --context "$CONTEXT" delete namespace "$NAMESPACE"
```

Namespace deletion can delete dynamically provisioned volumes according to the
StorageClass reclaim policy.

### Validated result

On 2026-09-28 this exact application-owned Basic-over-TLS manifest passed in
the disposable Kubernetes 1.34.0/linux-arm64 Kind environment with Operator
commit `fd9192e8bd0bdcd82b35df2eb4f9eecbc837e80f` and application image:

```text
docker.io/library/coherence-health-mutator-demo@sha256:7e17b851a0ed02885bb23915f108e023d45b9794eb79cf4c0d5b3f0750edd5ec
```

The run proved:

- generated application-owned kubelet probes used HTTPS on port 8089;
- an unauthenticated PUT returned 401;
- authenticated suspend and resume returned 200;
- invalid Operator CA trust blocked scale-to-zero with three ready members;
- restoring the CA completed the transition from three members to zero Pods;
- all three PVC identities were retained; and
- restoring three members recovered all 257 persistent canary entries.

The Kind environment had one Kubernetes worker/control-plane node. This proves
three Coherence members, not scheduling across three physical worker nodes.

### Other application-owned authentication modes

Coherence 15.1.1-0-5 also exposes `cert` and `cert+basic` mutator authentication
modes. They require an application-owned socket provider that requests and
validates client certificates.

The Operator's custom HTTP action supports server verification and
Secret-backed Basic credentials, but it has no client-certificate fields.
Therefore it cannot natively call `cert` or `cert+basic` endpoints. Use a
custom `suspendProbe.exec` command capable of presenting the client certificate
and key, or use managed Basic-over-one-way-TLS. Treat the custom exec path as
application-owned and qualify it separately; it is not covered by the demo's
Basic-over-TLS acceptance result.
