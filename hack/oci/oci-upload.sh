#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

#!/usr/bin/env sh

UPLOAD_TIMESTAMP=$(date +"%Y%m%d%H%M")
FILE_NAME=pipelinerun-${UPLOAD_TIMESTAMP}.tgz
FULL_PATH_NAME=${HOME}/${FILE_NAME}
PA_JSON=${HOME}/pa.json
EXPIRY_DATE=$(date -d '+1 days' --iso-8601=seconds)

tar -C build/_output -czf ${FULL_PATH_NAME} .
oci os object put --bucket-name coherence-cert-tests --file ${FULL_PATH_NAME}
oci os preauth-request create -bn coherence-cert-tests \
    --time-expires=${EXPIRY_DATE} \
    --access-type ObjectRead \
    --name pipelinerun-${UPLOAD_TIMESTAMP} \
    -on ${FILE_NAME} > ${PA_JSON}

cat ${PA_JSON}
PA_URI=$(cat ${PA_JSON} | jq -r '.data."access-uri"')
echo -n "https://objectstorage.${REGION}.oraclecloud.com${PA_URI}" | tee ${TASK_RESULT_PATH}
