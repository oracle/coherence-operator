#!/usr/bin/env bash
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

set -o -v errexit

GITHUB_REPO=$(git config --get remote.origin.url)
BUILD_OUTPUT=build/_output
RELEASE_ASSETS_DIR=${BUILD_OUTPUT}/github-release-assets

VERSION=$(cat ${BUILD_OUTPUT}/version.txt)

UNAME_S=$(uname -s)

echo ${GITHUB_REPO}
echo ${VERSION}
echo ${UNAME_S}

gh repo set-default ${GITHUB_REPO}

# create a temporary location for the release assets as we need to rename some
rm -rf ${RELEASE_ASSETS_DIR}
mkdir -p ${RELEASE_ASSETS_DIR}
cp ${BUILD_OUTPUT}/coherence-operator-manifests.tar.gz ${RELEASE_ASSETS_DIR}
cp ${BUILD_OUTPUT}/coherence-operator.yaml ${RELEASE_ASSETS_DIR}
cp ${BUILD_OUTPUT}/coherence-operator-restricted.yaml ${RELEASE_ASSETS_DIR}
cp ${BUILD_OUTPUT}/manifests/crd/coherence.oracle.com_coherence.yaml ${RELEASE_ASSETS_DIR}
cp ${BUILD_OUTPUT}/manifests/crd-small/coherence.oracle.com_coherence.yaml ${RELEASE_ASSETS_DIR}/coherence.oracle.com_coherence_small.yaml
cp ${BUILD_OUTPUT}/manifests/crd/coherencejob.oracle.com_coherence.yaml ${RELEASE_ASSETS_DIR}
cp ${BUILD_OUTPUT}/manifests/crd-small/coherencejob.oracle.com_coherence.yaml ${RELEASE_ASSETS_DIR}/coherencejob.oracle.com_coherence_small.yaml
cp ${BUILD_OUTPUT}/dashboards/${VERSION}/coherence-dashboards.tar.gz ${RELEASE_ASSETS_DIR}
cp ${BUILD_OUTPUT}/coherence-operator-bundle.tar.gz ${RELEASE_ASSETS_DIR}
cp ${BUILD_OUTPUT}/docs.zip ${RELEASE_ASSETS_DIR}

K8S_VERSIONS=$(yq '.jobs.build.strategy.matrix.matrixName[]' .github/workflows/k8s-matrix.yaml | tr '\n' ' ')
OS_VERSIONS=$(cat ${BUILD_OUTPUT}/openshift-version.txt)
cp hack/github/release-template.md ${RELEASE_ASSETS_DIR}/release-template.md
if [ "Darwin" = "${UNAME_S}" ]; then
  sed -i '' -e "s^VERSION_PLACEHOLDER^${VERSION}^g" ${RELEASE_ASSETS_DIR}/release-template.md
  sed -i '' -e "s^K8S_VERSIONS_PLACEHOLDER^${K8S_VERSIONS}^g" ${RELEASE_ASSETS_DIR}/release-template.md
  sed -i '' -e "s^OPENSHIFT_VERSIONS_PLACEHOLDER^${OS_VERSIONS}^g" ${RELEASE_ASSETS_DIR}/release-template.md
else
  sed -i -e "s^VERSION_PLACEHOLDER^${VERSION}^g" ${RELEASE_ASSETS_DIR}/release-template.md
  sed -i -e "s^K8S_VERSIONS_PLACEHOLDER^${K8S_VERSIONS}^g" ${RELEASE_ASSETS_DIR}/release-template.md
  sed -i -e "s^OPENSHIFT_VERSIONS_PLACEHOLDER^${OS_VERSIONS}^g" ${RELEASE_ASSETS_DIR}/release-template.md
fi

gh release delete v${VERSION} --yes || true
gh release create v${VERSION} --draft \
    --verify-tag \
    --generate-notes \
    --target ${RELEASE_BRANCH} \
    --latest \
    --notes-file "${RELEASE_ASSETS_DIR}/release-template.md"

# Upload all the various assets to the release
gh release upload v${VERSION} "${RELEASE_ASSETS_DIR}/coherence-operator-manifests.tar.gz#Coherence Operator all install manifest files (.tar.gz)"
gh release upload v${VERSION} "${RELEASE_ASSETS_DIR}/coherence-operator.yaml#Coherence Operator single install manifest file (.yaml)"
gh release upload v${VERSION} "${RELEASE_ASSETS_DIR}/coherence-operator-restricted.yaml#Coherence Operator single restricted install manifest file (.yaml)"
gh release upload v${VERSION} "${RELEASE_ASSETS_DIR}/coherence.oracle.com_coherence.yaml#Coherence full CRD (.yaml)"
gh release upload v${VERSION} "${RELEASE_ASSETS_DIR}/coherence.oracle.com_coherence_small.yaml#Coherence reduced size CRD (.yaml)"
gh release upload v${VERSION} "${RELEASE_ASSETS_DIR}/coherencejob.oracle.com_coherence.yaml#CoherenceJob full CRD (.yaml)"
gh release upload v${VERSION} "${RELEASE_ASSETS_DIR}/coherencejob.oracle.com_coherence_small.yaml#CoherenceJob reduced size CRD (.yaml)"
gh release upload v${VERSION} "${RELEASE_ASSETS_DIR}/coherence-dashboards.tar.gz#Coherence Grafana Dashboards (.tar.gz)"
gh release upload v${VERSION} "${RELEASE_ASSETS_DIR}/coherence-operator-bundle.tar.gz#Coherence Operator OLM bundle (.tar.gz)"
gh release upload v${VERSION} "${RELEASE_ASSETS_DIR}/docs.zip#Coherence Operator documentation (.zip)"

