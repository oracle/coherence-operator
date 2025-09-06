#!/usr/bin/env sh
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

SCRIPT_PATH="${0}"
SCRIPT_DIR=$(dirname -- "${SCRIPT_PATH}")
TEKTON_DIR=$(realpath -- "${SCRIPT_DIR}")

if [[ "${KUBECTL_CMD}" == "" ]]; then
  KUBECTL_CMD="kubectl"
fi

if [[ "${PIPELINE_NAMESPACE}" != "" ]]; then
  NS=${PIPELINE_NAMESPACE}
else
  NS=default
fi

echo "Installing Tekton resources ${NS}"
tkn task delete --namespace ${NS} --force git-clone || true
tkn hub install task --namespace ${NS} git-clone
tkn task delete --namespace ${NS} --force git-cli || true
tkn hub install task --namespace ${NS} git-cli


# Install Operator Tasks
${KUBECTL_CMD} --namespace ${NS} apply --filename ${TEKTON_DIR}/task-make.yaml
${KUBECTL_CMD} --namespace ${NS} apply --filename ${TEKTON_DIR}/task-setup-env.yaml
${KUBECTL_CMD} --namespace ${NS} apply --filename ${TEKTON_DIR}/task-buildah.yaml
${KUBECTL_CMD} --namespace ${NS} apply --filename ${TEKTON_DIR}/task-check-image.yaml
${KUBECTL_CMD} --namespace ${NS} apply --filename ${TEKTON_DIR}/task-oci-cli.yaml

# Install Operator Pipelines
${KUBECTL_CMD} --namespace ${NS} apply --filename ${TEKTON_DIR}/pipeline-operator-ci.yaml

# Install Operator test configmap
${KUBECTL_CMD} --namespace ${NS} apply --filename ${TEKTON_DIR}/os-cert-config.yaml
