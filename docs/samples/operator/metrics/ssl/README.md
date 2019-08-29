# Enable SSL for Metrics

By default when metrics are enabled via configuring Prometheus Operator as described [here](../enable-metrics/README.md), metrics utilize standard HTTP.

This sample shows you how to enable SSL for metrics capture only when configuring an external Prometheus to scrape the metrics.  

> **Note:** It is not supported to enable SSL for metrics using the out of the box Prometheus installed with the Coherence Operator.

> **Note:** Use of Prometheus and Grafana is available only when using the operator with Oracle Coherence 12.2.1.4.0 version.

[Return to Metrics samples](../) / [Return to Coherence Operator samples](../../) / [Return to samples](../../../README.md#list-of-samples)

## Sample files

* [src/main/java/com/oracle/coherence/examples/SampleMetricsClient.java](src/main/java/com/oracle/coherence/examples/SampleMetricsClient.java) -
  Client connect to management over REST through SSL
  
* [src/main/java/com/oracle/coherence/examples/HttpSSLHelper.java](src/main/java/com/oracle/coherence/examples/HttpSSLHelper.java) -
  Client connect to management over REST through SSL
  
## Prerequisites

Ensure you have already installed the Coherence Operator using the instructions [here](../../../README.md#install-the-coherence-operator).

## Installation Steps

1. Change to the `samples/operator/metrics/ssl` directory and ensure you have your Maven build environment set for JDK 8 and build the project:

   ```bash
   $ mvn clean compile
   ```
   
   > **Note**: This sample uses self-signed certificates and simple passwords. They are for sample purposes only and must **NOT** use these in a production environment. You must use and generate proper certificates with appropriate passwords.

1. Create the SSL secret:

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
      --set cluster=metrics-ssl-cluster \
      --set imagePullSecrets=sample-coherence-secret \
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
   
   > *Note:* If your version of the Coherence Operator does not default to using Coherence 12.2.1.4.0, then you need to replace `your-12.2.1.4.0-Coherence-image` with an appropriate 12.2.1.4.0 image.
   
1. Confirm that SSL is applied:

   ```bash
   $ kubectl logs storage-coherence-0 --namespace sample-coherence-ns | grep SSLSocketProviderDependencies
   ```
   ```console
   2019-06-17 02:15:01.525/11.176 Oracle Coherence GE 12.2.1.4.0 <D5> (thread=main, member=1): instantiated SSLSocketProviderDependencies: SSLSocketProvider(auth=two-way, identity=SunX509/file:/coherence/certs/metrics/icarus.jks, trust=SunX509/file:/coherence/certs/metrics/truststore-guardians.jks)
   ```
   
1. Start port forward for the metrics port:

   ```bash
   $ kubectl port-forward storage-coherence-0 -n sample-coherence-ns 9612:9612
   ```
   ```console
   Forwarding from [::1]:9612 -> 9612
   Forwarding from 127.0.0.1:9612 -> 9612
   ```   

1. (Optional) Configure Prometheus.

   Follow the instructions [here](../../../operator/metrics/own-prometheus/README.md) to configure Prometheus to point to your SSL endpoints.

## Uninstall the Charts

Use the following command to delete the two charts installed in this sample:

```bash
$ helm delete storage --purge
```

Delete the secret using the following:

```bash
$ kubectl delete secret ssl-secret --namespace sample-coherence-ns
```

Before starting another sample, ensure that all the pods are removed from previous sample.

If you want to remove the `coherence-operator`, then use the `helm delete` command.
