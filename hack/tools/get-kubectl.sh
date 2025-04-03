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

VERSION=$(curl -L -s https://dl.k8s.io/release/stable.txt)

if [ "Darwin" = "${UNAME_S}" ]; then
  if [ "x86_64" = "${UNAME_M}" ]; then
    echo "Downloading kubectl ${VERSION} ${UNAME_S} ${UNAME_M}"
    curl -Ls "https://dl.k8s.io/release/${version}/bin/darwin/amd64/kubectl" -o ${TOOLS_BIN}/kubectl
  else
    echo "Downloading kubectl ${VERSION} ${UNAME_S} ${UNAME_M}"
    curl -Ls "https://dl.k8s.io/release/${version}/bin/darwin/arm64/kubectl" -o ${TOOLS_BIN}/kubectl
  fi
else
  if [ "x86_64" = "${UNAME_M}" ]; then
    echo "Downloading kubectl ${VERSION} ${UNAME_S} ${UNAME_M}"
    curl -Ls "https://dl.k8s.io/release/${version}/bin/linux/amd64/kubectl" -o ${TOOLS_BIN}/kubectl
  else
    echo "Downloading kubectl ${VERSION} ${UNAME_S} ${UNAME_M}"
    curl -Ls "https://dl.k8s.io/release/${version}/bin/linux/arm64/kubectl" -o ${TOOLS_BIN}/kubectl
  fi
fi

chmod +x ${TOOLS_BIN}/kubectl
