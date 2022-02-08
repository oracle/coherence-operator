# ----------------------------------------------------------------------------------------------------------------------
# Copyright (c) 2019, 2022, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
# ----------------------------------------------------------------------------------------------------------------------
# This is the Makefile to build the Coherence Kubernetes Operator.
# ----------------------------------------------------------------------------------------------------------------------

# ======================================================================================================================
# Makefile Variables
#
# The following section contains all of the variables and properties used by other targets in the Makefile
# to set things like build directories, version numbers etc.
# ======================================================================================================================

# The version of the Operator being build - this should be a valid SemVer format
VERSION ?= 3.2.5
MVN_VERSION ?= $(VERSION)-SNAPSHOT

# The version number to be replaced by this release
PREV_VERSION ?= 3.2.4

# The operator version to use to run certification tests against
CERTIFICATION_VERSION ?= $(VERSION)

# The previous Operator version used to run the compatibility tests.
COMPATIBLE_VERSION  = 3.2.4
# The selector to use to find Operator Pods of the COMPATIBLE_VERSION (do not put in double quotes!!)
COMPATIBLE_SELECTOR = control-plane=coherence

# The GitHub project URL
PROJECT_URL = https://github.com/oracle/coherence-operator

# ----------------------------------------------------------------------------------------------------------------------
# The Coherence image to use for deployments that do not specify an image
# ----------------------------------------------------------------------------------------------------------------------
COHERENCE_VERSION ?= 21.12.1
COHERENCE_IMAGE ?= ghcr.io/oracle/coherence-ce:21.12.1
# This is the Coherence image that will be used in tests.
# Changing this variable will allow test builds to be run against different Coherence versions
# without altering the default image name.
TEST_COHERENCE_IMAGE ?= $(COHERENCE_IMAGE)
TEST_COHERENCE_VERSION ?= $(COHERENCE_VERSION)
TEST_COHERENCE_GID ?= com.oracle.coherence.ce

# The current working directory
CURRDIR := $(shell pwd)

# ----------------------------------------------------------------------------------------------------------------------
# By default we target amd64 as this is by far the most common local build environment
# We actually build images for amd64 and arm64
# ----------------------------------------------------------------------------------------------------------------------
IMAGE_ARCH      ?= amd64
ARCH            ?= amd64
OS              ?= linux
UNAME_S         := $(shell uname -s)
GOPROXY         ?= https://proxy.golang.org

# ----------------------------------------------------------------------------------------------------------------------
# Set the location of the Operator SDK executable
# ----------------------------------------------------------------------------------------------------------------------
UNAME_S               = $(shell uname -s)
UNAME_M               = $(shell uname -m)
OPERATOR_SDK_VERSION := v1.9.0

# ----------------------------------------------------------------------------------------------------------------------
# Options to append to the Maven command
# ----------------------------------------------------------------------------------------------------------------------
MAVEN_OPTIONS ?= -Dmaven.wagon.httpconnectionManager.ttlSeconds=25 -Dmaven.wagon.http.retryHandler.count=3
MAVEN_BUILD_OPTS :=$(USE_MAVEN_SETTINGS) -Drevision=$(MVN_VERSION) -Dcoherence.version=$(COHERENCE_VERSION) $(MAVEN_OPTIONS)

# ----------------------------------------------------------------------------------------------------------------------
# Operator image names
# ----------------------------------------------------------------------------------------------------------------------
RELEASE_IMAGE_PREFIX   ?= ghcr.io/oracle/
OPERATOR_IMAGE_REPO    := $(RELEASE_IMAGE_PREFIX)coherence-operator
OPERATOR_IMAGE         := $(OPERATOR_IMAGE_REPO):$(VERSION)
OPERATOR_IMAGE_DELVE   := $(OPERATOR_IMAGE_REPO):delve
OPERATOR_IMAGE_DEBUG   := $(OPERATOR_IMAGE_REPO):debug
UTILS_IMAGE            ?= $(OPERATOR_IMAGE_REPO):$(VERSION)-utils
TEST_BASE_IMAGE        ?= $(OPERATOR_IMAGE_REPO):$(VERSION)-test-base
# The Operator images to push
OPERATOR_RELEASE_REPO   ?= $(OPERATOR_IMAGE_REPO)
OPERATOR_RELEASE_IMAGE  := $(OPERATOR_RELEASE_REPO):$(VERSION)
UTILS_RELEASE_IMAGE     := $(OPERATOR_RELEASE_REPO):$(VERSION)-utils
TEST_BASE_RELEASE_IMAGE := $(OPERATOR_RELEASE_REPO):$(VERSION)-test-base
BUNDLE_RELEASE_IMAGE    := $(OPERATOR_RELEASE_REPO):$(VERSION)-bundle

GPG_PASSPHRASE :=

# ----------------------------------------------------------------------------------------------------------------------
# The test application images used in integration tests
# ----------------------------------------------------------------------------------------------------------------------
TEST_APPLICATION_IMAGE             := $(RELEASE_IMAGE_PREFIX)operator-test:1.0.0
TEST_COMPATIBILITY_IMAGE           := $(RELEASE_IMAGE_PREFIX)operator-test-compatibility:1.0.0
TEST_APPLICATION_IMAGE_CLIENT      := $(RELEASE_IMAGE_PREFIX)operator-test-client:1.0.0
TEST_APPLICATION_IMAGE_HELIDON     := $(RELEASE_IMAGE_PREFIX)operator-test-helidon:1.0.0
TEST_APPLICATION_IMAGE_SPRING      := $(RELEASE_IMAGE_PREFIX)operator-test-spring:1.0.0
TEST_APPLICATION_IMAGE_SPRING_FAT  := $(RELEASE_IMAGE_PREFIX)operator-test-spring-fat:1.0.0
TEST_APPLICATION_IMAGE_SPRING_CNBP := $(RELEASE_IMAGE_PREFIX)operator-test-spring-cnbp:1.0.0

# ----------------------------------------------------------------------------------------------------------------------
# Operator Lifecycle Manager properties
# ----------------------------------------------------------------------------------------------------------------------
# CHANNELS define the bundle channels used in the Operator Lifecycle Manager bundle.
CHANNELS := stable
# Add a new line here if you would like to change its default config. (E.g CHANNELS = "preview,fast,stable")
# To re-generate a bundle for other specific channels without changing the standard setup, you can:
# - use the CHANNELS as arg of the bundle target (e.g make bundle CHANNELS=preview,fast,stable)
# - use environment variables to overwrite this value (e.g export CHANNELS="preview,fast,stable")
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif

# DEFAULT_CHANNEL defines the default channel used in the bundle.
DEFAULT_CHANNEL := stable
# Add a new line here if you would like to change its default config. (E.g DEFAULT_CHANNEL = "stable")
# To re-generate a bundle for any other default channel without changing the default setup, you can:
# - use the DEFAULT_CHANNEL as arg of the bundle target (e.g make bundle DEFAULT_CHANNEL=stable)
# - use environment variables to overwrite this value (e.g export DEFAULT_CHANNEL="stable")
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

# BUNDLE_IMG defines the image:tag used for the bundle.
# You can use it as an arg. (E.g make bundle-build BUNDLE_IMG=<some-registry>/<project-name-bundle>:<tag>)
BUNDLE_IMG ?= $(OPERATOR_IMAGE_REPO)-bundle:$(VERSION)

# ----------------------------------------------------------------------------------------------------------------------
# Release build options
# ----------------------------------------------------------------------------------------------------------------------
RELEASE_DRY_RUN  ?= true
PRE_RELEASE      ?= true

# ----------------------------------------------------------------------------------------------------------------------
# Testing properties
# ----------------------------------------------------------------------------------------------------------------------
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

# default as in test/e2e/helper/proj_helpers.go
OPERATOR_NAMESPACE ?= operator-test
# the test client namespace
OPERATOR_NAMESPACE_CLIENT ?= operator-test-client
# the optional namespaces the operator should watch
WATCH_NAMESPACE ?=
# flag indicating whether the test namespace should be reset (deleted and recreated) before tests
CREATE_OPERATOR_NAMESPACE ?= true

# restart local storage for persistence
LOCAL_STORAGE_RESTART ?= false

# ----------------------------------------------------------------------------------------------------------------------
# Env variables used to create pull secrets
# This is required if building and testing in environments that need to pull or push
# images to private registries. For example building and testing with k8s in OCI.
# ----------------------------------------------------------------------------------------------------------------------
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

# Configure the image pull secrets information.
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

# ----------------------------------------------------------------------------------------------------------------------
# Build output directories
# ----------------------------------------------------------------------------------------------------------------------
override BUILD_OUTPUT        := $(CURRDIR)/build/_output
override BUILD_ASSETS        := pkg/data/assets
override BUILD_BIN           := $(CURRDIR)/bin
override BUILD_DEPLOY        := $(BUILD_OUTPUT)/config
override BUILD_HELM          := $(BUILD_OUTPUT)/helm-charts
override BUILD_MANIFESTS     := $(BUILD_OUTPUT)/manifests
override BUILD_MANIFESTS_PKG := $(BUILD_OUTPUT)/coherence-operator-manifests.tar.gz
override BUILD_PROPS         := $(BUILD_OUTPUT)/build.properties
override BUILD_TARGETS       := $(BUILD_OUTPUT)/targets
override TEST_LOGS_DIR       := $(BUILD_OUTPUT)/test-logs


# ----------------------------------------------------------------------------------------------------------------------
# Set the location of various build tools
# ----------------------------------------------------------------------------------------------------------------------
TOOLS_DIRECTORY   = $(CURRDIR)/build/tools
TOOLS_BIN         = $(TOOLS_DIRECTORY)/bin
OPERATOR_SDK_HOME = $(TOOLS_DIRECTORY)/sdk/$(UNAME_S)-$(UNAME_M)
OPERATOR_SDK      = $(OPERATOR_SDK_HOME)/operator-sdk
PROMETHEUS_HOME   = $(TOOLS_DIRECTORY)/prometheus

# ----------------------------------------------------------------------------------------------------------------------
# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
# ----------------------------------------------------------------------------------------------------------------------
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# ----------------------------------------------------------------------------------------------------------------------
# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
# ----------------------------------------------------------------------------------------------------------------------
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# ----------------------------------------------------------------------------------------------------------------------
# Capture the Git commit to add to the build information that is then embedded in the Go binary
# ----------------------------------------------------------------------------------------------------------------------
GITCOMMIT         ?= $(shell git rev-list -1 HEAD)
GITREPO           := https://github.com/oracle/coherence-operator.git
SOURCE_DATE_EPOCH := $(shell git show -s --format=format:%ct HEAD)
DATE_FMT          := "%Y-%m-%dT%H:%M:%SZ"
BUILD_DATE        := $(shell date -u -d "@$SOURCE_DATE_EPOCH" "+${DATE_FMT}" 2>/dev/null || date -u -r "${SOURCE_DATE_EPOCH}" "+${DATE_FMT}" 2>/dev/null || date -u "+${DATE_FMT}")
BUILD_USER        := $(shell whoami)

