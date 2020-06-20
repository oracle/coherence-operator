# ---------------------------------------------------------------------------
# Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
# ---------------------------------------------------------------------------
# This is the Makefile to build the Coherence Kubernetes Operator.
# ---------------------------------------------------------------------------

# The version of the Operator being build - this should be a valid SemVer format
VERSION ?= 3.0.0

# VERSION_SUFFIX is ann optional version suffix. For a full release this should be
# set to blank, for an interim release it should be set to a value to identify that
# release.
# For example if building the third release candidate this value might be
# set to VERSION_SUFFIX=RC3
# If VERSION_SUFFIX = DATE then the suffix will be a timestamp of the form yyMMddhhmm
# The default value for local and pipeline builds is "ci".
VERSION_SUFFIX ?= ci

# Set the full version string by combining the version and optional suffix
ifeq (, $(VERSION_SUFFIX))
VERSION_FULL := $(VERSION)
else
VERSION_FULL := $(VERSION)-$(VERSION_SUFFIX)
endif

# The operator version to use to run certification tests against
CERTIFICATION_VERSION ?= $(VERSION_FULL)

# A SPACE delimited list of previous Operator versions that are used to run the compatibility tests.
# These must be released versions as their released Helm charts will be downloaded prior to
# running the compatibility tests.
COMPATIBLE_VERSIONS = 2.1.0

# Capture the Git commit to add to the build information
GITCOMMIT       ?= $(shell git rev-list -1 HEAD)
GITREPO         := https://github.com/oracle/coherence-operator.git

CURRDIR         := $(shell pwd)

ARCH            ?= amd64
OS              ?= linux
UNAME_S         := $(shell uname -s)
GOPROXY         ?= https://proxy.golang.org

# Set the location of the Operator SDK executable
UNAME_S               = $(shell uname -s)
UNAME_M               = $(shell uname -m)
OPERATOR_SDK_VERSION := v0.18.0
OPERATOR_SDK          = $(CURRDIR)/etc/sdk/$(UNAME_S)-$(UNAME_M)/operator-sdk
OP_CHMOD             := $(shell chmod +x $(OPERATOR_SDK))

# The image prefix to use for Coherence images
COHERENCE_IMAGE_PREFIX ?= oraclecoherence/
# The Coherence image name to inject into the Helm chart
HELM_COHERENCE_IMAGE   ?= oraclecoherence/coherence-ce:14.1.1-0-1

# One may need to define RELEASE_IMAGE_PREFIX in the environment.
# For releases this will be docker.pkg.github.com/oracle/coherence-operator/
RELEASE_IMAGE_PREFIX ?= "oraclecoherence/"
OPERATOR_IMAGE_REPO  := $(RELEASE_IMAGE_PREFIX)coherence-operator
OPERATOR_IMAGE       := $(OPERATOR_IMAGE_REPO):$(VERSION_FULL)
UTILS_IMAGE          ?= $(OPERATOR_IMAGE_REPO):$(VERSION_FULL)-utils
TEST_USER_IMAGE      := $(RELEASE_IMAGE_PREFIX)operator-test-jib:$(VERSION_FULL)

RELEASE_DRY_RUN  ?= true
PRE_RELEASE      ?= true

# Extra arguments to pass to the go test command for the various test steps.
# For example, when running make e2e-test we can run just a single test such
# as the zone test using the go test -run=regex argument like this
#   make e2e-test GO_TEST_FLAGS='-run=^TestZone$$'
GO_TEST_FLAGS     ?= -timeout=20m
GO_TEST_FLAGS_E2E := -timeout=100m

# This is the Coherence image that will be used in the Go tests.
# Changing this variable will allow test builds to be run against different Coherence versions
TEST_COHERENCE_IMAGE ?= $(HELM_COHERENCE_IMAGE)

# default as in test/e2e/helper/proj_helpers.go
TEST_NAMESPACE ?= operator-test

CREATE_TEST_NAMESPACE ?= true

# Prometheus Operator settings
PROMETHEUS_INCLUDE_GRAFANA   ?= true
PROMETHEUS_OPERATOR_VERSION  ?= 8.13.7
GRAFANA_DASHBOARDS           ?= dashboards/grafana-legacy/

# Elasticsearch & Kibana settings
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

IMAGE_PULL_POLICY  ?=

# Env variable used by the kubectl test framework to locate the kubectl binary
TEST_ASSET_KUBECTL ?= $(shell which kubectl)

override BUILD_OUTPUT      := ./build/_output
override BUILD_PROPS       := $(BUILD_OUTPUT)/build.properties
override CHART_DIR         := $(BUILD_OUTPUT)/helm-charts
override PREV_CHART_DIR    := $(BUILD_OUTPUT)/previous-charts
override CRD_DIR           := deploy/crds
override TEST_LOGS_DIR     := $(BUILD_OUTPUT)/test-logs

GOS          = $(shell find pkg -type f -name "*.go" ! -name "*_test.go")
OPTESTGOS    = $(shell find cmd/optest -type f -name "*.go" ! -name "*_test.go")
UTILGOS      = $(shell find cmd/utilsinit -type f -name "*.go" ! -name "*_test.go")
COP_CHARTS   = $(shell find helm-charts/coherence-operator -type f)
DEPLOYS      = $(shell find deploy -type f -name "*.yaml")
CRD_VERSION  ?= "v1"

TEST_MANIFEST_DIR         := $(BUILD_OUTPUT)/manifest
TEST_MANIFEST_FILE        := test-manifest.yaml
TEST_LOCAL_MANIFEST_FILE  := local-manifest.yaml
TEST_GLOBAL_MANIFEST_FILE := global-manifest.yaml
TEST_SSL_SECRET           := coherence-ssl-secret

# ---------------------------------------------------------------------------
# Do a search and replace of properties in selected files in the Helm charts
# This is done because the Helm charts can be large and processing every file
# makes the build slower
# ---------------------------------------------------------------------------
define replaceprop
	for i in $(1); do \
		filename="$(CHART_DIR)/$${i}"; \
		echo "Replacing properties in file $${filename}"; \
		if [ -f $${filename} ]; then \
			temp_file=$(BUILD_OUTPUT)/temp.out; \
			awk -F'=' 'NR==FNR {a[$$1]=$$2;next} {for (i in a) {x = sprintf("\\$${%s}", i); gsub(x, a[i])}}1' $(BUILD_PROPS) $${filename} > $${temp_file}; \
			mv $${temp_file} $${filename}; \
		fi \
	done
endef

.PHONY: all
all: build-all-images

# ---------------------------------------------------------------------------
# Configure the build properties
# ---------------------------------------------------------------------------
$(BUILD_PROPS):
	# Ensures that build output directories exist
	@echo "Creating build directories"
	@mkdir -p $(BUILD_OUTPUT)
	@mkdir -p $(TEST_LOGS_DIR)
	@mkdir -p $(CHART_DIR)
	@mkdir -p $(PREV_CHART_DIR)
	# create build.properties
	rm -f $(BUILD_PROPS)
	printf "HELM_COHERENCE_IMAGE=$(HELM_COHERENCE_IMAGE)\n\
	UTILS_IMAGE=$(UTILS_IMAGE)\n\
	OPERATOR_IMAGE=$(OPERATOR_IMAGE)\n\
	VERSION_FULL=$(VERSION_FULL)\n\
	VERSION=$(VERSION)\n" > $(BUILD_PROPS)

# ---------------------------------------------------------------------------
# Builds the project, helm charts and Docker image
# ---------------------------------------------------------------------------
.PHONY: build-operator
build-operator: build-runner-artifacts $(BUILD_OUTPUT)/bin/operator

