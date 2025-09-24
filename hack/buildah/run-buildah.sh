#!/usr/bin/env bash
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

# ---------------------------------------------------------------------------
# This script determines whether Buildah is available locally and if it is
# runs the Coherence image builder script otherwise it starts Buildah inside
# a container and exports the images to the local Docker daemon.
# ---------------------------------------------------------------------------
set -x -e

BASEDIR=$(dirname "$0")

# Ensure the IMAGE_NAME has been set - this is the name of the image to build
if [ "${IMAGE_NAME}" == "" ]
then
  echo "ERROR: No IMAGE_NAME environment variable has been set"
  exit 1
fi

if [ "${IMAGE_NAME_AMD}" == "" ]
then
  IMAGE_NAME_AMD="${IMAGE_NAME}-amd64"
fi

if [ "${IMAGE_NAME_ARM}" == "" ]
then
  IMAGE_NAME_AMD="${IMAGE_NAME}-arm64"
fi

ARGS=

if [ "$1" == "PUSH" ]
then
  SCRIPT_NAME="${BASEDIR}/push-images.sh"
elif [ "$1" == "BUILD" ]
then
  SCRIPT_NAME="${BASEDIR}/build-coherence-images.sh"
elif [ "$1" == "EXEC" ]
then
  SCRIPT_NAME=
  ARGS=-it
fi

if [ "${SCRIPT_NAME}" != "" ]
then
  chmod +x "${SCRIPT_NAME}"
fi

BUILDAH=""
if [ "${LOCAL_BUILDAH}" != "false" ]
then
  BUILDAH=$(which buildah || true)
fi

if [ "${BUILDAH}" != "" ]
then
  echo "Running Buildah locally"
  if [ "${NO_DOCKER_DAEMON}" == "" ]
  then
    export NO_DOCKER_DAEMON=true
  fi
# we must run the script with Buildah unshare
  buildah unshare "${SCRIPT_NAME}"
else
  echo "Buildah not found locally - running in container"
  if [ "${NO_DOCKER_DAEMON}" == "" ]
  then
    NO_DOCKER_DAEMON=false
  fi

  $DOCKER_CMD rm -f buildah || true

  if [ "${BUILDAH_VOLUME}" == "" ]
  then
    export BUILDAH_VOLUME=buildah-containers-volume
  fi
  
  if ! $DOCKER_CMD volume inspect "${BUILDAH_VOLUME}";
  then
    $DOCKER_CMD volume create "${BUILDAH_VOLUME}"
  fi

  if [ "${MY_DOCKER_HOST}" == "" ]
  then
    DOCKER_HOST="${MY_DOCKER_HOST}"
  fi

  if [ "${DOCKER_HOST}" == "" ]
  then
    PDM=$(which podman || true)
    if [ "${USE_PODMAN}" != "false" && "${PDM}" != "" ]
    then
#     we're using Podman
      MY_UID=$(id -u)
      DOCKER_HOST=/run/user/${MY_UID}/podman/podman.sock
    else
#     we're using Docker
      DOCKER_HOST=/var/run/docker.sock
    fi
  fi

  if [ "${PROJECT_ROOT}" == "" ]
  then
    PROJECT_ROOT="${BASEDIR}"
  fi

  $DOCKER_CMD run --rm ${ARGS} -v "${PROJECT_ROOT}:${PROJECT_ROOT}" \
      -v ${DOCKER_HOST}:${DOCKER_HOST}  \
      --privileged --network host \
      -e PROJECT_ROOT="${PROJECT_ROOT}" \
      -e VERSION="${VERSION}" \
      -e REVISION="${REVISION}" \
      -e OCR_DOCKER_USERNAME="${OCR_DOCKER_USERNAME}" \
      -e OCR_DOCKER_USERNAME="${OCR_DOCKER_USERNAME}" \
      -e OCR_DOCKER_SERVER="${OCR_DOCKER_SERVER}" \
      -e REDHAT_REGISTRY_USERNAME="${REDHAT_REGISTRY_USERNAME}" \
      -e REDHAT_REGISTRY_PASSWORD="${REDHAT_REGISTRY_PASSWORD}" \
      -e COHERENCE_VERSION="${COHERENCE_VERSION}" \
      -e IMAGE_NAME="${IMAGE_NAME}" \
      -e AMD_BASE_IMAGE="${AMD_BASE_IMAGE}" \
      -e IMAGE_NAME_AMD="${IMAGE_NAME_AMD}" \
      -e ARM_BASE_IMAGE="${ARM_BASE_IMAGE}" \
      -e IMAGE_NAME_ARM="${IMAGE_NAME_ARM}" \
      -e IMAGE_ARCH="${IMAGE_ARCH}" \
      -e MAIN_CLASS="${MAIN_CLASS}" \
      -e PODMAN_IMPORT="${PODMAN_IMPORT}" \
      -e DOCKER_HOST="${DOCKER_HOST}" \
      -e NO_DOCKER_DAEMON="${NO_DOCKER_DAEMON}" \
      -e IMAGE_NAME_REGISTRY="${IMAGE_NAME_REGISTRY}" \
      -e REGISTRY_USERNAME="${REGISTRY_USERNAME}" \
      -e REGISTRY_PASSWORD="${REGISTRY_PASSWORD}" \
      -e BUILD_ALL_IMAGES="${BUILD_ALL_IMAGES}" \
      -e HTTP_PROXY="${HTTP_PROXY}" -e HTTPS_PROXY="${HTTPS_PROXY}" -e NO_PROXY="${NO_PROXY}" \
      -e http_proxy="${http_proxy}" -e https_proxy="${https_proxy}" -e no_proxy="${no_proxy}" \
      --name buildah \
      quay.io/buildah/stable:v1.37.1 "${SCRIPT_NAME}"
fi

