#/bin/bash -e

# This script runs commands in the CoreOS shell via vagrant ssh

vagrant ssh -c /bin/bash -c "$*"