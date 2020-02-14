# ---------------------------------------------------------------------------
# Copyright (c) 2019, 2020, Oracle and/or its affiliates. All rights reserved.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
# ---------------------------------------------------------------------------
# This is the Makefile to build the Coherence Kubernetes Operator.
# ---------------------------------------------------------------------------

# The version of the Operator being build - this should be a valid SemVer format
VERSION ?= 2.1.0

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

# Capture the Git commit to add to the build information
GITCOMMIT       ?= $(shell git rev-list -1 HEAD)
GITREPO         := https://github.com/oracle/coherence-operator.git

CURRDIR         := $(shell pwd)

ARCH            ?= amd64
OS              ?= linux
UNAME_S         := $(shell uname -s)
GOPROXY         ?= https://proxy.golang.org

# Set the location of the Operator SDK executable
UNAME_S      = $(shell uname -s)
UNAME_M      = $(shell uname -m)
OPERATOR_SDK = $(CURRDIR)/sdk/$(UNAME_S)-$(UNAME_M)/operator-sdk
OP_CHMOD     := $(shell chmod +x $(OPERATOR_SDK))

# The image prefix to use for Coherence images
COHERENCE_IMAGE_PREFIX ?= container-registry.oracle.com/middleware/
# The Coherence image name to inject into the Helm chart
HELM_COHERENCE_IMAGE   ?= container-registry.oracle.com/middleware/coherence:12.2.1.4.0

# One may need to define RELEASE_IMAGE_PREFIX in the environment.
RELEASE_IMAGE_PREFIX ?= "iad.ocir.io/odx-stateservice/test/$(USER)/"
OPERATOR_IMAGE       := $(RELEASE_IMAGE_PREFIX)oracle/coherence-operator:$(VERSION_FULL)
UTILS_IMAGE          ?= $(RELEASE_IMAGE_PREFIX)oracle/coherence-operator:$(VERSION_FULL)-utils
TEST_USER_IMAGE      := $(RELEASE_IMAGE_PREFIX)oracle/operator-test-image:$(VERSION_FULL)

RELEASE_DRY_RUN  ?= true
PRE_RELEASE      ?= true

# The version of the Prometheus Operator chart that is used as a sub-chart in the
# Coherence Operator chart
PROMETHEUS_HELMCHART_VERSION ?= 5.7.0

# The EFK images to use
ELASTICSEARCH_IMAGE ?= docker.elastic.co/elasticsearch/elasticsearch-oss:6.6.0
FLUENTD_IMAGE       ?= fluent/fluentd-kubernetes-daemonset:v1.3.3-debian-elasticsearch-1.3
KIBANA_IMAGE        ?= docker.elastic.co/kibana/kibana-oss:6.6.0

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

# restart local storage for persistence
LOCAL_STORAGE_RESTART ?= false

IMAGE_PULL_SECRETS ?=
IMAGE_PULL_POLICY  ?=

# Env variable used by the kubectl test framework to locate the kubectl binary
TEST_ASSET_KUBECTL ?= $(shell which kubectl)

override BUILD_OUTPUT  := ./build/_output
override BUILD_PROPS   := $(BUILD_OUTPUT)/build.properties
override CHART_DIR     := $(BUILD_OUTPUT)/helm-charts
override CRD_DIR       := deploy/crds
override TEST_LOGS_DIR := $(BUILD_OUTPUT)/test-logs

ifeq (, $(shell which ginkgo))
GO_TEST_CMD = go
else
GO_TEST_CMD = ginkgo
endif

GOS        = $(shell find pkg -type f -name "*.go" ! -name "*_test.go")
COPYGOS    = $(shell find cmd/copyartifacts -type f -name "*.go" ! -name "*_test.go")
OPTESTGOS  = $(shell find cmd/optest -type f -name "*.go" ! -name "*_test.go")
UTILGOS    = $(shell find cmd/utilsinit -type f -name "*.go" ! -name "*_test.go")
COH_CHARTS = $(shell find helm-charts/coherence -type f)
COP_CHARTS = $(shell find helm-charts/coherence-operator -type f)
DEPLOYS    = $(shell find deploy -type f -name "*.yaml")
CRDS       = $(shell find deploy/crds -name "*_crd.yaml")

