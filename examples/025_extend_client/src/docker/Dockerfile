# Copyright (c) 2022, Oracle and/or its affiliates.
#
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
FROM gcr.io/distroless/java11-debian11

# Set the entry point to be the Java command to run
ENTRYPOINT ["java", "-cp", "/app/classes:/app/libs/*", "com.oracle.coherence.examples.extend.Main"]

# Set the Coherence log level to debug logging
ENV COHERENCE_LOG_LEVEL=9
# Effectively disabled multicast cluster discovery, which does not work in containers
ENV COHERENCE_TTL=0

# Copy all the application files into the /app directory in the image
# This is the default structure supported by the Coherence Operator
COPY app app
