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

	storeHash, hashFound := storage.GetHash()
	if !hashFound || storeHash != hash || phase != coh.ConditionTypeReady {
		// Storage state was saved with no hash or a different hash so is not in the desired state
		// or the Coherence resource is not in the Ready state
		// Create the desired resources the deployment
		if desiredResources, err = deployment.CreateKubernetesResources(); err != nil {
			// there was an error creating the resources
			return desiredResources, err
		}

		if hashFound {
			// The "storeHash" is not "", so it must have been previously processed by the Operator (could have been a previous version).
			// The operator now uses the Coherence resource's generation instead of calculating a hash, so if the version is
			// 3.4.3 or earlier and hash difference is ignored.
			if deployment.IsBeforeOrSameVersion("3.4.3") {
				// There is an edge case where the Coherence resource could have legitimately been updated whilst
				// the Operator and web-hooks were uninstalled. In that case we would ignore the update until another
				// update is made. The simplest way for the customer to work around this is to add the
				// AnnotationOperatorVersion annotation with some value, which will then be overwritten by the web-hook
				// and the Coherence resource will be correctly processes.
				desiredResources = storage.GetLatest()
				log.Info("Ignoring hash difference for 3.4.3 or earlier resource", "hash", hash, "store", storeHash)
			}
		}
	} else {
		// storage state was saved with the current hash so is already in the desired state
		desiredResources = storage.GetLatest()
	}
	return desiredResources, err
}
