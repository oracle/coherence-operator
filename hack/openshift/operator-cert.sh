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

if [ -z "${SUBMIT_RESULTS:-}" ]; then
  SUBMIT_RESULTS=false
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

# If the OpenShift kubeconfig env var is set use that for the
# kubeconfig location, otherwise use the default KUBECONFIG
# that was set earlier
if [ -z "${OPENSHIFT_KUBECONFIG:-}" ]; then
  export OPENSHIFT_KUBECONFIG="${KUBECONFIG}"
fi
if oc get secret kubeconfig >/dev/null 2>&1; then
  oc delete secret kubeconfig
fi
oc create secret generic kubeconfig --from-file=kubeconfig="${OPENSHIFT_KUBECONFIG}"

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
COHERENCE_OPERATORS_REPO=coherence-community/certified-operators
GIT_REPO_URL=https://github.com/${COHERENCE_OPERATORS_REPO}.git
GIT_CERT_BRANCH=cert-temp
LATEST_RELEASE=""

# If this is a certification run to submit a new release then we always
# clone a new repo, even one exist locally
if [ "${USE_LATEST_OPERATOR_RELEASE}" = "true" ]; then
# delete the old local repo to force a new one to be cloned
  rm -rf "${BUILD_DIR}/certified-operators"
# Find the latest release of the Coherence Operator on GitHub
  LATEST_RELEASE=$(gh release list --repo oracle/coherence-operator --json name,isLatest --jq '.[] | select(.isLatest)|.name')
