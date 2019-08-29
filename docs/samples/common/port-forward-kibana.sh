#!/bin/bash

trap "exit" INT

namespace=default
if [ $# -eq 1 ] ; then
   namespace=$1
fi
   
export KIBANA_POD=$(kubectl get pods --namespace $namespace | grep kibana | awk '{print $1}')
echo "Port-forwarding $KIBANA_POD in $namespace"

while :
do
   kubectl port-forward --namespace $namespace $KIBANA_POD 5601:5601
done
