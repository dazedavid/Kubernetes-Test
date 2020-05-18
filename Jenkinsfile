pipeline {
   agent any
   stages {
      stage ('Checking out GIT Files') {
         steps {
            checkout scm
        }
      }
      stage ('Doing Test Jobs') {
         steps {
            script {
               withKubeConfig([
                    credentialsId: 'jenkins-robot',
                    caCertificate: '',
                    serverUrl: 'https://kubernetes.docker.internal:6443',
                    contextName: '',
                    clusterName: '',
                    namespace: ''
                    ]) {
               def root = tool name: 'Go'
               withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
               sh 'go version'
               sh 'pwd'
               sh 'kubectl get pods'
               sh 'go mod init "github.com/gruntwork-io/terratest/master/modules"'
               sh 'go test -v -tags kubernetes -run TestKubernetes'                  
                }
            }    
         }
        }   
      }   
   }
}
