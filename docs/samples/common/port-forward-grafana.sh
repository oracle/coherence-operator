#!/bin/bash

trap "exit" INT

namespace=default
if [ $# -eq 1 ] ; then
   namespace=$1
fi

export GRAFANA_POD=$(kubectl get pods --namespace $namespace | grep grafana | awk '{print $1}')
echo "Port-forwarding $GRAFANA_POD in $namespace"

while :
do
   kubectl port-forward --namespace $namespace $GRAFANA_POD 3000:3000
done
