FROM scratch

ARG version
ARG coherence_image
ARG utils_image

LABEL "com.oracle.coherence.application"="operator"
LABEL "com.oracle.coherence.version"="$version"

ENV COHERENCE_IMAGE=$coherence_image \
    UTILS_IMAGE=$utils_image

COPY build/_output/manager  .
ENTRYPOINT ["/manager"]
