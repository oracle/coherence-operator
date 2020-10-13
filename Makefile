# ----------------------------------------------------------------------------------------------------------------------
# Copyright (c) 2019, 2020, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
# ----------------------------------------------------------------------------------------------------------------------
# This is the Makefile to build the Coherence Kubernetes Operator.
# ----------------------------------------------------------------------------------------------------------------------

# The version of the Operator being build - this should be a valid SemVer format
VERSION ?= 3.1.2

# The version number to be replaced by this release
PREV_VERSION ?= 3.1.2

# The operator version to use to run certification tests against
CERTIFICATION_VERSION ?= $(VERSION)

# The previous Operator version used to run the compatibility tests.
COMPATIBLE_VERSION  = 3.0.2
# The selector to use to find Operator Pods of the COMPATIBLE_VERSION (do not put in double quotes!!)
COMPATIBLE_SELECTOR = component=coherence-operator

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
OPERATOR_SDK_VERSION := v1.0.0
OPERATOR_SDK          = $(CURRDIR)/hack/sdk/$(UNAME_S)-$(UNAME_M)/operator-sdk

# The Coherence image to use for deployments that do not specify an image
COHERENCE_IMAGE   ?= oraclecoherence/coherence-ce:20.06.1
# This is the Coherence image that will be used in tests.
# Changing this variable will allow test builds to be run against different Coherence versions
# without altering the default image name.
TEST_COHERENCE_IMAGE ?= $(COHERENCE_IMAGE)

# Operator image names
RELEASE_IMAGE_PREFIX   ?= ghcr.io/oracle/
OPERATOR_IMAGE_REPO    := $(RELEASE_IMAGE_PREFIX)coherence-operator
OPERATOR_IMAGE         := $(OPERATOR_IMAGE_REPO):$(VERSION)
UTILS_IMAGE            ?= $(OPERATOR_IMAGE_REPO):$(VERSION)-utils
# The Operator images to push
OPERATOR_RELEASE_REPO  ?= $(OPERATOR_IMAGE_REPO)
OPERATOR_RELEASE_IMAGE := $(OPERATOR_RELEASE_REPO):$(VERSION)
UTILS_RELEASE_IMAGE    := $(OPERATOR_RELEASE_REPO):$(VERSION)-utils
BUNDLE_RELEASE_IMAGE   := $(OPERATOR_RELEASE_REPO):$(VERSION)-bundle

# The test application images used in integration tests
TEST_APPLICATION_IMAGE             := $(RELEASE_IMAGE_PREFIX)operator-test:$(VERSION)
TEST_APPLICATION_IMAGE_SPRING      := $(RELEASE_IMAGE_PREFIX)operator-test:$(VERSION)-spring
TEST_APPLICATION_IMAGE_SPRING_FAT  := $(RELEASE_IMAGE_PREFIX)operator-test:$(VERSION)-spring-fat
TEST_APPLICATION_IMAGE_SPRING_CNBP := $(RELEASE_IMAGE_PREFIX)operator-test:$(VERSION)-spring-cnbp

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
#   make e2e-test GO_TEST_FLAGS_E2E=-run=^TestZone$
ifeq ($(origin RUN_ONE), undefined)
GO_TEST_FLAGS     ?= -timeout=20m
GO_TEST_FLAGS_E2E ?= -timeout=100m
else
GO_TEST_FLAGS     ?= -timeout=20m -run=^$(RUN_ONE)$$
GO_TEST_FLAGS_E2E ?= -timeout=100m -run=^$(RUN_ONE)$$
endif

.PHONY: jk
jk:
	echo "GO_TEST_FLAGS=$(GO_TEST_FLAGS)"
	echo "GO_TEST_FLAGS_E2E=$(GO_TEST_FLAGS_E2E)"

# default as in test/e2e/helper/proj_helpers.go
OPERATOR_NAMESPACE ?= operator-test
# the optional namespaces the operator should watch
WATCH_NAMESPACE ?=
# flag indicating whether the test namespace should be reset (deleted and recreated) before tests
CREATE_OPERATOR_NAMESPACE ?= true

# Prometheus Operator settings (used in integration tests)
PROMETHEUS_INCLUDE_GRAFANA   ?= true
PROMETHEUS_OPERATOR_VERSION  ?= 8.13.7
PROMETHEUS_ADAPTER_VERSION   ?= 2.5.0
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

override BUILD_OUTPUT        := ./build/_output
override BUILD_ASSETS        := $(BUILD_OUTPUT)/assets
override BUILD_BIN           := ./bin
override BUILD_DEPLOY        := $(BUILD_OUTPUT)/config
override BUILD_HELM          := $(BUILD_OUTPUT)/helm-charts
override BUILD_MANIFESTS     := $(BUILD_OUTPUT)/manifests
override BUILD_MANIFESTS_PKG := $(BUILD_OUTPUT)/coherence-operator-manifests-$(VERSION).tar.gz
override BUILD_PROPS         := $(BUILD_OUTPUT)/build.properties
override BUILD_TARGETS       := $(BUILD_OUTPUT)/targets
override TEST_LOGS_DIR       := $(BUILD_OUTPUT)/test-logs

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

SOURCE_DATE_EPOCH := $(shell git show -s --format=format:%ct HEAD)
DATE_FMT          := "%Y-%m-%dT%H:%M:%SZ"
BUILD_DATE        := $(shell date -u -d "@$SOURCE_DATE_EPOCH" "+${DATE_FMT}" 2>/dev/null || date -u -r "${SOURCE_DATE_EPOCH}" "+${DATE_FMT}" 2>/dev/null || date -u "+${DATE_FMT}")

