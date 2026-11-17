#!/usr/bin/env bash
# Copyright (c) 2026, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
set -euo pipefail
cd "$(dirname "$0")/../.."

# Restored timestamps predate a fresh checkout, and some generation dependencies
# are phony. Skip these completed targets explicitly.
# Validate first: an incomplete download must fail, rather than trigger a rebuild.
version=$(cat build/_output/version.txt)
chart="build/_output/helm-charts/coherence-operator-${version}.tgz"
completed=(
  build/_output/build.properties
  build/_output/targets/generate
  build/_output/targets/manifests
  build/_output/targets/java
  build/_output/targets/cli
  build/_output/targets/build-operator
  bin/runner
  "$chart"
)
required=(
  "${completed[@]}"
  api/v1/zz_generated.deepcopy.go
  pkg/data/assets/config.json
  config/crd/bases/coherence.oracle.com_coherence.yaml
  config/crd-small/bases/coherence.oracle.com_coherence.yaml
  build/_output/coherence-operator-manifests.tar.gz
  build/_output/coherence-operator.yaml
  build/_output/manifests/crd/coherence.oracle.com_coherence.yaml
  build/_output/helm-charts/coherence-operator/Chart.yaml
)
for file in "${required[@]}"; do
  test -f "$file" || { echo "Missing shared build artifact: $file" >&2; exit 1; }
done
test -x bin/runner
args=()
for file in "${completed[@]}"; do
  # The Makefile uses absolute paths for build output targets.
  args+=(--old-file="$PWD/$file")
done
# GNU make does not propagate --old-file to sub-makes. Route $(MAKE) through
# this wrapper too. MF is an existing recipe's manual MAKEFLAGS expansion; clear
# it here because flags/overrides already travel through make's environment.
exec make "${args[@]}" "MAKE=$PWD/hack/ci/make-prebuilt.sh" MF= "$@"
