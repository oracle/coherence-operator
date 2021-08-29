#!/bin/sh

REQUIRED_VERSION=$1
OPERATOR_SDK_HOME=$2

UNAME_S=$(uname -s)
OPERATOR_SDK=${OPERATOR_SDK_HOME}/operator-sdk
OK=0

if [ ! -f "${OPERATOR_SDK}" ]; then
#  Operator SDK does not exist
  echo "Operator SDK not found at ${OPERATOR_SDK}"
  mkdir -p ${OPERATOR_SDK_HOME} || true
  OK=1
else
#  Operator SDK exists, check its version
  echo "Operator SDK found, checking version - ${OPERATOR_SDK}"
  VERSION=$(${OPERATOR_SDK} version)
  echo "${VERSION}"
  echo ${VERSION} | grep "operator-sdk version: \"${REQUIRED_VERSION}\""
  OK=$?
fi

if [ ${OK} != 0 ]; then
  echo "Operator SDK not found or not correct version"

  if [ "${UNAME_S}" = "Darwin" ]; then
    URL="https://github.com/operator-framework/operator-sdk/releases/download/${REQUIRED_VERSION}/operator-sdk_darwin_amd64"
  else
    URL="https://github.com/operator-framework/operator-sdk/releases/download/${REQUIRED_VERSION}/operator-sdk-linux-amd64"
  fi

  echo "Downloading Operator SDK ${UNAME_S} version from ${URL}"
  curl -L ${URL} -o ${OPERATOR_SDK}
  chmod +x ${OPERATOR_SDK}
fi

