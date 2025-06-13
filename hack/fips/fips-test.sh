#!/usr/bin/env bash
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
set -o errexit

PODS=$(kubectl -n operator-test get pod -l control-plane=coherence -o name)

for POD in ${PODS}
do
  echo "Checking Operator Pod ${POD} is running in FIPS mode"
  kubectl -n operator-test logs ${POD} | grep "Operator is running with FIPS 140 Enabled"
  if [[ $? == 1 ]]
  then
    echo "Failed - did not find FIPS log message for Pod ${POD}"
    exit 1
  fi
done

