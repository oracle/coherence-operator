FROM scratch

ARG version
ARG coherence_image
ARG utils_image

LABEL "com.oracle.coherence.application"="operator"
LABEL "com.oracle.coherence.version"="$version"

ENV OPERATOR=/usr/local/bin/coherence-operator \
    HOME=/opt/coherence \
    HELM_COHERENCE_IMAGE=$coherence_image \
    UTILS_IMAGE=$utils_image \
    VERSION_FULL=$version

COPY deploy/crds/*_crd.yaml                ${HOME}/crds/
COPY image/scripts                         ${HOME}/operator-scripts
COPY build/_output/bin/operator            ${OPERATOR}

ENTRYPOINT ["/usr/local/bin/coherence-operator"]
CMD ["run", "--crd-files=/opt/coherence/crds", "--scripts-dir=/opt/coherence/operator-scripts"]