TEST_MANIFEST_DIR         := $(BUILD_OUTPUT)/manifest
TEST_MANIFEST_FILE        := test-manifest.yaml
TEST_LOCAL_MANIFEST_FILE  := local-manifest.yaml
TEST_GLOBAL_MANIFEST_FILE := global-manifest.yaml
TEST_MANIFEST_VALUES      ?= deploy/test-values.yaml
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
		if [[ -f $${filename} ]]; then \
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
	# create build.properties
	rm -f $(BUILD_PROPS)
	printf "HELM_COHERENCE_IMAGE=$(HELM_COHERENCE_IMAGE)\n\
	UTILS_IMAGE=$(UTILS_IMAGE)\n\
	OPERATOR_IMAGE=$(OPERATOR_IMAGE)\n\
	ELASTICSEARCH_IMAGE=$(ELASTICSEARCH_IMAGE)\n\
	FLUENTD_IMAGE=$(FLUENTD_IMAGE)\n\
	KIBANA_IMAGE=$(KIBANA_IMAGE)\n\
	PROMETHEUS_HELMCHART_VERSION=$(PROMETHEUS_HELMCHART_VERSION)\n\
	VERSION_FULL=$(VERSION_FULL)\n\
	VERSION=$(VERSION)\n" > $(BUILD_PROPS)

# ---------------------------------------------------------------------------
# Builds the project, helm charts and Docker image
# ---------------------------------------------------------------------------
.PHONY: build-operator
build-operator: $(BUILD_OUTPUT)/bin/operator

# ---------------------------------------------------------------------------
# Internal make step that builds the Operator Docker image and Helm charts
# ---------------------------------------------------------------------------
$(BUILD_OUTPUT)/bin/operator: export CGO_ENABLED = 0
$(BUILD_OUTPUT)/bin/operator: export GOARCH = $(ARCH)
$(BUILD_OUTPUT)/bin/operator: export GOOS = $(OS)
$(BUILD_OUTPUT)/bin/operator: export GO111MODULE = on
$(BUILD_OUTPUT)/bin/operator: $(GOS) $(DEPLOYS) $(CHART_DIR)/coherence $(CHART_DIR)/coherence-operator-$(VERSION_FULL).tgz
	@echo "Building: $(OPERATOR_IMAGE)"
	@echo "Running Operator SDK build"
	BUILD_INFO="$(VERSION_FULL)|$(GITCOMMIT)|$$(date -u | tr ' ' '.')"; \
	$(OPERATOR_SDK) build $(OPERATOR_IMAGE) --verbose --go-build-args "-o $(BUILD_OUTPUT)/bin/operator -ldflags -X=main.BuildInfo=$${BUILD_INFO}"

# ---------------------------------------------------------------------------
# Internal make step that builds the Operator copy artifacts utility
# ---------------------------------------------------------------------------
.PHONY: build-copy-artifacts
build-copy-artifacts: $(BUILD_OUTPUT)/bin/copy

$(BUILD_OUTPUT)/bin/copy: export CGO_ENABLED = 0
$(BUILD_OUTPUT)/bin/copy: export GOARCH = $(ARCH)
$(BUILD_OUTPUT)/bin/copy: export GOOS = $(OS)
$(BUILD_OUTPUT)/bin/copy: export GO111MODULE = on
$(BUILD_OUTPUT)/bin/copy: $(GOS) $(DEPLOYS) $(COPYGOS)
	go build -o $(BUILD_OUTPUT)/bin/copy ./cmd/copyartifacts

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
# Internal make step to build the Coherence Helm chart that is packaged
# inside the Operator Docker image.
# ---------------------------------------------------------------------------
$(CHART_DIR)/coherence: $(COH_CHARTS) $(BUILD_PROPS)
	rm -rf $(CHART_DIR)/coherence
	# Copy the Helm charts from their source location to the distribution folder
	cp -R ./helm-charts/coherence $(CHART_DIR)
	$(call replaceprop,coherence/Chart.yaml coherence/values.yaml)
	helm lint $(CHART_DIR)/coherence

