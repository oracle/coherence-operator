#!/usr/bin/env bash
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
set -e -v

ROOT_DIR=$(pwd)
BUILD_DIR=${ROOT_DIR}/build
PIPELINES_DIR=${BUILD_DIR}/operator-pipelines

if [ "${OPENSHIFT_VERSION}" = "" ]; then
  OPENSHIFT_VERSION=v4.19
fi

# RedHat OpenShift certification tests
# See: https://github.com/redhat-openshift-ecosystem/certification-releases/blob/main/4.9/ga/operator-cert-workflow.md

if [ -z "${GITHUB_TOKEN:-}" ]; then
  echo "Error: GITHUB_TOKEN is not set"
  exit 1
fi

# Step A - Get Project ID
if [ -z "${OPENSHIFT_PROJECT_ID:-}" ]; then
  echo "Error: OPENSHIFT_PROJECT_ID is not set"
  exit 1
fi

# Step B - Get API Key
if [ -z "${OPENSHIFT_API_KEY:-}" ]; then
  echo "Error: OPENSHIFT_API_KEY is not set"
  exit 1
fi

# Step C - Install Pipeline

# Step C.1 - Install OpenShift Pipelines Operator
#   This is already done in the OpenShift test cluster

# Step C.2 - Configure the OpenShift CLI tool (oc)

if [ -z "${KUBECONFIG:-}" ]; then
  export KUBECONFIG=${HOME}/.kube/config
  echo "KUBECONFIG is unset, using default ${KUBECONFIG}"
fi

# Step C.3 - Create an OpenShift Project (namespace) to work in

if [ -z "${PROJECT_NAME:-}" ]; then
  export PROJECT_NAME=operator-cert
fi
echo "Running OpenShift certification in project ${PROJECT_NAME}"

if [ "${RESET_PROJECT}" != "false" ]; then
  echo "Resetting project ${PROJECT_NAME}"
  # Switch to the default project
  oc project default
  # If the project we want exists, then delete it
  if oc get namespace ${PROJECT_NAME} >/dev/null 2>&1; then
    oc delete ns ${PROJECT_NAME} --force=true
  fi

  # Create a new project for testing
  oc adm new-project ${PROJECT_NAME}
else
    if oc get namespace ${PROJECT_NAME} >/dev/null 2>&1; then
      echo "Using existing project ${PROJECT_NAME}"
    else
      echo "Creating project ${PROJECT_NAME}"
      # Create a new project for testing
      oc adm new-project ${PROJECT_NAME}
    fi
fi

oc project ${PROJECT_NAME}

# Step C.4 - Add the Kubeconfig secret
if oc get secret kubeconfig >/dev/null 2>&1; then
  oc delete secret kubeconfig
fi
oc create secret generic kubeconfig --from-file=kubeconfig=${KUBECONFIG}

# Step C.5 - Import Red Hat Catalogs

oc import-image certified-operator-index:${OPENSHIFT_VERSION} \
  --from=registry.redhat.io/redhat/certified-operator-index:${OPENSHIFT_VERSION} \
  --reference-policy local \
  --scheduled \
  --confirm

# Step C.6 - Install the Certification Pipeline and dependencies into the cluster
if [ ! -e "${PIPELINES_DIR}" ]; then
  cd "${BUILD_DIR}"
  git clone --quiet https://github.com/redhat-openshift-ecosystem/operator-pipelines
fi

cd "${PIPELINES_DIR}"

