#!/usr/bin/sh
#
# Copyright (c) 2021., Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
#

BASEDIR=$(dirname "$0")

CM_NAMESPACE=cert-manager
COH_NAMESPACE=coherence-test

# Install Cert Manager
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.2.0/cert-manager.yaml

# Wait for Cert Manager Pod
POD=$(kubectl -n ${CM_NAMESPACE} get pod -l app=cert-manager -o name)
kubectl -n ${CM_NAMESPACE} wait --for=condition=Ready ${POD}

# Wait for Cert Manager CA Injector Pod
POD=$(kubectl -n ${CM_NAMESPACE} get pod -l app=cainjector -o name)
kubectl -n ${CM_NAMESPACE} wait --for=condition=Ready ${POD}

# Wait for Cert Manager Web-Hook Pod
POD=$(kubectl -n ${CM_NAMESPACE} get pod -l app=webhook -o name)
kubectl -n ${CM_NAMESPACE} wait --for=condition=Ready ${POD}

kubectl apply -f ${BASEDIR}/selfsigned-issuer.yaml
kubectl wait --for=condition=Ready clusterissuer/selfsigned-issuer

kubectl -n ${CM_NAMESPACE} apply -f manifests/ca-cert.yaml
kubectl -n ${CM_NAMESPACE} wait --for=condition=Ready certificate/ca-certificate

kubectl apply -f manifests/ca-issuer.yaml
kubectl wait --for=condition=Ready clusterissuer/ca-issuer


kubectl create ns ${COH_NAMESPACE}

kubectl -n ${COH_NAMESPACE} create secret generic server-keystore-secret --from-literal=password-key=password
kubectl -n ${COH_NAMESPACE} apply -f manifests/server-keystore.yaml
kubectl -n ${COH_NAMESPACE} wait --for=condition=Ready certificate/server-keystore

kubectl -n ${COH_NAMESPACE} create secret generic client-keystore-secret --from-literal=password-key=secret
kubectl -n ${COH_NAMESPACE} apply -f manifests/client-keystore.yaml
kubectl -n ${COH_NAMESPACE} wait --for=condition=Ready certificate/client-keystore

