# Coherence Operator Deployment Example

This example showcases how to deploy Coherence applications using the Coherence Operator.

The following scenarios are covered:

1. Installing the Coherence Operator 
1. Installing a Coherence cluster 
1. Deploying a Proxy tier
1. Deploying a storage-disabled application 
1. Enabling Active Persistence
1. Viewing Metrics via Grafana

After the initial install of the Coherence cluster, the following examples 
build on the previous ones by issuing a `kubectl apply` to modify
the install adding additional roles.

You can use `kubectl create` for any of the examples to install that one directly.

# Table of Contents

* [Prerequisites](#prerequisites)
  * [Coherence Operator Quick Start](#coherence-operator-quick-start) 
  * [Software Versions](#software-versions) 
  * [Create the example namespace](#create-the-example-namespace)
  * [Clone the GitHub repository](#clone-the-github-repository)
  * [Install the Coherence Operator](#install-the-coherence-operator)
  * [Port Forward and Access Grafana](#port-forward-and-access-grafana)
* [Run the Examples](#run-the-examples)
  * [Example 1 - Coherence cluster only](#example-1---coherence-cluster-only)
  * [Example 2 - Adding a Proxy role ](#example-2---adding-a-proxy-role) 
  * [Example 3 - Adding a User application role](#example-3---adding-a-user-application-role)
  * [Example 4 - Enabling Persistence](#example-4---enabling-persistence)
* [View Cluster Metrics via Grafana](#view-cluster-metrics-via-grafana)  
* [Cleaning Up](#cleaning-up)  

# Prerequisites

## Coherence Operator Quick Start

Ensure you have followed all the [Quick Start Guide](https://oracle.github.io/coherence-operator/docs/3.1.0/#/about/03_quickstart) including the
prerequisites and have been able to successfully install the Coherence Operator and a Coherence Cluster.

## Software Versions

Ensure you have the following software installed:

* Java 11+ JDK either [OpenJDK](https://adoptopenjdk.net/) or [Oracle JDK](https://www.oracle.com/java/technologies/javase-downloads.html)
* [Maven](https://maven.apache.org) version 3.6.0+
* [Docker](https://docs.docker.com/install/) version 17.03+.
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) version v1.13.0+ but ideally the latest version .
* Access to a Kubernetes v1.13.0+ cluster.
* [Helm](https://helm.sh/docs/intro/install/) version 3.2.4+ (2.14.3+ supported, see notes) 

> Note: Ensure that your local Kubernetes is enabled. If you are running Docker using Docker Desktop, select Enable Kubernetes in the Settings menu.

## Create the example namespace

You need to create the namespace for the first time to run any of the examples. Create your target namespace:

```bash
kubectl create namespace coherence-example

namespace/coherence-example created
```

> Note: In the examples, a Kubernetes namespace called `coherence-example` is used. 
> If you want to change this namespace, ensure that you change any references to this namespace 
> to match your selected namespace when running the examples.

## Clone the GitHub repository

The examples exist in the `examples` directory in the Coherence Operator GitHub repository - https://github.com/oracle/coherence-operator.

Clone the repository:

```bash
git clone https://github.com/oracle/coherence-operator

cd coherence-operator/examples
```

In the `examples` root directory, check the POM (`pom.xml`) and verify that the value of the `coherence.version` property 
matches the version of Coherence that you are actually using. For example, the  default value is
Coherence Community Edition (CE) version `14.1.1-0-1` and if you wish to change this then set the value of `coherence.version` 
use the `-Dcoherence.version=` argument for all invocations of `mvn`.

Ensure you have Docker running and your Maven and JDK11 build environment set and use the 
following command to build the projects and associated Docker images:

```bash
mvn clean install -P docker
```            

> Note: If you want to change the Coherence version from the default of 14.1.1-0-1 to 20.06, then 
> you can supply the -Dcoherence.version=20.06 for the mvn commands below. You must also use JDK11.

> Note: If you are running behind a corporate proxy and receive the following message building the 
> Docker image:
> `Connect to gcr.io:443 [gcr.io/172.217.212.82] failed: connect timed out` you must modify the build command 
> to add the proxy hosts and ports to be used by the `jib-maven-plugin` as shown below:
>
> ```bash
> mvn clean install -P docker -Dhttps.proxyHost=host -Dhttps.proxyPort=80 -Dhttp.proxyHost=host -Dhttp.proxyPort=80
> ```

This will result in the following Docker image being created which contains the configuration and server-side 
artifacts to be use by all deployments.

```console
deployment-example:3.1.1
```   

> Note: If you are running against a remote Kubernetes cluster, you need to tag and 
> push the Docker image to your repository accessible to that cluster. 
> You also need to prefix the image name in the `yaml` files below.

## Install the Coherence Operator

Issue the following command to install the Coherence Operator:

```bash 
helm install --namespace coherence-example coherence-operator coherence/coherence-operator
```

> Note: for Helm version 2, use the following:

```bash
helm install coherence/coherence-operator --namespace coherence-example --name coherence-operator
```

Confirm the operator is running:

```bash
kubectl get pods -n coherence-example

NAME                                  READY   STATUS    RESTARTS   AGE
coherence-operator-578497bb5b-w89kt   1/1     Running   0          29s
```
      
# Run the Examples

Change to the `examples/deployment` directory to run the following commands.

## Example 1 - Coherence cluster only

The first example uses the yaml file [src/main/yaml/example-cluster.yaml](src/main/yaml/example-cluster.yaml) which
defines a single role `storage` which will store cluster data. 

1.  Install the Coherence cluster `storage` role

    ```bash   
    kubectl -n coherence-example create -f src/main/yaml/example-cluster.yaml 
    ```       

    List the created Coherence cluster

    ```bash
    kubectl -n coherence-example get coherence    
    
    NAME                      CLUSTER           ROLE                      REPLICAS   READY   PHASE
    example-cluster-storage   example-cluster   example-cluster-storage   2                  Created
    
    NAME                                                         AGE
    coherencerole.coherence.oracle.com/example-cluster-storage   18s
    ```           

1.  View the running pods

    ```bash  
    kubectl -n coherence-example get pods
    
    NAME                        READY   STATUS    RESTARTS   AGE
    example-cluster-storage-0   1/1     Running   0          3m12s
    example-cluster-storage-1   1/1     Running   0          3m12s
    example-cluster-storage-2   1/1     Running   0          3m12s
    ```
   
    > Note: You may also use `kubectl get pods -n coherence-example` to view all pods including those from the Coherence Operator.

1.  Connect to the Coherence Console inside the cluster to add data

    > Note: Since we cannot yet access the cluster via Coherence*Extend, we will connect via Coherence console to add data. 
   
    ```bash
    kubectl exec -it -n coherence-example example-cluster-storage-0 /coherence-operator/utils/runner console
    ```

    At the prompt type the following to create a cache called `test`:
    
    ```bash
    cache test
    ```
    
    Use the following to create 10,000 entries of 100 bytes:
    
    ```bash
    bulkput 10000 100 0 100
    ```        
    
    Lastly issue the command `size` to verify the cache entry count.
    
    Type `bye` to exit the console.
    
1.  Scale the `storage` role to 6 members

    ```bash
    kubectl -n coherence-example scale sts/example-cluster-storage --replicas=6
    ```

    Use the following to verify all 6 nodes are Running and READY before continuing.
    
    ```bash
    kubectl -n coherence-example get pods
    
    NAME                        READY   STATUS    RESTARTS   AGE
    example-cluster-storage-0   1/1     Running   0          17m
    example-cluster-storage-1   1/1     Running   0          17m
    example-cluster-storage-2   1/1     Running   0          17m
    example-cluster-storage-3   1/1     Running   0          89s
    example-cluster-storage-4   1/1     Running   0          89s
    example-cluster-storage-5   1/1     Running   0          89s
    ```
    
1.  Confirm the cache count

    Re-run step 3 above and just use the `cache test` and `size` commands to confirm the number of entries is still 10,000.
  
    This confirms that the scale-out was done in a `safe` manner ensuring no data loss. 
    
1.  Scale the `storage` role back to 3 members
   
    ```bash
    kubectl -n coherence-example scale sts/example-cluster-storage --replicas=3
    ```                               
    
    By using the following, you will see that the number of members will gradually scale back to 
    3 during which the is done in a `safe` manner ensuring no data loss.
    
    ```bash
    kubectl -n coherence-example get pods  
    
    NAME                        READY   STATUS        RESTARTS   AGE
    example-cluster-storage-0   1/1     Running       0          19m
    example-cluster-storage-1   1/1     Running       0          19m
    example-cluster-storage-2   1/1     Running       0          19m
    example-cluster-storage-3   1/1     Running       0          3m41s
    example-cluster-storage-4   0/1     Terminating   0          3m41s                             
    ```

## Example 2 - Adding a Proxy role 

The second example uses the yaml file [src/main/yaml/example-cluster-proxy.yaml](src/main/yaml/example-cluster-proxy.yaml) which
adds a proxy server `example-cluster-proxy` to allow for Coherence*Extend connections via a Proxy server.

The additional yaml added below shows:

* A port called `proxy` being exposed on 20000
* The role being set as storage-disabled
* A different cache config being used which will start a Proxy Server. See [here](src/main/resources/proxy-cache-config.xml) for details

```yaml
apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: example-cluster-proxy
spec:
  cluster: example-cluster
  jvm:
    memory:
      heapSize: 512m
  ports:
    - name: metrics
      port: 9612
      serviceMonitor:
        enabled: true
    - name: proxy
      port: 20000
  coherence:
    cacheConfig: proxy-cache-config.xml
    storageEnabled: false
    metrics:
      enabled: true
  image: deployment-example:3.1.1
  imagePullPolicy: Always
  replicas: 1
```

1.  Install the `proxy` role

    ```bash
    kubectl -n coherence-example apply -f src/main/yaml/example-cluster-proxy.yaml 
    
    kubectl get coherence -n coherence-example
    
    NAME                      CLUSTER           ROLE                      REPLICAS   READY   PHASE
    example-cluster-proxy     example-cluster   example-cluster-proxy     1          1       Ready
    example-cluster-storage   example-cluster   example-cluster-storage   3          3       Ready
    ```      

1.  View the running pods

    ```bash  
    kubectl -n coherence-example get pods
    
    NAME                                  READY   STATUS    RESTARTS   AGE
    coherence-operator-578497bb5b-w89kt   1/1     Running   0          68m
    example-cluster-proxy-0               1/1     Running   0          2m41s
    example-cluster-storage-0             1/1     Running   0          29m
    example-cluster-storage-1             1/1     Running   0          29m
    example-cluster-storage-2             1/1     Running   0          2m43s
    ```    
    
    Ensure the `example-cluster-proxy-0` pod is Running and READY before continuing.
    
1.  Port forward the proxy port

    In a separate terminal, run the following:

    ```bash
    kubectl port-forward -n coherence-example example-cluster-proxy-0 20000:20000
    ``` 

1.  Connect via CohQL and add data

    In a separate terminal, change to the `examples/deployments` directory and run the following to 
    start Coherence Query Language (CohQL):
    
    ```bash
    mvn exec:java       
    
    Coherence Command Line Tool

    CohQL> 
    ```
    
    Run the following `CohQL` commands to view and insert data into the cluster.
    
    ```sql 
    CohQL> select count() from 'test';   
    
    Results
    10000
    
    CohQL> insert into 'test' key('key-1') value('value-1');
    
    CohQL> select key(), value() from 'test' where key() = 'key-1';
    Results
    ["key-1", "value-1"]       
    
    CohQL> select count() from 'test';
    Results
    10001       
    
    CohQL> quit
    ```
    
    The above results will show that you can see the data previously inserted and
    can add new data into the cluster using Coherence*Extend.
       
## Example 3 - Adding a User application role 

The third example uses the yaml file [src/main/yaml/example-cluster-app.yaml](src/main/yaml/example-cluster-app.yaml) which
adds a new role `rest`. This role defines a user application which uses [Helidon](https://helidon.io/) to create a
`/query` endpoint allowing the user to send CohQL commands via this endpoint.

The additional yaml added below shows:

* A port called `http` being exposed on 8080 for the application
* The role being set as storage-disabled
* Using the storage-cache-config.xml but as storage-disabled
* An alternate main class to run - `com.oracle.coherence.examples.Main`

```yaml
apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: example-cluster-rest
spec:
  cluster: example-cluster
  jvm:
    memory:
      heapSize: 512m
  ports:
    - name: metrics
      port: 9612
      serviceMonitor:
        enabled: true
    - name: http
      port: 8080
  coherence:
    cacheConfig: storage-cache-config.xml
    storageEnabled: false
    metrics:
      enabled: true
  image: deployment-example:3.1.1
  imagePullPolicy: Always
  application:
    main: com.oracle.coherence.examples.Main
```

1.  Install the `rest` role

    ```bash
    kubectl -n coherence-example apply -f src/main/yaml/example-cluster-app.yaml  
    
    kubectl get coherence -n coherence-example
    
    NAME                      CLUSTER           ROLE                      REPLICAS   READY   PHASE
    example-cluster-proxy     example-cluster   example-cluster-proxy     1          1       Ready
    example-cluster-rest      example-cluster   example-cluster-rest      1          1       Ready
    example-cluster-storage   example-cluster   example-cluster-storage   3          3       Ready
    ```      

1.  View the running pods

    ```bash  
    kubectl -n coherence-example get pods
    
    NAME                              READY   STATUS    RESTARTS   AGE
    coherence-operator-578497bb5b-w89kt   1/1     Running   0          90m
    example-cluster-proxy-0               1/1     Running   0          3m57s
    example-cluster-rest-0                1/1     Running   0          3m57s
    example-cluster-storage-0             1/1     Running   0          3m59s
    example-cluster-storage-1             1/1     Running   0          3m58s
    example-cluster-storage-2             1/1     Running   0          3m58s
    ```    
    
1.  Port forward the application port

    In a separate terminal, run the following:

    ```bash
    kubectl port-forward -n coherence-example example-cluster-rest-0 8080:8080
    ``` 

1.  Access the custom `/query` endpoint

    Use the various `CohQL` commands via the `/query` endpoint to access, and mutate data in the Coherence cluster.

    ```bash
    curl -i -w '\n' -X PUT http://127.0.0.1:8080/query -d '{"query":"create cache foo"}'
    
    HTTP/1.1 200 OK
    Date: Fri, 19 Jun 2020 06:29:40 GMT
    transfer-encoding: chunked
    connection: keep-alive
    
    curl -i -w '\n' -X PUT http://127.0.0.1:8080/query -d '{"query":"insert into foo key(\"foo\") value(\"bar\")"}'
    
    HTTP/1.1 200 OK
    Date: Fri, 19 Jun 2020 06:29:44 GMT
    transfer-encoding: chunked
    connection: keep-alive
    
    curl -i -w '\n' -X PUT http://127.0.0.1:8080/query -d '{"query":"select key(),value() from foo"}' 
    
    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Fri, 19 Jun 2020 06:29:55 GMT
    transfer-encoding: chunked
    connection: keep-alive
    
    {"result":"{foo=[foo, bar]}"}        
    
    curl -i -w '\n' -X PUT http://127.0.0.1:8080/query -d '{"query":"create cache test"}'
    
    HTTP/1.1 200 OK
    Date: Fri, 19 Jun 2020 06:30:00 GMT
    transfer-encoding: chunked
    connection: keep-alive
    
    curl -i -w '\n' -X PUT http://127.0.0.1:8080/query -d '{"query":"select count() from test"}'
    
    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Fri, 19 Jun 2020 06:30:20 GMT
    transfer-encoding: chunked
    connection: keep-alive
    
    {"result":"10001"}
    ```                    
    
## Example 4 - Enabling Persistence

The fourth example uses the yaml file [src/main/yaml/example-cluster-persistence.yaml](src/main/yaml/example-cluster-persistence.yaml) which
enabled Active Persistence for the `storage` role by adding a `persistence:` element.
  
The additional yaml added to the storage role below shows:

* Active Persistence being enabled via `persistence.enabled=true`
* Various Persistence Volume Claim (PVC) values being set under `persistentVolumeClaim`

```yaml
  coherence:
    cacheConfig: storage-cache-config.xml
    metrics:
      enabled: true
    persistence:
      enabled: true
      persistentVolumeClaim:
        accessModes:
          - ReadWriteOnce
        resources:
          requests:
            storage: 1Gi
```   

> Note: By default, when you enable Coherence Persistence, the required infrastructure in terms
> of persistent volumes (PV) and persistent volume claims (PVC) is set up automatically. Also, the persistence-mode
> is set to `active`. This allows the Coherence cluster to be restarted and the data to be retained.

1.  Delete the existing deployment

    We must first delete the existing deployment as we need to redeploy to enable Active Persistence.
    
    ```bash
    kubectl -n coherence-example delete -f src/main/yaml/example-cluster-app.yaml
    ```                                   
    
    Ensure all the pods have terminated before you continue.
                                
1.  Install the cluster with Persistence enabled
  
    ```bash
    kubectl -n coherence-example create -f src/main/yaml/example-cluster-persistence.yaml 
    ```                                                                      

1.  View the running pods and PVC's

    ```bash  
    kubectl -n coherence-example get pods
    
    NAME                            READY   STATUS    RESTARTS   AGE
    example-cluster-rest-0          1/1     Running   0          5s
    example-cluster-proxy-0         1/1     Running   0          5m1s
    example-cluster-storage-0       1/1     Running   0          5m3s
    example-cluster-storage-1       1/1     Running   0          5m3s
    example-cluster-storage-2       1/1     Running   0          5m3s
    ```       
    
    Check the Persistent Volumes and PVC are automatically created.

    ```bash
    kubectl get pvc -n coherence-example
    
    NAME                                           STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
    persistence-volume-example-cluster-storage-0   Bound    pvc-15b46996-eb35-11e9-9b4b-025000000001   1Gi        RWO            hostpath       55s
    persistence-volume-example-cluster-storage-1   Bound    pvc-15bd99e9-eb35-11e9-9b4b-025000000001   1Gi        RWO            hostpath       55s
    persistence-volume-example-cluster-storage-2   Bound    pvc-15e55b6b-eb35-11e9-9b4b-025000000001   1Gi        RWO            hostpath       55s
    ```                                                                                                                                             

    Wait until all  nodes are Running and READY before continuing.

1.  Check Active Persistence is enabled

    Use the following to view the logs of the `example-cluster-storage-0` pod and validate that
    Active Persistence is enabled.
    
    ```bash
    kubectl logs example-cluster-storage-0 -c coherence -n coherence-example | grep 'Created persistent'
    
    ...
    019-10-10 04:52:00.179/77.023 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=DistributedCache:PartitionedCache, member=4): Created persistent store /persistence/active/example-cluster/PartitionedCache/126-2-16db40199bc-4
    2019-10-10 04:52:00.247/77.091 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=DistributedCache:PartitionedCache, member=4): Created persistent store /persistence/active/example-cluster/PartitionedCache/127-2-16db40199bc-4
    ...
    ```   
    
    If you see output similar to above then Active Persistence is enabled.

1.  Connect to the Coherence Console to add data

    ```bash
    kubectl exec -it -n coherence-example example-cluster-storage-0 /coherence-operator/utils/runner console
    ```

    At the prompt type the following to create a cache called `test`:
    
    ```bash
    cache test
    ```
    
    Use the following to create 10,000 entries of 100 bytes:
    
    ```bash
    bulkput 10000 100 0 100
    ```        
    
    Lastly issue the command `size` to verify the cache entry count.
    
    Type `bye` to exit the console.

1.  Delete the cluster

    > Note: This will not delete the PVC's.

    ```bash
    kubectl -n coherence-example delete -f src/main/yaml/example-cluster-persistence.yaml 
    ```       
      
    Use `kubectl get pods -n coherence-example` to confirm the pods have terminated.

1.  Confirm the PVC's are still present

    ```bash
    kubectl get pvc -n coherence-example 
    
    NAME                                           STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
    persistence-volume-example-cluster-storage-0   Bound    pvc-730f86fe-eb19-11e9-9b4b-025000000001   1Gi        RWO            hostpath       116s
    persistence-volume-example-cluster-storage-1   Bound    pvc-73191751-eb19-11e9-9b4b-025000000001   1Gi        RWO            hostpath       116s
    persistence-volume-example-cluster-storage-2   Bound    pvc-73230889-eb19-11e9-9b4b-025000000001   1Gi        RWO            hostpath       116s
    ```       

1.  Re-install the cluster

    ```bash
    kubectl -n coherence-example create -f src/main/yaml/example-cluster-persistence.yaml 
    ```               

1.  Follow the logs for Persistence messages

    ```bash
    kubectl logs example-cluster-storage-0 -c coherence -n coherence-example -f
    ```

    You should see a message regarding recovering partitions, similar to the following:

    ```console
    2019-10-10 05:00:14.255/32.206 Oracle Coherence GE 12.2.1.4.0 <D5> (thread=DistributedCache:PartitionedCache, member=1): Recovering 86 partitions
    ...
    2019-10-10 05:00:17.417/35.368 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=DistributedCache:PartitionedCache, member=1): Created persistent store /persistence/active/example-cluster/PartitionedCache/50-3-16db409d035-1 from SafeBerkeleyDBStore(50-2-16db40199bc-4, /persistence/active/example-cluster/PartitionedCache/50-2-16db40199bc-4)
    ...
    ```    
    
    Finally you should see the following indicating active recovery has completed.
    
    ```console
    2019-10-10 08:18:04.870/59.565 Oracle Coherence GE 12.2.1.4.0 <Info> (thread=DistributedCache:PartitionedCache, member=1): 
       Recovered PartitionSet{172..256} from active persistent store
    ```
 
1.  Confirm the data has been recovered

     ```bash
    kubectl exec -it -n coherence-example example-cluster-storage-0 /coherence-operator/utils/runner console
    ```

    At the prompt type the following to create a cache called `test`:
    
    ```bash
    cache test
    ```
    
    Lastly issue the command `size` to verify the cache entry count is 10,000 meaning the data has been recovered.
    
    Type `bye` to exit the console.
  
# View Cluster Metrics via Grafana

If you wish to view metrics via Grafana, you must carry out the following steps before you
install any of the examples above.

## Install Prometheus Operator

1. Add the `stable` helm repository

    ```bash
   helm repo add stable https://kubernetes-charts.storage.googleapis.com/
   
   helm repo update 
   ```

1. Create Prometheus pre-requisites

    ```bash
    kubectl apply -f src/main/yaml/prometheus-rbac.yaml  
    ```
     
1. Create Config Maps for datasource and dashboards

    ```bash 
    kubectl -n coherence-example create -f src/main/yaml/grafana-datasource-config.yaml  
   
    kubectl -n coherence-example label configmap demo-grafana-datasource grafana_datasource=1  

    kubectl -n coherence-example create -f https://oracle.github.io/coherence-operator/dashboards/3.1.0/coherence-grafana-dashboards.yaml

    kubectl -n coherence-example label configmap coherence-grafana-dashboards grafana_dashboard=1
    ```        

1. Install Prometheus Operator

    > Note: If you have already installed Prometheus Operator before on this Kubernetes Cluster
    > then set `--set prometheusOperator.createCustomResource=false`.

    Issue the following command to install the Prometheus Operator using Helm:
    
    ```bash
    helm install --namespace coherence-example --version 8.13.9 \
        --set grafana.enabled=true \
        --set prometheusOperator.createCustomResource=true \
        --values src/main/yaml/prometheus-values.yaml prometheus stable/prometheus-operator  
    ```        
   
    > Note: for Helm version 2, use the following:
    
    ```bash
    helm install --namespace coherence-example --version 8.13.9 \
        --set grafana.enabled=true --name prometheus \
        --set prometheusOperator.createCustomResource=true \
        --values src/main/yaml/prometheus-values.yaml stable/prometheus-operator 
    ```
   
    Use the following to view the installed pods:
    
    ```bash
    kubectl get pods -n coherence-example
   
    NAME                                                   READY   STATUS    RESTARTS   AGE
    coherence-operator-578497bb5b-w89kt                    1/1     Running   0          136m
    prometheus-grafana-6bb6d86f86-rgsm6                    2/2     Running   0          85s
    prometheus-kube-state-metrics-5496457bd-vjqgd          1/1     Running   0          85s
    prometheus-prometheus-node-exporter-29lrp              1/1     Running   0          85s
    prometheus-prometheus-node-exporter-82b5w              1/1     Running   0          85s
    prometheus-prometheus-node-exporter-mbj2k              1/1     Running   0          85s
    prometheus-prometheus-oper-operator-6bc97bc4d7-67qjp   2/2     Running   0          85s
    prometheus-prometheus-prometheus-oper-prometheus-0     3/3     Running   1          68s 
   ```                
   
## Port Forward and Access Grafana

Port-forward the ports for these components using the scripts
in the `examples/bin/` directory.

*   Port-forward Grafana for viewing metrics
  
    ```bash
    ./port-forward-grafana.sh coherence-example   

    Port-forwarding coherence-operator-grafana-8454698bcf-5dqvm in coherence-example
    Open the following URL: http://127.0.0.1:3000/d/coh-main/coherence-dashboard-main
    Forwarding from 127.0.0.1:3000 -> 3000
    Forwarding from [::1]:3000 -> 3000
    ```      
    
    The default username is `admin` and the password is `prom-operator`.
      
*   Port-forward Kibana for viewing log messages
   
    ```bash
    ./port-forward-kibana.sh coherence-example       
  
    Port-forwarding kibana-67c4f74ffb-nspwz in coherence-example
    Open the following URL: http://127.0.0.1:5601/
    Forwarding from 127.0.0.1:5601 -> 5601
    Forwarding from [::1]:5601 -> 5601
    ```            

# Cleaning Up

1.  Delete the cluster

    ```bash
    kubectl -n coherence-example delete -f src/main/yaml/example-cluster-persistence.yaml  
    ```     

1.  Delete the PVC's

    Ensure all the pods have all terminated before you delete the PVC's.

    ```bash
    kubectl get pvc -n coherence-example | sed 1d | awk '{print $1}' | xargs kubectl delete pvc -n coherence-example
    ```

1.  Remove the Coherence Operator

    ```bash
    helm delete coherence-operator --namespace coherence-example
    ```                                                         
    
    > For Helm version 2 use the following:
                                                                                                                                                                                          
    ```bash
    helm delete coherence-operator --purge
    ```     
    
1. Delete Prometheus Operator

   ```bash
   helm delete prometheus --namespace coherence-example
   
   kubectl -n coherence-example delete -f src/main/yaml/grafana-datasource-config.yaml
   
   kubectl -n coherence-example delete configmap coherence-grafana-dashboards
   
   kubectl delete -f src/main/yaml/prometheus-rbac.yaml      
   ```    
   
   > For Helm version 2 use the following:
   
   ```bash
   helm delete prometheus --purge
   ```                                                                                                                                                                                                                                                                                                                                                                 
   
   > Note: You can optionally delete the Prometheus Operator Custom Resource Definitions
   > (CRD's) if you are not going to install Prometheus Operator again. 
   
   ```bash
   $ kubectl delete crd alertmanagers.monitoring.coreos.com 
   $ kubectl delete crd podmonitors.monitoring.coreos.com
   $ kubectl delete crd prometheuses.monitoring.coreos.com
   $ kubectl delete crd prometheusrules.monitoring.coreos.com 
   $ kubectl delete crd prometheusrules.monitoring.coreos.com 
   $ kubectl delete crd servicemonitors.monitoring.coreos.com 
   $ kubectl delete crd thanosrulers.monitoring.coreos.com 
   ```   
   
   A shorthand way of doing this if you are running Linux/Mac is:
   ```bash
   kubectl get crds -n coherence-example | grep monitoring.coreos.com | awk '{print $1}' | xargs kubectl delete crd
   ``` 

    
