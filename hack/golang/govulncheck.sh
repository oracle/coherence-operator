#!/bin/sh
#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#


go install golang.org/x/vuln/cmd/govulncheck@latest
make runner
govulncheck -mode binary -show traces,version,verbose  ./bin/runner
