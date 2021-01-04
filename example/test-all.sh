#!/bin/bash

# required external services:
# 1) start sshd on 5022 with username/password: app/app
# 2) start ipfs on 5001
# set -x

function test_server() {
    oci --port $1 --store $2 &

    ./test.sh $1 $2
}

##
declare -a arr=(
"memory:/tmp/oci-memory-$RANDOM"
"file:/tmp/oci-file-$RANDOM"
"scp://app:app@local.ipdr.io:5022/tmp/oci-scp-$RANDOM"
"ipfs://local.ipdr.io:5001/tmp/oci-ipfs-$RANDOM"
)

if ! command -v oci &> /dev/null
then
    echo "  *** oci not found, building it..."
    (cd ../../oci && go install ./cmd/oci)
fi

##
for store in "${arr[@]}"
do
   port=$((5000 + $RANDOM % 1000))
   echo "*** testing $port $store"
   time test_server $port $store
done

exit 0
