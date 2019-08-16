
VERSION         ?= 2.0.0-SNAPSHOT
GITCOMMIT       ?= $(shell git rev-list -1 HEAD)

ARCH            ?= amd64
OS              ?= linux
UNAME_S         := $(shell uname -s)

COHERENCE_IMAGE_PREFIX ?= container-registry.oracle.com/middleware/
HELM_COHERENCE_IMAGE   ?= $(COHERENCE_IMAGE_PREFIX)coherence:12.2.1.4.0-b74630

# One may need to define RELEASE_IMAGE_PREFIX in the environment.
OPERATOR_IMAGE   := $(RELEASE_IMAGE_PREFIX)oracle/coherence-operator:$(VERSION)
HELM_UTILS_IMAGE ?= $(RELEASE_IMAGE_PREFIX)oracle/coherence-operator:$(VERSION)-utils

PROMETHEUS_HELMCHART_VERSION ?= 5.7.0

# default as in test/e2e/helper/proj_helpers.go
TEST_NAMESPACE ?= operator-test

override BUILD_OUTPUT  := ./build/_output
override BUILD_PROPS   := $(BUILD_OUTPUT)/build.properties
override CHART_DIR     := $(BUILD_OUTPUT)/helm-charts
override TEST_LOGS_DIR := $(BUILD_OUTPUT)/test-logs

ifeq (, $(shell which ginkgo))
GO_TEST_CMD = go
else
GO_TEST_CMD = ginkgo
endif

# Do a search and replace of properties in selected files in the Helm charts
# This is done because the Helm charts can be large and processing every file
# makes the build slower
define replaceprop
	filename="$(CHART_DIR)/$(1)"; \
	echo "Replacing properties in file $${filename}"; \
	if [[ -f $${filename} ]]; then \
		temp_file=$(BUILD_OUTPUT)/temp.out; \
		awk -F'=' 'NR==FNR {a[$$1]=$$2;next} {for (i in a) {x = sprintf("\\$${%s}", i); gsub(x, a[i])}}1' $(BUILD_PROPS) $${filename} > $${temp_file}; \
		mv $${temp_file} $${filename}; \
	fi
endef

.PHONY: all build test e2e-local-test e2e-test install-crds uninstall-crds generate push clean

all: build

$(BUILD_PROPS):
	# Ensures that build output directories exist
	@echo "Creating build directories"
	@mkdir -p $(BUILD_OUTPUT)
	@mkdir -p $(TEST_LOGS_DIR)
	@mkdir -p $(CHART_DIR)
	# create build.properties
	rm -f $(BUILD_PROPS)
	printf "HELM_COHERENCE_IMAGE=$(HELM_COHERENCE_IMAGE)\n\
	HELM_UTILS_IMAGE=$(HELM_UTILS_IMAGE)\n\
	OPERATOR_IMAGE=$(OPERATOR_IMAGE)\n\
	PROMETHEUS_HELMCHART_VERSION=$(PROMETHEUS_HELMCHART_VERSION)\n\
	VERSION=$(VERSION)\n" > $(BUILD_PROPS)

# Builds the project, helm charts and Docker image
build: export CGO_ENABLED = 0
build: export GOARCH = $(ARCH)
build: export GOOS = $(OS)
build: export GO111MODULE = on
build: $(BUILD_OUTPUT)/bin/operator

GOS=$(shell find pkg -type f -name "*.go" ! -name "*_test.go")

$(BUILD_OUTPUT)/bin/operator: $(GOS) $(CHART_DIR)/coherence-$(VERSION).tar.gz $(CHART_DIR)/coherence-operator-$(VERSION).tar.gz
	@echo "Building: $(OPERATOR_IMAGE)"
	@echo "Running Operator SDK build"
	BUILD_INFO="$(VERSION)|$(GITCOMMIT)|$$(date -u | tr ' ' '.')"; \
	operator-sdk build $(OPERATOR_IMAGE) --verbose --go-build-args "-o $(BUILD_OUTPUT)/bin/operator -ldflags -X=main.BuildInfo=$${BUILD_INFO}"

COH_CHARTS=$(shell find helm-charts/coherence -type f)

$(CHART_DIR)/coherence-$(VERSION).tar.gz: $(COH_CHARTS) $(BUILD_PROPS)
	# Copy the Helm charts from their source location to the distribution folder
	cp -R ./helm-charts/coherence $(CHART_DIR)

	$(call replaceprop,coherence/Chart.yaml)
	$(call replaceprop,coherence/values.yaml)

	# For each Helm chart folder package the chart into a .tar.gz
	# Package the chart into a .tr.gz - we don't use helm package as the version might not be SEMVER
	echo "Creating Helm chart package $(CHART_DIR)/coherence"
	helm lint $(CHART_DIR)/coherence
	tar -czf $(CHART_DIR)/coherence-$(VERSION).tar.gz $(CHART_DIR)/coherence

COP_CHARTS=$(shell find helm-charts/coherence-operator -type f)

