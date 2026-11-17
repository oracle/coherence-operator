#!/bin/sh
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

set -eu

ISTIO_VERSION=${1-}
TOOLS_DIRECTORY=${2:?Usage: get-istio-latest.sh VERSION TOOLS_DIRECTORY}
SCRIPT_DIRECTORY=$(CDPATH='' cd -- "$(dirname "$0")" && pwd)
mkdir -p "$TOOLS_DIRECTORY"
TOOLS_DIRECTORY=$(CDPATH='' cd -- "$TOOLS_DIRECTORY" && pwd)

if [ -z "$ISTIO_VERSION" ] || [ "$ISTIO_VERSION" = latest ]; then
  ISTIO_VERSION=$("$SCRIPT_DIRECTORY/find-istio-version.sh" "$TOOLS_DIRECTORY/istio-latest.txt")
fi
case "$ISTIO_VERSION" in
  ''|*[!0-9.]*)
    echo "Invalid Istio version: $ISTIO_VERSION" >&2
    exit 1
    ;;
esac
if ! printf '%s\n' "$ISTIO_VERSION" | LC_ALL=C grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$'; then
  echo "Invalid Istio version: $ISTIO_VERSION" >&2
  exit 1
fi

echo "Getting Istio version $ISTIO_VERSION"
ISTIO_HOME=$TOOLS_DIRECTORY/istio-$ISTIO_VERSION

# A directory left by a failed download is not a usable installation.
if [ ! -x "$ISTIO_HOME/bin/istioctl" ]; then
  echo "Istio will be installed into $ISTIO_HOME"
  INSTALLER=$(mktemp "$TOOLS_DIRECTORY/istio-installer.XXXXXX")
  trap 'rm -f "$INSTALLER"' EXIT HUP INT TERM
  curl -fsSL --retry 3 https://istio.io/downloadIstio -o "$INSTALLER"
  cd "$TOOLS_DIRECTORY"
  ISTIO_VERSION=$ISTIO_VERSION sh "$INSTALLER"
  if [ ! -x "$ISTIO_HOME/bin/istioctl" ]; then
    echo "Istio installer did not create $ISTIO_HOME/bin/istioctl" >&2
    exit 1
  fi
fi
