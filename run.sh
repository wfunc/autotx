#!/bin/bash
set -e

if [ "$1" == "go" ];then
    go build -v .
    ./autotx
else
    ./build-go.sh
    Verbose=1 CodeURL=https://example.com ./build/server/autotx
fi
