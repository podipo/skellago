#!/bin/bash

# This script runs a docker image which then calls scripts/container/install_demo.sh

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

docker run -ti --rm -v /skellago:$WORKING_DIR $DKR_ENV $* /skellago/scripts/container/install_demo.sh
