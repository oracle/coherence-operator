# ----------------------------------------------------------------------------------------------------------------------
# Copyright (c) 2019, 2024, Oracle and/or its affiliates.
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
VERSION ?= 3.4.2
MVN_VERSION ?= $(VERSION)

# The version number to be replaced by this release
PREV_VERSION ?= 3.4.1
NEXT_VERSION := $(shell sh ./hack/next-version.sh "$(VERSION)")

# The operator version to use to run certification tests against
CERTIFICATION_VERSION ?= $(VERSION)

# The previous Operator version used to run the compatibility tests.
COMPATIBLE_VERSION  ?= 3.4.1
# The selector to use to find Operator Pods of the COMPATIBLE_VERSION (do not put in double quotes!!)
COMPATIBLE_SELECTOR ?= control-plane=coherence

# The GitHub project URL
PROJECT_URL = https://github.com/oracle/coherence-operator

KUBERNETES_DOC_VERSION=v1.30

# ----------------------------------------------------------------------------------------------------------------------
# The Coherence image to use for deployments that do not specify an image
# ----------------------------------------------------------------------------------------------------------------------
# The Coherence version to build against - must be a Java 8 compatible version
COHERENCE_VERSION     ?= 21.12.5
COHERENCE_VERSION_LTS ?= 14.1.2-0-0
# The default Coherence image the Operator will run if no image is specified
COHERENCE_IMAGE_REGISTRY ?= ghcr.io/oracle
COHERENCE_IMAGE_NAME     ?= coherence-ce
COHERENCE_IMAGE_TAG      ?= $(COHERENCE_VERSION_LTS)
COHERENCE_IMAGE          ?= $(COHERENCE_IMAGE_REGISTRY)/$(COHERENCE_IMAGE_NAME):$(COHERENCE_IMAGE_TAG)
COHERENCE_GROUP_ID       ?= com.oracle.coherence.ce
# The Java version that tests will be compiled to.
# This should match the version required by the COHERENCE_IMAGE version
BUILD_JAVA_VERSION        ?= 17
COHERENCE_TEST_BASE_IMAGE ?= gcr.io/distroless/java17-debian12

# This is the Coherence image that will be used in tests.
# Changing this variable will allow test builds to be run against different Coherence versions
# without altering the default image name.
TEST_COHERENCE_IMAGE ?= $(COHERENCE_IMAGE)
TEST_COHERENCE_VERSION ?= $(COHERENCE_VERSION)
TEST_COHERENCE_GID ?= com.oracle.coherence.ce

# The current working directory
CURRDIR := $(shell pwd)

GH_TOKEN ?=
ifeq ("$(GH_TOKEN)", "")
  GH_AUTH := 'Foo: Bar'
else
  GH_AUTH := 'authorization: Bearer $(GH_TOKEN)'
endif

# ----------------------------------------------------------------------------------------------------------------------
# By default we target amd64 as this is by far the most common local build environment
# We actually build images for amd64 and arm64
# ----------------------------------------------------------------------------------------------------------------------
UNAME_S         =  $(shell uname -s)
UNAME_M         =  $(shell uname -m)
ifeq (x86_64, $(UNAME_M))
	IMAGE_ARCH  = amd64
	ARCH        = amd64
else
	IMAGE_ARCH  = $(UNAME_M)
	ARCH        = $(UNAME_M)
endif

OS              ?= linux
GOPROXY         ?= https://proxy.golang.org

# ----------------------------------------------------------------------------------------------------------------------
# Set the location of the Operator SDK executable
# ----------------------------------------------------------------------------------------------------------------------
OPERATOR_SDK_VERSION := v1.9.0

# ----------------------------------------------------------------------------------------------------------------------
# Options to append to the Maven command
# ----------------------------------------------------------------------------------------------------------------------
MAVEN_OPTIONS ?= -Dmaven.wagon.httpconnectionManager.ttlSeconds=25 -Dmaven.wagon.http.retryHandler.count=3
MAVEN_BUILD_OPTS :=$(USE_MAVEN_SETTINGS) -Drevision=$(MVN_VERSION) -Dcoherence.version=$(COHERENCE_VERSION) -Dcoherence.version=$(COHERENCE_VERSION_LTS) -Dcoherence.groupId=$(COHERENCE_GROUP_ID) -Dcoherence.test.base.image=$(COHERENCE_TEST_BASE_IMAGE) -Dbuild.java.version=$(BUILD_JAVA_VERSION) $(MAVEN_OPTIONS)

# ----------------------------------------------------------------------------------------------------------------------
# Operator image names
# ----------------------------------------------------------------------------------------------------------------------
BASE_IMAGE_REGISTRY     ?= ghcr.io
BASE_IMAGE_REPO         ?= oracle
OPERATOR_IMAGE_REGISTRY ?= $(BASE_IMAGE_REGISTRY)/$(BASE_IMAGE_REPO)
RELEASE_IMAGE_PREFIX    ?= $(OPERATOR_IMAGE_REGISTRY)/
OPERATOR_IMAGE_NAME     := coherence-operator
OPERATOR_BASE_IMAGE     ?= scratch
OPERATOR_OL_BASE_IMAGE  ?= container-registry.oracle.com/java/jdk:17
OPERATOR_IMAGE_REPO     := $(RELEASE_IMAGE_PREFIX)$(OPERATOR_IMAGE_NAME)
OPERATOR_IMAGE          := $(OPERATOR_IMAGE_REPO):$(VERSION)
OPERATOR_IMAGE_DELVE    := $(OPERATOR_IMAGE_REPO):delve
OPERATOR_IMAGE_DEBUG    := $(OPERATOR_IMAGE_REPO):debug
TEST_BASE_IMAGE         ?= $(OPERATOR_IMAGE_REPO)-test-base:$(VERSION)
# The Operator images to push
OPERATOR_RELEASE_REPO   ?= $(OPERATOR_IMAGE_REPO)
OPERATOR_RELEASE_IMAGE  := $(OPERATOR_RELEASE_REPO):$(VERSION)
TEST_BASE_RELEASE_IMAGE := $(OPERATOR_RELEASE_REPO)-test-base:$(VERSION)
BUNDLE_RELEASE_IMAGE    := $(OPERATOR_RELEASE_REPO)-bundle:$(VERSION)

OPERATOR_PACKAGE_PREFIX := $(OPERATOR_IMAGE_REPO)-package
OPERATOR_PACKAGE_IMAGE  := $(OPERATOR_PACKAGE_PREFIX):$(VERSION)
OPERATOR_REPO_PREFIX    := $(OPERATOR_IMAGE_REPO)-repo
OPERATOR_REPO_IMAGE     := $(OPERATOR_REPO_PREFIX):$(VERSION)

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

# default test namespace, as in test/e2e/helper/proj_helpers.go
OPERATOR_NAMESPACE ?= operator-test
# default test cluster namespace, as in test/e2e/helper/proj_helpers.go
CLUSTER_NAMESPACE ?= coherence-test
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

ifeq (Darwin, $(UNAME_S))
	SED = sed -i ''
else
	SED = sed -i
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
override BUILD_BIN_AMD64     := $(BUILD_BIN)/linux/amd64
override BUILD_BIN_ARM64     := $(BUILD_BIN)/linux/arm64
override BUILD_DEPLOY        := $(BUILD_OUTPUT)/config
override BUILD_HELM          := $(BUILD_OUTPUT)/helm-charts
override BUILD_MANIFESTS     := $(BUILD_OUTPUT)/manifests
override BUILD_MANIFESTS_PKG := $(BUILD_OUTPUT)/coherence-operator-manifests.tar.gz
override BUILD_PROPS         := $(BUILD_OUTPUT)/build.properties
override BUILD_TARGETS       := $(BUILD_OUTPUT)/targets
override SCRIPTS_DIR         := $(CURRDIR)/hack
override EXAMPLES_DIR        := $(CURRDIR)/examples
override TEST_LOGS_DIR       := $(BUILD_OUTPUT)/test-logs
override TANZU_DIR           := $(BUILD_OUTPUT)/tanzu
override TANZU_PACKAGE_DIR   := $(BUILD_OUTPUT)/tanzu/package
override TANZU_REPO_DIR      := $(BUILD_OUTPUT)/tanzu/repo


# ----------------------------------------------------------------------------------------------------------------------
# Set the location of various build tools
# ----------------------------------------------------------------------------------------------------------------------
TOOLS_DIRECTORY   = $(CURRDIR)/build/tools
TOOLS_BIN         = $(TOOLS_DIRECTORY)/bin
TOOLS_MANIFESTS   = $(TOOLS_DIRECTORY)/manifests
OPERATOR_SDK_HOME = $(TOOLS_DIRECTORY)/sdk/$(UNAME_S)-$(UNAME_M)
OPERATOR_SDK      = $(OPERATOR_SDK_HOME)/operator-sdk

# ----------------------------------------------------------------------------------------------------------------------
# The ttl.sh images used in integration tests
# ----------------------------------------------------------------------------------------------------------------------
TTL_REGISTRY                       := ttl.sh
TTL_TIMEOUT                        := 1h
TTL_UUID_FILE                      := $(BUILD_OUTPUT)/ttl-uuid.txt
TTL_UUID                           := $(shell if [ -f $(TTL_UUID_FILE) ]; then cat $(TTL_UUID_FILE); else uuidgen | tr A-Z a-z > $(TTL_UUID_FILE) && cat $(TTL_UUID_FILE); fi)
TTL_OPERATOR_IMAGE                 := $(TTL_REGISTRY)/coherence/$(TTL_UUID)/$(OPERATOR_IMAGE_NAME):$(TTL_TIMEOUT)
TTL_PACKAGE_IMAGE                  := $(TTL_REGISTRY)/coherence/$(TTL_UUID)/$(OPERATOR_IMAGE_NAME)-package:$(TTL_TIMEOUT)
TTL_REPO_IMAGE                     := $(TTL_REGISTRY)/coherence/$(TTL_UUID)/$(OPERATOR_IMAGE_NAME)-repo:$(TTL_TIMEOUT)
TTL_APPLICATION_IMAGE              := $(TTL_REGISTRY)/coherence/$(TTL_UUID)/operator-test:$(TTL_TIMEOUT)
TTL_COMPATIBILITY_IMAGE            := $(TTL_REGISTRY)/coherence/$(TTL_UUID)/operator-test-compatibility:$(TTL_TIMEOUT)
TTL_APPLICATION_IMAGE_CLIENT       := $(TTL_REGISTRY)/coherence/$(TTL_UUID)/operator-test-client:$(TTL_TIMEOUT)
TTL_APPLICATION_IMAGE_HELIDON      := $(TTL_REGISTRY)/coherence/$(TTL_UUID)/operator-test-helidon:$(TTL_TIMEOUT)
TTL_APPLICATION_IMAGE_SPRING       := $(TTL_REGISTRY)/coherence/$(TTL_UUID)/operator-test-spring:$(TTL_TIMEOUT)
TTL_APPLICATION_IMAGE_SPRING_FAT   := $(TTL_REGISTRY)/coherence/$(TTL_UUID)/operator-test-spring-fat:$(TTL_TIMEOUT)
TTL_APPLICATION_IMAGE_SPRING_CNBP  := $(TTL_REGISTRY)/coherence/$(TTL_UUID)/operator-test-spring-cnbp:$(TTL_TIMEOUT)

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
JAVA_FILES       = $(shell find ./java -type f)

TEST_SSL_SECRET := coherence-ssl-secret

