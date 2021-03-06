#!/bin/bash

trap "exit" INT

if [ $# -ne 1 ] ; then
   echo "Usage: $0 [namespace]"
   exit 1
fi

namespace=$1

export PROM_POD=$(kubectl get pod --namespace $namespace | grep prometheus-0 | awk '{print $1}')
echo "Port-forwarding $PROM_POD in $namespace"

echo "Open the following URL: http://127.0.0.1:9090/targets"

while :
do
   kubectl port-forward --namespace $namespace  $PROM_POD 9090:9090
done
