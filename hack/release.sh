#
# Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
#!/bin/bash

# -----------------------------------------------------------------------------
# This script uses below set of environment variables for doing coherence
# operator release.
#
# Required Set of Enviornment Varibales:
# 1) BRANCH_NAME           : Name of the GIT branch to be used to run release.
# 2) RELEASE_SUFFIX        : The release suffix, e.g. RC1, alpha1 etc.
# 3) NEXT_VERSION          : An optional next version to use if doing a full
#                            release and bumping the version.
# 4) RELEASE_IMAGE_PREFIX  : Docker repository prefix to be used for
#                            Coherence Operator docker image.
# 5) DRY_RUN               : To indicate whether to run script in dry mode,
#                            defaults to true if not specified.
# -----------------------------------------------------------------------------

# -----------------------------------------------------------------------------
# Setup Release Branch
# -----------------------------------------------------------------------------
setupReleaseBranch()
  {
  if [ -z $(git ls-remote -q --tags | grep $RELEASE_TAG) ]; then

    pwd
    ls -ls

    git commit -a -m "Preparing for Release Version $RELEASE_VERSION"
    git tag $RELEASE_TAG

    if [ "false" = "$DRY_RUN" ]; then
      git push origin $RELEASE_TAG
    fi

    if [ 0 -eq $? ]; then
#     Only udate the version if NEXT_VERSION is set
      if [[ "${NEXT_VERSION}" != "" ]]; then
        awk '{sub(/^VERSION \?= [0-9]\.[0-9]\.[0-9]/,"VERSION ?= '${NEXT_VERSION}'")}1' Makefile > Makefile.temp
        mv Makefile.temp Makefile
        mvn -f java -DnewVersion=${NEXT_VERSION} versions:set versions:commit

        git commit -a -m "Preparing for Next Development version."
        if [ "false" = "$DRY_RUN" ]; then
          git push
        fi
        return $?
      fi
    else
      return 1
    fi
  else
    echo ""
    echo "Git tag $RELEASE_TAG already exists releasing from existing release tag $RELEASE_TAG ..."
    echo ""
  fi
  }

# -----------------------------------------------------------------------------
# Build Release Branch & Push Docker Images
# -----------------------------------------------------------------------------
buildReleaseBranch()
  {
  [ -n "${WORKSPACE}" ] || local WORKSPACE=`mktemp -d`
  export RELEASE_DIR=${WORKSPACE}/release-$RELEASE_TAG
  echo "RELEASE_DIR = $RELEASE_DIR"
  rm -fr $RELEASE_DIR
  mkdir -p $RELEASE_DIR

  git ls-tree HEAD --name-only | xargs -I'{}' sh -c 'cp -r $1 $RELEASE_DIR' -- {}
  cp -r .git $RELEASE_DIR
  cd $RELEASE_DIR

  git checkout $RELEASE_TAG
  git status
  pwd

  make build-all-images VERSION_SUFFIX="${RELEASE_SUFFIX}"

  if [ "false" = "$DRY_RUN" ]; then
    make push-all-images VERSION_SUFFIX="${RELEASE_SUFFIX}"
  fi

  STATUS=$?
  if [ 0 -eq "$STATUS" ]; then
    export COH_OP_CHART=$(find build/_output/helm-charts -regex '.*coherence-operator.*-helm.tar.gz' -print)
    echo COH_OP_CHART=$COH_OP_CHART
  fi

  return $STATUS
  }

# -----------------------------------------------------------------------------
# Check for required environment variables pointing to the coherence operator
# chart after building the release branch.
# -----------------------------------------------------------------------------
checkRequiredEnvVars()
  {
  if [[ -z "$COH_OP_CHART" ]]; then
    echo "Required envrionment variable COH_OP_CHART pointing to coherence-operator chart is not set."
    return 1
  fi
  }

# -----------------------------------------------------------------------------
# Publish helm charts to helm repo.
# -----------------------------------------------------------------------------
publishCharts()
  {
  if [ ! -d charts ]; then
    mkdir charts
  fi

  cp $COH_OP_CHART charts/

  git checkout gh-pages
  if [ 0 -ne "$?" ]; then
    echo "Failed to switch to the required gh-pages branch."
    return 1
  fi

  echo "Running helm repo index command ..."
  helm repo index charts --url https://oracle.github.io/coherence-operator/charts

  git status
  git add charts/*
  git clean -d -f
  git status

  git config user.name "Coherence Bot"
  git config user.email coherence-bot_ww@oracle.com
  git commit -m "Release coherence-operator helm chart version: $RELEASE_VERSION"
  if [ 0 -ne $? ]; then
    echo "Failed to commit the changes containing coherence-operator helm chart."
    return 1
  fi

  git log -1
  if [ "false" = "$DRY_RUN" ]; then
    git push origin gh-pages
  fi

  return $?
  }

# -----------------------------------------------------------------------------
# Display error message($1) with the given exit status($2)
# -----------------------------------------------------------------------------
errorMessage()
  {
  echo "$1 $2"
  return $2
  }

DRY_RUN=${DRY_RUN:-true}
echo "DRY_RUN = $DRY_RUN"

if [[ -n "$BRANCH_NAME" ]]; then

#  git checkout $BRANCH_NAME
  RELEASE_VERSION=$(make version VERSION_SUFFIX="${RELEASE_SUFFIX}")
  RELEASE_TAG=v$RELEASE_VERSION

  echo "RELEASE_SUFFIX = ${RELEASE_SUFFIX}"
  echo "RELEASE_VERSION = ${RELEASE_VERSION}"
  echo "RELEASE_TAG = ${RELEASE_TAG}"
  echo "NEXT_VERSION = ${NEXT_VERSION}"
  echo "RELEASE_IMAGE_PREFIX = ${RELEASE_IMAGE_PREFIX}"

  setupReleaseBranch
  SETUP_BRANCH_STATUS=$?
  echo "SETUP_BRANCH_STATUS == $SETUP_BRANCH_STATUS"

  if [ 0 -eq $SETUP_BRANCH_STATUS ]; then
    buildReleaseBranch
    BUILD_STATUS=$?
    echo "BUILD STATUS == $BUILD_STATUS"
    if [ 0 -ne $BUILD_STATUS ]; then
      errorMessage "Build process failed with exit " $BUILD_STATUS
    fi
  else
    errorMessage "Setting up release branch failed with exit " $SETUP_BRANCH_STATUS
  fi
else
  errorMessage "Required environment variable BRANCH_NAME is not set so exit with status " 1
fi

checkRequiredEnvVars

if [ 0 -eq $? ]; then
  echo "PWD = $(pwd)"
  publishCharts
fi