# ----------------------------------------------------------------------------------------------------------------------
# Prometheus Operator settings (used in integration tests)
# ----------------------------------------------------------------------------------------------------------------------
# The version of kube-prometheus to use (main = latest main branch from https://github.com/prometheus-operator/kube-prometheus)
PROMETHEUS_VERSION           ?= main
PROMETHEUS_HOME               = $(TOOLS_DIRECTORY)/prometheus/$(PROMETHEUS_VERSION)
PROMETHEUS_NAMESPACE         ?= monitoring
PROMETHEUS_ADAPTER_VERSION   ?= 2.5.0
GRAFANA_DASHBOARDS           ?= dashboards/grafana/

# ----------------------------------------------------------------------------------------------------------------------
# MetalLB load balancer settings
# ----------------------------------------------------------------------------------------------------------------------
METALLB_VERSION ?= v0.12.1

# ----------------------------------------------------------------------------------------------------------------------
# Istio settings
# ----------------------------------------------------------------------------------------------------------------------
# The version of Istio to install, leave empty for the latest
ISTIO_VERSION ?=

# ----------------------------------------------------------------------------------------------------------------------
# Tanzu settings
# ----------------------------------------------------------------------------------------------------------------------
# The version of Tanzu to install, leave empty for the latest
TANZU_VERSION ?=
TANZU =

# ======================================================================================================================
# Makefile targets start here
# ======================================================================================================================

# ----------------------------------------------------------------------------------------------------------------------
# Display the Makefile help - this is a list of the targets with a description.
# This target MUST be the first target in the Makefile so that it is run when running make with no arguments
# ----------------------------------------------------------------------------------------------------------------------
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

ttl-uuid:
	echo "TTL UUID: $(TTL_UUID)"

# ======================================================================================================================
# Build targets
# ======================================================================================================================
##@ Build

.PHONY: all
all: build-all-images helm-chart ## Build all the Coherence Operator artefacts and images

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
	@mkdir -p $(TOOLS_MANIFESTS)
	# create build.properties
	rm -f $(BUILD_PROPS)
	printf "COHERENCE_IMAGE=$(COHERENCE_IMAGE)\n\
	COHERENCE_IMAGE_REGISTRY=$(COHERENCE_IMAGE_REGISTRY)\n\
	COHERENCE_IMAGE_NAME=$(COHERENCE_IMAGE_NAME)\n\
	COHERENCE_IMAGE_TAG=$(COHERENCE_IMAGE_TAG)\n\
	OPERATOR_IMAGE_REGISTRY=$(OPERATOR_IMAGE_REGISTRY)\n\
	RELEASE_IMAGE_PREFIX=$(RELEASE_IMAGE_PREFIX)\n\
	OPERATOR_IMAGE_NAME=$(OPERATOR_IMAGE_NAME)\n\
	OPERATOR_IMAGE=$(OPERATOR_IMAGE)\n\
	VERSION=$(VERSION)\n\
	OPERATOR_PACKAGE_IMAGE=$(OPERATOR_PACKAGE_IMAGE)\n" > $(BUILD_PROPS)

