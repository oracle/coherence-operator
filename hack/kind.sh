#!/bin/sh
set -o errexit

if [ -h "${0}" ] ; then
    readonly SCRIPT_PATH="$(readlink "${0}")"
else
    readonly SCRIPT_PATH="${0}"
fi

readonly WS_DIR=$(dirname -- "${SCRIPT_PATH}")

# desired cluster name; default is "kind"
KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-operator}"

# create registry container unless it already exists
#reg_name='kind-registry'
#reg_port='5000'
#running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"
#if [ "${running}" != 'true' ]; then
#  docker run \
#    -d --restart=always -p "127.0.0.1:${reg_port}:5000" --name "${reg_name}" \
#    registry:2
#fi

if [ "" = "${KIND_CONFIG}" ]; then
  KIND_CONFIG="${WS_DIR}/kind-config.yaml"
fi

# create a cluster with the local registry enabled in containerd
kind create cluster --name "${KIND_CLUSTER_NAME}" --config "${KIND_CONFIG}" $@
