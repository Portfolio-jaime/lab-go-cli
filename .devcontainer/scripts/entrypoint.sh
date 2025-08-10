#!/bin/sh
set -e

# Start Docker daemon in the background
echo "Entrypoint: Starting Docker daemon..."
/usr/local/bin/start-docker.sh > /tmp/docker.log 2>&1 &

# Wait for Docker daemon to be ready
echo "Entrypoint: Waiting for Docker daemon..."
while ! docker info >/dev/null 2>&1; do
    sleep 1
done
echo "Entrypoint: Docker daemon is ready and running in the background."

# Execute the command passed into the container (e.g., sleep infinity)
exec "$@"