GIT_ORIGIN=$(git config remote.origin.url)
GITHUB_PUSH_UPSTREAM="${GIT_ORIGIN}"
if [ ! -z "${GITHUB_TOKEN:-}" ]; then
  GITHUB_PUSH_UPSTREAM=$(echo "${GIT_ORIGIN}" | sed -e s#://#://${GITHUB_USERNAME}:${GITHUB_TOKEN}@#)
  git config user.name "${GITHUB_USERNAME}"
  if [ ! -z "${GITHUB_USER_EMAIL:-}" ]; then
    git config user.email "${GITHUB_USER_EMAIL}"
  fi
fi

oc apply -R -f ansible/roles/operator-pipeline/templates/openshift/pipelines/operator-ci-pipeline.yml
oc apply -R -f ansible/roles/operator-pipeline/templates/openshift/tasks
# Create a new SCC
oc apply -f ansible/roles/operator-pipeline/templates/openshift/openshift-pipelines-custom-scc.yml
# Add SCC to a pipeline service account
oc adm policy add-scc-to-user pipelines-custom-scc -z pipeline

# Step C.7 - Configuration Steps for Submitting Results

# Add a GitHub API Token for the repo where the PR will be created
if oc get secret github-api-token >/dev/null 2>&1; then
  oc delete secret github-api-token
fi
oc create secret generic github-api-token --from-literal GITHUB_TOKEN="${GITHUB_TOKEN}"

# Add Red Hat Container API access key
if oc get secret pyxis-api-secret >/dev/null 2>&1; then
  oc delete secret pyxis-api-secret
fi
oc create secret generic pyxis-api-secret --from-literal pyxis_api_key="${OPENSHIFT_API_KEY}"

if [ "${GITHUB_SSL_KEY_SECRET}" != "" ]; then
  if [ -e "${GITHUB_SSL_KEY_SECRET}" ]
  then
    oc apply -f "${GITHUB_SSL_KEY_SECRET}"
  fi
fi

if oc get secret registry-dockerconfig-secret >/dev/null 2>&1; then
  oc delete secret registry-dockerconfig-secret
fi
oc create secret docker-registry registry-dockerconfig-secret \
    --docker-server="${REGISTRY_HOST}" \
    --docker-username="${REGISTRY_USERNAME}" \
    --docker-password="${REGISTRY_PASSWORD}" \
    --docker-email=someone@oracle.com

# Step E - Add Operator Bundle
# Checkout the certified-operators fork
OPERATOR_VERSION=$(cat "${BUILD_DIR}/_output/version.txt")
COHERENCE_OPERATORS_REPO=coherence-community/certified-operators
GIT_REPO_URL=https://github.com/${COHERENCE_OPERATORS_REPO}.git
GIT_CERT_BRANCH=cert-temp
BUNDLE_PATH=operators/oracle-coherence/${OPERATOR_VERSION}

if [ ! -e ${BUILD_DIR}/certified-operators ]; then
  cd ${BUILD_DIR}
  git clone --quiet ${GIT_REPO_URL} certified-operators
fi

cd ${BUILD_DIR}/certified-operators
git checkout main
git pull

GIT_ORIGIN=$(git config remote.origin.url)
GITHUB_PUSH_UPSTREAM="${GIT_ORIGIN}"
if [ ! -z "${GITHUB_TOKEN:-}" ]; then
  GITHUB_PUSH_UPSTREAM=$(echo "${GIT_ORIGIN}" | sed -e s#://#://${GITHUB_USERNAME}:${GITHUB_TOKEN}@#)
  git config user.name "${GITHUB_USERNAME}"
  if [ ! -z "${GITHUB_USER_EMAIL:-}" ]; then
    git config user.email "${GITHUB_USER_EMAIL}"
  fi
fi

git branch ${GIT_CERT_BRANCH} -d || true
git checkout -b ${GIT_CERT_BRANCH}
rm -rf ${BUNDLE_PATH}
mkdir -p ${BUNDLE_PATH}
cp -R ${BUILD_DIR}/_output/bundle/coherence-operator/ operators/oracle-coherence/
cp ${ROOT_DIR}/bundle/ci.yaml operators/oracle-coherence/
git add -A operators/oracle-coherence/*
git status
git commit -m "Adding Oracle Coherence Operator v${OPERATOR_VERSION}"
git push -u "${GITHUB_PUSH_UPSTREAM}" -f cert-temp

# Step F - Run Pipeline
cd "${ROOT_DIR}"

if [ -z "${UPSTREAM_REPO_NAME:-}" ]; then
  UPSTREAM_REPO_NAME=${COHERENCE_OPERATORS_REPO}
  echo "UPSTREAM_REPO_NAME is not set, defaulting to ${UPSTREAM_REPO_NAME}"
fi

oc apply --filename "${ROOT_DIR}/tekton/workspace-pv.yaml"

# Delete any old runs
oc delete $(tkn pipelinerun list -o name) || true

PIPELINE_TIMESTAMP=$(date +"%Y%m%d%H%M")
PIPELINE_RUN_NAME="operator-cert-run-${PIPELINE_TIMESTAMP}"
echo "Using PIPELINE_RUN_NAME ${PIPELINE_RUN_NAME}"
cp ${ROOT_DIR}/hack/openshift/pipeline-run.yaml "${ROOT_DIR}/run.yaml"
sed -i -e "s/NAME_PLACEHOLDER/${PIPELINE_RUN_NAME}/g" "${ROOT_DIR}/run.yaml"
sed -i -e "s^GIT_REPO_PLACEHOLDER^${GIT_REPO_URL}^g" "${ROOT_DIR}/run.yaml"
sed -i -e "s^GIT_CERT_BRANCH_PLACEHOLDER^${GIT_CERT_BRANCH}^g" "${ROOT_DIR}/run.yaml"
sed -i -e "s^BUNDLE_PATH_PLACEHOLDER^${BUNDLE_PATH}^g" "${ROOT_DIR}/run.yaml"
sed -i -e "s^GITHUB_USERNAME_PLACEHOLDER^${GITHUB_USERNAME}^g" "${ROOT_DIR}/run.yaml"
sed -i -e "s^GITHUB_USER_EMAIL_PLACEHOLDER^${GITHUB_USER_EMAIL}^g" "${ROOT_DIR}/run.yaml"
sed -i -e "s^UPSTREAM_REPO_NAME_PLACEHOLDER^${UPSTREAM_REPO_NAME}^g" "${ROOT_DIR}/run.yaml"
sed -i -e "s^REGISTRY_HOST_PLACEHOLDER^${REGISTRY_HOST}^g" "${ROOT_DIR}/run.yaml"
sed -i -e "s^REGISTRY_NAMESPACE_PLACEHOLDER^${REGISTRY_NAMESPACE}^g" "${ROOT_DIR}/run.yaml"

cat "${ROOT_DIR}/run.yaml"
oc create -f "${ROOT_DIR}/run.yaml"
rm run.yaml

tkn pipelinerun logs "${PIPELINE_RUN_NAME}" -n "${PIPELINE_NAMESPACE}" --follow
tkn pipelinerun describe "${PIPELINE_RUN_NAME}" -n "${PIPELINE_NAMESPACE}"
tkn pipelinerun describe "${PIPELINE_RUN_NAME}" -n "${PIPELINE_NAMESPACE}" -o jsonpath="{.status.conditions[0].reason}" > pipeline-result.txt
echo "Pipeline result"
cat pipeline-result.txt

#echo "Running full pipeline with digest pinning"
#tkn pipeline start operator-ci-pipeline \
#  --use-param-defaults \
#  --param git_repo_url=${GIT_REPO_URL} \
#  --param git_branch=${GIT_CERT_BRANCH} \
#  --param bundle_path=${BUNDLE_PATH} \
#  --param env=prod \
#  --param pin_digests=true \
#  --param git_username=${GITHUB_USERNAME} \
#  --param git_email=${GITHUB_USER_EMAIL} \
#  --param upstream_repo_name=${UPSTREAM_REPO_NAME} \
#  --param registry=${REGISTRY_HOST} \
#  --param image_namespace=${REGISTRY_NAMESPACE} \
#  --param submit=false \
#  --workspace name=pipeline,volumeClaimTemplateFile=${PIPELINES_DIR}/templates/workspace-template.yml \
#  --workspace name=ssh-dir,secret=github-ssh-credentials \
#  --workspace name=registry-credentials,secret=registry-dockerconfig-secret \
#  --showlog

