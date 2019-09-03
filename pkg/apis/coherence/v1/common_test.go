/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1_test

// Returns a pointer to an int32
func int32Ptr(x int32) *int32 {
	return &x
}

// Returns a pointer to an int32
func boolPtr(x bool) *bool {
	return &x
}

// Returns a pointer to a string
func stringPtr(x string) *string {
	return &x
}
