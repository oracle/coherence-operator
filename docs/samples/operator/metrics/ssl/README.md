# Enable SSL for Metrics

By default when metrics are enabled via configuring Prometheus Operator as described 
[here](../enable-metrics/README.md), metrics utilize standard HTTP. 

This sample shows you how to enable SSL for metrics capture only when configuring an
external Prometheus to scrape the metrics.  

> **Note:** It is not supported to enable SSL for metrics using the out of the box Prometheus 
> installed with the Coherence Operator. 

> **Note:** Use of Prometheus and Grafana is only available when using the
> operator with Coherence 12.2.1.4.

[Return to Metrics samples](../) / [Return to Coherence Operator samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/java/com/oracle/coherence/examples/SampleMetricsClient.java](src/main/java/com/oracle/coherence/examples/SampleMetricsClient.java) -
  client to connect to management over rest via SSL
  
* [src/main/java/com/oracle/coherence/examples/HttpSSLHelper.java](src/main/java/com/oracle/coherence/examples/HttpSSLHelper.java) -
  client to connect to management over rest via SSL
  
## Prerequisites

Ensure you have already installed the Coherence Operator by using the instructions [here](../../../README.md#install-the-coherence-operator).

> **Note:** You do *not() need to enable metrics capture via `--set prometheusoperator.enabled=true` as
> an external Prometheus will be used.

## Installation Steps

1. Change to the `samples/operator/metrics/ssl` directory and ensure you have your maven build environment set for JDK8 and build the project.

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
      --set cluster=metrics-ssl-cluster \
      --set imagePullSecrets=sample-coherence-secret \
      --set prometheusoperator.enabled=true \
      --set logCaptureEnabled=false \
      --set store.metrics.ssl.enabled=true \
      --set store.metrics.ssl.secrets=ssl-secret \
      --set store.metrics.ssl.keyStore=icarus.jks \
      --set store.metrics.ssl.keyStorePasswordFile=storepassword.txt \
      --set store.metrics.ssl.keyPasswordFile=keypassword.txt \
      --set store.metrics.ssl.keyStoreType=JKS \
      --set store.metrics.ssl.trustStore=truststore-guardians.jks \
      --set store.metrics.ssl.trustStorePasswordFile=trustpassword.txt \
      --set store.metrics.ssl.trustStoreType=JKS \
      --set store.metrics.ssl.requireClientCert=true \
      --set coherence.image=your-12.2.1.4.0-Coherence-image \
      coherence/coherence
   ```
   
   > *Note:* If your version of the Coherence Operator does not default to using Coherence
   > 12.2.1.4.0, then you will need to replace `your-12.2.1.4.0-Coherence-image` with an
   > appropriate 12.2.1.4.0 image.
   
1. Confirm SSL is applied

   ```bash
   $ kubectl logs storage-coherence-0 --namespace sample-coherence-ns | grep SSLSocketProviderDependencies
   2019-06-17 02:15:01.525/11.176 Oracle Coherence GE 12.2.1.4.0 <D5> (thread=main, member=1): instantiated SSLSocketProviderDependencies: SSLSocketProvider(auth=two-way, 
           identity=SunX509/file:/coherence/certs/metrics/icarus.jks,
           trust=SunX509/file:/coherence/certs/metrics/truststore-guardians.jks)
   ```
   
1. Port-Forward the metrics port

   ```bash
   $ kubectl port-forward storage-coherence-0 -n sample-coherence-ns 9095:9095
   Forwarding from [::1]:9095 -> 9095
   Forwarding from 127.0.0.1:9095 -> 9095
   ```   

1. (Optionally)Configure your Prometheus

Carry out the instructions [here](../../../operator/metrics/own-prometheus/README.md) to configure
Prometheus to point to your SSL endpoints.   
   
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

   