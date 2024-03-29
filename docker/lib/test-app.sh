#!/bin/bash -xeu

if [ "$(uname -m)" == "x86_64" ]; then
  arch="x86_64"
else
  arch="armv7hf"
fi

docker build -t fabrizio2210/photobook-api -f docker/x86_64/Dockerfile-api .
if [ "${arch}" == "x86_64" ]; then
  stack="docker/lib/stack-test-x86_64.yml"
else
  stack="docker/lib/stack-test-armv7hf.yml"
fi

docker_compose="docker-compose"
docker compose version && docker_compose="docker compose"

${docker_compose} -f ${stack} run api
${docker_compose} -f ${stack} down
