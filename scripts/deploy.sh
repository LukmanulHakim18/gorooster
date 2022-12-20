#!/usr/bin/env bash

set -e
set -o pipefail

if [ $# -eq 0 ]
  then
    echo "Usage: deploy.sh [version] [namespace] [service]"
    exit 1
fi

cat k8s/deployment.yaml | sed 's/\$BUILD_NUMBER'"/$1/g" | sed 's/\$NAMESPACE'"/$2/g" | kubectl apply -n $3 -f - --kubeconfig=kubeconfig.conf
kubectl apply -f k8s/service.yaml -n $3 --kubeconfig=kubeconfig.conf

# echo 'Update image version in docker-compose.yml'
# sed -E -i.bak.$(date +"%Y%m%d%H%M%S") "s/(.$2\/$3:).*/\1$1/" docker-compose.yml

# echo 'Create and start containers'
# docker-compose up -d
