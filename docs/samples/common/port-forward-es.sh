#!/bin/bash

trap "exit" INT

namespace=default
if [ $# -eq 1 ] ; then
   namespace=$1
fi
   
export ELASTICSEARCH_POD=$(kubectl get pods --namespace $namespace | grep elasticsearch | awk '{print $1}')
echo "Port-forwarding $ELASTICSEARCH_POD in $namespace"

while :
do
   kubectl port-forward --namespace $namespace $ELASTICSEARCH_POD 9200:9200
done