# ---------------------------------------------------------------------------
# Internal make step that builds the Operator Docker image and Helm charts
# ---------------------------------------------------------------------------
$(BUILD_OUTPUT)/bin/operator: export CGO_ENABLED = 0
$(BUILD_OUTPUT)/bin/operator: export GOARCH = $(ARCH)
$(BUILD_OUTPUT)/bin/operator: export GOOS = $(OS)
$(BUILD_OUTPUT)/bin/operator: export GO111MODULE = on
$(BUILD_OUTPUT)/bin/operator: export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
$(BUILD_OUTPUT)/bin/operator: export UTILS_IMAGE := $(UTILS_IMAGE)
$(BUILD_OUTPUT)/bin/operator: export VERSION_FULL := $(VERSION_FULL)
$(BUILD_OUTPUT)/bin/operator: $(GOS) $(DEPLOYS) $(CHART_DIR)/coherence-operator-$(VERSION_FULL).tgz ensure-sdk
	@echo "Building: $(OPERATOR_IMAGE)"
	@echo "Running Operator SDK build"
	BUILD_INFO="$(VERSION_FULL)|$(GITCOMMIT)|$$(date -u | tr ' ' '.')"; \
	$(OPERATOR_SDK) build $(OPERATOR_IMAGE) --verbose \
	    --image-build-args "--build-arg version=$(VERSION_FULL) --build-arg coherence_image=$(HELM_COHERENCE_IMAGE) --build-arg utils_image=$(UTILS_IMAGE)" \
	    --go-build-args "-o $(BUILD_OUTPUT)/bin/operator -ldflags -X=main.BuildInfo=$${BUILD_INFO}"

# ---------------------------------------------------------------------------
# Ensure Operator SDK is at the correct version
# ---------------------------------------------------------------------------
.PHONY: ensure-sdk
ensure-sdk:
	./hack/ensure-sdk.sh $(OPERATOR_SDK_VERSION)

# ---------------------------------------------------------------------------
# Internal make step that builds the Operator runner artifacts utility
# ---------------------------------------------------------------------------
.PHONY: build-runner-artifacts
build-runner-artifacts: $(BUILD_OUTPUT)/bin/runner

$(BUILD_OUTPUT)/bin/runner: export CGO_ENABLED = 0
$(BUILD_OUTPUT)/bin/runner: export GOARCH = $(ARCH)
$(BUILD_OUTPUT)/bin/runner: export GOOS = $(OS)
$(BUILD_OUTPUT)/bin/runner: export GO111MODULE = on
$(BUILD_OUTPUT)/bin/runner: $(GOS) $(DEPLOYS)
	BUILD_INFO="$(VERSION_FULL)|$(GITCOMMIT)|$$(date -u | tr ' ' '.')"; \
	go build -ldflags -X=main.BuildInfo=$${BUILD_INFO} -o $(BUILD_OUTPUT)/bin/runner ./cmd/runner

# ---------------------------------------------------------------------------
# Internal make step that builds the Operator test utility
# ---------------------------------------------------------------------------
.PHONY: build-op-test
build-op-test: $(BUILD_OUTPUT)/bin/op-test

$(BUILD_OUTPUT)/bin/op-test: export CGO_ENABLED = 0
$(BUILD_OUTPUT)/bin/op-test: export GOARCH = $(ARCH)
$(BUILD_OUTPUT)/bin/op-test: export GOOS = $(OS)
$(BUILD_OUTPUT)/bin/op-test: export GO111MODULE = on
$(BUILD_OUTPUT)/bin/op-test: $(GOS) $(DEPLOYS) $(OPTESTGOS)
	go build -o $(BUILD_OUTPUT)/bin/op-test ./cmd/optest

# ---------------------------------------------------------------------------
# Internal make step that builds the Operator utils init utility
# ---------------------------------------------------------------------------
.PHONY: build-utils-init
build-utils-init: $(BUILD_OUTPUT)/bin/utils-init

$(BUILD_OUTPUT)/bin/utils-init: export CGO_ENABLED = 0
$(BUILD_OUTPUT)/bin/utils-init: export GOARCH = $(ARCH)
$(BUILD_OUTPUT)/bin/utils-init: export GOOS = $(OS)
$(BUILD_OUTPUT)/bin/utils-init: export GO111MODULE = on
$(BUILD_OUTPUT)/bin/utils-init: $(GOS) $(DEPLOYS) $(UTILGOS)
	go build -o $(BUILD_OUTPUT)/bin/utils-init ./cmd/utilsinit

# ---------------------------------------------------------------------------
# Build the Coherence operator Helm chart and package it into a tar.gz
# ---------------------------------------------------------------------------
$(CHART_DIR)/coherence-operator-$(VERSION_FULL).tgz: $(COP_CHARTS) $(BUILD_PROPS)
	# Copy the Helm charts from their source location to the distribution folder
	@echo "Copying Operator chart to $(CHART_DIR)/coherence-operator"
	cp -R ./helm-charts/coherence-operator $(CHART_DIR)

	$(call replaceprop,coherence-operator/Chart.yaml coherence-operator/values.yaml coherence-operator/requirements.yaml coherence-operator/templates/deployment.yaml)

	# Package the chart into a .tr.gz - we don't use helm package as the version might not be SEMVER
	echo "Creating Helm chart package $(CHART_DIR)/coherence-operator"
	helm lint $(CHART_DIR)/coherence-operator
	tar -C $(CHART_DIR)/coherence-operator -czf $(CHART_DIR)/coherence-operator-$(VERSION_FULL).tgz .

# ---------------------------------------------------------------------------
# Build the Operator Helm chart and package it into a tar.gz
# ---------------------------------------------------------------------------
.PHONY: helm-chart
helm-chart: $(COP_CHARTS) $(BUILD_PROPS) $(CHART_DIR)/coherence-operator-$(VERSION_FULL).tgz


# ---------------------------------------------------------------------------
# Executes the Go unit tests that do not require a k8s cluster
# ---------------------------------------------------------------------------
.PHONY: test-operator
test-operator: export CGO_ENABLED = 0
test-operator: export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
test-operator: export UTILS_IMAGE := $(UTILS_IMAGE)
test-operator: build-operator
	@echo "Running operator tests"
	go test $(GO_TEST_FLAGS) -v ./cmd/... ./pkg/... \
	2>&1 | tee $(TEST_LOGS_DIR)/operator-test.out
	go run ./cmd/testreports/ -fail -suite-name-prefix=operator-test/ \
	    -input $(TEST_LOGS_DIR)/operator-test.out \
	    -output $(TEST_LOGS_DIR)/operator-test.xml

# ---------------------------------------------------------------------------
# Executes the Go end-to-end tests that require a k8s cluster using
# a LOCAL operator instance (i.e. the operator is not deployed to k8s).
# These tests will use whichever k8s cluster the local environment
# is pointing to.
# ---------------------------------------------------------------------------
.PHONY: e2e-local-test
e2e-local-test: export CGO_ENABLED = 0
e2e-local-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
e2e-local-test: export TEST_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_LOCAL_MANIFEST_FILE)
e2e-local-test: export TEST_GLOBAL_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_GLOBAL_MANIFEST_FILE)
e2e-local-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
e2e-local-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
e2e-local-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
e2e-local-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
e2e-local-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
e2e-local-test: export TEST_ASSET_KUBECTL := $(TEST_ASSET_KUBECTL)
e2e-local-test: export LOCAL_STORAGE_RESTART := $(LOCAL_STORAGE_RESTART)
e2e-local-test: export VERSION_FULL := $(VERSION_FULL)
e2e-local-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
e2e-local-test: export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
e2e-local-test: export UTILS_IMAGE := $(UTILS_IMAGE)
e2e-local-test: build-operator reset-namespace create-ssl-secrets operator-manifest uninstall-crds
	@echo "executing end-to-end tests"
	$(OPERATOR_SDK) test local ./test/e2e/local \
		--operator-namespace $(TEST_NAMESPACE) --watch-namespace  $(TEST_NAMESPACE)\
		--up-local --verbose --debug  --go-test-flags "$(GO_TEST_FLAGS_E2E)" \
		--local-operator-flags "--coherence-image=$(HELM_COHERENCE_IMAGE) --utils-image=$(UTILS_IMAGE)" \
		--namespaced-manifest=$(TEST_MANIFEST) \
		--global-manifest=$(TEST_GLOBAL_MANIFEST) \
		 2>&1 | tee $(TEST_LOGS_DIR)/operator-e2e-local-test.out
	$(MAKE) delete-namespace
	go run ./cmd/testreports/ -fail -suite-name-prefix=e2e-local-test/ \
	    -input $(TEST_LOGS_DIR)/operator-e2e-local-test.out \
	    -output $(TEST_LOGS_DIR)/operator-e2e-local-test.xml

