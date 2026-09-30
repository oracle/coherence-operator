#!/usr/bin/env bash
# Copyright (c) 2026, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
set -euo pipefail

context=${1:?Kubernetes context is required}
namespace=${2:?namespace is required}
image_digest=${3:?pinned test application image digest is required}
material=${4:?new material directory is required}

if [[ "$image_digest" != *@sha256:* ]]; then
  echo "The test application image must be pinned by digest" >&2
  exit 1
fi
if [[ -e "$material" ]]; then
  echo "Material directory already exists: $material" >&2
  exit 1
fi

script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
bash "$script_dir/prepare.sh" "$material" health.example

kubectl --context "$context" create namespace "$namespace" --dry-run=client -o yaml | kubectl --context "$context" apply -f -
kubectl --context "$context" -n "$namespace" create secret generic health-auth-v1 \
  --from-file=username="$material/username" \
  --from-file=password="$material/password" \
  --dry-run=client -o yaml | kubectl --context "$context" apply -f -
kubectl --context "$context" -n "$namespace" create secret generic health-stores-v1 \
  --from-file=server.jks="$material/server.jks" \
  --from-file=store-password="$material/store-password" \
  --dry-run=client -o yaml | kubectl --context "$context" apply -f -
kubectl --context "$context" -n "$namespace" create secret generic health-ca-v1 \
  --from-file=ca.pem="$material/ca.pem" \
  --dry-run=client -o yaml | kubectl --context "$context" apply -f -
for secret in health-auth-v1 health-stores-v1 health-ca-v1; do
  kubectl --context "$context" -n "$namespace" patch secret "$secret" --type=merge -p '{"immutable":true}'
done

export IMAGE_DIGEST=$image_digest
export NAMESPACE=$namespace
envsubst "\${IMAGE_DIGEST} \${NAMESPACE}" < "$script_dir/coherence.yaml" > "$material/coherence.yaml"
envsubst "\${IMAGE_DIGEST} \${NAMESPACE}" < "$script_dir/coherence-unconfigured.yaml" > "$material/coherence-unconfigured.yaml"
kubectl --context "$context" apply -f "$material/coherence-unconfigured.yaml"
for _ in $(seq 1 120); do
  if kubectl --context "$context" -n "$namespace" get statefulset health-secure >/dev/null 2>&1; then
    break
  fi
  sleep 2
done
if ! kubectl --context "$context" -n "$namespace" get statefulset health-secure >/dev/null 2>&1; then
  echo "StatefulSet was not created" >&2
  exit 1
fi
kubectl --context "$context" -n "$namespace" rollout status statefulset/health-secure --timeout=10m
