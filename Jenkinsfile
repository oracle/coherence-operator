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
        junit allowEmptyResults: true, testResults: "pkg/**/test-report.xml,test/**/test-report.xml,build/_output/test-logs/operator-e2e-local-test.xml,build/_output/test-logs/operator-e2e-helm-test.xml,build/_output/test-logs/operator-e2e-test.xml,java/**/surefire-reports/*.xml,java/**/failsafe-reports/*.xml"
        archiveArtifacts onlyIfSuccessful: false, allowEmptyArchive: true, artifacts: 'build/_output/helm-charts/**/*,build/_output/test-logs/**/*,deploy/**/*,java/utils/target/test-output/**/*,java/utils/target/surefire-reports/**/*,java/utils/target/failsafe-reports/**/*,java/functional-tests/target/test-output/**/*,java/functional-tests/target/surefire-reports/**/*,java/functional-tests/target/failsafe-reports/**/*'
        sh '''
            helm delete --purge $(helm ls --namespace $TEST_NAMESPACE --short) || true
            kubectl delete clusterrole $TEST_NAMESPACE-coherence-operator || true
            kubectl delete clusterrolebinding $TEST_NAMESPACE-coherence-operator-cluster || true
            kubectl delete namespace $TEST_NAMESPACE --force --grace-period=0 || true
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

        RELEASE_SUFFIX = """${sh(
                             returnStdout: true,
                             script: 'if [ "${RELEASE_SUFFIX}" == "DATE" ]; then echo $(date -u +%y%m%d%H%M); else echo ${RELEASE_SUFFIX}; fi'
                         )}"""
    }
    options {
        buildDiscarder logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '28', numToKeepStr: '')
        timeout(time: 4, unit: 'HOURS')
    }
    stages {
        stage('code-review') {
            steps {
                echo 'Code Review'
                script {
                    setBuildStatus("Code Review in Progress...", "PENDING", "${env.PROJECT_URL}", "${env.GIT_COMMIT}")
                }
                sh '''
                    if [ -z "$HTTP_PROXY" ]; then
                        unset HTTP_PROXY
                        unset HTTPS_PROXY
                        unset NO_PROXY
                    fi
                '''
                withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                    sh '''
                    make code-review
                    '''
                }
            }
        }
        stage('build') {
            steps {
                echo 'Build'
                script {
                    setBuildStatus("Build in Progress...", "PENDING", "${env.PROJECT_URL}", "${env.GIT_COMMIT}")
                }
                sh '''
                    if [ -z "$HTTP_PROXY" ]; then
                        unset HTTP_PROXY
                        unset HTTPS_PROXY
                        unset NO_PROXY
                    fi
                '''
                withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                    sh '''
                    docker rmi $(docker images -q -f "dangling=true") || true
                    make clean
                    export RELEASE_IMAGE_PREFIX=$(eval echo $TEST_IMAGE_PREFIX)
                    export TEST_MANIFEST_VALUES=deploy/oci-values.yaml
                    make build-all
                    '''
                }
            }
        }
        stage('test') {
            when {
                expression { env.SKIP_TESTS != 'true' }
            }
            steps {
                echo 'Tests'
                script {
                    setBuildStatus("Tests in Progress...", "PENDING", "${env.PROJECT_URL}", "${env.GIT_COMMIT}")
                }
                sh '''
                    if [ -z "$HTTP_PROXY" ]; then
                        unset HTTP_PROXY
                        unset HTTPS_PROXY
                        unset NO_PROXY
                    fi
                '''
                withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                    sh '''
                    export RELEASE_IMAGE_PREFIX=$(eval echo $TEST_IMAGE_PREFIX)
                    export TEST_MANIFEST_VALUES=deploy/oci-values.yaml
                    make test-all
                    '''
                }
            }
        }
        stage('build-images') {
            steps {
                echo 'Build Docker Images'
                script {
                    setBuildStatus("Building Docker images...", "PENDING", "${env.PROJECT_URL}", "${env.GIT_COMMIT}")
                }
                withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                    sh '''
                        export http_proxy=$HTTP_PROXY
                        export RELEASE_IMAGE_PREFIX=$(eval echo $TEST_IMAGE_PREFIX)
                        make build-all-images
                    '''
                }
            }
        }
        stage('push-images') {
            steps {
                echo 'Docker Push'
                script {
                    setBuildStatus("Pushing Docker images...", "PENDING", "${env.PROJECT_URL}", "${env.GIT_COMMIT}")
                }
                withCredentials([
                    string(credentialsId: 'coherence-operator-docker-password', variable: 'DOCKER_PASSWORD'),
                    string(credentialsId: 'coherence-operator-docker-username', variable: 'DOCKER_USERNAME'),
                    string(credentialsId: 'coherence-operator-docker-server',   variable: 'DOCKER_SERVER')]) {
                    withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                        sh '''
                            docker login $DOCKER_SERVER -u $DOCKER_USERNAME -p $DOCKER_PASSWORD
                            export http_proxy=$HTTP_PROXY
                            export RELEASE_IMAGE_PREFIX=$(eval echo $TEST_IMAGE_PREFIX)
                            make push-all-images
                        '''
                    }
                }
            }
        }
        stage('create-secrets') {
            steps {
                echo 'Create K8s secrets'
                script {
                    setBuildStatus("Creating K8s secrets...", "PENDING", "${env.PROJECT_URL}", "${env.GIT_COMMIT}")
                }
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
            when {
                expression { env.SKIP_TESTS != 'true' }
            }
            steps {
                echo 'Operator end-to-end local tests'
                script {
                    setBuildStatus("Running Operator end-to-end local tests...", "PENDING", "${env.PROJECT_URL}", "${env.GIT_COMMIT}")
                }
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
        stage('e2e-test') {
            when {
                expression { env.SKIP_TESTS != 'true' }
            }
            steps {
                echo 'Operator end-to-end tests'
                script {
                    setBuildStatus("Running Operator end-to-end tests...", "PENDING", "${env.PROJECT_URL}", "${env.GIT_COMMIT}")
                }
                sh '''
                    export http_proxy=$HTTP_PROXY
                    export CREATE_TEST_NAMESPACE=false
                    export IMAGE_PULL_POLICY=Always
                    export IMAGE_PULL_SECRETS=coherence-k8s-operator-development-secret,ocr-k8s-operator-development-secret
                    export RELEASE_IMAGE_PREFIX=$(eval echo $TEST_IMAGE_PREFIX)
                    export TEST_MANIFEST_VALUES=deploy/oci-values.yaml
                    make e2e-test
                '''
            }
        }
        stage('helm-test') {
            when {
                expression { env.SKIP_TESTS != 'true' }
            }
            steps {
                echo 'Operator Helm tests'
                script {
                    setBuildStatus("Running Operator Helm tests...", "PENDING", "${env.PROJECT_URL}", "${env.GIT_COMMIT}")
                }
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
        stage('release') {
            when {
                expression { env.RELEASE_ON_SUCCESS == 'true' }
            }
            steps {
                echo 'Release'
                sh '''
                    if [ -z "$HTTP_PROXY" ]; then
                        unset HTTP_PROXY
                        unset HTTPS_PROXY
                        unset NO_PROXY
                    fi
                '''
                withCredentials([
                    string(credentialsId: 'coherence-operator-release-dockerhub-username', variable: 'DOCKER_HUB_USERNAME'),
                    string(credentialsId: 'coherence-operator-release-dockerhub-password', variable: 'DOCKER_HUB_PASSWORD'),
                    string(credentialsId: 'coherence-operator-docker-password', variable: 'DOCKER_PASSWORD'),
                    string(credentialsId: 'coherence-operator-docker-username', variable: 'DOCKER_USERNAME'),
                    string(credentialsId: 'coherence-operator-docker-server',   variable: 'DOCKER_SERVER')]) {
                    withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                        sh '''
                        export RELEASE_IMAGE_PREFIX=$(eval echo $TEST_IMAGE_PREFIX)
                        docker login $DOCKER_SERVER -u $DOCKER_USERNAME -p $DOCKER_PASSWORD
                        docker login -u DOCKER_HUB_USERNAME -p DOCKER_HUB_PASSWORD
                        git config user.name "Coherence Bot"
                        git config user.email coherence-bot_ww@oracle.com
                        make clean
                        make build-all-images VERSION_SUFFIX=${RELEASE_SUFFIX} RELEASE_IMAGE_PREFIX=${RELEASE_IMAGE_REPO}
                        make release RELEASE_DRY_RUN=${DRY_RUN} RELEASE_IMAGE_PREFIX=${RELEASE_IMAGE_REPO} VERSION_SUFFIX=${RELEASE_SUFFIX}
                        '''
                    }
                }
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
