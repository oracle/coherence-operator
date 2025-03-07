#
# Copyright (c) 2021, 2023, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
FROM container-registry.oracle.com/middleware/coherence-operator:delve

ARG target
ARG version
ARG coherence_image
ARG operator_image

LABEL "com.oracle.coherence.application"="operator"
LABEL "com.oracle.coherence.version"="$version"

ENV COHERENCE_IMAGE=$coherence_image \
    OPERATOR_IMAGE=$operator_image

WORKDIR /

COPY bin/linux/$target/runner-debug                                    /files/runner
COPY java/coherence-operator/target/docker/lib/*.jar                   /files/lib/
COPY java/coherence-operator/target/docker/logging/logging.properties  /files/logging/logging.properties

ENTRYPOINT ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/files/runner", "--"]
