#!/bin/sh
set -o errexit

# desired cluster name; default is "kind"
KIND_CLUSTER_NAME="${KIND_CLUSTER_NAME:-operator}"

kind create cluster --name "${KIND_CLUSTER_NAME}" $@
