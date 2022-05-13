#!/bin/bash -eux

# Supposing to build on x86_64 architecture
docker build -t fabrizio2210/photobook-backend:armv7hf -f docker/armv7hf/Dockerfile-backend .
docker push fabrizio2210/photobook-backend:armv7hf
docker build -t fabrizio2210/photobook-frontend:armv7hf -f docker/armv7hf/Dockerfile-frontend .
docker push fabrizio2210/photobook-frontend:armv7hf
