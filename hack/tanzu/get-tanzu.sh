#!/bin/sh
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

TANZU_VERSION=$1
TOOLS_DIRECTORY=$2

if [ "${TANZU_VERSION}" = "" -o "${TANZU_VERSION}" = "latest" ]
then
  TANZU_VERSION="$(curl -s -H "Accept: application/vnd.github.v3+json" \
          https://api.github.com/repos/vmware-tanzu/community-edition/releases \
          | jq '.[0].tag_name' |  tr -d '"')"
  TANZU_VERSION="${TANZU_VERSION##*/}"
fi

echo "Getting Tanzu version ${TANZU_VERSION}"

cd "${TOOLS_DIRECTORY}" || exit

TANZU_HOME=${TOOLS_DIRECTORY}/tanzu/${TANZU_VERSION}
echo "Tanzu will be installed into ${TANZU_HOME}"

if [ ! -d "${TANZU_HOME}" ]; then
  mkdir -p ${TANZU_HOME} || true

  UNAME_S=$(uname -s)
  if [ "${UNAME_S}" = "Darwin" ]; then
    echo "Downloading MacOS amd64 Tanzu into ${TANZU_HOME}"
    URL="https://github.com/vmware-tanzu/community-edition/releases/download/${TANZU_VERSION}/tce-darwin-amd64-${TANZU_VERSION}.tar.gz"
  else
    echo "Downloading Linux amd64 Tanzu into ${TANZU_HOME}"
    URL="https://github.com/vmware-tanzu/community-edition/releases/download/${TANZU_VERSION}/tce-linux-amd64-${TANZU_VERSION}.tar.gz"
  fi

  curl -L ${URL} -o ${TANZU_HOME}/tanzu.tar.gz
  echo "Extracting Tanzu into ${TANZU_HOME}"

  cd "${TANZU_HOME}" || exit 1

  tar --strip-components=1 -xf tanzu.tar.gz
  rm tanzu.tar.gz

  ./install.sh

  tanzu version
fi

curl -L https://carvel.dev/install.sh | bash

exit 0
