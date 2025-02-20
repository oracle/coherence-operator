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
ARG release

LABEL "com.oracle.coherence.application"="operator"
LABEL "com.oracle.coherence.version"="$version"
LABEL "org.opencontainers.image.revision"="$release"
LABEL "org.opencontainers.image.description"="The Oracle Coherece Kubernetes Operator image ($target)"
LABEL "org.opencontainers.image.source"="https://github.com/oracle/coherence-operator"
LABEL "org.opencontainers.image.authors"="To contact the authors use this link https://github.com/oracle/coherence-operator/discussions"
LABEL "org.opencontainers.image.licenses"="UPL-1.0"
LABEL "org.opencontainers.image.description"="The Oracle Coherece Kubernetes Operator allows full lifecycle management of Oracle Coherence workloads in Kubernetes."

LABEL "name"="Oracle Coherence Kubernetes Operator"
LABEL "vendor"="Oracle"
LABEL "version"="$version"
LABEL "release"="$release"
LABEL "maintainer"="Oracle Coherence Engieering Team"
LABEL "summary"="A Kubernetes Operator for managing Oracle Coherence clusters"
LABEL "description"="The Oracle Coherece Kubernetes Operator allows full lifecycle management of Oracle Coherence workloads in Kubernetes."

ENV COHERENCE_IMAGE=$coherence_image \
    OPERATOR_IMAGE=$operator_image

COPY LICENSE.txt                                                       /licenses/LICENSE.txt
COPY bin/linux/$target/*                                               /files/
COPY java/coherence-operator/target/docker/lib/*.jar                   /files/lib/
COPY java/coherence-operator/target/docker/logging/logging.properties  /files/logging/logging.properties

USER 1000

ENTRYPOINT ["/files/runner"]
CMD ["-h"]
