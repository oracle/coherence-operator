/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package data

import "embed"

// Assets is the set of embedded files.
//go:embed assets/*
var Assets embed.FS
