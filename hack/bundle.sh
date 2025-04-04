#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

set -o errexit

ROOT_DIR=$(pwd)
BUILD_OUTPUT=${ROOT_DIR}/build/_output
BUILD_BUNDLE="${BUILD_OUTPUT}/bundle"
VERSION=$(cat "${BUILD_OUTPUT}/version.txt")

echo "ROOT_DIR="${ROOT_DIR}""
echo "GITHUB_REPO=${GITHUB_REPO}"
echo "VERSION=${VERSION}"

BUNDLE_DIR="${BUILD_BUNDLE}/coherence-operator"
mkdir -p "${BUNDLE_DIR}/${VERSION}/manifests" || true
mkdir -p "${BUNDLE_DIR}/${VERSION}/metadata" || true
mkdir -p "${BUNDLE_DIR}/${VERSION}/tests/scorecard" || true

for FILE in ${ROOT_DIR}/bundle/manifests/*; do
  NAME=$(basename ${FILE})
  echo "---" > "${BUNDLE_DIR}/${VERSION}/manifests/${NAME}"
  cat "${FILE}" >> "${BUNDLE_DIR}/${VERSION}/manifests/${NAME}"
done

for FILE in ${ROOT_DIR}/bundle/metadata/*; do
  NAME=$(basename ${FILE})
  echo "---" > "${BUNDLE_DIR}/${VERSION}/metadata/${NAME}"
  cat "${FILE}" >> "${BUNDLE_DIR}/${VERSION}/metadata/${NAME}"
done

for FILE in ${ROOT_DIR}/bundle/tests/scorecard/*; do
  NAME=$(basename ${FILE})
  echo "---" > "${BUNDLE_DIR}/${VERSION}/tests/scorecard/${NAME}"
  cat "${FILE}" >> "${BUNDLE_DIR}/${VERSION}/tests/scorecard/${NAME}"
done
