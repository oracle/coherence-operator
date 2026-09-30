<!--
Copyright (c) 2026, Oracle and/or its affiliates.
Licensed under the Universal Permissive License v 1.0 as shown at
https://oss.oracle.com/licenses/upl.
-->

# Authenticated Coherence health mutators

The generated client documentation publishes the complete first-run guide at
[`010_protected_health_endpoints.adoc`](010_protected_health_endpoints.adoc).
Reusable customer manifests and procedures are in the
[`096_protected_health_endpoints` example](../../examples/096_protected_health_endpoints/README.adoc).
The notes below summarize implementation boundaries; they are not a substitute
for that client workflow.

`spec.healthMutatorConnection` is an opt-in StatefulSet setting for Coherence runtimes that expose the hardened Coherence health handler. It does not switch health-server implementations. If `COHERENCE_OPERATOR_HEALTH_CHECK=true` selects the legacy Operator health server, managed secure mode fails during Pod configuration.

The supported managed path sends `PUT /suspend` with Secret-backed Basic credentials over verified HTTPS. It enables `coherence.health.http.mutators.enabled=true`, selects `coherence.health.http.auth=basic`, configures `HealthSSLProvider`, and uses the `CoherenceREST` JAAS context. Read-only health probes stay unauthenticated. Explicit custom suspend and kubelet probes remain customer-owned and are not rewritten.

One-way member TLS requires a server keystore and the keystore/key password files. The member-side provider does not load a truststore because it does not request client certificates. The separate `tls.caSecret` remains required for the Operator to verify each member's certificate. Older `serverTLS.trustStore*` fields are still accepted but are ignored by this health provider.

The initial supported packaging paths are normal Java/classpath and Spring Boot
2/3 with a configured fat JAR. The Cloud Native Buildpack launcher is rejected
explicitly. Startup fails closed if the selected application runtime cannot load
the hardened handler contract. Qualify the exact released Operator image,
Coherence artifact, Java version, application image, and Kubernetes version used
for a deployment.

## Protected endpoint configuration examples

See the [published customer example](../../examples/096_protected_health_endpoints/README.adoc)
for supported managed and application-owned configurations.

## Repository acceptance fixtures

The repository-only managed fixture and manual application-owned qualification
workflow now live beside their scripts and manifests in
[`test/health-mutator/README.md`](../../test/health-mutator/README.md). They use
test applications, generated credentials, disposable namespaces, failure
injection, and destructive cleanup; they are not production deployment
instructions.

Do not add GET or anonymous fallback, disable TLS verification, select the
Operator health server as a workaround, or remove old credential and certificate
material while an old Pod generation may still use it. Mixed-generation
credential or certificate rotation and no-downtime rotation remain unqualified.