# ---------------------------------------------------------------------------
# Run e2e local test in debug mode.
# This assumes that "make run-debug" has already been run so that a local
# Operator is running in debug mode and that the k8s namespace has been
# configured.
#
# Typically this step would be run with the GO_TEST_FLAGS variable set to
# run a specific test. For example to just run the ZoneTest...
#
# make debug-e2e-test GO_TEST_FLAGS='-run=^TestZone$$'
#
# ...where the -run argument is passed to go test and contains a reg-ex
# matching the name of the individual test to run.
#
# The reg-ex used to identify a test can also be used to run individual
# sub-tests. For example the status_ha_test.go file has a test called
# TestStatusHA that has a sub-test called HttpStatusHAHandler.
# To run this sub-test...
#
# make debug-e2e-test GO_TEST_FLAGS_E2E='-run=^TestStatusHA/HttpStatusHAHandler$$'
#
# ---------------------------------------------------------------------------
.PHONY: debug-e2e-local-test
debug-e2e-local-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
debug-e2e-local-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
debug-e2e-local-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
debug-e2e-local-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
debug-e2e-local-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
debug-e2e-local-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
debug-e2e-local-test: export TEST_ASSET_KUBECTL := $(TEST_ASSET_KUBECTL)
debug-e2e-local-test: export VERSION_FULL := $(VERSION_FULL)
debug-e2e-local-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
debug-e2e-local-test: export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
debug-e2e-local-test: export UTILS_IMAGE := $(UTILS_IMAGE)
debug-e2e-local-test:
	$(OPERATOR_SDK) test local ./test/e2e/local \
		--operator-namespace $(TEST_NAMESPACE) --watch-namespace $(TEST_NAMESPACE) \
		--verbose --debug  --go-test-flags \
		"$(GO_TEST_FLAGS_E2E)" --no-setup


# ---------------------------------------------------------------------------
# Executes the Go end-to-end tests that require a k8s cluster using
# a DEPLOYED operator instance (i.e. the operator Docker image is
# deployed to k8s). These tests will use whichever k8s cluster the
# local environment is pointing to.
# ---------------------------------------------------------------------------
.PHONY: e2e-test
e2e-test: export CGO_ENABLED = 0
e2e-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
e2e-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
e2e-test: export TEST_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_MANIFEST_FILE)
e2e-test: export TEST_GLOBAL_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_GLOBAL_MANIFEST_FILE)
e2e-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
e2e-test: export TEST_NAMESPACE := $(TEST_NAMESPACE)
e2e-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
e2e-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
e2e-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
e2e-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
e2e-test: export TEST_ASSET_KUBECTL := $(TEST_ASSET_KUBECTL)
e2e-test: export VERSION_FULL := $(VERSION_FULL)
e2e-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
e2e-test: export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
e2e-test: export UTILS_IMAGE := $(UTILS_IMAGE)
e2e-test: build-operator reset-namespace create-ssl-secrets operator-manifest uninstall-crds
	@echo "executing end-to-end tests"
	$(OPERATOR_SDK) test local ./test/e2e/remote --operator-namespace $(TEST_NAMESPACE) \
		--verbose --debug  --go-test-flags "$(GO_TEST_FLAGS_E2E)" \
		--namespaced-manifest=$(TEST_MANIFEST) \
		--global-manifest=$(TEST_GLOBAL_MANIFEST) \
		 2>&1 | tee $(TEST_LOGS_DIR)/operator-e2e-test.out
	$(MAKE) delete-namespace
	go run ./cmd/testreports/ -fail -suite-name-prefix=e2e-test/ \
	    -input $(TEST_LOGS_DIR)/operator-e2e-test.out \
	    -output $(TEST_LOGS_DIR)/operator-e2e-test.xml

# ---------------------------------------------------------------------------
# Run e2e test in debug mode.
# This assumes that "make run-debug" has already been run so that a local
# Operator is running in debug mode and that the k8s namespace has been
# configured.
#
# NOTE: The majority of e2e-test tests will FAIL if run woth a local operator
# due to the fact that either the Operator needs to access a Pod directly
# (for example in scaling tests) or the Pod needs to contact the Operator
# directly (for example in the zone tests). Due to the network constraints
# in k8s the Pods and Opererator cannot see eachother. For some debugging
# scenarios this may be OK but BEWARE!!
#
# Typically this step would be run with the GO_TEST_FLAGS variable set to
# run a specific test. For example to just run the ZoneTest...
#
# make debug-e2e-test GO_TEST_FLAGS='-run=^TestZone$$'
#
# ...where the -run argument is passed to go test and contains a reg-ex
# matching the name of the individual test to run.
#
# The reg-ex used to identify a test can also be used to run individual
# sub-tests. For example the scaling_test.go file has a test called
# TestScaling that has a sub-test called DownSafeScaling.
# To run this sub-test...
#
# make debug-e2e-test GO_TEST_FLAGS_E2E='-run=^TestScaling/DownSafeScaling$$'
#
# ---------------------------------------------------------------------------
.PHONY: debug-e2e-test
debug-e2e-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
debug-e2e-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
debug-e2e-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
debug-e2e-test: export TEST_NAMESPACE := $(TEST_NAMESPACE)
debug-e2e-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
debug-e2e-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
debug-e2e-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
debug-e2e-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
debug-e2e-test: export TEST_ASSET_KUBECTL := $(TEST_ASSET_KUBECTL)
debug-e2e-test: export VERSION_FULL := $(VERSION_FULL)
debug-e2e-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
debug-e2e-test: export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
debug-e2e-test: export UTILS_IMAGE := $(UTILS_IMAGE)
debug-e2e-test:
	$(OPERATOR_SDK) test local ./test/e2e/remote --operator-namespace $(TEST_NAMESPACE) \
		--verbose --debug  --go-test-flags "$(GO_TEST_FLAGS_E2E)" \
		--no-setup  2>&1 | tee $(TEST_LOGS_DIR)/operator-e2e-test.out