LDFLAGS          = -X main.Version=$(VERSION) -X main.Commit=$(GITCOMMIT) -X main.Date=$(BUILD_DATE) -X main.Author=$(BUILD_USER)
GOS              = $(shell find . -type f -name "*.go" ! -name "*_test.go")
HELM_FILES       = $(shell find helm-charts/coherence-operator -type f)
API_GO_FILES     = $(shell find . -type f -name "*.go" ! -name "*_test.go"  ! -name "zz*.go")
CRDV1_FILES      = $(shell find ./config/crd -type f)

TEST_SSL_SECRET := coherence-ssl-secret

# ----------------------------------------------------------------------------------------------------------------------
# Prometheus Operator settings (used in integration tests)
# ----------------------------------------------------------------------------------------------------------------------
PROMETHEUS_VERSION           ?= v0.8.0
PROMETHEUS_NAMESPACE         ?= monitoring
PROMETHEUS_ADAPTER_VERSION   ?= 2.5.0
GRAFANA_DASHBOARDS           ?= dashboards/grafana/

# ----------------------------------------------------------------------------------------------------------------------
# Elasticsearch & Kibana settings (used in integration tests)
# ----------------------------------------------------------------------------------------------------------------------
ELASTIC_VERSION ?= 7.6.2
KIBANA_INDEX_PATTERN := "6abb1220-3feb-11e9-a9a3-4b1c09db6e6a"

# ----------------------------------------------------------------------------------------------------------------------
# MetalLB load balancer settings
# ----------------------------------------------------------------------------------------------------------------------
METALLB_VERSION ?= v0.10.2

# ----------------------------------------------------------------------------------------------------------------------
# Istio settings
# ----------------------------------------------------------------------------------------------------------------------
# The version of Istio to install, leave empty for the latest
ISTIO_VERSION ?=

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

.PHONY: all
all: java-client build-all-images helm-chart ## Build all the Coherence Operator artefacts and images

# ----------------------------------------------------------------------------------------------------------------------
# Configure the build properties
# ----------------------------------------------------------------------------------------------------------------------
$(BUILD_PROPS):
	# Ensures that build output directories exist
	@echo "Creating build directories"
	@mkdir -p $(BUILD_OUTPUT)
	@mkdir -p $(BUILD_BIN)
	@mkdir -p $(BUILD_DEPLOY)
	@mkdir -p $(BUILD_HELM)
	@mkdir -p $(BUILD_MANIFESTS)
	@mkdir -p $(BUILD_TARGETS)
	@mkdir -p $(TEST_LOGS_DIR)
	@mkdir -p $(TOOLS_BIN)
	# create build.properties
	rm -f $(BUILD_PROPS)
	printf "COHERENCE_IMAGE=$(COHERENCE_IMAGE)\n\
	UTILS_IMAGE=$(UTILS_IMAGE)\n\
	OPERATOR_IMAGE=$(OPERATOR_IMAGE)\n\
	VERSION=$(VERSION)\n" > $(BUILD_PROPS)

# ----------------------------------------------------------------------------------------------------------------------
# Clean-up all of the build artifacts
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: clean
clean: ## Cleans the build 
	-rm -rf $(BUILD_OUTPUT)
	-rm -rf $(BUILD_BIN)
	-rm -rf bundle
	rm pkg/data/zz_generated_*.go || true
	./mvnw -f java clean $(MAVEN_BUILD_OPTS)
	./mvnw -f examples clean $(MAVEN_BUILD_OPTS)

# ----------------------------------------------------------------------------------------------------------------------
# Clean-up all of the locally downloaded build tools
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: clean-tools
clean-tools: ## Cleans the locally downloaded build tools (i.e. need a new tool version)
	-rm -rf $(TOOLS_BIN)


# ----------------------------------------------------------------------------------------------------------------------
# Builds the Operator
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-operator
build-operator: $(BUILD_TARGETS)/build-operator ## Build the Coherence Operator image

# ----------------------------------------------------------------------------------------------------------------------
# Build the Operator Image
# ----------------------------------------------------------------------------------------------------------------------
#   We copy the Dockerfile to $(BUILD_OUTPUT) only so that we can use it as a conditional build dependency in this Makefile
$(BUILD_TARGETS)/build-operator: $(BUILD_BIN)/manager $(BUILD_BIN)/runner
	docker build --no-cache --build-arg version=$(VERSION) \
		--build-arg coherence_image=$(COHERENCE_IMAGE) \
		--build-arg utils_image=$(UTILS_IMAGE) \
		--build-arg target=amd64 \
		. -t $(OPERATOR_IMAGE)-amd64
	docker build --no-cache --build-arg version=$(VERSION) \
		--build-arg coherence_image=$(COHERENCE_IMAGE) \
		--build-arg utils_image=$(UTILS_IMAGE) \
		--build-arg target=arm64 \
		. -t $(OPERATOR_IMAGE)-arm64
	docker tag $(OPERATOR_IMAGE)-$(IMAGE_ARCH) $(OPERATOR_IMAGE)
	touch $(BUILD_TARGETS)/build-operator

.PHONY: build-operator-debug
build-operator-debug: $(BUILD_BIN)/linux/amd64/manager-debug ## Build the Coherence Operator image with the Delve debugger
	docker build --no-cache --build-arg version=$(VERSION) \
		--build-arg BASE_IMAGE=$(OPERATOR_IMAGE_DELVE) \
		--build-arg coherence_image=$(COHERENCE_IMAGE) \
		--build-arg utils_image=$(UTILS_IMAGE) \
		--build-arg target=amd64 \
		-f debug/Dockerfile \
		. -t $(OPERATOR_IMAGE_DEBUG)

build-delve-image: ## Build the Coherence Operator Delve debugger base image
	docker build -f debug/Base.Dockerfile -t $(OPERATOR_IMAGE_DELVE) debug

$(BUILD_BIN)/linux/amd64/manager-debug: $(BUILD_PROPS) $(GOS) $(BUILD_TARGETS)/generate $(BUILD_TARGETS)/manifests
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -gcflags "-N -l" -ldflags "$(LDFLAGS)" -a -o $(BUILD_BIN)/linux/amd64/manager-debug main.go
	chmod +x $(BUILD_BIN)/linux/amd64/manager-debug

# ----------------------------------------------------------------------------------------------------------------------
# Build the Operator Utils Docker image
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-utils
build-utils: build-mvn $(BUILD_BIN)/runner  ## Build the Coherence Operator utils image
	cp -R $(BUILD_BIN)/linux  java/coherence-operator/target/docker
	docker build --no-cache --build-arg target=amd64 -t $(UTILS_IMAGE)-amd64 java/coherence-operator/target/docker
	docker build --no-cache --build-arg target=arm64 -t $(UTILS_IMAGE)-arm64 java/coherence-operator/target/docker
	docker tag $(UTILS_IMAGE)-$(IMAGE_ARCH) $(UTILS_IMAGE)

# ----------------------------------------------------------------------------------------------------------------------
# Build the Operator base test image
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-test-base
build-test-base: export ARTIFACT_DIR        := $(CURRDIR)/java/coherence-operator
build-test-base: export VERSION             := $(VERSION)
build-test-base: export IMAGE_NAME          := $(TEST_BASE_IMAGE)
build-test-base: export AMD_BASE_IMAGE      := gcr.io/distroless/java11
build-test-base: export ARM_BASE_IMAGE      := gcr.io/distroless/java11
build-test-base: export PROJECT_URL         := $(PROJECT_URL)
build-test-base: export PROJECT_VENDOR      := Oracle
build-test-base: export PROJECT_DESCRIPTION := Oracle Coherence base test image
build-test-base: build-mvn $(BUILD_BIN)/runner  ## Build the Coherence test base image
	cp -R $(BUILD_BIN)/linux  java/coherence-operator/target/docker
	$(CURRDIR)/java/coherence-operator/run-buildah.sh BUILD

# ----------------------------------------------------------------------------------------------------------------------
# Build the Operator images without the test images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-operator-images
build-operator-images: $(BUILD_TARGETS)/build-operator build-utils build-test-base ## Build all operator images

# ----------------------------------------------------------------------------------------------------------------------
# Build the Operator Test images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-test-images
build-test-images: build-mvn build-client-image build-basic-test-image ## Build all of the test images
	./mvnw -B -f java/operator-test-helidon package jib:dockerBuild -DskipTests -Djib.to.image=$(TEST_APPLICATION_IMAGE_HELIDON) $(MAVEN_BUILD_OPTS)
	./mvnw -B -f java/operator-test-spring package jib:dockerBuild -DskipTests -Djib.to.image=$(TEST_APPLICATION_IMAGE_SPRING) $(MAVEN_BUILD_OPTS)
	./mvnw -B -f java/operator-test-spring package spring-boot:build-image -DskipTests -Dcnbp-image-name=$(TEST_APPLICATION_IMAGE_SPRING_CNBP) $(MAVEN_BUILD_OPTS)
	docker build -f java/operator-test-spring/target/FatJar.Dockerfile -t $(TEST_APPLICATION_IMAGE_SPRING_FAT) java/operator-test-spring/target
	rm -rf java/operator-test-spring/target/spring || true && mkdir java/operator-test-spring/target/spring
	cp java/operator-test-spring/target/operator-test-spring-$(MVN_VERSION).jar java/operator-test-spring/target/spring/operator-test-spring-$(MVN_VERSION).jar
	cd java/operator-test-spring/target/spring && jar -xvf operator-test-spring-$(MVN_VERSION).jar && rm -f operator-test-spring-$(MVN_VERSION).jar
	docker build -f java/operator-test-spring/target/Dir.Dockerfile -t $(TEST_APPLICATION_IMAGE_SPRING) java/operator-test-spring/target

# ----------------------------------------------------------------------------------------------------------------------
# Build the basic Operator Test image
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-basic-test-image
build-basic-test-image: build-mvn ## Build the basic Operator test image
	./mvnw -B -f java/operator-test package jib:dockerBuild -DskipTests -Djib.to.image=$(TEST_APPLICATION_IMAGE) $(MAVEN_BUILD_OPTS)

.PHONY: build-client-image
build-client-image: ## Build the test client image
	./mvnw -B -f java/operator-test-client package jib:dockerBuild -DskipTests -Djib.to.image=$(TEST_APPLICATION_IMAGE_CLIENT) $(MAVEN_BUILD_OPTS)

