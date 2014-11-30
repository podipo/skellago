#!/bin/bash

MODULE=$1

COLLECT_DIR=/skellago/collect
MODULE_DIR=$COLLECT_DIR/$MODULE
OUTPUT_FILE=$COLLECT_DIR/${MODULE}-artifact.tar.gz

mkdir -p $MODULE_DIR/bin/$MODULE

# Gather together everything needed in the artifact into a build directory
cp -r /skellago/go/src/podipo.com/skellago/$MODULE/static $MODULE_DIR/
cp -r $GOBIN/* $MODULE_DIR/bin/$MODULE/

# Tar everything needed in the artifact
(cd $MODULE_DIR && tar -zcvf $OUTPUT_FILE *)
