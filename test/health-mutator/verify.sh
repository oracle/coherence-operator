#!/usr/bin/env bash
# Copyright (c) 2026, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
set -euo pipefail

context=${1:?Kubernetes context is required}
namespace=${2:?namespace is required}
material=${3:?material directory is required}
health_port=${HEALTH_MUTATOR_LOCAL_PORT:-18089}
rest_port=${HEALTH_MUTATOR_REST_PORT:-18080}
port_forward_pid=

stop_forward() {
  if [[ -n "$port_forward_pid" ]]; then
    kill "$port_forward_pid" >/dev/null 2>&1 || true
    wait "$port_forward_pid" >/dev/null 2>&1 || true
    port_forward_pid=
  fi
}
trap stop_forward EXIT

start_forward() {
  scheme=$1
  stop_forward
  pod=$(kubectl --context "$context" -n "$namespace" get pod \
    -l coherenceDeployment=health-secure \
    -o jsonpath='{.items[0].metadata.name}')
  kubectl --context "$context" -n "$namespace" port-forward "pod/$pod" \
    "$health_port:8089" "$rest_port:8080" > "$material/port-forward.log" 2>&1 &
  port_forward_pid=$!
  for _ in $(seq 1 60); do
    if [[ "$scheme" == https ]]; then
      if curl --noproxy '*' --silent --show-error --fail \
        --cacert "$material/ca.pem" \
        --resolve "health.example:$health_port:127.0.0.1" \
        "https://health.example:$health_port/ready" >/dev/null 2>&1; then
        return
      fi
    elif curl --noproxy '*' --silent --show-error --fail \
      "http://127.0.0.1:$health_port/ready" >/dev/null 2>&1; then
      return
    fi
    sleep 1
  done
  echo "Health endpoint did not become reachable; see $material/port-forward.log" >&2
  exit 1
}

status() {
  curl --noproxy '*' --silent --show-error --output /dev/null --write-out '%{http_code}' \
    --cacert "$material/ca.pem" \
    --resolve "health.example:$health_port:127.0.0.1" "$@"
}

expect_status() {
  expected=$1
  shift
  actual=$(status "$@")
  if [[ "$actual" != "$expected" ]]; then
    echo "Expected HTTP $expected but received $actual" >&2
    exit 1
  fi
}

start_forward http
unconfigured_status=$(curl --noproxy '*' --silent --show-error --output /dev/null --write-out '%{http_code}' \
  --request PUT "http://127.0.0.1:$health_port/suspend")
if [[ "$unconfigured_status" != 403 ]]; then
  echo "Expected the unconfigured hardened mutator to return HTTP 403 but received $unconfigured_status" >&2
  exit 1
fi
stop_forward
kubectl --context "$context" apply -f "$material/coherence.yaml"
for _ in $(seq 1 120); do
  annotation=$(kubectl --context "$context" -n "$namespace" get statefulset health-secure \
    -o jsonpath='{.spec.template.metadata.annotations.coherence\.oracle\.com/health-mutator-connection}' 2>/dev/null || true)
  if [[ "$annotation" == *health-auth-v1* ]]; then
    break
  fi
  sleep 2
done
if [[ "$annotation" != *health-auth-v1* ]]; then
  echo "StatefulSet did not adopt the managed health configuration" >&2
  exit 1
fi
kubectl --context "$context" -n "$namespace" rollout status statefulset/health-secure --timeout=10m

start_forward https
expect_status 200 "https://health.example:$health_port/ready"
expect_status 401 --request PUT "https://health.example:$health_port/suspend"
expect_status 401 --request PUT --user 'health-operator:wrong' "https://health.example:$health_port/suspend"
curl --noproxy '*' --silent --show-error --fail "http://127.0.0.1:$rest_port/canaryStart" >/dev/null
expect_status 200 --request PUT \
  --user "$(<"$material/username"):$(<"$material/password")" \
  "https://health.example:$health_port/suspend"
expect_status 200 --request PUT \
  --user "$(<"$material/username"):$(<"$material/password")" \
  "https://health.example:$health_port/resume"
