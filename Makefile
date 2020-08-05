# ----------------------------------------------------------------------------------------------------------------------
# Copyright (c) 2019, 2020, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
# ----------------------------------------------------------------------------------------------------------------------
# This is the Makefile to build the Coherence Kubernetes Operator.
# ----------------------------------------------------------------------------------------------------------------------

# The version of the Operator being build - this should be a valid SemVer format
VERSION ?= 3.1.0

# The operator version to use to run certification tests against
CERTIFICATION_VERSION ?= $(VERSION)

# A SPACE delimited list of previous Operator versions that are used to run the compatibility tests.
COMPATIBLE_VERSIONS = 3.0.0

# Capture the Git commit to add to the build information
GITCOMMIT       ?= $(shell git rev-list -1 HEAD)
GITREPO         := https://github.com/oracle/coherence-operator.git
BUILD_DATE      := $(shell date -u | tr ' ' '.')
BUILD_INFO      := "$(VERSION)|$(GITCOMMIT)|$(BUILD_DATE)"

CURRDIR         := $(shell pwd)

ARCH            ?= amd64
OS              ?= linux
UNAME_S         := $(shell uname -s)
GOPROXY         ?= https://proxy.golang.org

# Set the location of the Operator SDK executable
UNAME_S               = $(shell uname -s)
UNAME_M               = $(shell uname -m)
OPERATOR_SDK_VERSION := v0.19.0
OPERATOR_SDK          = $(CURRDIR)/etc/sdk/$(UNAME_S)-$(UNAME_M)/operator-sdk

# The Coherence image to use for deployments that do not specify an image
COHERENCE_IMAGE   ?= oraclecoherence/coherence-ce:14.1.1-0-1
# This is the Coherence image that will be used in tests.
# Changing this variable will allow test builds to be run against different Coherence versions
# without altering the default image name.
TEST_COHERENCE_IMAGE ?= $(COHERENCE_IMAGE)

# Operator image names
RELEASE_IMAGE_PREFIX   ?= container-registry.oracle.com/middleware/
OPERATOR_IMAGE_REPO    := $(RELEASE_IMAGE_PREFIX)coherence-operator
OPERATOR_IMAGE         := $(OPERATOR_IMAGE_REPO):$(VERSION)
UTILS_IMAGE            ?= $(OPERATOR_IMAGE_REPO):$(VERSION)-utils
# The Operator images to push
OPERATOR_RELEASE_REPO  ?= $(OPERATOR_IMAGE_REPO)
OPERATOR_RELEASE_IMAGE := $(OPERATOR_RELEASE_REPO):$(VERSION)
UTILS_RELEASE_IMAGE    := $(OPERATOR_RELEASE_REPO):$(VERSION)-utils

# The test application image used in integration tests
TEST_APPLICATION_IMAGE := $(RELEASE_IMAGE_PREFIX)operator-test-jib:$(VERSION)

# Default bundle image tag
BUNDLE_IMG ?= $(OPERATOR_IMAGE_REPO):$(VERSION)-bundle
# Options for 'bundle-build'
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# Release build options
RELEASE_DRY_RUN  ?= true
PRE_RELEASE      ?= true

# Extra arguments to pass to the go test command for the various test steps.
# For example, when running make e2e-test we can run just a single test such
# as the zone test using the go test -run=regex argument like this
#   make e2e-test GO_TEST_FLAGS_E2E='-run=^TestZone$$'
GO_TEST_FLAGS     ?= -timeout=20m
GO_TEST_FLAGS_E2E := -timeout=100m

# default as in test/e2e/helper/proj_helpers.go
TEST_NAMESPACE ?= operator-test
# flag indicating whether the test namespace should be reset (deleted and recreated) before tests
CREATE_TEST_NAMESPACE ?= true

# Prometheus Operator settings (used in integration tests)
PROMETHEUS_INCLUDE_GRAFANA   ?= true
PROMETHEUS_OPERATOR_VERSION  ?= 8.13.7
GRAFANA_DASHBOARDS           ?= dashboards/grafana-legacy/

# Elasticsearch & Kibana settings (used in integration tests)
ELASTIC_VERSION ?= 7.6.2
KIBANA_INDEX_PATTERN := "6abb1220-3feb-11e9-a9a3-4b1c09db6e6a"

# restart local storage for persistence
LOCAL_STORAGE_RESTART ?= false

# Env variables used to create pull secrets
DOCKER_SERVER ?=
DOCKER_USERNAME ?=
DOCKER_PASSWORD ?=
OCR_DOCKER_USERNAME ?=
OCR_DOCKER_PASSWORD ?=
MAVEN_USER ?=
MAVEN_PASSWORD ?=


ifneq ("$(MAVEN_SETTINGS)","")
	USE_MAVEN_SETTINGS = -s $(MAVEN_SETTINGS)
else
	USE_MAVEN_SETTINGS =
endif

# Configure the image pull secrets information
ifneq ("$(or $(DOCKER_USERNAME),$(DOCKER_PASSWORD))","")
DOCKER_SECRET = coherence-k8s-operator-development-secret
else
DOCKER_SECRET =
endif
ifneq ("$(or $(OCR_DOCKER_USERNAME),$(OCR_DOCKER_PASSWORD))","")
OCR_DOCKER_SECRET = ocr-k8s-operator-development-secret
else
OCR_DOCKER_SECRET =
endif

ifneq ("$(and $(DOCKER_SECRET),$(OCR_DOCKER_SECRET))","")
IMAGE_PULL_SECRETS ?= $(DOCKER_SECRET),$(OCR_DOCKER_SECRET)
else
ifneq ("$(DOCKER_SECRET)","")
IMAGE_PULL_SECRETS ?= $(DOCKER_SECRET)
else
ifneq ("$(OCR_DOCKER_SECRET)","")
IMAGE_PULL_SECRETS ?= $(OCR_DOCKER_SECRET)
else
IMAGE_PULL_SECRETS ?=
endif
endif
endif

IMAGE_PULL_POLICY  ?= IfNotPresent

# Env variable used by the kubectl test framework to locate the kubectl binary
TEST_ASSET_KUBECTL ?= $(shell which kubectl)

override BUILD_OUTPUT      := ./build/_output
override BUILD_BIN         := ./bin
override BUILD_CONFIG      := $(BUILD_OUTPUT)/config
override BUILD_ASSETS      := $(BUILD_OUTPUT)/assets
override BUILD_TARGETS     := $(BUILD_OUTPUT)/targets
override BUILD_PROPS       := $(BUILD_OUTPUT)/build.properties
override TEST_LOGS_DIR     := $(BUILD_OUTPUT)/test-logs

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GOS          = $(shell find . -type f -name "*.go" ! -name "*_test.go")
OPTESTGOS    = $(shell find cmd/optest -type f -name "*.go" ! -name "*_test.go")
API_GO_FILES = $(shell find api -type f -name "*.go" ! -name "*_test.go"  ! -name "zz*.go")
CRD_V1       ?= $(shell kubectl api-versions | grep '^apiextensions.k8s.io/v1$$')

TEST_SSL_SECRET := coherence-ssl-secret