# ----------------------------------------------------------------------------------------------------------------------
# Clean-up all of the build artifacts
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: clean
clean: ## Cleans the build
	-rm -rf $(BUILD_OUTPUT) || true
	-rm -rf $(BUILD_BIN) || true
	-rm -rf bundle || true
	rm config/crd/bases/*.yaml || true
	rm -rf config/crd-small || true
	rm pkg/data/zz_generated_*.go || true
	rm pkg/data/assets/*.yaml || true
	rm pkg/data/assets/*.json || true
	./mvnw -f java clean $(MAVEN_BUILD_OPTS)
	./mvnw -f examples clean $(MAVEN_BUILD_OPTS)

# ----------------------------------------------------------------------------------------------------------------------
# Clean-up all of the locally downloaded build tools
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: clean-tools
clean-tools: ## Cleans the locally downloaded build tools (i.e. need a new tool version)
	-rm -rf $(TOOLS_BIN) || true


# ----------------------------------------------------------------------------------------------------------------------
# Builds the Operator
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-operator
build-operator: $(BUILD_TARGETS)/build-operator ## Build the Coherence Operator image

$(BUILD_TARGETS)/build-operator: $(BUILD_BIN)/runner $(BUILD_TARGETS)/java $(BUILD_TARGETS)/cli
	docker build --platform linux/amd64 --no-cache --build-arg version=$(VERSION) \
		--build-arg BASE_IMAGE=$(OPERATOR_BASE_IMAGE) \
		--build-arg coherence_image=$(COHERENCE_IMAGE) \
		--build-arg operator_image=$(OPERATOR_IMAGE) \
		--build-arg target=amd64 \
		. -t $(OPERATOR_IMAGE)-amd64
	docker build --platform linux/arm64 --no-cache --build-arg version=$(VERSION) \
		--build-arg BASE_IMAGE=$(OPERATOR_BASE_IMAGE) \
		--build-arg coherence_image=$(COHERENCE_IMAGE) \
		--build-arg operator_image=$(OPERATOR_IMAGE) \
		--build-arg target=arm64 \
		. -t $(OPERATOR_IMAGE)-arm64
	docker tag $(OPERATOR_IMAGE)-$(IMAGE_ARCH) $(OPERATOR_IMAGE)
	touch $(BUILD_TARGETS)/build-operator

.PHONY: build-operator-with-tools
build-operator-with-tools: $(BUILD_BIN)/runner $(BUILD_TARGETS)/java ## Build the Coherence Operator image on OL-8 with debug tools
	mkdir -p $(BUILD_OUTPUT)/images || true
	cat Dockerfile debug/Tools.Dockerfile > $(BUILD_OUTPUT)/images/Dockerfile
	docker build --no-cache --build-arg version=$(VERSION) \
		--build-arg BASE_IMAGE=$(OPERATOR_OL_BASE_IMAGE) \
		--build-arg coherence_image=$(COHERENCE_IMAGE) \
		--build-arg operator_image=$(OPERATOR_IMAGE) \
		--build-arg target=amd64 \
		-f $(BUILD_OUTPUT)/images/Dockerfile \
		. -t $(OPERATOR_IMAGE)

.PHONY: build-operator-debug
build-operator-debug: $(BUILD_BIN)/linux/amd64/runner-debug $(BUILD_TARGETS)/java ## Build the Coherence Operator image with the Delve debugger
	docker build --no-cache --build-arg version=$(VERSION) \
		--build-arg BASE_IMAGE=$(OPERATOR_IMAGE_DELVE) \
		--build-arg coherence_image=$(COHERENCE_IMAGE) \
		--build-arg operator_image=$(OPERATOR_IMAGE) \
		--build-arg target=amd64 \
		-f debug/Dockerfile \
		. -t $(OPERATOR_IMAGE_DEBUG)

build-delve-image: ## Build the Coherence Operator Delve debugger base image
	docker build -f debug/Base.Dockerfile -t $(OPERATOR_IMAGE_DELVE) debug

$(BUILD_BIN)/linux/amd64/runner-debug: $(BUILD_PROPS) $(GOS) $(BUILD_TARGETS)/generate $(BUILD_TARGETS)/manifests
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -gcflags "-N -l" -ldflags "$(LDFLAGS)" -a -o $(BUILD_BIN)/linux/amd64/runner-debug ./runner
	chmod +x $(BUILD_BIN)/linux/amd64/runner-debug

# ----------------------------------------------------------------------------------------------------------------------
# Build the Operator images without the test images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-operator-images
build-operator-images: $(BUILD_TARGETS)/build-operator ## Build all operator images

# ----------------------------------------------------------------------------------------------------------------------
# Build the Operator Test images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-test-images
build-test-images: $(BUILD_TARGETS)/java build-client-image build-basic-test-image ## Build all of the test images
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
build-basic-test-image: $(BUILD_TARGETS)/java ## Build the basic Operator test image
	./mvnw -B -f java/operator-test clean package jib:dockerBuild -DskipTests -Djib.to.image=$(TEST_APPLICATION_IMAGE) $(MAVEN_BUILD_OPTS) -Dcoherence.version=$(COHERENCE_IMAGE_TAG)

.PHONY: build-client-image
build-client-image: ## Build the test client image
	./mvnw -B -f java/operator-test-client package jib:dockerBuild -DskipTests -Djib.to.image=$(TEST_APPLICATION_IMAGE_CLIENT) $(MAVEN_BUILD_OPTS)

# ----------------------------------------------------------------------------------------------------------------------
# Build all of the Docker images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-all-images
build-all-images: $(BUILD_TARGETS)/build-operator build-test-images build-compatibility-image ## Build all images (including tests)

# ----------------------------------------------------------------------------------------------------------------------
# Ensure Operator SDK is at the correct version
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: ensure-sdk
ensure-sdk:
	@echo "Ensuring Operator SDK is present at version $(OPERATOR_SDK_VERSION)"
	$(SCRIPTS_DIR)/ensure-sdk.sh $(OPERATOR_SDK_VERSION) $(OPERATOR_SDK_HOME)

# ----------------------------------------------------------------------------------------------------------------------
# Internal make step that builds the Operator runner artifacts utility
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-runner
build-runner: $(BUILD_BIN)/runner  ## Build the Coherence Operator runner binary

$(BUILD_BIN)/runner: $(BUILD_PROPS) $(GOS) $(BUILD_TARGETS)/generate $(BUILD_TARGETS)/manifests
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -trimpath -ldflags "$(LDFLAGS)" -o $(BUILD_BIN)/runner ./runner
	mkdir -p $(BUILD_BIN_AMD64) || true
	cp -f $(BUILD_BIN)/runner $(BUILD_BIN_AMD64)/runner
	mkdir -p $(BUILD_BIN_ARM64)/linux || true
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GO111MODULE=on go build -trimpath -ldflags "$(LDFLAGS)" -a -o $(BUILD_BIN_ARM64)/runner ./runner

# ----------------------------------------------------------------------------------------------------------------------
# Build the Java artifacts
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-mvn
build-mvn: $(BUILD_TARGETS)/java ## Build the Java artefacts

$(BUILD_TARGETS)/java: $(JAVA_FILES)
	./mvnw -B -f java clean install -DskipTests $(MAVEN_BUILD_OPTS)
	touch $(BUILD_TARGETS)/java


# ---------------------------------------------------------------------------
# Build the Coherence operator Helm chart and package it into a tar.gz
# ---------------------------------------------------------------------------
.PHONY: helm-chart
helm-chart: $(BUILD_PROPS) $(BUILD_HELM)/coherence-operator-$(VERSION).tgz   ## Build the Coherence Operator Helm chart

$(BUILD_HELM)/coherence-operator-$(VERSION).tgz: $(BUILD_PROPS) $(HELM_FILES) $(BUILD_TARGETS)/generate $(BUILD_TARGETS)/manifests $(TOOLS_BIN)/kustomize
# Copy the Helm chart from the source location to the distribution folder
	-mkdir -p $(BUILD_HELM)
	cp -R ./helm-charts/coherence-operator $(BUILD_HELM)
	$(call replaceprop,$(BUILD_HELM)/coherence-operator/Chart.yaml $(BUILD_HELM)/coherence-operator/values.yaml $(BUILD_HELM)/coherence-operator/templates/deployment.yaml)
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
		filename="$${i}"; \
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

config/crd/bases/coherence.oracle.com_coherence.yaml: $(TOOLS_BIN)/kustomize $(API_GO_FILES) $(TOOLS_BIN)/controller-gen
	$(CONTROLLER_GEN) "crd:crdVersions={v1}" \
	  rbac:roleName=manager-role paths="{./api/...,./controllers/...}" \
	  output:crd:dir=config/crd/bases
	cp -R config/crd/ config/crd-small
	$(CONTROLLER_GEN) "crd:crdVersions={v1},maxDescLen=0" \
	  rbac:roleName=manager-role paths="{./api/...,./controllers/...}" \
	  output:crd:dir=config/crd-small/bases
	cd config/crd && $(KUSTOMIZE) edit add label "app.kubernetes.io/version:$(VERSION)" -f
	$(KUSTOMIZE) build config/crd -o $(BUILD_ASSETS)/
	cd config/crd-small && $(KUSTOMIZE) edit add label "app.kubernetes.io/version:$(VERSION)" -f

# ----------------------------------------------------------------------------------------------------------------------
# Generate the config.json file used by the Operator for default configuration values
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: generate-config
generate-config: $(BUILD_PROPS) $(BUILD_OUTPUT)/config.json

$(BUILD_OUTPUT)/config.json:
	@echo "Generating Operator config"
	@printf "{\n \
	  \"coherence-image\": \"$(COHERENCE_IMAGE)\",\n \
	  \"operator-image\": \"$(OPERATOR_RELEASE_IMAGE)\"\n}\n" > $(BUILD_OUTPUT)/config.json
	cp $(BUILD_OUTPUT)/config.json $(BUILD_ASSETS)/config.json

# ----------------------------------------------------------------------------------------------------------------------
# Generate code, configuration and docs.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: generate
generate: $(BUILD_TARGETS)/generate  ## Run Kubebuilder code and configuration generation

$(BUILD_TARGETS)/generate: $(BUILD_PROPS) $(BUILD_OUTPUT)/config.json api/v1/zz_generated.deepcopy.go
	touch $(BUILD_TARGETS)/generate

api/v1/zz_generated.deepcopy.go: $(API_GO_FILES) $(TOOLS_BIN)/controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./api/..."

# ----------------------------------------------------------------------------------------------------------------------
# Generate API docs
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: api-doc-gen
api-doc-gen: docs/about/04_coherence_spec.adoc  ## Generate API documentation

docs/about/04_coherence_spec.adoc: export KUBERNETES_DOC_VERSION := $(KUBERNETES_DOC_VERSION)
docs/about/04_coherence_spec.adoc: $(API_GO_FILES) utils/docgen/main.go
	@echo "Generating CRD Doc"
	go run ./utils/docgen/ \
		api/v1/coherenceresourcespec_types.go \
		api/v1/coherence_types.go \
		api/v1/coherenceresource_types.go \
		> docs/about/04_coherence_spec.adoc

# ----------------------------------------------------------------------------------------------------------------------
# Generate the keys and certs used in tests.
# ----------------------------------------------------------------------------------------------------------------------
$(BUILD_OUTPUT)/certs:
	@echo "Generating test keys and certs"
	$(SCRIPTS_DIR)/keys.sh

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
	$(TOOLS_BIN)/golangci-lint run -v --timeout=5m --exclude='G402:' --exclude='G101:' --exclude='G114:' --skip-dirs=.*/fakes --skip-files=zz_.*,generated/*,pkg/data/assets... ./api/... ./controllers/... ./pkg/... ./runner/...
	$(TOOLS_BIN)/golangci-lint run -v --timeout=5m --exclude='G107:' --exclude='G101:' --exclude='G112:' --exclude='SA4005:' --exclude='should not use dot imports' ./test/... ./pkg/fakes/...


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
	  -X hack/install-cli.sh \
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
	  -X tanzu/package/package.yml \
	  -X tanzu/package/values.yml \
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
run: create-namespace ## run the Operator locally
	go run -ldflags "$(LDFLAGS)" ./runner/main.go operator --skip-service-suspend=true --coherence-dev-mode=true \
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
run-debug: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
run-debug: create-namespace ## run the Operator locally with Delve debugger
	dlv debug ./runner --headless --listen=:2345 --api-version=2 --accept-multiclient \
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
	$(SCRIPTS_DIR)/kill-local.sh

# ======================================================================================================================
# Targets related to Operator Lifecycle Manager and the Operator SDK
# ======================================================================================================================
##@ Operator Lifecycle Manager

# ----------------------------------------------------------------------------------------------------------------------
# Generate bundle manifests and metadata, then validate generated files.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: bundle
bundle: $(BUILD_PROPS) ensure-sdk $(TOOLS_BIN)/kustomize $(BUILD_TARGETS)/manifests  ## Generate OLM bundle manifests and metadata, then validate generated files.
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
	curl -sSLo $(OPM) https://github.com/operator-framework/operator-registry/releases/download/v1.15.1/$${OS}-$${ARCH}-opm --header $(GH_AUTH) ;\
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
test-operator: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
test-operator: $(BUILD_PROPS) $(BUILD_TARGETS)/manifests $(BUILD_TARGETS)/generate gotestsum  ## Run the Operator unit tests
	@echo "Running operator tests"
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-test.xml \
	  -- $(GO_TEST_FLAGS) -v ./api/... ./controllers/... ./pkg/...

# ----------------------------------------------------------------------------------------------------------------------
# Build and test the Java artifacts
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: test-mvn
test-mvn: $(BUILD_OUTPUT)/certs $(BUILD_TARGETS)/java  ## Run the Java artefact tests
	./mvnw -B -f java verify -Dtest.certs.location=$(BUILD_OUTPUT)/certs $(MAVEN_BUILD_OPTS)


# ----------------------------------------------------------------------------------------------------------------------
# Run all unit tests (both Go and Java)
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: test-all
test-all: test-mvn test-operator  ## Run all unit tests

ENVTEST_K8S_VERSION =  1.31.0
ENVTEST_VERSION     ?= release-0.19

.PHONY: envtest
envtest: $(TOOLS_BIN)/setup-envtest ## Download setup-envtest locally if necessary.

envtest-delete:
	$(TOOLS_BIN)/setup-envtest --bin-dir $(TOOLS_BIN) cleanup latest-on-disk
	rm -rf $(TOOLS_BIN)/k8s || true

$(TOOLS_BIN)/setup-envtest:
	test -s $(TOOLS_BIN)/setup-envtest || GOBIN=$(TOOLS_BIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@$(ENVTEST_VERSION)
	ls -al $(TOOLS_BIN)

k8stools: $(TOOLS_BIN)/k8s

$(TOOLS_BIN)/k8s: $(TOOLS_BIN)/setup-envtest
	mkdir -p $(TOOLS_BIN)/k8s || true
	$(TOOLS_BIN)/setup-envtest --bin-dir $(TOOLS_BIN) use $(ENVTEST_K8S_VERSION)



# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end tests that require a k8s cluster using
# a LOCAL operator instance (i.e. the operator is not deployed to k8s).
# These tests will use whichever k8s cluster the local environment
# is pointing to.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: e2e-local-test
e2e-local-test: export CGO_ENABLED = 0
e2e-local-test: export OPERATOR_NAMESPACE := $(OPERATOR_NAMESPACE)
e2e-local-test: export CLUSTER_NAMESPACE := $(CLUSTER_NAMESPACE)
e2e-local-test: export OPERATOR_NAMESPACE_CLIENT := $(OPERATOR_NAMESPACE_CLIENT)
e2e-local-test: export BUILD_OUTPUT := $(BUILD_OUTPUT)
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
run-e2e-test: export CLUSTER_NAMESPACE := $(CLUSTER_NAMESPACE)
run-e2e-test: export OPERATOR_NAMESPACE_CLIENT := $(OPERATOR_NAMESPACE_CLIENT)
run-e2e-test: export BUILD_OUTPUT := $(BUILD_OUTPUT)
run-e2e-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
run-e2e-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
run-e2e-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
run-e2e-test: export TEST_ASSET_KUBECTL := $(TEST_ASSET_KUBECTL)
run-e2e-test: export VERSION := $(VERSION)
run-e2e-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
run-e2e-test: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
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
# Executes the Go end-to-end tests that require a K3d cluster using
# a LOCAL operator instance (i.e. the operator is not deployed to k8s).
# These tests will use whichever K3d cluster the local environment
# is pointing to.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: e2e-k3d-test
e2e-k3d-test: export CGO_ENABLED = 0
e2e-k3d-test: export OPERATOR_NAMESPACE := $(OPERATOR_NAMESPACE)
e2e-k3d-test: export CLUSTER_NAMESPACE := $(CLUSTER_NAMESPACE)
e2e-k3d-test: export OPERATOR_NAMESPACE_CLIENT := $(OPERATOR_NAMESPACE_CLIENT)
e2e-k3d-test: export BUILD_OUTPUT := $(BUILD_OUTPUT)
e2e-k3d-test: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
e2e-k3d-test: export TEST_APPLICATION_IMAGE_CLIENT := $(TEST_APPLICATION_IMAGE_CLIENT)
e2e-k3d-test: export TEST_APPLICATION_IMAGE_HELIDON := $(TEST_APPLICATION_IMAGE_HELIDON)
e2e-k3d-test: export TEST_APPLICATION_IMAGE_SPRING := $(TEST_APPLICATION_IMAGE_SPRING)
e2e-k3d-test: export TEST_APPLICATION_IMAGE_SPRING_FAT := $(TEST_APPLICATION_IMAGE_SPRING_FAT)
e2e-k3d-test: export TEST_APPLICATION_IMAGE_SPRING_CNBP := $(TEST_APPLICATION_IMAGE_SPRING_CNBP)
e2e-k3d-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
e2e-k3d-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
e2e-k3d-test: export COH_SKIP_SITE := true
e2e-k3d-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
e2e-k3d-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
e2e-k3d-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
e2e-k3d-test: export TEST_ASSET_KUBECTL := $(TEST_ASSET_KUBECTL)
e2e-k3d-test: export LOCAL_STORAGE_RESTART := $(LOCAL_STORAGE_RESTART)
e2e-k3d-test: export VERSION := $(VERSION)
e2e-k3d-test: export MVN_VERSION := $(MVN_VERSION)
e2e-k3d-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
e2e-k3d-test: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
e2e-k3d-test: reset-namespace create-ssl-secrets install-crds gotestsum undeploy   ## Run the Operator end-to-end 'local' functional tests using a local Operator instance
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-k3d-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/e2e/large-cluster/...

# ----------------------------------------------------------------------------------------------------------------------
# Run the end-to-end Coherence client tests.
# ----------------------------------------------------------------------------------------------------------------------
e2e-client-test: export CGO_ENABLED = 0
e2e-client-test: export CLIENT_CLASSPATH := $(CURRDIR)/java/operator-test-client/target/operator-test-client-$(MVN_VERSION).jar:$(CURRDIR)/java/operator-test-client/target/lib/*
e2e-client-test: export OPERATOR_NAMESPACE := $(OPERATOR_NAMESPACE)
e2e-client-test: export OPERATOR_NAMESPACE_CLIENT := $(OPERATOR_NAMESPACE_CLIENT)
e2e-client-test: export BUILD_OUTPUT := $(BUILD_OUTPUT)
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
e2e-client-test: build-operator-images build-client-image reset-namespace create-ssl-secrets install-crds gotestsum undeploy   ## Run the end-to-end Coherence client tests using a local Operator deployment
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-client-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/e2e/clients/...


# ----------------------------------------------------------------------------------------------------------------------
# Run the end-to-end Helm chart tests.
# ----------------------------------------------------------------------------------------------------------------------
e2e-helm-test: export OPERATOR_IMAGE_REGISTRY := $(OPERATOR_IMAGE_REGISTRY)
e2e-helm-test: export OPERATOR_IMAGE_NAME := $(OPERATOR_IMAGE_NAME)
e2e-helm-test: export VERSION := $(VERSION)
e2e-helm-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
e2e-helm-test: export COHERENCE_IMAGE_REGISTRY := $(COHERENCE_IMAGE_REGISTRY)
e2e-helm-test: export COHERENCE_IMAGE_NAME := $(COHERENCE_IMAGE_NAME)
e2e-helm-test: export COHERENCE_IMAGE_TAG := $(COHERENCE_IMAGE_TAG)
e2e-helm-test: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
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
run-prometheus-test: gotestsum
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-prometheus-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/e2e/prometheus/...

# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end Operator backwards compatibility tests to ensure upgrades from previous Operator versions
# work and do not impact running clusters, etc.
# These tests will use whichever k8s cluster the local environment is pointing to.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: compatibility-test
compatibility-test: undeploy build-all-images $(BUILD_HELM)/coherence-operator-$(VERSION).tgz undeploy clean-namespace reset-namespace gotestsum just-compatibility-test  ## Run the Operator backwards compatibility tests

.PHONY: just-compatibility-test
just-compatibility-test: export CGO_ENABLED = 0
just-compatibility-test: export OPERATOR_NAMESPACE := $(OPERATOR_NAMESPACE)
just-compatibility-test: export CLUSTER_NAMESPACE := $(CLUSTER_NAMESPACE)
just-compatibility-test: export BUILD_OUTPUT := $(BUILD_OUTPUT)
just-compatibility-test: export TEST_APPLICATION_IMAGE := $(TEST_APPLICATION_IMAGE)
just-compatibility-test: export TEST_APPLICATION_IMAGE_CLIENT := $(TEST_APPLICATION_IMAGE_CLIENT)
just-compatibility-test: export TEST_APPLICATION_IMAGE_HELIDON := $(TEST_APPLICATION_IMAGE_HELIDON)
just-compatibility-test: export TEST_APPLICATION_IMAGE_SPRING := $(TEST_APPLICATION_IMAGE_SPRING)
just-compatibility-test: export TEST_APPLICATION_IMAGE_SPRING_FAT := $(TEST_APPLICATION_IMAGE_SPRING_FAT)
just-compatibility-test: export TEST_APPLICATION_IMAGE_SPRING_CNBP := $(TEST_APPLICATION_IMAGE_SPRING_CNBP)
just-compatibility-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
just-compatibility-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
just-compatibility-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
just-compatibility-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
just-compatibility-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
just-compatibility-test: export VERSION := $(VERSION)
just-compatibility-test: export COMPATIBLE_VERSION := $(COMPATIBLE_VERSION)
just-compatibility-test: export COMPATIBLE_SELECTOR := $(COMPATIBLE_SELECTOR)
just-compatibility-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
just-compatibility-test: export COHERENCE_IMAGE := $(COHERENCE_IMAGE)
just-compatibility-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
just-compatibility-test:  ## Run the Operator backwards compatibility tests WITHOUT building anything
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
	@echo "Running certification tests"
	$(MAKE) run-certification  $${MF} \
	; rc=$$? \
	; $(MAKE) cleanup-certification $${MF} \
	; exit $$rc


# ----------------------------------------------------------------------------------------------------------------------
# Install the Operator prior to running compatibility tests.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-certification
install-certification: $(BUILD_TARGETS)/build-operator prepare-network-policies reset-namespace create-ssl-secrets deploy-and-wait

# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end Operator certification tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# Note that the namespace will be created if it does not exist.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: run-certification
run-certification: export CGO_ENABLED = 0
run-certification: export OPERATOR_NAMESPACE := $(OPERATOR_NAMESPACE)
run-certification: export CLUSTER_NAMESPACE := $(CLUSTER_NAMESPACE)
run-certification: export BUILD_OUTPUT := $(BUILD_OUTPUT)
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
run-certification: gotestsum
	$(GOTESTSUM) --format standard-verbose --junitfile $(TEST_LOGS_DIR)/operator-e2e-certification-test.xml \
	  -- $(GO_TEST_FLAGS_E2E) ./test/certification/...

# ----------------------------------------------------------------------------------------------------------------------
# Clean up after to running compatibility tests.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: cleanup-certification
cleanup-certification: undeploy clean-namespace

# ----------------------------------------------------------------------------------------------------------------------
# Executes the Go end-to-end Operator Kubernetes network policy tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# Note that the namespace will be created if it does not exist.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: network-policy-test
network-policy-test: export MF = $(MAKEFLAGS)
network-policy-test: install-network-policy-tests     ## Run the Operator Kubernetes network policy tests
	$(MAKE) run-certification  $${MF} \
	; rc=$$? \
	; $(MAKE) cleanup-certification $${MF} \
	; exit $$rc

# ----------------------------------------------------------------------------------------------------------------------
# Install the Operator prior to running network policy tests.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-network-policy-tests
install-network-policy-tests: $(BUILD_TARGETS)/build-operator reset-namespace install-network-policies create-ssl-secrets deploy-and-wait

# ----------------------------------------------------------------------------------------------------------------------
# Install the network policies from the examples
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-network-policies
install-network-policies: install-operator-network-policies install-coherence-network-policies
	@echo "API Server info"
	kubectl get svc -o wide
	kubectl get endpoints kubernetes
	@echo "Network policies installed in $(OPERATOR_NAMESPACE)"
	kubectl get networkpolicy -n $(OPERATOR_NAMESPACE)
	@echo "Network policies installed in $(CLUSTER_NAMESPACE)"
	kubectl get networkpolicy -n $(CLUSTER_NAMESPACE)

# ----------------------------------------------------------------------------------------------------------------------
# Prepare a copy of the example network policies
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: prepare-network-policies
prepare-network-policies: export IP1=$(shell kubectl -n default get endpoints kubernetes -o jsonpath='{.subsets[0].addresses[0].ip}')
prepare-network-policies: export IP2=$(shell kubectl -n default get svc kubernetes -o jsonpath='{.spec.clusterIP}')
prepare-network-policies: export API_PORT=$(shell kubectl -n default get endpoints kubernetes -o jsonpath='{.subsets[0].ports[0].port}')
prepare-network-policies:
	mkdir -p $(BUILD_OUTPUT)/network-policies
	cp $(EXAMPLES_DIR)/095_network_policies/*.sh $(BUILD_OUTPUT)/network-policies
	cp -R $(EXAMPLES_DIR)/095_network_policies/manifests $(BUILD_OUTPUT)/network-policies
	$(SED) -e 's/172.18.0.2/${IP1}/g' $(BUILD_OUTPUT)/network-policies/manifests/allow-k8s-api-server.yaml
	$(SED) -e 's/10.96.0.1/${IP2}/g' $(BUILD_OUTPUT)/network-policies/manifests/allow-k8s-api-server.yaml
	$(SED) -e 's/6443/${API_PORT}/g' $(BUILD_OUTPUT)/network-policies/manifests/allow-k8s-api-server.yaml
	$(SED) -e 's/172.18.0.2/${IP1}/g' $(BUILD_OUTPUT)/network-policies/manifests/allow-webhook-ingress-from-api-server.yaml
	$(SED) -e 's/10.96.0.1/${IP2}/g' $(BUILD_OUTPUT)/network-policies/manifests/allow-webhook-ingress-from-api-server.yaml

# ----------------------------------------------------------------------------------------------------------------------
# Uninstall the network policies from the examples
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: uninstall-network-policies
uninstall-network-policies: uninstall-operator-network-policies uninstall-coherence-network-policies
	@echo "Network policies installed in $(OPERATOR_NAMESPACE)"
	kubectl get networkpolicy -n $(OPERATOR_NAMESPACE)
	@echo "Network policies installed in $(CLUSTER_NAMESPACE)"
	kubectl get networkpolicy -n $(CLUSTER_NAMESPACE)

# ----------------------------------------------------------------------------------------------------------------------
# Install the Operator network policies from the examples
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-operator-network-policies
install-operator-network-policies: export NAMESPACE := $(OPERATOR_NAMESPACE)
install-operator-network-policies: prepare-network-policies
	$(BUILD_OUTPUT)/network-policies/add-operator-policies.sh

# ----------------------------------------------------------------------------------------------------------------------
# Uninstall the Operator network policies from the examples
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: uninstall-operator-network-policies
uninstall-operator-network-policies: export NAMESPACE := $(OPERATOR_NAMESPACE)
uninstall-operator-network-policies: prepare-network-policies
	$(BUILD_OUTPUT)/network-policies/remove-operator-policies.sh

# ----------------------------------------------------------------------------------------------------------------------
# Install the Coherence cluster network policies from the examples
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: install-coherence-network-policies
install-coherence-network-policies: export NAMESPACE := $(CLUSTER_NAMESPACE)
install-coherence-network-policies: prepare-network-policies
	$(BUILD_OUTPUT)/network-policies/add-cluster-member-policies.sh

# ----------------------------------------------------------------------------------------------------------------------
# Uninstall the Coherence cluster network policies from the examples
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: uninstall-coherence-network-policies
uninstall-coherence-network-policies: export NAMESPACE := $(CLUSTER_NAMESPACE)
uninstall-coherence-network-policies: prepare-network-policies
	$(BUILD_OUTPUT)/network-policies/remove-cluster-member-policies.sh

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
	@echo "Uninstalling CRDs - calling prepare_deploy"
	$(call prepare_deploy,$(OPERATOR_IMAGE),$(OPERATOR_NAMESPACE))
	@echo "Uninstalling CRDs - executing deletion"
	$(KUSTOMIZE) build $(BUILD_DEPLOY)/crd | kubectl delete --force -f - || true
	@echo "Uninstall CRDs completed"

# ----------------------------------------------------------------------------------------------------------------------
# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: deploy-and-wait
deploy-and-wait: deploy wait-for-deploy   ## Deploy the Coherence Operator and wait for the Operator Pod to be ready

# The Operator is deployed HA by default
OPERATOR_HA ?= true

.PHONY: deploy
deploy: prepare-deploy create-namespace $(TOOLS_BIN)/kustomize   ## Deploy the Coherence Operator
ifneq (,$(WATCH_NAMESPACE))
	cd $(BUILD_DEPLOY)/manager && $(KUSTOMIZE) edit add configmap env-vars --from-literal WATCH_NAMESPACE=$(WATCH_NAMESPACE)
endif
ifeq (false,$(OPERATOR_HA))
	cd $(BUILD_DEPLOY)/manager && $(KUSTOMIZE) edit add patch --kind Deployment --name controller-manager --path single-replica-patch.yaml
endif
	kubectl -n $(OPERATOR_NAMESPACE) create secret generic coherence-webhook-server-cert || true
	$(KUSTOMIZE) build $(BUILD_DEPLOY)/default | kubectl apply -f -
	sleep 5

.PHONY: just-deploy
just-deploy: ## Deploy the Coherence Operator without rebuilding anything
	$(call do_deploy,$(OPERATOR_IMAGE),$(OPERATOR_NAMESPACE))

.PHONY: prepare-deploy
prepare-deploy: $(BUILD_TARGETS)/manifests $(BUILD_TARGETS)/build-operator $(TOOLS_BIN)/kustomize
	$(call prepare_deploy,$(OPERATOR_IMAGE),$(OPERATOR_NAMESPACE))

.PHONY: deploy-debug
deploy-debug: prepare-deploy-debug create-namespace $(TOOLS_BIN)/kustomize   ## Deploy the Coherence Operator running with Delve
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
prepare-deploy-debug: $(BUILD_TARGETS)/manifests build-operator-debug $(TOOLS_BIN)/kustomize
	$(call prepare_deploy,$(OPERATOR_IMAGE_DEBUG),$(OPERATOR_NAMESPACE))

.PHONY: wait-for-deploy
wait-for-deploy: export POD=$(shell kubectl -n $(OPERATOR_NAMESPACE) get pod -l control-plane=coherence -o name)
wait-for-deploy:
	sleep 30
	echo "Operator Pods:"
	kubectl -n $(OPERATOR_NAMESPACE) get pod -l control-plane=coherence
	echo "Waiting for Operator to be ready. Pod: $(POD)"
	kubectl -n $(OPERATOR_NAMESPACE) wait --for condition=ready --timeout 480s $(POD)

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
	cd $(BUILD_DEPLOY)/manager && $(KUSTOMIZE) edit add configmap env-vars --from-literal OPERATOR_IMAGE=$(1)
	cd $(BUILD_DEPLOY)/manager && $(KUSTOMIZE) edit set image controller=$(1)
	cd $(BUILD_DEPLOY)/default && $(KUSTOMIZE) edit set namespace $(2)
endef

define do_deploy
	$(call prepare_deploy,$(1),$(2))
	$(KUSTOMIZE) build $(BUILD_DEPLOY)/default | kubectl apply -f -
endef


# ----------------------------------------------------------------------------------------------------------------------
# Un-deploy controller from the configured Kubernetes cluster in ~/.kube/config
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: undeploy
undeploy: $(BUILD_PROPS) $(BUILD_TARGETS)/manifests $(TOOLS_BIN)/kustomize  ## Undeploy the Coherence Operator
	@echo "Undeploy Coherence Operator..."
	$(call prepare_deploy,$(OPERATOR_IMAGE),$(OPERATOR_NAMESPACE))
	$(KUSTOMIZE) build $(BUILD_DEPLOY)/default | kubectl delete -f - || true
	kubectl -n $(OPERATOR_NAMESPACE) delete secret coherence-webhook-server-cert || true
	kubectl delete mutatingwebhookconfiguration coherence-operator-mutating-webhook-configuration || true
	kubectl delete validatingwebhookconfiguration coherence-operator-validating-webhook-configuration || true
	@echo "Undeploy Coherence Operator completed"


# ----------------------------------------------------------------------------------------------------------------------
# Tail the deployed operator logs.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: tail-logs
tail-logs: export POD=$(shell kubectl -n $(OPERATOR_NAMESPACE) get pod -l control-plane=coherence -o name)
tail-logs:     ## Tail the Coherence Operator Pod logs (with follow)
	kubectl -n $(OPERATOR_NAMESPACE) logs $(POD) -c manager -f


$(BUILD_MANIFESTS_PKG): $(TOOLS_BIN)/kustomize $(TOOLS_BIN)/yq
	rm -rf $(BUILD_MANIFESTS) || true
	mkdir -p $(BUILD_MANIFESTS)/crd
	$(KUSTOMIZE) build config/crd > $(BUILD_MANIFESTS)/crd/temp.yaml
	mkdir -p $(BUILD_MANIFESTS)/crd-small
	$(KUSTOMIZE) build config/crd-small > $(BUILD_MANIFESTS)/crd-small/temp.yaml
	cp -R config/default/ $(BUILD_MANIFESTS)/default
	cp -R config/manager/ $(BUILD_MANIFESTS)/manager
	cp -R config/rbac/ $(BUILD_MANIFESTS)/rbac
	tar -C $(BUILD_OUTPUT) -czf $(BUILD_MANIFESTS_PKG) manifests/
	$(call prepare_deploy,$(OPERATOR_IMAGE),"coherence")
	cp config/namespace/namespace.yaml $(BUILD_OUTPUT)/coherence-operator.yaml
	$(KUSTOMIZE) build $(BUILD_DEPLOY)/default >> $(BUILD_OUTPUT)/coherence-operator.yaml
	$(SED) -e 's/name: coherence-operator-env-vars-.*/name: coherence-operator-env-vars/g' $(BUILD_OUTPUT)/coherence-operator.yaml
	cd $(BUILD_MANIFESTS)/crd && $(TOOLS_BIN)/yq --no-doc -s '.metadata.name + ".yaml"' temp.yaml
	rm $(BUILD_MANIFESTS)/crd/temp.yaml
	mv $(BUILD_MANIFESTS)/crd/coherence.coherence.oracle.com.yaml $(BUILD_MANIFESTS)/crd/coherence.oracle.com_coherence.yaml
	mv $(BUILD_MANIFESTS)/crd/coherencejob.coherence.oracle.com.yaml $(BUILD_MANIFESTS)/crd/coherencejob.oracle.com_coherence.yaml
	cd $(BUILD_MANIFESTS)/crd-small && $(TOOLS_BIN)/yq --no-doc -s '.metadata.name + ".yaml"' temp.yaml
	rm $(BUILD_MANIFESTS)/crd-small/temp.yaml
	mv $(BUILD_MANIFESTS)/crd-small/coherence.coherence.oracle.com.yaml $(BUILD_MANIFESTS)/crd-small/coherence.oracle.com_coherence.yaml
	mv $(BUILD_MANIFESTS)/crd-small/coherencejob.coherence.oracle.com.yaml $(BUILD_MANIFESTS)/crd-small/coherencejob.oracle.com_coherence.yaml

