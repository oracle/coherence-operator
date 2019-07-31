FROM oraclelinux:7-slim

ENV OPERATOR=/usr/local/bin/helm-operator \
    USER_UID=1001 \
    USER_NAME=helm \
    HOME=/opt/helm

COPY watches.yaml ${HOME}/watches.yaml

# install operator binaries
COPY build/_output/bin/operator ${OPERATOR}


COPY bin /usr/local/bin
RUN  /usr/local/bin/user_setup
COPY build/_output/helm-charts/ ${HOME}/helm-charts/

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}