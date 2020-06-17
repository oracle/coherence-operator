# Coherence Operator Deployment Example

This example showcases how to deploy Coherence applications using the Coherence Operator.

The following scenarios are covered:

1. Installing the Coherence Operator 
1. Installing a Coherence cluster 
1. Deploying a Proxy tier
1. Deploying an storage-disabled application 
1. Enabling Active Persistence

After the initial install of the Coherence cluster, the following examples 
build on the previous ones by issuing a `kubectl apply` to modify
the install adding additional roles.

You can use `kubectl create` for any of the examples to install that one directly.

# Table of Contents

* [Prerequisites](#prerequisites)
  * [Coherence Operator Quick Start](#coherence-operator-quick-start) 
  * [JDK and Maven Versions](#jdk-and-maven-versions) 
  * [Install Coherence into local Maven repository](#install-coherence-into-local-maven-repository)
  * [Create the example namespace](#create-the-example-namespace)
  * [Get Coherence Docker image](#get-coherence-docker-image)
  * [Create a secret](#create-a-secret)
  * [Clone the GitHub repository](#clone-the-github-repository)
  * [Install the Coherence Operator](#install-the-coherence-operator)
  * [Port Forward Grafana and Kibana ports](#port-forward-grafana-and-kibana-ports)
* [Run the Examples](#run-the-examples)
  * [Example 1 - Coherence cluster only](#example-1---coherence-cluster-only)
  * [Example 2 - Adding a Proxy role ](#example-2---adding-a-proxy-role) 
  * [Example 3 - Adding a User application role](#example-3---adding-a-user-application-role)
  * [Example 4 - Enabling Persistence](#example-4---enabling-persistence)
* [Cleaning Up](#cleaning-up)  

# Prerequisites

## Coherence Operator Quick Start

Ensure you have followed all the [Quick Start Guide](https://oracle.github.io/coherence-operator/docs/#/about/03_quickstart) including the
prerequisites and have been able to successfully install the Coherence Operator and a Coherence Cluster.

## JDK and Maven versions

Ensure you have the following software installed:

* [Java 8+ JDK](http://jdk.java.net/)
* [Maven](https://maven.apache.org) version 3.5+
* [Docker](https://docs.docker.com/install/) version 17.03+.
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) version v1.12.0+ .
* Access to a Kubernetes v1.12.0+ cluster.

> Note: You can use a later version of Java, for example, JDK11, as the
>`maven.compiler.source` and `target` are set to JDK 8 in the example `pom.xml` files.

> Note: Ensure that your local Kubernetes is enabled. If you are running Docker using Docker Desktop, select Enable Kubernetes in the Settings menu.

## Install Coherence into local Maven repository

1. Download and install Oracle Coherence 12.2.1.4 from [Oracle Technology Network](https://www.oracle.com/technetwork/middleware/coherence/downloads/index.html).

1. Ensure that the `COHERENCE_HOME` environment variable is set to point to the `coherence` directory under your install location containing the bin, lib, and doc directories. This is required only for the Maven `install-file` commands.

1. Install Coherence into your local Maven repository:

   ```bash
   mvn install:install-file -Dfile=$COHERENCE_HOME/lib/coherence.jar   \
       -DpomFile=$COHERENCE_HOME/plugins/maven/com/oracle/coherence/coherence/12.2.1/coherence.12.2.1.pom
  
   mvn install:install-file -Dfile=$COHERENCE_HOME/lib/coherence-metrics.jar \
       -DpomFile=$COHERENCE_HOME/plugins/maven/com/oracle/coherence/coherence-metrics/12.2.1/coherence-metrics.12.2.1.pom
   ```

## Create the example namespace

You need to create the namespace for the first time to run any of the examples. Create your target namespace:

```bash
kubectl create namespace coherence-example

namespace/coherence-example created
```

> Note: In the examples, a Kubernetes namespace called `coherence-example` is used. 
> If you want to change this namespace, ensure that you change any references to this namespace 
> to match your selected namespace when running the examples.

## Get Coherence Docker image

Get the Coherence Docker image from the Oracle Container Registry:

1. In a web browser, navigate to [Oracle Container Registry](https://container-registry.oracle.com) and click Sign In.
1. Enter your Oracle credentials or create an account if you don't have one.
1. Search for `coherence` in the Search Oracle Container Registry field.
1. Click coherence in the search result list.
1. On the Oracle Coherence page, select the language from the drop-down list and click Continue.
1. Click `Accept` on the Oracle Standard Terms and Conditions page.

> Note: You may pull this image and place in your own repository or pull directly from Oracle Container Registry.

## Create a secret

If all of your images can be pulled from public repositories, this step is not
required. Otherwise, you need to enable your Kubernetes cluster to pull
images from private repositories. You must create a secret to convey
the docker credentials to Kubernetes. In the examples,
the secret named `coherence-example-secret` is used in the namespace `coherence-example`.

```bash
kubectl -n coherence-example \
  create secret docker-registry coherence-example-secret \
  --docker-server=$DOCKER_REPO \                              
  --docker-username=$DOCKER_USERNAME \                        
  --docker-password=$DOCKER_PASSWORD \                        
  --docker-email=$DOCKER_EMAIL                                
```

* Replace `$DOCKER_REPO` with the name of the Docker repository that the images are to be pulled from.
* Replace `$DOCKER_USERNAME` with your username for that repository.
* Replace `$DOCKER_PASSWORD` with your password for that repository.
* Replace `$DOCKER_EMAIL` with your email (or even a fake email).

See the [Kubernetes documentation](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/)
on pull secrets for more details.

## Clone the GitHub repository

The examples exist in a directory `examples` in the Coherence Operator GitHub repository - https://github.com/oracle/coherence-operator.

Clone the repository:

```bash
git clone https://github.com/oracle/coherence-operator

cd coherence-operator/examples
```

In the `examples` root directory, check the pom.xm) and verify that the value of the `coherence.version` property 
matches the version of Coherence that you are actually using. For example, if you have Coherence 12.2.1.4.0, then the value of `coherence.version` must
be `12.2.1-4-0`.  If this value needs adjustment, use the `-Dcoherence.version=` argument for all invocations of `mvn`.

Ensure you have Docker running and your Maven and JDK11 build environment set and use the 
following command to build the projects and associated Docker images:

```bash
mvn clean install -P docker
```

This will result in the following Docker image being created which contains the configuration and server-side 
artifacts to be use by all deployments.

```console
deployment-example:2.1.2
```
## Install the Coherence Operator

For this example we are going to install the Coherence Operator with metrics enables as well as 
EFK (Elasticsearch, FluentD and Kiabana) configured for log capture and analysis in Kibana.

If you are short on resources, then you can set `installEFK=false` and `prometheusoperator.enabled=false`.

> Note: If this is the first time you have installed with Prometheus enabled, you should set `createCustomResource=true`.   


```bash
helm install coherence/coherence-operator \
--set prometheusoperator.enabled=true \
--set prometheusoperator.prometheusOperator.createCustomResource=false \
--set installEFK=true \
--namespace coherence-example \
--set imagePullSecrets[0].name=coherence-example-secret \
--name coherence-operator
```    

Confirm the operator and all associated pods are running:

```bash
kubectl get pods -n coherence-example

NAME                                                     READY   STATUS    RESTARTS   AGE
coherence-operator-65bc5c7b49-9l5br                      1/1     Running   0          80s
coherence-operator-grafana-8454698bcf-5dqvm              2/2     Running   0          80s
coherence-operator-kube-state-metrics-6dc8675d87-jz4d7   1/1     Running   0          80s
coherence-operator-prometh-operator-58d94ffbb8-lhgsn     1/1     Running   0          80s
coherence-operator-prometheus-node-exporter-lxtvs        1/1     Running   0          80s
elasticsearch-f978d6fdd-69nfd                            1/1     Running   0          80s
kibana-9964496fd-rbslk                                   1/1     Running   0          80s
prometheus-coherence-operator-prometh-prometheus-0       3/3     Running   1          66s
```
      
## Port Forward Grafana and Kibana ports

If you have enabled metrics or ELK integration you can port-forward the ports for these components using the scripts
in the `examples/bin/` directory.

*   Port-forward Grafana for viewing metrics
  
    ```bash
    ./port-forward-grafana.sh coherence-example   

    Port-forwarding coherence-operator-grafana-8454698bcf-5dqvm in coherence-example
    Open the following URL: http://127.0.0.1:3000/d/coh-main/coherence-dashboard-main
    Forwarding from 127.0.0.1:3000 -> 3000
    Forwarding from [::1]:3000 -> 3000
    ``` 
      
*   Port-forward Kibana for viewing log messages
   
    ```bash
    ./port-forward-kibana.sh coherence-example       
  
    Port-forwarding kibana-67c4f74ffb-nspwz in coherence-example
    Open the following URL: http://127.0.0.1:5601/
    Forwarding from 127.0.0.1:5601 -> 5601
    Forwarding from [::1]:5601 -> 5601
    ```            

# Run the Examples

Change to the `examples/deployment` directory to run the following commands.

## Example 1 - Coherence cluster only

The first example uses the yaml file [src/main/yaml/example-cluster.yaml](src/main/yaml/example-cluster.yaml) which
defines a single role `storage` which will store cluster data. The configuration under the main `spec` entry will be the defaults for all
roles unless they are overridden by a specific role.

This saves duplication of role configuration leading to less configuration errors.

1.  Install the Coherence cluster `storage` role

    ```bash   
    kubectl -n coherence-example create -f src/main/yaml/example-cluster.yaml 
    ```       

    List the created `coherenceclusters` and `coherenceroles` using:

    ```bash
    kubectl -n coherence-example get coherence    
    
    NAME                                                    AGE
    coherencecluster.coherence.oracle.com/example-cluster   18s
    
    NAME                                                         AGE
    coherencerole.coherence.oracle.com/example-cluster-storage   18s
    ```           

1.  View the running pods

    ```bash  
    kubectl -n coherence-example get pod -l coherenceCluster=example-cluster
    
    NAME                        READY   STATUS    RESTARTS   AGE
    example-cluster-storage-0   1/1     Running   0          3m12s
    example-cluster-storage-1   1/1     Running   0          3m12s
    example-cluster-storage-2   1/1     Running   0          3m12s
    ```
   
    > Note: You may also use `kubectl get pods -n coherence-example` to view all pods including those from the Coherence Operator.

1.  Connect to the Coherence Console inside the cluster to add data

    > Note: Since we cannot yet access the cluster via Coherence*Extend, we will connect via Coherence console to add data. 
   
    ```bash
    kubectl exec -it -n coherence-example example-cluster-storage-0 bash /scripts/startCoherence.sh console
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
    kubectl -n coherence-example scale coherencerole/example-cluster-storage --replicas=6
    ```

    Use the following to verify all 6 nodes are Running and READY before continuing.
    
    ```bash
    kubectl -n coherence-example get pod -l coherenceCluster=example-cluster  
    
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
    kubectl -n coherence-example scale coherencerole/example-cluster-storage --replicas=3
    ```                               
    
    By using the following, you will see that the number of members will gradually scale back to 
    3 during which the is done in a `safe` manner ensuring no data loss.
    
    ```bash
    kubectl -n coherence-example get pod -l coherenceCluster=example-cluster   
    
    NAME                        READY   STATUS        RESTARTS   AGE
    example-cluster-storage-0   1/1     Running       0          19m
    example-cluster-storage-1   1/1     Running       0          19m
    example-cluster-storage-2   1/1     Running       0          19m
    example-cluster-storage-3   1/1     Running       0          3m41s
    example-cluster-storage-4   0/1     Terminating   0          3m41s                             
    ```

## Example 2 - Adding a Proxy role 

The second example uses the yaml file [src/main/yaml/example-cluster-proxy.yaml](src/main/yaml/example-cluster-proxy.yaml) which
adds a new role `proxy` to allow for Coherence*Extend connections via a Proxy server.

The snippet of yaml added below shows:

* A port called `proxy` being exposed on 20000
* The role being set as storage-disabled
* A different cache config being used which will start a Proxy Server. See [here](src/main/resources/conf/proxy-cache-config.xml) for details

```yaml
    - role: proxy
      replicas: 1
      ports:
        - name: proxy
          port: 20000
      coherence:
        cacheConfig: proxy-cache-config.xml
        storageEnabled: false
```

1.  Install the `proxy` role

    ```bash
    kubectl -n coherence-example apply -f src/main/yaml/example-cluster-proxy.yaml 
    
    kubectl get coherenceroles -n coherence-example
    
    NAME                      AGE
    example-cluster-proxy     79s
    example-cluster-storage   7m46s
    ```      

1.  View the running pods

    ```bash  
    kubectl -n coherence-example get pod -l coherenceCluster=example-cluster
    
    NAME                        READY   STATUS    RESTARTS   AGE
    example-cluster-proxy-0     1/1     Running   0          28s
    example-cluster-storage-0   1/1     Running   0          34m
    example-cluster-storage-1   1/1     Running   0          34m
    example-cluster-storage-2   1/1     Running   0          34m
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

The snippet of yaml added below shows:

* A port called `http` being exposed on 8080 for the application
* The role being set as storage-disabled
* Using the storage-cache-config.xml but as storage-disabled
* An alternate main class to run - `com.oracle.coherence.examples.Main`

```yaml
    - role: rest
      replicas: 1
      ports:
        - name: http
          port: 8080
      coherence:
        cacheConfig: storage-cache-config.xml
        storageEnabled: false
      application:
        main: com.oracle.coherence.examples.Main
```

1.  Install the `helidon-app' role

    ```bash
    kubectl -n coherence-example apply -f src/main/yaml/example-cluster-app.yaml  
    
    kubectl get coherenceroles -n coherence-example
    
    NAME                          AGE
    example-cluster-rest          11s
    example-cluster-proxy         10m
    example-cluster-storage       22m
    ```      

1.  View the running pods

    ```bash  
    kubectl -n coherence-example get pod -l coherenceCluster=example-cluster
    
    NAME                            READY   STATUS    RESTARTS   AGE
    example-cluster-rest-0          1/1     Running   0          5s
    example-cluster-proxy-0         1/1     Running   0          5m1s
    example-cluster-storage-0       1/1     Running   0          5m3s
    example-cluster-storage-1       1/1     Running   0          5m3s
    example-cluster-storage-2       1/1     Running   0          5m3s
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
    Date: Thu, 18 Apr 2019 06:48:15 GMT
    transfer-encoding: chunked
    connection: keep-alive
    
    curl -i -w '\n' -X PUT http://127.0.0.1:8080/query -d '{"query":"insert into foo key(\"foo\") value(\"bar\")"}'
    
    HTTP/1.1 200 OK
    Date: Thu, 18 Apr 2019 06:48:40 GMT
    transfer-encoding: chunked
    connection: keep-alive
    
    curl -i -w '\n' -X PUT http://127.0.0.1:8080/query -d '{"query":"select key(),value() from foo"}' 
    
    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Thu, 18 Apr 2019 06:49:15 GMT
    transfer-encoding: chunked
    connection: keep-alive
    
    {"result":"{foo=[foo, bar]}"}        
    
    curl -i -w '\n' -X PUT http://127.0.0.1:8080/query -d '{"query":"create cache test"}'
    
    HTTP/1.1 200 OK
    Date: Thu, 18 Apr 2019 06:48:15 GMT
    transfer-encoding: chunked
    connection: keep-alive
    
    curl -i -w '\n' -X PUT http://127.0.0.1:8080/query -d '{"query":"select count() from test"}'
    
    HTTP/1.1 200 OK
    Content-Type: application/json
    Date: Thu, 10 Oct 2019 03:56:18 GMT
    transfer-encoding: chunked
    connection: keep-alive
    
    {"result":"10001"}
    ```                    
    
## Example 4 - Enabling Persistence

The fourth example uses the yaml file [src/main/yaml/example-cluster-persistenceyaml](src/main/yaml/example-cluster-persistence.yaml) which
enabled Active Persistence for the `storage` role by adding a `persistence:` element.
  
The snippet of yaml added below shows:

* Active Persistence being enabled via `persistence.enabled=true`
* Various Persistence Volume Claim (PVC) values being set under `persistentVolumeClaim`

```yaml
  roles:
    - role: storage
      replicas: 3
      coherence:
        cacheConfig: storage-cache-config.xml
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
    kubectl -n coherence-example get pod -l coherenceCluster=example-cluster
    
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
    kubectl exec -it -n coherence-example example-cluster-storage-0 bash /scripts/startCoherence.sh console
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
    kubectl exec -it -n coherence-example example-cluster-storage-0 bash /scripts/startCoherence.sh console
    ```

    At the prompt type the following to create a cache called `test`:
    
    ```bash
    cache test
    ```
    
    Lastly issue the command `size` to verify the cache entry count is 10,000 meaning the data has been recovered.
    
    Type `bye` to exit the console.
  
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
    helm delete coherence-operator --purge
    ```