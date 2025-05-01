#!/usr/bin/env bash
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

make bundle
make scorecard
make bundle-push
make catalog-build
make catalog-push

make olm-undeploy
make olm-undeploy-catalog

make olm-deploy-catalog
make wait-for-olm-catalog-deploy
make olm-deploy


