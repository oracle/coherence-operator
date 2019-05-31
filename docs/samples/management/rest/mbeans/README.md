# Modify Writable MBeans

Management over REST provides the ability to modify Coherence MBeans that are writable.

This sample will show how to do this using curl to modify the log level of a member as well
as the expiryDelay of a cache.

> **Note**: use of Management over REST is only available when using the
> operator with Coherence 12.2.1.4.

[Return to Management over REST samples](../)  [Return to Management samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Prerequisites

1. Install Coherence Cluster

   Follow the instructions [here](../standard/README.md#installation-steps) to install a Coherence cluster and port-forward the Management
   over REST port.

2. Install the Coherence JVisualVM plugin

   Follow the instructions [here](https://docs.oracle.com/middleware/12213/coherence/manage/using-jmx-manage-oracle-coherence.htm)
   to install the JVisualVM plugin.

## Installation Steps
   
1. Retrieve the current loggingLevel
   
   ```bash
   $ curl http://127.0.0.1:30000/management/coherence/cluster/members?fields=loggingLevel 2> /dev/null | json_pp | grep "loggingLevel"

   "loggingLevel" : 5,
   "loggingLevel" : 5,
   "loggingLevel" : 5,
   "loggingLevel" : 5,
   ``` 
   
   This shows all members are running at logging level = 5 (the default).  
   
   Ensure there are no D6 messages:
   
   ```bash
   $ kubectl logs storage-coherence-0 -n sample-coherence-ns | grep D6
   ```
   
   The above should not return anything.
   
1. Set the loggingLevel for each member to 9

   This must be done against each member:
   
   ```bash
   $ for i in 1 2 3 4; do 
      curl -X POST -H 'Content-type: application/json' http://127.0.0.1:30000/management/coherence/cluster/members/$i -d '{"loggingLevel": 9}'
   done
   {}{}{}{}
   ```
   
1. Re-Query the logging level

   ```bash
   $ curl http://127.0.0.1:30000/management/coherence/cluster/members?fields=loggingLevel 2> /dev/null | json_pp | grep "loggingLevel"

  "loggingLevel" : 9
  "loggingLevel" : 9
  "loggingLevel" : 9,
  "loggingLevel" : 9

1. Add data to the cluster via the Coherence Console

   Connect to the Coherence `console` using the following to create a cache.  This will also trigger log
   messages for the joining member.

   ```bash
   $ kubectl exec -it --namespace sample-coherence-ns storage-coherence-0 bash /scripts/startCoherence.sh console
   ```   
   
   At the `Map (?):` prompt, type `cache test`.  This will create a cache in the service `PartitionedCache`.
   
   Add and entry to the cache using
   
   ```bash
   put 1 one

   null
   ```
   
   Then type `bye` to exit the `console`.
   
1. Inspect the log files for D6 messages

   ```bash
   $ kubectl logs storage-coherence-0 -n sample-coherence-ns | grep D6
   2019-04-24 04:58:56.203/3687.142 Oracle Coherence GE 12.2.1.4.0 <D6> (thread=Cluster, member=1): TcpRing connected to Member(Id=5, Timestamp=2019-04-24 04:58:55.99, Address=10.1.4.147:32923, MachineId=30443, Location=site:coherence.sample-coherence-ns.svc.cluster.local,machine:docker-for-desktop,process:6020,member:storage-coherence-0, Role=CoherenceConsole)
   2019-04-24 04:58:56.204/3687.144 Oracle Coherence GE 12.2.1.4.0 <D6> (thread=Cluster, member=1): TcpRing connected to Member(Id=5, Timestamp=2019-04-24 04:58:55.99, Address=10.1.4.147:32923, MachineId=30443, Location=site:coherence.sample-coherence-ns.svc.cluster.local,machine:docker-for-desktop,process:6020,member:storage-coherence-0, Role=CoherenceConsole)
   2019-04-24 04:58:56.480/3687.420 Oracle Coherence GE 12.2.1.4.0 <D6> (thread=Transport:TransportService, member=1): Registered Connection {Peer=tmb://10.1.4.147:32923.64682, Service=TransportService, Member=5, Not established, State=CONNECTING, peer=tmb://10.1.4.147:32923.64682, state=OPEN, socket=MultiplexedSocket{Socket[addr=/10.1.4.147,port=32923,localport=57374]}, bytes(in=0, out=0), flushlock false, bufferedOut=0B, unflushed=0B, delivered(in=0, out=0), timeout(n/a), interestOps=0, unflushed receipt=0, receiptReturn 0, isReceiptFlushRequired false, bufferedIn(), msgs(in=0, out=0/0)}
   ```   
   
   You can see that there are now level D6 messages being displayed. These particular messages related to the cluster member (console)
   joining the cluster.
   
1. Retrieve the current expiryDelay for all members

   ```bash
   curl http://127.0.0.1:30000/management/coherence/cluster/services/PartitionedCache/caches/test/members?fields=expiryDelay 2> /dev/null | json_pp | grep expiryDelay
   "expiryDelay" : 0,
   "expiryDelay" : 0,
   "expiryDelay" : 0,
   "expiryDelay" : 0
   ```
   
1. Set the expiryDelay for each member to 60000ms

   ```bash
   $ for i in 1 2 3 4; do 
      curl -X POST -H 'Content-type: application/json' http://127.0.0.1:30000/management/coherence/cluster/services/PartitionedCache/caches/test/members/$i -d '{"expiryDelay": 60000}'
   done
   {}{}{}{}
   ```
   
1. Re-query the expiryDelay for all members 

   ```bash
   curl http://127.0.0.1:30000/management/coherence/cluster/services/PartitionedCache/caches/test/members?fields=expiryDelay 2> /dev/null | json_pp | grep expiryDelay
   "expiryDelay" : 60000,
   "expiryDelay" : 60000,
   "expiryDelay" : 60000,
   "expiryDelay" : 60000
   ```
   
   > **Note**: You can also update `highUnits` in a similar way to to `expiryDelay`.

## Uninstalling the Charts

Carry out the following commands to delete the chart installed in this sample.

```bash
$ helm delete storage --purge
```

Before starting another sample, ensure that all the pods are gone from previous samples.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above. 
      
  
