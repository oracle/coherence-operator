#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

#!/usr/bin/env bash

# --------------------------------------------------------------------
# This script exports various environment variables so that
# Make targets will be executed using Podman and OpenShift.
# --------------------------------------------------------------------

export DOCKER_CMD=podman
export JIB_EXECUTABLE=$(which podman)
