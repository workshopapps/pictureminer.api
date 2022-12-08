pipeline {
	agent any
	
	stages {

		stage("Build Backend"){

			steps {
                sh " rm -rf pictureminer.api"
			    sh "git clone https://github.com/workshopapps/pictureminer.api.git"
                dir('pictureminer.api') {
                    sh "go build"
                }
			}  
	    }
		stage("Deploy Backend"){

			steps {
                sh "sudo su - javi && whoami"
				
				sh "sudo cp -fr ${WORKSPACE}/pictureminer.api/* /home/javi/backend/"
				sh "sudo systemctl restart discripto_api.service"
			} 
	    }

		stage("Performance test"){

			steps{
				echo 'Installing k6'
                // sh 'sudo chmod +x setup_k6.sh'
                // sh 'sudo ./setup_k6.sh'
                echo 'Running K6 performance tests...'
				sh 'ls -a'
				sh "pwd"
                sh 'k6 run Performance_Test_Discriptob.js'
			}
		}
		
    }
    post{
        failure{
            emailext attachLog: true, 
            to: 'hng9.discripto@gmail.com',
            subject: '${BUILD_TAG} Build failed',
            body: '${BUILD_TAG} Build Failed \nMore Info can be found here: ${BUILD_URL} or in the log file below'
        }
    }
}
