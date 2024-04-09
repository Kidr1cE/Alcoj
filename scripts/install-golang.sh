# !/bin/bash
# usage: bash scripts/install-golang.sh 5

MASTER_ADDR=host.docker.internal:8081

bash scripts/start-master.sh golang 8081 7071

num_workers=$1

WORKER_PORT=50041

for ((i=1; i<=num_workers; i++)); do
    echo "Starting worker $i on port $WORKER_PORT, connecting to master at $MASTER_ADDR"
    bash scripts/start-worker.sh $WORKER_PORT $MASTER_ADDR
    sleep 1
    ((WORKER_PORT++))
done
