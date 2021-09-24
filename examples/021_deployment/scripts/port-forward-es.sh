#!/bin/bash

trap "exit" INT

if [ $# -ne 1 ] ; then
   echo "Usage: $0 [namespace]"
   exit 1
fi

namespace=$1

export ELASTICSEARCH_POD=$(kubectl get pods --namespace $namespace | grep elasticsearch | awk '{print $1}')
echo "Port-forwarding $ELASTICSEARCH_POD in $namespace"

while :
do
   kubectl port-forward --namespace $namespace $ELASTICSEARCH_POD 9200:9200
done
