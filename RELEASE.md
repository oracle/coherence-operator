<!--
Copyright 2018, Oracle Corporation and/or its affiliates.
All rights reserved.  Licensed under the Universal
Permissive License v 1.0 as shown at
http://oss.oracle.com/licenses/upl.

-->

The release process consists of 2 parts.

1. Creating a release for GitHub
2. Creating and pushing the Docker images

## Creating a release for GitHub

To create a release for GitHub the following steps need to be done

1. Make sure you are on a freshly checked out project in the branch you want to release from
2. Create a release branch (git checkout -b release)
3. Change the POMs to reflect the release version (mvn versions:set versions:commit)
4. Remove '${test.image.prefix}' from the release.image.prefix in the master POM
5. Commit the changes (git commit)
6. Create a Git tag (git tag vX.Y.Z)
7. Push the Git tag (git push origin vX.Y.Z)
8. Checkout master branch (git checkout master)
9. Remove release branch (git branch -D release)

## Create the Docker images (that can be pushed)

As a prerequisite the release (aka Git tag) should exist in GitHub 

1. Make sure JDK 11 is the active JDK (java -version)
2. Make sure Maven 3.6.0 is the active Maven version (mvn --version)
3. Make sure it is a clean checkout of the coherence-operator project (git clone)
4. Checkout the Git tag, it should be in the format of vX.Y.Z (git checkout vX.Y.Z)
5. Create the docker images (mvn -Pdocker -DskipTests=true -DskipITs=true clean install)
6. Push the created Docker images (docker push ...)
   a. oracle/coherence-operator:x.y.z
   b. oracle/coherence-operator:x.y.z-utils