.PHONY: all
all: build-all-images

# ----------------------------------------------------------------------------------------------------------------------
# Configure the build properties
# ----------------------------------------------------------------------------------------------------------------------
$(BUILD_PROPS):
	# Ensures that build output directories exist
	@echo "Creating build directories"
	@mkdir -p $(BUILD_OUTPUT)
	@mkdir -p $(BUILD_BIN)
	@mkdir -p $(BUILD_ASSETS)
	@mkdir -p $(BUILD_TARGETS)
	@mkdir -p $(TEST_LOGS_DIR)
	# create build.properties
	rm -f $(BUILD_PROPS)
	printf "COHERENCE_IMAGE=$(COHERENCE_IMAGE)\n\
	UTILS_IMAGE=$(UTILS_IMAGE)\n\
	OPERATOR_IMAGE=$(OPERATOR_IMAGE)\n\
	VERSION=$(VERSION)\n" > $(BUILD_PROPS)

# ----------------------------------------------------------------------------------------------------------------------
# Builds the Operator
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-operator
build-operator: $(BUILD_TARGETS)/build-operator

# ----------------------------------------------------------------------------------------------------------------------
# Build the Operator Image
# ----------------------------------------------------------------------------------------------------------------------
#   We copy the Dockerfile to $(BUILD_OUTPUT) only so that we can use it as a conditional build dependency in this Makefile
$(BUILD_TARGETS)/build-operator: $(BUILD_BIN)/manager $(BUILD_BIN)/runner
	docker build --build-arg version=$(VERSION) \
		--build-arg coherence_image=$(COHERENCE_IMAGE) \
		--build-arg utils_image=$(UTILS_IMAGE) \
		. -t $(OPERATOR_IMAGE)
	touch $(BUILD_TARGETS)/build-operator

# ----------------------------------------------------------------------------------------------------------------------
# Build the operator linux binary
# ----------------------------------------------------------------------------------------------------------------------
$(BUILD_BIN)/manager: $(BUILD_PROPS) $(GOS) $(BUILD_TARGETS)/generate $(BUILD_TARGETS)/manifests
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags -X=main.BuildInfo=$BuildInfo -a -o $(BUILD_BIN)/manager main.go

# ----------------------------------------------------------------------------------------------------------------------
# Ensure Operator SDK is at the correct version
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: ensure-sdk
ensure-sdk:
	@echo "Ensuring Operator SDK is present"
	./hack/ensure-sdk.sh $(OPERATOR_SDK_VERSION)

# ----------------------------------------------------------------------------------------------------------------------
# Internal make step that builds the Operator runner artifacts utility
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-runner
build-runner: $(BUILD_BIN)/runner

$(BUILD_BIN)/runner: $(BUILD_PROPS) $(GOS)
	@echo "Building Operator Runner"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags -X=main.BuildInfo=$(BUILD_INFO) -o $(BUILD_BIN)/runner ./cmd/runner

# ----------------------------------------------------------------------------------------------------------------------
# Internal make step that builds the Operator legacy converter
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: converter
converter: $(BUILD_BIN)/converter $(BUILD_BIN)/converter-linux-amd64 $(BUILD_BIN)/converter-darwin-amd64 $(BUILD_BIN)/converter-windows-amd64

$(BUILD_BIN)/converter: $(BUILD_PROPS) $(GOS)
	CGO_ENABLED=0 GO111MODULE=on GOOS=$(OS) GOARCH=$(ARCH) go build -o $(BUILD_BIN)/converter ./cmd/converter

$(BUILD_BIN)/converter-linux-amd64: $(BUILD_PROPS) $(GOS)
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o $(BUILD_BIN)/converter-linux-amd64 ./cmd/converter

$(BUILD_BIN)/converter-darwin-amd64: $(BUILD_PROPS) $(GOS)
	CGO_ENABLED=0 GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o $(BUILD_BIN)/converter-darwin-amd64 ./cmd/converter

$(BUILD_BIN)/converter-windows-amd64: $(BUILD_PROPS) $(GOS)
	CGO_ENABLED=0 GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o $(BUILD_BIN)/converter-windows-amd64 ./cmd/converter

# ----------------------------------------------------------------------------------------------------------------------
# Internal make step that builds the Operator test utility
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-op-test
build-op-test: $(BUILD_BIN)/op-test

$(BUILD_BIN)/op-test: $(BUILD_PROPS) $(GOS) $(OPTESTGOS)
	CGO_ENABLED=0 GO111MODULE=on GOOS=$(OS) GOARCH=$(ARCH) go build -o $(BUILD_BIN)/op-test ./cmd/optest

# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go unit tests that do not require a k8s cluster
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: test-operator
test-operator: export CGO_ENABLED = 0
test-operator: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
test-operator: export UTILS_IMAGE := $(UTILS_IMAGE)
test-operator: $(BUILD_TARGETS)/build-operator gotestsum
	@echo "Running operator tests"
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-test.xml \
	  -- $(GO_TEST_FLAGS) -v ./api/... ./controllers/... ./cmd/... ./pkg/...

# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end tests that require a k8s cluster using
# a LOCAL operator instance (i.e. the operator is not deployed to k8s).
# These tests will use whichever k8s cluster the local environment
# is pointing to.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: e2e-local-test
e2e-local-test: export CGO_ENABLED = 0
e2e-local-test: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
e2e-local-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
e2e-local-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
e2e-local-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
e2e-local-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
e2e-local-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
e2e-local-test: export TEST_ASSET_KUBECTL := $(TEST_ASSET_KUBECTL)
e2e-local-test: export LOCAL_STORAGE_RESTART := $(LOCAL_STORAGE_RESTART)
e2e-local-test: export VERSION := $(VERSION)
e2e-local-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
e2e-local-test: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
e2e-local-test: export UTILS_IMAGE := $(UTILS_IMAGE)
e2e-local-test: $(BUILD_TARGETS)/build-operator reset-namespace create-ssl-secrets install-crds gotestsum
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-local-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/e2e/local/...

# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end tests that require a k8s cluster using
# a DEPLOYED operator instance (i.e. the operator Docker image is
# deployed to k8s). These tests will use whichever k8s cluster the
# local environment is pointing to.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: e2e-test
e2e-test: $(BUILD_TARGETS)/build-operator reset-namespace create-ssl-secrets uninstall-crds deploy
	$(MAKE) $(MAKEFLAGS) run-e2e-test \
	; rc=$$? \
	; echo "E2E Tests completed with return code $$rc" \
	; $(MAKE) $(MAKEFLAGS) undeploy \
	; $(MAKE) $(MAKEFLAGS) delete-namespace \
	; exit $$rc

.PHONY: run-e2e-test
run-e2e-test: export CGO_ENABLED = 0
run-e2e-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
run-e2e-test: export TEST_NAMESPACE := $(TEST_NAMESPACE)
run-e2e-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
run-e2e-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
run-e2e-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
run-e2e-test: export TEST_ASSET_KUBECTL := $(TEST_ASSET_KUBECTL)
run-e2e-test: export VERSION := $(VERSION)
run-e2e-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
run-e2e-test: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
run-e2e-test: export UTILS_IMAGE := $(UTILS_IMAGE)
run-e2e-test: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
run-e2e-test: gotestsum
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/e2e/remote/...

# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end tests that require Prometheus in the k8s cluster
# using a LOCAL operator instance (i.e. the operator is not deployed to k8s).
#
# This target DOES NOT install Prometheus, use the e2e-prometheus-test target
# to fully reset the test namespace.
#
# These tests will use whichever k8s cluster the local environment
# is pointing to.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: run-prometheus-test
run-prometheus-test: export CGO_ENABLED = 0
run-prometheus-test: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
run-prometheus-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
run-prometheus-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
run-prometheus-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
run-prometheus-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
run-prometheus-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
run-prometheus-test: export TEST_ASSET_KUBECTL := $(TEST_ASSET_KUBECTL)
run-prometheus-test: export LOCAL_STORAGE_RESTART := $(LOCAL_STORAGE_RESTART)
run-prometheus-test: export VERSION := $(VERSION)
run-prometheus-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
run-prometheus-test: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
run-prometheus-test: export UTILS_IMAGE := $(UTILS_IMAGE)
run-prometheus-test: $(BUILD_TARGETS)/build-operator create-ssl-secrets gotestsum
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-prometheus-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/e2e/prometheus/... \
	  2>&1 | tee $(TEST_LOGS_DIR)/operator-e2e-prometheus-test.out

.PHONY: e2e-prometheus-test
e2e-prometheus-test: reset-namespace install-prometheus
	$(MAKE) $(MAKEFLAGS) run-prometheus-test \
	; rc=$$? \
	; $(MAKE) $(MAKEFLAGS) uninstall-prometheus \
	; $(MAKE) $(MAKEFLAGS) delete-namespace \
	; exit $$rc


# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end tests that require Elasticsearch in the k8s cluster
# using a LOCAL operator instance (i.e. the operator is not deployed to k8s).
#
# This target DOES NOT install Elasticsearch, use the e2e-elastic-test target
# to fully reset the test namespace.
#
# These tests will use whichever k8s cluster the local environment
# is pointing to.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: run-elastic-test
run-elastic-test: export CGO_ENABLED = 0
run-elastic-test: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
run-elastic-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
run-elastic-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
run-elastic-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
run-elastic-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
run-elastic-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
run-elastic-test: export TEST_ASSET_KUBECTL := $(TEST_ASSET_KUBECTL)
run-elastic-test: export LOCAL_STORAGE_RESTART := $(LOCAL_STORAGE_RESTART)
run-elastic-test: export VERSION := $(VERSION)
run-elastic-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
run-elastic-test: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
run-elastic-test: export UTILS_IMAGE := $(UTILS_IMAGE)
run-elastic-test: export KIBANA_INDEX_PATTERN := $(KIBANA_INDEX_PATTERN)
run-elastic-test: $(BUILD_TARGETS)/build-operator create-ssl-secrets gotestsum
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-elastic-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/e2e/elastic/... \
	  2>&1 | tee $(TEST_LOGS_DIR)/operator-e2e-elastic-test.out

.PHONY: e2e-elastic-test
e2e-elastic-test: reset-namespace install-elastic
	$(MAKE) $(MAKEFLAGS) run-elastic-test \
	; rc=$$? \
	; $(MAKE) $(MAKEFLAGS) uninstall-elastic \
	; $(MAKE) $(MAKEFLAGS) delete-namespace \
	; exit $$rc

# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end Operator Compatibility tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: compatibility-test
compatibility-test: export CGO_ENABLED = 0
compatibility-test: export TEST_NAMESPACE := $(TEST_NAMESPACE)
compatibility-test: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
compatibility-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
compatibility-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
compatibility-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
compatibility-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
compatibility-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
compatibility-test: export VERSION := $(VERSION)
compatibility-test: export COMPATIBLE_VERSIONS := $(COMPATIBLE_VERSIONS)
compatibility-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
compatibility-test: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
compatibility-test: export UTILS_IMAGE := $(UTILS_IMAGE)
compatibility-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
compatibility-test: $(BUILD_TARGETS)/build-operator clean-namespace reset-namespace create-ssl-secrets gotestsum
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-compatibility-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/compatibility/... \
	  2>&1 | tee $(TEST_LOGS_DIR)/operator-e2e-compatibility-test.out


# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end Operator certification tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# Note that the namespace will be created if it does not exist.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: certification-test
certification-test: export CGO_ENABLED = 0
certification-test: export TEST_NAMESPACE := $(TEST_NAMESPACE)
certification-test: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
certification-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
certification-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
certification-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
certification-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
certification-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
certification-test: export VERSION := $(VERSION)
certification-test: export CERTIFICATION_VERSION := $(CERTIFICATION_VERSION)
certification-test: export OPERATOR_IMAGE_REPO := $(OPERATOR_IMAGE_REPO)
certification-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
certification-test: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
certification-test: export UTILS_IMAGE := $(UTILS_IMAGE)
certification-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
certification-test: install-certification
	$(MAKE) $(MAKEFLAGS) run-certification \
	; rc=$$? \
	; $(MAKE) $(MAKEFLAGS) cleanup-certification \
	; exit $$rc


# ----------------------------------------------------------------------------------------------------------------------
# Install the Operator prior to running compatability tests.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-certification
install-certification: export TEST_NAMESPACE := $(TEST_NAMESPACE)
install-certification: export VERSION := $(VERSION)
install-certification: export VERSION := $(VERSION)
install-certification: export CERTIFICATION_VERSION := $(CERTIFICATION_VERSION)
install-certification: $(BUILD_TARGETS)/build-operator reset-namespace create-ssl-secrets
ifeq ("$(CERTIFICATION_VERSION)","$(VERSION)")
	$(MAKE) deploy
#else
#	helm repo add coherence https://oracle.github.io/coherence-operator/charts || true
#	helm repo update || true
#	helm install --atomic --namespace $(TEST_NAMESPACE) --wait --version $(CERTIFICATION_VERSION) operator ./helm-charts/coherence-operator
endif

# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end Operator certification tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# Note that the namespace will be created if it does not exist.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: run-certification
run-certification: export CGO_ENABLED = 0
run-certification: export TEST_NAMESPACE := $(TEST_NAMESPACE)
run-certification: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
run-certification: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
run-certification: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
run-certification: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
run-certification: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
run-certification: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
run-certification: export VERSION := $(VERSION)
run-certification: export CERTIFICATION_VERSION := $(CERTIFICATION_VERSION)
run-certification: export OPERATOR_IMAGE_REPO := $(OPERATOR_IMAGE_REPO)
run-certification: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
run-certification: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
run-certification: export UTILS_IMAGE := $(UTILS_IMAGE)
run-certification: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
run-certification: gotestsum
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-certification-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/certification/... \
	  2>&1 | tee $(TEST_LOGS_DIR)/operator-e2e-certification-test.out

# ----------------------------------------------------------------------------------------------------------------------
# Clean up after to running compatability tests.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: cleanup-certification
cleanup-certification: export TEST_NAMESPACE := $(TEST_NAMESPACE)
cleanup-certification:
	$(MAKE) deploy
	$(MAKE) uninstall-crds
	$(MAKE) delete-namespace


