FROM oraclelinux:7-slim

ENV OPERATOR=/usr/local/bin/coherence-operator \
    USER_UID=1001 \
    USER_NAME=coherence \
    HOME=/opt/coherence

COPY watches.yaml ${HOME}/watches.yaml

# install operator binaries
COPY build/_output/bin/operator ${OPERATOR}


COPY bin /usr/local/bin
RUN  /usr/local/bin/user_setup
COPY build/_output/helm-charts/coherence ${HOME}/helm-charts/coherence
COPY deploy/crds/*_crd.yaml ${HOME}/crds/

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