# ----------------------------------------------------------------------------------------------------------------------
# Build all of the Docker images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-all-images
build-all-images: $(BUILD_TARGETS)/build-operator build-utils build-test-base build-test-images ## Build all images (including tests)

# ----------------------------------------------------------------------------------------------------------------------
# Build the operator linux binary
# ----------------------------------------------------------------------------------------------------------------------
$(BUILD_BIN)/manager: $(BUILD_PROPS) $(GOS) $(BUILD_TARGETS)/generate $(BUILD_TARGETS)/manifests
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -trimpath -ldflags "$(LDFLAGS)" -a -o $(BUILD_BIN)/manager main.go
	mkdir -p $(BUILD_BIN)/linux/amd64 || true
	cp -f $(BUILD_BIN)/manager $(BUILD_BIN)/linux/amd64/manager
	mkdir -p $(BUILD_BIN)/linux/arm64 || true
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GO111MODULE=on go build -trimpath -ldflags "$(LDFLAGS)" -a -o $(BUILD_BIN)/linux/arm64/manager main.go

# ----------------------------------------------------------------------------------------------------------------------
# Ensure Operator SDK is at the correct version
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: ensure-sdk
ensure-sdk:
	@echo "Ensuring Operator SDK is present at version $(OPERATOR_SDK_VERSION)"
	./hack/ensure-sdk.sh $(OPERATOR_SDK_VERSION) $(OPERATOR_SDK_HOME)

# ----------------------------------------------------------------------------------------------------------------------
# Internal make step that builds the Operator runner artifacts utility
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-runner
build-runner: $(BUILD_BIN)/runner  ## Build the Coherence Operator runner binary

$(BUILD_BIN)/runner: $(BUILD_PROPS) $(GOS)
	@echo "Building Operator Runner"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -trimpath -ldflags "$(LDFLAGS)" -o $(BUILD_BIN)/runner ./runner
	mkdir -p $(BUILD_BIN)/linux/amd64 || true
	cp -f $(BUILD_BIN)/runner $(BUILD_BIN)/linux/amd64/runner
	mkdir -p $(BUILD_BIN)/linux/arm64 || true
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GO111MODULE=on go build -trimpath -ldflags "$(LDFLAGS)" -a -o $(BUILD_BIN)/linux/arm64/runner ./runner

# ----------------------------------------------------------------------------------------------------------------------
# Build the Java artifacts
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-mvn
build-mvn: ## Build the Java artefacts
	./mvnw -B -f java clean install -DskipTests $(MAVEN_BUILD_OPTS)

# ----------------------------------------------------------------------------------------------------------------------
# Build Java client
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: java-client
java-client: $(BUILD_PROPS) $(BUILD_TARGETS)/generate $(BUILD_TARGETS)/manifests prepare-deploy $(BUILD_OUTPUT)/java-client/java/gen/pom.xml build-mvn

# ---------------------------------------------------------------------------
# Build the Coherence operator Helm chart and package it into a tar.gz
# ---------------------------------------------------------------------------
.PHONY: helm-chart
helm-chart: $(BUILD_PROPS) $(BUILD_HELM)/coherence-operator-$(VERSION).tgz   ## Build the Coherence Operator Helm chart

$(BUILD_HELM)/coherence-operator-$(VERSION).tgz: $(BUILD_PROPS) $(HELM_FILES) $(BUILD_TARGETS)/generate $(BUILD_TARGETS)/manifests kustomize
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

# ======================================================================================================================
# General development related targets
# ======================================================================================================================
##@ Development

# ----------------------------------------------------------------------------------------------------------------------
# Generate manifests e.g. CRD, RBAC etc.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: manifests
manifests: $(BUILD_TARGETS)/manifests ## Generate the CustomResourceDefinition and other yaml manifests.

$(BUILD_TARGETS)/manifests: $(BUILD_PROPS) config/crd/bases/coherence.oracle.com_coherence.yaml docs/about/04_coherence_spec.adoc $(BUILD_MANIFESTS_PKG)
	touch $(BUILD_TARGETS)/manifests

config/crd/bases/coherence.oracle.com_coherence.yaml: kustomize $(API_GO_FILES) controller-gen
	$(CONTROLLER_GEN) "crd:crdVersions={v1}" \
	  rbac:roleName=manager-role paths="{./api/...,./controllers/...}" \
	  output:crd:artifacts:config=config/crd/bases
	cd config/crd && $(KUSTOMIZE) edit add label "app.kubernetes.io/version:$(VERSION)" -f
	$(KUSTOMIZE) build config/crd > $(BUILD_ASSETS)/crd_v1.yaml

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
	cp $(BUILD_OUTPUT)/config.json $(BUILD_ASSETS)/config.json

# ----------------------------------------------------------------------------------------------------------------------
# Generate code, configuration and docs.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: generate
generate: $(BUILD_TARGETS)/generate  ## Run Kubebuilder code and configuration generation

$(BUILD_TARGETS)/generate: $(BUILD_PROPS) $(BUILD_OUTPUT)/config.json api/v1/zz_generated.deepcopy.go
	touch $(BUILD_TARGETS)/generate

api/v1/zz_generated.deepcopy.go: $(API_GO_FILES) controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./api/..."

# ----------------------------------------------------------------------------------------------------------------------
# Generate API docs
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: api-doc-gen
api-doc-gen: docs/about/04_coherence_spec.adoc  ## Generate API documentation

docs/about/04_coherence_spec.adoc: $(API_GO_FILES)
	@echo "Generating CRD Doc"
	go run ./docgen/ \
		api/v1/coherenceresourcespec_types.go \
		api/v1/coherence_types.go \
		api/v1/coherenceresource_types.go \
		> docs/about/04_coherence_spec.adoc

# ----------------------------------------------------------------------------------------------------------------------
# Generate the keys and certs used in tests.
# ----------------------------------------------------------------------------------------------------------------------
$(BUILD_OUTPUT)/certs:
	@echo "Generating test keys and certs"
	./hack/keys.sh

# ----------------------------------------------------------------------------------------------------------------------
# Executes the code review targets.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: code-review
code-review: export MAVEN_USER := $(MAVEN_USER)
code-review: export MAVEN_PASSWORD := $(MAVEN_PASSWORD)
code-review: $(BUILD_TARGETS)/generate golangci copyright  ## Full code review and Checkstyle test
	./mvnw -B -f java validate -DskipTests -P checkstyle $(MAVEN_BUILD_OPTS)

