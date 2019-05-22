# Running the samples integration tests

## Prerequisites 

Please ensure that you have met all the prerequisites as described in the
[Samples Readme](README.md#confirm-quickstart-runtime-prerequisites).

## Run the tests

1. Create the test namespace

   ```bash
   export CI_BUILD_ID=test
   export NS=test-sample-${CI_BUILD_ID}

   kubectl create namespace $NS
   ```
   
1. Create the secrets 

   ```bash  
   kubectl create secret docker-registry coherence-k8s-operator-development-secret \
        --namespace $NS \
        --docker-server=your-docker-server \
        --docker-username=your-docker-username \
        --docker-password='your-docker-password'
        
   kubectl create secret docker-registry sample-coherence-secret \
        --namespace $NS \
        --docker-server=your-docker-server \
        --docker-username=your-docker-username \
        --docker-password='your-docker-password'   
   ```
   
   > Note: You must have push permissions for this repository if you not running Kubernetes 
   > locally.
   
1. Run the tests

   ```bash
   cd docs/samples
   mvn -Dcoherence.image.prefix=store/oracle/ -Dk8s.chart.test.versions=0.9.4 \
       -Dk8s.namespace=$NS -Dk8s.create.namespace=false -P docker,helm-test clean verify
   ```   
   
   > Note: If you are running against a remote Kubernetes cluster, you must also specify
   > the profile `dockerPush`
   
   > Note: You can also specify multiple versions of the chart to test: e.g. 
   > `-Dk8s.chart.test.versions=0.9.4,1.0.0`. 