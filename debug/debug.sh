#!/bin/bash -e
#
# Copyright (c) 2021, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec manager