# ----------------------------------------------------------------------------------------------------------------------
# Delete and re-create the test namespace
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: create-namespace
create-namespace: export KUBECONFIG_PATH := $(KUBECONFIG_PATH)
create-namespace: ## Create the test namespace
ifeq ($(CREATE_OPERATOR_NAMESPACE),true)
	kubectl get ns $(OPERATOR_NAMESPACE) -o name > /dev/null 2>&1 || kubectl create namespace $(OPERATOR_NAMESPACE)
	kubectl get ns $(OPERATOR_NAMESPACE_CLIENT) -o name > /dev/null 2>&1 || kubectl create namespace $(OPERATOR_NAMESPACE_CLIENT)
	kubectl get ns $(CLUSTER_NAMESPACE) -o name > /dev/null 2>&1 || kubectl create namespace $(CLUSTER_NAMESPACE)
endif
	kubectl label namespace $(OPERATOR_NAMESPACE) coherence.oracle.com/test=true --overwrite
	kubectl label namespace $(OPERATOR_NAMESPACE_CLIENT) coherence.oracle.com/test=true --overwrite
	kubectl label namespace $(CLUSTER_NAMESPACE) coherence.oracle.com/test=true --overwrite

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
	$(call delete_ns,$(OPERATOR_NAMESPACE))
	$(call delete_ns,$(OPERATOR_NAMESPACE_CLIENT))
	$(call delete_ns,$(CLUSTER_NAMESPACE))
