#!/bin/bash

# Supposing to deploy on x86_64 architecture
docker build -t fabrizio2210/light_cicd-backend -f docker/x86_64/Dockerfile-backend .
docker build -t fabrizio2210/light_cicd-frontend -f docker/x86_64/Dockerfile-frontend .
docker-compose -f docker/lib/stack.yml  up
