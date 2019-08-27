#!/usr/bin/env bash

export OPERATOR_NAME=coherence-operator

script_name=$0
script_full_path=$(dirname "$0")

./${script_full_path}/kill-local.sh

operator-sdk up local --namespace=test-jk --operator-flags="--watches-file=local-watches.yaml" --enable-delve  2>&1 | tee operator.out
