#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

#!/usr/bin/env bash

export OPERATOR_IMAGE_REGISTRY=iad.ocir.io/odx-stateservice/test
export DOCKER_CMD=podman
export DOCKER_HOST=unix://$XDG_RUNTIME_DIR/podman/podman.sock
export MY_DOCKER_HOST=${DOCKER_HOST}
export JIB_EXECUTABLE=$(which podman)
export USE_PODMAN=true
export LOCAL_BUILDAH=true
export DEPLOY_DOCKER_CONFIG_JSON=$XDG_RUNTIME_DIR/containers/auth.json
