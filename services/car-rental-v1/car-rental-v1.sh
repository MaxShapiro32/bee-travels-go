#!/bin/bash

if [ "$1" = "start" ]; then
    go env -w GOBIN=$(pwd)/bin
    go install
    bin/car-rental-v1
elif [ "$1" = "test" ]; then
    go test
fi
