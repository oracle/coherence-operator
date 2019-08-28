# ---------------------------------------------------------------------------
# Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
# ---------------------------------------------------------------------------
# This is the Makefile to build the Coherence Kubernetes Operator.
# ---------------------------------------------------------------------------

# The version of the Operator being build - this should be a valid SemVer format
VERSION ?= 2.0.0

# An optional version suffix. For a full release this should be set to blank,
# for an interim release it should be set to a value to identify that release.
# For example if building the third release candidate this value might be
# set to VERSION_SUFFIX=RC3
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

ARCH            ?= amd64
OS              ?= linux
UNAME_S         := $(shell uname -s)

# The image prefix to use for Coherence images
COHERENCE_IMAGE_PREFIX ?= iad.ocir.io/odx-stateservice/test/
# The Coherence image name to inject into the Helm chart
HELM_COHERENCE_IMAGE   ?= $(COHERENCE_IMAGE_PREFIX)coherence:12.2.1.4.0-b74630

# One may need to define RELEASE_IMAGE_PREFIX in the environment.
OPERATOR_IMAGE   := $(RELEASE_IMAGE_PREFIX)oracle/coherence-operator:$(VERSION_FULL)
UTILS_IMAGE      ?= $(RELEASE_IMAGE_PREFIX)oracle/coherence-operator:$(VERSION_FULL)-utils
TEST_USER_IMAGE  := $(RELEASE_IMAGE_PREFIX)oracle/operator-test-image:$(VERSION_FULL)

# The version of the Prometheus Operator chart that is used as a sub-chart in the
# Coherence Operator chart
PROMETHEUS_HELMCHART_VERSION ?= 5.7.0

# Extra arguments to pass to the go test command for the various test steps.
# For example, when running make e2e-test we can run just a single test such
# as the zone test using the go test -run=regex argument like this
#   make e2e-test GO_TEST_FLAGS='-run=^TestZone$$'
GO_TEST_FLAGS     ?=
GO_TEST_FLAGS_E2E := -timeout=100m $(GO_TEST_FLAGS)

# This is the Coherence image that will be used in the Go tests.
# Changing this variable will allow test builds to be run against differet Coherence versions
TEST_COHERENCE_IMAGE ?= $(HELM_COHERENCE_IMAGE)

# default as in test/e2e/helper/proj_helpers.go
TEST_NAMESPACE ?= operator-test

CREATE_TEST_NAMESPACE ?= true

IMAGE_PULL_SECRETS ?=
IMAGE_PULL_POLICY  ?=

override BUILD_OUTPUT  := ./build/_output
override BUILD_PROPS   := $(BUILD_OUTPUT)/build.properties
override CHART_DIR     := $(BUILD_OUTPUT)/helm-charts
override TEST_LOGS_DIR := $(BUILD_OUTPUT)/test-logs

ifeq (, $(shell which ginkgo))
GO_TEST_CMD = go
else
GO_TEST_CMD = ginkgo
endif

GOS=$(shell find pkg -type f -name "*.go" ! -name "*_test.go")
COH_CHARTS=$(shell find helm-charts/coherence -type f)
COP_CHARTS=$(shell find helm-charts/coherence-operator -type f)
DEPLOYS=$(shell find deploy -type f -name "*.yaml")
CRDS=$(shell find deploy/crds -name "*_crd.yaml")

TEST_MANIFEST_DIR         := $(BUILD_OUTPUT)/manifest
TEST_MANIFEST_FILE        := test-manifest.yaml
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

.PHONY: all build test-operator e2e-local-test e2e-test install-crds uninstall-crds generate push-operator-image clean operator-manifest reset-namespace create-ssl-secrets

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
	PROMETHEUS_HELMCHART_VERSION=$(PROMETHEUS_HELMCHART_VERSION)\n\
	VERSION_FULL=$(VERSION_FULL)\n\
	VERSION_SUFFIX=$(VERSION_SUFFIX)\n\
	VERSION=$(VERSION)\n" > $(BUILD_PROPS)

# ---------------------------------------------------------------------------
# Builds the project, helm charts and Docker image
# ---------------------------------------------------------------------------
build-operator: $(BUILD_OUTPUT)/bin/operator

