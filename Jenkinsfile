pipeline {
    agent none
    environment {
        HTTP_PROXY  = credentials('coherence-operator-http-proxy')
        HTTPS_PROXY = credentials('coherence-operator-https-proxy')
        NO_PROXY    = credentials('coherence-operator-no-proxy')
    }
    options {
        lock('kubernetes-stage1')
    }
    stages {
        stage('maven-build') {
            agent {
              label 'Kubernetes'
            }
            steps {
                sh 'rm PULL_SECRET_* || true'
                echo 'Maven Build'
                sh 'helm init --client-only --skip-refresh'
                echo 'Starting Maven process'
                withMaven(jdk: 'Jdk11', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
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
                    echo "proxy = $HTTP_PROXY" > ~/.curlrc
                    export http_proxy=$HTTP_PROXY
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
        stage('docker-build') {
            agent {
                label 'Kubernetes'
            }
            steps {
                echo 'Docker Build'
                sh 'helm init --client-only'
                sh 'docker swarm leave --force || true'
                sh 'docker swarm init'
                withMaven(jdk: 'Jdk11', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                    sh 'mvn generate-resources'
                    sh 'mvn -Pdocker clean install'
                }
            }
        }
        stage('docker-push') {
            agent {
                label 'Kubernetes'
            }
            steps {
                echo 'Docker Push'
                sh 'helm init --client-only'
                withMaven(jdk: 'Jdk11', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                    sh 'mvn -B -Dmaven.test.skip=true -P docker -P dockerPush clean install'
                }
            }
        }
        stage('kubernetes-tests') {
            agent {
                label 'Kubernetes'
            }
            steps {
                echo 'Kubernetes Tests'
                withCredentials([
                    string(credentialsId: 'coherence-operator-docker-pull-secret-email',    variable: 'PULL_SECRET_EMAIL'),
                    string(credentialsId: 'coherence-operator-docker-pull-secret-password', variable: 'PULL_SECRET_PASSWORD'),
                    string(credentialsId: 'coherence-operator-docker-pull-secret-username', variable: 'PULL_SECRET_USERNAME'),
                    string(credentialsId: 'coherence-operator-docker-pull-secret-server',   variable: 'PULL_SECRET_SERVER')]) {
                    sh '''
                        helm init --client-only
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
                    '''
                    withMaven(jdk: 'Jdk11', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                        sh '''
                            mvn -Dbedrock.helm=`type -ap helm` \
                                -Dop.image.pull.policy=Always \
                                -Dci.build=$BUILD_NUMBER \
                                -Dk8s.image.pull.secret=coherence-k8s-operator-development-secret \
                                -Dk8s.create.namespace=false \
                                -Dfailsafe.skipAfterFailureCount=1 \
                                -P pushTestImage -P helm-test clean install
                        '''
                    }
                }
            }
            post {
                always {
                    sh '''
                        kubectl delete namespace test-cop-$BUILD_NUMBER  || true
                        kubectl delete namespace test-cop2-$BUILD_NUMBER || true
                        kubectl delete crd --ignore-not-found=true alertmanagers.monitoring.coreos.com   || true
                        kubectl delete crd --ignore-not-found=true prometheuses.monitoring.coreos.com    || true
                        kubectl delete crd --ignore-not-found=true prometheusrules.monitoring.coreos.com || true
                        kubectl delete crd --ignore-not-found=true servicemonitors.monitoring.coreos.com || true
                    '''
                    deleteDir()
                }
            }
        }
    }
}
