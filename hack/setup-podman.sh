#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

#!/usr/bin/env bash

# --------------------------------------------------------------------
# This script exports various environment variables so that
# Make targets will be executed using Podman and OpenShift.
# --------------------------------------------------------------------

export DOCKER_CMD=podman
export PUSH_ARGS=--tls-verify=false
export REGISTRY=$(oc get route/default-route -n openshift-image-registry -o=jsonpath='{.spec.host}')
export OPERATOR_IMAGE_REGISTRY=${REGISTRY}/oracle
export PREFLIGHT_REGISTRY_CRED=$(echo -n bogus:$(oc whoami -t) | base64)

podman login -u bogus -p $(oc whoami -t) --tls-verify=false ${REGISTRY}
podman login -u bogus -p $(oc whoami -t) --tls-verify=false ${OPERATOR_IMAGE_REGISTRY}

#oc new-project oracle || true
#oc -n oracle create is coherence-operator || true
# Allow anyone to pull oracle images
#oc adm policy add-role-to-group system:image-puller system:authenticated --namespace=oracle