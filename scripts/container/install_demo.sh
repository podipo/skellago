#!/bin/bash -ex

go install -v podipo.com/skellago/demo
$GOBIN/demo
