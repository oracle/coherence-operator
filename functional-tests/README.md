# Coherence K8s Operator Integration Testing

This module builds and test the Coherence Kubernetes Operator.

## Docker / Kubernetes setup
- Insecure registries for docker registry if it is necessary
- Manual proxy configurations
- increase memory in docker / kubernetes if your deployment consumes a lot of memory

## Build and run the test
- Set the following Unix environment:
  http_proxy, https_proxy, no_proxy
- Set the following environment variables in your shell:
  http_proxy_host, http_proxy_port, http_no_proxy and CI_BUILD_ID.
  The CI_BUILD_ID is alphanumeric that is used to construct k8s namespaces used.
- If the kubectl / helm is not installed under /usr/local/bin, then one
  need to specify k8s.kubectl and bedrock.helm properties in the mvn command.
- Run the following command to build and run integration tests.
    ```
    mvn -Phelm-test clean install
    ```
  One can run integration tests with user specified kubectl and helm as follows:
    ```
    mvn -Dk8s.kubectl=/usr/bin/kubectl -Dbedrock.helm=/tools/bin/helm -Phelm-test clean install
    ```
- For debugging purposes, the namespaces can be created outside the mvn command
  and tests can be run as follows:
    ```
    mvn -Dk8s.kubectl=/usr/bin/kubectl -Dbedrock.helm=/tools/bin/helm -Dk8s.create.namespace=false -Phelm-test clean install
    ```

    Note: If you set system property 'k8s.create.namespace' to 'false' in above maven command, then first you must need to setup the kubernetes secrets before running the tests using below command (Replace <NAMESPACE> with the your kubernetes namespace in which tests are running)
    ```
    kubectl create secret docker-registry coherence-k8s-operator-development-secret --namespace <NAMESPACE> --docker-server=<SERVER>--docker-username=<USERNAME> --docker-password=<PASSWORD> --docker-email=<EMAIL>
    ```

  A specific test can be run by using it.test environment as follows:
    ```
    mvn -Dit.test=LogHelmChartIT -Dk8s.kubectl=/usr/bin/kubectl -Dbedrock.helm=/tools/bin/helm -Dk8s.create.namespace=false -Phelm-test clean install
    ```