endif
	kubectl delete clusterrole operator-test-coherence-operator --force --ignore-not-found=true --grace-period=0 && echo "deleted namespace" || true
	kubectl delete clusterrolebinding operator-test-coherence-operator --ignore-not-found=true --force --grace-period=0 && echo "deleted namespace" || true

define delete_ns
	if kubectl get ns $(1); then \
		echo "Deleting test namespace $(1)" ;\
		kubectl delete namespace $(1) --force --ignore-not-found=true --grace-period=0 --timeout=600s ;\
		echo "deleted namespace $(1)" || true ;\
	fi
endef

# ----------------------------------------------------------------------------------------------------------------------
# Delete all of the Coherence resources from the test namespace.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: delete-coherence-clusters
delete-coherence-clusters: ## Delete all running Coherence clusters in the test namespace
	for i in $$(kubectl -n  $(OPERATOR_NAMESPACE) get coherencejob.coherence.oracle.com -o name); do \
  		kubectl -n $(OPERATOR_NAMESPACE) patch $${i} -p '{"metadata":{"finalizers":[]}}' --type=merge || true ;\
		kubectl -n $(OPERATOR_NAMESPACE) delete $${i}; \
	done
	for i in $$(kubectl -n  $(CLUSTER_NAMESPACE) get coherencejob.coherence.oracle.com -o name); do \
  		kubectl -n $(CLUSTER_NAMESPACE) patch $${i} -p '{"metadata":{"finalizers":[]}}' --type=merge || true ;\
		kubectl -n $(CLUSTER_NAMESPACE) delete $${i}; \
	done
	for i in $$(kubectl -n  $(OPERATOR_NAMESPACE) get coherence.coherence.oracle.com -o name); do \
  		kubectl -n $(OPERATOR_NAMESPACE) patch $${i} -p '{"metadata":{"finalizers":[]}}' --type=merge || true ;\
		kubectl -n $(OPERATOR_NAMESPACE) delete $${i}; \
	done
	for i in $$(kubectl -n  $(CLUSTER_NAMESPACE) get coherence.coherence.oracle.com -o name); do \
  		kubectl -n $(CLUSTER_NAMESPACE) patch $${i} -p '{"metadata":{"finalizers":[]}}' --type=merge || true ;\
		kubectl -n $(CLUSTER_NAMESPACE) delete $${i}; \
	done

# ----------------------------------------------------------------------------------------------------------------------
# Delete all resource from the test namespace
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: clean-namespace
clean-namespace: delete-coherence-clusters   ## Clean-up deployments in the test namespace
	@echo "Cleaning Namespaces..."
	kubectl delete --all networkpolicy --namespace=$(OPERATOR_NAMESPACE) || true
	kubectl delete --all networkpolicy --namespace=$(CLUSTER_NAMESPACE) || true
	for i in $$(kubectl -n $(OPERATOR_NAMESPACE) get all -o name); do \
		echo "Deleting $${i} from test namespace $(OPERATOR_NAMESPACE)" \
		kubectl -n $(OPERATOR_NAMESPACE) delete $${i}; \
	done
	for i in $$(kubectl -n $(CLUSTER_NAMESPACE) get all -o name); do \
		echo "Deleting $${i} from test namespace $(CLUSTER_NAMESPACE)" \
		kubectl -n $(CLUSTER_NAMESPACE) delete $${i}; \
	done
	@echo "Cleaning Namespaces completed"

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

KIND_CLUSTER   ?= operator
KIND_IMAGE     ?= "kindest/node:v1.32.0@sha256:c48c62eac5da28cdadcf560d1d8616cfa6783b58f0d94cf63ad1bf49600cb027"
CALICO_TIMEOUT ?= 300s

# ----------------------------------------------------------------------------------------------------------------------
# Start a Kind cluster
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: kind
kind:   ## Run a default KinD cluster
	kind create cluster --name $(KIND_CLUSTER) --wait 10m --config $(SCRIPTS_DIR)/kind-config.yaml --image $(KIND_IMAGE)
	$(SCRIPTS_DIR)/kind-label-node.sh

.PHONY: kind-dual
kind-dual:   ## Run a KinD cluster configured for a dual stack IPv4 and IPv6 network
	kind create cluster --name $(KIND_CLUSTER) --wait 10m --config $(SCRIPTS_DIR)/kind-config-dual.yaml --image $(KIND_IMAGE)
	$(SCRIPTS_DIR)/kind-label-node.sh

# ----------------------------------------------------------------------------------------------------------------------
# Start a Kind cluster
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: kind-single-worker
kind-single-worker:   ## Run a KinD cluster with a single worker node
	kind create cluster --name $(KIND_CLUSTER) --wait 10m --config $(SCRIPTS_DIR)/kind-config-single.yaml --image $(KIND_IMAGE)
	$(SCRIPTS_DIR)/kind-label-node.sh

# ----------------------------------------------------------------------------------------------------------------------
# Start a Kind cluster with Calico
# ----------------------------------------------------------------------------------------------------------------------
CALICO_VERSION ?= v3.25.0

.PHONY: kind-calico
kind-calico: export KIND_CONFIG=$(SCRIPTS_DIR)/kind-config-calico.yaml
kind-calico:   ## Run a KinD cluster with Calico
	kind create cluster --name $(KIND_CLUSTER) --config $(SCRIPTS_DIR)/kind-config-calico.yaml --image $(KIND_IMAGE)
	$(SCRIPTS_DIR)/kind-label-node.sh
	kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/$(CALICO_VERSION)/manifests/calico.yaml
	kubectl -n kube-system set env daemonset/calico-node FELIX_IGNORELOOSERPF=true
	sleep 30
	kubectl -n kube-system wait --for condition=ready --timeout=$(CALICO_TIMEOUT) -l k8s-app=calico-node pod
	kubectl -n kube-system wait --for condition=ready --timeout=$(CALICO_TIMEOUT) -l k8s-app=kube-dns pod

