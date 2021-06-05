#!/bin/bash

export DOCKER_IMAGE_NAME="timescaledev/promscale-extension:latest-pg12"

export DOCKER_CONTAINER_NAME="timescaledb"
export DOCKER_BRIDGE_NETWORK="promscale-timescaledb"
export HOST_PORT="5433"
export NODE_PORT="5432"

export POSTGRES_DB_NAME="postgres"
export POSTGRES_PASSWORD="postgres"
export POSTGRES_PASSWORD="password"
export POSTGRES_CMD_ARGS="-csynchronous_commit=off"

