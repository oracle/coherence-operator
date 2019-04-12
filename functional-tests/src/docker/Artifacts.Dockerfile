#
# Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
FROM oraclelinux:7-slim

ARG  VERSION

COPY lib/coherence-operator-tests.jar  /files/lib/coherence-operator-tests.jar
COPY custom-logging.properties             /files/conf/custom-logging.properties

RUN  mkdir -p /files/conf && echo $VERSION > /files/conf/version.txt
