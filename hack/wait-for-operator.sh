#!/usr/bin/env bash

#
# Copyright (c) 2023, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

if [[ "${OPERATOR_NAMESPACE}" == "" ]]; then
  OPERATOR_NAMESPACE="operator-test"
fi

echo "Waiting for Operator to start in namespace ${OPERATOR_NAMESPACE}"

POD=$(kubectl -n "${OPERATOR_NAMESPACE}" get pod -l control-plane=coherence -o name)

echo "Operator Pods:"
kubectl -n "${OPERATOR_NAMESPACE}" get pod -l control-plane=coherence
echo "Waiting for Operator to be ready. Pod: $(POD)"

if ! kubectl -n "${OPERATOR_NAMESPACE}" wait --for condition=ready --timeout 480s "${POD}"; then
  echo "Operator Pod ${POD} failed to start"
  if [[ "${TEST_LOGS_DIR}" == "" ]]; then
    TEST_LOGS_DIR="build/_output/test-logs"
  fi
  mkdir -p "${TEST_LOGS_DIR}" || true
  kubectl -n "${OPERATOR_NAMESPACE}" logs "${POD}" >> "${TEST_LOGS_DIR}/${POD}.log"
  exit 1
fi
