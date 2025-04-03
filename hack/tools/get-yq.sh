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

YQ_VERSION = v4.45.1

if [ "Darwin" = "${UNAME_S}" ]; then
  if [ "x86_64" = "${UNAME_M}" ]; then
    echo "Downloading yq ${VERSION} ${UNAME_S} ${UNAME_M}"
  	curl -L https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_darwin_amd64 -o ${TOOLS_BIN}/yq
  else
    echo "Downloading yq ${VERSION} ${UNAME_S} ${UNAME_M}"
  	curl -L https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_darwin_arm64 -o ${TOOLS_BIN}/yq
  fi
else
  if [ "x86_64" = "${UNAME_M}" ]; then
    echo "Downloading yq ${VERSION} ${UNAME_S} ${UNAME_M}"
  	curl -L https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_linux_amd64 -o ${TOOLS_BIN}/yq
  else
    echo "Downloading yq ${VERSION} ${UNAME_S} ${UNAME_M}"
  	curl -L https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_linux_arm64 -o ${TOOLS_BIN}/yq
  fi
fi

chmod +x ${TOOLS_BIN}/yq
