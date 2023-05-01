#!/bin/bash

# Supposing to deploy on x86_64 architecture
docker build -t fabrizio2210/photobook-backend -f docker/x86_64/Dockerfile-backend .
docker build -t fabrizio2210/photobook-frontend -f docker/x86_64/Dockerfile-frontend .
docker build -t fabrizio2210/photobook-worker -f docker/x86_64/Dockerfile-worker .
docker build -t fabrizio2210/photobook-api -f docker/x86_64/Dockerfile-api .
docker build -t fabrizio2210/photobook-sse -f docker/x86_64/Dockerfile-sse .
docker compose -f docker/lib/stack.yml  up
