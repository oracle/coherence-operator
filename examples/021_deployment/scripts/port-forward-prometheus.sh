#!/bin/bash

trap "exit" INT

echo "Open the following URL: http://127.0.0.1:9090/targets"

while :
do
   kubectl --namespace monitoring port-forward svc/prometheus-k8s 9090
done
