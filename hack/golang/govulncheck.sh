#!/bin/sh
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
set -o errexit

ROOT_DIR=$(pwd)
TOOLS_BIN=${ROOT_DIR}/build/tools/bin

test -s ${TOOLS_BIN}/govulncheck || GOBIN=${TOOLS_BIN} go install golang.org/x/vuln/cmd/govulncheck@latest
chmod +x ${TOOLS_BIN}/govulncheck

go version

make build-operator-images

echo "INFO: govulncheck - Checking x84_64 runner"
${TOOLS_BIN}/govulncheck -mode binary -show traces,version,verbose  ./bin/linux/amd64/runner
echo "INFO: govulncheck - Checking x84_64 cohctl"
${TOOLS_BIN}/govulncheck -mode binary -show traces,version,verbose  ./bin/linux/amd64/cohctl

echo "INFO: govulncheck - Checking Arm64 runner"
${TOOLS_BIN}/govulncheck -mode binary -show traces,version,verbose  ./bin/linux/arm64/runner
echo "INFO: govulncheck - Checking Arm64 cohctl"
${TOOLS_BIN}/govulncheck -mode binary -show traces,version,verbose  ./bin/linux/arm64/cohctl