# ---------------------------------------------------------------------------
# Executes the Go end-to-end tests that require Prometheus in the k8s cluster
# using a LOCAL operator instance (i.e. the operator is not deployed to k8s).
#
# This target DOES NOT install Prometheus, use the e2e-prometheus-test target
# to fully reset the test namespace.
#
# These tests will use whichever k8s cluster the local environment
# is pointing to.
# ---------------------------------------------------------------------------
.PHONY: run-prometheus-test
run-prometheus-test: export CGO_ENABLED = 0
run-prometheus-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
run-prometheus-test: export TEST_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_LOCAL_MANIFEST_FILE)
run-prometheus-test: export TEST_GLOBAL_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_GLOBAL_MANIFEST_FILE)
run-prometheus-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
run-prometheus-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
run-prometheus-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
run-prometheus-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
run-prometheus-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
run-prometheus-test: export TEST_ASSET_KUBECTL := $(TEST_ASSET_KUBECTL)
run-prometheus-test: export LOCAL_STORAGE_RESTART := $(LOCAL_STORAGE_RESTART)
run-prometheus-test: export VERSION_FULL := $(VERSION_FULL)
run-prometheus-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
run-prometheus-test: export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
run-prometheus-test: export UTILS_IMAGE := $(UTILS_IMAGE)
run-prometheus-test: build-operator create-ssl-secrets operator-manifest
	@echo "executing Prometheus end-to-end tests"
	$(OPERATOR_SDK) test local ./test/e2e/prometheus \
		--operator-namespace $(TEST_NAMESPACE) --watch-namespace  $(TEST_NAMESPACE)\
		--up-local --verbose --debug  --go-test-flags "$(GO_TEST_FLAGS_E2E)" \
		--local-operator-flags "--coherence-image=$(HELM_COHERENCE_IMAGE) --utils-image=$(UTILS_IMAGE)" \
		--namespaced-manifest=$(TEST_MANIFEST) \
		--global-manifest=$(TEST_GLOBAL_MANIFEST) \
		 2>&1 | tee $(TEST_LOGS_DIR)/operator-e2e-prometheus-test.out
	go run ./cmd/testreports/ -fail -suite-name-prefix=e2e-prometheus-test/ \
	    -input $(TEST_LOGS_DIR)/operator-e2e-prometheus-test.out \
	    -output $(TEST_LOGS_DIR)/operator-e2e-prometheus-test.xml

.PHONY: e2e-prometheus-test
e2e-prometheus-test: reset-namespace install-prometheus
	$(MAKE) $(MAKEFLAGS) run-prometheus-test \
	; rc=$$? \
	; $(MAKE) $(MAKEFLAGS) uninstall-prometheus \
	; $(MAKE) $(MAKEFLAGS) delete-namespace \
	; exit $$rc


# ---------------------------------------------------------------------------
# Executes the Go end-to-end tests that require Elasticsearch in the k8s cluster
# using a LOCAL operator instance (i.e. the operator is not deployed to k8s).
#
# This target DOES NOT install Elasticsearch, use the e2e-elastic-test target
# to fully reset the test namespace.
#
# These tests will use whichever k8s cluster the local environment
# is pointing to.
# ---------------------------------------------------------------------------
.PHONY: run-elastic-test
run-elastic-test: export CGO_ENABLED = 0
run-elastic-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
run-elastic-test: export TEST_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_LOCAL_MANIFEST_FILE)
run-elastic-test: export TEST_GLOBAL_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_GLOBAL_MANIFEST_FILE)
run-elastic-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
run-elastic-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
run-elastic-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
run-elastic-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
run-elastic-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
run-elastic-test: export TEST_ASSET_KUBECTL := $(TEST_ASSET_KUBECTL)
run-elastic-test: export LOCAL_STORAGE_RESTART := $(LOCAL_STORAGE_RESTART)
run-elastic-test: export VERSION_FULL := $(VERSION_FULL)
run-elastic-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
run-elastic-test: export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
run-elastic-test: export UTILS_IMAGE := $(UTILS_IMAGE)
run-elastic-test: export KIBANA_INDEX_PATTERN := $(KIBANA_INDEX_PATTERN)
run-elastic-test: build-operator create-ssl-secrets operator-manifest
	@echo "executing Elasticsearch end-to-end tests"
	$(OPERATOR_SDK) test local ./test/e2e/elastic \
		--operator-namespace $(TEST_NAMESPACE) --watch-namespace  $(TEST_NAMESPACE)\
		--up-local --verbose --debug  --go-test-flags "$(GO_TEST_FLAGS_E2E)" \
		--local-operator-flags "--coherence-image=$(HELM_COHERENCE_IMAGE) --utils-image=$(UTILS_IMAGE)" \
		--namespaced-manifest=$(TEST_MANIFEST) \
		--global-manifest=$(TEST_GLOBAL_MANIFEST) \
		 2>&1 | tee $(TEST_LOGS_DIR)/operator-e2e-elastic-test.out
	go run ./cmd/testreports/ -fail -suite-name-prefix=e2e-elastic-test/ \
	    -input $(TEST_LOGS_DIR)/operator-e2e-elastic-test.out \
	    -output $(TEST_LOGS_DIR)/operator-e2e-elastic-test.xml

.PHONY: e2e-elastic-test
e2e-elastic-test: reset-namespace install-elastic
	$(MAKE) $(MAKEFLAGS) run-elastic-test \
	; rc=$$? \
	; $(MAKE) $(MAKEFLAGS) uninstall-elastic \
	; $(MAKE) $(MAKEFLAGS) delete-namespace \
	; exit $$rc


# ---------------------------------------------------------------------------
# Executes the Go end-to-end Operator Helm chart tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# Note that the namespace will be created if it does not exist.
# ---------------------------------------------------------------------------
.PHONY: helm-test
helm-test: export CGO_ENABLED = 0
helm-test: export TEST_NAMESPACE := $(TEST_NAMESPACE)
helm-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
helm-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
helm-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
helm-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
helm-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
helm-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
helm-test: export VERSION_FULL := $(VERSION_FULL)
helm-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
helm-test: export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
helm-test: export UTILS_IMAGE := $(UTILS_IMAGE)
helm-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
helm-test: build-operator reset-namespace create-ssl-secrets
	@echo "executing Operator Helm Chart end-to-end tests"
	$(OPERATOR_SDK) test local ./test/e2e/helm --operator-namespace $(TEST_NAMESPACE) \
		--verbose --debug  --go-test-flags "$(GO_TEST_FLAGS_E2E)" \
		--no-setup  2>&1 | tee $(TEST_LOGS_DIR)/operator-e2e-helm-test.out
	$(MAKE) uninstall-crds
	$(MAKE) delete-namespace
	go run ./cmd/testreports/ -fail -suite-name-prefix=e2e-helm-test/ \
	    -input $(TEST_LOGS_DIR)/operator-e2e-helm-test.out \
	    -output $(TEST_LOGS_DIR)/operator-e2e-helm-test.xml

# ---------------------------------------------------------------------------
# Executes the Go end-to-end Operator Compatibility tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# ---------------------------------------------------------------------------
.PHONY: helm-test
compatibility-test: export CGO_ENABLED = 0
compatibility-test: export TEST_NAMESPACE := $(TEST_NAMESPACE)
compatibility-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
compatibility-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
compatibility-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
compatibility-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
compatibility-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
compatibility-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
compatibility-test: export VERSION_FULL := $(VERSION_FULL)
compatibility-test: export COMPATIBLE_VERSIONS := $(COMPATIBLE_VERSIONS)
compatibility-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
compatibility-test: export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
compatibility-test: export UTILS_IMAGE := $(UTILS_IMAGE)
compatibility-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
compatibility-test: build-operator clean-namespace reset-namespace create-ssl-secrets get-previous
	@echo "executing Operator compatibility tests"
	$(OPERATOR_SDK) test local ./test/e2e/compatibility --operator-namespace $(TEST_NAMESPACE) \
		--verbose --debug  --go-test-flags "$(GO_TEST_FLAGS_E2E)" \
		--no-setup  2>&1 | tee $(TEST_LOGS_DIR)/operator-compatibility-test.out
	$(MAKE) uninstall-crds
	$(MAKE) delete-namespace
	go run ./cmd/testreports/ -fail -suite-name-prefix=compatibility-test/ \
	    -input $(TEST_LOGS_DIR)/operator-compatibility-test.out \
	    -output $(TEST_LOGS_DIR)/operator-compatibility-test.xml

