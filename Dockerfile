#
# Copyright (c) 2019, 2023, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
ARG BASE_IMAGE=scratch
FROM $BASE_IMAGE

ARG target
ARG version
ARG coherence_image
ARG operator_image

LABEL "com.oracle.coherence.application"="operator"
LABEL "com.oracle.coherence.version"="$version"

ENV COHERENCE_IMAGE=$coherence_image \
    OPERATOR_IMAGE=$operator_image

COPY --chmod=0755 bin/linux/$target/*                                  /files/
COPY java/coherence-operator/target/docker/lib/*.jar                   /files/lib/
COPY java/coherence-operator/target/docker/logging/logging.properties  /files/logging/logging.properties

ENTRYPOINT ["/files/runner"]
CMD ["-h"]
