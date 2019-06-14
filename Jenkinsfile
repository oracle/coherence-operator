pipeline {
    agent none
    environment {
        HTTP_PROXY  = credentials('coherence-operator-http-proxy')
        HTTPS_PROXY = credentials('coherence-operator-https-proxy')
        NO_PROXY    = credentials('coherence-operator-no-proxy')
    }
    options {
        lock(label: 'kubernetes-stage', quantity: 1)
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
                   sh '''
		       cd docs/samples 
		       mvn clean install
		   '''
                }
            }
        }
        stage('docker-build') {
            agent {
                label 'Kubernetes'
            }
            steps {
                echo 'Docker Build'
                sh '''
		    docker swarm leave --force || true
                    docker swarm init
                    if [ -z "$HTTP_PROXY" ]; then
                        unset HTTP_PROXY
                        unset HTTPS_PROXY
                        unset NO_PROXY
                    fi
                    helm init --client-only
                '''
                withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                    sh '''
                        cd docs/samples 
			env
		        mvn -P docker,docker-v1,docker-v2,dockerPush clean install
		    '''
                }
            }
        }
        stage('kubernetes-samples-tests') {
            agent {
                label 'Kubernetes'
            }
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
                        if [ -z "$HTTP_PROXY" ]; then
                            unset HTTP_PROXY
                            unset HTTPS_PROXY
                            unset NO_PROXY
                        fi

                        helm init --client-only
			export NS=test-sample-${BUILD_NUMBER}
                        helm repo add coherence https://oracle.github.io/coherence-operator/charts
		        helm repo update
                        kubectl create namespace $NS || true
                        kubectl create secret docker-registry coherence-k8s-operator-development-secret \
                           --namespace $NS \
                           --docker-server=$PULL_SECRET_SERVER \
                           --docker-username=$PULL_SECRET_USERNAME \
                           --docker-password="$PULL_SECRET_PASSWORD" \
                           --docker-email=$PULL_SECRET_EMAIL || true
                        kubectl create secret docker-registry sample-coherence-secret \
                           --namespace $NS \
                           --docker-server=$PULL_SECRET_SERVER \
                           --docker-username=$PULL_SECRET_USERNAME \
                           --docker-password="$PULL_SECRET_PASSWORD" \
                           --docker-email=$PULL_SECRET_EMAIL || true
                        kubectl create secret docker-registry ocr-k8s-operator-development-secret \
                           --namespace $NS \
                           --docker-server=$OCR_PULL_SECRET_SERVER \
                           --docker-username=$OCR_PULL_SECRET_USERNAME \
                           --docker-password="$OCR_PULL_SECRET_PASSWORD" \
                           --docker-email=$OCR_PULL_SECRET_EMAIL || true
                    '''
                    withMaven(jdk: 'JDK 11.0.3', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                        sh '''
			    env
                            export HELM_BINARY=`which helm`
                            export KUBECTL_BINARY=`which kubectl`
                            export NS=test-sample-${BUILD_NUMBER}
		            cd docs/samples 
                            mvn -Dbedrock.helm=''$HELM_BINARY'' \
                                -Dk8s.kubectl=''$KUBECTL_BINARY'' \
                                -Dop.image.pull.policy=Always \
                                -Dci.build=$BUILD_NUMBER \
                                -Dk8s.image.pull.secret=coherence-k8s-operator-development-secret,ocr-k8s-operator-development-secret \
                                -Dk8s.create.namespace=false \
				-Dk8s.chart.test.versions=0.9.8 \
				-Dk8s.namespace=$NS \
                                -P helm-test clean install
                        '''
                    }
                }
            }
            post {
                always {
                    sh '''
		        export NS=test-sample-${BUILD_NUMBER}
                        helm delete --purge $(helm ls --namespace $NS --short) || true
                        kubectl delete namespace $NS || true
                        kubectl delete crd --ignore-not-found=true alertmanagers.monitoring.coreos.com   || true
                        kubectl delete crd --ignore-not-found=true prometheuses.monitoring.coreos.com    || true
                        kubectl delete crd --ignore-not-found=true prometheusrules.monitoring.coreos.com || true
                        kubectl delete crd --ignore-not-found=true servicemonitors.monitoring.coreos.com || true
			helm repo remove coherence
                    '''
                    deleteDir()
                }
            }
        }
    }
}
