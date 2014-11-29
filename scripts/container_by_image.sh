#!/bin/bash

# performs a local docker command 
# Usage: container_by_image.sh <action like stop or rm> <image_name>

ACTION=$1
IMAGE_NAME=$2

docker ps -a | grep $IMAGE_NAME | awk '{print $1}' | xargs docker $ACTION