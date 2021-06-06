#
# Copyright (c) 2019, 2020 Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
FROM scratch

ARG target
ARG version
ARG coherence_image
ARG utils_image

LABEL "com.oracle.coherence.application"="operator"
LABEL "com.oracle.coherence.version"="$version"

ENV COHERENCE_IMAGE=$coherence_image \
    UTILS_IMAGE=$utils_image

COPY bin/linux/$target/manager  .
ENTRYPOINT ["/manager"]