# ----------------------------------------------------------------------------------------------------------------------
# Executes golangci-lint to perform various code review checks on the source.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: golangci
golangci: $(TOOLS_BIN)/golangci-lint ## Go code review
	$(TOOLS_BIN)/golangci-lint run -v --timeout=5m --exclude='G402:' --skip-dirs=.*/fakes --skip-files=zz_.*,generated/*,pkg/data/assets... ./api/... ./controllers/... ./pkg/... ./runner/...
	$(TOOLS_BIN)/golangci-lint run -v --timeout=5m --exclude='G107:' --exclude='should not use dot imports' ./test/... ./pkg/fakes/...


# ----------------------------------------------------------------------------------------------------------------------
# Performs a copyright check.
# To add exclusions add the file or folder pattern using the -X parameter.
# Add directories to be scanned at the end of the parameter list.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: copyright
copyright:  ## Check copyright headers
	@java -cp hack/glassfish-copyright-maven-plugin-2.1.jar \
	  org.glassfish.copyright.Copyright -C hack/copyright.txt \
	  -X .adoc \
	  -X bin/ \
	  -X build/ \
	  -X clientset/ \
	  -X dashboards/ \
	  -X /Dockerfile \
	  -X .Dockerfile \
	  -X docs/ \
	  -X examples/.mvn/ \
	  -X examples/helm/chart/templates/NOTES.txt \
	  -X .factories \
	  -X hack/copyright.txt \
	  -X hack/intellij-codestyle.xml \
	  -X hack/istio- \
	  -X hack/sdk/ \
	  -X go.mod \
	  -X go.sum \
	  -X .gradle/ \
	  -X gradle/ \
	  -X gradlew \
	  -X gradlew.bat \
	  -X HEADER.txt \
	  -X helm-charts/coherence-operator/templates/NOTES.txt \
	  -X .iml \
	  -X java/certs/ \
	  -X java/src/copyright/EXCLUDE.txt \
	  -X Jenkinsfile \
	  -X .jar \
	  -X jib-cache/ \
	  -X .jks \
	  -X .json \
	  -X LICENSE.txt \
	  -X Makefile \
	  -X .md \
	  -X meta/ \
	  -X .mvn/ \
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
	  -X .txt \
	  -X .yaml \
	  -X pkg/data/assets/ \
	  -X zz_generated.

# ----------------------------------------------------------------------------------------------------------------------
# Run the Operator locally.
#
# To exit out of the local Operator you can use ctrl-c or ctrl-z but
# sometimes this leaves orphaned processes on the local machine so
# ensure these are killed run "make debug-stop"
# ----------------------------------------------------------------------------------------------------------------------
run: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
run: export UTILS_IMAGE := $(UTILS_IMAGE)
run: create-namespace ## run the Operator locally
	go run -ldflags "$(LDFLAGS)" ./main.go --skip-service-suspend=true --coherence-dev-mode=true \
		--cert-type=self-signed --webhook-service=host.docker.internal \
	    2>&1 | tee $(TEST_LOGS_DIR)/operator-debug.out

# ----------------------------------------------------------------------------------------------------------------------
# Run the Operator locally after deleting and recreating the test namespace.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: run-clean
run-clean: reset-namespace run ## run the Operator locally after resetting the namespace

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
run-debug: create-namespace ## run the Operator locally with Delve debugger
	dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient \
		-- --skip-service-suspend=true --coherence-dev-mode=true \
		--cert-type=self-signed --webhook-service=host.docker.internal

# ----------------------------------------------------------------------------------------------------------------------
# Run the Operator locally in debug mode after deleting and recreating
# the test namespace.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: run-debug-clean
run-debug-clean: reset-namespace run-debug ## run the Operator locally with Delve debugger after resetting the namespace

# ----------------------------------------------------------------------------------------------------------------------
# Kill any locally running Operator
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: stop
stop: ## kill any locally running operator process
	./hack/kill-local.sh

# ======================================================================================================================
# Targets related to Operator Lifecycle Manager and the Operator SDK
# ======================================================================================================================
##@ Operator Lifecycle Manager

# ----------------------------------------------------------------------------------------------------------------------
# Generate bundle manifests and metadata, then validate generated files.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: bundle
bundle: $(BUILD_PROPS) ensure-sdk kustomize $(BUILD_TARGETS)/manifests  ## Generate OLM bundle manifests and metadata, then validate generated files.
	$(OPERATOR_SDK) generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(OPERATOR_IMAGE)
	$(KUSTOMIZE) build config/manifests | $(OPERATOR_SDK) generate bundle -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)
	$(OPERATOR_SDK) bundle validate ./bundle
	$(OPERATOR_SDK) bundle validate ./bundle --select-optional suite=operatorframework --optional-values=k8s-version=1.22
	$(OPERATOR_SDK) bundle validate ./bundle --select-optional name=community --optional-values=image-path=bundle.Dockerfile

# ----------------------------------------------------------------------------------------------------------------------
# Build the bundle image.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: bundle-build
bundle-build:  ## Build the OLM image
	docker build --no-cache -f bundle.Dockerfile -t $(BUNDLE_IMG) .

.PHONY: bundle-push
bundle-push: ## Push the OLM bundle image.
	$(MAKE) docker-push IMG=$(BUNDLE_IMG)

.PHONY: opm
OPM = ./bin/opm
opm: ## Download opm locally if necessary.
ifeq (,$(wildcard $(OPM)))
ifeq (,$(shell which opm 2>/dev/null))
	@{ \
	set -e ;\
	mkdir -p $(dir $(OPM)) ;\
	OS=$(shell go env GOOS) && ARCH=$(shell go env GOARCH) && \
	curl -sSLo $(OPM) https://github.com/operator-framework/operator-registry/releases/download/v1.15.1/$${OS}-$${ARCH}-opm ;\
	chmod +x $(OPM) ;\
	}
else
OPM = $(shell which opm)
endif
endif

# A comma-separated list of bundle images (e.g. make catalog-build BUNDLE_IMGS=example.com/operator-bundle:v0.1.0,example.com/operator-bundle:v0.2.0).
# These images MUST exist in a registry and be pull-able.
BUNDLE_IMGS ?= $(BUNDLE_IMG)

# The image tag given to the resulting catalog image (e.g. make catalog-build CATALOG_IMG=example.com/operator-catalog:v0.2.0).
CATALOG_IMG ?= $(IMAGE_TAG_BASE)-catalog:v$(VERSION)

# Set CATALOG_BASE_IMG to an existing catalog image tag to add $BUNDLE_IMGS to that image.
ifneq ($(origin CATALOG_BASE_IMG), undefined)
FROM_INDEX_OPT := --from-index $(CATALOG_BASE_IMG)
endif

# Build a catalog image by adding bundle images to an empty catalog using the operator package manager tool, 'opm'.
# This recipe invokes 'opm' in 'semver' bundle add mode. For more information on add modes, see:
# https://github.com/operator-framework/community-operators/blob/7f1438c/docs/packaging-operator.md#updating-your-existing-operator
.PHONY: catalog-build
catalog-build: opm ## Build an OLM catalog image.
	$(OPM) index add --container-tool docker --mode semver --tag $(CATALOG_IMG) --bundles $(BUNDLE_IMGS) $(FROM_INDEX_OPT)

# Push the catalog image.
.PHONY: catalog-push
catalog-push: ## Push an OLM catalog image.
	$(MAKE) docker-push IMG=$(CATALOG_IMG)


# ======================================================================================================================
# Targets to run various tests
# ======================================================================================================================
##@ Test

# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go unit tests that do not require a k8s cluster
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: test-operator
test-operator: export CGO_ENABLED = 0
test-operator: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
test-operator: export UTILS_IMAGE := $(UTILS_IMAGE)
test-operator: $(BUILD_PROPS) $(BUILD_TARGETS)/manifests $(BUILD_TARGETS)/generate gotestsum  ## Run the Operator unit tests
	@echo "Running operator tests"
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-test.xml \
	  -- $(GO_TEST_FLAGS) -v ./api/... ./controllers/... ./pkg/...

# ----------------------------------------------------------------------------------------------------------------------
# Build and test the Java artifacts
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: test-mvn
test-mvn: $(BUILD_OUTPUT)/certs build-mvn  ## Run the Java artefact tests
	./mvnw -B -f java verify -Dtest.certs.location=$(BUILD_OUTPUT)/certs $(MAVEN_BUILD_OPTS)


# ----------------------------------------------------------------------------------------------------------------------
# Run all unit tests (both Go and Java)
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: test-all
test-all: test-mvn test-operator  ## Run all unit tests

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
e2e-local-test: export TEST_APPLICATION_IMAGE_CLIENT := $(TEST_APPLICATION_IMAGE_CLIENT)
e2e-local-test: export TEST_APPLICATION_IMAGE_HELIDON := $(TEST_APPLICATION_IMAGE_HELIDON)
e2e-local-test: export TEST_APPLICATION_IMAGE_SPRING := $(TEST_APPLICATION_IMAGE_SPRING)
e2e-local-test: export TEST_APPLICATION_IMAGE_SPRING_FAT := $(TEST_APPLICATION_IMAGE_SPRING_FAT)
e2e-local-test: export TEST_APPLICATION_IMAGE_SPRING_CNBP := $(TEST_APPLICATION_IMAGE_SPRING_CNBP)
e2e-local-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
e2e-local-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
e2e-local-test: export COH_SKIP_SITE := true
e2e-local-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
e2e-local-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
e2e-local-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
e2e-local-test: export TEST_ASSET_KUBECTL := $(TEST_ASSET_KUBECTL)
e2e-local-test: export LOCAL_STORAGE_RESTART := $(LOCAL_STORAGE_RESTART)
e2e-local-test: export VERSION := $(VERSION)
e2e-local-test: export MVN_VERSION := $(MVN_VERSION)
e2e-local-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
e2e-local-test: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
e2e-local-test: export UTILS_IMAGE := $(UTILS_IMAGE)
e2e-local-test: $(BUILD_TARGETS)/build-operator reset-namespace create-ssl-secrets install-crds gotestsum undeploy   ## Run the Operator end-to-end 'local' functional tests using a local Operator instance
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
e2e-test: prepare-e2e-test ## Run the Operator end-to-end 'remote' functional tests using an Operator deployed in k8s
	$(MAKE) run-e2e-test $${MF} \
	; rc=$$? \
	; $(MAKE) undeploy $${MF} \
	; $(MAKE) delete-namespace $${MF} \
	; exit $$rc

.PHONY: prepare-e2e-test
prepare-e2e-test: $(BUILD_TARGETS)/build-operator reset-namespace create-ssl-secrets install-crds deploy-and-wait

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
run-e2e-test: export TEST_APPLICATION_IMAGE_CLIENT := $(TEST_APPLICATION_IMAGE_CLIENT)
run-e2e-test: export TEST_APPLICATION_IMAGE_HELIDON := $(TEST_APPLICATION_IMAGE_HELIDON)
run-e2e-test: export TEST_APPLICATION_IMAGE_SPRING := $(TEST_APPLICATION_IMAGE_SPRING)
run-e2e-test: export TEST_APPLICATION_IMAGE_SPRING_FAT := $(TEST_APPLICATION_IMAGE_SPRING_FAT)
run-e2e-test: export TEST_APPLICATION_IMAGE_SPRING_CNBP := $(TEST_APPLICATION_IMAGE_SPRING_CNBP)
run-e2e-test: gotestsum  ## Run the Operator 'remote' end-to-end functional tests using an ALREADY DEPLOYED Operator
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/e2e/remote/...


# ----------------------------------------------------------------------------------------------------------------------
# Run the end-to-end Coherence client tests.
# ----------------------------------------------------------------------------------------------------------------------
e2e-client-test: export CGO_ENABLED = 0
e2e-client-test: export CLIENT_CLASSPATH := $(CURRDIR)/java/operator-test-client/target/operator-test-client-$(MVN_VERSION).jar:$(CURRDIR)/java/operator-test-client/target/lib/*
e2e-client-test: export OPERATOR_NAMESPACE := $(OPERATOR_NAMESPACE)
e2e-client-test: export OPERATOR_NAMESPACE_CLIENT := $(OPERATOR_NAMESPACE_CLIENT)
e2e-client-test: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
e2e-client-test: export TEST_APPLICATION_IMAGE_CLIENT := $(TEST_APPLICATION_IMAGE_CLIENT)
e2e-client-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
e2e-client-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
e2e-client-test: export COH_SKIP_SITE := true
e2e-client-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
e2e-client-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
e2e-client-test: export VERSION := $(VERSION)
e2e-client-test: export MVN_VERSION := $(MVN_VERSION)
e2e-client-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
e2e-client-test: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
e2e-client-test: export UTILS_IMAGE := $(UTILS_IMAGE)
e2e-client-test: build-operator-images build-client-image reset-namespace create-ssl-secrets install-crds gotestsum undeploy   ## Run the end-to-end Coherence client tests using a local Operator deployment
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-client-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/e2e/clients/...


# ----------------------------------------------------------------------------------------------------------------------
# Run the end-to-end Helm chart tests.
# ----------------------------------------------------------------------------------------------------------------------
e2e-helm-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
e2e-helm-test: export UTILS_IMAGE := $(UTILS_IMAGE)
e2e-helm-test: $(BUILD_PROPS) $(BUILD_HELM)/coherence-operator-$(VERSION).tgz reset-namespace install-crds gotestsum  ## Run the Operator Helm chart end-to-end functional tests
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
e2e-prometheus-test: reset-namespace install-prometheus $(BUILD_TARGETS)/build-operator create-ssl-secrets install-crds deploy-and-wait   ## Run the Operator metrics/Prometheus end-to-end functional tests
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
run-prometheus-test: export TEST_APPLICATION_IMAGE_HELIDON := $(TEST_APPLICATION_IMAGE_HELIDON)
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
e2e-elastic-test: reset-namespace install-elastic $(BUILD_TARGETS)/build-operator create-ssl-secrets install-crds deploy-and-wait   ## Run the Operator logging/ElasticSearch end-to-end functional tests
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
run-elastic-test: export TEST_APPLICATION_IMAGE_HELIDON := $(TEST_APPLICATION_IMAGE_HELIDON)
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
# Executes the Go end-to-end Operator backwards compatibility tests to ensure upgrades from previous Operator versions
# work and do not impact running clusters, etc.
# These tests will use whichever k8s cluster the local environment is pointing to.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: compatibility-test
compatibility-test: export CGO_ENABLED = 0
compatibility-test: export OPERATOR_NAMESPACE := $(OPERATOR_NAMESPACE)
compatibility-test: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
compatibility-test: export TEST_APPLICATION_IMAGE_CLIENT := $(TEST_APPLICATION_IMAGE_CLIENT)
compatibility-test: export TEST_APPLICATION_IMAGE_HELIDON := $(TEST_APPLICATION_IMAGE_HELIDON)
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
compatibility-test: undeploy build-all-images $(BUILD_HELM)/coherence-operator-$(VERSION).tgz undeploy clean-namespace reset-namespace gotestsum    ## Run the Operator backwards compatibility tests
	helm repo add coherence https://oracle.github.io/coherence-operator/charts
	helm repo update
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-compatibility-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/e2e/compatibility/...


# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end Operator Kubernetes versions certification tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# Note that the namespace will be created if it does not exist.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: certification-test
certification-test: export MF = $(MAKEFLAGS)
certification-test: install-certification     ## Run the Operator Kubernetes versions certification tests
	$(MAKE) run-certification  $${MF} \
	; rc=$$? \
	; $(MAKE) cleanup-certification $${MF} \
	; exit $$rc


# ----------------------------------------------------------------------------------------------------------------------
# Install the Operator prior to running compatibility tests.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-certification
install-certification: $(BUILD_TARGETS)/build-operator reset-namespace create-ssl-secrets deploy-and-wait

# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end Operator certification tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# Note that the namespace will be created if it does not exist.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: run-certification
run-certification: export CGO_ENABLED = 0
run-certification: export OPERATOR_NAMESPACE := $(OPERATOR_NAMESPACE)
run-certification: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
run-certification: export TEST_APPLICATION_IMAGE_CLIENT := $(TEST_APPLICATION_IMAGE_CLIENT)
run-certification: export TEST_APPLICATION_IMAGE_HELIDON := $(TEST_APPLICATION_IMAGE_HELIDON)
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
# Clean up after to running compatibility tests.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: cleanup-certification
cleanup-certification: undeploy uninstall-crds clean-namespace

# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end Operator Coherence versions compatibility tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# Note that the namespace will be created if it does not exist.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: coherence-compatibility-test
coherence-compatibility-test: export MF = $(MAKEFLAGS)
coherence-compatibility-test: install-coherence-compatibility   ## Run the Operator Coherence versions compatibility tests
	$(MAKE) run-coherence-compatibility  $${MF} \
	; rc=$$? \
	; $(MAKE) cleanup-coherence-compatibility $${MF} \
	; exit $$rc


# ----------------------------------------------------------------------------------------------------------------------
# Install the Operator prior to running coherence compatibility tests.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-coherence-compatibility
install-coherence-compatibility: $(BUILD_TARGETS)/build-operator reset-namespace create-ssl-secrets deploy-and-wait

# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end Operator coherence compatibility tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# Note that the namespace will be created if it does not exist.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: run-coherence-compatibility
run-coherence-compatibility: export CGO_ENABLED = 0
run-coherence-compatibility: export OPERATOR_NAMESPACE := $(OPERATOR_NAMESPACE)
run-coherence-compatibility: export TEST_COMPATIBILITY_IMAGE := $(TEST_COMPATIBILITY_IMAGE)
run-coherence-compatibility: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
run-coherence-compatibility: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
run-coherence-compatibility: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
run-coherence-compatibility: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
run-coherence-compatibility: export VERSION := $(VERSION)
run-coherence-compatibility: export CERTIFICATION_VERSION := $(CERTIFICATION_VERSION)
run-coherence-compatibility: export OPERATOR_IMAGE_REPO := $(OPERATOR_IMAGE_REPO)
run-coherence-compatibility: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
run-coherence-compatibility: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
run-coherence-compatibility: export UTILS_IMAGE := $(UTILS_IMAGE)
run-coherence-compatibility: gotestsum $(BUILD_TARGETS)/generate
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-coherence-compatibility-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/coherence_compatibility/...

# ----------------------------------------------------------------------------------------------------------------------
# Clean up after to running backwards compatibility tests.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: cleanup-coherence-compatibility
cleanup-coherence-compatibility: undeploy uninstall-crds clean-namespace

# ======================================================================================================================
# Targets related to deploying the Operator into k8s for testing and debugging
# ======================================================================================================================
##@ Deployment

# ----------------------------------------------------------------------------------------------------------------------
# Install CRDs into Kubernetes.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-crds
install-crds: prepare-deploy uninstall-crds  ## Install the CRDs
	$(KUSTOMIZE) build $(BUILD_DEPLOY)/crd | kubectl create -f -

# ----------------------------------------------------------------------------------------------------------------------
# Uninstall CRDs from Kubernetes.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: uninstall-crds
uninstall-crds: $(BUILD_TARGETS)/manifests  ## Uninstall the CRDs
	$(call prepare_deploy,$(OPERATOR_IMAGE),$(OPERATOR_NAMESPACE))
	$(KUSTOMIZE) build $(BUILD_DEPLOY)/crd | kubectl delete -f - || true

# ----------------------------------------------------------------------------------------------------------------------
# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: deploy-and-wait
deploy-and-wait: deploy wait-for-deploy   ## Deploy the Coherence Operator and wait for the Operator Pod to be ready

.PHONY: deploy
deploy: prepare-deploy create-namespace kustomize   ## Deploy the Coherence Operator
ifneq (,$(WATCH_NAMESPACE))
	cd $(BUILD_DEPLOY)/manager && $(KUSTOMIZE) edit add configmap env-vars --from-literal WATCH_NAMESPACE=$(WATCH_NAMESPACE)
endif
	kubectl -n $(OPERATOR_NAMESPACE) create secret generic coherence-webhook-server-cert || true
	$(KUSTOMIZE) build $(BUILD_DEPLOY)/default | kubectl apply -f -
	sleep 5

.PHONY: just-deploy
just-deploy:
	$(call prepare_deploy,$(OPERATOR_IMAGE),$(OPERATOR_NAMESPACE))
	$(KUSTOMIZE) build $(BUILD_DEPLOY)/default | kubectl apply -f -

.PHONY: prepare-deploy
prepare-deploy: $(BUILD_TARGETS)/manifests $(BUILD_TARGETS)/build-operator kustomize
	$(call prepare_deploy,$(OPERATOR_IMAGE),$(OPERATOR_NAMESPACE))

.PHONY: deploy-debug
deploy-debug: prepare-deploy-debug create-namespace kustomize   ## Deploy the Coherence Operator running with Delve
ifneq (,$(WATCH_NAMESPACE))
	cd $(BUILD_DEPLOY)/manager && $(KUSTOMIZE) edit add configmap env-vars --from-literal WATCH_NAMESPACE=$(WATCH_NAMESPACE)
endif
	kubectl -n $(OPERATOR_NAMESPACE) create secret generic coherence-webhook-server-cert || true
	$(KUSTOMIZE) build $(BUILD_DEPLOY)/default | kubectl apply -f -
	sleep 5
	@echo ""
	@echo "Deployed a debug enabled Operator."
	@echo ""
	@echo "To allow the IDE to connect a remote debugger run the following command to forward port 2345"
	@echo ""
	@echo "make port-forward-debug"
	@echo ""
	@echo "To see tail the Operator logs during debugging you can run:"
	@echo ""
	@echo "make tail-logs"
	@echo ""


.PHONY: port-forward-debug
port-forward-debug: export POD=$(shell kubectl -n $(OPERATOR_NAMESPACE) get pod -l control-plane=coherence -o name)
port-forward-debug:  ## Run a port-forward process to forward localhost:2345 to port 2345 in the Operator Pod
	@echo "Starting port-forward to the Operator Pod on port 2345 - DO NOT stop this process until debugging is finished!"
	@echo "Connect your IDE debugger to localhost:2345 (which is the default remote debug setting in IDEs like Goland)"
	@echo "If your IDE immediately disconnects it may be that the Operator Pod was not yet started, so try again."
	@echo ""
	kubectl -n $(OPERATOR_NAMESPACE) port-forward $(POD) 2345:2345 || true

.PHONY: prepare-deploy-debug
prepare-deploy-debug: $(BUILD_TARGETS)/manifests build-operator-debug kustomize
	$(call prepare_deploy,$(OPERATOR_IMAGE_DEBUG),$(OPERATOR_NAMESPACE))

.PHONY: wait-for-deploy
wait-for-deploy: export POD=$(shell kubectl -n $(OPERATOR_NAMESPACE) get pod -l control-plane=coherence -o name)
wait-for-deploy:
	echo "Waiting for Operator to be ready"
	sleep 10
	kubectl -n $(OPERATOR_NAMESPACE) wait --for condition=ready --timeout 120s $(POD)

# ----------------------------------------------------------------------------------------------------------------------
# Prepare the deployment manifests - this is called by a number of other targets.
# Parameter #1 is the Operator Image Name
# Parameter #2 is the name of the namespace to deploy into
# ----------------------------------------------------------------------------------------------------------------------
define prepare_deploy
	-rm -r $(BUILD_DEPLOY)
	mkdir -p $(BUILD_DEPLOY)
	cp -R config $(BUILD_OUTPUT)
	cd $(BUILD_DEPLOY)/manager && $(KUSTOMIZE) edit add configmap env-vars --from-literal COHERENCE_IMAGE=$(COHERENCE_IMAGE)
	cd $(BUILD_DEPLOY)/manager && $(KUSTOMIZE) edit add configmap env-vars --from-literal UTILS_IMAGE=$(UTILS_IMAGE)
	cd $(BUILD_DEPLOY)/manager && $(KUSTOMIZE) edit set image controller=$(1)
	cd $(BUILD_DEPLOY)/default && $(KUSTOMIZE) edit set namespace $(2)
endef

# ----------------------------------------------------------------------------------------------------------------------
# Un-deploy controller from the configured Kubernetes cluster in ~/.kube/config
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: undeploy
undeploy: $(BUILD_PROPS) $(BUILD_TARGETS)/manifests kustomize  ## Undeploy the Coherence Operator
	$(call prepare_deploy,$(OPERATOR_IMAGE),$(OPERATOR_NAMESPACE))
	$(KUSTOMIZE) build $(BUILD_DEPLOY)/default | kubectl delete -f - || true
	kubectl -n $(OPERATOR_NAMESPACE) delete secret coherence-webhook-server-cert || true
	kubectl delete mutatingwebhookconfiguration coherence-operator-mutating-webhook-configuration || true
	kubectl delete validatingwebhookconfiguration coherence-operator-validating-webhook-configuration || true


# ----------------------------------------------------------------------------------------------------------------------
# Tail the deployed operator logs.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: tail-logs
tail-logs: export POD=$(shell kubectl -n $(OPERATOR_NAMESPACE) get pod -l control-plane=coherence -o name)
tail-logs:     ## Tail the Coherence Operator Pod logs (with follow)
	kubectl -n $(OPERATOR_NAMESPACE) logs $(POD) -c manager -f


$(BUILD_MANIFESTS_PKG): kustomize
	rm -rf $(BUILD_MANIFESTS) || true
	mkdir -p $(BUILD_MANIFESTS)/crd
	$(KUSTOMIZE) build config/crd > $(BUILD_MANIFESTS)/crd/coherence.oracle.com_coherence.yaml
	cp -R config/default/ $(BUILD_MANIFESTS)/default
	cp -R config/manager/ $(BUILD_MANIFESTS)/manager
	cp -R config/rbac/ $(BUILD_MANIFESTS)/rbac
	tar -C $(BUILD_OUTPUT) -czf $(BUILD_MANIFESTS_PKG) manifests/
	$(call prepare_deploy,$(OPERATOR_IMAGE),"coherence")
	cp config/namespace/namespace.yaml $(BUILD_OUTPUT)/coherence-operator.yaml
	$(KUSTOMIZE) build $(BUILD_DEPLOY)/default >> $(BUILD_OUTPUT)/coherence-operator.yaml
ifeq (Darwin, $(UNAME_S))
	sed -i '' -e 's/name: coherence-operator-env-vars-.*/name: coherence-operator-env-vars/g' $(BUILD_OUTPUT)/coherence-operator.yaml
