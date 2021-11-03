#!/bin/sh

VER=$1

test -z "$1" && {
    echo "Usage: build.sh version"
    exit 1
}

curl https://github.com/rockintest/rocktest-go/releases/download/$VER/rocktest-go-$VER-linux-amd64.tar.gz | tar xzf -
