#!/bin/sh
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

set -eu

ISTIO_VERSION_FILE=${1:?Usage: find-istio-version.sh CACHE_FILE}

valid_version() {
  case "$1" in
    ''|*[!0-9.]*) return 1 ;;
  esac
  printf '%s\n' "$1" | LC_ALL=C grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$'
}

if [ -f "$ISTIO_VERSION_FILE" ]; then
  ISTIO_VERSION=$(cat "$ISTIO_VERSION_FILE")
  if valid_version "$ISTIO_VERSION"; then
    printf '%s\n' "$ISTIO_VERSION"
    exit 0
  fi
fi

# Do not cache failed lookups or depend on GitHub's releases page HTML.
RELEASE=$(curl -fsSL --retry 3 -H 'Accept: application/vnd.github+json' \
  https://api.github.com/repos/istio/istio/releases/latest) || {
  echo "Unable to fetch the latest Istio release" >&2
  exit 1
}
ISTIO_VERSION=$(printf '%s\n' "$RELEASE" | jq -er '.tag_name') || {
  echo "Unable to read the latest Istio release tag" >&2
  exit 1
}
ISTIO_VERSION=${ISTIO_VERSION#v}
if ! valid_version "$ISTIO_VERSION"; then
  echo "Invalid Istio release version: $ISTIO_VERSION" >&2
  exit 1
fi

mkdir -p "$(dirname "$ISTIO_VERSION_FILE")"
CACHE_TEMP=$(mktemp "${ISTIO_VERSION_FILE}.XXXXXX")
trap 'rm -f "$CACHE_TEMP"' EXIT HUP INT TERM
printf '%s\n' "$ISTIO_VERSION" > "$CACHE_TEMP"
mv "$CACHE_TEMP" "$ISTIO_VERSION_FILE"
printf '%s\n' "$ISTIO_VERSION"
