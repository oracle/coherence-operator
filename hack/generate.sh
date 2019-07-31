#!/usr/bin/env bash

echo "Generating deep copy code"
operator-sdk generate k8s

echo "Generating Open API code and CRDs"
operator-sdk generate openapi