# ---------------------------------------------------------------------------
# Executes the Go unit tests that do not require a k8s cluster
# ---------------------------------------------------------------------------
.PHONY: test-operator
test-operator: export CGO_ENABLED = 0
test-operator: build-operator
	@echo "Running operator tests"
	go test $(GO_TEST_FLAGS) -v ./cmd/... ./pkg/... ./pkg/helm_test/... \
	2>&1 | tee $(TEST_LOGS_DIR)/operator-test.out
	go run ./cmd/testreports/ -fail -suite-name-prefix=operator-test/ \
	    -input $(TEST_LOGS_DIR)/operator-test.out \
	    -output $(TEST_LOGS_DIR)/operator-test.xml

# ---------------------------------------------------------------------------
# Executes the Go end-to-end tests that require a k8s cluster using
# a LOCAL operator instance (i.e. the operator is not deployed to k8s).
# These tests will use whichever k8s cluster the local environment
# is pointing to.
# These tests require the Operator CRDs and will install them before
# tests start and remove them afterwards.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
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
e2e-local-test: build-operator reset-namespace create-ssl-secrets operator-manifest uninstall-crds
	@echo "executing end-to-end tests"
	$(OPERATOR_SDK) test local ./test/e2e/local --namespace $(TEST_NAMESPACE) --up-local \
		--verbose --debug  --go-test-flags "$(GO_TEST_FLAGS_E2E)" \
		--local-operator-flags "--watches-file=local-watches.yaml --crd-files=$(CRD_DIR)" \
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
# make debug-e2e-test GO_TEST_FLAGS='-run=^TestStatusHA/HttpStatusHAHandler$$'
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
debug-e2e-local-test:
	$(OPERATOR_SDK) test local ./test/e2e/local \
	    --namespace $(TEST_NAMESPACE) \
		--verbose --debug  --go-test-flags \
		"$(GO_TEST_FLAGS_E2E)" --no-setup


# ---------------------------------------------------------------------------
# Executes the Go script tests that require a k8s cluster using
# a LOCAL operator instance (i.e. the operator is not deployed to k8s).
# These tests will use whichever k8s cluster the local environment
# is pointing to.
# These tests require the Operator CRDs and will install them before
# tests start and remove them afterwards.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
# ---------------------------------------------------------------------------
.PHONY: script-test
script-test: export CGO_ENABLED = 0
script-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
script-test: export TEST_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_LOCAL_MANIFEST_FILE)
script-test: export TEST_GLOBAL_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_GLOBAL_MANIFEST_FILE)
script-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
script-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
script-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
script-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
script-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
script-test: export TEST_ASSET_KUBECTL := $(TEST_ASSET_KUBECTL)
script-test: build-operator reset-namespace create-ssl-secrets operator-manifest uninstall-crds
	@echo "executing end-to-end tests"
	$(OPERATOR_SDK) test local ./test/e2e/script --namespace $(TEST_NAMESPACE) --up-local \
		--verbose --debug  --go-test-flags "$(GO_TEST_FLAGS_E2E)" \
		--local-operator-flags "--watches-file=local-watches.yaml --crd-files=$(CRD_DIR)" \
		--namespaced-manifest=$(TEST_MANIFEST) \
		--global-manifest=$(TEST_GLOBAL_MANIFEST) \
		 2>&1 | tee $(TEST_LOGS_DIR)/operator-script-test.out
	$(MAKE) delete-namespace
	go run ./cmd/testreports/ -fail -suite-name-prefix=script-test/ \
	    -input $(TEST_LOGS_DIR)/operator-script-test.out \
	    -output $(TEST_LOGS_DIR)/operator-script-test.xml


