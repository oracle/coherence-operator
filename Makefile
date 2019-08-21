
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

# Extra arguments to pass to the go test command for the various test steps.
# For example, when running make e2e-test we can run just a single test such
# as the zone test using the go test -run=regex argument like this
#   make e2e-test GO_TEST_FLAGS='-run=^TestZone$$'
GO_TEST_FLAGS ?=

# default as in test/e2e/helper/proj_helpers.go
TEST_NAMESPACE ?= operator-test

IMAGE_PULL_SECRETS ?=

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
RBAC=deploy/service_account.yaml deploy/role.yaml deploy/role_binding.yaml

TEST_MANIFEST_DIR    := $(BUILD_OUTPUT)/manifest
TEST_MANIFEST_FILE   := test-manifest.yaml
TEST_MANIFEST_VALUES ?= deploy/test-values.yaml
TEST_SSL_SECRET      := coherence-ssl-secret

# Do a search and replace of properties in selected files in the Helm charts
# This is done because the Helm charts can be large and processing every file
# makes the build slower
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
build: $(BUILD_OUTPUT)/bin/operator

$(BUILD_OUTPUT)/bin/operator: export CGO_ENABLED = 0
$(BUILD_OUTPUT)/bin/operator: export GOARCH = $(ARCH)
$(BUILD_OUTPUT)/bin/operator: export GOOS = $(OS)
$(BUILD_OUTPUT)/bin/operator: export GO111MODULE = on
$(BUILD_OUTPUT)/bin/operator: $(GOS) $(DEPLOYS) $(CHART_DIR)/coherence-$(VERSION).tar.gz $(CHART_DIR)/coherence-operator-$(VERSION).tar.gz
	@echo "Building: $(OPERATOR_IMAGE)"
	@echo "Running Operator SDK build"
	BUILD_INFO="$(VERSION)|$(GITCOMMIT)|$$(date -u | tr ' ' '.')"; \
	operator-sdk build $(OPERATOR_IMAGE) --verbose --go-build-args "-o $(BUILD_OUTPUT)/bin/operator -ldflags -X=main.BuildInfo=$${BUILD_INFO}"

$(CHART_DIR)/coherence-operator-$(VERSION).tar.gz: $(COP_CHARTS) $(BUILD_PROPS) $(RBAC)
	# Copy the Helm charts from their source location to the distribution folder
	cp -R ./helm-charts/coherence-operator $(CHART_DIR)
	for i in $(RBAC); do \
		f=`basename $${i}`; \
		cp $${i} $(CHART_DIR)/coherence-operator/templates/$${f}; \
	done

	$(call replaceprop,coherence-operator/Chart.yaml coherence-operator/values.yaml coherence-operator/requirements.yaml coherence-operator/templates/deployment.yaml)

	# For each Helm chart folder package the chart into a .tar.gz
	# Package the chart into a .tr.gz - we don't use helm package as the version might not be SEMVER
	echo "Creating Helm chart package $(CHART_DIR)/coherence-operator"
	helm lint $(CHART_DIR)/coherence-operator
	tar -czf $(CHART_DIR)/coherence-operator-$(VERSION).tar.gz $(CHART_DIR)/coherence-operator

$(CHART_DIR)/coherence-$(VERSION).tar.gz: $(COH_CHARTS) $(BUILD_PROPS)
	# Copy the Helm charts from their source location to the distribution folder
	cp -R ./helm-charts/coherence $(CHART_DIR)

	$(call replaceprop,coherence/Chart.yaml coherence/values.yaml)

	# For each Helm chart folder package the chart into a .tar.gz
	# Package the chart into a .tr.gz - we don't use helm package as the version might not be SEMVER
	echo "Creating Helm chart package $(CHART_DIR)/coherence"
	helm lint $(CHART_DIR)/coherence
	tar -czf $(CHART_DIR)/coherence-$(VERSION).tar.gz $(CHART_DIR)/coherence

# Executes the Go unit tests that do not require a k8s cluster
test: export CGO_ENABLED = 0
test: build
	@echo "Running operator tests"
	$(GO_TEST_CMD) test $(GO_TEST_FLAGS) -v ./cmd/... ./pkg/...

