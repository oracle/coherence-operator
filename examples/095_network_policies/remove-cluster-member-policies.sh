#!/usr/bin/env bash

SCRIPT_PATH=${0}
WORK_DIR=$(dirname -- "${SCRIPT_PATH}")

if [ "${NAMESPACE}" == "" ]
then
    NAMESPACE="datastore"
fi

kubectl -n ${NAMESPACE} delete -f "${WORK_DIR}/manifests/deny-all.yaml"
kubectl -n ${NAMESPACE} delete -f "${WORK_DIR}/manifests/allow-dns-kube-system.yaml"
kubectl -n ${NAMESPACE} delete -f "${WORK_DIR}/manifests/allow-cluster-member-access.yaml"
kubectl -n ${NAMESPACE} delete -f "${WORK_DIR}/manifests/allow-cluster-member-operator-access.yaml"
kubectl -n ${NAMESPACE} delete -f "${WORK_DIR}/manifests/allow-extend-ingress.yaml"
kubectl -n ${NAMESPACE} delete -f "${WORK_DIR}/manifests/allow-grpc-ingress.yaml"
kubectl -n ${NAMESPACE} delete -f "${WORK_DIR}/manifests/allow-metrics-ingress.yaml"


