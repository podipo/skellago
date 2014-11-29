#!/bin/bash

MODULE=$1

DIST_DIR=/skellago/dist
COLLECTION_DIR=$DIST_DIR/$MODULE
OUTPUT_FILE=$DIST_DIR/${MODULE}-artifact.tar.gz

mkdir -p $COLLECTION_DIR/bin/$MODULE

# Gather together everything needed in the artifact into a build directory
cp -r /skellago/go/src/podipo.com/skellago/$MODULE/static $COLLECTION_DIR/
cp -r $GOBIN/* $COLLECTION_DIR/bin/$MODULE/

# Tar everything needed in the artifact
(cd $COLLECTION_DIR && tar -zcvf $OUTPUT_FILE *)
