#!/bin/sh

#
# This is the prefix script which will be used during the make process to run commands in a Docker container
# The first argument should be the tag
#

terminal=""
if [ "$1" = "-t" ]; then
	terminal=" -ti"
	shift
fi

WORKING_DIR=/skellago
GOPATH=$WORKING_DIR/go
GOOS=linux
GOARCH=amd64
GOBIN=$GOPATH/bin/${GOOS}_${GOARCH}

DKR_ENV="--env GOPATH=$GOPATH \
	--env GOBIN=$GOBIN \
	--env GOARCH=$GOARCH \
	--env GOOS=$GOOS \
	--env PATH=/go/bin:/usr/src/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:$GOBIN"

MOUNTS="-v $PWD:$WORKING_DIR"

TAG=$1
shift

docker run -i $terminal --rm $DKR_ENV $MOUNTS $TAG sh -c "$*"