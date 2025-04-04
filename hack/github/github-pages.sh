#!/usr/bin/env bash
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

set -o -v errexit

ROOT_DIR=$(pwd)
GITHUB_REPO=$(git config --get remote.origin.url)
GITHUB_REPO=https://github.com/oracle/coherence-operator.git
BUILD_OUTPUT=${ROOT_DIR}/build/_output
BUILD_GH_PAGES="${ROOT_DIR}/gh-pages"
VERSION=$(cat "${BUILD_OUTPUT}/version.txt")

@echo "ROOT_DIR="${ROOT_DIR}""
@echo "GITHUB_REPO=${GITHUB_REPO}"
@echo "VERSION=${VERSION}"

rm -rf ${BUILD_GH_PAGES}
git clone -b gh-pages --single-branch ${GITHUB_REPO} "${ROOT_DIR}/gh-pages"
cd ${BUILD_GH_PAGES}

mkdir -p ${BUILD_GH_PAGES}/dashboards || true
rm -rf ${BUILD_GH_PAGES}/dashboards/${VERSION} || true
cp -R ${BUILD_OUTPUT}/dashboards/${VERSION} ${BUILD_GH_PAGES}/dashboards/
git add dashboards/${VERSION}/*
rm -rf ${BUILD_GH_PAGES}/dashboards/latest || true
cp -R ${BUILD_GH_PAGES}/dashboards/${VERSION} ${BUILD_GH_PAGES}/dashboards/latest
git add -A dashboards/latest/*

mkdir ${BUILD_GH_PAGES}/docs/${VERSION} || true
rm -rf ${BUILD_GH_PAGES}/docs/${VERSION}/ || true
cp -R ${BUILD_OUTPUT}/docs ${BUILD_GH_PAGES}/docs/${VERSION}/
rm -rf ${BUILD_GH_PAGES}/docs/latest
cp -R ${BUILD_GH_PAGES}/docs/${VERSION} ${BUILD_GH_PAGES}/docs/latest
git add -A docs/*

mkdir -p ${BUILD_GH_PAGES}/charts || true
cp ${BUILD_OUTPUT}/helm-charts/coherence-operator-${VERSION}.tgz ${BUILD_GH_PAGES}/charts/
helm repo index charts --url https://oracle.github.io/coherence-operator/charts
git add charts/coherence-operator-${VERSION}.tgz
git add charts/index.yaml

git clean -d -f

git status

git commit -m "Release Coherence Operator version: ${VERSION}"
git log -1
git push origin gh-pages

cd ${ROOT_DIR}