# ----------------------------------------------------------------------------------------------------------------------
# Stop and delete the Kind cluster
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: kind-stop
kind-stop:   ## Stop and delete the KinD cluster named "$(KIND_CLUSTER)"
	kind delete cluster --name $(KIND_CLUSTER)

# ----------------------------------------------------------------------------------------------------------------------
# Load images into Kind
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: kind-load
kind-load: kind-load-operator kind-load-coherence  ## Load all images into the KinD cluster
	kind load docker-image --name $(KIND_CLUSTER) $(TEST_APPLICATION_IMAGE) || true
	kind load docker-image --name $(KIND_CLUSTER) $(TEST_APPLICATION_IMAGE_CLIENT) || true
	kind load docker-image --name $(KIND_CLUSTER) $(TEST_APPLICATION_IMAGE_HELIDON) || true
	kind load docker-image --name $(KIND_CLUSTER) $(TEST_APPLICATION_IMAGE_SPRING) || true
	kind load docker-image --name $(KIND_CLUSTER) $(TEST_APPLICATION_IMAGE_SPRING_FAT) || true
	kind load docker-image --name $(KIND_CLUSTER) $(TEST_APPLICATION_IMAGE_SPRING_CNBP) || true

.PHONY: kind-load-coherence
kind-load-coherence:   ## Load the Coherence image into the KinD cluster
	docker pull $(COHERENCE_IMAGE)
	kind load docker-image --name $(KIND_CLUSTER) $(COHERENCE_IMAGE)

.PHONY: kind-load-operator
kind-load-operator:   ## Load the Operator images into the KinD cluster
	kind load docker-image --name $(KIND_CLUSTER) $(OPERATOR_IMAGE) || true

# ----------------------------------------------------------------------------------------------------------------------
# Load compatibility images into Kind
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: kind-load-compatibility
kind-load-compatibility:   ## Load the compatibility test images into the KinD cluster
	kind load docker-image --name $(KIND_CLUSTER) $(TEST_COMPATIBILITY_IMAGE) || true

# ======================================================================================================================
# Targets related to running k3d clusters
# ======================================================================================================================
##@ K3d

K3D_CLUSTER           ?= operator
K3D_REGISTRY          ?= myregistry
K3D_REGISTRY_PORT     ?= 12345
K3D_INTERNAL_REGISTRY := k3d-$(K3D_REGISTRY).localhost:$(K3D_REGISTRY_PORT)

.PHONY: k3d
k3d: $(TOOLS_BIN)/k3d k3d-create k3d-load-operator create-namespace  ## Run a default k3d cluster

.PHONY: k3d-create
k3d-create: $(TOOLS_BIN)/k3d ## Create the k3d cluster
	$(TOOLS_BIN)/k3d registry create myregistry.localhost --port 12345
	$(TOOLS_BIN)/k3d cluster create $(K3D_CLUSTER) --agents 5 \
		--registry-use $(K3D_INTERNAL_REGISTRY) --no-lb \
		--runtime-ulimit "nofile=64000:64000" --runtime-ulimit "nproc=64000:64000" \
		--api-port 127.0.0.1:6550

.PHONY: k3d-stop
k3d-stop: $(TOOLS_BIN)/k3d  ## Stop a default k3d cluster
	$(TOOLS_BIN)/k3d cluster delete $(K3D_CLUSTER)
	$(TOOLS_BIN)/k3d registry delete myregistry.localhost

.PHONY: k3d-load-operator
k3d-load-operator: $(TOOLS_BIN)/k3d  ## Load the Operator images into the k3d cluster
	$(TOOLS_BIN)/k3d image import $(OPERATOR_IMAGE) -c $(K3D_CLUSTER)

.PHONY: k3d-load-coherence
k3d-load-coherence: $(TOOLS_BIN)/k3d  ## Load the Coherence images into the k3d cluster
	$(TOOLS_BIN)/k3d image import $(COHERENCE_IMAGE) -c $(K3D_CLUSTER)

.PHONY: k3d-load-all
k3d-load-all: $(TOOLS_BIN)/k3d k3d-load-operator k3d-load-coherence ## Load all the test images into the k3d cluster

.PHONY: k3d-get
k3d-get: $(TOOLS_BIN)/k3d ## Install k3d

K3D_PATH = ${PATH}
$(TOOLS_BIN)/k3d:
	export K3D_INSTALL_DIR=$(TOOLS_BIN) \
		&& export USE_SUDO=false \
		&& export PATH="$(TOOLS_BIN):$(K3D_PATH)" \
		&& curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash

# ======================================================================================================================
# Targets related to running Minikube
# ======================================================================================================================
##@ Minikube

# the version of minikube to install
MINIKUBE_VERSION ?= latest
MINIKUBE_K8S     ?= 1.25.8

# ----------------------------------------------------------------------------------------------------------------------
# Start Minikube
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: minikube
minikube: minikube-install  ## Run a default minikube cluster with Calico
	$(MINIKUBE) start --driver docker --cni calico --kubernetes-version $(MINIKUBE_K8S)
	$(MINIKUBE) status
	kubectl get nodes

# ----------------------------------------------------------------------------------------------------------------------
# Stop Minikube
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: minikube-stop
minikube-stop:  ## Stop and delete the minikube cluster
	$(MINIKUBE) stop || true
	$(MINIKUBE) delete || true

# ----------------------------------------------------------------------------------------------------------------------
# Install Minikube
# ----------------------------------------------------------------------------------------------------------------------
MINIKUBE = $(TOOLS_BIN)/minikube
.PHONY: minikube-install
minikube-install: $(TOOLS_BIN)/minikube ## Install minikube (defaults to the latest version, can be changed by setting MINIKUBE_VERSION)
	$(MINIKUBE) version

$(TOOLS_BIN)/minikube:
ifeq (Darwin, $(UNAME_S))
ifeq (x86_64, $(UNAME_M))
	curl -LOs https://storage.googleapis.com/minikube/releases/$(MINIKUBE_VERSION)/minikube-darwin-amd64
	mkdir -p $(TOOLS_BIN) || true
	install minikube-darwin-amd64 $(TOOLS_BIN)/minikube
	rm minikube-darwin-amd64
else
	curl -LOs https://storage.googleapis.com/minikube/releases/$(MINIKUBE_VERSION)/minikube-darwin-arm64
	mkdir -p $(TOOLS_BIN) || true
	install minikube-darwin-arm64 $(TOOLS_BIN)/minikube
	rm minikube-darwin-arm64
endif
else
	curl -LOs https://storage.googleapis.com/minikube/releases/$(MINIKUBE_VERSION)/minikube-linux-amd64
	mkdir -p $(TOOLS_BIN) || true
	install minikube-linux-amd64 $(TOOLS_BIN)/minikube
	rm minikube-linux-amd64
endif

# ----------------------------------------------------------------------------------------------------------------------
# Install yq
# ----------------------------------------------------------------------------------------------------------------------
YQ         = $(TOOLS_BIN)/yq
YQ_VERSION = v4.44.3

.PHONY: yq-install
yq-install: $(TOOLS_BIN)/yq  ## Install yq (defaults to the latest version, can be changed by setting YQ_VERSION)
	$(YQ) version

$(TOOLS_BIN)/yq:
	mkdir -p $(TOOLS_BIN) || true
ifeq (Darwin, $(UNAME_S))
ifeq (x86_64, $(UNAME_M))
	curl -L https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_darwin_amd64 -o $(TOOLS_BIN)/yq
else
	curl -L https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_darwin_arm64 -o $(TOOLS_BIN)/yq
endif
else
	curl -L https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_linux_amd64 -o $(TOOLS_BIN)/yq
endif
	chmod +x $(TOOLS_BIN)/yq

# ======================================================================================================================
# Kubernetes Cert Manager targets
# ======================================================================================================================
##@ Cert Manager

CERT_MANAGER_VERSION ?= v1.8.0
# Get latest version...
#  curl -s -H "Accept: application/vnd.github.v3+json" --header $(GH_AUTH) https://api.github.com/repos/cert-manager/cert-manager/releases | jq '.[0].tag_name' |  tr -d '"'

.PHONY: install-cmctl
install-cmctl: $(TOOLS_BIN)/cmctl ## Install the Cert Manager CLI into $(TOOLS_BIN)

CMCTL = $(TOOLS_BIN)/cmctl
$(TOOLS_BIN)/cmctl:
	OS=$(shell go env GOOS) && ARCH=$(shell go env GOARCH) && \
	curl -sSL -o cmctl.tar.gz https://github.com/cert-manager/cert-manager/releases/download/$(CERT_MANAGER_VERSION)/cmctl-$${OS}-$${ARCH}.tar.gz  --header $(GH_AUTH)
	tar xzf cmctl.tar.gz
	mv cmctl $(TOOLS_BIN)
	rm cmctl.tar.gz

.PHONY: install-cert-manager
install-cert-manager: $(TOOLS_BIN)/cmctl ## Install Cert manager into the Kubernetes cluster
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.8.0/cert-manager.yam
	$(CMCTL) check api --wait=10m

.PHONY: uninstall-cert-manager
uninstall-cert-manager: ## Uninstall Cert manager from the Kubernetes cluster
	kubectl delete -f https://github.com/cert-manager/cert-manager/releases/download/$(CERT_MANAGER_VERSION)/cert-manager.yam


# ======================================================================================================================
# Tanzu related targets
# ======================================================================================================================
##@ Tanzu

TANZU = $(shell which tanzu)
.PHONY: get-tanzu
get-tanzu: $(BUILD_PROPS)
	$(SCRIPTS_DIR)/get-tanzu.sh "$(TANZU_VERSION)" "$(TOOLS_DIRECTORY)"

.PHONY: tanzu-create-cluster
tanzu-create-cluster: ## Create a local Tanzu unmanaged cluster named "$(KIND_CLUSTER)" (default "operator")
	rm -rf $(HOME)/.config/tanzu/tkg/unmanaged/$(KIND_CLUSTER)
	$(TANZU) uc create $(KIND_CLUSTER) --worker-node-count 2

.PHONY: tanzu-delete-cluster
tanzu-delete-cluster: ## Delete the local Tanzu unmanaged cluster named "$(KIND_CLUSTER)" (default "operator")
	$(TANZU) uc delete $(KIND_CLUSTER)

