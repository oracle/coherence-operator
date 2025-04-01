#!/bin/bash
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
VERSION=$(curl -s https://oracle.github.io/coherence-cli/stable.txt)

function set_arch() {
    if [ "$ARCH" == "x86_64" ] ; then
        ARCH="amd64"
    elif [ "$ARCH" == "aarch64" -o "$ARCH" == "arm64" ] ; then
        ARCH="arm64"
    else
        echo "Unsupported architecture: $ARCH"
        exit 1
    fi
}

function installed() {
    echo "Installed cohctl into ${COHCTL_HOME}"
}

echo "Installing Coherence CLI ${VERSION} for ${OS}/${ARCH} into ${COHCTL_HOME} ..."

if [ "$OS" == "Darwin" ]; then
    set_arch
    TEMP=`mktemp -d`
    PKG="Oracle-Coherence-CLI-${VERSION}-darwin-${ARCH}.pkg"
    DEST=${TEMP}/${PKG}
    echo "Downloading and opening ${DEST}"
    URL=https://github.com/oracle/coherence-cli/releases/download/${VERSION}/${PKG}
    curl -sLo  ${DEST} $URL && open ${DEST} && installed
elif [ "$OS" == "Linux" ]; then
    set_arch
    TEMP=`mktemp -d`
    URL=https://github.com/oracle/coherence-cli/releases/download/${VERSION}/cohctl-${VERSION}-linux-${ARCH}
    curl -sLo ${TEMP}/cohctl $URL && chmod u+x ${TEMP}/cohctl
    mv ${TEMP}/cohctl ${COHCTL_HOME} && installed
else
    echo "For all other platforms, please see: https://github.com/oracle/coherence-cli/releases"
    exit 1
fi