BUILD_INFO       = $(VERSION)|$(GITCOMMIT)|$(BUILD_DATE)
LDFLAGS          = -X main.Version=$(VERSION) -X main.Commit=$(GITCOMMIT) -X main.Date=$(BUILD_DATE)
GOS              = $(shell find . -type f -name "*.go" ! -name "*_test.go")
HELM_FILES       = $(shell find helm-charts/coherence-operator -type f)
API_GO_FILES     = $(shell find . -type f -name "*.go" ! -name "*_test.go"  ! -name "zz*.go")
CRDV1_FILES      = $(shell find ./config/crd -type f)
CRDV1BETA1_FILES = $(shell find ./config/crd-v1beta1 -type f)
CRD_V1           ?= $(shell kubectl api-versions | grep '^apiextensions.k8s.io/v1$$')

TEST_SSL_SECRET := coherence-ssl-secret

.PHONY: all
all: build-all-images helm-chart

# ----------------------------------------------------------------------------------------------------------------------
# Configure the build properties
# ----------------------------------------------------------------------------------------------------------------------
$(BUILD_PROPS):
	# Ensures that build output directories exist
	@echo "Creating build directories"
	@mkdir -p $(BUILD_OUTPUT)
	@mkdir -p $(BUILD_ASSETS)
	@mkdir -p $(BUILD_BIN)
	@mkdir -p $(BUILD_DEPLOY)
	@mkdir -p $(BUILD_HELM)
	@mkdir -p $(BUILD_MANIFESTS)
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
$(BUILD_BIN)/manager: $(BUILD_PROPS) $(GOS) generate manifests
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags "$(LDFLAGS)" -a -o $(BUILD_BIN)/manager main.go

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
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags "$(LDFLAGS)" -o $(BUILD_BIN)/runner ./runner

# ----------------------------------------------------------------------------------------------------------------------
# Internal make step that builds the Operator legacy converter
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: converter
converter: $(BUILD_BIN)/converter $(BUILD_BIN)/converter-linux-amd64 $(BUILD_BIN)/converter-darwin-amd64 $(BUILD_BIN)/converter-windows-amd64

$(BUILD_BIN)/converter: $(BUILD_PROPS) $(GOS)
	CGO_ENABLED=0 GO111MODULE=on GOOS=$(OS) GOARCH=$(ARCH) go build -o $(BUILD_BIN)/converter ./converter

$(BUILD_BIN)/converter-linux-amd64: $(BUILD_PROPS) $(GOS)
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o $(BUILD_BIN)/converter-linux-amd64 ./converter

$(BUILD_BIN)/converter-darwin-amd64: $(BUILD_PROPS) $(GOS)
	CGO_ENABLED=0 GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o $(BUILD_BIN)/converter-darwin-amd64 ./converter

$(BUILD_BIN)/converter-windows-amd64: $(BUILD_PROPS) $(GOS)
	CGO_ENABLED=0 GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o $(BUILD_BIN)/converter-windows-amd64 ./converter

# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go unit tests that do not require a k8s cluster
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: test-operator
test-operator: export CGO_ENABLED = 0
test-operator: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
test-operator: export UTILS_IMAGE := $(UTILS_IMAGE)
test-operator: $(BUILD_PROPS) $(BUILD_TARGETS)/manifests $(BUILD_TARGETS)/generate gotestsum
	@echo "Running operator tests"
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-test.xml \
	  -- $(GO_TEST_FLAGS) -v ./api/... ./controllers/... ./pkg/...

# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end tests that require a k8s cluster using
# a LOCAL operator instance (i.e. the operator is not deployed to k8s).
# These tests will use whichever k8s cluster the local environment
# is pointing to.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: e2e-local-test
e2e-local-test: export CGO_ENABLED = 0
e2e-local-test: export OPERATOR_NAMESPACE := $(OPERATOR_NAMESPACE)
e2e-local-test: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
e2e-local-test: export TEST_APPLICATION_IMAGE_SPRING := $(TEST_APPLICATION_IMAGE_SPRING)
e2e-local-test: export TEST_APPLICATION_IMAGE_SPRING_FAT := $(TEST_APPLICATION_IMAGE_SPRING_FAT)
e2e-local-test: export TEST_APPLICATION_IMAGE_SPRING_CNBP := $(TEST_APPLICATION_IMAGE_SPRING_CNBP)
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
e2e-test: export MF = $(MAKEFLAGS)
e2e-test: prepare-e2e-test
	$(MAKE) run-e2e-test $${MF} \
	; rc=$$? \
	; $(MAKE) undeploy $${MF} \
	; $(MAKE) delete-namespace $${MF} \
	; exit $$rc

.PHONY: prepare-e2e-test
prepare-e2e-test: $(BUILD_TARGETS)/build-operator reset-namespace create-ssl-secrets install-crds deploy

.PHONY: run-e2e-test
run-e2e-test: export CGO_ENABLED = 0
run-e2e-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
run-e2e-test: export OPERATOR_NAMESPACE := $(OPERATOR_NAMESPACE)
run-e2e-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
run-e2e-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
run-e2e-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
run-e2e-test: export TEST_ASSET_KUBECTL := $(TEST_ASSET_KUBECTL)
run-e2e-test: export VERSION := $(VERSION)
run-e2e-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
run-e2e-test: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
run-e2e-test: export UTILS_IMAGE := $(UTILS_IMAGE)
run-e2e-test: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
run-e2e-test: export TEST_APPLICATION_IMAGE_SPRING := $(TEST_APPLICATION_IMAGE_SPRING)
run-e2e-test: export TEST_APPLICATION_IMAGE_SPRING_FAT := $(TEST_APPLICATION_IMAGE_SPRING_FAT)
run-e2e-test: export TEST_APPLICATION_IMAGE_SPRING_CNBP := $(TEST_APPLICATION_IMAGE_SPRING_CNBP)
run-e2e-test: gotestsum
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/e2e/remote/...


