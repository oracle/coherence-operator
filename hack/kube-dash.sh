#!/bin/sh

# Display the admin tokens used by k8s - this makes it easy to then copy and paste one of the tokens into the Dashboard login page
kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep kube-proxy-token | awk '{print $1}')

echo ""
echo "Access the K8s Dashboard at http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/"
echo ""


kubectl proxy