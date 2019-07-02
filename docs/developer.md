# Oracle Coherence Kubernetes Operator Developer Guide

The Developer Guide provides information for developers who want to build, install, and test the operator.

After successfully completing the steps in this guide, you can execute the instructions in the [Quick Start Guide](quickstart.md)
and [User Guide](user-guide.md).

## Prerequisites

Refer to the [Requirements](quickstart.md) section in the Quick Start Guide.
In addition to the requirements defined in the Quick Start Guide, you require the following software versions for the build environment:

* Mac OS 10.13.6
* Docker Desktop 2.0.0.3 (31259).  Channel: stable, 8858db33c8, with Kubernetes
  v1.10.11.
* Oracle JDK 11.0.1 2018-10-16 LTS
* Apache Maven 3.5.4

> **Note**: You need to make the necessary adjustments to execute the steps in this guide on other operating systems with other Docker
> versions.



### Verify System Environment

Check and verify that your environment is properly configured with the following software for building and installing the operator:

| Software | Verify | Expected Output|
| ---------| ----------------------------|----------------|
| Docker   | `$ docker run hello world`  | `Hello from Docker!`|
| Kubernetes | `$ kubectl version` | `Client Version: version.Info{Major:"1", Minor:"13", GitVersion:"v1.13.3", GitCommit:"721bfa751924da8d1680787490c54b9179b1fed0", GitTreeState:"clean", BuildDate:"2019-02-04T04:49:22Z", GoVersion:"go1.11.5", Compiler:"gc", Platform:"darwin/amd64"}<br>Server Version: version.Info{Major:"1", Minor:"10", GitVersion:"v1.10.11", GitCommit:"637c7e288581ee40ab4ca210618a89a555b6e7e9", GitTreeState:"clean", BuildDate:"2018-11-26T14:25:46Z", GoVersion:"go1.9.3", Compiler:"gc", Platform:"linux/amd64"}` |
| Helm | `$ helm version` | `Client: &version.Version{SemVer:"v2.12.3", GitCommit:"eecf22f77df5f65c823aacd2dbd30ae6c65f186e", GitTreeState:"clean"}<br>Server: &version.Version{SemVer:"v2.12.3", GitCommit:"eecf22f77df5f65c823aacd2dbd30ae6c65f186e", GitTreeState:"clean"}` |
| Java | `java version` | `java version "11.0.1" 2018-10-16 LTS` |
| Maven | `mvn version` | `Apache Maven 3.5.4 (1edded0938998edf8bf061f1ceb3cfdeccf443fe; 2018-06-17T14:33:14-04:00) <br> Maven home: /Users/username/Downloads/apache-maven-3.5.4 <br> Java version: 11.0.1, vendor: Oracle Corporation, runtime: /Library/Java/JavaVirtualMachines/jdk-11.0.1.jdk/Contents/Home` |

## Build the Operator

To build the operator without running any tests, do the following:

1. Clone and check out the current version of the operator from the [GitHub repository](https://github.com/oracle/coherence-operator).
2. Create a maven `settings.xml` file. If you already have one, ensure that the following settings are included in your default profile. All of the maven commands in this guide use this `settings.xml` file.

   ```xml
   <properties>
       <test.image.prefix>DOCKER_REPO_HOSTNAME/DOCKER_REPO_PREFIX/dev/DEV_USERNAME/</test.image.prefix>
   </properties>
   ```
   In the `settings.xml` file:
   * `DOCKER_REPO_HOSTNAME` is the hostname of the docker repo in which you will push your built docker images.  
   >**Note**: You are not required to
     push any images when executing the steps in this document.

   * `DOCKER_REPO_PREFIX` is some prefix within that repo.

   * `DEV_USERNAME` is a username unique to your development environment.
  *  In this example, `YOUR_test.image.prefix_VALUE` is the
   value used for `test.image.prefix` property in the `settings.xml` file.

3. Obtain a Coherence 12.2.1.3.2 Docker image and tag it correctly.
  * Refer to the section [Obtain Images from Oracle Container Registry](quickstart,md) to pull the Coherence Docker image from the Oracle Container Registry.

    `docker pull container-registry.oracle.com/middleware/coherence:12.2.1.3.2`

4. Tag the obtained image in the way it is required to build the operator.  

      1. Obtain the image hash for the Coherence 12.2.1.3.2 Docker image.

         ```bash
         $ docker images | grep 12.2.1.3.2`
         ```

         In this example, it is assumed that the Coherence image hash is `COHERENCE_IMAGE_HASH`.

         ```bash
         docker tag COHERENCE_IMAGE_HASH YOUR_test.image.prefix_VALUE/oracle/coherence:12.2.1.3.2`
         ```

      2. After executing this command, again execute the command to list the docker images which will list the image hash:

         ```bash
         $ docker images | grep 12.2.1.3.2

         YOUR_test.image.prefix_VALUE/oracle/coherence 12.2.1.3.2 7e7feca04384 2 months ago 547MB
         ```
6. From the top level directory of the `coherence-operator` repository in GitHub, do the following.

   ```bash
   $ mvn -DskipTests clean install
   ```

   This produces output similar to the following:

   ```bash
   ...
   [INFO] ------------------------------------------------------------------------
   [INFO] Reactor Summary:
   [INFO]
   [INFO] coherence-operator parent VERSION ........... SUCCESS [  2.487 s]
   [INFO] coherence-operator ................................. SUCCESS [ 21.651 s]
   [INFO] coherence-utils .................................... SUCCESS [ 22.868 s]
   [INFO] coherence-operator-tests VERSION ............ SUCCESS [ 11.468 s]
   [INFO] ------------------------------------------------------------------------
   [INFO] BUILD SUCCESS
   [INFO] ------------------------------------------------------------------------
   [INFO] Total time: 58.756 s
   [INFO] Finished at: 2019-04-17T18:35:14-04:00
   [INFO] ------------------------------------------------------------------------
   ```

   >**Note**: The `VERSION` in the output will be similar to `1.0.0-SNAPSHOT`.

7. `mvn -DskipTests generate-resources`

   This produces output similar to the output of the preceding step.

8. `mvn -DskipTests -Pdocker clean install`

    The output of this command has the following message:
   ```bash
   ...
   Successfully built af61471e4774
   Successfully tagged YOUR_test.image.prefix_VALUE/oracle/coherence-operator:VERSION
   ...
   Successfully built 88495a497a16
   Successfully tagged YOUR_test.image.prefix_VALUE/oracle/coherence-utils:VERSION
   ```

   >**Note**: The `VERSION` in the output will be similar to `1.0.0-SNAPSHOT`.

9. Verify that the docker images have been built and are accessible to your
   local docker server.

   ```bash
   $ docker images | grep YOUR_test.image.prefix_VALUE
   ```

   This produces output similar to the following:

   ```bash
   YOUR_test.image.prefix_VALUE/oracle/coherence-utils    VERSION 88495a497a16 14 minutes ago 124MB
   YOUR_test.image.prefix_VALUE/oracle/coherence-operator VERSION af61471e4774 14 minutes ago 537MB
   YOUR_test.image.prefix_VALUE/oracle/coherence          12.2.1.3.2       7e7feca04384 2 months ago 547MB
   ```

   >**Note**: The `VERSION` in the output will be similar to `1.0.0-SNAPSHOT`.

10. Verify that the Coherence Helm chart and the Coherence Operator Helm chart have been built and are accessible in your work area.

   ```bash
   $ ls -la operator/target | grep "drw" | grep coherence
   ```

   This produces output similar to the following:

   ```bash
   drwxr-xr-x   3 username  staff      96 Apr 17 18:38 coherence-VERSION-helm
   drwxr-xr-x   3 username  staff      96 Apr 17 18:38 coherence-operator-VERSION-helm
   ```
   If you want to use the build image as the source while executing the steps in the [Quick Start Guide](quickstart.md) and [User Guide](user-guide.md), replace the Helm repository prefix with the full qualified path:

   * `coherence/coherence` - `/Users/username/workareas/coherence-operator/target/coherence-1.0.0-SNAPSHOT-helm/coherence`

   * `coherence/coherence-operator` - `/Users/username/workareas/coherence-operator/target/coherence-operator-1.0.0-SNAPSHOT-helm/coherence-operator`

   > **Note:** It is assumed that the Coherence Operator is built within `/Users/username/workareas/coherence-operator`. The `VERSION` in the output will be similar to `1.0.0-SNAPSHOT`.