# ----------------------------------------------------------------------------------------------------------------------
# Run the end-to-end Helm chart tests.
# ----------------------------------------------------------------------------------------------------------------------
e2e-helm-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
e2e-helm-test: export UTILS_IMAGE := $(UTILS_IMAGE)
e2e-helm-test: $(BUILD_PROPS) $(BUILD_HELM)/coherence-operator-$(VERSION).tgz reset-namespace gotestsum
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-helm-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/e2e/helm/...


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
.PHONY: e2e-prometheus-test
e2e-prometheus-test: export MF = $(MAKEFLAGS)
e2e-prometheus-test: reset-namespace install-prometheus $(BUILD_TARGETS)/build-operator create-ssl-secrets install-crds deploy
	$(MAKE) run-prometheus-test $${MF} \
	; rc=$$? \
	; $(MAKE) uninstall-prometheus $${MF} \
	; $(MAKE) undeploy $${MF} \
	; $(MAKE) delete-namespace $${MF} \
	; exit $$rc

.PHONY: run-prometheus-test
run-prometheus-test: export CGO_ENABLED = 0
run-prometheus-test: export OPERATOR_NAMESPACE := $(OPERATOR_NAMESPACE)
run-prometheus-test: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
run-prometheus-test: export TEST_APPLICATION_IMAGE_SPRING := $(TEST_APPLICATION_IMAGE_SPRING)
run-prometheus-test: export TEST_APPLICATION_IMAGE_SPRING_FAT := $(TEST_APPLICATION_IMAGE_SPRING_FAT)
run-prometheus-test: export TEST_APPLICATION_IMAGE_SPRING_CNBP := $(TEST_APPLICATION_IMAGE_SPRING_CNBP)
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
run-prometheus-test: gotestsum
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-prometheus-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/e2e/prometheus/...


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
.PHONY: e2e-elastic-test
e2e-elastic-test: export MF = $(MAKEFLAGS)
e2e-elastic-test: reset-namespace install-elastic $(BUILD_TARGETS)/build-operator create-ssl-secrets install-crds deploy
	$(MAKE) run-elastic-test $${MF} \
	; rc=$$? \
	; $(MAKE) uninstall-elastic $${MF} \
	; $(MAKE) undeploy $${MF} \
	; $(MAKE) delete-namespace $${MF} \
	; exit $$rc

.PHONY: run-elastic-test
run-elastic-test: export CGO_ENABLED = 0
run-elastic-test: export OPERATOR_NAMESPACE := $(OPERATOR_NAMESPACE)
run-elastic-test: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
run-elastic-test: export TEST_APPLICATION_IMAGE_SPRING := $(TEST_APPLICATION_IMAGE_SPRING)
run-elastic-test: export TEST_APPLICATION_IMAGE_SPRING_FAT := $(TEST_APPLICATION_IMAGE_SPRING_FAT)
run-elastic-test: export TEST_APPLICATION_IMAGE_SPRING_CNBP := $(TEST_APPLICATION_IMAGE_SPRING_CNBP)
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
run-elastic-test: gotestsum
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-elastic-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/e2e/elastic/...

# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end Operator Compatibility tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: compatibility-test
compatibility-test: export CGO_ENABLED = 0
compatibility-test: export OPERATOR_NAMESPACE := $(OPERATOR_NAMESPACE)
compatibility-test: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
compatibility-test: export TEST_APPLICATION_IMAGE_SPRING := $(TEST_APPLICATION_IMAGE_SPRING)
compatibility-test: export TEST_APPLICATION_IMAGE_SPRING_FAT := $(TEST_APPLICATION_IMAGE_SPRING_FAT)
compatibility-test: export TEST_APPLICATION_IMAGE_SPRING_CNBP := $(TEST_APPLICATION_IMAGE_SPRING_CNBP)
compatibility-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
compatibility-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
compatibility-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
compatibility-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
compatibility-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
compatibility-test: export VERSION := $(VERSION)
compatibility-test: export COMPATIBLE_VERSION := $(COMPATIBLE_VERSION)
compatibility-test: export COMPATIBLE_SELECTOR := $(COMPATIBLE_SELECTOR)
compatibility-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
compatibility-test: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
compatibility-test: export UTILS_IMAGE := $(UTILS_IMAGE)
compatibility-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
compatibility-test: undeploy build-all-images $(BUILD_HELM)/coherence-operator-$(VERSION).tgz undeploy clean-namespace reset-namespace gotestsum
	helm repo add coherence https://oracle.github.io/coherence-operator/charts
	helm repo update
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-compatibility-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/e2e/compatibility/...


# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end Operator certification tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# Note that the namespace will be created if it does not exist.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: certification-test
certification-test: export MF = $(MAKEFLAGS)
certification-test: install-certification
	$(MAKE) run-certification  $${MF} \
	; rc=$$? \
	; $(MAKE) cleanup-certification $${MF} \
	; exit $$rc


# ----------------------------------------------------------------------------------------------------------------------
# Install the Operator prior to running compatability tests.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-certification
install-certification: $(BUILD_TARGETS)/build-operator reset-namespace create-ssl-secrets deploy

# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end Operator certification tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# Note that the namespace will be created if it does not exist.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: run-certification
run-certification: export CGO_ENABLED = 0
run-certification: export OPERATOR_NAMESPACE := $(OPERATOR_NAMESPACE)
run-certification: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
run-certification: export TEST_APPLICATION_IMAGE_SPRING := $(TEST_APPLICATION_IMAGE_SPRING)
run-certification: export TEST_APPLICATION_IMAGE_SPRING_FAT := $(TEST_APPLICATION_IMAGE_SPRING_FAT)
run-certification: export TEST_APPLICATION_IMAGE_SPRING_CNBP := $(TEST_APPLICATION_IMAGE_SPRING_CNBP)
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
run-certification: gotestsum
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-certification-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/certification/...

# ----------------------------------------------------------------------------------------------------------------------
# Clean up after to running compatability tests.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: cleanup-certification
cleanup-certification: undeploy uninstall-crds clean-namespace

# ---------------------------------------------------------------------------
# Build the Coherence operator Helm chart and package it into a tar.gz
# ---------------------------------------------------------------------------
.PHONY: helm-chart
helm-chart: $(BUILD_PROPS) $(BUILD_HELM)/coherence-operator-$(VERSION).tgz

