#!/usr/bin/env bash
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

ROOT_DIR=$(pwd)
BUILD_DIR="${ROOT_DIR}/build"
OUTPUT_DIR="${BUILD_DIR}/_output"

mkdir -p "${OUTPUT_DIR}" || true

if [ -z "${SUBMIT_RESULTS:-}" ]; then
  SUBMIT_RESULTS=false
fi

if [ "${SUBMIT_RESULTS}" = "true" ]; then
  if [ -z "${OPENSHIFT_API_KEY:-}" ]; then
    echo "Error: SUBMIT_RESULTS is set to 'true' but OPENSHIFT_API_KEY is not set"
    exit 1
  fi
  if [ -z "${OPENSHIFT_IMAGE_COMPONENT_ID:-}" ]; then
    OPENSHIFT_IMAGE_COMPONENT_ID="67bdf00eb9f79dcdb25aa8e2"
  fi
  EXTRA_ARGS="--pyxis-api-token=${OPENSHIFT_API_KEY} --certification-component-id=${OPENSHIFT_IMAGE_COMPONENT_ID}"
else
  EXTRA_ARGS=""
fi

if [ "${USE_LATEST_OPERATOR_RELEASE}" = "true" ]; then
  echo "Run preflight on latest release"
# Find the latest release of the Coherence Operator on GitHub
  LATEST_RELEASE=$(gh release list --repo oracle/coherence-operator --json name,isLatest --jq '.[] | select(.isLatest)|.name')
# Strip the v from the front of the release to give the Operator version
  OPERATOR_VERSION=${LATEST_RELEASE#"v"}
  echo "Latest Operator version is ${OPERATOR_VERSION}"
# Check the latest release image exists on OCR
  COHERENCE_OPERATOR_IMAGE="container-registry.oracle.com/middleware/coherence-operator:${OPERATOR_VERSION}"
  echo "Checking Oracle Container Registry for image ${COHERENCE_OPERATOR_IMAGE}"
  docker manifest inspect "${COHERENCE_OPERATOR_IMAGE}" > /dev/null
  if [ $? -ne 0 ]; then
    echo "ERROR: Image ${COHERENCE_OPERATOR_IMAGE} does not exist on OCR."
    exit 1
  fi
else
  echo "Run preflight on latest build"
# We are just testing a local build, so use the current version
  OPERATOR_VERSION=$(cat "${BUILD_DIR}/_output/version.txt")
# We will not be submitting results
  SUBMIT_RESULTS=false
  COHERENCE_OPERATOR_IMAGE="${REGISTRY_HOST}/${REGISTRY_NAMESPACE}/coherence-operator:${OPERATOR_VERSION}"
fi

echo "Running preflight on ${COHERENCE_OPERATOR_IMAGE}"

PREFLIGHT_LOG="${OUTPUT_DIR}/preflight.log"
preflight check container --submit="${SUBMIT_RESULTS}" --logfile="${PREFLIGHT_LOG}" ${EXTRA_ARGS} ${COHERENCE_OPERATOR_IMAGE}
