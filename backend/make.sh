#!/bin/bash

# Define variables
IMAGE_NAME="mongoserver"
CONTAINER_NAME="mongo-backend-service"


# Checking and killing if running...
echo "Stopping the Docker container..."
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME

# Create workflow for container

echo "Building the Docker image..."
docker build -t $IMAGE_NAME .

echo "Creating the Docker container..."
docker create --name $CONTAINER_NAME $IMAGE_NAME

echo "Starting the Docker container..."
docker start $CONTAINER_NAME

docker ps -a
