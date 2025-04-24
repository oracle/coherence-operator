#
# Copyright (c) 2020, 2025, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
set -o errexit


ROOT_DIR=$(pwd)
TOOLS_BIN=${ROOT_DIR}/build/tools/bin

UNAME_S=$(uname -s)
UNAME_M=$(uname -m)

if [ "Darwin" = "${UNAME_S}" ]; then
  if [ "x86_64" = "${UNAME_M}" ]; then
    echo "Downloading OpenShift OC CLI ${UNAME_S} ${UNAME_M}"
  	curl -Ls https://mirror.openshift.com/pub/openshift-v4/x86_64/clients/ocp/stable/openshift-client-mac.tar.gz -o openshift-client.tar.gz
  else
    echo "Downloading OpenShift OC CLI ${UNAME_S} ${UNAME_M}"
    curl -Ls https://mirror.openshift.com/pub/openshift-v4/aarch64/clients/ocp/stable/openshift-client-mac-arm64.tar.gz -o openshift-client.tar.gz
  fi
else
  if [ "x86_64" = "${UNAME_M}" ]; then
    echo "Downloading OpenShift OC CLI ${UNAME_S} ${UNAME_M}"
    curl -Ls https://mirror.openshift.com/pub/openshift-v4/x86_64/clients/ocp/stable/openshift-client-linux.tar.gz -o openshift-client.tar.gz
  else
    echo "Downloading OpenShift OC CLI ${UNAME_S} ${UNAME_M}"
    curl -Ls https://mirror.openshift.com/pub/openshift-v4/aarch64/clients/ocp/stable/openshift-client-linux.tar.gz -o openshift-client.tar.gz
  fi
fi

tar -zxvf openshift-client.tar.gz oc
mv oc ${TOOLS_BIN}/oc
chmod +x ${TOOLS_BIN}/oc
rm openshift-client.tar.gz
