#!/usr/bin/env bash

kubectl label node k3d-operator-agent-0 topology.kubernetes.io/zone=zone-one --overwrite
kubectl label node k3d-operator-agent-0 topology.kubernetes.io/region=one --overwrite
kubectl label node k3d-operator-agent-0 oci.oraclecloud.com/fault-domain=fd-one --overwrite
kubectl label node k3d-operator-agent-0 coherence.oracle.com/test=test-one --overwrite

kubectl label node k3d-operator-agent-1 topology.kubernetes.io/zone=zone-one --overwrite
kubectl label node k3d-operator-agent-1 topology.kubernetes.io/region=one --overwrite
kubectl label node k3d-operator-agent-1 oci.oraclecloud.com/fault-domain=fd-one --overwrite
kubectl label node k3d-operator-agent-1 coherence.oracle.com/test=test-one --overwrite

kubectl label node k3d-operator-agent-2 topology.kubernetes.io/zone=zone-two --overwrite || true
kubectl label node k3d-operator-agent-2 topology.kubernetes.io/region=two --overwrite || true
kubectl label node k3d-operator-agent-2 oci.oraclecloud.com/fault-domain=fd-two --overwrite || true
kubectl label node k3d-operator-agent-2 coherence.oracle.com/test=test-two --overwrite || true

kubectl label node k3d-operator-agent-3 topology.kubernetes.io/zone=zone-two --overwrite || true
kubectl label node k3d-operator-agent-3 topology.kubernetes.io/region=two --overwrite || true
kubectl label node k3d-operator-agent-3 oci.oraclecloud.com/fault-domain=fd-two --overwrite || true
kubectl label node k3d-operator-agent-3 coherence.oracle.com/test=test-two --overwrite || true

kubectl label node k3d-operator-agent-4 topology.kubernetes.io/zone=zone-three --overwrite || true
kubectl label node k3d-operator-agent-4 topology.kubernetes.io/region=three --overwrite || true
kubectl label node k3d-operator-agent-4 oci.oraclecloud.com/fault-domain=fd-three --overwrite || true
kubectl label node k3d-operator-agent-4 coherence.oracle.com/test=test-three --overwrite || true

kubectl label node k3d-operator-server-0 topology.kubernetes.io/zone=zone-three --overwrite || true
kubectl label node k3d-operator-server-0 topology.kubernetes.io/region=three --overwrite || true
kubectl label node k3d-operator-server-0 oci.oraclecloud.com/fault-domain=fd-three --overwrite || true
kubectl label node k3d-operator-server-0 coherence.oracle.com/test=test-three --overwrite || true

