#!/usr/bin/env bash

kubectl label node operator-worker topology.kubernetes.io/zone=zone-one --overwrite
kubectl label node operator-worker topology.kubernetes.io/region=one --overwrite
kubectl label node operator-worker2 failure-domain.beta.kubernetes.io/zone=zone-two --overwrite
kubectl label node operator-worker2 failure-domain.beta.kubernetes.io/region=two --overwrite
kubectl label node operator-worker3 topology.kubernetes.io/zone=zone-three --overwrite
kubectl label node operator-worker3 topology.kubernetes.io/region=three --overwrite
