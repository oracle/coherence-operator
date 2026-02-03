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

PREFLIGHT_VERSION=1.16.0
PREFLIGHT_ROOT_URL="https://github.com/redhat-openshift-ecosystem/openshift-preflight/releases/download/${PREFLIGHT_VERSION}"

if [ "Darwin" = "${UNAME_S}" ]; then
  if [ "x86_64" = "${UNAME_M}" ]; then
    echo "Downloading OpenShift Preflight CLI ${UNAME_S} ${UNAME_M}"
  	curl -Ls ${PREFLIGHT_ROOT_URL}/preflight-darwin-amd64 -o ${TOOLS_BIN}/preflight
  else
    echo "Downloading OpenShift Preflight CLI ${UNAME_S} ${UNAME_M}"
    curl -Ls ${PREFLIGHT_ROOT_URL}/preflight-darwin-arm64 -o ${TOOLS_BIN}/preflight
  fi
else
  if [ "x86_64" = "${UNAME_M}" ]; then
    echo "Downloading OpenShift Preflight CLI ${UNAME_S} ${UNAME_M}"
    curl -Ls ${PREFLIGHT_ROOT_URL}/preflight-linux-amd64 -o ${TOOLS_BIN}/preflight
  else
    echo "Downloading OpenShift Preflight CLI ${UNAME_S} ${UNAME_M}"
    curl -Ls ${PREFLIGHT_ROOT_URL}/preflight-linux-arm64 -o ${TOOLS_BIN}/preflight
  fi
fi

chmod +x ${TOOLS_BIN}/preflight
preflight

