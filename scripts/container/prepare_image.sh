#!/bin/sh

# Usage prepare_image.sh <source dir> <path to dist artifact> <deployment dir>
# Example: prepare_image.sh containers/foo dist/artifacts/foo-amd64-latest.tar.gz deploy/containers/foo 

cd /skellago

SOURCE_DIR=$1
DIST_ARTIFACT=$2
DEPLOY_DIR=$3

# Clean out the deploy dir and then copy into it the source files and the dist artifact
rm -rf $DEPLOY_DIR
mkdir -p $DEPLOY_DIR
cp -r  $DIST_ARTIFACT $SOURCE_DIR/* $DEPLOY_DIR

