#!/bin/bash

# set -x
set -e
set -o pipefail

#
export OCI_PORT=${1:-"5000"}

export DOCKER_REGISTRY_HOST="local.dhnt.io:$OCI_PORT"
export IPFS_PATH=${IPFS_PATH:-"$HOME/.dhnt/ipfs/data"}

#
function build_run {
    # build Docker image
    docker build --quiet -t $1 --build-arg REPO_NAME=$1 --build-arg DOCKER_REGISTRY_HOST=$DOCKER_REGISTRY_HOST .

    # test run
    docker run $1
}

function cleanup {
    docker rmi -f $(docker image ls -q $1)
}

#
my=$DOCKER_REGISTRY_HOST

###
repo_name="hello/docker-cli:v0.0.1-b$RANDOM"
echo "*** push/pull $repo_name using docker cli..."
build_run $repo_name

# push
docker tag $repo_name $my/$repo_name
docker push --quiet $my/$repo_name

# pull
docker pull --quiet $my/$repo_name

# run image
docker run $my/$repo_name

# clean up
cleanup $repo_name

echo "***test complete."
