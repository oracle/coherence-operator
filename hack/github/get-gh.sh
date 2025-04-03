#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
set -o errexit


ROOT_DIR=$(pwd)
TOOLS_BIN=${ROOT_DIR}/build/tools/bin

UNAME_S=$(uname -s)
UNAME_M=$(uname -m)

GH_CLI_VERSION=2.69.0

if [ "Darwin" = "${UNAME_S}" ]; then
  if [ "x86_64" = "${UNAME_M}" ]; then
    echo "Downloading GitHub CLI ${VERSION} ${UNAME_S} ${UNAME_M}"
    curl -Ls https://github.com/cli/cli/releases/download/v${GH_CLI_VERSION}/gh_${GH_CLI_VERSION}_macos_amd64.zip -o gh.zip
    unzip gh.zip
    mv gh_${GH_CLI_VERSION}_macOS_amd64/bin/gh ${TOOLS_BIN}/
    rm -rf gh_${GH_CLI_VERSION}_macOS_amd64
    rm gh.zip
  else
    echo "Downloading GitHub CLI ${VERSION} ${UNAME_S} ${UNAME_M}"
    curl -Ls https://github.com/cli/cli/releases/download/v${GH_CLI_VERSION}/gh_${GH_CLI_VERSION}_macos_arm64.zip -o gh.zip
    unzip gh.zip
    mv gh_${GH_CLI_VERSION}_macOS_arm64/bin/gh ${TOOLS_BIN}/
    rm -rf gh_${GH_CLI_VERSION}_macOS_arm64
    rm gh.zip
  fi
else
  if [ "x86_64" = "${UNAME_M}" ]; then
    echo "Downloading GitHub CLI ${VERSION} ${UNAME_S} ${UNAME_M}"
    curl -Ls https://github.com/cli/cli/releases/download/v${GH_CLI_VERSION}/gh_${GH_CLI_VERSION}_linux_amd64.tar.gz -o gh.tar.gz
    tar -xvf gh.tar.gz
    mv gh_${GH_CLI_VERSION}_linux_amd64/bin/gh ${TOOLS_BIN}/
    rm -rf gh_${GH_CLI_VERSION}_linux_amd64
  else
    echo "Downloading GitHub CLI ${VERSION} ${UNAME_S} ${UNAME_M}"
    curl -Ls https://github.com/cli/cli/releases/download/v${GH_CLI_VERSION}/gh_${GH_CLI_VERSION}_linux_arm64.tar.gz -o gh.tar.gz
    tar -xvf gh.tar.gz
    mv gh_${GH_CLI_VERSION}_linux_arm64/bin/gh ${TOOLS_BIN}/
    rm -rf gh_${GH_CLI_VERSION}_linux_arm64
    rm gh.tar.gz
  fi
fi

chmod +x ${TOOLS_BIN}/gh
