def setBuildStatus(String message, String state, String project_url, String sha) {
    step([
        $class: "GitHubCommitStatusSetter",
        reposSource: [$class: "ManuallyEnteredRepositorySource", url: project_url],
        contextSource: [$class: "ManuallyEnteredCommitContextSource", context: "ci/jenkins/build-status"],
        errorHandlers: [[$class: "ChangingBuildStatusErrorHandler", result: "UNSTABLE"]],
        commitShaSource: [$class: "ManuallyEnteredShaSource", sha: sha ],
        statusBackrefSource: [$class: "ManuallyEnteredBackrefSource", backref: project_url + "/commit/" + sha],
        statusResultSource: [$class: "ConditionalStatusResultSource", results: [[$class: "AnyBuildResult", message: message, state: state]] ]
    ]);
}

def testStep(String wsdir, String nspId, String additionalArgument) {
    echo 'Kubernetes Tests'
    withCredentials([
        string(credentialsId: 'coherence-operator-docker-email',    variable: 'DOCKER_EMAIL'),
        string(credentialsId: 'coherence-operator-docker-password', variable: 'DOCKER_PASSWORD'),
        string(credentialsId: 'coherence-operator-docker-username', variable: 'DOCKER_USERNAME'),
        string(credentialsId: 'coherence-operator-docker-server',   variable: 'DOCKER_SERVER'),
        string(credentialsId: 'ocr-docker-email',    variable: 'OCR_DOCKER_EMAIL'),
        string(credentialsId: 'ocr-docker-password', variable: 'OCR_DOCKER_PASSWORD'),
        string(credentialsId: 'ocr-docker-username', variable: 'OCR_DOCKER_USERNAME'),
        string(credentialsId: 'ocr-docker-server',   variable: 'OCR_DOCKER_SERVER')]) {
        sh """
            for i in test-cop-$nspId test-cop2-$nspId
            do
                kubectl create namespace \$i || true
                kubectl create secret docker-registry coherence-k8s-operator-development-secret \
                    --namespace \$i \
                    --docker-server=\$DOCKER_SERVER \
                    --docker-username=\$DOCKER_USERNAME \
                    --docker-password="\$DOCKER_PASSWORD" \
                    --docker-email=\$DOCKER_EMAIL || true
                kubectl create secret docker-registry ocr-k8s-operator-development-secret \
                    --namespace \$i \
                    --docker-server=\$OCR_DOCKER_SERVER \
                    --docker-username=\$OCR_DOCKER_USERNAME \
                    --docker-password="\$OCR_DOCKER_PASSWORD" \
                    --docker-email=\$OCR_DOCKER_EMAIL || true
            done
        """
        withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
            sh """
                cd $wsdir
                export HELM_BINARY=`which helm`
                export KUBECTL_BINARY=`which kubectl`
                mvn -Dbedrock.helm="\$HELM_BINARY" \
                    -Dk8s.kubectl="\$KUBECTL_BINARY" \
                    -Dop.image.pull.policy=Always \
                    -Dci.build=$nspId \
                    -Dk8s.image.pull.secret=coherence-k8s-operator-development-secret,ocr-k8s-operator-development-secret \
                    -Dk8s.create.namespace=false \
                    -Dhelm.install.maxRetry=6 \
                    install -P helm-test -pl functional-tests $additionalArgument
            """
        }
    }
}

