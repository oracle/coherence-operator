# ----------------------------------------------------------------------------------------------------------------------
# Copyright (c) 2025, Oracle and/or its affiliates.
#
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
CURRDIR         	 	:= $(shell pwd)
NAMESPACE        		 ?= coherence-demo
IMAGE_PREFIX             ?=
IMAGE_VERSION            ?= 1.0.0
PLATFORM                 ?= linux/arm64
GO_IMAGE                 ?= $(IMAGE_PREFIX)polyglot-client-go:$(IMAGE_VERSION)
PY_IMAGE                 ?= $(IMAGE_PREFIX)polyglot-client-py:$(IMAGE_VERSION)
JS_IMAGE                 ?= $(IMAGE_PREFIX)polyglot-client-js:$(IMAGE_VERSION)


# ======================================================================================================================
# Makefile targets start here
# ======================================================================================================================

# ----------------------------------------------------------------------------------------------------------------------
# Display the Makefile help - this is a list of the targets with a description.
# This target MUST be the first target in the Makefile so that it is run when running make with no arguments
# ----------------------------------------------------------------------------------------------------------------------
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


# ======================================================================================================================
# Build targets
# ======================================================================================================================
##@ Build

.PHONY: create-namespace
create-namespace: ## Create the namespace
	kubectl create namespace $(NAMESPACE) || true

.PHONY: delete-namespace
delete-namespace: ## Delete the namespace
	kubectl delete namespace $(NAMESPACE) || true

.PHONY: create-go-image
create-go-image: ## Create the Go Docker image
	cd go && docker buildx build --platform $(PLATFORM) -t $(GO_IMAGE) .

.PHONY: create-py-image
create-py-image: ## Create the Python Docker image
	cd py && docker buildx build --platform $(PLATFORM) -t $(PY_IMAGE) .

.PHONY: create-js-image
create-js-image: ## Create the JavaScript Docker image
	cd js && docker buildx build --platform $(PLATFORM) -t $(JS_IMAGE) .

.PHONY: kind-load-images
kind-load-images: ## Load images into kind
	kind load docker-image $(GO_IMAGE)
	kind load docker-image $(PY_IMAGE)
	kind load docker-image $(JS_IMAGE)

.PHONY: create-all-images
create-all-images: create-go-image create-js-image create-py-image ## Create all images

# ======================================================================================================================
# Deploy targets
# ======================================================================================================================
##@ Deploy

.PHONY: deploy-operator
deploy-operator: ## Deploy the Coherence Operator
	kubectl apply -f https://github.com/oracle/coherence-operator/releases/download/v3.4.3/coherence-operator.yaml

.PHONY: undeploy-operator
undeploy-operator: ## Undeploy the Coherence Operator
	kubectl delete -f https://github.com/oracle/coherence-operator/releases/download/v3.4.3/coherence-operator.yaml

.PHONY: deploy-coherence
deploy-coherence: ## Deploy the Coherence Cluster
	kubectl -n $(NAMESPACE) apply -f yaml/coherence-cluster.yaml
	sleep 5
	kubectl -n $(NAMESPACE) get pods

.PHONY: undeploy-coherence
undeploy-coherence: ## Deploy the Coherence Cluster
	kubectl -n $(NAMESPACE) delete -f yaml/coherence-cluster.yaml

.PHONY: deploy-go-client
deploy-go-client: ## Deploy the Go client
	kubectl -n $(NAMESPACE) apply -f yaml/go-client.yaml

.PHONY: undeploy-go-client
undeploy-go-client: ## Undeploy the Go client
	kubectl -n $(NAMESPACE) delete -f yaml/go-client.yaml

.PHONY: deploy-py-client
deploy-py-client: ## Deploy the Python client
	kubectl -n $(NAMESPACE) apply -f yaml/py-client.yaml

.PHONY: undeploy-py-client
undeploy-py-client: ## Undeploy the Python client
	kubectl -n $(NAMESPACE) delete -f yaml/py-client.yaml

.PHONY: deploy-js-client
deploy-js-client: ## Deploy the JavaScript client
	kubectl -n $(NAMESPACE) apply -f yaml/js-client.yaml

.PHONY: undeploy-js-client
undeploy-js-client: ## Deploy the JavaScript client
	kubectl -n $(NAMESPACE) delete -f yaml/js-client.yaml

.PHONY: deploy-all-clients
deploy-all-clients: deploy-go-client deploy-js-client deploy-py-client ## Deploy all clients

.PHONY: undeploy-all-clients
undeploy-all-clients: undeploy-go-client undeploy-js-client undeploy-py-client ## Undeploy all clients