$(BUILD_HELM)/coherence-operator-$(VERSION).tgz: $(BUILD_PROPS) $(HELM_FILES) generate $(BUILD_TARGETS)/manifests $(GOBIN)/kustomize
# Copy the Helm chart from the source location to the distribution folder
	-mkdir -p $(BUILD_HELM)
	cp -R ./helm-charts/coherence-operator $(BUILD_HELM)
	$(call replaceprop,coherence-operator/Chart.yaml coherence-operator/values.yaml coherence-operator/templates/deployment.yaml)
# Package the chart into a .tr.gz - we don't use helm package as the version might not be SEMVER
	helm lint $(BUILD_HELM)/coherence-operator
	tar -C $(BUILD_HELM)/coherence-operator -czf $(BUILD_HELM)/coherence-operator-$(VERSION).tgz .

# ---------------------------------------------------------------------------
# Do a search and replace of properties in selected files in the Helm charts.
# This is done because the Helm charts can be large and processing every file
# makes the build slower.
# ---------------------------------------------------------------------------
define replaceprop
	for i in $(1); do \
		filename="$(BUILD_HELM)/$${i}"; \
		echo "Replacing properties in file $${filename}"; \
		if [ -f $${filename} ]; then \
			temp_file=$(BUILD_OUTPUT)/temp.out; \
			awk -F'=' 'NR==FNR {a[$$1]=$$2;next} {for (i in a) {x = sprintf("\\$${%s}", i); gsub(x, a[i])}}1' $(BUILD_PROPS) $${filename} > $${temp_file}; \
			mv $${temp_file} $${filename}; \
		fi \
	done
endef

# ----------------------------------------------------------------------------------------------------------------------
# Install CRDs into Kubernetes.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-crds
install-crds: prepare-deploy uninstall-crds
ifeq ("$(CRD_V1)","apiextensions.k8s.io/v1")
	$(GOBIN)/kustomize build $(BUILD_DEPLOY)/crd | kubectl create -f -
else
	$(GOBIN)/kustomize build $(BUILD_DEPLOY)/crd-v1beta1 | kubectl create -f --validate=false -
endif

# ----------------------------------------------------------------------------------------------------------------------
# Uninstall CRDs from Kubernetes.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: uninstall-crds
uninstall-crds: prepare-deploy
ifeq ("$(CRD_V1)","apiextensions.k8s.io/v1")
	$(GOBIN)/kustomize build $(BUILD_DEPLOY)/crd | kubectl delete -f - || true
else
	$(GOBIN)/kustomize build $(BUILD_DEPLOY)/crd-v1beta1 | kubectl delete -f - || true
endif

# ----------------------------------------------------------------------------------------------------------------------
# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: deploy
deploy: prepare-deploy
ifneq (,$(WATCH_NAMESPACE))
	cd $(BUILD_DEPLOY)/manager && $(GOBIN)/kustomize edit add configmap env-vars --from-literal WATCH_NAMESPACE=$(WATCH_NAMESPACE)
endif
	kubectl -n $(OPERATOR_NAMESPACE) create secret generic coherence-webhook-server-cert || true
	$(GOBIN)/kustomize build $(BUILD_DEPLOY)/default | kubectl apply -f -

.PHONY: prepare-deploy
prepare-deploy: manifests $(BUILD_TARGETS)/build-operator $(GOBIN)/kustomize
	$(call prepare_deploy,$(OPERATOR_IMAGE),$(OPERATOR_NAMESPACE))

# ----------------------------------------------------------------------------------------------------------------------
# Prepare the deployment manifests - this is called by a number of other targets.
# Parameter #1 is the Operator Image Name
# Parameter #2 is the name of the namespace to deploy into
# ----------------------------------------------------------------------------------------------------------------------
define prepare_deploy
	-rm -r $(BUILD_DEPLOY)
	mkdir -p $(BUILD_DEPLOY)
	cp -R config $(BUILD_OUTPUT)
	cd $(BUILD_DEPLOY)/manager && $(GOBIN)/kustomize edit add configmap env-vars --from-literal COHERENCE_IMAGE=$(COHERENCE_IMAGE)
	cd $(BUILD_DEPLOY)/manager && $(GOBIN)/kustomize edit add configmap env-vars --from-literal UTILS_IMAGE=$(UTILS_IMAGE)
	cd $(BUILD_DEPLOY)/manager && $(GOBIN)/kustomize edit set image controller=$(1)
	cd $(BUILD_DEPLOY)/default && $(GOBIN)/kustomize edit set namespace $(2)
endef

# ----------------------------------------------------------------------------------------------------------------------
# Un-deploy controller from the configured Kubernetes cluster in ~/.kube/config
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: undeploy
undeploy: $(BUILD_PROPS) $(BUILD_TARGETS)/manifests $(GOBIN)/kustomize
	$(call prepare_deploy,$(OPERATOR_IMAGE),$(OPERATOR_NAMESPACE))
	$(GOBIN)/kustomize build $(BUILD_DEPLOY)/default | kubectl delete -f - || true
	kubectl -n $(OPERATOR_NAMESPACE) delete secret coherence-webhook-server-cert || true
	kubectl delete mutatingwebhookconfiguration coherence-operator-mutating-webhook-configuration || true
	kubectl delete validatingwebhookconfiguration coherence-operator-validating-webhook-configuration || true


# ----------------------------------------------------------------------------------------------------------------------
# Tail the deployed operator logs.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: tail-logs
tail-logs: export POD=$(shell kubectl -n $(OPERATOR_NAMESPACE) get pod -l component=coherence-operator -o name)
tail-logs:
	kubectl -n $(OPERATOR_NAMESPACE) logs $(POD) -c manager -f

# ----------------------------------------------------------------------------------------------------------------------
# Generate manifests e.g. CRD, RBAC etc.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: manifests
manifests: $(BUILD_TARGETS)/manifests $(BUILD_MANIFESTS_PKG)

