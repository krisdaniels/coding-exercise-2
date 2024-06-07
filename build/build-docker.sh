#!/bin/bash

DOCKER_BUILDKIT=1 docker build -t coding-exercise -f package/docker/Dockerfile ../
