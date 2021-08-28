/*
 * Copyright (c) 2021, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package v1

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"k8s.io/apimachinery/pkg/util/rand"
	hashutil "k8s.io/kubernetes/pkg/util/hash"
)

func EnsureHashLabel(c *Coherence) (string, bool) {
	labels := c.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	hashNew := ComputeHash(&c.Spec, nil)
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
func ComputeHash(template *CoherenceResourceSpec, collisionCount *int32) string {
	podTemplateSpecHasher := fnv.New32a()
	hashutil.DeepHashObject(podTemplateSpecHasher, *template)

	// Add collisionCount in the hash if it exists.
	if collisionCount != nil {
		collisionCountBytes := make([]byte, 8)
		binary.LittleEndian.PutUint32(collisionCountBytes, uint32(*collisionCount))
		_, _ = podTemplateSpecHasher.Write(collisionCountBytes)
	}

	return rand.SafeEncodeString(fmt.Sprint(podTemplateSpecHasher.Sum32()))
}
