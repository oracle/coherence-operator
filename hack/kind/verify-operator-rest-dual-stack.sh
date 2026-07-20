#!/usr/bin/env bash

#
# Copyright (c) 2026, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

set -o errexit
set -o nounset
set -o pipefail

namespace="${1:-operator-test}"
service_name="${2:-coherence-operator-rest}"
expected_policy="${3:-RequireDualStack}"
probe_name="coherence-operator-rest-dual-stack-probe"
missing_resource="coherence-operator-rest-verification-does-not-exist"
endpoint_slices=""

if (( $# >= 4 )); then
  expected_families=("${@:4}")
else
  expected_families=(IPv4 IPv6)
fi

fail() {
  echo "$*" >&2
  exit 1
}

ip_family() {
  if [[ "$1" == *:* ]]; then
    echo IPv6
  else
    echo IPv4
  fi
}

assert_dual_stack_pod() {
  local pod_name="$1"
  local pod_ips=()
  local ip
  local has_ipv4=false
  local has_ipv6=false

  mapfile -t pod_ips < <(
    kubectl --namespace "${namespace}" get pod "${pod_name}" \
      -o jsonpath='{range .status.podIPs[*]}{.ip}{"\n"}{end}'
  )
  for ip in "${pod_ips[@]}"; do
    if [[ "$(ip_family "${ip}")" == IPv4 ]]; then
      has_ipv4=true
    else
      has_ipv6=true
    fi
  done
  if [[ "${has_ipv4}" != true || "${has_ipv6}" != true ]]; then
    fail "Expected Pod ${pod_name} to have IPv4 and IPv6 addresses, found ${pod_ips[*]:-<none>}"
  fi
}

probe_http() {
  local description="$1"
  local expected_ip="$2"
  shift 2
  local output=""
  local metadata=""
  local body=""
  local http_code=""
  local remote_ip=""
  local attempt

  for attempt in 1 2 3 4 5; do
    if output="$(
      kubectl --namespace "${namespace}" exec "${probe_name}" -- \
        curl --noproxy '*' --silent --show-error --connect-timeout 5 --max-time 10 \
          --write-out $'\n%{http_code}\t%{remote_ip}' "$@"
    )"; then
      metadata="${output##*$'\n'}"
      body="${output%$'\n'*}"
      http_code="${metadata%%$'\t'*}"
      remote_ip="${metadata#*$'\t'}"
      if [[ "${http_code}" == 404 && "${remote_ip}" == "${expected_ip}" && "${body}" == *'"Actual": "NotFound"'* ]]; then
        return 0
      fi
      echo "${description} returned HTTP ${http_code}, remote IP ${remote_ip}, body: ${body}" >&2
    fi
    if (( attempt < 5 )); then
      echo "${description} failed on attempt ${attempt}; retrying in 2 seconds" >&2
      sleep 2
    fi
  done
  fail "${description} failed after 5 attempts"
}

wait_for_endpoint_slices() {
  local attempt
  local family
  local ready_count
  local ready
  local endpoint_families=()

  for attempt in {1..30}; do
    endpoint_slices="$(
      kubectl --namespace "${namespace}" get endpointslices \
        --selector "kubernetes.io/service-name=${service_name}" -o json
    )"
    mapfile -t endpoint_families < <(
      jq -r '[.items[].addressType] | unique[]' <<<"${endpoint_slices}"
    )
    ready=true
    if [[ "${endpoint_families[*]}" != "${sorted_expected_families[*]}" ]]; then
      ready=false
    fi
    for family in "${expected_families[@]}"; do
      ready_count="$(
        jq -r --arg family "${family}" '
          [.items[] | select(.addressType == $family) |
           .endpoints[] | select(.conditions.ready != false) | .addresses[]] | length
        ' <<<"${endpoint_slices}"
      )"
      if (( ready_count == 0 )); then
        ready=false
      fi
    done
    if [[ "${ready}" == true ]]; then
      return 0
    fi
    if (( attempt < 30 )); then
      sleep 2
    fi
  done
  fail "EndpointSlices did not publish ready endpoints for families ${expected_families[*]} within 60 seconds"
}

