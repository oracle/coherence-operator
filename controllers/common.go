/*
 * Copyright (c) 2020, 2025, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package controllers

import (
	"github.com/go-logr/logr"
	coh "github.com/oracle/coherence-operator/api/v1"
	"github.com/oracle/coherence-operator/pkg/utils"
)

func checkCoherenceHash(deployment *coh.Coherence, storage utils.Storage, log logr.Logger) (coh.Resources, error) {
	return checkHash(deployment, deployment.Status.Phase, storage, log)
}

func checkJobHash(deployment *coh.CoherenceJob, storage utils.Storage, log logr.Logger) (coh.Resources, error) {
	return checkHash(deployment, deployment.Status.Phase, storage, log)
}

func checkHash(deployment coh.CoherenceResource, phase coh.ConditionType, storage utils.Storage, log logr.Logger) (coh.Resources, error) {
	hash := deployment.GetLabels()[coh.LabelCoherenceHash]
	var desiredResources coh.Resources
	var err error

	storeHash, found := storage.GetHash()
	if !found || storeHash != hash || phase != coh.ConditionTypeReady {
		// Storage state was saved with no hash or a different hash so is not in the desired state
		// or the Coherence resource is not in the Ready state
		// Create the desired resources the deployment
		if desiredResources, err = deployment.CreateKubernetesResources(); err != nil {
			return desiredResources, err
		}

		if found {
			// The "storeHash" is not "", so it must have been processed by the Operator (could have been a previous version).
			// There was a bug prior to 3.4.2 where the hash was calculated at the wrong point in the defaulting web-hook,
			// and the has used only a portion of the spec, so the "currentHash" may be wrong, and hence differ from the
			// recalculated "hash".
			if deployment.IsBeforeVersion("3.4.2") {
				// the AnnotationOperatorVersion annotation was added in the 3.2.8 web-hook, so if it is missing
				// the Coherence resource was added or updated prior to 3.2.8, or the version is present but is
				// prior to 3.4.2. In this case we just ignore the difference in hash.
				// There is an edge case where the Coherence resource could have legitimately been updated whilst
				// the Operator and web-hooks were uninstalled. In that case we would ignore the update until another
				// update is made. The simplest way for the customer to work around this is to add the
				// AnnotationOperatorVersion annotation with some value, which will then be overwritten by the web-hook
				// and the Coherence resource will be correctly processes.
				desiredResources = storage.GetLatest()
				log.Info("Ignoring hash difference for pre-3.4.2 resource", "hash", hash, "store", storeHash)
			}
		}
	} else {
		// storage state was saved with the current hash so is already in the desired state
		desiredResources = storage.GetLatest()
	}
	return desiredResources, err
}
