#!/usr/bin/env bash

BUILD_VER=${1}

rm wasmclient.wasm

echo "GOOS=js GOARCH=wasm go build -o wasmclient.wasm -ldflags -X main.Ver=${BUILD_VER}"
GOOS=js GOARCH=wasm go build -o wasmclient.wasm -ldflags "-X main.Ver=${BUILD_VER}" wasmclient.go

echo "move files"
mv wasmclient.wasm clientdata/
