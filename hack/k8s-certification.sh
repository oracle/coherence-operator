#!/usr/bin/env bash

#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

# --------------------------------------------------------------------------------
# This script will build the Operator and run the K8s certification test suite.
#
# Pre-requisites:
# * A K8s cluster must be available to run the tests against.
#   The test suite will run using whatever is the current default K8s context
#   on the machine being used.
#
# * Helm v3 must be installed on the machine running the tests as Helm is used
#   by the test suite to install the operator.
#
# * The Operator image and its test images wil be pushed to an image registry.
#   The registry to use should be set using the OPERATOR_IMAGE_REGISTRY environment
#   variable. If the registry requires credentials to push to then the
#   docker login command should already have been executed before this script.
#
# * The MAVEN_USER and MAVEN_PASSWORD environment variables must have been set
#   with credentials to use the https://nexus.synoki.io/repository/maven/ Maven
#   repository to pull down build dependencies.
#
# --------------------------------------------------------------------------------

export OPERATOR_NAMESPACE=coherence

echo "Building Operator"
if ! make build-operator;
then
  exit 1
fi

echo "Building Test imager"
if ! make build-basic-test-image;
then
  exit 1
fi

echo "Building Compatibility Test imager"
if ! make build-compatibility-image;
then
  exit 1
fi

if ! make helm-chart;
then
  exit 1
fi

if [[ "$LOAD_KIND" == "true" ]]; then
  echo "Loading Images to Kind"
  if ! make kind-load;
  then
    exit 1
  fi
fi

echo "Running Certification Tests"
if ! make certification-test;
then
  exit 1
fi

if [[ "$RUN_NET_TEST" != "false" ]]
then
  echo "Running Network Policy Tests"
  if ! make network-policy-test;
  then
    exit 1
  fi

  echo "Running Network Policy Tests - Single Namespace"
  export OPERATOR_NAMESPACE=coherence
  export CLUSTER_NAMESPACE=coherence
  if ! make network-policy-test;
  then
    exit 1
  fi
else
  echo "Skipping Network Policy Tests"
fi

