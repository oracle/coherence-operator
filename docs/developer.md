# Coherence Operator Developer Guide

This document describes:

* how to build the operator, without running any tests.

* how to locally "install" the built artifacts so that the operator
  can be trialed.  Built artifacts include:

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

* Obtain the Coherence 12.2.1.3.2 Docker image and tag it correctly.

   1. Download [Oracle Coherence 12.2.1.3.0 Standalone](https://www.oracle.com/technetwork/middleware/coherence/downloads/index.html).  Download the `Coherence Stand-Alone Install`.
   
   2. Git Clone the Oracle `docker-images` repository.
   
      `git clone git@github.com:fryp/docker-images.git`

   3. Within that repository, cd to `OracleCoherence/dockerfiles/12.2.1.3.0`.
   
   4. Make it so the `fmw_12.2.1.3.0_coherence_Disk1_1of1.zip`
      downloaded in step 1 is in that directory.
      
   5. Build the docker image and tag it as
      `oracle/coherence:12.2.1.3.0-standalone`.
      
      `docker build -f Dockerfile.standalone -t oracle/coherence:12.2.1.3.0-standalone`
      
   6. Verify that it built correctly
   
      `docker images | grep 12.2.1.3.0-standalone`
      
      This should show output similar to the following:
      
      `oracle/coherence 12.2.1.3.0-standalone c6dbeed01b35 22 seconds ago 622MB`
      
   7. cd to `OracleCoherence/samples/122132-patch` within the
      `docker-images` cloned repository.
      
   8. Follow the steps in [these
      instructions](https://github.com/fryp/docker-images/blob/master/OracleCoherence/samples/122132-patch/README.md)
      to create a Coherence 12.2.1.3.2 docker image.  
      
   9. Obtain the image hash for the resultant Docker image.

   `docker images | grep 12.2.1.3.2` 

   For discussion, let's call this `COHERENCE_IMAGE_HASH`.

   `docker tag COHERENCE_IMAGE_HASH YOUR_test.image.prefix_VALUE/oracle/coherence:12.2.1.3.2`

   After this command successfully completes, you must be able to say
   
   `docker images | grep 12.2.1.3.2` 

   and see the expected COHERENCE_IMAGE_HASH.  For example:

   ```
   YOUR_test.image.prefix_VALUE/oracle/coherence 12.2.1.3.2 7e7feca04384 2 months ago 547MB
   ```

* `mvn -DskipTests clean install`

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
   In addition the the output must contain output similar to the following,
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
   YOUR_test.image.prefix_VALUE/oracle/coherence          12.2.1.3.2       7e7feca04384 2 months ago 547MB
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
