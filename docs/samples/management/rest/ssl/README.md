# Enable SSL with Management over REST Endpoint

By default, when the Coherence chart is installed, the Management over REST endpoint is exposed at port 30000 (through HTTP) on each of the pods.

This sample shows how you can access and configure the management over REST endpoint to use SSL.

> **Note**: Use of management over REST endpoint is available only when using the operator with Oracle Coherence 12.2.1.4.0 version.

[Return to Management over REST samples](../)  [Return to Management samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/java/com/oracle/coherence/examples/SampleRESTClient.java](src/main/java/com/oracle/coherence/examples/SampleRESTClient.java) -
  Client to connect to management over REST through SSL
  
* [src/main/java/com/oracle/coherence/examples/HttpSSLHelper.java](src/main/java/com/oracle/coherence/examples/HttpSSLHelper.java) -
  Client to connect to management over Rest through SSL

## Prerequisites

Ensure you have already installed the Coherence Operator using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/management/rest/ssl/src/main/resources/certs` directory and ensure that you have Maven build environment set for JDK 8 and build the project:

   ```bash
   $ mvn clean compile
   ```
   
   > **Note**: This sample uses self-signed certificates and simple passwords. They are for sample purposes only and must not be used in a production environment. You must use and generate proper certificates with appropriate passwords.

1. Create the SSL Secret:

   ```bash   
   $ cd /src/main/resources/certs

   $ kubectl -n sample-coherence-ns create secret generic ssl-secret \
        --from-file icarus.jks \
        --from-file truststore-guardians.jks \
        --from-literal keypassword.txt=password \
        --from-literal storepassword.txt=password \
        --from-literal trustpassword.txt=secret
   ```

1. Install the Coherence cluster:
   
   ```bash
   $ helm install \
      --namespace sample-coherence-ns \
      --name storage \
      --set clusterSize=3 \
      --set cluster=rest-ssl-cluster \
      --set imagePullSecrets=sample-coherence-secret \
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
   
   > *Note:* If your version of the Coherence Operator does not default to using Oracle Coherence 12.2.1.4.0, then you must replace `your-12.2.1.4.0-Coherence-image` with an appropriate 12.2.1.4.0 image.
   
1. Confirm that the SSL is applied:

   ```bash
   $ kubectl logs storage-coherence-0 --namespace sample-coherence-ns | grep SSLSocketProviderDependencies
   ```
   ```console
   2019-06-17 03:31:13.256/7.535 Oracle Coherence GE 12.2.1.4.0 <D5> (thread=main, member=1): instantiated SSLSocketProviderDependencies: SSLSocketProvider(auth=two-way, 
     identity=SunX509/file:/coherence/certs/management/icarus.jks,
     trust=SunX509/file:/coherence/certs/management/truststore-guardians.jks)
   ```
   
1. Port forward the management over REST endpoint port:

   ```bash
   $ kubectl port-forward storage-coherence-0 -n sample-coherence-ns 30000:30000
   ```
   ```console
   Forwarding from [::1]:30000 -> 30000
   Forwarding from 127.0.0.1:30000 -> 30000
   ```   

1. Access management over REST endpoint

   Use the following command to run the `SampleRESTClient` which connects through SSL:

   ```bash
   $ mvn exec:java
   ```
   
   This results in the output with the additional content from the REST endpoint:
   
   ```console
   Success, HTTP Response code is 200
   ```
   
## Uninstall the Charts

Use the following commands to delete the two charts installed in this sample:

```bash
$ helm delete storage --purge
```

Delete the secret using the following:

```bash
$ kubectl delete secret ssl-secret --namespace sample-coherence-ns
```

Before starting another sample, ensure that all the pods are removed from previous sample.

If you want to remove the `coherence-operator`, then use the `helm delete` command.
