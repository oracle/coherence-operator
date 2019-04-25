# Coherence Operator Developer Guide

This document describes:

* how to build the operator, without running any tests.

* how to locally "install" the built artifacts so that the operator
  can be tried out.  Built artifacts include:

   * Docker images

   * Helm charts

* how to validate that the build was successful.

Upon successfully completing these steps, you should be able to
successfully execute the instructions in the [quickstart](quickstart.md)
and [user-guide](user-guide.md).

## Prerequisites

* Ensure the prerequisites [listed in the quickstart](quickstart.md#prerequisites) are all met.

These instructions have been validated with the following software and
versions.

* Mac OS 10.13.6

* Docker Desktop 2.0.0.3 (31259).  Channel: stable, 8858db33c8, with Kubernetes
  v1.10.11.

* Oracle JDK 11.0.1 2018-10-16 LTS

* Apache Maven 3.5.4

> You will need to make the necessary adjustments to execute the steps in
> this developer guide on other operating systems with other Docker
> versions.

### Prerequisite validation

The following commands allow you to validate the satisfaction of the
prerequisites.  Failure of any of these validation steps means that you
will not be able to successfully perform the steps in the remainder of
this developer guide.  You must make it so you can perform these
validation steps before continuing.

#### 1 Validate Docker is running successfully

`docker run hello-world`

This should produce output similar to the following:

```
Hello from Docker!
This message shows that your installation appears to be working correctly.
...
```

#### 2 Validate Kubernetes is Installed Correctly

`kubectl version`

This should produce output similar to the following:

```
Client Version: version.Info{Major:"1", Minor:"13", GitVersion:"v1.13.3", GitCommit:"721bfa751924da8d1680787490c54b9179b1fed0", GitTreeState:"clean", BuildDate:"2019-02-04T04:49:22Z", GoVersion:"go1.11.5", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"1", Minor:"10", GitVersion:"v1.10.11", GitCommit:"637c7e288581ee40ab4ca210618a89a555b6e7e9", GitTreeState:"clean", BuildDate:"2018-11-26T14:25:46Z", GoVersion:"go1.9.3", Compiler:"gc", Platform:"linux/amd64"}
```

#### 3 Validate Helm is installed Correctly

`helm version`

This should produce output similar to the following:

```
Client: &version.Version{SemVer:"v2.12.3", GitCommit:"eecf22f77df5f65c823aacd2dbd30ae6c65f186e", GitTreeState:"clean"}
Server: &version.Version{SemVer:"v2.12.3", GitCommit:"eecf22f77df5f65c823aacd2dbd30ae6c65f186e", GitTreeState:"clean"}
```

#### 4 Validate Java is Installed Correctly

`java -version`

This should produce output similar to the following:

```
java version "11.0.1" 2018-10-16 LTS
...
```

#### 5 Validate Maven is Installed Correctly

`mvn -version`

This should produce output similar to the following:

```
Apache Maven 3.5.4 (1edded0938998edf8bf061f1ceb3cfdeccf443fe; 2018-06-17T14:33:14-04:00)
Maven home: /Users/username/Downloads/apache-maven-3.5.4
Java version: 11.0.1, vendor: Oracle Corporation, runtime: /Library/Java/JavaVirtualMachines/jdk-11.0.1.jdk/Contents/Home
...
```

## How to Build the Operator, Without Running Any Tests

* Check out the `1.0` branch of the [GitHub
repository](https://github.com/oracle/coherence-operator).

* If you do not have a maven `settings.xml` file, create one.  If you
have one, make sure the following is included in your `default` profile.
All of the maven commands in this document are assumed to use this
`settings.xml` file.

   ```
   <properties>
       <test.image.prefix>DOCKER_REPO_HOSTNAME/DOCKER_REPO_PREFIX/dev/DEV_USERNAME/</test.image.prefix>
   </properties>
   ```

   `DOCKER_REPO_HOSTNAME` is the hostname of the docker repo that you may
   eventually push your built docker images to.  **You are not required to
   push any images when executing the steps in this document.**

   `DOCKER_REPO_PREFIX` is some prefix within that repo.

   `DEV_USERNAME` is a username unique to your development environment.

   In the remainder of this document, `YOUR_test.image.prefix_VALUE` is the
   value of your `test.image.prefix` property in your `settings.xml` file.

* Obtain a Coherence 12.2.1.3.2 Docker image and tag it correctly.

   1. The process is to get the 12.2.1.3 Docker image and apply a patch
     to derive a Docker image that contains Coherence 12.2.1.3.2.
     First, let's get the 12.2.1.3 Docker image and tag it correctly
     
      1. Go to [store.docker.com](https://store.docker.com/).

      2. Search for "Oracle Coherence".

      3. Choose "Developer Plan (12.2.1.3)".

      4. Choose "Proceed to Checkout".

      5. Create a Docker Id, or log in with it if you have one already.

      6. Check the `I agree that my use of each program in this Content,
         including any subsequent updates or upgrades...` box.

      7. Check the `I acknowledge and allow Docker to share my personal
         information linked to my Docker ID with this Publisher.` box.

      8. Consider whether or not you want to check the `Please keep me
         informed of products, services and solutions from this
         Publisher` box.

      9. At the command line, do `docker login` with your Docker store
         credentials.

      10. At the command line do `docker pull store/oracle/coherence:12.2.1.3`
     
      11. Provide a tag that effectively removes the `store` prefix: `docker store/oracle/coherence:12.2.1.3 oracle/coherence:12.2.1.3`

   2. Now that we have `oracle/coherence:12.2.1.3` in our local Docker
     server, applyt the patch to derive 12.2.1.3.2 from it.
     
      1. Clone the Oracle Docker Images git repository: `git clone git@github.com:oracle/docker-images.git`
      
      2. Change directory to `OracleCoherence/samples/122132-patch-for-k8s`

      3. Follow the steps in [these instructions](https://github.com/oracle/docker-images/blob/master/OracleCoherence/samples/122132-patch-for-k8s/README.md)
         to create a Coherence 12.2.1.3.2 docker image.
         
      
   3. Tag the 12.2.1.3.2 image in the way the operator build expects.  
   
      1. Obtain the image hash for the resultant Docker image.

         `docker images | grep 12.2.1.3.1` 

          For discussion, let's call this `COHERENCE_IMAGE_HASH`.

          `docker tag COHERENCE_IMAGE_HASH YOUR_test.image.prefix_VALUE/oracle/coherence:12.2.1.3.1`

          After this command successfully completes, you must be able to say

          `docker images | grep 12.2.1.3.1` 

          and see the expected COHERENCE_IMAGE_HASH.  For example:

          ```
          YOUR_test.image.prefix_VALUE/oracle/coherence 12.2.1.3.1 7e7feca04384 2 months ago 547MB
          ```
* From the top level directory of the `coherence-operator` repository,
  on the `1.0` branch, do the following.

   `mvn -DskipTests clean install`

   This should produce output similar to the following:

   ```
   ...
   [INFO] ------------------------------------------------------------------------
   [INFO] Reactor Summary:
   [INFO]
   [INFO] coherence-operator parent OPERATOR_VERSION ........... SUCCESS [  2.487 s]
   [INFO] coherence-operator ................................. SUCCESS [ 21.651 s]
   [INFO] coherence-utils .................................... SUCCESS [ 22.868 s]
   [INFO] coherence-operator-tests OPERATOR_VERSION ............ SUCCESS [ 11.468 s]
   [INFO] ------------------------------------------------------------------------
   [INFO] BUILD SUCCESS
   [INFO] ------------------------------------------------------------------------
   [INFO] Total time: 58.756 s
   [INFO] Finished at: 2019-04-17T18:35:14-04:00
   [INFO] ------------------------------------------------------------------------
   ```

   Note that `OPERATOR_VERSION` will actually be something like `1.0.0-SNAPSHOT`.

* `mvn -DskipTests generate-resources`

   This should produce output similar to the output of the preceding step.

* `mvn -DskipTests -Pdocker clean install`

   This should produce output similar to the output of the preceding step.
   In addition the output must contain messages similar to the following,
   somewhere in the middle of the output.

   ```
   ...
   Successfully built af61471e4774
   Successfully tagged YOUR_test.image.prefix_VALUE/oracle/coherence-operator:OPERATOR_VERSION
   ...
   Successfully built 88495a497a16
   Successfully tagged YOUR_test.image.prefix_VALUE/oracle/coherence-utils:OPERATOR_VERSION
   ```

   Note that `OPERATOR_VERSION` will actually be something like `1.0.0-SNAPSHOT`.

* Verify the docker images have been built and are accessible to your
  local docker server.
  
   `docker images | grep YOUR_test.image.prefix_VALUE`

   This should produce output similar to the following:

   ```
   YOUR_test.image.prefix_VALUE/oracle/coherence-utils    OPERATOR_VERSION 88495a497a16 14 minutes ago 124MB
   YOUR_test.image.prefix_VALUE/oracle/coherence-operator OPERATOR_VERSION af61471e4774 14 minutes ago 537MB
   YOUR_test.image.prefix_VALUE/oracle/coherence          12.2.1.3.1       7e7feca04384 2 months ago 547MB
   ```

   Note that `OPERATOR_VERSION` will actually be something like `1.0.0-SNAPSHOT`.

* Verify the Helm charts have been built and are accessible in your
  workarea.
  
   `ls -la operator/target | grep "drw" | grep coherence`

   This should produce output similar to the following:

   ```
   drwxr-xr-x   3 username  staff      96 Apr 17 18:38 coherence-OPERATOR_VERSION-helm
   drwxr-xr-x   3 username  staff      96 Apr 17 18:38 coherence-operator-OPERATOR_VERSION-helm
   ```

   When executing the steps in the [quickstart](quickstart.md) and
   [user-guide](user-guide.md), replace `HELM_PREFIX` with the fully
   qualified path to the parent directory of those two above directories.
   For example:
   `/Users/username/workareas/coherence-operator/operator/target`.

   Note that `OPERATOR_VERSION` will actually be something like `1.0.0-SNAPSHOT`.
