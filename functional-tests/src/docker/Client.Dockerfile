#
# Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
FROM openjdk:11-oracle

RUN mkdir -p /files/libs
RUN mkdir -p /files/conf

COPY run.sh     /run.sh
COPY lib/*      /files/lib/

RUN chmod +x /run.sh

ENTRYPOINT ["/run.sh"]
