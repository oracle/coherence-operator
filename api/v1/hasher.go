/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"k8s.io/apimachinery/pkg/util/rand"
)

func EnsureCoherenceHashLabel(c *Coherence) (string, bool) {
	labels := c.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	spec := c.Spec
	hashNew := ComputeHash(&spec, nil)
	hashCurrent, found := labels[LabelCoherenceHash]
	if !found || hashCurrent != hashNew {
		labels[LabelCoherenceHash] = hashNew
		c.SetLabels(labels)
		return hashNew, true
	}
	return hashCurrent, false
}

func EnsureJobHashLabel(c *CoherenceJob) (string, bool) {
	labels := c.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	spec := c.Spec
	hashNew := ComputeHash(&spec, nil)
	hashCurrent, found := labels[LabelCoherenceHash]
	if !found || hashCurrent != hashNew {
		labels[LabelCoherenceHash] = hashNew
		c.SetLabels(labels)
		return hashNew, true
	}
	return hashCurrent, false
}

// ComputeHash returns a hash value calculated from Coherence spec and
// The hash will be safe encoded to avoid bad words.
func ComputeHash(in interface{}, collisionCount *int32) string {
	hasher := fnv.New32a()
	b, _ := json.Marshal(in)
	_, _ = hasher.Write(b)

	// Add collisionCount in the hash if it exists.
	if collisionCount != nil {
		collisionCountBytes := make([]byte, 8)
		binary.LittleEndian.PutUint32(collisionCountBytes, uint32(*collisionCount))
		_, _ = hasher.Write(collisionCountBytes)
	}

	return rand.SafeEncodeString(fmt.Sprint(hasher.Sum32()))
}