if curl --noproxy '*' --silent --show-error --fail \
  --cacert "$material/password" \
  --resolve "health.example:$health_port:127.0.0.1" \
  "https://health.example:$health_port/ready" >/dev/null 2>&1; then
  echo "Invalid CA unexpectedly passed TLS verification" >&2
  exit 1
fi
stop_forward

kubectl --context "$context" -n "$namespace" create secret generic health-ca-failure-v1 \
  --from-file=ca.pem="$material/ca.pem" --dry-run=client -o yaml | \
  kubectl --context "$context" apply -f -
kubectl --context "$context" -n "$namespace" patch coherence.coherence.oracle.com health-secure --type=merge \
  -p '{"spec":{"healthMutatorConnection":{"tls":{"caSecret":{"name":"health-ca-failure-v1","key":"ca.pem"},"serverName":"health.example"}}}}'
for _ in $(seq 1 120); do
  annotation=$(kubectl --context "$context" -n "$namespace" get statefulset health-secure \
    -o jsonpath='{.spec.template.metadata.annotations.coherence\.oracle\.com/health-mutator-connection}' 2>/dev/null || true)
  if [[ "$annotation" == *health-ca-failure-v1* ]]; then
    break
  fi
  sleep 2
done
if [[ "$annotation" != *health-ca-failure-v1* ]]; then
  echo "StatefulSet did not adopt the failure-injection CA reference" >&2
  exit 1
fi
kubectl --context "$context" -n "$namespace" rollout status statefulset/health-secure --timeout=10m
printf 'not a certificate\n' > "$material/bad-ca.pem"
kubectl --context "$context" -n "$namespace" create secret generic health-ca-failure-v1 \
  --from-file=ca.pem="$material/bad-ca.pem" --dry-run=client -o yaml | \
  kubectl --context "$context" apply -f -
kubectl --context "$context" -n "$namespace" patch coherence.coherence.oracle.com health-secure --type=merge -p '{"spec":{"replicas":0}}'
for _ in $(seq 1 120); do
  warning_events=$(kubectl --context "$context" -n "$namespace" get events \
    --field-selector involvedObject.kind=Coherence,involvedObject.name=health-secure,type=Warning \
    -o jsonpath='{range .items[?(@.reason=="Scaling")]}{.message}{"\n"}{end}' 2>/dev/null || true)
  if [[ "$warning_events" == *"failed suspending Coherence services"* ]]; then
    break
  fi
  sleep 2
done
if [[ "$warning_events" != *"failed suspending Coherence services"* ]]; then
  echo "Operator did not report a failed suspension with invalid trust" >&2
  exit 1
fi
if ! kubectl --context "$context" -n "$namespace" get statefulset health-secure >/dev/null 2>&1; then
  echo "Scale-to-zero was not blocked by invalid trust" >&2
  exit 1
fi

kubectl --context "$context" -n "$namespace" create secret generic health-ca-failure-v1 \
  --from-file=ca.pem="$material/ca.pem" --dry-run=client -o yaml | \
  kubectl --context "$context" apply -f -
for _ in $(seq 1 120); do
  if ! kubectl --context "$context" -n "$namespace" get statefulset health-secure >/dev/null 2>&1; then
    break
  fi
  sleep 2
done
if kubectl --context "$context" -n "$namespace" get statefulset health-secure >/dev/null 2>&1; then
  echo "StatefulSet did not reach scale-to-zero after restoring trust" >&2
  exit 1
fi

kubectl --context "$context" -n "$namespace" patch coherence.coherence.oracle.com health-secure --type=merge -p '{"spec":{"replicas":2}}'
for _ in $(seq 1 120); do
  if kubectl --context "$context" -n "$namespace" get statefulset health-secure >/dev/null 2>&1; then
    break
  fi
  sleep 2
done
if ! kubectl --context "$context" -n "$namespace" get statefulset health-secure >/dev/null 2>&1; then
  echo "StatefulSet was not recreated" >&2
  exit 1
fi
kubectl --context "$context" -n "$namespace" rollout status statefulset/health-secure --timeout=10m
start_forward https
curl --noproxy '*' --silent --show-error --fail "http://127.0.0.1:$rest_port/canaryCheck" >/dev/null
echo "Unconfigured 403, authenticated suspension, fail-closed TLS, scale-to-zero, restart, and persistent canary verification passed"
