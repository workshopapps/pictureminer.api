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
		
    }
}