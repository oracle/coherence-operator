#!/usr/bin/env bash

#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

# Apply the custom node labels required for site and rack tests
NODES=$(kubectl get nodes -o name)
for NODE in $NODES; do
  kubectl label $NODE coherence.oracle.com/site=test-site --overwrite
  kubectl label $NODE coherence.oracle.com/rack=test-rack --overwrite
done
