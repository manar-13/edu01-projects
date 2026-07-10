#!/bin/bash

# Remove previous container if exists
docker rm -f ascii-art-web 2>/dev/null

# Remove existing image if any
docker rmi -f ascii-art-web-image 2>/dev/null

# Build the Docker image
docker build -t ascii-art-web-image .

# Run the container on port 8080
docker run -p 8080:8080 --detach --name ascii-art-web ascii-art-web-image

# Show running containers
docker ps