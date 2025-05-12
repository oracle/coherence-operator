#!/bin/sh
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

set -o errexit

OS=$(go env GOOS)
ARCH=$(go env GOARCH)

if [ "CERT_MANAGER_VERSION" = "" ];
then
  echo "CERT_MANAGER_VERSION is not set"
  exit 1
fi

helm uninstall trust-manager -n cert-manager
kubectl delete crd bundles.trust.cert-manager.io

${KUBECTL_CMD} delete -f https://github.com/cert-manager/cert-manager/releases/download/${CERT_MANAGER_VERSION}/cert-manager.yaml
