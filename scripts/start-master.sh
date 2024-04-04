# !/bin/bash
# Usage: ./start-master.sh [LANGUAGE] [MASTER_PORT] [MASTER_WEBSOCKET_PORT]

LANGUAGE=python     # Default language python/golang
MASTER_PORT=8080
MASTER_WEBSOCKET_PORT=7070


if [[ "$#" -ge 3 ]]; then
        LANGUAGE="$1"
        MASTER_PORT="$2"
        MASTER_WEBSOCKET_PORT="$3"
elif [[ "$#" -eq 1 ]]; then
        LANGUAGE="$1"
fi

MASTER_NAME=${LANGUAGE}-master

docker run --privileged \
    --name $MASTER_NAME \
    -e PORT=$MASTER_PORT \
    -e WEBSOCKET_PORT=$MASTER_WEBSOCKET_PORT \
    -e LANGUAGE=$LANGUAGE \
    -p $MASTER_PORT:$MASTER_PORT \
    -p 7070:7070 \
    -v sandbox:/sandbox \
    -it master:v0.0.1
