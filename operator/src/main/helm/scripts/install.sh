#
# Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
#!/bin/bash -e

# Install Helm
if [ "$(uname -s)" = "Darwin" ]; then
    # We're on Mac OS
    OS_SUFFIX="darwin_amd64"
else
    # We're on linux
    OS_SUFFIX="linux-amd64"
fi

if [ ! -f /usr/local/bin/helm ]; then
    echo "Installing Helm..."
    HELM_LATEST_VERSION="v2.11.0"
    HELM_INSTALL=helm-${HELM_LATEST_VERSION}-${OS_SUFFIX}.tar.gz
    HELM_URL="http://storage.googleapis.com/kubernetes-helm/${HELM_INSTALL}"
    echo "Downloading Helm from ${HELM_URL}..."
    curl -fsSL -o ${HELM_INSTALL} ${HELM_URL}
    tar -xvf ${HELM_INSTALL}
    sudo mv linux-amd64/helm /usr/local/bin
    rm -f ${HELM_INSTALL}
    rm -rf linux-amd64

    # Setup Helm so that it will work with helm dep commands. Only the client
    # needs to be setup. In addition, the incubator repo needs to be
    # available for charts that depend on it.
    helm init -c
    #helm repo add incubator https://kubernetes-charts-incubator.storage.googleapis.com/
fi

# Install A YAML Linter
# Pinning to a version for consistency
echo "Installing yamllint..."
sudo pip install --proxy=${http_proxy} yamllint==1.8.1

# Install YQ YAML Command line reader
if [ ! -f /usr/local/bin/yq ]; then
    echo "Installing YQ..."
    YQ_VERSION=1.14.0
    YQ_INSTALL=yq_${OS_SUFFIX}
    YQ_URL="https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/${YQ_INSTALL}"
    echo "Downloading YQ from ${YQ_URL}..."
    curl -fsSL -o ${YQ_INSTALL} ${YQ_URL}
    chmod +x ${YQ_INSTALL}
    sudo mv ${YQ_INSTALL} /usr/local/bin/yq
fi

# Install SemVer testing tool
if [ ! -f /usr/local/bin/vert ]; then
    if [ "$(uname -s)" != "Darwin" ]; then
        echo "Installing Vert..."
        VERT_VERSION=v0.1.0
        VERT_INSTALL=vert-${VERT_VERSION}-${OS_SUFFIX}
        VERT_URL="http://github.com/Masterminds/vert/releases/download/${VERT_VERSION}/${VERT_INSTALL}"
        echo "Downloading Vert from ${VERT_URL}..."
        curl -fsSL -o ${VERT_INSTALL} ${VERT_URL}
        chmod +x ${VERT_INSTALL}
        sudo mv ${VERT_INSTALL} /usr/local/bin/vert
    fi
fi