# Strip the v from the front of the release to give the Operator version
  OPERATOR_VERSION=${LATEST_RELEASE#"v"}
  echo "Latest Operator version is ${OPERATOR_VERSION}"
# Check the latest release image exists on OCR
  COHERENCE_OPERATOR_IMAGE="container-registry.oracle.com/middleware/coherence-operator:${OPERATOR_VERSION}"
  echo "Checking Oracle Container Registry for image ${OCR_COHERENCE_IMAGE}"
  podman manifest inspect "${OCR_COHERENCE_IMAGE}" > /dev/null
  if [ $? -ne 0 ]; then
    echo "ERROR: Image ${OCR_COHERENCE_IMAGE} does not exist on OCR."
    exit 1
  fi
# Use a proper name for the git branch
  GIT_CERT_BRANCH="release-${OPERATOR_VERSION}"
# Set the upstream repo for the pull request to be the official RedHat repo
  UPSTREAM_REPO_NAME=redhat-openshift-ecosystem/certified-operators
# make sure the certified-operators repo in Coherence Community is sync'ed with the RedHat repo
  gh repo sync "${COHERENCE_OPERATORS_REPO}"
else
# We are just testing a local build, so use the current version
  OPERATOR_VERSION=$(cat "${BUILD_DIR}/_output/version.txt")
# We will not be submitting results
  SUBMIT_RESULTS=false
  COHERENCE_OPERATOR_IMAGE="${REGISTRY_HOST}/${REGISTRY_NAMESPACE}/coherence-operator:${OPERATOR_VERSION}"
fi

if [ -z "${OPERATOR_VERSION:-}" ]; then
  echo "Error: OPERATOR_VERSION has not been set"
  exit 1
fi
BUNDLE_PATH=operators/oracle-coherence/${OPERATOR_VERSION}

# If the certified-operators repo does not exist locally then clone it
if [ ! -e "${BUILD_DIR}/certified-operators" ]; then
  echo "Cloning repo ${GIT_REPO_URL}"
  cd "${BUILD_DIR}"
  git clone --quiet ${GIT_REPO_URL} certified-operators
fi

# cd to the certified-operators local repo and refresh
cd "${BUILD_DIR}/certified-operators"
git checkout main
git pull

# Configure Git in the local repo so we can push to it
GIT_ORIGIN=$(git config remote.origin.url)
GITHUB_PUSH_UPSTREAM="${GIT_ORIGIN}"
if [ ! -z "${GITHUB_TOKEN:-}" ]; then
  GITHUB_PUSH_UPSTREAM=$(echo "${GIT_ORIGIN}" | sed -e s#://#://${GITHUB_USERNAME}:${GITHUB_TOKEN}@#)
  git config user.name "${GITHUB_USERNAME}"
  if [ ! -z "${GITHUB_USER_EMAIL:-}" ]; then
    git config user.email "${GITHUB_USER_EMAIL}"
  fi
fi

# Delete the branch from GitHub (if it exists)
git push -u "${GITHUB_PUSH_UPSTREAM}" -d ${GIT_CERT_BRANCH} || true
# Delete the pinned branch from GitHub (if it exists)
git push -u "${GITHUB_PUSH_UPSTREAM}" -d ${GIT_CERT_BRANCH}-pinned || true
# Delete the local branch (if it exists)
git branch ${GIT_CERT_BRANCH} -D || true
# Delete the local pinned branch (if it exists)
git branch ${GIT_CERT_BRANCH}-pinned -D || true
# Create a new local branch
git checkout -b ${GIT_CERT_BRANCH}

if [ "${USE_LATEST_OPERATOR_RELEASE}" = "true" ]; then
# We are certifying a real release so make sure the latest release does not already exist
# in the certified-operators repo
  DIR_NAME=operators/coherence-operator/${OPERATOR_VERSION}
  if [ -d ${DIR_NAME} ]; then
    echo "Coherence Operator ${OPERATOR_VERSION} is already submitted to ${GIT_REPO_URL}"
    exit 1
  fi
# download the bundle tar.gz from the Operator release on GitHub
  gh release download ${LATEST_RELEASE} --repo oracle/coherence-operator --pattern coherence-operator-bundle.tar.gz
# unpack the tar.gz into a temp location
  rm -rf bundle-temp || true
  TEMP_BUNDLE_DIR=bundle-temp
  mkdir -p "${TEMP_BUNDLE_DIR}"
  tar -xvf coherence-operator-bundle.tar.gz -C "${TEMP_BUNDLE_DIR}/"
  rm coherence-operator-bundle.tar.gz
# copy the bundle contents to the actual location in the certified-operators repo
  mkdir -p "${BUNDLE_PATH}"
  cp -R ""${TEMP_BUNDLE_DIR}"/coherence-operator/${OPERATOR_VERSION}" operators/oracle-coherence/
#  make sure the ci.yaml file exists
  echo "cert_project_id: ${OPENSHIFT_PROJECT_ID}" > operators/oracle-coherence/ci.yaml
  rm -rf "${TEMP_BUNDLE_DIR}" || true
else
#  we are testing a local build, so copy the local bundle folder to the certified-operators repo
  rm -rf "${BUNDLE_PATH}"
  mkdir -p "${BUNDLE_PATH}"
  cp -R "${BUILD_DIR}/_output/bundle/coherence-operator/${OPERATOR_VERSION}" operators/oracle-coherence/
  cp "${ROOT_DIR}/bundle/ci.yaml" operators/oracle-coherence/
fi

# Add the new bundle files to git, commit and push them
git add -A operators/oracle-coherence/*
git status
git commit -m "Adding Oracle Coherence Operator v${OPERATOR_VERSION}"
git push -u "${GITHUB_PUSH_UPSTREAM}" -f "${GIT_CERT_BRANCH}"

# Step F - Run Pipeline
cd "${ROOT_DIR}"

if [ -z "${UPSTREAM_REPO_NAME:-}" ]; then
  UPSTREAM_REPO_NAME=${COHERENCE_OPERATORS_REPO}
  echo "UPSTREAM_REPO_NAME is not set, defaulting to ${UPSTREAM_REPO_NAME}"
fi

# Delete any old runs
oc delete $(tkn pipelinerun list -o name) || true

GITHUB_REPO_URL="git@github.com:${COHERENCE_OPERATORS_REPO}.git"

PIPELINE_TIMESTAMP=$(date +"%Y%m%d%H%M")
PIPELINE_RUN_NAME="operator-cert-run-${PIPELINE_TIMESTAMP}"
echo "Using PIPELINE_RUN_NAME ${PIPELINE_RUN_NAME}"
cp ${ROOT_DIR}/hack/openshift/pipeline-run.yaml "${ROOT_DIR}/run.yaml"
sed -i -e "s/PIPELINE_NAME_PLACEHOLDER/${PIPELINE_RUN_NAME}/g" "${ROOT_DIR}/run.yaml"
sed -i -e "s^GIT_REPO_PLACEHOLDER^${GITHUB_REPO_URL}^g" "${ROOT_DIR}/run.yaml"
sed -i -e "s^GIT_CERT_BRANCH_PLACEHOLDER^${GIT_CERT_BRANCH}^g" "${ROOT_DIR}/run.yaml"
sed -i -e "s^BUNDLE_PATH_PLACEHOLDER^${BUNDLE_PATH}^g" "${ROOT_DIR}/run.yaml"
sed -i -e "s^GITHUB_USERNAME_PLACEHOLDER^${GITHUB_USERNAME}^g" "${ROOT_DIR}/run.yaml"
sed -i -e "s^GITHUB_USER_EMAIL_PLACEHOLDER^${GITHUB_USER_EMAIL}^g" "${ROOT_DIR}/run.yaml"
sed -i -e "s^UPSTREAM_REPO_NAME_PLACEHOLDER^${UPSTREAM_REPO_NAME}^g" "${ROOT_DIR}/run.yaml"
sed -i -e "s^REGISTRY_HOST_PLACEHOLDER^${REGISTRY_HOST}^g" "${ROOT_DIR}/run.yaml"
sed -i -e "s^REGISTRY_NAMESPACE_PLACEHOLDER^${REGISTRY_NAMESPACE}^g" "${ROOT_DIR}/run.yaml"
sed -i -e "s^SUBMIT_RESULTS_PLACEHOLDER^${SUBMIT_RESULTS}^g" "${ROOT_DIR}/run.yaml"

cat "${ROOT_DIR}/run.yaml"
oc create -f "${ROOT_DIR}/run.yaml"
rm run.yaml

tkn pipelinerun logs "${PIPELINE_RUN_NAME}" -n "${PROJECT_NAME}" --follow
tkn pipelinerun describe "${PIPELINE_RUN_NAME}" -n "${PROJECT_NAME}"
tkn pipelinerun describe "${PIPELINE_RUN_NAME}" -n "${PROJECT_NAME}" -o jsonpath="{.status.conditions[0].reason}" > pipeline-result.txt
echo "Pipeline result"
cat pipeline-result.txt

