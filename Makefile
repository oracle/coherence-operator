
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

VARS := $(foreach V,$\
		$(sort $(.VARIABLES)),$\
			$(if $(filter-out environment% default automatic, $(origin $V)),$\
				$V=$($V)\
				 	) )

.PHONY: all
all: build

# Ensures that build output directories exist
.PHONY: build-dirs
build-dirs:
	@echo "Creating build directories"
	@mkdir -p build/_output

# Builds the project, helm charts and Docker image
.PHONY: build
build: build-dirs Makefile
	@echo "Building: $(OPERATOR_IMAGE)"
	@./hack/build.sh ${VARS}

# Executes the Go unit tests that do not require a k8s cluster
.PHONY: test
test: build-dirs Makefile
	./hack/test.sh

# This step will run the Operator SDK code generators.
# These commands will generate the CRD files from the API structs and will
# also generate the Go DeepCopy code for the API structs.
# This step would require running if any of the structs in the files under
# the pkg/apis directory have been changed.
.PHONY: generate
generate:
	@echo "Running Operator SDK generate k8s
	@operator-sdk generate k8s
	@echo "Running Operator SDK generate openapi
	@operator-sdk generate openapi
