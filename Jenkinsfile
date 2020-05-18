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
               def root = tool name: 'Go'
               withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
               sh 'go version'
               sh 'pwd'
               sh 'go get -u -d github.com/gruntwork-io/terratest/modules/http-helper'
               sh 'go get -u -d github.com/schollz/progressbar/v3'
               sh 'go mod init "github.com/gruntwork-io/terratest/master/modules"'
               sh 'go test -v -tags kubernetes -run TestKubernetes'                  
               }
            }
         }
      }   
   }
}
