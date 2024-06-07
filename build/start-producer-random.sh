#!/bin/bash

RABBITMQ_EP="amqp://guest:guest@host.docker.internal:5672/"
docker run \
    --add-host=host.docker.internal:host-gateway \
    -ti \
    --rm \
    coding-exercise \
    --rabbitmq_endpoint=$RABBITMQ_EP \
    producer \
    random