# ---------------------------------------------------------------------------
# Executes the Go end-to-end tests that require a k8s cluster using
# a DEPLOYED operator instance (i.e. the operator Docker image is
# deployed to k8s). These tests will use whichever k8s cluster the
# local environment is pointing to.
# These tests require the Operator CRDs and will install them before
# tests start and remove them afterwards.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
# ---------------------------------------------------------------------------
.PHONY: e2e-test
e2e-test: export CGO_ENABLED = 0
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
e2e-test: build-operator reset-namespace create-ssl-secrets operator-manifest uninstall-crds
	@echo "executing end-to-end tests"
	$(OPERATOR_SDK) test local ./test/e2e/remote --namespace $(TEST_NAMESPACE) \
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
# make debug-e2e-test GO_TEST_FLAGS='-run=^TestScaling/DownSafeScaling$$'
#
# ---------------------------------------------------------------------------
.PHONY: debug-e2e-test
debug-e2e-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
debug-e2e-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
debug-e2e-test: export TEST_NAMESPACE := $(TEST_NAMESPACE)
debug-e2e-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
debug-e2e-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
debug-e2e-test: export TEST_STORAGE_CLASS := $(TEST_STORAGE_CLASS)
debug-e2e-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
debug-e2e-test: export TEST_ASSET_KUBECTL := $(TEST_ASSET_KUBECTL)
debug-e2e-test:
	$(OPERATOR_SDK) test local ./test/e2e/remote --namespace $(TEST_NAMESPACE) \
		--verbose --debug  --go-test-flags "$(GO_TEST_FLAGS_E2E)" \
		--no-setup  2>&1 | tee $(TEST_LOGS_DIR)/operator-e2e-test.out


# ---------------------------------------------------------------------------
# Executes the Go end-to-end Operator Helm chart tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# These tests require the Operator CRDs and will install them before tests start
# and remove them afterwards.
# Note that the namespace will be created if it does not exist.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
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
helm-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
helm-test: build-operator reset-namespace create-ssl-secrets install-crds
	@echo "executing Operator Helm Chart end-to-end tests"
	$(OPERATOR_SDK) test local ./test/e2e/helm --namespace $(TEST_NAMESPACE) \
		--verbose --debug  --go-test-flags "$(GO_TEST_FLAGS_E2E)" \
		--no-setup  2>&1 | tee $(TEST_LOGS_DIR)/operator-e2e-helm-test.out
	$(MAKE) uninstall-crds
	$(MAKE) delete-namespace
	go run ./cmd/testreports/ -fail -suite-name-prefix=e2e-helm-test/ \
	    -input $(TEST_LOGS_DIR)/operator-e2e-helm-test.out \
	    -output $(TEST_LOGS_DIR)/operator-e2e-helm-test.xml

# ---------------------------------------------------------------------------
# Install CRDs into Kubernetes.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
# ---------------------------------------------------------------------------
.PHONY: install-crds
install-crds: uninstall-crds
	@echo "Installing CRDs"
	for i in $(CRDS); do \
		kubectl create -f $${i}; \
	done

# ---------------------------------------------------------------------------
# Uninstall CRDs from Kubernetes.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
# ---------------------------------------------------------------------------
.PHONY: uninstall-crds
uninstall-crds:
	@echo "Removing CRDs"
	for i in $(CRDS); do \
		(kubectl delete -f $${i} & ); \
	done
	kubectl patch crd coherenceinternals.coherence.oracle.com -p '{"metadata":{"finalizers":[]}}' --type=merge || true


# ---------------------------------------------------------------------------
# This step will run the Operator SDK code generators.
# These commands will generate the CRD files from the API structs and will
# also generate the Go DeepCopy code for the API structs.
# This step would require running if any of the structs in the files under
# the pkg/apis directory have been changed.
# ---------------------------------------------------------------------------
.PHONY: generate
generate:
	@echo "Generating deep copy code"
	$(OPERATOR_SDK) generate k8s
	@echo "Generating CRDs"
	$(OPERATOR_SDK) generate crds

# ---------------------------------------------------------------------------
# Clean-up all of the build artifacts
# ---------------------------------------------------------------------------
.PHONY: clean
clean:
	rm -rf build/_output
	mvn -f java clean
	mvn -f examples clean

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
operator-manifest: export TEST_MANIFEST_VALUES := $(TEST_MANIFEST_VALUES)
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
reset-namespace: delete-namespace
ifeq ($(CREATE_TEST_NAMESPACE),true)
	@echo "Creating test namespace $(TEST_NAMESPACE)"
	kubectl create namespace $(TEST_NAMESPACE)
