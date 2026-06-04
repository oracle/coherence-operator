#!/usr/bin/env bash

SCRIPT_PATH=${0}
WORK_DIR=$(dirname -- "${SCRIPT_PATH}")

if [ "${NAMESPACE}" == "" ]
then
    NAMESPACE="coherence"
fi

kubectl -n ${NAMESPACE} apply -f "${WORK_DIR}/manifests/deny-all.yaml"
kubectl -n ${NAMESPACE} apply -f "${WORK_DIR}/manifests/allow-dns-kube-system.yaml"
# The operator must keep API-server egress after deny-all is applied; otherwise
# it cannot watch resources or become ready during the network policy tests.
kubectl -n ${NAMESPACE} apply -f "${WORK_DIR}/manifests/allow-k8s-api-server.yaml"
kubectl -n ${NAMESPACE} apply -f "${WORK_DIR}/manifests/allow-operator-rest-ingress.yaml"
kubectl -n ${NAMESPACE} apply -f "${WORK_DIR}/manifests/allow-operator-cluster-member-egress.yaml"