def archiveAndCleanup() {
    dir (env.WORKSPACE) {
        archiveArtifacts onlyIfSuccessful: false, allowEmptyArchive: true, artifacts: 'build/**/*,deploy/**/*,coherence-utils/utils/target/test-output/**/*,coherence-utils/utils/target/surefire-reports/**/*,coherence-utils/utils/target/failsafe-reports/**/*,coherence-utils/functional-tests/target/test-output/**/*,coherence-utils/functional-tests/target/surefire-reports/**/*,coherence-utils/functional-tests/target/failsafe-reports/**/*'
    }
    //dir (env.WORKSPACE) {
    //    archiveArtifacts onlyIfSuccessful: false, allowEmptyArchive: true, artifacts: 'functional-tests/target/test-output/**/*,functional-tests/target/surefire-reports/**/*,functional-tests/target/failsafe-reports/**/*,_ws2/functional-tests/target/test-output/**/*,_ws2/functional-tests/target/surefire-reports/**/*,_ws2/functional-tests/target/failsafe-reports/**/*'
    //}
    //sh '''
    //    for i in test-cop-$BUILD_NUMBER test-cop2-$BUILD_NUMBER test-cop-${BUILD_NUMBER}-2 test-cop2-${BUILD_NUMBER}-2
    //    do
    //        helm delete --purge $(helm ls --namespace $i --short) || true
    //        kubectl delete namespace $i || true
    //    done
    //'''
}

