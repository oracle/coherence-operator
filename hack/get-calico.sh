#!/bin/sh
#
# Copyright (c) 2022, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

CALICO_VERSION=$1

if [ "${CALICO_VERSION}" = "" -o "${CALICO_VERSION}" = "latest" ]
then
  CALICO_VERSION="$(curl -sL https://github.com/projectcalico/calico/releases | \
      grep -o 'releases/tag/v[0-9]*.[0-9]*.[0-9]' | \
      sort --version-sort | \
      tail -1 | awk -F'/' '{ print $3}' | awk -F'.' '{ print $1"."$2}')"

  CALICO_VERSION="${CALICO_VERSION##*/}"
fi

echo "Getting Calico version ${CALICO_VERSION}"

kubectl apply -f "https://docs.projectcalico.org/${CALICO_VERSION}/manifests/calico.yaml"
kubectl -n kube-system set env daemonset/calico-node FELIX_IGNORELOOSERPF=true
kubectl -n kube-system wait --for condition=ready --timeout=300s -l k8s-app=calico-node pod
kubectl -n kube-system wait --for condition=ready --timeout=300s -l k8s-app=kube-dns pod