$(BUILD_MANIFESTS_PKG):
	rm -rf $(BUILD_MANIFESTS) || true
	mkdir -p $(BUILD_MANIFESTS)
	cp -R config/default/ $(BUILD_MANIFESTS)/default
	cp -R config/manager/ $(BUILD_MANIFESTS)/manager
	cp -R config/rbac/ $(BUILD_MANIFESTS)/rbac
	tar -C $(BUILD_OUTPUT) -czf $(BUILD_MANIFESTS_PKG) manifests/

$(BUILD_TARGETS)/manifests: $(BUILD_PROPS) config/crd/bases/coherence.oracle.com_coherence.yaml config/crd-v1beta1/bases/coherence.oracle.com_coherence.yaml docs/about/04_coherence_spec.adoc
	touch $(BUILD_TARGETS)/manifests

config/crd/bases/coherence.oracle.com_coherence.yaml: $(API_GO_FILES) $(GOBIN)/controller-gen
	$(GOBIN)/controller-gen "crd:trivialVersions=true,crdVersions={v1}" \
	  rbac:roleName=manager-role paths="{./api/...,./controllers/...}" \
	  output:crd:artifacts:config=config/crd/bases

config/crd-v1beta1/bases/coherence.oracle.com_coherence.yaml: $(API_GO_FILES) $(GOBIN)/controller-gen
	@echo "Generating CRD v1beta1"
	$(GOBIN)/controller-gen "crd:trivialVersions=true,crdVersions={v1beta1}" \
	  rbac:roleName=manager-role paths="{./api/...,./controllers/...}" \
	  output:crd:artifacts:config=config/crd-v1beta1/bases

# ----------------------------------------------------------------------------------------------------------------------
# Generate the config.json file used by the Operator for default configuration values
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: generate-config
generate-config: $(BUILD_PROPS) $(BUILD_OUTPUT)/config.json

$(BUILD_OUTPUT)/config.json:
	@echo "Generating Operator config"
	@printf "{\n \
	  \"coherence-image\": \"$(COHERENCE_IMAGE)\",\n \
	  \"utils-image\": \"$(UTILS_RELEASE_IMAGE)\"\n}\n" > $(BUILD_OUTPUT)/config.json

# ----------------------------------------------------------------------------------------------------------------------
# Generate code, configuration and docs.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: generate
generate: $(BUILD_TARGETS)/generate

$(BUILD_TARGETS)/generate: $(BUILD_PROPS) $(BUILD_OUTPUT)/config.json api/v1/zz_generated.deepcopy.go pkg/data/zz_generated_assets.go
	touch $(BUILD_TARGETS)/generate

api/v1/zz_generated.deepcopy.go: $(API_GO_FILES) $(GOBIN)/controller-gen
	$(GOBIN)/controller-gen object:headerFile="./hack/boilerplate.go.txt" paths="./api/..."

pkg/data/zz_generated_assets.go: $(BUILD_OUTPUT)/config.json $(CRDV1_FILES) $(CRDV1BETA1_FILES) $(GOBIN)/kustomize
	echo "Embedding configuration and CRD files"
	cp $(BUILD_OUTPUT)/config.json $(BUILD_ASSETS)/config.json
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
# find or download yq
# ----------------------------------------------------------------------------------------------------------------------
$(GOBIN)/yq:
	@{ \
	set -e ;\
	YQ_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$YQ_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get github.com/mikefarah/yq/v3 ;\
	rm -rf $$YQ_GEN_TMP_DIR ;\
	}
	YQ=$(GOBIN)/yq

# ----------------------------------------------------------------------------------------------------------------------
# Generate bundle manifests and metadata, then validate generated files.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: bundle
bundle: $(BUILD_PROPS) $(GOBIN)/kustomize $(BUILD_TARGETS)/manifests
	$(OPERATOR_SDK) generate kustomize manifests -q
	cd config/manager && $(GOBIN)/kustomize edit set image controller=$(OPERATOR_IMAGE)
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
	go run ./docgen/ \
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
	rm pkg/data/zz_generated_assets.go || true
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
ifeq ($(CREATE_OPERATOR_NAMESPACE),true)
	@echo "Creating test namespace $(OPERATOR_NAMESPACE)"
	kubectl create namespace $(OPERATOR_NAMESPACE)
endif
ifneq ($(DOCKER_SERVER),)
	@echo "Creating pull secrets for $(DOCKER_SERVER)"
	kubectl create secret docker-registry coherence-k8s-operator-development-secret \
								--namespace $(OPERATOR_NAMESPACE) \
								--docker-server "$(DOCKER_SERVER)" \
								--docker-username "$(DOCKER_USERNAME)" \
								--docker-password "$(DOCKER_PASSWORD)" \
								--docker-email="docker@dummy.com"
endif
ifneq ("$(or $(OCR_DOCKER_USERNAME),$(OCR_DOCKER_PASSWORD))","")
	@echo "Creating pull secrets for container-registry.oracle.com"
	kubectl create secret docker-registry ocr-k8s-operator-development-secret \
								--namespace $(OPERATOR_NAMESPACE) \
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
ifeq ($(CREATE_OPERATOR_NAMESPACE),true)
	@echo "Deleting test namespace $(OPERATOR_NAMESPACE)"
	kubectl delete namespace $(OPERATOR_NAMESPACE) --force --grace-period=0 && echo "deleted namespace" || true
endif
	kubectl delete clusterrole operator-test-coherence-operator --force --grace-period=0 && echo "deleted namespace" || true
	kubectl delete clusterrolebinding operator-test-coherence-operator --force --grace-period=0 && echo "deleted namespace" || true


# ----------------------------------------------------------------------------------------------------------------------
# Delete all resource from the test namespace
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: clean-namespace
clean-namespace: delete-coherence-clusters
	for i in $$(kubectl -n $(OPERATOR_NAMESPACE) get all -o name); do \
		echo "Deleting $${i} from test namespace $(OPERATOR_NAMESPACE)" \
		kubectl -n $(OPERATOR_NAMESPACE) delete $${i}; \
	done

