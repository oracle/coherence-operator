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

func getDesiredResources(deployment *coh.Coherence, storage utils.Storage, log logr.Logger) (coh.Resources, error) {
	return checkHash(deployment, deployment.Status.Phase, storage, log)
}

func getDesiredJobResources(deployment *coh.CoherenceJob, storage utils.Storage, log logr.Logger) (coh.Resources, error) {
	return checkHash(deployment, deployment.Status.Phase, storage, log)
}

func checkHash(deployment coh.CoherenceResource, phase coh.ConditionType, storage utils.Storage, log logr.Logger) (coh.Resources, error) {
	hash := deployment.GetGenerationString()
	var desiredResources coh.Resources
	var err error

	storeHash := storage.GetHash()
	if storeHash == "" || storeHash != hash || phase != coh.ConditionTypeReady {
		// Storage state was saved with no hash or a different hash so is not in the desired state
		// or the Coherence resource is not in the Ready state
		// Create the desired resources the deployment
		if desiredResources, err = deployment.CreateKubernetesResources(); err != nil {
			return desiredResources, err
		}

		if storeHash != "" {
			// The "storeHash" is not "", so it must have been processed by the Operator (could have been a previous version).
			// Prior to 3.5.0 the hash was calculated based on the resource spec but this was unreliable and since 3.5.0
			// the Coherence resource metadata generation is used instead of a hash.
			if deployment.IsBeforeVersion("3.5.0") {
				desiredResources = storage.GetLatest()
				log.Info("Ignoring hash difference for pre-3.5.0 resource", "hash", hash, "store", storeHash)
			}

		}
	} else {
		// storage state was saved with the current hash so is already in the desired state
		desiredResources = storage.GetLatest()
	}
	return desiredResources, err
}
