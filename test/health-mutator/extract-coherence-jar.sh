#!/usr/bin/env bash
# Copyright (c) 2026, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
set -euo pipefail

image=${1:?Coherence image is required}
version=${2:?Coherence version is required}
output=${3:?output jar path is required}
container=$(docker create "$image")
cleanup() {
  docker rm -f "$container" >/dev/null 2>&1 || true
}
trap cleanup EXIT

source="/app/libs/coherence-${version}.jar"
mkdir -p "$(dirname "$output")"
if ! docker cp "${container}:${source}" "$output"; then
  echo "No Coherence jar found at ${source} in ${image}" >&2
  exit 1
fi
test -s "$output"
