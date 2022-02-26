#!/bin/bash

set -e
set -u
set -x
#############
# Environment

registryCredential='docker-login'
buildDir='/tmp/build'
debPackageStash='deb'
venvPackageStash='venv'

runRepository="$(mktemp -d)"
workspace=$(dirname $0)
if ! echo $workspace |  grep "^/" ;  then
  workspace="$(pwd)/$workspace"
fi
repository='/tmp'
changedFiles="$(git diff --name-only HEAD^1 HEAD)"

docker/lib/test-app.sh
exit 0
