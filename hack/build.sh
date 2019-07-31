#!/usr/bin/env bash

set -o errexit
set -o nounset

# Convert all the script args to variables
# The arguments passed in from the Makefile are a list of key=value pairs
for v in ${@}; do
    IFS='=' read -ra ADDR <<< "${v}"
    ARRLEN=${#ADDR[@]}
    if [[ ${ARRLEN} -gt 1 ]]
    then
        declare "${ADDR[0]}"="${ADDR[1]}"
    fi
done

export CGO_ENABLED=0
export GOARCH="${ARCH}"
export GOOS="${OS}"
export GO111MODULE=on

# Main build steps
main()
    {
    createBuildProperties

    createHelmCharts

    echo "Creating CRD distribution"
    copy "./deploy/crds/*_crd.yaml" "./build/_output/crds/"

    echo "Creating deployment yaml files"
    copy "./deploy/*.yaml" "./build/_output/yaml/"

    goBuild
    }

# Build the Go components
goBuild()
    {
    echo "Running Operator SDK build"
    DATE=$(date -u)
    BUILD_INFO="${VERSION}|${GITCOMMIT}|${DATE// /.}"

    operator-sdk build ${OPERATOR_IMAGE} --verbose --go-build-args "-o ./build/_output/bin/operator -ldflags -X=main.BuildInfo=${BUILD_INFO}"
    }

# Create the Helm chart distributions
createHelmCharts()
    {
    echo "Creating Helm chart distributions"
    local CHART_DIR=./build/_output/helm-charts

    # Copy the Helm charts from their source location to the distribution folder
    copy "./helm-charts/"                "build/_output/helm-charts"
    cp   "./deploy/role.yaml"            "build/_output/helm-charts/coherence-operator/templates/role.yaml"
    cp   "./deploy/role_binding.yaml"    "build/_output/helm-charts/coherence-operator/templates/role_binding.yaml"
    cp   "./deploy/service_account.yaml" "build/_output/helm-charts/coherence-operator/templates/service_account.yaml"

    # Do a search and replace of properties in selected files in the Helm charts
    # This is done because the Helm charts can be large and processing every file
    # makes the build slower
    replaceProperties build/_output/helm-charts/coherence/Chart.yaml
    replaceProperties build/_output/helm-charts/coherence/values.yaml
    replaceProperties build/_output/helm-charts/coherence-operator/Chart.yaml
    replaceProperties build/_output/helm-charts/coherence-operator/values.yaml
    replaceProperties build/_output/helm-charts/coherence-operator/requirements.yaml
    replaceProperties build/_output/helm-charts/coherence-operator/templates/deployment.yaml

    # For each Helm chart folder package the chart into a .tar.gz
    local chart_list=$(ls -1ad ./build/_output/helm-charts/*)
    for chart in ${chart_list}; do
        local chartname=$(basename ${chart})
        echo "Creating Helm chart package ${chart}"

        # Run helm lint to verify the chart is valid
        helm lint ${chart}
        # Package the chart into a .tr.gz - we don;t use helm package as the version might not be SEMVER
        tar -czf ${CHART_DIR}/${chartname}-${VERSION}.tar.gz ${chart}
    done
    }

# Replace values in all files in a folder
# Values are replaced from the build.properties file
copyAndReplaceFiles()
    {
    local SRC=${1}
    local DEST=${2}

    copy "${SRC}" "${DEST}"

    find ${DEST} -type f | while read filename
    do
        replaceProperties ${filename}
    done
    }

# Copy a folder and its contents overwriting the destination
copy()
    {
    local SRC=${1}
    local DEST=${2}

    echo "Copying files from '${SRC}' to '${DEST}'"

    rm -rf ${DEST}
    mkdir ${DEST}
    cp -R ${SRC} ${DEST}
    }

# Replace any properties of the format ${name} with the
# corresponding value from the build properties file,
# which includes environment variables and variables
# passed in from the Makefile
replaceProperties()
    {
    echo "Replacing properties in file ${1}"
    if [[ -f ${1} ]]
    then
        local filename=${1}
        local temp_file=./build/_output/temp.out
        awk -F'=' 'NR==FNR {a[$1]=$2;next} {for (i in a) {x = sprintf("\\${%s}",  i);  gsub(x, a[i])}}1' ${BUILD_PROPS} ${filename} > ${temp_file}
        mv ${temp_file} ${filename}
    fi
    }

createBuildProperties()
    {
    # Combine ALL script and environment variables into a single list.
    # This is probably too much as it includes a lot of variables, we probably
    # want to do something to build a more selective list of build properties.
    PROPS=$(set)
    BUILD_PROPS=./build/_output/build.properties

    if [[ -f ${BUILD_PROPS} ]]
    then
        rm ${BUILD_PROPS}
    fi

    if [[ -f "build.properties" ]]
    then
        cp build.properties ${BUILD_PROPS}
    fi
    printf "\n${PROPS}" >> ${BUILD_PROPS}
    }

main