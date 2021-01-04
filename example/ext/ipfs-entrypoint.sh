#!/bin/sh
set -e

##
export IPFS_PATH=${IPFS_PATH:-/data/ipfs}
export IPFS_PROFILE=${IPFS_PROFILE:-lowpower}

function config_ipfs() {
	ipfs config Addresses.API /ip4/0.0.0.0/tcp/5001
  	ipfs config Addresses.Gateway /ip4/0.0.0.0/tcp/8080
	# ipfs config --json Swarm.EnableAutoRelay true
	# ipfs config --json Experimental.Libp2pStreamMounting true
	# ipfs config --json Experimental.FilestoreEnabled true
	ipfs config --json API.HTTPHeaders.Access-Control-Allow-Origin '["http://local.dhnt.io:5001", "http://127.0.0.1:5001", "https://webui.ipfs.io"]'
	ipfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "GET", "POST"]'
}

#
mkdir -p $IPFS_PATH

ipfs version
if [ ! -e "$IPFS_PATH/config" ]; then
  ipfs init --profile=$IPFS_PROFILE
fi

export IPFS_ID=$(ipfs id "--format=<id>")
echo "configuring $IPFS_ID ..."
config_ipfs

ipfs daemon --migrate=true

##
# export IPDR_CID_STORE=${IPDR_CID_STORE:-"/data/ipdr/cids"}
# export IPDR_CID_RESOLVER=${IPDR_CID_RESOLVER:-"file:$IPDR_CID_STORE"}

# mkdir -p $IPDR_CID_STORE

# ipdr server --port 5000 --cid-store $IPDR_CID_STORE --cid-resolver $IPDR_CID_RESOLVER &

# wait
exec "$@"


