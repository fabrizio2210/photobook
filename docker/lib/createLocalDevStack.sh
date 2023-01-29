#!/bin/bash

# Supposing to deploy on x86_64 architecture
docker build -t fabrizio2210/photobook-backend-dev -f docker/x86_64/Dockerfile-backend-dev .
docker build -t fabrizio2210/photobook-frontend-dev -f docker/x86_64/Dockerfile-frontend-dev .
docker build -t fabrizio2210/photobook-worker-dev -f docker/x86_64/Dockerfile-worker-dev .
docker compose -f docker/lib/stack-dev.yml  up
