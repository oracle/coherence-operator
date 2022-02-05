#!/usr/bin/env bash
#
# Copyright (c) 2022, Oracle and/or its affiliates.
#
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

# ---------------------------------------------------------------------------
# This script uses Buildah to build a multi-architecture Coherence image.
# The architectures built are linux/amd64 and linux/arm64.
# The images are pushed to the local Docker daemon unless NO_DAEMON=true.
# ---------------------------------------------------------------------------
set -x

BASEDIR=$(dirname "$0")

# Ensure the ARTIFACT_DIR has been set - this is the root directory of the artifacts to copy to the image
if [ "${ARTIFACT_DIR}" == "" ]
then
  echo "ERROR: No ARTIFACT_DIR environment variable has been set"
  exit 1
fi
# Ensure the IMAGE_NAME has been set - this is the name of the image to build
if [ "${IMAGE_NAME}" == "" ]
then
  echo "ERROR: No IMAGE_NAME environment variable has been set"
  exit 1
fi
# Ensure the AMD_BASE_IMAGE has been set - this is the name of the base image for amd64
if [ "${AMD_BASE_IMAGE}" == "" ]
then
  echo "ERROR: No AMD_BASE_IMAGE environment variable has been set"
  exit 1
fi
# Ensure the ARM_BASE_IMAGE has been set - this is the name of the base image for arm64
if [ "${ARM_BASE_IMAGE}" == "" ]
then
  echo "ERROR: No ARM_BASE_IMAGE environment variable has been set"
  exit 1
fi

# Ensure there is a default architecture - if not set we assume amd64
if [ "${IMAGE_ARCH}" == "" ]
then
  IMAGE_ARCH="amd64"
fi

# we must use docker format to use health checks
export BUILDAH_FORMAT=docker

# Build the entrypoint command line.
ENTRY_POINT="/coherence-operator/utils/runner"

# The command line
CMD="server"

# The health check command line
HEALTH_CMD="ready"

# The image creation date
CREATED=$(date)

# Common image builder function
# param 1: the image architecture, e.g. amd64 or arm64
# param 2: the image o/s e.g. linux
# param 3: the base image
# param 4: the image name
common_image(){
  # Create the container from the base image, setting the architecture and O/S
  buildah from --arch "${1}" --os "${2}" --name "container-${1}" ${3}

  # Add the configuration, entrypoint, etc...
  buildah config --healthcheck-start-period 10s --healthcheck-interval 10s --healthcheck "CMD ${ENTRY_POINT} ${HEALTH_CMD}" container-${1}

  buildah config --arch "${1}" --os "${2}" \
      --entrypoint "[\"${ENTRY_POINT}\"]" --cmd "${CMD}" \
      --annotation "org.opencontainers.image.created=${CREATED}" \
      --annotation "org.opencontainers.image.url=${PROJECT_URL}" \
      --annotation "org.opencontainers.image.version=${VERSION}" \
      --annotation "org.opencontainers.image.source=http://github.com/oracle/coherence-operator" \
      --annotation "org.opencontainers.image.vendor=${PROJECT_VENDOR}" \
      --annotation "org.opencontainers.image.title=${PROJECT_DESCRIPTION} ${VERSION}" \
      --label "org.opencontainers.image.url=${PROJECT_URL}" \
      --label "org.opencontainers.image.version=${VERSION}" \
      --label "org.opencontainers.image.source=http://github.com/oracle/coherence-operator" \
      --label "org.opencontainers.image.vendor=${PROJECT_VENDOR}" \
      --label "org.opencontainers.image.title=Oracle Coherence ${VERSION}" \
      "container-${1}"

  # Copy files into the container
  buildah copy "container-${1}" "${ARTIFACT_DIR}/target/docker/linux/${1}/runner" /coherence-operator/utils/runner
  buildah copy "container-${1}" "${ARTIFACT_DIR}/target/docker/lib"               /app/libs

  echo
  buildah inspect container-${1}
  echo

  # Commit the container to an image
  buildah commit "container-${1}" "coherence-operator:${1}"

  # Export the image to the Docker daemon unless NO_DAEMON is true
  if [ "${NO_DAEMON}" != "true" ]
  then
    buildah push -f v2s2 "coherence-operator:${1}" "docker-daemon:${4}"
    echo "Pushed ${2}/${1} image ${4} to Docker daemon"
  fi
}

buildah version

if [ "${DOCKER_HUB_USERNAME}" != "" ] && [ "${DOCKER_HUB_PASSWORD}" != "" ]
then
  buildah login -u "${DOCKER_HUB_USERNAME}" -p "${DOCKER_HUB_PASSWORD}" "docker.io"
fi

if [ "${DOCKER_REGISTRY}" != "" ] && [ "${DOCKER_USERNAME}" != "" ] && [ "${DOCKER_PASSWORD}" != "" ]
then
  buildah login -u "${DOCKER_USERNAME}" -p "${DOCKER_PASSWORD}" "${DOCKER_REGISTRY}"
fi

# Build the amd64 image
common_image amd64 linux "${AMD_BASE_IMAGE}" "${IMAGE_NAME}-amd64"

# Build the arm64 image
common_image arm64 linux "${ARM_BASE_IMAGE}" "${IMAGE_NAME}-arm64"

# Push the relevant image to the docker daemon base on the build machine's o/s architecture
if [ "${NO_DAEMON}" != "true" ]
then
  buildah push -f v2s2 "coherence-operator:${IMAGE_ARCH}" "docker-daemon:${IMAGE_NAME}"
  echo "Pushed linux/${IMAGE_ARCH} image ${IMAGE_NAME} to Docker daemon"
fi

# Clean-up
buildah rm container-amd64
buildah rm container-arm64


