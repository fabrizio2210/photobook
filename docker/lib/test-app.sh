#!/bin/bash -xeu

echo $(pwd)
docker build -t fabrizio2210/photobook-backend -f docker/x86_64/Dockerfile-backend .
if [ "$(uname -m)" == "x86_64" ]; then
  stack="docker/lib/stack-test-x86_64.yml"
else
  stack="docker/lib/stack-test-armv7hf.yml"
fi

docker-compose -f ${stack} run flask
docker-compose -f ${stack} down