# ----------------------------------------------------------------------------------------------------------------------
# Install CRDs into Kubernetes.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-crds
install-crds: uninstall-crds $(BUILD_TARGETS)/manifests $(GOBIN)/kustomize
ifeq ("$(CRD_V1)","apiextensions.k8s.io/v1")
	$(GOBIN)/kustomize build config/crd | kubectl create -f -
else
	$(GOBIN)/kustomize build config/crd-v1beta1 | kubectl create -f -
endif

# ----------------------------------------------------------------------------------------------------------------------
# Uninstall CRDs from Kubernetes.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: uninstall-crds
uninstall-crds: $(BUILD_TARGETS)/manifests $(GOBIN)/kustomize
ifeq ("$(CRD_V1)","apiextensions.k8s.io/v1")
	$(GOBIN)/kustomize build config/crd | kubectl delete -f - || true
else
	$(GOBIN)/kustomize build config/crd-v1beta1 | kubectl delete -f -
endif

# ----------------------------------------------------------------------------------------------------------------------
# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: deploy
deploy: $(BUILD_TARGETS)/manifests $(GOBIN)/kustomize
	cp -R config/ $(BUILD_CONFIG)
#   Uncomment to watch a single namespace
#	cd $(BUILD_CONFIG)/manager && $(GOBIN)/kustomize edit add configmap env-vars --from-literal WATCH_NAMESPACE=$(TEST_NAMESPACE)
	cd $(BUILD_CONFIG)/default && $(GOBIN)/kustomize edit set namespace $(TEST_NAMESPACE)
	cd $(BUILD_CONFIG)/manager && $(GOBIN)/kustomize edit set image controller=$(OPERATOR_IMAGE)
	$(GOBIN)/kustomize build $(BUILD_CONFIG)/default | kubectl create -f -

# ----------------------------------------------------------------------------------------------------------------------
# Un-deploy controller from the configured Kubernetes cluster in ~/.kube/config
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: undeploy
undeploy: $(BUILD_TARGETS)/manifests $(GOBIN)/kustomize
	cp -R config/ $(BUILD_CONFIG)
	cd $(BUILD_CONFIG)/default && $(GOBIN)/kustomize edit add configmap source-vars --from-literal OPERATOR_NAMESPACE=$(TEST_NAMESPACE)
	cd $(BUILD_CONFIG)/default && $(GOBIN)/kustomize edit set namespace $(TEST_NAMESPACE)
	cd $(BUILD_CONFIG)/manager && $(GOBIN)/kustomize edit set image controller=$(OPERATOR_IMAGE)
	$(GOBIN)/kustomize build $(BUILD_CONFIG)/default | kubectl delete -f -

# ----------------------------------------------------------------------------------------------------------------------
# Generate manifests e.g. CRD, RBAC etc.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: manifests
manifests: $(BUILD_TARGETS)/manifests

$(BUILD_TARGETS)/manifests: $(BUILD_PROPS) config/crd/bases/coherence.oracle.com_coherences.yaml config/crd-v1beta1/bases/coherence.oracle.com_coherences.yaml docs/about/04_coherence_spec.adoc
	touch $(BUILD_TARGETS)/manifests

config/crd/bases/coherence.oracle.com_coherences.yaml: $(API_GO_FILES) $(GOBIN)/controller-gen
	$(GOBIN)/controller-gen "crd:trivialVersions=true,crdVersions={v1}" \
	  rbac:roleName=manager-role webhook paths="{./api/...,./controllers/...}" \
	  output:crd:artifacts:config=config/crd/bases

config/crd-v1beta1/bases/coherence.oracle.com_coherences.yaml: $(API_GO_FILES) $(GOBIN)/controller-gen
	@echo "Generating CRD v1beta1"
	cp -R config/crd/bases config/crd-v1beta1/bases
	$(GOBIN)/controller-gen "crd:trivialVersions=true,crdVersions={v1beta1}" \
	  rbac:roleName=manager-role webhook paths="{./api/...,./controllers/...}" \
	  output:crd:artifacts:config=config/crd-v1beta1/bases

# ----------------------------------------------------------------------------------------------------------------------
# Generate the data.json file used by the Operator for default configuration values
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: generate-config
generate-config:  $(BUILD_PROPS)
	@echo "Generating Operator config"
	@printf "{\n\
	  \"CoherenceImage\": \"$(COHERENCE_IMAGE)\",\n\
	  \"UtilsImage\": \"$(UTILS_RELEASE_IMAGE)\"\n\
	}" > config/operator/new-data.json
# If the new file is different to the old file replace the old with the new
# This ensures that Git only thinks there is a file update if ghe contents have actually changed
	@if ! diff config/operator/new-data.json config/operator/data.json; then \
	  cp config/operator/new-data.json config/operator/data.json ; \
	fi
	rm config/operator/new-data.json \

# ----------------------------------------------------------------------------------------------------------------------
# Generate code, configuration and docs.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: generate
generate: $(BUILD_TARGETS)/generate
	touch $(BUILD_TARGETS)/generate

$(BUILD_TARGETS)/generate: $(BUILD_PROPS) api/v1/zz_generated.deepcopy.go pkg/data/zz_generated_assets.go
	touch $(BUILD_TARGETS)/generate

api/v1/zz_generated.deepcopy.go: $(API_GO_FILES) $(GOBIN)/controller-gen
	$(GOBIN)/controller-gen object:headerFile="./hack/boilerplate.go.txt" paths="./api/..."

pkg/data/zz_generated_assets.go: config/operator/data.json config/crd/bases/coherence.oracle.com_coherences.yaml config/crd-v1beta1/bases/coherence.oracle.com_coherences.yaml $(GOBIN)/kustomize
	echo "Embedding configuration and CRD files"
	cp config/operator/data.json $(BUILD_ASSETS)/data.json
	echo "Embedding v1 CRD files"
	$(GOBIN)/kustomize build config/crd > $(BUILD_ASSETS)/crd_v1.yaml
	echo "Embedding v1beat1 CRD files"
	$(GOBIN)/kustomize build config/crd-v1beta1 > $(BUILD_ASSETS)/crd_v1beta1.yaml
	go get -u github.com/shurcooL/vfsgen
	go run ./pkg/generate/assets_generate.go

# ----------------------------------------------------------------------------------------------------------------------
# find or download controller-gen
# ----------------------------------------------------------------------------------------------------------------------
$(GOBIN)/controller-gen:
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.3.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}

# ----------------------------------------------------------------------------------------------------------------------
# find or download kustomize
# ----------------------------------------------------------------------------------------------------------------------
$(GOBIN)/kustomize:
	@{ \
	set -e ;\
	KUSTOMIZE_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$KUSTOMIZE_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/kustomize/kustomize/v3@v3.5.4 ;\
	rm -rf $$KUSTOMIZE_GEN_TMP_DIR ;\
	}
	KUSTOMIZE=$(GOBIN)/kustomize

