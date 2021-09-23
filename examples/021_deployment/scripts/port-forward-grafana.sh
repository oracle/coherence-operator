#!/bin/bash

trap "exit" INT

echo "Open the following URL: http://localhost:3000/d/coh-main/coherence-dashboard-main"

while :
do
   kubectl --namespace monitoring port-forward svc/grafana 3000
done
