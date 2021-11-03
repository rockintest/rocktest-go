#!/bin/sh

VER=$1

test -z "$1" && {
    echo "Usage: build.sh version"
    exit 1
}

echo "Use release $VER : https://github.com/rockintest/rocktest-go/releases/download/$VER/rocktest-go-$VER-linux-amd64.tar.gz"
curl -L https://github.com/rockintest/rocktest-go/releases/download/$VER/rocktest-go-$VER-linux-amd64.tar.gz | tar xzf -

docker build -t rocktest-go .

rm rocktest-go

docker tag rocktest-go rockintest/rocktest-go:latest
docker push rockintest/rocktest-go:latest

docker tag rocktest-go rockintest/rocktest-go:$1
docker push rockintest/rocktest-go:$1

