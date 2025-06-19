#!/usr/bin/env bash
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

set -o -v errexit

if [ "${WORKSPACE}" == "" ]
then
  WORKSPACE=$(pwd)
fi

if [ "${OPERATOR_VERSION}" == "" ]
then
  OPERATOR_VERSION=$(make version)
fi

GITHUB_REPO=https://github.com/oracle/coherence-operator.git
BUILD_OUTPUT=${WORKSPACE}/build/_output
BUILD_GH_PAGES="${WORKSPACE}/gh-pages"

@echo "WORKSPACE="${WORKSPACE}""
@echo "GITHUB_REPO=${GITHUB_REPO}"
@echo "OPERATOR_VERSION=${OPERATOR_VERSION}"

rm -rf ${BUILD_GH_PAGES}
git clone -b gh-pages --single-branch ${GITHUB_REPO} "${WORKSPACE}/gh-pages"
cd ${BUILD_GH_PAGES}

GIT_ORIGIN=$(git config remote.origin.url)
GIT_URL=$(echo ${GIT_ORIGIN} | sed -e s#://#://${GITHUB_USERNAME}:${GITHUB_TOKEN}@#)
git remote set-url origin "${GIT_URL}"

mkdir -p ${BUILD_GH_PAGES}/dashboards || true
rm -rf ${BUILD_GH_PAGES}/dashboards/${OPERATOR_VERSION} || true
cp -R ${BUILD_OUTPUT}/dashboards/${OPERATOR_VERSION} ${BUILD_GH_PAGES}/dashboards/
git add dashboards/${OPERATOR_VERSION}/*
rm -rf ${BUILD_GH_PAGES}/dashboards/latest || true
cp -R ${BUILD_GH_PAGES}/dashboards/${OPERATOR_VERSION} ${BUILD_GH_PAGES}/dashboards/latest
git add -A dashboards/latest/*

mkdir -p ${BUILD_GH_PAGES}/charts || true
cp ${BUILD_OUTPUT}/helm-charts/coherence-operator-${OPERATOR_VERSION}.tgz ${BUILD_GH_PAGES}/charts/
ls -al ${BUILD_GH_PAGES}/charts
helm repo index ${BUILD_GH_PAGES}/charts --url https://oracle.github.io/coherence-operator/charts
git add ${BUILD_GH_PAGES}/charts/coherence-operator-${OPERATOR_VERSION}.tgz
git add ${BUILD_GH_PAGES}/charts/index.yaml

git clean -d -f

git status

git commit -m "Release Coherence Operator version: ${OPERATOR_VERSION}"
git log -1
git push origin gh-pages

cd ${WORKSPACE}