# Executes the Go end-to-end tests that require a k8s cluster using
# a local operator instance (i.e. the operator is not deployed to k8s).
# These tests will use whichever k8s cluster the local environment
# is pointing to.
# These tests require the Operator CRDs and will install them before
# tests start and remove them afterwards.
e2e-local-test: export CGO_ENABLED = 0
e2e-local-test: export TEST_LOGS = $(TEST_LOGS_DIR)
e2e-local-test: export TEST_USER_IMAGE = $(RELEASE_IMAGE_PREFIX)oracle/operator-test-image:$(VERSION)
e2e-local-test: export TEST_NAMESPACE := $(TEST_NAMESPACE)
e2e-local-test: build reset-namespace
	@echo "executing end-to-end tests"
	operator-sdk test local ./test/e2e/local --namespace $(TEST_NAMESPACE) --up-local \
		--verbose --debug  --go-test-flags "-timeout=60m $(GO_TEST_FLAGS)" \
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
e2e-test: export TEST_USER_IMAGE = $(RELEASE_IMAGE_PREFIX)oracle/operator-test-image:$(VERSION)
e2e-test: export TEST_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_MANIFEST_FILE)
e2e-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
e2e-test: export TEST_NAMESPACE := $(TEST_NAMESPACE)
e2e-test: build reset-namespace create-ssl-secrets operator-manifest
	@echo "executing end-to-end tests"
	operator-sdk test local ./test/e2e/remote --namespace $(TEST_NAMESPACE) \
		--image $(OPERATOR_IMAGE) --go-test-flags "-timeout=60m $(GO_TEST_FLAGS)" \
		--verbose --debug --namespaced-manifest=$(TEST_MANIFEST) \
		 2>&1 | tee $(TEST_LOGS)/operator-e2e-test.out
	@echo "deleting test namespace"
	kubectl delete namespace $(TEST_NAMESPACE)

# Executes the Go end-to-end Operator Helm chart tests.
# These tests will use whichever k8s cluster the local environment is pointing to.
# These tests require the Operator CRDs and will install them before tests start
# and remove them afterwards.
# Note that the namespace will be created by Helm if it does not exist.
helm-test: export CGO_ENABLED = 0
helm-test: export TEST_LOGS = $(TEST_LOGS_DIR)
helm-test: export TEST_NAMESPACE := $(TEST_NAMESPACE)
helm-test: export TEST_USER_IMAGE = $(RELEASE_IMAGE_PREFIX)oracle/operator-test-image:$(VERSION)
helm-test: export IMAGE_PULL_SECRETS := $(IMAGE_PULL_SECRETS)
helm-test: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
helm-test: build reset-namespace create-ssl-secrets
	@echo "Installing CRDs"
	$(MAKE) install-crds
	@echo "executing Operator Helm Chart end-to-end tests"
	$(GO_TEST_CMD) test $(GO_TEST_FLAGS) -v ./test/e2e/helm/...
	@echo "Removing CRDs"
	$(MAKE) uninstall-crds
	@echo "deleting test namespace"
	kubectl delete namespace $(TEST_NAMESPACE)

# Install CRDs
install-crds: uninstall-crds
	for i in $(CRDS); do \
		kubectl create -f $${i}; \
	done

# Uninstall CRDs
uninstall-crds:
	for i in $(CRDS); do \
		kubectl delete -f $${i} || true; \
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

# Create the k8s yaml manifest that will be used by the Operator SDK to install the Operator when running e2e tests.
operator-manifest: export TEST_NAMESPACE := $(TEST_NAMESPACE)
operator-manifest: export TEST_MANIFEST_DIR := $(TEST_MANIFEST_DIR)
operator-manifest: export TEST_MANIFEST := $(TEST_MANIFEST_DIR)/$(TEST_MANIFEST_FILE)
operator-manifest: export TEST_MANIFEST_VALUES := $(TEST_MANIFEST_VALUES)
operator-manifest: $(CHART_DIR)/coherence-operator-$(VERSION).tar.gz
	@mkdir -p $(TEST_MANIFEST_DIR)
	go run ./cmd/helmutil/

# Generate the keys and certs used in tests.
$(BUILD_OUTPUT)/certs:
	@echo "Generating test keys and certs"
	./hack/keys.sh

# Delete and re-create the test namespace
reset-namespace: export TEST_NAMESPACE := $(TEST_NAMESPACE)
reset-namespace:
	@echo "Deleting test namespace $(TEST_NAMESPACE)"
	kubectl delete namespace $(TEST_NAMESPACE) && echo "deleted namespace" || echo ""
	@echo "Creating test namespace $(TEST_NAMESPACE)"
	kubectl create namespace $(TEST_NAMESPACE)


# Create the k8s secret to use in SSL/TLS testing.
create-ssl-secrets: export TEST_NAMESPACE := $(TEST_NAMESPACE)
create-ssl-secrets: export TEST_SSL_SECRET := $(TEST_SSL_SECRET)
create-ssl-secrets: $(BUILD_OUTPUT)/certs
	@echo "Deleting SSL secret $(TEST_SSL_SECRET)"
	kubectl --namespace $(TEST_NAMESPACE) delete secret $(TEST_SSL_SECRET) && echo "secret deleted" || echo ""
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

