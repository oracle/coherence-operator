def setBuildStatus(String message, String state, String target_url, String sha) {
    step([
        $class: "GitHubCommitStatusSetter",
        reposSource: [$class: "ManuallyEnteredRepositorySource", url: "https://github.com/oracle/coherence-operator"],
        contextSource: [$class: "ManuallyEnteredCommitContextSource", context: "ci/jenkins/build-status"],
        errorHandlers: [[$class: "ChangingBuildStatusErrorHandler", result: "UNSTABLE"]],
        commitShaSource: [$class: "ManuallyEnteredShaSource", sha: sha ],
        statusBackrefSource: [$class: "ManuallyEnteredBackrefSource", backref: target_url],
        statusResultSource: [$class: "ConditionalStatusResultSource", results: [[$class: "AnyBuildResult", message: message, state: state]] ]
    ]);
}

pipeline {
    agent none
    environment {
        HTTP_PROXY  = credentials('coherence-operator-http-proxy')
        HTTPS_PROXY = credentials('coherence-operator-https-proxy')
        NO_PROXY    = credentials('coherence-operator-no-proxy')
        PROJECT_URL = "https://github.com/oracle/coherence-operator"
        COMMIT_URL = "${PROJECT_URL}" + "/commit/" + "${GIT_COMMIT}"
    }
    options {
        buildDiscarder logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '28', numToKeepStr: '')
        lock(label: 'kubernetes-stage', quantity: 1)
        timeout(time: 4, unit: 'HOURS')
    }
    stages {
        stage('maven-build') {
            agent {
              label 'Kubernetes'
            }
            post {
                always {
                    setBuildStatus("Build in Progress...", "PENDING", "${COMMIT_URL}" ,"${env.GIT_COMMIT}");
                }
            }
            steps {
                echo 'Maven Build'
                echo "============= COMMIT_URL"
                echo "${env.COMMIT_URL}"
                echo "============= GIT_COMMIT"
                echo "${env.GIT_COMMIT}"
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
                    setBuildStatus("Build failed", "FAILURE", "${COMMIT_URL}" ,"${env.GIT_COMMIT}");
                }
            }
        }
        stage("docker-build-push-tests") {
            agent {
                label 'Kubernetes'
            }
            stages{
                stage('docker-build') {
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
                            sh 'mvn generate-resources'
                            sh 'mvn -Pdocker clean install'
                        }
                    }
                }
                stage('docker-push') {
                    steps {
                        echo 'Docker Push'
                        withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                            sh 'mvn -B -Dmaven.test.skip=true -P docker -P dockerPush clean install'
                        }
                    }
                }
                stage('kubernetes-tests') {
                    steps {
                        echo 'Kubernetes Tests'
                        withCredentials([
                            string(credentialsId: 'coherence-operator-docker-pull-secret-email',    variable: 'PULL_SECRET_EMAIL'),
                            string(credentialsId: 'coherence-operator-docker-pull-secret-password', variable: 'PULL_SECRET_PASSWORD'),
                            string(credentialsId: 'coherence-operator-docker-pull-secret-username', variable: 'PULL_SECRET_USERNAME'),
                            string(credentialsId: 'coherence-operator-docker-pull-secret-server',   variable: 'PULL_SECRET_SERVER'),
                            string(credentialsId: 'ocr-docker-pull-secret-email',    variable: 'OCR_PULL_SECRET_EMAIL'),
                            string(credentialsId: 'ocr-docker-pull-secret-password', variable: 'OCR_PULL_SECRET_PASSWORD'),
                            string(credentialsId: 'ocr-docker-pull-secret-username', variable: 'OCR_PULL_SECRET_USERNAME'),
                            string(credentialsId: 'ocr-docker-pull-secret-server',   variable: 'OCR_PULL_SECRET_SERVER')]) {
                            sh '''
                                kubectl create namespace test-cop-$BUILD_NUMBER  || true
                                kubectl create namespace test-cop2-$BUILD_NUMBER || true
                                kubectl create secret docker-registry coherence-k8s-operator-development-secret \
                                   --namespace test-cop-$BUILD_NUMBER \
                                   --docker-server=$PULL_SECRET_SERVER \
                                   --docker-username=$PULL_SECRET_USERNAME \
                                   --docker-password="$PULL_SECRET_PASSWORD" \
                                   --docker-email=$PULL_SECRET_EMAIL || true
                                kubectl create secret docker-registry coherence-k8s-operator-development-secret \
                                   --namespace test-cop2-$BUILD_NUMBER \
                                   --docker-server=$PULL_SECRET_SERVER \
                                   --docker-username=$PULL_SECRET_USERNAME \
                                   --docker-password="$PULL_SECRET_PASSWORD" \
                                   --docker-email=$PULL_SECRET_EMAIL || true
                                kubectl create secret docker-registry ocr-k8s-operator-development-secret \
                                   --namespace test-cop-$BUILD_NUMBER \
                                   --docker-server=$OCR_PULL_SECRET_SERVER \
                                   --docker-username=$OCR_PULL_SECRET_USERNAME \
                                   --docker-password="$OCR_PULL_SECRET_PASSWORD" \
                                   --docker-email=$OCR_PULL_SECRET_EMAIL || true
                                kubectl create secret docker-registry ocr-k8s-operator-development-secret \
                                   --namespace test-cop2-$BUILD_NUMBER \
                                   --docker-server=$OCR_PULL_SECRET_SERVER \
                                   --docker-username=$OCR_PULL_SECRET_USERNAME \
                                   --docker-password="$OCR_PULL_SECRET_PASSWORD" \
                                   --docker-email=$OCR_PULL_SECRET_EMAIL || true
                            '''
                            withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                                sh '''
                                    export HELM_BINARY=`which helm`
                                    export KUBECTL_BINARY=`which kubectl`
                                    mvn -Dbedrock.helm=''$HELM_BINARY'' \
                                        -Dk8s.kubectl=''$KUBECTL_BINARY'' \
                                        -Dop.image.pull.policy=Always \
                                        -Dci.build=$BUILD_NUMBER \
                                        -Dk8s.image.pull.secret=coherence-k8s-operator-development-secret,ocr-k8s-operator-development-secret \
                                        -Dk8s.create.namespace=false \
                                        -P pushTestImage -P helm-test clean install
                                '''
                            }
                        }
                    }
                    post {
                        always {
                            archiveArtifacts onlyIfSuccessful: false, allowEmptyArchive: true, artifacts: 'functional-tests/target/test-output/**/*,functional-tests/target/surefire-reports/**/*,functional-tests/target/failsafe-reports/**/*'
                            sh '''
                                helm delete --purge $(helm ls --namespace test-cop-$BUILD_NUMBER --short) || true
                                kubectl delete namespace test-cop-$BUILD_NUMBER  || true
                                kubectl delete namespace test-cop2-$BUILD_NUMBER || true
                            '''
                        }
                        failure {
                            setBuildStatus("Build failed", "FAILURE", "${COMMIT_URL}" ,"${env.GIT_COMMIT}");
                        }
                    }
                }
                stage('kubernetes-tests-latestCoherenceReleasedImage') {
                    steps {
                        echo 'Kubernetes Tests'
                        withCredentials([
                            string(credentialsId: 'coherence-operator-docker-pull-secret-email',    variable: 'PULL_SECRET_EMAIL'),
                            string(credentialsId: 'coherence-operator-docker-pull-secret-password', variable: 'PULL_SECRET_PASSWORD'),
                            string(credentialsId: 'coherence-operator-docker-pull-secret-username', variable: 'PULL_SECRET_USERNAME'),
                            string(credentialsId: 'coherence-operator-docker-pull-secret-server',   variable: 'PULL_SECRET_SERVER'),
                            string(credentialsId: 'ocr-docker-pull-secret-email',    variable: 'OCR_PULL_SECRET_EMAIL'),
                            string(credentialsId: 'ocr-docker-pull-secret-password', variable: 'OCR_PULL_SECRET_PASSWORD'),
                            string(credentialsId: 'ocr-docker-pull-secret-username', variable: 'OCR_PULL_SECRET_USERNAME'),
                            string(credentialsId: 'ocr-docker-pull-secret-server',   variable: 'OCR_PULL_SECRET_SERVER')]) {
                            sh '''
                                kubectl create namespace test-cop-$BUILD_NUMBER  || true
                                kubectl create namespace test-cop2-$BUILD_NUMBER || true
                                kubectl create secret docker-registry coherence-k8s-operator-development-secret \
                                   --namespace test-cop-$BUILD_NUMBER \
                                   --docker-server=$PULL_SECRET_SERVER \
                                   --docker-username=$PULL_SECRET_USERNAME \
                                   --docker-password="$PULL_SECRET_PASSWORD" \
                                   --docker-email=$PULL_SECRET_EMAIL || true
                                kubectl create secret docker-registry coherence-k8s-operator-development-secret \
                                   --namespace test-cop2-$BUILD_NUMBER \
                                   --docker-server=$PULL_SECRET_SERVER \
                                   --docker-username=$PULL_SECRET_USERNAME \
                                   --docker-password="$PULL_SECRET_PASSWORD" \
                                   --docker-email=$PULL_SECRET_EMAIL || true
                                kubectl create secret docker-registry ocr-k8s-operator-development-secret \
                                   --namespace test-cop-$BUILD_NUMBER \
                                   --docker-server=$OCR_PULL_SECRET_SERVER \
                                   --docker-username=$OCR_PULL_SECRET_USERNAME \
                                   --docker-password="$OCR_PULL_SECRET_PASSWORD" \
                                   --docker-email=$OCR_PULL_SECRET_EMAIL || true
                                kubectl create secret docker-registry ocr-k8s-operator-development-secret \
                                   --namespace test-cop2-$BUILD_NUMBER \
                                   --docker-server=$OCR_PULL_SECRET_SERVER \
                                   --docker-username=$OCR_PULL_SECRET_USERNAME \
                                   --docker-password="$OCR_PULL_SECRET_PASSWORD" \
                                   --docker-email=$OCR_PULL_SECRET_EMAIL || true
                            '''
                            withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                                sh '''
                                    export HELM_BINARY=`which helm`
                                    export KUBECTL_BINARY=`which kubectl`
                                    mvn -Dbedrock.helm=''$HELM_BINARY'' \
                                        -Dk8s.kubectl=''$KUBECTL_BINARY'' \
                                        -Dop.image.pull.policy=Always \
                                        -Dci.build=$BUILD_NUMBER \
                                        -Dk8s.image.pull.secret=coherence-k8s-operator-development-secret,ocr-k8s-operator-development-secret \
                                        -Dk8s.create.namespace=false \
                                        -P testLatestCoherenceReleasedImage \
                                        -P pushTestImage -P helm-test clean install
                                '''
                            }
                        }
                    }
                    post {
                        always {
                            archiveArtifacts onlyIfSuccessful: false, allowEmptyArchive: true, artifacts: 'functional-tests/target/test-output/**/*,functional-tests/target/surefire-reports/**/*,functional-tests/target/failsafe-reports/**/*'
                            sh '''
                                helm delete --purge $(helm ls --namespace test-cop-$BUILD_NUMBER --short) || true
                                kubectl delete namespace test-cop-$BUILD_NUMBER  || true
                                kubectl delete namespace test-cop2-$BUILD_NUMBER || true
                            '''
                        }
                        success {
                            setBuildStatus("Build succeeded !", "SUCCESS", "${COMMIT_URL}" ,"${env.GIT_COMMIT}");
                        }
                        failure {
                            setBuildStatus("Build failed", "FAILURE", "${COMMIT_URL}" ,"${env.GIT_COMMIT}");
                        }
                    }
                }
            }
            post {
                always {
                    deleteDir()
                }
            }
        }
    }
}