endif

# ---------------------------------------------------------------------------
# Delete the test namespace
# ---------------------------------------------------------------------------
.PHONY: delete-namespace
delete-namespace:
ifeq ($(CREATE_TEST_NAMESPACE),true)
	@echo "Deleting test namespace $(TEST_NAMESPACE)"
	kubectl delete namespace $(TEST_NAMESPACE) && echo "deleted namespace" || true
endif

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
	mvn -f java package -DskipTests
	mvn -f examples package -DskipTests

# ---------------------------------------------------------------------------
# Build and test the Java artifacts
# ---------------------------------------------------------------------------
.PHONY: test-mvn
test-mvn: build-mvn
	mvn -f java verify

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
build-utils-image: build-mvn build-copy-artifacts build-utils-init build-op-test
	cp $(BUILD_OUTPUT)/bin/copy java/coherence-utils/target/docker/copy
	cp $(BUILD_OUTPUT)/bin/op-test java/coherence-utils/target/docker/op-test
	cp $(BUILD_OUTPUT)/bin/utils-init java/coherence-utils/target/docker/utils-init
	docker build -t $(UTILS_IMAGE) java/coherence-utils/target/docker

# ---------------------------------------------------------------------------
# Push the Operator Utils Docker image
# ---------------------------------------------------------------------------
.PHONY: push-utils-image
push-utils-image:
	@echo "Pushing $(UTILS_IMAGE)"
	docker push $(UTILS_IMAGE)

# ---------------------------------------------------------------------------
# Build the Operator Test Docker image
# ---------------------------------------------------------------------------
.PHONY: build-test-image
build-test-image: build-mvn
	docker build -t $(TEST_USER_IMAGE) java/operator-test/target/docker


# ---------------------------------------------------------------------------
# Push the Operator Utils Docker image
# ---------------------------------------------------------------------------
.PHONY: push-test-image
push-test-image:
	@echo "Pushing $(TEST_USER_IMAGE)"
	docker push $(TEST_USER_IMAGE)

# ---------------------------------------------------------------------------
# Build all of the Docker images
# ---------------------------------------------------------------------------
.PHONY: build-all-images
build-all-images: build-operator build-utils-image build-test-image

# ---------------------------------------------------------------------------
# Push all of the Docker images
# ---------------------------------------------------------------------------
.PHONY: push-all-images
push-all-images: push-operator-image push-utils-image push-test-image

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
# Run the Operator in locally.
#
# To exit out of the local Operator you can use ctrl-c or ctrl-z but
# sometimes this leaves orphaned processes on the local machine so
# ensure these are killed run "make debug-stop"
# ---------------------------------------------------------------------------
.PHONY: run
run: $(CHART_DIR)/coherence reset-namespace create-ssl-secrets
	$(OPERATOR_SDK) run --local --namespace=$(TEST_NAMESPACE) \
	--operator-flags="--watches-file=local-watches.yaml --crd-files=$(CRD_DIR)" \
	2>&1 | tee $(TEST_LOGS_DIR)/operator-debug.out

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
run-debug: $(CHART_DIR)/coherence reset-namespace create-ssl-secrets
	$(OPERATOR_SDK) run --local --namespace=$(TEST_NAMESPACE) \
	--operator-flags="--watches-file=local-watches.yaml --crd-files=$(CRD_DIR)" \
	--enable-delve \
	2>&1 | tee $(TEST_LOGS_DIR)/operator-debug.out

