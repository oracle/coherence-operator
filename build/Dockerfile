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

COPY watches.yaml                          ${HOME}/watches.yaml
COPY build/_output/helm-charts/coherence   ${HOME}/helm-charts/coherence
COPY deploy/crds/*_crd.yaml                ${HOME}/crds/
COPY build/_output/bin/operator            ${OPERATOR}

ENTRYPOINT ["/usr/local/bin/coherence-operator", "run", "helm", "--watches-file=/opt/coherence/watches.yaml", "--crd-files=/opt/coherence/crds"]
