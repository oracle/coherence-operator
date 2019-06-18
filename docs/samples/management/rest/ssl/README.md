# Enable SSL with management over REST

By default when the Coherence chart is installed the Management over REST endpoint will be exposed
as port 30000 (via HTTP) on each of the Pods. 

This sample shows how you can access configure the Management over REST endpoint to use SSL.

> **Note**: Use of Management over REST is only available when using the
> operator with Coherence 12.2.1.4.

[Return to Management over REST samples](../)  [Return to Management samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/java/com/oracle/coherence/examples/SampleRESTClient.java](src/main/java/com/oracle/coherence/examples/SampleRESTClient.java) -
  client to connect to management over rest via SSL
  
* [src/main/java/com/oracle/coherence/examples/HttpSSLHelper.java](src/main/java/com/oracle/coherence/examples/HttpSSLHelper.java) -
  client to connect to management over rest via SSL

## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/management/rest/ssl/src/main/resources/certs` directory and ensure you have your maven build environment set for JDK8 and build the project.

   ```bash
   $ mvn clean compile
   ```
   
   > **Note**: This sample uses self-signed certificates and simple passwords. They are for sample
   > purposes only and should **NOT** use these in a production environment.
   > You should use and generate proper certificates with appropriate passwords.

1. Create the SSL Secret 

   ```bash   
   $ cd /src/main/resources/certs

   $ kubectl -n sample-coherence-ns create secret generic ssl-secret \
        --from-file icarus.jks \
        --from-file truststore-guardians.jks \
        --from-literal keypassword.txt=password \
        --from-literal storepassword.txt=password \
        --from-literal trustpassword.txt=secret
   ```

1. Install the Coherence cluster
   
   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=rest-ssl-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=true \
      --set logCaptureEnabled=false \
      --set store.management.ssl.enabled=true \
      --set store.management.ssl.secrets=ssl-secret \
      --set store.management.ssl.keyStore=icarus.jks \
      --set store.management.ssl.keyStorePasswordFile=storepassword.txt \
      --set store.management.ssl.keyPasswordFile=keypassword.txt \
      --set store.management.ssl.keyStoreType=JKS \
      --set store.management.ssl.trustStore=truststore-guardians.jks \
      --set store.management.ssl.trustStorePasswordFile=trustpassword.txt \
      --set store.management.ssl.trustStoreType=JKS \
      --set store.management.ssl.requireClientCert=false \
      --set coherence.image=your-12.2.1.4.0-Coherence-image \
      coherence/coherence
   ```
   
   > *Note:* If your version of the Coherence Operator does not default to using Coherence
   > 12.2.1.4.0, then you will need to replace `your-12.2.1.4.0-Coherence-image` with an
   > appropriate 12.2.1.4.0 image.
   
1. Confirm SSL is applied

   ```bash
   $ kubectl logs storage-coherence-0 --namespace sample-coherence-ns | grep SSLSocketProviderDependencies
   2019-06-17 03:31:13.256/7.535 Oracle Coherence GE 12.2.1.4.0 <D5> (thread=main, member=1): instantiated SSLSocketProviderDependencies: SSLSocketProvider(auth=two-way, 
     identity=SunX509/file:/coherence/certs/management/icarus.jks,
     trust=SunX509/file:/coherence/certs/management/truststore-guardians.jks)
   ```
   
1. Port-Forward the Management over REST port

   ```bash
   $ kubectl port-forward storage-coherence-0 -n sample-coherence-ns 30000:30000
   Forwarding from [::1]:30000 -> 30000
   Forwarding from 127.0.0.1:30000 -> 30000
   ```   

1. Access Management Over REST

   Issue the following command to run the `SampleRESTClient` which will connect via SSL.

   ```bash
   $ mvn exec:java
   ```
   
   This should result in the output below with the additional content from the REST endpoint:
   
   ```bash
   Success, HTTP Response code is 200
   ```
   
## Uninstalling the Charts

Carry out the following commands to delete the two charts installed in this sample.

```bash
$ helm delete storage --purge
```

Delete the secret using the following:

```bash
$ kubectl delete secret ssl-secret --namespace sample-coherence-ns
```

Before starting another sample, ensure that all the pods are gone from previous sample.

If you wish to remove the `coherence-operator`, then include it in the `helm delete` command above.