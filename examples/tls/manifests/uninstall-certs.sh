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

kubectl -n ${COH_NAMESPACE} delete -f ${BASEDIR}/client-keystore.yaml
kubectl -n ${COH_NAMESPACE} delete secret client-keystore-secret
kubectl -n ${COH_NAMESPACE} delete -f ${BASEDIR}/server-keystore.yaml
kubectl -n ${COH_NAMESPACE} delete secret server-keystore-secret

kubectl delete -f ${BASEDIR}/ca-issuer.yaml

kubectl -n ${CM_NAMESPACE} delete -f ${BASEDIR}/ca-cert.yaml
kubectl delete -f ${BASEDIR}/selfsigned-issuer.yaml
