#!/bin/bash -e 

# Stop and remove all of the images in the docker container
docker ps -a | grep -v CONTAINER | awk '{print $1}' | xargs docker stop
docker ps -a | grep -v CONTAINER | awk '{print $1}' | xargs docker rm
