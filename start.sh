#!/bin/bash

echo "Importing the environment variables"
source "config.sh"

echo "Running tsdb in docker container"
docker start "${DOCKER_CONTAINER_NAME}"

