pipeline {
   agent any
   stages {
      stage ('Checking out GIT Files') {
         steps {
            checkout scm
        }
      }
      stage ('Path assigning') {
         steps {
               sh 'PATH="/usr/local/bin:$PATH"'
         }
      }
      stage ('Doing Test Jobs') {
         steps {
            script {
               withEnv(["KUBECONFIG=$HOME/.kube/kubeconfig"]){
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
