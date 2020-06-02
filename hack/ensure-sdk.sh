#!/bin/sh

REQUIRED_VERSION=$1

CURRDIR=$(pwd)
UNAME_S=$(uname -s)
UNAME_M=$(uname -m)
OPERATOR_SDK=${CURRDIR}/etc/sdk/${UNAME_S}-${UNAME_M}/operator-sdk

mkdir -p ${CURRDIR}/etc/sdk/${UNAME_S}-${UNAME_M} || true

VERSION=$(${OPERATOR_SDK} version)
echo "${VERSION}"
echo ${VERSION} | grep "operator-sdk version: \"${REQUIRED_VERSION}\""
if [[ $? == 1 ]]; then
  if [[ "${UNAME_S}" == "Linux" ]]; then
    curl -L https://github.com/operator-framework/operator-sdk/releases/download/${REQUIRED_VERSION}/operator-sdk-${REQUIRED_VERSION}-x86_64-linux-gnu -o ${OPERATOR_SDK}
  fi
  if [[ "${UNAME_S}" == "Darwin" ]]; then
    curl -L https://github.com/operator-framework/operator-sdk/releases/download/${REQUIRED_VERSION}/operator-sdk-${REQUIRED_VERSION}-x86_64-apple-darwin -o ${OPERATOR_SDK}
  fi
  chmod +x ${OPERATOR_SDK}
fi