pipeline {
    agent none
    environment {
        HTTP_PROXY  = credentials('coherence-operator-http-proxy')
        HTTPS_PROXY = credentials('coherence-operator-https-proxy')
        NO_PROXY    = credentials('coherence-operator-no-proxy')
        PROJECT_URL = "https://github.com/oracle/coherence-operator"

        COHERENCE_IMAGE_PREFIX = credentials('coherence-operator-coherence-image-prefix')
        TEST_IMAGE_PREFIX      = credentials('coherence-operator-test-image-prefix')
    }
    options {
        buildDiscarder logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '28', numToKeepStr: '')
        timeout(time: 4, unit: 'HOURS')
    }
    stages {
    /*
        stage('maven-build') {
            agent {
              label 'Kubernetes'
            }
            steps {
                echo 'Maven Build'
                script {
                    setBuildStatus("Build in Progress...", "PENDING", "${env.PROJECT_URL}", "${env.GIT_COMMIT}")
                }
                sh '''
                    if [ -z "$HTTP_PROXY" ]; then
                        unset HTTP_PROXY
                        unset HTTPS_PROXY
                        unset NO_PROXY
                    fi
                    helm init --client-only --skip-refresh
                '''
                withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                   sh 'mvn clean install'
                }
                archiveArtifacts 'operator/target/*.tar.gz'
                stash includes: 'operator/target/*.tar.gz', name: 'helm-chart'
            }
            post {
                failure {
                    setBuildStatus("Build failed", "FAILURE", "${env.PROJECT_URL}", "${env.GIT_COMMIT}");
                }
            }
        }
        stage('helm-verify') {
            agent {
                docker {
                    image 'circleci/python:3.6.4'
                    args '-u root'
                    label 'Docker'
                }
            }
            steps {
                echo 'Helm Verify'
                unstash 'helm-chart'
                sh '''
                    if [ -z "$HTTP_PROXY" ]; then
                        unset HTTP_PROXY
                        unset HTTPS_PROXY
                        unset NO_PROXY
                    else
                        echo "proxy = $HTTP_PROXY" > ~/.curlrc
                        export http_proxy=$HTTP_PROXY
                    fi
                    sh operator/src/main/helm/scripts/install.sh
                    mkdir -p operator/target/temp
                    echo "Contents of operator/target"
                    ls operator/target
                    export COH_CHART=$(find operator/target -regex '.*coherence-[0-9].*-helm.tar.gz' -print)
                    echo COH_CHART=$COH_CHART
                    export COH_OP_CHART=$(find operator/target -regex '.*coherence-operator.*-helm.tar.gz' -print)
                    echo COH_OP_CHART=$COH_OP_CHART
                    tar -xf $COH_CHART -C operator/target/temp
                    tar -xf $COH_OP_CHART -C operator/target/temp
                    sh operator/src/main/helm/scripts/lint.sh operator/target/temp/coherence/
                    sh operator/src/main/helm/scripts/lint.sh operator/target/temp/coherence-operator/
                '''
            }
            post {
                always {
                    sh 'rm -rf operator/target/temp/'
                }
                failure {
                    setBuildStatus("Build failed", "FAILURE", "${env.PROJECT_URL}", "${env.GIT_COMMIT}");
                }
            }
        }
        */
        stage("docker-build-push-tests") {
            agent {
                label 'Kubernetes'
            }
            stages{
                stage('build-utils') {
                    steps {
                        echo 'Docker Build'
                        sh '''
                            if [ -z "$HTTP_PROXY" ]; then
                                unset HTTP_PROXY
                                unset HTTPS_PROXY
                                unset NO_PROXY
                            fi
                            helm init --client-only
                        '''
                        sh 'docker swarm leave --force || true'
                        sh 'docker swarm init'
                        withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                            sh '''
                                cd coherence-utils
                                mvn generate-resources
                                mvn -Pdocker clean install
                            '''
                        }
                    }
                }
                stage('docker-push-utils') {
                    steps {
                        echo 'Docker Push'
                        withCredentials([
                            string(credentialsId: 'coherence-operator-docker-password', variable: 'DOCKER_PASSWORD'),
                            string(credentialsId: 'coherence-operator-docker-username', variable: 'DOCKER_USERNAME'),
                            string(credentialsId: 'coherence-operator-docker-server',   variable: 'DOCKER_SERVER')]) {
                            withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                                sh '''
                                    docker login $DOCKER_SERVER -u $DOCKER_USERNAME -p $DOCKER_PASSWORD
                                    cd coherence-utils
                                    mvn -B -Dmaven.test.skip=true -P docker -P docker-push clean install
                                '''
                            }
                        }
                    }
                }
                stage('docker-build-push-test') {
                    steps {
                        withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                            sh '''
                                cd coherence-utils
                                mvn -B clean install -P helm-test -P push-test-image -Dmaven.test.skip=true
                            '''
                        }
                    }
                }
                stage('build-go-code') {
                    steps {
                        sh '''
                            export http_proxy=$HTTP_PROXY
                            export RELEASE_IMAGE_PREFIX=$(eval echo $TEST_IMAGE_PREFIX)
                            make build
                        '''
                    }
                }
                stage('test-go-code') {
                    steps {
                        sh 'make test'
                    }
                }
                stage('push-operator') {
                    steps {
                        sh '''
                            export http_proxy=$HTTP_PROXY
                            export RELEASE_IMAGE_PREFIX=$(eval echo $TEST_IMAGE_PREFIX)
                            make push
                        '''
                    }
                }
                /*
                stage('copy-workspace') {
                    steps {
                        sh '''
                            rm -rf _ws2
                            mkdir _ws2
                            find . -type d -not -name _ws2 -not -name . -maxdepth 1 -exec cp -R \\{} _ws2 \\;
                            find . -type f -maxdepth 1 -exec cp \\{} _ws2 \\;
                        '''
                    }
                }
                stage('maven-build-2') {
                    steps {
                        withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                            sh '''
                                cd _ws2
                                mvn -B install -P testLatestCoherenceReleasedImage
                            '''
                        }
                    }
                }
                stage('kubernetes-tests-parallel') {
                    parallel {
                        stage('test-1') {
                            steps {
                                testStep('.', env.BUILD_NUMBER, '')
                            }
                        }
                        stage('test-2') {
                            steps {
                                testStep('./_ws2', env.BUILD_NUMBER + '-2', '-P testLatestCoherenceReleasedImage')
                            }
                        }
                    }
                }*/
            }
            post {
                always {
                    script {
                        archiveAndCleanup()
                    }
                    deleteDir()
                }
                success {
                    setBuildStatus("Build succeeded", "SUCCESS", "${env.PROJECT_URL}", "${env.GIT_COMMIT}");
                }
                failure {
                    setBuildStatus("Build failed", "FAILURE", "${env.PROJECT_URL}", "${env.GIT_COMMIT}");
                }
            }
        }
    }
}
