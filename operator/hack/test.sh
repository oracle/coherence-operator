#!/usr/bin/env bash

set -o errexit
set -o nounset

export CGO_ENABLED=0

echo "Running operator tests"


GINKGO=$(which ginkgo)
if [[ "${GINKGO}" == "" ]]
then
    CMD="go"
else
    CMD="ginkgo"
fi

exec ${CMD} test -v ./cmd/...  ./pkg/...