# ----------------------------------------------------------------------------------------------------------------------
# Create the k8s secret to use in SSL/TLS testing.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: create-ssl-secrets
create-ssl-secrets: $(BUILD_OUTPUT)/certs
	@echo "Deleting SSL secret $(TEST_SSL_SECRET)"
	kubectl --namespace $(OPERATOR_NAMESPACE) delete secret $(TEST_SSL_SECRET) && echo "secret deleted" || true
	@echo "Creating SSL secret $(TEST_SSL_SECRET)"
	kubectl create secret generic $(TEST_SSL_SECRET) \
		--namespace $(OPERATOR_NAMESPACE) \
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
	mvn $(USE_MAVEN_SETTINGS) -B -f ./examples package -DskipTests -P docker

# ----------------------------------------------------------------------------------------------------------------------
# Build and test the examples
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: test-examples
test-examples: build-examples
	mvn $(USE_MAVEN_SETTINGS) -B -f ./examples verify

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
build-utils-image: build-mvn $(BUILD_BIN)/runner
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
# Build the Operator Test images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-test-images
build-test-images: build-mvn
	mvn $(USE_MAVEN_SETTINGS) -B -f java/operator-test package jib:dockerBuild -DskipTests -Djib.to.image=$(TEST_APPLICATION_IMAGE)
	mvn $(USE_MAVEN_SETTINGS) -B -f java/operator-test-spring package jib:dockerBuild -DskipTests -Djib.to.image=$(TEST_APPLICATION_IMAGE_SPRING)
	mvn $(USE_MAVEN_SETTINGS) -B -f java/operator-test-spring package spring-boot:build-image -DskipTests -Dcnbp-image-name=$(TEST_APPLICATION_IMAGE_SPRING_CNBP)
	docker build -f java/operator-test-spring/target/FatJar.Dockerfile -t $(TEST_APPLICATION_IMAGE_SPRING_FAT) java/operator-test-spring/target
	rm -rf java/operator-test-spring/target/spring || true && mkdir java/operator-test-spring/target/spring
	cp java/operator-test-spring/target/operator-test-spring-$(VERSION).jar java/operator-test-spring/target/spring/operator-test-spring-$(VERSION).jar
	cd java/operator-test-spring/target/spring && jar -xvf operator-test-spring-$(VERSION).jar && rm -f operator-test-spring-$(VERSION).jar
	docker build -f java/operator-test-spring/target/Dir.Dockerfile -t $(TEST_APPLICATION_IMAGE_SPRING) java/operator-test-spring/target

# ----------------------------------------------------------------------------------------------------------------------
# Push the Operator JIB Test Docker images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-test-images
push-test-images:
	docker push $(TEST_APPLICATION_IMAGE)
	docker push $(TEST_APPLICATION_IMAGE_SPRING)
	docker push $(TEST_APPLICATION_IMAGE_SPRING_FAT)
	docker push $(TEST_APPLICATION_IMAGE_SPRING_CNBP)

# ----------------------------------------------------------------------------------------------------------------------
# Build all of the Docker images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-all-images
build-all-images: $(BUILD_TARGETS)/build-operator build-utils-image build-test-images

# ----------------------------------------------------------------------------------------------------------------------
# Push all of the Docker images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-all-images
push-all-images: push-operator-image push-utils-image push-test-images

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
	go run -ldflags "$(LDFLAGS)" ./main.go --skip-service-suspend=true --coherence-dev-mode=true \
		--cert-type=self-signed --webhook-service=host.docker.internal \
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
run-debug: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
run-debug: export UTILS_IMAGE := $(UTILS_IMAGE)
run-debug:
	dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient \
		-- --skip-service-suspend=true --coherence-dev-mode=true \
		--cert-type=self-signed --webhook-service=host.docker.internal

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
	kind load docker-image --name operator $(TEST_APPLICATION_IMAGE_SPRING)|| true
	kind load docker-image --name operator $(TEST_APPLICATION_IMAGE_SPRING_FAT)|| true
	kind load docker-image --name operator $(TEST_APPLICATION_IMAGE_SPRING_CNBP)|| true
	kind load docker-image --name operator gcr.io/kubebuilder/kube-rbac-proxy:v0.5.0 || true

# ----------------------------------------------------------------------------------------------------------------------
# Install Prometheus
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-prometheus
install-prometheus:
	kubectl create ns $(OPERATOR_NAMESPACE) || true
	kubectl create -f hack/prometheus-rbac.yaml
	helm repo add stable https://kubernetes-charts.storage.googleapis.com/ || true
	@echo "Create Grafana Dashboards ConfigMap:"
	kubectl -n $(OPERATOR_NAMESPACE) create configmap coherence-grafana-dashboards --from-file=$(GRAFANA_DASHBOARDS)
	kubectl -n $(OPERATOR_NAMESPACE) label configmap coherence-grafana-dashboards grafana_dashboard=1
	@echo "Getting Helm Version:"
	helm version
	@echo "Installing stable/prometheus-operator:"
	helm install --atomic --namespace $(OPERATOR_NAMESPACE) --version $(PROMETHEUS_OPERATOR_VERSION) --wait \
		--set grafana.enabled=$(PROMETHEUS_INCLUDE_GRAFANA) \
		--values hack/prometheus-values.yaml prometheus stable/prometheus-operator
	@echo "Installing Prometheus instance:"
	kubectl -n $(OPERATOR_NAMESPACE) apply -f hack/prometheus.yaml
	sleep 10
	kubectl -n $(OPERATOR_NAMESPACE) wait --for=condition=Ready pod/prometheus-prometheus-0

# ----------------------------------------------------------------------------------------------------------------------
# Uninstall Prometheus
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: uninstall-prometheus
uninstall-prometheus:
	kubectl -n $(OPERATOR_NAMESPACE) delete -f hack/prometheus.yaml || true
	kubectl -n $(OPERATOR_NAMESPACE) delete configmap coherence-grafana-dashboards || true
	helm --namespace $(OPERATOR_NAMESPACE) delete prometheus || true
	kubectl delete -f hack/prometheus-rbac.yaml || true

