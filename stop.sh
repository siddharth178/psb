#!/bin/bash

echo "Importing the environment variables"
source "config.sh"

echo "Stopping tsdb docker container"
docker stop "${DOCKER_CONTAINER_NAME}"
