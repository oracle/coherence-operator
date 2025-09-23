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

if [ -z "${REDHAT_EXAMPLE_IMAGE}" ]; then
    echo "Error: REDHAT_EXAMPLE_IMAGE should be set to the name of the Red Hat Coherence image"
    exit 1
fi

if [ "${SUBMIT_RESULTS}" = "true" ]; then
  if [ -z "${OPENSHIFT_API_KEY:-}" ]; then
    echo "Error: SUBMIT_RESULTS is set to 'true' but OPENSHIFT_API_KEY is not set"
    exit 1
  fi
  if [ -z "${OPENSHIFT_COHERENCE_COMPONENT_ID:-}" ]; then
    OPENSHIFT_COHERENCE_COMPONENT_ID="68d28054a49e977fe49f4234"
  fi
  EXTRA_ARGS="--pyxis-api-token=${OPENSHIFT_API_KEY} --certification-component-id=${OPENSHIFT_COHERENCE_COMPONENT_ID}"
else
  EXTRA_ARGS=""
fi

echo "Running preflight on ${REDHAT_EXAMPLE_IMAGE}"

PREFLIGHT_LOG="${OUTPUT_DIR}/preflight.log"
preflight check container --submit="${SUBMIT_RESULTS}" --logfile="${PREFLIGHT_LOG}" ${EXTRA_ARGS} ${REDHAT_EXAMPLE_IMAGE}
