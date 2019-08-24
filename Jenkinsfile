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

def archiveAndCleanup() {
    dir (env.WORKSPACE) {
        archiveArtifacts onlyIfSuccessful: false, allowEmptyArchive: true, artifacts: 'build/**/*,deploy/**/*,coherence-utils/utils/target/test-output/**/*,coherence-utils/utils/target/surefire-reports/**/*,coherence-utils/utils/target/failsafe-reports/**/*,coherence-utils/functional-tests/target/test-output/**/*,coherence-utils/functional-tests/target/surefire-reports/**/*,coherence-utils/functional-tests/target/failsafe-reports/**/*'
        sh '''
            helm delete --purge $(helm ls --namespace $TEST_NAMESPACE --short) || true
            kubectl delete clusterrole test-cop-sw-coherence-operator-cluster-role || true
            kubectl delete clusterrolebinding test-cop-sw-coherence-operator-cluster-role-binding || true
            kubectl delete namespace $TEST_NAMESPACE || true
            make uninstall-crds || true
        '''
    }
}

pipeline {
    agent {
        label 'Kubernetes'
    }
    environment {
        HTTP_PROXY  = credentials('coherence-operator-http-proxy')
        HTTPS_PROXY = credentials('coherence-operator-https-proxy')
        NO_PROXY    = credentials('coherence-operator-no-proxy')
        PROJECT_URL = "https://github.com/oracle/coherence-operator"

        COHERENCE_IMAGE_PREFIX = credentials('coherence-operator-coherence-image-prefix')
        TEST_IMAGE_PREFIX      = credentials('coherence-operator-test-image-prefix')

        TEST_NAMESPACE = "test-cop-${env.BUILD_NUMBER}"
    }
    options {
        buildDiscarder logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '28', numToKeepStr: '')
        timeout(time: 4, unit: 'HOURS')
    }
    stages {
        stage('build-utils') {
            steps {
                echo 'Docker Build utils'
                script {
                    setBuildStatus("Build in Progress...", "PENDING", "${env.PROJECT_URL}", "${env.GIT_COMMIT}")
                }
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
        stage('push-utils') {
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
        stage('build-push-test') {
            steps {
                withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                    sh '''
                        cd coherence-utils
                        mvn -B clean install -P helm-test -P docker -P push-test-image -Dmaven.test.skip=true
                    '''
                }
            }
        }
        stage('build-go') {
            steps {
                sh '''
                    export http_proxy=$HTTP_PROXY
                    export RELEASE_IMAGE_PREFIX=$(eval echo $TEST_IMAGE_PREFIX)
                    export TEST_MANIFEST_VALUES=deploy/oci-values.yaml
                    make build
                '''
            }
        }
        stage('test') {
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
        stage('create-secrets') {
            steps {
                withCredentials([
                    string(credentialsId: 'coherence-operator-docker-email',    variable: 'DOCKER_EMAIL'),
                    string(credentialsId: 'coherence-operator-docker-password', variable: 'DOCKER_PASSWORD'),
                    string(credentialsId: 'coherence-operator-docker-username', variable: 'DOCKER_USERNAME'),
                    string(credentialsId: 'coherence-operator-docker-server',   variable: 'DOCKER_SERVER'),
                    string(credentialsId: 'ocr-docker-email',    variable: 'OCR_DOCKER_EMAIL'),
                    string(credentialsId: 'ocr-docker-password', variable: 'OCR_DOCKER_PASSWORD'),
                    string(credentialsId: 'ocr-docker-username', variable: 'OCR_DOCKER_USERNAME'),
                    string(credentialsId: 'ocr-docker-server',   variable: 'OCR_DOCKER_SERVER')]) {
                    sh '''
                        kubectl create namespace $TEST_NAMESPACE || true
                        kubectl create secret docker-registry coherence-k8s-operator-development-secret \
                            --namespace $TEST_NAMESPACE \
                            --docker-server=$DOCKER_SERVER \
                            --docker-username=$DOCKER_USERNAME \
                            --docker-password="$DOCKER_PASSWORD" \
                            --docker-email=$DOCKER_EMAIL || true
                        kubectl create secret docker-registry ocr-k8s-operator-development-secret \
                            --namespace $TEST_NAMESPACE \
                            --docker-server=$OCR_DOCKER_SERVER \
                            --docker-username=$OCR_DOCKER_USERNAME \
                            --docker-password="$OCR_DOCKER_PASSWORD" \
                            --docker-email=$OCR_DOCKER_EMAIL || true
                    '''
                }
            }
        }
        stage('e2e-local-test') {
            steps {
                sh '''
                    export http_proxy=$HTTP_PROXY
                    export CREATE_TEST_NAMESPACE=false
                    export IMAGE_PULL_SECRETS=coherence-k8s-operator-development-secret,ocr-k8s-operator-development-secret
                    export IMAGE_PULL_POLICY=Always
                    export RELEASE_IMAGE_PREFIX=$(eval echo $TEST_IMAGE_PREFIX)
                    export TEST_MANIFEST_VALUES=deploy/oci-values.yaml
                    make e2e-local-test
                '''
            }
        }
        stage('helm-test') {
            steps {
                sh '''
                    export http_proxy=$HTTP_PROXY
                    export CREATE_TEST_NAMESPACE=false
                    export IMAGE_PULL_POLICY=Always
                    export IMAGE_PULL_SECRETS=coherence-k8s-operator-development-secret,ocr-k8s-operator-development-secret
                    export RELEASE_IMAGE_PREFIX=$(eval echo $TEST_IMAGE_PREFIX)
                    export TEST_MANIFEST_VALUES=deploy/oci-values.yaml
                    make helm-test
                '''
            }
        }
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
