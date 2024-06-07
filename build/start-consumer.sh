#!/bin/bash

RABBITMQ_EP="amqp://guest:guest@host.docker.internal:5672/"
mkdir -p $PWD/out

docker run \
    --add-host=host.docker.internal:host-gateway \
    -ti \
    --rm \
    -v $PWD/out:/out \
    coding-exercise \
    --rabbitmq_endpoint=$RABBITMQ_EP \
    consumer \
    --outputfolder=/out/