# ---------------------------------------------------------------------------
# Obtain the previous versions of the Operator Helm chart that will be used
# torun compatibiity tests.
# ---------------------------------------------------------------------------
.PHONY: get-previous
get-previous: $(BUILD_PROPS)
	for i in $(COMPATIBLE_VERSIONS); do \
      FILE=$(PREV_CHART_DIR)/coherence-operator-$${i}.tgz; \
      DIR=$(PREV_CHART_DIR)/coherence-operator-$${i}; \
      if [ ! -f "$${FILE}" ]; then \
	    echo "Downloading Operator Helm chart version $${i} to file $${FILE}"; \
	    curl -X GET https://oracle.github.io/coherence-operator/charts/coherence-operator-$${i}.tgz -o $${FILE}; \
      fi; \
	  echo "Unpacking Operator Helm chart version $${i} to $${DIR}"; \
      rm -rf $${DIR}; \
      mkdir $${DIR}; \
      tar -C $${DIR} -xzf $${FILE}; \
    done

# ---------------------------------------------------------------------------
# Executes the Go end-to-end Operator certification tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# Note that the namespace will be created if it does not exist.
# ---------------------------------------------------------------------------
.PHONY: certification-test
certification-test: export CGO_ENABLED = 0
certification-test: export TEST_NAMESPACE := $(TEST_NAMESPACE)
certification-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
certification-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
certification-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
certification-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
certification-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
certification-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
certification-test: export VERSION_FULL := $(VERSION_FULL)
certification-test: export CERTIFICATION_VERSION := $(CERTIFICATION_VERSION)
certification-test: export OPERATOR_IMAGE_REPO := $(OPERATOR_IMAGE_REPO)
certification-test: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
certification-test: export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
certification-test: export UTILS_IMAGE := $(UTILS_IMAGE)
certification-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
certification-test: install-certification
	$(MAKE) $(MAKEFLAGS) run-certification \
	; rc=$$? \
	; $(MAKE) $(MAKEFLAGS) cleanup-certification \
	; exit $$rc


# ---------------------------------------------------------------------------
# Install the Operator prior to running compatability tests.
# ---------------------------------------------------------------------------
.PHONY: install-certification
install-certification: export TEST_NAMESPACE := $(TEST_NAMESPACE)
install-certification: export VERSION := $(VERSION)
install-certification: export VERSION_FULL := $(VERSION_FULL)
install-certification: export CERTIFICATION_VERSION := $(CERTIFICATION_VERSION)
install-certification: build-operator reset-namespace create-ssl-secrets
ifeq ("$(CERTIFICATION_VERSION)","$(VERSION_FULL)")
	helm install --atomic --namespace $(TEST_NAMESPACE) --wait operator $(CHART_DIR)/coherence-operator
else
	helm repo add coherence https://oracle.github.io/coherence-operator/charts || true
	helm repo update || true
	helm install --atomic --namespace $(TEST_NAMESPACE) --wait --version operator ./helm-charts/coherence-operator
endif

# ---------------------------------------------------------------------------
# Executes the Go end-to-end Operator certification tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# Note that the namespace will be created if it does not exist.
# ---------------------------------------------------------------------------
.PHONY: run-certification
run-certification: export CGO_ENABLED = 0
run-certification: export TEST_NAMESPACE := $(TEST_NAMESPACE)
run-certification: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
run-certification: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
run-certification: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
run-certification: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
run-certification: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
run-certification: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
run-certification: export VERSION_FULL := $(VERSION_FULL)
run-certification: export CERTIFICATION_VERSION := $(CERTIFICATION_VERSION)
run-certification: export OPERATOR_IMAGE_REPO := $(OPERATOR_IMAGE_REPO)
run-certification: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
run-certification: export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
run-certification: export UTILS_IMAGE := $(UTILS_IMAGE)
run-certification: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
run-certification:
	@echo "Executing Operator certification tests"
	$(OPERATOR_SDK) test local ./test/certification --operator-namespace $(TEST_NAMESPACE) \
		--verbose --debug  --go-test-flags "$(GO_TEST_FLAGS_E2E)" \
		--no-setup  2>&1 | tee $(TEST_LOGS_DIR)/operator-certification-test.out
	go run ./cmd/testreports/ -fail -suite-name-prefix=certification-test/ \
	    -input $(TEST_LOGS_DIR)/operator-certification-test.out \
	    -output $(TEST_LOGS_DIR)/operator-certification-test.xml

# ---------------------------------------------------------------------------
# Clean up after to running compatability tests.
# ---------------------------------------------------------------------------
.PHONY: cleanup-certification
cleanup-certification: export TEST_NAMESPACE := $(TEST_NAMESPACE)
cleanup-certification:
	helm delete --namespace $(TEST_NAMESPACE) operator || true
	$(MAKE) uninstall-crds
	$(MAKE) delete-namespace


# ---------------------------------------------------------------------------
# Install CRDs into Kubernetes.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
# ---------------------------------------------------------------------------
.PHONY: install-crds
install-crds: uninstall-crds
	@echo "Installing CRDs $(CRD_VERSION)"
ifeq ("$(CRD_VERSION)","v1beta1")
	kubectl --validate=false create -f deploy/crds/v1beta1/coherence.oracle.com_coherence_crd.yaml || true
else
	kubectl --validate=false create -f deploy/crds/coherence.oracle.com_coherence_crd.yaml || true
endif

# ---------------------------------------------------------------------------
# Uninstall CRDs from Kubernetes.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
# ---------------------------------------------------------------------------
.PHONY: uninstall-crds
uninstall-crds: $(BUILD_PROPS)
	@echo "Removing CRDs"
	kubectl delete crd coherence.coherence.oracle.com || true

# ---------------------------------------------------------------------------
# This step will run the Operator SDK code generators.
# These commands will generate the CRD files from the API structs and will
# also generate the Go DeepCopy code for the API structs.
# This step would require running if any of the structs in the files under
# the pkg/apis directory have been changed.
# ---------------------------------------------------------------------------
.PHONY: generate
generate: export GOPATH := $(GOPATH)
generate:
	@echo "Generating deep copy code"
	$(OPERATOR_SDK) generate k8s --verbose
	@echo "Generating v1beta1 CRDs"
	$(OPERATOR_SDK) generate crds --verbose --crd-version v1beta1
	mv deploy/crds/coherence.oracle.com_coherence_crd.yaml deploy/crds/v1beta1/
	@echo "Generating v1 CRDs"
	$(OPERATOR_SDK) generate crds --verbose --crd-version v1
	@echo "Generating OpenAPI"
	@which $(BUILD_OUTPUT)/bin/openapi-gen > /dev/null || go build -o $(BUILD_OUTPUT)/bin/openapi-gen k8s.io/kube-openapi/cmd/openapi-gen
	$(BUILD_OUTPUT)/bin/openapi-gen --logtostderr=true \
		-i ./pkg/apis/coherence/v1 \
		-o "" \
		-O zz_generated.openapi \
		-p ./pkg/apis/coherence/v1 \
		-h ./hack/boilerplate.go.txt \
		-r "-"
	@echo "Getting kustomize"
	go get sigs.k8s.io/kustomize || true
	@echo "Applying kustomize to v1 CRDs"
	$(GOPATH)/bin/kustomize build deploy/crds -o deploy/crds/coherence.oracle.com_coherence_crd.yaml
	$(MAKE) api-doc-gen
	@echo "Getting go-bindata"
	go get -u github.com/shuLhan/go-bindata/... || true
	@echo "Embedding CRDs"
	$(GOPATH)/bin/go-bindata -o pkg/operator/zz_generated.assets.go -ignore .DS_Store deploy/crds/...

# ---------------------------------------------------------------------------
# Generate API docs
# ---------------------------------------------------------------------------
.PHONY: api-doc-gen
api-doc-gen:
	go run ./cmd/docgen/ \
		pkg/apis/coherence/v1/coherenceresourcespec_types.go \
		pkg/apis/coherence/v1/coherence_types.go \
		pkg/apis/coherence/v1/coherenceresource_types.go \
		> docs/about/04_coherence_spec.adoc

