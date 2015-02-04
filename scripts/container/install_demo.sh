#!/bin/bash -ex

go install -v podipo.com/skellago/demo
go install -v example.com/api/example_demo
$GOBIN/demo
$GOBIN/example_demo