# ---------------------------------------------------------------------------
# Internal make step that builds the Operator Docker image and Helm charts
# ---------------------------------------------------------------------------
$(BUILD_OUTPUT)/bin/operator: export CGO_ENABLED = 0
$(BUILD_OUTPUT)/bin/operator: export GOARCH = $(ARCH)
$(BUILD_OUTPUT)/bin/operator: export GOOS = $(OS)
$(BUILD_OUTPUT)/bin/operator: export GO111MODULE = on
$(BUILD_OUTPUT)/bin/operator: $(GOS) $(DEPLOYS) $(CHART_DIR)/coherence-$(VERSION_FULL).tar.gz $(CHART_DIR)/coherence-operator-$(VERSION_FULL).tar.gz
	@echo "Building: $(OPERATOR_IMAGE)"
	@echo "Running Operator SDK build"
	BUILD_INFO="$(VERSION_FULL)|$(GITCOMMIT)|$$(date -u | tr ' ' '.')"; \
	operator-sdk build $(OPERATOR_IMAGE) --verbose --go-build-args "-o $(BUILD_OUTPUT)/bin/operator -ldflags -X=main.BuildInfo=$${BUILD_INFO}"

# ---------------------------------------------------------------------------
# Build the COperator Helm chart
# ---------------------------------------------------------------------------
$(CHART_DIR)/coherence-operator-$(VERSION_FULL).tar.gz: $(COP_CHARTS) $(BUILD_PROPS)
	# Copy the Helm charts from their source location to the distribution folder
	cp -R ./helm-charts/coherence-operator $(CHART_DIR)

	$(call replaceprop,coherence-operator/Chart.yaml coherence-operator/values.yaml coherence-operator/requirements.yaml coherence-operator/templates/deployment.yaml)

	# For each Helm chart folder package the chart into a .tar.gz
	# Package the chart into a .tr.gz - we don't use helm package as the version might not be SEMVER
	echo "Creating Helm chart package $(CHART_DIR)/coherence-operator"
	helm lint $(CHART_DIR)/coherence-operator

# ---------------------------------------------------------------------------
# Build the Operator Helm chart and package it into a tar.gz
# ---------------------------------------------------------------------------
helm-chart: $(COP_CHARTS) $(BUILD_PROPS) $(CHART_DIR)/coherence-operator-$(VERSION_FULL).tar.gz
	tar -czf $(CHART_DIR)/coherence-operator-$(VERSION_FULL).tar.gz $(CHART_DIR)/coherence-operator

# ---------------------------------------------------------------------------
# Internal make step to build the Coherence Helm chart that is packaged
# inside the Operator Docker image.
# ---------------------------------------------------------------------------
$(CHART_DIR)/coherence-$(VERSION_FULL).tar.gz: $(COH_CHARTS) $(BUILD_PROPS)
	# Copy the Helm charts from their source location to the distribution folder
	cp -R ./helm-charts/coherence $(CHART_DIR)

	$(call replaceprop,coherence/Chart.yaml coherence/values.yaml)

	# For each Helm chart folder package the chart into a .tar.gz
	# Package the chart into a .tr.gz - we don't use helm package as the version might not be SEMVER
	echo "Creating Helm chart package $(CHART_DIR)/coherence"
	helm lint $(CHART_DIR)/coherence
	tar -czf $(CHART_DIR)/coherence-$(VERSION_FULL).tar.gz $(CHART_DIR)/coherence