# ---------------------------------------------------------------------------
# Clean-up all of the build artifacts
# ---------------------------------------------------------------------------
.PHONY: clean
clean:
	rm -rf build/_output
	mvn $(USE_MAVEN_SETTINGS) -f java clean
	mvn $(USE_MAVEN_SETTINGS) -f examples clean

# ---------------------------------------------------------------------------
# Create the k8s yaml manifest that will be used by the Operator SDK to
# install the Operator when running e2e tests.
# ---------------------------------------------------------------------------
.PHONY: operator-manifest
operator-manifest: export TEST_NAMESPACE := $(TEST_NAMESPACE)
operator-manifest: export TEST_MANIFEST_DIR := $(TEST_MANIFEST_DIR)
operator-manifest: export TEST_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_MANIFEST_FILE)
operator-manifest: export TEST_LOCAL_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_LOCAL_MANIFEST_FILE)
operator-manifest: export TEST_GLOBAL_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_GLOBAL_MANIFEST_FILE)
operator-manifest: $(CHART_DIR)/coherence-operator-$(VERSION_FULL).tgz
	@mkdir -p $(TEST_MANIFEST_DIR)
	go run ./cmd/manifestutil/

# ---------------------------------------------------------------------------
# Generate the keys and certs used in tests.
# ---------------------------------------------------------------------------
$(BUILD_OUTPUT)/certs:
	@echo "Generating test keys and certs"
	./hack/keys.sh

# ---------------------------------------------------------------------------
# Delete and re-create the test namespace
# ---------------------------------------------------------------------------
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

# ---------------------------------------------------------------------------
# Delete the test namespace
# ---------------------------------------------------------------------------
.PHONY: delete-namespace
delete-namespace: clean-namespace
ifeq ($(CREATE_TEST_NAMESPACE),true)
	@echo "Deleting test namespace $(TEST_NAMESPACE)"
	kubectl delete namespace $(TEST_NAMESPACE) --force --grace-period=0 && echo "deleted namespace" || true
endif
	kubectl delete clusterrole operator-test-coherence-operator --force --grace-period=0 && echo "deleted namespace" || true
	kubectl delete clusterrolebinding operator-test-coherence-operator --force --grace-period=0 && echo "deleted namespace" || true


# ---------------------------------------------------------------------------
# Delete all resource from the test namespace
# ---------------------------------------------------------------------------
.PHONY: clean-namespace
clean-namespace: delete-coherence-clusters
	for i in $$(kubectl -n $(TEST_NAMESPACE) get all -o name); do \
		echo "Deleting $${i} from test namespace $(TEST_NAMESPACE)" \
		kubectl -n $(TEST_NAMESPACE) delete $${i}; \
	done

# ---------------------------------------------------------------------------
# Create the k8s secret to use in SSL/TLS testing.
# ---------------------------------------------------------------------------
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

# ---------------------------------------------------------------------------
# Build the Java artifacts
# ---------------------------------------------------------------------------
.PHONY: build-mvn
build-mvn:
	mvn $(USE_MAVEN_SETTINGS) -B -f java package -DskipTests

# ---------------------------------------------------------------------------
# Build and test the Java artifacts
# ---------------------------------------------------------------------------
.PHONY: test-mvn
test-mvn: build-mvn
	mvn $(USE_MAVEN_SETTINGS) -B -f java verify

# ---------------------------------------------------------------------------
# Build the examples
# ---------------------------------------------------------------------------
.PHONY: build-examples
build-examples:
	mvn $(USE_MAVEN_SETTINGS) -B -f examples package -DskipTests

# ---------------------------------------------------------------------------
# Build and test the examples
# ---------------------------------------------------------------------------
.PHONY: test-examples
test-examples: build-examples
	mvn $(USE_MAVEN_SETTINGS) -B -f examples verify

# ---------------------------------------------------------------------------
# Run all unit tests (both Go and Java)
# ---------------------------------------------------------------------------
.PHONY: test-all
test-all: test-mvn test-operator

# ---------------------------------------------------------------------------
# Push the Operator Docker image
# ---------------------------------------------------------------------------
.PHONY: push-operator-image
push-operator-image: build-operator
	@echo "Pushing $(OPERATOR_IMAGE)"
	docker push $(OPERATOR_IMAGE)

# ---------------------------------------------------------------------------
# Build the Operator Utils Docker image
# ---------------------------------------------------------------------------
.PHONY: build-utils-image
build-utils-image: build-mvn build-runner-artifacts build-utils-init build-op-test
	cp $(BUILD_OUTPUT)/bin/op-test    java/coherence-utils/target/docker/op-test
	cp $(BUILD_OUTPUT)/bin/utils-init java/coherence-utils/target/docker/utils-init
	cp $(BUILD_OUTPUT)/bin/runner     java/coherence-utils/target/docker/runner
	cp -r image/config                java/coherence-utils/target/docker/config
	cp -r image/logging               java/coherence-utils/target/docker/logging
	docker build -t $(UTILS_IMAGE) java/coherence-utils/target/docker

# ---------------------------------------------------------------------------
# Push the Operator Utils Docker image
# ---------------------------------------------------------------------------
.PHONY: push-utils-image
push-utils-image:
	@echo "Pushing $(UTILS_IMAGE)"
	docker push $(UTILS_IMAGE)

# ---------------------------------------------------------------------------
# Build the Operator JIB Test image
# ---------------------------------------------------------------------------
.PHONY: build-jib-image
build-jib-image: build-mvn
	mvn $(USE_MAVEN_SETTINGS) -B -f java package jib:dockerBuild -DskipTests -Djib.to.image=$(TEST_USER_IMAGE)

# ---------------------------------------------------------------------------
# Push the Operator JIB Test Docker images
# ---------------------------------------------------------------------------
.PHONY: push-jib-image
push-jib-image:
	@echo "Pushing $(TEST_USER_IMAGE)"
	docker push $(TEST_USER_IMAGE)

# ---------------------------------------------------------------------------
# Build all of the Docker images
# ---------------------------------------------------------------------------
.PHONY: build-all-images
build-all-images: build-operator build-utils-image build-jib-image

# ---------------------------------------------------------------------------
# Push all of the Docker images
# ---------------------------------------------------------------------------
.PHONY: push-all-images
push-all-images: push-operator-image push-utils-image push-jib-image

# ---------------------------------------------------------------------------
# Push all of the Docker images that are released
# ---------------------------------------------------------------------------
.PHONY: push-release-images
push-release-images: push-operator-image push-utils-image

# ---------------------------------------------------------------------------
# Build everything
# ---------------------------------------------------------------------------
.PHONY: build-all
build-all: build-mvn build-operator


# ---------------------------------------------------------------------------
# Run the Operator locally.
#
# To exit out of the local Operator you can use ctrl-c or ctrl-z but
# sometimes this leaves orphaned processes on the local machine so
# ensure these are killed run "make debug-stop"
# ---------------------------------------------------------------------------
.PHONY: run
run: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
run: export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
run: export UTILS_IMAGE := $(UTILS_IMAGE)
run: export VERSION_FULL := $(VERSION_FULL)
run: export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
run: export UTILS_IMAGE := $(UTILS_IMAGE)
run: 
	BUILD_INFO="$(VERSION_FULL)|$(GITCOMMIT)|$$(date -u | tr ' ' '.')"; \
	$(OPERATOR_SDK) run local --watch-namespace=$(TEST_NAMESPACE) \
	--go-ldflags="-X=main.BuildInfo=$${BUILD_INFO}" \
	--operator-flags="--coherence-image=$(HELM_COHERENCE_IMAGE) \
	                  --utils-image=$(UTILS_IMAGE)" \
	2>&1 | tee $(TEST_LOGS_DIR)/operator-debug.out

# ---------------------------------------------------------------------------
# Run the Operator locally after deleting and recreating the test namespace.
# ---------------------------------------------------------------------------
.PHONY: run-clean
run-clean: reset-namespace run

