#!/usr/bin/env bash
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
set -x -e -v

buildah version

if [ "${IMAGE_NAME_REGISTRY}" != "" ] && [ "${REGISTRY_USERNAME}" != "" ] && [ "${REGISTRY_PASSWORD}" != "" ]
then
  buildah login -u "${REGISTRY_USERNAME}" -p "${REGISTRY_PASSWORD}" "${IMAGE_NAME_REGISTRY}"
fi

if [ "${IMAGE_NAME_AMD}" == "" ]
then
  IMAGE_NAME_AMD="${IMAGE_NAME}-amd64"
fi

if [ "${IMAGE_NAME_ARM}" == "" ]
then
  IMAGE_NAME_ARM="${IMAGE_NAME}-arm64"
fi

if [ "${NO_DOCKER_DAEMON}" != "true" ]
then
  buildah pull "docker-daemon:${IMAGE_NAME_AMD}"
  buildah pull "docker-daemon:${IMAGE_NAME_ARM}"
fi

buildah rmi "${IMAGE_NAME}" || true
buildah manifest rm "${IMAGE_NAME}" || true

buildah manifest create \
    --annotation org.opencontainers.image.source=https://github.com/oracle/coherence-operator \
    --annotation org.opencontainers.image.licenses="UPL-1.0" \
    --annotation org.opencontainers.image.version="${VERSION}" \
    --annotation org.opencontainers.image.revision="${REVISION}" \
    "${IMAGE_NAME}"

buildah manifest add --arch amd64 --os linux \
    --annotation org.opencontainers.image.source=https://github.com/oracle/coherence-operator \
    --annotation org.opencontainers.image.licenses="UPL-1.0" \
    --annotation org.opencontainers.image.version="${VERSION}" \
    --annotation org.opencontainers.image.revision="${REVISION}" \
    "${IMAGE_NAME}" "${IMAGE_NAME_AMD}"

buildah manifest add --arch arm64 --os linux \
    --annotation org.opencontainers.image.source=https://github.com/oracle/coherence-operator \
    --annotation "org.opencontainers.image.licenses"="UPL-1.0" \
    --annotation org.opencontainers.image.version="${VERSION}" \
    --annotation org.opencontainers.image.revision="${REVISION}" \
    "${IMAGE_NAME}" "${IMAGE_NAME_ARM}"

buildah manifest inspect "${IMAGE_NAME}"

buildah manifest push --all -f v2s2 "${IMAGE_NAME}" "docker://${IMAGE_NAME}"