# ----------------------------------------------------------------------------------------------------------------------
# find or download gotestsum
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: gotestsum
gotestsum:
ifeq (, $(shell which gotestsum))
	@{ \
	set -e ;\
	GOTESTSUM_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$GOTESTSUM_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get gotest.tools/gotestsum@v0.5.2 ;\
	rm -rf $$GOTESTSUM_GEN_TMP_DIR ;\
	}
GOTESTSUM=$(GOBIN)/gotestsum
else
GOTESTSUM=$(shell which gotestsum)
endif

# ----------------------------------------------------------------------------------------------------------------------
# Generate bundle manifests and metadata, then validate generated files.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: bundle
bundle: $(BUILD_TARGETS)/manifests
	$(OPERATOR_SDK) generate $(GOBIN)/kustomize manifests -q
	kustomize build config/manifests | $(OPERATOR_SDK) generate bundle -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)
	$(OPERATOR_SDK) bundle validate ./bundle

# ----------------------------------------------------------------------------------------------------------------------
# Build the bundle image.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: bundle-build
bundle-build:
	docker build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

# ----------------------------------------------------------------------------------------------------------------------
# Generate API docs
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: api-doc-gen
api-doc-gen: docs/about/04_coherence_spec.adoc

docs/about/04_coherence_spec.adoc: $(API_GO_FILES)
	@echo "Generating CRD Doc"
	go run ./cmd/docgen/ \
		api/v1/coherenceresourcespec_types.go \
		api/v1/coherence_types.go \
		api/v1/coherenceresource_types.go \
		> docs/about/04_coherence_spec.adoc