else
	sed -i 's/name: coherence-operator-env-vars-.*/name: coherence-operator-env-vars/g' $(BUILD_OUTPUT)/coherence-operator.yaml
endif


# ----------------------------------------------------------------------------------------------------------------------
# Delete and re-create the test namespace
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: create-namespace
create-namespace: export KUBECONFIG_PATH := $(KUBECONFIG_PATH)
create-namespace: ## Create the test namespace
ifeq ($(CREATE_OPERATOR_NAMESPACE),true)
	kubectl get ns $(OPERATOR_NAMESPACE) -o name > /dev/null 2>&1 || kubectl create namespace $(OPERATOR_NAMESPACE)
	kubectl get ns $(OPERATOR_NAMESPACE_CLIENT) -o name > /dev/null 2>&1 || kubectl create namespace $(OPERATOR_NAMESPACE_CLIENT)
endif

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
reset-namespace: delete-namespace create-namespace      ## Reset the test namespace
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
delete-namespace: clean-namespace  ## Delete the test namespace
ifeq ($(CREATE_OPERATOR_NAMESPACE),true)
	@echo "Deleting test namespace $(OPERATOR_NAMESPACE)"
	kubectl delete namespace $(OPERATOR_NAMESPACE) --force --grace-period=0 && echo "deleted namespace $(OPERATOR_NAMESPACE)" || true
	kubectl delete namespace $(OPERATOR_NAMESPACE_CLIENT) --force --grace-period=0 && echo "deleted namespace $(OPERATOR_NAMESPACE_CLIENT)" || true
