def testStep(String additionalArgument) {
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
        sh '''
            for i in test-cop-$BUILD_NUMBER test-cop2-$BUILD_NUMBER
            do
                kubectl create namespace $i || true
                kubectl create secret docker-registry coherence-k8s-operator-development-secret \
                    --namespace $i \
                    --docker-server=$DOCKER_SERVER \
                    --docker-username=$DOCKER_USERNAME \
                    --docker-password="$DOCKER_PASSWORD" \
                    --docker-email=$DOCKER_EMAIL || true
                kubectl create secret docker-registry ocr-k8s-operator-development-secret \
                    --namespace $i \
                    --docker-server=$OCR_DOCKER_SERVER \
                    --docker-username=$OCR_DOCKER_USERNAME \
                    --docker-password="$OCR_DOCKER_PASSWORD" \
                    --docker-email=$OCR_DOCKER_EMAIL || true
            done
        '''
        withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
            sh """
                export HELM_BINARY=`which helm`
                export KUBECTL_BINARY=`which kubectl`
                mvn -Dbedrock.helm="\$HELM_BINARY" \
                    -Dk8s.kubectl="\$KUBECTL_BINARY" \
                    -Dop.image.pull.policy=Always \
                    -Dci.build=$BUILD_NUMBER \
                    -Dk8s.image.pull.secret=coherence-k8s-operator-development-secret,ocr-k8s-operator-development-secret \
                    -Dk8s.create.namespace=false \
                    $additionalArgument \
                    -P pushTestImage -P helm-test clean install
            """
        }
    }
}

def archiveAndCleanup() {
    archiveArtifacts onlyIfSuccessful: false, allowEmptyArchive: true, artifacts: 'functional-tests/target/test-output/**/*,functional-tests/target/surefire-reports/**/*,functional-tests/target/failsafe-reports/**/*'
    sh '''
        helm delete --purge $(helm ls --namespace test-cop-$BUILD_NUMBER --short) || true
        kubectl delete namespace test-cop-$BUILD_NUMBER  || true
        kubectl delete namespace test-cop2-$BUILD_NUMBER || true
    '''
}

pipeline {
    agent none
    environment {
        HTTP_PROXY  = credentials('coherence-operator-http-proxy')
        HTTPS_PROXY = credentials('coherence-operator-https-proxy')
        NO_PROXY    = credentials('coherence-operator-no-proxy')
    }
    options {
        buildDiscarder logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '28', numToKeepStr: '')
        timeout(time: 4, unit: 'HOURS')
    }
    stages {
        stage('maven-build') {
            agent {
              label 'Kubernetes'
            }
            steps {
                echo 'Maven Build'
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
                        withCredentials([
                            string(credentialsId: 'coherence-operator-docker-password', variable: 'DOCKER_PASSWORD'),
                            string(credentialsId: 'coherence-operator-docker-username', variable: 'DOCKER_USERNAME'),
                            string(credentialsId: 'coherence-operator-docker-server',   variable: 'DOCKER_SERVER')]) {
                            withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                                sh '''
                                    docker login $DOCKER_SERVER -u $DOCKER_USERNAME -p $DOCKER_PASSWORD
                                    mvn -B -Dmaven.test.skip=true -P docker -P dockerPush clean install
                                '''
                            }
                        }
                    }
                }
                stage('kubernetes-tests') {
                    steps {
                        script {
                            testStep('')
                        }
                    }
                    post {
                        always {
                            script {
                                archiveAndCleanup()
                            }
                        }
                    }
                }
                stage('kubernetes-tests-latestCoherenceReleasedImage') {
                    steps {
                        script {
                            testStep('-P testLatestCoherenceReleasedImage')
                        }
                    }
                    post {
                        always {
                            script {
                                archiveAndCleanup()
                            }
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
