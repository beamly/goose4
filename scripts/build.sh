#!/bin/sh

echo "PATH is: ${PATH}"
echo "GOPATH is: ${GOPATH}"

echo "==== Printing all env vars ===="
printenv
echo "==============================="

make clean
make updatedeps
make test
