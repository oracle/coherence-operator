FROM oraclelinux:7-slim AS builder

RUN echo "coherence:x:$(id -u):$(id -g):coherence user:/opt/coherence:/sbin/nologin" >> /etc/passwd

FROM scratch

ENV OPERATOR=/usr/local/bin/coherence-operator \
    USER_UID=1001 \
    USER_NAME=coherence \
    HOME=/opt/coherence

LABEL "com.oracle.coherence.application"="operator"
LABEL "com.oracle.coherence.version"="${VERSION_FULL}"

COPY --from=builder /etc/passwd            /etc/passwd
COPY watches.yaml                          ${HOME}/watches.yaml
COPY build/_output/helm-charts/coherence   ${HOME}/helm-charts/coherence
COPY deploy/crds/*_crd.yaml                ${HOME}/crds/
COPY build/_output/bin/operator            ${OPERATOR}

ENTRYPOINT ["/usr/local/bin/coherence-operator", "run", "helm", "--watches-file=/opt/coherence/watches.yaml", "--crd-files=/opt/coherence/crds"]

USER ${USER_UID}