# ----------------------------------------------------------------------------------------------------------------------
# Clean-up all of the build artifacts
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: clean
clean:
	-rm -rf build/_output
	-rm -f bin/*
	mvn $(USE_MAVEN_SETTINGS) -f java clean
	mvn $(USE_MAVEN_SETTINGS) -f examples clean

# ----------------------------------------------------------------------------------------------------------------------
# Generate the keys and certs used in tests.
# ----------------------------------------------------------------------------------------------------------------------
$(BUILD_OUTPUT)/certs:
	@echo "Generating test keys and certs"
	./hack/keys.sh

# ----------------------------------------------------------------------------------------------------------------------
# Delete and re-create the test namespace
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: reset-namespace
reset-namespace: export KUBECONFIG_PATH := $(KUBECONFIG_PATH)
reset-namespace: export DOCKER_SERVER := $(DOCKER_SERVER)
reset-namespace: export DOCKER_USERNAME := $(DOCKER_USERNAME)
reset-namespace: export DOCKER_PASSWORD := $(DOCKER_PASSWORD)
reset-namespace: export OCR_DOCKER_USERNAME := $(OCR_DOCKER_USERNAME)
reset-namespace: export OCR_DOCKER_PASSWORD := $(OCR_DOCKER_PASSWORD)
reset-namespace: delete-namespace
ifeq ($(CREATE_TEST_NAMESPACE),true)
	@echo "Creating test namespace $(TEST_NAMESPACE)"
	kubectl create namespace $(TEST_NAMESPACE)
endif
ifneq ($(DOCKER_SERVER),)
	@echo "Creating pull secrets for $(DOCKER_SERVER)"
	kubectl create secret docker-registry coherence-k8s-operator-development-secret \
								--namespace $(TEST_NAMESPACE) \
								--docker-server "$(DOCKER_SERVER)" \
								--docker-username "$(DOCKER_USERNAME)" \
								--docker-password "$(DOCKER_PASSWORD)" \
								--docker-email="docker@dummy.com"
endif
ifneq ("$(or $(OCR_DOCKER_USERNAME),$(OCR_DOCKER_PASSWORD))","")
	@echo "Creating pull secrets for container-registry.oracle.com"
	kubectl create secret docker-registry ocr-k8s-operator-development-secret \
								--namespace $(TEST_NAMESPACE) \
								--docker-server container-registry.oracle.com \
								--docker-username "$(OCR_DOCKER_USERNAME)" \
								--docker-password "$(OCR_DOCKER_PASSWORD)" \
								--docker-email "docker@dummy.com"
endif

# ----------------------------------------------------------------------------------------------------------------------
# Delete the test namespace
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: delete-namespace
delete-namespace: clean-namespace
ifeq ($(CREATE_TEST_NAMESPACE),true)
	@echo "Deleting test namespace $(TEST_NAMESPACE)"
	kubectl delete namespace $(TEST_NAMESPACE) --force --grace-period=0 && echo "deleted namespace" || true
endif
	kubectl delete clusterrole operator-test-coherence-operator --force --grace-period=0 && echo "deleted namespace" || true
	kubectl delete clusterrolebinding operator-test-coherence-operator --force --grace-period=0 && echo "deleted namespace" || true


# ----------------------------------------------------------------------------------------------------------------------
# Delete all resource from the test namespace
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: clean-namespace
clean-namespace: delete-coherence-clusters
	for i in $$(kubectl -n $(TEST_NAMESPACE) get all -o name); do \
		echo "Deleting $${i} from test namespace $(TEST_NAMESPACE)" \
		kubectl -n $(TEST_NAMESPACE) delete $${i}; \
	done

# ----------------------------------------------------------------------------------------------------------------------
# Create the k8s secret to use in SSL/TLS testing.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: create-ssl-secrets
create-ssl-secrets: $(BUILD_OUTPUT)/certs
	@echo "Deleting SSL secret $(TEST_SSL_SECRET)"
	kubectl --namespace $(TEST_NAMESPACE) delete secret $(TEST_SSL_SECRET) && echo "secret deleted" || true
	@echo "Creating SSL secret $(TEST_SSL_SECRET)"
	kubectl create secret generic $(TEST_SSL_SECRET) \
		--namespace $(TEST_NAMESPACE) \
		--from-file=keystore.jks=build/_output/certs/icarus.jks \
		--from-file=storepass.txt=build/_output/certs/storepassword.txt \
		--from-file=keypass.txt=build/_output/certs/keypassword.txt \
		--from-file=truststore.jks=build/_output/certs/truststore-guardians.jks \
		--from-file=trustpass.txt=build/_output/certs/trustpassword.txt \
		--from-file=operator.key=build/_output/certs/icarus.key \
		--from-file=operator.crt=build/_output/certs/icarus.crt \
		--from-file=operator-ca.crt=build/_output/certs/guardians-ca.crt

# ----------------------------------------------------------------------------------------------------------------------
# Build the Java artifacts
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-mvn
build-mvn:
	mvn $(USE_MAVEN_SETTINGS) -B -f java package -DskipTests

# ----------------------------------------------------------------------------------------------------------------------
# Build and test the Java artifacts
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: test-mvn
test-mvn: build-mvn
	mvn $(USE_MAVEN_SETTINGS) -B -f java verify

# ----------------------------------------------------------------------------------------------------------------------
# Build the examples
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-examples
build-examples:
	mvn $(USE_MAVEN_SETTINGS) -B -f examples package -DskipTests

# ----------------------------------------------------------------------------------------------------------------------
# Build and test the examples
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: test-examples
test-examples: build-examples
	mvn $(USE_MAVEN_SETTINGS) -B -f examples verify

# ----------------------------------------------------------------------------------------------------------------------
# Run all unit tests (both Go and Java)
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: test-all
test-all: test-mvn test-operator

# ----------------------------------------------------------------------------------------------------------------------
# Push the Operator Docker image
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-operator-image
push-operator-image: $(BUILD_TARGETS)/build-operator
ifeq ($(OPERATOR_RELEASE_IMAGE), $(OPERATOR_IMAGE))
	@echo "Pushing $(OPERATOR_IMAGE)"
	docker push $(OPERATOR_IMAGE)
else
	@echo "Tagging $(OPERATOR_IMAGE) as $(OPERATOR_RELEASE_IMAGE)"
	docker tag $(OPERATOR_IMAGE) $(OPERATOR_RELEASE_IMAGE)
	@echo "Pushing $(OPERATOR_RELEASE_IMAGE)"
	docker push $(OPERATOR_RELEASE_IMAGE)
endif

# ----------------------------------------------------------------------------------------------------------------------
# Build the Operator Utils Docker image
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-utils-image
build-utils-image: build-mvn $(BUILD_BIN)/runner $(BUILD_BIN)/op-test
	cp $(BUILD_BIN)/op-test java/coherence-utils/target/docker/op-test
	cp $(BUILD_BIN)/runner  java/coherence-utils/target/docker/runner
	docker build -t $(UTILS_IMAGE) java/coherence-utils/target/docker

# ----------------------------------------------------------------------------------------------------------------------
# Push the Operator Utils Docker image
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-utils-image
push-utils-image:
ifeq ($(UTILS_RELEASE_IMAGE), $(UTILS_IMAGE))
	@echo "Pushing $(UTILS_IMAGE)"
	docker push $(UTILS_IMAGE)
else
	@echo "Tagging $(UTILS_IMAGE) as $(UTILS_RELEASE_IMAGE)"
	docker tag $(UTILS_IMAGE) $(UTILS_RELEASE_IMAGE)
	@echo "Pushing $(UTILS_RELEASE_IMAGE)"
	docker push $(UTILS_RELEASE_IMAGE)
endif

# ----------------------------------------------------------------------------------------------------------------------
# Build the Operator JIB Test image
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-jib-image
build-jib-image: build-mvn
	mvn $(USE_MAVEN_SETTINGS) -B -f java package jib:dockerBuild -DskipTests -Djib.to.image=$(TEST_APPLICATION_IMAGE)

# ----------------------------------------------------------------------------------------------------------------------
# Push the Operator JIB Test Docker images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-jib-image
push-jib-image:
	@echo "Pushing $(TEST_APPLICATION_IMAGE)"
	docker push $(TEST_APPLICATION_IMAGE)

# ----------------------------------------------------------------------------------------------------------------------
# Build all of the Docker images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-all-images
build-all-images: $(BUILD_TARGETS)/build-operator build-utils-image build-jib-image

# ----------------------------------------------------------------------------------------------------------------------
# Push all of the Docker images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-all-images
push-all-images: push-operator-image push-utils-image push-jib-image

# ----------------------------------------------------------------------------------------------------------------------
# Push all of the Docker images that are released
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-release-images
push-release-images: push-operator-image push-utils-image

# ----------------------------------------------------------------------------------------------------------------------
# Run the Operator locally.
#
# To exit out of the local Operator you can use ctrl-c or ctrl-z but
# sometimes this leaves orphaned processes on the local machine so
# ensure these are killed run "make debug-stop"
# ----------------------------------------------------------------------------------------------------------------------
run: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
run: export UTILS_IMAGE := $(UTILS_IMAGE)
run:
	go run -ldflags='-X=main.BuildInfo=$(BUILD_INFO)' ./main.go \
	    2>&1 | tee $(TEST_LOGS_DIR)/operator-debug.out

# ----------------------------------------------------------------------------------------------------------------------
# Run the Operator locally after deleting and recreating the test namespace.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: run-clean
run-clean: reset-namespace run

# ----------------------------------------------------------------------------------------------------------------------
# Run the Operator in locally debug mode,
# Running this task will start the Operator and pause it until a Delve
# is attached.
#
# To exit out of the local Operator you can use ctrl-c or ctrl-z but
# sometimes this leaves orphaned processes on the local machine so
# ensure these are killed run "make debug-stop"
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: run-debug
run-debug: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
run-debug: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
run-debug: export UTILS_IMAGE := $(UTILS_IMAGE)
run-debug: export VERSION := $(VERSION)
run-debug: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
run-debug: export UTILS_IMAGE := $(UTILS_IMAGE)
run-debug:
	$(OPERATOR_SDK) run local --watch-namespace=$(TEST_NAMESPACE) \
	--go-ldflags="-X=main.BuildInfo=$(BUILD_INFO)" \
	--operator-flags="--coherence-image=$(COHERENCE_IMAGE) --utils-image=$(UTILS_IMAGE)" \
	--enable-delve \
	2>&1 | tee $(TEST_LOGS_DIR)/operator-debug.out

# ----------------------------------------------------------------------------------------------------------------------
# Run the Operator locally in debug mode after deleting and recreating
# the test namespace.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: run-debug-clean
run-debug-clean: reset-namespace run-debug

# ----------------------------------------------------------------------------------------------------------------------
# Kill any locally running Operator
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: stop
stop:
	./hack/kill-local.sh

# ----------------------------------------------------------------------------------------------------------------------
# Start a Kind cluster
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: kind
kind: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
kind:
	./hack/kind.sh
	docker pull $(COHERENCE_IMAGE)
	kind load docker-image --name operator $(COHERENCE_IMAGE)

# ----------------------------------------------------------------------------------------------------------------------
# Start a Kind 1.12 cluster
# ----------------------------------------------------------------------------------------------------------------------
kind-12: kind-12-start kind-load

kind-12-start:
	./hack/kind.sh --image "kindest/node:v1.12.10@sha256:faeb82453af2f9373447bb63f50bae02b8020968e0889c7fa308e19b348916cb"
	docker pull $(COHERENCE_IMAGE) || true
	kind load docker-image --name operator $(COHERENCE_IMAGE) || true

# ----------------------------------------------------------------------------------------------------------------------
# Start a Kind 1.16 cluster
# ----------------------------------------------------------------------------------------------------------------------
kind-16: kind-16-start kind-load

kind-16-start:
	./hack/kind.sh --image "kindest/node:v1.16.9@sha256:7175872357bc85847ec4b1aba46ed1d12fa054c83ac7a8a11f5c268957fd5765"
	docker pull $(COHERENCE_IMAGE) || true
	kind load docker-image --name operator $(COHERENCE_IMAGE) || true


# ----------------------------------------------------------------------------------------------------------------------
# Start a Kind 1.18 cluster
# ----------------------------------------------------------------------------------------------------------------------
kind-18: kind-18-start kind-load

kind-18-start:
	./hack/kind.sh --image "kindest/node:v1.18.2@sha256:7b27a6d0f2517ff88ba444025beae41491b016bc6af573ba467b70c5e8e0d85f"
	docker pull $(COHERENCE_IMAGE) || true
	kind load docker-image --name operator $(COHERENCE_IMAGE) || true

# ----------------------------------------------------------------------------------------------------------------------
# Load images into Kind
# ----------------------------------------------------------------------------------------------------------------------
kind-load:
	kind load docker-image --name operator $(OPERATOR_IMAGE)|| true
	kind load docker-image --name operator $(UTILS_IMAGE)|| true
	kind load docker-image --name operator $(TEST_APPLICATION_IMAGE)|| true

# ----------------------------------------------------------------------------------------------------------------------
# Install Prometheus
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-prometheus
install-prometheus:
	kubectl create ns $(TEST_NAMESPACE) || true
	kubectl create -f etc/prometheus-rbac.yaml
	helm repo add stable https://kubernetes-charts.storage.googleapis.com/ || true
	@echo "Create Grafana Dashboards ConfigMap:"
	kubectl -n $(TEST_NAMESPACE) create configmap coherence-grafana-dashboards --from-file=$(GRAFANA_DASHBOARDS)
	kubectl -n $(TEST_NAMESPACE) label configmap coherence-grafana-dashboards grafana_dashboard=1
	@echo "Getting Helm Version:"
	helm version
	@echo "Installing stable/prometheus-operator:"
	helm install --atomic --namespace $(TEST_NAMESPACE) --version $(PROMETHEUS_OPERATOR_VERSION) --wait \
		--set grafana.enabled=$(PROMETHEUS_INCLUDE_GRAFANA) \
		--values etc/prometheus-values.yaml prometheus stable/prometheus-operator
	@echo "Installing Prometheus instance:"
	kubectl -n $(TEST_NAMESPACE) apply -f etc/prometheus.yaml
	sleep 10
	kubectl -n $(TEST_NAMESPACE) wait --for=condition=Ready pod/prometheus-prometheus-0

# ----------------------------------------------------------------------------------------------------------------------
# Uninstall Prometheus
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: uninstall-prometheus
uninstall-prometheus:
	kubectl -n $(TEST_NAMESPACE) delete -f etc/prometheus.yaml || true
	kubectl -n $(TEST_NAMESPACE) delete configmap coherence-grafana-dashboards || true
	helm --namespace $(TEST_NAMESPACE) delete prometheus || true
	kubectl delete -f etc/prometheus-rbac.yaml || true

# ----------------------------------------------------------------------------------------------------------------------
# Start a port-forward process to the Grafana Pod.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: port-forward-grafana
port-forward-grafana: export GRAFANA_POD := $(shell kubectl -n $(TEST_NAMESPACE) get pod -l app.kubernetes.io/name=grafana -o name)
port-forward-grafana:
	@echo "Reach Grafana on http://127.0.0.1:3000"
	@echo "User: admin Password: prom-operator"
	kubectl -n $(TEST_NAMESPACE) port-forward $(GRAFANA_POD) 3000:3000

# ----------------------------------------------------------------------------------------------------------------------
# Install Elasticsearch & Kibana
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-elastic
install-elastic: helm-install-elastic kibana-import

.PHONY: helm-install-elastic
helm-install-elastic:
	kubectl create ns $(TEST_NAMESPACE) || true
#   Create the ConfigMap containing the Coherence Kibana dashboards
	kubectl -n $(TEST_NAMESPACE) delete secret coherence-kibana-dashboard || true
	kubectl -n $(TEST_NAMESPACE) create secret generic --from-file dashboards/kibana/kibana-dashboard-data.json coherence-kibana-dashboard
#   Create the ConfigMap containing the Coherence Kibana dashboards import script
	kubectl -n $(TEST_NAMESPACE) delete secret coherence-kibana-import || true
	kubectl -n $(TEST_NAMESPACE) create secret generic --from-file etc/coherence-dashboard-import.sh coherence-kibana-import
#   Set-up the Elastic Helm repo
	@echo "Getting Helm Version:"
	helm version
	helm repo add elastic https://helm.elastic.co || true
#   Install Elasticsearch
	helm install --atomic --namespace $(TEST_NAMESPACE) --version $(ELASTIC_VERSION) --wait --timeout=10m \
		--debug --values etc/elastic-values.yaml elasticsearch elastic/elasticsearch
#   Install Kibana
	helm install --atomic --namespace $(TEST_NAMESPACE) --version $(ELASTIC_VERSION) --wait --timeout=10m \
		--debug --values etc/kibana-values.yaml kibana elastic/kibana \

.PHONY: kibana-import
kibana-import:
	KIBANA_POD=$$(kubectl -n $(TEST_NAMESPACE) get pod -l app=kibana -o name) \
	; kubectl -n $(TEST_NAMESPACE) exec -it $${KIBANA_POD} /bin/bash /usr/share/kibana/data/coherence/scripts/coherence-dashboard-import.sh

# ----------------------------------------------------------------------------------------------------------------------
# Uninstall Elasticsearch & Kibana
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: uninstall-elastic
uninstall-elastic:
	helm uninstall --namespace $(TEST_NAMESPACE) kibana || true
	helm uninstall --namespace $(TEST_NAMESPACE) elasticsearch || true
	kubectl -n $(TEST_NAMESPACE) delete pvc elasticsearch-master-elasticsearch-master-0 || true

# ----------------------------------------------------------------------------------------------------------------------
# Start a port-forward process to the Kibana Pod.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: port-forward-kibana
port-forward-kibana: export KIBANA_POD := $(shell kubectl -n $(TEST_NAMESPACE) get pod -l app=kibana -o name)
port-forward-kibana:
	@echo "Reach Kibana on http://127.0.0.1:5601"
	kubectl -n $(TEST_NAMESPACE) port-forward $(KIBANA_POD) 5601:5601

# ----------------------------------------------------------------------------------------------------------------------
# Start a port-forward process to the Elasticsearch Pod.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: port-forward-es
port-forward-es: export ES_POD := $(shell kubectl -n $(TEST_NAMESPACE) get pod -l app=elasticsearch-master -o name)
port-forward-es:
	@echo "Reach Elasticsearch on http://127.0.0.1:9200"
	kubectl -n $(TEST_NAMESPACE) port-forward $(ES_POD) 9200:9200


# ----------------------------------------------------------------------------------------------------------------------
# Delete all of the Coherence resources from the test namespace.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: delete-coherence-clusters
delete-coherence-clusters:
	for i in $$(kubectl -n  $(TEST_NAMESPACE) get coherence -o name); do \
		kubectl -n $(TEST_NAMESPACE) delete $${i}; \
	done

# ----------------------------------------------------------------------------------------------------------------------
# Obtain the golangci-lint binary
# ----------------------------------------------------------------------------------------------------------------------
$(BUILD_BIN)/golangci-lint:
	@mkdir -p $(BUILD_BIN)
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BUILD_BIN) v1.29.0

# ----------------------------------------------------------------------------------------------------------------------
# Executes golangci-lint to perform various code review checks on the source.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: golangci
golangci: $(BUILD_BIN)/golangci-lint
	$(BUILD_BIN)/golangci-lint run -v --timeout=5m --skip-files=zz_.*,generated/* ./api/... ./controllers/... ./pkg/... ./cmd/...


# ----------------------------------------------------------------------------------------------------------------------
# Performs a copyright check.
# To add exclusions add the file or folder pattern using the -X parameter.
# Add directories to be scanned at the end of the parameter list.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: copyright
copyright:
	@java -cp etc/glassfish-copyright-maven-plugin-2.1.jar \
	  org.glassfish.copyright.Copyright -C etc/copyright.txt \
	  -X .adoc \
	  -X bin/ \
	  -X build/_output/ \
	  -X clientset/ \
	  -X dashboards/grafana/ \
	  -X dashboards/kibana/ \
	  -X /Dockerfile \
	  -X docs/ \
	  -X etc/copyright.txt \
	  -X etc/intellij-codestyle.xml \
	  -X etc/sdk/ \
	  -X go.mod \
	  -X go.sum \
	  -X HEADER.txt \
	  -X .iml \
	  -X java/src/copyright/EXCLUDE.txt \
	  -X Jenkinsfile \
	  -X .jar \
	  -X .jks \
	  -X .json \
	  -X LICENSE.txt \
	  -X Makefile \
	  -X .md \
	  -X PROJECT \
	  -X .sh \
	  -X temp/ \
	  -X temp/olm/ \
	  -X /test-report.xml \
	  -X THIRD_PARTY_LICENSES.txt \
	  -X tools.go \
	  -X .tpl \
	  -X .yaml \
	  -X pkg/apis/coherence/legacy/zz_generated.deepcopy.go \
	  -X pkg/data/zz_generated_assets.go \
	  -X zz_generated.

# ----------------------------------------------------------------------------------------------------------------------
# Executes the code review targets.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: code-review
code-review: export MAVEN_USER := $(MAVEN_USER)
code-review: export MAVEN_PASSWORD := $(MAVEN_PASSWORD)
code-review: golangci copyright
	mvn $(USE_MAVEN_SETTINGS) -B -f java validate -DskipTests -P checkstyle
	mvn $(USE_MAVEN_SETTINGS) -B -f examples validate -DskipTests -P checkstyle

# ----------------------------------------------------------------------------------------------------------------------
# Display the full version string for the artifacts that would be built.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: version
version:
	@echo ${VERSION}

# ----------------------------------------------------------------------------------------------------------------------
# Build the documentation.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: docs
docs:
	mvn $(USE_MAVEN_SETTINGS) -B -f java install -P docs -pl docs -DskipTests -Doperator.version=$(VERSION)

# ----------------------------------------------------------------------------------------------------------------------
# Start a local web server to serve the documentation.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: serve-docs
serve-docs:
	@echo "Serving documentation on http://localhost:8080"
	cd $(BUILD_OUTPUT)/docs; \
	python -m SimpleHTTPServer 8080

# ----------------------------------------------------------------------------------------------------------------------
# Release the Coherence Operator dashboards
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: release-dashboards
release-dashboards:
	@echo "Releasing Dashboards $(VERSION)"
	mkdir -p $(BUILD_OUTPUT)/dashboards/$(VERSION) || true
	tar -czvf $(BUILD_OUTPUT)/dashboards/$(VERSION)/coherence-dashboards.tar.gz  dashboards/
	kubectl create configmap coherence-grafana-dashboards --from-file=dashboards/grafana \
		--dry-run -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION)/coherence-grafana-dashboards.yaml
	kubectl create configmap coherence-grafana-dashboards --from-file=dashboards/grafana-legacy \
		--dry-run -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION)/coherence-grafana-legacy-dashboards.yaml
	kubectl create configmap coherence-kibana-dashboards --from-file=dashboards/kibana \
		--dry-run -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION)/coherence-kibana-dashboards.yaml
	mkdir -p dashboards || true
	mv $(BUILD_OUTPUT)/dashboards/$(VERSION)/ dashboards/

# ----------------------------------------------------------------------------------------------------------------------
# Release the Coherence Operator to the gh-pages branch.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: release-ghpages
release-ghpages:  docs
	@echo "Releasing Dashboards $(VERSION)"
	mkdir -p $(BUILD_OUTPUT)/dashboards/$(VERSION) || true
	tar -czvf $(BUILD_OUTPUT)/dashboards/$(VERSION)/coherence-dashboards.tar.gz  dashboards/
	kubectl create configmap coherence-grafana-dashboards --from-file=dashboards/grafana \
		--dry-run -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION)/coherence-grafana-dashboards.yaml
	kubectl create configmap coherence-grafana-dashboards --from-file=dashboards/grafana-legacy \
		--dry-run -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION)/coherence-grafana-legacy-dashboards.yaml
	kubectl create configmap coherence-kibana-dashboards --from-file=dashboards/kibana \
		--dry-run -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION)/coherence-kibana-dashboards.yaml
	cp hack/docs-unstable-index.sh $(BUILD_OUTPUT)/docs-unstable-index.sh
	git stash save --keep-index --include-untracked || true
	git stash drop || true
	git checkout gh-pages
	git pull
	mkdir -p dashboards || true
	rm -rf dashboards/$(VERSION) || true
	mv $(BUILD_OUTPUT)/dashboards/$(VERSION)/ dashboards/
	git add dashboards/$(VERSION)/*
ifeq (true, $(PRE_RELEASE))
	mkdir -p docs-unstable || true
	rm -rf docs-unstable/$(VERSION)/ || true
	mv $(BUILD_OUTPUT)/docs/ docs-unstable/$(VERSION)/
	sh $(BUILD_OUTPUT)/docs-unstable-index.sh
	ls -ls docs-unstable

	git status
	git add docs-unstable/*
else
	mkdir docs/$(VERSION) || true
	rm -rf docs/$(VERSION)/ || true
	mv $(BUILD_OUTPUT)/docs/ docs/$(VERSION)/
	ls -ls docs

	git status
	git add docs/*
endif
	git clean -d -f
	git status
	git commit -m "adding Coherence Operator docs version: $(VERSION)"
	git log -1
ifeq (true, $(RELEASE_DRY_RUN))
	@echo "release dry-run - would have pushed docs $(VERSION) to gh-pages"
else
	git push origin gh-pages
endif


# ----------------------------------------------------------------------------------------------------------------------
# Tag Git for the release.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: release-tag
release-tag:
ifeq (true, $(RELEASE_DRY_RUN))
	@echo "release dry-run - would have created release tag v$(VERSION)"
else
	@echo "creating release tag v$(VERSION)"
	git push origin :refs/tags/v$(VERSION)
	git tag -f -a -m "built $(VERSION)" v$(VERSION)
	git push origin --tags
endif

# ----------------------------------------------------------------------------------------------------------------------
# Release the Coherence Operator.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: release
release:
ifeq (true, $(RELEASE_DRY_RUN))
release: build-all-images release-tag release-ghpages
	@echo "release dry-run: would have pushed images"
else
release: build-all-images release-tag release-ghpages push-release-images
endif


# ----------------------------------------------------------------------------------------------------------------------
# List all of the targets in the Makefile
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: list
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
