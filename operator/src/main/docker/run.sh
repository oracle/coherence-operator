#!/bin/bash
# Copyright 2017, 2018, Oracle Corporation and/or its affiliates.  All rights reserved.
# Licensed under the Universal Permissive License v 1.0 as shown at http://oss.oracle.com/licenses/upl.

# set up a logging.properties file that has a FileHandler in it, and have it
# write to /logs/operator.log
LOGGING_CONFIG="/operator/logging.properties"
cp ${LOGGING_CONFIG}.template ${LOGGING_CONFIG}

# if the java logging level has been customized and is a valid value, update logging.properties to match
if [[ ! -z "$JAVA_LOGGING_LEVEL" ]]; then
  SEVERE="SEVERE"
  WARNING="WARNING"
  INFO="INFO"
  CONFIG="CONFIG"
  FINE="FINE"
  FINER="FINER"
  FINEST="FINEST"
  if [ $JAVA_LOGGING_LEVEL != $SEVERE  ] && \
     [ $JAVA_LOGGING_LEVEL != $WARNING ] && \
     [ $JAVA_LOGGING_LEVEL != $INFO    ] && \
     [ $JAVA_LOGGING_LEVEL != $CONFIG  ] && \
     [ $JAVA_LOGGING_LEVEL != $FINE    ] && \
     [ $JAVA_LOGGING_LEVEL != $FINER   ] && \
     [ $JAVA_LOGGING_LEVEL != $FINEST  ]; then
    echo "WARNING: Ignoring invalid JAVA_LOGGING_LEVEL: \"${JAVA_LOGGING_LEVEL}\". Valid values are $SEVERE, $WARNING, $INFO, $CONFIG, $FINE, $FINER and $FINEST."
  else
    sed -i -e "s|\(.*\.level=\).*|\1${JAVA_LOGGING_LEVEL}|g" $LOGGING_CONFIG
  fi
fi

if [ "$ENABLE_FILE_LOGGING" = "true" ]; then
  LOG_HANDLERS="java.util.logging.ConsoleHandler,java.util.logging.FileHandler"
  sed -i -e "s|^\(handlers=\).*$|\1${LOG_HANDLERS}|" $LOGGING_CONFIG
fi

LOGGING="-Djava.util.logging.config.file=${LOGGING_CONFIG}"
mkdir -m 777 -p /logs

# Start operator
java $LOGGING -cp "/operator/lib/*"  com.oracle.coherence.k8s.operator.CoherenceOperator
