#!/usr/bin/env bash

if [[ "${OPERATOR_NAMESPACE}" == "" ]]
then
    OPERATOR_NAMESPACE="default"
fi

kubectl -n ${OPERATOR_NAMESPACE} delete -f deploy/role_binding.yaml
kubectl -n ${OPERATOR_NAMESPACE} delete -f deploy/role.yaml
kubectl -n ${OPERATOR_NAMESPACE} delete -f deploy/service_account.yaml

echo "Cleaning up Coherence CRDs"
kubectl -n ${OPERATOR_NAMESPACE} delete -f deploy/crds/coherence_v1_coherencerole_crd.yaml
kubectl -n ${OPERATOR_NAMESPACE} delete -f deploy/crds/coherence_v1_coherencecluster_crd.yaml
kubectl -n ${OPERATOR_NAMESPACE} delete -f deploy/crds/coherence_v1_coherenceinternal_crd.yaml

echo "Remaining CRDs:"
kubectl get crd

