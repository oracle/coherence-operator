#!/usr/bin/env bash

kubectl label node operator-worker topology.kubernetes.io/zone=zone-one --overwrite
kubectl label node operator-worker topology.kubernetes.io/region=one --overwrite
kubectl label node operator-worker oci.oraclecloud.com/fault-domain=fd-one --overwrite
kubectl label node operator-worker coherence.oracle.com/test=test-one --overwrite
kubectl label node operator-worker2 topology.kubernetes.io/zone=zone-two --overwrite || true
kubectl label node operator-worker2 topology.kubernetes.io/region=two --overwrite || true
kubectl label node operator-worker2 oci.oraclecloud.com/fault-domain=fd-two --overwrite || true
kubectl label node operator-worker2 coherence.oracle.com/test=test-two --overwrite || true
kubectl label node operator-worker3 topology.kubernetes.io/zone=zone-three --overwrite || true
kubectl label node operator-worker3 topology.kubernetes.io/region=three --overwrite || true
kubectl label node operator-worker3 oci.oraclecloud.com/fault-domain=fd-three --overwrite || true
kubectl label node operator-worker3 coherence.oracle.com/test=test-three --overwrite || true

