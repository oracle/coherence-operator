#!/bin/bash

trap "exit" INT

if [ $# -ne 1 ] ; then
   echo "Usage: $0 [namespace]"
   exit 1
fi

namespace=$1
   
export KIBANA_POD=$(kubectl get pods --namespace $namespace | grep kibana | awk '{print $1}')
echo "Port-forwarding $KIBANA_POD in $namespace"
echo "Open the following URL: http://127.0.0.1:5601/"

while :
do
   kubectl port-forward --namespace $namespace $KIBANA_POD 5601:5601
done