# ---------------------------------------------------------------------------
# Kill any locally running Operator
# ---------------------------------------------------------------------------
.PHONY: stop
stop:
	./hack/kill-local.sh


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
# Delete all of the CoherenceClusters from the test namespace.
# This step will patch the finalizers on all CoherenceInternal
# resources to ensure that they are properly deleted.
# ---------------------------------------------------------------------------
.PHONY: delete-coherence-clusters
delete-coherence-clusters:
	for i in $$(kubectl -n  $(TEST_NAMESPACE) get coherencecluster -o name); do \
		kubectl -n $(TEST_NAMESPACE) delete $${i}; \
	done
	for i in $$(kubectl -n  $(TEST_NAMESPACE) get coherenceinternal -o name); do \
		kubectl -n $(TEST_NAMESPACE) patch $${i}  -p '{"metadata":{"finalizers": []}}' --type=merge; \
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
	$(BUILD_OUTPUT)/bin/golangci-lint run -v --timeout=5m  ./pkg/... ./cmd/...


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
	  -X /Dockerfile \
	  -X docs/ \
	  -X etc/copyright.txt \
	  -X etc/intellij-codestyle.xml \
	  -X go.mod \
	  -X go.sum \
	  -X HEADER.txt \
	  -X helm-charts/coherence/templates/NOTES.txt \
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
	  -X sdk/ \
	  -X .sh \
	  -X temp/ \
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
code-review: golangci copyright
	mvn -f java package -DskipTests -P checkstyle
	mvn -f examples package -DskipTests -P checkstyle

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
	mvn -f java install -P docs -pl docs -DskipTests -Doperator.version=$(VERSION_FULL)

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
.PHONY: release-ghpages
release-ghpages: helm-chart docs
	@echo "Releasing Helm chart $(VERSION_FULL)"
	cp hack/docs-unstable-index.sh $(BUILD_OUTPUT)/docs-unstable-index.sh
	git reset
	git checkout gh-pages
	git pull
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
# Initialise a Kind k8s cluster
# ---------------------------------------------------------------------------
.PHONY: kind-up
kind-up: kind-create-cluster kind-upload-coherence kind-upload-operator kind-upload-efk

.PHONY: kind-create-cluster
export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
export UTILS_IMAGE := $(UTILS_IMAGE)
export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
kind-create-cluster:
	./hack/start-kind.sh

.PHONY: kind-upload-coherence
export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
export UTILS_IMAGE := $(UTILS_IMAGE)
export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
kind-upload-coherence:
	docker pull "${HELM_COHERENCE_IMAGE}"
	kind load docker-image --name operator-test --nodes operator-test-worker,operator-test-worker2,operator-test-worker3 "${HELM_COHERENCE_IMAGE}"

.PHONY: kind-upload-operator
export HELM_COHERENCE_IMAGE := $(HELM_COHERENCE_IMAGE)
export OPERATOR_IMAGE := $(OPERATOR_IMAGE)
export UTILS_IMAGE := $(UTILS_IMAGE)
export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
kind-upload-operator:
	kind load docker-image --name operator-test --nodes operator-test-worker,operator-test-worker2,operator-test-worker3 "${OPERATOR_IMAGE}"
	kind load docker-image --name operator-test --nodes operator-test-worker,operator-test-worker2,operator-test-worker3 "${UTILS_IMAGE}"
	kind load docker-image --name operator-test --nodes operator-test-worker,operator-test-worker2,operator-test-worker3 "${TEST_USER_IMAGE}"

.PHONY: kind-upload-efk
export ELASTICSEARCH_IMAGE := $(ELASTICSEARCH_IMAGE)
export FLUENTD_IMAGE := $(FLUENTD_IMAGE)
export KIBANA_IMAGE := $(KIBANA_IMAGE)
kind-upload-efk:
#	docker pull ${ELASTICSEARCH_IMAGE}
#	kind load docker-image --name operator-test --nodes operator-test-worker,operator-test-worker2,operator-test-worker3 "${ELASTICSEARCH_IMAGE}"
	docker pull ${FLUENTD_IMAGE}
	kind load docker-image --name operator-test --nodes operator-test-worker,operator-test-worker2,operator-test-worker3 "${FLUENTD_IMAGE}"
#	docker pull ${KIBANA_IMAGE}
#	kind load docker-image --name operator-test --nodes operator-test-worker,operator-test-worker2,operator-test-worker3 "${KIBANA_IMAGE}"

# ---------------------------------------------------------------------------
# Clean-up a Kind k8s cluster
# ---------------------------------------------------------------------------
.PHONY: kind-init
kind-down:
	kind delete cluster --name operator-test || true
	unset KUBECONFIG || true

# ---------------------------------------------------------------------------
# List all of the targets in the Makefile
# ---------------------------------------------------------------------------
.PHONY: list
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
