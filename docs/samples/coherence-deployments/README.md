# Coherence Deployments

# Table of Contents     

1. [Add application jars/Config to a Coherence deployment](sidecar)
1. [Accessing Coherence via Coherence*Extend](extend)
   1. [Access Coherence via default proxy port](extend/default)
   1. [Access Coherence via separate proxy tier](extend/proxy-tier)
   1. [Enabling SSL for Proxy Servers](extend/ssl)
   1. [Using multiple Coherence*Extend proxies](extend/multiple)
1. [Accessing Coherence via storage-disabled clients](storage-disabled)
   1. [Storage-disabled client in cluster via interceptor](storage-disabled/interceptor)
   1. [Storage-disabled client in cluster as separate user image](storage-disabled/other)
1. [Federation](federation)   
   1. [Within a single Kubernetes cluster](federation/within-cluster)
   1. [Across across separate Kubernets clusters](federation/across-clusters)
1. [Persistence](persistence)
   1. [Use default persistent volume claim](persistence/default)
   1. [Use a specific persistent volume](persistence/pvc)
   1. [Specify a separate snapshot location for active persistence](persistence/snapshot)
   1. [Specifying an archiver](persistence/archiver)
1. [Elastic Data](elastic-data)
   1. [Deploy using default FlashJournal locations](elastic-data/default)
   1. [Deploy using external volume mapped to the host](elastic-data/pvc)
1. [Installing Multiple Coherence clusters with one Operator](multiple-clusters)   
   
[Return to samples](../README.md#list-of-samples)   