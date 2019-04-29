#
# Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
#!/bin/bash

exitCode=0

# Run is a wrapper around the execution of functions. It captures non-zero exit
# codes and remembers an error happened. This enables running all the linters
# and capturing if any of them failed.
run() {
  $@
  local ret=$?
  if [ $ret -ne 0 ]; then
    exitCode=1
  fi
}

# Lint the Chart.yaml and values.yaml files for Helm
yamllinter() {
  printf "\nLinting the Chart.yaml and values.yaml files at ${1}\n"

  # If a Chart.yaml file is present lint it. Otherwise report an error
  # because one should exist
  if [ -e $1/Chart.yaml ]; then
    run yamllint -c ${curDir}/lintconf.yaml $1/Chart.yaml
  else
    echo "Error $1/Chart.yaml file is missing"
    exitCode=1
  fi

  # If a values.yaml file is present lint it. Otherwise report an error
  # because one should exist
  if [ -e $1/values.yaml ]; then
    run yamllint -c ${curDir}/lintconf.yaml $1/values.yaml
  else
    echo "Error $1/values.yaml file is missing"
    exitCode=1
  fi
}

# include the semvercompare function
# curDir="$(dirname "$0")"
source "operator/src/main/helm/scripts/semvercompare.sh"

directory=$1

printf "\nRunning helm dep build on the chart at ${directory}\n"
run helm dep build ${directory}

printf "\nRunning helm lint on the chart at ${directory}\n"
run helm lint ${directory}

yamllinter ${directory}

# Skip version check
#semvercompare ${directory}

# Check for the existence of the NOTES.txt file. This is required for charts
# in this repo.
if [ ! -f ${directory}/templates/NOTES.txt ]; then
  echo "Error NOTES.txt template not found. Please create one."
  echo "For more information see https://docs.helm.sh/developing_charts/#chart-license-readme-and-notes"
  exitCode=1
fi

exit $exitCode
