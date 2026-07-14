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
probe_name="coherence-operator-rest-dual-stack-probe"

policy="$(kubectl --namespace "${namespace}" get service "${service_name}" -o jsonpath='{.spec.ipFamilyPolicy}')"
if [[ "${policy}" != "RequireDualStack" ]]; then
  echo "Expected ipFamilyPolicy RequireDualStack, found ${policy}" >&2
  exit 1
fi

mapfile -t families < <(
  kubectl --namespace "${namespace}" get service "${service_name}" \
    -o jsonpath='{range .spec.ipFamilies[*]}{@}{"\n"}{end}'
)
mapfile -t cluster_ips < <(
  kubectl --namespace "${namespace}" get service "${service_name}" \
    -o jsonpath='{range .spec.clusterIPs[*]}{@}{"\n"}{end}'
)

expected_families=(IPv4 IPv6)
if [[ "${#families[@]}" -ne 2 || "${families[0]}" != "${expected_families[0]}" || "${families[1]}" != "${expected_families[1]}" ]]; then
  echo "Expected ipFamilies IPv4,IPv6, found ${families[*]:-<none>}" >&2
  exit 1
fi
if [[ "${#cluster_ips[@]}" -ne 2 ]]; then
  echo "Expected two cluster IPs, found ${cluster_ips[*]:-<none>}" >&2
  exit 1
fi

for index in 0 1; do
  ip="${cluster_ips[${index}]}"
  if [[ "${families[${index}]}" == "IPv4" && "${ip}" == *:* ]]; then
    echo "Expected IPv4 cluster IP at index ${index}, found ${ip}" >&2
    exit 1
  fi
  if [[ "${families[${index}]}" == "IPv6" && "${ip}" != *:* ]]; then
    echo "Expected IPv6 cluster IP at index ${index}, found ${ip}" >&2
    exit 1
  fi
done

kubectl --namespace "${namespace}" delete pod "${probe_name}" --ignore-not-found
trap 'kubectl --namespace "${namespace}" delete pod "${probe_name}" --ignore-not-found --wait=false' EXIT
kubectl --namespace "${namespace}" run "${probe_name}" \
  --image=curlimages/curl:8.15.0 \
  --restart=Never \
  --command -- sleep 600
kubectl --namespace "${namespace}" wait --for=condition=Ready "pod/${probe_name}" --timeout=2m

for ip in "${cluster_ips[@]}"; do
  if [[ "${ip}" == *:* ]]; then
    url="http://[${ip}]:8000/"
  else
    url="http://${ip}:8000/"
  fi

  # The root path may return any HTTP status; receiving a response is sufficient
  # because this check verifies REST listener reachability through each Service IP.
  reachable=false
  for attempt in 1 2 3 4 5; do
    if kubectl --namespace "${namespace}" exec "${probe_name}" -- \
      curl --noproxy '*' --silent --show-error --connect-timeout 5 --max-time 10 --output /dev/null "${url}"; then
      reachable=true
      break
    fi
    if ((attempt < 5)); then
      echo "REST probe to ${url} failed on attempt ${attempt}; retrying in 2 seconds" >&2
      sleep 2
    fi
  done
  if [[ "${reachable}" != "true" ]]; then
    echo "REST probe to ${url} failed after 5 attempts" >&2
    exit 1
  fi
done

echo "Verified ${service_name} through IPv4 ${cluster_ips[0]} and IPv6 ${cluster_ips[1]}"
