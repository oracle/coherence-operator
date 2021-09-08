#!/bin/sh
#
# Copyright (c) 2021, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

ISTIO_VERSION=$1
TOOLS_DIRECTORY=$2

if [ "${ISTIO_VERSION}" = "" ]
then
  ISTIO_VERSION="$(curl -sL https://github.com/istio/istio/releases | \
                    grep -o 'releases/[0-9]*.[0-9]*.[0-9]*/' | sort --version-sort | \
                    tail -1 | awk -F'/' '{ print $2}')"
  ISTIO_VERSION="${ISTIO_VERSION##*/}"
fi

echo "Getting Istio version ${ISTIO_VERSION}"

cd "${TOOLS_DIRECTORY}" || exit

ISTIO_HOME=${TOOLS_DIRECTORY}/istio-${ISTIO_VERSION}

if [ ! -d "${ISTIO_HOME}" ]; then
  echo "Istio will be installed into ${ISTIO_HOME}"
  mkdir -p ${ISTIO_HOME} || true
  curl -sL https://istio.io/downloadIstio | ISTIO_VERSION=${ISTIO_VERSION} sh -
fi

