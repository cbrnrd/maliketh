#!/bin/bash

# This script is used to create the database and admin operator config.

# Get the container id of the `server_operator` container.
container_id=$(docker ps | grep server_operator | awk '{print $1}')

# Confirm the user wants to do this
echo -n "This will delete the database and create a new one. Are you sure you want to do this? (y/n) "
read -r response
if [ "$response" != "y" ]; then
    echo "Aborting."
    exit 1
fi

# Run make_db.py
docker exec $container_id python3 make_db.py