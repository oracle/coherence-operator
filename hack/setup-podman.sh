#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

#!/usr/bin/env bash

# --------------------------------------------------------------------
# This script exports various environment variables so that
# Make targets will be executed using Podman.
# --------------------------------------------------------------------

export DOCKER_CMD=podman
export DOCKER_HOST=unix://$XDG_RUNTIME_DIR/podman/podman.sock
export JIB_EXECUTABLE=$(which podman)
export MY_DOCKER_HOST=${DOCKER_HOST}
