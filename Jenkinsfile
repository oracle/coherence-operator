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
                echo 'Maven Build'
                sh '''
                    if [ -z "$HTTP_PROXY" ]; then
                        unset HTTP_PROXY
                        unset HTTPS_PROXY
                        unset NO_PROXY
                    fi
                    helm init --client-only --skip-refresh
                '''
                withMaven(jdk: 'Jdk11', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                   sh 'cd docs/samples && mvn clean install'
                }
            }
        }
        stage('docker-build') {
            agent {
                label 'Kubernetes'
            }
            steps {
                echo 'Docker Build'
                sh 'docker swarm leave --force || true'
                sh 'docker swarm init'
                sh '''
                    if [ -z "$HTTP_PROXY" ]; then
                        unset HTTP_PROXY
                        unset HTTPS_PROXY
                        unset NO_PROXY
                    fi
                '''
                withMaven(jdk: 'Jdk11', maven: 'Maven3.6.0', mavenSettingsConfig: 'coherence-operator-maven-settings', tempBinDir: '') {
                    sh '''
		    export http_proxy_host=`echo $HTTP_PROXY | sed -e 's,http://,,' | cut -d : -f1`
		    export http_proxy_port=`echo $HTTP_PROXY | sed  -e 's,/$,,' -e 's,http://,,' | cut -d : -f2`
		    export https_proxy_host=`echo $HTTPS_PROXY | sed -e 's,http.*://,,' | cut -d : -f1`
		    export https_proxy_port=`echo $HTTPS_PROXY | sed -e 's,/$,,' -e 's,http.*://,,' | cut -d : -f2`
                    cd docs/samples && mvn -Dhttp.proxyHost=$http_proxy_host -Dhttps.proxyHost=$https_proxy_host -Dhttp.proxyPort=$http_proxy_port -Dhttps.proxyPort=$https_proxy_port -Pdocker,docker-v1,docker-v2 clean install '
                    '''
                }
            }
        }
        stage('docker-push') {
            agent {
                label 'Kubernetes'
            }
            steps {
                echo 'Docker Push - SKIP '
            }
        }
        stage('kubernetes-tests') {
            agent {
                label 'Kubernetes'
            }
            steps {
                echo 'Kubernetes Tests - SKIP'
		sh 'helm repo add coherence https://oracle.github.io/coherence-operator/charts'
		sh 'helm repo update'
            }
            post {
                always {
                    sh '''
                        echo 'The End'
                    '''
                    deleteDir()
                }
            }
        }
    }
}
