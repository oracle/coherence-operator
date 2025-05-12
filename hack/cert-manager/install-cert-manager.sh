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

curl -fsSL -o cmctl https://github.com/cert-manager/cmctl/releases/latest/download/cmctl_${OS}_${ARCH}
chmod +x cmctl
mv cmctl ${TOOLS_BIN}

${KUBECTL_CMD} apply -f https://github.com/cert-manager/cert-manager/releases/download/${CERT_MANAGER_VERSION}/cert-manager.yaml
${TOOLS_BIN}/cmctl check api --wait=10m

helm repo add jetstack https://charts.jetstack.io --force-update

helm upgrade trust-manager jetstack/trust-manager \
  --install \
  --namespace cert-manager \
  --wait
