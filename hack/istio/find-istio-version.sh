#!/bin/sh
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

ISTIO_VERSION_FILE=$1
ISTIO_VERSION=""

if [ -e $1 ]
then
  ISTIO_VERSION=$(cat $1)
else
if [ "${ISTIO_VERSION}" = "" -o "${ISTIO_VERSION}" = "latest" ]
then
  ISTIO_VERSION="$(curl -sL https://github.com/istio/istio/releases | \
                    grep -o 'releases/[0-9]*.[0-9]*.[0-9]*/' | sort --version-sort | \
                    tail -1 | awk -F'/' '{ print $2}')"
  ISTIO_VERSION="${ISTIO_VERSION##*/}"
  echo ${ISTIO_VERSION} > $1
fi
fi

echo ${ISTIO_VERSION}
