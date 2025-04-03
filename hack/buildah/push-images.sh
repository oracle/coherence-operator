#!/usr/bin/env bash
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
set -x -e

buildah version

if [ "${OPERATOR_IMAGE_REGISTRY}" != "" ] && [ "${REGISTRY_USERNAME}" != "" ] && [ "${REGISTRY_PASSWORD}" != "" ]
then
  buildah login -u "${REGISTRY_USERNAME}" -p "${REGISTRY_PASSWORD}" "${OPERATOR_IMAGE_REGISTRY}"
fi

if [ "${NO_DOCKER_DAEMON}" != "true" ]
then
  buildah pull "docker-daemon:${OPERATOR_IMAGE_AMD}"
  buildah pull "docker-daemon:${OPERATOR_IMAGE_ARM}"
fi

DESCR='The Oracle Coherence Kubernetes Operator allows full lifecycle management of Oracle Coherence workloads in Kubernetes.'

buildah rmi "${OPERATOR_IMAGE}" || true
buildah manifest rm "${OPERATOR_IMAGE}" || true
buildah manifest create \
    --annotation org.opencontainers.image.source=https://github.com/oracle/coherence-operator \
    --annotation org.opencontainers.image.description="\"${DESCR}\"" \
    --annotation org.opencontainers.image.licenses="UPL-1.0" \
    --annotation org.opencontainers.image.version="${VERSION}" \
    --annotation org.opencontainers.image.revision="${REVISION}" \
    "${OPERATOR_IMAGE}"
buildah manifest add --arch amd64 --os linux \
    --annotation org.opencontainers.image.source=https://github.com/oracle/coherence-operator \
    --annotation org.opencontainers.image.description="\"${DESCR}\"" \
    --annotation org.opencontainers.image.licenses="UPL-1.0" \
    --annotation org.opencontainers.image.version="${VERSION}" \
    --annotation org.opencontainers.image.revision="${REVISION}" \
    "${OPERATOR_IMAGE}" "${OPERATOR_IMAGE_AMD}"
buildah manifest add --arch arm64 --os linux \
    --annotation org.opencontainers.image.source=https://github.com/oracle/coherence-operator \
    --annotation org.opencontainers.image.description="\"${DESCR}\"" \
    --annotation "org.opencontainers.image.licenses"="UPL-1.0" \
    --annotation org.opencontainers.image.version="${VERSION}" \
    --annotation org.opencontainers.image.revision="${REVISION}" \
    "${OPERATOR_IMAGE}" "${OPERATOR_IMAGE_ARM}"
buildah manifest inspect "${OPERATOR_IMAGE}"

buildah manifest push --all -f v2s2 "${OPERATOR_IMAGE}" "docker://${OPERATOR_IMAGE}"