# ---------------------------------------------------------------------------
# Run the Operator in locally debug mode,
# Running this task will start the Operator and pause it until a Delve
# is attached.
#
# To exit out of the local Operator you can use ctrl-c or ctrl-z but
# sometimes this leaves orphaned processes on the local machine so
# ensure these are killed run "make debug-stop"
# ---------------------------------------------------------------------------
.PHONY: run-debug
run-debug: export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
run-debug: export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
run-debug: export UTILS_IMAGE := $(UTILS_IMAGE)
run-debug: export VERSION_FULL := $(VERSION_FULL)
run-debug: export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
run-debug: export UTILS_IMAGE := $(UTILS_IMAGE)
run-debug: 
	BUILD_INFO="$(VERSION_FULL)|$(GITCOMMIT)|$$(date -u | tr ' ' '.')"; \
	$(OPERATOR_SDK) run local --watch-namespace=$(TEST_NAMESPACE) \
	--go-ldflags="-X=main.BuildInfo=$${BUILD_INFO}" \
	--operator-flags="--coherence-image=$(HELM_COHERENCE_IMAGE) --utils-image=$(UTILS_IMAGE)" \
	--enable-delve \
	2>&1 | tee $(TEST_LOGS_DIR)/operator-debug.out

# ---------------------------------------------------------------------------
# Run the Operator locally in debug mode after deleting and recreating
# the test namespace.
# ---------------------------------------------------------------------------
.PHONY: run-debug-clean
run-debug-clean: reset-namespace run-debug

# ---------------------------------------------------------------------------
# Kill any locally running Operator
# ---------------------------------------------------------------------------
.PHONY: stop
stop:
	./hack/kill-local.sh


# ---------------------------------------------------------------------------
# Start a Kind cluster
# ---------------------------------------------------------------------------
.PHONY: kind
kind: export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
kind:
	./hack/kind.sh
	docker pull $(HELM_COHERENCE_IMAGE)
	kind load docker-image --name operator $(HELM_COHERENCE_IMAGE)

# ---------------------------------------------------------------------------
# Start a Kind 1.12 cluster
# ---------------------------------------------------------------------------
.PHONY: kind-12
kind-12: kind-12-start kind-load

.PHONY: kind-12-start
kind-12-start:
	./hack/kind.sh --image "kindest/node:v1.12.10@sha256:faeb82453af2f9373447bb63f50bae02b8020968e0889c7fa308e19b348916cb"
	docker pull $(HELM_COHERENCE_IMAGE) || true
	kind load docker-image --name operator $(HELM_COHERENCE_IMAGE) || true

# ---------------------------------------------------------------------------
# Start a Kind 1.18 cluster
# ---------------------------------------------------------------------------
.PHONY: kind-18
kind-18: kind-18-start kind-load

.PHONY: kind-18-start
kind-18-start:
	./hack/kind.sh --image "kindest/node:v1.18.2@sha256:7b27a6d0f2517ff88ba444025beae41491b016bc6af573ba467b70c5e8e0d85f"
	docker pull $(HELM_COHERENCE_IMAGE) || true
	kind load docker-image --name operator $(HELM_COHERENCE_IMAGE) || true

# ---------------------------------------------------------------------------
# Load images into Kind
# ---------------------------------------------------------------------------
.PHONY: kind-load
kind-load:
	kind load docker-image --name operator $(OPERATOR_IMAGE)|| true
	kind load docker-image --name operator $(UTILS_IMAGE)|| true
	kind load docker-image --name operator $(TEST_USER_IMAGE)|| true

# ---------------------------------------------------------------------------
# Install the Operator Helm chart.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
# ---------------------------------------------------------------------------
.PHONY: operator-helm-install
operator-helm-install: operator-helm-delete build-operator reset-namespace create-ssl-secrets
	helm install --name operator --namespace $(TEST_NAMESPACE) $(CHART_DIR)/coherence-operator


# ---------------------------------------------------------------------------
# Uninstall the Operator Helm chart.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
# ---------------------------------------------------------------------------
.PHONY: operator-helm-delete
operator-helm-delete:
	helm delete --purge operator || true

# ---------------------------------------------------------------------------
# Install Prometheus
# ---------------------------------------------------------------------------
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

# ---------------------------------------------------------------------------
# Uninstall Prometheus
# ---------------------------------------------------------------------------
.PHONY: uninstall-prometheus
uninstall-prometheus:
	kubectl -n $(TEST_NAMESPACE) delete -f etc/prometheus.yaml || true
	kubectl -n $(TEST_NAMESPACE) delete configmap coherence-grafana-dashboards || true
	helm --namespace $(TEST_NAMESPACE) delete prometheus || true
	kubectl delete -f etc/prometheus-rbac.yaml || true

# ---------------------------------------------------------------------------
# Start a port-forward process to the Grafana Pod.
# ---------------------------------------------------------------------------
.PHONY: port-forward-grafana
port-forward-grafana: export GRAFANA_POD := $(shell kubectl -n $(TEST_NAMESPACE) get pod -l app.kubernetes.io/name=grafana -o name)
port-forward-grafana:
	@echo "Reach Grafana on http://127.0.0.1:3000"
	@echo "User: admin Password: prom-operator"
	kubectl -n $(TEST_NAMESPACE) port-forward $(GRAFANA_POD) 3000:3000

# ---------------------------------------------------------------------------
# Install Elasticsearch & Kibana
# ---------------------------------------------------------------------------
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

# ---------------------------------------------------------------------------
# Uninstall Elasticsearch & Kibana
# ---------------------------------------------------------------------------
.PHONY: uninstall-elastic
uninstall-elastic:
	helm uninstall --namespace $(TEST_NAMESPACE) kibana || true
	helm uninstall --namespace $(TEST_NAMESPACE) elasticsearch || true
	kubectl -n $(TEST_NAMESPACE) delete pvc elasticsearch-master-elasticsearch-master-0 || true

# ---------------------------------------------------------------------------
# Start a port-forward process to the Kibana Pod.
# ---------------------------------------------------------------------------
.PHONY: port-forward-kibana
port-forward-kibana: export KIBANA_POD := $(shell kubectl -n $(TEST_NAMESPACE) get pod -l app=kibana -o name)
port-forward-kibana:
	@echo "Reach Kibana on http://127.0.0.1:5601"
	kubectl -n $(TEST_NAMESPACE) port-forward $(KIBANA_POD) 5601:5601

# ---------------------------------------------------------------------------
# Start a port-forward process to the Elasticsearch Pod.
# ---------------------------------------------------------------------------
.PHONY: port-forward-es
port-forward-es: export ES_POD := $(shell kubectl -n $(TEST_NAMESPACE) get pod -l app=elasticsearch-master -o name)
port-forward-es:
	@echo "Reach Elasticsearch on http://127.0.0.1:9200"
	kubectl -n $(TEST_NAMESPACE) port-forward $(ES_POD) 9200:9200


# ---------------------------------------------------------------------------
# Delete all of the Coherence resources from the test namespace.
# ---------------------------------------------------------------------------
.PHONY: delete-coherence-clusters
delete-coherence-clusters:
	for i in $$(kubectl -n  $(TEST_NAMESPACE) get coherence -o name); do \
		kubectl -n $(TEST_NAMESPACE) delete $${i}; \
	done

# ---------------------------------------------------------------------------
# Obtain the golangci-lint binary
# ---------------------------------------------------------------------------
$(BUILD_OUTPUT)/bin/golangci-lint:
	@mkdir -p $(BUILD_OUTPUT)/bin
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BUILD_OUTPUT)/bin v1.21.0

