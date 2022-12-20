#!/usr/bin/env bash

set -e
set -o pipefail

if [ $# -eq 0 ]
  then
    echo "Usage: testing.sh [version]"
    exit 1
fi


docker build --no-cache -t asia.gcr.io/$PROJECT/$SERVICE:$1 --build-arg BRANCH_NAME=$2 --build-arg BUILD_NUMBER=$3 --build-arg TOKEN=$4 --build-arg SSH_KEY="$(cat ~/.ssh/id_ecdsa)" -f Dockerfile.test --network=host .
docker rmi asia.gcr.io/$PROJECT/$SERVICE:$1