.PHONY: install-prometheus-adapter
install-prometheus-adapter:
	kubectl create ns $(OPERATOR_NAMESPACE) || true
	helm repo add stable https://kubernetes-charts.storage.googleapis.com/ || true
	helm install --atomic --namespace $(OPERATOR_NAMESPACE) --version $(PROMETHEUS_ADAPTER_VERSION) --wait \
		--set prometheus.url=http://prometheus.$(OPERATOR_NAMESPACE).svc \
		--values hack/prometheus-adapter-values.yaml prometheus-adapter stable/prometheus-adapter

.PHONY: uninstall-prometheus-adapter
uninstall-prometheus-adapter:
	helm --namespace $(OPERATOR_NAMESPACE) delete prometheus-adapter || true

# ----------------------------------------------------------------------------------------------------------------------
# Start a port-forward process to the Grafana Pod.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: port-forward-grafana
port-forward-grafana: export GRAFANA_POD := $(shell kubectl -n $(OPERATOR_NAMESPACE) get pod -l app.kubernetes.io/name=grafana -o name)
port-forward-grafana:
	@echo "Reach Grafana on http://127.0.0.1:3000"
	@echo "User: admin Password: prom-operator"
	kubectl -n $(OPERATOR_NAMESPACE) port-forward $(GRAFANA_POD) 3000:3000

# ----------------------------------------------------------------------------------------------------------------------
# Install Elasticsearch & Kibana
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-elastic
install-elastic: helm-install-elastic kibana-import

.PHONY: helm-install-elastic
helm-install-elastic:
	kubectl create ns $(OPERATOR_NAMESPACE) || true
#   Create the ConfigMap containing the Coherence Kibana dashboards
	kubectl -n $(OPERATOR_NAMESPACE) delete secret coherence-kibana-dashboard || true
	kubectl -n $(OPERATOR_NAMESPACE) create secret generic --from-file dashboards/kibana/kibana-dashboard-data.json coherence-kibana-dashboard
#   Create the ConfigMap containing the Coherence Kibana dashboards import script
	kubectl -n $(OPERATOR_NAMESPACE) delete secret coherence-kibana-import || true
	kubectl -n $(OPERATOR_NAMESPACE) create secret generic --from-file hack/coherence-dashboard-import.sh coherence-kibana-import
#   Set-up the Elastic Helm repo
	@echo "Getting Helm Version:"
	helm version
	helm repo add elastic https://helm.elastic.co || true
#   Install Elasticsearch
	helm install --atomic --namespace $(OPERATOR_NAMESPACE) --version $(ELASTIC_VERSION) --wait --timeout=10m \
		--debug --values hack/elastic-values.yaml elasticsearch elastic/elasticsearch
#   Install Kibana
	helm install --atomic --namespace $(OPERATOR_NAMESPACE) --version $(ELASTIC_VERSION) --wait --timeout=10m \
		--debug --values hack/kibana-values.yaml kibana elastic/kibana \

.PHONY: kibana-import
kibana-import:
	KIBANA_POD=$$(kubectl -n $(OPERATOR_NAMESPACE) get pod -l app=kibana -o name) \
	; kubectl -n $(OPERATOR_NAMESPACE) exec -it $${KIBANA_POD} /bin/bash /usr/share/kibana/data/coherence/scripts/coherence-dashboard-import.sh

# ----------------------------------------------------------------------------------------------------------------------
# Uninstall Elasticsearch & Kibana
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: uninstall-elastic
uninstall-elastic:
	helm uninstall --namespace $(OPERATOR_NAMESPACE) kibana || true
	helm uninstall --namespace $(OPERATOR_NAMESPACE) elasticsearch || true
	kubectl -n $(OPERATOR_NAMESPACE) delete pvc elasticsearch-master-elasticsearch-master-0 || true

# ----------------------------------------------------------------------------------------------------------------------
# Start a port-forward process to the Kibana Pod.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: port-forward-kibana
port-forward-kibana: export KIBANA_POD := $(shell kubectl -n $(OPERATOR_NAMESPACE) get pod -l app=kibana -o name)
port-forward-kibana:
	@echo "Reach Kibana on http://127.0.0.1:5601"
	kubectl -n $(OPERATOR_NAMESPACE) port-forward $(KIBANA_POD) 5601:5601

# ----------------------------------------------------------------------------------------------------------------------
# Start a port-forward process to the Elasticsearch Pod.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: port-forward-es
port-forward-es: export ES_POD := $(shell kubectl -n $(OPERATOR_NAMESPACE) get pod -l app=elasticsearch-master -o name)
port-forward-es:
	@echo "Reach Elasticsearch on http://127.0.0.1:9200"
	kubectl -n $(OPERATOR_NAMESPACE) port-forward $(ES_POD) 9200:9200


