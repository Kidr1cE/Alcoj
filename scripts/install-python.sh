# !/bin/bash

MASTER_ADDR=host.docker.internal:8080

bash scripts/start-master.sh python 8080 7070

num_workers=$1

WORKER_PORT=50051

for ((i=1; i<=num_workers; i++)); do
    echo "Starting worker $i on port $WORKER_PORT, connecting to master at $MASTER_ADDR"
    bash scripts/start-worker.sh $WORKER_PORT $MASTER_ADDR
    sleep 1
    ((WORKER_PORT++))
done
