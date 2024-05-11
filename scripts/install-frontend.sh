# !/bin/bash

docker run --privileged \
        --name alcoj-frontend \
        -p 3000:3000 \
        -d frontend:v0.0.1
