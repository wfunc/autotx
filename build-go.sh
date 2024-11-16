#!/bin/bash
set -e

pkg_ver=`git rev-parse --abbrev-ref HEAD`

if [ "$1" == "docker" ];then
    cd `dirname ${0}`
    docker build -t autotx:$pkg_ver -f DockerfileGo .
else
    mkdir -p build/server
    go build -v -o build/server/autotx .
fi