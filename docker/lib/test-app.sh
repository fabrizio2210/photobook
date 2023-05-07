#!/bin/bash -xeu

if [ "$(uname -m)" == "x86_64" ]; then
  arch="x86_64"
else
  arch="armv7hf"
fi

#docker build -t fabrizio2210/photobook-backend:${arch} -f docker/x86_64/Dockerfile-backend .
docker build -t fabrizio2210/photobook-api -f docker/x86_64/Dockerfile-api .
if [ "${arch}" == "x86_64" ]; then
  stack="docker/lib/stack-test-x86_64.yml"
else
  stack="docker/lib/stack-test-armv7hf.yml"
fi

docker compose -f ${stack} run api
#docker compose -f ${stack} run flask
docker compose -f ${stack} down
