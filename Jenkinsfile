pipeline {
   agent any
   stages {
      stage ('Checking out GIT') {
         steps {
            checkout scm
        }
      }
      stage ('Doing Test Jobs') {
         steps {
            script {
               def root = tool name: 'Go'
               withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
               sh 'go version'
                     sh 'pwd'
                     sh 'go mod init "github.com/gruntwork-io/terratest/master/modules"'
                     sh 'go test -v -tags kubernetes -run TestKubernetes > test.html'     
                      publishHTML (target: [
                           allowMissing: false,
                           alwaysLinkToLastBuild: false,
                           keepAll: true,
                           reportDir: './',
                           reportFiles: 'test.html',
                           reportName: 'Kafka Report'
                         ])
                  }
               }
             }
          }
       }
    }   