if (( ${#expected_families[@]} < 1 || ${#expected_families[@]} > 2 )); then
  fail "Expected one or two IP families, found ${expected_families[*]:-<none>}"
fi
for family in "${expected_families[@]}"; do
  if [[ "${family}" != IPv4 && "${family}" != IPv6 ]]; then
    fail "Unsupported expected IP family ${family}"
  fi
done
command -v jq >/dev/null || fail "jq is required to verify EndpointSlices"

policy="$(kubectl --namespace "${namespace}" get service "${service_name}" -o jsonpath='{.spec.ipFamilyPolicy}')"
if [[ "${policy}" != "${expected_policy}" ]]; then
  fail "Expected ipFamilyPolicy ${expected_policy}, found ${policy}"
fi

mapfile -t families < <(
  kubectl --namespace "${namespace}" get service "${service_name}" \
    -o jsonpath='{range .spec.ipFamilies[*]}{@}{"\n"}{end}'
)
mapfile -t cluster_ips < <(
  kubectl --namespace "${namespace}" get service "${service_name}" \
    -o jsonpath='{range .spec.clusterIPs[*]}{@}{"\n"}{end}'
)

if [[ "${families[*]}" != "${expected_families[*]}" ]]; then
  fail "Expected ipFamilies ${expected_families[*]}, found ${families[*]:-<none>}"
fi
if [[ "${#cluster_ips[@]}" -ne "${#expected_families[@]}" ]]; then
  fail "Expected ${#expected_families[@]} cluster IPs, found ${cluster_ips[*]:-<none>}"
fi

declare -A seen_cluster_ips=()
for index in "${!expected_families[@]}"; do
  ip="${cluster_ips[${index}]}"
  if [[ "$(ip_family "${ip}")" != "${families[${index}]}" ]]; then
    fail "Expected ${families[${index}]} cluster IP at index ${index}, found ${ip}"
  fi
  if [[ -n "${seen_cluster_ips[${ip}]:-}" ]]; then
    fail "Expected unique cluster IPs, found duplicate ${ip}"
  fi
  seen_cluster_ips["${ip}"]=true
done

mapfile -t sorted_expected_families < <(printf '%s\n' "${expected_families[@]}" | sort -u)
wait_for_endpoint_slices

operator_pods=()
for family in "${expected_families[@]}"; do
  mapfile -t endpoint_addresses < <(
    jq -r --arg family "${family}" '
      .items[] | select(.addressType == $family) |
      .endpoints[] | select(.conditions.ready != false) | .addresses[]
    ' <<<"${endpoint_slices}"
  )
  if (( ${#endpoint_addresses[@]} == 0 )); then
    fail "Expected at least one ready ${family} endpoint for Service ${service_name}"
  fi
  for ip in "${endpoint_addresses[@]}"; do
    if [[ "$(ip_family "${ip}")" != "${family}" ]]; then
      fail "EndpointSlice declared ${family} but contained address ${ip}"
    fi
  done
  mapfile -t family_pods < <(
    jq -r --arg family "${family}" '
      .items[] | select(.addressType == $family) |
      .endpoints[] | select(.conditions.ready != false) |
      select(.targetRef.kind == "Pod") | .targetRef.name
    ' <<<"${endpoint_slices}"
  )
  operator_pods+=("${family_pods[@]}")
done
mapfile -t operator_pods < <(printf '%s\n' "${operator_pods[@]}" | sed '/^$/d' | sort -u)
if (( ${#operator_pods[@]} == 0 )); then
  fail "Expected EndpointSlices for ${service_name} to target at least one Pod"
fi
for pod_name in "${operator_pods[@]}"; do
  assert_dual_stack_pod "${pod_name}"
done

kubectl --namespace "${namespace}" delete pod "${probe_name}" --ignore-not-found
trap 'kubectl --namespace "${namespace}" delete pod "${probe_name}" --ignore-not-found --wait=false' EXIT
kubectl --namespace "${namespace}" run "${probe_name}" \
  --image=curlimages/curl:8.15.0 \
  --restart=Never \
  --command -- sleep 600
kubectl --namespace "${namespace}" wait --for=condition=Ready "pod/${probe_name}" --timeout=2m
assert_dual_stack_pod "${probe_name}"

status_path="/status/${namespace}/${missing_resource}"
service_host="${service_name}.${namespace}.svc"
for index in "${!expected_families[@]}"; do
  family="${expected_families[${index}]}"
  ip="${cluster_ips[${index}]}"
  if [[ "${family}" == IPv6 ]]; then
    direct_url="http://[${ip}]:8000${status_path}"
    curl_family=-6
  else
    direct_url="http://${ip}:8000${status_path}"
    curl_family=-4
  fi
  probe_http "Direct ${family} REST probe to ${direct_url}" "${ip}" "${direct_url}"
  dns_url="http://${service_host}:8000${status_path}"
  probe_http "DNS ${family} REST probe to ${dns_url}" "${ip}" "${curl_family}" "${dns_url}"
done

echo "Verified ${service_name}: policy=${policy}, families=${families[*]}, clusterIPs=${cluster_ips[*]}"