endif
	kubectl delete clusterrole operator-test-coherence-operator --force --grace-period=0 && echo "deleted namespace" || true
	kubectl delete clusterrolebinding operator-test-coherence-operator --force --grace-period=0 && echo "deleted namespace" || true

# ----------------------------------------------------------------------------------------------------------------------
# Delete all of the Coherence resources from the test namespace.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: delete-coherence-clusters
delete-coherence-clusters: ## Delete all running Coherence clusters in the test namespace
	for i in $$(kubectl -n  $(OPERATOR_NAMESPACE) get coherence.coherence.oracle.com -o name); do \
  		kubectl -n $(OPERATOR_NAMESPACE) patch $${i} -p '{"metadata":{"finalizers":[]}}' --type=merge || true ;\
		kubectl -n $(OPERATOR_NAMESPACE) delete $${i}; \
	done

# ----------------------------------------------------------------------------------------------------------------------
# Delete all resource from the test namespace
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: clean-namespace
clean-namespace: delete-coherence-clusters   ## Clean-up deployments in the test namespace
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


# ======================================================================================================================
# Targets related to running KinD clusters
# ======================================================================================================================
##@ KinD

KIND_IMAGE ?= "kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6"

# ----------------------------------------------------------------------------------------------------------------------
# Start a Kind cluster
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: kind
kind:   ## Run a default KinD cluster
	./hack/kind.sh --image $(KIND_IMAGE)
	./hack/kind-label-node.sh
	docker pull $(COHERENCE_IMAGE)
	kind load docker-image --name operator $(COHERENCE_IMAGE)

# ----------------------------------------------------------------------------------------------------------------------
# Stop and delete the Kind cluster
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: kind-stop
kind-stop:   ## Stop and delete the KinD cluster named "operator"
	kind delete cluster --name operator

# ----------------------------------------------------------------------------------------------------------------------
# Load images into Kind
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: kind-load
kind-load: kind-load-operator  ## Load all images into the KinD cluster
	kind load docker-image --name operator $(TEST_APPLICATION_IMAGE) || true
	kind load docker-image --name operator $(TEST_APPLICATION_IMAGE_CLIENT) || true
	kind load docker-image --name operator $(TEST_APPLICATION_IMAGE_HELIDON) || true
	kind load docker-image --name operator $(TEST_APPLICATION_IMAGE_SPRING) || true
	kind load docker-image --name operator $(TEST_APPLICATION_IMAGE_SPRING_FAT) || true
	kind load docker-image --name operator $(TEST_APPLICATION_IMAGE_SPRING_CNBP) || true
	kind load docker-image --name operator gcr.io/kubebuilder/kube-rbac-proxy:v0.5.0 || true
	kind load docker-image --name operator docker.elastic.co/elasticsearch/elasticsearch:7.6.2 || true
	kind load docker-image --name operator docker.elastic.co/kibana/kibana:7.6.2 || true

.PHONY: kind-load-operator
kind-load-operator:   ## Load the Operator images into the KinD cluster
	kind load docker-image --name operator $(OPERATOR_IMAGE) || true
	kind load docker-image --name operator $(UTILS_IMAGE) || true

# ----------------------------------------------------------------------------------------------------------------------
# Load compatibility images into Kind
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: kind-load-compatibility
kind-load-compatibility:   ## Load the compatibility test images into the KinD cluster
	kind load docker-image --name operator $(TEST_COMPATIBILITY_IMAGE) || true

# ======================================================================================================================
# Miscellaneous targets
# ======================================================================================================================
##@ Miscellaneous

# ----------------------------------------------------------------------------------------------------------------------
# find or download controller-gen
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: controller-gen
CONTROLLER_GEN = $(TOOLS_BIN)/controller-gen
controller-gen: ## Download controller-gen locally if necessary.
	$(call go-get-tool,$(CONTROLLER_GEN),sigs.k8s.io/controller-tools/cmd/controller-gen@v0.7.0)

# ----------------------------------------------------------------------------------------------------------------------
# find or download kustomize
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: kustomize
KUSTOMIZE = $(TOOLS_BIN)/kustomize
kustomize: ## Download kustomize locally if necessary.
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

# ----------------------------------------------------------------------------------------------------------------------
# find or download gotestsum
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: gotestsum
GOTESTSUM = $(TOOLS_BIN)/gotestsum
gotestsum: ## Download gotestsum locally if necessary.
	$(call go-get-tool,$(GOTESTSUM),gotest.tools/gotestsum@v0.5.2)


# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2) into $(TOOLS_BIN)" ;\
GOBIN=$(TOOLS_BIN) go get $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef


# ----------------------------------------------------------------------------------------------------------------------
# Deploy the Java artifacts
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: mvn-deploy
mvn-deploy: java-client
	./mvnw $(MAVEN_BUILD_OPTS) -s ./.mvn/settings.xml -B -f java clean deploy -DskipTests -DskipTests -Prelease -Dgpg.passphrase=$(GPG_PASSPHRASE)

# ----------------------------------------------------------------------------------------------------------------------
# Build the examples
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-examples
build-examples:
	./mvnw -B -f ./examples package jib:dockerBuild -DskipTests $(MAVEN_BUILD_OPTS)

# ----------------------------------------------------------------------------------------------------------------------
# Build and test the examples
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: test-examples
test-examples: build-examples
	./mvnw -B -f ./examples verify $(MAVEN_BUILD_OPTS)

# ----------------------------------------------------------------------------------------------------------------------
# Push the Operator Docker image
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-operator-image
push-operator-image: $(BUILD_TARGETS)/build-operator
ifeq ($(OPERATOR_RELEASE_IMAGE), $(OPERATOR_IMAGE))
	@echo "Pushing $(OPERATOR_IMAGE)"
	docker push $(OPERATOR_IMAGE)-amd64
	docker push $(OPERATOR_IMAGE)-arm64
	docker manifest create $(OPERATOR_IMAGE) \
		--amend $(OPERATOR_IMAGE)-amd64 \
		--amend $(OPERATOR_IMAGE)-arm64
	docker manifest annotate $(OPERATOR_IMAGE) $(OPERATOR_IMAGE)-arm64 --arch arm64
	docker manifest push $(OPERATOR_IMAGE) 
else
	@echo "Tagging $(OPERATOR_IMAGE)-amd64 as $(OPERATOR_RELEASE_IMAGE)-amd64"
	docker tag $(OPERATOR_IMAGE)-amd64 $(OPERATOR_RELEASE_IMAGE)-amd64
	@echo "Pushing $(OPERATOR_RELEASE_IMAGE)-amd64"
	docker push $(OPERATOR_RELEASE_IMAGE)-amd64
	@echo "Tagging $(OPERATOR_IMAGE)-arm64 as $(OPERATOR_RELEASE_IMAGE)-arm64"
	docker tag $(OPERATOR_IMAGE)-arm64 $(OPERATOR_RELEASE_IMAGE)-arm64
	@echo "Pushing $(OPERATOR_RELEASE_IMAGE)-arm64"
	docker push $(OPERATOR_RELEASE_IMAGE)-arm64
	docker manifest create $(OPERATOR_RELEASE_IMAGE) \
		--amend $(OPERATOR_RELEASE_IMAGE)-amd64 \
		--amend $(OPERATOR_RELEASE_IMAGE)-arm64
	docker manifest annotate $(OPERATOR_RELEASE_IMAGE) $(OPERATOR_RELEASE_IMAGE)-arm64 --arch arm64
	docker manifest push $(OPERATOR_RELEASE_IMAGE)
endif

# ----------------------------------------------------------------------------------------------------------------------
# Push the Operator Utils Docker image
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-utils-image
push-utils-image:
ifeq ($(UTILS_RELEASE_IMAGE), $(UTILS_IMAGE))
	@echo "Pushing $(UTILS_IMAGE)-amd64"
	docker push $(UTILS_IMAGE)-amd64
	@echo "Pushing $(UTILS_IMAGE)-arm64"
	docker push $(UTILS_IMAGE)-arm64
	@echo "Creating $(UTILS_IMAGE) manifest"
	docker manifest create $(UTILS_IMAGE) \
		--amend $(UTILS_IMAGE)-amd64 \
		--amend $(UTILS_IMAGE)-arm64
	docker manifest annotate $(UTILS_IMAGE) $(UTILS_IMAGE)-arm64 --arch arm64
	@echo "Pushing $(UTILS_IMAGE) manifest"
	docker manifest push $(UTILS_IMAGE)
