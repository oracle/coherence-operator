#!/bin/sh -e

if [ -f "${COHERENCE_HOME}/lib/coherence.jar" ]; then
  echo "Copying ${COHERENCE_HOME}/lib/coherence.jar"
  cp "${COHERENCE_HOME}/lib/coherence.jar" /app/libs/coherence.jar
fi
