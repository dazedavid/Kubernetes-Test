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
                     sh 'rm -rf /tools/golang/golang-1.13.3/go/src/github.com/gruntwork-io'
                     sh 'rm -rf /home/jenkins/src/github.com/gruntwork-io'
                     sh 'go mod init github.com/gruntwork-io/terratest/master/modules'
                     sh 'go get -u -d github.com/fatih/color'    
                     sh 'go get -u -d github.com/gruntwork-io/terratest/modules/http-helper'
                     sh 'go get github.com/schollz/progressbar' 
                     sh 'go get -u -d github.com/reiver/go-telnet'   
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