# ----------------------------------------------------------------------------------------------------------------------
# Delete all of the Coherence resources from the test namespace.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: delete-coherence-clusters
delete-coherence-clusters:
	for i in $$(kubectl -n  $(OPERATOR_NAMESPACE) get coherence.coherence.oracle.com -o name); do \
  		kubectl -n $(OPERATOR_NAMESPACE) patch $${i} -p '{"metadata":{"finalizers":[]}}' --type=merge || true ;\
		kubectl -n $(OPERATOR_NAMESPACE) delete $${i}; \
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
	$(BUILD_BIN)/golangci-lint run -v --timeout=5m --skip-files=zz_.*,generated/* ./api/... ./controllers/... ./pkg/... ./runner/... ./converter/...


# ----------------------------------------------------------------------------------------------------------------------
# Performs a copyright check.
# To add exclusions add the file or folder pattern using the -X parameter.
# Add directories to be scanned at the end of the parameter list.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: copyright
copyright:
	@java -cp hack/glassfish-copyright-maven-plugin-2.1.jar \
	  org.glassfish.copyright.Copyright -C hack/copyright.txt \
	  -X .adoc \
	  -X bin/ \
	  -X build/_output/ \
	  -X clientset/ \
	  -X dashboards/grafana/ \
	  -X dashboards/kibana/ \
	  -X /Dockerfile \
	  -X .Dockerfile \
	  -X docs/ \
	  -X examples/.mvn/ \
	  -X .factories \
	  -X hack/copyright.txt \
	  -X hack/intellij-codestyle.xml \
	  -X hack/sdk/ \
	  -X go.mod \
	  -X go.sum \
	  -X HEADER.txt \
	  -X helm-charts/coherence-operator/templates/NOTES.txt \
	  -X .iml \
	  -X java/src/copyright/EXCLUDE.txt \
	  -X Jenkinsfile \
	  -X .jar \
	  -X .jks \
	  -X .json \
	  -X LICENSE.txt \
	  -X Makefile \
	  -X .md \
	  -X mvnw \
	  -X mvnw.cmd \
	  -X .png \
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
code-review: generate golangci copyright
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
# Pre-Release Tasks
# Update the version numbers in the documentation to be the version about to be released
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: pre-release
pre-release:
	sed -i '' 's/$(subst .,\.,$(PREV_VERSION))/$(VERSION)/g' README.md
	find docs \( -name '*.adoc' -o -name '*.md' \) -exec sed -i '' 's/$(subst .,\.,$(PREV_VERSION))/$(VERSION)/g' {} +
	find examples \( -name '*.adoc' -o -name '*.md' -o -name '*.yaml' \) -exec sed -i '' 's/$(subst .,\.,$(PREV_VERSION))/$(VERSION)/g' {} +

# ----------------------------------------------------------------------------------------------------------------------
# Post-Release Tasks
# Update the version numbers
#post-release: check-new-version new-version manifests generate build-all-images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: post-release
post-release: check-new-version new-version

.PHONY: check-new-version
check-new-version:
ifeq (, $(NEW_VERSION))
	@echo "You must specify the NEW_VERSION parameter"
	exit 1
else
	@echo "Updating version from $(VERSION) to $(NEW_VERSION)"
endif

.PHONY: new-version
new-version:
	sed -i '' 's/$(subst .,\.,$(PREV_VERSION))/$(VERSION)/g' MAKEFILE
	sed -i '' 's/$(subst .,\.,$(VERSION))/$(NEW_VERSION)/g' MAKEFILE
	find config \( -name '*.yaml' -o -name '*.json' \) -exec sed -i '' 's/$(subst .,\.,$(VERSION))/$(NEW_VERSION)/g' {} +
	find java \( -name 'pom.xml' \) -exec sed -i '' 's/<version>$(subst .,\.,$(VERSION))<\/version>/<version>$(NEW_VERSION)<\/version>/g' {} +


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
release-ghpages:  helm-chart docs
	@echo "Releasing Dashboards $(VERSION)"
	mkdir -p $(BUILD_OUTPUT)/dashboards/$(VERSION) || true
	tar -czvf $(BUILD_OUTPUT)/dashboards/$(VERSION)/coherence-dashboards.tar.gz  dashboards/
	kubectl create configmap coherence-grafana-dashboards --from-file=dashboards/grafana \
		--dry-run -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION)/coherence-grafana-dashboards.yaml
	kubectl create configmap coherence-grafana-dashboards --from-file=dashboards/grafana-legacy \
		--dry-run -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION)/coherence-grafana-legacy-dashboards.yaml
	kubectl create configmap coherence-kibana-dashboards --from-file=dashboards/kibana \
		--dry-run -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION)/coherence-kibana-dashboards.yaml
	git stash save --keep-index --include-untracked || true
	git stash drop || true
	git checkout --track origin/gh-pages
	git config pull.rebase true
	git pull
	mkdir -p dashboards || true
	rm -rf dashboards/$(VERSION) || true
	mv $(BUILD_OUTPUT)/dashboards/$(VERSION)/ dashboards/
	git add dashboards/$(VERSION)/*
ifeq (true, $(PRE_RELEASE))
	mkdir -p docs-unstable || true
	rm -rf docs-unstable/$(VERSION)/ || true
	mv $(BUILD_OUTPUT)/docs/ docs-unstable/$(VERSION)/
	ls -ls docs-unstable

	mkdir -p charts-unstable || true
	cp $(BUILD_HELM)/coherence-operator-$(VERSION).tgz charts-unstable/
	helm repo index charts-unstable --url https://oracle.github.io/coherence-operator/charts-unstable
	git add charts-unstable/coherence-operator-$(VERSION).tgz
	git add charts-unstable/index.yaml
	ls -ls charts-unstable

	git add docs-unstable/*
	git status
else
	mkdir docs/$(VERSION) || true
	rm -rf docs/$(VERSION)/ || true
	mv $(BUILD_OUTPUT)/docs/ docs/$(VERSION)/
	ls -ls docs

	mkdir -p charts || true
	cp $(BUILD_HELM)/coherence-operator-$(VERSION).tgz charts/
	helm repo index charts --url https://oracle.github.io/coherence-operator/charts
	git add charts/coherence-operator-$(VERSION).tgz
	git add charts/index.yaml
	ls -ls charts

	git add docs/*
	git status
endif
	git clean -d -f
	git status
	git commit -m "adding Coherence Operator Helm chart and docs version: $(VERSION)"
	git log -1
ifeq (true, $(RELEASE_DRY_RUN))
	@echo "release dry-run - would have pushed Helm chart and docs $(VERSION) to gh-pages"
else
	git push origin gh-pages
endif


# ----------------------------------------------------------------------------------------------------------------------
# Release the Coherence Operator.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: release
release:
ifeq (true, $(RELEASE_DRY_RUN))
release: build-all-images release-ghpages
	@echo "release dry-run: would have pushed images"
else
release: build-all-images release-ghpages push-release-images
endif


# ----------------------------------------------------------------------------------------------------------------------
# List all of the targets in the Makefile
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: list
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
