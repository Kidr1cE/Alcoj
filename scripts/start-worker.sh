# !/bin/bash
# Usage: ./start-worker.sh [WORKER_PORT] [MASTER_ADDR]

CONTAINER_UUID=$(uuidgen)
echo "Starting worker with id: $CONTAINER_UUID"

WORKER_PORT=50051
MASTER_ADDR=host.docker.internal:8080

# Check if command line arguments are provided for WORKER_PORT and MASTER_ADDR
if [[ "$#" -ge 2 ]]; then
        # If at least two command line arguments are present, assign them to WORKER_PORT and MASTER_ADDR
        WORKER_PORT="$1"
        MASTER_ADDR="$2"
elif [[ "$#" -eq 1 ]]; then
        # If only one command line argument is present, assume it's for WORKER_PORT
        WORKER_PORT="$1"
fi

docker run --privileged \
        --name ${CONTAINER_UUID} \
        -e WORKER_ID=${CONTAINER_UUID} \
        -e PORT=$WORKER_PORT \
        -e MASTER_ADDR=$MASTER_ADDR \
        -p $WORKER_PORT:$WORKER_PORT \
        -v /var/run/docker.sock:/var/run/docker.sock \
        -d worker:v0.0.3
