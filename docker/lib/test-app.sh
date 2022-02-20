#!/bin/bash

# Supposing to deploy on x86_64 architecture
docker build -t fabrizio2210/photobook-backend-dev -f docker/x86_64/Dockerfile-backend-dev .
docker-compose -f docker/lib/stack-test.yml run flask