# ---------------------------------------------------------------------------
# Executes the Go unit tests that do not require a k8s cluster
# ---------------------------------------------------------------------------
test-operator: export CGO_ENABLED = 0
test-operator: build-operator
	@echo "Running operator tests"
	$(GO_TEST_CMD) test $(GO_TEST_FLAGS) -v ./cmd/... ./pkg/...

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
e2e-local-test: export CGO_ENABLED = 0
e2e-local-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
e2e-local-test: export TEST_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_MANIFEST_FILE)
e2e-local-test: export TEST_GLOBAL_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_GLOBAL_MANIFEST_FILE)
e2e-local-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
e2e-local-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
e2e-local-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
e2e-local-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
e2e-local-test: build-operator reset-namespace create-ssl-secrets operator-manifest uninstall-crds
	@echo "executing end-to-end tests"
	operator-sdk test local ./test/e2e/local/... --namespace $(TEST_NAMESPACE) --up-local \
		--verbose --debug  --go-test-flags "$(GO_TEST_FLAGS_E2E)" \
		--local-operator-flags "--watches-file=local-watches.yaml" \
		--namespaced-manifest=$(TEST_MANIFEST) \
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
debug-e2e-local-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
debug-e2e-local-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
debug-e2e-local-test: export TEST_NAMESPACE := $(TEST_NAMESPACE)
debug-e2e-local-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
debug-e2e-local-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
debug-e2e-local-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
debug-e2e-local-test:
	operator-sdk test local ./test/e2e/local/... \
	    --namespace $(TEST_NAMESPACE) \
		--verbose --debug  --go-test-flags \
		"$(GO_TEST_FLAGS_E2E)" --no-setup


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
e2e-test: export CGO_ENABLED = 0
e2e-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
e2e-test: export TEST_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_MANIFEST_FILE)
e2e-test: export TEST_GLOBAL_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_GLOBAL_MANIFEST_FILE)
e2e-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
e2e-test: export TEST_NAMESPACE := $(TEST_NAMESPACE)
e2e-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
e2e-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
e2e-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
e2e-test: build-operator reset-namespace create-ssl-secrets operator-manifest uninstall-crds
	@echo "executing end-to-end tests"
	operator-sdk test local ./test/e2e/remote/... --namespace $(TEST_NAMESPACE) \
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
debug-e2e-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
debug-e2e-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
debug-e2e-test: export TEST_NAMESPACE := $(TEST_NAMESPACE)
debug-e2e-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
debug-e2e-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
debug-e2e-test: export GO_TEST_FLAGS_E2E := $(strip $(GO_TEST_FLAGS_E2E))
debug-e2e-test:
	operator-sdk test local ./test/e2e/remote/... --namespace $(TEST_NAMESPACE) \
		--verbose --debug  --go-test-flags "$(GO_TEST_FLAGS_E2E)" \
		--no-setup  2>&1 | tee $(TEST_LOGS_DIR)/operator-e2e-test.out


# ---------------------------------------------------------------------------
# Executes the Go end-to-end Operator Helm chart tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# These tests require the Operator CRDs and will install them before tests start
# and remove them afterwards.
# Note that the namespace will be created by Helm if it does not exist.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
# ---------------------------------------------------------------------------
helm-test: export CGO_ENABLED = 0
helm-test: export TEST_NAMESPACE := $(TEST_NAMESPACE)
helm-test: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
helm-test: export TEST_COHERENCE_IMAGE := $(TEST_COHERENCE_IMAGE)
helm-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
helm-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
helm-test: export TEST_IMAGE_PULL_POLICY := $(IMAGE_PULL_POLICY)
helm-test: build-operator reset-namespace create-ssl-secrets
	$(MAKE) install-crds
	@echo "executing Operator Helm Chart end-to-end tests"
	$(GO_TEST_CMD) test $(GO_TEST_FLAGS) -v ./test/e2e/helm/...
	$(MAKE) uninstall-crds
	$(MAKE) delete-namespace

# ---------------------------------------------------------------------------
# Install CRDs into Kubernetes.
# This step will use whatever Kubeconfig the current environment is
# configured to use.
# ---------------------------------------------------------------------------
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
generate:
	@echo "Generating deep copy code"
	operator-sdk generate k8s
	@echo "Generating Open API code and CRDs"
	operator-sdk generate openapi

# ---------------------------------------------------------------------------
# Clean-up all of the build artifacts
# ---------------------------------------------------------------------------
clean:
	rm -rf build/_output
	mvn -f java clean

# ---------------------------------------------------------------------------
# Create the k8s yaml manifest that will be used by the Operator SDK to
# install the Operator when running e2e tests.
# ---------------------------------------------------------------------------
operator-manifest: export TEST_NAMESPACE := $(TEST_NAMESPACE)
operator-manifest: export TEST_MANIFEST_DIR := $(TEST_MANIFEST_DIR)
operator-manifest: export TEST_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_MANIFEST_FILE)
operator-manifest: export TEST_GLOBAL_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_GLOBAL_MANIFEST_FILE)
operator-manifest: export TEST_MANIFEST_VALUES := $(TEST_MANIFEST_VALUES)
operator-manifest: $(CHART_DIR)/coherence-operator-$(VERSION_FULL).tar.gz
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
reset-namespace: delete-namespace
ifeq ($(CREATE_TEST_NAMESPACE),true)
	@echo "Creating test namespace $(TEST_NAMESPACE)"
	kubectl create namespace $(TEST_NAMESPACE)
