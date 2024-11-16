#!/bin/bash
set -e

if [ "$1" == "go" ];then
    go build -v .
    ./autotx
else
    ./build-go.sh
    ./build/server/autotx
fi
