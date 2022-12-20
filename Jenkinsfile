def Tag_Release() {
    sh 'git describe --tags --abbrev=0 > Tag_Release'
    def Tag_Release = readFile('Tag_Release').trim()
    sh 'rm Tag_Release'
    Tag_Release
}

pipeline {
    agent {
        node {
            label 'slave-00 || slave-01'
            customWorkspace "workspace/${env.BRANCH_NAME}/src/git.bluebird.id/mybb/order-query"
        }
    }
    environment {
        SERVICE = 'order-query'
        TEAMS_MICROSOFT = credentials('56abe4aa-f20c-4509-81b3-2d427abc8565')
        PROJECT= "${env.SERVICE}"
        TESTING = "${env.EXECUTOR_NUMBER}-${env.BUILD_NUMBER}"
        BRANCH_NAME = "${env.BRANCH_NAME}"
        BUILD_NUMBER = "${env.BUILD_NUMBER}"
    }
    options {
        buildDiscarder(logRotator(daysToKeepStr: env.BRANCH_NAME == 'master' ? '90' : '30'))
    }
    stages {
        stage('Checkout') {
            when {
                anyOf { branch 'master'; branch 'develop'; branch 'staging' }
            }
            steps {
                echo 'Checking out from Git'
                checkout scm
            }
        }

        stage('Prepare') {
            steps {
                withCredentials([file(credentialsId: '3521ab7f-3916-4e56-a41e-c0dedd2e98e9', variable: 'sa')]) {
                sh "cp $sa service-account.json"
                sh "chmod 644 service-account.json"
                sh "docker login -u _json_key --password-stdin https://asia.gcr.io < service-account.json"
                }
            }
        }
        stage('Testing and Code Review'){
            environment{
                 GOLANG_PROTOBUF_REGISTRATION_CONFLICT= "warn"
            }
             steps {
                withCredentials([string(credentialsId: '04398f9c-36e4-4161-b6b2-9098e7c26ad9', variable: 'TOKEN')]) {
                    sh 'chmod +x testing.sh'
                    sh './testing.sh $TESTING $BRANCH_NAME $BUILD_NUMBER $TOKEN'
                }
            }
        }
        
        stage('Build and Deploy') {
            environment {
                GOPATH = "${env.JENKINS_HOME}/workspace/${env.BRANCH_NAME}"
                PATH = "${env.GOPATH}/bin:${env.PATH}"
                VERSION_PREFIX = '1.0'
            }
            stages {
                stage('Deploy to development') {
                    when {
                        branch 'develop'
                    }
                    environment {
                        ALPHA = "${env.VERSION_PREFIX}-alpha${env.BUILD_NUMBER}"
                        NAMESPACE="mybluebird-dev-1202"
                    }
                    steps {
                        withCredentials([file(credentialsId: '4f1ac961-5456-463e-a756-122c54957d59', variable: 'kubeconfig')]) {
                        sh "cp $kubeconfig kubeconfig.conf"
                        sh "chmod 644 kubeconfig.conf"
                        sh "gcloud auth activate-service-account --key-file service-account.json"
                        sh 'chmod +x ./build.sh'
                        sh './build.sh $ALPHA'
                        sh 'chmod +x ./deploy.sh'
                        sh './deploy.sh $ALPHA $NAMESPACE default'
                        sh 'rm kubeconfig.conf service-account.json'
                        }
                    }
                }
                stage('Deploy to staging') {
                    when {
                        branch 'staging'
                    }
                    environment {
                        BETA = "${env.VERSION_PREFIX}-beta${env.BUILD_NUMBER}"
                        NAMESPACE="mybluebird"
                    }
                    steps {
                        withCredentials([file(credentialsId: '4f1ac961-5456-463e-a756-122c54957d59', variable: 'kubeconfig')]) {
                            sh "cp $kubeconfig kubeconfig.conf"
                            sh "chmod 644 kubeconfig.conf"
                            sh "gcloud auth activate-service-account --key-file service-account.json"
                            sh 'chmod +x ./scripts/build.sh'
                            sh './scripts/build.sh $BETA'
                            sh 'chmod +x ./scripts/deploy.sh'
                            sh './scripts/deploy.sh $BETA $NAMESPACE staging'
                            sh 'rm kubeconfig.conf service-account.json'
                        }
                    }
                }
                stage('Deploy to production') {
                    when {
                        tag "v*"
                    }
                    environment {
                        VERSION = "${env.TAG_NAME}"
                        NAMESPACE="mybluebird-prd-2003"
                    }
                    steps {
                        withCredentials([file(credentialsId: '4ad85739-181e-426f-ba88-5f6523421945', variable: 'kubeconfig')]) {
                            sh "cp $kubeconfig kubeconfig.conf"
                            sh "chmod 644 kubeconfig.conf"
                            sh "gcloud auth activate-service-account --key-file service-account.json"
                            sh 'chmod +x build.sh'
                            sh './build.sh $VERSION'
                            sh 'chmod +x deploy.sh'
                            sh './deploy.sh $VERSION $NAMESPACE default'
                            sh 'rm kubeconfig.conf service-account.json'
                        }

                    }
                }
            }
        }
    }
    post {
        success {
            office365ConnectorSend webhookUrl: "$TEAMS_MICROSOFT",
                message: "Application mybb Service $SERVICE has been [deployed]",
                color: "05b222",
                status: 'Success'
        }
        failure {
            office365ConnectorSend webhookUrl: "$TEAMS_MICROSOFT",
                message: "Application mybb Service $SERVICE has been [Failed]",
                color: "d00000",
                status: 'Failed'
        }
    }

}
