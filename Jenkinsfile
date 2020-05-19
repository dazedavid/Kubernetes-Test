ppipeline {
   agent any
   stages {
      stage ('Checking out GIT Files') {
         steps {
            checkout scm
        }
      }
      stage ('Checking Go Version') {
         steps {
            scrpit {
               sh 'go version'
               }
            }
         }
      }
      stage ('Preparing Go Test') {
         steps {
            script {
              sh 'go mod init "github.com/gruntwork-io/terratest/master/modules"'                
            }
         }
      }
      stage ('Doing Test') {
         steps {
            script {
               sh 'go test -v -tags kubernetes -run TestKubernetes' 
            }
         }
      }
   }
}
