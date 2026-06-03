#!/bin/sh
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

set -o errexit

# desired cluster name; default is "kind"
KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-operator}"
# Allow callers to use the project-pinned kind binary so this wrapper creates
# clusters with the same CLI compatibility guarantees as the Makefile targets.
KIND_CMD="${KIND_CMD:-kind}"

"${KIND_CMD}" create cluster --name "${KIND_CLUSTER_NAME}" "$@"
