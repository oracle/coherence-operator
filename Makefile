
VERSION         ?= 2.0.0-SNAPSHOT
GITCOMMIT       ?= $(shell git rev-list -1 HEAD)

ARCH            ?= amd64
OS              ?= linux
UNAME_S         := $(shell uname -s)

REGISTRY        ?= iad.ocir.io
TENANT          ?= odx-stateservice/test/oracle
OPERATOR_IMAGE  := ${REGISTRY}/${TENANT}/coherence-operator:${VERSION}

HELM_COHERENCE_IMAGE ?= container-registry.oracle.com/middleware/coherence:12.2.1.3.2
HELM_UTILS_IMAGE     ?= iad.ocir.io/odx-stateservice/test/oracle/coherence-operator:1.1.0-SNAPSHOT-utils

PROMETHEUS_HELMCHART_VERSION ?= 5.7.0

override BUILD_OUTPUT := ./build/_output
override BUILD_PROPS  := $(BUILD_OUTPUT)/build.properties
override CHART_DIR    := $(BUILD_OUTPUT)/helm-charts

all: build

# Ensures that build output directories exist
build-dirs:
	@echo "Creating build directories"
	@mkdir -p $(BUILD_OUTPUT)

# Builds the project, helm charts and Docker image
build: export CGO_ENABLED = 0
build: export GOARCH = $(ARCH)
build: export GOOS = $(OS)
build: export GO111MODULE = on
build: build-dirs
	@echo "Building: $(OPERATOR_IMAGE)"
	# create build.properties
	rm -f $(BUILD_PROPS)
	printf "HELM_COHERENCE_IMAGE=$(HELM_COHERENCE_IMAGE)\n\
	HELM_UTILS_IMAGE=$(HELM_UTILS_IMAGE)\n\
	OPERATOR_IMAGE=$(OPERATOR_IMAGE)\n\
	PROMETHEUS_HELMCHART_VERSION=$(PROMETHEUS_HELMCHART_VERSION)\n\
	VERSION=$(VERSION)\n" > $(BUILD_PROPS)

	# create Helm charts
	@echo "Creating Helm chart distributions"
	# Copy the Helm charts from their source location to the distribution folder
	rm -rf $(CHART_DIR); mkdir $(CHART_DIR); cp -R ./helm-charts/* $(CHART_DIR)
	cp ./deploy/role.yaml            $(CHART_DIR)/coherence-operator/templates/role.yaml
	cp ./deploy/role_binding.yaml    $(CHART_DIR)/coherence-operator/templates/role_binding.yaml
	cp ./deploy/service_account.yaml $(CHART_DIR)/coherence-operator/templates/service_account.yaml

	# Do a search and replace of properties in selected files in the Helm charts
	# This is done because the Helm charts can be large and processing every file
	# makes the build slower
	for i in coherence/Chart.yaml coherence/values.yaml coherence-operator/Chart.yaml coherence-operator/values.yaml \
			coherence-operator/requirements.yaml coherence-operator/templates/deployment.yaml; do \
		filename="$(CHART_DIR)/$${i}"; \
		echo "Replacing properties in file $${filename}"; \
		if [[ -f $${filename} ]]; then \
			temp_file=$(BUILD_OUTPUT)/temp.out; \
			awk -F'=' 'NR==FNR {a[$$1]=$$2;next} {for (i in a) {x = sprintf("\\$${%s}", i); gsub(x, a[i])}}1' ${BUILD_PROPS} $${filename} > $${temp_file}; \
			mv $${temp_file} $${filename}; \
		fi; \
	done

	# For each Helm chart folder package the chart into a .tar.gz
	# Package the chart into a .tr.gz - we don't use helm package as the version might not be SEMVER
	for chart in `ls -1ad $(BUILD_OUTPUT)/helm-charts/*`; do \
		chartname=$$(basename $${chart}); \
		echo "Creating Helm chart package $${chart}"; \
		helm lint $${chart}; \
		tar -czf $(CHART_DIR)/$${chartname}-$(VERSION).tar.gz $${chart}; \
	done

	@echo "Creating CRD distribution"
	rm -rf $(BUILD_OUTPUT)/crds/; mkdir $(BUILD_OUTPUT)/crds/; cp -R ./deploy/crds/*_crd.yaml $(BUILD_OUTPUT)/crds/
	@echo "Creating deployment yaml files"
	rm -rf $(BUILD_OUTPUT)/yaml/; mkdir $(BUILD_OUTPUT)/yaml/; cp -R ./deploy/*.yaml $(BUILD_OUTPUT)/yaml/

	@echo "Running Operator SDK build"
	BUILD_INFO="$(VERSION)|$(GITCOMMIT)|$$(date -u | tr ' ' '.')"; \
	operator-sdk build $(OPERATOR_IMAGE) --verbose --go-build-args "-o $(BUILD_OUTPUT)/bin/operator -ldflags -X=main.BuildInfo=$${BUILD_INFO}"

# Executes the Go unit tests that do not require a k8s cluster
test: export CGO_ENABLED = 0
test: build-dirs
	@echo "Running operator tests"
	if [ -z `which ginkgo` ]; then \
		CMD=go; \
	else \
		CMD=ginkgo; \
	fi; \
	$$CMD test -v ./cmd/... ./pkg/...

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

.PHONY: all build-dirs build test generate