.PHONY: tanzu-package-internal
tanzu-package-internal: $(BUILD_PROPS) $(BUILD_TARGETS)/generate $(BUILD_TARGETS)/manifests $(TOOLS_BIN)/kustomize
	rm -r $(TANZU_PACKAGE_DIR) || true
	mkdir -p $(TANZU_PACKAGE_DIR)/config $(TANZU_PACKAGE_DIR)/.imgpkg || true
	cp -vR tanzu/package/* $(TANZU_PACKAGE_DIR)/config/
	ls -al $(TANZU_PACKAGE_DIR)/config/
	$(call prepare_deploy,$(OPERATOR_IMAGE),tanzu-namespace)

.PHONY: tanzu-package
tanzu-package: tanzu-package-internal ## Create the Tanzu package files.
	$(KUSTOMIZE) build $(BUILD_DEPLOY)/default >> $(TANZU_PACKAGE_DIR)/config/package.yml
	$(SED) -e 's/tanzu-namespace/#@ data.values.namespace/g' $(TANZU_PACKAGE_DIR)/config/package.yml
	$(call pushTanzuPackage,$(OPERATOR_PACKAGE_IMAGE))

.PHONY: tanzu-ttl-package
tanzu-ttl-package: tanzu-package-internal ## Create the Tanzu package files using images from ttl.sh
	$(KUSTOMIZE) build $(BUILD_DEPLOY)/default >> $(TANZU_PACKAGE_DIR)/config/package.yml
	$(SED) -e 's/tanzu-namespace/#@ data.values.namespace/g' $(TANZU_PACKAGE_DIR)/config/package.yml
	$(SED) -e 's,$(OPERATOR_IMAGE),$(TTL_OPERATOR_IMAGE),g' $(TANZU_PACKAGE_DIR)/config/package.yml
	$(call pushTanzuPackage,$(TTL_PACKAGE_IMAGE))

define pushTanzuPackage
	kbld -f $(TANZU_PACKAGE_DIR)/config/ --imgpkg-lock-output $(TANZU_PACKAGE_DIR)/.imgpkg/images.yml
	tar -czvf $(TANZU_DIR)/tanzu-package.tar.gz  $(TANZU_PACKAGE_DIR)/
	imgpkg push -b $(1) -f $(TANZU_PACKAGE_DIR)/
endef

.PHONY: tanzu-repo-internal
tanzu-repo-internal:
	rm -r $(TANZU_REPO_DIR) || true
	mkdir -p $(TANZU_REPO_DIR)/.imgpkg $(TANZU_REPO_DIR)/packages/coherence-operator.oracle.github.com
	cp ./tanzu/repo/metadata.yaml $(TANZU_REPO_DIR)/packages/coherence-operator.oracle.github.com/metadata.yaml
	cp ./tanzu/repo/version.yaml $(TANZU_REPO_DIR)/packages/coherence-operator.oracle.github.com/v$(VERSION).yaml
	$(call replaceprop,$(TANZU_REPO_DIR)/packages/coherence-operator.oracle.github.com/v$(VERSION).yaml)

.PHONY: tanzu-repo
tanzu-repo: tanzu-package tanzu-repo-internal ## Create the Tanzu repo files
	$(call pushTanzuRepo,$(OPERATOR_REPO_IMAGE))

.PHONY: tanzu-ttl-repo
tanzu-ttl-repo: tanzu-ttl-package tanzu-repo-internal ## Create the Tanzu repo files using images from ttl.sh
	$(SED) -e 's,$(OPERATOR_PACKAGE_IMAGE),$(TTL_PACKAGE_IMAGE),g' $(TANZU_REPO_DIR)/packages/coherence-operator.oracle.github.com/v$(VERSION).yaml
	$(call pushTanzuRepo,$(TTL_REPO_IMAGE))

define pushTanzuRepo
	kbld -f $(TANZU_REPO_DIR)/packages/ --imgpkg-lock-output $(TANZU_REPO_DIR)/.imgpkg/images.yml
	tar -czvf $(TANZU_DIR)/tanzu-repo.tar.gz  $(TANZU_REPO_DIR)/
	imgpkg push -b $(1) -f $(TANZU_REPO_DIR)/
endef

.PHONY: tanzu-install-repo
tanzu-install-repo: ## Install the Coherence package repo into Tanzu
	$(call tanzuInstallRepo,$(OPERATOR_REPO_IMAGE))

.PHONY: tanzu-ttl-install-repo
tanzu-ttl-install-repo: ## Install the Coherence package repo into Tanzu using images from ttl.sh
	$(call tanzuInstallRepo,$(TTL_REPO_IMAGE))

.PHONY: tanzu-delete-repo
tanzu-delete-repo: ## Delete the Coherence package repo into Tanzu
	$(TANZU) package repository delete coherence-repo -y --namespace coherence

define tanzuInstallRepo
	$(TANZU) package repository add coherence-repo \
		--url $(1) \
		--namespace coherence \
		--create-namespace
	$(TANZU) package repository list --namespace coherence
	$(TANZU) package available list --namespace coherence
endef

.PHONY: tanzu-install
tanzu-install: ## Install the Coherence Operator package into Tanzu
	$(TANZU) package install coherence \
		--package-name coherence-operator.oracle.github.com \
		--version $(VERSION) \
		--namespace coherence
	$(TANZU) package installed list --namespace coherence

# ======================================================================================================================
# Miscellaneous targets
# ======================================================================================================================
##@ Miscellaneous

TRIVY_CACHE ?=

.PHONY: trivy-scan
trivy-scan: build-operator-images $(TOOLS_BIN)/trivy ## Scan the Operator image using Trivy
ifeq (Darwin, $(UNAME_S))
	$(TOOLS_BIN)/trivy --exit-code 1 --severity CRITICAL,HIGH --cache-dir $(HOME)/Library/Caches/trivy image $(OPERATOR_IMAGE)
else
ifdef TRIVY_CACHE
	$(TOOLS_BIN)/trivy --exit-code 1 --severity CRITICAL,HIGH --cache-dir $(TRIVY_CACHE) image $(OPERATOR_IMAGE)
else
	$(TOOLS_BIN)/trivy --exit-code 1 --severity CRITICAL,HIGH image $(OPERATOR_IMAGE)
endif
endif

.PHONY: get-trivy
get-trivy: $(TOOLS_BIN)/trivy

$(TOOLS_BIN)/trivy:
	test -s $(TOOLS_BIN)/trivy || curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b $(TOOLS_BIN) v0.56.2

# ----------------------------------------------------------------------------------------------------------------------
# find or download controller-gen
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: controller-gen
CONTROLLER_GEN = $(TOOLS_BIN)/controller-gen
controller-gen: $(TOOLS_BIN)/controller-gen ## Download controller-gen locally if necessary.

$(TOOLS_BIN)/controller-gen:
	@echo "Downloading controller-gen"
	test -s $(TOOLS_BIN)/controller-gen || GOBIN=$(TOOLS_BIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.17.0
	ls -al $(TOOLS_BIN)

# ----------------------------------------------------------------------------------------------------------------------
# find or download kustomize
# ----------------------------------------------------------------------------------------------------------------------
KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
KUSTOMIZE_VERSION ?= v5.4.2

.PHONY: kustomize
KUSTOMIZE = $(TOOLS_BIN)/kustomize
kustomize: $(TOOLS_BIN)/kustomize ## Download kustomize locally if necessary.

$(TOOLS_BIN)/kustomize:
	test -s $(TOOLS_BIN)/kustomize || { curl -Ss $(KUSTOMIZE_INSTALL_SCRIPT) --header $(GH_AUTH) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(TOOLS_BIN); }

# ----------------------------------------------------------------------------------------------------------------------
# find or download the Coherence CLI
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: coherence-cli
coherence-cli: $(BUILD_TARGETS)/cli ## Download the Coherence CLI locally if necessary.

$(BUILD_TARGETS)/cli: $(BUILD_BIN_AMD64)/cohctl $(BUILD_BIN_ARM64)/cohctl
	touch $(BUILD_TARGETS)/cli

$(BUILD_BIN_AMD64)/cohctl: export COHCTL_HOME=$(BUILD_BIN_AMD64)
$(BUILD_BIN_AMD64)/cohctl: export OS=Linux
$(BUILD_BIN_AMD64)/cohctl: export ARCH=x86_64
$(BUILD_BIN_AMD64)/cohctl:
	./hack/install-cli.sh
	chmod +x $(BUILD_BIN_AMD64)/cohctl

$(BUILD_BIN_ARM64)/cohctl: export COHCTL_HOME=$(BUILD_BIN_ARM64)
$(BUILD_BIN_ARM64)/cohctl: export OS=Linux
$(BUILD_BIN_ARM64)/cohctl: export ARCH=arm64
$(BUILD_BIN_ARM64)/cohctl:
	./hack/install-cli.sh
	chmod +x $(BUILD_BIN_ARM64)/cohctl

# ----------------------------------------------------------------------------------------------------------------------
# find or download gotestsum
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: gotestsum
GOTESTSUM = $(TOOLS_BIN)/gotestsum
gotestsum: ## Download gotestsum locally if necessary.
	test -s $(TOOLS_BIN)/gotestsum || GOBIN=$(TOOLS_BIN) go install gotest.tools/gotestsum@v1.8.2

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
# Push the Operator Test images to ttl.sh
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-ttl-test-images
push-ttl-test-images:
	docker tag $(TEST_APPLICATION_IMAGE) $(TTL_APPLICATION_IMAGE)
	docker push $(TTL_APPLICATION_IMAGE)
	docker tag $(TEST_APPLICATION_IMAGE_CLIENT) $(TTL_APPLICATION_IMAGE_CLIENT)
	docker push $(TTL_APPLICATION_IMAGE_CLIENT)
	docker tag $(TEST_APPLICATION_IMAGE_HELIDON) $(TTL_APPLICATION_IMAGE_HELIDON)
	docker push $(TTL_APPLICATION_IMAGE_HELIDON)
	docker tag $(TEST_APPLICATION_IMAGE_SPRING) $(TTL_APPLICATION_IMAGE_SPRING)
	docker push $(TTL_APPLICATION_IMAGE_SPRING)
	docker tag $(TEST_APPLICATION_IMAGE_SPRING_FAT) $(TTL_APPLICATION_IMAGE_SPRING_FAT)
	docker push $(TTL_APPLICATION_IMAGE_SPRING_FAT)
	docker tag $(TEST_APPLICATION_IMAGE_SPRING_CNBP) $(TTL_APPLICATION_IMAGE_SPRING_CNBP)
	docker push $(TTL_APPLICATION_IMAGE_SPRING_CNBP)

# ----------------------------------------------------------------------------------------------------------------------
# Build the Operator Test images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: build-compatibility-image
build-compatibility-image: $(BUILD_TARGETS)/java
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
push-compatibility-image: build-compatibility-image
	docker push $(TEST_COMPATIBILITY_IMAGE)

# ----------------------------------------------------------------------------------------------------------------------
# Push the Operator JIB Test Docker images to ttl.sh
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-ttl-compatibility-image
push-ttl-compatibility-image:
	docker tag $(TEST_COMPATIBILITY_IMAGE) $(TTL_COMPATIBILITY_IMAGE)
	docker push $(TTL_COMPATIBILITY_IMAGE)

# ----------------------------------------------------------------------------------------------------------------------
# Push all of the Docker images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-all-images
push-all-images: push-test-images push-operator-image

# ----------------------------------------------------------------------------------------------------------------------
# Push the Operator images to ttl.sh
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-ttl-operator-images
push-ttl-operator-images:
	docker tag $(OPERATOR_IMAGE) $(TTL_OPERATOR_IMAGE)
	docker push $(TTL_OPERATOR_IMAGE)

# ----------------------------------------------------------------------------------------------------------------------
# Push all the images to ttl.sh
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-all-ttl-images
push-all-ttl-images:  push-ttl-operator-images push-ttl-test-images

# ----------------------------------------------------------------------------------------------------------------------
# Push all of the Docker images that are released
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-release-images
push-release-images: push-operator-image tanzu-repo

# ----------------------------------------------------------------------------------------------------------------------
# Install Prometheus
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: get-prometheus
get-prometheus: $(PROMETHEUS_HOME)/$(PROMETHEUS_VERSION).txt ## Download Prometheus Operator kube-prometheus

$(PROMETHEUS_HOME)/$(PROMETHEUS_VERSION).txt: $(BUILD_PROPS)
ifeq (main, $(PROMETHEUS_VERSION))
	curl -sL  https://github.com/prometheus-operator/kube-prometheus/archive/main.tar.gz -o $(BUILD_OUTPUT)/prometheus.tar.gz  --header $(GH_AUTH)
else
	curl -sL https://github.com/prometheus-operator/kube-prometheus/archive/refs/tags/$(PROMETHEUS_VERSION).tar.gz -o $(BUILD_OUTPUT)/prometheus.tar.gz  --header $(GH_AUTH)
endif
	mkdir -p $(PROMETHEUS_HOME)
	tar -zxf $(BUILD_OUTPUT)/prometheus.tar.gz --directory $(PROMETHEUS_HOME) --strip-components=1
	rm $(BUILD_OUTPUT)/prometheus.tar.gz
	touch $(PROMETHEUS_HOME)/$(PROMETHEUS_VERSION).txt

.PHONY: install-prometheus
install-prometheus: get-prometheus ## Install Prometheus and Grafana
	kubectl create -f $(PROMETHEUS_HOME)/manifests/setup
	sleep 10
	until kubectl get servicemonitors --all-namespaces ; do date; sleep 1; echo ""; done
#   We create additional custom RBAC rules because the defaults do not work
#   in an RBAC enabled cluster such as KinD
#   See: https://prometheus-operator.dev/docs/operator/rbac/
	kubectl create -f hack/prometheus-rbac.yaml
	kubectl create -f $(PROMETHEUS_HOME)/manifests
	sleep 10
	kubectl -n monitoring get all
	@echo "Waiting for Prometheus StatefulSet to be ready"
	until kubectl -n monitoring get statefulset/prometheus-k8s ; do date; sleep 1; echo ""; done
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
install-istio: get-istio ## Install the latest version of Istio into k8s (or override the version using the ISTIO_VERSION env var)
	$(eval ISTIO_HOME := $(shell find $(TOOLS_DIRECTORY) -maxdepth 1 -type d | grep istio))
	$(ISTIO_HOME)/bin/istioctl install --set profile=demo -y
	kubectl -n istio-system wait --for condition=available deployment.apps/istiod
	kubectl -n istio-system wait --for condition=available deployment.apps/istio-ingressgateway
	kubectl -n istio-system wait --for condition=available deployment.apps/istio-egressgateway
	kubectl apply -f ./hack/istio-strict.yaml
	kubectl -n $(OPERATOR_NAMESPACE) apply -f ./hack/istio-operator.yaml
	kubectl label namespace $(OPERATOR_NAMESPACE) istio-injection=enabled --overwrite=true
	kubectl label namespace $(OPERATOR_NAMESPACE_CLIENT) istio-injection=enabled --overwrite=true
	kubectl label namespace $(CLUSTER_NAMESPACE) istio-injection=enabled --overwrite=true
	kubectl apply -f $(ISTIO_HOME)/samples/addons

# ----------------------------------------------------------------------------------------------------------------------
# Uninstall Istio
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: uninstall-istio
uninstall-istio: get-istio ## Uninstall Istio from k8s
	kubectl -n $(OPERATOR_NAMESPACE) delete -f ./hack/istio-operator.yaml || true
	kubectl delete -f ./hack/istio-strict.yaml
	$(eval ISTIO_HOME := $(shell find $(TOOLS_DIRECTORY) -maxdepth 1 -type d | grep istio))
	$(ISTIO_HOME)/bin/istioctl uninstall --purge -y


# ----------------------------------------------------------------------------------------------------------------------
# Get the latest Istio version
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: get-istio
get-istio: $(BUILD_PROPS)
	$(SCRIPTS_DIR)/get-istio-latest.sh "$(ISTIO_VERSION)" "$(TOOLS_DIRECTORY)"
	$(eval ISTIO_HOME := $(shell find $(TOOLS_DIRECTORY) -maxdepth 1 -type d | grep istio))
	@echo "Istio installed at $(ISTIO_HOME)"

# ----------------------------------------------------------------------------------------------------------------------
# Obtain the golangci-lint binary
# ----------------------------------------------------------------------------------------------------------------------
$(TOOLS_BIN)/golangci-lint:
	@mkdir -p $(TOOLS_BIN)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh --header $(GH_AUTH) | sh -s -- -b $(TOOLS_BIN) v1.63.1

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
docs: api-doc-gen
	./mvnw -B -f java install -P docs -pl docs -DskipTests \
		-Doperator.version=$(VERSION) \
		-Doperator.image=$(OPERATOR_IMAGE) \
		-Dcoherence.image=$(COHERENCE_IMAGE) \
		-Dk8s-doc-version=$(KUBERNETES_DOC_VERSION) \
		$(MAVEN_OPTIONS)
	mkdir -p $(BUILD_OUTPUT)/docs/images/images
	cp -R docs/images/* build/_output/docs/images/
	find examples/ -name \*.png -exec cp {} build/_output/docs/images/images/ \;

# ----------------------------------------------------------------------------------------------------------------------
# Test the documentation.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: test-docs
test-docs: docs
	go run ./utils/linkcheck/ --file $(BUILD_OUTPUT)/docs/pages/... \
		--exclude 'https://oracle.github.io/coherence-operator/charts' \
		--exclude 'https://github.com/oracle/coherence-operator/releases' \
		--exclude 'https://oracle.github.io/coherence-operator/docs/latest/' \
		--exclude 'http://proxyserver' \
		--exclude 'http://&lt;pod-ip' \
		--exclude 'http://elasticsearch-master' \
		--exclude 'http://payments' \
 		2>&1 | tee $(TEST_LOGS_DIR)/doc-link-check.log

# ----------------------------------------------------------------------------------------------------------------------
# Start a local web server to serve the documentation.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: serve-docs
serve-docs:
	@echo "Serving documentation on http://localhost:8080"
	cd $(BUILD_OUTPUT)/docs; \
	python -m SimpleHTTPServer 8080

# ======================================================================================================================
# Release targets
# ======================================================================================================================
##@ Release Targets

# ----------------------------------------------------------------------------------------------------------------------
# Pre-Release Tasks
# Update the version numbers in the documentation to be the version about to be released
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: pre-release
pre-release:
	$(SED) 's/$(subst .,\.,$(PREV_VERSION))/$(VERSION)/g' README.md
	find docs \( -name '*.adoc' -o -name '*.md' \) -exec $(SED) 's/$(subst .,\.,$(PREV_VERSION))/$(VERSION)/g' {} +
	find examples \( -name '*.adoc' -o -name '*.md' -o -name '*.yaml' \) -exec $(SED) 's/$(subst .,\.,$(PREV_VERSION))/$(VERSION)/g' {} +

# ----------------------------------------------------------------------------------------------------------------------
# Post-Release Tasks
# Update the version numbers
#post-release: check-new-version new-version manifests generate build-all-images
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: post-release
post-release: check-new-version new-version

# ----------------------------------------------------------------------------------------------------------------------
# Release the Coherence Operator dashboards
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: release-dashboards
release-dashboards:
	@echo "Releasing Dashboards $(VERSION)"
	mkdir -p $(BUILD_OUTPUT)/dashboards/$(VERSION) || true
	tar -czvf $(BUILD_OUTPUT)/dashboards/$(VERSION)/coherence-dashboards.tar.gz  dashboards/
	kubectl create configmap coherence-grafana-dashboards --from-file=dashboards/grafana \
		--dry-run=client -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION)/coherence-grafana-dashboards.yaml
	kubectl create configmap coherence-kibana-dashboards --from-file=dashboards/kibana \
		--dry-run=client -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION)/coherence-kibana-dashboards.yaml

# ----------------------------------------------------------------------------------------------------------------------
# Release the Coherence Operator to the gh-pages branch.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: release-ghpages
release-ghpages:  helm-chart docs release-dashboards
	mkdir -p /tmp/coherence-operator || true
	cp -R $(BUILD_OUTPUT) /tmp/coherence-operator
	cp $(BUILD_OUTPUT)/dashboards/$(VERSION)/coherence-dashboards.tar.gz /tmp/coherence-operator/_output/coherence-dashboards.tar.gz
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
# Release the Coherence Operator snapshot documentation.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: push-snapshot-docs
push-snapshot-docs: $(BUILD_TARGETS)/generate $(BUILD_TARGETS)/manifests docs
	rm -rf /tmp/coherence-operator || true
	mkdir -p /tmp/coherence-operator || true
	cp -R $(BUILD_OUTPUT)/ /tmp/coherence-operator
	ls -al /tmp/coherence-operator
	git stash save --keep-index --include-untracked || true
	git stash drop || true
	git checkout --track origin/gh-pages
	git config pull.rebase true
	git pull
	rm -rf docs/snapshot
	mv /tmp/coherence-operator/_output/docs/ docs/snapshot/
	ls -al docs/
	git add -A docs/snapshot/*
	git status
	git clean -d -f
	git status
	git commit -m "Release Coherence Operator snapshot docs $(VERSION)"
	git log -1
	git push origin gh-pages


# ----------------------------------------------------------------------------------------------------------------------
# Release the Coherence Operator.
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: release
release: ## Release the Operator
ifeq (true, $(RELEASE_DRY_RUN))
release: build-all-images release-ghpages
	@echo "release dry-run: would have pushed images"
else
release: build-all-images push-release-images release-ghpages
endif

# ----------------------------------------------------------------------------------------------------------------------
# Update the Operator version and all references to the previous version
# ----------------------------------------------------------------------------------------------------------------------
.PHONY: new-version
new-version: ## Update the Operator Version (must be run with NEXT_VERSION=x.y.z specified)
	$(SED) 's/$(subst .,\.,$(VERSION))/$(NEXT_VERSION)/g' Makefile
	$(SED) 's/$(subst .,\.,$(PREV_VERSION))/$(VERSION)/g' Makefile
	find docs \( -name '*.adoc' -o -name '*.yaml' \) -exec $(SED) 's/$(subst .,\.,$(VERSION))/$(NEXT_VERSION)/g' {} +
	find examples \( -name 'pom.xml' \) -exec $(SED) 's/<version>$(subst .,\.,$(VERSION))<\/version>/<version>$(NEXT_VERSION)<\/version>/g' {} +
	find examples \( -name '*.adoc' -o -name 'Dockerfile' \) -exec $(SED) 's/$(subst .,\.,$(VERSION))/$(NEXT_VERSION)/g' {} +
	find examples \( -name '*.yaml' -o -name '*.json' \) -exec $(SED) 's/$(subst .,\.,$(VERSION))/$(NEXT_VERSION)/g' {} +
	find config \( -name '*.yaml' -o -name '*.json' \) -exec $(SED) 's/$(subst .,\.,$(VERSION))/$(NEXT_VERSION)/g' {} +
	find helm-charts \( -name '*.yaml' -o -name '*.json' \) -exec $(SED) 's/$(subst .,\.,$(VERSION))/$(NEXT_VERSION)/g' {} +
	$(SED) -e 's/<revision>$(subst .,\.,$(VERSION))<\/revision>/<revision>$(NEXT_VERSION)<\/revision>/g' java/pom.xml

GIT_BRANCH="version-update-$(VERSION)"
GIT_LABEL="version-update"

.PHONY: new-version-pr
new-version-pr: ## Create a PR to update the version
	git config user.email "action@github.com"
	git config user.name "GitHub Action"
	git checkout -b $(GIT_BRANCH)
	git commit -am "Version update to $(VERSION)"
	git push --set-upstream origin $(GIT_BRANCH)

	gh label create "$(GIT_LABEL)" \
		--description "Pull requests with version update" \
		--force \
	|| true

	gh pr create \
		--title "Version update to $(VERSION)" \
		--body "Current pull request contains version update to version $(VERSION)" \
		--label "$(GIT_LABEL)" \
		--head $(GIT_BRANCH)

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
	curl -sSL https://github.com/github/licensed/releases/download/2.14.4/licensed-2.14.4-darwin-x64.tar.gz --header $(GH_AUTH) > licensed.tar.gz
else
	curl -sSL https://github.com/github/licensed/releases/download/2.14.4/licensed-2.14.4-linux-x64.tar.gz --header $(GH_AUTH) > licensed.tar.gz
endif
	tar -xzf licensed.tar.gz
	rm -f licensed.tar.gz
	mv ./licensed $(TOOLS_BIN)/licensed
	chmod +x $(TOOLS_BIN)/licensed
