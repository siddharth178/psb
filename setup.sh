#!/bin/bash

echo "Importing the environment variables"
source "config.sh"

echo "Destroying existing database and docker containers"
docker rm -f "${DOCKER_CONTAINER_NAME}" || true

echo "Creating docker network bridge (ignore error, if run the second time)"
docker network create -d bridge "${DOCKER_BRIDGE_NETWORK}" || true

echo "Running docker image"
docker run -d\
    --name "${DOCKER_CONTAINER_NAME}" \
    -e POSTGRES_PASSWORD="${POSTGRES_PASSWORD}" \
    -e POSTGRES_USER="${POSTGRES_USER}" \
    -p "${HOST_PORT}:${NODE_PORT}" \
    --network "${DOCKER_BRIDGE_NETWORK}" \
    "${DOCKER_IMAGE_NAME}" "${POSTGRES_DB_NAME}" "${POSTGRES_CMD_ARGS}"

