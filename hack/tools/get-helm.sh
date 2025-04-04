#!/usr/bin/env bash
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
set -o errexit


ROOT_DIR=$(pwd)
TOOLS_BIN=${ROOT_DIR}/build/tools/bin

UNAME_S=$(uname -s)
UNAME_M=$(uname -m)

HELM_VERSION=3.17.2

if [ "Darwin" = "${UNAME_S}" ]; then
  if [ "x86_64" = "${UNAME_M}" ]; then
    echo "Downloading helm ${VERSION} ${UNAME_S} ${UNAME_M}"
    curl -Ls https://get.helm.sh/helm-v${HELM_VERSION}-darwin-amd64.tar.gz -o helm.tar.gz
    tar -xvf helm.tar.gz
    mv darwin-amd64/helm ${TOOLS_BIN}/
    rm -rf darwin-amd64
    rm helm.tar.gz
  else
    echo "Downloading helm ${VERSION} ${UNAME_S} ${UNAME_M}"
    curl -Ls https://get.helm.sh/helm-v${HELM_VERSION}-darwin-arm64.tar.gz -o helm.tar.gz
    tar -xvf helm.tar.gz
    mv darwin-arm64/helm ${TOOLS_BIN}/
    rm -rf darwin-arm64
    rm helm.tar.gz
  fi
else
  if [ "x86_64" = "${UNAME_M}" ]; then
    echo "Downloading helm ${VERSION} ${UNAME_S} ${UNAME_M}"
    curl -Ls https://get.helm.sh/helm-v${HELM_VERSION}-linux-amd64.tar.gz -o helm.tar.gz
    tar -xvf helm.tar.gz
    mv linux-amd64/helm ${TOOLS_BIN}/
    rm -rf linux-amd64
    rm helm.tar.gz
  else
    echo "Downloading helm ${VERSION} ${UNAME_S} ${UNAME_M}"
    curl -Ls https://get.helm.sh/helm-v${HELM_VERSION}-linux-arm64.tar.gz -o helm.tar.gz
    tar -xvf helm.tar.gz
    mv linux-arm64/helm ${TOOLS_BIN}/
    rm -rf linux-arm64
    rm helm.tar.gz
  fi
fi

chmod +x ${TOOLS_BIN}/helm
