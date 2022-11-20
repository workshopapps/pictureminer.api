# pictureminer.api

## Requirements:
**This Project requires a version of GO 1.19.3 or later**

**A Code Editor**

## Instructions:

**step 1: clone the DEV Branch using "git clone https://github.com/workshopapps/pictureminer.api.git "**

**step 2: Branch out to your feature branch " git checkout -b "your-feature-branch" "**

**step 3: Set upstream to DEv " git branch --set-upstream-to=origin/dev feature "**

**step 4: pull latest work " git pull "**

> you can start working N/B: "once your done making your feature please pull again before commiting"

**step 5: " git add . && git commit -m "Your commit message" "**

**step 6: " git push origin "Your-feature-branch":dev "**

## Link on Instructions on how to update go:

### https://buildvirtual.net/how-to-upgrade-go-on-ubuntu/#:~:text=How%20to%20Upgrade%20Go%20on%20Ubuntu%201%20Remove,named%20hello.go%2C%20which%20contains%20the%20following%20lines%3A%20



## Building Docker Image
Stop the running app and build docker image:
    docker build -t <image-name> .
    
## View Built Image
    docker image ls

## Run Docker Image
    docker run -p port:port <image-name>
# Check the output at http://localhost:port/ or http://0.0.0.0:port/ or http://127.0.0.1:port/

## Check running containers
    docker ps
  
# Stop the container

## Tag locally before pushing to the Dockerhub

    docker tag <image-name> <dockerhub-username>/<docker-image:version-number>

## Push Image
    docker login

    docker push <dockerhub-username>/<docker-image:version-number>
  
# Check the image in your Dockerhub online at https://hub.docker.com/repository/docker/
