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

if [ "${TEKTON_VERSION}" == "" ]
then
  TEKTON_VERSION=0.40.0
fi

if [ "Darwin" = "${UNAME_S}" ]; then
  echo "Downloading Tekton ${VERSION} ${UNAME_S} ${UNAME_M}"
	curl -Ls https://github.com/tektoncd/cli/releases/download/v${TEKTON_VERSION}/tkn_${TEKTON_VERSION}_Darwin_all.tar.gz -o tekton.tar.gz
else
  if [ "x86_64" = "${UNAME_M}" ]; then
    echo "Downloading Tekton ${VERSION} ${UNAME_S} ${UNAME_M}"
  	curl -Ls https://github.com/tektoncd/cli/releases/download/v${TEKTON_VERSION}/tkn_${TEKTON_VERSION}_Linux_x86_64.tar.gz -o tekton.tar.gz
  else
    echo "Downloading Tekton ${VERSION} ${UNAME_S} ${UNAME_M}"
	curl -Ls https://github.com/tektoncd/cli/releases/download/v${TEKTON_VERSION}/tkn_${TEKTON_VERSION}_Linux_aarch64.tar.gz -o tekton.tar.gz
  fi
fi

tar -C ${TOOLS_BIN}/ -xvf tekton.tar.gz
rm tekton.tar.gz
rm ${TOOLS_BIN}/LICENSE || true
rm ${TOOLS_BIN}/README.md || true
chmod +x ${TOOLS_BIN}/tkn
