/*
 * Copyright (c) 2019, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

// The local package contains end-to-end Operator tests that require the
// Operator to be deployed into the K8s cluster.
// These tests use the operator-sdk end-to-end framework and must be run
// using the "operator-sdk test local" command but without the "--up-local"
// parameter.
package remote
