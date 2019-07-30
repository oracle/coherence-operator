#!/usr/bin/env bash

script_name=$0
script_full_path=$(dirname "$0")

if [[ "${OPERATOR_NAMESPACE}" == "" ]]
then
    OPERATOR_NAMESPACE="default"
fi

./${script_full_path}/cleanup.sh

kubectl -n ${OPERATOR_NAMESPACE} create -f build/_output/crds/coherence_v1_coherencerole_crd.yaml
kubectl -n ${OPERATOR_NAMESPACE} create -f build/_output/crds/coherence_v1_coherencecluster_crd.yaml
kubectl -n ${OPERATOR_NAMESPACE} create -f build/_output/crds/coherence_v1_coherenceinternal_crd.yaml

kubectl -n ${OPERATOR_NAMESPACE} create -f build/_output/yaml/service_account.yaml
kubectl -n ${OPERATOR_NAMESPACE} create -f build/_output/yaml/role.yaml
kubectl -n ${OPERATOR_NAMESPACE} create -f build/_output/yaml/role_binding.yaml

