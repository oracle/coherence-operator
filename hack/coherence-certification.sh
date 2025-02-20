#!/usr/bin/env bash

# --------------------------------------------------------------------------------
# This script will build the Operator and run the Coherence certification test suite.
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

echo "Building Operator"
make clean
make build-operator-images
if [[ $? != 0 ]]; then
  exit 1
fi

echo "Building Coherence Compatibility Image"
make build-compatibility-image
if [[ $? != 0 ]]; then
  exit 1
fi

echo "Pushing Images"
make push-release-images
if [[ $? != 0 ]]; then
  exit 1
fi
make push-compatibility-image
if [[ $? != 0 ]]; then
  exit 1
fi

if [[ "$LOAD_KIND" == "true" ]]; then
  echo "Loading Images to Kind"
  make kind-load
  make kind-load-compatibility
fi

echo "Running Coherence Compatibility Tests"
make coherence-compatibility-test
if [[ $? != 0 ]]; then
  exit 1
fi