$(CHART_DIR)/coherence-operator-$(VERSION).tar.gz: $(COP_CHARTS) $(BUILD_PROPS)
	# Copy the Helm charts from their source location to the distribution folder
	cp -R ./helm-charts/coherence-operator $(CHART_DIR)
	for i in role.yaml role_binding.yaml service_account.yaml; do \
		cp ./deploy/$${i} $(CHART_DIR)/coherence-operator/templates/$${i}; \
	done

	$(call replaceprop,coherence-operator/Chart.yaml)
	$(call replaceprop,coherence-operator/values.yaml)
	$(call replaceprop,coherence-operator/requirements.yaml)
	$(call replaceprop,coherence-operator/templates/deployment.yaml)

	# For each Helm chart folder package the chart into a .tar.gz
	# Package the chart into a .tr.gz - we don't use helm package as the version might not be SEMVER
	echo "Creating Helm chart package $(CHART_DIR)/coherence-operator"
	helm lint $(CHART_DIR)/coherence-operator
	tar -czf $(CHART_DIR)/coherence-operator-$(VERSION).tar.gz $(CHART_DIR)/coherence-operator

# Executes the Go unit tests that do not require a k8s cluster
test: export CGO_ENABLED = 0
test: build
	@echo "Running operator tests"
	$(GO_TEST_CMD) test -v ./cmd/... ./pkg/...

# Executes the Go end-to-end tests that require a k8s cluster using
# a local operator instance (i.e. the operator is not deployed to k8s).
# These tests will use whichever k8s cluster the local environment
# is pointing to.
# These tests require the Operator CRDs and will install them before
# tests start and remove them afterwards.
e2e-local-test: export CGO_ENABLED = 0
e2e-local-test: export TEST_LOGS = $(TEST_LOGS_DIR)
e2e-local-test: build
	@echo "creating test namespace"
	kubectl create namespace $(TEST_NAMESPACE)
	@echo "executing end-to-end tests"
	operator-sdk test local ./test/e2e/local --namespace $(TEST_NAMESPACE) --up-local \
		--verbose --debug \
		--local-operator-flags "--watches-file=local-watches.yaml" \
		 2>&1 | tee $(TEST_LOGS)/operator-e2e-local-test.out
	@echo "deleting test namespace"
	kubectl delete namespace $(TEST_NAMESPACE)

# Executes the Go end-to-end tests that require a k8s cluster using
# a deployed operator instance (i.e. the operator Docker image is
# deployed to k8s). These tests will use whichever k8s cluster the
# local environment is pointing to.
# These tests require the Operator CRDs and will install them before
# tests start and remove them afterwards.
e2e-test: export CGO_ENABLED = 0
e2e-test: export TEST_LOGS = $(TEST_LOGS_DIR)
e2e-test: build
	@echo "creating test namespace"
	kubectl create namespace $(TEST_NAMESPACE)
	@echo "executing end-to-end tests"
	operator-sdk test local ./test/e2e/remote --namespace $(TEST_NAMESPACE) \
		--image iad.ocir.io/odx-stateservice/test/oracle/coherence-operator:2.0.0-SNAPSHOT \
		--verbose --debug \
		 2>&1 | tee $(TEST_LOGS)/operator-e2e-test.out
	@echo "deleting test namespace"
	kubectl delete namespace $(TEST_NAMESPACE)

# Executes the Go end-to-end Operator Helm chart tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# These tests require the Operator CRDs and will install them before tests start
# and remove them afterwards.
helm-test: export CGO_ENABLED = 0
helm-test: export TEST_LOGS = $(TEST_LOGS_DIR)
helm-test: export TEST_NAMESPACE := $(TEST_NAMESPACE)
helm-test: build
	@echo "creating test namespace"
	kubectl create namespace $(TEST_NAMESPACE)
	@echo "Installing CRDs"
	$(MAKE) install-crds
	@echo "executing Operator Helm Chart end-to-end tests"
	$(GO_TEST_CMD) test -v ./test/e2e/helm/...
	@echo "Removing CRDs"
	$(MAKE) uninstall-crds
	@echo "deleting test namespace"
	kubectl delete namespace $(TEST_NAMESPACE)

# Install CRDs
install-crds: uninstall-crds
	for i in coherence_v1_coherencerole_crd.yaml coherence_v1_coherencecluster_crd.yaml coherence_v1_coherenceinternal_crd.yaml; do \
		kubectl create -f deploy/crds/$${i}; \
	done

# Uninstall CRDs
uninstall-crds:
	for i in coherence_v1_coherenceinternal_crd.yaml coherence_v1_coherencerole_crd.yaml coherence_v1_coherencecluster_crd.yaml; do \
		kubectl delete -f deploy/crds/$${i} || true; \
	done

# This step will run the Operator SDK code generators.
# These commands will generate the CRD files from the API structs and will
# also generate the Go DeepCopy code for the API structs.
# This step would require running if any of the structs in the files under
# the pkg/apis directory have been changed.
generate:
	@echo "Generating deep copy code"
	operator-sdk generate k8s
	@echo "Generating Open API code and CRDs"
	operator-sdk generate openapi

# This step push the operator image to registry.
push:
	@echo "Pushing $(OPERATOR_IMAGE)"
	docker push $(OPERATOR_IMAGE)

clean:
	rm -rf build/_output
