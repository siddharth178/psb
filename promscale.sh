#!/bin/bash

echo "Importing the environment variables"
source "config.sh"

promscale-0.4.1 --db-name "${POSTGRES_DB_NAME}" --db-password "${POSTGRES_PASSWORD}" --db-ssl-mode allow --db-host localhost --db-port "${HOST_PORT}"