else
	@echo "Tagging $(UTILS_IMAGE)-amd64 as $(UTILS_RELEASE_IMAGE)-amd64"
	docker tag $(UTILS_IMAGE)-amd64 $(UTILS_RELEASE_IMAGE)-amd64
	@echo "Pushing $(UTILS_RELEASE_IMAGE)-amd64"
	docker push $(UTILS_RELEASE_IMAGE)-amd64
	@echo "Tagging $(UTILS_IMAGE)-arm64 as $(UTILS_RELEASE_IMAGE)-arm64"
	docker tag $(UTILS_IMAGE)-arm64 $(UTILS_RELEASE_IMAGE)-arm64
	@echo "Pushing $(UTILS_RELEASE_IMAGE)-arm64"
	docker push $(UTILS_RELEASE_IMAGE)-arm64
	@echo "Creating $(UTILS_RELEASE_IMAGE) manifest"
	docker manifest create $(UTILS_RELEASE_IMAGE) \
		--amend $(UTILS_RELEASE_IMAGE)-amd64 \
		--amend $(UTILS_RELEASE_IMAGE)-arm64
	docker manifest annotate $(UTILS_RELEASE_IMAGE) $(UTILS_RELEASE_IMAGE)-arm64 --arch arm64
	@echo "Pushing $(UTILS_RELEASE_IMAGE) manifest"
	docker manifest push $(UTILS_RELEASE_IMAGE)
endif

# ----------------------------------------------------------------------------------------------------------------------
# Push the test base images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-test-base-images
push-test-base-images:
ifeq ($(TEST_BASE_RELEASE_IMAGE), $(TEST_BASE_IMAGE))
	@echo "Pushing $(TEST_BASE_IMAGE)-amd64"
	docker push $(TEST_BASE_IMAGE)-amd64
	@echo "Pushing $(TEST_BASE_IMAGE)-arm64"
	docker push $(TEST_BASE_IMAGE)-arm64
	@echo "Creating $(TEST_BASE_IMAGE) manifest"
	docker manifest create $(TEST_BASE_IMAGE) \
		--amend $(TEST_BASE_IMAGE)-amd64 \
		--amend $(TEST_BASE_IMAGE)-arm64
	docker manifest annotate $(TEST_BASE_IMAGE) $(TEST_BASE_IMAGE)-arm64 --arch arm64
	@echo "Pushing $(TEST_BASE_IMAGE) manifest"
	docker manifest push $(TEST_BASE_IMAGE)
else
	@echo "Tagging $(TEST_BASE_IMAGE)-amd64 as $(TEST_BASE_RELEASE_IMAGE)-amd64"
	docker tag $(TEST_BASE_IMAGE)-amd64 $(TEST_BASE_RELEASE_IMAGE)-amd64
	@echo "Pushing $(TEST_BASE_RELEASE_IMAGE)-amd64"
	docker push $(TEST_BASE_RELEASE_IMAGE)-amd64
	@echo "Tagging $(TEST_BASE_IMAGE)-arm64 as $(TEST_BASE_RELEASE_IMAGE)-arm64"
	docker tag $(TEST_BASE_IMAGE)-arm64 $(TEST_BASE_RELEASE_IMAGE)-arm64
	@echo "Pushing $(TEST_BASE_RELEASE_IMAGE)-arm64"
	docker push $(TEST_BASE_RELEASE_IMAGE)-arm64
	@echo "Creating $(TEST_BASE_RELEASE_IMAGE) manifest"
	docker manifest create $(TEST_BASE_RELEASE_IMAGE) \
		--amend $(TEST_BASE_RELEASE_IMAGE)-amd64 \
		--amend $(TEST_BASE_RELEASE_IMAGE)-arm64
	docker manifest annotate $(TEST_BASE_RELEASE_IMAGE) $(TEST_BASE_RELEASE_IMAGE)-arm64 --arch arm64
	@echo "Pushing $(TEST_BASE_RELEASE_IMAGE) manifest"
	docker manifest push $(TEST_BASE_RELEASE_IMAGE)
endif

# ----------------------------------------------------------------------------------------------------------------------
# Push the Operator JIB Test Docker images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-test-images
push-test-images:
	docker push $(TEST_APPLICATION_IMAGE)
	docker push $(TEST_APPLICATION_IMAGE_CLIENT)
	docker push $(TEST_APPLICATION_IMAGE_HELIDON)
	docker push $(TEST_APPLICATION_IMAGE_SPRING)
	docker push $(TEST_APPLICATION_IMAGE_SPRING_FAT)
	docker push $(TEST_APPLICATION_IMAGE_SPRING_CNBP)

# ----------------------------------------------------------------------------------------------------------------------
# Build the Operator Test images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-compatibility-image
build-compatibility-image: build-mvn
	./mvnw -B -f java/operator-compatibility package -DskipTests \
	    -Dcoherence.compatibility.image.name=$(TEST_COMPATIBILITY_IMAGE) \
	    -Dcoherence.compatibility.coherence.image=$(COHERENCE_IMAGE) $(MAVEN_BUILD_OPTS)
	./mvnw -B -f java/operator-compatibility exec:exec \
	    -Dcoherence.compatibility.image.name=$(TEST_COMPATIBILITY_IMAGE) \
	    -Dcoherence.compatibility.coherence.image=$(COHERENCE_IMAGE) $(MAVEN_BUILD_OPTS)

# ----------------------------------------------------------------------------------------------------------------------
# Push the Operator JIB Test Docker images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-compatibility-image
push-compatibility-image:
	docker push $(TEST_COMPATIBILITY_IMAGE)

# ----------------------------------------------------------------------------------------------------------------------
# Push all of the Docker images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-all-images
push-all-images: push-test-images push-test-base-images push-utils-image push-operator-image

# ----------------------------------------------------------------------------------------------------------------------
# Push all of the Docker images that are released
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-release-images
push-release-images: push-test-base-images push-utils-image push-operator-image

# ----------------------------------------------------------------------------------------------------------------------
# Install Prometheus
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: get-prometheus
get-prometheus: $(PROMETHEUS_HOME)/$(PROMETHEUS_VERSION).txt ## Download Prometheus Operator kube-prometheus

$(PROMETHEUS_HOME)/$(PROMETHEUS_VERSION).txt:
	curl -sL https://github.com/prometheus-operator/kube-prometheus/archive/refs/tags/$(PROMETHEUS_VERSION).tar.gz -o $(BUILD_OUTPUT)/prometheus.tar.gz
	mkdir $(PROMETHEUS_HOME)
	tar -zxf $(BUILD_OUTPUT)/prometheus.tar.gz --directory $(PROMETHEUS_HOME) --strip-components=1
	rm $(BUILD_OUTPUT)/prometheus.tar.gz
	touch $(PROMETHEUS_HOME)/$(PROMETHEUS_VERSION).txt

.PHONY: install-prometheus
install-prometheus: get-prometheus ## Install Prometheus and Grafana
	kubectl create -f $(PROMETHEUS_HOME)/manifests/setup
	until kubectl get servicemonitors --all-namespaces ; do date; sleep 1; echo ""; done
#   We create additional custom RBAC rules because the defaults do not work
#   in an RBAC enabled cluster such as KinD
#   See: https://prometheus-operator.dev/docs/operator/rbac/
	kubectl create -f hack/prometheus-rbac.yaml
	kubectl create -f $(PROMETHEUS_HOME)/manifests
	sleep 10
	@echo "Waiting for Prometheus StatefulSet to be ready"
	kubectl -n monitoring rollout status statefulset/prometheus-k8s --timeout=5m
	@echo "Waiting for Grafana Deployment to be ready"
	kubectl -n monitoring rollout status deployment/grafana --timeout=5m

# ----------------------------------------------------------------------------------------------------------------------
# Uninstall Prometheus
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: uninstall-prometheus
uninstall-prometheus: get-prometheus ## Uninstall Prometheus and Grafana
	kubectl delete --ignore-not-found=true -f $(PROMETHEUS_HOME)/manifests || true
	kubectl delete --ignore-not-found=true -f $(PROMETHEUS_HOME)/manifests/setup || true
	kubectl delete --ignore-not-found=true -f hack/prometheus-rbac.yaml

# ----------------------------------------------------------------------------------------------------------------------
# Install Prometheus Adapter used for k8s metrics and Horizontal Pod Autoscaler
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-prometheus-adapter
install-prometheus-adapter:
	kubectl create ns $(OPERATOR_NAMESPACE) || true
	helm repo add stable https://kubernetes-charts.storage.googleapis.com/ || true
	helm install --atomic --namespace $(OPERATOR_NAMESPACE) --version $(PROMETHEUS_ADAPTER_VERSION) --wait \
		--set prometheus.url=http://prometheus.$(OPERATOR_NAMESPACE).svc \
		--values hack/prometheus-adapter-values.yaml prometheus-adapter stable/prometheus-adapter

# ----------------------------------------------------------------------------------------------------------------------
# Uninstall Prometheus Adapter used for k8s metrics and Horizontal Pod Autoscaler
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: uninstall-prometheus-adapter
uninstall-prometheus-adapter:
	helm --namespace $(OPERATOR_NAMESPACE) delete prometheus-adapter || true

# ----------------------------------------------------------------------------------------------------------------------
# Start a port-forward process to the Grafana Pod.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: port-forward-grafana
port-forward-grafana: ## Run a port-forward to Grafana on http://127.0.0.1:3000
	@echo "Reach Grafana on http://127.0.0.1:3000"
	@echo "User: admin Password: admin"
	kubectl --namespace monitoring port-forward svc/grafana 3000

# ----------------------------------------------------------------------------------------------------------------------
# Install Elasticsearch & Kibana
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-elastic
install-elastic: helm-install-elastic kibana-import ## Install Elastic Stack

.PHONY: helm-install-elastic
helm-install-elastic:
	kubectl create ns $(OPERATOR_NAMESPACE) || true
#   Create the ConfigMap containing the Coherence Kibana dashboards
	kubectl -n $(OPERATOR_NAMESPACE) delete secret coherence-kibana-dashboard || true
	kubectl -n $(OPERATOR_NAMESPACE) create secret generic --from-file dashboards/kibana/kibana-dashboard-data.json coherence-kibana-dashboard
#   Create the ConfigMap containing the Coherence Kibana dashboards import script
	kubectl -n $(OPERATOR_NAMESPACE) delete secret coherence-kibana-import || true
	kubectl -n $(OPERATOR_NAMESPACE) create secret generic --from-file hack/kibana-import.sh coherence-kibana-import
