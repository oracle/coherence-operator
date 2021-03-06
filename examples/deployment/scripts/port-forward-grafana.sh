#!/bin/bash

trap "exit" INT

if [ $# -ne 1 ] ; then
   echo "Usage: $0 [namespace]"
   exit 1
fi

namespace=$1

export GRAFANA_POD=$(kubectl get pods --namespace $namespace | grep prometheus-grafana | grep -v test | awk '{print $1}')
echo "Port-forwarding $GRAFANA_POD in $namespace"

echo "Open the following URL: http://127.0.0.1:3000/d/coh-main/coherence-dashboard-main"

while :
do
   kubectl port-forward --namespace $namespace $GRAFANA_POD 3000:3000
done