# ---------------------------------------------------------------------------
# Executes golangci-lint to perform various code review checks on the source.
# ---------------------------------------------------------------------------
.PHONY: golangci
golangci: $(BUILD_OUTPUT)/bin/golangci-lint
	$(BUILD_OUTPUT)/bin/golangci-lint run -v --timeout=5m --skip-files=zz_.*,generated/*  ./pkg/... ./cmd/...


# ---------------------------------------------------------------------------
# Performs a copyright check.
# To add exclusions add the file or folder pattern using the -X parameter.
# Add directories to be scanned at the end of the parameter list.
# ---------------------------------------------------------------------------
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
	  -X helm-charts/coherence-operator/charts/prometheus-operator/ \
	  -X helm-charts/coherence-operator/templates/NOTES.txt \
	  -X .iml \
	  -X Jenkinsfile \
	  -X .jar \
	  -X .jks \
	  -X .json \
	  -X LICENSE.txt \
	  -X Makefile \
	  -X .md \
	  -X .sh \
	  -X temp/ \
	  -X temp/olm/ \
	  -X /test-report.xml \
	  -X THIRD_PARTY_LICENSES.txt \
	  -X tools.go \
	  -X .tpl \
	  -X .yaml \
	  -X zz_generated.

# ---------------------------------------------------------------------------
# Executes the code review targets.
# ---------------------------------------------------------------------------
.PHONY: code-review
code-review: export MAVEN_USER := $(MAVEN_USER)
code-review: export MAVEN_PASSWORD := $(MAVEN_PASSWORD)
code-review: golangci copyright
	mvn $(USE_MAVEN_SETTINGS) -B -f java validate -DskipTests -P checkstyle
	mvn $(USE_MAVEN_SETTINGS) -B -f examples validate -DskipTests -P checkstyle

# ---------------------------------------------------------------------------
# Display the full version string for the artifacts that would be built.
# ---------------------------------------------------------------------------
.PHONY: version
version:
	@echo ${VERSION_FULL}

# ---------------------------------------------------------------------------
# Build the documentation.
# ---------------------------------------------------------------------------
.PHONY: docs
docs:
	mvn $(USE_MAVEN_SETTINGS) -B -f java install -P docs -pl docs -DskipTests -Doperator.version=$(VERSION_FULL)

# ---------------------------------------------------------------------------
# Start a local web server to serve the documentation.
# ---------------------------------------------------------------------------
.PHONY: serve-docs
serve-docs:
	@echo "Serving documentation on http://localhost:8080"
	cd $(BUILD_OUTPUT)/docs; \
	python -m SimpleHTTPServer 8080

# ---------------------------------------------------------------------------
# Release the Coherence Operator documentation and Helm chart to the
# gh-pages branch.
# ---------------------------------------------------------------------------
.PHONY: release-dashboards
release-dashboards:
	@echo "Releasing Dashboards $(VERSION_FULL)"
	mkdir -p $(BUILD_OUTPUT)/dashboards/$(VERSION_FULL) || true
	tar -czvf $(BUILD_OUTPUT)/dashboards/$(VERSION_FULL)/coherence-dashboards.tar.gz  dashboards/
	kubectl create configmap coherence-grafana-dashboards --from-file=dashboards/grafana \
		--dry-run -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION_FULL)/coherence-grafana-dashboards.yaml
	kubectl create configmap coherence-grafana-dashboards --from-file=dashboards/grafana-legacy \
		--dry-run -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION_FULL)/coherence-grafana-legacy-dashboards.yaml
	kubectl create configmap coherence-kibana-dashboards --from-file=dashboards/kibana \
		--dry-run -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION_FULL)/coherence-kibana-dashboards.yaml
	mkdir -p dashboards || true
	mv $(BUILD_OUTPUT)/dashboards/$(VERSION_FULL)/ dashboards/

# ---------------------------------------------------------------------------
# Release the Coherence Operator documentation and Helm chart to the
# gh-pages branch.
# ---------------------------------------------------------------------------
.PHONY: release-ghpages
release-ghpages: helm-chart docs
	@echo "Releasing Dashboards $(VERSION_FULL)"
	mkdir -p $(BUILD_OUTPUT)/dashboards/$(VERSION_FULL) || true
	tar -czvf $(BUILD_OUTPUT)/dashboards/$(VERSION_FULL)/coherence-dashboards.tar.gz  dashboards/
	kubectl create configmap coherence-grafana-dashboards --from-file=dashboards/grafana \
		--dry-run -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION_FULL)/coherence-grafana-dashboards.yaml
	kubectl create configmap coherence-grafana-dashboards --from-file=dashboards/grafana-legacy \
		--dry-run -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION_FULL)/coherence-grafana-legacy-dashboards.yaml
	kubectl create configmap coherence-kibana-dashboards --from-file=dashboards/kibana \
		--dry-run -o yaml > $(BUILD_OUTPUT)/dashboards/$(VERSION_FULL)/coherence-kibana-dashboards.yaml
	cp hack/docs-unstable-index.sh $(BUILD_OUTPUT)/docs-unstable-index.sh
	git stash save --keep-index --include-untracked || true
	git stash drop || true
	git checkout gh-pages
	git pull
	mkdir -p dashboards || true
	mv $(BUILD_OUTPUT)/dashboards/$(VERSION_FULL)/ dashboards/
	git add dashboards/$(VERSION_FULL)/*
	@echo "Releasing Helm chart $(VERSION_FULL)"
ifeq (true, $(PRE_RELEASE))
	mkdir -p docs-unstable || true
	rm -rf docs-unstable/$(VERSION_FULL)/ || true
	mv $(BUILD_OUTPUT)/docs/ docs-unstable/$(VERSION_FULL)/
	sh $(BUILD_OUTPUT)/docs-unstable-index.sh
	ls -ls docs-unstable

	mkdir -p charts-unstable || true
	cp $(CHART_DIR)/coherence-operator-$(VERSION_FULL).tgz charts-unstable/
	helm repo index charts-unstable --url https://oracle.github.io/coherence-operator/charts-unstable
	ls -ls charts-unstable

	git status
	git add docs-unstable/*
	git add charts-unstable/*
else
	mkdir docs/$(VERSION_FULL) || true
	rm -rf docs/$(VERSION_FULL)/ || true
	mv $(BUILD_OUTPUT)/docs/ docs/$(VERSION_FULL)/
	ls -ls docs

	mkdir -p charts || true
	cp $(CHART_DIR)/coherence-operator-$(VERSION_FULL).tgz charts/
	helm repo index charts --url https://oracle.github.io/coherence-operator/charts
	ls -ls charts

	git status
	git add docs/*
	git add charts/*
endif
	git clean -d -f
	git status
	git commit -m "adding Coherence Operator docs and helm chart version: $(VERSION_FULL)"
	git log -1
ifeq (true, $(RELEASE_DRY_RUN))
	@echo "release dry-run - would have pushed docs and Helm chart $(VERSION_FULL) to gh-pages"
else
	git push origin gh-pages
endif


# ---------------------------------------------------------------------------
# Tag Git for the release.
# ---------------------------------------------------------------------------
.PHONY: release-tag
release-tag:
ifeq (true, $(RELEASE_DRY_RUN))
	@echo "release dry-run - would have created release tag v$(VERSION_FULL)"
else
	@echo "creating release tag v$(VERSION_FULL)"
	git push origin :refs/tags/v$(VERSION_FULL)
	git tag -f -a -m "built $(VERSION_FULL)" v$(VERSION_FULL)
	git push origin --tags
endif

# ---------------------------------------------------------------------------
# Release the Coherence Operator.
# ---------------------------------------------------------------------------
.PHONY: release
release:

ifeq (true, $(RELEASE_DRY_RUN))
release: build-all-images release-tag release-ghpages
	@echo "release dry-run: would have pushed images"
else
release: build-all-images release-tag release-ghpages push-release-images
endif


# ---------------------------------------------------------------------------
# List all of the targets in the Makefile
# ---------------------------------------------------------------------------
.PHONY: list
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
