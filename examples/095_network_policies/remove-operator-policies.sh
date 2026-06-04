#!/usr/bin/env bash

SCRIPT_PATH=${0}
WORK_DIR=$(dirname -- "${SCRIPT_PATH}")

if [ "${NAMESPACE}" == "" ]
then
    NAMESPACE="coherence"
fi

kubectl -n ${NAMESPACE} delete -f "${WORK_DIR}/manifests/deny-all.yaml"
kubectl -n ${NAMESPACE} delete -f "${WORK_DIR}/manifests/allow-dns-kube-system.yaml"
# Keep this list aligned with add-operator-policies.sh so teardown only removes
# policies that are still installed by the current examples.
kubectl -n ${NAMESPACE} delete -f "${WORK_DIR}/manifests/allow-k8s-api-server.yaml"
kubectl -n ${NAMESPACE} delete -f "${WORK_DIR}/manifests/allow-operator-rest-ingress.yaml"
kubectl -n ${NAMESPACE} delete -f "${WORK_DIR}/manifests/allow-operator-cluster-member-egress.yaml"
