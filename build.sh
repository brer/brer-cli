#!/usr/bin/bash

set -e

mkdir -p ./build
rm -f ./build/*

arch=amd64
oss=(darwin linux windows)

for os in ${oss[@]}
do
	env GOOS=${os} GOARCH=${arch} go build -o ./build/brer-cli_${os}_${arch}
done
