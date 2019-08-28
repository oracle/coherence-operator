#!/usr/bin/env sh

#!/bin/sh -e -x -u

if [ -z "${WKA}" ] || [ -z "${LISTEN_PORT}" ]; then
  echo "Required WKA environment is not set."
  exit 1
fi

JAVA_OPTS="-Dcoherence.wka=${WKA} \
           -Dcoherence.wka.port=${LISTEN_PORT} \
           -Dcoherence.pof.config=test-pof-config.xml \
           -Dcoherence.cacheconfig=extend-client-cache-config.xml"

if [ -n "${CLUSTER}" ]; then
  JAVA_OPTS="${JAVA_OPTS} -Dcoherence.cluster=${CLUSTER}"
fi

CLASSPATH="/files/conf:/files/lib/*"

CMD="java -cp ${CLASSPATH} ${JAVA_OPTS} ${CLI_OPTS} custom.CloudClient"

echo "Invoking command: ${CMD}"
exec ${CMD}