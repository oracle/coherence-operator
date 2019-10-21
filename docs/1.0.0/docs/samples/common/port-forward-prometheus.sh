#!/bin/bash

trap "exit" INT

namespace=default
if [ $# -eq 1 ] ; then
   namespace=$1
fi

export PROM_POD=$(kubectl get pod --namespace $namespace | grep prometheus-0 | awk '{print $1}')
echo "Port-forwarding $PROM_POD in $namespace"

while :
do
   kubectl port-forward --namespace $namespace  $PROM_POD 9090:9090
done