#   Set-up the Elastic Helm repo
	@echo "Getting Helm Version:"
	helm version
	helm repo add elastic https://helm.elastic.co || true
	helm repo update || true
#   Install Elasticsearch
	helm install --atomic --namespace $(OPERATOR_NAMESPACE) --version $(ELASTIC_VERSION) --wait --timeout=10m \
		--debug --values hack/elastic-values.yaml elasticsearch elastic/elasticsearch
#   Install Kibana
	helm install --atomic --namespace $(OPERATOR_NAMESPACE) --version $(ELASTIC_VERSION) --wait --timeout=10m \
		--debug --values hack/kibana-values.yaml kibana elastic/kibana \

.PHONY: kibana-import
kibana-import:
	KIBANA_POD=$$(kubectl -n $(OPERATOR_NAMESPACE) get pod -l app=kibana -o name) \
	; kubectl -n $(OPERATOR_NAMESPACE) exec -it $${KIBANA_POD} -- /bin/bash /usr/share/kibana/data/coherence/scripts/kibana-import.sh

# ----------------------------------------------------------------------------------------------------------------------
# Uninstall Elasticsearch & Kibana
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: uninstall-elastic
uninstall-elastic: ## Uninstall Elastic Stack
	helm uninstall --namespace $(OPERATOR_NAMESPACE) kibana || true
	helm uninstall --namespace $(OPERATOR_NAMESPACE) elasticsearch || true
	kubectl -n $(OPERATOR_NAMESPACE) delete pvc elasticsearch-master-elasticsearch-master-0 || true

# ----------------------------------------------------------------------------------------------------------------------
# Start a port-forward process to the Kibana Pod.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: port-forward-kibana
port-forward-kibana: export KIBANA_POD := $(shell kubectl -n $(OPERATOR_NAMESPACE) get pod -l app=kibana -o name)
port-forward-kibana: ## Run a port-forward to Kibana on http://127.0.0.1:5601
	@echo "Reach Kibana on http://127.0.0.1:5601"
	kubectl -n $(OPERATOR_NAMESPACE) port-forward $(KIBANA_POD) 5601:5601

# ----------------------------------------------------------------------------------------------------------------------
# Start a port-forward process to the Elasticsearch Pod.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: port-forward-es
port-forward-es: export ES_POD := $(shell kubectl -n $(OPERATOR_NAMESPACE) get pod -l app=elasticsearch-master -o name)
port-forward-es: ## Run a port-forward to Elasticsearch on http://127.0.0.1:9200
	@echo "Reach Elasticsearch on http://127.0.0.1:9200"
	kubectl -n $(OPERATOR_NAMESPACE) port-forward $(ES_POD) 9200:9200


# ----------------------------------------------------------------------------------------------------------------------
# Install MetalLB
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-metallb
install-metallb: ## Install MetalLB to allow services of type LoadBalancer
	kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/$(METALLB_VERSION)/manifests/namespace.yaml
	kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/$(METALLB_VERSION)/manifests/metallb.yaml
	kubectl apply -f hack/metallb-config.yaml

# ----------------------------------------------------------------------------------------------------------------------
# Uninstall MetalLB
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: uninstall-metallb
uninstall-metallb: ## Uninstall MetalLB
	kubectl delete -f hack/metallb-config.yaml || true
	kubectl delete -f https://raw.githubusercontent.com/metallb/metallb/$(METALLB_VERSION)/manifests/metallb.yaml || true
	kubectl delete -f https://raw.githubusercontent.com/metallb/metallb/$(METALLB_VERSION)/manifests/namespace.yaml || true


# ----------------------------------------------------------------------------------------------------------------------
# Install the latest Istio version
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-istio
install-istio: get-istio ## Install the latest version of Istio into k8s
	$(eval ISTIO_HOME := $(shell find $(TOOLS_DIRECTORY) -maxdepth 1 -type d | grep istio))
	$(ISTIO_HOME)/bin/istioctl install --set profile=demo -y
	sleep 10
	kubectl -n istio-system get pod -l app=istio-egressgateway -o name | xargs \
		kubectl -n istio-system wait --for condition=ready --timeout 300s
	kubectl -n istio-system get pod -l app=istio-ingressgateway -o name | xargs \
		kubectl -n istio-system wait --for condition=ready --timeout 300s
	kubectl -n istio-system get pod -l app=istiod -o name | xargs \
		kubectl -n istio-system wait --for condition=ready --timeout 300s
	kubectl label namespace $(OPERATOR_NAMESPACE) istio-injection=enabled --overwrite=true

# ----------------------------------------------------------------------------------------------------------------------
# Uninstall Istio
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: uninstall-istio
uninstall-istio: get-istio ## Uninstall Istio from k8s
	$(eval ISTIO_HOME := $(shell find $(TOOLS_DIRECTORY) -maxdepth 1 -type d | grep istio))
	$(ISTIO_HOME)/bin/istioctl x uninstall --purge -y


# ----------------------------------------------------------------------------------------------------------------------
# Get the latest Istio version
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: get-istio
get-istio:
	./hack/get-istio-latest.sh "$(ISTIO_VERSION)" "$(TOOLS_DIRECTORY)"
	$(eval ISTIO_HOME := $(shell find $(TOOLS_DIRECTORY) -maxdepth 1 -type d | grep istio))
	@echo "Istio installed at $(ISTIO_HOME)"

# ----------------------------------------------------------------------------------------------------------------------
# Obtain the golangci-lint binary
# ----------------------------------------------------------------------------------------------------------------------
$(TOOLS_BIN)/golangci-lint:
	@mkdir -p $(TOOLS_BIN)
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(TOOLS_BIN) v1.29.0

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
	./mvnw -B -f java install -P docs -pl docs -DskipTests \
		-Doperator.version=$(VERSION) \
		-Doperator.image=$(OPERATOR_IMAGE) \
		-Doperator.utils.image=$(UTILS_IMAGE) \
		$(MAVEN_OPTIONS)
	mkdir -p $(BUILD_OUTPUT)/docs/images/images
	cp -R docs/images/* build/_output/docs/images/
	find examples/ -name \*.png -exec cp {} build/_output/docs/images/images/ \;

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
	kubectl create configmap coherence-kibana-dashboards --from-file=dashboards/kibana \
		--dry-run -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION)/coherence-kibana-dashboards.yaml

# ----------------------------------------------------------------------------------------------------------------------
# Release the Coherence Operator to the gh-pages branch.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: release-ghpages
release-ghpages:  helm-chart docs release-dashboards
	mkdir -p /tmp/coherence-operator || true
	cp -R $(BUILD_OUTPUT) /tmp/coherence-operator
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
	rm -rf dashboards/latest || true
	cp -R dashboards/$(VERSION) dashboards/latest
	git add -A dashboards/latest/*
	mkdir docs/$(VERSION) || true
	rm -rf docs/$(VERSION)/ || true
	mv $(BUILD_OUTPUT)/docs/ docs/$(VERSION)/
	rm -rf docs/latest
	cp -R docs/$(VERSION) docs/latest
	git add -A docs/*

	ls -ls docs

	mkdir -p charts || true
	cp $(BUILD_HELM)/coherence-operator-$(VERSION).tgz charts/
	helm repo index charts --url https://oracle.github.io/coherence-operator/charts
	git add charts/coherence-operator-$(VERSION).tgz
	git add charts/index.yaml
	ls -ls charts

	git status
endif
	git clean -d -f
	git status
	git commit -m "Release Coherence Operator version: $(VERSION)"
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
release: build-all-images push-release-images release-ghpages
endif


# ----------------------------------------------------------------------------------------------------------------------
# Create the third-party license file
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: license
license: $(TOOLS_BIN)/licensed
	mkdir .licenses || true
	touch .licenses/NOTICE
	$(TOOLS_BIN)/licensed cache
	$(TOOLS_BIN)/licensed notice
	cp .licenses/NOTICE THIRD_PARTY_LICENSES.txt


$(TOOLS_BIN)/licensed:
ifeq (Darwin, $(UNAME_S))
	curl -sSL https://github.com/github/licensed/releases/download/2.14.4/licensed-2.14.4-darwin-x64.tar.gz > licensed.tar.gz
else
	curl -sSL https://github.com/github/licensed/releases/download/2.14.4/licensed-2.14.4-linux-x64.tar.gz > licensed.tar.gz
endif
	tar -xzf licensed.tar.gz
	rm -f licensed.tar.gz
	mv ./licensed $(TOOLS_BIN)/licensed
	chmod +x $(TOOLS_BIN)/licensed


# ----------------------------------------------------------------------------------------------------------------------
# Generate Java client
# ----------------------------------------------------------------------------------------------------------------------
$(BUILD_OUTPUT)/java-client/java/gen/pom.xml: export LOCAL_MANIFEST_FILE := $(BUILD_OUTPUT)/java-client/crds/coherence.oracle.com_coherence.yaml
$(BUILD_OUTPUT)/java-client/java/gen/pom.xml: $(BUILD_TARGETS)/generate $(BUILD_TARGETS)/manifests $(KUSTOMIZE)
	docker pull ghcr.io/yue9944882/crd-model-gen:v1.0.6 || true
	rm -rf $(BUILD_OUTPUT)/java-client || true
	mkdir -p $(BUILD_OUTPUT)/java-client/crds
	mkdir -p $(BUILD_OUTPUT)/java-client/java/gen
	cp $(CURRDIR)/client/generate.sh $(BUILD_OUTPUT)/java-client/java/generate.sh
	chmod +x $(BUILD_OUTPUT)/java-client/java/generate.sh
	cp $(CURRDIR)/client/Dockerfile $(BUILD_OUTPUT)/java-client/java/Dockerfile
	docker build -f $(BUILD_OUTPUT)/java-client/java/Dockerfile -t crd-model-gen:custom $(BUILD_OUTPUT)/java-client/java
	$(KUSTOMIZE) build $(BUILD_DEPLOY)/crd > $(LOCAL_MANIFEST_FILE)
	docker run --rm --network host \
	  -v "$(LOCAL_MANIFEST_FILE)":"$(LOCAL_MANIFEST_FILE)" \
	  -v /var/run/docker.sock:/var/run/docker.sock \
	  -v "$(BUILD_OUTPUT)/java-client/java":"$(BUILD_OUTPUT)/java-client/java" \
	  crd-model-gen:custom \
	  /generate.sh \
	  -u $(LOCAL_MANIFEST_FILE) -n com.oracle.coherence -p com.oracle.coherence.k8s.client -o "$(BUILD_OUTPUT)/java-client/java"
	kind delete cluster || true