endif

# ---------------------------------------------------------------------------
# Delete the test namespace
# ---------------------------------------------------------------------------
delete-namespace:
ifeq ($(CREATE_TEST_NAMESPACE),true)
	@echo "Deleting test namespace $(TEST_NAMESPACE)"
	kubectl delete namespace $(TEST_NAMESPACE) && echo "deleted namespace" || true
endif

# ---------------------------------------------------------------------------
# Create the k8s secret to use in SSL/TLS testing.
# ---------------------------------------------------------------------------
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
build-mvn:
	mvn -f java package -DskipTests

# ---------------------------------------------------------------------------
# Build and test the Java artifacts
# ---------------------------------------------------------------------------
test-mvn: build-mvn
	mvn -f java verify

# ---------------------------------------------------------------------------
# Run all unit tests (both Go and Java)
# ---------------------------------------------------------------------------
test-all: test-mvn test-operator

# ---------------------------------------------------------------------------
# Push the Operator Docker image
# ---------------------------------------------------------------------------
push-operator-image: build-operator
	@echo "Pushing $(OPERATOR_IMAGE)"
	docker push $(OPERATOR_IMAGE)

# ---------------------------------------------------------------------------
# Build the Operator Utils Docker image
# ---------------------------------------------------------------------------
build-utils-image: export UTILS_IMAGE := $(UTILS_IMAGE)
build-utils-image: build-mvn
	docker build -t $(UTILS_IMAGE) java/coherence-utils/target/docker

# ---------------------------------------------------------------------------
# Push the Operator Utils Docker image
# ---------------------------------------------------------------------------
push-utils-image: export UTILS_IMAGE := $(UTILS_IMAGE)
push-utils-image:
	@echo "Pushing $(UTILS_IMAGE)"
	docker push $(UTILS_IMAGE)

# ---------------------------------------------------------------------------
# Build the Operator Test Docker image
# ---------------------------------------------------------------------------
build-test-image: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
build-test-image: build-mvn
	docker build -t $(TEST_USER_IMAGE) java/operator-test/target/docker


# ---------------------------------------------------------------------------
# Push the Operator Utils Docker image
# ---------------------------------------------------------------------------
push-test-image: export TEST_USER_IMAGE := $(TEST_USER_IMAGE)
push-test-image:
	@echo "Pushing $(TEST_USER_IMAGE)"
	docker push $(TEST_USER_IMAGE)

# ---------------------------------------------------------------------------
# Build all of the Docker images
# ---------------------------------------------------------------------------
build-all-images: build-operator build-utils-image build-test-image

# ---------------------------------------------------------------------------
# Push all of the Docker images
# ---------------------------------------------------------------------------
push-all-images: push-operator-image push-utils-image push-test-image

# ---------------------------------------------------------------------------
# Build everything
# ---------------------------------------------------------------------------
build-all: build-mvn build-operator


# ---------------------------------------------------------------------------
# Run the Operator in debug mode
# Running this task will start the Operator and pause it until a Delve
# is attached.
#
# To exit out of the local Operator you can use ctrl-c or ctrl-z but
# sometimes this leaves orphaned processes on the local machine so
# ensure these are killed run "make debug-stop"
# ---------------------------------------------------------------------------
run-debug: export TEST_NAMESPACE := $(TEST_NAMESPACE)
run-debug: $(CHART_DIR)/coherence-$(VERSION_FULL).tar.gz reset-namespace create-ssl-secrets uninstall-crds install-crds
	operator-sdk up local --namespace=$(TEST_NAMESPACE) \
	--operator-flags="--watches-file=local-watches.yaml" \
	--enable-delve \
	2>&1 | tee $(TEST_LOGS_DIR)/operator-debug.out

# ---------------------------------------------------------------------------
# Kill any locally running Operator
# ---------------------------------------------------------------------------
debug-stop:
	./hack/kill-local.sh
